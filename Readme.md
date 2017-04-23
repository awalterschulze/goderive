# goderive

[![Build Status](https://travis-ci.org/awalterschulze/goderive.svg?branch=master)](https://travis-ci.org/awalterschulze/goderive)

goderive parses your go code for functions which are not implemented and then generates these functions for you by deriving their implementations from the parameter types. Functions that are currently supported include:

  - [Equal](https://github.com/awalterschulze/goderive#equal)
  - [Compare](https://github.com/awalterschulze/goderive#compare)

More functions are in the works:

  - Keys
  - Sorted
  - Fmap
  - Join

Functions which have been previously derived will be regenerated to keep them up to date with the latest modifications to your types.  This keeps these functions, which are truly mundane to write, maintainable.

Distinguishing between which function (Equal, Compare, ...) should be derived is done using a customizable prefix, see command line flags.

## Equal

The `deriveEqual` function is a faster alternative to `reflect.DeepEqual`.

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
  - Function
  - Struct without an Equal method
  - Unnamed Structs, which are not comparable with `==`

## Compare

The `deriveCompare` function is a maintainable way to implement Less functions.

### Example

In the following code the `deriveCompare` function will be spotted as a function that was not implemented (or was previously derived) and has a prefix `deriveCompare`.

```go
package main

type MyStruct struct {
	Int64  int64
	String string
}

func (this *MyStruct) Less(that *MyStruct) bool {
	return deriveCompare(this, that) < 0
}
```

### Unsupported Types

  - Chan
  - Interface
  - Function
  - Struct without a Compare method
  - Unnamed Structs

## Keys

The `deriveKeys` function returns a map's keys as a slice.

## Sorted

The `deriveSorted` function is useful for deterministically ranging over maps when used with `deriveKeys`.
This feature requires Go 1.8

### TODO

  - more tests
  - update readme example
  - add example to example package

## Fmap

The `deriveFmap` function applies a given function to each element of a list, returning a list of results in the same order.

### Example

In the following code the `deriveFmap` function will be spotted as a function that was not implemented (or was previously derived) and has a prefix `deriveFmap`.

```go
func main() {
	list := []int{1, 2, 3}
	list = deriveFmap(func(i int) int { return i+1 }, list)
	for _, e := range list {
		fmt.Printf("%d", list)
	}
	// print 234
}
```

goderive will then generate the following code in a `derived.gen.go` file in the same package:

```go
func deriveFmap(f func(int) int, list []int) []int {
	out := make([]int, len(list))
	for i, elem := range list {
		out[i] = f(elem)
	}
	return out
}
```

### TODO

  - currently only slices are supported, think about supporting other types and not just slices
  - think about functions without a return type

## Join

The `deriveJoin` function applies a given joins a slice of slices into a single slice.

### Example

In the following code the `deriveJoin` function will be spotted as a function that was not implemented (or was previously derived) and has a prefix `deriveJoin`.

```go
func main() {
	ss := []string{"a,b", "c,d"}
	split := func(s string) []string {
		return strings.Split(s, ",")
	}
	list := deriveJoin(deriveFmap(split, ss))
	for _, e := range list {
		fmt.Printf("%s", list)
	}
	// print abcd
}
```

goderive will then generate the following code in a `derived.gen.go` file in the same package:

```go
func deriveJoin(list [][]string) []string {
	if list == nil {
		return nil
	}
	res := []string{}
	for _, elem := range list {
		res = append(res, elem...)
	}
	return res
}
```

### TODO

  - currently only slices are supported, think about supporting other types and not just slices
  - what about []string and not just [][]string as in the current example.
