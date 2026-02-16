package main

import "fmt"

// 04. Iterate - Iterate over map (order is random!)

func main() {
	m := map[string]int{
		"apple":  1,
		"banana": 2,
		"cherry": 3,
	}

	fmt.Println("Iteration (order varies):")
	for k, v := range m {
		fmt.Printf("  %s: %d\n", k, v)
	}

	// Key only
	fmt.Print("Keys: ")
	for k := range m {
		fmt.Print(k, " ")
	}
	fmt.Println()

	// Value only (discard key)
	fmt.Print("Values: ")
	for _, v := range m {
		fmt.Print(v, " ")
	}
	fmt.Println()
}
