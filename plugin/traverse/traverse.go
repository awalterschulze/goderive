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

// Package traverse contains the implementation of the traverse plugin, which generates the deriveTraverse function.
//
// The deriveTraverse function applies a given function to each element of a list, returning a list of results in the same order or an error.
//   deriveTraverse(func(A) (B, error), []A) ([]B, error)
package traverse

import (
	"fmt"
	"go/types"

	"github.com/ndeloof/goderive/derive"
)

// NewPlugin creates a new traverse plugin.
// This function returns the plugin name, default prefix and a constructor for the traverse code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("traverse", "deriveTraverse", New)
}

// New is a constructor for the traverse code generator.
// This generator should be reconstructed for each package.
func New(typesMap derive.TypesMap, p derive.Printer, deps map[string]derive.Dependency) derive.Generator {
	return &gen{
		TypesMap: typesMap,
		printer:  p,
		tuple:    deps["tuple"],
	}
}

type gen struct {
	derive.TypesMap
	printer derive.Printer
	tuple   derive.Dependency
}

func (g *gen) Add(name string, typs []types.Type) (string, error) {
	if len(typs) != 2 {
		return "", fmt.Errorf("%s does not have two arguments", name)
	}
	switch typs[1].(type) {
	case *types.Slice:
		_, _, err := g.sliceInOut(name, typs)
		if err != nil {
			return "", err
		}
		return g.SetFuncName(name, typs...)
	}
	return "", fmt.Errorf("unsupported type %s, not a slice", typs[1])
}

func (g *gen) sliceInOut(name string, typs []types.Type) (inTyp types.Type, outTyp types.Type, err error) {
	sliceTyp, ok := typs[1].(*types.Slice)
	if !ok {
		return nil, nil, fmt.Errorf("%s, the second argument, %s, is not of type slice", name, g.TypeString(typs[1]))
	}
	sig, ok := typs[0].(*types.Signature)
	if !ok {
		return nil, nil, fmt.Errorf("%s, the first argument, %s, is not of type function", name, g.TypeString(typs[0]))
	}
	params := sig.Params()
	if params.Len() != 1 {
		return nil, nil, fmt.Errorf("%s, the first argument is a function, but wanted a function with one argument", name)
	}
	elemTyp := sliceTyp.Elem()
	inTyp = params.At(0).Type()
	if !types.Identical(inTyp, elemTyp) {
		return nil, nil, fmt.Errorf("%s the function input type and slice element type are different %s != %s",
			name, inTyp, elemTyp)
	}
	res := sig.Results()
	if res.Len() != 2 {
		return nil, nil, fmt.Errorf("%s, the function argument does not have a single result, but has %d resulting parameters", name, res.Len())
	}
	if !derive.IsError(res.At(1).Type()) {
		return nil, nil, fmt.Errorf("%s, the function's second result is not an error, but %s", name, g.TypeString(res.At(1).Type()))
	}
	outTyp = res.At(0).Type()
	return inTyp, outTyp, nil
}

func (g *gen) Generate(typs []types.Type) error {
	switch typs[1].(type) {
	case *types.Slice:
		return g.genSlice(typs)
	}
	return fmt.Errorf("unsupported type %s, not a slice or a string", typs[1])
}

func (g *gen) genSlice(typs []types.Type) error {
	name := g.GetFuncName(typs...)
	in, out, err := g.sliceInOut(name, typs)
	if err != nil {
		return err
	}
	g.Generating(typs...)
	p := g.printer
	inStr := g.TypeString(in)
	outStr := g.TypeString(out)
	p.P("")
	p.P("// %s returns a list where each element of the input list has been morphed by the input function or an error.", name)
	p.P("func %s(f func(%s) (%s, error), list []%s) ([]%s, error) {", name, inStr, outStr, inStr, outStr)
	p.In()
	p.P("out := make([]%s, len(list))", outStr)
	p.P("var err error")
	p.P("for i, elem := range list {")
	p.In()
	p.P("out[i], err = f(elem)")
	p.P("if err != nil {")
	p.In()
	p.P("return nil, err")
	p.Out()
	p.P("}")
	p.Out()
	p.P("}")
	p.P("return out, nil")
	p.Out()
	p.P("}")
	return nil
}
