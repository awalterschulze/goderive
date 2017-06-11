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

package equal

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/awalterschulze/goderive/derive"
)

func NewPlugin() derive.Plugin {
	return derive.NewPlugin("equal", "deriveEqual", New)
}

func New(typesMap derive.TypesMap, p derive.Printer, deps map[string]derive.Dependency) derive.Generator {
	return &equal{
		TypesMap:   typesMap,
		printer:    p,
		bytesPkg:   p.NewImport("bytes"),
		reflectPkg: p.NewImport("reflect"),
		unsafePkg:  p.NewImport("unsafe"),
	}
}

type equal struct {
	derive.TypesMap
	printer    derive.Printer
	bytesPkg   derive.Import
	reflectPkg derive.Import
	unsafePkg  derive.Import
}

func (this *equal) Add(name string, typs []types.Type) (string, error) {
	if len(typs) != 2 {
		return "", fmt.Errorf("%s does not have two arguments", name)
	}
	if !types.Identical(typs[0], typs[1]) {
		return "", fmt.Errorf("%s has two arguments, but they are of different types %s != %s",
			name, this.TypeString(typs[0]), this.TypeString(typs[1]))
	}
	return this.SetFuncName(name, typs[0])
}

func (this *equal) Generate() error {
	for _, typs := range this.ToGenerate() {
		if err := this.genFunc(typs[0]); err != nil {
			return err
		}
	}
	for _, typs := range this.ToGenerate() {
		if err := this.genFunc(typs[0]); err != nil {
			return err
		}
	}
	return nil
}

func (g *equal) genFunc(typ types.Type) error {
	p := g.printer
	g.Generating(typ)
	typeStr := g.TypeString(typ)
	p.P("")
	p.P("func %s(this, that %s) bool {", g.GetFuncName(typ), typeStr)
	p.In()
	if err := g.genStatement(typ, "this", "that"); err != nil {
		return nil
	}
	p.Out()
	p.P("}")
	return nil
}

func (g *equal) genStatement(typ types.Type, this, that string) error {
	p := g.printer
	switch ttyp := typ.(type) {
	case *types.Basic:
		fieldStr, err := g.field(this, that, typ)
		if err != nil {
			return err
		}
		p.P("return " + fieldStr)
		return nil
	case *types.Pointer:
		thisref, thatref := "*"+this, "*"+that
		reftyp := ttyp.Elem()
		named, ok := reftyp.(*types.Named)
		if !ok {
			p.P("if %s == nil && %s == nil {", this, that)
			p.In()
			p.P("return true")
			p.Out()
			p.P("}")
			p.P("if %s != nil && %s != nil {", this, that)
			p.In()
			if err := g.genStatement(reftyp, thisref, thatref); err != nil {
				return err
			}
			p.Out()
			p.P("}")
			p.P("return false")
			return nil
		}
		fields := derive.Fields(g.TypesMap, named.Underlying().(*types.Struct))
		if len(fields.Fields) == 0 {
			p.P("return %s == nil && %s == nil", this, that)
			return nil
		}
		if fields.Reflect {
			p.P(`thisv := `+g.reflectPkg()+`.Indirect(`+g.reflectPkg()+`.ValueOf(%s))`, this)
			p.P(`thatv := `+g.reflectPkg()+`.Indirect(`+g.reflectPkg()+`.ValueOf(%s))`, that)
		}
		p.P("return (%s == nil && %s == nil) ||", this, that)
		p.In()
		p.P("%s != nil && %s != nil &&", this, that)
		for i, field := range fields.Fields {
			fieldType := field.Type
			var thisField, thatField string
			if !field.Private() {
				thisField, thatField = field.Name(this, nil), field.Name(that, nil)
			} else {
				thisField, thatField = field.Name("thisv", g.unsafePkg), field.Name("thatv", g.unsafePkg)
			}
			fieldStr, err := g.field(thisField, thatField, fieldType)
			if err != nil {
				return err
			}
			if (i + 1) != len(fields.Fields) {
				fieldStr += " &&"
			}
			if i == 0 {
				p.In()
			}
			p.P(fieldStr)
		}
		p.Out()
		p.Out()
		return nil
	case *types.Struct:
		if canEqual(ttyp) {
			p.P("return %s == %s", this, that)
			return nil
		}
	case *types.Named:
		fieldStr, err := g.field("&"+this, "&"+that, types.NewPointer(ttyp))
		if err != nil {
			return err
		}
		p.P("return " + fieldStr)
		return nil
	case *types.Slice:
		p.P("if %s == nil || %s == nil {", this, that)
		p.In()
		p.P("return %s == nil && %s == nil", this, that)
		p.Out()
		p.P("}")
		p.P("if len(%s) != len(%s) {", this, that)
		p.In()
		p.P("return false")
		p.Out()
		p.P("}")
		p.P("for i := 0; i < len(%s); i++ {", this)
		p.In()
		thisElem, thatElem := wrap(this)+"[i]", wrap(that)+"[i]"
		eqStr, err := g.field(thisElem, thatElem, ttyp.Elem())
		if err != nil {
			return err
		}
		p.P("if %s {", not(eqStr))
		p.In()
		p.P("return false")
		p.Out()
		p.P("}")
		p.Out()
		p.P("}")
		p.P("return true")
		return nil
	case *types.Array:
		p.P("for i := 0; i < len(%s); i++ {", this)
		p.In()
		thisElem, thatElem := wrap(this)+"[i]", wrap(that)+"[i]"
		eqStr, err := g.field(thisElem, thatElem, ttyp.Elem())
		if err != nil {
			return err
		}
		p.P("if %s {", not(eqStr))
		p.In()
		p.P("return false")
		p.Out()
		p.P("}")
		p.Out()
		p.P("}")
		p.P("return true")
		return nil
	case *types.Map:
		p.P("if %s == nil || %s == nil {", this, that)
		p.In()
		p.P("return %s == nil && %s == nil", this, that)
		p.Out()
		p.P("}")
		p.P("if len(%s) != len(%s) {", this, that)
		p.In()
		p.P("return false")
		p.Out()
		p.P("}")
		p.P("for k, v := range %s {", this)
		p.In()
		p.P("thatv, ok := %s[k]", wrap(that))
		p.P("if !ok {")
		p.In()
		p.P("return false")
		p.Out()
		p.P("}")
		eqStr, err := g.field("v", "thatv", ttyp.Elem())
		if err != nil {
			return err
		}
		p.P("if %s {", not(eqStr))
		p.In()
		p.P("return false")
		p.Out()
		p.P("}")
		p.Out()
		p.P("}")
		p.P("return true")
		return nil
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

func canEqual(tt types.Type) bool {
	t := tt.Underlying()
	switch typ := t.(type) {
	case *types.Basic:
		return typ.Kind() != types.UntypedNil
	case *types.Struct:
		for i := 0; i < typ.NumFields(); i++ {
			f := typ.Field(i)
			ft := f.Type()
			if !canEqual(ft) {
				return false
			}
		}
		return true
	case *types.Array:
		return canEqual(typ.Elem())
	}
	return false
}

func hasEqualMethod(typ *types.Named) bool {
	for i := 0; i < typ.NumMethods(); i++ {
		meth := typ.Method(i)
		if meth.Name() != "Equal" {
			continue
		}
		sig, ok := meth.Type().(*types.Signature)
		if !ok {
			// impossible, but lets check anyway
			continue
		}
		if sig.Params().Len() != 1 {
			continue
		}
		res := sig.Results()
		if res.Len() != 1 {
			continue
		}
		b, ok := res.At(0).Type().(*types.Basic)
		if !ok {
			continue
		}
		if b.Kind() != types.Bool {
			continue
		}
		return true
	}
	return false
}

func (this *equal) field(thisField, thatField string, fieldType types.Type) (string, error) {
	if canEqual(fieldType) {
		return fmt.Sprintf("%s == %s", thisField, thatField), nil
	}
	switch typ := fieldType.(type) {
	case *types.Pointer:
		ref := typ.Elem()
		if named, ok := ref.(*types.Named); ok {
			if hasEqualMethod(named) {
				return fmt.Sprintf("%s.Equal(%s)", wrap(thisField), thatField), nil
			} else {
				return fmt.Sprintf("%s(%s, %s)", this.GetFuncName(typ), thisField, thatField), nil
			}
		}
		eqStr, err := this.field("*("+thisField+")", "*("+thatField+")", ref)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("((%[1]s == nil && %[2]s == nil) || (%[1]s != nil && %[2]s != nil && %[3]s))", thisField, thatField, eqStr), nil
	case *types.Array:
		return fmt.Sprintf("%s(%s, %s)", this.GetFuncName(typ), thisField, thatField), nil
	case *types.Slice:
		if b, ok := typ.Elem().(*types.Basic); ok && b.Kind() == types.Byte {
			return fmt.Sprintf("%s.Equal(%s, %s)", this.bytesPkg(), thisField, thatField), nil
		}
		return fmt.Sprintf("%s(%s, %s)", this.GetFuncName(typ), thisField, thatField), nil
	case *types.Map:
		return fmt.Sprintf("%s(%s, %s)", this.GetFuncName(typ), thisField, thatField), nil
	case *types.Named:
		if hasEqualMethod(typ) {
			return fmt.Sprintf("%s.Equal(&%s)", thisField, thatField), nil
		} else {
			return this.field("&"+thisField, "&"+thatField, types.NewPointer(fieldType))
		}
	default: // *Chan, *Tuple, *Signature, *Interface, *types.Basic.Kind() == types.UntypedNil, *Struct
		return "", fmt.Errorf("unsupported type %#v", fieldType)
	}
}
