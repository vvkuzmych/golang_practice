package main

import (
	"fmt"
	"time"
)

// 12. Ticker - Періодичне виконання завдань

func main() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
			}
		}
	}()

	// Працюємо 2 секунди
	time.Sleep(2 * time.Second)
	done <- true
	fmt.Println("Ticker stopped")
}
