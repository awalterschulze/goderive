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

import "testing"
import "reflect"

type DeriveTheDerived struct {
	Field int
}

func inefficientEqual(this, that *DeriveTheDerived) bool {
	return deriveEqualInefficientDeriveTheDerived(
		deriveCompareDeriveTheDerived(this, that),
		-1*deriveCompareDeriveTheDerived(this, that),
	)
}

func TestInefficientEqual(t *testing.T) {
	this := random(&DeriveTheDerived{}).(*DeriveTheDerived)
	if !inefficientEqual(this, this) {
		t.Fatal("not equal")
	}
}

func TestFmapKeys(t *testing.T) {
	m := map[int]string{1: "a", 2: "b", 3: "c"}
	got := deriveFmapForKeys(func(i int) string { return m[i] }, deriveSortedInts(deriveKeysForFmap(m)))
	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v != want %#v", got, want)
	}
}
