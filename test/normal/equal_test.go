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

func equal(this, that interface{}) bool {
	method := reflect.ValueOf(this).MethodByName("Equal")
	res := method.Call([]reflect.Value{reflect.ValueOf(that)})
	return res[0].Interface().(bool)
}

func TestEqualStructs(t *testing.T) {
	structs := []interface{}{
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
		&UnnamedStruct{},
		&StructWithStructFieldWithoutEqualMethod{},
		&StructWithStructWithFromAnotherPackage{},
		&FieldWithStructWithPrivateFields{},
		&Enums{},
		&NamedTypes{},
		&Time{},
		&Duration{},
		&Nickname{},
		&PrivateEmbedded{},
	}
	for _, this := range structs {
		desc := reflect.TypeOf(this).Elem().Name()
		t.Run(desc, func(t *testing.T) {
			for i := 0; i < 100; i++ {
				this = random(this)
				if want, got := true, equal(this, this); want != got {
					t.Fatalf("want %v got %v\n this = %#v\n", want, got, this)
				}
				that := random(this)
				if want, got := reflect.DeepEqual(this, that), equal(this, that); want != got {
					t.Fatalf("want %v got %v\n this = %#v\n that = %#v", want, got, this, that)
				}
			}
		})
	}
}

func TestEqualInline(t *testing.T) {
	t.Run("intslices", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			this := random([]int{}).([]int)
			if want, got := true, deriveEqualSliceOfint(this, this); want != got {
				t.Fatalf("want %v got %v\n this = %#v\n", want, got, this)
			}
			that := random(this).([]int)
			if want, got := reflect.DeepEqual(this, that), deriveEqualSliceOfint(this, that); want != got {
				t.Fatalf("want %v got %v\n this = %#v\n that = %#v", want, got, this, that)
			}
		}
	})
	t.Run("mapinttoint", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			this := random(map[int]int{}).(map[int]int)
			if want, got := true, deriveEqualMapOfintToint(this, this); want != got {
				t.Fatalf("want %v got %v\n this = %#v\n", want, got, this)
			}
			that := random(this).(map[int]int)
			if want, got := reflect.DeepEqual(this, that), deriveEqualMapOfintToint(this, that); want != got {
				t.Fatalf("want %v got %v\n this = %#v\n that = %#v", want, got, this, that)
			}
		}
	})
	t.Run("intptr", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			var intptr *int
			this := random(intptr).(*int)
			if want, got := true, deriveEqualPtrToint(this, this); want != got {
				t.Fatalf("want %v got %v\n this = %#v\n", want, got, this)
			}
			that := random(this).(*int)
			if want, got := reflect.DeepEqual(this, that), deriveEqualPtrToint(this, that); want != got {
				t.Fatalf("want %v got %v\n this = %#v\n that = %#v", want, got, this, that)
			}
		}
	})
	t.Run("ptrtoslice", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			var intptr *[]int
			this := random(intptr).(*[]int)
			if want, got := true, deriveEqualPtrToSliceOfint(this, this); want != got {
				t.Fatalf("want %v got %v\n this = %#v\n", want, got, this)
			}
			that := random(this).(*[]int)
			if want, got := reflect.DeepEqual(this, that), deriveEqualPtrToSliceOfint(this, that); want != got {
				t.Fatalf("want %v got %v\n this = %#v\n that = %#v", want, got, this, that)
			}
		}
	})
	t.Run("ptrtoarray", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			var intptr *[10]int
			this := random(intptr).(*[10]int)
			if want, got := true, deriveEqualPtrToArray10Ofint(this, this); want != got {
				t.Fatalf("want %v got %v\n this = %#v\n", want, got, this)
			}
			that := random(this).(*[10]int)
			if want, got := reflect.DeepEqual(this, that), deriveEqualPtrToArray10Ofint(this, that); want != got {
				t.Fatalf("want %v got %v\n this = %#v\n that = %#v", want, got, this, that)
			}
		}
	})
	t.Run("ptrtomap", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			var intptr *map[int]int
			this := random(intptr).(*map[int]int)
			if want, got := true, deriveEqualPtrToMapOfintToint(this, this); want != got {
				t.Fatalf("want %v got %v\n this = %#v\n", want, got, this)
			}
			that := random(this).(*map[int]int)
			if want, got := reflect.DeepEqual(this, that), deriveEqualPtrToMapOfintToint(this, that); want != got {
				t.Fatalf("want %v got %v\n this = %#v\n that = %#v", want, got, this, that)
			}
		}
	})
	t.Run("structnoptr", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			var strct BuiltInTypes
			this := random(strct).(BuiltInTypes)
			if want, got := true, deriveEqual1(this, this); want != got {
				t.Fatalf("want %v got %v\n this = %#v\n", want, got, this)
			}
			that := random(strct).(BuiltInTypes)
			if want, got := reflect.DeepEqual(this, that), deriveEqual1(this, that); want != got {
				t.Fatalf("want %v got %v\n this = %#v\n that = %#v", want, got, this, that)
			}
		}
	})
}

func TestCurriedEqual(t *testing.T) {
	item := &BuiltInTypes{
		Int:     100,
		String:  "abc",
		Float64: r.Float64(),
	}
	items := random([]*BuiltInTypes{}).([]*BuiltInTypes)
	contains := deriveAnyEqualCurry(deriveEqualCurry(item), items)
	if contains {
		t.Fatalf("should not be contained")
	}
	items = append(items, item)
	eq := deriveEqualCurry(item)
	contains = deriveAnyEqualCurry(eq, items)
	if !contains {
		t.Fatalf("expected to be contained")
	}
}
