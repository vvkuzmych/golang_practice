# Exercise 2: Worker Pool

## 🎯 Мета

Створити **worker pool** для паралельної обробки завдань з обмеженим числом одночасних workers.

---

## 📋 Завдання

Реалізуйте worker pool з наступними характеристиками:

### Параметри:
- **Workers:** 5 workers
- **Jobs:** 100 jobs для обробки
- **Processing:** Кожен job - це число, результат = `job * job` (квадрат)
- **Error handling:** 10% jobs мають "fail" (повертати помилку)

### Вимоги:

- ✅ Створити `Job` struct з полем `ID int`
- ✅ Створити `Result` struct з полями `JobID int`, `Value int`, `Error error`
- ✅ Використати buffered channels для jobs та results
- ✅ Використати `WaitGroup` для синхронізації workers
- ✅ Зібрати та показати статистику: успішних, помилкових, total results

---

## 💡 Підказки

### Struct Definitions:

```go
type Job struct {
    ID int
}

type Result struct {
    JobID  int
    Value  int
    Error  error
}
```

### Worker Function:

```go
func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
    defer wg.Done()
    
    for job := range jobs {
        // Симулювати помилку (10% jobs)
        if job.ID%10 == 0 {
            results <- Result{
                JobID: job.ID,
                Error: fmt.Errorf("processing failed for job %d", job.ID),
            }
            continue
        }
        
        // Успішна обробка
        time.Sleep(10 * time.Millisecond) // Симуляція роботи
        results <- Result{
            JobID: job.ID,
            Value: job.ID * job.ID,
        }
    }
}
```

### Main Function Structure:

```go
func main() {
    const numJobs = 100
    const numWorkers = 5
    
    jobs := make(chan Job, numJobs)
    results := make(chan Result, numJobs)
    var wg sync.WaitGroup
    
    // 1. Start workers
    for w := 1; w <= numWorkers; w++ {
        wg.Add(1)
        go worker(w, jobs, results, &wg)
    }
    
    // 2. Send jobs
    go func() {
        for j := 1; j <= numJobs; j++ {
            jobs <- Job{ID: j}
        }
        close(jobs)
    }()
    
    // 3. Close results after all workers finish
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // 4. Collect results
    successCount := 0
    errorCount := 0
    for result := range results {
        if result.Error != nil {
            errorCount++
        } else {
            successCount++
        }
    }
    
    // 5. Print statistics
    fmt.Printf("Total: %d\n", numJobs)
    fmt.Printf("Success: %d\n", successCount)
    fmt.Printf("Errors: %d\n", errorCount)
}
```

---

## 🎓 Ключові концепції

1. **Worker Pool** - обмежена кількість concurrent workers
2. **Buffered channels** - для jobs та results
3. **WaitGroup** - синхронізація workers
4. **Error handling** - через Result struct
5. **Channel closure** - коректне закриття jobs і results

---

## ✅ Критерії успіху

- [ ] 5 workers обробляють 100 jobs
- [ ] Результат для кожного job: `job.ID * job.ID`
- [ ] ~10 jobs з помилками (кожен 10-й: 10, 20, 30, ...)
- [ ] Всі результати зібрані (без втрат!)
- [ ] Статистика виведена коректно
- [ ] Немає deadlock

---

## 🚀 Очікуваний результат

```
Worker 1: processing job 1
Worker 2: processing job 2
Worker 3: processing job 3
...
Worker 5: processing job 100

=== Statistics ===
Total:   100
Success: 90
Errors:  10
```

---

## 🔥 Бонус (опціонально)

### Бонус 1: Progress Bar
Додайте прогрес-бар який показує скільки jobs оброблено:

```
Processing: [=========>        ] 54/100 (54%)
```

**Підказка:** Використайте окрему goroutine для моніторингу results channel.

### Бонус 2: Rate Limiting
Обмежте кількість одночасних операцій до **2** (незалежно від кількості workers):

```go
semaphore := make(chan struct{}, 2)  // Max 2 concurrent

// В worker:
semaphore <- struct{}{}  // Acquire
// ... робота ...
<-semaphore  // Release
```

### Бонус 3: Retry Failed Jobs
Додайте механізм retry для jobs з помилками (максимум 3 спроби):

```go
type Job struct {
    ID      int
    Retries int  // Кількість спроб
}

// В worker:
if shouldFail && job.Retries < 3 {
    job.Retries++
    jobs <- job  // Відправляємо назад для retry
}
```

### Бонус 4: Dynamic Worker Scaling
Додайте можливість динамічно збільшувати/зменшувати кількість workers:

```go
addWorkerCh := make(chan bool)
removeWorkerCh := make(chan bool)

// Додати worker
addWorkerCh <- true

// Видалити worker
removeWorkerCh <- true
```

---

## 📚 Корисні посилання

- Theory: `week_5/theory/02_channels.md` - buffered channels
- Practice: `week_5/practice/worker_pool/main.go` - приклади worker pool
- Solution: `week_5/solutions/solution_2.go` (після виконання)

---

**Удачі! 🎉**

**Час виконання:** 45-60 хвилин
