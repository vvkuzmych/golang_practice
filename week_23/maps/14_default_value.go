package main

import "fmt"

// 14. Default Value - Map with default values when key missing

func getWithDefault[K comparable, V any](m map[K]V, key K, defaultVal V) V {
	if v, ok := m[key]; ok {
		return v
	}
	return defaultVal
}

func main() {
	config := map[string]string{
		"host": "localhost",
		"port": "8080",
	}

	// Get with default
	host := getWithDefault(config, "host", "0.0.0.0")
	proto := getWithDefault(config, "protocol", "http")
	fmt.Println("Host:", host, "Proto:", proto)

	// Count with default 0
	counts := map[string]int{"a": 1, "b": 2}
	fmt.Println("a:", getWithDefault(counts, "a", 0))
	fmt.Println("c:", getWithDefault(counts, "c", 0))

	// Nested default helper
	cache := make(map[string]int)
	getOrSet := func(k string, d int) int {
		if v, ok := cache[k]; ok {
			return v
		}
		cache[k] = d
		return d
	}
	fmt.Println("getOrSet('x', 10):", getOrSet("x", 10))
	fmt.Println("getOrSet('x', 99):", getOrSet("x", 99))
}
