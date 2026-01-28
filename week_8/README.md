# Week 8: Debugging & Race Conditions

## 🎯 Мета

Навчитися знаходити та виправляти **race conditions** і **goroutine leaks** в concurrent Go коді.

---

## 📚 Теорія

### 1. Race Conditions
**Файл:** `theory/01_race_conditions.md`

- Що таке race condition?
- Як знайти з `go test -race`?
- Рішення: Mutex, Atomic, Channels
- Best practices

### 2. Goroutine Leaks
**Файл:** `theory/02_goroutine_leaks.md`

- Що таке goroutine leak?
- Типи leaks (channel, mutex, waitgroup)
- Як знайти (pprof, goleak)
- Patterns для clean shutdown

---

## 💻 Практика

### 1. Race Conditions
**Директорія:** `practice/01_race_conditions/`

**Broken (з багами):**
- `broken_counter.go` - Race на counter
- `broken_map.go` - Race на map
- `broken_slice.go` - Race на slice

**Fixed (виправлено):**
- `fixed_counter.go` - 3 способи: Mutex, Atomic, Channel
- `fixed_map.go` - sync.Map
- `fixed_slice.go` - Mutex

**Як запускати:**

```bash
cd practice/01_race_conditions

# Знайти race condition
go run -race broken_counter.go

# Перевірити fix
go run fixed_counter.go
```

### 2. Goroutine Leaks
**Директорія:** `practice/02_goroutine_leaks/`

**Broken (з багами):**
- `broken_channel_leak.go` - Channel leaks
- `broken_waitgroup_leak.go` - WaitGroup leak
- `broken_http_leak.go` - HTTP request leak

**Fixed (виправлено):**
- `fixed_channel_leak.go` - 6 способів fix
- `fixed_waitgroup_leak.go` - Proper WaitGroup usage
- `fixed_http_leak.go` - Context & timeout

**Як запускати:**

```bash
cd practice/02_goroutine_leaks

# Показати leak
go run broken_channel_leak.go

# Перевірити fix
go run fixed_channel_leak.go
```

---

## 🔧 Debugging Tools

### 1. Race Detector

```bash
# Build з race detector
go build -race

# Test з race detector
go test -race ./...

# Run з race detector
go run -race main.go
```

**⚠️ Важливо:** Race detector має ~10x overhead. НЕ використовуй в production!

### 2. Check Goroutine Count

```go
import "runtime"

before := runtime.NumGoroutine()
// ... run code ...
after := runtime.NumGoroutine()

if after > before {
    fmt.Printf("Leaked %d goroutines!\n", after-before)
}
```

### 3. pprof для Goroutines

```go
import _ "net/http/pprof"

func main() {
    go http.ListenAndServe("localhost:6060", nil)
    // Your code
}
```

```bash
# Check goroutines
go tool pprof http://localhost:6060/debug/pprof/goroutine

# In pprof:
top      # Top functions
list     # Source code
traces   # Stack traces
```

### 4. goleak (Uber's tool)

```go
import "go.uber.org/goleak"

func TestMain(m *testing.M) {
    goleak.VerifyTestMain(m)
}

func TestSomething(t *testing.T) {
    defer goleak.VerifyNone(t)
    // Test code
}
```

---

## 📊 Workflow: Знайти → Виправити

### Step 1: Запусти з `-race`

```bash
go test -race ./...
# або
go run -race main.go
```

### Step 2: Проаналізуй Output

```
==================
WARNING: DATA RACE
Read at 0x... by goroutine 7:
  main.Counter.Get()
      counter.go:15 +0x38

Previous write at 0x... by goroutine 6:
  main.Counter.Increment()
      counter.go:11 +0x4e
==================
```

**Шукай:**
- Line numbers (counter.go:15, counter.go:11)
- Operation (Read, Write)
- Goroutine IDs

### Step 3: Вибери Fix

**Для race conditions:**
- Mutex → General purpose
- Atomic → Simple operations (counter, flag)
- Channel → Communication

**Для goroutine leaks:**
- Context → Cancellation
- Done channel → Shutdown signal
- Timeout → Prevent forever blocking
- Close channel → Unblock receivers

### Step 4: Перевір Fix

```bash
# Race detector не повинен скаржитись
go test -race ./...

# Goroutines не повинні рости
runtime.NumGoroutine()
```

---

## ✅ Чеклист: Безпечний Concurrent Code

### Race Conditions

- [ ] Всі shared variables захищені (mutex/atomic/channel)
- [ ] Maps використовують `sync.Map` або mutex
- [ ] Slices protected під час append
- [ ] Struct fields мають mutex якщо потрібен concurrent access
- [ ] `go test -race` проходить без помилок

### Goroutine Leaks

- [ ] Кожна goroutine має exit strategy
- [ ] Використовується context для cancellation
- [ ] Channels закриваються (producer side)
- [ ] WaitGroup має `defer wg.Done()`
- [ ] HTTP requests мають timeout
- [ ] `runtime.NumGoroutine()` не росте

---

## 🎯 Практичні Завдання

### Завдання 1: Знайди Race Condition

1. Запусти `broken_counter.go` з `-race`
2. Проаналізуй output
3. Знайди проблемну лінію
4. Подивись на `fixed_counter.go`
5. Зрозумій чому кожен fix працює

### Завдання 2: Знайди Goroutine Leak

1. Запусти `broken_channel_leak.go`
2. Подивись на `runtime.NumGoroutine()`
3. Зрозумій чому goroutines не виходять
4. Подивись на `fixed_channel_leak.go`
5. Запусти і перевір що leaks виправлені

### Завдання 3: Власний Код

Створи свій код з:
1. Race condition на map
2. Goroutine leak на channel
3. Виправ обидва

---

## 📈 Performance Comparison

```bash
cd practice/01_race_conditions
go test -bench=. -benchmem
```

**Expected results:**

```
BenchmarkMutex-8       10000000    120 ns/op
BenchmarkAtomic-8      50000000     30 ns/op
BenchmarkChannel-8      5000000    300 ns/op
```

**Висновок:**
- Atomic найшвидший
- Mutex універсальний
- Channel для communication

---

## 🐛 Common Bugs

### Bug 1: Forgotten Lock

```go
// ❌ BAD
func (c *Counter) Get() int {
    return c.value  // Race!
}

// ✅ GOOD
func (c *Counter) Get() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}
```

### Bug 2: Defer в Loop

```go
// ❌ BAD: Defers accumulate
for i := 0; i < 1000; i++ {
    mu.Lock()
    defer mu.Unlock()  // ❌ Won't unlock until function exits!
    // work
}

// ✅ GOOD: Unlock in loop
for i := 0; i < 1000; i++ {
    mu.Lock()
    // work
    mu.Unlock()
}
```

### Bug 3: Blocking Forever

```go
// ❌ BAD: No exit
ch := make(chan int)
go func() {
    <-ch  // Blocks forever!
}()

// ✅ GOOD: Can exit
ch := make(chan int)
done := make(chan struct{})
go func() {
    select {
    case <-ch:
    case <-done:
        return
    }
}()
```

---

## 🎓 Key Takeaways

### Race Conditions

1. **Always test with `-race`**
2. **Protect shared state** (mutex/atomic/channel)
3. **Keep critical sections small**
4. **Prefer channels for communication**

### Goroutine Leaks

1. **Every goroutine needs exit strategy**
2. **Use context for cancellation**
3. **Close channels (producer side)**
4. **Monitor goroutine count**

---

## 📖 Ресурси

### Документація

- [Race Detector](https://go.dev/doc/articles/race_detector)
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [goleak](https://github.com/uber-go/goleak)

### Tools

```bash
# Race detector
go test -race

# Goroutine profiling
go tool pprof http://localhost:6060/debug/pprof/goroutine

# Leak detection
go get -u go.uber.org/goleak
```

---

## 🚀 Quick Start

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_8

# Читати теорію
cat theory/01_race_conditions.md
cat theory/02_goroutine_leaks.md

# Запускати практику
cd practice/01_race_conditions
go run -race broken_counter.go
go run fixed_counter.go

cd ../02_goroutine_leaks
go run broken_channel_leak.go
go run fixed_channel_leak.go
```

---

**"The two hardest problems in concurrent programming: race conditions and race conditions." 🐛**

**Status:** Week 8 Materials Complete ✅  
**Created:** 2026-01-28
