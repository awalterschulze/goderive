package traverse

import "strconv"

func toInts(ss []string) ([]int, error) {
	return deriveTraverse(strconv.Atoi, ss)
}
