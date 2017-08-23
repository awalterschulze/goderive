package prefix

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func random() *MyStruct {
	v, _ := quick.Value(reflect.TypeOf(&MyStruct{}), r)
	return v.Interface().(*MyStruct)
}

func TestGoGenerate(t *testing.T) {
	this := random()
	that := random()
	if want, got := reflect.DeepEqual(this, that), this.Equal(that); want != got {
		t.Fatalf("want %v got %v", want, got)
	}
}
