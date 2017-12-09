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
	"strconv"
	"testing"
)

func toInts(ss []string) ([]int, error) {
	return deriveTraverse(strconv.Atoi, ss)
}

func TestTraverse(t *testing.T) {
	is, err := toInts([]string{"1", "2", "3"})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(is, []int{1, 2, 3}) {
		t.Fatalf("not equal")
	}
	if _, err := toInts([]string{"1", "a", "3"}); err == nil {
		t.Fatal("expected error")
	}
}
