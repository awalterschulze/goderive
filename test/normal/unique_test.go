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

func TestUniqueInt64s(t *testing.T) {
	input := []int64{1, 2, 3, 2, 1}
	want := []int64{1, 2, 3}
	got := deriveUniqueInt64s(input)
	if len(got) != len(want) {
		t.Fatalf("got too long: %#v", got)
	}
	for _, g := range got {
		if !deriveContainsInt64s(want, g) {
			t.Fatalf("did not get %d", g)
		}
	}
}

func TestUniqueStructs(t *testing.T) {
	var b *BuiltInTypes
	b1 := random(b).(*BuiltInTypes)
	b2 := random(b).(*BuiltInTypes)
	b3 := random(b).(*BuiltInTypes)
	input := []*BuiltInTypes{b1, b2, b3, b2, b1}
	want := []*BuiltInTypes{b1, b2, b3}
	got := deriveUniqueStructs(input)
	if len(got) != len(want) {
		t.Fatalf("got too long: %#v", got)
	}
	for _, g := range got {
		if !deriveContainsStruct(want, g) {
			t.Fatalf("did not get %#v", g)
		}
	}
}

func TestUniqueStructsWithoutPointers(t *testing.T) {
	var b *PtrToBuiltInTypes
	b1 := random(b).(*PtrToBuiltInTypes)
	b2 := random(b).(*PtrToBuiltInTypes)
	b3 := random(b).(*PtrToBuiltInTypes)
	input := []PtrToBuiltInTypes{*b1, *b2, *b3, *b2, *b1}
	want := []PtrToBuiltInTypes{*b1, *b2, *b3}
	got := deriveUniqueStructsPtrs(input)
	if len(got) != len(want) {
		t.Fatalf("got too long: %#v", got)
	}
	for _, g := range got {
		if !deriveContainsStructPtr(want, g) {
			t.Fatalf("did not get %#v", g)
		}
	}
}
