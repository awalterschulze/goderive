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

// Package join contains the implementation of the join plugin, which generates the deriveJoin function.
//
// The deriveJoin function joins a slice of slices into a single slice.
//    deriveJoin([][]T) []T
//    deriveJoin([]string) string
//
// The deriveJoin function also joins two tuples, both with errors, into a single tuple with a single error.
//    deriveJoin(func() (T, error), error) func() (T, error)
//    deriveJoin(func() error, error) func() error
//    deriveJoin(func() (T, ..., error), error) func() (T, ..., error)
//
// The deriveJoin function can also join channels
//    deriveJoin(<-chan <-chan T) <-chan T
//    deriveJoin(chan <-chan T) <-chan T
//    deriveJoin([]<-chan T) <-chan T
//    deriveJoin([]chan T) <-chan T
// deriveJoin immediately return the output channel and start up a go routine to process the main incoming channel.
// It will then start up a go routine to listen on every new incoming channel and send those events to the outgoing channel.
//    deriveJoin(chan T, chan T, ...) <-chan T
// deriveJoin with a variable number of channels as parameter will do a select over those channels, until all are closed.
package join

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/awalterschulze/goderive/derive"
)

// NewPlugin creates a new join plugin.
// This function returns the plugin name, default prefix and a constructor for the join code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("join", "deriveJoin", New)
}

// New is a constructor for the join code generator.
// This generator should be reconstructed for each package.
func New(typesMap derive.TypesMap, p derive.Printer, deps map[string]derive.Dependency) derive.Generator {
	return &gen{
		TypesMap:   typesMap,
		printer:    p,
		stringsPkg: p.NewImport("strings", "strings"),
		syncPkg:    p.NewImport("sync", "sync"),
	}
}

type gen struct {
	derive.TypesMap
	printer    derive.Printer
	stringsPkg derive.Import
	syncPkg    derive.Import
}

func (g *gen) Add(name string, typs []types.Type) (string, error) {
	if len(typs) == 0 {
		return "", fmt.Errorf("%s does not have at least one argument", name)
	}
	switch t := typs[0].(type) {
	case *types.Slice:
		switch t.Elem().(type) {
		case *types.Slice:
			_, err := g.sliceType(name, typs)
			if err != nil {
				return "", err
			}
			return g.SetFuncName(name, typs...)
		case *types.Basic:
			err := g.stringType(name, typs)
			if err != nil {
				return "", err
			}
			return g.SetFuncName(name, typs...)
		case *types.Chan:
			_, _, err := g.sliceOfChanType(name, typs)
			if err != nil {
				return "", err
			}
			return g.SetFuncName(name, typs...)
		}
	case *types.Signature:
		_, err := g.errorType(name, typs)
		if err != nil {
			return "", err
		}
		return g.SetFuncName(name, typs...)
	case *types.Tuple:
		if t.Len() == 2 {
			ts := make([]types.Type, 2)
			ts[0] = t.At(0).Type()
			ts[1] = t.At(1).Type()
			_, err := g.errorType(name, ts)
			if err != nil {
				return "", err
			}
			return g.SetFuncName(name, ts...)
		}
	case *types.Chan:
		switch t.Elem().(type) {
		case *types.Chan:
			_, _, err := g.chanType(name, typs)
			if err != nil {
				return "", err
			}
			return g.SetFuncName(name, typs...)
		default:
			_, _, err := g.chanVariantTypes(name, typs)
			if err != nil {
				return "", err
			}
			return g.SetFuncName(name, typs...)
		}
	}
	return "", fmt.Errorf("unsupported type %s, not (a slice of slices) or (a slice of string)", typs[0])
}

func (g *gen) errorType(name string, typs []types.Type) ([]types.Type, error) {
	if len(typs) != 2 {
		return nil, fmt.Errorf("%s does not have two arguments", name)
	}
	sig, ok := typs[0].(*types.Signature)
	if !ok {
		return nil, fmt.Errorf("%s, the first argument, %s, is not of type function", name, typs[0])
	}
	if !derive.IsError(typs[1]) {
		return nil, fmt.Errorf("%s, the second argument, %s, is not of type error", name, typs[1])
	}
	if sig.Params().Len() != 0 {
		return nil, fmt.Errorf("%s, the first argument is a function, but it has parameters %v", name, sig.Params())
	}
	res := sig.Results()
	if res.Len() == 0 {
		return nil, fmt.Errorf("%s, the first argument is a function, but it has no results", name)
	}
	last := res.At(res.Len() - 1)
	if !derive.IsError(last.Type()) {
		return nil, fmt.Errorf("%s, the first argument is a function, but its last result is not an error: %v", name, last.Type())
	}
	outTyps := make([]types.Type, res.Len()-1)
	for i := range outTyps {
		outTyps[i] = res.At(i).Type()
	}
	return outTyps, nil
}

func (g *gen) chanType(name string, typs []types.Type) (types.Type, types.ChanDir, error) {
	if len(typs) != 1 {
		return nil, types.SendRecv, fmt.Errorf("%s does not have one argument", name)
	}
	chanTyp, ok := typs[0].(*types.Chan)
	if !ok {
		return nil, types.SendRecv, fmt.Errorf("%s, the argument, %s, is not of type chan", name, typs[0])
	}
	chanOfChanTyp, ok := chanTyp.Elem().(*types.Chan)
	if !ok {
		return nil, types.SendRecv, fmt.Errorf("%s, the argument, %s, is not of type chan of chan", name, typs[0])
	}
	elemType := chanOfChanTyp.Elem()
	return elemType, chanTyp.Dir(), nil
}

func (g *gen) chanVariantTypes(name string, typs []types.Type) ([]types.Type, []types.ChanDir, error) {
	if len(typs) < 2 {
		return nil, nil, fmt.Errorf("%s does not have at least two arguments", name)
	}
	chanTyps := make([]types.Type, len(typs))
	dirs := make([]types.ChanDir, len(typs))
	for i := range typs {
		chanTyp, ok := typs[i].(*types.Chan)
		if !ok {
			return nil, nil, fmt.Errorf("%s, the argument, %s, is not of type chan", name, typs[0])
		}
		chanTyps[i] = chanTyp.Elem()
		if i != 0 {
			if !types.Identical(chanTyps[i], chanTyps[i-1]) {
				return nil, nil, fmt.Errorf("%s, channel types are different %s != %s", name, typs[i-1], typs[i])
			}
		}
		dirs[i] = chanTyp.Dir()
	}
	return chanTyps, dirs, nil
}

func (g *gen) sliceOfChanType(name string, typs []types.Type) (types.Type, types.ChanDir, error) {
	if len(typs) != 1 {
		return nil, types.SendRecv, fmt.Errorf("%s does not have one argument", name)
	}
	sliceTyp, ok := typs[0].(*types.Slice)
	if !ok {
		return nil, types.SendRecv, fmt.Errorf("%s, the argument, %s, is not of type slice", name, typs[0])
	}
	sliceOfChanTyp, ok := sliceTyp.Elem().(*types.Chan)
	if !ok {
		return nil, types.SendRecv, fmt.Errorf("%s, the argument, %s, is not of type slice of chan", name, typs[0])
	}
	elemType := sliceOfChanTyp.Elem()
	return elemType, sliceOfChanTyp.Dir(), nil
}

func (g *gen) sliceType(name string, typs []types.Type) (types.Type, error) {
	if len(typs) != 1 {
		return nil, fmt.Errorf("%s does not have one argument", name)
	}
	sliceTyp, ok := typs[0].(*types.Slice)
	if !ok {
		return nil, fmt.Errorf("%s, the argument, %s, is not of type slice", name, typs[0])
	}
	sliceOfSliceTyp, ok := sliceTyp.Elem().(*types.Slice)
	if !ok {
		return nil, fmt.Errorf("%s, the argument, %s, is not of type slice of slice", name, typs[0])
	}
	elemType := sliceOfSliceTyp.Elem()
	return elemType, nil
}

func (g *gen) stringType(name string, typs []types.Type) error {
	if len(typs) != 1 {
		return fmt.Errorf("%s does not have one argument", name)
	}
	sliceTyp, ok := typs[0].(*types.Slice)
	if !ok {
		return fmt.Errorf("%s, the argument, %s, is not of type slice", name, typs[0])
	}
	basic, ok := sliceTyp.Elem().(*types.Basic)
	if !ok {
		return fmt.Errorf("%s, the argument, %s, is not of a slice of type basic", name, typs[0])
	}
	if basic.Kind() != types.String {
		return fmt.Errorf("%s, the argument, %s, is not of a slice of string", name, typs[0])
	}
	return nil
}

func (g *gen) Generate(typs []types.Type) error {
	switch t := typs[0].(type) {
	case *types.Slice:
		switch t.Elem().(type) {
		case *types.Slice:
			return g.genSlice(typs)
		case *types.Basic:
			return g.genString(typs)
		case *types.Chan:
			return g.genSliceOfChan(typs)
		}
	case *types.Signature:
		return g.genError(typs)
	case *types.Tuple:
		if t.Len() == 2 {
			ts := make([]types.Type, 2)
			ts[0] = t.At(0).Type()
			ts[1] = t.At(1).Type()
			return g.genError(ts)
		}
	case *types.Chan:
		switch t.Elem().(type) {
		case *types.Chan:
			return g.genChan(typs)
		default:
			return g.genChanVariant(typs)
		}
	}
	return fmt.Errorf("unsupported type %s, not (a slice of slices) or (a slice of string) or (a function and error)", typs[0])
}

func (g *gen) genError(typs []types.Type) error {
	p := g.printer
	g.Generating(typs...)
	name := g.GetFuncName(typs...)
	outTyps, err := g.errorType(name, typs)
	if err != nil {
		return err
	}
	p.P("")
	if len(outTyps) == 0 {
		p.P("func %s(f func() error, err error) error {", name)
		p.In()
		p.P("if err != nil {")
		p.In()
		p.P("return err")
		p.Out()
		p.P("}")
		p.P("return f()")
		p.Out()
		p.P("}")
	} else {
		outs := make([]string, len(outTyps))
		zeros := make([]string, len(outTyps))
		for i := range outTyps {
			outs[i] = g.TypeString(outTyps[i])
			zeros[i] = derive.Zero(outTyps[i])
		}
		outStr := strings.Join(outs, ", ")
		p.P("func %s(f func() (%s, error), err error) (%s, error) {", name, outStr, outStr)
		p.In()
		p.P("if err != nil {")
		p.In()
		p.P("return %s, err", strings.Join(zeros, ", "))
		p.Out()
		p.P("}")
		p.P("return f()")
		p.Out()
		p.P("}")
	}
	return nil
}

func (g *gen) genChan(typs []types.Type) error {
	p := g.printer
	g.Generating(typs...)
	name := g.GetFuncName(typs...)
	elemTyp, dir, err := g.chanType(name, typs)
	if err != nil {
		return err
	}
	dirStr := ""
	if dir == types.RecvOnly {
		dirStr = "<-"
	}
	typStr := g.TypeString(elemTyp)
	p.P("")
	p.P("func %s(in %schan (<-chan %s)) <-chan %s {", name, dirStr, typStr, typStr)
	p.In()
	p.P("out := make(chan %s)", typStr)
	p.P("go func() {")
	p.In()
	p.P("wait := %s.WaitGroup{}", g.syncPkg())
	p.P("for c := range in {")
	p.In()
	p.P("wait.Add(1)")
	p.P("res := c")
	p.P("go func() {")
	p.In()
	p.P("for r := range res {")
	p.In()
	p.P("out <- r")
	p.Out()
	p.P("}")
	p.P("wait.Done()")
	p.Out()
	p.P("}()")
	p.Out()
	p.P("}")
	p.P("wait.Wait()")
	p.P("close(out)")
	p.Out()
	p.P("}()")
	p.P("return out")
	p.Out()
	p.P("}")
	return nil
}

func (g *gen) genChanVariant(typs []types.Type) error {
	p := g.printer
	g.Generating(typs...)
	name := g.GetFuncName(typs...)
	elemTyps, dirs, err := g.chanVariantTypes(name, typs)
	if err != nil {
		return err
	}
	dirstrs := make([]string, len(typs))
	elemstrs := make([]string, len(typs))
	pairs := make([]string, len(typs))
	typStr := g.TypeString(elemTyps[0])
	csnil := make([]string, len(typs))
	for i := range typs {
		if dirs[i] == types.RecvOnly {
			dirstrs[i] = "<-"
		}
		elemstrs[i] = g.TypeString(elemTyps[i])
		pairs[i] = fmt.Sprintf("c%d %schan %s", i, dirstrs[i], elemstrs[i])
		csnil[i] = fmt.Sprintf("c%d != nil", i)
	}
	p.P("")
	p.P("func %s(%s) <-chan %s {", name, strings.Join(pairs, ", "), typStr)
	p.In()
	p.P("out := make(chan %s)", typStr)
	p.P("go func() {")
	p.In()
	p.P("for %s {", strings.Join(csnil, " || "))
	p.In()
	p.P("select {")
	for i := range typs {
		p.P("case v%d, ok%d := <-c%d:", i, i, i)
		p.In()
		p.P("if !ok%d {", i)
		p.In()
		p.P("c%d = nil", i)
		p.Out()
		p.P("} else {")
		p.In()
		p.P("out <- v%d", i)
		p.Out()
		p.P("}")
		p.Out()
	}
	p.P("}")
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

func (g *gen) genSliceOfChan(typs []types.Type) error {
	p := g.printer
	g.Generating(typs...)
	name := g.GetFuncName(typs...)
	elemTyp, dir, err := g.sliceOfChanType(name, typs)
	if err != nil {
		return err
	}
	typStr := g.TypeString(elemTyp)
	dirStr := ""
	if dir == types.RecvOnly {
		dirStr = "<-"
	}
	p.P("")
	p.P("func %s(in []%schan %s) <-chan %s {", name, dirStr, typStr, typStr)
	p.In()
	p.P("out := make(chan %s)", typStr)
	p.P("go func() {")
	p.In()
	p.P("wait := %s.WaitGroup{}", g.syncPkg())
	p.P("for _, c := range in {")
	p.In()
	p.P("wait.Add(1)")
	p.P("res := c")
	p.P("go func() {")
	p.In()
	p.P("for r := range res {")
	p.In()
	p.P("out <- r")
	p.Out()
	p.P("}")
	p.P("wait.Done()")
	p.Out()
	p.P("}()")
	p.Out()
	p.P("}")
	p.P("wait.Wait()")
	p.P("close(out)")
	p.Out()
	p.P("}()")
	p.P("return out")
	p.Out()
	p.P("}")
	return nil
}

func (g *gen) genSlice(typs []types.Type) error {
	p := g.printer
	g.Generating(typs...)
	name := g.GetFuncName(typs...)
	elemTyp, err := g.sliceType(name, typs)
	if err != nil {
		return err
	}
	typStr := g.TypeString(elemTyp)
	p.P("")
	p.P("func %s(list [][]%s) []%s {", name, typStr, typStr)
	p.In()
	p.P("if list == nil {")
	p.In()
	p.P("return nil")
	p.Out()
	p.P("}")
	p.P("l := 0")
	p.P("for _, elem := range list {")
	p.In()
	p.P("l += len(elem)")
	p.Out()
	p.P("}")
	p.P("res := make([]%s, 0, l)", typStr)
	p.P("for _, elem := range list {")
	p.In()
	p.P("res = append(res, elem...)")
	p.Out()
	p.P("}")
	p.P("return res")
	p.Out()
	p.P("}")
	return nil
}

func (g *gen) genString(typs []types.Type) error {
	p := g.printer
	g.Generating(typs...)
	name := g.GetFuncName(typs...)
	p.P("")
	p.P("func %s(list []string) string {", name)
	p.In()
	p.P("return %s.Join(list, \"\")", g.stringsPkg())
	p.Out()
	p.P("}")
	return nil
}
