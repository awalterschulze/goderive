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
	"testing"
)

func TestEqualA(t *testing.T) {
	a := &A{}
	if !a.Equal(a) {
		t.Fatal("empty A not equal to itself")
	}
	a = NewRandA()
	if !a.Equal(a) {
		t.Fatal("random A not equal to itself")
	}
	a2 := NewRandA()
	if a.Equal(a2) {
		t.Fatalf("random a %#v equal to another random a %#v", a, a2)
	}
}

func TestEqualB(t *testing.T) {
	b := &B{}
	if !b.Equal(b) {
		t.Fatal("empty B not equal to itself")
	}
	b = NewRandB()
	if !b.Equal(b) {
		t.Fatal("random B not equal to itself")
	}
	b2 := NewRandB()
	if b2.Equal(b) {
		t.Fatalf("random B %#v equal to another random B %#v", b, b2)
	}
}
