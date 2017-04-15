//  Copyright 2017 Walter Schulze
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/loader"
)

var fmapPrefix = flag.String("fmap.prefix", "deriveFmap", "set the prefix for fmap functions that should be derived.")

type fmap struct {
	TypesMap
	qual     types.Qualifier
	printer  Printer
	bytesPkg Import
}

func newFmap(p Printer, qual types.Qualifier, typesMap TypesMap) *fmap {
	return &fmap{
		TypesMap: typesMap,
		qual:     qual,
		printer:  p,
	}
}

func (this *fmap) Generate(pkgInfo *loader.PackageInfo, prefix string, call *ast.CallExpr) (bool, error) {
	fn, ok := call.Fun.(*ast.Ident)
	if !ok {
		return false, nil
	}
	if !strings.HasPrefix(fn.Name, prefix) {
		return false, nil
	}
	if len(call.Args) != 2 {
		return false, fmt.Errorf("%s does not have two arguments", fn.Name)
	}
	t0 := pkgInfo.TypeOf(call.Args[0])
	t1 := pkgInfo.TypeOf(call.Args[1])
	if t0 == nil {
		return false, nil
	}
	if t1 == nil {
		return false, nil
	}
	sliceTyp, ok := t1.(*types.Slice)
	if !ok {
		return false, fmt.Errorf("%s, the second argument, %s, is not of type slice", fn.Name, t1)
	}
	sig, ok := t0.(*types.Signature)
	if !ok {
		return false, fmt.Errorf("%s, the second argument, %s, is not of type function", fn.Name, t0)
	}
	params := sig.Params()
	if params.Len() != 1 {
		return false, fmt.Errorf("%s, the second argument is a function, but wanted a function with one argument", fn.Name)
	}
	elemTyp := sliceTyp.Elem()
	inTyp := params.At(0).Type()
	if !types.Identical(inTyp, elemTyp) {
		return false, fmt.Errorf("%s the function input type and slice element type are different %s != %s",
			fn.Name, inTyp, elemTyp)
	}
	res := sig.Results()
	if res.Len() != 1 {
		return false, fmt.Errorf("%s, the function argument does not have a single result, but has %d resulting parameters", fn.Name, res.Len())
	}
	outTyp := res.At(0).Type()
	if err := this.SetFuncName(fn.Name, inTyp, outTyp); err != nil {
		return false, err
	}
	for _, typs := range this.ToGenerate() {
		if err := this.genFuncFor(typs[0], typs[1]); err != nil {
			return false, err
		}
	}
	return true, nil
}

func (this *fmap) genFuncFor(in, out types.Type) error {
	p := this.printer
	this.Generating(in, out)
	inStr := types.TypeString(in, this.qual)
	outStr := types.TypeString(out, this.qual)
	p.P("")
	p.P("func %s(f func(%s) %s, list []%s) []%s {", this.GetFuncName(in, out), inStr, outStr, inStr, outStr)
	p.In()
	p.P("out := make([]%s, len(list))", outStr)
	p.P("for i, elem := range list {")
	p.In()
	p.P("out[i] = f(elem)")
	p.Out()
	p.P("}")
	p.P("return out")
	p.Out()
	p.P("}")
	return nil
}
