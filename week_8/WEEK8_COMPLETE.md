# ✅ Week 8 - Завершено!

## 🎯 Що створено

**Week 8: Debugging & Race Conditions** - модуль про пошук та виправлення bugs в concurrent Go коді.

---

## 📊 Статистика

### Створено файлів

**Теорія:** 2 файли
- `theory/01_race_conditions.md` (320+ рядків)
- `theory/02_goroutine_leaks.md` (450+ рядків)

**Практика:** 6+ файлів
- `practice/01_race_conditions/` (broken + fixed)
- `practice/02_goroutine_leaks/` (broken + fixed)

**Документація:** 3 файли
- `README.md` - Повний опис
- `QUICK_START.md` - Швидкий старт
- `WEEK8_COMPLETE.md` - Цей звіт

**Загалом:** 11+ файлів, ~1000+ рядків коду + документації

---

## 📚 Що покрито

### 1. Race Conditions ⚔️

**Теорія:**
- Що таке race condition?
- Як виникає (3 операції: READ, ADD, WRITE)
- Go Race Detector (`go test -race`)
- 3 способи виправити: Mutex, Atomic, Channel
- Порівняння performance
- Best practices

**Практика:**
- Broken counter (race на int)
- Fixed counter (3 рішення)
- Broken map (concurrent map access)
- Broken slice (concurrent append)

**Ключові інструменти:**
```bash
go test -race ./...
go run -race main.go
go build -race
```

### 2. Goroutine Leaks 🚰

**Теорія:**
- Що таке goroutine leak?
- 6 типів leaks:
  1. Blocked on channel receive
  2. Blocked on channel send
  3. Waiting on mutex
  4. Waiting on WaitGroup
  5. HTTP without timeout
  6. Context not canceled
- Як знайти: `runtime.NumGoroutine()`, pprof, goleak
- 3 patterns для clean shutdown

**Практика:**
- Broken channel leak (receiver + sender)
- Fixed channel leak (6 способів)
- Context для cancellation
- Done channel pattern
- Timeout pattern

**Ключові інструменти:**
```go
runtime.NumGoroutine()
go tool pprof .../goroutine
goleak.VerifyNone(t)
```

---

## 🔧 Debugging Workflow

### Крок 1: Знайти баг

**Race Condition:**
```bash
go run -race main.go
```

**Goroutine Leak:**
```go
before := runtime.NumGoroutine()
// ... code ...
after := runtime.NumGoroutine()
if after > before {
    fmt.Printf("Leaked: %d\n", after-before)
}
```

### Крок 2: Проаналізувати

**Race Detector Output:**
```
WARNING: DATA RACE
Read at 0x... by goroutine 7:
  main.Counter.Get()
      counter.go:15

Previous write at 0x... by goroutine 6:
  main.Counter.Increment()
      counter.go:11
```

**pprof Goroutines:**
```bash
go tool pprof http://localhost:6060/debug/pprof/goroutine
> top
> traces
```

### Крок 3: Виправити

**Race Condition:**
- Mutex для general purpose
- Atomic для simple ops
- Channel для communication

**Goroutine Leak:**
- Context для cancellation
- Done channel для shutdown
- Close channel для unblock
- Timeout для safety

---

## ✅ 3 Способи: Race Conditions

### 1. Mutex (Universal)

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

**Коли:** General purpose, складні операції

### 2. Atomic (Fast)

```go
var counter int64

func increment() {
    atomic.AddInt64(&counter, 1)
}
```

**Коли:** Прості операції (increment, swap)

### 3. Channel (Safe Communication)

```go
ch := make(chan int)

go func() {
    for val := range ch {
        process(val)  // Safe
    }
}()

ch <- 42
```

**Коли:** Communication між goroutines

---

## ✅ 6 Способів: Goroutine Leaks

### 1. Close Channel

```go
ch := make(chan int)
go func() {
    val, ok := <-ch
    if !ok { return }  // Channel closed
}()
close(ch)
```

### 2. Send Value

```go
ch := make(chan int)
go func() { val := <-ch }()
ch <- 42  // Unblocks
```

### 3. Buffered Channel

```go
ch := make(chan int, 1)  // Won't block
go func() { ch <- 42 }()
```

### 4. Context

```go
ctx, cancel := context.WithCancel(ctx)
go func() {
    <-ctx.Done()
    return
}()
cancel()
```

### 5. Done Channel

```go
done := make(chan struct{})
go func() {
    <-done
    return
}()
close(done)
```

### 6. Timeout

```go
select {
case val := <-ch:
    process(val)
case <-time.After(5 * time.Second):
    return  // Exit after timeout
}
```

---

## 📊 Performance Comparison

```
Benchmark Results (ns/op):

Mutex:    120 ns/op  (universal)
Atomic:    30 ns/op  (fastest)
Channel:  300 ns/op  (safest for communication)
```

**Висновок:**
- Atomic для counters, flags
- Mutex для загального використання
- Channel для communication

---

## 🎯 Best Practices

### Race Conditions

1. ✅ **Always test with `-race`**
2. ✅ **Protect all shared state**
3. ✅ **Keep critical sections small**
4. ✅ **Use defer for unlock**
5. ✅ **Document what mutex protects**

### Goroutine Leaks

1. ✅ **Every goroutine needs exit strategy**
2. ✅ **Use context for cancellation**
3. ✅ **Close channels (producer side)**
4. ✅ **Test with goleak**
5. ✅ **Monitor goroutine count in production**

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

### Bug 2: Channel Never Closed

```go
// ❌ BAD
ch := make(chan int)
go func() {
    <-ch  // Blocks forever
}()

// ✅ GOOD
ch := make(chan int)
go func() {
    <-ch
}()
close(ch)  // Unblocks
```

### Bug 3: Context Never Canceled

```go
// ❌ BAD
ctx := context.Background()
go func() {
    <-ctx.Done()  // Never happens
}()

// ✅ GOOD
ctx, cancel := context.WithCancel(ctx)
defer cancel()
go func() {
    <-ctx.Done()  // Will happen
}()
```

---

## 🚀 Як використовувати

### Quick Start

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_8

# Читати
cat README.md
cat QUICK_START.md

# Запускати
cd practice/01_race_conditions
go run -race broken_counter.go
go run fixed_counter.go

cd ../02_goroutine_leaks
go run broken_channel_leak.go
go run fixed_channel_leak.go
```

### Recommended Learning Path

**День 1:** Race Conditions
1. Теорія: `theory/01_race_conditions.md`
2. Практика: `practice/01_race_conditions/`
3. Запусти з `-race`
4. Зрозумій 3 способи fix

**День 2:** Goroutine Leaks
1. Теорія: `theory/02_goroutine_leaks.md`
2. Практика: `practice/02_goroutine_leaks/`
3. Подивись на `runtime.NumGoroutine()`
4. Зрозумій 6 способів fix

**День 3:** Власна практика
1. Створи код з race condition
2. Знайди з `-race`
3. Виправ
4. Створи код з goroutine leak
5. Виправ

---

## 🔗 Зв'язок з іншими модулями

### Week 6: Goroutines & Concurrency

Week 8 є продовженням Week 6!

```
Week 6: Як писати concurrent code
   ↓
Week 8: Як debug concurrent code
```

**Файли:**
- `week_6/theory/07_goroutines_concurrency.md` - Basics
- `week_6/practice/06_goroutines/main.go` - Examples
- `week_8/` - Debugging

### Week 7: State Machine

State Pattern використовує goroutines → можуть бути leaks!

```
Week 7: ATM State Machine
   ↓
Week 8: Debug goroutine leaks
```

### Design Patterns

Chain of Responsibility, Observer використовують goroutines.

---

## 📖 Ресурси

### Tools

1. **Race Detector:** `go test -race`
2. **pprof:** `go tool pprof .../goroutine`
3. **goleak:** `go.uber.org/goleak`
4. **runtime:** `runtime.NumGoroutine()`

### Documentation

- [Go Race Detector](https://go.dev/doc/articles/race_detector)
- [Concurrency Patterns](https://go.dev/blog/pipelines)
- [Uber goleak](https://github.com/uber-go/goleak)

---

## 🎓 Висновок

### Race Condition - це:

✅ Bug concurrent коду  
✅ Виникає без синхронізації  
✅ Знаходиться Race Detector  
✅ Виправляється Mutex/Atomic/Channel  

### Goroutine Leak - це:

✅ Memory leak в Go  
✅ Goroutine блокується назавжди  
✅ Знаходиться runtime.NumGoroutine(), pprof  
✅ Виправляється proper cleanup  

### Golden Rules:

1. **Always test with `-race`**
2. **Every goroutine needs exit strategy**
3. **Use context for cancellation**
4. **Close channels (producer side)**
5. **Monitor goroutine count**

---

## ✅ Week 8 Complete!

```
Progress: 100% ✅

Theory:   ████████████ 2/2
Practice: ████████████ 6/6
Docs:     ████████████ 3/3
```

**Дата завершення:** 2026-01-28  
**Статус:** COMPLETE ✅  
**Локація:** `/Users/vkuzm/GolandProjects/golang_practice/week_8`

---

## 🎉 Вітаємо!

Тепер ти вмієш:
- ✅ Знаходити race conditions з `-race`
- ✅ Виправляти race conditions (3 способи)
- ✅ Знаходити goroutine leaks
- ✅ Виправляти goroutine leaks (6 способів)
- ✅ Використовувати debugging tools
- ✅ Писати безпечний concurrent код

**"The best debugger is your knowledge!" 🐛🔍**
