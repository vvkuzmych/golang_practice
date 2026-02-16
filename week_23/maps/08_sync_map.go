package main

import (
	"fmt"
	"sync"
)

// 08. Sync Map - Concurrent-safe map with sync.Map

func main() {
	var m sync.Map

	// Store values
	m.Store("key1", "value1")
	m.Store("key2", 42)
	m.Store("key3", 3.14)

	// Load value
	if v, ok := m.Load("key1"); ok {
		fmt.Println("key1:", v)
	}

	// LoadOrStore - get or set
	if v, loaded := m.LoadOrStore("key4", "new"); loaded {
		fmt.Println("key4 existed:", v)
	} else {
		fmt.Println("key4 stored:", v)
	}

	// Delete
	m.Delete("key2")

	// Range over all entries
	fmt.Print("All entries: ")
	m.Range(func(k, v interface{}) bool {
		fmt.Printf("%v=%v ", k, v)
		return true
	})
	fmt.Println()
}
