package main

import (
	"fmt"
	"sync"
	"time"
)

// 14. Semaphore Pattern - Обмеження кількості одночасних горутин

func worker(id int, sem chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	// Acquire semaphore
	sem <- struct{}{}
	defer func() { <-sem }() // Release semaphore

	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(time.Second)
	fmt.Printf("Worker %d done\n", id)
}

func main() {
	maxConcurrent := 3
	sem := make(chan struct{}, maxConcurrent) // Semaphore
	var wg sync.WaitGroup

	// Запускаємо 10 workers, але тільки 3 працюють одночасно
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go worker(i, sem, &wg)
	}

	wg.Wait()
	fmt.Println("All workers done")
}
