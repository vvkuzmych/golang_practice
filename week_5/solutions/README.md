# Solutions — Week 5

## 📋 Рішення вправ

Тут містяться рішення для всіх вправ тижня 5.

---

## 📂 Файли

| Solution | Exercise | Тема | Складність |
|----------|----------|------|------------|
| `solution_1.go` | Exercise 1 | Pipeline with Goroutines | ⭐⭐ |
| `solution_2.go` | Exercise 2 | Worker Pool | ⭐⭐⭐ |
| `solution_3.go` | Exercise 3 | Graceful Shutdown | ⭐⭐⭐⭐ |

---

## 🚀 Як запускати

### Solution 1: Pipeline

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_5/solutions
go run solution_1.go
```

**Очікуваний результат:**
- 20 чисел оброблених через pipeline
- Generator → Processor → Consumer
- Коректне закриття channels

---

### Solution 2: Worker Pool

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_5/solutions
go run solution_2.go
```

**Очікуваний результат:**
- 100 jobs оброблені 5 workers
- ~10 jobs з помилками (кожен 10-й)
- Статистика: Success/Error count

---

### Solution 3: Graceful Shutdown

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_5/solutions
go run solution_3.go

# Press Ctrl+C to trigger graceful shutdown
```

**Альтернативно (auto-shutdown demo):**
```bash
timeout 5 go run solution_3.go
```

**Очікуваний результат:**
- Сервіс працює і генерує jobs
- Ctrl+C зупиняє прийом нових jobs
- Workers завершують поточні jobs
- Коректний cleanup
- Статистика виведена

---

## 📊 Покриття концепцій

### Solution 1: Pipeline
- ✅ Goroutines
- ✅ Unbuffered channels
- ✅ Channel closure
- ✅ Range over channel
- ✅ Unidirectional channels

### Solution 2: Worker Pool
- ✅ Buffered channels
- ✅ WaitGroup
- ✅ Multiple workers
- ✅ Error handling
- ✅ Result aggregation

### Solution 3: Graceful Shutdown
- ✅ Signal handling (SIGINT, SIGTERM)
- ✅ Context cancellation
- ✅ Multi-stage shutdown
- ✅ Timeout pattern
- ✅ Statistics tracking

---

## 💡 Ключові takeaways

### 1. Pipeline Pattern

**Принцип:** Кожен етап - окрема goroutine, channels для communication

```go
// ✅ Good
numbers := generator()
processed := processor(numbers)
consumer(processed)

// ❌ Bad (все в одній goroutine)
for i := range numbers {
    processed := process(i)
    consume(processed)
}
```

**Переваги:**
- Паралелізм
- Модульність
- Легко масштабується

---

### 2. Worker Pool Pattern

**Принцип:** Обмежена кількість workers, unbounded jobs

```go
// ✅ Good (bounded concurrency)
for w := 1; w <= numWorkers; w++ {
    go worker(jobs, results)
}

// ❌ Bad (unbounded goroutines)
for job := range jobs {
    go process(job)  // Може створити мільйони goroutines!
}
```

**Переваги:**
- Контроль concurrency
- Ефективне використання ресурсів
- Уникає resource exhaustion

---

### 3. Graceful Shutdown

**Принцип:** Multi-stage shutdown з timeout

```go
// Stage 1: Stop accepting new work
cancel()

// Stage 2: Wait for existing work (with timeout)
select {
case <-done:
    // Success
case <-time.After(timeout):
    // Force exit
}

// Stage 3: Cleanup
close(channels)
```

**Переваги:**
- Не втрачаємо дані
- Коректний cleanup
- Force exit якщо timeout

---

## ⚠️ Типові помилки

### ❌ Помилка 1: Забути close() channel

```go
// ❌ Bad
for i := 1; i <= 10; i++ {
    ch <- i
}
// Забули close(ch)!

for v := range ch {  // DEADLOCK!
    fmt.Println(v)
}
```

**Виправлення:**
```go
// ✅ Good
for i := 1; i <= 10; i++ {
    ch <- i
}
close(ch)  // ✅ Закриваємо!

for v := range ch {
    fmt.Println(v)
}
```

---

### ❌ Помилка 2: Не використовувати WaitGroup

```go
// ❌ Bad
for w := 1; w <= 3; w++ {
    go worker(w)
}
// Main може завершитись до workers!
```

**Виправлення:**
```go
// ✅ Good
var wg sync.WaitGroup
for w := 1; w <= 3; w++ {
    wg.Add(1)
    go worker(w, &wg)
}
wg.Wait()  // ✅ Чекаємо всіх!
```

---

### ❌ Помилка 3: Не закривати results channel

```go
// ❌ Bad
go func() {
    for j := range jobs {
        results <- process(j)
    }
    // Забули close(results)!
}()

for r := range results {  // Буде чекати вічно!
    fmt.Println(r)
}
```

**Виправлення:**
```go
// ✅ Good
go func() {
    for j := range jobs {
        results <- process(j)
    }
    close(results)  // ✅ Закриваємо!
}()

for r := range results {
    fmt.Println(r)
}
```

---

## 📚 Додаткові ресурси

### Theory Files:
- `theory/01_goroutine_basics.md` - goroutine lifecycle
- `theory/02_channels.md` - buffered vs unbuffered
- `theory/03_select_statement.md` - select usage
- `theory/04_deadlock.md` - deadlock scenarios
- `theory/05_channel_vs_queue.md` - channel vs queue

### Practice Examples:
- `practice/goroutine_basics/` - goroutine приклади
- `practice/channel_patterns/` - channel patterns
- `practice/worker_pool/` - worker pool приклади
- `practice/graceful_shutdown/` - shutdown patterns

---

## ✅ Self-check

Після вивчення solutions, ви повинні вміти:

- [ ] Створити pipeline з goroutines та channels
- [ ] Реалізувати worker pool з bounded concurrency
- [ ] Коректно закривати channels (без deadlock)
- [ ] Використовувати WaitGroup для синхронізації
- [ ] Обробляти помилки в concurrent code
- [ ] Реалізувати graceful shutdown з signal handling
- [ ] Використовувати context для cancellation
- [ ] Додавати timeout для operations

---

## 🎓 Наступні кроки

Після вивчення solutions:

1. **Спробуйте написати з нуля** (без підглядання)
2. **Експериментуйте:**
   - Змініть кількість workers
   - Додайте логування
   - Додайте metrics
3. **Бонус завдання** з exercises
4. **Створіть власний проект** з використанням цих patterns

---

**Удачі з вивченням! 🚀**
