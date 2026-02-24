# Week 25 - Ruby Threads vs Go Goroutines 🔄

Порівняння concurrency в Ruby та Go.

---

## 📁 Структура

```
week_25/
├── 01_basic_threads.rb          Ruby: Базові threads
├── 01_basic_goroutines.go       Go: Базові goroutines
├── 02_mutex.rb                  Ruby: Mutex синхронізація
├── 02_mutex.go                  Go: Mutex синхронізація
├── 03_thread_pool.rb            Ruby: Thread pool
├── 03_worker_pool.go            Go: Worker pool (channels)
├── 04_producer_consumer.rb      Ruby: Queue pattern
├── 04_producer_consumer.go      Go: Channel pattern
├── 05_parallel_http.rb          Ruby: Parallel requests
├── 05_parallel_http.go          Go: Parallel requests
├── 06_timeout.rb                Ruby: Timeout pattern
├── 06_timeout.go                Go: Context timeout
├── 07_race_condition.rb         Ruby: Race condition demo
├── 07_race_condition.go         Go: Race condition demo
└── COMPARISON.md                 Детальне порівняння
```

---

## 🚀 Запуск

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_25

# Ruby examples
ruby 01_basic_threads.rb
ruby 03_thread_pool.rb

# Go examples
go run 01_basic_goroutines.go
go run 03_worker_pool.go
```

---

## ⚖️ Ruby vs Go - Ключові відмінності

| Feature | Ruby Threads | Go Goroutines |
|---------|--------------|---------------|
| **Вага** | OS threads (важкі) | Lightweight (легкі) |
| **Кількість** | Сотні-тисячі | Мільйони |
| **GIL** | ✅ Є (обмежує CPU) | ❌ Немає |
| **Синхронізація** | Mutex, Queue | Channels, Mutex |
| **Комунікація** | Shared memory | Channels (CSP) |
| **Створення** | `Thread.new` | `go func()` |
| **CPU-bound** | ❌ Погано (GIL) | ✅ Добре |
| **I/O-bound** | ✅ Добре | ✅ Відмінно |

---

## 🎯 Коли використовувати

### Ruby Threads ✅
- I/O операції (HTTP, DB, files)
- Background tasks
- Concurrent requests
- Web servers (Puma)

### Ruby Threads ❌
- CPU-intensive задачі → використовуйте **процеси**

### Go Goroutines ✅
- Все вище +
- CPU-intensive задачі
- Massive concurrency (мільйони)
- Real-time systems

---

## 📚 Patterns

### 7 основних patterns:
1. **Basic** - Створення і запуск
2. **Mutex** - Синхронізація доступу
3. **Pool** - Обмежена кількість workers
4. **Producer/Consumer** - Queue/Channel pattern
5. **Parallel HTTP** - Concurrent requests
6. **Timeout** - Обмеження часу
7. **Race Condition** - Демо проблеми

**Кожен pattern в Ruby і Go для порівняння!**

---

## 💡 Quick Examples

### Ruby
```ruby
# Thread
thread = Thread.new { puts "Hello" }
thread.join

# Mutex
mutex = Mutex.new
mutex.synchronize { counter += 1 }

# Queue
queue = Queue.new
queue << "task"
task = queue.pop
```

### Go
```go
// Goroutine
go func() { fmt.Println("Hello") }()

// Mutex
var mu sync.Mutex
mu.Lock()
counter++
mu.Unlock()

// Channel
ch := make(chan string)
ch <- "task"
task := <-ch
```

---

**14 Examples + Detailed Comparison!** 🎉
