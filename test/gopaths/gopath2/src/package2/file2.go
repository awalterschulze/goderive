package package2

import (
	"package1"
)

type Type2 struct {
	Field1 *package1.Type1
}

func (this *Type2) Equal(that *Type2) bool {
	return deriveEqual(this, that)
}
