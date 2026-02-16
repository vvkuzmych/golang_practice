package main

import (
	"fmt"
	"sync"
)

// 02b. Fan-Out with Split Channels - Розподіл одного каналу на кілька

// Паттерн: Один канал → Split → Кілька каналів (Round-Robin)
//
//             ┌→ Channel1 → Consumer1
// Input ──────┼→ Channel2 → Consumer2
//             └→ Channel3 → Consumer3

// SplitChannels розподіляє один канал в кілька (Round-Robin)
func splitChannels[T any](input <-chan T, n int) []<-chan T {
	outputs := make([]chan T, n)
	for i := 0; i < n; i++ {
		outputs[i] = make(chan T)
	}

	go func() {
		defer func() {
			// Закриваємо всі канали ПІСЛЯ завершення розподілу
			for _, ch := range outputs {
				close(ch)
			}
		}()

		idx := 0
		for task := range input {
			outputs[idx] <- task
			idx = (idx + 1) % n
		}
	}()

	// Конвертуємо в read-only канали
	results := make([]<-chan T, n)
	for i := 0; i < n; i++ {
		results[i] = outputs[i]
	}
	return results
}

func main() {
	fmt.Println("=== Fan-Out with Split Channels ===")
	fmt.Println()

	// Producer: генеруємо завдання
	input := make(chan int)
	go func() {
		defer close(input)
		for i := 1; i <= 10; i++ {
			fmt.Printf("Producing: %d\n", i)
			input <- i
		}
	}()

	// SplitChannels: розподіляємо між 3 канали
	channels := splitChannels(input, 3)

	var wg sync.WaitGroup
	wg.Add(3)

	// Consumer 1
	go func() {
		defer wg.Done()
		for val := range channels[0] {
			fmt.Printf("  Channel 1 received: %d\n", val)
		}
		fmt.Println("  Channel 1 done")
	}()

	// Consumer 2
	go func() {
		defer wg.Done()
		for val := range channels[1] {
			fmt.Printf("  Channel 2 received: %d\n", val)
		}
		fmt.Println("  Channel 2 done")
	}()

	// Consumer 3
	go func() {
		defer wg.Done()
		for val := range channels[2] {
			fmt.Printf("  Channel 3 received: %d\n", val)
		}
		fmt.Println("  Channel 3 done")
	}()

	wg.Wait()

	fmt.Println()
	fmt.Println("✅ All tasks split between 3 channels (Round-Robin)")
	fmt.Println()
	fmt.Println("Distribution:")
	fmt.Println("  Ch1: 1, 4, 7, 10")
	fmt.Println("  Ch2: 2, 5, 8")
	fmt.Println("  Ch3: 3, 6, 9")
}

// Use cases:
// - Load balancing between multiple consumers
// - Distributing work evenly
// - Partitioning data streams
// - Round-robin task distribution
