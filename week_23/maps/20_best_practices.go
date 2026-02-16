package main

import "fmt"

// 20. Map Best Practices - Summary of Go map idioms

func main() {
	// 1. Initialize with make or literal
	m1 := make(map[string]int)
	m2 := map[string]int{"a": 1, "b": 2}

	// 2. Always check existence for zero-value keys
	_, ok := m1["missing"]
	fmt.Println("Key exists:", ok)

	// 3. Preallocate when size known
	knownSize := 100
	m3 := make(map[int]string, knownSize)
	_ = m3

	// 4. Maps are reference types - no need for pointer
	modifyMap(m2)
	fmt.Println("Modified m2:", m2)

	// 5. Use map for set: map[T]struct{}
	set := make(map[string]struct{})
	set["a"] = struct{}{}
	_, inSet := set["a"]
	fmt.Println("In set:", inSet)

	// 6. nil map - can read but not write
	var nilMap map[string]int
	fmt.Println("nil map len:", len(nilMap))
	// nilMap["x"] = 1 // panic!
}

func modifyMap(m map[string]int) {
	m["c"] = 3
}
