# Week 9: Concurrency Patterns

## 🎯 Мета

Опанувати **3 головні concurrency patterns** для Go: Worker Pool, Fan-In/Fan-Out, та Pipeline.

---

## 📚 Теорія

### 1. Worker Pool
**Файл:** `theory/01_worker_pool.md`

- Що таке Worker Pool?
- Controlled concurrency
- Rate limiting
- Context для cancellation
- Real-world use cases

**Pattern:**
```
Jobs Queue → [Worker 1]
           → [Worker 2]
           → [Worker 3]
                ↓
            Results
```

### 2. Fan-In / Fan-Out
**Файл:** `theory/02_fan_in_fan_out.md`

- Fan-Out: розподіл роботи
- Fan-In: збір результатів
- Multiplexing з select
- Bounded concurrency
- Priority queues

**Pattern:**
```
        Fan-Out
Input → [Worker 1] →
      → [Worker 2] →  Fan-In → Output
      → [Worker 3] →
```

### 3. Pipeline
**Файл:** `theory/03_pipeline.md`

- Series of stages
- Connected by channels
- Composable architecture
- Error handling
- Buffering strategies

**Pattern:**
```
Input → [Stage 1] → [Stage 2] → [Stage 3] → Output
```

---

## 💻 Практика

### 1. Worker Pool
**Директорія:** `practice/01_worker_pool/`

**Файл:** `simple_worker_pool.go`

**3 приклади:**
1. Simple worker pool
2. Worker pool з timeout
3. Dynamic cancellation

**Як запускати:**
```bash
cd practice/01_worker_pool
go run simple_worker_pool.go
```

### 2. Fan-In / Fan-Out
**Директорія:** `practice/02_fan_in_fan_out/`

**Файл:** `fan_pattern.go`

**2 приклади:**
1. Basic fan-out/fan-in
2. Large workload processing

**Як запускати:**
```bash
cd practice/02_fan_in_fan_out
go run fan_pattern.go
```

### 3. Pipeline з Context
**Директорія:** `practice/03_pipeline/`

**Файл:** `three_stage_pipeline.go`

**3 приклади:**
1. Normal completion
2. Timeout cancellation
3. Manual cancellation

**Pipeline stages:**
- Stage 1: Generate numbers
- Stage 2: Square numbers
- Stage 3: Filter even numbers

**Як запускати:**
```bash
cd practice/03_pipeline
go run three_stage_pipeline.go
```

---

## 🎯 Key Patterns Summary

### Worker Pool

```go
// Fixed workers, job queue
numWorkers := 5
jobs := make(chan Job, 100)
results := make(chan Result, 100)

for w := 1; w <= numWorkers; w++ {
    go worker(w, jobs, results)
}
```

**Use cases:**
- HTTP request processing
- Database batch operations
- Image processing
- File processing

### Fan-Out / Fan-In

```go
// Fan-Out: Split work
workers := fanOut(input, numWorkers)

// Fan-In: Merge results
results := fanIn(workers...)
```

**Use cases:**
- Parallel processing
- Distributed computation
- Multi-source aggregation

### Pipeline

```go
// Connect stages
stage1 := generator(data)
stage2 := transform(stage1)
stage3 := filter(stage2)
output := collect(stage3)
```

**Use cases:**
- ETL pipelines
- Log processing
- Stream processing
- Data transformation

---

## 🔧 Context Usage

### Timeout

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

stage := processWithContext(ctx, input)
```

### Manual Cancellation

```go
ctx, cancel := context.WithCancel(context.Background())

go func() {
    // Cancel on some condition
    if shouldCancel() {
        cancel()
    }
}()

stage := processWithContext(ctx, input)
```

### Graceful Shutdown

```go
ctx, cancel := context.WithCancel(context.Background())

// Handle signals
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, os.Interrupt)

go func() {
    <-sigChan
    cancel()
}()
```

---

## ✅ Best Practices

### 1. Always Close Channels

```go
go func() {
    defer close(out)  // ✅ Producer closes
    for _, val := range data {
        out <- val
    }
}()
```

### 2. Use Context для Cancellation

```go
select {
case <-ctx.Done():
    return  // ✅ Respect cancellation
case val := <-input:
    process(val)
}
```

### 3. Proper WaitGroup Usage

```go
var wg sync.WaitGroup
wg.Add(numWorkers)  // ✅ Before starting goroutines

for i := 0; i < numWorkers; i++ {
    go func() {
        defer wg.Done()
        // Work
    }()
}

wg.Wait()
```

### 4. Buffer Channels для Performance

```go
// For high throughput
jobs := make(chan Job, 100)      // ✅ Buffered
results := make(chan Result, 100) // ✅ Buffered
```

---

## 🐛 Common Mistakes

### Mistake 1: Not Closing Channels

```go
// ❌ BAD: Channel never closes
go func() {
    for val := range input {
        output <- process(val)
    }
    // Forgot close(output)!
}()
```

### Mistake 2: Goroutine Leaks

```go
// ❌ BAD: No exit strategy
go func() {
    for {
        doWork()  // Forever!
    }
}()

// ✅ GOOD: Use context
go func() {
    for {
        select {
        case <-ctx.Done():
            return  // ✅ Can exit
        default:
            doWork()
        }
    }
}()
```

### Mistake 3: Deadlock on Unbuffered Channels

```go
// ❌ BAD: Blocks forever
ch := make(chan int)
ch <- 42  // Blocks! No receiver

// ✅ GOOD: Use goroutine
go func() {
    ch <- 42
}()
val := <-ch
```

---

## 📊 Performance Tips

### Tune Worker Count

```go
import "runtime"

// CPU-bound
numWorkers := runtime.NumCPU()

// I/O-bound
numWorkers := runtime.NumCPU() * 2

// Mixed
numWorkers := runtime.NumCPU() + 2
```

### Buffer Size

```go
// Small (tight control)
ch := make(chan T, numWorkers)

// Medium (balance)
ch := make(chan T, numWorkers * 2)

// Large (throughput)
ch := make(chan T, 1000)
```

### Rate Limiting

```go
// Limit concurrent operations
semaphore := make(chan struct{}, maxConcurrent)

for job := range jobs {
    semaphore <- struct{}{}  // Acquire
    
    go func() {
        defer func() { <-semaphore }()  // Release
        process(job)
    }()
}
```

---

## 🎯 Combining Patterns

### Worker Pool + Pipeline

```go
// Pipeline stage with worker pool
func parallelStage(input <-chan T, numWorkers int) <-chan R {
    out := make(chan R)
    var wg sync.WaitGroup
    
    // Worker pool for this stage
    for w := 0; w < numWorkers; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for val := range input {
                result := process(val)
                out <- result
            }
        }()
    }
    
    go func() {
        wg.Wait()
        close(out)
    }()
    
    return out
}
```

### Fan-Out + Pipeline

```go
// Multi-stage with parallel processing
stage1 := generator(data)
stage2Workers := fanOut(stage1, 5)  // Parallel processing
stage2 := fanIn(stage2Workers...)
stage3 := filter(stage2)
output := collect(stage3)
```

---

## 🔍 Debugging

### Check Goroutine Count

```go
import "runtime"

before := runtime.NumGoroutine()
// Run code
after := runtime.NumGoroutine()

if after > before {
    fmt.Printf("Leaked %d goroutines\n", after-before)
}
```

### Use Race Detector

```bash
go run -race main.go
go test -race ./...
```

### pprof для Profiling

```go
import _ "net/http/pprof"

go http.ListenAndServe("localhost:6060", nil)
```

```bash
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

---

## 🎓 Key Takeaways

### Worker Pool:

1. **Fixed concurrency** (control resource usage)
2. **Reuse goroutines** (efficient)
3. **Job queue pattern** (decoupling)

### Fan-Out/Fan-In:

1. **Parallel processing** (maximize throughput)
2. **Merge results** (single output stream)
3. **Dynamic workers** (flexible)

### Pipeline:

1. **Composable stages** (modularity)
2. **Streaming processing** (memory efficient)
3. **Each stage = one responsibility** (clean code)

### Context:

1. **Always use context** (for cancellation)
2. **Propagate through pipeline** (all stages)
3. **Graceful shutdown** (clean cleanup)

---

## 📖 Ресурси

### Документація

- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Advanced Go Concurrency](https://go.dev/blog/io2013-talk-concurrency)
- Week 6: `theory/07_goroutines_concurrency.md`
- Week 8: `theory/01_race_conditions.md`

### Recommended Reading

- "Concurrency in Go" by Katherine Cox-Buday
- "Go in Action" by William Kennedy

---

## 🚀 Quick Start

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_9

# Read theory
cat theory/01_worker_pool.md
cat theory/02_fan_in_fan_out.md
cat theory/03_pipeline.md

# Run examples
cd practice/01_worker_pool && go run simple_worker_pool.go
cd ../02_fan_in_fan_out && go run fan_pattern.go
cd ../03_pipeline && go run three_stage_pipeline.go
```

---

## 🎯 Learning Path

### День 1: Worker Pool

1. Читай `theory/01_worker_pool.md`
2. Запускай `practice/01_worker_pool/simple_worker_pool.go`
3. Зрозумій controlled concurrency
4. Напиши власний worker pool

### День 2: Fan-In/Fan-Out

1. Читай `theory/02_fan_in_fan_out.md`
2. Запускай `practice/02_fan_in_fan_out/fan_pattern.go`
3. Зрозумій parallel processing
4. Напиши власний fan pattern

### День 3: Pipeline

1. Читай `theory/03_pipeline.md`
2. Запускай `practice/03_pipeline/three_stage_pipeline.go`
3. Зрозумій composable stages
4. Напиши власний pipeline

### День 4: Combine Patterns

1. Combine Worker Pool + Pipeline
2. Add context cancellation
3. Add error handling
4. Build production-ready system

---

**"Concurrency is not parallelism, but it enables parallelism!" 🔄**

**Status:** Week 9 Materials Complete ✅  
**Created:** 2026-01-28
