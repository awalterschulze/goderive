# goderive

[![Build Status](https://travis-ci.org/awalterschulze/goderive.svg?branch=master)](https://travis-ci.org/awalterschulze/goderive)

`goderive` parses your go code for functions which are not implemented and then generates these functions for you by deriving their implementations from the parameter types. 

Deep Functions:

  - [Equal](http://godoc.org/github.com/awalterschulze/goderive/plugin/equal) `deriveEqual(T, T) bool`
  - [Compare](http://godoc.org/github.com/awalterschulze/goderive/plugin/compare) `deriveCompare(T, T) int`
  - [CopyTo](http://godoc.org/github.com/awalterschulze/goderive/plugin/copyto) `deriveCopyTo(src *T, dst *T)`

Tool Functions:

  - [Keys](http://godoc.org/github.com/awalterschulze/goderive/plugin/keys) `deriveKeys(map[K]V) []K`
  - [Sort](http://godoc.org/github.com/awalterschulze/goderive/plugin/sort) `deriveSort([]T) []T`
  - [Unique](http://godoc.org/github.com/awalterschulze/goderive/plugin/unique) `deriveUnique([]T) []T`
  - [Set](http://godoc.org/github.com/awalterschulze/goderive/plugin/set) `deriveSet([]T) map[T]struct{}`
  - [Min](http://godoc.org/github.com/awalterschulze/goderive/plugin/min) `deriveMin(list []T, default T) (min T)` or `deriveMin(T, T) T`
  - [Max](http://godoc.org/github.com/awalterschulze/goderive/plugin/max) `deriveMax(list []T, default T) (max T)` or `deriveMax(T, T) T`
  - [Contains](http://godoc.org/github.com/awalterschulze/goderive/plugin/contains) `deriveContains([]T, T) bool`
  - [Intersect](http://godoc.org/github.com/awalterschulze/goderive/plugin/intersect) `deriveIntersect(a, b []T) []T` or `deriveIntersect(a, b map[T]struct{}) map[T]struct{}`
  - [Union](http://godoc.org/github.com/awalterschulze/goderive/plugin/union) `deriveUnion(a, b []T) []T` or `deriveUnion(a, b map[T]struct{}) map[T]struct{}`

Functional Functions:

  - [Fmap](http://godoc.org/github.com/awalterschulze/goderive/plugin/fmap) `deriveFmap(f(A) B, []A) []B`
  - [Join](http://godoc.org/github.com/awalterschulze/goderive/plugin/join) `deriveJoin([][]T) []T`
  - [Filter](http://godoc.org/github.com/awalterschulze/goderive/plugin/filter) `deriveFilter(pred func(T) bool, []T) []T`
  - [TakeWhile](http://godoc.org/github.com/awalterschulze/goderive/plugin/takewhile) `deriveTakeWhile(pred func(T) bool, []T) []T`
  - [Flip](http://godoc.org/github.com/awalterschulze/goderive/plugin/flip) `deriveFlip(f func(A, B, ...) T) func(B, A, ...) T`
  - [Curry](http://godoc.org/github.com/awalterschulze/goderive/plugin/curry) `deriveCurry(f func(A, B, ...) T) func(A) func(B, ...) T`
  - [Uncurry](http://godoc.org/github.com/awalterschulze/goderive/plugin/uncurry) `deriveUncurry(f func(A) func(B, ...) T) func(A, B, ...) T`

When goderive walks over your code it is looking for a function that:
  - was not implemented (or was previously derived) and
  - has a predefined prefix.

Functions which have been previously derived will be regenerated to keep them up to date with the latest modifications to your types.  This keeps these functions, which are truly mundane to write, maintainable.  

For example when someone in your team adds a new field to a struct and forgets to update the CopyTo method.  This problem is solved by goderive, by generating generated functions given the new types.

Function prefixes are by default `deriveCamelCaseFunctionName`, for example `deriveEqual`.
These are customizable using command line flags.

Let `goderive` edit your function names in your source code, by enabling `autoname` and `dedup` using the command line flags.
These flags respectively makes sure than your functions have unique names and that you don't generate multiple functions that do the same thing.

## Examples

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
			((this.StringPtr == nil && that.StringPtr == nil) || 
        (this.StringPtr != nil && that.StringPtr != nil && *(this.StringPtr) == *(that.StringPtr)))
}
```

More Examples:

  - [Equal](https://github.com/awalterschulze/goderive/tree/master/example/plugin/equal)
  - [Compare](https://github.com/awalterschulze/goderive/tree/master/example/plugin/compare)
  - [CopyTo](https://github.com/awalterschulze/goderive/tree/master/example/plugin/copyto)
  - [Keys](https://github.com/awalterschulze/goderive/tree/master/example/plugin/keys)
  - [Sort](https://github.com/awalterschulze/goderive/tree/master/example/plugin/sort)
  - [Unique](https://github.com/awalterschulze/goderive/tree/master/example/plugin/unique)
  - [Set](https://github.com/awalterschulze/goderive/tree/master/example/plugin/set)
  - [Min](https://github.com/awalterschulze/goderive/tree/master/example/plugin/min)

## How to run

goderive can be run from the command line:

`goderive ./...`

, using the same path semantics as the go tool.

[You can also run goderive using go generate](https://github.com/awalterschulze/goderive/blob/master/example/gogenerate/example.go) 

[And you can customize function prefixes](https://github.com/awalterschulze/goderive/blob/master/example/customprefix/Makefile)

You can let goderive rename your functions using the `-autoname` and `-dedup` flags.
If these flags are not used, goderive will not touch your code and rather return an error.

## Customization

The derive package allows you to create your own code generator plugins, see all the current plugins for examples.

You can also create your own vanity binary.
Including your own generators and/or customization of function prefixes, etc.
This should be easy to figure out by looking at [main.go](https://github.com/awalterschulze/goderive/blob/master/main.go)

## Inspired By

Haskell's deriving

## Users

These projects use goderive:

  - [katydid](https://github.com/katydid/katydid)
