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

// Package filter contains the implementation of the filter plugin, which generates the deriveFilter function.
//
// The deriveFilter function applies a predicate to each element of a list, returning a list of filtered results in the same order.
//   func deriveFilter(func (T) bool, []T) []T
package filter

import (
	"fmt"
	"go/types"

	"github.com/awalterschulze/goderive/derive"
)

// NewPlugin creates a new filter plugin.
// This function returns the plugin name, default prefix and a constructor for the filter code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("filter", "deriveFilter", New)
}

// New is a constructor for the filter code generator.
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
	if !types.Identical(outTyp, types.Typ[types.Bool]) {
		return "", fmt.Errorf("%s, the function argument has a single result, but %s is not a bool", name, outTyp)
	}
	return this.SetFuncName(name, inTyp)
}

func (this *gen) Generate(typs []types.Type) error {
	return this.genFuncFor(typs[0])
}

func (this *gen) genFuncFor(in types.Type) error {
	p := this.printer
	this.Generating(in)
	inStr := this.TypeString(in)
	p.P("")
	p.P("func %s(pred func(%s) bool, list []%s) []%s {", this.GetFuncName(in), inStr, inStr, inStr)
	p.In()
	p.P("out := make([]%s, 0, len(list))", inStr)
	p.P("for i, elem := range list {")
	p.In()
	p.P("if pred(elem) {")
	p.In()
	p.P("out = append(out, list[i])")
	p.Out()
	p.P("}")
	p.Out()
	p.P("}")
	p.P("return out")
	p.Out()
	p.P("}")
	return nil
}
