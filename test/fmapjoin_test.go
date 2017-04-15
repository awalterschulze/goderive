package test

import (
	"reflect"
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

func TestJoin(t *testing.T) {
	got := deriveJoin([][]int{{1, 2}, {3, 4}})
	want := []int{1, 2, 3, 4}
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
