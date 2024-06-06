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

// Package unique contains the implementation of the unique plugin, which generates the deriveUnique function.
//
// The deriveUnique function returns a list of unique elements.
//   deriveUnique([]T) []T
//
// Example: https://github.com/ndeloof/goderive/tree/master/example/plugin/unique
//
// deriveUnique mutates the list in place.
package unique

import (
	"fmt"
	"go/types"

	"github.com/ndeloof/goderive/derive"
)

// NewPlugin creates a new unique plugin.
// This function returns the plugin name, default prefix and a constructor for the unique code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("unique", "deriveUnique", New)
}

// New is a constructor for the unique code generator.
// This generator should be reconstructed for each package.
func New(typesMap derive.TypesMap, p derive.Printer, deps map[string]derive.Dependency) derive.Generator {
	return &gen{
		TypesMap: typesMap,
		printer:  p,
		hash:     deps["hash"],
		keys:     deps["keys"],
		set:      deps["set"],
		equal:    deps["equal"],
	}
}

type gen struct {
	derive.TypesMap
	printer derive.Printer
	hash    derive.Dependency
	keys    derive.Dependency
	set     derive.Dependency
	equal   derive.Dependency
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
	p.P("// %s returns a list containing only the unique items from the input list.", name)
	p.P("// It does this by reusing the input list.")
	p.P("func %s(list %s) %s {", name, typeStr, typeStr)
	p.In()
	p.P("if len(list) == 0 {")
	p.In()
	p.P("return nil")
	p.Out()
	p.P("}")
	if derive.IsComparable(typ.Elem()) {
		maptyp := types.NewMap(typ.Elem(), types.NewStruct(nil, nil))
		p.P("return %s(%s(list))", g.keys.GetFuncName(maptyp), g.set.GetFuncName(typ))
		p.Out()
		p.P("}")
		return nil
	}
	p.P("table := make(map[uint64][]int)")
	p.P("u := 0")
	p.P("for i := 0; i < len(list); i++ {")
	p.In()
	p.P("contains := false")
	p.P("hash := %s(list[i])", g.hash.GetFuncName(typ.Elem()))
	p.P("indexes := table[hash]")
	p.P("for _, index := range indexes {")
	p.In()
	p.P("if %s(list[index], list[i]) {", g.equal.GetFuncName(typ.Elem(), typ.Elem()))
	p.In()
	p.P("contains = true")
	p.P("break")
	p.Out()
	p.P("}")
	p.Out()
	p.P("}")
	p.P("if contains {")
	p.In()
	p.P("continue")
	p.Out()
	p.P("}")
	p.P("if i != u {")
	p.In()
	p.P("list[u] = list[i]")
	p.Out()
	p.P("}")
	p.P("table[hash] = append(table[hash], u)")
	p.P("u++")
	p.Out()
	p.P("}")
	p.P("return list[:u]")
	p.Out()
	p.P("}")
	return nil
}
