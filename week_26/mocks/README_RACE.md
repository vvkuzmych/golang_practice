# 🔍 Як перевірити Race Condition - Швидкий Start

## ⚡ TL;DR (Найважливіше)

```bash
# Перевірити race condition:
go test -race
go run -race main.go
```

---

## 🎯 Швидка демонстрація

### 1. Запусти простий приклад:
```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_26/mocks

# БЕЗ race detector - побачиш втрачені інкременти
go run simple_race_example.go

# З race detector - побачиш ДЕ саме проблема
go run -race simple_race_example.go
```

### Результат БЕЗ `-race`:
```
Очікували: 1000
Отримали:  938
❌ Втрачено 62 інкрементів!
```

### Результат З `-race`:
```
==================
WARNING: DATA RACE
Read at 0x00c000094038 by goroutine 8:
  main.main.func1()
      /path/to/simple_race_example.go:31 +0x68

Previous write at 0x00c000094038 by goroutine 7:
  main.main.func1()
      /path/to/simple_race_example.go:31 +0x78
==================

Found 2 data race(s)
```

**Race detector показує:**
- ✅ Точний номер рядка (line 31)
- ✅ Які goroutines конфліктують (goroutine 7 і 8)
- ✅ Що саме робили (Read/Write)
- ✅ Адресу пам'яті (0x00c000094038)

---

## 📋 Основні команди

| Команда | Що робить |
|---------|-----------|
| `go run -race main.go` | Запустити програму з race detector |
| `go test -race` | Запустити тести з race detector |
| `go test -race -v` | Детальний вивід тестів |
| `go test -race ./...` | Перевірити весь проект |
| `go build -race` | Білд з race detector (тільки для тестів!) |

---

## ❌ Приклад небезпечного коду

```go
counter := 0  // ❌ Shared state

for i := 0; i < 1000; i++ {
    go func() {
        counter++  // ❌ RACE CONDITION!
    }()
}
```

**Проблема:** Багато goroutines одночасно читають/пишуть `counter`

---

## ✅ Як виправити

### Варіант 1: Mutex
```go
var mu sync.Mutex
counter := 0

for i := 0; i < 1000; i++ {
    go func() {
        mu.Lock()
        counter++
        mu.Unlock()
    }()
}
```

### Варіант 2: Atomic
```go
var counter int64

for i := 0; i < 1000; i++ {
    go func() {
        atomic.AddInt64(&counter, 1)
    }()
}
```

### Варіант 3: Channel
```go
ch := make(chan int, 1000)

for i := 0; i < 1000; i++ {
    go func() {
        ch <- 1
    }()
}

counter := 0
for i := 0; i < 1000; i++ {
    counter += <-ch
}
```

---

## 📁 Файли в проекті

| Файл | Опис | Команда |
|------|------|---------|
| `simple_race_example.go` | Простий приклад (1 хвилина) | `go run simple_race_example.go` |
| `race_condition_demo.go` | Детальна демонстрація (5 хвилин) | `go run race_condition_demo.go` |
| `race_test.go` | Тести для URLService | `go test -race -v` |
| `RACE_CONDITION_GUIDE.md` | Повний гайд | Читати |
| `README_RACE.md` | Цей файл | Читати |

---

## 🎓 Де race detector НЕ працює

Race detector знаходить тільки race conditions які **відбулись під час виконання**.

❌ **Не знайде:**
```go
// Якщо цей код не виконався під час тестів
if someRareCondition {
    counter++ // Race condition тут, але не виконається
}
```

✅ **Рішення:** Пиши тести які покривають всі сценарії!

---

## ⚠️ Важливо знати

1. **Race detector - тільки для розробки!**
   - ❌ НЕ використовуй в production
   - ✅ Використовуй в тестах і CI/CD

2. **Performance:**
   - ⏱️ Час виконання: ~10x повільніше
   - 💾 Пам'ять: ~10x більше

3. **Coverage:**
   - Знаходить тільки те що виконалось
   - Потребує хороші тести

---

## 🚀 Практика

### Завдання 1: Знайди проблему
```bash
go run simple_race_example.go
# Скільки інкрементів втрачено?
```

### Завдання 2: Побач де проблема
```bash
go run -race simple_race_example.go
# На якому рядку race condition?
```

### Завдання 3: Перевір свій код
```bash
go test -race -v
# Чи є race conditions в твоєму коді?
```

---

## 💡 Коротке резюме

### Як перевірити:
```bash
go run -race main.go   # програма
go test -race          # тести
```

### Що шукати:
- ❌ Shared state без синхронізації
- ❌ Запис в map з різних goroutines
- ❌ Append в slice з різних goroutines

### Як виправити:
1. `sync.Mutex` - universal solution
2. `sync/atomic` - для простих операцій
3. `channels` - для комунікації
4. `sync.Map` - для concurrent maps

---

## 📚 Більше інформації

- 📖 Детальний гайд: `RACE_CONDITION_GUIDE.md`
- 🎥 Демонстрація: `race_condition_demo.go`
- 🧪 Тести: `race_test.go`
- 🌐 [Офіційна документація](https://go.dev/doc/articles/race_detector)

---

## ✅ Checklist

Перед commit:

- [ ] Запустив `go test -race`
- [ ] Перевірив всі concurrent code paths
- [ ] Додав тести для concurrent scenarios
- [ ] Використав mutex/channels де потрібно
- [ ] Code review пройдено

---

**Готовий перевірити свій код? 🚀**

```bash
go test -race -v
```
