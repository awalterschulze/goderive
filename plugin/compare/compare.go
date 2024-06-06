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

// Package compare contains the implementation of the compare plugin, which generates the deriveCompare function.
//
// The deriveCompare function is a maintainable and fast way to implement fast Less functions.
//   deriveCompare(T, T) bool
//   deriveCompare(T) func(T) bool
//
// When goderive walks over your code it is looking for a function that:
//  - was not implemented (or was previously derived) and
//  - has a predefined prefix.
//
// In the following code the deriveCompare function will be found, because
// it was not implemented and it has a prefix deriveCompare.
// This prefix is configurable.
//
//	package main
//
//	import "sort"
//
//	type MyStruct struct {
//		Int64     int64
//		StringPtr *string
//	}
//
//	func sortStructs(ss []*MyStruct) {
//		sort.Slice(ss,  func(i, j int) bool {
//				return deriveCompare(ss[i], ss[j]) < 0
//		})
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
// https://github.com/ndeloof/goderive/tree/master/example/plugin/equal
//
// This plugin has been tested thoroughly.
package compare

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/ndeloof/goderive/derive"
)

// NewPlugin creates a new compare plugin.
// This function returns the plugin name, default prefix and a constructor for the compare code generator.
func NewPlugin() derive.Plugin {
	return derive.NewPlugin("compare", "deriveCompare", New)
}

// New is a constructor for the compare code generator.
// This generator should be reconstructed for each package.
func New(typesMap derive.TypesMap, p derive.Printer, deps map[string]derive.Dependency) derive.Generator {
	return &gen{
		TypesMap:   typesMap,
		printer:    p,
		bytesPkg:   p.NewImport("bytes", "bytes"),
		stringsPkg: p.NewImport("strings", "strings"),
		reflectPkg: p.NewImport("reflect", "reflect"),
		unsafePkg:  p.NewImport("unsafe", "unsafe"),
		keys:       deps["keys"],
		sort:       deps["sort"],
	}
}

type gen struct {
	derive.TypesMap
	printer    derive.Printer
	bytesPkg   derive.Import
	stringsPkg derive.Import
	reflectPkg derive.Import
	unsafePkg  derive.Import
	keys       derive.Dependency
	sort       derive.Dependency
}

func (g *gen) Add(name string, typs []types.Type) (string, error) {
	if len(typs) != 1 && len(typs) != 2 {
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

func compareMethodInputParam(typ *types.Named) *types.Type {
	for i := 0; i < typ.NumMethods(); i++ {
		meth := typ.Method(i)
		if meth.Name() != "Compare" {
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
		if b.Kind() != types.Int {
			continue
		}
		inputType := sig.Params().At(0).Type()
		return &inputType
	}
	return nil
}

func (g *gen) genCurriedFunc(typ types.Type) error {
	p := g.printer
	g.Generating(typ)
	typeStr := g.TypeString(typ)
	p.P("")
	p.P("// %s returns a curried compare function, which returns:", g.GetFuncName(typ))
	p.P("//   * 0 if this and that are equal,")
	p.P("//   * -1 is this is smaller and")
	p.P("//   * +1 is this is bigger.")
	p.P("func %s(this %s) func(%s) int {", g.GetFuncName(typ), typeStr, typeStr)
	p.In()
	p.P("return func(that %s) int {", typeStr)
	p.In()
	if err := g.genStatement(typ, "this", "that"); err != nil {
		return err
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
	p.P("// %s returns:", g.GetFuncName(typs...))
	p.P("//   * 0 if this and that are equal,")
	p.P("//   * -1 is this is smaller and")
	p.P("//   * +1 is this is bigger.")
	p.P("func %s(this, that %s) int {", g.GetFuncName(typs...), typeStr)
	p.In()
	if err := g.genStatement(typs[0], "this", "that"); err != nil {
		return err
	}
	p.Out()
	p.P("}")
	return nil
}

func (g *gen) genStatement(typ types.Type, this, that string) error {
	p := g.printer
	switch ttyp := typ.Underlying().(type) {
	case *types.Pointer:
		p.P("if %s == nil {", this)
		p.In()
		p.P("if %s == nil {", that)
		p.In()
		p.P("return 0")
		p.Out()
		p.P("}")
		p.P("return -1")
		p.Out()
		p.P("}")
		p.P("if %s == nil {", that)
		p.In()
		p.P("return 1")
		p.Out()
		p.P("}")
		reftyp := ttyp.Elem()
		named, isNamed := reftyp.(*types.Named)
		strct, isStruct := reftyp.Underlying().(*types.Struct)
		if !isStruct || !isNamed {
			p.P("return %s(*%s, *%s)", g.GetFuncName(reftyp, reftyp), this, that)
			return nil
		}
		external := g.TypesMap.IsExternal(named)
		fields := derive.Fields(g.TypesMap, strct, external)
		if fields.Reflect {
			p.P(`thisv := ` + g.reflectPkg() + `.Indirect(` + g.reflectPkg() + `.ValueOf(` + this + `))`)
			p.P(`thatv := ` + g.reflectPkg() + `.Indirect(` + g.reflectPkg() + `.ValueOf(` + that + `))`)
		}
		for _, field := range fields.Fields {
			fieldType := field.Type
			var thisField, thatField string
			if field.Private() && external {
				thisField = field.Name("thisv", g.unsafePkg)
				thatField = field.Name("thatv", g.unsafePkg)
			} else {
				thisField = field.Name(this, nil)
				thatField = field.Name(that, nil)
			}
			fieldStr, err := g.field(thisField, thatField, fieldType)
			if err != nil {
				return err
			}
			p.P("if c := %s; c != 0 {", fieldStr)
			p.In()
			p.P("return c")
			p.Out()
			p.P("}")
		}
		p.P("return 0")
		return nil
	case *types.Basic:
		switch ttyp.Kind() {
		case types.String:
			p.P("return %s.Compare(%s, %s)", g.stringsPkg(), this, that)
		case types.Complex128, types.Complex64:
			p.P("if thisr, thatr := real(%s), real(%s); thisr == thatr {", this, that)
			p.In()
			p.P("if thisi, thati := imag(%s), imag(%s); thisi == thati {", this, that)
			p.In()
			p.P("return 0")
			p.Out()
			p.P(`} else if thisi < thati {`)
			p.In()
			p.P("return -1")
			p.Out()
			p.P(`} else {`)
			p.In()
			p.P("return 1")
			p.Out()
			p.P(`}`)
			p.Out()
			p.P(`} else if thisr < thatr {`)
			p.In()
			p.P("return -1")
			p.Out()
			p.P(`} else {`)
			p.In()
			p.P("return 1")
			p.Out()
			p.P(`}`)
		case types.Bool:
			p.P("if %s == %s {", this, that)
			p.In()
			p.P("return 0")
			p.Out()
			p.P("}")
			p.P("if %s {", that)
			p.In()
			p.P("return -1")
			p.Out()
			p.P("}")
			p.P("return 1")
		default:
			p.P("if %s != %s {", this, that)
			p.In()
			p.P("if %s < %s {", this, that)
			p.In()
			p.P("return -1")
			p.Out()
			p.P("} else {")
			p.In()
			p.P("return 1")
			p.Out()
			p.P("}")
			p.Out()
			p.P("}")
			p.P("return 0")
		}
		return nil
	case *types.Struct:
		if _, isNamed := typ.(*types.Named); isNamed {
			fieldStr, err := g.field("&"+this, "&"+that, types.NewPointer(typ))
			if err != nil {
				return err
			}
			p.P("return " + fieldStr)
			return nil
		}
	case *types.Slice:
		p.P("if %s == nil {", this)
		p.In()
		p.P("if %s == nil {", that)
		p.In()
		p.P("return 0")
		p.Out()
		p.P("}")
		p.P("return -1")
		p.Out()
		p.P("}")
		p.P("if %s == nil {", that)
		p.In()
		p.P("return 1")
		p.Out()
		p.P("}")
		p.P("if len(%s) != len(%s) {", this, that)
		p.In()
		p.P("if len(%s) < len(%s) {", this, that)
		p.In()
		p.P("return -1")
		p.Out()
		p.P("}")
		p.P("return 1")
		p.Out()
		p.P("}")
		p.P("for i := 0; i < len(%s); i++ {", this)
		p.In()
		cmpStr, err := g.field(this+"[i]", that+"[i]", ttyp.Elem())
		if err != nil {
			return err
		}
		p.P("if c := %s; c != 0 {", cmpStr)
		p.In()
		p.P("return c")
		p.Out()
		p.P("}")
		p.Out()
		p.P("}")
		p.P("return 0")
		return nil
	case *types.Array:
		p.P("if len(%s) != len(%s) {", this, that)
		p.In()
		p.P("if len(%s) < len(%s) {", this, that)
		p.In()
		p.P("return -1")
		p.Out()
		p.P("}")
		p.P("return 1")
		p.Out()
		p.P("}")
		p.P("for i := 0; i < len(%s); i++ {", this)
		p.In()
		cmpStr, err := g.field(this+"[i]", that+"[i]", ttyp.Elem())
		if err != nil {
			return err
		}
		p.P("if c := %s; c != 0 {", cmpStr)
		p.In()
		p.P("return c")
		p.Out()
		p.P("}")
		p.Out()
		p.P("}")
		p.P("return 0")
		return nil
	case *types.Map:
		p.P("if %s == nil {", this)
		p.In()
		p.P("if %s == nil {", that)
		p.In()
		p.P("return 0")
		p.Out()
		p.P("}")
		p.P("return -1")
		p.Out()
		p.P("}")
		p.P("if %s == nil {", that)
		p.In()
		p.P("return 1")
		p.Out()
		p.P("}")
		p.P("if len(%s) != len(%s) {", this, that)
		p.In()
		p.P("if len(%s) < len(%s) {", this, that)
		p.In()
		p.P("return -1")
		p.Out()
		p.P("}")
		p.P("return 1")
		p.Out()
		p.P("}")
		p.P("thiskeys := %s(%s(%s))", g.sort.GetFuncName(types.NewSlice(ttyp.Key())), g.keys.GetFuncName(typ), this)
		p.P("thatkeys := %s(%s(%s))", g.sort.GetFuncName(types.NewSlice(ttyp.Key())), g.keys.GetFuncName(typ), that)
		p.P("for i, thiskey := range thiskeys {")
		p.In()
		p.P("thatkey := thatkeys[i]")
		p.P("if thiskey == thatkey {")
		p.In()
		p.P("thisvalue := this[thiskey]")
		p.P("thatvalue := that[thatkey]")
		cmpStr, err := g.field("thisvalue", "thatvalue", ttyp.Elem())
		if err != nil {
			return err
		}
		p.P("if c := %s; c != 0 {", cmpStr)
		p.In()
		p.P(`return c`)
		p.Out()
		p.P(`}`)
		p.Out()
		p.P(`} else {`)
		p.In()
		cmpStr2, err := g.field("thiskey", "thatkey", ttyp.Key())
		if err != nil {
			return err
		}
		p.P("if c := %s; c != 0 {", cmpStr2)
		p.In()
		p.P(`return c`)
		p.Out()
		p.P(`}`)
		p.Out()
		p.P(`}`)
		p.Out()
		p.P(`}`)
		p.P(`return 0`)
		return nil
	}
	return fmt.Errorf("unsupported compare type: %s", g.TypeString(typ))
}

func wrap(value string) string {
	if strings.HasPrefix(value, "*") || strings.HasPrefix(value, "&") {
		return "(" + value + ")"
	}
	return value
}

func (g *gen) field(thisField, thatField string, fieldType types.Type) (string, error) {
	if named, isNamed := fieldType.(*types.Named); isNamed {
		inputType := compareMethodInputParam(named)
		if inputType != nil {
			ityp := *inputType
			if _, ok := ityp.(*types.Pointer); ok {
				return fmt.Sprintf("%s.Compare(&%s)", wrap(thisField), thatField), nil
			} else if _, ok := ityp.(*types.Interface); ok {
				return fmt.Sprintf("%s.Compare(&%s)", wrap(thisField), thatField), nil
			} else {
				return fmt.Sprintf("%s.Compare(%s)", wrap(thisField), thatField), nil
			}
		}
	}
	switch typ := fieldType.Underlying().(type) {
	case *types.Basic:
		if typ.Kind() == types.String {
			return fmt.Sprintf("%s.Compare(%s, %s)", g.stringsPkg(), thisField, thatField), nil
		}
		return fmt.Sprintf("%s(%s, %s)", g.GetFuncName(fieldType, fieldType), thisField, thatField), nil
	case *types.Pointer:
		ref := typ.Elem()
		if named, ok := ref.(*types.Named); ok {
			inputType := compareMethodInputParam(named)
			if inputType != nil {
				ityp := *inputType
				if _, ok := ityp.(*types.Pointer); ok {
					return fmt.Sprintf("%s.Compare(%s)", wrap(thisField), thatField), nil
				} else if _, ok := ityp.(*types.Interface); ok {
					return fmt.Sprintf("%s.Compare(%s)", wrap(thisField), thatField), nil
				} else {
					// fall through to deferencing of pointers
				}
			} else {
				return fmt.Sprintf("%s(%s, %s)", g.GetFuncName(typ, typ), thisField, thatField), nil
			}
		}
		return fmt.Sprintf("%s(%s, %s)", g.GetFuncName(typ, typ), thisField, thatField), nil
	case *types.Array, *types.Map:
		return fmt.Sprintf("%s(%s, %s)", g.GetFuncName(typ, typ), thisField, thatField), nil
	case *types.Slice:
		if b, ok := typ.Elem().(*types.Basic); ok && b.Kind() == types.Byte {
			return fmt.Sprintf("%s.Compare(%s, %s)", g.bytesPkg(), thisField, thatField), nil
		}
		return fmt.Sprintf("%s(%s, %s)", g.GetFuncName(typ, typ), thisField, thatField), nil
	case *types.Struct:
		return g.field("&"+thisField, "&"+thatField, types.NewPointer(fieldType))
	default: // *Chan, *Tuple, *Signature, *Interface, *types.Basic.Kind() == types.UntypedNil, *Struct
		return "", fmt.Errorf("unsupported field type %s", g.TypeString(fieldType))
	}
}
