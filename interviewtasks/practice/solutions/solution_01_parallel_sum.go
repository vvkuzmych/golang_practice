package main

import (
	"fmt"
	"sync"
)

// ParallelSum розраховує суму slice паралельно
func ParallelSum(numbers []int, workers int) int {
	if len(numbers) == 0 {
		return 0
	}

	// Якщо workers більше ніж елементів, використовуємо len(numbers)
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
		// Визначити start та end для цього worker
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

func main() {
	// Test 1
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	result := ParallelSum(numbers, 2)
	fmt.Printf("Sum: %d (expected 55)\n\n", result)

	// Test 2
	numbers = []int{10, 20, 30, 40}
	result = ParallelSum(numbers, 4)
	fmt.Printf("Sum: %d (expected 100)\n\n", result)

	// Test 3: Large array
	numbers = make([]int, 1000000)
	for i := range numbers {
		numbers[i] = 1
	}
	result = ParallelSum(numbers, 8)
	fmt.Printf("Sum: %d (expected 1000000)\n", result)
}
