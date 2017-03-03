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

package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/loader"
)

var comparePrefix = flag.String("compare.prefix", "deriveCompare", "set the prefix for compare functions that should be derived.")

func generateCompare(p Printer, pkgInfo *loader.PackageInfo, prefix string, strict bool, calls []*ast.CallExpr) error {
	qual := types.RelativeTo(pkgInfo.Pkg)
	typesMap := newTypesMap(qual, prefix)

	for _, call := range calls {
		fn, ok := call.Fun.(*ast.Ident)
		if !ok {
			continue
		}
		if !strings.HasPrefix(fn.Name, prefix) {
			continue
		}
		if len(call.Args) != 2 {
			return fmt.Errorf("%s does not have two arguments\n", fn.Name)
		}
		t0 := pkgInfo.TypeOf(call.Args[0])
		t1 := pkgInfo.TypeOf(call.Args[1])
		if !types.Identical(t0, t1) {
			return fmt.Errorf("%s has two arguments, but they are of different types %s != %s\n",
				fn.Name, t0, t1)
		}

		if err := typesMap.SetFuncName(t0, fn.Name); err != nil {
			return err
		}
	}

	eq := newCompare(p, typesMap, qual)

	moreToGenerate := true
	for moreToGenerate {
		moreToGenerate = false
		for _, typ := range typesMap.List() {
			if typesMap.IsGenerated(typ) {
				continue
			}
			moreToGenerate = true
			if err := eq.genFuncFor(typ); err != nil {
				return err
			}
		}
	}
	return nil
}

func newCompare(printer Printer, typesMap TypesMap, qual types.Qualifier) *compare {
	return &compare{
		printer:    printer,
		typesMap:   typesMap,
		qual:       qual,
		bytesPkg:   printer.NewImport("bytes"),
		stringsPkg: printer.NewImport("strings"),
	}
}

type compare struct {
	printer    Printer
	typesMap   TypesMap
	qual       types.Qualifier
	bytesPkg   Import
	stringsPkg Import
}

func (this *compare) genFuncFor(typ types.Type) error {
	p := this.printer
	this.typesMap.Generating(typ)
	typeStr := types.TypeString(typ, this.qual)
	p.P("")
	p.P("func %s(this, that %s) int {", this.typesMap.GetFuncName(typ), typeStr)
	p.In()
	switch ttyp := typ.Underlying().(type) {
	case *types.Pointer:
		p.P("if this == nil {")
		p.In()
		p.P("if that == nil {")
		p.In()
		p.P("return 0")
		p.Out()
		p.P("}")
		p.P("return -1")
		p.Out()
		p.P("}")
		p.P("if that == nil {")
		p.In()
		p.P("return 1")
		p.Out()
		p.P("}")
		ref := ttyp.Elem()
		p.P("return %s(*this, *that)", this.typesMap.GetFuncName(ref))
	case *types.Basic:
		switch ttyp.Kind() {
		case types.String:
			p.P("return %s.Compare(this, that)", this.stringsPkg())
		case types.Complex128, types.Complex64:
			p.P("return 0 //TODO")
		case types.Bool:
			p.P("if this == that {")
			p.In()
			p.P("return 0")
			p.Out()
			p.P("}")
			p.P("if that {")
			p.In()
			p.P("return -1")
			p.Out()
			p.P("}")
			p.P("return 1")
		default:
			p.P("if this != that {")
			p.In()
			p.P("if this < that {")
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
	case *types.Struct:
		numFields := ttyp.NumFields()
		for i := 0; i < numFields; i++ {
			field := ttyp.Field(i)
			fieldType := field.Type()
			fieldName := field.Name()
			thisField := "this." + fieldName
			thatField := "that." + fieldName
			fieldStr, err := this.field(thisField, thatField, fieldType)
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
	case *types.Slice:
		p.P("if this == nil {")
		p.In()
		p.P("if that == nil {")
		p.In()
		p.P("return 0")
		p.Out()
		p.P("}")
		p.P("return -1")
		p.Out()
		p.P("}")
		p.P("if that == nil {")
		p.In()
		p.P("return 1")
		p.Out()
		p.P("}")
		p.P("if len(this) != len(that) {")
		p.In()
		p.P("if len(this) < len(that) {")
		p.In()
		p.P("return -1")
		p.Out()
		p.P("}")
		p.P("return 1")
		p.Out()
		p.P("}")
		p.P("for i := 0; i < len(this); i++ {")
		p.In()
		cmpStr, err := this.field("this[i]", "that[i]", ttyp.Elem())
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
	case *types.Array:
		p.P("if len(this) != len(that) {")
		p.In()
		p.P("if len(this) < len(that) {")
		p.In()
		p.P("return -1")
		p.Out()
		p.P("}")
		p.P("return 1")
		p.Out()
		p.P("}")
		p.P("for i := 0; i < len(this); i++ {")
		p.In()
		cmpStr, err := this.field("this[i]", "that[i]", ttyp.Elem())
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
	// case *types.Map:
	// 	p.P("if this == nil {")
	// 	p.In()
	// 	p.P("if that == nil {")
	// 	p.In()
	// 	p.P("return 0")
	// 	p.Out()
	// 	p.P("}")
	// 	p.P("return -1")
	// 	p.Out()
	// 	p.P("}")
	// 	p.P("if that == nil {")
	// 	p.In()
	// 	p.P("return 1")
	// 	p.Out()
	// 	p.P("}")
	// 	p.P("if len(this) != len(that) {")
	// 	p.In()
	// 	p.P("if len(this) < len(that) {")
	// 	p.In()
	// 	p.P("return -1")
	// 	p.Out()
	// 	p.P("}")
	// 	p.P("return 1")
	// 	p.Out()
	// 	p.P("}")
	// 	p.P("for k, v := range this {")
	// 	p.In()
	// 	p.P("thatv, ok := that[k]")
	// 	p.P("if !ok {")
	// 	p.In()
	// 	p.P("return false")
	// 	p.Out()
	// 	p.P("}")
	// 	cmpStr, err := this.field("v", "thatv", ttyp.Elem())
	// 	if err != nil {
	// 		return err
	// 	}
	// 	p.P("if %s {", cmpStr)
	// 	p.In()
	// 	p.P("return false")
	// 	p.Out()
	// 	p.P("}")
	// 	p.Out()
	// 	p.P("}")
	// 	p.P("return true")
	// TODO create sorted map keys plugin
	default:
		return fmt.Errorf("unsupported compare type: %#v", typ)
	}
	p.Out()
	p.P("}")
	return nil
}

func (this *compare) field(thisField, thatField string, fieldType types.Type) (string, error) {
	switch typ := fieldType.(type) {
	case *types.Basic:
		if typ.Kind() == types.String {
			return fmt.Sprintf("%s.Compare(%s, %s)", this.stringsPkg(), thisField, thatField), nil
		}
		return fmt.Sprintf("%s(%s, %s)", this.typesMap.GetFuncName(fieldType), thisField, thatField), nil
	case *types.Pointer:
		ref := typ.Elem()
		if _, ok := ref.(*types.Named); ok {
			return fmt.Sprintf("%s.Compare(%s)", thisField, thatField), nil
		}
		return fmt.Sprintf("%s(%s, %s)", this.typesMap.GetFuncName(typ), thisField, thatField), nil
	case *types.Array, *types.Map:
		return fmt.Sprintf("%s(%s, %s)", this.typesMap.GetFuncName(typ), thisField, thatField), nil
	case *types.Slice:
		if b, ok := typ.Elem().(*types.Basic); ok && b.Kind() == types.Byte {
			return fmt.Sprintf("%s.Compare(%s, %s)", this.bytesPkg(), thisField, thatField), nil
		}
		return fmt.Sprintf("%s(%s, %s)", this.typesMap.GetFuncName(typ), thisField, thatField), nil
	case *types.Named:
		return fmt.Sprintf("%s.Compare(&%s)", thisField, thatField), nil
	default: // *Chan, *Tuple, *Signature, *Interface, *types.Basic.Kind() == types.UntypedNil, *Struct
		return "", fmt.Errorf("unsupported field type %#v", fieldType)
	}
}
