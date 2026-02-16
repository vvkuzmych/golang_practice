package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// 20. Graceful Shutdown - Корректне завершення горутин

func worker(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d: shutting down gracefully\n", id)
			// Cleanup resources
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("Worker %d: cleanup done\n", id)
			return
		default:
			fmt.Printf("Worker %d: working...\n", id)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	// Створюємо context з можливістю cancel
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Запускаємо workers
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(ctx, i, &wg)
	}

	// Ловимо сигнали завершення
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Чекаємо сигнал
	<-sigChan
	fmt.Println("\nReceived shutdown signal, gracefully shutting down...")

	// Відправляємо cancel signal
	cancel()

	// Чекаємо завершення всіх workers
	wg.Wait()
	fmt.Println("All workers shut down, exiting")
}
