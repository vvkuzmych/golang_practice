package main

import (
	"fmt"
	"time"
)

// 11. Rate Limiting - Обмеження частоти виконання

func main() {
	// Ticker для rate limiting (5 запитів на секунду)
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	requests := make(chan int, 10)
	for i := 1; i <= 10; i++ {
		requests <- i
	}
	close(requests)

	// Обробляємо запити з rate limiting
	for req := range requests {
		<-ticker.C // Чекаємо ticker
		fmt.Printf("Processing request %d at %v\n", req, time.Now())
	}
}
