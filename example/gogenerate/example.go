//go:generate goderive .

package gogenerate

type MyStruct struct {
	Int64  int64
	String string
}

func (this *MyStruct) Equal(that *MyStruct) bool {
	return deriveEqualPtrToMyStruct(this, that)
}
