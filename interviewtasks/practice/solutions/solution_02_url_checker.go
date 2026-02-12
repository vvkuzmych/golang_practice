package main

import (
	"fmt"
	"net/http"
	"time"
)

type URLStatus struct {
	URL        string
	StatusCode int
	Error      error
}

// CheckURLs перевіряє доступність URLs паралельно
func CheckURLs(urls []string) []URLStatus {
	if len(urls) == 0 {
		return []URLStatus{}
	}

	results := make([]URLStatus, len(urls))
	resultChan := make(chan URLStatus, len(urls))

	// HTTP client з timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Запустити goroutine для кожного URL
	for i, url := range urls {
		go func(index int, u string) {
			status := URLStatus{URL: u}

			// HTTP GET запит
			resp, err := client.Get(u)
			if err != nil {
				status.Error = err
				status.StatusCode = 0
			} else {
				status.StatusCode = resp.StatusCode
				resp.Body.Close()
			}

			// Надіслати результат в channel
			resultChan <- status
		}(i, url)
	}

	// Зібрати результати з channel
	statusMap := make(map[string]URLStatus)
	for i := 0; i < len(urls); i++ {
		status := <-resultChan
		statusMap[status.URL] = status
	}

	// Відновити порядок (важливо!)
	for i, url := range urls {
		results[i] = statusMap[url]
	}

	return results
}

func main() {
	urls := []string{
		"https://google.com",
		"https://github.com",
		"https://golang.org",
	}

	fmt.Println("Checking URLs...")
	start := time.Now()

	results := CheckURLs(urls)

	elapsed := time.Since(start)

	// Вивести результати
	for _, r := range results {
		if r.Error != nil {
			fmt.Printf("❌ %s - Error: %v\n", r.URL, r.Error)
		} else {
			fmt.Printf("✅ %s - Status: %d\n", r.URL, r.StatusCode)
		}
	}

	fmt.Printf("\nTotal time: %v\n", elapsed)
	fmt.Println("(Sequential would take 3x longer)")
}
