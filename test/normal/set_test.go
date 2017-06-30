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

func TestSet(t *testing.T) {
	var m map[int64]struct{}
	m = random(m).(map[int64]struct{})
	m2 := deriveSetInt64s(deriveKeysForInt64s(m))
	if len(m2) != len(m) {
		t.Fatalf("length of keys: want %d got %d", len(m), len(m2))
	}
	for key, _ := range m2 {
		if _, ok := m[key]; !ok {
			t.Fatalf("key %v does not exist in %#v", key, m)
		}
	}
	for key, _ := range m {
		if _, ok := m2[key]; !ok {
			t.Fatalf("key %v does not exist in %#v", key, m2)
		}
	}
}
