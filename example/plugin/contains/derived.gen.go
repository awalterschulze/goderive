// Code generated by goderive DO NOT EDIT.

package contains

// deriveContains returns whether the item is contained in the list.
//
// Deprecated: In favour of generics.
func deriveContains(list []boat, item boat) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}
