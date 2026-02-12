# Task 5: Race Condition Detection & Fix

**Level:** Advanced  
**Time:** 15 minutes  
**Topics:** Race Conditions, Mutex, Thread Safety

---

## 📝 Task

Дано **багатопотоковий лічильник з race condition**. Твоє завдання:

1. Знайти race condition
2. Виправити його
3. Написати тест, який детектує race condition

---

## 🐛 Buggy Code

```go
type Counter struct {
    value int
}

func (c *Counter) Increment() {
    c.value++
}

func (c *Counter) GetValue() int {
    return c.value
}

// 1000 goroutines increment counter
func TestCounter() {
    counter := &Counter{}
    
    for i := 0; i < 1000; i++ {
        go counter.Increment()
    }
    
    time.Sleep(1 * time.Second)
    
    fmt.Println("Counter:", counter.GetValue())
    // Expected: 1000
    // Actual: varies (often < 1000) ❌
}
```

---

## ❓ Questions

1. **Чому результат НЕ 1000?**
2. **Де саме race condition?**
3. **Як виправити?**

---

## ✅ Requirements

Створи `SafeCounter` з такими методами:

```go
type SafeCounter interface {
    Increment()
    Decrement()
    GetValue() int
    Reset()
}
```

**Вимоги:**
- Thread-safe (без race conditions)
- Підтримка concurrent read/write
- Використай `sync.Mutex` або `sync.RWMutex`
- Напиши тест з `go test -race` для детекції race conditions

---

## 🧪 Test Cases

```go
// Test 1: Concurrent increments
func TestConcurrentIncrements(t *testing.T) {
    counter := NewSafeCounter()
    
    var wg sync.WaitGroup
    wg.Add(1000)
    
    for i := 0; i < 1000; i++ {
        go func() {
            defer wg.Done()
            counter.Increment()
        }()
    }
    
    wg.Wait()
    
    assert.Equal(t, 1000, counter.GetValue())
}

// Test 2: Concurrent increments and decrements
func TestConcurrentIncrementsAndDecrements(t *testing.T) {
    counter := NewSafeCounter()
    
    var wg sync.WaitGroup
    wg.Add(2000)
    
    // 1000 increments
    for i := 0; i < 1000; i++ {
        go func() {
            defer wg.Done()
            counter.Increment()
        }()
    }
    
    // 1000 decrements
    for i := 0; i < 1000; i++ {
        go func() {
            defer wg.Done()
            counter.Decrement()
        }()
    }
    
    wg.Wait()
    
    assert.Equal(t, 0, counter.GetValue())
}

// Test 3: Concurrent reads and writes
func TestConcurrentReadsAndWrites(t *testing.T) {
    counter := NewSafeCounter()
    
    var wg sync.WaitGroup
    wg.Add(2000)
    
    // 1000 writes
    for i := 0; i < 1000; i++ {
        go func() {
            defer wg.Done()
            counter.Increment()
        }()
    }
    
    // 1000 reads
    for i := 0; i < 1000; i++ {
        go func() {
            defer wg.Done()
            _ = counter.GetValue()
        }()
    }
    
    wg.Wait()
    
    assert.Equal(t, 1000, counter.GetValue())
}

// Test 4: Reset
func TestReset(t *testing.T) {
    counter := NewSafeCounter()
    counter.Increment()
    counter.Increment()
    counter.Increment()
    
    counter.Reset()
    
    assert.Equal(t, 0, counter.GetValue())
}
```

---

## 💡 Hints

### Race Condition Explanation

```go
// ❌ NOT thread-safe
func (c *Counter) Increment() {
    c.value++  // This is actually 3 operations:
               // 1. Read c.value
               // 2. Add 1
               // 3. Write c.value
               
    // Two goroutines can interleave:
    // Goroutine A: Read c.value (0)
    // Goroutine B: Read c.value (0)
    // Goroutine A: Write c.value (1)
    // Goroutine B: Write c.value (1)
    // Result: 1 instead of 2 ❌
}
```

### Solution 1: Mutex

```go
type SafeCounter struct {
    value int
    mu    sync.Mutex
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}
```

### Solution 2: RWMutex (Better for many reads)

```go
type SafeCounter struct {
    value int
    mu    sync.RWMutex
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

func (c *SafeCounter) GetValue() int {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.value
}
```

---

## 🔍 How to Detect Race Conditions

```bash
# Run tests with race detector
go test -race ./...

# Run program with race detector
go run -race main.go
```

**Race detector output:**
```
WARNING: DATA RACE
Read at 0x00c000012088 by goroutine 7:
  main.(*Counter).GetValue()
      /path/to/counter.go:10 +0x3a

Previous write at 0x00c000012088 by goroutine 6:
  main.(*Counter).Increment()
      /path/to/counter.go:6 +0x4e
```

---

## 🎯 Real-World Examples

Race conditions часто зустрічаються в:

1. **HTTP request counters**
```go
type Server struct {
    requestCount int  // ❌ Race condition!
}
```

2. **Caching layers**
```go
type Cache struct {
    data map[string]interface{}  // ❌ Race condition!
}
```

3. **Connection pools**
```go
type Pool struct {
    activeConnections int  // ❌ Race condition!
}
```

---

**Рішення:** `solutions/solution_05_race_condition.go`

**Good luck!** 🚀
