// Example: customprefix shows how to defined a derived function that does not
// have to start with default "deriveEqual" prefix.
// in the Makefile we can see the goderive command being called:
//
//   goderive --equalprefix="eq" ./...
//
// This sets the new prefix to "eq".
package customprefix

type MyStruct struct {
	Int64  int64
	String string
}

func (this *MyStruct) Equal(that *MyStruct) bool {
	return eqPtrToMyStruct(this, that)
}
