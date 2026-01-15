# Exercise 3: Graceful Shutdown

## 🎯 Мета

Створити HTTP-подібний сервіс з **graceful shutdown** - коректним завершенням роботи при отриманні сигналу (Ctrl+C).

---

## 📋 Завдання

Реалізуйте сервіс з наступними характеристиками:

### Функціональність:
- **Workers:** 3 workers обробляють jobs
- **Jobs:** Безкінечний потік jobs (генеруються кожні 500ms)
- **Shutdown:** При отриманні SIGINT (Ctrl+C):
  1. Зупинити прийом нових jobs
  2. Дочекатись завершення поточних jobs
  3. Cleanup (закрити channels, показати статистику)

### Вимоги:

- ✅ Signal handling для SIGINT та SIGTERM
- ✅ Context для керування lifecycle
- ✅ Multi-stage shutdown (stop new → finish existing → cleanup)
- ✅ Timeout на shutdown (максимум 5 секунд)
- ✅ Статистика: total jobs processed, total time running

---

## 💡 Підказки

### Signal Handling Setup:

```go
// Setup signal channel
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

// Wait for signal
go func() {
    sig := <-sigChan
    fmt.Printf("\nReceived signal: %v\n", sig)
    cancel() // Cancel context
}()
```

### Job Generator:

```go
func jobGenerator(ctx context.Context, jobs chan<- int) {
    jobID := 1
    ticker := time.NewTicker(500 * time.Millisecond)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            fmt.Println("Generator: stopping (no more jobs)")
            close(jobs)
            return
        case <-ticker.C:
            select {
            case jobs <- jobID:
                fmt.Printf("Generated job %d\n", jobID)
                jobID++
            case <-ctx.Done():
                close(jobs)
                return
            }
        }
    }
}
```

### Worker:

```go
func worker(id int, ctx context.Context, jobs <-chan int, wg *sync.WaitGroup) {
    defer wg.Done()
    
    for {
        select {
        case job, ok := <-jobs:
            if !ok {
                fmt.Printf("Worker %d: jobs channel closed, exiting\n", id)
                return
            }
            fmt.Printf("Worker %d: processing job %d\n", id, job)
            time.Sleep(1 * time.Second) // Simulate work
            fmt.Printf("Worker %d: finished job %d\n", id, job)
            
        case <-ctx.Done():
            fmt.Printf("Worker %d: context cancelled, finishing current job\n", id)
            // Закінчуємо поточний job, потім exit
            return
        }
    }
}
```

### Main Function Structure:

```go
func main() {
    // 1. Setup context and signal handling
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    // 2. Start components
    jobs := make(chan int, 10)
    var wg sync.WaitGroup
    
    go jobGenerator(ctx, jobs)
    
    for w := 1; w <= 3; w++ {
        wg.Add(1)
        go worker(w, ctx, jobs, &wg)
    }
    
    // 3. Wait for signal
    sig := <-sigChan
    fmt.Printf("\nReceived signal: %v\n", sig)
    fmt.Println("Initiating graceful shutdown...")
    
    // 4. Cancel context (stop accepting new jobs)
    cancel()
    
    // 5. Wait for workers with timeout
    done := make(chan bool)
    go func() {
        wg.Wait()
        done <- true
    }()
    
    select {
    case <-done:
        fmt.Println("Graceful shutdown completed")
    case <-time.After(5 * time.Second):
        fmt.Println("Shutdown timeout! Force exit...")
    }
    
    // 6. Print statistics
    fmt.Println("\n=== Statistics ===")
    // ... показати статистику
}
```

---

## 🎓 Ключові концепції

1. **Signal handling** - SIGINT, SIGTERM
2. **Context cancellation** - для координації shutdown
3. **WaitGroup** - очікування завершення workers
4. **Timeout pattern** - для force shutdown
5. **Multi-stage shutdown** - stop new → finish existing → cleanup

---

## ✅ Критерії успіху

- [ ] Сервіс запускається і генерує jobs кожні 500ms
- [ ] 3 workers обробляють jobs (по 1 секунді кожен)
- [ ] Ctrl+C коректно зупиняє сервіс
- [ ] Generator зупиняється одразу (no new jobs)
- [ ] Workers закінчують поточні jobs
- [ ] Показується статистика (total jobs, running time)
- [ ] Timeout працює (force exit після 5 секунд)

---

## 🚀 Очікуваний результат

```
Generated job 1
Worker 1: processing job 1
Generated job 2
Worker 2: processing job 2
Generated job 3
Worker 3: processing job 3
Worker 1: finished job 1
Generated job 4
Worker 1: processing job 4
^C
Received signal: interrupt
Initiating graceful shutdown...
Generator: stopping (no more jobs)
Worker 1: finishing current job
Worker 2: finishing current job
Worker 3: finishing current job
Worker 1: finished job 4
Worker 2: finished job 2
Worker 3: finished job 3
Worker 1: jobs channel closed, exiting
Worker 2: jobs channel closed, exiting
Worker 3: jobs channel closed, exiting
Graceful shutdown completed

=== Statistics ===
Total jobs processed: 4
Total running time: 5.2s
```

---

## 🔥 Бонус (опціонально)

### Бонус 1: HTTP Server
Замініть job generator на реальний HTTP сервер:

```go
srv := &http.Server{Addr: ":8080"}

// Graceful shutdown
<-sigChan
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
srv.Shutdown(ctx)
```

### Бонус 2: Multi-stage with Metrics
Додайте детальні метрики для кожного етапу:

```
Stage 1: Stop accepting new jobs (0.1s)
Stage 2: Finish existing jobs (3.2s)
Stage 3: Cleanup resources (0.3s)
Total shutdown time: 3.6s
```

### Бонус 3: Configurable Timeout
Зробіть timeout configurable через env variable або flag:

```go
shutdownTimeout := 5 * time.Second
if timeoutStr := os.Getenv("SHUTDOWN_TIMEOUT"); timeoutStr != "" {
    if d, err := time.ParseDuration(timeoutStr); err == nil {
        shutdownTimeout = d
    }
}
```

### Бонус 4: Handle Multiple Signals
Додайте обробку повторного сигналу (force exit):

```go
// First signal: graceful shutdown
// Second signal: force exit

signalCount := 0
for sig := range sigChan {
    signalCount++
    if signalCount == 1 {
        fmt.Println("Graceful shutdown initiated...")
        cancel()
    } else {
        fmt.Println("Force exit!")
        os.Exit(1)
    }
}
```

---

## 📚 Корисні посилання

- Theory: `week_5/theory/04_context_basics.md` - context usage
- Practice: `week_5/practice/graceful_shutdown/main.go` - приклади shutdown
- Solution: `week_5/solutions/solution_3.go` (після виконання)

---

## 🧪 Тестування

```bash
# Run the service
go run solution_3.go

# In another terminal, send signal:
# Ctrl+C or:
kill -SIGINT <pid>
kill -SIGTERM <pid>

# Test timeout (make workers slow, shutdown should timeout)
```

---

**Удачі! 🎉**

**Час виконання:** 60-90 хвилин

**Примітка:** Це найскладніше завдання тижня, але найкорисніше для реальних production систем!
