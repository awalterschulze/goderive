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
	"strconv"
	"strings"
)

type TypesMap interface {
	SetFuncName(typ types.Type, name string) error
	GetFuncName(typ types.Type) string
	Generating(typ types.Type)
	ToGenerate() []types.Type
	Done() bool
}

type typesMap struct {
	qual      types.Qualifier
	prefix    string
	generated map[string]bool
	funcToTyp map[string]string
	typToFunc map[string]string
	typs      []types.Type
}

func newTypesMap(qual types.Qualifier, prefix string) TypesMap {
	return &typesMap{
		qual:      qual,
		prefix:    prefix,
		generated: make(map[string]bool),
		funcToTyp: make(map[string]string),
		typToFunc: make(map[string]string),
		typs:      nil,
	}
}

func (this *typesMap) SetFuncName(typ types.Type, funcName string) error {
	typName := this.nameOf(typ)
	if fName, ok := this.typToFunc[typName]; ok {
		if fName == funcName {
			return nil
		}
		return fmt.Errorf("ambigious function names for type %s = (%s | %s)", typ, fName, funcName)
	}
	if tName, ok := this.funcToTyp[funcName]; ok {
		if tName == typName {
			return nil
		}
		return fmt.Errorf("duplicate function name %s = (%s | %s)", funcName, tName, typName)
	}
	if _, ok := this.generated[typName]; !ok {
		this.generated[typName] = false
	}
	this.typToFunc[typName] = funcName
	this.funcToTyp[funcName] = typName
	this.typs = append(this.typs, typ)
	return nil
}

func (this *typesMap) GetFuncName(typ types.Type) string {
	name := this.nameOf(typ)
	if f, ok := this.typToFunc[name]; ok {
		return f
	}
	funcName := this.funcOf(typ)
	_, exists := this.funcToTyp[funcName]
	for exists {
		funcName += "_"
		_, exists = this.funcToTyp[funcName]
	}
	this.SetFuncName(typ, funcName)
	return funcName
}

func (this *typesMap) Generating(typ types.Type) {
	name := this.nameOf(typ)
	this.generated[name] = true
}

func (this *typesMap) isGenerated(typ types.Type) bool {
	name := this.nameOf(typ)
	return this.generated[name]
}

func (this *typesMap) ToGenerate() []types.Type {
	typs := make([]types.Type, 0, len(this.typs))
	for i, typ := range this.typs {
		if !this.isGenerated(typ) {
			typs = append(typs, this.typs[i])
		}
	}
	return typs
}

func (this *typesMap) Done() bool {
	for _, typ := range this.typs {
		if !this.isGenerated(typ) {
			return false
		}
	}
	return true
}

func (this *typesMap) nameOf(typ types.Type) string {
	return typeName(typ, this.qual)
}

func (this *typesMap) funcOf(typ types.Type) string {
	return this.prefix + strings.Replace(typeName(typ, this.qual), "$", "", -1)
}

func typeName(typ types.Type, qual types.Qualifier) string {
	switch t := typ.(type) {
	case *types.Pointer:
		return "PtrTo" + typeName(t.Elem(), qual)
	case *types.Array:
		sizeStr := strconv.Itoa(int(t.Len()))
		return "Array" + sizeStr + "Of" + typeName(t.Elem(), qual)
	case *types.Slice:
		return "SliceOf" + typeName(t.Elem(), qual)
	case *types.Map:
		return "MapOf" + typeName(t.Key(), qual) + "To" + typeName(t.Elem(), qual)
	}
	// The dollar helps to make sure that typenames cannot be faked by the user.
	return "$" + types.TypeString(typ, qual)
}
