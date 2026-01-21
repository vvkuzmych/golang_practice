# 📚 All Channel Types in Go

**Повний довідник з усіма типами каналів та операціями**

---

## 📦 Файли

```
channel_patterns/
├── all_channel_types.go      # Всі приклади (запускаємо)
├── CHANNEL_CHEATSHEET.md     # Швидка шпаргалка
└── README_CHANNELS.md         # Цей файл
```

---

## 🚀 Як запустити

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_5/practice/channel_patterns

# Запустити всі приклади
go run all_channel_types.go
```

**Вивід:**
```
╔══════════════════════════════════════════════════════════════════╗
║          📚 All Channel Types in Go                             ║
╚══════════════════════════════════════════════════════════════════╝

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
1️⃣  UNBUFFERED CHANNEL
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Goroutine: Sending 42...
  Main: Receiving...
  Main: Received 42
  Goroutine: Sent!

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
2️⃣  BUFFERED CHANNEL
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Channel buffer: len=3, cap=3
  Reading: first
  Reading: second
  Reading: third

... (і так далі для всіх 10 типів)
```

---

## 📚 Що включено

### 1️⃣ Unbuffered Channel (Небуферизований)
- Створення: `make(chan T)`
- Завжди блокує запис і читання
- Синхронна комунікація

### 2️⃣ Buffered Channel (Буферизований)
- Створення: `make(chan T, capacity)`
- Блокує лише коли буфер повний/порожній
- `len(ch)` та `cap(ch)`

### 3️⃣ Send-Only Channel (Тільки запис)
- Тип: `chan<- T`
- Можна тільки писати
- Використовується в параметрах функцій

### 4️⃣ Receive-Only Channel (Тільки читання)
- Тип: `<-chan T`
- Можна тільки читати
- Використовується в параметрах функцій

### 5️⃣ Closing Channel (Закриття)
- `close(ch)`
- Після close() не можна писати (panic)
- Можна читати (повертає zero value)
- Comma-ok idiom: `val, ok := <-ch`

### 6️⃣ Range Over Channel (Ітерація)
- `for val := range ch`
- Зупиняється після `close(ch)`
- Автоматичне читання всіх значень

### 7️⃣ Select (Вибір з кількох каналів)
- Вибирає перший готовий case
- Випадковий вибір, якщо кілька готові
- `default` для non-blocking операцій
- Таймаути з `time.After()`

### 8️⃣ Nil Channel
- `var ch chan T` (не ініціалізований)
- Операції блокують навічно
- Корисно в select для "вимкнення" каналу

### 9️⃣ Channel of Channels (Канал каналів)
- `chan chan T`
- Request-response pattern
- Динамічна комунікація

### 🔟 Bidirectional vs Unidirectional
- `chan T` - можна read & write
- `chan<- T` - тільки write
- `<-chan T` - тільки read
- Автоматична конвертація

---

## 🎯 Quick Syntax Reference

```go
// Створення
ch := make(chan int)        // unbuffered
ch := make(chan int, 10)    // buffered

// Запис
ch <- 42

// Читання
value := <-ch

// Перевірка закриття
value, ok := <-ch
if !ok {
    // канал закритий
}

// Закриття
close(ch)

// Range
for val := range ch {
    // ...
}

// Select
select {
case val := <-ch1:
    // ...
case ch2 <- 42:
    // ...
default:
    // ...
}

// Send-only параметр
func send(ch chan<- int) {
    ch <- 42
}

// Receive-only параметр
func receive(ch <-chan int) int {
    return <-ch
}
```

---

## 📊 Таблиця операцій

| Операція | Unbuffered | Buffered (є місце) | Buffered (повний) | Closed | Nil |
|----------|------------|-------------------|-------------------|--------|-----|
| **ch <- v** | блокує | OK | блокує | panic | блокує навічно |
| **v := <-ch** | блокує | OK | блокує (якщо порожній) | zero value | блокує навічно |
| **close(ch)** | OK | OK | OK | panic | panic |
| **len(ch)** | 0 | 0-cap | cap | 0 | 0 |
| **cap(ch)** | 0 | cap | cap | cap | 0 |

---

## 🎓 Best Practices

### ✅ Роби так:

1. **Закривай канал на стороні відправника:**
```go
go func() {
    for i := 0; i < 10; i++ {
        ch <- i
    }
    close(ch) // ✅ sender закриває
}()
```

2. **Використовуй range для читання всіх значень:**
```go
for val := range ch {
    process(val)
}
```

3. **Використовуй select для таймаутів:**
```go
select {
case result := <-ch:
    return result
case <-time.After(1 * time.Second):
    return timeout
}
```

4. **Використовуй буферизовані канали для продуктивності:**
```go
results := make(chan Result, 100) // буфер для batch обробки
```

---

### ❌ Не роби так:

1. **Не закривай канал на стороні receiver:**
```go
// ❌ BAD
go func() {
    val := <-ch
    close(ch) // може бути ще інші readers!
}()
```

2. **Не пиши в закритий канал:**
```go
close(ch)
ch <- 42 // ❌ panic: send on closed channel
```

3. **Не закривай nil канал:**
```go
var ch chan int
close(ch) // ❌ panic: close of nil channel
```

4. **Не забувай про deadlocks:**
```go
ch := make(chan int)
ch <- 42 // ❌ deadlock! (unbuffered, немає reader)
```

---

## 🚨 Common Mistakes

### 1. Deadlock через unbuffered channel
```go
// ❌ Помилка
ch := make(chan int)
ch <- 42 // deadlock!

// ✅ Виправлення 1: використай buffered
ch := make(chan int, 1)
ch <- 42

// ✅ Виправлення 2: використай горутину
ch := make(chan int)
go func() {
    ch <- 42
}()
```

---

### 2. Range без close
```go
// ❌ Помилка
go func() {
    ch <- 1
    ch <- 2
    // забули close(ch)
}()

for val := range ch {
    fmt.Println(val) // зависне після всіх значень!
}

// ✅ Виправлення
go func() {
    ch <- 1
    ch <- 2
    close(ch) // обов'язково!
}()
```

---

### 3. Подвійне закриття
```go
// ❌ Помилка
close(ch)
close(ch) // panic: close of closed channel

// ✅ Виправлення: використай sync.Once
var once sync.Once
once.Do(func() { close(ch) })
```

---

### 4. Goroutine leak
```go
// ❌ Помилка (горутина зависне)
func leak() {
    ch := make(chan int)
    go func() {
        val := <-ch // ніхто ніколи не пише в ch!
        fmt.Println(val)
    }()
    // ch ніколи не використовується
}

// ✅ Виправлення: передай done канал
func noLeak(done <-chan struct{}) {
    ch := make(chan int)
    go func() {
        select {
        case val := <-ch:
            fmt.Println(val)
        case <-done:
            return // вихід при закритті done
        }
    }()
}
```

---

## 🔬 Тестування

Перевір, що все працює:

```bash
# Запусти приклади
go run all_channel_types.go

# Якщо все OK, побачиш:
# ╔══════════════════════════════════════════════════════════════════╗
# ║                     ✅ Completed!                                ║
# ╚══════════════════════════════════════════════════════════════════╝
```

---

## 📖 Додаткові ресурси

- [Effective Go - Channels](https://go.dev/doc/effective_go#channels)
- [Go by Example - Channels](https://gobyexample.com/channels)
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)

---

## 💡 Практика

Після вивчення спробуй:

1. **Worker Pool:** `week_5/practice/channel_patterns/main.go`
2. **Pipeline Pattern:** `week_5/practice/channel_patterns/`
3. **Fan-Out/Fan-In:** `week_5/practice/channel_patterns/`

---

**Створено:** 2026-01-19  
**Week:** 5 - Goroutines & Channels  
**Автор:** AI Assistant
