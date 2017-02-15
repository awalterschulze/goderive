# goderive

[![Build Status](https://travis-ci.org/awalterschulze/goderive.svg?branch=master)](https://travis-ci.org/awalterschulze/goderive)

goderive parses your go code and generates functions that are derived from the given types. Functions that are currently supported include:

  - Equal
  - Compare (TODO)

## Example

The `deriveEqualForPtrToA` function will be spotted as a function that needs to generated, because it has a prefix `deriveEqualFor` and specifies the type `PtrToA` in the suffix.

```go
package main

type A struct {
	B []byte
	// ... lots of other fields
}

func main() {
	a1 := &A{B: []byte("abc")}
	a2 := &A{B: []byte("cde")}
	if !deriveEqualForPtrToA(a1, a2) {
		println("SUCCESS")
	}
}
```

This way only the used functions are generated, keeping generated code to a minimum.
