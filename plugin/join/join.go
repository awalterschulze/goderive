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

// Package join contains the implementation of the join plugin, which generates the deriveJoin function.
//
// The deriveJoin function joins a slice of slices into a single slice.
//
// More things to come:
//	- currently only slices are supported, think about supporting other types and not just slices
//	- what about []string and not just [][]string as in the current example.
package join

import (
	"fmt"
	"go/types"

	"github.com/awalterschulze/goderive/derive"
)

// NewPlugin creates a new join plugin.
// This function returns the plugin name, default prefix and a constructor for the join code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("join", "deriveJoin", New)
}

// New is a constructor for the join code generator.
// This generator should be reconstructed for each package.
func New(typesMap derive.TypesMap, p derive.Printer, deps map[string]derive.Dependency) derive.Generator {
	return &gen{
		TypesMap: typesMap,
		printer:  p,
	}
}

type gen struct {
	derive.TypesMap
	printer  derive.Printer
	bytesPkg derive.Import
}

func (this *gen) Add(name string, typs []types.Type) (string, error) {
	if len(typs) != 1 {
		return "", fmt.Errorf("%s does not have one argument", name)
	}
	sliceTyp, ok := typs[0].(*types.Slice)
	if !ok {
		return "", fmt.Errorf("%s, the argument, %s, is not of type slice", name, typs[0])
	}
	sliceOfSliceTyp, ok := sliceTyp.Elem().(*types.Slice)
	if !ok {
		return "", fmt.Errorf("%s, the argument, %s, is not of type slice of slice", name, typs[0])
	}
	elemType := sliceOfSliceTyp.Elem()
	return this.SetFuncName(name, elemType)
}

func (this *gen) Generate(typs []types.Type) error {
	return this.genFuncFor(typs[0])
}

func (this *gen) genFuncFor(typ types.Type) error {
	p := this.printer
	this.Generating(typ)
	typStr := this.TypeString(typ)
	p.P("")
	p.P("func %s(list [][]%s) []%s {", this.GetFuncName(typ), typStr, typStr)
	p.In()
	p.P("if list == nil {")
	p.In()
	p.P("return nil")
	p.Out()
	p.P("}")
	p.P("l := 0")
	p.P("for _, elem := range list {")
	p.In()
	p.P("l += len(elem)")
	p.Out()
	p.P("}")
	p.P("res := make([]%s, 0, l)", typStr)
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
