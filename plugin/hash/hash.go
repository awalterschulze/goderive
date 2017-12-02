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

// Package hash contains the implementation of the hash plugin, which generates the deriveHash function.
//
// The deriveHash function is returns a hash of the input object.
//   deriveHash(T) uint64
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
//	- unnamed structs, which are not comparable with the == operator
//
// Example output can be found here:
// https://github.com/awalterschulze/goderive/tree/master/example/plugin/hash
//
// This plugin has been tested thoroughly.
package hash

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/awalterschulze/goderive/derive"
)

// NewPlugin creates a new hash plugin.
// This function returns the plugin name, default prefix and a constructor for the hash code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("hash", "deriveHash", New)
}

// New is a constructor for the hash code generator.
// This generator should be reconstructed for each package.
func New(typesMap derive.TypesMap, p derive.Printer, deps map[string]derive.Dependency) derive.Generator {
	return &gen{
		TypesMap: typesMap,
		printer:  p,
		mathPkg:  p.NewImport("math", "math"),
		keys:     deps["keys"],
		sort:     deps["sort"],
	}
}

type gen struct {
	derive.TypesMap
	printer derive.Printer
	mathPkg derive.Import
	keys    derive.Dependency
	sort    derive.Dependency
}

func (g *gen) Add(name string, typs []types.Type) (string, error) {
	if len(typs) != 1 {
		return "", fmt.Errorf("%s does not have one argument", name)
	}
	return g.SetFuncName(name, typs...)
}

func (g *gen) Generate(typs []types.Type) error {
	return g.genFunc(typs)
}

func (g *gen) genFunc(typs []types.Type) error {
	p := g.printer
	g.Generating(typs...)
	typeStr := g.TypeString(typs[0])
	name := g.GetFuncName(typs...)
	p.P("")
	p.P("// %s returns the hash of the object.", name)
	if strct, ok := typs[0].(*types.Struct); ok {
		fields := derive.GetStructFields(strct)
		fieldStrs, err := g.FieldStrings(fields)
		if err != nil {
			return err
		}
		p.P("func %s(object struct {", name)
		p.In()
		for _, fieldStr := range fieldStrs {
			p.P(fieldStr)
		}
		p.Out()
		p.P("}) uint64 {")
	} else {
		p.P("func %s(object %s) uint64 {", name, typeStr)
	}
	p.In()
	if err := g.genStatement("object", typs[0]); err != nil {
		return nil
	}
	p.Out()
	p.P("}")
	return nil
}

func (g *gen) genStatement(o string, typ types.Type) error {
	p := g.printer
	switch ttyp := typ.Underlying().(type) {
	case *types.Basic:
		switch ttyp.Kind() {
		case types.Bool, types.UntypedBool:
			p.P("if %s {", o)
			p.In()
			p.P("return 1")
			p.Out()
			p.P("}")
			p.P("return 0")
			return nil
		case types.String, types.UntypedString:
			p.P("var h uint64")
			p.P("for _, c := range %s {", o)
			p.In()
			p.P("h = 31*h + uint64(c)")
			p.Out()
			p.P("}")
			p.P("return h")
			return nil
		}
		fieldStr, err := g.field(o, typ)
		if err != nil {
			return err
		}
		p.P("return " + fieldStr)
		return nil
	case *types.Pointer:
		ref := "*" + o
		reftyp := ttyp.Elem()
		named, isNamed := reftyp.(*types.Named)
		strct, isStruct := reftyp.Underlying().(*types.Struct)
		p.P("if %s == nil {", o)
		p.In()
		p.P("return 0")
		p.Out()
		p.P("}")
		if isStruct && isNamed {
			external := g.TypesMap.IsExternal(named)
			fields := derive.Fields(g.TypesMap, strct, external)
			if len(fields.Fields) == 0 {
				p.P("return 17")
				return nil
			}
			p.P("h := uint64(17)")
			for _, field := range fields.Fields {
				fieldType := field.Type
				if field.Private() && external {
					continue
				}
				fieldName := field.Name(o, nil)
				fieldStr, err := g.field(fieldName, fieldType)
				if err != nil {
					return err
				}
				p.P("h = 31*h + %s", fieldStr)
			}
			p.P("return h")
			return nil
		} else {
			fieldStr, err := g.field(ref, reftyp)
			if err != nil {
				return err
			}
			p.P("return (31 * 17) + %s", fieldStr)
			return nil
		}
	case *types.Struct:
		if _, isNamed := typ.(*types.Named); isNamed {
			fieldStr, err := g.field("&"+o, types.NewPointer(typ))
			if err != nil {
				return err
			}
			p.P("return " + fieldStr)
			return nil
		} else {
			fields := derive.Fields(g.TypesMap, ttyp, false)
			if len(fields.Fields) == 0 {
				p.P("return 17")
				return nil
			}
			p.P("h := uint64(17)")
			for _, field := range fields.Fields {
				fieldType := field.Type
				fieldName := field.Name(o, nil)
				fieldStr, err := g.field(fieldName, fieldType)
				if err != nil {
					return err
				}
				p.P("h = 31*h + %s", fieldStr)
			}
			p.P("return h")
			return nil
		}
	case *types.Slice:
		p.P("if %s == nil {", o)
		p.In()
		p.P("return 0")
		p.Out()
		p.P("}")
		p.P("h := uint64(17)")
		p.P("for i := 0; i < len(%s); i++ {", o)
		p.In()
		elem := wrap(o) + "[i]"
		fieldStr, err := g.field(elem, ttyp.Elem())
		if err != nil {
			return err
		}
		p.P("h = 31*h + %s", fieldStr)
		p.Out()
		p.P("}")
		p.P("return h")
		return nil
	case *types.Array:
		p.P("h := uint64(17)")
		p.P("for i := 0; i < len(%s); i++ {", o)
		p.In()
		elem := wrap(o) + "[i]"
		fieldStr, err := g.field(elem, ttyp.Elem())
		if err != nil {
			return err
		}
		p.P("h = 31*h + %s", fieldStr)
		p.Out()
		p.P("}")
		p.P("return h")
		return nil
	case *types.Map:
		p.P("if %s == nil {", o)
		p.In()
		p.P("return 0")
		p.Out()
		p.P("}")
		p.P("h := uint64(17)")
		p.P("for _, k := range %s(%s(%s)) {", g.sort.GetFuncName(types.NewSlice(ttyp.Key())), g.keys.GetFuncName(typ), o)
		p.In()
		keyStr, err := g.field("k", ttyp.Key())
		if err != nil {
			return err
		}
		p.P("h = 31*h + %s", keyStr)
		valStr, err := g.field(o+"[k]", ttyp.Elem())
		if err != nil {
			return err
		}
		p.P("h = 31*h + %s", valStr)
		p.Out()
		p.P("}")
		p.P("return h")
		return nil
	}
	return fmt.Errorf("unsupported type: %#v", typ)
}

func wrap(value string) string {
	if strings.HasPrefix(value, "*") || strings.HasPrefix(value, "&") {
		return "(" + value + ")"
	}
	return value
}

func hasHashMethod(typ *types.Named) bool {
	for i := 0; i < typ.NumMethods(); i++ {
		meth := typ.Method(i)
		if meth.Name() != "Hash" {
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
		if b.Kind() != types.Int32 {
			continue
		}
		return true
	}
	return false
}

func (g *gen) field(fieldName string, fieldType types.Type) (string, error) {
	switch typ := fieldType.Underlying().(type) {
	case *types.Basic:
		switch typ.Kind() {
		case types.UntypedNil:
			return "0", nil
		case types.Bool, types.UntypedBool:
			return fmt.Sprintf("%s(%s)", g.GetFuncName(fieldType), fieldName), nil
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
			types.Uint, types.Uint8, types.Uint16, types.Uint32,
			types.Uintptr, types.UnsafePointer, types.UntypedInt:
			return fmt.Sprintf("uint64(%s)", fieldName), nil
		case types.Uint64:
			return fmt.Sprintf("%s", fieldName), nil
		case types.Float32:
			return fmt.Sprintf("uint64(%s.Float32bits(%s))", g.mathPkg(), fieldName), nil
		case types.Float64:
			return fmt.Sprintf("%s.Float64bits(%s)", g.mathPkg(), fieldName), nil
		case types.Complex64:
			return fmt.Sprintf("(31 * ((31 * 17) + uint64(%s.Float32bits(real(%s))))) + uint64(%s.Float32bits(imag(%s)))", g.mathPkg(), fieldName, g.mathPkg(), fieldName), nil
		case types.Complex128:
			return fmt.Sprintf("(31 * ((31 * 17) + %s.Float64bits(real(%s)))) + %s.Float64bits(imag(%s))", g.mathPkg(), fieldName, g.mathPkg(), fieldName), nil
		case types.String, types.UntypedString:
			return fmt.Sprintf("%s(%s)", g.GetFuncName(fieldType), fieldName), nil
		}
	case *types.Pointer:
		ref := typ.Elem()
		if named, ok := ref.(*types.Named); ok {
			if hasHashMethod(named) {
				return fmt.Sprintf("%s.Hash()", wrap(fieldName)), nil
			}
		}
		return fmt.Sprintf("%s(%s)", g.GetFuncName(fieldType), fieldName), nil
	case *types.Array:
		return fmt.Sprintf("%s(%s)", g.GetFuncName(fieldType), fieldName), nil
	case *types.Slice:
		return fmt.Sprintf("%s(%s)", g.GetFuncName(fieldType), fieldName), nil
	case *types.Map:
		return fmt.Sprintf("%s(%s)", g.GetFuncName(fieldType), fieldName), nil
	case *types.Struct:
		if named, isNamed := fieldType.(*types.Named); isNamed {
			if hasHashMethod(named) {
				return fmt.Sprintf("%s.Hash()", wrap(fieldName)), nil
			}
		}
		return fmt.Sprintf("%s(%s)", g.GetFuncName(fieldType), fieldName), nil
	}
	// *Chan, *Tuple, *Signature, *Interface, *types.Basic.Kind() == types.UntypedNil, *Struct
	return "", fmt.Errorf("unsupported type %#v", fieldType)
}
