# Quick Start - Ruby vs Go Concurrency 🚀

Швидкий старт для запуску всіх прикладів.

---

## ⚡ One-liner для всіх прикладів

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_25

# Запустити ВСІ Ruby приклади
for f in *.rb; do echo "=== $f ==="; ruby "$f"; echo; done

# Запустити ВСІ Go приклади
for f in *.go; do echo "=== $f ==="; go run "$f"; echo; done
```

---

## 📝 Окремі приклади

### 1. Basic Threads/Goroutines
```bash
# Ruby
ruby 01_basic_threads.rb

# Go
go run 01_basic_goroutines.go
```

**Що показує:** Створення, параметри, return values

---

### 2. Mutex Synchronization
```bash
# Ruby
ruby 02_mutex.rb

# Go
go run 02_mutex.go
```

**Що показує:** Race condition, mutex.synchronize, bank account

---

### 3. Thread Pool / Worker Pool
```bash
# Ruby
ruby 03_thread_pool.rb

# Go
go run 03_worker_pool.go
```

**Що показує:** Fixed workers, task queue, graceful shutdown

---

### 4. Producer/Consumer
```bash
# Ruby
ruby 04_producer_consumer.rb

# Go
go run 04_producer_consumer.go
```

**Що показує:** Queue/Channel pattern, multiple producers/consumers

---

### 5. Parallel HTTP Requests
```bash
# Ruby (потрібен інтернет)
ruby 05_parallel_http.rb

# Go (потрібен інтернет)
go run 05_parallel_http.go
```

**Що показує:** Concurrent HTTP, error handling, speedup

---

### 6. Timeout Pattern
```bash
# Ruby
ruby 06_timeout.rb

# Go
go run 06_timeout.go
```

**Що показує:** Timeout.timeout, context.WithTimeout, cancellation

---

### 7. Race Condition Demo
```bash
# Ruby
ruby 07_race_condition.rb

# Go
go run 07_race_condition.go

# Go with race detector
go run -race 07_race_condition.go
```

**Що показує:** Race conditions, як їх виявити та виправити

---

## 🔬 Race Detector (Go only)

```bash
# Запустити всі Go приклади з race detector
for f in *.go; do 
  echo "=== $f with race detector ===" 
  go run -race "$f" 2>&1 | head -20
  echo
done
```

---

## 📊 Performance Comparison

### CPU-bound (обчислення)
```bash
# Ruby (1 thread)
time ruby -e '1_000_000.times { Math.sqrt(rand) }'

# Ruby (4 threads)
time ruby -e '
threads = 4.times.map { Thread.new { 250_000.times { Math.sqrt(rand) } } }
threads.each(&:join)
'

# Go (1 goroutine)
time go run -e 'package main; import ("math"; "math/rand"); func main() { for i := 0; i < 1_000_000; i++ { _ = math.Sqrt(rand.Float64()) } }'

# Go (4 goroutines)
# Create temp file and run
```

**Очікуваний результат:**
- Ruby: 4 threads ≈ 1 thread (GIL!)
- Go: 4 goroutines ≈ 4x faster

---

### I/O-bound (HTTP)

```bash
# Ruby
time ruby 05_parallel_http.rb

# Go
time go run 05_parallel_http.go
```

**Очікуваний результат:**
- Ruby: добре працює (GIL releases on I/O)
- Go: трохи швидше (lightweight goroutines)

---

## 🎯 Швидкий тест

```bash
# Все в одній команді
cd /Users/vkuzm/GolandProjects/golang_practice/week_25 && \
echo "=== RUBY EXAMPLES ===" && \
ruby 01_basic_threads.rb && \
echo && \
echo "=== GO EXAMPLES ===" && \
go run 01_basic_goroutines.go && \
echo && \
echo "✅ Done! Check COMPARISON.md for details"
```

---

## 📖 Документація

- `README.md` - Загальний огляд
- `COMPARISON.md` - **Детальне порівняння Ruby vs Go**
- Коментарі в кожному файлі

---

## 💡 Tips

1. **Почніть з COMPARISON.md** - там є вся теорія
2. **Запускайте side-by-side** - Ruby vs Go для кожного pattern
3. **Експериментуйте** - змінюйте кількість threads/goroutines
4. **Race detector** - завжди використовуйте для Go (`go run -race`)
5. **Network access** - для HTTP прикладів потрібен інтернет

---

## ⚠️ Notes

### Ruby
- Потребує Ruby >= 2.7
- GIL обмежує CPU паралелізм
- Добре для I/O tasks

### Go
- Потребує Go >= 1.19
- Справжній паралелізм
- Використовуйте `-race` для виявлення проблем

---

**🎉 Готово до експериментів!**
