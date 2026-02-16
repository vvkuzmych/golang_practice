package main

import "fmt"

// 03. Channel Directions - Напрямки каналів

// Функція приймає канал тільки для відправки
func sender(ch chan<- string) {
	ch <- "Hello"
	// ch <- "World" // Не можемо читати: invalid operation
}

// Функція приймає канал тільки для читання
func receiver(ch <-chan string) {
	msg := <-ch
	fmt.Println(msg)
	// ch <- "test" // Не можемо відправляти
}

func main() {
	messages := make(chan string, 1)

	sender(messages)
	receiver(messages)
}
