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

// Package set contains the implementation of the set plugin, which generates the deriveSet function.
//   func deriveSet([]T) map[T]struct{}
//
// Example: https://github.com/awalterschulze/goderive/tree/master/example/plugin/set
package set

import (
	"fmt"
	"go/types"

	"github.com/awalterschulze/goderive/derive"
)

// NewPlugin creates a new keys plugin.
// This function returns the plugin name, default prefix and a constructor for the set code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("set", "deriveSet", New)
}

// New is a constructor for the set code generator.
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
		return fmt.Errorf("%s, the first argument, %s, is not of type slice", g.GetFuncName(typ), typ)
	}
	return g.genFuncFor(sliceType)
}

func (g *gen) genFuncFor(typ *types.Slice) error {
	p := g.printer
	g.Generating(typ)
	name := g.GetFuncName(typ)
	typeStr := g.TypeString(typ.Elem())
	p.P("")
	p.P("// %s returns the input list as a map with the items of the list as the keys of the map.", name)
	p.P("func %s(list []%s) map[%s]struct{} {", name, typeStr, typeStr)
	p.In()
	p.P("set := make(map[%s]struct{}, len(list))", typeStr)
	p.P("for _, v := range list {")
	p.In()
	p.P("set[v] = struct{}{}")
	p.Out()
	p.P("}")
	p.P("return set")
	p.Out()
	p.P("}")
	return nil
}
