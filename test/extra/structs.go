package extra

import (
	"math/rand"
	"reflect"
)

type StructWithoutEqualMethod struct {
	Number int64
}

type PrivateFieldAndNoEqualMethod struct {
	number    int64
	numbers   []int64
	ptr       *int64
	numberpts []*int64
	strct     *StructWithoutEqualMethod
}

func (this *PrivateFieldAndNoEqualMethod) Generate(rand *rand.Rand, size int) reflect.Value {
	if size == 0 {
		this = nil
		return reflect.ValueOf(this)
	}
	this = &PrivateFieldAndNoEqualMethod{}
	this.number = rand.Int63()
	if size == 1 {
		return reflect.ValueOf(this)
	}
	this.numbers = make([]int64, size/2)
	for i := 0; i < len(this.numbers); i++ {
		this.numbers[i] = rand.Int63()
	}
	if size == 2 {
		return reflect.ValueOf(this)
	}
	n := rand.Int63()
	this.ptr = &n
	if size == 3 {
		return reflect.ValueOf(this)
	}
	this.numberpts = make([]*int64, size/2)
	for i := 0; i < len(this.numberpts); i++ {
		n := rand.Int63()
		this.numberpts[i] = &n
	}
	if size == 4 {
		return reflect.ValueOf(this)
	}
	this.strct = &StructWithoutEqualMethod{}
	if size == 5 {
		return reflect.ValueOf(this)
	}
	this.strct.Number = rand.Int63()
	return reflect.ValueOf(this)
}
