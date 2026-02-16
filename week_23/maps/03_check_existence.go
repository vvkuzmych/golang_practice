package main

import "fmt"

// 03. Check Existence - Check if key exists in map

func main() {
	ages := map[string]int{
		"Alice": 30,
		"Bob":   25,
		"Carol": 35,
	}

	// Idiomatic check: value, ok
	if age, ok := ages["Alice"]; ok {
		fmt.Println("Alice's age:", age)
	}

	if age, ok := ages["Dave"]; ok {
		fmt.Println("Dave's age:", age)
	} else {
		fmt.Println("Dave not found")
	}

	// Zero value when key missing - careful with 0, false, ""
	score := map[string]int{"a": 0, "b": 100}
	if v, ok := score["a"]; ok {
		fmt.Println("a scored:", v) // 0 - exists!
	}
	if v, ok := score["c"]; ok {
		fmt.Println("c scored:", v)
	} else {
		fmt.Println("c not in map")
	}
}
