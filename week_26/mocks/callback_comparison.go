package main

import (
	"context"
	"fmt"
	"time"
)

// Приклад 1: БЕЗ goroutine (СИНХРОННИЙ)
func (s *URLService) ProcessWithCallbackSync(
	ctx context.Context,
	urls []string,
	callback func([]string, error),
) {
	// ❌ НЕ використовуємо goroutine
	res, err := s.sorter.Sort(ctx, urls) // <- ТУТ БЛОКУЄМОСЬ!
	callback(res, err)
	// Функція повертається ТІЛЬКИ після виконання всієї роботи
}

// Приклад 2: З goroutine (АСИНХРОННИЙ)
func (s *URLService) ProcessWithCallbackAsync(
	ctx context.Context,
	urls []string,
	callback func([]string, error),
) {
	// ✅ Використовуємо goroutine
	go func() {
		res, err := s.sorter.Sort(ctx, urls)
		callback(res, err)
	}()
	// Функція повертається МИТТЄВО, робота продовжується в фоні
}

// ═══════════════════════════════════════════════════════════════════
// ДЕМОНСТРАЦІЯ РІЗНИЦІ
// ═══════════════════════════════════════════════════════════════════

func demoSyncVsAsync() {
	fmt.Println("╔═══════════════════════════════════════════════════╗")
	fmt.Println("║   ДЕМОНСТРАЦІЯ: Sync vs Async Callback           ║")
	fmt.Println("╚═══════════════════════════════════════════════════╝")
	fmt.Println()

	// Підготовка
	slowSorter := &SlowSorter{delay: 2 * time.Second}
	svc := NewURLService(slowSorter)
	urls := []string{"a.com", "b.com", "c.com"}

	// ─────────────────────────────────────────────────────────────
	// ТЕСТ 1: СИНХРОННИЙ ВИКЛИК (без goroutine)
	// ─────────────────────────────────────────────────────────────
	fmt.Println("🔴 ТЕСТ 1: Синхронний виклик (БЕЗ goroutine)")
	fmt.Println("─────────────────────────────────────────────────────")

	start1 := time.Now()
	fmt.Printf("[%s] 1. Викликаємо ProcessWithCallbackSync...\n", time.Now().Format("15:04:05.000"))

	svc.ProcessWithCallbackSync(context.Background(), urls, func(result []string, err error) {
		fmt.Printf("[%s] 3. ✓ Callback виконався! Result: %v\n", time.Now().Format("15:04:05.000"), result)
	})

	fmt.Printf("[%s] 2. ✗ Функція повернулась (але вже після всієї роботи!)\n", time.Now().Format("15:04:05.000"))
	fmt.Printf("⏱️  Загальний час: %v\n", time.Since(start1))
	fmt.Println()
	fmt.Println("❌ ПРОБЛЕМА: Основний потік БЛОКУВАВСЯ 2 секунди!")
	fmt.Println()

	// ─────────────────────────────────────────────────────────────
	// ТЕСТ 2: АСИНХРОННИЙ ВИКЛИК (з goroutine)
	// ─────────────────────────────────────────────────────────────
	fmt.Println("🟢 ТЕСТ 2: Асинхронний виклик (З goroutine)")
	fmt.Println("─────────────────────────────────────────────────────")

	start2 := time.Now()
	fmt.Printf("[%s] 1. Викликаємо ProcessWithCallbackAsync...\n", time.Now().Format("15:04:05.000"))

	done := make(chan bool)
	svc.ProcessWithCallbackAsync(context.Background(), urls, func(result []string, err error) {
		fmt.Printf("[%s] 4. ✓ Callback виконався! Result: %v\n", time.Now().Format("15:04:05.000"), result)
		done <- true
	})

	fmt.Printf("[%s] 2. ✓ Функція повернулась МИТТЄВО!\n", time.Now().Format("15:04:05.000"))
	fmt.Printf("[%s] 3. 💪 Можемо робити іншу роботу...\n", time.Now().Format("15:04:05.000"))

	// Симулюємо корисну роботу
	for i := 1; i <= 3; i++ {
		time.Sleep(500 * time.Millisecond)
		fmt.Printf("[%s]    → Виконую корисну роботу #%d\n", time.Now().Format("15:04:05.000"), i)
	}

	<-done // Чекаємо завершення
	fmt.Printf("⏱️  Загальний час: %v\n", time.Since(start2))
	fmt.Println()
	fmt.Println("✅ ПЕРЕВАГА: Основний потік НЕ блокувався, робота йшла паралельно!")
}

// ═══════════════════════════════════════════════════════════════════
// TIMELINE ВІЗУАЛІЗАЦІЯ
// ═══════════════════════════════════════════════════════════════════

func printTimeline() {
	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              TIMELINE: Sync vs Async                          ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	fmt.Println("🔴 БЕЗ goroutine (Синхронний):")
	fmt.Println("───────────────────────────────────────────────────────────────")
	fmt.Println("Main Thread:")
	fmt.Println("  │")
	fmt.Println("  ├─ Call ProcessWithCallbackSync()")
	fmt.Println("  │")
	fmt.Println("  ├─ Execute sorter.Sort()    ← БЛОКУЄТЬСЯ ТУТ!")
	fmt.Println("  │   ⏳ Чекаємо 2 секунди...")
	fmt.Println("  │   ⏳ Чекаємо 2 секунди...")
	fmt.Println("  │")
	fmt.Println("  ├─ Call callback()")
	fmt.Println("  │")
	fmt.Println("  └─ Return                   ← Повертається ПІСЛЯ 2 секунд")
	fmt.Println()
	fmt.Println("❌ Не можемо робити іншу роботу під час сортування!")
	fmt.Println()

	fmt.Println("🟢 З goroutine (Асинхронний):")
	fmt.Println("───────────────────────────────────────────────────────────────")
	fmt.Println("Main Thread                      Goroutine Thread")
	fmt.Println("────────────                     ────────────────")
	fmt.Println("  │")
	fmt.Println("  ├─ Call ProcessWithCallbackAsync()")
	fmt.Println("  │")
	fmt.Println("  ├─ Spawn goroutine          →  │")
	fmt.Println("  │                                ├─ Execute sorter.Sort()")
	fmt.Println("  └─ Return МИТТЄВО ✓             │   ⏳ Сортування...")
	fmt.Println("  │                                │   ⏳ Сортування...")
	fmt.Println("  ├─ Робимо іншу роботу #1 ✓      │")
	fmt.Println("  ├─ Робимо іншу роботу #2 ✓      │")
	fmt.Println("  ├─ Робимо іншу роботу #3 ✓      │")
	fmt.Println("  │                                ├─ Call callback()")
	fmt.Println("  │                                └─ Goroutine ends")
	fmt.Println("  │")
	fmt.Println("  └─ Continue...")
	fmt.Println()
	fmt.Println("✅ Можемо робити корисну роботу паралельно!")
	fmt.Println()
}

// ═══════════════════════════════════════════════════════════════════
// HELPER: SlowSorter для демонстрації
// ═══════════════════════════════════════════════════════════════════

type SlowSorter struct {
	delay time.Duration
}

func (s *SlowSorter) Sort(ctx context.Context, urls []string) ([]string, error) {
	// Симулюємо повільну операцію
	time.Sleep(s.delay)
	return urls, nil
}

// ═══════════════════════════════════════════════════════════════════
// ЧОМУ ПОТРІБНА goroutine?
// ═══════════════════════════════════════════════════════════════════

func explainWhy() {
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║          ЧОМУ ПОТРІБНА goroutine в Callback?                  ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	reasons := []struct {
		emoji  string
		reason string
		detail string
	}{
		{
			"⚡",
			"Асинхронність (Non-blocking)",
			"Функція повертається миттєво, не чекаючи завершення роботи.",
		},
		{
			"🔄",
			"Паралельність",
			"Основний потік може продовжувати роботу, поки callback виконується в фоні.",
		},
		{
			"📱",
			"Відповідність патерну",
			"Callback pattern означає \"викличи мене пізніше\", а не \"зачекай тут\".",
		},
		{
			"🎯",
			"Узгодженість з ProcessAsync()",
			"Обидва методи (channel і callback) працюють асинхронно.",
		},
		{
			"💪",
			"Кращий User Experience",
			"UI не зависає, сервер може обробляти інші запити.",
		},
	}

	for i, r := range reasons {
		fmt.Printf("%d. %s %s\n", i+1, r.emoji, r.reason)
		fmt.Printf("   └─ %s\n", r.detail)
		fmt.Println()
	}
}

// ═══════════════════════════════════════════════════════════════════
// РЕАЛЬНИЙ ПРИКЛАД
// ═══════════════════════════════════════════════════════════════════

func realWorldExample() {
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              РЕАЛЬНИЙ ПРИКЛАД (Web Server)                    ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	fmt.Println("🌐 HTTP Handler:")
	fmt.Println("───────────────────────────────────────────────────────────────")
	fmt.Println()

	fmt.Println("❌ БЕЗ goroutine:")
	fmt.Println("func HandleRequest(w http.ResponseWriter, r *http.Request) {")
	fmt.Println("    service.ProcessWithCallbackSync(ctx, urls, func(result, err) {")
	fmt.Println("        w.Write(result)  // Відповідь готова")
	fmt.Println("    })")
	fmt.Println("    // ← Тут БЛОКУЄТЬСЯ вся goroutine веб-сервера!")
	fmt.Println("    // ← Інші запити чекають!")
	fmt.Println("}")
	fmt.Println()
	fmt.Println("📊 Результат: Сервер може обробити 1 запит за раз")
	fmt.Println()

	fmt.Println("✅ З goroutine:")
	fmt.Println("func HandleRequest(w http.ResponseWriter, r *http.Request) {")
	fmt.Println("    service.ProcessWithCallbackAsync(ctx, urls, func(result, err) {")
	fmt.Println("        w.Write(result)  // Відповідь готова")
	fmt.Println("    })")
	fmt.Println("    // ← Одразу повернулись!")
	fmt.Println("    // ← Можемо обробляти наступний запит!")
	fmt.Println("}")
	fmt.Println()
	fmt.Println("📊 Результат: Сервер може обробити БАГАТО запитів паралельно")
	fmt.Println()
}

// ═══════════════════════════════════════════════════════════════════
// MAIN (для запуску демонстрації)
// ═══════════════════════════════════════════════════════════════════

/*
func main() {
	explainWhy()
	fmt.Println()
	printTimeline()
	fmt.Println()
	demoSyncVsAsync()
	fmt.Println()
	realWorldExample()
}
*/
