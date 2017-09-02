package package1

type Type1 struct {
	StringPtr *string
}

func (this *Type1) Equal(that *Type1) bool {
	return deriveEqual(this, that)
}
