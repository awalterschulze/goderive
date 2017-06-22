# goderive

[![Build Status](https://travis-ci.org/awalterschulze/goderive.svg?branch=master)](https://travis-ci.org/awalterschulze/goderive)

`goderive` parses your go code for functions which are not implemented and then generates these functions for you by deriving their implementations from the parameter types. 

Functions that are currently supported include:

  - [Equal](http://godoc.org/github.com/awalterschulze/goderive/plugin/equal)
  - [Compare](http://godoc.org/github.com/awalterschulze/goderive/plugin/compare)
  - [Keys](https://github.com/awalterschulze/goderive#keys)
  - [Sorted](https://github.com/awalterschulze/goderive#sorted)
  - [Sort](https://github.com/awalterschulze/goderive#sort)

More functions are in the works:

  - Fmap
  - Join

Functions which have been previously derived will be regenerated to keep them up to date with the latest modifications to your types.  This keeps these functions, which are truly mundane to write, maintainable.

Distinguishing between which function (`Equal`, `Compare`, ...) should be derived is done using a customizable prefix, see command line flags.

Let `goderive` edit your function names in your source code, by enabling `autoname` and `dedup` using the command line flags.
These flags respectively makes sure than your functions have unique names and that you don't generate multiple functions that do the same thing.

## Example

In the following code the `deriveEqual` function will be spotted as a function that was not implemented (or was previously derived) and has a prefix `deriveEqual`.

```go
package main

type MyStruct struct {
	Int64     int64
	StringPtr *string
}

func (this *MyStruct) Equal(that *MyStruct) bool {
	return deriveEqual(this, that)
}
```

goderive will then generate the following code in a `derived.gen.go` file in the same package:

```go
func deriveEqual(this, that *MyStruct) bool {
	return (this == nil && that == nil) ||
		this != nil && that != nil &&
			this.Int64 == that.Int64 &&
			((this.StringPtr == nil && that.StringPtr == nil) || (this.StringPtr != nil && that.StringPtr != nil && *(this.StringPtr) == *(that.StringPtr)))
}
```

## Keys

The `deriveKeys` function returns a map's keys as a slice.

## Sort

This feature requires Go 1.8

The `deriveSort` function is useful for deterministically ranging over maps when used with `deriveKeys`.

`deriveSort` supports only the types that `deriveCompare` supports, since it uses it for sorting.

## Sorted

This feature requires Go 1.8

`deriveSorted` is like `deriveSort`, but sacrifices efficiency for immutability by creating a copy of its input.

## Fmap

The `deriveFmap` function applies a given function to each element of a list, returning a list of results in the same order.

TODO:
  - currently only slices are supported, think about supporting other types and not just slices
  - think about functions without a return type

## Join

The `deriveJoin` function applies a given joins a slice of slices into a single slice.

TODO:
  - currently only slices are supported, think about supporting other types and not just slices
  - what about []string and not just [][]string as in the current example.
