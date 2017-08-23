package gostring

type MyStruct struct {
	Int64     int64
	StringPtr *string
}

func (this *MyStruct) GoString() string {
	return deriveGoString(this)
}
