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
	"fmt"
	"go/types"
	"os"
)

const eqFuncPrefix = "deriveEqual"

func equalFuncName(typ types.Type, qual types.Qualifier) string {
	return eqFuncPrefix + typeName(typ, qual)
}

func genEqual(p Printer, m TypesMap, qual types.Qualifier, typ types.Type) {
	m.Set(typ, true)
	typeStr := types.TypeString(typ, qual)
	p.P("")
	p.P("func %s(this, that %s) bool {", equalFuncName(typ, qual), typeStr)
	p.In()
	switch ttyp := typ.(type) {
	case *types.Pointer:
		ref := ttyp.Elem()
		switch tttyp := ref.Underlying().(type) {
		case *types.Struct:
			numFields := tttyp.NumFields()
			if numFields == 0 {
				p.P("return (this == nil && that == nil) || (this != nil) && (that != nil)")
			} else {
				p.P("return (this == nil && that == nil) || (this != nil) && (that != nil) &&")
			}
			p.In()
			for i := 0; i < numFields; i++ {
				field := tttyp.Field(i)
				fieldType := field.Type()
				fieldName := field.Name()
				this := "this." + fieldName
				that := "that." + fieldName
				fieldStr, err := equalField(m, qual, this, that, fieldType)
				if err != nil {
					fmt.Fprintf(os.Stderr, err.Error())
					return
				}
				if (i + 1) != numFields {
					p.P(fieldStr + " &&")
				} else {
					p.P(fieldStr)
				}
			}
			p.Out()
		default:
			fmt.Fprintf(os.Stderr, "unsupported: pointer is not a named struct, but %#v\n", ref)
			return
		}
	case *types.Slice:
		p.P("if this == nil || that == nil {")
		p.In()
		p.P("return this == nil && that == nil")
		p.Out()
		p.P("}")
		p.P("if len(this) != len(that) {")
		p.In()
		p.P("return false")
		p.Out()
		p.P("}")
		p.P("for i := 0; i < len(this); i++ {")
		p.In()
		eqStr, err := equalField(m, qual, "this[i]", "that[i]", ttyp.Elem())
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			return
		}
		p.P("if !(%s) {", eqStr)
		p.In()
		p.P("return false")
		p.Out()
		p.P("}")
		p.Out()
		p.P("}")
		p.P("return true")
	// TODO case *types.Array:

	case *types.Map:
		p.P("if this == nil || that == nil {")
		p.In()
		p.P("return this == nil && that == nil")
		p.Out()
		p.P("}")
		p.P("if len(this) != len(that) {")
		p.In()
		p.P("return false")
		p.Out()
		p.P("}")
		p.P("for k, v := range this {")
		p.In()
		p.P("thatv, ok := that[k]")
		p.P("if !ok {")
		p.In()
		p.P("return false")
		p.Out()
		p.P("}")
		eqStr, err := equalField(m, qual, "v", "thatv", ttyp.Elem())
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			return
		}
		p.P("if !(%s) {", eqStr)
		p.In()
		p.P("return false")
		p.Out()
		p.P("}")
		p.Out()
		p.P("}")
		p.P("return true")
	default:
		fmt.Fprintf(os.Stderr, "unsupported type: %#v", typ)
		return
	}
	p.Out()
	p.P("}")
}

func equalField(m TypesMap, qual types.Qualifier, this, that string, fieldType types.Type) (string, error) {
	if isComparable(fieldType) {
		return fmt.Sprintf("%s == %s", this, that), nil
	}
	switch typ := fieldType.(type) {
	case *types.Pointer:
		ref := typ.Elem()
		if _, ok := ref.(*types.Named); ok {
			return fmt.Sprintf("%s.Equal(%s)", this, that), nil
		}
		eqStr, err := equalField(m, qual, "*"+this, "*"+that, ref)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("((%[1]s == nil && %[2]s == nil) || (%[1]s != nil && %[2]s != nil && %[3]s))", this, that, eqStr), nil
	// TODO case *types.Array:
	// 	p.newFunc(typ)
	// 	return fmt.Sprintf("%s(%s, %s)", equalFuncName(typ, p.qual), this, that), nil
	case *types.Slice:
		if b, ok := typ.Elem().(*types.Basic); ok && b.Kind() == types.Byte {
			return fmt.Sprintf("bytes.Equal(%s, %s)", this, that), nil
		}
		m.Set(typ, false)
		return fmt.Sprintf("%s(%s, %s)", equalFuncName(typ, qual), this, that), nil
	case *types.Map:
		m.Set(typ, false)
		return fmt.Sprintf("%s(%s, %s)", equalFuncName(typ, qual), this, that), nil
	case *types.Named:
		return fmt.Sprintf("%s.Equal(&%s)", this, that), nil
	default: // *Chan, *Tuple, *Signature, *Interface, *types.Basic.Kind() == types.UntypedNil, *Struct
		return "", fmt.Errorf("unsupported type %#v", fieldType)
	}
}
