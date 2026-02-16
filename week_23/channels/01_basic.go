package main

import "fmt"

// 01. Basic Channel - Базове використання каналу

func main() {
	// Створюємо канал
	messages := make(chan string)

	// Відправляємо в горутині
	go func() {
		messages <- "Hello"
		messages <- "World"
	}()

	// Читаємо з каналу
	msg1 := <-messages
	msg2 := <-messages

	fmt.Println(msg1)
	fmt.Println(msg2)
}
