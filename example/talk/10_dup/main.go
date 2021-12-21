package main

import (
	"fmt"
)

func main() {
	xint, yint := int(1), int(1)
	mint := Min(xint, yint)
	fmt.Printf("min int = %v\n", mint)
	xint2, yint2 := int(3), int(4)
	// rename Min to Min2, run goderive and see the two functions being given the same name.
	mint2 := Min(xint2, yint2)
	fmt.Printf("min int 2 = %v\n", mint2)
}
