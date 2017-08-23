package deepcopy

type MyStruct struct {
	Int64     int64
	StringPtr *string
}

func (m *MyStruct) Clone() *MyStruct {
	if m == nil {
		return nil
	}
	n := &MyStruct{}
	deriveDeepCopy(n, m)
	return n
}
