# Week 9 - Quick Start 🚀

## 🎯 Мета

Опанувати **3 concurrency patterns**: Worker Pool, Fan-In/Fan-Out, Pipeline.

---

## ⚡ 5-хвилинний старт

### 1. Worker Pool

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_9/practice/01_worker_pool
go run simple_worker_pool.go
```

**Що побачиш:**
- 3 workers обробляють 10 jobs
- Controlled concurrency
- Timeout cancellation
- Manual cancellation

### 2. Fan-In / Fan-Out

```bash
cd ../02_fan_in_fan_out
go run fan_pattern.go
```

**Що побачиш:**
- Fan-Out до кількох workers
- Parallel processing
- Fan-In results
- Збір в один stream

### 3. Pipeline з Context

```bash
cd ../03_pipeline
go run three_stage_pipeline.go
```

**Що побачиш:**
- 3 stages: Generate → Square → Filter
- Normal completion
- Timeout cancellation
- Manual cancellation

---

## 📚 Читати теорію

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_9

# Worker Pool
cat theory/01_worker_pool.md

# Fan-In/Fan-Out
cat theory/02_fan_in_fan_out.md

# Pipeline
cat theory/03_pipeline.md
```

---

## 🎯 3 Головні Patterns

### Pattern 1: Worker Pool

```go
numWorkers := 5
jobs := make(chan Job, 100)
results := make(chan Result, 100)

// Start workers
for w := 1; w <= numWorkers; w++ {
    go worker(w, jobs, results)
}

// Send jobs
for j := 1; j <= 1000; j++ {
    jobs <- Job{ID: j}
}
close(jobs)
```

**Use case:** HTTP processing, batch operations

### Pattern 2: Fan-Out / Fan-In

```go
// Fan-Out: розподіл
workers := fanOut(input, numWorkers)

// Fan-In: збір
results := fanIn(workers...)
```

**Use case:** Parallel processing, aggregation

### Pattern 3: Pipeline

```go
// 3 stages connected
stage1 := generate(data)
stage2 := transform(stage1)
stage3 := filter(stage2)

// Process
for result := range stage3 {
    fmt.Println(result)
}
```

**Use case:** ETL, log processing, streaming

---

## 🔧 Context для Cancellation

### Timeout

```go
ctx, cancel := context.WithTimeout(
    context.Background(), 
    5*time.Second
)
defer cancel()

stage := process(ctx, input)
```

### Manual Cancel

```go
ctx, cancel := context.WithCancel(context.Background())

// Cancel on condition
if shouldStop {
    cancel()
}

stage := process(ctx, input)
```

### Check Cancellation

```go
select {
case <-ctx.Done():
    return  // ✅ Stop
case val := <-input:
    process(val)
}
```

---

## ✅ 3 Golden Rules

### Rule 1: Close Channels

```go
go func() {
    defer close(out)  // ✅ Always
    for val := range input {
        out <- val
    }
}()
```

### Rule 2: Use Context

```go
func stage(ctx context.Context, in <-chan T) <-chan R {
    out := make(chan R)
    go func() {
        defer close(out)
        for {
            select {
            case <-ctx.Done():
                return  // ✅ Respect cancellation
            case val, ok := <-in:
                if !ok { return }
                out <- process(val)
            }
        }
    }()
    return out
}
```

### Rule 3: WaitGroup Sync

```go
var wg sync.WaitGroup

for i := 0; i < numWorkers; i++ {
    wg.Add(1)  // ✅ Before goroutine
    go func() {
        defer wg.Done()
        // Work
    }()
}

wg.Wait()
```

---

## 🐛 Common Patterns

### Pattern 1: Fixed Worker Pool

```go
func workerPool(numWorkers int, jobs <-chan Job) <-chan Result {
    results := make(chan Result)
    var wg sync.WaitGroup
    
    for w := 0; w < numWorkers; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                results <- process(job)
            }
        }()
    }
    
    go func() {
        wg.Wait()
        close(results)
    }()
    
    return results
}
```

### Pattern 2: Fan-Out/Fan-In

```go
func fanOut(in <-chan T, n int) []<-chan R {
    outs := make([]<-chan R, n)
    for i := 0; i < n; i++ {
        ch := make(chan R)
        outs[i] = ch
        go func(out chan<- R) {
            defer close(out)
            for val := range in {
                out <- process(val)
            }
        }(ch)
    }
    return outs
}

func fanIn(channels ...<-chan T) <-chan T {
    out := make(chan T)
    var wg sync.WaitGroup
    
    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan T) {
            defer wg.Done()
            for val := range c {
                out <- val
            }
        }(ch)
    }
    
    go func() {
        wg.Wait()
        close(out)
    }()
    
    return out
}
```

### Pattern 3: Pipeline Stage

```go
func stage(in <-chan T) <-chan R {
    out := make(chan R)
    go func() {
        defer close(out)
        for val := range in {
            out <- process(val)
        }
    }()
    return out
}

// Build pipeline
s1 := stage1(input)
s2 := stage2(s1)
s3 := stage3(s2)
```

---

## 📖 Структура

```
week_9/
├── README.md           # Повний опис
├── QUICK_START.md      # Цей файл
├── theory/
│   ├── 01_worker_pool.md
│   ├── 02_fan_in_fan_out.md
│   └── 03_pipeline.md
├── practice/
│   ├── 01_worker_pool/
│   │   └── simple_worker_pool.go
│   ├── 02_fan_in_fan_out/
│   │   └── fan_pattern.go
│   └── 03_pipeline/
│       └── three_stage_pipeline.go
└── exercises/          # TODO
```

---

## 🎯 Рекомендований порядок

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
4. Напиши власний pipeline з 4+ stages

---

## 🎓 Key Concepts

### Worker Pool
- **Fixed workers** (control concurrency)
- **Job queue** (decouple producer/consumer)
- **Reuse goroutines** (efficiency)

### Fan-Out/Fan-In
- **Fan-Out:** split work to many
- **Fan-In:** merge results to one
- **Parallel processing** (speed)

### Pipeline
- **Series of stages** (composable)
- **Connected by channels** (streaming)
- **Each stage = one function** (SRP)

---

## 🔍 Debug Tips

```go
// Check goroutine count
before := runtime.NumGoroutine()
// ... code ...
after := runtime.NumGoroutine()
fmt.Printf("Goroutines: %d -> %d\n", before, after)

// Race detector
go run -race main.go

// pprof
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

---

## 📖 Ресурси

- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- Week 6: Goroutines basics
- Week 8: Debugging & Race

---

**"Master these 3 patterns, master Go concurrency!" 🔄**
