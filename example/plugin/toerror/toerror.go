package toerror

import (
	"fmt"
	"reflect"
)

func getJsonTag(f reflect.StructTag) (string, error) {
	return deriveToError(fmt.Errorf("json tag not exists"), f.Lookup)("json")
}
