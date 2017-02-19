//go:generate goderive .

// Example: gogenerate shows us how to call goderive using gogenerate instead of using a Makefile.
package gogenerate

type MyStruct struct {
	Int64  int64
	String string
}

func (this *MyStruct) Equal(that *MyStruct) bool {
	return deriveEqual(this, that)
}
