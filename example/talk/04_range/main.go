package main

import "fmt"

type Person struct {
	name string
	age  int
}

func main() {
	x := Person{"Donna", 19}
	y := Person{"Ron", 39}
	people := map[string]Person{x.name: x, y.name: y}
	// My most popular use case: Deterministic range over a map in one line.
	for _, name := range Sort(Keys(people)) {
		fmt.Printf("%v\n", people[name])
	}
}
