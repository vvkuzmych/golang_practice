# Task 3: Worker Pool Pattern

**Level:** Intermediate-Advanced  
**Time:** 20 minutes  
**Topics:** Worker Pool, Channels, Context

---

## 📝 Task

Імплементуй **Worker Pool** pattern для обробки jobs паралельно з обмеженою кількістю workers.

Worker pool має фіксовану кількість goroutines, які обробляють jobs з черги.

---

## 📥 Function Signature

```go
type Job struct {
    ID   int
    Data string
}

type Result struct {
    JobID  int
    Output string
    Error  error
}

// ProcessFunc - функція для обробки job
type ProcessFunc func(job Job) (string, error)

func WorkerPool(jobs []Job, numWorkers int, process ProcessFunc) []Result
```

**Parameters:**
- `jobs` - список jobs для обробки
- `numWorkers` - кількість worker goroutines
- `process` - функція обробки job

**Returns:**
- `[]Result` - результати в **тому самому порядку** що jobs

---

## 💡 Examples

```go
// Simple processing function
processFunc := func(job Job) (string, error) {
    // Simulate work
    time.Sleep(100 * time.Millisecond)
    return strings.ToUpper(job.Data), nil
}

jobs := []Job{
    {ID: 1, Data: "hello"},
    {ID: 2, Data: "world"},
    {ID: 3, Data: "golang"},
}

results := WorkerPool(jobs, 2, processFunc)

// results[0] = Result{JobID: 1, Output: "HELLO", Error: nil}
// results[1] = Result{JobID: 2, Output: "WORLD", Error: nil}
// results[2] = Result{JobID: 3, Output: "GOLANG", Error: nil}
```

---

## ✅ Requirements

- Створи **фіксовану кількість workers** (не по worker на кожен job!)
- Використай **buffered channel** для jobs queue
- Використай **channel** для збору results
- Збережи порядок результатів
- Обробляй помилки (якщо `process` повертає error)
- Workers мають завершитись після обробки всіх jobs

---

## 🏗️ Architecture

```
                   Jobs Channel
                  (buffered queue)
                        │
        ┌───────────────┼───────────────┐
        │               │               │
   [Worker 1]      [Worker 2]      [Worker 3]
        │               │               │
        └───────────────┼───────────────┘
                        │
                  Results Channel
                        ↓
                   Main Thread
```

---

## 🧪 Test Cases

```go
// Test 1: Normal processing
processFunc := func(job Job) (string, error) {
    return strings.ToUpper(job.Data), nil
}
jobs := []Job{{ID: 1, Data: "hello"}, {ID: 2, Data: "world"}}
results := WorkerPool(jobs, 2, processFunc)
assert.Equal(t, "HELLO", results[0].Output)
assert.Equal(t, "WORLD", results[1].Output)

// Test 2: Processing with error
processFunc := func(job Job) (string, error) {
    if job.Data == "error" {
        return "", fmt.Errorf("processing failed")
    }
    return job.Data, nil
}
jobs := []Job{{ID: 1, Data: "ok"}, {ID: 2, Data: "error"}}
results := WorkerPool(jobs, 2, processFunc)
assert.Nil(t, results[0].Error)
assert.NotNil(t, results[1].Error)

// Test 3: More jobs than workers
jobs := make([]Job, 100)
for i := range jobs {
    jobs[i] = Job{ID: i, Data: fmt.Sprintf("job%d", i)}
}
results := WorkerPool(jobs, 5, processFunc)  // 5 workers, 100 jobs
assert.Equal(t, 100, len(results))

// Test 4: More workers than jobs
jobs := []Job{{ID: 1, Data: "test"}}
results := WorkerPool(jobs, 10, processFunc)  // 10 workers, 1 job
assert.Equal(t, 1, len(results))

// Test 5: Empty jobs
jobs := []Job{}
results := WorkerPool(jobs, 3, processFunc)
assert.Equal(t, 0, len(results))
```

---

## 💡 Hints

### Approach 1: Jobs Channel + Results Channel

```go
1. Створити buffered channel для jobs
2. Створити channel для results
3. Запустити N workers (кожен читає з jobs channel)
4. Надіслати всі jobs в jobs channel
5. Закрити jobs channel (workers зупиняться після обробки всіх)
6. Зібрати results з results channel
7. Відновити порядок (використати map з JobID)
```

### Approach 2: WaitGroup + Results Slice

```go
1. Створити buffered channel для jobs
2. Створити slice для results
3. Запустити N workers з WaitGroup
4. Надіслати всі jobs
5. wg.Wait() для очікування
6. Results вже в правильному порядку
```

---

## 🎯 Challenge (Bonus)

Додай можливість **cancellation** через `context.Context`:

```go
func WorkerPoolWithContext(ctx context.Context, jobs []Job, numWorkers int, process ProcessFunc) ([]Result, error)
```

Якщо `ctx.Done()`, всі workers мають зупинитись і повернути необроблені jobs.

---

**Рішення:** `solutions/solution_03_worker_pool.go`

**Good luck!** 🚀
