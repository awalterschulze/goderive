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

type SomeComplexTypes struct {
	J []*RecursiveType
	K []RecursiveType
	L *RecursiveType
	M RecursiveType
	N map[int]RecursiveType
	O map[string]*RecursiveType
	P map[int64]string
}

func (this *SomeComplexTypes) Equal(that *SomeComplexTypes) bool {
	return deriveEqualPtrToSomeComplexTypes(this, that)
}

type RecursiveType struct {
	Bytes []byte
	N     map[int]RecursiveType
}

func (this *RecursiveType) Equal(that *RecursiveType) bool {
	return deriveEqualPtrToRecursiveType(this, that)
}
