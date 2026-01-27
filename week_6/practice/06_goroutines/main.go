package main

import (
	"fmt"
	"sync"
	"time"
)

// ===== 1. Базові Goroutines =====

func simpleGoroutine() {
	fmt.Println("\n=== 1. Simple Goroutines ===")

	for i := 1; i <= 3; i++ {
		go func(n int) {
			fmt.Printf("Goroutine %d running\n", n)
		}(i)
	}

	time.Sleep(100 * time.Millisecond)
}

// ===== 2. Channels =====

func channelExample() {
	fmt.Println("\n=== 2. Channels ===")

	ch := make(chan string)

	go func() {
		ch <- "Hello from goroutine!"
	}()

	message := <-ch
	fmt.Println(message)
}

// ===== 3. WaitGroup =====

func waitGroupExample() {
	fmt.Println("\n=== 3. WaitGroup ===")

	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Worker %d starting\n", id)
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("Worker %d done\n", id)
		}(i)
	}

	wg.Wait()
	fmt.Println("All workers completed")
}

// ===== 4. Worker Pool =====

func workerPoolExample() {
	fmt.Println("\n=== 4. Worker Pool ===")

	jobs := make(chan int, 10)
	results := make(chan int, 10)

	// Запускаємо 3 workers
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	// Відправляємо 9 jobs
	for j := 1; j <= 9; j++ {
		jobs <- j
	}
	close(jobs)

	// Збираємо результати
	for r := 1; r <= 9; r++ {
		result := <-results
		fmt.Printf("Result: %d\n", result)
	}
}

func worker(id int, jobs <-chan int, results chan<- int) {
	for job := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, job)
		time.Sleep(50 * time.Millisecond)
		results <- job * 2
	}
}

// ===== 5. Select =====

func selectExample() {
	fmt.Println("\n=== 5. Select ===")

	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "from channel 1"
	}()

	go func() {
		time.Sleep(200 * time.Millisecond)
		ch2 <- "from channel 2"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("Received:", msg1)
		case msg2 := <-ch2:
			fmt.Println("Received:", msg2)
		case <-time.After(300 * time.Millisecond):
			fmt.Println("Timeout!")
		}
	}
}

// ===== 6. Mutex =====

func mutexExample() {
	fmt.Println("\n=== 6. Mutex (Safe Counter) ===")

	counter := &SafeCounter{}
	var wg sync.WaitGroup

	// 100 goroutines інкрементують counter
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()
	fmt.Printf("Final count: %d\n", counter.Value())
}

type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

// ===== 7. Pipeline =====

func pipelineExample() {
	fmt.Println("\n=== 7. Pipeline ===")

	// Pipeline: generate → square → sum
	nums := generate(1, 2, 3, 4, 5)
	squared := square(nums)

	sum := 0
	for n := range squared {
		sum += n
	}
	fmt.Printf("Sum of squares: %d\n", sum) // 1 + 4 + 9 + 16 + 25 = 55
}

func generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

// ===== Main =====

func main() {
	fmt.Println("🚀 Goroutines & Concurrency Examples")
	fmt.Println("=====================================")

	simpleGoroutine()
	channelExample()
	waitGroupExample()
	workerPoolExample()
	selectExample()
	mutexExample()
	pipelineExample()

	fmt.Println("\n✅ All examples completed!")
}
