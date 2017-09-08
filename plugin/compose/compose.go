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

// Package compose contains the implementation of the compose plugin, which generates the deriveCompose function.
//
// The deriveCompose function composes multiple functions that return an error into one function.
//    deriveCompose(func() (A, error), func(A) (B, error)) func() (B, error)
//    deriveCompose(func(A) (B, error), func(B) (C, error)) func(A) (C, error)
//    deriveCompose(func(A...) (B..., error), func(B...) (C..., error)) func(A...) (C..., error)
//    deriveCompose(func(A...) (B..., error), ..., func(C...) (D..., error)) func(A...) (D..., error)
//
// Example output can be found here:
// https://github.com/awalterschulze/goderive/tree/master/example/plugin/compose
package compose

import (
	"fmt"
	"go/types"
	"strconv"
	"strings"

	"github.com/awalterschulze/goderive/derive"
)

// NewPlugin creates a new compose plugin.
// This function returns the plugin name, default prefix and a constructor for the compose code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("compose", "deriveCompose", New)
}

// New is a constructor for the compose code generator.
// This generator should be reconstructed for each package.
func New(typesMap derive.TypesMap, p derive.Printer, deps map[string]derive.Dependency) derive.Generator {
	return &gen{
		TypesMap: typesMap,
		printer:  p,
		fmap:     deps["fmap"],
		join:     deps["join"],
	}
}

type gen struct {
	derive.TypesMap
	printer derive.Printer
	fmap    derive.Dependency
	join    derive.Dependency
}

func (g *gen) Add(name string, typs []types.Type) (string, error) {
	if len(typs) <= 1 {
		return "", fmt.Errorf("%s does not have more than one argument", name)
	}
	switch typs[0].(type) {
	case *types.Signature:
		_, _, err := g.errorType(name, typs)
		if err != nil {
			return "", err
		}
		return g.SetFuncName(name, typs...)
	}
	return "", fmt.Errorf("unsupported type %s", typs[0])
}

func (g *gen) errorType(name string, typs []types.Type) ([][]types.Type, [][]types.Type, error) {
	if len(typs) <= 1 {
		return nil, nil, fmt.Errorf("%s does not have at least two arguments", name)
	}
	params := make([][]types.Type, len(typs))
	results := make([][]types.Type, len(typs))
	for i, typ := range typs {
		sig, ok := typs[i].(*types.Signature)
		if !ok {
			return nil, nil, fmt.Errorf("%s, argument number %d, %s, is not of type function", name, i, typ)
		}
		params[i] = make([]types.Type, sig.Params().Len())
		for j := range params[i] {
			params[i][j] = sig.Params().At(j).Type()
		}
		if sig.Results().Len() == 0 {
			return nil, nil, fmt.Errorf("%s, function number %d, %s, does not return any parameters", name, i, typ)
		}
		errType := sig.Results().At(sig.Results().Len() - 1).Type()
		if !derive.IsError(errType) {
			return nil, nil, fmt.Errorf("%s, function number %d's last result, %s, is not of type error", name, i, errType)
		}
		results[i] = make([]types.Type, sig.Results().Len()-1)
		for j := range results[i] {
			results[i][j] = sig.Results().At(j).Type()
		}
	}

	for i := range params {
		if i == 0 {
			continue
		}
		if len(results[i-1]) != len(params[i]) {
			return nil, nil, fmt.Errorf("%s, function number %d and number %d has a different number of results and parameters respectively", name, i-1, i)
		}
		for j := range params[i] {
			if !types.AssignableTo(results[i-1][j], params[i][j]) {
				return nil, nil, fmt.Errorf("%s, function number %d and function number %d's results and parameter number %d's type is not assignable", name, i-1, i, j)
			}
		}
	}
	return params, results, nil
}

func (g *gen) Generate(typs []types.Type) error {
	switch typs[0].(type) {
	case *types.Signature:
		return g.genError(typs)
	}
	return fmt.Errorf("unsupported type %s, not (a slice of slices) or (a slice of string) or (a function and error)", typs[0])
}

func (g *gen) typeStrings(typs []types.Type) []string {
	ss := make([]string, len(typs))
	for i := range typs {
		ss[i] = g.TypeString(typs[i])
	}
	return ss
}

func wrap(s string) string {
	if strings.Contains(s, ",") {
		return "(" + s + ")"
	}
	return s
}

func vars(prefix string, num int) []string {
	ss := make([]string, num)
	for i := range ss {
		ss[i] = prefix + strconv.Itoa(i)
	}
	return ss
}

func zip(ss, rr []string) []string {
	qq := make([]string, len(ss))
	for i := range ss {
		qq[i] = ss[i] + " " + rr[i]
	}
	return qq
}

func (g *gen) genError(typs []types.Type) error {
	p := g.printer
	g.Generating(typs...)
	name := g.GetFuncName(typs...)
	params, results, err := g.errorType(name, typs)
	if err != nil {
		return err
	}

	paramStrs := make([][]string, len(params))
	fs := make([]string, len(params))
	resultStrs := make([][]string, len(results))
	fVarType := make([]string, len(params))
	vars := make([][]string, len(params)+1)
	for i := range params {
		paramStrs[i] = g.typeStrings(params[i])
		fs[i] = "f" + strconv.Itoa(i)
		resultStrs[i] = append(g.typeStrings(results[i]), "error")
		fVarType[i] = fmt.Sprintf("%s func(%s) %s", fs[i], strings.Join(paramStrs[i], ", "), wrap(strings.Join(resultStrs[i], ", ")))
		vars[i] = make([]string, len(params[i]))
		for j := range vars[i] {
			vars[i][j] = "v_" + strconv.Itoa(i) + "_" + strconv.Itoa(j)
		}
	}
	firstVarTypes := make([]string, len(params[0]))
	for i := range params[0] {
		firstVarTypes[i] = vars[0][i] + " " + g.TypeString(params[0][i])
	}
	zeros := make([]string, len(results[len(results)-1]))
	vars[len(vars)-1] = make([]string, len(results[len(results)-1]))
	for i, r := range results[len(results)-1] {
		zeros[i] = derive.Zero(r)
		vars[len(vars)-1][i] = "v_" + strconv.Itoa(len(vars)-1) + "_" + strconv.Itoa(i)
	}
	resFuncType := fmt.Sprintf("func(%s) %s", strings.Join(paramStrs[0], ", "), wrap(strings.Join(resultStrs[len(resultStrs)-1], ", ")))
	p.P("")
	p.P("func %s(%s) %s {", name, strings.Join(fVarType, ", "), resFuncType)
	p.In()
	p.P("return func(%s) %s {", strings.Join(firstVarTypes, ", "), wrap(strings.Join(resultStrs[len(resultStrs)-1], ", ")))
	p.In()
	for i := range params {
		p.P("%s, err%d := %s(%s)", strings.Join(vars[i+1], ", "), i, fs[i], strings.Join(vars[i], ", "))
		p.P("if err%d != nil {", i)
		p.In()
		p.P("return %s, err%d", strings.Join(zeros, ", "), i)
		p.Out()
		p.P("}")
	}
	p.P("return %s, nil", strings.Join(vars[len(vars)-1], ", "))
	p.Out()
	p.P("}")
	p.Out()
	p.P("}")

	return nil
}
