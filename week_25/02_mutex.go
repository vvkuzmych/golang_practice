package main

import (
	"fmt"
	"sync"
)

// 02. Mutex in Go - Synchronization

func main() {
	fmt.Println("=== Go Mutex - Synchronization ===")
	fmt.Println()

	// Problem: Race condition
	fmt.Println("1. Race condition (WITHOUT mutex):")
	counter := 0
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				counter++ // Race condition!
			}
		}()
	}

	wg.Wait()
	fmt.Printf("  Expected: 10000\n")
	fmt.Printf("  Got:      %d\n", counter) // Може бути менше!
	fmt.Println()

	// Solution: Mutex
	fmt.Println("2. Thread-safe (WITH mutex):")
	counter = 0
	var mu sync.Mutex

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				mu.Lock()
				counter++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	fmt.Printf("  Expected: 10000\n")
	fmt.Printf("  Got:      %d\n", counter) // Завжди 10000!
	fmt.Println()

	// RWMutex (Read-Write Mutex)
	fmt.Println("3. RWMutex (multiple readers, single writer):")
	var rwmu sync.RWMutex
	data := 0

	// Multiple readers
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			rwmu.RLock()
			fmt.Printf("  Reader %d: %d\n", n, data)
			rwmu.RUnlock()
		}(i)
	}

	// Single writer
	wg.Add(1)
	go func() {
		defer wg.Done()
		rwmu.Lock()
		data = 42
		fmt.Println("  Writer: updated to 42")
		rwmu.Unlock()
	}()

	wg.Wait()
	fmt.Println()

	// Real-world example: Bank account
	fmt.Println("4. Real-world: Bank account")
	account := NewBankAccount(1000)

	// Multiple goroutines accessing account
	for i := 0; i < 10; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			account.Deposit(100)
		}()
		go func() {
			defer wg.Done()
			account.Withdraw(50)
		}()
	}

	wg.Wait()
	fmt.Printf("  Final balance: $%d\n", account.Balance())
	fmt.Println()

	fmt.Println("✅ Go Mutex complete")
}

// BankAccount with mutex
type BankAccount struct {
	balance int
	mu      sync.Mutex
}

func NewBankAccount(balance int) *BankAccount {
	return &BankAccount{balance: balance}
}

func (a *BankAccount) Balance() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.balance
}

func (a *BankAccount) Deposit(amount int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.balance += amount
}

func (a *BankAccount) Withdraw(amount int) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.balance < amount {
		return false
	}
	a.balance -= amount
	return true
}

// Key points:
// - sync.Mutex for exclusive access
// - mu.Lock() / mu.Unlock()
// - defer mu.Unlock() pattern
// - sync.RWMutex for read-heavy workloads
// - No GIL - true parallelism
