# Швидкий старт - Тиждень 4

## ⚡ Одна команда для запуску

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_4 && go run main.go
```

---

## 📁 Структура

```
week_4/
├── theory/          → Читай матеріали
├── practice/        → Запускай приклади
├── exercises/       → Виконуй завдання
└── solutions/       → Порівнюй рішення
```

---

## 🎯 План на тиждень

### 📖 День 1-2: Теорія
```bash
# Error Interface
cat theory/01_error_interface.md

# Error Wrapping
cat theory/02_error_wrapping.md

# errors.Is/As
cat theory/03_errors_is_as.md

# Context
cat theory/04_context_basics.md
```

### 💻 День 3-4: Практика
```bash
# Error Basics
cd practice/error_basics && go run main.go

# Error Wrapping
cd ../error_wrapping && go run main.go

# Context Timeout
cd ../context_timeout && go run main.go

# Context Cancellation
cd ../context_cancellation && go run main.go
```

### ✏️ День 5-6: Завдання
```bash
# Читай завдання
cd exercises
cat exercise_1.md
cat exercise_2.md
cat exercise_3.md

# Створюй файли для рішень
touch my_solution_1.go
touch my_solution_2.go
touch my_solution_3.go

# Порівнюй з рішеннями
cd ../solutions
cat solution_1.go
```

---

## 🚀 Швидкі приклади

### Error Basics
```go
package main

import (
    "errors"
    "fmt"
)

var ErrNotFound = errors.New("not found")

func main() {
    err := findUser(999)
    if errors.Is(err, ErrNotFound) {
        fmt.Println("User not found!")
    }
}

func findUser(id int) error {
    return fmt.Errorf("user %d: %w", id, ErrNotFound)
}
```

### Context Timeout
```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    result := make(chan string)
    go slowOperation(ctx, result)

    select {
    case <-ctx.Done():
        fmt.Println("Timeout!")
    case res := <-result:
        fmt.Println("Result:", res)
    }
}

func slowOperation(ctx context.Context, result chan<- string) {
    time.Sleep(5 * time.Second)
    result <- "done"
}
```

---

## ⚠️ Головне правило тижня

### ❌ НІКОЛИ НЕ РОБІТЬ ТАК:
```go
type Service struct {
    ctx context.Context  // ❌ НІ!
}
```

### ✅ ЗАВЖДИ РОБІТЬ ТАК:
```go
type Service struct {
    db *sql.DB  // ✅ Так!
}

func (s *Service) Process(ctx context.Context) error {
    // ctx як параметр!
}
```

**Чому?**
- Context має lifetime прив'язаний до операції
- Зберігання в struct → memory leaks
- Порушує Go idioms

---

## 📊 Checklist на кінець тижня

- [ ] Читав всю теорію
- [ ] Запустив всі practice приклади
- [ ] Виконав exercise_1 (Custom errors)
- [ ] Виконав exercise_2 (Error wrapping)
- [ ] Виконав exercise_3 (Context timeout)
- [ ] Розумію чому context НЕ в struct
- [ ] Можу пояснити errors.Is vs ==
- [ ] Знаю коли використовувати %w

---

## 🆘 Якщо застряг

1. **Проблема з errors:**
   - Перечитай `theory/02_error_wrapping.md`
   - Подивись `practice/error_wrapping/main.go`
   - Використовуй `%w` для wrapping

2. **Проблема з context:**
   - Перечитай `theory/04_context_basics.md`
   - Подивись `practice/context_timeout/main.go`
   - Пам'ятай: завжди `defer cancel()`

3. **Не розумію щось:**
   - Запитай в чаті/форумі
   - Перечитай теорію
   - Експериментуй з кодом

---

## 💡 Корисні команди

```bash
# Запустити конкретний приклад
go run practice/error_basics/main.go

# Перевірити код на помилки
go vet ./...

# Форматування коду
go fmt ./...

# Тести (якщо є)
go test ./...

# Детальний вивід
go run -v main.go
```

---

## 📚 Додаткове читання

- [Go by Example: Errors](https://gobyexample.com/errors)
- [Go by Example: Context](https://gobyexample.com/context)
- [Effective Go: Errors](https://go.dev/doc/effective_go#errors)

---

**Готовий? Почни з README.md! 🚀**
