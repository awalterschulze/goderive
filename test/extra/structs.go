package extra

import (
	"math/rand"
	"reflect"
)

type StructWithoutEqualMethod struct {
	Number int64
}

type PrivateFieldAndNoEqualMethod struct {
	number int64
}

func (this *PrivateFieldAndNoEqualMethod) Generate(rand *rand.Rand, size int) reflect.Value {
	if size == 0 {
		this = nil
		return reflect.ValueOf(this)
	}
	this = &PrivateFieldAndNoEqualMethod{}
	this.number = rand.Int63()
	return reflect.ValueOf(this)
}
