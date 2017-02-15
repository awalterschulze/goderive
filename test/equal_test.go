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

func TestEqual(t *testing.T) {
	structs := []interface{}{
		&BuiltInTypes{},
		&PtrToBuiltInTypes{},
		&SliceOfBuiltInTypes{},
		&SliceOfPtrToBuiltInTypes{},
		&ArrayOfBuiltInTypes{},
		&ArrayOfPtrToBuiltInTypes{},

		&SliceToSlice{},
		&SomeComplexTypes{},
		&RecursiveType{},
	}
	for _, this := range structs {
		desc := reflect.TypeOf(this).Elem().Name()
		t.Run(desc, func(t *testing.T) {
			for i := 0; i < 100; i++ {
				if !equal(this, this) {
					t.Fatal("empty not equal to itself")
				}
				this = random(this)
				if !equal(this, this) {
					t.Fatal("random not equal to itself")
				}
				that := random(this)
				for reflect.ValueOf(that).IsNil() {
					that = random(this)
				}
				if equal(this, that) {
					t.Fatalf("random %#v equal to another random %#v", this, that)
				}
			}
		})
	}
}
