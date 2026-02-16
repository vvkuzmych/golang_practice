package main

import (
	"context"
	"fmt"
	"time"
)

// 08. Context Cancellation - Скасування горутин через context

func worker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d: cancelled\n", id)
			return
		default:
			fmt.Printf("Worker %d: working...\n", id)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// Запускаємо 3 workers
	for i := 1; i <= 3; i++ {
		go worker(ctx, i)
	}

	// Працюємо 2 секунди
	time.Sleep(2 * time.Second)

	// Скасовуємо всі горутини
	fmt.Println("Cancelling all workers...")
	cancel()

	// Даємо час на завершення
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Main done")
}
