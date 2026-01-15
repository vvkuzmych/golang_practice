package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Демонстрація основних концепцій Тижня 4

func main() {
	fmt.Println("=== ТИЖДЕНЬ 4: Error Handling + Context ===\n")

	// 1. Базові помилки
	fmt.Println("1️⃣ Error Basics:")
	demoErrorBasics()

	// 2. Error Wrapping
	fmt.Println("\n2️⃣ Error Wrapping:")
	demoErrorWrapping()

	// 3. errors.Is / errors.As
	fmt.Println("\n3️⃣ errors.Is / errors.As:")
	demoErrorsIsAs()

	// 4. Context Timeout
	fmt.Println("\n4️⃣ Context Timeout:")
	demoContextTimeout()

	// 5. Context Cancellation
	fmt.Println("\n5️⃣ Context Cancellation:")
	demoContextCancellation()

	fmt.Println("\n✅ Все працює!")
}

// ===== 1. Error Basics =====

var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
)

func demoErrorBasics() {
	err := findUser(999)
	if err != nil {
		fmt.Printf("   Error: %v\n", err)

		// Перевірка sentinel error
		if errors.Is(err, ErrNotFound) {
			fmt.Println("   ✓ Detected: User not found")
		}
	}
}

func findUser(id int) error {
	// Симуляція: користувач не знайдений
	return ErrNotFound
}

// ===== 2. Error Wrapping =====

func demoErrorWrapping() {
	err := processUser(42)
	if err != nil {
		fmt.Printf("   Error chain: %v\n", err)

		// Unwrap для отримання оригінальної помилки
		if errors.Is(err, ErrNotFound) {
			fmt.Println("   ✓ Original error preserved in chain")
		}
	}
}

func processUser(id int) error {
	err := loadUserFromDB(id)
	if err != nil {
		// Wrapping з додаванням контексту
		return fmt.Errorf("failed to process user %d: %w", id, err)
	}
	return nil
}

func loadUserFromDB(id int) error {
	// Симуляція DB помилки
	return fmt.Errorf("database query failed: %w", ErrNotFound)
}

// ===== 3. errors.Is / errors.As =====

type ValidationError struct {
	Field string
	Value string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation failed: field '%s' has invalid value '%s'", e.Field, e.Value)
}

func demoErrorsIsAs() {
	err := validateEmail("invalid-email")
	if err != nil {
		fmt.Printf("   Error: %v\n", err)

		// errors.As для type assertion
		var valErr ValidationError
		if errors.As(err, &valErr) {
			fmt.Printf("   ✓ Validation error detected: field=%s\n", valErr.Field)
		}
	}
}

func validateEmail(email string) error {
	if email == "invalid-email" {
		return ValidationError{
			Field: "email",
			Value: email,
		}
	}
	return nil
}

// ===== 4. Context Timeout =====

func demoContextTimeout() {
	// Timeout через 2 секунди
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result := make(chan string, 1)

	// Запускаємо довгу операцію
	go func() {
		time.Sleep(500 * time.Millisecond) // Операція займе 500мс
		result <- "completed"
	}()

	// Чекаємо результат або timeout
	select {
	case <-ctx.Done():
		fmt.Println("   ✗ Operation timed out")
	case res := <-result:
		fmt.Printf("   ✓ Operation %s before timeout\n", res)
	}
}

// ===== 5. Context Cancellation =====

func demoContextCancellation() {
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan bool)

	// Запускаємо горутину
	go worker(ctx, done)

	// Через 1 секунду скасовуємо
	time.Sleep(1 * time.Second)
	fmt.Println("   Cancelling operation...")
	cancel()

	// Чекаємо завершення
	<-done
	fmt.Println("   ✓ Worker stopped gracefully")
}

func worker(ctx context.Context, done chan<- bool) {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Context скасований
			done <- true
			return
		case <-ticker.C:
			// Продовжуємо роботу
		}
	}
}
