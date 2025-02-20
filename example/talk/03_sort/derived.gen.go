// Code generated by goderive DO NOT EDIT.

package main

import (
	"sort"
	"strings"
)

// Sort sorts the slice inplace and also returns it.
//
// Deprecated: In favour of generics.
func Sort(list []Person) []Person {
	sort.Slice(list, func(i, j int) bool { return Compare(list[i], list[j]) < 0 })
	return list
}

// Compare returns:
//   - 0 if this and that are equal,
//   - -1 is this is smaller and
//   - +1 is this is bigger.
func Compare(this, that Person) int {
	return Compare_(&this, &that)
}

// Compare_ returns:
//   - 0 if this and that are equal,
//   - -1 is this is smaller and
//   - +1 is this is bigger.
func Compare_(this, that *Person) int {
	if this == nil {
		if that == nil {
			return 0
		}
		return -1
	}
	if that == nil {
		return 1
	}
	if c := strings.Compare(this.name, that.name); c != 0 {
		return c
	}
	if c := Compare_i(this.age, that.age); c != 0 {
		return c
	}
	return 0
}

// Compare_i returns:
//   - 0 if this and that are equal,
//   - -1 is this is smaller and
//   - +1 is this is bigger.
func Compare_i(this, that int) int {
	if this != that {
		if this < that {
			return -1
		} else {
			return 1
		}
	}
	return 0
}
