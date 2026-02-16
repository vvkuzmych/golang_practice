package main

import "fmt"

// 16. Slice to Map - Convert slice to map (by key or index)

func main() {
	// Slice to map by index
	words := []string{"a", "b", "c"}
	byIndex := make(map[int]string)
	for i, w := range words {
		byIndex[i] = w
	}
	fmt.Println("By index:", byIndex)

	// Slice to map by value as key
	byValue := make(map[string]bool)
	for _, w := range words {
		byValue[w] = true
	}
	fmt.Println("By value:", byValue)

	// Slice of structs to map by field
	type Item struct {
		ID   string
		Name string
	}
	items := []Item{{"1", "Apple"}, {"2", "Banana"}, {"3", "Cherry"}}
	byID := make(map[string]Item)
	for _, item := range items {
		byID[item.ID] = item
	}
	fmt.Println("By ID:", byID)
}
