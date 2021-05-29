package test

import (
	"reflect"

	"vendortest"
)

type intAlias int

type stringAlias string

type stringKeyAlias string

type float64Alias float64

type LocalType struct{}

func true5(lt *LocalType) (*LocalType, bool) {
	return lt, true
}

type DeriveTheDerived struct {
	Field int
}

type SomeJson struct {
	Name  string
	Other KeyValue
}

type KeyValue map[string]interface{}

func (kv KeyValue) Equal(that KeyValue) bool {
	return reflect.DeepEqual(kv, that)
}

type Visitor struct {
	UserName   *string
	RemoteAddr string
}

type UseVendor struct {
	Vendors []*vendortest.AVendoredObject
}

type Adder struct {
	Int int
}
