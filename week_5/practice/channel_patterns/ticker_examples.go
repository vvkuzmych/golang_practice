package main

import (
	"fmt"
	"time"
)

// ═══════════════════════════════════════════════════════════════════
// TIME TICKER З КАНАЛАМИ - ВСІ СПОСОБИ
// ═══════════════════════════════════════════════════════════════════

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════════╗")
	fmt.Println("║          ⏰ Time Ticker with Channels                           ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// 1. Basic Ticker
	example1_BasicTicker()

	// 2. Ticker with Stop
	example2_TickerWithStop()

	// 3. Ticker in Select
	example3_TickerInSelect()

	// 4. Multiple Tickers
	example4_MultipleTickers()

	// 5. Ticker vs Timer
	example5_TickerVsTimer()

	// 6. Rate Limiting with Ticker
	example6_RateLimiting()

	// 7. Ticker with Done Channel
	example7_TickerWithDone()

	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                     ✅ Completed!                                ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════╝")
}

// ═══════════════════════════════════════════════════════════════════
// 1. BASIC TICKER
// ═══════════════════════════════════════════════════════════════════

func example1_BasicTicker() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("1️⃣  BASIC TICKER")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Створення ticker, який "тікає" кожні 500ms
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop() // ⚠️ ВАЖЛИВО: завжди зупиняй ticker!

	fmt.Println("  Ticker created (500ms interval)")

	// ticker.C - це КАНАЛ (<-chan time.Time)
	// Він відправляє поточний час кожні 500ms

	count := 0
	for t := range ticker.C {
		count++
		fmt.Printf("  Tick %d at %s\n", count, t.Format("15:04:05.000"))

		if count >= 3 {
			break // Зупиняємо після 3 тіків
		}
	}

	fmt.Println("  ✅ Ticker stopped")
	fmt.Println()
}

// ═══════════════════════════════════════════════════════════════════
// 2. TICKER WITH STOP
// ═══════════════════════════════════════════════════════════════════

func example2_TickerWithStop() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("2️⃣  TICKER WITH STOP")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	ticker := time.NewTicker(200 * time.Millisecond)

	fmt.Println("  Starting ticker for 1 second...")

	// Запускаємо в горутині
	go func() {
		for t := range ticker.C {
			fmt.Printf("    Tick: %s\n", t.Format("15:04:05.000"))
		}
		fmt.Println("    Ticker channel closed")
	}()

	// Чекаємо 1 секунду
	time.Sleep(1 * time.Second)

	// Зупиняємо ticker
	ticker.Stop()
	fmt.Println("  ✅ Ticker stopped after 1 second")
	fmt.Println()

	time.Sleep(300 * time.Millisecond) // Дамо час закінчити
}

// ═══════════════════════════════════════════════════════════════════
// 3. TICKER IN SELECT
// ═══════════════════════════════════════════════════════════════════

func example3_TickerInSelect() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("3️⃣  TICKER IN SELECT")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	done := make(chan bool)
	work := make(chan string, 3)

	// Горутина відправляє роботу
	go func() {
		work <- "task 1"
		time.Sleep(500 * time.Millisecond)
		work <- "task 2"
		time.Sleep(500 * time.Millisecond)
		work <- "task 3"
		close(work)
		done <- true
	}()

	fmt.Println("  Waiting for work or ticker...")

	for {
		select {
		case task, ok := <-work:
			if !ok {
				fmt.Println("  All work done!")
				return
			}
			fmt.Printf("  📦 Got work: %s\n", task)

		case t := <-ticker.C:
			fmt.Printf("  ⏰ Tick at %s (still waiting...)\n", t.Format("15:04:05.000"))

		case <-done:
			fmt.Println("  ✅ Done channel received")
			fmt.Println()
			return
		}
	}
}

// ═══════════════════════════════════════════════════════════════════
// 4. MULTIPLE TICKERS
// ═══════════════════════════════════════════════════════════════════

func example4_MultipleTickers() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("4️⃣  MULTIPLE TICKERS")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Два ticker з різними інтервалами
	ticker1 := time.NewTicker(400 * time.Millisecond)
	ticker2 := time.NewTicker(700 * time.Millisecond)
	defer ticker1.Stop()
	defer ticker2.Stop()

	timeout := time.After(2 * time.Second)

	fmt.Println("  Two tickers with different intervals:")
	fmt.Println("    Ticker 1: 400ms")
	fmt.Println("    Ticker 2: 700ms")
	fmt.Println()

	for {
		select {
		case t := <-ticker1.C:
			fmt.Printf("  🟢 Ticker 1: %s\n", t.Format("15:04:05.000"))

		case t := <-ticker2.C:
			fmt.Printf("  🔵 Ticker 2: %s\n", t.Format("15:04:05.000"))

		case <-timeout:
			fmt.Println()
			fmt.Println("  ⏱️  Timeout! Stopping both tickers")
			fmt.Println()
			return
		}
	}
}

// ═══════════════════════════════════════════════════════════════════
// 5. TICKER VS TIMER
// ═══════════════════════════════════════════════════════════════════

func example5_TickerVsTimer() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("5️⃣  TICKER VS TIMER")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	fmt.Println("  Ticker - періодичний (кожні N мілісекунд):")
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	for i := 0; i < 3; i++ {
		t := <-ticker.C
		fmt.Printf("    Ticker tick %d: %s\n", i+1, t.Format("15:04:05.000"))
	}

	fmt.Println()
	fmt.Println("  Timer - одноразовий (після N мілісекунд):")

	timer := time.NewTimer(500 * time.Millisecond)
	defer timer.Stop()

	t := <-timer.C
	fmt.Printf("    Timer fired: %s\n", t.Format("15:04:05.000"))

	// Спроба прочитати знову - заблокується назавжди!
	// t = <-timer.C // ❌ deadlock!

	fmt.Println()
	fmt.Println("  Summary:")
	fmt.Println("    Ticker: 🔄 repeats (periodic)")
	fmt.Println("    Timer:  1️⃣  fires once (one-shot)")
	fmt.Println()
}

// ═══════════════════════════════════════════════════════════════════
// 6. RATE LIMITING WITH TICKER
// ═══════════════════════════════════════════════════════════════════

func example6_RateLimiting() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("6️⃣  RATE LIMITING WITH TICKER")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Обмеження: максимум 1 запит кожні 400ms
	limiter := time.NewTicker(400 * time.Millisecond)
	defer limiter.Stop()

	requests := []string{"req1", "req2", "req3", "req4", "req5"}

	fmt.Println("  Rate limiting: 1 request per 400ms")
	fmt.Println()

	for i, req := range requests {
		<-limiter.C // Чекаємо на tick (rate limit)
		fmt.Printf("  [%d] Processing: %s at %s\n",
			i+1, req, time.Now().Format("15:04:05.000"))
	}

	fmt.Println()
	fmt.Println("  ✅ All requests processed with rate limiting")
	fmt.Println()
}

// ═══════════════════════════════════════════════════════════════════
// 7. TICKER WITH DONE CHANNEL (Graceful Shutdown)
// ═══════════════════════════════════════════════════════════════════

func example7_TickerWithDone() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("7️⃣  TICKER WITH DONE CHANNEL (Graceful Shutdown)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	done := make(chan bool)

	// Worker з ticker
	go func() {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()

		count := 0
		for {
			select {
			case t := <-ticker.C:
				count++
				fmt.Printf("  [Worker] Tick %d at %s\n", count, t.Format("15:04:05.000"))

			case <-done:
				fmt.Println("  [Worker] Received shutdown signal, stopping...")
				return
			}
		}
	}()

	// Даємо попрацювати 1 секунду
	fmt.Println("  Worker started, will stop after 1 second...")
	time.Sleep(1 * time.Second)

	// Graceful shutdown
	fmt.Println()
	fmt.Println("  Sending shutdown signal...")
	done <- true

	// Даємо час закінчити
	time.Sleep(100 * time.Millisecond)
	fmt.Println("  ✅ Worker stopped gracefully")
	fmt.Println()
}

// ═══════════════════════════════════════════════════════════════════
// BONUS: Ticker Pattern - Periodic Task
// ═══════════════════════════════════════════════════════════════════

// PeriodicTask виконує task кожні interval мілісекунд
func PeriodicTask(interval time.Duration, task func(), done <-chan bool) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Виконати одразу
	task()

	for {
		select {
		case <-ticker.C:
			task()
		case <-done:
			return
		}
	}
}

// Приклад використання PeriodicTask
func examplePeriodicTask() {
	done := make(chan bool)

	count := 0
	task := func() {
		count++
		fmt.Printf("Task executed: %d\n", count)
	}

	go PeriodicTask(500*time.Millisecond, task, done)

	time.Sleep(2 * time.Second)
	done <- true
}
