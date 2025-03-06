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

package test

import (
	"reflect"
	"testing"
)

func deepcopy(this, that interface{}) {
	method := reflect.ValueOf(this).MethodByName("DeepCopy")
	method.Call([]reflect.Value{reflect.ValueOf(that)})
}

func TestDeepCopyStructs(t *testing.T) {
	structs := []interface{}{
		&Empty{},
		&BuiltInTypes{},
		&PtrToBuiltInTypes{},
		&SliceOfBuiltInTypes{},
		&SliceOfPtrToBuiltInTypes{},
		&ArrayOfBuiltInTypes{},
		&ArrayOfPtrToBuiltInTypes{},
		&MapsOfBuiltInTypes{},
		&MapsOfSimplerBuiltInTypes{},
		&SliceToSlice{},
		&PtrTo{},
		&Structs{},
		&MapWithStructs{},
		&RecursiveType{},
		&EmbeddedStruct1{},
		&EmbeddedStruct2{},
		&UnnamedStruct{},
		&StructWithStructFieldWithoutEqualMethod{},
		&StructWithStructWithFromAnotherPackage{},
		&FieldWithStructWithPrivateFields{},
		&Enums{},
		&NamedTypes{},
		&Duration{},
		&Nickname{},
		&PrivateEmbedded{},
		&StructOfStructs{},
	}
	for _, this := range structs {
		desc := reflect.TypeOf(this).Elem().Name()
		t.Run(desc, func(t *testing.T) {
			for i := 0; i < 100; i++ {
				this = random(this)
				for reflect.ValueOf(this).IsNil() {
					this = random(this)
				}
				that := random(this)
				for reflect.ValueOf(that).IsNil() {
					that = random(that)
				}
				deepcopy(this, that)
				if want, got := true, reflect.DeepEqual(this, that); want != got {
					t.Fatalf("want %v got %v\n this = %#v, that = %#v\n", want, got, this, that)
				}
			}
		})
	}
}

func TestDeepCopyMapNilEntry(t *testing.T) {
	this := &MapWithStructs{StringToPtrToName: map[string]*Name{
		"a": nil,
	}}
	that := &MapWithStructs{}
	deepcopy(this, that)
	if want, got := true, reflect.DeepEqual(this, that); want != got {
		t.Fatalf("want %v got %v\n this = %#v, that = %#v\n", want, got, this, that)
	}
}

func TestDeepCopyCustomMap(t *testing.T) {
	this := customMap{
		"a": "A",
	}
	that := customMap{}
	deepcopy(this, that)
	if that["a"] != "copy of A" {
		t.Fatalf("expected use of customMap.DeepCopy")
	}
}

type customMap map[string]string

func (c customMap) DeepCopy(to customMap) {
	for k, v := range c {
		to[k] = "copy of " + v
	}
}

type SimpleStruct struct {
	Level int
}

func (this *SimpleStruct) DeepCopy(that *SimpleStruct) {
	deriveDeepCopySimpleStruct(that, this)
}

type aliasToStruct = SimpleStruct

func TestDeepCopyAlias(t *testing.T) {
	this := &aliasToStruct{
		Level: 99,
	}

	that := &aliasToStruct{}

	deepcopy(this, that)

	if that.Level != this.Level {
		t.Error("expected level to use copied value")
	}
}
