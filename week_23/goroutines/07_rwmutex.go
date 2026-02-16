package main

import (
	"fmt"
	"sync"
	"time"
)

// 07. RWMutex - Читання/Запис mutex

type Cache struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]string),
	}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock() // Read lock
	defer c.mu.RUnlock()
	val, ok := c.data[key]
	return val, ok
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock() // Write lock
	defer c.mu.Unlock()
	c.data[key] = value
}

func main() {
	cache := NewCache()
	var wg sync.WaitGroup

	// Writer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			cache.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Multiple readers
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				val, ok := cache.Get(fmt.Sprintf("key%d", j))
				if ok {
					fmt.Printf("Reader %d: %s\n", id, val)
				}
				time.Sleep(50 * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
}
