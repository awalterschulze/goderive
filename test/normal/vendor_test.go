package test

import (
	"testing"
	"vendortest"
)

type UseVendor struct {
	Vendors []*vendortest.AVendoredObject
}

func TestVendor(t *testing.T) {
	if !deriveEqual(&UseVendor{}, &UseVendor{}) {
		t.Fatal("not equal")
	}
}
