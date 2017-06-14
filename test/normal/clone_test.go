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

func clone(this interface{}) interface{} {
	method := reflect.ValueOf(this).MethodByName("Clone")
	res := method.Call([]reflect.Value{})
	return res[0].Interface()
}

func TestCloneStructs(t *testing.T) {
	structs := []interface{}{
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
		// &MapWithStructs{},
		// &RecursiveType{},
		// &EmbeddedStruct1{},
		// &EmbeddedStruct2{},
		// &StructWithStructFieldWithoutEqualMethod{},
		// &StructWithStructWithFromAnotherPackage{},
		// &FieldWithStructWithPrivateFields{},
	}
	for _, this := range structs {
		desc := reflect.TypeOf(this).Elem().Name()
		t.Run(desc, func(t *testing.T) {
			for i := 0; i < 100; i++ {
				this = random(this)
				that := clone(this)
				if want, got := true, reflect.DeepEqual(this, that); want != got {
					t.Fatalf("want %v got %v\n this = %#v, that = %#v\n", want, got, this, that)
				}
			}
		})
	}
}
