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

// Package fmap contains the implementation of the fmap plugin, which generates the deriveFmap function.
//
// The deriveFmap function applies a given function to each element of a list, returning a list of results in the same order.
//   deriveFmap(func(A) B, []A) []B
//   deriveFmap(func(rune) B, string) []B
//
// deriveFmap can also be applied to a function that returns a value and an error.
//   deriveFmap(func(A) B, func() (A, error)) (B, error)
//   deriveFmap(func(A) (B, error), func() (A, error)) (func() (B, error), error)
//   deriveFmap(func(A), func() (A, error)) error
//   deriveFmap(func(A) (B, c, d, ...), func() (A, error)) (func() (B, c, d, ...), error)
// deriveFmap will propagate the error and not apply the first function to the result of the second function, if the second function returns an error.
//
// deriveFmap can also be applied to a channel.
//   deriveFmap(func(A) B, <-chan A) <-chan B
// deriveFmap will return the output channel immediately and start up a go routine in the background to process the incoming channel.
package fmap

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/ndeloof/goderive/derive"
)

// NewPlugin creates a new fmap plugin.
// This function returns the plugin name, default prefix and a constructor for the fmap code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("fmap", "deriveFmap", New)
}

// New is a constructor for the fmap code generator.
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
	case *types.Basic:
		_, err := g.stringOut(name, typs)
		if err != nil {
			return "", err
		}
		return g.SetFuncName(name, typs...)
	case *types.Signature:
		_, _, err := g.errorInOut(name, typs)
		if err != nil {
			return "", err
		}
		return g.SetFuncName(name, typs...)
	case *types.Chan:
		_, _, err := g.chanInOut(name, typs)
		if err != nil {
			return "", err
		}
		return g.SetFuncName(name, typs...)
	}
	return "", fmt.Errorf("unsupported type %s, not a slice or a string", typs[1])
}

func (g *gen) errorInOut(name string, typs []types.Type) (inTyp types.Type, outs *types.Tuple, err error) {
	esig, ok := typs[1].(*types.Signature)
	if !ok {
		return nil, nil, fmt.Errorf("%s, the second argument, %s, is not of type function", name, g.TypeString(typs[1]))
	}
	eparams := esig.Params()
	if eparams.Len() != 0 {
		return nil, nil, fmt.Errorf("%s, the second argument is a function, but wanted a function with zero arguments", name)
	}
	eres := esig.Results()
	if eres.Len() != 2 {
		return nil, nil, fmt.Errorf("%s, the second function argument does not have two results, but has %d resulting parameters", name, eres.Len())
	}
	if !derive.IsError(eres.At(1).Type()) {
		return nil, nil, fmt.Errorf("%s, the second argument is a function, but its second argument is not an error: %s", name, eres.At(1).Type())
	}
	elemTyp := eres.At(0).Type()
	sig, ok := typs[0].(*types.Signature)
	if !ok {
		return nil, nil, fmt.Errorf("%s, the first argument, %s, is not of type function", name, g.TypeString(typs[0]))
	}
	params := sig.Params()
	if params.Len() != 1 {
		return nil, nil, fmt.Errorf("%s, the first argument is a function, but wanted a function with one argument", name)
	}
	inTyp = params.At(0).Type()
	if !types.Identical(inTyp, elemTyp) {
		return nil, nil, fmt.Errorf("%s the function input type is not of type rune != %s",
			name, elemTyp)
	}
	res := sig.Results()
	return inTyp, res, nil
}

func (g *gen) stringOut(name string, typs []types.Type) (outTyp types.Type, err error) {
	typs[1] = types.Default(typs[1])
	basic, ok := typs[1].(*types.Basic)
	if !ok {
		return nil, fmt.Errorf("%s, the second argument, %s, is not of type basic", name, g.TypeString(typs[1]))
	}
	if basic.Kind() != types.String {
		return nil, fmt.Errorf("%s, the second argument, %v, is not a string", name, basic)
	}
	sig, ok := typs[0].(*types.Signature)
	if !ok {
		return nil, fmt.Errorf("%s, the first argument, %s, is not of type function", name, g.TypeString(typs[0]))
	}
	params := sig.Params()
	if params.Len() != 1 {
		return nil, fmt.Errorf("%s, the first argument is a function, but wanted a function with one argument", name)
	}
	elemTyp := types.Typ[types.Rune]
	inTyp := params.At(0).Type()
	if !types.Identical(inTyp, elemTyp) {
		return nil, fmt.Errorf("%s the function input type is not of type rune != %s",
			name, elemTyp)
	}
	res := sig.Results()
	if res.Len() != 1 {
		return nil, fmt.Errorf("%s, the function argument does not have a single result, but has %d resulting parameters", name, res.Len())
	}
	outTyp = res.At(0).Type()
	return outTyp, nil
}

func (g *gen) chanInOut(name string, typs []types.Type) (inTyp, outTyp types.Type, err error) {
	chanType, ok := typs[1].(*types.Chan)
	if !ok {
		return nil, nil, fmt.Errorf("%s, the second argument, %s, is not of type chan", name, g.TypeString(typs[1]))
	}
	sig, ok := typs[0].(*types.Signature)
	if !ok {
		return nil, nil, fmt.Errorf("%s, the first argument, %s, is not of type function", name, g.TypeString(typs[0]))
	}
	params := sig.Params()
	if params.Len() != 1 {
		return nil, nil, fmt.Errorf("%s, the first argument is a function, but wanted a function with one argument", name)
	}
	elemTyp := chanType.Elem()
	inTyp = params.At(0).Type()
	if !types.Identical(inTyp, elemTyp) {
		return nil, nil, fmt.Errorf("%s the function input type and chan element type are different %s != %s",
			name, inTyp, elemTyp)
	}
	res := sig.Results()
	if res.Len() != 1 {
		return nil, nil, fmt.Errorf("%s, the function argument does not have a single result, but has %d resulting parameters", name, res.Len())
	}
	outTyp = res.At(0).Type()
	return inTyp, outTyp, nil
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
	if res.Len() != 1 {
		return nil, nil, fmt.Errorf("%s, the function argument does not have a single result, but has %d resulting parameters", name, res.Len())
	}
	outTyp = res.At(0).Type()
	return inTyp, outTyp, nil
}

func (g *gen) Generate(typs []types.Type) error {
	switch typs[1].(type) {
	case *types.Slice:
		return g.genSlice(typs)
	case *types.Basic:
		return g.genString(typs)
	case *types.Signature:
		return g.genError(typs)
	case *types.Chan:
		return g.genChan(typs)
	}
	return fmt.Errorf("unsupported type %s, not a slice or a string", typs[1])
}

func (g *gen) genChan(typs []types.Type) error {
	name := g.GetFuncName(typs...)
	in, out, err := g.chanInOut(name, typs)
	if err != nil {
		return err
	}
	g.Generating(typs...)
	p := g.printer
	inStr := g.TypeString(in)
	outStr := g.TypeString(out)
	outerStr := outStr
	if strings.HasPrefix(outStr, "<-") {
		outerStr = "(" + outStr + ")"
	}
	p.P("")
	p.P("// %s returns an output channel where the items are the result of the input function being applied to the items on the input channel.", name)
	p.P("func %s(f func(%s) %s, in <-chan %s) <-chan %s {", name, inStr, outStr, inStr, outerStr)
	p.In()
	p.P("out := make(chan %s, cap(in))", outerStr)
	p.P("go func() {")
	p.In()
	p.P("for a := range in {")
	p.In()
	p.P("b := f(a)")
	p.P("out <- b")
	p.Out()
	p.P("}")
	p.P("close(out)")
	p.Out()
	p.P("}()")
	p.P("return out")
	p.Out()
	p.P("}")
	return nil
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
	p.P("// %s returns a list where each element of the input list has been morphed by the input function.", name)
	p.P("func %s(f func(%s) %s, list []%s) []%s {", name, inStr, outStr, inStr, outStr)
	p.In()
	p.P("out := make([]%s, len(list))", outStr)
	p.P("for i, elem := range list {")
	p.In()
	p.P("out[i] = f(elem)")
	p.Out()
	p.P("}")
	p.P("return out")
	p.Out()
	p.P("}")
	return nil
}

func (g *gen) genString(typs []types.Type) error {
	name := g.GetFuncName(typs...)
	out, err := g.stringOut(name, typs)
	if err != nil {
		return err
	}
	g.Generating(typs...)
	p := g.printer
	outStr := g.TypeString(out)
	p.P("")
	p.P("// %s morphs a string into list by apply the input function to each rune.", name)
	p.P("func %s(f func(rune) %s, ss string) []%s {", name, outStr, outStr)
	p.In()
	p.P("out := make([]%s, len([]rune(ss)))", outStr)
	p.P("for i, elem := range ss {")
	p.In()
	p.P("out[i] = f(elem)")
	p.Out()
	p.P("}")
	p.P("return out")
	p.Out()
	p.P("}")
	return nil
}

func (g *gen) genError(typs []types.Type) error {
	name := g.GetFuncName(typs...)
	in, out, err := g.errorInOut(name, typs)
	if err != nil {
		return err
	}
	g.Generating(typs...)
	p := g.printer
	inStr := g.TypeString(in)
	p.P("")
	switch out.Len() {
	case 0:
		p.P("// %s returns an error if g returns one, otherwise it applies f to g's result.", name)
		p.P("func %s(f func(%s), g func() (%s, error)) error {", name, inStr, inStr)
		p.In()
		p.P("v, err := g()")
		p.P("if err != nil {")
		p.In()
		p.P("return err")
		p.Out()
		p.P("}")
		p.P("f(v)")
		p.P("return nil")
		p.Out()
		p.P("}")
		return nil
	case 1:
		t := out.At(0).Type()
		outStr := g.TypeString(t)
		zeroStr := derive.Zero(t)
		p.P("// %s returns an error if g returns one, otherwise it applies f to g's result and returns it.", name)
		p.P("func %s(f func(%s) %s, g func() (%s, error)) (%s, error) {", name, inStr, outStr, inStr, outStr)
		p.In()
		p.P("v, err := g()")
		p.P("if err != nil {")
		p.In()
		p.P("return %s, err", zeroStr)
		p.Out()
		p.P("}")
		p.P("return f(v), nil")
		p.Out()
		p.P("}")
	default:
		outTyps := make([]types.Type, out.Len())
		outTypStrs := make([]string, out.Len())
		for i := range outTyps {
			outTyps[i] = out.At(i).Type()
			outTypStrs[i] = g.TypeString(outTyps[i])
		}
		outStr := strings.Join(outTypStrs, ", ")
		p.P("// %s returns an error if g returns one, otherwise it applies f to g's result and returns it.", name)
		p.P("func %s(f func(%s) (%s), g func() (%s, error)) (func() (%s), error) {", name, inStr, outStr, inStr, outStr)
		p.In()
		p.P("v, err := g()")
		p.P("if err != nil {")
		p.In()
		p.P("return nil, err")
		p.Out()
		p.P("}")
		p.P("return %s(f(v)), nil", g.tuple.GetFuncName(outTyps...))
		p.Out()
		p.P("}")
	}
	return nil
}
