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

// Package uncurry contains the implementation of the uncurry plugin, which generates the deriveUncurry function.
//
// The deriveUncurry function uncurries the input function.
//   deriveUncurry(f func(A) func(B, ...) T) func(A, B, ...) T
package uncurry

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/awalterschulze/goderive/derive"
)

// NewPlugin creates a new uncurry plugin.
// This function returns the plugin name, default prefix and a constructor for the uncurry code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("uncurry", "deriveUncurry", New)
}

// New is a constructor for the uncurry code generator.
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
	if params.Len() != 1 {
		return "", fmt.Errorf("%s, the first argument is a function, but wanted a function with one argument", name)
	}
	if sig.Results().Len() != 1 {
		return "", fmt.Errorf("%s, expected 1 result for the input function", name)
	}
	retVar := sig.Results().At(0)
	retSig, ok := retVar.Type().(*types.Signature)
	if !ok {
		return "", fmt.Errorf("%s, does not return a function", name)
	}
	retSig = derive.RenameBlankIdentifierWith(retSig, "innerParam_")
	newTup := types.NewTuple(types.NewVar(retVar.Pos(), retVar.Pkg(), retVar.Name(), retSig))
	sig = types.NewSignature(sig.Recv(), sig.Params(), newTup, sig.Variadic())
	return g.SetFuncName(name, derive.RenameBlankIdentifier(sig))
}

func (g *gen) Generate(typs []types.Type) error {
	sig, ok := typs[0].(*types.Signature)
	if !ok {
		return fmt.Errorf("%s, the first argument, %s, is not of type function", g.GetFuncName(typs[0]), g.TypeString(typs[0]))
	}
	return g.genFuncFor(sig)
}

func uncurrySig(sig *types.Signature) (*types.Signature, *types.Tuple) {
	params := vars(sig.Params())
	res := vars(sig.Results())[0]
	ressig := res.Type().(*types.Signature)
	resparams := vars(ressig.Params())
	newvars := append(params, resparams...)
	newparams := types.NewTuple(newvars...)
	f := types.NewSignature(nil, newparams, ressig.Results(), ressig.Variadic())
	return f, ressig.Params()
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
	gtyp, styp := uncurrySig(ftyp)
	gStr := g.TypeString(gtyp)
	firstStr := varnames(ftyp.Params())
	secondStr := varnames(styp)
	p.P("")
	p.P("// %s combines a function that returns a function, into one function.", name)
	p.P("func %s(f %s) %s {", name, fStr, gStr)
	p.In()
	p.P("return %s {", gStr)
	p.In()
	p.P("return f(%s)(%s)", strings.Join(firstStr, ", "), strings.Join(secondStr, ", "))
	p.Out()
	p.P("}")
	p.Out()
	p.P("}")
	return nil
}
