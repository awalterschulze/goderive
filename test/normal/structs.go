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
	"time"

	"github.com/awalterschulze/goderive/test/extra"
	"github.com/awalterschulze/goderive/test/nickname"
)

type Empty struct{}

func (this *Empty) Equal(that *Empty) bool {
	return deriveEqualPtrToEmpty(this, that)
}

func (this *Empty) Compare(that *Empty) int {
	return deriveComparePtrToEmpty(this, that)
}

func (this *Empty) DeepCopy(that *Empty) {
	deriveDeepCopyPtrToEmpty(that, this)
}

func (this *Empty) GoString() string {
	return deriveGoStringEmpty(this)
}

func (this *Empty) Clone() *Empty {
	return deriveCloneEmpty(this)
}

func (this *Empty) Hash() uint64 {
	return deriveHashEmpty(this)
}

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

func (this *BuiltInTypes) DeepCopy(that *BuiltInTypes) {
	deriveDeepCopyPtrToBuiltInTypes(that, this)
}

func (this *BuiltInTypes) GoString() string {
	return deriveGoStringBuiltInTypes(this)
}

func (this *BuiltInTypes) Clone() *BuiltInTypes {
	return deriveCloneBuiltInTypes(this)
}

func (this *BuiltInTypes) Hash() uint64 {
	return deriveHashBuiltInTypes(this)
}

type PrivateBuiltInTypes struct {
	privateBool       bool
	privateByte       byte
	privateComplex128 complex128
	privateComplex64  complex64
	privateFloat64    float64
	privateFloat32    float32
	privateInt        int
	privateInt16      int16
	privateInt32      int32
	privateInt64      int64
	privateInt8       int8
	privateRune       rune
	privateString     string
	privateUint       uint
	privateUint16     uint16
	privateUint32     uint32
	privateUint64     uint64
	privateUint8      uint8
	privateUintPtr    uintptr
}

func (this *PrivateBuiltInTypes) Equal(that *PrivateBuiltInTypes) bool {
	return deriveEqualPtrToPrivateBuiltInTypes(this, that)
}

func (this *PrivateBuiltInTypes) Compare(that *PrivateBuiltInTypes) int {
	return deriveComparePtrToPrivateBuiltInTypes(this, that)
}

func (this *PrivateBuiltInTypes) DeepCopy(that *PrivateBuiltInTypes) {
	deriveDeepCopyPtrToPrivateBuiltInTypes(that, this)
}

func (this *PrivateBuiltInTypes) Hash() uint64 {
	return deriveHashPtrToPrivateBuiltInTypes(this)
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

func (this *PtrToBuiltInTypes) DeepCopy(that *PtrToBuiltInTypes) {
	deriveDeepCopyPtrToPtrToBuiltInTypes(that, this)
}

func (this *PtrToBuiltInTypes) GoString() string {
	return deriveGoStringPtrToBuiltInTypes(this)
}

func (this *PtrToBuiltInTypes) Clone() *PtrToBuiltInTypes {
	return deriveClonePtrToBuiltInTypes(this)
}

func (this *PtrToBuiltInTypes) Hash() uint64 {
	return deriveHashPtrToBuiltInTypes(this)
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

func (this *SliceOfBuiltInTypes) DeepCopy(that *SliceOfBuiltInTypes) {
	deriveDeepCopyPtrToSliceOfBuiltInTypes(that, this)
}

func (this *SliceOfBuiltInTypes) GoString() string {
	return deriveGoStringSliceOfBuiltInTypes(this)
}

func (this *SliceOfBuiltInTypes) Hash() uint64 {
	return deriveHashSliceOfBuiltInTypes(this)
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

func (this *SliceOfPtrToBuiltInTypes) DeepCopy(that *SliceOfPtrToBuiltInTypes) {
	deriveDeepCopyPtrToSliceOfPtrToBuiltInTypes(that, this)
}

func (this *SliceOfPtrToBuiltInTypes) GoString() string {
	return deriveGoStringSliceOfPtrToBuiltInTypes(this)
}

func (this *SliceOfPtrToBuiltInTypes) Hash() uint64 {
	return deriveHashSliceOfPtrToBuiltInTypes(this)
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

func (this *ArrayOfBuiltInTypes) DeepCopy(that *ArrayOfBuiltInTypes) {
	deriveDeepCopyPtrToArrayOfBuiltInTypes(that, this)
}

func (this *ArrayOfBuiltInTypes) GoString() string {
	return deriveGoStringArrayOfBuiltInTypes(this)
}

func (this *ArrayOfBuiltInTypes) Hash() uint64 {
	return deriveHashArrayOfBuiltInTypes(this)
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

func (this *ArrayOfPtrToBuiltInTypes) DeepCopy(that *ArrayOfPtrToBuiltInTypes) {
	deriveDeepCopyPtrToArrayOfPtrToBuiltInTypes(that, this)
}

func (this *ArrayOfPtrToBuiltInTypes) GoString() string {
	return deriveGoStringArrayOfPtrToBuiltInTypes(this)
}

func (this *ArrayOfPtrToBuiltInTypes) Hash() uint64 {
	return deriveHashArrayOfPtrToBuiltInTypes(this)
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

func (this *MapsOfSimplerBuiltInTypes) DeepCopy(that *MapsOfSimplerBuiltInTypes) {
	deriveDeepCopyPtrToMapsOfSimplerBuiltInTypes(that, this)
}

func (this *MapsOfSimplerBuiltInTypes) GoString() string {
	return deriveGoStringMapsOfSimplerBuiltInTypes(this)
}

func (this *MapsOfSimplerBuiltInTypes) Hash() uint64 {
	return deriveHashMapsOfSimplerBuiltInTypes(this)
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

func (this *MapsOfBuiltInTypes) DeepCopy(that *MapsOfBuiltInTypes) {
	deriveDeepCopyPtrToMapsOfBuiltInTypes(that, this)
}

func (this *MapsOfBuiltInTypes) GoString() string {
	return deriveGoStringMapsOfBuiltInTypes(this)
}

func (this *MapsOfBuiltInTypes) Hash() uint64 {
	return deriveHashMapsOfBuiltInTypes(this)
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

func (this *SliceToSlice) DeepCopy(that *SliceToSlice) {
	deriveDeepCopyPtrToSliceToSlice(that, this)
}

func (this *SliceToSlice) GoString() string {
	return deriveGoStringSliceToSlice(this)
}

func (this *SliceToSlice) Hash() uint64 {
	return deriveHashSliceToSlice(this)
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

func (this *PtrTo) DeepCopy(that *PtrTo) {
	deriveDeepCopyPtrToPtrTo(that, this)
}

func (this *PtrTo) GoString() string {
	return deriveGoStringPtrTo(this)
}

func (this *PtrTo) Hash() uint64 {
	return deriveHashPtrTo(this)
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

func (this *Name) DeepCopy(that *Name) {
	deriveDeepCopyPtrToName(that, this)
}

func (this *Name) GoString() string {
	return deriveGoStringName(this)
}

func (this *Name) Hash() uint64 {
	return deriveHashName(this)
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

func (this *Structs) DeepCopy(that *Structs) {
	deriveDeepCopyPtrToStructs(that, this)
}

func (this *Structs) GoString() string {
	return deriveGoStringStructs(this)
}

func (this *Structs) Hash() uint64 {
	return deriveHashStructs(this)
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

func (this *MapWithStructs) DeepCopy(that *MapWithStructs) {
	deriveDeepCopyPtrToMapWithStructs(that, this)
}

func (this *MapWithStructs) GoString() string {
	return deriveGoStringMapWithStructs(this)
}

func (this *MapWithStructs) Hash() uint64 {
	return deriveHashMapWithStructs(this)
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

func (this *RecursiveType) DeepCopy(that *RecursiveType) {
	deriveDeepCopyPtrToRecursiveType(that, this)
}

func (this *RecursiveType) GoString() string {
	return deriveGoStringRecursiveType(this)
}

func (this *RecursiveType) Hash() uint64 {
	return deriveHashRecursiveType(this)
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

func (this *EmbeddedStruct1) DeepCopy(that *EmbeddedStruct1) {
	deriveDeepCopyPtrToEmbeddedStruct1(that, this)
}

func (this *EmbeddedStruct1) GoString() string {
	return deriveGoStringEmbeddedStruct1(this)
}

func (this *EmbeddedStruct1) Hash() uint64 {
	return deriveHashEmbeddedStruct1(this)
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

func (this *EmbeddedStruct2) DeepCopy(that *EmbeddedStruct2) {
	deriveDeepCopyPtrToEmbeddedStruct2(that, this)
}

func (this *EmbeddedStruct2) GoString() string {
	return deriveGoStringEmbeddedStruct2(this)
}

func (this *EmbeddedStruct2) Hash() uint64 {
	return deriveHashEmbeddedStruct2(this)
}

type UnnamedStruct struct {
	Unnamed struct {
		String string
	}
}

func (this *UnnamedStruct) Equal(that *UnnamedStruct) bool {
	return deriveEqualPtrToUnnamedStruct(this, that)
}

func (this *UnnamedStruct) DeepCopy(that *UnnamedStruct) {
	deriveDeepCopyPtrToUnnamedStruct(that, this)
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

func (this *StructWithStructFieldWithoutEqualMethod) DeepCopy(that *StructWithStructFieldWithoutEqualMethod) {
	deriveDeepCopyPtrToStructWithStructFieldWithoutEqualMethod(that, this)
}

func (this *StructWithStructFieldWithoutEqualMethod) GoString() string {
	return deriveGoStringStructWithStructFieldWithoutEqualMethod(this)
}

func (this *StructWithStructFieldWithoutEqualMethod) Hash() uint64 {
	return deriveHashStructWithStructFieldWithoutEqualMethod(this)
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

func (this *StructWithStructWithFromAnotherPackage) DeepCopy(that *StructWithStructWithFromAnotherPackage) {
	deriveDeepCopyPtrToStructWithStructWithFromAnotherPackage(that, this)
}

func (this *StructWithStructWithFromAnotherPackage) GoString() string {
	return deriveGoStringStructWithStructWithFromAnotherPackage(this)
}

func (this *StructWithStructWithFromAnotherPackage) Hash() uint64 {
	return deriveHashStructWithStructWithFromAnotherPackage(this)
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

func (this *FieldWithStructWithPrivateFields) DeepCopy(that *FieldWithStructWithPrivateFields) {
	deriveDeepCopyPtrToFieldWithStructWithPrivateFields(that, this)
}

type Enums struct {
	Enum             MyEnum
	PtrToEnum        *MyEnum
	SliceToEnum      []MyEnum
	SliceToPtrToEnum []*MyEnum
	MapToEnum        map[int32]MyEnum
	EnumToMap        map[MyEnum]int32
	ArrayEnum        [2]MyEnum
}

type MyEnum int32

func (this *Enums) Equal(that *Enums) bool {
	return deriveEqualPtrToEnums(this, that)
}

func (this *Enums) Compare(that *Enums) int {
	return deriveComparePtrToEnums(this, that)
}

func (this *Enums) DeepCopy(that *Enums) {
	deriveDeepCopyPtrToEnums(that, this)
}

func (this *Enums) GoString() string {
	return deriveGoStringEnums(this)
}

func (this *Enums) Hash() uint64 {
	return deriveHashEnums(this)
}

type NamedTypes struct {
	Slice        MySlice
	PtrToSlice   *MySlice
	SliceToSlice []MySlice
}

type MySlice []int64

func (this *NamedTypes) Equal(that *NamedTypes) bool {
	return deriveEqualPtrToNamedTypes(this, that)
}

func (this *NamedTypes) Compare(that *NamedTypes) int {
	return deriveComparePtrToNamedTypes(this, that)
}

func (this *NamedTypes) DeepCopy(that *NamedTypes) {
	deriveDeepCopyPtrToNamedTypes(that, this)
}

func (this *NamedTypes) GoString() string {
	return deriveGoStringNamedTypes(this)
}

func (this *NamedTypes) Hash() uint64 {
	return deriveHashNamedTypes(this)
}

type Time struct {
	T time.Time
	P *time.Time
	// Ts  []time.Time
	// TPs []*time.Time
	// MT  map[int]time.Time
}

func (this *Time) Generate(rand *rand.Rand, size int) reflect.Value {
	if size == 0 {
		this = nil
		return reflect.ValueOf(this)
	}
	this = &Time{}
	if size == 1 {
		this.T = time.Unix(0, rand.Int63())
		return reflect.ValueOf(this)
	}
	this.T = time.Unix(0, rand.Int63())
	t := time.Unix(0, rand.Int63())
	this.P = &t
	return reflect.ValueOf(this)
}

func (this *Time) Equal(that *Time) bool {
	return deriveEqualPtrToTime(this, that)
}

func (this *Time) DeepCopy(that *Time) {
	deriveDeepCopyPtrToTime(that, this)
}

type Duration struct {
	D   time.Duration
	P   *time.Duration
	Ds  []time.Duration
	DPs []*time.Duration
	MD  map[int]time.Duration
}

func (this *Duration) Equal(that *Duration) bool {
	return deriveEqualPtrToDuration(this, that)
}

func (this *Duration) Compare(that *Duration) int {
	return deriveComparePtrToDuration(this, that)
}

func (this *Duration) DeepCopy(that *Duration) {
	deriveDeepCopyPtrToDuration(that, this)
}

func (this *Duration) GoString() string {
	return deriveGoStringDuration(this)
}

func (this *Duration) Hash() uint64 {
	return deriveHashDuration(this)
}

type Nickname struct {
	Alias map[string][]*pickle.Rick
}

func (this *Nickname) Equal(that *Nickname) bool {
	return deriveEqualPtrToNickname(this, that)
}

func (this *Nickname) Compare(that *Nickname) int {
	return deriveComparePtrToNickname(this, that)
}

func (this *Nickname) DeepCopy(that *Nickname) {
	deriveDeepCopyPtrToNickname(that, this)
}

func (this *Nickname) GoString() string {
	return deriveGoStringNickname(this)
}

func (this *Nickname) Hash() uint64 {
	return deriveHashNickname(this)
}

type privateStruct struct {
	ptrfield *int
}

type PrivateEmbedded struct {
	privateStruct
}

func (this *PrivateEmbedded) Generate(rand *rand.Rand, size int) reflect.Value {
	if size == 0 {
		this = nil
		return reflect.ValueOf(this)
	}
	this = &PrivateEmbedded{}
	if size == 1 {
		return reflect.ValueOf(this)
	}
	i := rand.Int()
	this.ptrfield = &i
	return reflect.ValueOf(this)
}

func (this *PrivateEmbedded) Equal(that *PrivateEmbedded) bool {
	return deriveEqualPtrToPrivateEmbedded(this, that)
}

func (this *PrivateEmbedded) Compare(that *PrivateEmbedded) int {
	return deriveComparePtrToPrivateEmbedded(this, that)
}

func (this *PrivateEmbedded) DeepCopy(that *PrivateEmbedded) {
	deriveDeepCopyPtrToPrivateEmbedded(that, this)
}

func (this *PrivateEmbedded) GoString() string {
	return deriveGoStringPrivateEmbedded(this)
}

func (this *PrivateEmbedded) Hash() uint64 {
	return deriveHashPrivateEmbedded(this)
}

type StructOfStructs struct {
	S1, S2 Structs
}

func (this *StructOfStructs) DeepCopy(that *StructOfStructs) {
	deriveDeepCopyPtrToStructOfStructs(that, this)
}
