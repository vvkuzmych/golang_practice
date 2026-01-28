# Pipeline Pattern

## 🎯 Що таке Pipeline?

**Pipeline** - це series of stages connected by channels, де кожна stage є група goroutines, що виконує одну функцію.

```
Input → [Stage 1] → [Stage 2] → [Stage 3] → Output
         (goroutines)  (goroutines)  (goroutines)
```

**Кожна stage:**
- Receives values від upstream через inbound channels
- Performs some function на даних
- Sends values downstream через outbound channels

---

## 📊 Базовий Pipeline (3 Stages)

```go
package main

import "fmt"

// Stage 1: Generate numbers
func generator(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            out <- n
        }
    }()
    return out
}

// Stage 2: Square numbers
func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            out <- n * n
        }
    }()
    return out
}

// Stage 3: Sum numbers
func sum(in <-chan int) int {
    total := 0
    for n := range in {
        total += n
    }
    return total
}

func main() {
    // Build pipeline
    nums := generator(1, 2, 3, 4, 5)
    squares := square(nums)
    result := sum(squares)
    
    fmt.Printf("Sum of squares: %d\n", result)  // 55
}
```

---

## 🔄 Pipeline з Multiple Stages

```go
package main

import (
    "fmt"
    "strings"
)

type Data struct {
    Value string
}

// Stage 1: Generate data
func generate(values []string) <-chan Data {
    out := make(chan Data)
    go func() {
        defer close(out)
        for _, v := range values {
            out <- Data{Value: v}
        }
    }()
    return out
}

// Stage 2: Transform to uppercase
func uppercase(in <-chan Data) <-chan Data {
    out := make(chan Data)
    go func() {
        defer close(out)
        for data := range in {
            data.Value = strings.ToUpper(data.Value)
            out <- data
        }
    }()
    return out
}

// Stage 3: Add prefix
func prefix(in <-chan Data, pre string) <-chan Data {
    out := make(chan Data)
    go func() {
        defer close(out)
        for data := range in {
            data.Value = pre + data.Value
            out <- data
        }
    }()
    return out
}

// Stage 4: Filter (only длина > 5)
func filter(in <-chan Data) <-chan Data {
    out := make(chan Data)
    go func() {
        defer close(out)
        for data := range in {
            if len(data.Value) > 5 {
                out <- data
            }
        }
    }()
    return out
}

func main() {
    input := []string{"hello", "world", "go", "pipeline"}
    
    // Build pipeline
    stage1 := generate(input)
    stage2 := uppercase(stage1)
    stage3 := prefix(stage2, ">> ")
    stage4 := filter(stage3)
    
    // Process results
    for data := range stage4 {
        fmt.Println(data.Value)
    }
}
```

**Output:**
```
>> HELLO
>> WORLD
>> PIPELINE
```

---

## 🎯 Pipeline з Fan-Out/Fan-In

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Stage 1: Generator
func generator(nums []int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            out <- n
        }
    }()
    return out
}

// Stage 2: Heavy processing (Fan-Out)
func process(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            // Simulate heavy work
            time.Sleep(100 * time.Millisecond)
            out <- n * n
        }
    }()
    return out
}

// Fan-Out to multiple workers
func fanOut(in <-chan int, numWorkers int) []<-chan int {
    outputs := make([]<-chan int, numWorkers)
    for i := 0; i < numWorkers; i++ {
        outputs[i] = process(in)
    }
    return outputs
}

// Fan-In from multiple workers
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

// Stage 3: Collect results
func collect(in <-chan int) []int {
    results := make([]int, 0)
    for val := range in {
        results = append(results, val)
    }
    return results
}

func main() {
    nums := []int{1, 2, 3, 4, 5, 6, 7, 8}
    
    // Stage 1: Generate
    stage1 := generator(nums)
    
    // Stage 2: Fan-Out to 3 workers
    workers := fanOut(stage1, 3)
    
    // Stage 3: Fan-In results
    stage2 := fanIn(workers...)
    
    // Stage 4: Collect
    results := collect(stage2)
    
    fmt.Printf("Results: %v\n", results)
}
```

---

## 🔒 Pipeline з Context (Cancellation)

```go
package main

import (
    "context"
    "fmt"
    "time"
)

// Stage with context
func generator(ctx context.Context, nums []int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            select {
            case <-ctx.Done():
                fmt.Println("Generator cancelled")
                return
            case out <- n:
            }
        }
    }()
    return out
}

func square(ctx context.Context, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for {
            select {
            case <-ctx.Done():
                fmt.Println("Square cancelled")
                return
            case n, ok := <-in:
                if !ok {
                    return
                }
                
                // Simulate work
                select {
                case <-ctx.Done():
                    return
                case <-time.After(100 * time.Millisecond):
                    out <- n * n
                }
            }
        }
    }()
    return out
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
    defer cancel()
    
    nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    
    stage1 := generator(ctx, nums)
    stage2 := square(ctx, stage1)
    
    for val := range stage2 {
        fmt.Printf("Result: %d\n", val)
    }
    
    fmt.Println("Done (may be cancelled)")
}
```

---

## 🎯 Real-World Example: Log Processing Pipeline

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "time"
)

type LogEntry struct {
    Timestamp time.Time
    Level     string
    Message   string
}

// Stage 1: Read log file
func readLogs(filename string) <-chan string {
    out := make(chan string)
    go func() {
        defer close(out)
        file, err := os.Open(filename)
        if err != nil {
            return
        }
        defer file.Close()
        
        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            out <- scanner.Text()
        }
    }()
    return out
}

// Stage 2: Parse log lines
func parseLogs(in <-chan string) <-chan LogEntry {
    out := make(chan LogEntry)
    go func() {
        defer close(out)
        for line := range in {
            parts := strings.Split(line, " | ")
            if len(parts) < 3 {
                continue
            }
            
            timestamp, _ := time.Parse(time.RFC3339, parts[0])
            out <- LogEntry{
                Timestamp: timestamp,
                Level:     parts[1],
                Message:   parts[2],
            }
        }
    }()
    return out
}

// Stage 3: Filter by level
func filterByLevel(in <-chan LogEntry, level string) <-chan LogEntry {
    out := make(chan LogEntry)
    go func() {
        defer close(out)
        for entry := range in {
            if entry.Level == level {
                out <- entry
            }
        }
    }()
    return out
}

// Stage 4: Format output
func formatLogs(in <-chan LogEntry) <-chan string {
    out := make(chan string)
    go func() {
        defer close(out)
        for entry := range in {
            formatted := fmt.Sprintf("[%s] %s: %s",
                entry.Timestamp.Format("15:04:05"),
                entry.Level,
                entry.Message)
            out <- formatted
        }
    }()
    return out
}

func main() {
    // Build pipeline
    stage1 := readLogs("app.log")
    stage2 := parseLogs(stage1)
    stage3 := filterByLevel(stage2, "ERROR")
    stage4 := formatLogs(stage3)
    
    // Process results
    for formatted := range stage4 {
        fmt.Println(formatted)
    }
}
```

---

## 📊 Image Processing Pipeline

```go
package main

import (
    "fmt"
    "image"
    "os"
)

type ImageJob struct {
    Filename string
    Image    image.Image
}

// Stage 1: Load images
func loadImages(filenames []string) <-chan ImageJob {
    out := make(chan ImageJob)
    go func() {
        defer close(out)
        for _, filename := range filenames {
            file, err := os.Open(filename)
            if err != nil {
                continue
            }
            
            img, _, err := image.Decode(file)
            file.Close()
            if err != nil {
                continue
            }
            
            out <- ImageJob{Filename: filename, Image: img}
        }
    }()
    return out
}

// Stage 2: Resize images
func resizeImages(in <-chan ImageJob, width, height int) <-chan ImageJob {
    out := make(chan ImageJob)
    go func() {
        defer close(out)
        for job := range in {
            // resized := resize(job.Image, width, height)
            // job.Image = resized
            out <- job
        }
    }()
    return out
}

// Stage 3: Apply filter
func applyFilter(in <-chan ImageJob, filter string) <-chan ImageJob {
    out := make(chan ImageJob)
    go func() {
        defer close(out)
        for job := range in {
            fmt.Printf("Applying %s filter to %s\n", filter, job.Filename)
            // Apply filter logic
            out <- job
        }
    }()
    return out
}

// Stage 4: Save images
func saveImages(in <-chan ImageJob, outputDir string) <-chan string {
    out := make(chan string)
    go func() {
        defer close(out)
        for job := range in {
            outputPath := outputDir + "/" + job.Filename
            // Save image logic
            fmt.Printf("Saved to %s\n", outputPath)
            out <- outputPath
        }
    }()
    return out
}

func main() {
    filenames := []string{"img1.jpg", "img2.jpg", "img3.jpg"}
    
    // Build pipeline
    stage1 := loadImages(filenames)
    stage2 := resizeImages(stage1, 800, 600)
    stage3 := applyFilter(stage2, "sepia")
    stage4 := saveImages(stage3, "./output")
    
    // Collect results
    for path := range stage4 {
        fmt.Printf("Processed: %s\n", path)
    }
}
```

---

## 🔄 Buffered Pipeline (Performance)

```go
func bufferedGenerator(nums []int, bufSize int) <-chan int {
    out := make(chan int, bufSize)  // Buffered!
    go func() {
        defer close(out)
        for _, n := range nums {
            out <- n
        }
    }()
    return out
}

func bufferedSquare(in <-chan int, bufSize int) <-chan int {
    out := make(chan int, bufSize)  // Buffered!
    go func() {
        defer close(out)
        for n := range in {
            out <- n * n
        }
    }()
    return out
}

// Usage
nums := bufferedGenerator(data, 100)
squares := bufferedSquare(nums, 100)
```

**Benefits:**
- Reduce blocking
- Better throughput
- Smooth out processing spikes

---

## ✅ Pipeline Patterns

### Pattern 1: Linear Pipeline

```
Input → [A] → [B] → [C] → Output
```

```go
stage1 := stageA(input)
stage2 := stageB(stage1)
output := stageC(stage2)
```

### Pattern 2: Branching Pipeline

```
        → [B1] →
Input → [A] → [B2] → [C] → Output
        → [B3] →
```

```go
stageA := processA(input)
branches := fanOut(stageA, 3)
stageB := fanIn(branches...)
output := processC(stageB)
```

### Pattern 3: Diamond Pipeline

```
        → [B] →
Input → [A]     [D] → Output
        → [C] →
```

```go
stageA := processA(input)
stageB := processB(stageA)
stageC := processC(stageA)
output := merge(stageB, stageC)
```

---

## 🐛 Common Mistakes

### Mistake 1: Not Closing Channels

```go
// ❌ BAD: Channel never closes
func stage(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * 2
        }
        // Forgot close(out)!
    }()
    return out
}

// Next stage will block forever!
```

### Mistake 2: Blocking on Send

```go
// ❌ BAD: Can deadlock
func stage(in <-chan int) <-chan int {
    out := make(chan int)  // Unbuffered
    for n := range in {
        out <- n * 2  // Blocks if no receiver!
    }
    return out
}

// ✅ GOOD: Use goroutine
func stage(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            out <- n * 2
        }
    }()
    return out
}
```

### Mistake 3: Goroutine Leak on Error

```go
// ❌ BAD: Goroutine leaks if error
func stage(in <-chan int) (<-chan int, error) {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            if n < 0 {
                return  // ❌ Goroutine exits but caller doesn't know!
            }
            out <- n * 2
        }
    }()
    return out, nil
}

// ✅ GOOD: Use error channel
func stage(in <-chan int) (<-chan int, <-chan error) {
    out := make(chan int)
    errs := make(chan error, 1)
    
    go func() {
        defer close(out)
        for n := range in {
            if n < 0 {
                errs <- fmt.Errorf("negative number: %d", n)
                return
            }
            out <- n * 2
        }
    }()
    
    return out, errs
}
```

---

## ✅ Best Practices

### 1. Each Stage Owns Output Channel

```go
func stage(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)  // ✅ Stage closes its output
        // Process...
    }()
    return out
}
```

### 2. Use Context for Cancellation

```go
func stage(ctx context.Context, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for {
            select {
            case <-ctx.Done():
                return  // ✅ Respect cancellation
            case n, ok := <-in:
                if !ok {
                    return
                }
                out <- process(n)
            }
        }
    }()
    return out
}
```

### 3. Buffer для Performance

```go
// For high-throughput pipelines
out := make(chan int, 100)  // ✅ Buffered
```

### 4. Error Handling

```go
type Result struct {
    Value int
    Error error
}

func stage(in <-chan int) <-chan Result {
    out := make(chan Result)
    go func() {
        defer close(out)
        for n := range in {
            val, err := process(n)
            out <- Result{Value: val, Error: err}
        }
    }()
    return out
}
```

---

## 🎓 Висновок

### Pipeline Pattern:

✅ Series of stages  
✅ Connected by channels  
✅ Each stage = group of goroutines  
✅ Composable and reusable  

### Benefits:

- **Separation of concerns** (кожна stage - одна функція)
- **Parallelism** (stages run concurrently)
- **Composability** (легко додавати stages)
- **Testability** (кожну stage можна тестувати окремо)

### When to Use:

- Data processing (ETL)
- Log processing
- Image/video processing
- Network request processing
- Any multi-step transformation

---

## 📖 Далі

- `practice/03_pipeline/` - Практичні приклади
- `04_context_cancel.md` - Context patterns
- Combine all patterns для production systems

**"Pipeline = Composable Data Processing!" 🔄**
