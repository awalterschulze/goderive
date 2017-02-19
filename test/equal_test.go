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
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func equal(this, that interface{}) bool {
	eqMethod := reflect.ValueOf(this).MethodByName("Equal")
	res := eqMethod.Call([]reflect.Value{reflect.ValueOf(that)})
	return res[0].Interface().(bool)
}

func random(this interface{}) interface{} {
	v, ok := quick.Value(reflect.TypeOf(this), r)
	if !ok {
		panic(fmt.Sprintf("unable to generate value for type: %T", this))
	}
	return v.Interface()
}

func TestEqualStructs(t *testing.T) {
	structs := []interface{}{
		&BuiltInTypes{},
		&PtrToBuiltInTypes{},
		&SliceOfBuiltInTypes{},
		&SliceOfPtrToBuiltInTypes{},
		&ArrayOfBuiltInTypes{},
		&ArrayOfPtrToBuiltInTypes{},
		&MapsOfBuiltInTypes{},
		&SliceToSlice{},
		&Structs{},
		&MapWithStructs{},
		&RecursiveType{},
		&EmbeddedStruct1{},
		&EmbeddedStruct2{},
		&UnnamedStruct{},
	}
	for _, this := range structs {
		desc := reflect.TypeOf(this).Elem().Name()
		t.Run(desc, func(t *testing.T) {
			properRandom := false
			for i := 0; i < 100; i++ {
				this = random(this)
				if want, got := true, equal(this, this); want != got {
					t.Fatalf("want %v got %v\n this = %#v\n", want, got, this)
				}
				that := random(this)
				if want, got := reflect.DeepEqual(this, that), equal(this, that); want != got {
					t.Fatalf("want %v got %v\n this = %#v\n that = %#v", want, got, this, that)
				}
				if !properRandom && !equal(this, that) {
					properRandom = true
				}
			}
			if !properRandom {
				t.Fatal("random is not working")
			}
		})
	}
}

func TestEqualInline(t *testing.T) {
	t.Run("intslices", func(t *testing.T) {
		this := random([]int{}).([]int)
		for i := 0; i < 100; i++ {
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
		this := random(map[int]int{}).(map[int]int)
		for i := 0; i < 100; i++ {
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
		var intptr *int
		this := random(intptr).(*int)
		for i := 0; i < 100; i++ {
			if want, got := true, deriveEqualPtrToint(this, this); want != got {
				t.Fatalf("want %v got %v\n this = %#v\n", want, got, this)
			}
			that := random(this).(*int)
			if want, got := reflect.DeepEqual(this, that), deriveEqualPtrToint(this, that); want != got {
				t.Fatalf("want %v got %v\n this = %#v\n that = %#v", want, got, this, that)
			}
		}
	})
}
