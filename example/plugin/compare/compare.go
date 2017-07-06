package compare

import "sort"

type MyStruct struct {
	Int64     int64
	StringPtr *string
}

func sortStructs(ss []*MyStruct) {
	sort.Slice(ss, func(i, j int) bool {
		return deriveCompare(ss[i], ss[j]) < 0
	})
}
