package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// ============= Examples =============

// Example 1: Basic Signal Handling
func example1_BasicSignalHandling() {
	fmt.Println("1️⃣ Basic Signal Handling")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("Press Ctrl+C to stop...")

	// Create signal channel
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for signal
	sig := <-sigChan
	fmt.Printf("\n✓ Received signal: %v\n", sig)
	fmt.Println("✓ Shutting down gracefully...\n")
}

// Example 2: Worker with Done Channel
func example2_WorkerWithDoneChannel() {
	fmt.Println("2️⃣ Worker with Done Channel")
	fmt.Println("─────────────────────────────────────────")

	done := make(chan bool)
	finished := make(chan bool)

	// Worker
	go func() {
		for {
			select {
			case <-done:
				fmt.Println("   Worker: received shutdown signal")
				fmt.Println("   Worker: cleaning up...")
				time.Sleep(100 * time.Millisecond)
				fmt.Println("   Worker: shutdown complete")
				finished <- true
				return
			default:
				fmt.Println("   Worker: working...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	// Simulate running
	time.Sleep(2 * time.Second)

	// Signal shutdown
	fmt.Println("\nMain: sending shutdown signal...")
	close(done)

	// Wait for worker to finish
	<-finished
	fmt.Println("✓ All workers stopped\n")
}

// Example 3: Context-based Shutdown
func example3_ContextBasedShutdown() {
	fmt.Println("3️⃣ Context-based Shutdown")
	fmt.Println("─────────────────────────────────────────")

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Worker function
	worker := func(id int, ctx context.Context) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("   Worker %d: context cancelled, stopping...\n", id)
				return
			default:
				fmt.Printf("   Worker %d: working...\n", id)
				time.Sleep(500 * time.Millisecond)
			}
		}
	}

	// Start workers
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i, ctx)
	}

	// Simulate running
	time.Sleep(2 * time.Second)

	// Cancel context (shutdown signal)
	fmt.Println("\nMain: cancelling context...")
	cancel()

	// Wait for all workers
	wg.Wait()
	fmt.Println("✓ All workers stopped\n")
}

// Example 4: Graceful Shutdown with Timeout
func example4_ShutdownWithTimeout() {
	fmt.Println("4️⃣ Graceful Shutdown with Timeout")
	fmt.Println("─────────────────────────────────────────")

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Slow worker (simulates long cleanup)
	slowWorker := func(id int, ctx context.Context) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("   Worker %d: starting cleanup (takes 3s)...\n", id)
				time.Sleep(3 * time.Second)
				fmt.Printf("   Worker %d: cleanup done\n", id)
				return
			default:
				time.Sleep(500 * time.Millisecond)
			}
		}
	}

	wg.Add(1)
	go slowWorker(1, ctx)

	time.Sleep(1 * time.Second)

	// Cancel with timeout
	fmt.Println("\nMain: initiating shutdown (timeout: 2s)...")
	cancel()

	// Wait with timeout
	shutdownComplete := make(chan bool)
	go func() {
		wg.Wait()
		shutdownComplete <- true
	}()

	select {
	case <-shutdownComplete:
		fmt.Println("✓ Shutdown completed in time")
	case <-time.After(2 * time.Second):
		fmt.Println("⚠️  Shutdown timeout! Force exit...")
	}
	fmt.Println()
}

// Example 5: Multi-stage Shutdown
func example5_MultiStageShutdown() {
	fmt.Println("5️⃣ Multi-stage Shutdown")
	fmt.Println("─────────────────────────────────────────")

	var wg sync.WaitGroup
	stopAcceptingJobs := make(chan bool)
	shutdownWorkers := make(chan bool)

	// Job producer
	jobs := make(chan int, 10)
	go func() {
		i := 1
		for {
			select {
			case <-stopAcceptingJobs:
				fmt.Println("   Producer: stopped accepting new jobs")
				close(jobs)
				return
			default:
				jobs <- i
				fmt.Printf("   Producer: job %d created\n", i)
				i++
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()

	// Worker
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case job, ok := <-jobs:
				if !ok {
					fmt.Println("   Worker: no more jobs, finishing...")
					return
				}
				fmt.Printf("   Worker: processing job %d\n", job)
				time.Sleep(200 * time.Millisecond)
			case <-shutdownWorkers:
				fmt.Println("   Worker: emergency shutdown!")
				return
			}
		}
	}()

	// Simulate running
	time.Sleep(2 * time.Second)

	// Stage 1: Stop accepting new jobs
	fmt.Println("\n🛑 Stage 1: Stop accepting new jobs...")
	close(stopAcceptingJobs)

	time.Sleep(1 * time.Second)

	// Stage 2: Wait for workers to finish existing jobs (or timeout)
	fmt.Println("\n🛑 Stage 2: Waiting for workers to finish existing jobs...")
	workersDone := make(chan bool)
	go func() {
		wg.Wait()
		workersDone <- true
	}()

	select {
	case <-workersDone:
		fmt.Println("✓ All workers finished gracefully")
	case <-time.After(3 * time.Second):
		fmt.Println("⚠️  Workers timeout! Forcing shutdown...")
		close(shutdownWorkers)
		wg.Wait()
	}
	fmt.Println()
}

// Example 6: Signal Handling with Graceful Shutdown (DEMO mode)
func example6_SignalHandlingDemo() {
	fmt.Println("6️⃣ Signal Handling with Graceful Shutdown (DEMO)")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("⚠️  In real app, run and press Ctrl+C")
	fmt.Println("⚠️  DEMO mode: simulating signal after 2s...\n")

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Worker
	worker := func(id int) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("   Worker %d: shutting down...\n", id)
				time.Sleep(100 * time.Millisecond) // Cleanup
				fmt.Printf("   Worker %d: done\n", id)
				return
			default:
				fmt.Printf("   Worker %d: working...\n", id)
				time.Sleep(500 * time.Millisecond)
			}
		}
	}

	// Start workers
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go worker(i)
	}

	// DEMO: Simulate signal after 2 seconds
	go func() {
		time.Sleep(2 * time.Second)
		sigChan <- syscall.SIGINT
	}()

	// Wait for signal
	sig := <-sigChan
	fmt.Printf("\n🛑 Received signal: %v\n", sig)
	fmt.Println("🛑 Initiating graceful shutdown...")

	// Cancel context
	cancel()

	// Wait for workers with timeout
	done := make(chan bool)
	go func() {
		wg.Wait()
		done <- true
	}()

	select {
	case <-done:
		fmt.Println("✓ Graceful shutdown completed")
	case <-time.After(5 * time.Second):
		fmt.Println("⚠️  Shutdown timeout! Force exit...")
	}
	fmt.Println()
}

// ============= Main =============

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║      Graceful Shutdown Patterns          ║")
	fmt.Println("╚══════════════════════════════════════════╝")
	fmt.Println()

	// Uncomment to test real signal handling
	// example1_BasicSignalHandling()
	// return

	example2_WorkerWithDoneChannel()
	example3_ContextBasedShutdown()
	example4_ShutdownWithTimeout()
	example5_MultiStageShutdown()
	example6_SignalHandlingDemo()

	fmt.Println("📝 Висновки:")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("✅ Signal handling: SIGINT (Ctrl+C), SIGTERM")
	fmt.Println("✅ Done channel для простих cases")
	fmt.Println("✅ Context для складних scenarios")
	fmt.Println("✅ Timeout для force shutdown")
	fmt.Println("✅ Multi-stage: stop new → finish existing → force")
	fmt.Println("✅ WaitGroup для очікування всіх workers")
}
