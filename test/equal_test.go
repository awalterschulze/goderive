//  Copyright 2017 Walter Schulze
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package test

import (
	"testing"
)

func TestEqualBuiltInTypes(t *testing.T) {
	this := &BuiltInTypes{}
	if !this.Equal(this) {
		t.Fatal("empty not equal to itself")
	}
	this = NewRandBuiltInTypes()
	if !this.Equal(this) {
		t.Fatal("random not equal to itself")
	}
	that := NewRandBuiltInTypes()
	if this.Equal(that) {
		t.Fatalf("random %#v equal to another random %#v", this, that)
	}
}

func TestEqualPtrToBuiltInTypes(t *testing.T) {
	this := &PtrToBuiltInTypes{}
	if !this.Equal(this) {
		t.Fatal("empty not equal to itself")
	}
	this = NewRandPtrToBuiltInTypes()
	if !this.Equal(this) {
		t.Fatal("random not equal to itself")
	}
	that := NewRandPtrToBuiltInTypes()
	if this.Equal(that) {
		t.Fatalf("random %#v equal to another random %#v", this, that)
	}
}

func TestEqualSliceOfBuiltInTypes(t *testing.T) {
	this := &SliceOfBuiltInTypes{}
	if !this.Equal(this) {
		t.Fatal("empty not equal to itself")
	}
	this = NewRandSliceOfBuiltInTypes()
	if !this.Equal(this) {
		t.Fatal("random not equal to itself")
	}
	that := NewRandSliceOfBuiltInTypes()
	if this.Equal(that) {
		t.Fatalf("random %#v equal to another random %#v", this, that)
	}
}

func TestEqualSomeComplexTypes(t *testing.T) {
	this := &SomeComplexTypes{}
	if !this.Equal(this) {
		t.Fatal("empty not equal to itself")
	}
	this = NewRandSomeComplexTypes()
	if !this.Equal(this) {
		t.Fatal("random not equal to itself")
	}
	that := NewRandSomeComplexTypes()
	if this.Equal(that) {
		t.Fatalf("random %#v equal to another random %#v", this, that)
	}
}

func TestEqualRecursiveType(t *testing.T) {
	this := &RecursiveType{}
	if !this.Equal(this) {
		t.Fatal("empty not equal to itself")
	}
	this = NewRandRecursiveType()
	if !this.Equal(this) {
		t.Fatal("random not equal to itself")
	}
	that := NewRandRecursiveType()
	if this.Equal(that) {
		t.Fatalf("random %#v equal to another random %#v", this, that)
	}
}
