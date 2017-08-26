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
//  limitations under the License

package derive

import (
	"go/types"
	"strings"
)

// Named describes a named struct with a list of fields.
type Named struct {
	Fields  []*Field
	Reflect bool
}

// Field describes a struct field.
type Field struct {
	name     string
	external bool
	Type     types.Type
	typeStr  func() string
}

// Name returns the field name, given the receiver and the unsafe import, if needed.
func (f *Field) Name(recv string, unsafePkg Import) string {
	if !f.Private() || !f.external {
		return recv + "." + f.name
	}
	return `*(*` + f.typeStr() + `)(` + unsafePkg() + `.Pointer(` + recv + `.FieldByName("` + f.name + `").UnsafeAddr()))`
}

// DebugName simply returns the field name, without the receiver or any unsafe magic.
func (f *Field) DebugName() string {
	return f.name
}

// Private whether the field is private
func (f *Field) Private() bool {
	return strings.ToLower(f.name[0:1]) == f.name[0:1]
}

// Fields returns a new Named object containing a list of Fields for a given input struct.
func Fields(typesMap TypesMap, typ *types.Struct, external bool) *Named {
	numFields := typ.NumFields()
	n := &Named{
		Fields: make([]*Field, numFields),
	}
	for i := 0; i < numFields; i++ {
		field := typ.Field(i)
		fieldType := field.Type()
		fieldName := field.Name()
		n.Fields[i] = &Field{
			name:     fieldName,
			external: external,
			Type:     fieldType,
			typeStr: func() string {
				return typesMap.TypeString(fieldType)
			},
		}
		if n.Fields[i].Private() {
			if external {
				n.Reflect = true
			}
		}
	}
	return n
}
