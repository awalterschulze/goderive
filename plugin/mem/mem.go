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

// Package mem contains the implementation of the mem plugin, which generates the deriveMem function.
//
// The deriveMem function returns a memoized version of the input function.
//   func deriveMem(func(A) B) func(A) B
//
package mem

import (
	"fmt"
	"go/token"
	"go/types"
	"strconv"
	"strings"

	"github.com/ndeloof/goderive/derive"
)

// NewPlugin creates a new mem plugin.
// This function returns the plugin name, default prefix and a constructor for the mem code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("mem", "deriveMem", New)
}

// New is a constructor for the mem code generator.
// This generator should be reconstructed for each package.
func New(typesMap derive.TypesMap, p derive.Printer, deps map[string]derive.Dependency) derive.Generator {
	return &gen{
		TypesMap: typesMap,
		printer:  p,
		equal:    deps["equal"],
		hash:     deps["hash"],
	}
}

type gen struct {
	derive.TypesMap
	printer derive.Printer
	equal   derive.Dependency
	hash    derive.Dependency
}

func (g *gen) Add(name string, typs []types.Type) (string, error) {
	if len(typs) != 1 {
		return "", fmt.Errorf("%s does not have one argument", name)
	}
	if _, ok := typs[0].(*types.Signature); !ok {
		return "", fmt.Errorf("%s, the argument, %s, is not of type func", name, typs[1])
	}
	return g.SetFuncName(name, typs[0])
}

func (g *gen) Generate(typs []types.Type) error {
	typ := typs[0]
	sigType, ok := typ.(*types.Signature)
	if !ok {
		return fmt.Errorf("%s, the argument, %s, is not of type func", g.GetFuncName(typ), typ)
	}
	return g.genFunc(sigType)
}

func (g *gen) typeStrings(typs []types.Type) []string {
	ss := make([]string, len(typs))
	for i := range typs {
		ss[i] = g.TypeString(typs[i])
	}
	return ss
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

func (g *gen) genFunc(typ *types.Signature) error {
	p := g.printer
	g.Generating(typ)
	name := g.GetFuncName(typ)
	typeStr := g.TypeString(typ)

	paramTypes := make([]types.Type, typ.Params().Len())
	paramFields := make([]*types.Var, typ.Params().Len())
	for i := 0; i < typ.Params().Len(); i++ {
		paramTypes[i] = typ.Params().At(i).Type()
		paramFields[i] = types.NewField(token.NoPos, nil, "Param"+strconv.Itoa(i), paramTypes[i], false)
	}
	paramTypeStrs := g.typeStrings(paramTypes)
	paramVars := vars("param", typ.Params().Len())
	params := zip(paramVars, paramTypeStrs)

	paramStruct := types.NewStruct(paramFields, nil)

	resTypes := make([]types.Type, typ.Results().Len())
	resFields := make([]*types.Var, typ.Results().Len())
	for i := 0; i < typ.Results().Len(); i++ {
		resTypes[i] = typ.Results().At(i).Type()
		resFields[i] = types.NewField(token.NoPos, nil, "Res"+strconv.Itoa(i), resTypes[i], false)
	}
	resTypeStrs := g.typeStrings(resTypes)
	resVars := vars("res", typ.Results().Len())
	res := zip(resVars, resTypeStrs)
	resStr := strings.Join(resTypeStrs, ", ")
	if len(resTypeStrs) > 1 {
		resStr = "(" + resStr + ")"
	}

	p.P("")

	p.P("// %s returns a memoized version of the input function.", name)
	p.P("func %s(f %s) %s {", name, typeStr, typeStr)
	p.In()

	if len(paramTypes) == 0 {
		p.P("memoized := false")
		for _, r := range res {
			p.P("var %s", r)
		}
		p.P("return func(%s) %s {", strings.Join(params, ", "), resStr)
		p.In()
		p.P("if !memoized {")
		p.In()
		if len(resTypes) == 0 {
			p.P("f(%s)", strings.Join(paramVars, ", "))
		} else {
			p.P("%s = f(%s)", strings.Join(resVars, ", "), strings.Join(paramVars, ", "))
		}
		p.P("memoized = true")
		p.Out()
		p.P("}")
		if len(resTypes) == 0 {
			p.P("return")
		} else {
			p.P("return %s", strings.Join(resVars, ", "))
		}
		p.Out()
		p.P("}")
	} else if len(paramTypes) == 1 && derive.IsComparable(paramTypes[0]) {
		if len(resTypes) == 0 {
			p.P("m := make(map[%s]struct{})", paramTypeStrs[0])
			p.P("return func(%s) %s {", strings.Join(params, ", "), resStr)
			p.In()
			p.P("if _, ok := m[%s]; ok {", paramVars[0])
			p.In()
			p.P("return")
			p.Out()
			p.P("}")
			p.P("f(%s)", paramVars[0])
			p.P("m[%s] = struct{}{}", paramVars[0])
			p.P("return")
			p.Out()
			p.P("}")
		} else if len(resTypes) == 1 {
			p.P("m := make(map[%s]%s)", paramTypeStrs[0], resTypeStrs[0])
			p.P("return func(%s) %s {", strings.Join(params, ", "), resStr)
			p.In()
			p.P("if v, ok := m[%s]; ok {", paramVars[0])
			p.In()
			p.P("return v")
			p.Out()
			p.P("}")
			p.P("v := f(%s)", paramVars[0])
			p.P("m[%s] = v", paramVars[0])
			p.P("return v")
			p.Out()
			p.P("}")
		} else {
			p.P("type output struct {")
			p.In()
			outFields, err := g.FieldStrings(resFields)
			if err != nil {
				return err
			}
			for _, f := range outFields {
				p.P("%s", f)
			}
			p.Out()
			p.P("}")
			p.P("m := make(map[%s]output)", paramTypeStrs[0])
			p.P("return func(%s) %s {", strings.Join(params, ", "), resStr)
			p.In()
			p.P("if o, ok := m[%s]; ok {", paramVars[0])
			p.In()
			p.P("return %s", strings.Join(vars("o.Res", typ.Results().Len()), ", "))
			p.Out()
			p.P("}")
			p.P("%s := f(%s)", strings.Join(resVars, ", "), paramVars[0])
			p.P("m[%s] = output{%s}", paramVars[0], strings.Join(resVars, ", "))
			p.P("return %s", strings.Join(resVars, ", "))
			p.Out()
			p.P("}")
		}
	} else if derive.IsComparable(paramStruct) {
		p.P("type input struct {")
		p.In()
		inFields, err := g.FieldStrings(paramFields)
		if err != nil {
			return err
		}
		for _, f := range inFields {
			p.P("%s", f)
		}
		p.Out()
		p.P("}")
		if len(resTypes) == 0 {
			p.P("m := make(map[input]struct{})")
			p.P("return func(%s) %s {", strings.Join(params, ", "), resStr)
			p.In()
			p.P("in := input{%s}", strings.Join(paramVars, ", "))
			p.P("if _, ok := m[in]; ok {")
			p.In()
			p.P("return")
			p.Out()
			p.P("}")
			p.P("f(%s)", strings.Join(paramVars, ", "))
			p.P("m[in] = struct{}{}")
			p.P("return")
			p.Out()
			p.P("}")
		} else if len(resTypes) == 1 {
			p.P("m := make(map[input]%s)", resTypeStrs[0])
			p.P("return func(%s) %s {", strings.Join(params, ", "), resStr)
			p.In()
			p.P("in := input{%s}", strings.Join(paramVars, ", "))
			p.P("if v, ok := m[in]; ok {")
			p.In()
			p.P("return v")
			p.Out()
			p.P("}")
			p.P("v := f(%s)", strings.Join(paramVars, ", "))
			p.P("m[in] = v")
			p.P("return v")
			p.Out()
			p.P("}")
		} else {
			p.P("type output struct {")
			p.In()
			outFields, err := g.FieldStrings(resFields)
			if err != nil {
				return err
			}
			for _, f := range outFields {
				p.P("%s", f)
			}
			p.Out()
			p.P("}")
			p.P("m := make(map[input]output)")
			p.P("return func(%s) %s {", strings.Join(params, ", "), resStr)
			p.In()
			p.P("in := input{%s}", strings.Join(paramVars, ", "))
			p.P("if o, ok := m[in]; ok {")
			p.In()
			p.P("return %s", strings.Join(vars("o.Res", typ.Results().Len()), ", "))
			p.Out()
			p.P("}")
			p.P("%s := f(%s)", strings.Join(resVars, ", "), strings.Join(paramVars, ", "))
			p.P("m[in] = output{%s}", strings.Join(resVars, ", "))
			p.P("return %s", strings.Join(resVars, ", "))
			p.Out()
			p.P("}")
		}
	} else {
		if len(paramTypes) >= 2 {
			p.P("type input struct {")
			p.In()
			inFields, err := g.FieldStrings(paramFields)
			if err != nil {
				return err
			}
			for _, f := range inFields {
				p.P("%s", f)
			}
			p.Out()
			p.P("}")
		}

		if len(resTypes) >= 2 {
			p.P("type output struct {")
			p.In()
			outFields, err := g.FieldStrings(resFields)
			if err != nil {
				return err
			}
			for _, f := range outFields {
				p.P("%s", f)
			}
			p.Out()
			p.P("}")
		}

		p.P("type mem struct {")
		p.In()
		if len(paramTypes) == 1 {
			p.P("in  %s", g.TypeString(paramTypes[0]))
		} else if len(paramTypes) >= 2 {
			p.P("in  input")
		}
		if len(resTypes) == 1 {
			p.P("out %s", g.TypeString(resTypes[0]))
		} else if len(resTypes) >= 2 {
			p.P("out output")
		}
		p.Out()
		p.P("}")

		p.P("m := make(map[uint64][]mem)")
		p.P("return func(%s) %s {", strings.Join(params, ", "), resStr)
		p.In()

		if len(paramTypes) == 1 {
			p.P("h := %s(%s)", g.hash.GetFuncName(paramTypes[0]), paramVars[0])
		} else {
			p.P("in := input{%s}", strings.Join(paramVars, ", "))
			p.P("h := %s(in)", g.hash.GetFuncName(paramStruct))
		}
		p.P("vs, ok := m[h]")
		p.P("if ok {")
		p.In()
		p.P("for _, v := range vs {")
		p.In()
		if len(paramTypes) == 1 {
			p.P("if %s(v.in, %s) {", g.equal.GetFuncName(paramTypes[0], paramTypes[0]), paramVars[0])
		} else {
			p.P("if %s(v.in, in) {", g.equal.GetFuncName(paramStruct, paramStruct))
		}
		p.In()
		if len(resTypes) == 1 {
			p.P("return v.out")
		} else {
			p.P("return %s", strings.Join(vars("v.out.Res", typ.Results().Len()), ", "))
		}
		p.Out()
		p.P("}")
		p.Out()
		p.P("}")
		p.Out()
		p.P("}")
		p.P("%s := f(%s)", strings.Join(resVars, ", "), strings.Join(paramVars, ", "))
		if len(resTypes) == 1 {
			if len(paramTypes) == 1 {
				p.P("m[h] = append(m[h], mem{%s, %s})", paramVars[0], resVars[0])
			} else {
				p.P("m[h] = append(m[h], mem{in, %s})", resVars[0])
			}
		} else {
			if len(paramTypes) == 1 {
				p.P("m[h] = append(m[h], mem{%s, output{%s}})", paramVars[0], strings.Join(resVars, ", "))
			} else {
				p.P("m[h] = append(m[h], mem{in, output{%s}})", strings.Join(resVars, ", "))
			}
		}
		p.P("return %s", strings.Join(resVars, ", "))
		p.Out()
		p.P("}")
	}
	p.Out()
	p.P("}")
	return nil

}
