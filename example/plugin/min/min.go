package min

func positive(i, j int) bool {
	return deriveMin(i, j) >= 0
}

type boat struct {
	length int
}

func removeMin(boats []boat) []boat {
	if len(boats) == 0 {
		return boats
	}
	m := deriveMins(boats, boats[0])
	return deriveFilter(func(b boat) bool {
		return b != m
	}, boats)
}
