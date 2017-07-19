package contains

import (
	"fmt"
)

type boat struct {
	id int
}

func appendUnique(olds, news []boat) ([]boat, error) {
	for i := range news {
		if deriveContains(olds, news[i]) {
			return nil, fmt.Errorf("duplicate found")
		}
	}
	return append(olds, news...), nil
}
