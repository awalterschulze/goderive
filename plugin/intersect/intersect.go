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

// Package intersect contains the implementation of the intersect plugin, which generates the deriveIntersect function.
//   func deriveIntersect([]T, []T) []T
//   func deriveIntersect(map[T]struct{}, map[T]struct{}) map[T]struct{}
//
// Example: https://github.com/awalterschulze/goderive/tree/master/example/plugin/intersect
package intersect

import (
	"fmt"
	"go/types"

	"github.com/awalterschulze/goderive/derive"
)

// NewPlugin creates a new keys plugin.
// This function returns the plugin name, default prefix and a constructor for the intersect code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("intersect", "deriveIntersect", New)
}

// New is a constructor for the intersect code generator.
// This generator should be reconstructed for each package.
func New(typesMap derive.TypesMap, p derive.Printer, deps map[string]derive.Dependency) derive.Generator {
	return &gen{
		TypesMap: typesMap,
		printer:  p,
		contains: deps["contains"],
		min:      deps["min"],
	}
}

type gen struct {
	derive.TypesMap
	printer  derive.Printer
	contains derive.Dependency
	min      derive.Dependency
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
	typeStr := g.TypeString(typ.Key())
	p.P("")
	p.P("func %s(this, that map[%s]struct{}) map[%s]struct{} {", g.GetFuncName(typ), typeStr, typeStr)
	p.In()
	minFunc := g.min.GetFuncName(types.Typ[types.Int], types.Typ[types.Int])
	p.P("intersect := make(map[%s]struct{}, %s(len(this), len(that)))", typeStr, minFunc)
	p.P("for k := range this {")
	p.In()
	p.P("if _, ok := that[k]; ok {")
	p.In()
	p.P("intersect[k] = struct{}{}")
	p.Out()
	p.P("}")
	p.Out()
	p.P("}")
	p.P("return intersect")
	p.Out()
	p.P("}")
	return nil
}

func (g *gen) genSlice(typ *types.Slice) error {
	p := g.printer
	g.Generating(typ)
	typeStr := g.TypeString(typ.Elem())
	p.P("")
	p.P("func %s(this, that []%s) []%s {", g.GetFuncName(typ), typeStr, typeStr)
	p.In()
	minFunc := g.min.GetFuncName(types.Typ[types.Int], types.Typ[types.Int])
	p.P("intersect := make([]%s, 0, %s(len(this), len(that)))", typeStr, minFunc)
	p.P("for i, v := range this {")
	p.In()
	p.P("if %s(that, v) {", g.contains.GetFuncName(typ))
	p.In()
	p.P("intersect = append(intersect, this[i])")
	p.Out()
	p.P("}")
	p.Out()
	p.P("}")
	p.P("return intersect")
	p.Out()
	p.P("}")
	return nil
}
