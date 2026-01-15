package main

import (
	"context"
	"fmt"
	"time"
)

// ============= Examples =============

func example1_BasicCancellation() {
	fmt.Println("1️⃣ Базове Cancellation")
	fmt.Println("─────────────────────────────────────────")

	ctx, cancel := context.WithCancel(context.Background())

	go worker(ctx, "Worker-1")

	// Через 2 секунди скасовуємо
	time.Sleep(2 * time.Second)
	fmt.Println("→ Cancelling context...")
	cancel()

	// Даємо час worker завершитись
	time.Sleep(500 * time.Millisecond)
	fmt.Println()
}

func worker(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("✓ %s stopped gracefully\n", name)
			return
		default:
			fmt.Printf("  %s working...\n", name)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func example2_MultipleWorkers() {
	fmt.Println("2️⃣ Кілька Workers (cascading cancellation)")
	fmt.Println("─────────────────────────────────────────")

	ctx, cancel := context.WithCancel(context.Background())

	// Запускаємо 3 workers
	go worker(ctx, "Worker-1")
	go worker(ctx, "Worker-2")
	go worker(ctx, "Worker-3")

	// Через 1.5 секунди скасовуємо ВСЕ одразу
	time.Sleep(1500 * time.Millisecond)
	fmt.Println("→ Cancelling all workers...")
	cancel()

	time.Sleep(500 * time.Millisecond)
	fmt.Println()
}

func example3_ParentChildCancellation() {
	fmt.Println("3️⃣ Parent-Child Cancellation")
	fmt.Println("─────────────────────────────────────────")

	// Parent context
	parentCtx, parentCancel := context.WithCancel(context.Background())
	defer parentCancel()

	// Child context від parent
	childCtx, childCancel := context.WithCancel(parentCtx)
	defer childCancel()

	go func() {
		<-childCtx.Done()
		fmt.Println("✓ Child context cancelled")
	}()

	// Якщо скасувати parent - child теж скасується!
	time.Sleep(500 * time.Millisecond)
	fmt.Println("→ Cancelling parent context...")
	parentCancel()

	time.Sleep(200 * time.Millisecond)
	fmt.Println()
}

func example4_CancellationReason() {
	fmt.Println("4️⃣ Причина Cancellation")
	fmt.Println("─────────────────────────────────────────")

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-ctx.Done()
		// ctx.Err() поверне причину
		fmt.Printf("Context cancelled: %v\n", ctx.Err())
		// Output: "context canceled"
	}()

	time.Sleep(500 * time.Millisecond)
	cancel()

	time.Sleep(200 * time.Millisecond)
	fmt.Println()
}

func example5_CleanupOnCancel() {
	fmt.Println("5️⃣ Cleanup при Cancellation")
	fmt.Println("─────────────────────────────────────────")

	ctx, cancel := context.WithCancel(context.Background())

	go workerWithCleanup(ctx, "Worker")

	time.Sleep(1 * time.Second)
	fmt.Println("→ Cancelling...")
	cancel()

	time.Sleep(500 * time.Millisecond)
	fmt.Println()
}

func workerWithCleanup(ctx context.Context, name string) {
	// Відкриваємо "ресурс"
	fmt.Printf("%s: acquiring resources...\n", name)
	resource := "database connection"

	// Defer cleanup
	defer func() {
		fmt.Printf("%s: cleaning up %s\n", name, resource)
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s: context cancelled, stopping...\n", name)
			return
		default:
			fmt.Printf("%s: working with %s\n", name, resource)
			time.Sleep(300 * time.Millisecond)
		}
	}
}

func example6_LongRunningTask() {
	fmt.Println("6️⃣ Long-Running Task з Cancellation")
	fmt.Println("─────────────────────────────────────────")

	ctx, cancel := context.WithCancel(context.Background())

	// Симуляція користувача що чекає
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("→ User cancelled request")
		cancel()
	}()

	err := processLargeDataset(ctx)
	if err != nil {
		fmt.Printf("❌ Processing cancelled: %v\n", err)
	} else {
		fmt.Println("✓ Processing completed")
	}
	fmt.Println()
}

func processLargeDataset(ctx context.Context) error {
	items := 100

	for i := 1; i <= items; i++ {
		// Перевірка cancellation перед кожним item
		select {
		case <-ctx.Done():
			return fmt.Errorf("cancelled at item %d/%d: %w", i, items, ctx.Err())
		default:
		}

		// Обробка item
		if i%20 == 0 {
			fmt.Printf("  Processed %d/%d items\n", i, items)
		}
		time.Sleep(50 * time.Millisecond)
	}

	return nil
}

func example7_SelectWithDefault() {
	fmt.Println("7️⃣ Select з Default (non-blocking check)")
	fmt.Println("─────────────────────────────────────────")

	ctx, cancel := context.WithCancel(context.Background())

	// Перевірка без блокування
	select {
	case <-ctx.Done():
		fmt.Println("Context cancelled")
	default:
		fmt.Println("✓ Context still active")
	}

	// Скасовуємо
	cancel()

	// Перевірка знову
	select {
	case <-ctx.Done():
		fmt.Println("✓ Context cancelled (detected)")
	default:
		fmt.Println("Context still active")
	}
	fmt.Println()
}

func example8_ContextPropagation() {
	fmt.Println("8️⃣ Context Propagation через Layers")
	fmt.Println("─────────────────────────────────────────")

	ctx, cancel := context.WithCancel(context.Background())

	// Симулюємо cancellation після 1.5s
	go func() {
		time.Sleep(1500 * time.Millisecond)
		cancel()
	}()

	err := layerA(ctx)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	}
	fmt.Println()
}

func layerA(ctx context.Context) error {
	fmt.Println("→ Layer A")
	return layerB(ctx)
}

func layerB(ctx context.Context) error {
	fmt.Println("  → Layer B")
	return layerC(ctx)
}

func layerC(ctx context.Context) error {
	fmt.Println("    → Layer C (heavy work)")

	for i := 0; i < 5; i++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("layer C cancelled at step %d: %w", i, ctx.Err())
		default:
		}

		fmt.Printf("      Step %d...\n", i+1)
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

// ============= Main =============

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║     Context Cancellation Examples        ║")
	fmt.Println("╚══════════════════════════════════════════╝")
	fmt.Println()

	example1_BasicCancellation()
	example2_MultipleWorkers()
	example3_ParentChildCancellation()
	example4_CancellationReason()
	example5_CleanupOnCancel()
	example6_LongRunningTask()
	example7_SelectWithDefault()
	example8_ContextPropagation()

	fmt.Println("📝 Висновки:")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("✅ context.WithCancel() для manual control")
	fmt.Println("✅ cancel() скасовує всі дочірні contexts")
	fmt.Println("✅ Перевіряйте <-ctx.Done() в loops")
	fmt.Println("✅ ctx.Err() поверне context.Canceled")
	fmt.Println("✅ Використовуйте defer для cleanup")
	fmt.Println("✅ Parent cancel → child cancel автоматично")
}
