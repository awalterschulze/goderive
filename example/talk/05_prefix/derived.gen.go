// Code generated by goderive DO NOT EDIT.

package main

import (
	"sort"
)

// deriveSort sorts the slice inplace and also returns it.
//
// Deprecated: In favour of generics.
func deriveSort(list []string) []string {
	sort.Strings(list)
	return list
}

// deriveKeys returns the keys of the input map as a slice.
//
// Deprecated: In favour of generics.
func deriveKeys(m map[string]Person) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}
