package main

import "fmt"

// 07. Range and Close - Ітерація по каналу

func main() {
	queue := make(chan string, 3)

	queue <- "one"
	queue <- "two"
	queue <- "three"
	close(queue) // Важливо закрити канал!

	// Range автоматично завершується при закритті каналу
	for elem := range queue {
		fmt.Println(elem)
	}

	// Перевірка чи канал закритий
	val, ok := <-queue
	if !ok {
		fmt.Println("Channel closed, no more values")
	} else {
		fmt.Println("Received:", val)
	}
}
