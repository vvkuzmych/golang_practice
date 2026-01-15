# 🚀 QUICK START — Week 5

## Швидкий старт для тижня 5

### 1️⃣ Запустити демо

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_5
go run main.go
```

### 2️⃣ Почати з теорії

```bash
# Goroutine Basics
cat theory/01_goroutine_basics.md

# Channels (buffered vs unbuffered)
cat theory/02_channels.md

# Select statement
cat theory/03_select_statement.md

# Deadlock scenarios
cat theory/04_deadlock.md

# Channel vs Queue (ВАЖЛИВО!)
cat theory/05_channel_vs_queue.md
```

### 3️⃣ Запустити practice examples

```bash
# Goroutine basics
go run practice/goroutine_basics/main.go

# Channel patterns
go run practice/channel_patterns/main.go

# Worker pool
go run practice/worker_pool/main.go

# Graceful shutdown
timeout 8 go run practice/graceful_shutdown/main.go
```

### 4️⃣ Виконати exercises

```bash
# Exercise 1: Pipeline with goroutines
cat exercises/exercise_1.md

# Exercise 2: Worker pool
cat exercises/exercise_2.md

# Exercise 3: Graceful shutdown
cat exercises/exercise_3.md
```

### 5️⃣ Перевірити solutions

```bash
# Solution 1: Pipeline
go run solutions/solution_1.go

# Solution 2: Worker pool
go run solutions/solution_2.go

# Solution 3: Graceful shutdown (потребує Ctrl+C)
timeout 10 go run solutions/solution_3.go
```

---

## 📊 Структура навчання

```
День 1-2: Goroutines & Channels
  ├── theory/01_goroutine_basics.md
  ├── theory/02_channels.md
  └── practice/goroutine_basics/ + channel_patterns/

День 3-4: Select & Advanced
  ├── theory/03_select_statement.md
  ├── theory/04_deadlock.md
  ├── theory/05_channel_vs_queue.md
  └── practice/worker_pool/ + graceful_shutdown/

День 5-6: Exercises
  ├── exercises/exercise_1.md
  ├── exercises/exercise_2.md
  ├── exercises/exercise_3.md
  └── Порівняти з solutions/

День 7: Контроль
  └── Відповісти на контрольні питання
```

---

## ⚡ Ключові питання для контролю

### 1. Коли виникає deadlock?
- [ ] Всі goroutines заблоковані
- [ ] Відправка в unbuffered channel без receiver
- [ ] Range loop на channel що не закритий
- [ ] Циклічне очікування між goroutines

### 2. Чому channel — не queue?
- [ ] Channel для комунікації, Queue для зберігання
- [ ] Channel блокуючий (by design), Queue ні
- [ ] Channel - shared communication, Queue - shared state

---

## 🎯 Критерії готовності

Ви готові до контролю якщо:

✅ Розумію goroutine lifecycle
✅ Знаю різницю buffered vs unbuffered channel
✅ Вмію використовувати select
✅ Можу пояснити deadlock scenarios
✅ Розумію channel vs queue
✅ Реалізував worker pool
✅ Реалізував graceful shutdown

---

**Почніть з `go run main.go` і читайте теорію!** 🚀
