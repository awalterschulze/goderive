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

// Package pipeline contains the implementation of the pipeline plugin, which generates the derivePipeline function.
//
// The derivePipeline starts up a concurrent pipeline of the given functions.
//   derivePipeline(func(A) <-chan B, func(B) <-chan C) func(A) <-chan C
//
// Example output can be found here:
// https://github.com/awalterschulze/goderive/tree/master/example/plugin/pipeline
package pipeline

import (
	"fmt"
	"go/types"

	"github.com/awalterschulze/goderive/derive"
)

// NewPlugin creates a new pipeline plugin.
// This function returns the plugin name, default prefix and a constructor for the pipeline code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("pipeline", "derivePipeline", New)
}

// New is a constructor for the pipeline code generator.
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
	if len(typs) != 2 {
		return "", fmt.Errorf("%s expected two arguments", name)
	}
	_, b1, err := g.funcInChanOut(name, typs[0])
	if err != nil {
		return "", err
	}
	b2, _, err := g.funcInChanOut(name, typs[1])
	if err != nil {
		return "", err
	}
	if !types.Identical(b1, b2) {
		return "", fmt.Errorf("%s function one's output %s is not the same as function two's input %s", name, b1, b2)
	}
	return g.SetFuncName(name, typs...)
}

func (g *gen) funcInChanOut(name string, typ types.Type) (inTyp, outTyp types.Type, err error) {
	sig, ok := typ.(*types.Signature)
	if !ok {
		return nil, nil, fmt.Errorf("%s is not a function: %s", name, typ)
	}
	params := sig.Params()
	results := sig.Results()
	if params.Len() != 1 {
		return nil, nil, fmt.Errorf("%s, the function has more than one parameter", name)
	}
	if results.Len() != 1 {
		return nil, nil, fmt.Errorf("%s, the function has more than one result", name)
	}
	resType := results.At(0).Type()
	chanType, ok := resType.(*types.Chan)
	if !ok {
		return nil, nil, fmt.Errorf("%s, the result, %s, is not of type chan", name, g.TypeString(resType))
	}
	return params.At(0).Type(), chanType.Elem(), nil
}

func (g *gen) Generate(typs []types.Type) error {
	name := g.GetFuncName(typs...)
	a, b1, err := g.funcInChanOut(name, typs[0])
	if err != nil {
		return err
	}
	_, c, err := g.funcInChanOut(name, typs[1])
	if err != nil {
		return err
	}
	g.Generating(typs...)
	p := g.printer
	cc := types.NewChan(types.RecvOnly, types.NewChan(types.RecvOnly, c))
	t0str := g.TypeString(typs[0])
	t1str := g.TypeString(typs[1])
	astr := g.TypeString(a)
	cstr := g.TypeString(c)
	ccstr := g.join.GetFuncName(cc)

	fmapFunc := g.fmap.GetFuncName(typs[1], types.NewChan(types.RecvOnly, b1))
	p.P("")
	p.P("// %s composes f and g into a concurrent pipeline.", name)
	p.P("func %s(f %s, g %s) func(%s) <-chan %s {", name, t0str, t1str, astr, cstr)
	p.In()
	p.P("return func(a %s) <-chan %s {", astr, cstr)
	p.In()
	p.P("b := f(a)")
	p.P("return %s(%s(g, b))", ccstr, fmapFunc)
	p.Out()
	p.P("}")
	p.Out()
	p.P("}")
	return nil
}
