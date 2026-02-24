package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// 05. Parallel HTTP Requests in Go

type Post struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func main() {
	fmt.Println("=== Go Parallel HTTP Requests ===")
	fmt.Println()

	urls := []string{
		"https://jsonplaceholder.typicode.com/posts/1",
		"https://jsonplaceholder.typicode.com/posts/2",
		"https://jsonplaceholder.typicode.com/posts/3",
		"https://jsonplaceholder.typicode.com/posts/4",
		"https://jsonplaceholder.typicode.com/posts/5",
	}

	// Sequential (slow)
	fmt.Println("1. Sequential requests:")
	start := time.Now()

	for _, url := range urls {
		resp, _ := http.Get(url)
		fmt.Printf("  %s: %s\n", url, resp.Status)
		resp.Body.Close()
	}

	sequentialTime := time.Since(start)
	fmt.Printf("  Time: %.2fs\n", sequentialTime.Seconds())
	fmt.Println()

	// Parallel (fast)
	fmt.Println("2. Parallel requests (with goroutines):")
	start = time.Now()

	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			resp, _ := http.Get(u)
			fmt.Printf("  %s: %s\n", u, resp.Status)
			resp.Body.Close()
		}(url)
	}

	wg.Wait()
	parallelTime := time.Since(start)

	fmt.Printf("  Time: %.2fs\n", parallelTime.Seconds())
	fmt.Printf("  Speedup: %.2fx\n", sequentialTime.Seconds()/parallelTime.Seconds())
	fmt.Println()

	// With error handling
	fmt.Println("3. With error handling:")

	type Result struct {
		Success bool
		Title   string
		URL     string
		Error   string
	}

	resultCh := make(chan Result, len(urls))

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()

			resp, err := http.Get(u)
			if err != nil {
				resultCh <- Result{Success: false, Error: err.Error(), URL: u}
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				resultCh <- Result{Success: false, Error: resp.Status, URL: u}
				return
			}

			body, _ := io.ReadAll(resp.Body)
			var post Post
			if err := json.Unmarshal(body, &post); err != nil {
				resultCh <- Result{Success: false, Error: err.Error(), URL: u}
				return
			}

			resultCh <- Result{Success: true, Title: post.Title, URL: u}
		}(url)
	}

	// Close result channel after all goroutines done
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Collect results
	for result := range resultCh {
		if result.Success {
			fmt.Printf("  ✓ %s\n", result.URL)
			title := result.Title
			if len(title) > 50 {
				title = title[:50] + "..."
			}
			fmt.Printf("    Title: %s\n", title)
		} else {
			fmt.Printf("  ✗ %s: %s\n", result.URL, result.Error)
		}
	}

	fmt.Println()
	fmt.Println("✅ Parallel HTTP with error handling complete")
}

// Key points:
// - Goroutines are lightweight
// - Channels to collect results
// - close(ch) after wg.Wait()
// - Struct for result with error
// - Very efficient for I/O
// - No GIL - true parallelism
