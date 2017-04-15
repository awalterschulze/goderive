package test

import (
	"reflect"
	"testing"
)

func TestFmap(t *testing.T) {
	got := deriveFmap(func(i int) int { return i + 1 }, []int{1, 2})
	want := []int{2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
