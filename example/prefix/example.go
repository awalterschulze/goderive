// Example: prefix shows how to defined a derived function that does not
// have to start with default "deriveEqual" prefix.
// in the Makefile we can see the goderive command being called:
//
//	goderive --prefix="generate" ./...
//
// This sets the new prefix to "generate".
package prefix

type MyStruct struct {
	Int64     int64
	StringPtr *string
}

func (this *MyStruct) Equal(that *MyStruct) bool {
	return generateEqual(this, that)
}
