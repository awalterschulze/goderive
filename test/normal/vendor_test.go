package test

import (
	"testing"
)

func TestVendor(t *testing.T) {
	if !deriveEqual(&UseVendor{}, &UseVendor{}) {
		t.Fatal("not equal")
	}
}
