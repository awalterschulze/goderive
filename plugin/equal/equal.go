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

// Package equal contains the implementation of the equal plugin, which generates the deriveEqual function.
//
// The deriveEqual function is a faster alternative to reflect.DeepEqual.
//   deriveEqual(T, T) bool
//   deriveEqual(T) func(T) bool
//
// When goderive walks over your code it is looking for a function that:
//  - was not implemented (or was previously derived) and
//  - has a predefined prefix.
//
// In the following code the deriveEqual function will be found, because
// it was not implemented and it has a prefix deriveEqual.
// This prefix is configurable.
//
//	package main
//
//	type MyStruct struct {
//		Int64     int64
//		StringPtr *string
//	}
//
//	func (this *MyStruct) Equal(that *MyStruct) bool {
//		return deriveEqual(this, that)
//	}
//
// goderive will then generate the following code in a derived.gen.go file in the same package:
//
//	func deriveEqual(this, that *MyStruct) bool {
//		return (this == nil && that == nil) ||
//			this != nil && that != nil &&
//			this.Int64 == that.Int64 &&
//			((this.StringPtr == nil && that.StringPtr == nil) ||
//				(this.StringPtr != nil && that.StringPtr != nil && *(this.StringPtr) == *(that.StringPtr)))
//	}
//
// Supported types:
//	- basic types
//	- named structs
//	- slices
//	- maps
//	- pointers to these types
//	- private fields of structs in external packages (using reflect and unsafe)
//	- and many more
// Unsupported types:
//	- chan
//	- interface
//	- function
//	- unnamed structs, which are not comparable with the == operator
//
// Example output can be found here:
// https://github.com/awalterschulze/goderive/tree/master/example/plugin/equal
//
// This plugin has been tested thoroughly.
package equal

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/awalterschulze/goderive/derive"
)

// NewPlugin creates a new equal plugin.
// This function returns the plugin name, default prefix and a constructor for the equal code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("equal", "deriveEqual", New)
}

// New is a constructor for the equal code generator.
// This generator should be reconstructed for each package.
func New(typesMap derive.TypesMap, p derive.Printer, deps map[string]derive.Dependency) derive.Generator {
	return &gen{
		TypesMap:   typesMap,
		printer:    p,
		bytesPkg:   p.NewImport("bytes", "bytes"),
		reflectPkg: p.NewImport("reflect", "reflect"),
		unsafePkg:  p.NewImport("unsafe", "unsafe"),
	}
}

type gen struct {
	derive.TypesMap
	printer    derive.Printer
	bytesPkg   derive.Import
	reflectPkg derive.Import
	unsafePkg  derive.Import
}

func (g *gen) Add(name string, typs []types.Type) (string, error) {
	if len(typs) != 2 && len(typs) != 1 {
		return "", fmt.Errorf("%s does not have one or two arguments", name)
	}
	if len(typs) == 2 {
		if !types.Identical(typs[0], typs[1]) {
			return "", fmt.Errorf("%s has two arguments, but they are of different types %s != %s",
				name, g.TypeString(typs[0]), g.TypeString(typs[1]))
		}
	}
	return g.SetFuncName(name, typs...)
}

func (g *gen) Generate(typs []types.Type) error {
	if len(typs) == 1 {
		return g.genCurriedFunc(typs[0])
	}
	return g.genFunc(typs)
}

func (g *gen) genCurriedFunc(typ types.Type) error {
	p := g.printer
	g.Generating(typ)
	typeStr := g.TypeString(typ)
	p.P("")
	p.P("func %s(this %s) func(%s) bool {", g.GetFuncName(typ), typeStr, typeStr)
	p.In()
	p.P("return func(that %s) bool {", typeStr)
	p.In()
	if err := g.genStatement(typ, "this", "that"); err != nil {
		return nil
	}
	p.Out()
	p.P("}")
	p.Out()
	p.P("}")
	return nil
}

func (g *gen) genFunc(typs []types.Type) error {
	p := g.printer
	g.Generating(typs...)
	typeStr := g.TypeString(typs[0])
	p.P("")
	p.P("func %s(this, that %s) bool {", g.GetFuncName(typs...), typeStr)
	p.In()
	if err := g.genStatement(typs[0], "this", "that"); err != nil {
		return nil
	}
	p.Out()
	p.P("}")
	return nil
}

func (g *gen) genStatement(typ types.Type, this, that string) error {
	p := g.printer
	switch ttyp := typ.Underlying().(type) {
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
		named, isNamed := reftyp.(*types.Named)
		strct, isStruct := reftyp.Underlying().(*types.Struct)
		if !isStruct {
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
		if isNamed {
			external := g.TypesMap.IsExternal(named)
			fields := derive.Fields(g.TypesMap, strct, external)
			if len(fields.Fields) == 0 {
				p.P("return (%s == nil && %s == nil) || (%s != nil) && (%s != nil)", this, that, this, that)
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
				if field.Private() && external {
					thisField, thatField = field.Name("thisv", g.unsafePkg), field.Name("thatv", g.unsafePkg)
				} else {
					thisField, thatField = field.Name(this, nil), field.Name(that, nil)
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
		}
	case *types.Struct:
		if canEqual(ttyp) {
			p.P("return %s == %s", this, that)
			return nil
		}
		if _, isNamed := typ.(*types.Named); isNamed {
			fieldStr, err := g.field("&"+this, "&"+that, types.NewPointer(ttyp))
			if err != nil {
				return err
			}
			p.P("return " + fieldStr)
			return nil
		}
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

func equalMethodInputParam(typ *types.Named) *types.Type {
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
		inputType := sig.Params().At(0).Type()
		return &inputType
	}
	return nil
}

func (g *gen) field(thisField, thatField string, fieldType types.Type) (string, error) {
	if canEqual(fieldType) {
		return fmt.Sprintf("%s == %s", thisField, thatField), nil
	}
	switch typ := fieldType.Underlying().(type) {
	case *types.Pointer:
		ref := typ.Elem()
		if named, ok := ref.(*types.Named); ok {
			inputType := equalMethodInputParam(named)
			if inputType != nil {
				ityp := *inputType
				if _, ok := ityp.(*types.Pointer); ok {
					return fmt.Sprintf("%s.Equal(%s)", wrap(thisField), thatField), nil
				} else if _, ok := ityp.(*types.Interface); ok {
					return fmt.Sprintf("%s.Equal(%s)", wrap(thisField), thatField), nil
				} else {
					// fall through to deferencing of pointers
				}
			} else {
				return fmt.Sprintf("%s(%s, %s)", g.GetFuncName(typ, typ), thisField, thatField), nil
			}
		}
		eqStr, err := g.field("*("+thisField+")", "*("+thatField+")", ref)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("((%[1]s == nil && %[2]s == nil) || (%[1]s != nil && %[2]s != nil && %[3]s))", thisField, thatField, eqStr), nil
	case *types.Array:
		return fmt.Sprintf("%s(%s, %s)", g.GetFuncName(typ, typ), thisField, thatField), nil
	case *types.Slice:
		if b, ok := typ.Elem().(*types.Basic); ok && b.Kind() == types.Byte {
			return fmt.Sprintf("%s.Equal(%s, %s)", g.bytesPkg(), thisField, thatField), nil
		}
		return fmt.Sprintf("%s(%s, %s)", g.GetFuncName(typ, typ), thisField, thatField), nil
	case *types.Map:
		return fmt.Sprintf("%s(%s, %s)", g.GetFuncName(typ, typ), thisField, thatField), nil
	case *types.Struct:
		if named, isNamed := fieldType.(*types.Named); isNamed {
			inputType := equalMethodInputParam(named)
			if inputType != nil {
				ityp := *inputType
				if _, ok := ityp.(*types.Pointer); ok {
					return fmt.Sprintf("%s.Equal(&%s)", wrap(thisField), thatField), nil
				} else if _, ok := ityp.(*types.Interface); ok {
					return fmt.Sprintf("%s.Equal(&%s)", wrap(thisField), thatField), nil
				} else {
					return fmt.Sprintf("%s.Equal(%s)", wrap(thisField), thatField), nil
				}
			}
		}
		return g.field("&"+thisField, "&"+thatField, types.NewPointer(fieldType))
	default: // *Chan, *Tuple, *Signature, *Interface, *types.Basic.Kind() == types.UntypedNil, *Struct
		return "", fmt.Errorf("unsupported type %#v", fieldType)
	}
}
