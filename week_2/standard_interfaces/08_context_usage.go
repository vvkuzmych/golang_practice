package main

import (
	"context"
	"fmt"
	"time"
)

// ============= context.Context Interface =============

// type Context interface {
//     Deadline() (deadline time.Time, ok bool)
//     Done() <-chan struct{}
//     Err() error
//     Value(key interface{}) interface{}
// }

// ============= Simple Examples =============

// DoWork симулює тривалу роботу
func DoWork(ctx context.Context, name string, duration time.Duration) error {
	fmt.Printf("🔹 %s: starting work (%v)\n", name, duration)

	select {
	case <-time.After(duration):
		fmt.Printf("✅ %s: completed\n", name)
		return nil
	case <-ctx.Done():
		fmt.Printf("❌ %s: cancelled (%v)\n", name, ctx.Err())
		return ctx.Err()
	}
}

// FetchData симулює отримання даних
func FetchData(ctx context.Context, id int) (string, error) {
	// Симуляція затримки
	select {
	case <-time.After(2 * time.Second):
		return fmt.Sprintf("Data for ID %d", id), nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// ProcessRequest обробляє запит з таймаутом
func ProcessRequest(ctx context.Context, requestID string) error {
	fmt.Printf("📥 Processing request: %s\n", requestID)

	// Створюємо sub-context з таймаутом
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Симуляція роботи
	select {
	case <-time.After(2 * time.Second):
		fmt.Printf("✅ Request %s completed\n", requestID)
		return nil
	case <-ctx.Done():
		fmt.Printf("❌ Request %s timeout\n", requestID)
		return ctx.Err()
	}
}

// ============= Context with Values =============

type contextKey string

const (
	userIDKey    contextKey = "userID"
	requestIDKey contextKey = "requestID"
	traceIDKey   contextKey = "traceID"
)

// GetUserID отримує userID з context
func GetUserID(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(userIDKey).(int)
	return userID, ok
}

// GetRequestID отримує requestID з context
func GetRequestID(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(requestIDKey).(string)
	return requestID, ok
}

// AuthenticatedHandler обробляє запит з автентифікацією
func AuthenticatedHandler(ctx context.Context) {
	if userID, ok := GetUserID(ctx); ok {
		fmt.Printf("👤 User ID: %d\n", userID)
	} else {
		fmt.Println("❌ No user ID in context")
	}

	if reqID, ok := GetRequestID(ctx); ok {
		fmt.Printf("📋 Request ID: %s\n", reqID)
	}
}

// ============= Pipeline Example =============

// Stage1 перший етап pipeline
func Stage1(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for num := range in {
			select {
			case out <- num * 2:
				fmt.Printf("Stage1: %d -> %d\n", num, num*2)
			case <-ctx.Done():
				fmt.Println("Stage1: cancelled")
				return
			}
		}
	}()

	return out
}

// Stage2 другий етап pipeline
func Stage2(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for num := range in {
			select {
			case out <- num + 10:
				fmt.Printf("Stage2: %d -> %d\n", num, num+10)
			case <-ctx.Done():
				fmt.Println("Stage2: cancelled")
				return
			}
		}
	}()

	return out
}

// ============= Main =============

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║       context.Context Interface          ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	// ===== context.Background =====
	fmt.Println("\n🔹 context.Background")
	fmt.Println("─────────────────────────────────────────")

	ctx := context.Background()
	fmt.Printf("Context: %v\n", ctx)
	fmt.Println("✅ Базовий context (коренєвий)")

	// ===== context.WithCancel =====
	fmt.Println("\n🔹 context.WithCancel")
	fmt.Println("─────────────────────────────────────────")

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		DoWork(ctx, "Worker1", 5*time.Second)
	}()

	go func() {
		DoWork(ctx, "Worker2", 5*time.Second)
	}()

	// Скасувати через 2 секунди
	time.Sleep(2 * time.Second)
	fmt.Println("🛑 Cancelling context...")
	cancel()

	time.Sleep(1 * time.Second)

	// ===== context.WithTimeout =====
	fmt.Println("\n🔹 context.WithTimeout")
	fmt.Println("─────────────────────────────────────────")

	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go func() {
		DoWork(ctx, "ShortTask", 1*time.Second)
	}()

	go func() {
		DoWork(ctx, "LongTask", 5*time.Second)
	}()

	time.Sleep(4 * time.Second)

	// ===== context.WithDeadline =====
	fmt.Println("\n🔹 context.WithDeadline")
	fmt.Println("─────────────────────────────────────────")

	deadline := time.Now().Add(2 * time.Second)
	ctx, cancel = context.WithDeadline(context.Background(), deadline)
	defer cancel()

	fmt.Printf("Deadline: %s\n", deadline.Format("15:04:05"))

	err := DoWork(ctx, "DeadlineTask", 3*time.Second)
	if err != nil {
		fmt.Printf("Task error: %v\n", err)
	}

	time.Sleep(1 * time.Second)

	// ===== context.WithValue =====
	fmt.Println("\n🔹 context.WithValue")
	fmt.Println("─────────────────────────────────────────")

	ctx = context.Background()
	ctx = context.WithValue(ctx, userIDKey, 12345)
	ctx = context.WithValue(ctx, requestIDKey, "req-abc-123")
	ctx = context.WithValue(ctx, traceIDKey, "trace-xyz-789")

	AuthenticatedHandler(ctx)

	// ===== Nested Contexts =====
	fmt.Println("\n🔹 Nested Contexts")
	fmt.Println("─────────────────────────────────────────")

	parentCtx, parentCancel := context.WithCancel(context.Background())

	childCtx, childCancel := context.WithCancel(parentCtx)

	go func() {
		<-childCtx.Done()
		fmt.Println("Child context cancelled")
	}()

	go func() {
		<-parentCtx.Done()
		fmt.Println("Parent context cancelled")
	}()

	fmt.Println("Cancelling parent (дитина теж скасується)...")
	parentCancel()

	time.Sleep(100 * time.Millisecond)
	childCancel()

	// ===== Context in HTTP-style Request =====
	fmt.Println("\n🔹 HTTP-style Request Processing")
	fmt.Println("─────────────────────────────────────────")

	handleRequest := func(requestID string, timeout time.Duration) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, requestIDKey, requestID)

		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		ProcessRequest(ctx, requestID)
	}

	handleRequest("REQ-001", 5*time.Second)
	handleRequest("REQ-002", 1*time.Second) // timeout

	// ===== Pipeline with Context =====
	fmt.Println("\n🔹 Pipeline з Context")
	fmt.Println("─────────────────────────────────────────")

	pipelineCtx, pipelineCancel := context.WithCancel(context.Background())

	// Input channel
	input := make(chan int)

	// Pipeline stages
	stage1Out := Stage1(pipelineCtx, input)
	stage2Out := Stage2(pipelineCtx, stage1Out)

	// Producer
	go func() {
		for i := 1; i <= 5; i++ {
			select {
			case input <- i:
				time.Sleep(100 * time.Millisecond)
			case <-pipelineCtx.Done():
				close(input)
				return
			}
		}
		close(input)
	}()

	// Consumer
	go func() {
		for result := range stage2Out {
			fmt.Printf("Final result: %d\n", result)
		}
	}()

	// Дати час на обробку
	time.Sleep(300 * time.Millisecond)
	fmt.Println("🛑 Cancelling pipeline...")
	pipelineCancel()

	time.Sleep(500 * time.Millisecond)

	// ===== Multiple Operations =====
	fmt.Println("\n🔹 Multiple Operations з одним Context")
	fmt.Println("─────────────────────────────────────────")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Запустити кілька операцій
	done := make(chan bool)

	operations := []string{"DB Query", "API Call", "Cache Check"}
	for i, op := range operations {
		op := op
		dur := time.Duration(i+1) * time.Second
		go func() {
			DoWork(ctx, op, dur)
			done <- true
		}()
	}

	// Дочекатись всіх
	for i := 0; i < len(operations); i++ {
		<-done
	}

	// ===== Best Practices =====
	fmt.Println("\n🔹 Best Practices")
	fmt.Println("─────────────────────────────────────────")

	fmt.Println(`
✅ Добре:
  • Завжди передавати context як перший параметр
    func DoWork(ctx context.Context, data string) error
  
  • Завжди викликати cancel() функцію
    ctx, cancel := context.WithTimeout(...)
    defer cancel()
  
  • Використовувати ctx.Done() для скасування
    select {
    case <-ctx.Done():
        return ctx.Err()
    case result := <-ch:
        return result
    }

❌ Погано:
  • НЕ зберігати context в struct
    type Worker struct {
        ctx context.Context  // ❌ погано
    }
  
  • НЕ передавати nil context
    DoWork(nil, data)  // ❌ погано
    
  • НЕ ігнорувати ctx.Done()
    // просто чекаємо без перевірки ctx ❌
	`)

	// ===== Context Values Guidelines =====
	fmt.Println("\n🔹 Context Values Guidelines")
	fmt.Println("─────────────────────────────────────────")

	fmt.Println(`
✅ Добре використовувати для:
  • Request ID
  • Trace ID
  • User authentication info
  • Deadline/timeout info

❌ НЕ використовувати для:
  • Обов'язкових параметрів функції
  • Конфігурації
  • Dependency injection
  • Бізнес-логіки

💡 Правило: 
   Якщо без цього значення функція не може працювати
   → передавати як звичайний параметр, не через context
	`)

	// ===== Summary =====
	fmt.Println("\n\n📝 ВИСНОВКИ")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("✅ context.Context для:")
	fmt.Println("   • Скасування операцій")
	fmt.Println("   • Таймаути і дедлайни")
	fmt.Println("   • Передача метаданих (request ID, trace ID)")
	fmt.Println()
	fmt.Println("💡 Типи context:")
	fmt.Println("   • Background() - коренєвий context")
	fmt.Println("   • TODO() - placeholder")
	fmt.Println("   • WithCancel() - ручне скасування")
	fmt.Println("   • WithTimeout() - таймаут")
	fmt.Println("   • WithDeadline() - конкретний час")
	fmt.Println("   • WithValue() - передача даних")
	fmt.Println()
	fmt.Println("🔗 Context chain:")
	fmt.Println("   Parent context → Child context")
	fmt.Println("   Скасування parent скасовує всіх children")
	fmt.Println()
	fmt.Println("⚠️  Важливо:")
	fmt.Println("   • Перший параметр функції")
	fmt.Println("   • Завжди defer cancel()")
	fmt.Println("   • Не зберігати в struct")
	fmt.Println("   • Перевіряти ctx.Done()")
}
