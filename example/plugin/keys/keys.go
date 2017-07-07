package keys

import (
	"strconv"
)

func printMap(m map[string]int) {
	for _, k := range deriveSort(deriveKeys(m)) {
		println(k + ":" + strconv.Itoa(m[k]))
	}
}
