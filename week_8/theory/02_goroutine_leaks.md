# Goroutine Leaks в Go

## 🎯 Що таке Goroutine Leak?

**Goroutine Leak** - це коли goroutine залишається "живою" але не робить корисної роботи, займаючи пам'ять та ресурси.

## ⚠️ Проблема

```go
func leak() {
    ch := make(chan int)
    
    go func() {
        val := <-ch  // ❌ Блокується назавжди!
        fmt.Println(val)
    }()
    
    // Channel ніколи не закривається
    // Goroutine чекає вічно
}  // ← Goroutine leaked!
```

### Чому це проблема?

```
1 leaked goroutine = ~2KB memory
10,000 leaked goroutines = ~20MB memory
100,000 leaked goroutines = ~200MB memory

+ Stack traces
+ Scheduler overhead
+ GC pressure
```

---

## 🐛 Типи Goroutine Leaks

### 1. Blocked on Channel Receive

```go
// ❌ BAD: Never sends
func leak() {
    ch := make(chan int)
    
    go func() {
        val := <-ch  // Blocks forever!
        process(val)
    }()
    
    // Forgot to send!
}
```

**Fix:**

```go
// ✅ GOOD: Send or close
func fixed() {
    ch := make(chan int)
    
    go func() {
        val := <-ch
        process(val)
    }()
    
    ch <- 42      // Send value
    // Or: close(ch)  // Close channel
}
```

### 2. Blocked on Channel Send

```go
// ❌ BAD: Unbuffered channel, no receiver
func leak() {
    ch := make(chan int)
    
    go func() {
        ch <- 42  // Blocks forever!
    }()
    
    // No receiver!
}
```

**Fix:**

```go
// ✅ GOOD: Buffered channel or receiver
func fixed() {
    ch := make(chan int, 1)  // Buffered
    
    go func() {
        ch <- 42  // Won't block
    }()
}

// ✅ BETTER: With receiver
func better() {
    ch := make(chan int)
    
    go func() {
        ch <- 42
    }()
    
    val := <-ch  // Receive
    fmt.Println(val)
}
```

### 3. Waiting on Mutex

```go
// ❌ BAD: Lock never released
func leak() {
    var mu sync.Mutex
    mu.Lock()
    
    go func() {
        mu.Lock()  // Blocks forever!
        defer mu.Unlock()
        // Do work
    }()
    
    // Never unlocks!
}
```

**Fix:**

```go
// ✅ GOOD: Always unlock
func fixed() {
    var mu sync.Mutex
    
    mu.Lock()
    // Do work
    mu.Unlock()  // ✅ Release
    
    go func() {
        mu.Lock()
        defer mu.Unlock()
        // Do work
    }()
}
```

### 4. Waiting on WaitGroup

```go
// ❌ BAD: Never calls Done()
func leak() {
    var wg sync.WaitGroup
    
    wg.Add(1)
    go func() {
        // Forgot defer wg.Done()!
        doWork()
    }()
    
    wg.Wait()  // Blocks forever!
}
```

**Fix:**

```go
// ✅ GOOD: Always call Done()
func fixed() {
    var wg sync.WaitGroup
    
    wg.Add(1)
    go func() {
        defer wg.Done()  // ✅ Always
        doWork()
    }()
    
    wg.Wait()
}
```

### 5. HTTP Request Without Timeout

```go
// ❌ BAD: No timeout
func leak() {
    go func() {
        resp, _ := http.Get("https://slow-server.com")
        // Може чекати вічно!
        defer resp.Body.Close()
    }()
}
```

**Fix:**

```go
// ✅ GOOD: With timeout
func fixed() {
    go func() {
        client := &http.Client{
            Timeout: 10 * time.Second,  // ✅ Timeout
        }
        
        resp, err := client.Get("https://slow-server.com")
        if err != nil {
            return
        }
        defer resp.Body.Close()
    }()
}
```

### 6. Context Not Canceled

```go
// ❌ BAD: Background context
func leak() {
    ctx := context.Background()  // Never cancels!
    
    go func() {
        <-ctx.Done()  // Blocks forever!
    }()
}
```

**Fix:**

```go
// ✅ GOOD: Cancelable context
func fixed() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()  // ✅ Cleanup
    
    go func() {
        <-ctx.Done()  // Will unblock!
    }()
    
    // Do work...
    cancel()  // Signal done
}
```

---

## 🔍 Як знайти Goroutine Leaks?

### 1. runtime.NumGoroutine()

```go
func TestForLeaks(t *testing.T) {
    before := runtime.NumGoroutine()
    
    // Run code that might leak
    suspiciousFunction()
    
    time.Sleep(100 * time.Millisecond)  // Give time to clean up
    
    after := runtime.NumGoroutine()
    
    if after > before {
        t.Errorf("Goroutine leak: %d -> %d", before, after)
    }
}
```

### 2. pprof

```go
import _ "net/http/pprof"

func main() {
    go func() {
        http.ListenAndServe("localhost:6060", nil)
    }()
    
    // Your app code
}
```

```bash
# Check goroutines
go tool pprof http://localhost:6060/debug/pprof/goroutine

# In pprof:
top      # Top functions creating goroutines
list     # Source code
traces   # Stack traces
```

### 3. goleak (Uber's Tool)

```go
import "go.uber.org/goleak"

func TestMain(m *testing.M) {
    goleak.VerifyTestMain(m)
}

func TestSomething(t *testing.T) {
    defer goleak.VerifyNone(t)
    
    // Test code
}
```

---

## ✅ Як запобігти Goroutine Leaks?

### 1. Завжди використовуй Context

```go
func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return  // ✅ Cleanup
        case <-time.After(1 * time.Second):
            doWork()
        }
    }
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    go worker(ctx)
    
    // Do work...
    cancel()  // Signal cleanup
}
```

### 2. Buffered Channels для "Fire and Forget"

```go
// ❌ BAD: Unbuffered
func bad() {
    ch := make(chan Result)
    
    go func() {
        result := compute()
        ch <- result  // Може блокуватись!
    }()
    
    // Якщо не читаємо з ch - leak!
}

// ✅ GOOD: Buffered
func good() {
    ch := make(chan Result, 1)  // Buffer = 1
    
    go func() {
        result := compute()
        ch <- result  // Won't block
    }()
    
    // Can ignore ch if needed
}
```

### 3. Select з Timeout

```go
func worker(ch chan int) {
    for {
        select {
        case val := <-ch:
            process(val)
        case <-time.After(5 * time.Minute):
            // Cleanup if idle too long
            return
        }
    }
}
```

### 4. Proper Channel Closing

```go
func producer(ch chan int) {
    defer close(ch)  // ✅ Always close
    
    for i := 0; i < 10; i++ {
        ch <- i
    }
}

func consumer(ch chan int) {
    for val := range ch {  // ✅ Range over channel
        process(val)
    }
    // Loop exits when channel closed
}
```

---

## 📊 Patterns для Clean Shutdown

### Pattern 1: Done Channel

```go
type Worker struct {
    done chan struct{}
}

func (w *Worker) Start() {
    w.done = make(chan struct{})
    
    go func() {
        for {
            select {
            case <-w.done:
                return  // ✅ Cleanup
            default:
                // Do work
            }
        }
    }()
}

func (w *Worker) Stop() {
    close(w.done)  // Signal shutdown
}
```

### Pattern 2: Context

```go
type Worker struct {
    ctx    context.Context
    cancel context.CancelFunc
}

func (w *Worker) Start() {
    w.ctx, w.cancel = context.WithCancel(context.Background())
    
    go func() {
        for {
            select {
            case <-w.ctx.Done():
                return  // ✅ Cleanup
            default:
                // Do work
            }
        }
    }()
}

func (w *Worker) Stop() {
    w.cancel()
}
```

### Pattern 3: WaitGroup + Done

```go
type Worker struct {
    wg   sync.WaitGroup
    done chan struct{}
}

func (w *Worker) Start(n int) {
    w.done = make(chan struct{})
    
    for i := 0; i < n; i++ {
        w.wg.Add(1)
        go func() {
            defer w.wg.Done()
            
            for {
                select {
                case <-w.done:
                    return  // ✅ Cleanup
                default:
                    // Do work
                }
            }
        }()
    }
}

func (w *Worker) Stop() {
    close(w.done)  // Signal all
    w.wg.Wait()    // Wait for cleanup
}
```

---

## 🎯 Best Practices

### 1. Goroutines Повинні Мати Exit Strategy

```go
// ❌ BAD: No way to stop
go func() {
    for {
        doWork()  // Forever!
    }
}()

// ✅ GOOD: Can be stopped
ctx, cancel := context.WithCancel(ctx)
go func() {
    for {
        select {
        case <-ctx.Done():
            return  // ✅ Exit
        default:
            doWork()
        }
    }
}()
```

### 2. Test для Leaks

```go
func TestNoLeak(t *testing.T) {
    defer goleak.VerifyNone(t)
    
    // Code that shouldn't leak
}
```

### 3. Monitor в Production

```go
import "github.com/prometheus/client_golang/prometheus"

var goroutineCount = prometheus.NewGauge(...)

func monitor() {
    ticker := time.NewTicker(10 * time.Second)
    for range ticker.C {
        count := runtime.NumGoroutine()
        goroutineCount.Set(float64(count))
        
        if count > 10000 {
            alert("Too many goroutines!")
        }
    }
}
```

### 4. Always Close Producers

```go
// Producer closes channel
func produce(ch chan int) {
    defer close(ch)  // ✅
    
    for i := 0; i < 10; i++ {
        ch <- i
    }
}

// Consumer ranges over channel
func consume(ch chan int) {
    for val := range ch {  // ✅ Exits when closed
        process(val)
    }
}
```

---

## 🔍 Debugging Goroutine Leaks

### 1. Dump All Goroutines

```go
import "runtime/pprof"

func dumpGoroutines() {
    pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
}
```

### 2. HTTP Endpoint

```bash
# Start pprof server
import _ "net/http/pprof"
http.ListenAndServe("localhost:6060", nil)

# Access goroutine profile
curl http://localhost:6060/debug/pprof/goroutine?debug=1
```

### 3. Analyze Stack Traces

```
goroutine 42 [chan receive]:
main.worker()
    /path/to/main.go:100 +0x50
created by main.startWorkers
    /path/to/main.go:50 +0x3c
```

Look for:
- `chan receive` / `chan send` - Blocked on channel
- `select` - Waiting in select
- `sync.Mutex.Lock` - Waiting on mutex

---

## 📈 Real-World Examples

### HTTP Server with Leaks

```go
// ❌ BAD: Leaks goroutines
func handleRequest(w http.ResponseWriter, r *http.Request) {
    ch := make(chan Result)
    
    go func() {
        result := slowOperation()
        ch <- result  // Може блокуватись!
    }()
    
    // If client disconnects, goroutine leaks!
}

// ✅ GOOD: Respect client context
func handleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    ch := make(chan Result, 1)  // Buffered
    
    go func() {
        result := slowOperation()
        select {
        case ch <- result:
        case <-ctx.Done():
            return  // Client disconnected
        }
    }()
    
    select {
    case result := <-ch:
        json.NewEncoder(w).Encode(result)
    case <-ctx.Done():
        http.Error(w, "Request canceled", 499)
    }
}
```

---

## 🎓 Висновок

### Goroutine Leak - це:

✅ Memory leak specific до Go  
✅ Goroutine блокується назавжди  
✅ Знаходиться через `runtime.NumGoroutine()`, pprof, goleak  
✅ Виправляється proper cleanup  

### Golden Rules:

1. **Every goroutine needs an exit strategy**
2. **Use context for cancellation**
3. **Always close channels (producer side)**
4. **Test for leaks (goleak)**
5. **Monitor goroutine count in production**

---

## 📖 Далі

- `practice/02_goroutine_leaks/` - Приклади з багами
- `03_debugging.md` - Debugging tools
- `go tool pprof` - Profiling

**"A goroutine starts cheap but runs forever expensive!" 🚀**
