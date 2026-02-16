package main

import (
	"context"
	"fmt"
	"time"
)

// 10. Context Deadline - Встановлення конкретного дедлайну

func processWithDeadline(ctx context.Context) {
	deadline, ok := ctx.Deadline()
	if ok {
		fmt.Printf("Deadline: %v\n", deadline)
		fmt.Printf("Time until deadline: %v\n", time.Until(deadline))
	}

	select {
	case <-time.After(2 * time.Second):
		fmt.Println("Processing completed")
	case <-ctx.Done():
		fmt.Println("Deadline exceeded:", ctx.Err())
	}
}

func main() {
	// Deadline через 1 секунду від зараз
	deadline := time.Now().Add(1 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	processWithDeadline(ctx)
}
