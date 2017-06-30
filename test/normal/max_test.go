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

func TestMaxInt64(t *testing.T) {
	var list []int64
	list = random(list).([]int64)
	max := deriveMaxInt64s(list, int64(math.MaxInt64))
	sorted := deriveSortInt64s(list)
	othermax := int64(math.MinInt64)
	if len(sorted) > 0 {
		othermax = sorted[len(sorted)-1]
	}
	if max != othermax {
		t.Fatalf("%v != %v", max, othermax)
	}
}

func TestMaxStruct(t *testing.T) {
	var list []*BuiltInTypes
	list = random(list).([]*BuiltInTypes)
	max := deriveMaxStructs(list, nil)
	sorted := deriveSortStructs(list)
	var othermax *BuiltInTypes = nil
	if len(sorted) > 0 {
		othermax = sorted[len(sorted)-1]
	}
	if !max.Equal(othermax) {
		t.Fatalf("%v != %v", max, othermax)
	}
}
