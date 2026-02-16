package main

import (
	"fmt"
	"sync"
	"time"
)

// 12. Fan-Out Pattern - Розподіл роботи між workers

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

func square(in <-chan int, id int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			fmt.Printf("Worker %d processing %d\n", id, n)
			time.Sleep(100 * time.Millisecond)
			out <- n * n
		}
		close(out)
	}()
	return out
}

func main() {
	// Producer
	input := producer(1, 2, 3, 4, 5, 6, 7, 8)

	// Fan-Out: 3 workers обробляють паралельно
	c1 := square(input, 1)
	c2 := square(input, 2)
	c3 := square(input, 3)

	// Збираємо результати
	var wg sync.WaitGroup
	wg.Add(3)

	consume := func(ch <-chan int) {
		defer wg.Done()
		for val := range ch {
			fmt.Println("Result:", val)
		}
	}

	go consume(c1)
	go consume(c2)
	go consume(c3)

	wg.Wait()
}
