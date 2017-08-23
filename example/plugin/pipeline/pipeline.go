package pipeline

import "strings"

var lines = []string{
	"my name is judge",
	"welcome judy welcome judy",
	"welcome hello welcome judy",
	"welcome goodbye welcome judy",
}

func toChan(lines []string) <-chan string {
	c := make(chan string)
	go func() {
		for _, line := range lines {
			c <- line
		}
		close(c)
	}()
	return c
}

func wordsize(line string) <-chan int {
	c := make(chan int)
	go func() {
		words := strings.Split(line, " ")
		for _, word := range words {
			c <- len(word)
		}
		close(c)
	}()
	return c
}

func totalWordSizes() int {
	sizes := derivePipeline(toChan, wordsize)
	total := 0
	for size := range sizes(lines) {
		total += size
	}
	return total
}
