package main

import "fmt"

func do(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Printf("Int: %v\n", v*2)
	case string:
		fmt.Printf("String: %v\n", v+" world")
	default:
		fmt.Printf("Unknown type: %T\n", v)
	}
}
func main() {
	do(21)
	do("hello")
	do(true)
}
