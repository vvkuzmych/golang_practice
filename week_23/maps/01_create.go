package main

import "fmt"

func main() {
	m := make(map[string]int)
	m["age"] = 30
	fmt.Println(m)
	m2 := map[string]int{"a": 1, "b": 2}
	fmt.Println(m2)
}
