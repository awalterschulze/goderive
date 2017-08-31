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

// Package gostring contains the implementation of the gostring plugin, which generates the deriveGoString function.
//
// The deriveGoString function returns a string that reproduces the argument's value in valid go syntax.
// The deriveGoString function does a recursive print, even printing pointer values, unlike the default %#v operand.
//
// When goderive walks over your code it is looking for a function that:
//  - was not implemented (or was previously derived) and
//  - has a predefined prefix.
//
// In the following code the deriveGoString function will be found, because
// it was not implemented and it has a prefix deriveGoString.
// This prefix is configurable.
//
//	package main
//
//	type MyStruct struct {
//		Int64     int64
//		StringPtr *string
//	}
//
//	func (this *MyStruct) GoString() string {
//		return deriveGoString(this)
//	}
//
//	func main () {
//		s := "abc"
//		fmt.Printf("%#v", &MyStruct{StringPtr: &s})
//	}
//
// The %#v fmt operand will then invoke the GoString method.
// GoString does a recursive print, even printing pointer values, unlike the default %#v operand.
//
// Supported types:
//	- basic types
//	- named structs
//	- slices
//	- maps
//	- pointers to these types
//	- and many more
// Unsupported types:
//	- chan
//	- interface
//	- function
//	- private fields
//	- unnamed structs
//
// Example output can be found here:
// https://github.com/awalterschulze/goderive/tree/master/example/plugin/gostring
//
// This plugin has been tested thoroughly.
package gostring

import (
	"fmt"
	"go/types"

	"github.com/awalterschulze/goderive/derive"
)

// NewPlugin creates a new gostring plugin.
// This function returns the plugin name, default prefix and a constructor for the gostring code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("gostring", "deriveGoString", New)
}

// New is a constructor for the gostring code generator.
// This generator should be reconstructed for each package.
func New(typesMap derive.TypesMap, p derive.Printer, deps map[string]derive.Dependency) derive.Generator {
	return &gen{
		TypesMap:   typesMap,
		printer:    p,
		strconvPkg: p.NewImport("strconv", "strconv"),
		bytesPkg:   p.NewImport("bytes", "bytes"),
		fmtPkg:     p.NewImport("fmt", "fmt"),
	}
}

type gen struct {
	derive.TypesMap
	printer    derive.Printer
	strconvPkg derive.Import
	bytesPkg   derive.Import
	fmtPkg     derive.Import
}

func (g *gen) Add(name string, typs []types.Type) (string, error) {
	if len(typs) != 1 {
		return "", fmt.Errorf("%s does not have two arguments", name)
	}
	return g.SetFuncName(name, typs[0])
}

func (g *gen) Generate(typs []types.Type) error {
	return g.genFunc(typs[0])
}

func (g *gen) TypeString(typ types.Type) string {
	return g.TypesMap.(bypass).TypeStringBypass(typ)
}

type bypass interface {
	TypeStringBypass(types.Type) string
}

func (g *gen) genFunc(typ types.Type) error {
	p := g.printer
	g.Generating(typ)
	typeStr := g.TypesMap.TypeString(typ)
	gotypeStr := g.TypeString(typ)
	p.P("")
	p.P("func %s(this %s) string {", g.GetFuncName(typ), typeStr)
	p.In()
	p.P("buf := %s.NewBuffer(nil)", g.bytesPkg())
	p.P("%s.Fprintf(buf, \"func() %s {\\n\")", g.fmtPkg(), gotypeStr)
	if err := g.genStatement(typ, "this"); err != nil {
		return err
	}
	p.P("%s.Fprintf(buf, \"}()\\n\")", g.fmtPkg())
	p.P("return buf.String()")
	p.Out()
	p.P("}")
	return nil
}

func (g *gen) W(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	g.printer.P("%s.Fprintf(buf, \"%s\\n\")", g.fmtPkg(), s)
}

func (g *gen) P(format string, a ...interface{}) {
	g.printer.P(format, a...)
}

func isBasicPointer(typ types.Type) bool {
	p, ok := typ.Underlying().(*types.Pointer)
	if !ok {
		return false
	}
	_, ok = p.Elem().(*types.Basic)
	return ok
}

func (g *gen) genStatement(typ types.Type, this string) error {
	p := g.printer
	switch ttyp := typ.Underlying().(type) {
	case *types.Basic:
		p.P("%s.Fprintf(buf, \"return %s\\n\", %s)", g.fmtPkg(), "%#v", this)
		return nil
	case *types.Pointer:
		p.P("if %s == nil {", this)
		p.In()
		g.W("return nil")
		p.Out()
		p.P("} else {")
		p.In()
		reftyp := ttyp.Elem()
		thisref := "*" + this
		named, isNamed := reftyp.(*types.Named)
		strct, isStruct := reftyp.Underlying().(*types.Struct)
		if !isStruct {
			g.W("%s := new(%s)", this, g.TypeString(reftyp))
			g.genField(reftyp, thisref)
			g.W("return %s", this)
		} else {
			gotypeStr := g.TypeString(reftyp)
			external := isNamed && g.TypesMap.IsExternal(named)
			fields := derive.Fields(g.TypesMap, strct, external)
			if len(fields.Fields) == 0 {
				g.W("return &%s{}", gotypeStr)
			} else {
				g.W("%s := &%s{}", this, gotypeStr)
				for _, field := range fields.Fields {
					if field.Private() {
						return fmt.Errorf("private fields not supported, found %s in %v", field.DebugName(), g.TypeString(typ))
					}
					thisField := field.Name(this, nil)
					if err := g.genField(field.Type, thisField); err != nil {
						return err
					}
				}
				g.W("return %s", this)
			}
		}
		p.Out()
		p.P("}")
		return nil
	case *types.Struct:
		fields := derive.Fields(g.TypesMap, ttyp, false)
		gotypeStr := g.TypeString(typ)
		g.W("%s := &%s{}", this, gotypeStr)
		for _, field := range fields.Fields {
			if field.Private() {
				return fmt.Errorf("private fields not supported, found %s in %v", field.DebugName(), g.TypeString(typ))
			}
			thisField := field.Name(this, nil)
			if err := g.genField(field.Type, thisField); err != nil {
				return err
			}
		}
		g.W("return *%s", this)
		return nil
	case *types.Slice:
		p.P("if %s == nil {", this)
		p.In()
		g.W("return nil")
		p.Out()
		p.P("} else {")
		p.In()
		elmTyp := ttyp.Elem()
		if _, isBasic := elmTyp.(*types.Basic); isBasic {
			p.P("%s.Fprintf(buf, \"return %s\\n\", %s)", g.fmtPkg(), "%#v", this)
		} else {
			gotypeStr := g.TypeString(ttyp)
			p.P("%s.Fprintf(buf, \"%s := make(%s, %s)\\n\", %s)", g.fmtPkg(), this, gotypeStr, "%d", "len("+this+")")
			p.P("for i := range %s {", this)
			p.In()
			p.P("%s.Fprintf(buf, \"%s[%s] = %s\\n\", %s, %s)", g.fmtPkg(), this, "%d", "%s", "i", g.GetFuncName(elmTyp)+"("+this+"[i])")
			p.Out()
			p.P("}")
			g.W("return %s", this)
		}
		p.Out()
		p.P("}")
		return nil
	case *types.Array:
		elmTyp := ttyp.Elem()
		if _, isBasic := elmTyp.(*types.Basic); isBasic {
			p.P("%s.Fprintf(buf, \"return %s\\n\", %s)", g.fmtPkg(), "%#v", this)
		} else {
			gotypeStr := g.TypeString(typ)
			p.P("%s.Fprintf(buf, \"%s := %s{}\\n\")", g.fmtPkg(), this, gotypeStr)
			p.P("for i := range %s {", this)
			p.In()
			p.P("%s.Fprintf(buf, \"%s[%s] = %s\\n\", %s, %s)", g.fmtPkg(), this, "%d", "%s", "i", g.GetFuncName(elmTyp)+"("+this+"[i])")
			p.Out()
			p.P("}")
			g.W("return %s", this)
		}
		return nil
	case *types.Map:
		p.P("if %s == nil {", this)
		p.In()
		g.W("return nil")
		p.Out()
		p.P("} else {")
		p.In()
		elmTyp := ttyp.Elem()
		keyTyp := ttyp.Key()
		_, isBasicElm := elmTyp.(*types.Basic)
		_, isBasicKey := keyTyp.(*types.Basic)
		if isBasicElm && isBasicKey {
			p.P("%s.Fprintf(buf, \"return %s\\n\", %s)", g.fmtPkg(), "%#v", this)
		} else if isBasicKey {
			gotypeStr := g.TypeString(typ)
			p.P("%s.Fprintf(buf, \"%s := make(%s)\\n\")", g.fmtPkg(), this, gotypeStr)
			p.P("for k, v := range %s {", this)
			p.In()
			p.P("%s.Fprintf(buf, \"%s[%s] = %s\\n\", %s, %s)", g.fmtPkg(), this, "%#v", "%s", "k", g.GetFuncName(elmTyp)+"(v)")
			p.Out()
			p.P("}")
			g.W("return %s", this)
		} else {
			gotypeStr := g.TypeString(typ)
			p.P("%s.Fprintf(buf, \"%s := make(%s)\\n\")", g.fmtPkg(), this, gotypeStr)
			p.P("i := 0")
			p.P("for k, v := range %s {", this)
			p.In()
			p.P("%s.Fprintf(buf, \"key%s := %s\\n\", %s, %s)", g.fmtPkg(), "%d", "%s", "i", g.GetFuncName(keyTyp)+"(k)")
			p.P("%s.Fprintf(buf, \"%s[key%s] = %s\\n\", %s, %s)", g.fmtPkg(), this, "%d", "%s", "i", g.GetFuncName(elmTyp)+"(v)")
			p.P("i++")
			p.Out()
			p.P("}")
			g.W("return %s", this)
		}
		p.Out()
		p.P("}")
		return nil
	}
	return fmt.Errorf("unsupported root type: %#v", typ)
}

func (g *gen) genField(fieldType types.Type, this string) error {
	p := g.printer
	switch typ := fieldType.Underlying().(type) {
	case *types.Basic:
		p.P("%s.Fprintf(buf, \"%s = %s\\n\", %s)", g.fmtPkg(), this, "%#v", this)
		return nil
	case *types.Pointer:
		p.P("if %s != nil {", this)
		p.In()
		if b, ok := typ.Elem().(*types.Basic); ok {
			p.P("%s.Fprintf(buf, \"%s = func (v %s) *%s { return &v }(%s)\\n\", %s)", g.fmtPkg(), this, g.TypeString(b), g.TypeString(b), "%#v", "*"+this)
		} else {
			p.P("%s.Fprintf(buf, \"%s = %s\\n\", %s)", g.fmtPkg(), this, "%s", g.GetFuncName(typ)+"("+this+")")
		}
		p.Out()
		p.P("}")
		return nil
	case *types.Slice:
		p.P("if %s != nil {", this)
		p.In()
		if _, ok := typ.Elem().(*types.Basic); ok {
			p.P("%s.Fprintf(buf, \"%s = %s\\n\", %s)", g.fmtPkg(), this, "%#v", this)
		} else {
			p.P("%s.Fprintf(buf, \"%s = %s\\n\", %s)", g.fmtPkg(), this, "%s", g.GetFuncName(fieldType)+"("+this+")")
		}
		p.Out()
		p.P("}")
		return nil
	case *types.Array:
		if _, ok := typ.Elem().(*types.Basic); ok {
			p.P("%s.Fprintf(buf, \"%s = %s\\n\", %s)", g.fmtPkg(), this, "%#v", this)
		} else {
			p.P("%s.Fprintf(buf, \"%s = %s\\n\", %s)", g.fmtPkg(), this, "%s", g.GetFuncName(fieldType)+"("+this+")")
		}
		return nil
	case *types.Map:
		p.P("if %s != nil {", this)
		p.In()
		elmTyp := typ.Elem()
		keyTyp := typ.Key()
		_, isBasicElm := elmTyp.(*types.Basic)
		_, isBasicKey := keyTyp.(*types.Basic)
		if isBasicElm && isBasicKey {
			p.P("%s.Fprintf(buf, \"%s = %s\\n\", %s)", g.fmtPkg(), this, "%#v", this)
		} else {
			p.P("%s.Fprintf(buf, \"%s = %s\\n\", %s)", g.fmtPkg(), this, "%s", g.GetFuncName(fieldType)+"("+this+")")
		}
		p.Out()
		p.P("}")
		return nil
	case *types.Struct:
		p.P("%s.Fprintf(buf, \"%s = %s\\n\", %s)", g.fmtPkg(), this, "%s", g.GetFuncName(fieldType)+"("+this+")")
		return nil
	}
	return fmt.Errorf("unsupported field type %#v", fieldType)
}
