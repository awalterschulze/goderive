package test

import (
	"sort"
	"testing"
)

func TestSortedIntAlias(t *testing.T) {
	var aliased []intAlias
	aliased = random(aliased).([]intAlias)
	got := deriveSortedSliceIntAlias(aliased)
	if len(aliased) != len(got) {
		t.Fatalf("length of keys: want %d got %d", len(aliased), len(got))
	}
	var unaliased []int
	for _, val := range aliased {
		unaliased = append(unaliased, int(val))
	}
	if !sort.IntsAreSorted(unaliased) {
		t.Fatalf("slice are not sorted %v", got)
	}
}

func TestSortedStringAlias(t *testing.T) {
	var aliased []stringAlias
	aliased = random(aliased).([]stringAlias)
	got := deriveSortedSliceStringAlias(aliased)
	if len(aliased) != len(got) {
		t.Fatalf("length of keys: want %d got %d", len(aliased), len(got))
	}
	var unaliased []string
	for _, val := range aliased {
		unaliased = append(unaliased, string(val))
	}
	if !sort.StringsAreSorted(unaliased) {
		t.Fatalf("slice are not sorted %v", got)
	}
}

func TestSortedFloat64Alias(t *testing.T) {
	var aliased []float64Alias
	aliased = random(aliased).([]float64Alias)
	got := deriveSortedSliceFloat64Alias(aliased)
	if len(aliased) != len(got) {
		t.Fatalf("length of keys: want %d got %d", len(aliased), len(got))
	}
	var unaliased []float64
	for _, val := range aliased {
		unaliased = append(unaliased, float64(val))
	}
	if !sort.Float64sAreSorted(unaliased) {
		t.Fatalf("slice are not sorted %v", got)
	}
}
