package main

import "fmt"

// 09. Set - Map as set implementation (map[T]struct{})

func main() {
	// Set of strings - struct{} uses zero memory
	seen := make(map[string]struct{})

	words := []string{"a", "b", "a", "c", "b", "a"}
	for _, w := range words {
		seen[w] = struct{}{}
	}
	fmt.Println("Unique words:", len(seen))

	// Check membership
	if _, ok := seen["a"]; ok {
		fmt.Println("'a' is in set")
	}

	// Set of ints
	numbers := map[int]struct{}{
		1: {},
		2: {},
		3: {},
	}
	numbers[4] = struct{}{}
	fmt.Println("Numbers set size:", len(numbers))

	// Set operations
	setA := map[string]struct{}{"a": {}, "b": {}, "c": {}}
	setB := map[string]struct{}{"b": {}, "c": {}, "d": {}}
	intersection := make(map[string]struct{})
	for k := range setA {
		if _, ok := setB[k]; ok {
			intersection[k] = struct{}{}
		}
	}
	fmt.Println("AcapB:", intersection)
}
