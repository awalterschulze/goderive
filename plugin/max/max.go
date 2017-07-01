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

// Package max contains the implementation of the max plugin, which generates the deriveMax function.
// The deriveMax function returns the maximum value in a slice.
//   func deriveMax(default T, list []T) (max T)
package max

import (
	"fmt"
	"go/types"

	"github.com/awalterschulze/goderive/derive"
)

// NewPlugin creates a new max plugin.
// This function returns the plugin name, default prefix and a constructor for the max code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("max", "deriveMax", New)
}

// New is a constructor for the max code generator.
// This generator should be reconstructed for each package.
func New(typesMap derive.TypesMap, p derive.Printer, deps map[string]derive.Dependency) derive.Generator {
	return &gen{
		TypesMap: typesMap,
		printer:  p,
		compare:  deps["compare"],
	}
}

type gen struct {
	derive.TypesMap
	printer derive.Printer
	compare derive.Dependency
}

func (this *gen) Add(name string, typs []types.Type) (string, error) {
	if len(typs) != 2 {
		return "", fmt.Errorf("%s does not have two arguments", name)
	}
	sliceType, ok := typs[0].(*types.Slice)
	if !ok {
		return "", fmt.Errorf("%s, the first argument, %s, is not of type slice", name, typs[1])
	}
	if !types.AssignableTo(typs[1], sliceType.Elem()) {
		return "", fmt.Errorf("%s, the second argument, %s, is not is assignable to an element that of the slice type %s", name, typs[1], typs[0])
	}
	return this.SetFuncName(name, typs[0])
}

func (this *gen) Generate(typs []types.Type) error {
	typ := typs[0]
	sliceType, ok := typ.(*types.Slice)
	if !ok {
		return fmt.Errorf("%s, the first argument, %s, is not of type slice", this.GetFuncName(typ), typ)
	}
	return this.genFuncFor(sliceType)
}

func (this *gen) genFuncFor(typ *types.Slice) error {
	p := this.printer
	this.Generating(typ)
	etyp := typ.Elem()
	typeStr := this.TypeString(etyp)
	p.P("")
	p.P("func %s(list []%s, def %s) %s {", this.GetFuncName(typ), typeStr, typeStr, typeStr)
	p.In()
	p.P("if len(list) == 0 {")
	p.In()
	p.P("return def")
	p.Out()
	p.P("}")
	p.P("m := list[0]")
	p.P("list = list[1:]")
	p.P("for i, v := range list {")
	p.In()
	switch etyp.(type) {
	case *types.Basic:
		p.P("if v > m {")
	default:
		p.P("if %s(v, m) > 0 {", this.compare.GetFuncName(etyp))
	}
	p.In()
	p.P("m = list[i]")
	p.Out()
	p.P("}")
	p.Out()
	p.P("}")
	p.P("return m")
	p.Out()
	p.P("}")
	return nil
}
