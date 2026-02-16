package main

import (
	"fmt"
	"sync"
)

// 16. sync.Once - Виконання функції тільки один раз

var (
	instance *Singleton
	once     sync.Once
)

type Singleton struct {
	value string
}

func GetInstance() *Singleton {
	once.Do(func() {
		fmt.Println("Creating singleton instance...")
		instance = &Singleton{value: "I'm singleton"}
	})
	return instance
}

func main() {
	var wg sync.WaitGroup

	// 10 горутин намагаються створити singleton
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			s := GetInstance()
			fmt.Printf("Goroutine %d: %s\n", id, s.value)
		}(i)
	}

	wg.Wait()
}
