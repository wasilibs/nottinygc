package main

import "strings"

// We cannot really run microbenchmarks with tinygo so we create a simple loop and measure
// time taken in the host.

//export run
func run() {
	for i := 0; i < 1000; i++ {
		_ = strings.Repeat("a", 100000)
		_ = strings.Repeat("b", 100000)
		_ = strings.Repeat("c", 100000)
	}
}

func main() {}
