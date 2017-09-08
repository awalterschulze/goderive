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

// Package clone contains the implementation of the clone plugin, which generates the deriveClone function.
//
// The deriveClone function is a maintainable and fast way to implement fast"ish" clone functions.
//   func deriveClone(T) T
// I say fast"ish", since deriveClone creates a totally new copy of the value, whereas deepcopy reuses as much as of the memory that has been allocated by the destintation value.
//
// Supported types:
//	- basic types
//	- named structs
//	- slices
//	- maps
//	- pointers to these types
//	- private fields of structs in external packages (using reflect and unsafe)
//	- and many more
// Unsupported types:
//	- chan
//	- interface
//	- function
//	- unnamed structs, which are not comparable with the == operator
package clone

import (
	"fmt"
	"go/types"

	"github.com/awalterschulze/goderive/derive"
)

// NewPlugin creates a new clone plugin.
// This function returns the plugin name, default prefix and a constructor for the clone code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("clone", "deriveClone", New)
}

// New is a constructor for the clone code generator.
// This generator should be reconstructed for each package.
func New(typesMap derive.TypesMap, p derive.Printer, deps map[string]derive.Dependency) derive.Generator {
	return &gen{
		TypesMap: typesMap,
		printer:  p,
		deepcopy: deps["deepcopy"],
	}
}

type gen struct {
	derive.TypesMap
	printer  derive.Printer
	deepcopy derive.Dependency
}

func (g *gen) Add(name string, typs []types.Type) (string, error) {
	if len(typs) != 1 {
		return "", fmt.Errorf("%s does not have one argument", name)
	}
	return g.SetFuncName(name, typs[0])
}

func (g *gen) Generate(typs []types.Type) error {
	return g.genFuncFor(typs[0])
}

func (g *gen) genFuncFor(in types.Type) error {
	p := g.printer
	g.Generating(in)
	inStr := g.TypeString(in)
	p.P("")
	p.P("// %s returns a clone of the src parameter.", g.GetFuncName(in))
	p.P("func %s(src %s) %s {", g.GetFuncName(in), inStr, inStr)
	p.In()
	switch ttyp := in.Underlying().(type) {
	case *types.Pointer:
		p.P("if src == nil {")
		p.In()
		p.P("return nil")
		p.Out()
		p.P("}")
		p.P("dst := new(%s)", g.TypeString(ttyp.Elem()))
		p.P("%s(dst, src)", g.deepcopy.GetFuncName(in))
		p.P("return dst")
	case *types.Slice:
		p.P("if src == nil {")
		p.In()
		p.P("return nil")
		p.Out()
		p.P("}")
		p.P("dst := make(%s, len(src))", g.TypeString(in))
		p.P("%s(dst, src)", g.deepcopy.GetFuncName(in))
		p.P("return dst")
	case *types.Map:
		p.P("if src == nil {")
		p.In()
		p.P("return nil")
		p.Out()
		p.P("}")
		p.P("dst := make(%s)", g.TypeString(in))
		p.P("%s(dst, src)", g.deepcopy.GetFuncName(in))
		p.P("return dst")
	default:
		p.P("dst := new(%s)", inStr)
		p.P("%s(dst, &src)", g.deepcopy.GetFuncName(types.NewPointer(in)))
		p.P("return *dst")
	}
	p.Out()
	p.P("}")
	return nil
}
