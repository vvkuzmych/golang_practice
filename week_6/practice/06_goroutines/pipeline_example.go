package main

import (
	"fmt"
	"sync"
	"time"
)

// ===== Simple Pipeline =====

func simplePipeline() {
	fmt.Println("\n=== Simple Pipeline ===")

	// Stage 1 → Stage 2 → Stage 3
	nums := generateNums(1, 2, 3, 4, 5)
	squared := squareNums(nums)
	evenOnly := filterEven(squared)

	// Collect results
	for result := range evenOnly {
		fmt.Printf("Result: %d\n", result)
	}
}

func generateNums(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			fmt.Printf("Generate: %d\n", n)
			out <- n
		}
	}()
	return out
}

func squareNums(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			result := n * n
			fmt.Printf("Square: %d → %d\n", n, result)
			out <- result
		}
	}()
	return out
}

func filterEven(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			if n%2 == 0 {
				fmt.Printf("Filter: %d (passed)\n", n)
				out <- n
			} else {
				fmt.Printf("Filter: %d (rejected)\n", n)
			}
		}
	}()
	return out
}

// ===== Fan-Out / Fan-In =====

func fanOutFanIn() {
	fmt.Println("\n=== Fan-Out / Fan-In ===")

	// Input
	input := generateNums(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	// Fan-Out: 3 workers працюють паралельно
	worker1 := slowSquare(input, 1)
	worker2 := slowSquare(input, 2)
	worker3 := slowSquare(input, 3)

	// Fan-In: об'єднуємо результати
	results := merge(worker1, worker2, worker3)

	// Output
	for result := range results {
		fmt.Printf("Final result: %d\n", result)
	}
}

func slowSquare(in <-chan int, workerID int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			time.Sleep(100 * time.Millisecond) // Simulate slow work
			result := n * n
			fmt.Printf("Worker %d: %d → %d\n", workerID, n, result)
			out <- result
		}
	}()
	return out
}

func merge(channels ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for n := range c {
				out <- n
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// ===== Real-World: Data Processing =====

type Record struct {
	ID    int
	Value string
}

func dataProcessing() {
	fmt.Println("\n=== Data Processing Pipeline ===")

	// Pipeline: Generate → Validate → Transform → Save
	records := generateRecords(10)
	validated := validateRecords(records)
	transformed := transformRecords(validated)

	// Save results
	count := 0
	for record := range transformed {
		fmt.Printf("Saved: ID=%d, Value=%s\n", record.ID, record.Value)
		count++
	}
	fmt.Printf("Total processed: %d records\n", count)
}

func generateRecords(count int) <-chan Record {
	out := make(chan Record)
	go func() {
		defer close(out)
		for i := 1; i <= count; i++ {
			record := Record{
				ID:    i,
				Value: fmt.Sprintf("data_%d", i),
			}
			fmt.Printf("Generate: %+v\n", record)
			out <- record
		}
	}()
	return out
}

func validateRecords(in <-chan Record) <-chan Record {
	out := make(chan Record)
	go func() {
		defer close(out)
		for record := range in {
			// Validate (skip odd IDs for demo)
			if record.ID%2 == 0 {
				fmt.Printf("Validate: %+v (passed)\n", record)
				out <- record
			} else {
				fmt.Printf("Validate: %+v (rejected)\n", record)
			}
		}
	}()
	return out
}

func transformRecords(in <-chan Record) <-chan Record {
	out := make(chan Record)
	go func() {
		defer close(out)
		for record := range in {
			// Transform
			record.Value = record.Value + "_transformed"
			fmt.Printf("Transform: %+v\n", record)
			out <- record
		}
	}()
	return out
}

// ===== Main =====

func main() {
	fmt.Println("🔄 Pipeline Pattern Examples")
	fmt.Println("==============================")

	simplePipeline()
	fanOutFanIn()
	dataProcessing()

	fmt.Println("\n✅ All pipeline examples completed!")
}
