package main

import (
	"fmt"
	"time"
)

// 06. Timeout - Таймаут для каналів

func main() {
	ch := make(chan string, 1)

	go func() {
		time.Sleep(2 * time.Second)
		ch <- "result"
	}()

	// Чекаємо максимум 1 секунду
	select {
	case res := <-ch:
		fmt.Println("Received:", res)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout!")
	}
}
