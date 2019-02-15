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
	if results.Len() != 2 {
		return "", fmt.Errorf("%s, return type of the given function must be a tuple length of 2", name)
	}
	if !types.Identical(results.At(1).Type(), types.Typ[types.Bool]) {
		return "", fmt.Errorf("%s expected secondary bool return type for the given function. (got %v)", name, results.String())
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

func (g *gen) genFuncFor(deriveFuncName string, etyp *types.Named, ftyp *types.Signature) error {
	p := g.printer
	newResultType := types.NewTuple(ftyp.Results().At(0), types.NewVar(1, nil, "e", etyp))
	newSigType := types.NewSignature(nil, ftyp.Params(), newResultType, ftyp.Variadic())

	p.P("")
	p.P("// %s is...", deriveFuncName)
	p.P("func %s(err error, f %s) %s {", deriveFuncName, ftyp.String(), newSigType.String())
	p.In()
	p.P("return %s {", newSigType.String())
	p.In()
	as := varnames(newSigType.Params())
	p.P("out, success := f(%s)", strings.Join(as, ", "))
	p.P("if success {")
	p.In()
	p.P("return out, nil")
	p.Out()
	p.P("}")
	p.P("return out, err")
	p.Out()
	p.P("}")
	p.Out()
	p.P("}")
	return nil
}
