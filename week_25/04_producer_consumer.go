package main

import (
	"fmt"
	"sync"
	"time"
)

// 04. Producer/Consumer Pattern in Go

func main() {
	fmt.Println("=== Go Producer/Consumer ===")
	fmt.Println()

	dataCh := make(chan int, 5) // Buffered channel
	var wg sync.WaitGroup
	counter := 0
	var mu sync.Mutex

	// Producer goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			time.Sleep(100 * time.Millisecond)
			dataCh <- i
			fmt.Printf("Producer: added %d\n", i)
		}
		close(dataCh) // Signal done
		fmt.Println("Producer: finished")
	}()

	// Consumer goroutines
	for id := 0; id < 3; id++ {
		wg.Add(1)
		go func(consumerID int) {
			defer wg.Done()
			for item := range dataCh {
				fmt.Printf("  Consumer %d: processing %d\n", consumerID, item)
				time.Sleep(200 * time.Millisecond)

				mu.Lock()
				counter++
				mu.Unlock()
			}
			fmt.Printf("  Consumer %d: shutting down\n", consumerID)
		}(id)
	}

	// Wait for completion
	wg.Wait()

	fmt.Println()
	fmt.Printf("✅ Processed %d items\n", counter)
	fmt.Println()

	// Advanced: Multiple producers, multiple consumers
	fmt.Println("=== Multiple Producers & Consumers ===")
	fmt.Println()

	taskCh := make(chan string, 10)
	itemsProduced := 0
	itemsConsumed := 0
	var prodMu, consMu sync.Mutex

	// Multiple producers
	for id := 0; id < 3; id++ {
		wg.Add(1)
		go func(producerID int) {
			defer wg.Done()
			for i := 0; i < 5; i++ {
				item := fmt.Sprintf("P%d-%d", producerID, i)
				taskCh <- item
				prodMu.Lock()
				itemsProduced++
				prodMu.Unlock()
				fmt.Printf("Producer %d: %s\n", producerID, item)
				time.Sleep(50 * time.Millisecond)
			}
		}(id)
	}

	// Close channel after all producers done
	go func() {
		wg.Wait()
		close(taskCh)
	}()

	// Multiple consumers
	for id := 0; id < 2; id++ {
		wg.Add(1)
		go func(consumerID int) {
			defer wg.Done()
			for item := range taskCh {
				fmt.Printf("  Consumer %d: %s\n", consumerID, item)
				consMu.Lock()
				itemsConsumed++
				consMu.Unlock()
				time.Sleep(100 * time.Millisecond)
			}
		}(id)
	}

	wg.Wait()

	fmt.Println()
	fmt.Printf("Produced: %d, Consumed: %d\n", itemsProduced, itemsConsumed)
	fmt.Println()
	fmt.Println("✅ Multiple producers/consumers complete")
}

// Key points:
// - Channels are built-in and thread-safe
// - close(ch) to signal completion
// - range over channel to consume
// - Buffered channels for backpressure
// - Cleaner than Ruby Queue pattern
// - No sentinel values needed
