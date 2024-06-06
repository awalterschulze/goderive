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
//
// The deriveMax function returns the maximum of two arguments.
//   func deriveMax(T, T) T
//
// deriveMax is a generic version of
//   math.Max(x, y float64) float64
//
// deriveMax is preferable over abusing math.Max, for not float64 types:
// https://mrekucci.blogspot.nl/2015/07/dont-abuse-mathmax-mathmin.html
//
// It can also return the maximum element in a list.
//   func deriveMax(list []T, default T) (max T)
//
// A default value is provided for the empty list.
//
// Example: https://github.com/ndeloof/goderive/tree/master/example/plugin/max
package max

import (
	"fmt"
	"go/types"

	"github.com/ndeloof/goderive/derive"
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

func (g *gen) Add(name string, typs []types.Type) (string, error) {
	if len(typs) != 2 {
		return "", fmt.Errorf("%s does not have two arguments", name)
	}
	if types.Identical(typs[0], typs[1]) {
		return g.SetFuncName(name, typs[0], typs[1])
	}
	sliceType, ok := typs[0].(*types.Slice)
	if !ok {
		return "", fmt.Errorf("%s, the first argument, %s, is not of type slice", name, typs[0])
	}
	if !types.AssignableTo(typs[1], sliceType.Elem()) {
		return "", fmt.Errorf("%s, the second argument, %s, is not is assignable to an element that of the slice type %s", name, typs[1], typs[0])
	}
	return g.SetFuncName(name, typs[0], typs[1])
}

func (g *gen) Generate(typs []types.Type) error {
	if types.Identical(typs[0], typs[1]) {
		return g.genTwo(typs[0], typs[1])
	}
	sliceType, ok := typs[0].(*types.Slice)
	if !ok {
		return fmt.Errorf("%s, the first argument, %s, is not of type slice", g.GetFuncName(typs[0], typs[1]), typs[0])
	}
	return g.genSlice(sliceType, typs[1])
}

func (g *gen) genTwo(typ, typ2 types.Type) error {
	p := g.printer
	g.Generating(typ, typ2)
	typeStr := g.TypeString(typ)
	name := g.GetFuncName(typ, typ2)
	p.P("")
	p.P("// %s returns the maximum of the two input values.", name)
	p.P("func %s(a, b %s) %s {", name, typeStr, typeStr)
	p.In()
	switch typ.(type) {
	case *types.Basic:
		p.P("if a > b {")
	default:
		p.P("if %s(a, b) > 0 {", g.compare.GetFuncName(typ, typ))
	}
	p.In()
	p.P("return a")
	p.Out()
	p.P("}")
	p.P("return b")
	p.Out()
	p.P("}")
	return nil
}

func (g *gen) genSlice(typ *types.Slice, typ2 types.Type) error {
	p := g.printer
	g.Generating(typ, typ2)
	name := g.GetFuncName(typ, typ2)
	etyp := typ.Elem()
	typeStr := g.TypeString(etyp)
	p.P("")
	p.P("// %s returns the maximum value from the input list and the default value, if the list is empty.", name)
	p.P("func %s(list []%s, def %s) %s {", name, typeStr, typeStr, typeStr)
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
		p.P("if %s(v, m) > 0 {", g.compare.GetFuncName(etyp, etyp))
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
