package main

import (
	"fmt"
)

func main() {
	xint, yint := int(1), int(1)
	mint := Min(xint, yint)
	fmt.Printf("min int = %v\n", mint)
	xuint, yuint := uint(3), uint(4)
	// rename Min_ to Min and run goderive and see it being rename, such that the two different input types can be supported.
	muint := Min_(xuint, yuint)
	fmt.Printf("min uint = %v\n", muint)
}
