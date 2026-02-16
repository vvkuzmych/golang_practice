package main

import (
	"fmt"
	"sort"
)

// 06. Sort Map - Sort map by key or value

func main() {
	m := map[string]int{
		"banana": 3,
		"apple":  1,
		"cherry": 2,
	}

	// Sort by key
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	fmt.Println("By key:")
	for _, k := range keys {
		fmt.Printf("  %s: %d\n", k, m[k])
	}

	// Sort by value - need pair struct
	type pair struct {
		key   string
		value int
	}
	pairs := make([]pair, 0, len(m))
	for k, v := range m {
		pairs = append(pairs, pair{k, v})
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].value < pairs[j].value })
	fmt.Println("By value:")
	for _, p := range pairs {
		fmt.Printf("  %s: %d\n", p.key, p.value)
	}
}
