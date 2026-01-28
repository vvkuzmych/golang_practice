# Fan-In / Fan-Out Pattern

## 🎯 Що таке Fan-Out і Fan-In?

**Fan-Out** - розподіл роботи між множину workers  
**Fan-In** - збір результатів від множини workers в один channel

```
        Fan-Out (розподіл)
            
Input → [Worker 1] →
      → [Worker 2] →  Fan-In (збір) → Output
      → [Worker 3] →
```

---

## 📊 Fan-Out Pattern

### Що таке Fan-Out?

Один producer → Many consumers

```go
func fanOut(input <-chan int, numWorkers int) []<-chan int {
    outputs := make([]<-chan int, numWorkers)
    
    for i := 0; i < numWorkers; i++ {
        ch := make(chan int)
        outputs[i] = ch
        
        go func(out chan<-int) {
            for val := range input {
                result := process(val)  // Heavy work
                out <- result
            }
            close(out)
        }(ch)
    }
    
    return outputs
}
```

**Use case:** Паралельна обробка даних

---

## 📊 Fan-In Pattern

### Що таке Fan-In?

Many producers → One consumer

```go
func fanIn(channels ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    
    // Start goroutine for each input channel
    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for val := range c {
                out <- val
            }
        }(ch)
    }
    
    // Close output when all inputs done
    go func() {
        wg.Wait()
        close(out)
    }()
    
    return out
}
```

**Use case:** Збір результатів з кількох джерел

---

## 🎯 Повний приклад: Fan-Out + Fan-In

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Fan-Out: Розподіл роботи
func fanOut(input <-chan int, numWorkers int) []<-chan int {
    outputs := make([]<-chan int, numWorkers)
    
    for i := 0; i < numWorkers; i++ {
        out := make(chan int)
        outputs[i] = out
        
        go func(workerID int, out chan<- int) {
            defer close(out)
            
            for val := range input {
                fmt.Printf("Worker %d processing %d\n", workerID, val)
                
                // Simulate heavy work
                time.Sleep(100 * time.Millisecond)
                result := val * 2
                
                out <- result
            }
            
            fmt.Printf("Worker %d done\n", workerID)
        }(i, out)
    }
    
    return outputs
}

// Fan-In: Збір результатів
func fanIn(channels ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    
    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
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

// Generator: Create input stream
func generator(nums []int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, num := range nums {
            out <- num
        }
    }()
    return out
}

func main() {
    // Input data
    input := generator([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
    
    // Fan-Out to 3 workers
    workers := fanOut(input, 3)
    
    // Fan-In results
    results := fanIn(workers...)
    
    // Collect results
    for result := range results {
        fmt.Printf("Result: %d\n", result)
    }
    
    fmt.Println("All done!")
}
```

---

## 🔄 Fan-In з Select (Multiplexing)

```go
func fanInSelect(ch1, ch2 <-chan int) <-chan int {
    out := make(chan int)
    
    go func() {
        defer close(out)
        
        for {
            select {
            case val, ok := <-ch1:
                if !ok {
                    ch1 = nil  // Disable this case
                    continue
                }
                out <- val
            case val, ok := <-ch2:
                if !ok {
                    ch2 = nil  // Disable this case
                    continue
                }
                out <- val
            }
            
            // Both channels closed
            if ch1 == nil && ch2 == nil {
                return
            }
        }
    }()
    
    return out
}
```

**Use case:** Merge два streams в один

---

## 🎯 Real-World Example: HTTP Fetcher

```go
package main

import (
    "fmt"
    "io"
    "net/http"
    "sync"
    "time"
)

type URLRequest struct {
    URL string
}

type URLResponse struct {
    URL    string
    Body   string
    Error  error
}

// Fan-Out: Multiple HTTP fetchers
func httpFanOut(urls <-chan URLRequest, numWorkers int) []<-chan URLResponse {
    outputs := make([]<-chan URLResponse, numWorkers)
    
    client := &http.Client{Timeout: 5 * time.Second}
    
    for i := 0; i < numWorkers; i++ {
        out := make(chan URLResponse)
        outputs[i] = out
        
        go func(workerID int, out chan<- URLResponse) {
            defer close(out)
            
            for req := range urls {
                fmt.Printf("Worker %d fetching %s\n", workerID, req.URL)
                
                resp, err := client.Get(req.URL)
                if err != nil {
                    out <- URLResponse{URL: req.URL, Error: err}
                    continue
                }
                
                body, _ := io.ReadAll(resp.Body)
                resp.Body.Close()
                
                out <- URLResponse{
                    URL:  req.URL,
                    Body: string(body),
                }
            }
        }(i, out)
    }
    
    return outputs
}

// Fan-In: Collect responses
func httpFanIn(channels ...<-chan URLResponse) <-chan URLResponse {
    out := make(chan URLResponse)
    var wg sync.WaitGroup
    
    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan URLResponse) {
            defer wg.Done()
            for resp := range c {
                out <- resp
            }
        }(ch)
    }
    
    go func() {
        wg.Wait()
        close(out)
    }()
    
    return out
}

func main() {
    urls := []string{
        "https://google.com",
        "https://github.com",
        "https://golang.org",
        "https://stackoverflow.com",
    }
    
    // Create input channel
    input := make(chan URLRequest)
    go func() {
        defer close(input)
        for _, url := range urls {
            input <- URLRequest{URL: url}
        }
    }()
    
    // Fan-Out to 2 workers
    workers := httpFanOut(input, 2)
    
    // Fan-In results
    results := httpFanIn(workers...)
    
    // Process results
    for resp := range results {
        if resp.Error != nil {
            fmt.Printf("Error fetching %s: %v\n", resp.URL, resp.Error)
        } else {
            fmt.Printf("Fetched %s: %d bytes\n", resp.URL, len(resp.Body))
        }
    }
}
```

---

## 🎯 Fan-Out з Context (Cancellation)

```go
func fanOutWithContext(ctx context.Context, input <-chan int, numWorkers int) []<-chan int {
    outputs := make([]<-chan int, numWorkers)
    
    for i := 0; i < numWorkers; i++ {
        out := make(chan int)
        outputs[i] = out
        
        go func(workerID int, out chan<- int) {
            defer close(out)
            
            for {
                select {
                case <-ctx.Done():
                    fmt.Printf("Worker %d cancelled\n", workerID)
                    return
                case val, ok := <-input:
                    if !ok {
                        return
                    }
                    
                    // Process with cancellation check
                    select {
                    case <-ctx.Done():
                        return
                    default:
                        result := process(val)
                        out <- result
                    }
                }
            }
        }(i, out)
    }
    
    return outputs
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    input := generator([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
    
    workers := fanOutWithContext(ctx, input, 3)
    results := fanIn(workers...)
    
    for result := range results {
        fmt.Printf("Result: %d\n", result)
    }
}
```

---

## 📊 Bounded Fan-Out (Limited Concurrency)

```go
type BoundedFanOut struct {
    numWorkers int
    semaphore  chan struct{}
}

func NewBoundedFanOut(maxConcurrent int) *BoundedFanOut {
    return &BoundedFanOut{
        numWorkers: maxConcurrent,
        semaphore:  make(chan struct{}, maxConcurrent),
    }
}

func (b *BoundedFanOut) Process(input <-chan int) <-chan int {
    out := make(chan int)
    
    go func() {
        defer close(out)
        var wg sync.WaitGroup
        
        for val := range input {
            // Acquire semaphore
            b.semaphore <- struct{}{}
            
            wg.Add(1)
            go func(v int) {
                defer wg.Done()
                defer func() { <-b.semaphore }()  // Release
                
                result := process(v)
                out <- result
            }(val)
        }
        
        wg.Wait()
    }()
    
    return out
}
```

---

## 🎯 Fan-In з Priority

```go
type PriorityMessage struct {
    Priority int
    Data     interface{}
}

func priorityFanIn(high, low <-chan PriorityMessage) <-chan PriorityMessage {
    out := make(chan PriorityMessage)
    
    go func() {
        defer close(out)
        
        for {
            select {
            case msg, ok := <-high:
                if !ok {
                    high = nil
                    continue
                }
                out <- msg
            default:
                // Only read from low if high is empty
                select {
                case msg, ok := <-high:
                    if !ok {
                        high = nil
                        continue
                    }
                    out <- msg
                case msg, ok := <-low:
                    if !ok {
                        low = nil
                        continue
                    }
                    out <- msg
                }
            }
            
            if high == nil && low == nil {
                return
            }
        }
    }()
    
    return out
}
```

---

## 📈 Performance Patterns

### Pattern 1: Dynamic Workers

```go
func dynamicFanOut(input <-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    
    go func() {
        defer close(out)
        
        for val := range input {
            wg.Add(1)
            go func(v int) {
                defer wg.Done()
                result := process(v)
                out <- result
            }(val)
        }
        
        wg.Wait()
    }()
    
    return out
}
```

**Use case:** Коли не знаєш скільки workers потрібно

### Pattern 2: Fixed Workers (Worker Pool)

```go
func fixedFanOut(input <-chan int, numWorkers int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    
    for i := 0; i < numWorkers; i++ {
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

**Use case:** Контрольована concurrency

---

## ✅ Best Practices

### 1. Close Channels Properly

```go
// ✅ Producer closes input
go func() {
    for _, val := range data {
        input <- val
    }
    close(input)
}()

// ✅ Workers close outputs when done
go func() {
    defer close(out)
    for val := range input {
        out <- process(val)
    }
}()
```

### 2. Use WaitGroup для Fan-In

```go
func fanIn(channels ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    
    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for val := range c {
                out <- val
            }
        }(ch)
    }
    
    // ✅ Close output after all inputs done
    go func() {
        wg.Wait()
        close(out)
    }()
    
    return out
}
```

### 3. Handle Context для Cancellation

```go
select {
case <-ctx.Done():
    return  // ✅ Respect cancellation
case val := <-input:
    process(val)
}
```

---

## 🐛 Common Mistakes

### Mistake 1: Not Closing Output Channel

```go
// ❌ BAD: Output never closes
func fanIn(ch1, ch2 <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for val := range ch1 {
            out <- val
        }
    }()
    go func() {
        for val := range ch2 {
            out <- val
        }
    }()
    // Forgot to close(out)!
    return out
}
```

### Mistake 2: Race on Channel Nil

```go
// ❌ BAD: Can panic
select {
case val := <-ch1:
    ch1 = nil  // ❌ While still in select!
}

// ✅ GOOD: Check before disabling
case val, ok := <-ch1:
    if !ok {
        ch1 = nil
        continue
    }
```

---

## 🎓 Висновок

### Fan-Out:

✅ Distribute work to multiple workers  
✅ Parallel processing  
✅ Use for CPU-intensive or I/O tasks  

### Fan-In:

✅ Merge multiple streams  
✅ Collect results from workers  
✅ Use WaitGroup to sync  

### Combined Pattern:

```
Input → Fan-Out → [Workers] → Fan-In → Output
```

### Key Points:

1. **Fan-Out** для паралелізації
2. **Fan-In** для збору результатів
3. **Always close channels**
4. **Use context** для cancellation
5. **WaitGroup** для синхронізації

---

## 📖 Далі

- `practice/02_fan_in_fan_out/` - Практичні приклади
- `03_pipeline.md` - Pipeline pattern
- Combine patterns для powerful processing

**"Fan-Out/Fan-In = Parallel Processing!" 🔄**
