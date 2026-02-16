package main

import (
	"fmt"
	"sort"
)

// 15. Map to Slice - Convert map to slice (keys, values, or pairs)

func main() {
	m := map[string]int{"b": 2, "a": 1, "c": 3}

	// Keys to slice
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	fmt.Println("Keys:", keys)

	// Values to slice
	values := make([]int, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	fmt.Println("Values:", values)

	// Key-value pairs to slice
	type Pair struct {
		Key   string
		Value int
	}
	pairs := make([]Pair, 0, len(m))
	for k, v := range m {
		pairs = append(pairs, Pair{k, v})
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].Key < pairs[j].Key })
	fmt.Println("Pairs:", pairs)
}
