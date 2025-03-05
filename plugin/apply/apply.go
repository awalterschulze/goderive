//  Copyright 2021 Jake Son
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

// Package apply contains the implementation of the apply plugin, which generates the deriveApply function.
//
// The deriveApply function applies the given argument to a given function and returns a function which requires filling in the other arguments.
//
//	deriveApply(f func(...A, B) C, B) func(...A) C
package apply

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/awalterschulze/goderive/derive"
)

// NewPlugin creates a new apply plugin.
// This function returns the plugin name, default prefix and a constructor for the apply code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("apply", "deriveApply", New)
}

// New is a constructor for the apply code generator.
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
	sig, lastArg, err := g.parseTypes(name, typs)
	if err != nil {
		return "", err
	}
	return g.SetFuncName(name, derive.RenameBlankIdentifier(sig), lastArg)
}

func (g *gen) parseTypes(name string, typs []types.Type) (sig *types.Signature, lastArg types.Type, err error) {
	if len(typs) != 2 {
		return nil, nil, fmt.Errorf("%s does not have two arguments", name)
	}
	sig, ok := typs[0].(*types.Signature)
	if !ok {
		return nil, nil, fmt.Errorf("%s, the first argument, %s, is not of type function", name, g.TypeString(typs[0]))
	}
	params := sig.Params()
	if params.Len() < 1 {
		return nil, nil, fmt.Errorf("%s, the first argument is a function, but wanted a function with at least one argument", name)
	}
	lastArg = params.At(params.Len() - 1).Type()
	if !types.AssignableTo(typs[1], lastArg) {
		return nil, nil, fmt.Errorf("%s, the second argument, %s, is not is assignable to the last argument of the function, type %s", name, g.TypeString(typs[1]), g.TypeString(lastArg))
	}
	return sig, lastArg, nil
}

func applySig(sig *types.Signature) (last *types.Var, returnFunc *types.Signature) {
	vs := vars(sig.Params())
	last = vs[len(vs)-1]
	second := types.NewTuple(vs[:len(vs)-1]...)
	returnFunc = types.NewSignature(nil, second, sig.Results(), sig.Variadic())
	return
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

func (g *gen) Generate(typs []types.Type) error {
	name := g.GetFuncName(typs...)
	ftyp, lastArg, err := g.parseTypes(name, typs)
	if err != nil {
		return err
	}
	g.Generating(typs...)
	p := g.printer
	fStr := g.TypeString(ftyp)
	last, gtyp := applySig(ftyp)
	gStr := g.TypeString(gtyp)
	p.P("")
	p.P("// %s applies the second argument to a given function's last argument and returns a function which which takes the rest of the parameters as input and finally returns the original input function's results.", name)
	p.P("func %s(f %s, %s %s) %s {", name, fStr, last.Name(), g.TypeString(lastArg), gStr)
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
