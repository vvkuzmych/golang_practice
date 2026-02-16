package main

import (
	"context"
	"fmt"
	"time"
)

// 09. Context Timeout - Автоматичне скасування через timeout

func slowOperation(ctx context.Context) error {
	select {
	case <-time.After(3 * time.Second):
		fmt.Println("Operation completed")
		return nil
	case <-ctx.Done():
		fmt.Println("Operation cancelled:", ctx.Err())
		return ctx.Err()
	}
}

func main() {
	// Context з timeout 1 секунда
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	fmt.Println("Starting operation...")
	err := slowOperation(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
