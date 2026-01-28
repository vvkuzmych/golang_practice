# Worker Pool Pattern

## 🎯 Що таке Worker Pool?

**Worker Pool** - це pattern, де фіксована кількість goroutines (workers) обробляють tasks з спільної черги.

```
Jobs Queue → [Worker 1]
           → [Worker 2]
           → [Worker 3]
           → [Worker 4]
                ↓
            Results
```

---

## ⚠️ Проблема

Якщо створювати goroutine для кожного task:

```go
// ❌ BAD: Unbounded goroutines
for i := 0; i < 1000000; i++ {
    go processTask(i)  // 1 million goroutines!
}
```

**Проблеми:**
- Занадто багато goroutines
- Overwhelm scheduler
- Memory overhead (~2KB per goroutine)
- Resource exhaustion

---

## ✅ Рішення: Worker Pool

Обмежена кількість workers:

```go
// ✅ GOOD: Fixed workers
numWorkers := 10
jobs := make(chan Job, 100)
results := make(chan Result, 100)

// Start workers
for w := 1; w <= numWorkers; w++ {
    go worker(w, jobs, results)
}

// Send jobs
for i := 0; i < 1000000; i++ {
    jobs <- Job{ID: i}
}
close(jobs)
```

---

## 📊 Базова реалізація

### Версія 1: Simple Worker Pool

```go
package main

import (
    "fmt"
    "time"
)

type Job struct {
    ID int
}

type Result struct {
    Job    Job
    Result int
}

func worker(id int, jobs <-chan Job, results chan<- Result) {
    for job := range jobs {
        fmt.Printf("Worker %d started job %d\n", id, job.ID)
        
        // Simulate work
        time.Sleep(100 * time.Millisecond)
        result := job.ID * 2
        
        results <- Result{Job: job, Result: result}
    }
}

func main() {
    numJobs := 10
    numWorkers := 3
    
    jobs := make(chan Job, numJobs)
    results := make(chan Result, numJobs)
    
    // Start workers
    for w := 1; w <= numWorkers; w++ {
        go worker(w, jobs, results)
    }
    
    // Send jobs
    for j := 1; j <= numJobs; j++ {
        jobs <- Job{ID: j}
    }
    close(jobs)
    
    // Collect results
    for a := 1; a <= numJobs; a++ {
        result := <-results
        fmt.Printf("Job %d result: %d\n", result.Job.ID, result.Result)
    }
}
```

---

## 🏗️ Advanced: Worker Pool з WaitGroup

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type Job struct {
    ID   int
    Data string
}

func worker(id int, jobs <-chan Job, wg *sync.WaitGroup) {
    defer wg.Done()
    
    for job := range jobs {
        fmt.Printf("Worker %d processing job %d\n", id, job.ID)
        time.Sleep(100 * time.Millisecond)
        fmt.Printf("Worker %d finished job %d\n", id, job.ID)
    }
}

func main() {
    numWorkers := 3
    jobs := make(chan Job, 10)
    var wg sync.WaitGroup
    
    // Start workers
    for w := 1; w <= numWorkers; w++ {
        wg.Add(1)
        go worker(w, jobs, &wg)
    }
    
    // Send jobs
    for j := 1; j <= 10; j++ {
        jobs <- Job{ID: j, Data: fmt.Sprintf("task-%d", j)}
    }
    close(jobs)
    
    // Wait for all workers to finish
    wg.Wait()
    fmt.Println("All jobs completed!")
}
```

---

## 🎯 Worker Pool з Context (Cancellation)

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"
)

type Job struct {
    ID int
}

func worker(ctx context.Context, id int, jobs <-chan Job, wg *sync.WaitGroup) {
    defer wg.Done()
    
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %d cancelled\n", id)
            return
        case job, ok := <-jobs:
            if !ok {
                fmt.Printf("Worker %d: no more jobs\n", id)
                return
            }
            
            fmt.Printf("Worker %d processing job %d\n", id, job.ID)
            
            // Simulate work with cancellation check
            select {
            case <-ctx.Done():
                fmt.Printf("Worker %d: job %d cancelled\n", id, job.ID)
                return
            case <-time.After(100 * time.Millisecond):
                fmt.Printf("Worker %d finished job %d\n", id, job.ID)
            }
        }
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    numWorkers := 3
    jobs := make(chan Job, 10)
    var wg sync.WaitGroup
    
    // Start workers
    for w := 1; w <= numWorkers; w++ {
        wg.Add(1)
        go worker(ctx, w, jobs, &wg)
    }
    
    // Send jobs
    go func() {
        for j := 1; j <= 20; j++ {
            jobs <- Job{ID: j}
        }
        close(jobs)
    }()
    
    // Wait for completion or timeout
    wg.Wait()
    fmt.Println("Done!")
}
```

---

## 📈 Worker Pool з Rate Limiting

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type RateLimitedPool struct {
    workers   int
    jobs      chan Job
    results   chan Result
    semaphore chan struct{}
    wg        sync.WaitGroup
}

type Job struct {
    ID int
}

type Result struct {
    JobID  int
    Result int
}

func NewRateLimitedPool(workers int, rateLimit int) *RateLimitedPool {
    return &RateLimitedPool{
        workers:   workers,
        jobs:      make(chan Job, 100),
        results:   make(chan Result, 100),
        semaphore: make(chan struct{}, rateLimit),
    }
}

func (p *RateLimitedPool) Start() {
    for w := 1; w <= p.workers; w++ {
        p.wg.Add(1)
        go p.worker(w)
    }
}

func (p *RateLimitedPool) worker(id int) {
    defer p.wg.Done()
    
    for job := range p.jobs {
        // Rate limiting
        p.semaphore <- struct{}{}
        
        fmt.Printf("Worker %d processing job %d\n", id, job.ID)
        time.Sleep(100 * time.Millisecond)
        
        p.results <- Result{JobID: job.ID, Result: job.ID * 2}
        
        <-p.semaphore
    }
}

func (p *RateLimitedPool) Submit(job Job) {
    p.jobs <- job
}

func (p *RateLimitedPool) Close() {
    close(p.jobs)
}

func (p *RateLimitedPool) Wait() {
    p.wg.Wait()
    close(p.results)
}

func main() {
    pool := NewRateLimitedPool(5, 2) // 5 workers, max 2 concurrent
    pool.Start()
    
    // Submit jobs
    go func() {
        for i := 1; i <= 10; i++ {
            pool.Submit(Job{ID: i})
        }
        pool.Close()
    }()
    
    // Collect results
    go func() {
        for result := range pool.results {
            fmt.Printf("Result for job %d: %d\n", result.JobID, result.Result)
        }
    }()
    
    pool.Wait()
    time.Sleep(100 * time.Millisecond)
}
```

---

## 🎯 Use Cases

### 1. HTTP Request Processing

```go
type HTTPJob struct {
    URL string
}

func fetchWorker(id int, jobs <-chan HTTPJob, results chan<- string) {
    client := &http.Client{Timeout: 5 * time.Second}
    
    for job := range jobs {
        resp, err := client.Get(job.URL)
        if err != nil {
            results <- fmt.Sprintf("Error: %v", err)
            continue
        }
        
        body, _ := io.ReadAll(resp.Body)
        resp.Body.Close()
        
        results <- string(body)
    }
}

// Usage: Pool of 10 workers fetching URLs
```

### 2. Database Batch Processing

```go
type DBJob struct {
    Query  string
    Params []interface{}
}

func dbWorker(id int, jobs <-chan DBJob, db *sql.DB) {
    for job := range jobs {
        _, err := db.Exec(job.Query, job.Params...)
        if err != nil {
            log.Printf("Worker %d error: %v", id, err)
        }
    }
}

// Usage: Pool of 5 workers for bulk inserts
```

### 3. Image Processing

```go
type ImageJob struct {
    InputPath  string
    OutputPath string
}

func imageWorker(id int, jobs <-chan ImageJob) {
    for job := range jobs {
        img := loadImage(job.InputPath)
        resized := resize(img, 800, 600)
        saveImage(resized, job.OutputPath)
    }
}

// Usage: Pool of CPU-core workers for parallel processing
```

---

## 📊 Performance Tuning

### Оптимальна кількість workers

```go
import "runtime"

// CPU-bound tasks
numWorkers := runtime.NumCPU()

// I/O-bound tasks (HTTP, DB)
numWorkers := runtime.NumCPU() * 2  // Or more

// Mixed workload
numWorkers := runtime.NumCPU() + 2
```

### Buffer Size

```go
// Small buffer (tight control)
jobs := make(chan Job, numWorkers)

// Medium buffer (balance)
jobs := make(chan Job, numWorkers * 2)

// Large buffer (throughput)
jobs := make(chan Job, 1000)
```

---

## ✅ Best Practices

### 1. Always Close Jobs Channel

```go
// Producer closes
go func() {
    for _, job := range jobs {
        jobsChan <- job
    }
    close(jobsChan)  // ✅ Signal workers to stop
}()
```

### 2. Use WaitGroup для Cleanup

```go
var wg sync.WaitGroup
for w := 0; w < numWorkers; w++ {
    wg.Add(1)
    go worker(&wg, jobs)
}
wg.Wait()  // ✅ Wait for all workers
```

### 3. Handle Panics в Workers

```go
func worker(jobs <-chan Job) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Worker panic: %v", r)
        }
    }()
    
    for job := range jobs {
        processJob(job)
    }
}
```

### 4. Graceful Shutdown з Context

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Signal shutdown
cancel()

// Wait for workers with timeout
done := make(chan struct{})
go func() {
    wg.Wait()
    close(done)
}()

select {
case <-done:
    // Clean shutdown
case <-time.After(5 * time.Second):
    // Force shutdown
}
```

---

## 🐛 Common Mistakes

### Mistake 1: Not Closing Jobs Channel

```go
// ❌ BAD: Workers wait forever
go func() {
    for i := 0; i < 10; i++ {
        jobs <- Job{ID: i}
    }
    // Forgot close(jobs)!
}()
```

### Mistake 2: Wrong WaitGroup Usage

```go
// ❌ BAD: wg.Add() inside goroutine
go func() {
    wg.Add(1)  // ❌ Race condition!
    defer wg.Done()
}()

// ✅ GOOD: wg.Add() before goroutine
wg.Add(1)
go func() {
    defer wg.Done()
}()
```

### Mistake 3: Blocking on Results

```go
// ❌ BAD: Deadlock if results buffer full
for i := 0; i < 1000; i++ {
    jobs <- Job{ID: i}
}
close(jobs)

for i := 0; i < 1000; i++ {
    <-results  // May block if buffer small
}

// ✅ GOOD: Read results concurrently
go func() {
    for result := range results {
        processResult(result)
    }
}()
```

---

## 🎓 Висновок

### Worker Pool - це:

✅ Fixed number of goroutines  
✅ Process jobs from queue  
✅ Control concurrency  
✅ Reuse goroutines  

### Коли використовувати:

- HTTP request processing
- Database batch operations
- Image/video processing
- File processing
- Any I/O-bound or CPU-bound tasks

### Key Points:

1. **Limit concurrency** (don't create unlimited goroutines)
2. **Close jobs channel** (signal workers to stop)
3. **Use WaitGroup** (wait for completion)
4. **Add context** (enable cancellation)
5. **Tune workers count** (CPU-bound vs I/O-bound)

---

## 📖 Далі

- `practice/01_worker_pool/` - Практичні приклади
- `02_fan_in_fan_out.md` - Fan-In/Fan-Out pattern
- `03_pipeline.md` - Pipeline pattern

**"Worker Pool = Controlled Concurrency!" 🔄**
