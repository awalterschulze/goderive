package sort

type MyStruct struct {
	Int64     int64
	StringPtr *string
}

func equivalent(this, that []*MyStruct) bool {
	deriveSort(this)
	deriveSort(that)
	return deriveCompare(this, that) == 0
}
