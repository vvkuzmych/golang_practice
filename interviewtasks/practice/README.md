# Go Concurrency Practice Tasks 🚀

Практичні завдання на **goroutines, channels, context, race conditions** для підготовки до співбесід.

---

## 📚 Table of Contents

| Task | Level | Time | Topics |
|------|-------|------|--------|
| [Task 1: Parallel Sum](#task-1-parallel-sum) | Beginner | 10 min | Goroutines, WaitGroup, Mutex |
| [Task 2: URL Checker](#task-2-url-checker) | Intermediate | 15 min | Goroutines, Channels, Error Handling |
| [Task 3: Worker Pool](#task-3-worker-pool) | Intermediate-Advanced | 20 min | Worker Pool, Channels, Context |
| [Task 4: Context Timeout](#task-4-context-timeout) | Advanced | 15 min | Context, Timeout, Cancellation |
| [Task 5: Race Condition](#task-5-race-condition) | Advanced | 15 min | Race Conditions, Mutex, Thread Safety |

---

## 🎯 Learning Objectives

Після виконання цих завдань ти будеш впевнено:

- ✅ Створювати goroutines та синхронізувати їх за допомогою `sync.WaitGroup`
- ✅ Використовувати channels для комунікації між goroutines
- ✅ Імплементувати Worker Pool pattern
- ✅ Працювати з `context` для cancellation та timeouts
- ✅ Детектувати та фіксити race conditions з `sync.Mutex` та `sync.RWMutex`
- ✅ Писати thread-safe код

---

## 📖 Task Details

### Task 1: Parallel Sum

**Опис:** Розрахунок суми slice чисел паралельно за допомогою N workers.

**Що вивчиш:**
- Goroutines
- `sync.WaitGroup`
- `sync.Mutex`
- Розбиття роботи на chunks

**Файли:**
- 📝 Task: `tasks/TASK_01_parallel_sum.md`
- ✅ Solution: `solutions/solution_01_parallel_sum.go`

---

### Task 2: URL Checker

**Опис:** Паралельна перевірка доступності списку URLs з збереженням порядку результатів.

**Що вивчиш:**
- Goroutines для паралельних HTTP запитів
- Channels для збору результатів
- Error handling в concurrent code
- Збереження порядку результатів

**Файли:**
- 📝 Task: `tasks/TASK_02_url_checker.md`
- ✅ Solution: `solutions/solution_02_url_checker.go`

---

### Task 3: Worker Pool

**Опис:** Імплементація Worker Pool pattern з фіксованою кількістю workers.

**Що вивчиш:**
- Worker Pool pattern
- Buffered channels
- Job queue
- Resource management

**Файли:**
- 📝 Task: `tasks/TASK_03_worker_pool.md`
- ✅ Solution: `solutions/solution_03_worker_pool.go`

---

### Task 4: Context Timeout

**Опис:** HTTP запити з підтримкою timeout та cancellation через context.

**Що вивчиш:**
- `context.Context`
- `context.WithTimeout`
- `context.WithCancel`
- `context.WithDeadline`
- Graceful cancellation

**Файли:**
- 📝 Task: `tasks/TASK_04_context_timeout.md`
- ✅ Solution: `solutions/solution_04_context_timeout.go`

---

### Task 5: Race Condition

**Опис:** Детекція та виправлення race conditions в багатопотоковому коді.

**Що вивчиш:**
- Race conditions
- `sync.Mutex`
- `sync.RWMutex`
- Race detector (`go test -race`)
- Thread-safe patterns

**Файли:**
- 📝 Task: `tasks/TASK_05_race_condition.md`
- ✅ Solution: `solutions/solution_05_race_condition.go`

---

## 🚀 How to Use

### 1. Read the Task

```bash
cat tasks/TASK_01_parallel_sum.md
```

### 2. Try to Solve It Yourself

Create your own solution:

```bash
touch my_solution_01.go
```

### 3. Run Your Solution

```bash
go run my_solution_01.go
```

### 4. Check the Official Solution

```bash
cat solutions/solution_01_parallel_sum.go
go run solutions/solution_01_parallel_sum.go
```

### 5. Run with Race Detector (Important!)

```bash
go run -race solutions/solution_05_race_condition.go
```

---

## 🧪 Running Tests

Some solutions have tests. Run them:

```bash
# Run all tests
go test ./...

# Run tests with race detector
go test -race ./...

# Run specific test
go test -v -run TestConcurrentIncrements
```

---

## 📊 Difficulty Progression

```
Beginner          Intermediate       Advanced
   ↓                   ↓                 ↓
Task 1  ───────→  Task 2  ───────→  Task 4
                     ↓                  ↓
                  Task 3  ───────→  Task 5
```

**Recommended order:**
1. Task 1 (основи goroutines + sync)
2. Task 2 (channels + error handling)
3. Task 3 (worker pool pattern)
4. Task 4 (context для production)
5. Task 5 (race conditions - must know!)

---

## 💡 Tips

### For Beginners

- Почни з Task 1
- Використовуй `fmt.Println` для дебагу
- Не поспішай з channels, спочатку освой WaitGroup

### For Intermediate

- Завжди використовуй `defer wg.Done()`
- Закривай channels після надсилання всіх даних
- Пам'ятай про deadlocks (channel без читача)

### For Advanced

- Завжди запускай `go test -race`
- Використовуй `context` для cancellation в production
- Вивчи різницю між `Mutex` та `RWMutex`

---

## 🔍 Common Mistakes

### ❌ Mistake 1: Forgetting WaitGroup

```go
// ❌ Bad
for i := 0; i < 10; i++ {
    go doWork()
}
// Program exits before goroutines finish
```

```go
// ✅ Good
var wg sync.WaitGroup
wg.Add(10)
for i := 0; i < 10; i++ {
    go func() {
        defer wg.Done()
        doWork()
    }()
}
wg.Wait()
```

### ❌ Mistake 2: Race Condition

```go
// ❌ Bad
counter := 0
for i := 0; i < 10; i++ {
    go func() {
        counter++  // Race condition!
    }()
}
```

```go
// ✅ Good
var mu sync.Mutex
counter := 0
for i := 0; i < 10; i++ {
    go func() {
        mu.Lock()
        counter++
        mu.Unlock()
    }()
}
```

### ❌ Mistake 3: Goroutine Loop Variable Capture

```go
// ❌ Bad
for i := 0; i < 10; i++ {
    go func() {
        fmt.Println(i)  // All goroutines see i=10
    }()
}
```

```go
// ✅ Good
for i := 0; i < 10; i++ {
    go func(idx int) {
        fmt.Println(idx)
    }(i)  // Pass i as parameter
}
```

---

## 📚 Additional Resources

### Official Go Documentation

- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)
- [Go Tour - Concurrency](https://go.dev/tour/concurrency/1)
- [Context Package](https://pkg.go.dev/context)

### Articles

- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Advanced Go Concurrency Patterns](https://go.dev/blog/io2013-talk-concurrency)
- [Share Memory By Communicating](https://go.dev/blog/codelab-share)

### Videos

- [Concurrency is not Parallelism by Rob Pike](https://www.youtube.com/watch?v=oV9rvDllKEg)

---

## 🎯 Interview Questions

Після вивчення цих tasks ти зможеш відповісти на:

1. **Що таке goroutine і чим вона відрізняється від thread?**
2. **Як працює channel? Що таке buffered vs unbuffered?**
3. **Що таке race condition і як його детектувати?**
4. **Коли використовувати Mutex vs RWMutex?**
5. **Як працює context для cancellation?**
6. **Що таке Worker Pool і навіщо він потрібен?**
7. **Як уникнути goroutine leaks?**
8. **Що таке WaitGroup і коли його використовувати?**

---

## 🏆 Next Steps

Після завершення цих tasks:

1. ✅ Створи свої власні варіації завдань
2. ✅ Поглибся в [Advanced Concurrency Patterns](/Users/vkuzm/GolandProjects/golang_practice/interviewtasks/main.go)
3. ✅ Вивчи [Context Patterns](https://go.dev/blog/context)
4. ✅ Попрактикуй на [exercism.io](https://exercism.io/tracks/go)

---

## 📞 Support

Маєш питання або знайшов помилку?

- Перечитай task description
- Подивись на solution code
- Запусти з `-race` для детекції race conditions
- Подивись на [офіційну документацію](https://go.dev/doc/)

---

**Happy coding!** 🚀

*Remember: The best way to learn concurrency is to write buggy code, find the bug, and fix it yourself.* 😄
