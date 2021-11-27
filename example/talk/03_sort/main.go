package main

import "fmt"

type Person struct {
	name string
	age  int
}

func main() {
	x := Person{"Donna", 19}
	y := Person{"Ron", 39}
	people := []Person{y, x}
	// If we have compare, then we can sort
	fmt.Printf("%v\n", Sort(people))
}
