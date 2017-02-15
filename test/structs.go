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
	Float64    float64
	Float32    float32
	Int        int
	Int16      int16
	Int32      int32
	Int64      int64
	Int8       int8
	Rune       rune
	String     string
	Uint       uint
	Uint16     uint16
	Uint32     uint32
	Uint64     uint64
	Uint8      uint8
	UintPtr    uintptr
}

func (this *BuiltInTypes) Equal(that *BuiltInTypes) bool {
	return deriveEqualPtrToBuiltInTypes(this, that)
}

type PtrToBuiltInTypes struct {
	Bool       *bool
	Byte       *byte
	Complex128 *complex128
	Complex64  *complex64
	Float64    *float64
	Float32    *float32
	Int        *int
	Int16      *int16
	Int32      *int32
	Int64      *int64
	Int8       *int8
	Rune       *rune
	String     *string
	Uint       *uint
	Uint16     *uint16
	Uint32     *uint32
	Uint64     *uint64
	Uint8      *uint8
	UintPtr    *uintptr
}

func (this *PtrToBuiltInTypes) Equal(that *PtrToBuiltInTypes) bool {
	return deriveEqualPtrToPtrToBuiltInTypes(this, that)
}

type SliceOfBuiltInTypes struct {
	Bool       []bool
	Byte       []byte
	Complex128 []complex128
	Complex64  []complex64
	Float64    []float64
	Float32    []float32
	Int        []int
	Int16      []int16
	Int32      []int32
	Int64      []int64
	Int8       []int8
	Rune       []rune
	String     []string
	Uint       []uint
	Uint16     []uint16
	Uint32     []uint32
	Uint64     []uint64
	Uint8      []uint8
	UintPtr    []uintptr
}

func (this *SliceOfBuiltInTypes) Equal(that *SliceOfBuiltInTypes) bool {
	return deriveEqualPtrToSliceOfBuiltInTypes(this, that)
}

type SliceOfPtrToBuiltInTypes struct {
	Bool       []*bool
	Byte       []*byte
	Complex128 []*complex128
	Complex64  []*complex64
	Float64    []*float64
	Float32    []*float32
	Int        []*int
	Int16      []*int16
	Int32      []*int32
	Int64      []*int64
	Int8       []*int8
	Rune       []*rune
	String     []*string
	Uint       []*uint
	Uint16     []*uint16
	Uint32     []*uint32
	Uint64     []*uint64
	Uint8      []*uint8
	UintPtr    []*uintptr
}

func (this *SliceOfPtrToBuiltInTypes) Equal(that *SliceOfPtrToBuiltInTypes) bool {
	return deriveEqualPtrToSliceOfPtrToBuiltInTypes(this, that)
}

type ArrayOfBuiltInTypes struct {
	Bool       [1]bool
	Byte       [2]byte
	Complex128 [3]complex128
	Complex64  [4]complex64
	Float64    [5]float64
	Float32    [6]float32
	Int        [7]int
	Int16      [8]int16
	Int32      [9]int32
	Int64      [10]int64
	Int8       [11]int8
	Rune       [12]rune
	String     [13]string
	Uint       [14]uint
	Uint16     [15]uint16
	Uint32     [16]uint32
	Uint64     [17]uint64
	Uint8      [18]uint8
	UintPtr    [19]uintptr
}

func (this *ArrayOfBuiltInTypes) Equal(that *ArrayOfBuiltInTypes) bool {
	return deriveEqualPtrToArrayOfBuiltInTypes(this, that)
}

type ArrayOfPtrToBuiltInTypes struct {
	Bool       [1]*bool
	Byte       [2]*byte
	Complex128 [3]*complex128
	Complex64  [4]*complex64
	Float64    [5]*float64
	Float32    [6]*float32
	Int        [7]*int
	Int16      [8]*int16
	Int32      [9]*int32
	Int64      [10]*int64
	Int8       [11]*int8
	Rune       [12]*rune
	String     [13]*string
	Uint       [14]*uint
	Uint16     [15]*uint16
	Uint32     [16]*uint32
	Uint64     [17]*uint64
	Uint8      [18]*uint8
	UintPtr    [19]*uintptr
}

func (this *ArrayOfPtrToBuiltInTypes) Equal(that *ArrayOfPtrToBuiltInTypes) bool {
	return deriveEqualPtrToArrayOfPtrToBuiltInTypes(this, that)
}

type SliceToSlice struct {
	Ints    [][]int
	Strings [][]string
	IntPtrs [][]*int
}

func (this *SliceToSlice) Equal(that *SliceToSlice) bool {
	return deriveEqualPtrToSliceToSlice(this, that)
}

type Name struct {
	Name string
}

func (this *Name) Equal(that *Name) bool {
	return deriveEqualPtrToName(this, that)
}

type SomeComplexTypes struct {
	J []*Name
	K []Name
	L *Name
	M Name
	N map[int]Name
	O map[string]*Name
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
