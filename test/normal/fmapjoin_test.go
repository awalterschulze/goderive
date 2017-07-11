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
	"strings"
	"testing"
)

func TestFmapJoin(t *testing.T) {
	ss := []string{"a,b", "c,d"}
	split := func(s string) []string {
		return strings.Split(s, ",")
	}
	got := deriveJoinSS(deriveFmapSS(split, ss))
	want := []string{"a", "b", "c", "d"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestFmapJoinError(t *testing.T) {
	read := func() (string, error) {
		return "1", nil
	}
	parseInt := func(i string) (int64, error) {
		ii, err := strconv.ParseInt(i, 10, 64)
		return int64(ii), err
	}
	got, err := deriveJoinEE(deriveFmapEE64(parseInt, read))
	if err != nil {
		t.Fatal(err)
	}
	want := int64(1)
	if got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}
