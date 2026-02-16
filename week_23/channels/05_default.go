package main

import "fmt"

// 05. Default Case - Non-blocking select

func main() {
	messages := make(chan string)
	signals := make(chan bool)

	// Non-blocking receive
	select {
	case msg := <-messages:
		fmt.Println("Received message:", msg)
	default:
		fmt.Println("No message received")
	}

	// Non-blocking send
	msg := "hi"
	select {
	case messages <- msg:
		fmt.Println("Sent message:", msg)
	default:
		fmt.Println("No message sent")
	}

	// Multi-way non-blocking select
	select {
	case msg := <-messages:
		fmt.Println("Received:", msg)
	case sig := <-signals:
		fmt.Println("Received signal:", sig)
	default:
		fmt.Println("No activity")
	}
}
