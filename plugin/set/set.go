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

func (this *gen) Add(name string, typs []types.Type) (string, error) {
	if len(typs) != 1 {
		return "", fmt.Errorf("%s does not have one argument", name)
	}
	return this.SetFuncName(name, typs[0])
}

func (this *gen) Generate() error {
	for _, typs := range this.ToGenerate() {
		typ := typs[0]
		sliceType, ok := typ.(*types.Slice)
		if !ok {
			return fmt.Errorf("%s, the first argument, %s, is not of type slice", this.GetFuncName(typ), typ)
		}
		if err := this.genFuncFor(sliceType); err != nil {
			return err
		}
	}
	return nil
}

func (this *gen) genFuncFor(typ *types.Slice) error {
	p := this.printer
	this.Generating(typ)
	typeStr := this.TypeString(typ.Elem())
	p.P("")
	p.P("func %s(list []%s) map[%s]struct{} {", this.GetFuncName(typ), typeStr, typeStr)
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
