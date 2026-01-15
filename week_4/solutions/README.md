# Solutions - Week 4

Рішення для вправ Тижня 4 (Error Handling + Context).

---

## 📁 Файли

### 1. solution_1.go - ValidationError System

**Тема:** Custom Error Types

**Що демонструє:**
- ✅ Custom error type (ValidationError)
- ✅ Unwrap() implementation
- ✅ errors.Is() для перевірки sentinel errors
- ✅ errors.As() для витягування деталей
- ✅ Множинні помилки валідації

**Запуск:**
```bash
cd solutions
go run solution_1.go
```

**Output:**
```
✓ Valid User: all validations passed
❌ Invalid User: 4 validation errors
✓ errors.Is() detects ErrOutOfRange
✓ errors.As() extracts ValidationError details
```

---

### 2. solution_2.go - Error Wrapping Chain

**Тема:** Multi-Layer Architecture

**Що демонструє:**
- ✅ 3-рівнева архітектура (Database → Repository → Service)
- ✅ Error wrapping з %w на кожному рівні
- ✅ Контекст додається до помилок
- ✅ errors.Is() працює через весь chain
- ✅ Error traversal (розгортання ланцюжка)

**Запуск:**
```bash
cd solutions
go run solution_2.go
```

**Output:**
```
Scenario 1: User Not Found
  Error chain: 4 levels deep
  ✓ Original ErrNotFound detected

Scenario 2: Database Connection Error
  ✓ ErrConnection detected through wrapping

Scenario 3: Successful Operation
  ✓ User created

Scenario 4: Duplicate Key Error
  ✓ ErrDuplicateKey detected
```

---

### 3. solution_3.go - HTTP Service with Context

**Тема:** Context Timeout & Cancellation

**Що демонструє:**
- ✅ Context в HTTP handlers
- ✅ WithTimeout для кожного request
- ✅ context.DeadlineExceeded handling
- ✅ **Context НЕ в struct** (передається як параметр)
- ✅ Graceful cancellation
- ✅ Demo mode з різними scenarios

**Запуск:**
```bash
cd solutions
go run solution_3.go

# Для real server mode (закоментуйте demo код):
# go run solution_3.go
# curl http://localhost:8080/users/1
```

**Output:**
```
Scenario 1: Fast Request (500ms)
  ✓ Success (200 OK)

Scenario 2: Slow Request (4s)
  ✗ Timeout after 3s (504 Gateway Timeout)

Scenario 3: Medium Request (2s)
  ✓ Success (200 OK)
```

---

## 🎯 Навчальні цілі

Після виконання всіх solutions ви будете вміти:

### Solution 1:
- Створювати custom error types
- Реалізовувати Unwrap() метод
- Використовувати errors.Is/As для перевірки
- Збирати множинні помилки

### Solution 2:
- Проектувати багаторівневу архітектуру
- Правильно wrapping помилок на кожному рівні
- Додавати корисний контекст
- Дебажити error chains

### Solution 3:
- Використовувати context в HTTP
- Встановлювати timeouts
- Обробляти cancellation
- Розуміти чому context НЕ в struct

---

## 🚀 Швидкий запуск всіх solutions

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_4/solutions

# Solution 1
echo "=== SOLUTION 1 ===" && go run solution_1.go

# Solution 2
echo "=== SOLUTION 2 ===" && go run solution_2.go

# Solution 3 (з timeout через demo)
echo "=== SOLUTION 3 ===" && timeout 8 go run solution_3.go
```

---

## 📊 Порівняння з вправами

| Exercise | Solution | Тема | Складність |
|----------|----------|------|------------|
| exercise_1.md | solution_1.go | Validation System | 🟡 Medium |
| exercise_2.md | solution_2.go | Error Wrapping | 🟡 Medium |
| exercise_3.md | solution_3.go | Context + HTTP | 🟠 Advanced |

---

## 💡 Ключові відмінності від exercises

### Solution 1:
- Додано error type detection demo
- Показано як працює errors.Is/As
- Додано password masking

### Solution 2:
- Додано error chain traversal
- Демонстрація всіх 4 scenarios
- Performance timing для операцій

### Solution 3:
- Додано demo mode з автоматичними тестами
- Різні затримки для різних user ID
- Показано як timeout працює в реальності

---

## ⚠️ Важливі моменти

### 1. Context НЕ в Struct

**❌ Погано:**
```go
type Service struct {
    ctx context.Context  // НІ!
}
```

**✅ Добре:**
```go
type Service struct {
    db *Database
}

func (s *Service) Process(ctx context.Context) error {
    // Context як параметр!
}
```

**Всі solutions дотримуються цього правила!**

### 2. Завжди %w для Wrapping

```go
// ✅ Solution 2 демонструє:
return fmt.Errorf("service: failed: %w", err)
```

### 3. errors.Is/As працюють через Wrapping

```go
// ✅ Solution 1 і 2 демонструють:
if errors.Is(err, ErrNotFound) {
    // Спрацює навіть після кількох wraps!
}
```

---

## 🎓 Подальше вдосконалення

### Solution 1:
- Додайте підтримку вкладених структур
- Реалізуйте MultiError type
- Додайте JSON serialization для помилок

### Solution 2:
- Додайте retry logic
- Реалізуйте circuit breaker
- Додайте structured logging

### Solution 3:
- Додайте graceful shutdown
- Реалізуйте middleware chain
- Додайте distributed tracing
- Реалізуйте rate limiting

---

## 🔧 Модифікація solutions

### Для локального тестування:

#### Solution 3 - Real Server Mode:

1. Відкрийте `solution_3.go`
2. В `main()` закоментуйте:
   ```go
   // runDemoRequests()
   ```
3. Додайте:
   ```go
   select {} // Блокуємо main
   ```
4. Запустіть:
   ```bash
   go run solution_3.go
   ```
5. В іншому терміналі:
   ```bash
   curl http://localhost:8080/users/1
   curl http://localhost:8080/users/2
   curl http://localhost:8080/health
   ```

---

## 📚 Зв'язок з теорією

| Solution | Theory Files |
|----------|--------------|
| solution_1.go | 01_error_interface.md<br>03_errors_is_as.md |
| solution_2.go | 02_error_wrapping.md<br>03_errors_is_as.md |
| solution_3.go | 04_context_basics.md |

---

## ✅ Checklist перед завершенням

- [ ] Запустив всі 3 solutions
- [ ] Розумію як працює ValidationError
- [ ] Розумію error wrapping chain
- [ ] Розумію чому context НЕ в struct
- [ ] Можу пояснити errors.Is vs errors.As
- [ ] Розумію коли використовувати WithTimeout
- [ ] Можу створити власний error type з Unwrap()

---

**Готово? Переходь до наступного тижня! 🎉**
