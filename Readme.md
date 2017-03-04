# goderive

[![Build Status](https://travis-ci.org/awalterschulze/goderive.svg?branch=master)](https://travis-ci.org/awalterschulze/goderive)

goderive parses your go code for functions which are not implemented and then generates these functions for you by deriving their implementations from the parameter types. Functions that are currently supported include:

  - Equal
  - SortedMapKeys (requires go1.8)
  - Compare (TODO)

Functions which have been previously derived will be regenerated to keep them up to date with the latest modifications to your types.  This keeps these functions, which are truly mundane to write, maintainable.

Distinguishing between which function (Equal, Compare, ...) should be derived is done using a customizable prefix, see command line flags.

## Equal

The `derivEqual` function is a faster alternative to `reflect.DeepEqual`.

### Example

In the following code the `deriveEqual` function will be spotted as a function that was not implemented (or was previously derived) and has a prefix `deriveEqual`.

```go
package main

type MyStruct struct {
	Int64  int64
	String string
}

func (this *MyStruct) Equal(that *MyStruct) bool {
	return deriveEqual(this, that)
}
```

goderive will then generate the following code in a `derived.gen.go` file in the same package:

```go
func deriveEqual(this, that *MyStruct) bool {
	return (this == nil && that == nil) || (this != nil) && (that != nil) &&
		this.Int64 == that.Int64 &&
		this.String == that.String
}
```

### Unsupported Types

  - Chan
  - Interface
  - Unnamed Structs, which are not comparable with `==`

## SortedMapKeys (Alpha)

The `deriveSortedKeys` function is useful for deterministically ranging over maps.

### Example

In the following code the `deriveSortedKeys` function will be spotted as a function that was not implemented (or was previously derived) and has a prefix `deriveSortedKeys`.

```go
func main() {
	m := map[int]string{
		1: "a",
		3: "c",
		2: "b",
	}
	for k, v := range deriveSortedKeys(m) {
		fmt.Printf("%d", k)
	}
	// print 123
}
```

goderive will then generate the following code in a `derived.gen.go` file in the same package:

```go
func deriveSortedKeys(m map[int]int) []int {
	var keys []int
	for key, _ := range m {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}
```

### TODO

  - complex64, complex128
  - structs
  - bools
  - optimize strings and ints
  - more tests
  - update readme example
  - add example to example package

