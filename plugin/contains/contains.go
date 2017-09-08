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

// Package contains contains the implementation of the contains plugin, which generates the deriveContains function.
//
// The deriveContains function returns whether a value is contained in a slice.
//   func deriveContains([]T, T) bool
//
// Example: https://github.com/awalterschulze/goderive/tree/master/example/plugin/contains
package contains

import (
	"fmt"
	"go/types"

	"github.com/awalterschulze/goderive/derive"
)

// NewPlugin creates a new contains plugin.
// This function returns the plugin name, default prefix and a constructor for the contains code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("contains", "deriveContains", New)
}

// New is a constructor for the contains code generator.
// This generator should be reconstructed for each package.
func New(typesMap derive.TypesMap, p derive.Printer, deps map[string]derive.Dependency) derive.Generator {
	return &gen{
		TypesMap: typesMap,
		printer:  p,
		equal:    deps["equal"],
	}
}

type gen struct {
	derive.TypesMap
	printer derive.Printer
	equal   derive.Dependency
}

func (g *gen) Add(name string, typs []types.Type) (string, error) {
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
	return g.SetFuncName(name, typs[0])
}

func (g *gen) Generate(typs []types.Type) error {
	typ := typs[0]
	sliceType, ok := typ.(*types.Slice)
	if !ok {
		return fmt.Errorf("%s, the first argument, %s, is not of type slice", g.GetFuncName(typ), typ)
	}
	return g.genFuncFor(sliceType)
}

func canEqual(tt types.Type) bool {
	t := tt.Underlying()
	switch typ := t.(type) {
	case *types.Basic:
		return typ.Kind() != types.UntypedNil
	case *types.Struct:
		for i := 0; i < typ.NumFields(); i++ {
			f := typ.Field(i)
			ft := f.Type()
			if !canEqual(ft) {
				return false
			}
		}
		return true
	case *types.Array:
		return canEqual(typ.Elem())
	}
	return false
}

func (g *gen) genFuncFor(typ *types.Slice) error {
	p := g.printer
	g.Generating(typ)
	name := g.GetFuncName(typ)
	etyp := typ.Elem()
	typeStr := g.TypeString(etyp)
	p.P("")
	p.P("// %s returns whether the item is contained in the list.", name)
	p.P("func %s(list []%s, item %s) bool {", name, typeStr, typeStr)
	p.In()
	p.P("for _, v := range list {")
	p.In()
	if canEqual(etyp) {
		p.P("if v == item {")
	} else {
		p.P("if %s(v, item) {", g.equal.GetFuncName(etyp, etyp))
	}
	p.In()
	p.P("return true")
	p.Out()
	p.P("}")
	p.Out()
	p.P("}")
	p.P("return false")
	p.Out()
	p.P("}")
	return nil
}
