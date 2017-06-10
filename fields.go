package main

import (
	"go/types"
	"strings"
)

type Named struct {
	Fields  []*Field
	Reflect bool
}

type Field struct {
	name    string
	Type    types.Type
	typeStr string
}

func (f *Field) Name(recv string, unsafePkg Import) string {
	if !f.Private() {
		return recv + "." + f.name
	}
	return `*(*` + f.typeStr + `)(` + unsafePkg() + `.Pointer(` + recv + `.FieldByName("` + f.name + `").UnsafeAddr()))`
}

func (f *Field) Private() bool {
	return strings.ToLower(f.name[0:1]) == f.name[0:1]
}

func Fields(typesMap TypesMap, typ *types.Struct) *Named {
	numFields := typ.NumFields()
	n := &Named{
		Fields: make([]*Field, numFields),
	}
	for i := 0; i < numFields; i++ {
		field := typ.Field(i)
		fieldType := field.Type()
		fieldName := field.Name()
		n.Fields[i] = &Field{
			name:    fieldName,
			Type:    fieldType,
			typeStr: typesMap.TypeString(fieldType),
		}
		if n.Fields[i].Private() {
			n.Reflect = true
		}
	}
	return n
}
