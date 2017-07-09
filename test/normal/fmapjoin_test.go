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
	"errors"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestFmap(t *testing.T) {
	got := deriveFmap(func(i int) int { return i + 1 }, []int{1, 2})
	want := []int{2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestFmapString(t *testing.T) {
	got := deriveFmapString(func(r rune) bool { return r == 'a' }, "abc")
	want := []bool{true, false, false}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestJoin(t *testing.T) {
	got := deriveJoin([][]int{{1, 2}, {3, 4}})
	want := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestJoinString(t *testing.T) {
	got := deriveJoinString([]string{"abc", "cde"})
	want := "abccde"
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

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

func TestFmapError(t *testing.T) {
	f := func() (int, error) {
		return 1, nil
	}
	add := func(i int) int64 {
		return int64(i + 1)
	}
	got, err := deriveFmapError(add, f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := int64(2)
	if got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

func TestFmapErrorError(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		num := func() (string, error) {
			return "1", nil
		}
		gotf, err := deriveFmapEE(strconv.Atoi, num)
		if err != nil {
			t.Fatal(err)
		}
		got, err := gotf()
		if err != nil {
			t.Fatal(err)
		}
		want := 1
		if got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
	})
	t.Run("first error", func(t *testing.T) {
		num := func() (string, error) {
			return "", errors.New("hey")
		}
		gotf, err := deriveFmapEE(strconv.Atoi, num)
		if err == nil {
			t.Fatal("expected error")
		}
		if gotf != nil {
			t.Fatal("expected nil func")
		}
	})
	t.Run("second error", func(t *testing.T) {
		num := func() (string, error) {
			return "a", nil
		}
		gotf, err := deriveFmapEE(strconv.Atoi, num)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := gotf(); err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestFmapZeroError(t *testing.T) {
	num := func() (string, error) {
		return "1", nil
	}
	got := ""
	print := func(s string) {
		got = s
	}
	err := deriveFmapPrint(print, num)
	if err != nil {
		t.Fatal(err)
	}
	want := "1"
	if got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestFmapMoreParamsError(t *testing.T) {
	num := func() (string, error) {
		return "1", nil
	}
	conv := func(s string) (int, string, error) {
		i, err := strconv.Atoi(s)
		return i, s, err
	}
	gotf, err := deriveFmapMore(conv, num)
	if err != nil {
		t.Fatal(err)
	}
	goti, got, err := gotf()
	if err != nil {
		t.Fatal(err)
	}
	want := "1"
	wanti := 1
	if got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if goti != wanti {
		t.Fatalf("got %d, want %d", got, want)
	}
}
