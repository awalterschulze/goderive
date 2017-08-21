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

type Named struct {
	Fields  []*Field
	Reflect bool
}

type Field struct {
	name     string
	external bool
	Type     types.Type
	typeStr  func() string
}

func (f *Field) Name(recv string, unsafePkg Import) string {
	if !f.Private() || !f.external {
		return recv + "." + f.name
	}
	return `*(*` + f.typeStr() + `)(` + unsafePkg() + `.Pointer(` + recv + `.FieldByName("` + f.name + `").UnsafeAddr()))`
}

func (f *Field) DebugName() string {
	return f.name
}

func (f *Field) Private() bool {
	return strings.ToLower(f.name[0:1]) == f.name[0:1]
}

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
