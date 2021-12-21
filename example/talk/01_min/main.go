package main

import "fmt"

func main() {
	x := 19
	y := 39
	// shortest way to write a min function in Go. A Go question from 2011
	m := Min(x, y)
	fmt.Printf("%v\n", m)
	// you call a function that doesn't exist and goderive generates it for you.
}
