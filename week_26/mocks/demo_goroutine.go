package main

import (
	"context"
	"fmt"
	"time"
)

// URLSorter interface (копія з main.go для незалежного запуску)
type URLSorter1 interface {
	Sort(ctx context.Context, urls []string) ([]string, error)
}

// Повільний сортер для демонстрації
type SlowSorter1 struct {
	delay time.Duration
}

func (s *SlowSorter1) Sort(ctx context.Context, urls []string) ([]string, error) {
	fmt.Printf("  [Goroutine] Початок сортування...\n")
	time.Sleep(s.delay) // Симулюємо повільну роботу
	fmt.Printf("  [Goroutine] Сортування завершено!\n")
	return urls, nil
}

// Сервіс з ДВОМА версіями callback - для порівняння
type DemoService struct {
	sorter URLSorter1
}

// ❌ ВЕРСІЯ 1: БЕЗ goroutine (СИНХРОННА)
func (s *DemoService) ProcessSYNC(ctx context.Context, urls []string, callback func([]string, error)) {
	fmt.Println("  [Main] Викликаю ProcessSYNC...")
	res, err := s.sorter.Sort(ctx, urls) // ← БЛОКУЄ тут!
	callback(res, err)
	fmt.Println("  [Main] ProcessSYNC повернувся")
}

// ✅ ВЕРСІЯ 2: З goroutine (АСИНХРОННА)
func (s *DemoService) ProcessASYNC(ctx context.Context, urls []string, callback func([]string, error)) {
	fmt.Println("  [Main] Викликаю ProcessASYNC...")
	go func() {
		res, err := s.sorter.Sort(ctx, urls)
		callback(res, err)
	}()
	fmt.Println("  [Main] ProcessASYNC повернувся МИТТЄВО!")
}

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║  ДЕМО: Чому потрібна goroutine в ProcessWithCallback?   ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()

	service := &DemoService{
		sorter: &SlowSorter1{delay: 2 * time.Second}, // 2 секунди затримка
	}
	urls := []string{"a.com", "b.com", "c.com"}
	ctx := context.Background()

	// ═══════════════════════════════════════════════════════════════
	// ТЕСТ 1: БЕЗ goroutine (СИНХРОННИЙ)
	// ═══════════════════════════════════════════════════════════════
	fmt.Println("🔴 ТЕСТ 1: БЕЗ goroutine (Синхронний)")
	fmt.Println("─────────────────────────────────────────────────────────")

	start1 := time.Now()
	fmt.Printf("[%s] 1. Основний потік: Починаю...\n", formatTime())

	service.ProcessSYNC(ctx, urls, func(result []string, err error) {
		fmt.Printf("[%s] 4. Callback виконався! Результат: %v\n", formatTime(), result)
	})

	fmt.Printf("[%s] 2. Основний потік: Функція повернулась\n", formatTime())
	fmt.Printf("[%s] 3. Основний потік: Спробую зробити іншу роботу...\n", formatTime())
	fmt.Println("        ...але вже пізно, все завершилось! 😢")
	fmt.Printf("\n⏱️  Загальний час: %v\n", time.Since(start1))

	fmt.Println("\n❌ ПРОБЛЕМА:")
	fmt.Println("   • Функція повернулась тільки через 2 секунди")
	fmt.Println("   • Основний потік був ЗАБЛОКОВАНИЙ весь цей час")
	fmt.Println("   • Не могли робити іншу роботу паралельно")
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	// Пауза перед наступним тестом
	time.Sleep(1 * time.Second)

	// ═══════════════════════════════════════════════════════════════
	// ТЕСТ 2: З goroutine (АСИНХРОННИЙ)
	// ═══════════════════════════════════════════════════════════════
	fmt.Println("🟢 ТЕСТ 2: З goroutine (Асинхронний)")
	fmt.Println("─────────────────────────────────────────────────────────")

	start2 := time.Now()
	fmt.Printf("[%s] 1. Основний потік: Починаю...\n", formatTime())

	done := make(chan bool)
	service.ProcessASYNC(ctx, urls, func(result []string, err error) {
		fmt.Printf("[%s] 6. Callback виконався! Результат: %v\n", formatTime(), result)
		done <- true
	})

	fmt.Printf("[%s] 2. Основний потік: Функція повернулась\n", formatTime())
	fmt.Printf("[%s] 3. Основний потік: Можу робити іншу роботу! 💪\n", formatTime())

	// Симулюємо корисну роботу
	for i := 1; i <= 3; i++ {
		time.Sleep(600 * time.Millisecond)
		fmt.Printf("[%s]    4.%d. Виконую корисну роботу #%d ✓\n", formatTime(), i, i)
	}

	fmt.Printf("[%s] 5. Основний потік: Чекаю callback...\n", formatTime())
	<-done // Чекаємо завершення

	fmt.Printf("\n⏱️  Загальний час: %v\n", time.Since(start2))

	fmt.Println("\n✅ ПЕРЕВАГИ:")
	fmt.Println("   • Функція повернулась МИТТЄВО (~0мс)")
	fmt.Println("   • Основний потік НЕ блокувався")
	fmt.Println("   • Змогли виконати іншу роботу ПАРАЛЕЛЬНО")
	fmt.Println("   • Callback викликався коли робота завершилась")
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	// ═══════════════════════════════════════════════════════════════
	// ВИСНОВОК
	// ═══════════════════════════════════════════════════════════════
	fmt.Println("💡 ВИСНОВОК:")
	fmt.Println("─────────────────────────────────────────────────────────")
	fmt.Println()
	fmt.Println("Goroutine в ProcessWithCallback потрібна для:")
	fmt.Println()
	fmt.Println("  1. ⚡ Асинхронності - функція не блокує")
	fmt.Println("  2. 🔄 Паралельності - можна робити іншу роботу")
	fmt.Println("  3. 🚀 Продуктивності - кращий UX і throughput")
	fmt.Println("  4. 📱 Семантики - callback = \"викличи мене пізніше\"")
	fmt.Println()
	fmt.Println("БЕЗ goroutine це була б синхронна функція, що блокує!")
	fmt.Println()
}

// Helper для форматування часу
func formatTime() string {
	return time.Now().Format("15:04:05.000")
}
