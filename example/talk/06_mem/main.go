package main

import (
	"fmt"
	"regexp"
)

// You can memoize function results to you do not redo expensive operations.
// Works a lot like functools.cache in python.
var re = deriveMem(func(r string) *regexp.Regexp {
	fmt.Printf("compiling regex <%s>\n", r)
	return regexp.MustCompile(r)
})

// var re = deriveMem(regexp.MustCompile)

func main() {
	fmt.Printf("%v\n", re("ab.*").MatchString("abc"))
	fmt.Printf("%v\n", re("cd.*").MatchString("cde"))
	fmt.Printf("%v\n", re("ab.*").MatchString("cde"))
}
