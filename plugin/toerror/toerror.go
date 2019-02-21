//  Copyright 2019 Ingun Jon
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

// Package toerror contains the implementation of the toerror plugin, which generates the deriveToError function.
//
// The deriveToError function transforms return type of a function from (a, bool) into (a, error).
//   deriveToError(e error, f func(...) (T, bool)) func(...) (T, error)
package toerror

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/awalterschulze/goderive/derive"
)

// NewPlugin creates a new curry plugin.
// This function returns the plugin name, default prefix and a constructor for the curry code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("toerror", "deriveToError", New)
}

// New is a constructor for the curry code generator.
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
	if len(typs) != 2 {
		return "", fmt.Errorf("%s does not have two argument. Expected (error, func (a) (b, bool))", name)
	}
	errTyp := typs[0]
	if !derive.IsError(errTyp) {
		return "", fmt.Errorf("First parameter should be of type error")
	}
	funcTyp := typs[1]
	sig, ok := funcTyp.(*types.Signature)
	if !ok {
		return "", fmt.Errorf("%s, 2nd param %s, is not of type function", name, g.TypeString(typs[1]))
	}
	results := sig.Results()
	if results.Len() <= 0 {
		return "", fmt.Errorf("%s, given function must return at least 1 bool type", name)
	}
	if !types.Identical(results.At(results.Len()-1).Type(), types.Typ[types.Bool]) {
		return "", fmt.Errorf("%s, given function must return bool as last return type. (got %v)", name, results.String())
	}
	return g.SetFuncName(name, errTyp, funcTyp)
}

func (g *gen) Generate(typs []types.Type) error {
	g.Generating(typs...)
	errTyp := typs[0].(*types.Named)
	funcTyp := typs[1].(*types.Signature)
	return g.genFuncFor(g.GetFuncName(typs...), errTyp, funcTyp)
}

func varnames(tup *types.Tuple) []string {
	as := make([]string, tup.Len())
	for i := 0; i < tup.Len(); i++ {
		as[i] = tup.At(i).Name()
	}
	return as
}
func stripVarName(v *types.Var) *types.Var {
	return types.NewVar(v.Pos(), v.Pkg(), "", v.Type())
}
func outs(num int, last string) string {
	outs := make([]string, num, num)
	for i := 0; i < num-1; i++ {
		outs[i] = fmt.Sprintf("out%d", i)
	}
	outs[num-1] = last
	return strings.Join(outs, ", ")
}
func (g *gen) genFuncFor(deriveFuncName string, etyp *types.Named, ftyp *types.Signature) error {
	p := g.printer
	rlen := ftyp.Results().Len()
	mutable := make([]*types.Var, rlen, rlen)
	for i := 0; i < rlen-1; i++ {
		mutable[i] = stripVarName(ftyp.Results().At(i))
	}
	mutable[rlen-1] = types.NewVar(ftyp.Results().At(rlen-1).Pos(), nil, "", etyp)
	newResultType := types.NewTuple(mutable...)
	newSigType := types.NewSignature(nil, ftyp.Params(), newResultType, ftyp.Variadic())

	p.P("")
	p.P("// %s transforms sum-bool type into sum-error type. Main purpose is to make the given function composable. It returns given error when the result of the function is false.", deriveFuncName)
	p.P("func %s(err error, f %s) %s {", deriveFuncName, g.TypeString(ftyp), g.TypeString(newSigType))
	p.In()
	p.P("return %s {", g.TypeString(newSigType))
	p.In()
	as := varnames(newSigType.Params())
	p.P("%s := f(%s)", outs(rlen, "success"), strings.Join(as, ", "))
	p.P("if success {")
	p.In()
	p.P("return %s", outs(rlen, "nil"))
	p.Out()
	p.P("}")
	p.P("return %s", outs(rlen, "err"))
	p.Out()
	p.P("}")
	p.Out()
	p.P("}")
	return nil
}
