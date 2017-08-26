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

func TestIntersectMap(t *testing.T) {
	var m map[int64]struct{}
	m1 := random(m).(map[int64]struct{})
	m2 := random(m).(map[int64]struct{})
	intersection := deriveIntersectSetOfInt64s(m1, m2)
	if len(intersection) > len(m1) {
		t.Fatalf("length of intersection is bigger than the first original set: %d > %d", len(intersection), len(m1))
	}
	if len(intersection) > len(m2) {
		t.Fatalf("length of intersection is bigger than the second original set: %d > %d", len(intersection), len(m2))
	}
	for key := range intersection {
		if _, ok := m1[key]; !ok {
			t.Fatalf("key %v does not exist in set 1 %#v", key, m1)
		}
		if _, ok := m2[key]; !ok {
			t.Fatalf("key %v does not exist in set 2 %#v", key, m2)
		}
	}
}

func TestIntersectSlice(t *testing.T) {
	var m []int64
	m1 := random(m).([]int64)
	m2 := random(m).([]int64)
	intersection := deriveIntersectOfInt64s(m1, m2)
	if len(intersection) > len(m1) {
		t.Fatalf("length of intersection is bigger than the first original set: %d > %d", len(intersection), len(m1))
	}
	if len(intersection) > len(m2) {
		t.Fatalf("length of intersection is bigger than the second original set: %d > %d", len(intersection), len(m2))
	}
	for _, key := range intersection {
		if !deriveContainsInt64s(m1, key) {
			t.Fatalf("key %v does not exist in set 1 %#v", key, m1)
		}
		if !deriveContainsInt64s(m2, key) {
			t.Fatalf("key %v does not exist in set 2 %#v", key, m2)
		}
	}
}
