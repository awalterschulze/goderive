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

// Package union contains the implementation of the union plugin, which generates the deriveUnion function.
//   func deriveUnion([]T, []T) []T
//   func deriveUnion(map[T]struct{}, map[T]struct{}) map[T]struct{}
//
// Example: https://github.com/awalterschulze/goderive/tree/master/example/plugin/union
package union

import (
	"fmt"
	"go/types"

	"github.com/awalterschulze/goderive/derive"
)

// NewPlugin creates a new keys plugin.
// This function returns the plugin name, default prefix and a constructor for the union code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("union", "deriveUnion", New)
}

// New is a constructor for the union code generator.
// This generator should be reconstructed for each package.
func New(typesMap derive.TypesMap, p derive.Printer, deps map[string]derive.Dependency) derive.Generator {
	return &gen{
		TypesMap: typesMap,
		printer:  p,
		contains: deps["contains"],
	}
}

type gen struct {
	derive.TypesMap
	printer  derive.Printer
	contains derive.Dependency
}

func (g *gen) Add(name string, typs []types.Type) (string, error) {
	if len(typs) != 2 {
		return "", fmt.Errorf("%s does not have two arguments", name)
	}
	if !types.Identical(typs[0], typs[1]) {
		return "", fmt.Errorf("%s's two input types are not identical", name)
	}
	switch typ := typs[0].(type) {
	case *types.Slice:
	case *types.Map:
		if !types.Identical(typ.Elem(), types.NewStruct(nil, nil)) {
			return "", fmt.Errorf("%s takes an unsupported type: %s, map must be of type map[T]struct{}", name, typ)
		}
	default:
		return "", fmt.Errorf("%s takes an unsupported type: %s", name, typ)
	}
	return g.SetFuncName(name, typs[0])
}

func (g *gen) Generate(typs []types.Type) error {
	switch typ := typs[0].(type) {
	case *types.Slice:
		return g.genSlice(typ)
	case *types.Map:
		return g.genMap(typ)
	}
	return fmt.Errorf("%s, the argument type, %s, is not of type slice or map[T]struct{}", g.GetFuncName(typs[0]), typs[0])
}

func (g *gen) genMap(typ *types.Map) error {
	p := g.printer
	g.Generating(typ)
	name := g.GetFuncName(typ)
	typeStr := g.TypeString(typ.Key())
	p.P("")
	p.P("// %s returns the union of two maps, with respect to the keys.", name)
	p.P("// It does this by adding the keys to the first map.")
	p.P("func %s(union, that map[%s]struct{}) map[%s]struct{} {", name, typeStr, typeStr)
	p.In()
	p.P("for k := range that {")
	p.In()
	p.P("union[k] = struct{}{}")
	p.Out()
	p.P("}")
	p.P("return union")
	p.Out()
	p.P("}")
	return nil
}

func (g *gen) genSlice(typ *types.Slice) error {
	p := g.printer
	g.Generating(typ)
	name := g.GetFuncName(typ)
	typeStr := g.TypeString(typ.Elem())
	p.P("")
	p.P("// %s returns the union of the items of the two input lists.", name)
	p.P("// It does this by append items to the first list.")
	p.P("func %s(this, that []%s) []%s {", name, typeStr, typeStr)
	p.In()
	p.P("for i, v := range that {")
	p.In()
	p.P("if !%s(this, v) {", g.contains.GetFuncName(typ))
	p.In()
	p.P("this = append(this, that[i])")
	p.Out()
	p.P("}")
	p.Out()
	p.P("}")
	p.P("return this")
	p.Out()
	p.P("}")
	return nil
}
