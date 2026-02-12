package main

import (
	"fmt"
	"sync"
)

func ParallelSum(numbers []int, workers int) int {
	if len(numbers) == 0 {
		return 0
	}

	if workers > len(numbers) {
		workers = len(numbers)
	}
	var (
		totalSum int
		mu       sync.Mutex
		wg       sync.WaitGroup
	)

	// Розмір chunk для кожного worker
	chunkSize := len(numbers) / workers

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		start := i * chunkSize
		end := start + chunkSize

		// Останній worker бере всі що залишились
		if i == workers-1 {
			end = len(numbers)
		}

		// Запустити worker
		go func(workerID int, s, e int) {
			defer wg.Done()

			// Локальна сума для цього worker
			localSum := 0
			for j := s; j < e; j++ {
				localSum += numbers[j]
			}

			// Додати до загальної суми (потрібен lock!)
			mu.Lock()
			totalSum += localSum
			mu.Unlock()

			fmt.Printf("Worker %d: processed %d elements, sum=%d\n", workerID, e-s, localSum)
		}(i, start, end)
	}

	// Чекати всі workers
	wg.Wait()

	return totalSum

}

//func main() {
//s := ParallelSum([]int{1, 2, 3}, 2)
//fmt.Println(s)
//}
