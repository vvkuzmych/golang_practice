# ТИЖДЕНЬ 4 — Error Handling + Context

**Ціль:** production-підхід до помилок і контролю виконання

---

## 📚 Структура тижня

```
week_4/
├── README.md              # Цей файл
├── theory/                # Теоретичні матеріали
│   ├── 01_error_interface.md
│   ├── 02_error_wrapping.md
│   ├── 03_errors_is_as.md
│   └── 04_context_basics.md
├── practice/              # Практичні приклади
│   ├── error_basics/      # Базова робота з помилками
│   ├── error_wrapping/    # Wrapping і unwrapping
│   ├── context_timeout/   # Context з timeout
│   └── context_cancellation/  # Скасування операцій
├── exercises/             # Завдання для виконання
│   ├── exercise_1.md      # Custom errors
│   ├── exercise_2.md      # Error wrapping chain
│   └── exercise_3.md      # Context timeout
└── solutions/             # Рішення завдань
    ├── solution_1.go
    ├── solution_2.go
    └── solution_3.go
```

---

## 📖 Теорія

### Що потрібно вивчити:

1. **error як interface** (`theory/01_error_interface.md`)
   - Що таке error interface
   - Створення власних помилок
   - Error() string method
   - nil errors
   - Sentinel errors

2. **Error Wrapping** (`theory/02_error_wrapping.md`)
   - fmt.Errorf з %w
   - Чому wrapping важливий
   - Збереження контексту
   - Stack trace альтернативи
   - Unwrap() метод

3. **errors.Is / errors.As** (`theory/03_errors_is_as.md`)
   - errors.Is() для перевірки типу
   - errors.As() для type assertion
   - Порівняння з ==
   - Wrapped errors chain
   - Best practices

4. **Context** (`theory/04_context_basics.md`)
   - context.Background()
   - context.WithCancel()
   - context.WithTimeout()
   - context.WithDeadline()
   - context.WithValue()
   - Правила роботи з context
   - ⚠️ Чому НЕ зберігати context в struct

---

## 💻 Практика

### Практика 1: Error Basics
**Папка:** `practice/error_basics/`

Демонстрація:
- Створення custom errors
- Sentinel errors pattern
- Error checking
- fmt.Errorf usage
- Nil error handling

### Практика 2: Error Wrapping
**Папка:** `practice/error_wrapping/`

Демонстрація:
- Wrapping з fmt.Errorf("%w")
- Unwrapping errors
- errors.Is() в дії
- errors.As() type assertion
- Error chains

### Практика 3: Context Timeout
**Папка:** `practice/context_timeout/`

Демонстрація:
- Функції з timeout
- context.WithTimeout()
- Обробка ctx.Done()
- select з context
- Graceful shutdown

### Практика 4: Context Cancellation
**Папка:** `practice/context_cancellation/`

Демонстрація:
- context.WithCancel()
- Manual cancellation
- Cascading cancellation
- Cleanup на cancel
- Context propagation

---

## ✅ Контроль знань

Ви повинні вміти пояснити:

### 1. error Interface
- Що таке error interface?
- Як створити власну помилку?
- Що таке sentinel error?
- Чому error може бути nil?

### 2. Error Wrapping
- Навіщо wrapping errors?
- Різниця між %v і %w?
- Як працює errors.Is()?
- Коли використовувати errors.As()?

### 3. Context
- Що таке context.Context?
- Різниця між Background() і TODO()?
- Як працює cancellation?
- Коли використовувати WithTimeout?

### 4. ⚠️ Важливе питання
- **Чому НЕ можна зберігати context в struct?**
  - Context має lifetime
  - Context прив'язаний до операції, не до об'єкта
  - Може призвести до memory leaks
  - Порушує ідіоматичність Go
  - **Правило:** context завжди передається як перший параметр функції

---

## 🎯 Як проходити тиждень

### День 1-2: Теорія
1. Прочитати `theory/01_error_interface.md`
2. Прочитати `theory/02_error_wrapping.md`
3. Прочитати `theory/03_errors_is_as.md`
4. Прочитати `theory/04_context_basics.md`
5. Запустити приклади з теорії

### День 3-4: Практика
1. Вивчити `practice/error_basics/`
2. Вивчити `practice/error_wrapping/`
3. Вивчити `practice/context_timeout/`
4. Вивчити `practice/context_cancellation/`
5. Експериментувати з кодом

### День 5-6: Вправи
1. Виконати `exercises/exercise_1.md` (ValidationError)
2. Виконати `exercises/exercise_2.md` (DB Error Chain)
3. Виконати `exercises/exercise_3.md` (API with Timeout)
4. Порівняти з рішеннями

### День 7: Контроль
1. Відповісти на питання контролю
2. Створити власний error type
3. Написати функцію з context timeout
4. Пояснити чому context не в struct

---

## 📝 Критерії успіху

✅ Розумію error як interface
✅ Вмію створювати custom errors
✅ Знаю як wrapping працює
✅ Можу використовувати errors.Is/As
✅ Розумію context lifecycle
✅ Вмію працювати з timeouts
✅ Знаю правила роботи з context
✅ Можу пояснити чому context не в struct

---

## 🚀 Почати навчання

```bash
# Перейти в theory
cd /Users/vkuzm/GolandProjects/golang_practice/week_4/theory
cat 01_error_interface.md

# Запустити перший приклад
cd ../practice/error_basics
go run main.go

# Спробувати context timeout
cd ../context_timeout
go run main.go

# Виконати завдання
cd ../../exercises
cat exercise_1.md
```

---

## 💡 Ключові концепції

### error Interface
```go
// error - це просто interface
type error interface {
    Error() string
}

// Власна помилка
type MyError struct {
    Code int
    Msg  string
}

func (e MyError) Error() string {
    return fmt.Sprintf("error %d: %s", e.Code, e.Msg)
}
```

### Error Wrapping
```go
// ❌ Погано - втрачаємо оригінальну помилку
if err != nil {
    return fmt.Errorf("failed to open file: %v", err)
}

// ✅ Добре - зберігаємо оригінальну помилку
if err != nil {
    return fmt.Errorf("failed to open file: %w", err)
}

// Перевірка wrapped error
if errors.Is(err, os.ErrNotExist) {
    // файл не існує
}
```

### Context Timeout
```go
// Context з timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Використання в функції
result, err := fetchDataWithContext(ctx, url)
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        // timeout!
    }
}
```

### Context Cancellation
```go
ctx, cancel := context.WithCancel(context.Background())

go func() {
    // Довга операція
    select {
    case <-ctx.Done():
        // Операція скасована!
        return
    case result := <-workChan:
        // Обробка результату
    }
}()

// Скасувати операцію
cancel()
```

---

## ⚠️ Поширені помилки

### Помилка 1: Ігнорування помилок
```go
// ❌ ДУЖЕ ПОГАНО
data, _ := os.ReadFile("config.json")

// ✅ Добре
data, err := os.ReadFile("config.json")
if err != nil {
    return fmt.Errorf("read config: %w", err)
}
```

### Помилка 2: Wrapping без контексту
```go
// ❌ Погано - немає контексту
if err != nil {
    return err
}

// ✅ Добре - додаємо контекст
if err != nil {
    return fmt.Errorf("processing user %d: %w", userID, err)
}
```

### Помилка 3: Context в struct
```go
// ❌ ДУ ЖЕ ПОГАНО!
type Service struct {
    ctx context.Context  // НІ! НІ! НІ!
    db  *sql.DB
}

// ✅ Добре - context як параметр
type Service struct {
    db *sql.DB
}

func (s *Service) Process(ctx context.Context, data string) error {
    // використовуємо ctx тут
}
```

### Помилка 4: Не перевіряти ctx.Done()
```go
// ❌ Погано - ігноруємо cancellation
func longOperation(ctx context.Context) error {
    for i := 0; i < 1000; i++ {
        doWork(i)
    }
}

// ✅ Добре - перевіряємо ctx.Done()
func longOperation(ctx context.Context) error {
    for i := 0; i < 1000; i++ {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            doWork(i)
        }
    }
}
```

---

## 🎓 Production Best Practices

### 1. Завжди Wrapping
```go
// Додавайте контекст до кожної помилки
return fmt.Errorf("save user %d to DB: %w", userID, err)
```

### 2. Sentinel Errors
```go
// Визначте публічні errors для API
var (
    ErrNotFound     = errors.New("not found")
    ErrUnauthorized = errors.New("unauthorized")
    ErrInvalidInput = errors.New("invalid input")
)
```

### 3. Context Propagation
```go
// Передавайте context вниз по call stack
func HandleRequest(ctx context.Context) error {
    data, err := fetchData(ctx)
    if err != nil {
        return err
    }
    return processData(ctx, data)
}
```

### 4. Timeouts для зовнішніх викликів
```go
// Завжди встановлюйте timeout для HTTP, DB, gRPC
ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
defer cancel()

resp, err := http.Get(ctx, url)
```

---

## 📊 Context Lifecycle

```
Request → context.Background()
    ↓
    WithTimeout(5s) → функція A
    ↓
    WithValue("userID", 123) → функція B
    ↓
    WithCancel() → горутина
    ↓
    Done() ← timeout / cancel / success
```

**Правила:**
1. Context створюється на початку операції
2. Передається як перший параметр (`ctx context.Context`)
3. Може бути скасований у будь-який момент
4. Дочірні contexts автоматично скасовуються при cancel батька
5. **НІКОЛИ** не зберігайте context в struct!

---

## 🔥 Чому context НЕ в struct?

### Причина 1: Lifetime
```go
// ❌ Context прив'язаний до операції, не до об'єкта
type UserService struct {
    ctx context.Context  // Який саме context? Від якого request?
}
```

### Причина 2: Memory Leaks
```go
// ❌ Context може жити довше ніж потрібно
service := &UserService{
    ctx: ctx,  // ctx може бути вже cancelled!
}
```

### Причина 3: Ідіоматичність Go
```go
// ✅ Правильний підхід
func (s *UserService) GetUser(ctx context.Context, id int) (*User, error) {
    // ctx приходить з request
    // ctx живе тільки під час цієї операції
    // ctx автоматично cleanup після завершення
}
```

**Винятки:** Тільки якщо у вас є дуже вагома причина і ви розумієте всі ризики. В 99.9% випадків - передавайте context як параметр!

---

## 🎓 Після тижня 4

Ви будете знати:
- Як правильно обробляти помилки
- Як створювати зрозумілі error messages
- Як використовувати context для control flow
- Чому context завжди параметр, ніколи поле struct

**Наступний крок:** Тиждень 5 - Concurrency (goroutines, channels)

---

## 📚 Корисні ресурси

- [Go Blog: Error handling and Go](https://go.dev/blog/error-handling-and-go)
- [Go Blog: Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors)
- [Go Blog: Context](https://go.dev/blog/context)
- [Effective Go: Errors](https://go.dev/doc/effective_go#errors)

---

**Удачі! 🎉**
