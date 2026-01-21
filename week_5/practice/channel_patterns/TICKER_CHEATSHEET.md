# ⏰ Time Ticker - Швидка Шпаргалка

## 📝 Основи

```go
// Створення (тікає кожні 500ms)
ticker := time.NewTicker(500 * time.Millisecond)
defer ticker.Stop() // ⚠️ ЗАВЖДИ!

// ticker.C - це канал типу <-chan time.Time
```

---

## 🔄 Читання

```go
// Спосіб 1: Range
for t := range ticker.C {
    fmt.Println("Tick:", t)
}

// Спосіб 2: Select
for {
    select {
    case t := <-ticker.C:
        fmt.Println("Tick:", t)
    }
}
```

---

## ⚡ Швидкі Команди

| Команда | Опис |
|---------|------|
| `time.NewTicker(d)` | Створити ticker з інтервалом `d` |
| `ticker.C` | Канал типу `<-chan time.Time` |
| `ticker.Stop()` | Зупинити ticker (звільнити ресурси) |
| `<-ticker.C` | Прочитати час з каналу |

---

## 🎯 Use Cases

| Use Case | Інтервал | Приклад |
|----------|----------|---------|
| **UI Updates** | 16ms | `time.NewTicker(16 * time.Millisecond)` |
| **Polling** | 100ms-1s | `time.NewTicker(500 * time.Millisecond)` |
| **Heartbeat** | 1s-10s | `time.NewTicker(5 * time.Second)` |
| **Cleanup** | 1h-24h | `time.NewTicker(1 * time.Hour)` |

---

## 🔥 Типові Patterns

### Pattern 1: Ticker + Done

```go
ticker := time.NewTicker(1 * time.Second)
defer ticker.Stop()

for {
    select {
    case <-ticker.C:
        doWork()
    case <-done:
        return
    }
}
```

### Pattern 2: Rate Limiting

```go
limiter := time.NewTicker(100 * time.Millisecond)
defer limiter.Stop()

for _, req := range requests {
    <-limiter.C // Wait 100ms
    process(req)
}
```

### Pattern 3: Multiple Tickers

```go
fast := time.NewTicker(1 * time.Second)
slow := time.NewTicker(10 * time.Second)
defer fast.Stop()
defer slow.Stop()

for {
    select {
    case <-fast.C:
        quickTask()
    case <-slow.C:
        heavyTask()
    }
}
```

---

## 📊 Ticker vs Timer

| | Ticker | Timer |
|---|--------|-------|
| **Тип** | Періодичний | Одноразовий |
| **Створення** | `time.NewTicker(d)` | `time.NewTimer(d)` |
| **Поведінка** | Тікає КОЖНІ `d` | Тікає ОДИН РАЗ через `d` |
| **Use Case** | Повторювані задачі | Таймаути |

```go
// Ticker - повторюється
ticker := time.NewTicker(1 * time.Second)
<-ticker.C // Спрацює
<-ticker.C // Спрацює знову через 1s
<-ticker.C // І знову через 1s

// Timer - одноразовий
timer := time.NewTimer(1 * time.Second)
<-timer.C // Спрацює один раз
//<-timer.C // ❌ deadlock!
```

---

## ✅ DO / ❌ DON'T

### ✅ DO

```go
// ✅ Завжди Stop()
ticker := time.NewTicker(1 * time.Second)
defer ticker.Stop()

// ✅ Select з done
select {
case <-ticker.C:
    work()
case <-done:
    return
}

// ✅ Горутина для довгих операцій
for t := range ticker.C {
    go process(t)
}
```

### ❌ DON'T

```go
// ❌ Забули Stop() - memory leak!
ticker := time.NewTicker(1 * time.Second)
// ... без ticker.Stop()

// ❌ time.Tick - не можна зупинити!
for range time.Tick(1 * time.Second) {
    // memory leak!
}

// ❌ Блокування ticker
for t := range ticker.C {
    time.Sleep(10 * time.Second) // Пропустимо тіки!
}
```

---

## 🚨 Common Mistakes

### Mistake 1: Забув Stop()

```go
// ❌ BAD
func bad() {
    ticker := time.NewTicker(1 * time.Second)
    // Забули ticker.Stop()!
}

// ✅ GOOD
func good() {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
}
```

### Mistake 2: Повільна обробка

```go
// ❌ BAD - пропустиш тіки
ticker := time.NewTicker(100 * time.Millisecond)
for t := range ticker.C {
    time.Sleep(1 * time.Second) // Блокує!
}

// ✅ GOOD - не блокує
ticker := time.NewTicker(100 * time.Millisecond)
for t := range ticker.C {
    go process(t) // В горутині
}
```

### Mistake 3: Горутина без cleanup

```go
// ❌ BAD - як зупинити?
go func() {
    ticker := time.NewTicker(1 * time.Second)
    for range ticker.C {
        work()
    }
}()

// ✅ GOOD - з done каналом
go func() {
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
}()
```

---

## 📐 Ticker.C Буфер

**ticker.C має буфер = 1:**

```go
ticker := time.NewTicker(1 * time.Second)
time.Sleep(5 * time.Second)
<-ticker.C // Отримаєш ТІЛЬКИ 1 тік (останній)
// Решта 4 тіки - пропали!
```

**Висновок:** Читай регулярно!

---

## 🎓 Complete Example

```go
package main

import (
    "fmt"
    "time"
)

func worker(done <-chan bool) {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    count := 0
    for {
        select {
        case t := <-ticker.C:
            count++
            fmt.Printf("Tick %d at %s\n", count, t.Format("15:04:05"))
        case <-done:
            fmt.Println("Stopping worker...")
            return
        }
    }
}

func main() {
    done := make(chan bool)

    go worker(done)

    // Працюємо 5 секунд
    time.Sleep(5 * time.Second)

    // Graceful stop
    done <- true

    time.Sleep(100 * time.Millisecond)
    fmt.Println("Done!")
}
```

**Вивід:**
```
Tick 1 at 14:30:01
Tick 2 at 14:30:02
Tick 3 at 14:30:03
Tick 4 at 14:30:04
Tick 5 at 14:30:05
Stopping worker...
Done!
```

---

## 📚 Related

```go
// Блокує горутину
time.Sleep(1 * time.Second)

// Канал (одноразовий)
<-time.After(1 * time.Second)

// Timer (одноразовий, можна зупинити)
timer := time.NewTimer(1 * time.Second)
<-timer.C
timer.Stop()

// Ticker (періодичний, ЗАВЖДИ зупиняй)
ticker := time.NewTicker(1 * time.Second)
<-ticker.C // Спрацює
<-ticker.C // Спрацює знову
ticker.Stop()
```

---

## 🚀 Run Examples

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_5/practice/channel_patterns

# Run all ticker examples
go run ticker_examples.go
```

**7 Examples:**
1. ✅ Basic Ticker
2. ✅ Ticker with Stop
3. ✅ Ticker in Select
4. ✅ Multiple Tickers
5. ✅ Ticker vs Timer
6. ✅ Rate Limiting
7. ✅ Graceful Shutdown

---

**⭐ Key Takeaway:**
- Ticker = Періодичний (repeating)
- Timer = Одноразовий (one-shot)
- ЗАВЖДИ викликай `ticker.Stop()`!

---

**Файли:**
- `ticker_examples.go` - Код з прикладами
- `TICKER_GUIDE.md` - Повна документація
- `TICKER_CHEATSHEET.md` - Ця шпаргалка
