package main

import "fmt"

// 11. Merge - Merge multiple maps

func merge[K comparable, V any](maps ...map[K]V) map[K]V {
	result := make(map[K]V)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

func main() {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 20, "c": 3}
	m3 := map[string]int{"c": 30, "d": 4}

	merged := merge(m1, m2, m3)
	fmt.Println("Merged:", merged)
	// Later maps overwrite: b=20, c=30

	// Merge with overwrite control
	result := make(map[string]int)
	for k, v := range m1 {
		result[k] = v
	}
	for k, v := range m2 {
		if _, exists := result[k]; !exists {
			result[k] = v
		}
	}
	fmt.Println("Merge (no overwrite):", result)
}
