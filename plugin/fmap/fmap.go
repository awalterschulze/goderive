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

// Package fmap contains the implementation of the fmap plugin, which generates the deriveFmap function.
//
// The deriveFmap function applies a given function to each element of a list, returning a list of results in the same order.
//   deriveFmap(func(A) B, []A) []B
//
// More things to come:
//	- currently only slices are supported, think about supporting other types and not just slices
//	- think about functions without a return type
package fmap

import (
	"fmt"
	"go/types"

	"github.com/awalterschulze/goderive/derive"
)

// NewPlugin creates a new fmap plugin.
// This function returns the plugin name, default prefix and a constructor for the fmap code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("fmap", "deriveFmap", New)
}

// New is a constructor for the fmap code generator.
// This generator should be reconstructed for each package.
func New(typesMap derive.TypesMap, p derive.Printer, deps map[string]derive.Dependency) derive.Generator {
	return &gen{
		TypesMap: typesMap,
		printer:  p,
	}
}

type gen struct {
	derive.TypesMap
	printer derive.Printer
}

func (this *gen) Add(name string, typs []types.Type) (string, error) {
	if len(typs) != 2 {
		return "", fmt.Errorf("%s does not have two arguments", name)
	}
	switch typs[1].(type) {
	case *types.Slice:
		_, _, err := this.sliceInOut(name, typs)
		if err != nil {
			return "", err
		}
		return this.SetFuncName(name, typs...)
	case *types.Basic:
		_, err := this.stringOut(name, typs)
		if err != nil {
			return "", err
		}
		return this.SetFuncName(name, typs...)
	}
	return "", fmt.Errorf("unsupported type %s, not a slice or a string", typs[1])
}

func (this *gen) stringOut(name string, typs []types.Type) (outTyp types.Type, err error) {
	typs[1] = types.Default(typs[1])
	basic, ok := typs[1].(*types.Basic)
	if !ok {
		return nil, fmt.Errorf("%s, the second argument, %s, is not of type basic", name, this.TypeString(typs[1]))
	}
	if basic.Kind() != types.String {
		return nil, fmt.Errorf("%s, the second argument, %v, is not a string", name, basic)
	}
	sig, ok := typs[0].(*types.Signature)
	if !ok {
		return nil, fmt.Errorf("%s, the second argument, %s, is not of type function", name, this.TypeString(typs[0]))
	}
	params := sig.Params()
	if params.Len() != 1 {
		return nil, fmt.Errorf("%s, the second argument is a function, but wanted a function with one argument", name)
	}
	elemTyp := types.Typ[types.Rune]
	inTyp := params.At(0).Type()
	if !types.Identical(inTyp, elemTyp) {
		return nil, fmt.Errorf("%s the function input type is not of type rune != %s",
			name, elemTyp)
	}
	res := sig.Results()
	if res.Len() != 1 {
		return nil, fmt.Errorf("%s, the function argument does not have a single result, but has %d resulting parameters", name, res.Len())
	}
	outTyp = res.At(0).Type()
	return outTyp, nil
}

func (this *gen) sliceInOut(name string, typs []types.Type) (inTyp types.Type, outTyp types.Type, err error) {
	sliceTyp, ok := typs[1].(*types.Slice)
	if !ok {
		return nil, nil, fmt.Errorf("%s, the second argument, %s, is not of type slice", name, this.TypeString(typs[1]))
	}
	sig, ok := typs[0].(*types.Signature)
	if !ok {
		return nil, nil, fmt.Errorf("%s, the second argument, %s, is not of type function", name, this.TypeString(typs[0]))
	}
	params := sig.Params()
	if params.Len() != 1 {
		return nil, nil, fmt.Errorf("%s, the second argument is a function, but wanted a function with one argument", name)
	}
	elemTyp := sliceTyp.Elem()
	inTyp = params.At(0).Type()
	if !types.Identical(inTyp, elemTyp) {
		return nil, nil, fmt.Errorf("%s the function input type and slice element type are different %s != %s",
			name, inTyp, elemTyp)
	}
	res := sig.Results()
	if res.Len() != 1 {
		return nil, nil, fmt.Errorf("%s, the function argument does not have a single result, but has %d resulting parameters", name, res.Len())
	}
	outTyp = res.At(0).Type()
	return inTyp, outTyp, nil
}

func (this *gen) Generate(typs []types.Type) error {
	switch typs[1].(type) {
	case *types.Slice:
		return this.genSlice(typs)
	case *types.Basic:
		return this.genString(typs)
	}
	return fmt.Errorf("unsupported type %s, not a slice or a string", typs[1])
}

func (this *gen) genSlice(typs []types.Type) error {
	name := this.GetFuncName(typs...)
	in, out, err := this.sliceInOut(name, typs)
	if err != nil {
		return err
	}
	this.Generating(typs...)
	p := this.printer
	inStr := this.TypeString(in)
	outStr := this.TypeString(out)
	p.P("")
	p.P("func %s(f func(%s) %s, list []%s) []%s {", name, inStr, outStr, inStr, outStr)
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

func (this *gen) genString(typs []types.Type) error {
	name := this.GetFuncName(typs...)
	out, err := this.stringOut(name, typs)
	if err != nil {
		return err
	}
	this.Generating(typs...)
	p := this.printer
	outStr := this.TypeString(out)
	p.P("")
	p.P("func %s(f func(rune) %s, ss string) []%s {", name, outStr, outStr)
	p.In()
	p.P("out := make([]%s, len([]rune(ss)))", outStr)
	p.P("for i, elem := range ss {")
	p.In()
	p.P("out[i] = f(elem)")
	p.Out()
	p.P("}")
	p.P("return out")
	p.Out()
	p.P("}")
	return nil
}
