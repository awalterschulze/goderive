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
package gostring

import (
	"fmt"
	"go/types"
	"strings"

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
		strconvPkg: p.NewImport("strconv"),
		bytesPkg:   p.NewImport("bytes"),
		fmtPkg:     p.NewImport("fmt"),
	}
}

type gen struct {
	derive.TypesMap
	printer    derive.Printer
	strconvPkg derive.Import
	bytesPkg   derive.Import
	fmtPkg     derive.Import
}

func (this *gen) Add(name string, typs []types.Type) (string, error) {
	if len(typs) != 1 {
		return "", fmt.Errorf("%s does not have two arguments", name)
	}
	return this.SetFuncName(name, typs[0])
}

func (this *gen) Generate(typs []types.Type) error {
	return this.genFunc(typs[0])
}

func (g *gen) genFunc(typ types.Type) error {
	p := g.printer
	g.Generating(typ)
	typeStr := g.TypeString(typ)
	p.P("")
	p.P("func %s(this %s) string {", g.GetFuncName(typ), typeStr)
	p.In()
	p.P("buf := %s.NewBuffer(nil)", g.bytesPkg())
	if err := g.genStatement(typ, "this"); err != nil {
		return err
	}
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

func (g *gen) genStatement(typ types.Type, this string) error {
	p := g.printer
	switch ttyp := typ.Underlying().(type) {
	case *types.Basic:
		g.genField(typ, this)
		return nil
	case *types.Pointer:
		p.P("if %s == nil {", this)
		p.In()
		g.W("nil")
		p.Out()
		p.P("} else {")
		p.In()
		reftyp := ttyp.Elem()
		thisref := "*" + this
		named, isNamed := reftyp.(*types.Named)
		strct, isStruct := reftyp.Underlying().(*types.Struct)
		if !isStruct {
			g.genStatement(reftyp, thisref)
		} else if isNamed {
			external := g.TypesMap.IsExternal(named)
			fields := derive.Fields(g.TypesMap, strct, external)
			if len(fields.Fields) == 0 {
				g.W("&%s{}", g.TypeString(reftyp))
			} else {
				g.W("%s := &%s{}", this, g.TypeString(reftyp))
				for _, field := range fields.Fields {
					// fieldType := field.Type
					if field.Private() {
						return fmt.Errorf("private fields not supported, found %s in %v", field.Name("", nil), named)
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

	case *types.Slice:

	case *types.Array:

	case *types.Map:

	}
	return fmt.Errorf("unsupported type: %#v", typ)
}

var replacer = strings.NewReplacer(".", "_", "*", "_", "(", "_", ")", "_", "&", "_")

func newVar(this string) string {
	that := replacer.Replace(this)
	return that + "_tmp"
}

func hasGoStringMethod(typ *types.Named) bool {
	for i := 0; i < typ.NumMethods(); i++ {
		meth := typ.Method(i)
		if meth.Name() != "GoString" {
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
		b, ok := res.At(0).Type().(*types.Basic)
		if !ok {
			continue
		}
		if b.Kind() != types.String {
			continue
		}
		return true
	}
	return false
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
		// ref := typ.Elem()
		_ = typ
		tmpvar := newVar(this)
		// if named, ok := ref.(*types.Named); ok {
		// 	if hasGoStringMethod(named) {
		// 		p.P("%s.Fprintf(buf, \"%s = %s\\n\", %s)", g.fmtPkg(), tmpvar, "%s", this+".GoString()")
		// 		p.Out()
		// 		p.P("}")
		// 		return nil
		// 	}
		// }
		p.P("%s.Fprintf(buf, \"%s := %s\\n\", %s)", g.fmtPkg(), tmpvar, "%#v", "*"+this)
		p.P("%s.Fprintf(buf, \"%s = %s\\n\")", g.fmtPkg(), this, "*"+tmpvar)
		p.Out()
		p.P("}")
		return nil
		// case *types.Array:
		// 	return fmt.Sprintf("%s(%s)", this.GetFuncName(typ), thisField), nil
		// case *types.Slice:
		// 	return fmt.Sprintf("%s(%s)", this.GetFuncName(typ), thisField), nil
		// case *types.Map:
		// 	return fmt.Sprintf("%s(%s)", this.GetFuncName(typ), thisField), nil
		// case *types.Struct:
		// 	if named, isNamed := fieldType.(*types.Named); isNamed {
		// 		if hasGoStringMethod(named) {
		// 			return fmt.Sprintf("%s.GoString()", thisField), nil
		// 		}
		// 	}
		// 	return fmt.Sprintf("%s(%s)", this.GetFuncName(typ), thisField), nil
		// }
		// *Chan, *Tuple, *Signature, *Interface, *Struct
	}
	return fmt.Errorf("unsupported type %#v", fieldType)
}
