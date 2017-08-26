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
	"math"
	"testing"
)

func TestMinInt64(t *testing.T) {
	var list []int64
	list = random(list).([]int64)
	min := deriveMinInt64s(list, int64(math.MinInt64))
	sorted := deriveSortInt64s(list)
	othermin := int64(math.MinInt64)
	if len(sorted) > 0 {
		othermin = sorted[0]
	}
	if min != othermin {
		t.Fatalf("%v != %v", min, othermin)
	}
}

func TestMin2Int64(t *testing.T) {
	if m := deriveMinInt(1, 2); m != 1 {
		t.Fatalf("min should be 1, but its %d", m)
	}
	var v int
	a := random(v).(int)
	b := random(v).(int)
	if deriveMinInt(a, b) != deriveMinInt(b, a) {
		t.Fatal("min is unsemetric")
	}
}

func TestMinStruct(t *testing.T) {
	var list []*BuiltInTypes
	list = random(list).([]*BuiltInTypes)
	min := deriveMinStructs(list, nil)
	sorted := deriveSortStructs(list)
	var othermin *BuiltInTypes
	if len(sorted) > 0 {
		othermin = sorted[0]
	}
	if !min.Equal(othermin) {
		t.Fatalf("%v != %v", min, othermin)
	}
}
