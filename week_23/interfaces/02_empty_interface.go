package main

import "fmt"

func describe(i interface{}) {
	fmt.Printf("(%v, %T)\n", i, i)
}
func main() {
	describe(42)
	describe("hello")
	describe(true)
	describe([]int{1, 2, 3})
}
