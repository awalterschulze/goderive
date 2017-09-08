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

// Package tuple contains the implementation of the tuple plugin, which generates the deriveTuple function.
//
// The deriveTuple function takes its input parameters and returns a function that returns those parameters.
//   deriveTuple(A, B, ...) func() (A, B, ...)
// deriveTuple is useful, since a tuple is not a first class citizen in go, but a function that returns a tuple is.
package tuple

import (
	"fmt"
	"go/types"
	"strconv"
	"strings"

	"github.com/awalterschulze/goderive/derive"
)

// NewPlugin creates a new tuple plugin.
// This function returns the plugin name, default prefix and a constructor for the tuple code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("tuple", "deriveTuple", New)
}

// New is a constructor for the tuple code generator.
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
	if len(typs) == 0 {
		return "", fmt.Errorf("%s has zero arguments", name)
	}
	if len(typs) == 1 {
		if tup, ok := typs[0].(*types.Tuple); ok {
			tuptypes := make([]types.Type, tup.Len())
			for i := range tuptypes {
				tuptypes[i] = tup.At(i).Type()
			}
			return g.SetFuncName(name, tuptypes...)
		}
	}
	return g.SetFuncName(name, typs...)
}

func (g *gen) Generate(typs []types.Type) error {
	return g.genFuncFor(typs)
}

func (g *gen) genFuncFor(typs []types.Type) error {
	p := g.printer
	g.Generating(typs...)
	typStrs := make([]string, len(typs))
	paramStrs := make([]string, len(typs))
	varStrs := make([]string, len(typs))
	for i, t := range typs {
		typStrs[i] = g.TypeString(t)
		varStrs[i] = "v" + strconv.Itoa(i)
		paramStrs[i] = varStrs[i] + " " + typStrs[i]
	}
	typStr := strings.Join(typStrs, ", ")
	if len(typs) > 1 {
		typStr = "(" + typStr + ")"
	}
	name := g.GetFuncName(typs...)
	p.P("")
	p.P("// %s returns a function, which returns the input values.", name)
	p.P("// Since tuples are not first class citizens in Go, this is a way to fake it, because functions that return tuples are first class citizens.")
	p.P("func %s(%s) func() %s {", name, strings.Join(paramStrs, ", "), typStr)
	p.In()
	p.P("return func() %s {", typStr)
	p.In()
	p.P("return %s", strings.Join(varStrs, ", "))
	p.Out()
	p.P("}")
	p.Out()
	p.P("}")
	return nil
}
