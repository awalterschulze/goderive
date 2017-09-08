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

// Package min contains the implementation of the min plugin, which generates the deriveMin function.
//
// The deriveMin function returns the minimum of two arguments.
//   func deriveMin(T, T) T
//
// deriveMin is a generic version of
//   math.Min(x, y float64) float64
//
// deriveMin is preferable over abusing math.Min, for not float64 types:
// https://mrekucci.blogspot.nl/2015/07/dont-abuse-mathmax-mathmin.html
//
// It can also return the minimum element in a list.
//   func deriveMin(list []T, default T) (min T)
//
// A default value is provided for the empty list.
//
// Example: https://github.com/awalterschulze/goderive/tree/master/example/plugin/min
package min

import (
	"fmt"
	"go/types"

	"github.com/awalterschulze/goderive/derive"
)

// NewPlugin creates a new min plugin.
// This function returns the plugin name, default prefix and a constructor for the min code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("min", "deriveMin", New)
}

// New is a constructor for the min code generator.
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
	name := g.GetFuncName(typ, typ2)
	typeStr := g.TypeString(typ)
	p.P("")
	p.P("// %s returns the mimimum of the two input values.", name)
	p.P("func %s(a, b %s) %s {", name, typeStr, typeStr)
	p.In()
	switch typ.(type) {
	case *types.Basic:
		p.P("if a < b {")
	default:
		p.P("if %s(a, b) < 0 {", g.compare.GetFuncName(typ, typ))
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
	etyp := typ.Elem()
	name := g.GetFuncName(typ, typ2)
	typeStr := g.TypeString(etyp)
	p.P("")
	p.P("// %s returns the minimum value from the list, or the default value if the list is empty.", name)
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
		p.P("if v < m {")
	default:
		p.P("if %s(v, m) < 0 {", g.compare.GetFuncName(etyp, etyp))
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
