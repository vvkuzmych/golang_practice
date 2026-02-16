package main

import (
	"fmt"
	"sync"
)

// 04. Goroutine with Closure - Горутина з замиканням

func main() {
	var wg sync.WaitGroup

	// НЕПРАВИЛЬНО - всі горутини побачать i=5
	fmt.Println("Wrong way:")
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(i) // Всі виведуть 5 або 6
		}()
	}
	wg.Wait()

	fmt.Println("\nCorrect way:")
	// ПРАВИЛЬНО - передаємо i як параметр
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			fmt.Println(n)
		}(i)
	}
	wg.Wait()
}
