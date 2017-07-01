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

// Package flip contains the implementation of the flip plugin, which generates the deriveFlip function.
// The deriveFlip function flips the first two parameters of the input function.
//   deriveFlip(f func(A, B, ...) T) func(B, A, ...) T
package flip

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/awalterschulze/goderive/derive"
)

// NewPlugin creates a new flip plugin.
// This function returns the plugin name, default prefix and a constructor for the flip code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("flip", "deriveFlip", New)
}

// New is a constructor for the flip code generator.
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

func (this *gen) Add(name string, typs []types.Type) (string, error) {
	if len(typs) != 1 {
		return "", fmt.Errorf("%s does not have one argument", name)
	}
	sig, ok := typs[0].(*types.Signature)
	if !ok {
		return "", fmt.Errorf("%s, the first argument, %s, is not of type function", name, this.TypeString(typs[0]))
	}
	params := sig.Params()
	if params.Len() < 2 {
		return "", fmt.Errorf("%s, the first argument is a function, but wanted a function with more than one argument", name)
	}
	return this.SetFuncName(name, sig)
}

func (this *gen) Generate(typs []types.Type) error {
	sig, ok := typs[0].(*types.Signature)
	if !ok {
		return fmt.Errorf("%s, the first argument, %s, is not of type function", this.GetFuncName(typs[0]), this.TypeString(typs[0]))
	}
	return this.genFuncFor(sig)
}

func (this *gen) genFuncFor(ftyp *types.Signature) error {
	p := this.printer
	this.Generating(ftyp)
	fStr := this.TypeString(ftyp)
	funcName := this.GetFuncName(ftyp)
	params := ftyp.Params()
	a0 := params.At(0)
	a1 := params.At(1)
	*a0, *a1 = *a1, *a0
	gStr := this.TypeString(ftyp)
	*a0, *a1 = *a1, *a0
	p.P("")
	p.P("func %s(f %s) %s {", funcName, fStr, gStr)
	p.In()
	p.P("return %s {", gStr)
	p.In()
	as := make([]string, params.Len())
	for i := 0; i < params.Len(); i++ {
		as[i] = params.At(i).Name()
	}
	p.P("return f(%s)", strings.Join(as, ", "))
	p.Out()
	p.P("}")
	p.Out()
	p.P("}")
	return nil
}
