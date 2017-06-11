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

package join

import (
	"fmt"
	"go/types"

	"github.com/awalterschulze/goderive/derive"
)

func NewPlugin() derive.Plugin {
	return derive.NewPlugin("join", "deriveJoin", New)
}

func New(typesMap derive.TypesMap, p derive.Printer, deps map[string]derive.Dependency) derive.Generator {
	return &join{
		TypesMap: typesMap,
		printer:  p,
	}
}

type join struct {
	derive.TypesMap
	printer  derive.Printer
	bytesPkg derive.Import
}

func (this *join) Add(name string, typs []types.Type) (string, error) {
	if len(typs) != 1 {
		return "", fmt.Errorf("%s does not have one argument", name)
	}
	sliceTyp, ok := typs[0].(*types.Slice)
	if !ok {
		return "", fmt.Errorf("%s, the argument, %s, is not of type slice", name, typs[0])
	}
	sliceOfSliceTyp, ok := sliceTyp.Elem().(*types.Slice)
	if !ok {
		return "", fmt.Errorf("%s, the argument, %s, is not of type slice of slice", name, typs[0])
	}
	elemType := sliceOfSliceTyp.Elem()
	return this.SetFuncName(name, elemType)
}

func (this *join) Generate() error {
	for _, typs := range this.ToGenerate() {
		if err := this.genFuncFor(typs[0]); err != nil {
			return err
		}
	}
	return nil
}

func (this *join) genFuncFor(typ types.Type) error {
	p := this.printer
	this.Generating(typ)
	typStr := this.TypeString(typ)
	p.P("")
	p.P("func %s(list [][]%s) []%s {", this.GetFuncName(typ), typStr, typStr)
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
