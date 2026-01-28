package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// ✅ FIXED Solution 1: Mutex
type MutexCounter struct {
	mu    sync.Mutex
	value int
}

func (c *MutexCounter) Increment() {
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

func (c *MutexCounter) Get() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// ✅ FIXED Solution 2: Atomic
type AtomicCounter struct {
	value int64
}

func (c *AtomicCounter) Increment() {
	atomic.AddInt64(&c.value, 1)
}

func (c *AtomicCounter) Get() int64 {
	return atomic.LoadInt64(&c.value)
}

// ✅ FIXED Solution 3: Channel
type ChannelCounter struct {
	ch chan int
}

func NewChannelCounter() *ChannelCounter {
	c := &ChannelCounter{
		ch: make(chan int),
	}

	go func() {
		counter := 0
		for cmd := range c.ch {
			if cmd == 1 {
				counter++
			} else if cmd == 0 {
				c.ch <- counter
			}
		}
	}()

	return c
}

func (c *ChannelCounter) Increment() {
	c.ch <- 1
}

func (c *ChannelCounter) Get() int {
	c.ch <- 0
	return <-c.ch
}

func main() {
	fmt.Println("=== FIXED CODE (no race conditions) ===\n")

	// Solution 1: Mutex
	fmt.Println("Solution 1: Mutex")
	mutexCounter := &MutexCounter{}
	var wg1 sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg1.Add(1)
		go func() {
			defer wg1.Done()
			mutexCounter.Increment()
		}()
	}
	wg1.Wait()
	fmt.Printf("Mutex Counter: %d ✅\n\n", mutexCounter.Get())

	// Solution 2: Atomic
	fmt.Println("Solution 2: Atomic")
	atomicCounter := &AtomicCounter{}
	var wg2 sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			atomicCounter.Increment()
		}()
	}
	wg2.Wait()
	fmt.Printf("Atomic Counter: %d ✅\n\n", atomicCounter.Get())

	// Solution 3: Channel
	fmt.Println("Solution 3: Channel")
	channelCounter := NewChannelCounter()
	var wg3 sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg3.Add(1)
		go func() {
			defer wg3.Done()
			channelCounter.Increment()
		}()
	}
	wg3.Wait()
	fmt.Printf("Channel Counter: %d ✅\n", channelCounter.Get())
}
