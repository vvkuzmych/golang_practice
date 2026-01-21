# ⏰ Time Ticker з Каналами - Повний Гід

**Все про time.Ticker та канали в Go**

---

## 🎯 Що таке Ticker?

**Ticker** - це механізм, який відправляє поточний час у канал **періодично** (кожні N мілісекунд).

```go
ticker := time.NewTicker(500 * time.Millisecond)
// ticker.C - це канал типу <-chan time.Time
```

---

## 📚 Основи

### Створення Ticker

```go
// Створення ticker, який "тікає" кожні 500ms
ticker := time.NewTicker(500 * time.Millisecond)
defer ticker.Stop() // ⚠️ ЗАВЖДИ зупиняй!

// ticker.C - це RECEIVE-ONLY КАНАЛ (<-chan time.Time)
```

### Читання з Ticker

```go
// Спосіб 1: Range (простіший)
for t := range ticker.C {
    fmt.Println("Tick at:", t)
}

// Спосіб 2: Select (гнучкіший)
for {
    select {
    case t := <-ticker.C:
        fmt.Println("Tick at:", t)
    }
}
```

### ⚠️ ВАЖЛИВО: Завжди Stop()!

```go
ticker := time.NewTicker(1 * time.Second)
defer ticker.Stop() // Звільняє ресурси!

// Без Stop() - memory leak! Ticker продовжить працювати!
```

---

## 🔍 Ticker.C - Що це?

```go
type Ticker struct {
    C <-chan time.Time // RECEIVE-ONLY канал
    // ... інші поля
}
```

**Характеристики:**
- 📨 Тип: `<-chan time.Time` (receive-only)
- ⏰ Відправляє `time.Time` кожні N мілісекунд
- 🔄 Періодичний (безкінечний, поки не зупиниш)
- 📦 Буфер: 1 елемент (не блокує якщо не читаєш одразу)

---

## 📊 Ticker vs Timer

| Аспект | Ticker | Timer |
|--------|--------|-------|
| **Тип** | Періодичний (повторюється) | Одноразовий |
| **Створення** | `time.NewTicker(d)` | `time.NewTimer(d)` |
| **Канал** | `ticker.C` | `timer.C` |
| **Поведінка** | Тікає кожні `d` мілісекунд | Тікає ОДИН РАЗ через `d` мілісекунд |
| **Зупинка** | `ticker.Stop()` | `timer.Stop()` |
| **Use case** | Періодичні задачі | Таймаути, затримки |

### Приклад:

```go
// Ticker - повторюється
ticker := time.NewTicker(1 * time.Second)
for i := 0; i < 3; i++ {
    <-ticker.C // Спрацює 3 рази: 1s, 2s, 3s
}

// Timer - одноразовий
timer := time.NewTimer(1 * time.Second)
<-timer.C // Спрацює ОДИН РАЗ через 1s
//<-timer.C // ❌ deadlock! (не спрацює знову)
```

---

## 🎯 Use Cases

### 1️⃣ Періодичне виконання задачі

```go
ticker := time.NewTicker(5 * time.Second)
defer ticker.Stop()

for {
    select {
    case <-ticker.C:
        // Виконується кожні 5 секунд
        checkHealth()
        syncData()
    }
}
```

---

### 2️⃣ Rate Limiting (Обмеження швидкості)

```go
limiter := time.NewTicker(100 * time.Millisecond)
defer limiter.Stop()

for _, request := range requests {
    <-limiter.C // Чекаємо 100ms між запитами
    processRequest(request)
}
```

**Результат:** Максимум 10 запитів в секунду

---

### 3️⃣ Heartbeat Pattern

```go
func worker(done <-chan bool) {
    ticker := time.NewTicker(2 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            fmt.Println("❤️ Worker is alive")
        case <-done:
            return
        }
    }
}
```

---

### 4️⃣ Timeout з прогресом

```go
ticker := time.NewTicker(1 * time.Second)
timeout := time.After(5 * time.Second)
defer ticker.Stop()

for {
    select {
    case <-ticker.C:
        fmt.Println("Still waiting...")
    case <-timeout:
        fmt.Println("Timeout!")
        return
    case result := <-workCh:
        fmt.Println("Got result:", result)
        return
    }
}
```

---

### 5️⃣ Periodic Cleanup

```go
func cleanup(done <-chan bool) {
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            deleteOldFiles()
            clearCache()
            compactDatabase()
        case <-done:
            return
        }
    }
}
```

---

## 🔥 Patterns з Select

### Pattern 1: Ticker + Work + Done

```go
ticker := time.NewTicker(500 * time.Millisecond)
defer ticker.Stop()

for {
    select {
    case work := <-workCh:
        process(work)
    case <-ticker.C:
        fmt.Println("Heartbeat")
    case <-done:
        return
    }
}
```

---

### Pattern 2: Multiple Tickers

```go
ticker1 := time.NewTicker(1 * time.Second)   // Швидкий
ticker2 := time.NewTicker(5 * time.Second)   // Повільний
defer ticker1.Stop()
defer ticker2.Stop()

for {
    select {
    case <-ticker1.C:
        quickTask()
    case <-ticker2.C:
        heavyTask()
    }
}
```

---

### Pattern 3: Ticker з Graceful Shutdown

```go
func worker(done <-chan bool) {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            // Робота
            doWork()
        case <-done:
            // Graceful shutdown
            fmt.Println("Stopping...")
            cleanup()
            return
        }
    }
}

// Використання:
done := make(chan bool)
go worker(done)

// Пізніше...
done <- true // Graceful stop
```

---

## ⚡ Advanced: time.Tick (НЕ рекомендується!)

```go
// ❌ BAD: time.Tick - не можна зупинити!
for t := range time.Tick(1 * time.Second) {
    fmt.Println(t)
    // Ticker продовжує працювати навіть після виходу!
}

// ✅ GOOD: time.NewTicker - можна зупинити
ticker := time.NewTicker(1 * time.Second)
defer ticker.Stop()
for t := range ticker.C {
    fmt.Println(t)
}
```

**Чому `time.Tick` погано:**
- ❌ Не можна зупинити
- ❌ Memory leak
- ❌ Ticker працює вічно

**Використовуй `time.NewTicker` завжди!**

---

## 🚨 Common Mistakes

### ❌ Mistake 1: Забув Stop()

```go
// ❌ Memory leak!
func bad() {
    ticker := time.NewTicker(1 * time.Second)
    // ... робота ...
    // Забули ticker.Stop()!
}

// ✅ Правильно
func good() {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop() // Завжди!
    // ... робота ...
}
```

---

### ❌ Mistake 2: Блокування Ticker

```go
// ❌ Повільна обробка блокує ticker
ticker := time.NewTicker(100 * time.Millisecond)
for t := range ticker.C {
    time.Sleep(1 * time.Second) // Блокує!
    // Пропустимо багато тіків!
}

// ✅ Обробляй в окремій горутині
ticker := time.NewTicker(100 * time.Millisecond)
for t := range ticker.C {
    go process(t) // Не блокує
}
```

---

### ❌ Mistake 3: Ticker в горутині без cleanup

```go
// ❌ Горутина + ticker без cleanup
go func() {
    ticker := time.NewTicker(1 * time.Second)
    for range ticker.C {
        // Як зупинити???
    }
}()

// ✅ З done каналом
go func() {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    for {
        select {
        case <-ticker.C:
            // робота
        case <-done:
            return
        }
    }
}()
```

---

## 📐 Буфер у Ticker.C

**Ticker.C має буфер розміром 1:**

```go
ticker := time.NewTicker(1 * time.Second)

// Якщо не читаєш - ticker НЕ блокує відправника
time.Sleep(5 * time.Second)
// Ticker відправив 5 тіків, але в буфері тільки 1!

<-ticker.C // Отримаєш ОДИН тік (останній)
// Решта 4 тіки пропали!
```

**Висновок:** Читай з ticker.C регулярно, інакше пропустиш тіки!

---

## 🎓 Best Practices

### ✅ DO:

1. **Завжди зупиняй ticker:**
```go
ticker := time.NewTicker(1 * time.Second)
defer ticker.Stop()
```

2. **Використовуй select для flexibility:**
```go
select {
case <-ticker.C:
    // робота
case <-done:
    return
}
```

3. **Обробляй довгі операції в горутинах:**
```go
for t := range ticker.C {
    go heavyWork(t) // Не блокує ticker
}
```

4. **Додавай done канал для graceful shutdown:**
```go
func worker(done <-chan bool) {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    for {
        select {
        case <-ticker.C:
            work()
        case <-done:
            return
        }
    }
}
```

---

### ❌ DON'T:

1. **Не використовуй time.Tick:**
```go
// ❌ BAD
for range time.Tick(1 * time.Second) {
    // не можна зупинити!
}
```

2. **Не забувай Stop():**
```go
// ❌ Memory leak
ticker := time.NewTicker(1 * time.Second)
// ... робота без ticker.Stop()
```

3. **Не блокуй ticker:**
```go
// ❌ BAD
for t := range ticker.C {
    time.Sleep(10 * time.Second) // Блокує!
}
```

---

## 📊 Performance Tips

### Вибір інтервалу

```go
// ✅ Для UI updates
ticker := time.NewTicker(16 * time.Millisecond) // ~60 FPS

// ✅ Для polling
ticker := time.NewTicker(100 * time.Millisecond) // 10 разів/сек

// ✅ Для heartbeat
ticker := time.NewTicker(1 * time.Second)

// ✅ Для cleanup
ticker := time.NewTicker(1 * time.Hour)
```

### Мінімальний інтервал

```go
// ⚠️ Дуже короткий інтервал = висока нагрузка
ticker := time.NewTicker(1 * time.Millisecond) // Обережно!

// Мінімальний розумний інтервал: 10-100ms
ticker := time.NewTicker(10 * time.Millisecond) // OK
```

---

## 🔬 Приклади у файлі

Запусти повні приклади:

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_5/practice/channel_patterns

# Замінити main() на mainTicker() у ticker_examples.go
# Або запустити окремі функції
go run ticker_examples.go
```

**7 прикладів включено:**
1. Basic Ticker
2. Ticker with Stop
3. Ticker in Select
4. Multiple Tickers
5. Ticker vs Timer
6. Rate Limiting
7. Graceful Shutdown

---

## 📚 Related

- `time.Timer` - одноразовий (one-shot)
- `time.After(d)` - канал, який відправляє час через d
- `time.Sleep(d)` - блокує горутину на d

```go
// Порівняння
time.Sleep(1 * time.Second)           // Блокує
<-time.After(1 * time.Second)         // Канал (одноразовий)
<-time.NewTimer(1 * time.Second).C    // Timer (одноразовий)
<-time.NewTicker(1 * time.Second).C   // Ticker (періодичний)
```

---

**Створено:** 2026-01-19  
**Week:** 5 - Goroutines & Channels  
**Файл з прикладами:** `ticker_examples.go`
