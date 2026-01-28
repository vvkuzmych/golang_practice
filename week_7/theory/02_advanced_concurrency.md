# Advanced Concurrency Patterns

## Worker Pool Pattern

```go
type Job struct {
    ID   int
    Data string
}

type Result struct {
    Job   Job
    Value int
    Err   error
}

func worker(id int, jobs <-chan Job, results chan<- Result) {
    for job := range jobs {
        // Process
        results <- Result{Job: job, Value: len(job.Data)}
    }
}

func main() {
    jobs := make(chan Job, 100)
    results := make(chan Result, 100)
    
    // Start workers
    for w := 1; w <= 5; w++ {
        go worker(w, jobs, results)
    }
    
    // Send jobs
    for j := 1; j <= 10; j++ {
        jobs <- Job{ID: j, Data: "data"}
    }
    close(jobs)
    
    // Collect results
    for i := 1; i <= 10; i++ {
        result := <-results
        fmt.Printf("Result: %+v\n", result)
    }
}
```

## Rate Limiting

```go
import "golang.org/x/time/rate"

func main() {
    limiter := rate.NewLimiter(10, 5) // 10 req/sec, burst 5
    
    for i := 0; i < 100; i++ {
        if err := limiter.Wait(context.Background()); err != nil {
            log.Fatal(err)
        }
        handleRequest()
    }
}
```

## Circuit Breaker

```go
type CircuitBreaker struct {
    maxFailures int
    timeout     time.Duration
    failures    int
    lastFail    time.Time
    state       string // "closed", "open", "half-open"
    mu          sync.Mutex
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    cb.mu.Lock()
    defer cb.mu.Unlock()
    
    if cb.state == "open" {
        if time.Since(cb.lastFail) > cb.timeout {
            cb.state = "half-open"
        } else {
            return errors.New("circuit breaker open")
        }
    }
    
    err := fn()
    if err != nil {
        cb.failures++
        cb.lastFail = time.Now()
        if cb.failures >= cb.maxFailures {
            cb.state = "open"
        }
        return err
    }
    
    cb.failures = 0
    cb.state = "closed"
    return nil
}
```

## Semaphore Pattern

```go
type Semaphore chan struct{}

func NewSemaphore(max int) Semaphore {
    return make(Semaphore, max)
}

func (s Semaphore) Acquire() {
    s <- struct{}{}
}

func (s Semaphore) Release() {
    <-s
}

func main() {
    sem := NewSemaphore(3) // Max 3 concurrent
    
    for i := 0; i < 10; i++ {
        go func(id int) {
            sem.Acquire()
            defer sem.Release()
            
            // Do work
            time.Sleep(time.Second)
            fmt.Printf("Task %d done\n", id)
        }(i)
    }
}
```
