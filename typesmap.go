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
	"go/types"
)

type TypesMap interface {
	Set(typ types.Type, value bool)
	Get(typ types.Type) (value bool)
	List() []types.Type
}

type typesMap struct {
	qual types.Qualifier
	m    map[string]bool
	typs []types.Type
}

func newTypesMap(qual types.Qualifier) TypesMap {
	return &typesMap{qual, make(map[string]bool), nil}
}

func (this *typesMap) Set(typ types.Type, value bool) {
	name := typeName(typ, this.qual)
	if _, ok := this.m[name]; !ok {
		this.typs = append(this.typs, typ)
	}
	this.m[name] = value
}

func (this *typesMap) Get(typ types.Type) (value bool) {
	name := typeName(typ, this.qual)
	return this.m[name]
}

func (this *typesMap) List() []types.Type {
	return this.typs
}
