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
	"testing"
)

func TestContainsInt64(t *testing.T) {
	var list []int64
	list = random(list).([]int64)
	if len(list) == 0 {
		return
	}
	item := list[rand.Intn(len(list))]
	if !deriveContainsInt64s(list, item) {
		t.Fatalf("%v is not contained in %v", item, list)
	}
	s := deriveSetInt64s(list)
	delete(s, item)
	l := deriveKeysForInt64s(s)
	if deriveContainsInt64s(l, item) {
		t.Fatalf("%v is contained in %v", item, l)
	}
}

func TestContainsStruct(t *testing.T) {
	var list []*BuiltInTypes
	list = random(list).([]*BuiltInTypes)
	if len(list) == 0 {
		return
	}
	item := list[rand.Intn(len(list))]
	if !deriveContainsStruct(list, item) {
		t.Fatalf("%v is not contained in %v", item, list)
	}
	var newitem *BuiltInTypes
	newitem = random(newitem).(*BuiltInTypes)
	if deriveContainsStruct(list, newitem) {
		t.Fatalf("%v is contained in %v", newitem, list)
	}
}
