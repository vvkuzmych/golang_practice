package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ============= Examples =============

func example1_BasicTimeout() {
	fmt.Println("1️⃣ Базовий Timeout")
	fmt.Println("─────────────────────────────────────────")

	// Створюємо context з timeout 2 секунди
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // Завжди викликаємо cancel!

	result := make(chan string, 1)

	// Запускаємо операцію
	go func() {
		time.Sleep(1 * time.Second) // Швидко - встигне!
		result <- "completed"
	}()

	// Чекаємо результат або timeout
	select {
	case <-ctx.Done():
		fmt.Println("❌ Timeout!")
	case res := <-result:
		fmt.Printf("✓ %s\n", res)
	}
	fmt.Println()
}

func example2_TimeoutExceeded() {
	fmt.Println("2️⃣ Timeout Exceeded")
	fmt.Println("─────────────────────────────────────────")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	result := make(chan string, 1)

	go func() {
		time.Sleep(3 * time.Second) // Занадто повільно!
		result <- "completed"
	}()

	select {
	case <-ctx.Done():
		fmt.Printf("❌ Timeout exceeded: %v\n", ctx.Err())
		// ctx.Err() поверне context.DeadlineExceeded
	case res := <-result:
		fmt.Printf("✓ %s\n", res)
	}
	fmt.Println()
}

func example3_WithDeadline() {
	fmt.Println("3️⃣ WithDeadline (фіксований час)")
	fmt.Println("─────────────────────────────────────────")

	// Deadline через 2 секунди від зараз
	deadline := time.Now().Add(2 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	fmt.Printf("Deadline: %s\n", deadline.Format("15:04:05"))
	fmt.Printf("Current:  %s\n", time.Now().Format("15:04:05"))

	time.Sleep(1 * time.Second)

	select {
	case <-ctx.Done():
		fmt.Println("❌ Deadline reached")
	default:
		fmt.Println("✓ Still running")
	}
	fmt.Println()
}

func example4_CheckingTimeout() {
	fmt.Println("4️⃣ Перевірка Timeout в операції")
	fmt.Println("─────────────────────────────────────────")

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	err := longOperation(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("❌ Operation timed out")
		} else {
			fmt.Printf("❌ Error: %v\n", err)
		}
	}
	fmt.Println()
}

func longOperation(ctx context.Context) error {
	for i := 0; i < 10; i++ {
		// Перевірка перед кожною ітерацією
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		fmt.Printf("  Step %d...\n", i+1)
		time.Sleep(200 * time.Millisecond)
	}
	return nil
}

func example5_PropagatingTimeout() {
	fmt.Println("5️⃣ Propagation Timeout")
	fmt.Println("─────────────────────────────────────────")

	// Parent context з timeout
	parentCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Викликаємо функцію яка передає context далі
	err := serviceA(parentCtx)
	if err != nil {
		fmt.Printf("❌ Service A failed: %v\n", err)
	}
	fmt.Println()
}

func serviceA(ctx context.Context) error {
	fmt.Println("→ Service A started")

	// Передаємо context далі
	err := serviceB(ctx)
	if err != nil {
		return fmt.Errorf("serviceA: %w", err)
	}

	fmt.Println("✓ Service A completed")
	return nil
}

func serviceB(ctx context.Context) error {
	fmt.Println("  → Service B started")

	// Повільна операція
	time.Sleep(3 * time.Second)

	// Перевірка cancellation
	select {
	case <-ctx.Done():
		return fmt.Errorf("serviceB: %w", ctx.Err())
	default:
	}

	fmt.Println("  ✓ Service B completed")
	return nil
}

func example6_MultipleOperations() {
	fmt.Println("6️⃣ Кілька операцій з одним timeout")
	fmt.Println("─────────────────────────────────────────")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	start := time.Now()

	// Операція 1
	if err := operation1(ctx); err != nil {
		fmt.Printf("Op1 failed: %v (after %v)\n", err, time.Since(start))
		return
	}

	// Операція 2
	if err := operation2(ctx); err != nil {
		fmt.Printf("Op2 failed: %v (after %v)\n", err, time.Since(start))
		return
	}

	// Операція 3
	if err := operation3(ctx); err != nil {
		fmt.Printf("Op3 failed: %v (after %v)\n", err, time.Since(start))
		return
	}

	fmt.Printf("✓ All operations completed in %v\n\n", time.Since(start))
}

func operation1(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(1 * time.Second):
		fmt.Println("  ✓ Operation 1 completed")
		return nil
	}
}

func operation2(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(1 * time.Second):
		fmt.Println("  ✓ Operation 2 completed")
		return nil
	}
}

func operation3(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(1 * time.Second):
		fmt.Println("  ✓ Operation 3 completed")
		return nil
	}
}

func example7_TimeoutRemaining() {
	fmt.Println("7️⃣ Перевірка залишку часу")
	fmt.Println("─────────────────────────────────────────")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Перевіряємо deadline
	deadline, ok := ctx.Deadline()
	if ok {
		remaining := time.Until(deadline)
		fmt.Printf("Time remaining: %v\n", remaining)
		fmt.Printf("Deadline: %s\n", deadline.Format("15:04:05"))
	}

	// Чекаємо трохи
	time.Sleep(2 * time.Second)

	remaining := time.Until(deadline)
	fmt.Printf("After 2s, remaining: %v\n\n", remaining)
}

// ============= Main =============

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║       Context Timeout Examples           ║")
	fmt.Println("╚══════════════════════════════════════════╝")
	fmt.Println()

	example1_BasicTimeout()
	example2_TimeoutExceeded()
	example3_WithDeadline()
	example4_CheckingTimeout()
	example5_PropagatingTimeout()
	example6_MultipleOperations()
	example7_TimeoutRemaining()

	fmt.Println("📝 Висновки:")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("✅ context.WithTimeout() для обмеження часу")
	fmt.Println("✅ Завжди defer cancel() після створення")
	fmt.Println("✅ Перевіряйте ctx.Done() в loops")
	fmt.Println("✅ errors.Is(err, context.DeadlineExceeded)")
	fmt.Println("✅ Propagate context через функції")
	fmt.Println("✅ Можна отримати deadline через ctx.Deadline()")
}
