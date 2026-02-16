package main

import (
	"fmt"
	"time"
)

// 20. Request/Response Pattern - Запит/Відповідь

type Request struct {
	ID       int
	Data     string
	Response chan Response
}

type Response struct {
	ID     int
	Result string
}

func worker(requests <-chan Request) {
	for req := range requests {
		// Обробляємо запит
		time.Sleep(100 * time.Millisecond)
		result := fmt.Sprintf("Processed: %s", req.Data)

		// Відправляємо відповідь
		req.Response <- Response{
			ID:     req.ID,
			Result: result,
		}
	}
}

func main() {
	requests := make(chan Request)

	// Запускаємо worker
	go worker(requests)

	// Відправляємо кілька запитів
	for i := 1; i <= 5; i++ {
		// Створюємо канал відповіді для кожного запиту
		responseChan := make(chan Response)

		// Відправляємо запит
		requests <- Request{
			ID:       i,
			Data:     fmt.Sprintf("Request %d", i),
			Response: responseChan,
		}

		// Чекаємо на відповідь
		response := <-responseChan
		fmt.Printf("Got response for request %d: %s\n", response.ID, response.Result)
	}

	close(requests)
}
