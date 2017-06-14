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

import "github.com/awalterschulze/goderive/test/extra"

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

func (this *BuiltInTypes) Compare(that *BuiltInTypes) int {
	return deriveComparePtrToBuiltInTypes(this, that)
}

func (this *BuiltInTypes) Clone() *BuiltInTypes {
	return deriveClonePtrToBuiltInTypes(this)
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

func (this *PtrToBuiltInTypes) Compare(that *PtrToBuiltInTypes) int {
	return deriveComparePtrToPtrToBuiltInTypes(this, that)
}

func (this *PtrToBuiltInTypes) Clone() *PtrToBuiltInTypes {
	return deriveClonePtrToPtrToBuiltInTypes(this)
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

func (this *SliceOfBuiltInTypes) Compare(that *SliceOfBuiltInTypes) int {
	return deriveComparePtrToSliceOfBuiltInTypes(this, that)
}

func (this *SliceOfBuiltInTypes) Clone() *SliceOfBuiltInTypes {
	return deriveClonePtrToSliceOfBuiltInTypes(this)
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

func (this *SliceOfPtrToBuiltInTypes) Compare(that *SliceOfPtrToBuiltInTypes) int {
	return deriveComparePtrToSliceOfPtrToBuiltInTypes(this, that)
}

func (this *SliceOfPtrToBuiltInTypes) Clone() *SliceOfPtrToBuiltInTypes {
	return deriveClonePtrToSliceOfPtrToBuiltInTypes(this)
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

	AnotherBoolOfDifferentSize [10]bool
}

func (this *ArrayOfBuiltInTypes) Equal(that *ArrayOfBuiltInTypes) bool {
	return deriveEqualPtrToArrayOfBuiltInTypes(this, that)
}

func (this *ArrayOfBuiltInTypes) Compare(that *ArrayOfBuiltInTypes) int {
	return deriveComparePtrToArrayOfBuiltInTypes(this, that)
}

func (this *ArrayOfBuiltInTypes) Clone() *ArrayOfBuiltInTypes {
	return deriveClonePtrToArrayOfBuiltInTypes(this)
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

	AnotherBoolOfDifferentSize [10]*bool
}

func (this *ArrayOfPtrToBuiltInTypes) Equal(that *ArrayOfPtrToBuiltInTypes) bool {
	return deriveEqualPtrToArrayOfPtrToBuiltInTypes(this, that)
}

func (this *ArrayOfPtrToBuiltInTypes) Compare(that *ArrayOfPtrToBuiltInTypes) int {
	return deriveComparePtrToArrayOfPtrToBuiltInTypes(this, that)
}

func (this *ArrayOfPtrToBuiltInTypes) Clone() *ArrayOfPtrToBuiltInTypes {
	return deriveClonePtrToArrayOfPtrToBuiltInTypes(this)
}

type MapsOfSimplerBuiltInTypes struct {
	StringToUint32 map[string]uint32
	Uint64ToInt64  map[uint8]int64
}

func (this *MapsOfSimplerBuiltInTypes) Equal(that *MapsOfSimplerBuiltInTypes) bool {
	return deriveEqualPtrToMapsOfSimplerBuiltInTypes(this, that)
}

func (this *MapsOfSimplerBuiltInTypes) Compare(that *MapsOfSimplerBuiltInTypes) int {
	return deriveComparePtrToMapsOfSimplerBuiltInTypes(this, that)
}

func (this *MapsOfSimplerBuiltInTypes) Clone() *MapsOfSimplerBuiltInTypes {
	return deriveClonePtrToMapsOfSimplerBuiltInTypes(this)
}

type MapsOfBuiltInTypes struct {
	BoolToString          map[bool]string
	StringToBool          map[string]bool
	Complex128ToComplex64 map[complex128]complex64
	Float64ToUint32       map[float64]uint32
	Uint16ToUint8         map[uint16]uint8
}

func (this *MapsOfBuiltInTypes) Equal(that *MapsOfBuiltInTypes) bool {
	return deriveEqualPtrToMapsOfBuiltInTypes(this, that)
}

func (this *MapsOfBuiltInTypes) Compare(that *MapsOfBuiltInTypes) int {
	return deriveComparePtrToMapsOfBuiltInTypes(this, that)
}

func (this *MapsOfBuiltInTypes) Clone() *MapsOfBuiltInTypes {
	return deriveClonePtrToMapsOfBuiltInTypes(this)
}

type SliceToSlice struct {
	Ints    [][]int
	Strings [][]string
	IntPtrs [][]*int
}

func (this *SliceToSlice) Equal(that *SliceToSlice) bool {
	return deriveEqualPtrToSliceToSlice(this, that)
}

func (this *SliceToSlice) Compare(that *SliceToSlice) int {
	return deriveComparePtrToSliceToSlice(this, that)
}

func (this *SliceToSlice) Clone() *SliceToSlice {
	return deriveClonePtrToSliceToSlice(this)
}

type PtrTo struct {
	Basic *int
	Slice *[]int
	Array *[4]int
	Map   *map[int]int
}

func (this *PtrTo) Equal(that *PtrTo) bool {
	return deriveEqualPtrToPtrTo(this, that)
}

func (this *PtrTo) Compare(that *PtrTo) int {
	return deriveComparePtrToPtrTo(this, that)
}

func (this *PtrTo) Clone() *PtrTo {
	return deriveClonePtrToPtrTo(this)
}

type Name struct {
	Name string
}

func (this *Name) Equal(that *Name) bool {
	return deriveEqualPtrToName(this, that)
}

func (this *Name) Compare(that *Name) int {
	return deriveComparePtrToName(this, that)
}

type Structs struct {
	Struct             Name
	PtrToStruct        *Name
	SliceOfStructs     []Name
	SliceToPtrOfStruct []*Name
}

func (this *Structs) Equal(that *Structs) bool {
	return deriveEqualPtrToStructs(this, that)
}

func (this *Structs) Compare(that *Structs) int {
	return deriveComparePtrToStructs(this, that)
}

func (this *Structs) Clone() *Structs {
	return deriveClonePtrToStructs(this)
}

type MapWithStructs struct {
	NameToString             map[Name]string
	StringToName             map[string]Name
	StringToPtrToName        map[string]*Name
	StringToSliceOfName      map[string][]Name
	StringToSliceOfPtrToName map[string][]*Name
}

func (this *MapWithStructs) Equal(that *MapWithStructs) bool {
	return deriveEqualPtrToMapWithStructs(this, that)
}

func (this *MapWithStructs) Compare(that *MapWithStructs) int {
	return deriveComparePtrToMapWithStructs(this, that)
}

func (this *MapWithStructs) Clone() *MapWithStructs {
	return deriveClonePtrToMapWithStructs(this)
}

type RecursiveType struct {
	Bytes []byte
	N     map[int]RecursiveType
}

func (this *RecursiveType) Equal(that *RecursiveType) bool {
	return deriveEqualPtrToRecursiveType(this, that)
}

func (this *RecursiveType) Compare(that *RecursiveType) int {
	return deriveComparePtrToRecursiveType(this, that)
}

type EmbeddedStruct1 struct {
	Name
	*Structs
}

func (this *EmbeddedStruct1) Equal(that *EmbeddedStruct1) bool {
	return deriveEqualPtrToEmbeddedStruct1(this, that)
}

func (this *EmbeddedStruct1) Compare(that *EmbeddedStruct1) int {
	return deriveComparePtrToEmbeddedStruct1(this, that)
}

type EmbeddedStruct2 struct {
	Structs
	*Name
}

func (this *EmbeddedStruct2) Equal(that *EmbeddedStruct2) bool {
	return deriveEqualPtrToEmbeddedStruct2(this, that)
}

func (this *EmbeddedStruct2) Compare(that *EmbeddedStruct2) int {
	return deriveComparePtrToEmbeddedStruct2(this, that)
}

type UnnamedStruct struct {
	Unnamed struct {
		String string
	}
}

func (this *UnnamedStruct) Equal(that *UnnamedStruct) bool {
	return deriveEqualPtrToUnnamedStruct(this, that)
}

type StructWithStructFieldWithoutEqualMethod struct {
	A *StructWithoutEqualMethod
	B StructWithoutEqualMethod
}

func (this *StructWithStructFieldWithoutEqualMethod) Equal(that *StructWithStructFieldWithoutEqualMethod) bool {
	return deriveEqualPtrToStructWithStructFieldWithoutEqualMethod(this, that)
}

func (this *StructWithStructFieldWithoutEqualMethod) Compare(that *StructWithStructFieldWithoutEqualMethod) int {
	return deriveComparePtrToStructWithStructFieldWithoutEqualMethod(this, that)
}

type StructWithoutEqualMethod struct {
	Num int64
}

type StructWithStructWithFromAnotherPackage struct {
	A *extra.StructWithoutEqualMethod
	B extra.StructWithoutEqualMethod
}

func (this *StructWithStructWithFromAnotherPackage) Equal(that *StructWithStructWithFromAnotherPackage) bool {
	return deriveEqualPtrToStructWithStructWithFromAnotherPackage(this, that)
}

func (this *StructWithStructWithFromAnotherPackage) Compare(that *StructWithStructWithFromAnotherPackage) int {
	return deriveComparePtrToStructWithStructWithFromAnotherPackage(this, that)
}

type FieldWithStructWithPrivateFields struct {
	A *extra.PrivateFieldAndNoEqualMethod
}

func (this *FieldWithStructWithPrivateFields) Equal(that *FieldWithStructWithPrivateFields) bool {
	return deriveEqualPtrToFieldWithStructWithPrivateFields(this, that)
}

func (this *FieldWithStructWithPrivateFields) Compare(that *FieldWithStructWithPrivateFields) int {
	return deriveComparePtrToFieldWithStructWithPrivateFields(this, that)
}
