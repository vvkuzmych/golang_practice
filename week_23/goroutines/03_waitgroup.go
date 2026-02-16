package main

import (
	"fmt"
	"sync"
	"time"
)

// 03. WaitGroup - Чекаємо завершення всіх горутин

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done() // Викликаємо Done() при завершенні

	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(time.Second)
	fmt.Printf("Worker %d done\n", id)
}

func main() {
	var wg sync.WaitGroup

	// Запускаємо 3 workers
	for i := 1; i <= 3; i++ {
		wg.Add(1) // Додаємо до WaitGroup
		go worker(i, &wg)
	}

	wg.Wait() // Чекаємо завершення всіх
	fmt.Println("All workers done")
}
