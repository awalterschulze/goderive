# goderive

[![Build Status](https://travis-ci.org/awalterschulze/goderive.svg?branch=master)](https://travis-ci.org/awalterschulze/goderive)

`goderive` parses your go code for functions which are not implemented and then generates these functions for you by deriving their implementations from the parameter types. 

Deep Functions:

  - [Equal](http://godoc.org/github.com/awalterschulze/goderive/plugin/equal)
  - [Compare](http://godoc.org/github.com/awalterschulze/goderive/plugin/compare)
  - [CopyTo](http://godoc.org/github.com/awalterschulze/goderive/plugin/copyto)

Tool Functions:

  - [Keys](http://godoc.org/github.com/awalterschulze/goderive/plugin/keys)
  - [Sort](http://godoc.org/github.com/awalterschulze/goderive/plugin/sort)

Functional Functions:

  - [Fmap](http://godoc.org/github.com/awalterschulze/goderive/plugin/fmap)
  - [Join](http://godoc.org/github.com/awalterschulze/goderive/plugin/join)

When goderive walks over your code it is looking for a function that:
  - was not implemented (or was previously derived) and
  - has a prefix predefined prefix.

Functions which have been previously derived will be regenerated to keep them up to date with the latest modifications to your types.  This keeps these functions, which are truly mundane to write, maintainable.  

For example when someone in your team adds a new field to a struct and forgets to update the CopyTo method.  This problem is solved by goderive, by generating generated functions given the new types.

Function prefixes are by default `deriveCamelCaseFunctionName`, for example `deriveEqual`.
These are customizable using command line flags.

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
			((this.StringPtr == nil && that.StringPtr == nil) || 
        (this.StringPtr != nil && that.StringPtr != nil && *(this.StringPtr) == *(that.StringPtr)))
}
```
