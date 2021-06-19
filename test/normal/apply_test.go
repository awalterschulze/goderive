//  Copyright 2021 Jake Son
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

func TestApplySingle(t *testing.T) {
	applied := deriveApplyMarshal(json.Marshal, map[string]int{"value": 10})
	want := `{"value":10}`
	got, err := applied()
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestApplyMultiple(t *testing.T) {
	f := func(a int, b string, c bool) string {
		return fmt.Sprintf("%d%s%v", a, b, c)
	}
	applied := deriveApplyMultiple(f, true)
	want := `1atrue`
	got := applied(1, "a")
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

func TestApplyApplied(t *testing.T) {
	f := func(a string, b int, c bool) string {
		return fmt.Sprintf("%s%d%v", a, b, c)
	}
	applied := deriveApply3(f, true)
	want := `a1true`
	applyapplied := deriveApplyApplied(applied, 1)
	got := applyapplied("a")
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}
