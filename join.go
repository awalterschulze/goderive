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

var joinPrefix = flag.String("join.prefix", "deriveJoin", "set the prefix for join functions that should be derived.")

type join struct {
	TypesMap
	qual     types.Qualifier
	printer  Printer
	bytesPkg Import
}

func newJoin(p Printer, qual types.Qualifier, typesMap TypesMap) *join {
	return &join{
		TypesMap: typesMap,
		qual:     qual,
		printer:  p,
	}
}

func (this *join) Generate(pkgInfo *loader.PackageInfo, prefix string, call *ast.CallExpr) (bool, error) {
	fn, ok := call.Fun.(*ast.Ident)
	if !ok {
		return false, nil
	}
	if !strings.HasPrefix(fn.Name, prefix) {
		return false, nil
	}
	if len(call.Args) != 1 {
		return false, fmt.Errorf("%s does not have one argument", fn.Name)
	}
	t0 := pkgInfo.TypeOf(call.Args[0])
	if t0 == nil {
		return false, nil
	}
	sliceTyp, ok := t0.(*types.Slice)
	if !ok {
		return false, fmt.Errorf("%s, the argument, %s, is not of type slice", fn.Name, t0)
	}
	sliceOfSliceTyp, ok := sliceTyp.Elem().(*types.Slice)
	if !ok {
		return false, fmt.Errorf("%s, the argument, %s, is not of type slice of slice", fn.Name, t0)
	}
	elemType := sliceOfSliceTyp.Elem()
	if err := this.SetFuncName(fn.Name, elemType); err != nil {
		return false, err
	}
	for _, typs := range this.ToGenerate() {
		if err := this.genFuncFor(typs[0]); err != nil {
			return false, err
		}
	}
	return true, nil
}

func (this *join) genFuncFor(typ types.Type) error {
	p := this.printer
	this.Generating(typ)
	typStr := types.TypeString(typ, this.qual)
	p.P("")
	p.P("func %s(list [][]%s) []%s {", this.GetFuncName(typ), typStr, typStr)
	p.In()
	p.P("if list == nil {")
	p.In()
	p.P("return nil")
	p.Out()
	p.P("}")
	p.P("res := []%s{}", typStr)
	p.P("for _, elem := range list {")
	p.In()
	p.P("res = append(res, elem...)")
	p.Out()
	p.P("}")
	p.P("return res")
	p.Out()
	p.P("}")
	return nil
}
