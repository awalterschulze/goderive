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

// Package flip contains the implementation of the flip plugin, which generates the deriveFlip function.
//
// The deriveFlip function flips the first two parameters of the input function.
//   deriveFlip(f func(A, B, ...) T) func(B, A, ...) T
package flip

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/awalterschulze/goderive/derive"
)

// NewPlugin creates a new flip plugin.
// This function returns the plugin name, default prefix and a constructor for the flip code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("flip", "deriveFlip", New)
}

// New is a constructor for the flip code generator.
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
	sig, ok := typs[0].(*types.Signature)
	if !ok {
		return "", fmt.Errorf("%s, the first argument, %s, is not of type function", name, g.TypeString(typs[0]))
	}
	params := sig.Params()
	if params.Len() < 2 {
		return "", fmt.Errorf("%s, the first argument is a function, but wanted a function with more than one argument", name)
	}
	return g.SetFuncName(name, sig)
}

func (g *gen) Generate(typs []types.Type) error {
	sig, ok := typs[0].(*types.Signature)
	if !ok {
		return fmt.Errorf("%s, the first argument, %s, is not of type function", g.GetFuncName(typs[0]), g.TypeString(typs[0]))
	}
	return g.genFuncFor(sig)
}

func flipSig(sig *types.Signature) *types.Signature {
	flipped := vars(sig.Params())
	flipped[0], flipped[1] = flipped[1], flipped[0]
	ps := types.NewTuple(flipped...)
	return types.NewSignature(sig.Recv(), ps, sig.Results(), sig.Variadic())
}

func vars(tup *types.Tuple) []*types.Var {
	vars := make([]*types.Var, tup.Len())
	for i := range vars {
		vars[i] = tup.At(i)
	}
	return vars
}

func varnames(tup *types.Tuple) []string {
	as := make([]string, tup.Len())
	for i := 0; i < tup.Len(); i++ {
		as[i] = tup.At(i).Name()
	}
	return as
}

func (g *gen) genFuncFor(ftyp *types.Signature) error {
	p := g.printer
	g.Generating(ftyp)
	fStr := g.TypeString(ftyp)
	name := g.GetFuncName(ftyp)
	gtyp := flipSig(ftyp)
	gStr := g.TypeString(gtyp)
	p.P("")
	p.P("// %s returns the input function, but where first two parameters are flipped.", name)
	p.P("func %s(f %s) %s {", name, fStr, gStr)
	p.In()
	p.P("return %s {", gStr)
	p.In()
	as := varnames(ftyp.Params())
	p.P("return f(%s)", strings.Join(as, ", "))
	p.Out()
	p.P("}")
	p.Out()
	p.P("}")
	return nil
}
