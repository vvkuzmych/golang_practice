package main

import (
	"fmt"
	"sync"
	"time"
)

// =============================================================================
// 1. Goroutine Basics
// =============================================================================

func example1_BasicGoroutine() {
	fmt.Println("1️⃣ Basic Goroutine")
	fmt.Println("─────────────────────────────────────────")

	// Звичайна функція (виконується синхронно)
	fmt.Println("Main: Before goroutine")

	// Запуск goroutine (виконується асинхронно)
	go func() {
		fmt.Println("   Goroutine: Hello from goroutine!")
	}()

	// Даємо час goroutine виконатись
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Main: After goroutine\n")
}

func example2_WaitGroup() {
	fmt.Println("2️⃣ WaitGroup для синхронізації")
	fmt.Println("─────────────────────────────────────────")

	var wg sync.WaitGroup

	// Запускаємо 3 goroutines
	for i := 1; i <= 3; i++ {
		wg.Add(1) // Збільшуємо counter
		go func(id int) {
			defer wg.Done() // Зменшуємо counter при завершенні
			fmt.Printf("   Worker %d: started\n", id)
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("   Worker %d: finished\n", id)
		}(i)
	}

	// Чекаємо поки всі goroutines завершаться
	wg.Wait()
	fmt.Println("✓ All workers finished\n")
}

// =============================================================================
// 2. Channels (Unbuffered vs Buffered)
// =============================================================================

func example3_UnbufferedChannel() {
	fmt.Println("3️⃣ Unbuffered Channel (блокує відправника)")
	fmt.Println("─────────────────────────────────────────")

	ch := make(chan string) // Unbuffered channel

	// Receiver goroutine
	go func() {
		msg := <-ch // Блокується до отримання
		fmt.Printf("   Received: %s\n", msg)
	}()

	// Sender (блокується до receiver)
	fmt.Println("Sending...")
	ch <- "Hello" // Блокується тут до receiver
	fmt.Println("✓ Sent!\n")
}

func example4_BufferedChannel() {
	fmt.Println("4️⃣ Buffered Channel (не блокує до заповнення)")
	fmt.Println("─────────────────────────────────────────")

	ch := make(chan int, 2) // Buffered channel (capacity 2)

	// Відправка без блокування (buffer не повний)
	fmt.Println("Sending 1...")
	ch <- 1
	fmt.Println("✓ Sent 1 (no blocking)")

	fmt.Println("Sending 2...")
	ch <- 2
	fmt.Println("✓ Sent 2 (no blocking)")

	// Buffer повний! Наступна відправка заблокує
	fmt.Println("Buffer is full (2/2)")

	// Отримання
	fmt.Printf("Received: %d\n", <-ch)
	fmt.Printf("Received: %d\n\n", <-ch)
}

func example5_CloseAndRange() {
	fmt.Println("5️⃣ Close & Range over Channel")
	fmt.Println("─────────────────────────────────────────")

	ch := make(chan int, 3)

	// Sender goroutine
	go func() {
		for i := 1; i <= 3; i++ {
			ch <- i
		}
		close(ch) // ВАЖЛИВО: закриваємо channel після відправки всіх даних
	}()

	// Receiver: range автоматично завершиться після close()
	for value := range ch {
		fmt.Printf("   Received: %d\n", value)
	}
	fmt.Println("✓ Channel closed, range finished\n")
}

// =============================================================================
// 3. Select Statement
// =============================================================================

func example6_SelectBasic() {
	fmt.Println("6️⃣ Select Statement (multiple channels)")
	fmt.Println("─────────────────────────────────────────")

	ch1 := make(chan string)
	ch2 := make(chan string)

	// Goroutine 1
	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "from ch1"
	}()

	// Goroutine 2
	go func() {
		time.Sleep(50 * time.Millisecond)
		ch2 <- "from ch2"
	}()

	// Select: чекає першої готової операції
	select {
	case msg1 := <-ch1:
		fmt.Printf("   Received: %s\n", msg1)
	case msg2 := <-ch2:
		fmt.Printf("   Received: %s\n", msg2) // Це виконається (ch2 швидший)
	}
	fmt.Println()
}

func example7_SelectWithTimeout() {
	fmt.Println("7️⃣ Select with Timeout")
	fmt.Println("─────────────────────────────────────────")

	ch := make(chan string)

	// Повільна goroutine (2 секунди)
	go func() {
		time.Sleep(2 * time.Second)
		ch <- "data"
	}()

	// Select з timeout (1 секунда)
	select {
	case msg := <-ch:
		fmt.Printf("   Received: %s\n", msg)
	case <-time.After(1 * time.Second):
		fmt.Println("   ⏱️  Timeout! (1 second elapsed)")
	}
	fmt.Println()
}

func example8_SelectWithDefault() {
	fmt.Println("8️⃣ Select with Default (non-blocking)")
	fmt.Println("─────────────────────────────────────────")

	ch := make(chan int)

	// Спроба прочитати без блокування
	select {
	case val := <-ch:
		fmt.Printf("   Received: %d\n", val)
	default:
		fmt.Println("   No data available (non-blocking)")
	}
	fmt.Println()
}

// =============================================================================
// 4. Deadlock Scenarios
// =============================================================================

func example9_DeadlockScenario1() {
	fmt.Println("9️⃣ Deadlock Scenario 1: Unbuffered channel без receiver")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("❌ Закоментовано! Розкоментуйте щоб побачити deadlock:")
	fmt.Println()
	fmt.Println("   ch := make(chan int)")
	fmt.Println("   ch <- 42  // DEADLOCK! Ніхто не читає")
	fmt.Println()

	// Uncomment to see deadlock:
	// ch := make(chan int)
	// ch <- 42  // DEADLOCK!
}

func example10_DeadlockScenario2() {
	fmt.Println("🔟 Deadlock Scenario 2: Забули close() в range")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("❌ Закоментовано! Розкоментуйте щоб побачити deadlock:")
	fmt.Println()
	fmt.Println("   ch := make(chan int, 1)")
	fmt.Println("   ch <- 1")
	fmt.Println("   // Забули close(ch)!")
	fmt.Println("   for v := range ch {")
	fmt.Println("       fmt.Println(v)")
	fmt.Println("   } // DEADLOCK! Range чекає на close()")
	fmt.Println()

	// Uncomment to see deadlock:
	// ch := make(chan int, 1)
	// ch <- 1
	// // Забули close(ch)!
	// for v := range ch {
	// 	fmt.Println(v)
	// } // DEADLOCK!
}

// =============================================================================
// 5. Channel vs Queue
// =============================================================================

func example11_ChannelVsQueue() {
	fmt.Println("1️⃣1️⃣ Channel vs Queue - Ключові різниці")
	fmt.Println("─────────────────────────────────────────")

	fmt.Println("📦 CHANNEL:")
	fmt.Println("   • Призначення: синхронізація та комунікація")
	fmt.Println("   • Блокування: блокуючий (by design)")
	fmt.Println("   • Ownership: shared communication")
	fmt.Println("   • Буфер: оптимізація, не основна ціль")
	fmt.Println()

	fmt.Println("📚 QUEUE:")
	fmt.Println("   • Призначення: зберігання даних")
	fmt.Println("   • Блокування: non-blocking (зазвичай)")
	fmt.Println("   • Ownership: shared state")
	fmt.Println("   • Буфер: основна функція")
	fmt.Println()

	fmt.Println("🎯 ВИСНОВОК:")
	fmt.Println("   Channel — це інструмент для COMMUNICATION,")
	fmt.Println("   а не для DATA STORAGE!")
	fmt.Println()
}

// =============================================================================
// 6. Worker Pool Pattern
// =============================================================================

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("   Worker %d: processing job %d\n", id, job)
		time.Sleep(100 * time.Millisecond) // Симуляція роботи
		results <- job * 2
	}
	fmt.Printf("   Worker %d: finished\n", id)
}

func example12_WorkerPool() {
	fmt.Println("1️⃣2️⃣ Worker Pool Pattern")
	fmt.Println("─────────────────────────────────────────")

	const numJobs = 5
	const numWorkers = 2

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	var wg sync.WaitGroup

	// Запускаємо workers
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// Відправляємо jobs
	fmt.Println("Sending jobs...")
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs) // Закриваємо jobs channel

	// Чекаємо завершення workers
	wg.Wait()
	close(results) // Закриваємо results channel

	// Збираємо results
	fmt.Println("\nResults:")
	for result := range results {
		fmt.Printf("   Result: %d\n", result)
	}
	fmt.Println()
}

// =============================================================================
// Main Function
// =============================================================================

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║   ТИЖДЕНЬ 5: Goroutines & Channels       ║")
	fmt.Println("╚══════════════════════════════════════════╝")
	fmt.Println()

	// 1. Goroutine Basics
	example1_BasicGoroutine()
	example2_WaitGroup()

	// 2. Channels
	example3_UnbufferedChannel()
	example4_BufferedChannel()
	example5_CloseAndRange()

	// 3. Select
	example6_SelectBasic()
	example7_SelectWithTimeout()
	example8_SelectWithDefault()

	// 4. Deadlock Scenarios (commented out)
	example9_DeadlockScenario1()
	example10_DeadlockScenario2()

	// 5. Channel vs Queue
	example11_ChannelVsQueue()

	// 6. Worker Pool
	example12_WorkerPool()

	fmt.Println("📝 Висновки:")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("✅ Goroutines - легкі та швидкі")
	fmt.Println("✅ WaitGroup для синхронізації")
	fmt.Println("✅ Unbuffered channel - блокує sender")
	fmt.Println("✅ Buffered channel - буферизує до capacity")
	fmt.Println("✅ Select для роботи з кількома channels")
	fmt.Println("✅ Завжди close() channel після відправки всіх даних")
	fmt.Println("✅ Channel для communication, Queue для storage")
	fmt.Println("✅ Worker pool для bounded concurrency")
}
