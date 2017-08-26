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

// Package do contains the implementation of the do plugin, which generates the deriveDo function.
//
// The deriveDo function executes a list of functions concurrently and returns their results.
//   deriveDo(func() (A, error), func (B, error)) (A, B, error)
// Each function is executed in a go routine and the first error is returned.
// It waits for all functions to complete.
//
// The concept is stolen from applicative do in haskell or rather haxl.
// http://simonmar.github.io/bib/papers/applicativedo.pdf
// The applicative do rewrites the monadic do notation:
//   do {
//       a <- f
//       b <- g
//       return (f, g)
//   }
// To:
//   (,) <$> f <*> g
// When it detects that the functions do not depend on one another.
// Haskell type signatures that will hopefully help to explain.
// Fmap:
//   <$> :: (a -> b) -> m a -> m b
// Ap:
//   <*> :: m (a -> b) -> m a -> m b
//
// In go this could be:
//   func newTuple(a A, b B) func() (A, B) {
//       return func() (A, B) {
//           return a, b
//       }
//   }
//   deriveAp(
//       deriveFmap(
//           newTuple,
//           f,
//       )
//       g,
//   )
//   func deriveFmap(
//       newTuple func(A, B) func() (A, B),
//       f func() (A, error),
//       ) func(B) (func() (A, B), error) {
//           return func(b B) (func() (A, B), error) {
//               a, err := f()
//               if err != nil {
//                   return nil, err
//               }
//               return newTuple(a, b), nil
//           }
//   }
//   func deriveAp(fmapped func(B) (func() (A, B), error), g func() (B, error)) (func() (A, B), error) {
//       b, err := g()
//       if err != nil {
//           return nil, err
//       }
//       return fmapped(b)
//   }
// derviveDo builds on this, but requires the programmer to explicitly call deriveDo
//
// Example output can be found here:
// https://github.com/awalterschulze/goderive/tree/master/example/plugin/do
package do

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/awalterschulze/goderive/derive"
)

// NewPlugin creates a new do plugin.
// This function returns the plugin name, default prefix and a constructor for the do code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("do", "deriveDo", New)
}

// New is a constructor for the do code generator.
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
	if len(typs) < 2 {
		return "", fmt.Errorf("%s expected at least two arguments", name)
	}
	for i, typ := range typs {
		sig, ok := typ.(*types.Signature)
		if !ok {
			return "", fmt.Errorf("%s's argument number %d is not a function, but %s", name, i, typ)
		}
		if _, err := g.errorOut(name, sig); err != nil {
			return "", err
		}
	}
	return g.SetFuncName(name, typs...)
}

func (g *gen) errorOut(name string, sig *types.Signature) (typ types.Type, err error) {
	params := sig.Params()
	if params.Len() != 0 {
		return nil, fmt.Errorf("%s, the function argument does not take zero parameters", g.TypeString(sig))
	}
	res := sig.Results()
	if res.Len() != 2 {
		return nil, fmt.Errorf("%s, the function argument does not have two results, but has %d resulting parameters", name, res.Len())
	}
	if !derive.IsError(res.At(1).Type()) {
		return nil, fmt.Errorf("%s, the function's second return parameter is not an error: %s", name, res.At(1).Type())
	}
	elemTyp := res.At(0).Type()
	return elemTyp, nil
}

func (g *gen) Generate(typs []types.Type) error {
	name := g.GetFuncName(typs...)
	outs := make([]types.Type, len(typs))
	outstrs := make([]string, len(typs))
	funcstrs := make([]string, len(typs))
	vars := make([]string, len(typs))
	for i, typ := range typs {
		out, err := g.errorOut(name, typ.(*types.Signature))
		if err != nil {
			return err
		}
		outs[i] = out
		outstrs[i] = g.TypeString(out)
		funcstrs[i] = fmt.Sprintf("f%d func() (%s, error)", i, outstrs[i])
		vars[i] = fmt.Sprintf("v%d", i)
	}
	outstrs = append(outstrs, "error")
	g.Generating(typs...)
	p := g.printer
	p.P("")
	p.P("func %s(%s) (%s) {", name, strings.Join(funcstrs, ", "), strings.Join(outstrs, ", "))
	p.In()
	p.P("errChan := make(chan error)")
	for i := range typs {
		p.P("var %s %s", vars[i], g.TypeString(outs[i]))
		p.P("go func() {")
		p.In()
		p.P("var %serr error", vars[i])
		p.P("%s, %serr = f%d()", vars[i], vars[i], i)
		p.P("errChan <- %serr", vars[i])
		p.Out()
		p.P("}()")
	}
	p.P("var err error")
	p.P("for i := 0; i < %d; i++ {", len(typs))
	p.In()
	p.P("errc := <-errChan")
	p.P("if errc != nil {")
	p.In()
	p.P("if err == nil {")
	p.In()
	p.P("err = errc")
	p.Out()
	p.P("}")
	p.Out()
	p.P("}")
	p.Out()
	p.P("}")
	p.P("return %s, err", strings.Join(vars, ", "))
	p.Out()
	p.P("}")
	return nil
}
