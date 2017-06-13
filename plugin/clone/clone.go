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

package clone

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/awalterschulze/goderive/derive"
)

func NewPlugin() derive.Plugin {
	return derive.NewPlugin("clone", "deriveClone", New)
}

func New(typesMap derive.TypesMap, p derive.Printer, deps map[string]derive.Dependency) derive.Generator {
	return &gen{
		TypesMap:   typesMap,
		printer:    p,
		bytesPkg:   p.NewImport("bytes"),
		reflectPkg: p.NewImport("reflect"),
		unsafePkg:  p.NewImport("unsafe"),
	}
}

type gen struct {
	derive.TypesMap
	printer    derive.Printer
	bytesPkg   derive.Import
	reflectPkg derive.Import
	unsafePkg  derive.Import
}

func (this *gen) Add(name string, typs []types.Type) (string, error) {
	if len(typs) != 1 {
		return "", fmt.Errorf("%s does not have one argument", name)
	}
	return this.SetFuncName(name, typs[0])
}

func (this *gen) Generate() error {
	for _, typs := range this.ToGenerate() {
		if err := this.genFunc(typs[0]); err != nil {
			return err
		}
	}
	return nil
}

func (g *gen) genFunc(typ types.Type) error {
	p := g.printer
	g.Generating(typ)
	typeStr := g.TypeString(typ)
	p.P("")
	p.P("func %s(this %s) %s {", g.GetFuncName(typ), typeStr, typeStr)
	p.In()
	p.P("var that %s", typeStr)
	if err := g.genStatement(typ, "this", "that"); err != nil {
		return err
	}
	p.P("return that")
	p.Out()
	p.P("}")
	return nil
}

func (g *gen) genStatement(typ types.Type, this, that string) error {
	p := g.printer
	if canClone(typ) {
		p.P("%s = %s", that, this)
		return nil
	}
	switch ttyp := typ.(type) {
	case *types.Pointer:
		p.P("if %s != nil {", this)
		p.In()
		reftyp := ttyp.Elem()
		g.TypeString(reftyp)
		p.P("%s = new(%s)", that, g.TypeString(reftyp))
		thisref, thatref := "*"+this, "*"+that
		named, ok := reftyp.(*types.Named)
		if !ok {
			if err := g.genStatement(reftyp, thisref, thatref); err != nil {
				return err
			}
		} else {
			fields := derive.Fields(g.TypesMap, named.Underlying().(*types.Struct))
			if len(fields.Fields) > 0 {
				if fields.Reflect {
					p.P(`thisv := `+g.reflectPkg()+`.Indirect(`+g.reflectPkg()+`.ValueOf(%s))`, this)
					p.P(`thatv := `+g.reflectPkg()+`.Indirect(`+g.reflectPkg()+`.ValueOf(%s))`, that)
				}
				for _, field := range fields.Fields {
					fieldType := field.Type
					var thisField, thatField string
					if !field.Private() {
						thisField, thatField = field.Name(this, nil), field.Name(that, nil)
					} else {
						thisField, thatField = field.Name("thisv", g.unsafePkg), field.Name("thatv", g.unsafePkg)
					}
					if err := g.genStatement(fieldType, thisField, thatField); err != nil {
						return err
					}
				}
			}
		}
		p.Out()
		p.P("}")
		return nil
	case *types.Named:

	case *types.Slice:
		p.P("if %s != nil {", this)
		p.In()
		p.P("%s = make(%s, len(%s))", that, g.TypeString(typ), this)
		elmType := ttyp.Elem()
		if canClone(elmType) {
			p.P("copy(%s, %s)", that, this)
		} else {
			p.P("for i, thisvalue := range %s {", this)
			p.In()
			g.genStatement(elmType, "thisvalue", that+"[i]")
			p.Out()
			p.P("}")
		}
		p.Out()
		p.P("}")
		return nil
	case *types.Array:
		elmType := ttyp.Elem()
		p.P("for i, thisvalue := range %s {", this)
		p.In()
		g.genStatement(elmType, "thisvalue", that+"[i]")
		p.Out()
		p.P("}")
		return nil
	case *types.Map:

	}
	return fmt.Errorf("unsupported type: %#v", typ)
}

func not(s string) string {
	if strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")") {
		return "!" + s
	}
	return "!(" + s + ")"
}

func wrap(value string) string {
	if strings.HasPrefix(value, "*") || strings.HasPrefix(value, "&") {
		return "(" + value + ")"
	}
	return value
}

func canClone(tt types.Type) bool {
	t := tt.Underlying()
	switch typ := t.(type) {
	case *types.Basic:
		return typ.Kind() != types.UntypedNil
	case *types.Struct:
		for i := 0; i < typ.NumFields(); i++ {
			f := typ.Field(i)
			ft := f.Type()
			if !canClone(ft) {
				return false
			}
		}
		return true
	case *types.Array:
		return canClone(typ.Elem())
	}
	return false
}

func hasCloneMethod(typ *types.Named) bool {
	for i := 0; i < typ.NumMethods(); i++ {
		meth := typ.Method(i)
		if meth.Name() != "Clone" {
			continue
		}
		sig, ok := meth.Type().(*types.Signature)
		if !ok {
			// impossible, but lets check anyway
			continue
		}
		if sig.Params().Len() != 0 {
			continue
		}
		res := sig.Results()
		if res.Len() != 1 {
			continue
		}
		return true
	}
	return false
}
