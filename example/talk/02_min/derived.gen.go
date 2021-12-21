// Code generated by goderive DO NOT EDIT.

package main

import (
	"strings"
)

// Min returns the minimum of the two input values.
func Min(a, b *Person) *Person {
	if Compare(a, b) < 0 {
		return a
	}
	return b
}

// Compare returns:
//   * 0 if this and that are equal,
//   * -1 is this is smaller and
//   * +1 is this is bigger.
func Compare(this, that *Person) int {
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
	if c := Compare_(this.age, that.age); c != 0 {
		return c
	}
	return 0
}

// Compare_ returns:
//   * 0 if this and that are equal,
//   * -1 is this is smaller and
//   * +1 is this is bigger.
func Compare_(this, that int) int {
	if this != that {
		if this < that {
			return -1
		} else {
			return 1
		}
	}
	return 0
}
