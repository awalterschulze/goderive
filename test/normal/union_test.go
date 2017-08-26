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

func TestUnionMap(t *testing.T) {
	var m map[int64]struct{}
	m1 := random(m).(map[int64]struct{})
	m2 := random(m).(map[int64]struct{})
	union := deriveUnionSetOfInt64s(m1, m2)
	for key := range m1 {
		if _, ok := union[key]; !ok {
			t.Fatalf("key %v does not exist in union %#v", key, union)
		}
	}
	for key := range m2 {
		if _, ok := union[key]; !ok {
			t.Fatalf("key %v does not exist in union %#v", key, union)
		}
	}
}

func TestUnionSlice(t *testing.T) {
	var m []int64
	m1 := random(m).([]int64)
	m2 := random(m).([]int64)
	union := deriveUnionOfInt64s(m1, m2)
	for _, key := range m1 {
		if !deriveContainsInt64s(m1, key) {
			t.Fatalf("key %v does not exist in union %#v", key, union)
		}
	}
	for _, key := range m2 {
		if !deriveContainsInt64s(m2, key) {
			t.Fatalf("key %v does not exist in union %#v", key, union)
		}
	}
}
