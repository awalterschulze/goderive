package equal

type MyStruct struct {
	Int64     int64
	StringPtr *string
	Foo       *Foo
}

func (this *MyStruct) Equal(that *MyStruct) bool {
	return deriveEqual(this, that)
}

type Foo struct {
	Name  string
	other string
}

func (this *Foo) Equal(that *Foo) bool {
	return this.Name == that.Name
}
