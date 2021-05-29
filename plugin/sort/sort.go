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

// Package sort contains the implementation of the sort plugin, which generates the deriveSort function.
//
// The deriveSort function is useful for deterministically ranging over maps when used with deriveKeys.
// deriveSort supports only the types that deriveCompare supports, since it uses it for sorting.
//
// Example: https://github.com/awalterschulze/goderive/tree/master/example/plugin/sort
//
// Even though sort returns a list it also mutates the input list.
package sort

import (
	"fmt"
	"go/types"

	"github.com/awalterschulze/goderive/derive"
)

// NewPlugin creates a new sort plugin.
// This function returns the plugin name, default prefix and a constructor for the sort code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("sort", "deriveSort", New)
}

// New is a constructor for the sort code generator.
// This generator should be reconstructed for each package.
func New(typesMap derive.TypesMap, p derive.Printer, deps map[string]derive.Dependency) derive.Generator {
	return &gen{
		TypesMap: typesMap,
		printer:  p,
		sortPkg:  p.NewImport("sort", "sort"),
		compare:  deps["compare"],
	}
}

type gen struct {
	derive.TypesMap
	printer derive.Printer
	sortPkg derive.Import
	compare derive.Dependency
}

func (g *gen) Add(name string, typs []types.Type) (string, error) {
	if len(typs) != 1 {
		return "", fmt.Errorf("%s does not have one argument", name)
	}
	return g.SetFuncName(name, typs[0])
}

func (g *gen) Generate(typs []types.Type) error {
	typ := typs[0]
	sliceType, ok := typ.(*types.Slice)
	if !ok {
		return fmt.Errorf("%s, the first argument, %s, is not of type slice", g.GetFuncName(typ), g.TypeString(typ))
	}
	return g.genFuncFor(sliceType)
}

func (g *gen) genFuncFor(typ *types.Slice) error {
	p := g.printer
	g.Generating(typ)
	name := g.GetFuncName(typ)
	typeStr := g.TypeString(typ)
	p.P("")
	p.P("// %s sorts the slice inplace and also returns it.", name)
	p.P("func %s(list %s) %s {", name, typeStr, typeStr)
	p.In()
	if err := g.printSortFunc(typ); err != nil {
		return err
	}
	p.P("return list")
	p.Out()
	p.P("}")
	return nil
}

func (g *gen) printSortFunc(typ *types.Slice) error {
	p := g.printer
	etyp := typ.Elem()
	switch ttyp := etyp.(type) {
	case *types.Basic:
		switch ttyp.Kind() {
		case types.String:
			p.P(g.sortPkg() + ".Strings(list)")
			return nil
		case types.Float64:
			p.P(g.sortPkg() + ".Float64s(list)")
			return nil
		case types.Int:
			p.P(g.sortPkg() + ".Ints(list)")
			return nil
		}
	}

	switch ttyp := etyp.Underlying().(type) {
	case *types.Basic:
		switch ttyp.Kind() {
		case types.Complex64, types.Complex128, types.Bool:
			p.P(g.sortPkg() + ".Slice(list, func(i, j int) bool { return " + g.compare.GetFuncName(ttyp, ttyp) + "(list[i], list[j]) < 0 })")
		default:
			p.P(g.sortPkg() + ".Slice(list, func(i, j int) bool { return list[i] < list[j] })")
		}
	case *types.Pointer, *types.Struct, *types.Slice, *types.Array, *types.Map:
		p.P(g.sortPkg() + ".Slice(list, func(i, j int) bool { return " + g.compare.GetFuncName(etyp, etyp) + "(list[i], list[j]) < 0 })")
	default:
		return fmt.Errorf("unsupported compare type: %s", g.TypeString(typ))
	}

	return nil
}
