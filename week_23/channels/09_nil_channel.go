package main

import (
	"fmt"
	"time"
)

// 09. Nil Channel - Поведінка nil каналу

func main() {
	var ch chan string

	// Nil канал блокує назавжди
	go func() {
		fmt.Println("Trying to receive from nil channel...")
		msg := <-ch                   // Блокується назавжди
		fmt.Println("Received:", msg) // Ніколи не виконається
	}()

	// Використовуємо timeout
	select {
	case <-ch:
		fmt.Println("Received from nil channel")
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout: nil channel blocks forever")
	}

	// Корисно для вимикання каналу в select
	done := make(chan bool)

	go func() {
		time.Sleep(500 * time.Millisecond)
		done <- true
	}()

	// Після отримання done, встановлюємо його в nil
	select {
	case <-done:
		fmt.Println("Done received, disabling channel")
		done = nil
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout")
	}

	// Тепер done = nil і не буде обиратись в select
	select {
	case <-done:
		fmt.Println("This won't execute")
	case <-time.After(100 * time.Millisecond):
		fmt.Println("Done channel is nil, timeout occurred")
	}
}
