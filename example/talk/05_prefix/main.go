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
	// Usually the prefix derive is included if you do not pass any command line parameters.
	for _, name := range deriveSort(deriveKeys(people)) {
		fmt.Printf("%v\n", people[name])
	}
}
