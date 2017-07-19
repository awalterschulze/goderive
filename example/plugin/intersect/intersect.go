package intersect

func bestChoices(strategy1, strategy2 []int) []int {
	intersection := deriveIntersect(strategy1, strategy2)
	if len(intersection) == 0 {
		return append(strategy1, strategy2...)
	}
	return intersection
}
