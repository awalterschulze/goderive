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

package fmap

import (
	"flag"
	"fmt"
	"go/types"

	"github.com/awalterschulze/goderive/derive"
)

var Prefix = flag.String("fmap.prefix", "deriveFmap", "set the prefix for fmap functions that should be derived.")

type fmap struct {
	derive.TypesMap
	printer  derive.Printer
	bytesPkg derive.Import
}

func New(typesMap derive.TypesMap, p derive.Printer) *fmap {
	return &fmap{
		TypesMap: typesMap,
		printer:  p,
	}
}

func (this *fmap) Name() string {
	return "fmap"
}

func (this *fmap) Add(name string, typs []types.Type) (string, error) {
	if len(typs) != 2 {
		return "", fmt.Errorf("%s does not have two arguments", name)
	}
	sliceTyp, ok := typs[1].(*types.Slice)
	if !ok {
		return "", fmt.Errorf("%s, the second argument, %s, is not of type slice", name, this.TypeString(typs[1]))
	}
	sig, ok := typs[0].(*types.Signature)
	if !ok {
		return "", fmt.Errorf("%s, the second argument, %s, is not of type function", name, this.TypeString(typs[0]))
	}
	params := sig.Params()
	if params.Len() != 1 {
		return "", fmt.Errorf("%s, the second argument is a function, but wanted a function with one argument", name)
	}
	elemTyp := sliceTyp.Elem()
	inTyp := params.At(0).Type()
	if !types.Identical(inTyp, elemTyp) {
		return "", fmt.Errorf("%s the function input type and slice element type are different %s != %s",
			name, inTyp, elemTyp)
	}
	res := sig.Results()
	if res.Len() != 1 {
		return "", fmt.Errorf("%s, the function argument does not have a single result, but has %d resulting parameters", name, res.Len())
	}
	outTyp := res.At(0).Type()
	return this.SetFuncName(name, inTyp, outTyp)
}

func (this *fmap) Generate() error {
	for _, typs := range this.ToGenerate() {
		if err := this.genFuncFor(typs[0], typs[1]); err != nil {
			return err
		}
	}
	return nil
}

func (this *fmap) genFuncFor(in, out types.Type) error {
	p := this.printer
	this.Generating(in, out)
	inStr := this.TypeString(in)
	outStr := this.TypeString(out)
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
