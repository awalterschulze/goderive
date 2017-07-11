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
		stringsPkg: p.NewImport("strings"),
	}
}

type gen struct {
	derive.TypesMap
	printer    derive.Printer
	stringsPkg derive.Import
}

func (this *gen) Add(name string, typs []types.Type) (string, error) {
	if len(typs) == 0 {
		return "", fmt.Errorf("%s does not have at least one argument", name)
	}
	switch t := typs[0].(type) {
	case *types.Slice:
		switch t.Elem().(type) {
		case *types.Slice:
			_, err := this.sliceType(name, typs)
			if err != nil {
				return "", err
			}
			return this.SetFuncName(name, typs...)
		case *types.Basic:
			err := this.stringType(name, typs)
			if err != nil {
				return "", err
			}
			return this.SetFuncName(name, typs...)
		}
	case *types.Signature:
		_, err := this.errorType(name, typs)
		if err != nil {
			return "", err
		}
		return this.SetFuncName(name, typs...)
	case *types.Tuple:
		if t.Len() == 2 {
			ts := make([]types.Type, 2)
			ts[0] = t.At(0).Type()
			ts[1] = t.At(1).Type()
			_, err := this.errorType(name, ts)
			if err != nil {
				return "", err
			}
			return this.SetFuncName(name, ts...)
		}
	}
	return "", fmt.Errorf("unsupported type %s, not (a slice of slices) or (a slice of string)", typs[0])
}

func (this *gen) errorType(name string, typs []types.Type) ([]types.Type, error) {
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

func (this *gen) sliceType(name string, typs []types.Type) (types.Type, error) {
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

func (this *gen) stringType(name string, typs []types.Type) error {
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

func (this *gen) Generate(typs []types.Type) error {
	switch t := typs[0].(type) {
	case *types.Slice:
		switch t.Elem().(type) {
		case *types.Slice:
			return this.genSlice(typs)
		case *types.Basic:
			return this.genString(typs)
		}
	case *types.Signature:
		return this.genError(typs)
	case *types.Tuple:
		if t.Len() == 2 {
			ts := make([]types.Type, 2)
			ts[0] = t.At(0).Type()
			ts[1] = t.At(1).Type()
			return this.genError(ts)
		}
	}
	return fmt.Errorf("unsupported type %s, not (a slice of slices) or (a slice of string) or (a function and error)", typs[0])
}

func (this *gen) genError(typs []types.Type) error {
	p := this.printer
	this.Generating(typs...)
	name := this.GetFuncName(typs...)
	outTyps, err := this.errorType(name, typs)
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
			outs[i] = this.TypeString(outTyps[i])
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

func (this *gen) genSlice(typs []types.Type) error {
	p := this.printer
	this.Generating(typs...)
	name := this.GetFuncName(typs...)
	elemTyp, err := this.sliceType(name, typs)
	if err != nil {
		return err
	}
	typStr := this.TypeString(elemTyp)
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

func (this *gen) genString(typs []types.Type) error {
	p := this.printer
	this.Generating(typs...)
	name := this.GetFuncName(typs...)
	p.P("")
	p.P("func %s(list []string) string {", name)
	p.In()
	p.P("return %s.Join(list, \"\")", this.stringsPkg())
	p.Out()
	p.P("}")
	return nil
}
