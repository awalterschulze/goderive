package apply

import (
	"strings"
)

func splitThreeBySeps(s string) ([]string, []string) {
	splitThree := deriveApply(strings.SplitN, 3)

	commaSplit := splitThree(s, ",")
	colonSplit := splitThree(s, ":")

	return commaSplit, colonSplit
}
