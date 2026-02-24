# Ruby Threads vs Go Goroutines - Детальне порівняння 🔬

---

## 📊 Основні відмінності

| Характеристика | Ruby Threads | Go Goroutines |
|----------------|--------------|---------------|
| **Тип** | OS threads (системні) | User-space threads (зелені) |
| **Вага** | ~1MB стек на thread | ~2KB стек на goroutine |
| **Максимальна кількість** | Сотні-тисячі | Мільйони (10M+) |
| **Швидкість створення** | Повільно (~ms) | Дуже швидко (~μs) |
| **Контекстний перемикач** | Дорого (OS рівень) | Дешево (runtime рівень) |
| **GIL (Global Lock)** | ✅ Так (MRI Ruby) | ❌ Ні |
| **CPU паралелізм** | ❌ Обмежений GIL | ✅ Повний |
| **I/O паралелізм** | ✅ Так | ✅ Так |

---

## 🧵 1. Створення та запуск

### Ruby
```ruby
# Створення thread
thread = Thread.new do
  puts "Hello"
end

# З параметрами
thread = Thread.new(5, "test") do |num, text|
  puts "#{num}: #{text}"
end

# Очікування завершення
thread.join

# Отримання результату
result = thread.value
```

### Go
```go
// Створення goroutine
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    fmt.Println("Hello")
}()

// З параметрами
wg.Add(1)
go func(num int, text string) {
    defer wg.Done()
    fmt.Printf("%d: %s\n", num, text)
}(5, "test")

// Очікування завершення
wg.Wait()

// Результат через channel
result := make(chan int)
go func() {
    result <- 42
}()
value := <-result
```

**Висновок:** Go вимагає явного керування через `WaitGroup` або channels, Ruby простіший у використанні.

---

## 🔒 2. Синхронізація (Mutex)

### Ruby
```ruby
mutex = Mutex.new
counter = 0

mutex.synchronize do
  counter += 1
end

# Або manually
mutex.lock
counter += 1
mutex.unlock
```

### Go
```go
var mu sync.Mutex
counter := 0

mu.Lock()
counter++
mu.Unlock()

// З defer (рекомендовано)
func increment() {
    mu.Lock()
    defer mu.Unlock()
    counter++
}
```

**Висновок:** Схожі підходи, але Go частіше використовує channels замість mutex.

---

## 📡 3. Комунікація між threads/goroutines

### Ruby (Shared Memory)
```ruby
# Queue (thread-safe)
queue = Queue.new
queue << "task"
task = queue.pop

# Shared variable (потрібен Mutex!)
mutex = Mutex.new
data = []

mutex.synchronize do
  data << "item"
end
```

### Go (Channels - CSP)
```go
// Buffered channel
ch := make(chan string, 10)
ch <- "task"
task := <-ch

// Unbuffered (synchronous)
ch := make(chan string)
go func() {
    ch <- "data"
}()
data := <-ch

// Select для multiple channels
select {
case msg := <-ch1:
    fmt.Println(msg)
case msg := <-ch2:
    fmt.Println(msg)
}
```

**Висновок:** Go channels — більш потужний та виразний механізм. Ruby покладається на shared memory + Mutex.

---

## 🏊 4. Thread Pool / Worker Pool

### Ruby
```ruby
queue = Queue.new
workers = 5.times.map do |i|
  Thread.new do
    loop do
      task = queue.pop
      break if task == :stop
      # Process task
    end
  end
end

# Add tasks
queue << "task1"
queue << "task2"

# Stop workers
5.times { queue << :stop }
workers.each(&:join)
```

### Go
```go
tasks := make(chan string, 10)
var wg sync.WaitGroup

// Start workers
for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        for task := range tasks {
            // Process task
        }
    }(i)
}

// Add tasks
tasks <- "task1"
tasks <- "task2"
close(tasks) // Signal workers to stop

wg.Wait()
```

**Висновок:** Go простіше — `close(ch)` сигналізує всім workers одразу. Ruby потребує sentinel значення (`:stop`).

---

## 🌐 5. HTTP Паралельні запити

### Ruby
```ruby
urls = ["url1", "url2", "url3"]

threads = urls.map do |url|
  Thread.new(url) do |u|
    Net::HTTP.get_response(URI(u))
  end
end

responses = threads.map(&:value)
```

### Go
```go
urls := []string{"url1", "url2", "url3"}
results := make(chan *http.Response, len(urls))

for _, url := range urls {
    go func(u string) {
        resp, _ := http.Get(u)
        results <- resp
    }(url)
}

// Collect results
for i := 0; i < len(urls); i++ {
    resp := <-results
    // Use resp
}
```

**Висновок:** Обидва добре працюють з I/O. Go швидший для великої кількості запитів.

---

## ⏱️ 6. Timeout / Cancellation

### Ruby
```ruby
require 'timeout'

begin
  Timeout.timeout(5) do
    # Long operation
  end
rescue Timeout::Error
  puts "Timed out!"
end

# Або з Thread
thread = Thread.new { sleep 10 }
thread.join(5) || thread.kill
```

### Go
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

done := make(chan bool)
go func() {
    // Long operation
    done <- true
}()

select {
case <-done:
    fmt.Println("Completed")
case <-ctx.Done():
    fmt.Println("Timed out!")
}
```

**Висновок:** Go's `context` — більш потужний, підтримує cancellation propagation через всю програму.

---

## 🏎️ 7. Performance Comparison

### CPU-bound задачі (Обчислення)

```ruby
# Ruby - SLOW (GIL blocks parallel CPU)
threads = 4.times.map do
  Thread.new do
    1_000_000.times { Math.sqrt(rand) }
  end
end
threads.each(&:join)
# ~4s (same as single thread!)
```

```go
// Go - FAST (true parallelism)
var wg sync.WaitGroup
for i := 0; i < 4; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        for j := 0; j < 1_000_000; j++ {
            _ = math.Sqrt(rand.Float64())
        }
    }()
}
wg.Wait()
// ~1s (4x faster on 4 cores!)
```

**Результат:**
- **Ruby:** GIL не дає використовувати multiple CPU cores
- **Go:** Справжній паралелізм, використовує всі cores

---

### I/O-bound задачі (HTTP, DB, Files)

```ruby
# Ruby - GOOD (GIL releases during I/O)
threads = 100.times.map do
  Thread.new { Net::HTTP.get(URI("...")) }
end
threads.each(&:join)
# ~2s
```

```go
// Go - EXCELLENT (lightweight goroutines)
var wg sync.WaitGroup
for i := 0; i < 100; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        http.Get("...")
    }()
}
wg.Wait()
// ~1.5s
```

**Результат:**
- **Ruby:** Добре працює, GIL не блокує I/O
- **Go:** Трохи швидше через легковісні goroutines

---

## 🔢 8. Масштабованість

### Ruby Threads
```ruby
# Maximum ~5,000-10,000 threads before OOM
threads = []
10_000.times do
  threads << Thread.new { sleep 1 }
end
threads.each(&:join)
# Memory: ~10GB
```

### Go Goroutines
```ruby
# Can handle 1,000,000+ goroutines
var wg sync.WaitGroup
for i := 0; i < 1_000_000; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        time.Sleep(1 * time.Second)
    }()
}
wg.Wait()
# Memory: ~2-3GB
```

**Результат:**
- **Ruby:** ~10K max threads
- **Go:** ~1M+ goroutines

---

## 🛠️ 9. Debugging та Tools

### Ruby
```ruby
# Thread list
Thread.list

# Thread status
thread.status  # "run", "sleep", false, nil

# Thread backtrace
thread.backtrace

# Kill thread (dangerous!)
thread.kill
```

### Go
```bash
# Race detector
go run -race program.go

# CPU profiling
go run -cpuprofile=cpu.prof program.go

# Goroutine dump
GODEBUG=schedtrace=1000 ./program

# pprof for goroutine analysis
import _ "net/http/pprof"
```

**Висновок:** Go має набагато кращі tools для debugging concurrency (race detector, pprof).

---

## ✅ Коли використовувати Ruby Threads

1. **I/O операції** (HTTP requests, Database queries, File I/O)
2. **Background jobs** (Sidekiq, delayed_job)
3. **Web servers** (Puma використовує threads)
4. **Concurrent API calls**
5. **Невеликий масштаб** (до 100 threads)

**Не підходить для:**
- CPU-intensive обчислення (використовуйте процеси або Go)
- Massive concurrency (мільйони requests)

---

## ✅ Коли використовувати Go Goroutines

1. **Все, що підходить для Ruby threads** +
2. **CPU-intensive задачі** (справжній паралелізм)
3. **Massive concurrency** (мільйони goroutines)
4. **Real-time systems** (trading, gaming, streaming)
5. **Microservices** (багато concurrent requests)
6. **Network servers** (TCP/UDP)

**Підходить для всього!**

---

## 🎯 Підсумок

| Use Case | Ruby Threads | Go Goroutines | Переможець |
|----------|--------------|---------------|------------|
| **I/O операції** | ✅ Добре | ✅ Відмінно | Go |
| **CPU обчислення** | ❌ Погано (GIL) | ✅ Відмінно | Go |
| **Простота коду** | ✅ Простіше | ⚠️ Складніше | Ruby |
| **Масштабованість** | ⚠️ До 10K | ✅ До 1M+ | Go |
| **Web servers** | ✅ Puma | ✅ net/http | Обидва |
| **Debugging** | ⚠️ Обмежено | ✅ Чудові tools | Go |
| **Ecosystem** | ✅ Rails, gems | ✅ Stdlib | Обидва |

---

## 💡 Рекомендації

### Використовуйте Ruby якщо:
- Ви вже працюєте з Rails/Ruby
- Задачі в основному I/O-bound
- Кількість concurrent operations < 1000
- Простота важливіша за raw performance

### Використовуйте Go якщо:
- Потрібен справжній CPU паралелізм
- Massive concurrency (1000s-1M+ operations)
- Performance критичний
- Real-time або low-latency системи
- Microservices

---

## 📚 Resources

### Ruby
- [Ruby Thread docs](https://ruby-doc.org/core/Thread.html)
- [Mutex docs](https://ruby-doc.org/core/Mutex.html)
- [Queue docs](https://ruby-doc.org/stdlib/libdoc/thread/rdoc/Queue.html)

### Go
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)
- [sync package](https://pkg.go.dev/sync)
- [context package](https://pkg.go.dev/context)

---

**🎉 Обидві мови мають свої сильні сторони!**
- **Ruby:** Простота, Rails ecosystem, I/O tasks
- **Go:** Performance, масштабованість, CPU tasks
