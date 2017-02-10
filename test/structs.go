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
	"math/rand"
	"reflect"
	"testing/quick"
	"time"
)

type A struct {
	B string
	C *string
	D int64
	E *int64
	I []bool
	J []*B
	K []B
	L *B
	M B
	N map[int]B
	O map[string]*B
	P map[int64]string
}

func (this *A) Equal(that *A) bool {
	return derivEqualPtrToA(this, that)
}

type B struct {
	Bytes []byte
	N     map[int]B
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

var typeOfB = reflect.TypeOf(new(B))

func NewRandB() *B {
	v, ok := quick.Value(typeOfB, r)
	if !ok {
		panic("unable to generate value")
	}
	return v.Interface().(*B)
}

func (this *B) Equal(that *B) bool {
	return derivEqualPtrToB(this, that)
}

type DontGenerateEqualMethodForMe struct {
	A string
}

var typeOfA = reflect.TypeOf(new(A))

func NewRandA() *A {
	v, ok := quick.Value(typeOfA, r)
	if !ok {
		panic("unable to generate value")
	}
	return v.Interface().(*A)
}
