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
	"encoding/json"
	"fmt"
	"testing"
)

func TestFlip2(t *testing.T) {
	flipped := deriveFlipMarshal(json.Unmarshal)
	got := ""
	want := `string`
	if err := flipped(&got, []byte(`"`+want+`"`)); err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestFlip3(t *testing.T) {
	f := func(a int, b string, c bool) string {
		return fmt.Sprintf("%d%s%v", a, b, c)
	}
	flipped := deriveFlip3(f)
	want := `1atrue`
	got := flipped("a", 1, true)
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestFlipBlankIdentifier(t *testing.T) {
	f := func(a string, _ bool, c int) string {
		return fmt.Sprintf("%s%v%d", a, true, c)
	}
	flipped := deriveFlipBlankIdentifier(f)
	want := `atrue1`
	got := flipped(false, "a", 1)
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}
