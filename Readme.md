# goderive

goderive derives (generates) golang functions:

  - Equal
  - Compare (TODO)

## Example

The derivEqualForPtrToA function will be generated given the following code:

```go
package main

type A struct {
	B []byte
	// ... lots of other fields
}

func main() {
	a1 := &A{B: []byte("abc")}
	a2 := &A{B: []byte("cde")}
	if !derivEqualForPtrToA(a1, a2) {
		println("SUCCESS")
	}
}
```

This way only the used functions are generated, keeping generated code to a minimum.
