package main

var fib func(uint) uint

func init() {
	fib = deriveMem(func(i uint) uint {
		if i == 0 || i == 1 {
			return i
		}
		return fib(i-1) + fib(i-2)
	})
}

func main() {
	println(fib(1))
	println(fib(5))
	println(fib(64))
	// yes it really works.
	println(fib(1000))
}
