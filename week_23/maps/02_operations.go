package main

import "fmt"

func main() {
	m := map[string]int{"a": 1, "b": 2}
	fmt.Println(m["a"])
	m["c"] = 3
	delete(m, "b")
	fmt.Println(m)
}
