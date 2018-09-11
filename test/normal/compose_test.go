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
	"fmt"
	"strconv"
	"testing"
)

func TestCompose(t *testing.T) {
	read := func() (string, error) {
		return "1", nil
	}
	parseFloat := func(i string) (float64, error) {
		return strconv.ParseFloat(i, 64)
	}
	got, err := deriveCompose(read, parseFloat)()
	if err != nil {
		t.Fatal(err)
	}
	want := float64(1)
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestComposeA(t *testing.T) {
	read := func(s string) (string, error) {
		return s, nil
	}
	parseFloat := func(i string) (float64, error) {
		return strconv.ParseFloat(i, 64)
	}
	parse := deriveComposeA(read, parseFloat)
	got, err := parse("1")
	if err != nil {
		t.Fatal(err)
	}
	want := float64(1)
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestCompose2(t *testing.T) {
	read := func(s string, z string) ([]string, string, error) {
		return []string{s, z}, s + z, nil
	}
	parseFloat := func(ss []string, s string) (float64, error) {
		if ss[0]+ss[1] != s {
			return 0, fmt.Errorf("wtf")
		}
		return strconv.ParseFloat(s, 64)
	}
	parse := deriveCompose2(read, parseFloat)
	got, err := parse("1", "2")
	if err != nil {
		t.Fatal(err)
	}
	want := float64(12)
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestComposeRetBoolSuccess(t *testing.T) {
	read := func(s string) (string, error) {
		return s, nil
	}
	lenLessThan2 := func(i string) (bool, error) {
		result := len(i) < 2
		return result, nil
	}
	check := deriveComposeRetBool(read, lenLessThan2)
	got, err := check("1")
	if err != nil {
		t.Fatal(err)
	}
	want := true
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestComposeRetBoolFailed(t *testing.T) {
	read := func(s string) (string, error) {
		if s == "" {
			return s, fmt.Errorf("empty string")
		}
		return s, nil
	}
	lenLessThan2 := func(i string) (bool, error) {
		result := len(i) < 2
		return result, nil
	}
	check := deriveComposeRetBool(read, lenLessThan2)
	got, err := check("") // passing empty string will fail
	if err == nil {
		t.Fatalf("Expected error from empty string")
	}
	want := false
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestComposeVariadic(t *testing.T) {
	read := func(s string) (string, error) {
		return s, nil
	}
	parseFloat := func(i string) (float64, error) {
		return strconv.ParseFloat(i, 64)
	}
	toInt := func(f float64) (int, error) {
		i := int(f)
		if float64(i) != f {
			return 0, fmt.Errorf("%f is not a whole number", f)
		}
		return i, nil
	}
	parse := deriveComposeVariadic(read, parseFloat, toInt)
	got, err := parse("1")
	if err != nil {
		t.Fatal(err)
	}
	want := 1
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}
