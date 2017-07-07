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
// The deriveUnique function returns a list of unique elements.
//   deriveUnique([]T) []T
// Example: https://github.com/awalterschulze/goderive/tree/master/example/plugin/unique
// deriveUnique also mutates the list in place.
package unique

import (
	"fmt"
	"go/types"

	"github.com/awalterschulze/goderive/derive"
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
		contains: deps["contains"],
	}
}

type gen struct {
	derive.TypesMap
	printer  derive.Printer
	contains derive.Dependency
}

func (this *gen) Add(name string, typs []types.Type) (string, error) {
	if len(typs) != 1 {
		return "", fmt.Errorf("%s does not have one argument", name)
	}
	return this.SetFuncName(name, typs[0])
}

func (this *gen) Generate(typs []types.Type) error {
	typ := typs[0]
	sliceType, ok := typ.(*types.Slice)
	if !ok {
		return fmt.Errorf("%s, the first argument, %s, is not of type slice", this.GetFuncName(typ), this.TypeString(typ))
	}
	return this.genFuncFor(sliceType)
}

func (this *gen) genFuncFor(typ *types.Slice) error {
	p := this.printer
	this.Generating(typ)
	typeStr := this.TypeString(typ)
	p.P("")
	p.P("func %s(list %s) %s {", this.GetFuncName(typ), typeStr, typeStr)
	p.In()
	p.P("if len(list) == 0 {")
	p.In()
	p.P("return nil")
	p.Out()
	p.P("}")
	p.P("u := 1")
	p.P("for i := 1; i < len(list); i++ {")
	p.In()
	p.P("if !%s(list[:u], list[i]) {", this.contains.GetFuncName(typ))
	p.In()
	p.P("if i != u {")
	p.In()
	p.P("list[u] = list[i]")
	p.Out()
	p.P("}")
	p.P("u++")
	p.Out()
	p.P("}")
	p.Out()
	p.P("}")
	p.P("return list[:u]")
	p.Out()
	p.P("}")
	return nil
}
