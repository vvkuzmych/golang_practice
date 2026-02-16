package main

import (
	"fmt"
	"sync"
)

// 18. Cache - Simple in-memory cache implementation

type Cache struct {
	mu    sync.RWMutex
	items map[string]string
}

func NewCache() *Cache {
	return &Cache{items: make(map[string]string)}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.items[key]
	return v, ok
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = value
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func main() {
	cache := NewCache()

	cache.Set("user:1", "Alice")
	cache.Set("user:2", "Bob")

	if v, ok := cache.Get("user:1"); ok {
		fmt.Println("user:1:", v)
	}

	cache.Delete("user:2")
	if _, ok := cache.Get("user:2"); !ok {
		fmt.Println("user:2 deleted")
	}
}
