package union

type Person struct {
	Name string
	Vote *string
}

func ratio(survey, database []*Person) float64 {
	union := deriveUnion(deriveUnique(database), survey)
	if len(union) == 0 {
		return 0
	}
	voted := deriveFilter(func(p *Person) bool {
		return p.Vote != nil
	}, union)
	return float64(len(voted)) / float64(len(union))
}
