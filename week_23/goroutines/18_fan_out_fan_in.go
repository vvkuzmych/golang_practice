package main

import (
	"fmt"
	"sync"
)

// 18. Fan-Out / Fan-In Pattern - Розподіл і збір результатів

func producer(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

// Fan-Out: кілька workers обробляють дані
func worker(in <-chan int, id int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			fmt.Printf("Worker %d processing %d\n", id, n)
			out <- n * 2
		}
		close(out)
	}()
	return out
}

// Fan-In: об'єднуємо результати з кількох каналів
func merge(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Для кожного каналу запускаємо горутину
	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			out <- n
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go output(c)
	}

	// Закриваємо out після завершення всіх
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	// Producer
	nums := producer(1, 2, 3, 4, 5, 6, 7, 8)

	// Fan-Out: 3 workers
	w1 := worker(nums, 1)
	w2 := worker(nums, 2)
	w3 := worker(nums, 3)

	// Fan-In: merge results
	results := merge(w1, w2, w3)

	// Print results
	for result := range results {
		fmt.Printf("Result: %d\n", result)
	}
}
