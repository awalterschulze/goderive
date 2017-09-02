package package1

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func random() *Type1 {
	v, ok := quick.Value(reflect.TypeOf(&Type1{}), r)
	if !ok {
		panic("unable to generate value")
	}
	return v.Interface().(*Type1)
}

func TestEqual1(t *testing.T) {
	a := random()
	if !a.Equal(a) {
		t.Fatalf("expected equal")
	}
	b := random()
	if want, got := reflect.DeepEqual(a, b), a.Equal(b); want != got {
		t.Fatalf("want %v got %v\n this = %#v\n that = %#v", want, got, a, b)
	}
}
