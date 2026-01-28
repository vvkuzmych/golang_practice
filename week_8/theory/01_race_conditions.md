# Race Conditions в Go

## 🎯 Що таке Race Condition?

**Race Condition** виникає коли дві або більше goroutines намагаються одночасно читати і писати в одну і ту ж змінну без синхронізації.

## ⚠️ Проблема

```go
var counter = 0

func increment() {
    counter++  // ❌ НЕ атомарна операція!
}

func main() {
    for i := 0; i < 1000; i++ {
        go increment()
    }
    time.Sleep(time.Second)
    fmt.Println(counter)  // Очікується 1000, але буде менше!
}
```

### Чому це проблема?

```
counter++ це насправді 3 операції:

1. READ:  temp = counter
2. ADD:   temp = temp + 1  
3. WRITE: counter = temp

Goroutine 1:      Goroutine 2:
READ (0)          
                  READ (0)
ADD (1)           
                  ADD (1)
WRITE (1)         
                  WRITE (1)  ← Очікувалось 2, але 1!
```

---

## 🔍 Як знайти Race Condition?

### 1. Go Race Detector

```bash
# Build з race detector
go build -race

# Test з race detector
go test -race

# Run з race detector
go run -race main.go
```

### 2. Race Detector Output

```
==================
WARNING: DATA RACE
Read at 0x00c000014098 by goroutine 7:
  main.increment()
      /path/to/main.go:10 +0x38

Previous write at 0x00c000014098 by goroutine 6:
  main.increment()
      /path/to/main.go:10 +0x4e

Goroutine 7 (running) created at:
  main.main()
      /path/to/main.go:15 +0x7e
==================
```

---

## ✅ Рішення Race Conditions

### 1. Mutex (Mutual Exclusion)

```go
var (
    counter int
    mu      sync.Mutex
)

func increment() {
    mu.Lock()
    counter++
    mu.Unlock()
}
```

### 2. Atomic Operations

```go
var counter int64

func increment() {
    atomic.AddInt64(&counter, 1)
}
```

### 3. Channels

```go
func worker(ch chan int) {
    for {
        val := <-ch
        // Process val safely
    }
}

func main() {
    ch := make(chan int)
    go worker(ch)
    
    for i := 0; i < 1000; i++ {
        ch <- i  // Safe communication
    }
}
```

---

## 📊 Порівняння рішень

| Метод | Performance | Use Case |
|-------|-------------|----------|
| **Mutex** | Medium | General purpose |
| **RWMutex** | Better для reads | Багато reads, мало writes |
| **Atomic** | Fast | Прості операції (increment, swap) |
| **Channels** | Slower | Communication between goroutines |

---

## 🐛 Типові Race Conditions

### 1. Shared Counter

```go
// ❌ BAD
var counter int

func increment() {
    counter++  // Race!
}

// ✅ GOOD
var counter int64

func increment() {
    atomic.AddInt64(&counter, 1)
}
```

### 2. Map Access

```go
// ❌ BAD
var m = make(map[string]int)

func update(key string) {
    m[key]++  // Race!
}

// ✅ GOOD
var (
    m  = make(map[string]int)
    mu sync.Mutex
)

func update(key string) {
    mu.Lock()
    m[key]++
    mu.Unlock()
}

// ✅ BETTER (Go 1.9+)
var m sync.Map

func update(key string, val int) {
    m.Store(key, val)
}
```

### 3. Slice Append

```go
// ❌ BAD
var slice []int

func append(val int) {
    slice = append(slice, val)  // Race!
}

// ✅ GOOD
var (
    slice []int
    mu    sync.Mutex
)

func append(val int) {
    mu.Lock()
    slice = append(slice, val)
    mu.Unlock()
}
```

### 4. Struct Fields

```go
type Counter struct {
    value int  // ❌ Not safe for concurrent access
}

func (c *Counter) Increment() {
    c.value++  // Race!
}

// ✅ GOOD
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

---

## 🔒 sync.Mutex vs sync.RWMutex

### Mutex (Exclusive Lock)

```go
var mu sync.Mutex

// Тільки 1 goroutine може тримати lock
mu.Lock()
// Critical section
mu.Unlock()
```

### RWMutex (Read/Write Lock)

```go
var mu sync.RWMutex

// Багато readers одночасно
mu.RLock()
data := sharedData  // Read
mu.RUnlock()

// Тільки 1 writer
mu.Lock()
sharedData = newData  // Write
mu.Unlock()
```

**Коли використовувати RWMutex:**
- Багато reads, мало writes
- Reads дорогі (наприклад, copying великої структури)

---

## ⚡ Atomic Operations

### Підтримувані операції

```go
import "sync/atomic"

var counter int64

// Add
atomic.AddInt64(&counter, 1)
atomic.AddInt64(&counter, -1)

// Load
val := atomic.LoadInt64(&counter)

// Store
atomic.StoreInt64(&counter, 100)

// Swap
old := atomic.SwapInt64(&counter, 200)

// Compare and Swap (CAS)
swapped := atomic.CompareAndSwapInt64(&counter, 100, 200)
```

### Підтримувані типи

- `int32`, `int64`
- `uint32`, `uint64`
- `uintptr`
- `unsafe.Pointer`

---

## 🎯 Best Practices

### 1. Завжди використовуй Race Detector

```bash
# В CI/CD
go test -race ./...

# Локально
go test -race -short ./...
```

### 2. Принцип: "Do not communicate by sharing memory; share memory by communicating"

```go
// ❌ BAD: Sharing memory
var sharedData int
mu.Lock()
sharedData = 42
mu.Unlock()

// ✅ GOOD: Communicating
ch <- 42
data := <-ch
```

### 3. Keep Critical Sections Small

```go
// ❌ BAD: Lock тримається занадто довго
mu.Lock()
data := fetchData()       // Slow operation
processData(data)         // Slow operation
saveData(data)           // Slow operation
mu.Unlock()

// ✅ GOOD: Мінімальна critical section
data := fetchData()
processedData := processData(data)

mu.Lock()
saveData(processedData)  // Тільки це під lock
mu.Unlock()
```

### 4. Defer Unlock

```go
func update() error {
    mu.Lock()
    defer mu.Unlock()  // ✅ Завжди unlock, навіть якщо panic
    
    if err := validate(); err != nil {
        return err  // Unlock автоматично
    }
    
    // Update data
    return nil
}
```

---

## 🔍 Debugging Race Conditions

### 1. Включи Race Detector

```bash
go test -race ./...
```

### 2. Проаналізуй Output

```
WARNING: DATA RACE
Read at 0x00c000014098 by goroutine 7:
  main.increment()
      /path/to/main.go:10 +0x38
                          ↑ Line number

Previous write at 0x00c000014098 by goroutine 6:
  main.increment()
      /path/to/main.go:10 +0x4e
```

### 3. Знайди Access Pattern

- Які goroutines?
- Який shared state?
- Яка операція (read/write)?

### 4. Додай Synchronization

- Mutex для загального випадку
- Atomic для простих операцій
- Channels для communication

---

## 📈 Performance Impact

### Race Detector Overhead

```
CPU: ~10x slower
Memory: ~10x more

⚠️ НЕ використовуй в production!
Тільки для testing і development.
```

### Synchronization Overhead

```
No sync:     1x (baseline, але з races!)
Atomic:      ~2-3x
Mutex:       ~5-10x
Channel:     ~10-20x

Choose wisely!
```

---

## 🎓 Висновок

### Race Condition - це:

✅ Баг concurrency  
✅ Виникає без синхронізації  
✅ Знаходиться Race Detector  
✅ Виправляється Mutex/Atomic/Channels  

### Golden Rules:

1. **Always use Race Detector in tests**
2. **Minimize critical sections**
3. **Prefer channels for communication**
4. **Use atomic for simple counters**
5. **Document which mutex protects what**

---

## 📖 Далі

- `practice/01_race_conditions/` - Практичні приклади з багами
- `02_goroutine_leaks.md` - Goroutine leaks
- `go test -race` - Як використовувати

**Race Detector - ваш найкращий друг в concurrent Go!** 🔍
