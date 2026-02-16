package main

import (
	"fmt"
	"sort"
)

// 05. Keys and Values - Get all keys and values as slices

func keys[K comparable, V any](m map[K]V) []K {
	result := make([]K, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

func values[K comparable, V any](m map[K]V) []V {
	result := make([]V, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result
}

func main() {
	m := map[string]int{"b": 2, "a": 1, "c": 3}

	fmt.Println("Keys:", keys(m))
	fmt.Println("Values:", values(m))

	// Sort keys for deterministic iteration
	keyList := keys(m)
	sort.Strings(keyList)
	fmt.Println("Sorted keys:", keyList)
	for _, k := range keyList {
		fmt.Printf("  %s: %d\n", k, m[k])
	}
}
