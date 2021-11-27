package main

import (
	"fmt"
	"regexp"
)

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
