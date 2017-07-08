package max

func negative(i, j int) bool {
	return deriveMax(i, j) < 0
}

type boat struct {
	length int
}

func addIfBigger(boats []boat, big boat) []boat {
	if len(boats) == 0 {
		return append(boats, big)
	}
	m := deriveMaxs(boats, boats[0])
	if deriveCompare(big, m) > 0 {
		return append(boats, big)
	}
	return boats
}
