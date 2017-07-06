package equal

type MyStruct struct {
	Int64     int64
	StringPtr *string
}

func (this *MyStruct) Equal(that *MyStruct) bool {
	return deriveEqual(this, that)
}
