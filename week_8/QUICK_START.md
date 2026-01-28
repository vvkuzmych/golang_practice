# Week 8 - Quick Start 🚀

## 🎯 Мета

Навчитися знаходити та виправляти **race conditions** і **goroutine leaks**.

---

## ⚡ 5-хвилинний старт

### 1. Race Conditions

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_8/practice/01_race_conditions

# Запусти broken code з race detector
go run -race broken_counter.go
```

**Що побачиш:**

```
WARNING: DATA RACE
Read at 0x... by goroutine 7:
  main.(*BrokenCounter).Get()
      broken_counter.go:15

Previous write at 0x... by goroutine 6:
  main.(*BrokenCounter).Increment()
      broken_counter.go:11
```

**Виправи:**

```bash
# Запусти fixed version
go run fixed_counter.go
```

### 2. Goroutine Leaks

```bash
cd ../02_goroutine_leaks

# Запусти broken code
go run broken_channel_leak.go
```

**Що побачиш:**

```
Goroutines before: 1
Goroutines after leakyReceiver: 2 (leaked: 1) ❌
Goroutines after leakySender: 3 (leaked: 2) ❌
⚠️ Memory leak!
```

**Виправи:**

```bash
# Запусти fixed version
go run fixed_channel_leak.go
```

---

## 📚 Читати теорію

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_8

# Race conditions
cat theory/01_race_conditions.md

# Goroutine leaks
cat theory/02_goroutine_leaks.md
```

---

## 🔧 Основні команди

### Race Detector

```bash
# Test з race detector
go test -race ./...

# Run з race detector
go run -race main.go

# Build з race detector (для testing)
go build -race
```

### Goroutine Count

```go
import "runtime"

count := runtime.NumGoroutine()
fmt.Println("Goroutines:", count)
```

### pprof

```bash
# Goroutine profile
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

---

## ✅ 3 Способи виправити Race Condition

### 1. Mutex

```go
var mu sync.Mutex
mu.Lock()
counter++
mu.Unlock()
```

### 2. Atomic

```go
atomic.AddInt64(&counter, 1)
```

### 3. Channel

```go
ch <- value  // Safe
```

---

## ✅ 3 Способи виправити Goroutine Leak

### 1. Context

```go
ctx, cancel := context.WithCancel(ctx)
defer cancel()

go func() {
    <-ctx.Done()
    return
}()
```

### 2. Done Channel

```go
done := make(chan struct{})

go func() {
    <-done
    return
}()

close(done)
```

### 3. Close Channel

```go
ch := make(chan int)

go func() {
    for val := range ch {  // Exits when closed
        process(val)
    }
}()

close(ch)  // Unblocks receiver
```

---

## 📖 Структура

```
week_8/
├── README.md           # Повний опис
├── QUICK_START.md      # Цей файл
├── theory/
│   ├── 01_race_conditions.md
│   └── 02_goroutine_leaks.md
├── practice/
│   ├── 01_race_conditions/
│   │   ├── broken_counter.go
│   │   └── fixed_counter.go
│   └── 02_goroutine_leaks/
│       ├── broken_channel_leak.go
│       └── fixed_channel_leak.go
└── exercises/          # TODO

```

---

## 🎯 Рекомендований порядок

### День 1: Race Conditions

1. Читай `theory/01_race_conditions.md`
2. Запускай `broken_counter.go` з `-race`
3. Аналізуй output
4. Запускай `fixed_counter.go`
5. Зрозумій 3 способи fix

### День 2: Goroutine Leaks

1. Читай `theory/02_goroutine_leaks.md`
2. Запускай `broken_channel_leak.go`
3. Подивись на `runtime.NumGoroutine()`
4. Запускай `fixed_channel_leak.go`
5. Зрозумій 6 способів fix

### День 3: Практика

1. Створи свій код з race condition
2. Знайди з `-race`
3. Виправ
4. Створи свій код з goroutine leak
5. Знайди з `runtime.NumGoroutine()`
6. Виправ

---

## 🐛 Common Patterns

### Pattern 1: Safe Counter

```go
type Counter struct {
    mu    sync.Mutex
    value int
}

func (c *Counter) Inc() {
    c.mu.Lock()
    c.value++
    c.mu.Unlock()
}
```

### Pattern 2: Worker with Context

```go
func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return  // Clean exit
        default:
            // Work
        }
    }
}
```

### Pattern 3: Producer/Consumer

```go
func producer(ch chan int) {
    defer close(ch)  // Close when done
    for i := 0; i < 10; i++ {
        ch <- i
    }
}

func consumer(ch chan int) {
    for val := range ch {  // Exits when closed
        process(val)
    }
}
```

---

## 🎓 Key Rules

### Race Conditions

1. **Always test with `-race`**
2. **Protect all shared state**
3. **Use atomic for simple ops**
4. **Channels for communication**

### Goroutine Leaks

1. **Every goroutine needs exit**
2. **Use context for cleanup**
3. **Close channels (producer)**
4. **Monitor goroutine count**

---

## 📖 Ресурси

- [Go Race Detector](https://go.dev/doc/articles/race_detector)
- [Uber's goleak](https://github.com/uber-go/goleak)
- Week 6: `theory/07_goroutines_concurrency.md`

---

**"No concurrent code is safe until `-race` tested!" 🔍**
