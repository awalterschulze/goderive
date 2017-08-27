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

// Package dup contains the implementation of the dup plugin, which generates the deriveDup function.
//
// The deriveDup duplicates messages received on c to both c1 and c2.
//   deriveDup(c <-chan T) (c1, c2 <-chan T)
package dup

import (
	"fmt"
	"go/types"

	"github.com/awalterschulze/goderive/derive"
)

// NewPlugin creates a new dup plugin.
// This function returns the plugin name, default prefix and a constructor for the dup code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("dup", "deriveDup", New)
}

// New is a constructor for the dup code generator.
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
	if len(typs) != 1 {
		return "", fmt.Errorf("%s expected one argument", name)
	}
	_, _, err := g.chanOut(name, typs[0])
	if err != nil {
		return "", err
	}
	return g.SetFuncName(name, typs...)
}

func (g *gen) chanOut(name string, typ types.Type) (types.Type, types.ChanDir, error) {
	chanType, ok := typ.(*types.Chan)
	if !ok {
		return nil, types.SendRecv, fmt.Errorf("%s is not a channel: %s", name, typ)
	}
	return chanType.Elem(), chanType.Dir(), nil
}

func (g *gen) Generate(typs []types.Type) error {
	name := g.GetFuncName(typs...)
	elemTyp, dir, err := g.chanOut(name, typs[0])
	if err != nil {
		return err
	}
	g.Generating(typs...)
	p := g.printer
	dirstr := ""
	if dir == types.RecvOnly {
		dirstr = "<-"
	}
	typstr := g.TypeString(elemTyp)
	p.P("")
	p.P("// %s duplicates messages received on c to both c1 and c2.", name)
	p.P("func %s(c %schan %s) (c1, c2 <-chan %s) {", name, dirstr, typstr, typstr)
	p.In()
	p.P("cc1, cc2 := make(chan %s, cap(c)), make(chan %s, cap(c))", typstr, typstr)
	p.P("go func() {")
	p.In()
	p.P("for v := range c {")
	p.In()
	p.P("cc1 <- v")
	p.P("cc2 <- v")
	p.Out()
	p.P("}")
	p.P("close(cc1)")
	p.P("close(cc2)")
	p.Out()
	p.P("}()")
	p.P("return cc1, cc2")
	p.Out()
	p.P("}")
	return nil
}
