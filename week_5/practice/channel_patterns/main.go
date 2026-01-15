package main

import (
	"fmt"
	"time"
)

// ============= Examples =============

func example1_UnbufferedChannel() {
	fmt.Println("1️⃣ Unbuffered Channel (synchronous)")
	fmt.Println("─────────────────────────────────────────")

	ch := make(chan string) // Unbuffered (capacity 0)

	// Receiver
	go func() {
		msg := <-ch // Блокується до sender
		fmt.Printf("   Received: %s\n", msg)
	}()

	fmt.Println("Sending...")
	ch <- "Hello" // Блокується до receiver
	fmt.Println("✓ Sent (receiver was ready)\n")
	time.Sleep(50 * time.Millisecond)
}

func example2_BufferedChannel() {
	fmt.Println("2️⃣ Buffered Channel (asynchronous)")
	fmt.Println("─────────────────────────────────────────")

	ch := make(chan int, 3) // Buffered (capacity 3)

	// Відправка без блокування
	fmt.Println("Sending 3 values...")
	ch <- 1
	fmt.Println("   Sent: 1")
	ch <- 2
	fmt.Println("   Sent: 2")
	ch <- 3
	fmt.Println("   Sent: 3")
	fmt.Println("✓ All sent (no blocking, buffer has space)")

	// Buffer FULL! Наступна відправка заблокує
	fmt.Println("\nReceiving...")
	fmt.Printf("   Received: %d\n", <-ch)
	fmt.Printf("   Received: %d\n", <-ch)
	fmt.Printf("   Received: %d\n", <-ch)
	fmt.Println()
}

func example3_CloseChannel() {
	fmt.Println("3️⃣ Close Channel")
	fmt.Println("─────────────────────────────────────────")

	ch := make(chan int, 2)

	// Відправка
	ch <- 1
	ch <- 2
	close(ch) // Закриваємо channel

	// Читання з закритого channel
	v1, ok1 := <-ch
	fmt.Printf("   Value: %d, OK: %t\n", v1, ok1) // 1, true

	v2, ok2 := <-ch
	fmt.Printf("   Value: %d, OK: %t\n", v2, ok2) // 2, true

	v3, ok3 := <-ch
	fmt.Printf("   Value: %d, OK: %t (closed, zero value)\n", v3, ok3) // 0, false

	fmt.Println("✓ Можна читати з закритого channel (отримаємо zero value)\n")
}

func example4_RangeOverChannel() {
	fmt.Println("4️⃣ Range over Channel")
	fmt.Println("─────────────────────────────────────────")

	ch := make(chan int, 5)

	// Producer
	go func() {
		for i := 1; i <= 5; i++ {
			ch <- i
		}
		close(ch) // ВАЖЛИВО: закриваємо після відправки всіх даних
	}()

	// Consumer: range автоматично завершується після close()
	for value := range ch {
		fmt.Printf("   Received: %d\n", value)
	}
	fmt.Println("✓ Range завершився після close()\n")
}

func example5_UnidirectionalChannels() {
	fmt.Println("5️⃣ Unidirectional Channels")
	fmt.Println("─────────────────────────────────────────")

	ch := make(chan int, 2)

	// Send-only channel
	sendOnly := func(ch chan<- int) {
		ch <- 10
		ch <- 20
		// value := <-ch  // ❌ Compilation error! Send-only
		close(ch)
	}

	// Receive-only channel
	receiveOnly := func(ch <-chan int) {
		for v := range ch {
			fmt.Printf("   Received: %d\n", v)
		}
		// ch <- 30  // ❌ Compilation error! Receive-only
	}

	go sendOnly(ch)
	receiveOnly(ch)

	fmt.Println("✓ Unidirectional channels забезпечують type safety\n")
}

func example6_SelectBasic() {
	fmt.Println("6️⃣ Select Statement (multiple channels)")
	fmt.Println("─────────────────────────────────────────")

	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "from channel 1"
	}()

	go func() {
		time.Sleep(50 * time.Millisecond)
		ch2 <- "from channel 2"
	}()

	// Select: виконає ПЕРШУ готову операцію
	select {
	case msg1 := <-ch1:
		fmt.Printf("   Received: %s\n", msg1)
	case msg2 := <-ch2:
		fmt.Printf("   Received: %s (faster!)\n", msg2)
	}
	fmt.Println()
}

func example7_SelectDefault() {
	fmt.Println("7️⃣ Select with Default (non-blocking)")
	fmt.Println("─────────────────────────────────────────")

	ch := make(chan int)

	// Non-blocking receive
	select {
	case v := <-ch:
		fmt.Printf("   Received: %d\n", v)
	default:
		fmt.Println("   No data available (non-blocking)")
	}

	// Non-blocking send
	select {
	case ch <- 42:
		fmt.Println("   Sent: 42")
	default:
		fmt.Println("   Cannot send (no receiver, non-blocking)")
	}
	fmt.Println()
}

func example8_SelectTimeout() {
	fmt.Println("8️⃣ Select with Timeout")
	fmt.Println("─────────────────────────────────────────")

	ch := make(chan string)

	// Slow operation
	go func() {
		time.Sleep(2 * time.Second)
		ch <- "data"
	}()

	// Timeout after 1 second
	select {
	case msg := <-ch:
		fmt.Printf("   Received: %s\n", msg)
	case <-time.After(1 * time.Second):
		fmt.Println("   ⏱️  Timeout! (1 second)")
	}
	fmt.Println()
}

func example9_NilChannel() {
	fmt.Println("9️⃣ Nil Channel (always blocks)")
	fmt.Println("─────────────────────────────────────────")

	var ch chan int // nil channel

	// ❌ Ці операції ЗАВЖДИ блокують:
	// ch <- 42    // Блокується назавжди
	// <-ch        // Блокується назавжди

	// ✅ Використання в select:
	select {
	case ch <- 42:
		fmt.Println("   Sent (never happens)")
	case <-ch:
		fmt.Println("   Received (never happens)")
	default:
		fmt.Println("   Nil channel always blocks (default executed)")
	}
	fmt.Println()
}

func example10_Pipeline() {
	fmt.Println("🔟 Pipeline Pattern")
	fmt.Println("─────────────────────────────────────────")

	// Generator
	generator := func() <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for i := 1; i <= 5; i++ {
				out <- i
			}
		}()
		return out
	}

	// Processor (square)
	square := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for n := range in {
				out <- n * n
			}
		}()
		return out
	}

	// Pipeline: generator → square → print
	nums := generator()
	squares := square(nums)

	for s := range squares {
		fmt.Printf("   %d\n", s)
	}
	fmt.Println("✓ Pipeline завершився\n")
}

func example11_FanOut() {
	fmt.Println("1️⃣1️⃣ Fan-Out Pattern (one producer, many workers)")
	fmt.Println("─────────────────────────────────────────")

	jobs := make(chan int, 10)
	results := make(chan int, 10)

	// Producer: один генератор
	go func() {
		for i := 1; i <= 10; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	// Workers: 3 паралельні обробники
	for w := 1; w <= 3; w++ {
		go func(id int) {
			for job := range jobs {
				fmt.Printf("   Worker %d: processing job %d\n", id, job)
				results <- job * 2
			}
		}(w)
	}

	// Collector
	go func() {
		time.Sleep(100 * time.Millisecond) // Даємо час workers
		close(results)
	}()

	fmt.Println("Results:")
	for result := range results {
		fmt.Printf("   → %d\n", result)
	}
	fmt.Println()
}

func example12_FanIn() {
	fmt.Println("1️⃣2️⃣ Fan-In Pattern (many producers, one consumer)")
	fmt.Println("─────────────────────────────────────────")

	// Функція для об'єднання каналів
	merge := func(channels ...<-chan int) <-chan int {
		out := make(chan int)
		for _, ch := range channels {
			ch := ch // Захоплюємо для goroutine
			go func() {
				for v := range ch {
					out <- v
				}
			}()
		}
		return out
	}

	// 3 producer channels
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	go func() {
		for i := 1; i <= 3; i++ {
			ch1 <- i
		}
		close(ch1)
	}()

	go func() {
		for i := 10; i <= 12; i++ {
			ch2 <- i
		}
		close(ch2)
	}()

	go func() {
		for i := 100; i <= 102; i++ {
			ch3 <- i
		}
		close(ch3)
	}()

	// Merge всі channels
	merged := merge(ch1, ch2, ch3)

	// Даємо час producers
	time.Sleep(100 * time.Millisecond)

	// Read з merged channel
	fmt.Println("Merged results:")
	timeout := time.After(200 * time.Millisecond)
	for {
		select {
		case v, ok := <-merged:
			if !ok {
				return
			}
			fmt.Printf("   → %d\n", v)
		case <-timeout:
			fmt.Println("✓ Timeout reached\n")
			return
		}
	}
}

// ============= Main =============

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║       Channel Patterns Examples          ║")
	fmt.Println("╚══════════════════════════════════════════╝")
	fmt.Println()

	example1_UnbufferedChannel()
	example2_BufferedChannel()
	example3_CloseChannel()
	example4_RangeOverChannel()
	example5_UnidirectionalChannels()
	example6_SelectBasic()
	example7_SelectDefault()
	example8_SelectTimeout()
	example9_NilChannel()
	example10_Pipeline()
	example11_FanOut()
	example12_FanIn()

	fmt.Println("📝 Висновки:")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("✅ Unbuffered - синхронний (блокує)")
	fmt.Println("✅ Buffered - асинхронний (до заповнення)")
	fmt.Println("✅ Close channel після відправки всіх даних")
	fmt.Println("✅ Range over channel для простого читання")
	fmt.Println("✅ Unidirectional для type safety")
	fmt.Println("✅ Select для роботи з кількома channels")
	fmt.Println("✅ Pipeline/Fan-Out/Fan-In patterns")
}
