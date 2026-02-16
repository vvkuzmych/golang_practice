package main

import (
	"fmt"
	"sync"
)

// 02. Fan-Out Pattern - Розподіл роботи між workers

// Паттерн: Один producer → Кілька workers
//
//              ┌→ Worker1
// Producer ────┼→ Worker2  ──→ Results
//              └→ Worker3

//func producer(tasks ...int) <-chan int {
//	out := make(chan int)
//	go func() {
//		defer close(out)
//		for _, task := range tasks {
//			out <- task
//		}
//	}()
//	return out
//}
//
//func worker(id int, tasks <-chan int, results chan<- int, wg *sync.WaitGroup) {
//	defer wg.Done()
//
//	for task := range tasks {
//		fmt.Printf("Worker %d processing task %d\n", id, task)
//		time.Sleep(200 * time.Millisecond)
//		results <- task * 2
//	}
//	fmt.Printf("Worker %d finished\n", id)
//}

// Use cases:
// - Processing large datasets
// - Parallel HTTP requests
// - Image processing
// - Database batch operations
// - Task queue workers

// === Helper Functions (Advanced) ===

// MergeAllChannels об'єднує кілька каналів в один (Fan-In pattern)
//func MergeAllChannels[T any](channels ...<-chan T) <-chan T {
//	outputCh := make(chan T)
//	var wg sync.WaitGroup
//	wg.Add(len(channels))
//	for _, channel := range channels {
//		go func(ch <-chan T) {
//			defer wg.Done()
//			for task := range ch {
//				outputCh <- task
//			}
//		}(channel)
//	}
//	go func() {
//		wg.Wait()
//		close(outputCh)
//	}()
//	return outputCh
//}

// SplitChannels розподіляє один канал в кілька (Round-Robin)
func SplitChannels[T any](inputCh <-chan T, n int) []<-chan T {
	//формують кілька каналів
	//---------------------------------------
	outputChs := make([]chan T, n)
	for i := 0; i < n; i++ {
		outputChs[i] = make(chan T)
	}
	//-----------------------------------------

	//опрацьовує канал по раунд робін алгоритму
	//------------------------------------------
	go func() {
		idx := 0
		for task := range inputCh {
			outputChs[idx] <- task
			idx = (idx + 1) % n
		}
		for _, channel := range outputChs {
			close(channel)
		}
	}()
	//-----------------------------------------

	//повертають результат каналів в новостворених для читання
	//-----------------------------------------
	results := make([]<-chan T, n)
	for i := 0; i < n; i++ {
		results[i] = outputChs[i]
	}
	return results
	//-----------------------------------------
}

func main() {
	//fmt.Println("=== Fan-Out Pattern ===")
	//fmt.Println()
	//
	//// Producer: генеруємо завдання
	//tasks := producer(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	//
	//// Fan-Out: розподіляємо між 3 workers
	//results := make(chan int, 10)
	//var wg sync.WaitGroup
	//
	//numWorkers := 3
	//for w := 1; w <= numWorkers; w++ {
	//	wg.Add(1)
	//	go worker(w, tasks, results, &wg)
	//}
	//
	//// Закриваємо results після завершення всіх workers
	//go func() {
	//	wg.Wait()
	//	close(results)
	//}()
	//
	//// Збираємо результати
	//fmt.Println("\nResults:")
	//for result := range results {
	//	fmt.Printf("  Result: %d\n", result)
	//}
	//
	//fmt.Println()
	//fmt.Println("✅ All tasks processed by", numWorkers, "workers")

	//записують все в один канал
	//-----------------------------------------------
	channel := make(chan int)
	go func() {
		defer close(channel)
		for i := 0; i < 10; i++ {
			channel <- i
		}
	}()
	//-----------------------------------------------

	//розділяє канал на 2 канала
	//-----------------------------------------------
	channels := SplitChannels(channel, 2)
	//-----------------------------------------------

	var wg sync.WaitGroup
	wg.Add(len(channels))

	//зчитує з першого канала
	//-----------------------------------------------
	go func() {
		defer wg.Done()
		for task := range channels[0] {
			fmt.Println("Channel 1:", task)
		}
	}()
	//------------------------------------------------

	//зчитує з другого канала
	//-----------------------------------------------
	go func() {
		defer wg.Done()
		for task := range channels[1] {
			fmt.Println("Channel 2:", task)
		}
	}()
	//-----------------------------------------------

	wg.Wait()
}
