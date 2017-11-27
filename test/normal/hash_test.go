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

type hashable interface {
	Hash() uint64
}

func TestHashStructs(t *testing.T) {
	structs := []hashable{
		&Empty{},
		&BuiltInTypes{},
		&PtrToBuiltInTypes{},
		&SliceOfBuiltInTypes{},
		&SliceOfPtrToBuiltInTypes{},
		&ArrayOfBuiltInTypes{},
		&ArrayOfPtrToBuiltInTypes{},
		&MapsOfSimplerBuiltInTypes{},
		&MapsOfBuiltInTypes{},
		&SliceToSlice{},
		&PtrTo{},
		&Structs{},
		&MapWithStructs{},
		&RecursiveType{},
		&EmbeddedStruct1{},
		&EmbeddedStruct2{},
		// &UnnamedStruct{},
		&StructWithStructFieldWithoutEqualMethod{},
		&StructWithStructWithFromAnotherPackage{},
		// &FieldWithStructWithPrivateFields{},
		&Enums{},
		&NamedTypes{},
		// &Time{},
		&Duration{},
		&Nickname{},
		&PrivateEmbedded{},
	}
	for _, this := range structs {
		desc := reflect.TypeOf(this).Elem().Name()
		t.Run(desc, func(t *testing.T) {
			for i := 0; i < 100; i++ {
				this = random(this).(hashable)
				want := this.Hash()
				got := this.Hash()
				if want != got {
					t.Fatalf("want %v got %v\n this = %#v\n", want, got, this)
				}
			}
		})
	}
}

func TestHashInline(t *testing.T) {
	t.Run("intslices", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			this := random([]int{}).([]int)
			if want, got := deriveHashSliceOfint(this), deriveHashSliceOfint(this); want != got {
				t.Fatalf("want %v got %v\n this = %#v\n", want, got, this)
			}
		}
	})
	t.Run("mapinttoint", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			this := random(map[int]int{}).(map[int]int)
			if want, got := deriveHashMapOfintToint(this), deriveHashMapOfintToint(this); want != got {
				t.Fatalf("want %v got %v\n this = %#v\n", want, got, this)
			}
		}
	})
	t.Run("intptr", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			var intptr *int
			this := random(intptr).(*int)
			if want, got := deriveHashPtrToint(this), deriveHashPtrToint(this); want != got {
				t.Fatalf("want %v got %v\n this = %#v\n", want, got, this)
			}
		}
	})
	t.Run("ptrtoslice", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			var intptr *[]int
			this := random(intptr).(*[]int)
			if want, got := deriveHashPtrToSliceOfint(this), deriveHashPtrToSliceOfint(this); want != got {
				t.Fatalf("want %v got %v\n this = %#v\n", want, got, this)
			}
		}
	})
	t.Run("ptrtoarray", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			var intptr *[10]int
			this := random(intptr).(*[10]int)
			if want, got := deriveHashPtrToArray10Ofint(this), deriveHashPtrToArray10Ofint(this); want != got {
				t.Fatalf("want %v got %v\n this = %#v\n", want, got, this)
			}
		}
	})
	t.Run("ptrtomap", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			var intptr *map[int]int
			this := random(intptr).(*map[int]int)
			if want, got := deriveHashPtrToMapOfintToint(this), deriveHashPtrToMapOfintToint(this); want != got {
				t.Fatalf("want %v got %v\n this = %#v\n", want, got, this)
			}
		}
	})
	t.Run("structnoptr", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			var strct BuiltInTypes
			this := random(strct).(BuiltInTypes)
			if want, got := deriveHash1(this), deriveHash1(this); want != got {
				t.Fatalf("want %v got %v\n this = %#v\n", want, got, this)
			}
		}
	})
}
