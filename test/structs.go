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
	"math/rand"
	"reflect"
	"testing/quick"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

type BuiltInTypes struct {
	Bool       bool
	Byte       byte
	Complex128 complex128
	Complex64  complex64
	//Error error
	Float64 float64
	Float32 float32
	Int     int
	Int16   int16
	Int32   int32
	Int64   int64
	Int8    int8
	Rune    rune
	String  string
	Uint    uint
	Uint16  uint16
	Uint32  uint32
	Uint64  uint64
	Uint8   uint8
	UintPtr uintptr
}

func (this *BuiltInTypes) Equal(that *BuiltInTypes) bool {
	return deriveEqualPtrToBuiltInTypes(this, that)
}

var typeOfBuiltInTypes = reflect.TypeOf(new(BuiltInTypes))

func NewRandBuiltInTypes() *BuiltInTypes {
	v, ok := quick.Value(typeOfBuiltInTypes, r)
	if !ok {
		panic("unable to generate value")
	}
	return v.Interface().(*BuiltInTypes)
}

type PtrToBuiltInTypes struct {
	Bool       *bool
	Byte       *byte
	Complex128 *complex128
	Complex64  *complex64
	//Error error
	Float64 *float64
	Float32 *float32
	Int     *int
	Int16   *int16
	Int32   *int32
	Int64   *int64
	Int8    *int8
	Rune    *rune
	String  *string
	Uint    *uint
	Uint16  *uint16
	Uint32  *uint32
	Uint64  *uint64
	Uint8   *uint8
	UintPtr *uintptr
}

func (this *PtrToBuiltInTypes) Equal(that *PtrToBuiltInTypes) bool {
	return deriveEqualPtrToPtrToBuiltInTypes(this, that)
}

var typeOfPtrToBuiltInTypes = reflect.TypeOf(new(PtrToBuiltInTypes))

func NewRandPtrToBuiltInTypes() *PtrToBuiltInTypes {
	v, ok := quick.Value(typeOfPtrToBuiltInTypes, r)
	if !ok {
		panic("unable to generate value")
	}
	return v.Interface().(*PtrToBuiltInTypes)
}

type SliceOfBuiltInTypes struct {
	Bool       []bool
	Byte       []byte
	Complex128 []complex128
	Complex64  []complex64
	//Error error
	Float64 []float64
	Float32 []float32
	Int     []int
	Int16   []int16
	Int32   []int32
	Int64   []int64
	Int8    []int8
	Rune    []rune
	String  []string
	Uint    []uint
	Uint16  []uint16
	Uint32  []uint32
	Uint64  []uint64
	Uint8   []uint8
	UintPtr []uintptr
}

func (this *SliceOfBuiltInTypes) Equal(that *SliceOfBuiltInTypes) bool {
	return deriveEqualPtrToSliceOfBuiltInTypes(this, that)
}

var typeOfSliceOfBuiltInTypes = reflect.TypeOf(new(SliceOfBuiltInTypes))

func NewRandSliceOfBuiltInTypes() *SliceOfBuiltInTypes {
	v, ok := quick.Value(typeOfSliceOfBuiltInTypes, r)
	if !ok {
		panic("unable to generate value")
	}
	return v.Interface().(*SliceOfBuiltInTypes)
}

type SomeComplexTypes struct {
	J []*RecursiveType
	K []RecursiveType
	L *RecursiveType
	M RecursiveType
	N map[int]RecursiveType
	O map[string]*RecursiveType
	P map[int64]string
}

var typeOfSomeComplexTypes = reflect.TypeOf(new(SomeComplexTypes))

func NewRandSomeComplexTypes() *SomeComplexTypes {
	v, ok := quick.Value(typeOfSomeComplexTypes, r)
	if !ok {
		panic("unable to generate value")
	}
	return v.Interface().(*SomeComplexTypes)
}

func (this *SomeComplexTypes) Equal(that *SomeComplexTypes) bool {
	return deriveEqualPtrToSomeComplexTypes(this, that)
}

type RecursiveType struct {
	Bytes []byte
	N     map[int]RecursiveType
}

var typeOfRecursiveType = reflect.TypeOf(new(RecursiveType))

func NewRandRecursiveType() *RecursiveType {
	v, ok := quick.Value(typeOfRecursiveType, r)
	if !ok {
		panic("unable to generate value")
	}
	return v.Interface().(*RecursiveType)
}

func (this *RecursiveType) Equal(that *RecursiveType) bool {
	return deriveEqualPtrToRecursiveType(this, that)
}
