# ТИЖДЕНЬ 5 — Goroutines & Channels

**Ціль:** освоїти основи конкурентності в Go

---

## 📚 Структура тижня

```
week_5/
├── README.md              # Цей файл
├── QUICK_START.md         # Швидкий старт
├── STATUS.md              # Статус проєкту
├── main.go                # Демонстраційний файл
├── theory/                # Теоретичні матеріали
│   ├── 01_goroutine_basics.md
│   ├── 02_channels.md
│   ├── 03_select_statement.md
│   ├── 04_deadlock.md
│   └── 05_channel_vs_queue.md
├── practice/              # Практичні приклади
│   ├── goroutine_basics/
│   ├── channel_patterns/
│   ├── worker_pool/
│   └── graceful_shutdown/
├── exercises/             # Завдання для виконання
│   ├── exercise_1.md      # Goroutines & Channels
│   ├── exercise_2.md      # Worker Pool
│   └── exercise_3.md      # Graceful Shutdown
└── solutions/             # Рішення завдань
    ├── README.md
    ├── solution_1.go
    ├── solution_2.go
    └── solution_3.go
```

---

## 📖 Теорія

### Що потрібно вивчити:

1.  **Goroutine Basics** (`theory/01_goroutine_basics.md`)
    -   Що таке goroutine
    -   Життєвий цикл goroutine
    -   `go` keyword
    -   Goroutines vs OS threads
    -   M:N scheduling
    -   GOMAXPROCS
    -   Race conditions
    -   WaitGroup для синхронізації

2.  **Channels** (`theory/02_channels.md`)
    -   Що таке channel
    -   Створення каналів (`make(chan T)`)
    -   Buffered vs unbuffered channels
    -   Відправка та отримання даних
    -   Закриття каналів (`close()`)
    -   Range over channel
    -   Unidirectional channels (`<-chan`, `chan<-`)
    -   Nil channels

3.  **Select Statement** (`theory/03_select_statement.md`)
    -   Синтаксис `select`
    -   Multiple channel operations
    -   Default case (non-blocking)
    -   Timeout patterns
    -   Context integration
    -   Select with closed channels

4.  **Deadlock** (`theory/04_deadlock.md`)
    -   Що таке deadlock
    -   Коли виникає deadlock
    -   Типові сценарії deadlock
    -   Як виявити deadlock
    -   Як уникнути deadlock
    -   Go runtime deadlock detection

5.  **Channel vs Queue** (`theory/05_channel_vs_queue.md`)
    -   **Чому channel — не queue**
    -   Різниця між channel та queue
    -   Синхронізація vs буферизація
    -   Communication vs data storage
    -   Best practices

---

## 💻 Практика

### Практика 1: Goroutine Basics
**Папка:** `practice/goroutine_basics/`

Демонстрація:
- Створення goroutines
- WaitGroup для очікування
- Race conditions
- Goroutine leak prevention

### Практика 2: Channel Patterns
**Папка:** `practice/channel_patterns/`

Демонстрація:
- Unbuffered channels
- Buffered channels
- Channel directions
- Range and close
- Select statement

### Практика 3: Worker Pool
**Папка:** `practice/worker_pool/`

Демонстрація:
- Job queue
- Worker goroutines
- Result collection
- Bounded concurrency

### Практика 4: Graceful Shutdown
**Папка:** `practice/graceful_shutdown/`

Демонстрація:
- Signal handling
- Context cancellation
- Worker cleanup
- Timeout on shutdown

---

## 📝 Вправи

### Exercise 1: Goroutines & Channels
**Файл:** `exercises/exercise_1.md`

**Завдання:**
- Створити pipeline з goroutines
- Generator → Processor → Consumer
- Використати unbuffered channels
- Коректно закрити канали

**Concepts:**
- Goroutine communication
- Channel lifecycle
- Data flow patterns

### Exercise 2: Worker Pool
**Файл:** `exercises/exercise_2.md`

**Завдання:**
- Реалізувати worker pool
- 5 workers обробляють 100 jobs
- Збирати результати
- Обробляти помилки

**Concepts:**
- Bounded concurrency
- Job distribution
- Result aggregation

### Exercise 3: Graceful Shutdown
**Файл:** `exercises/exercise_3.md`

**Завдання:**
- HTTP сервер з graceful shutdown
- Signal handling (SIGINT, SIGTERM)
- Context cancellation
- Cleanup before exit

**Concepts:**
- Signal handling
- Context propagation
- Resource cleanup

---

## ✅ Контроль знань

Ви повинні вміти пояснити:

### 1. Goroutines
- Що таке goroutine і як вона працює?
- В чому різниця між goroutine та OS thread?
- Що таке M:N scheduling?
- Як уникнути goroutine leak?
- Коли використовувати WaitGroup?

### 2. Channels
- В чому різниця між buffered та unbuffered channel?
- Що станеться якщо відправити в unbuffered channel без receiver?
- Коли закривати channel?
- Що станеться якщо читати з закритого channel?
- Що таке unidirectional channel і навіщо він потрібен?

### 3. Select
- Як працює `select` з кількома channels?
- Що робить `default` case?
- Як реалізувати timeout з `select`?
- Що станеться якщо всі cases заблоковані?

### 4. Deadlock
- **Коли виникає deadlock?** (Основне питання!)
    1. Всі goroutines заблоковані
    2. Немає можливості розблокуватись
    3. Програма не може продовжити виконання
- Типові сценарії deadlock:
    - Відправка в unbuffered channel без receiver
    - Циклічне очікування між goroutines
    - Забули закрити channel в range loop
- Як виявити: Go runtime викине `fatal error: all goroutines are asleep - deadlock!`

### 5. Channel vs Queue
- **Чому channel — не queue?** (Критичне питання!)
    1. **Призначення:** Channel для синхронізації та комунікації, Queue для зберігання
    2. **Семантика:** Channel блокуючий (by design), Queue не блокує
    3. **Ownership:** Channel - shared communication, Queue - shared state
    4. **Buffering:** Buffer в channel - оптимізація, не основна ціль
- Коли використовувати channel?
    - Для комунікації між goroutines
    - Для сигналів (done, stop)
    - Для координації роботи
- Коли використовувати queue?
    - Для зберігання великої кількості даних
    - Для persistence
    - Коли потрібна складна логіка (priority, requeue)

---

## 🎯 Як проходити тиждень

### День 1-2: Теорія Goroutines & Channels
1.  Прочитати `theory/01_goroutine_basics.md`
2.  Прочитати `theory/02_channels.md`
3.  Запустити приклади з `main.go` (секції 1-2)
4.  Запустити `practice/goroutine_basics/`
5.  Запустити `practice/channel_patterns/`

### День 3-4: Теорія Select & Advanced Topics
1.  Прочитати `theory/03_select_statement.md`
2.  Прочитати `theory/04_deadlock.md`
3.  Прочитати `theory/05_channel_vs_queue.md` (ВАЖЛИВО!)
4.  Запустити приклади з `main.go` (секції 3-5)

### День 5-6: Вправи
1.  Виконати `exercises/exercise_1.md` (Pipeline)
2.  Виконати `exercises/exercise_2.md` (Worker Pool)
3.  Виконати `exercises/exercise_3.md` (Graceful Shutdown)
4.  Порівняти з рішеннями в `solutions/`

### День 7: Контроль
1.  Відповісти на питання контролю знань
2.  **Пояснити словами:**
    - Коли виникає deadlock
    - Чому channel — не queue (3 причини!)
3.  Створити власний worker pool з нуля
4.  Переглянути слабкі місця

---

## 📝 Критерії успіху

✅ Розумію як працюють goroutines (lifecycle, scheduling)
✅ Розумію різницю між buffered та unbuffered channels
✅ Вмію використовувати WaitGroup для синхронізації
✅ Вмію використовувати `select` для роботи з кількома channels
✅ Розумію коли закривати channel
✅ Вмію реалізувати worker pool pattern
✅ Вмію реалізувати graceful shutdown
✅ **Можу пояснити коли виникає deadlock (3+ сценарії)**
✅ **Можу пояснити чому channel — не queue (3 причини)**
✅ Розумію race conditions та як їх уникати
✅ Написав працюючий concurrent код

---

## 🚀 Почати навчання

```bash
# Перейти в week_5
cd /Users/vkuzm/GolandProjects/golang_practice/week_5

# Запустити демонстраційний файл
go run main.go

# Прочитати README
cat README.md

# Почати вивчення теорії
cat theory/01_goroutine_basics.md
```

---

## 💡 Ключові концепції

### Goroutine
```go
// Запуск goroutine
go func() {
    fmt.Println("Hello from goroutine!")
}()

// WaitGroup для очікування
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    // робота
}()
wg.Wait()
```

### Channels
```go
// Unbuffered - блокує до receiver
ch := make(chan int)

// Buffered - блокує тільки коли повний
ch := make(chan int, 10)

// Відправка
ch <- 42

// Отримання
value := <-ch

// Закриття
close(ch)

// Range
for v := range ch {
    fmt.Println(v)
}
```

### Select
```go
select {
case msg := <-ch1:
    fmt.Println("Received from ch1:", msg)
case ch2 <- 42:
    fmt.Println("Sent to ch2")
case <-time.After(1 * time.Second):
    fmt.Println("Timeout!")
default:
    fmt.Println("No operations ready")
}
```

### Worker Pool
```go
func worker(id int, jobs <-chan int, results chan<- int) {
    for job := range jobs {
        results <- job * 2
    }
}

func main() {
    jobs := make(chan int, 100)
    results := make(chan int, 100)

    // Запуск workers
    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }

    // Відправка jobs
    for j := 1; j <= 9; j++ {
        jobs <- j
    }
    close(jobs)

    // Отримання results
    for a := 1; a <= 9; a++ {
        <-results
    }
}
```

---

## ⚠️ ВАЖЛИВІ МОМЕНТИ

### ❌ Deadlock - Коли виникає:

```go
// 1. Unbuffered channel без receiver
ch := make(chan int)
ch <- 42  // DEADLOCK! Ніхто не читає

// 2. Чекаємо на закритий channel в range
ch := make(chan int)
go func() {
    ch <- 1
    // Забули close(ch)!
}()
for v := range ch {  // DEADLOCK! Range чекає на close()
    fmt.Println(v)
}

// 3. Циклічне очікування
ch1 := make(chan int)
ch2 := make(chan int)
go func() {
    ch1 <- <-ch2  // Чекає ch2
}()
go func() {
    ch2 <- <-ch1  // Чекає ch1
}()
// DEADLOCK! Обидві goroutines чекають одна одну
```

### 📦 Channel vs Queue - Ключові різниці:

| Аспект | Channel | Queue |
|--------|---------|-------|
| **Призначення** | Синхронізація/комунікація | Зберігання даних |
| **Блокування** | Блокуючий (by design) | Non-blocking (зазвичай) |
| **Ownership** | Shared communication | Shared state |
| **Буфер** | Оптимізація, не ціль | Основна функція |
| **Use case** | Координація goroutines | Accumulation, processing |

**Висновок:** Channel — це інструмент для **communication**, а не для **data storage**!

---

## 🎓 Після тижня 5

Ви будете знати:
- Як писати concurrent код з goroutines
- Як використовувати channels для комунікації
- Як реалізовувати типові concurrency patterns
- Як уникати deadlock та race conditions
- Коли використовувати channel, а коли queue
- Як правильно завершувати concurrent програми

**Наступний крок:** Тиждень 6 - Advanced Concurrency Patterns

---

**Удачі! 🎉**
