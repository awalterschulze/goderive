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
	out := method.Call(nil)
	return out[0].Interface()
}

func TestCloneStructs(t *testing.T) {
	structs := []interface{}{
		&Empty{},
		&BuiltInTypes{},
		&PtrToBuiltInTypes{},
	}
	for _, this := range structs {
		desc := reflect.TypeOf(this).Elem().Name()
		t.Run(desc, func(t *testing.T) {
			for i := 0; i < 100; i++ {
				this = random(this)
				for reflect.ValueOf(this).IsNil() {
					this = random(this)
				}
				that := clone(this)
				if want, got := true, reflect.DeepEqual(this, that); want != got {
					t.Fatalf("want %v got %v\n this = %#v, that = %#v\n", want, got, this, that)
				}
			}
		})
	}
}

func TestCloneInline(t *testing.T) {
	t.Run("intslices", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			want := random([]int{}).([]int)
			got := deriveCloneSliceOfint(want)
			if !reflect.DeepEqual(want, got) {
				t.Fatalf("wanted %#v, but got %#v", want, got)
			}
		}
	})
	t.Run("mapinttoint", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			want := random(map[int]int{}).(map[int]int)
			got := deriveCloneMapOfintToint(want)
			if !reflect.DeepEqual(want, got) {
				t.Fatalf("wanted %#v, but got %#v", want, got)
			}
		}
	})
	t.Run("intptr", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			var intptr *int
			want := random(intptr).(*int)
			got := deriveClonePtrToint(want)
			if !reflect.DeepEqual(want, got) {
				t.Fatalf("wanted %#v, but got %#v", want, got)
			}
		}
	})
	t.Run("ptrtoslice", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			var intptr *[]int
			want := random(intptr).(*[]int)
			got := deriveClonePtrToSliceOfint(want)
			if !reflect.DeepEqual(want, got) {
				t.Fatalf("wanted %#v, but got %#v", want, got)
			}
		}
	})
	t.Run("ptrtoarray", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			var intptr *[10]int
			want := random(intptr).(*[10]int)
			got := deriveClonePtrToArray10Ofint(want)
			if !reflect.DeepEqual(want, got) {
				t.Fatalf("wanted %#v, but got %#v", want, got)
			}
		}
	})
	t.Run("ptrtomap", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			var intptr *map[int]int
			want := random(intptr).(*map[int]int)
			got := deriveClonePtrToMapOfintToint(want)
			if !reflect.DeepEqual(want, got) {
				t.Fatalf("wanted %#v, but got %#v", want, got)
			}
		}
	})
	t.Run("structnoptr", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			var strct BuiltInTypes
			want := random(strct).(BuiltInTypes)
			got := deriveClone1(want)
			if !reflect.DeepEqual(want, got) {
				t.Fatalf("wanted %#v, but got %#v", want, got)
			}
		}
	})
}
