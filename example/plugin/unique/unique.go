package unique

type Visitor struct {
	UserName   *string
	RemoteAddr string
}

func uniqueVisitors(vs []*Visitor) int {
	return len(deriveUnique(vs))
}
