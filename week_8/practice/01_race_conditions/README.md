# Race Conditions - Практика

## 🎯 Мета

Навчитися знаходити та виправляти race conditions за допомогою Go Race Detector.

---

## 📁 Файли

### 1. Broken Code (з багами)

- `broken_counter.go` - Race condition на counter
- `broken_map.go` - Race condition на map
- `broken_slice.go` - Race condition на slice

### 2. Fixed Code (виправлено)

- `fixed_counter.go` - 3 способи fix: Mutex, Atomic, Channel
- `fixed_map.go` - sync.Map і Mutex
- `fixed_slice.go` - Mutex для slice

---

## 🚀 Як запускати

### Крок 1: Запусти broken code з Race Detector

```bash
cd practice/01_race_conditions

# Broken counter
go run -race broken_counter.go
```

### Крок 2: Проаналізуй output

```
==================
WARNING: DATA RACE
Read at 0x... by goroutine 7:
  main.(*BrokenCounter).Get()
      broken_counter.go:15 +0x38

Previous write at 0x... by goroutine 6:
  main.(*BrokenCounter).Increment()
      broken_counter.go:11 +0x4e
==================
```

### Крок 3: Запусти fixed code

```bash
go run fixed_counter.go
```

---

## 🔍 Що шукати в Race Detector Output?

### 1. Line Numbers

```
broken_counter.go:15  ← Яка лінія коду?
```

### 2. Operation Type

```
Read at...       ← Читання
Previous write...  ← Запис
```

### 3. Goroutine IDs

```
by goroutine 7  ← Яка goroutine?
```

---

## ✅ 3 Способи виправити Race Condition

### 1. Mutex (General Purpose)

```go
type Counter struct {
    mu    sync.Mutex
    value int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    c.value++
    c.mu.Unlock()
}
```

**Коли використовувати:**
- Загальний випадок
- Складні операції

### 2. Atomic (Simple Operations)

```go
type Counter struct {
    value int64
}

func (c *Counter) Increment() {
    atomic.AddInt64(&c.value, 1)
}
```

**Коли використовувати:**
- Прості операції (increment, swap)
- Максимальна performance

### 3. Channel (Communication)

```go
ch := make(chan int)

go func() {
    for val := range ch {
        process(val)
    }
}()

ch <- 42  // Safe
```

**Коли використовувати:**
- Communication між goroutines
- Async processing

---

## 📊 Performance Порівняння

```bash
# Benchmark всіх методів
go test -bench=. -benchmem
```

**Expected results:**

```
BenchmarkMutex-8     10000000    120 ns/op
BenchmarkAtomic-8    50000000     30 ns/op
BenchmarkChannel-8    5000000    300 ns/op
```

**Висновок:**
- Atomic найшвидший для простих операцій
- Mutex середній, але універсальний
- Channel повільніший, але найбезпечніший для communication

---

## 🐛 Завдання

### Завдання 1: Знайди і виправ race condition

```bash
go run -race broken_counter.go
```

1. Запусти з `-race`
2. Прочитай output
3. Знайди проблемну лінію
4. Подивись на fixed версію
5. Зрозумій чому fix працює

### Завдання 2: Benchmark

```bash
go test -bench=. -benchmem
```

Порівняй performance різних методів.

### Завдання 3: Власний приклад

Створи свій код з race condition і виправ його всіма 3 способами.

---

## 🎯 Key Takeaways

1. **Завжди використовуй `-race` під час тестування**
2. **Atomic для простих операцій (increment, swap)**
3. **Mutex для складних операцій**
4. **Channel для communication**
5. **Race Detector знаходить майже всі race conditions**

---

## 📖 Далі

- Запусти всі broken приклади з `-race`
- Проаналізуй output
- Запусти fixed версії
- Зрозумій різницю

**"No code is safe until tested with -race!" 🔍**
