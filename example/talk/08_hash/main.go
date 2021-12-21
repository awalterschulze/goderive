package main

import (
	"fmt"
	"strings"
)

// expensive function with two parameters.
func findLastName(firstname string, age *int) *string {
	// Mock: Search for a possible last name via an API
	if age == nil {
		return nil
	}
	lastname := strings.Map(func(r rune) rune {
		return r + rune(*age)
	}, firstname)
	return &lastname
}

// A complicated enough function will need a hash function.
// All simpler input parameters generate code that doesn't require a hash function.
// This way we generate the code you would have written yourself, instead of always generating the generic code.
var getLastName = deriveMem(findLastName)

func main() {
	age := 2
	fmt.Printf("%v\n", getLastName("Donna", nil))
	fmt.Printf("%v\n", *getLastName("Ron", &age))
	fmt.Printf("%v\n", getLastName("Donna", nil))
	// If age was not a pointer, we wouldn't need a hash function:
	// https://github.com/awalterschulze/goderive/tree/master/example/plugin/mem
}
