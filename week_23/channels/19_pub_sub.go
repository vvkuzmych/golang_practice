package main

import (
	"fmt"
	"sync"
	"time"
)

// 19. Pub/Sub Pattern - Publish/Subscribe

type PubSub struct {
	mu          sync.RWMutex
	subscribers map[string][]chan string
}

func NewPubSub() *PubSub {
	return &PubSub{
		subscribers: make(map[string][]chan string),
	}
}

func (ps *PubSub) Subscribe(topic string) <-chan string {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan string, 10)
	ps.subscribers[topic] = append(ps.subscribers[topic], ch)
	return ch
}

func (ps *PubSub) Publish(topic string, msg string) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, ch := range ps.subscribers[topic] {
		go func(c chan string) {
			c <- msg
		}(ch)
	}
}

func main() {
	ps := NewPubSub()

	// Підписуємось на топіки
	news := ps.Subscribe("news")
	sports := ps.Subscribe("sports")

	// Subscriber 1
	go func() {
		for msg := range news {
			fmt.Println("News subscriber 1:", msg)
		}
	}()

	// Subscriber 2
	go func() {
		for msg := range sports {
			fmt.Println("Sports subscriber:", msg)
		}
	}()

	// Публікуємо повідомлення
	ps.Publish("news", "Breaking news!")
	ps.Publish("sports", "Team won!")
	ps.Publish("news", "Weather update")

	time.Sleep(500 * time.Millisecond)
}
