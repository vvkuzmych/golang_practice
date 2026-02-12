package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// FetchWithContext робить HTTP запит з підтримкою context cancellation
func FetchWithContext(ctx context.Context, url string) (string, error) {
	// Перевірити чи context вже cancelled
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	// Створити HTTP request з context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Виконати запит
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// Context cancelled під час запиту
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Прочитати response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return string(body), nil
}

func main() {
	// Test 1: Successful request
	fmt.Println("=== Test 1: Normal Request ===")
	ctx := context.Background()
	_, err := FetchWithContext(ctx, "https://google.com")
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Println("✅ Success: Got response")
	}

	// Test 2: Request with timeout (will timeout)
	fmt.Println("\n=== Test 2: Timeout (1 second timeout, 3 second delay) ===")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	start := time.Now()
	_, err = FetchWithContext(ctx, "https://httpbin.org/delay/3")
	elapsed := time.Since(start)

	if err != nil {
		if err == context.DeadlineExceeded {
			fmt.Printf("❌ Timeout as expected after %v\n", elapsed)
		} else {
			fmt.Printf("❌ Error: %v\n", err)
		}
	} else {
		fmt.Println("✅ Success (unexpected)")
	}

	// Test 3: Already cancelled context
	fmt.Println("\n=== Test 3: Already Cancelled Context ===")
	ctx, cancel = context.WithCancel(context.Background())
	cancel() // Cancel одразу

	_, err = FetchWithContext(ctx, "https://google.com")
	if err == context.Canceled {
		fmt.Println("❌ Context was already cancelled (as expected)")
	} else {
		fmt.Printf("Unexpected result: %v\n", err)
	}

	// Test 4: Manual cancellation during request
	fmt.Println("\n=== Test 4: Manual Cancellation During Request ===")
	ctx, cancel = context.WithCancel(context.Background())

	// Cancel після 500ms
	go func() {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("Cancelling context...")
		cancel()
	}()

	start = time.Now()
	_, err = FetchWithContext(ctx, "https://httpbin.org/delay/5")
	elapsed = time.Since(start)

	if err != nil {
		fmt.Printf("❌ Request cancelled after %v: %v\n", elapsed, err)
	} else {
		fmt.Println("✅ Success (unexpected)")
	}

	// Test 5: Deadline
	fmt.Println("\n=== Test 5: Absolute Deadline ===")
	deadline := time.Now().Add(500 * time.Millisecond)
	ctx, cancel = context.WithDeadline(context.Background(), deadline)
	defer cancel()

	start = time.Now()
	_, err = FetchWithContext(ctx, "https://httpbin.org/delay/2")
	elapsed = time.Since(start)

	if err == context.DeadlineExceeded {
		fmt.Printf("❌ Deadline exceeded after %v (as expected)\n", elapsed)
	} else {
		fmt.Printf("Unexpected result: %v\n", err)
	}

	// Test 6: Success with generous timeout
	fmt.Println("\n=== Test 6: Success with Generous Timeout ===")
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	start = time.Now()
	_, err = FetchWithContext(ctx, "https://httpbin.org/delay/1")
	elapsed = time.Since(start)

	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Success after %v\n", elapsed)
	}
}
