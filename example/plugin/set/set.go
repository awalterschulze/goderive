package set

func subset(set, sub []int) bool {
	s := deriveSet(set)
	for _, k := range sub {
		if _, ok := s[k]; !ok {
			return false
		}
	}
	return true
}
