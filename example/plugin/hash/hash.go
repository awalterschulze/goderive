package hash

import "reflect"

type Entry struct {
	Name    string
	Numbers []int
}

func hasDuplicates(es []*Entry) bool {
	m := make(map[uint64][]*Entry)
	for i, e := range es {
		h := deriveHash(e)
		for _, f := range m[h] {
			if reflect.DeepEqual(e, f) {
				return true
			}
		}
		if _, ok := m[h]; !ok {
			m[h] = []*Entry{es[i]}
		}
	}
	return false
}
