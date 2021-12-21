package main

import "fmt"

type Person struct {
	name string
	age  int
}

func main() {
	x := &Person{"Donna", 19}
	y := &Person{"Ron", 39}
	// Min doesn't just work for ints, it generates a compare function *if needed*.
	fmt.Printf("%v\n", Min(x, y))
	// goderive inspects the input types to decide which code to generate.
}
