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
	"sort"
	"testing"
)

func TestSortedMapKeysStrings(t *testing.T) {
	var m map[string]string
	m = random(m).(map[string]string)
	keys := deriveSortedStrings(deriveKeysForMapStringToString(m))
	if len(keys) != len(m) {
		t.Fatalf("length of keys: want %d got %d", len(m), len(keys))
	}
	for _, key := range keys {
		if _, ok := m[key]; !ok {
			t.Fatalf("key %v does not exist in %#v", key, m)
		}
	}
	if !sort.StringsAreSorted(keys) {
		t.Fatalf("keys are not sorted %v", keys)
	}
}

func TestSortedMapKeysTypeAlias(t *testing.T) {
	type aliased map[string]string
	var unaliased map[string]string
	m := aliased(random(unaliased).(map[string]string))
	keys := deriveSortedStrings(deriveKeysForMapStringToString(m))
	if len(keys) != len(m) {
		t.Fatalf("length of keys: want %d got %d", len(m), len(keys))
	}
	for _, key := range keys {
		if _, ok := m[key]; !ok {
			t.Fatalf("key %v does not exist in %#v", key, m)
		}
	}
	if !sort.StringsAreSorted(keys) {
		t.Fatalf("keys are not sorted %v", keys)
	}
}

func TestSortedMapKeysInt(t *testing.T) {
	var m map[int]int64
	m = random(m).(map[int]int64)
	keys := deriveSortedInts(deriveKeysForMapIntToInt64(m))
	if len(keys) != len(m) {
		t.Fatalf("length of keys: want %d got %d", len(m), len(keys))
	}
	for _, key := range keys {
		if _, ok := m[key]; !ok {
			t.Fatalf("key %v does not exist in %#v", key, m)
		}
	}
	if !sort.IntsAreSorted(keys) {
		t.Fatalf("keys are not sorted %v", keys)
	}
}

func TestSortedMapKeysInt64s(t *testing.T) {
	var m map[int64]int64
	m = random(m).(map[int64]int64)
	keys := deriveSortInt64s(deriveKeysForMapInt64ToInt64(m))
	if len(keys) != len(m) {
		t.Fatalf("length of keys: want %d got %d", len(m), len(keys))
	}
	for _, key := range keys {
		if _, ok := m[key]; !ok {
			t.Fatalf("key %v does not exist in %#v", key, m)
		}
	}
	want := []int64{}
	for k := range m {
		want = append(want, k)
	}
	sort.Slice(want, func(i, j int) bool { return want[i] < want[j] })
	if !reflect.DeepEqual(keys, want) {
		t.Fatalf("want %v got %d", want, keys)
	}
}
