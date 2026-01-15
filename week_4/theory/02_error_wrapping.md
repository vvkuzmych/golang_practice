# Error Wrapping в Go

## Проблема: Втрата контексту

Уявіть таку ситуацію:

```go
func processOrder(orderID int) error {
    err := validateOrder(orderID)
    if err != nil {
        return err  // Ми втрачаємо контекст!
    }
    return nil
}

func validateOrder(orderID int) error {
    return errors.New("invalid order")
}

// Output: "invalid order"
// Але ми не знаємо:
// - Який orderID?
// - Де це сталося?
// - Яка була причина?
```

---

## Рішення: Error Wrapping

### До Go 1.13 (старий спосіб)

```go
if err != nil {
    return fmt.Errorf("failed to process order %d: %v", orderID, err)
}
```

**Проблема:** Оригінальна помилка втрачена! Не можна порівнювати з `errors.Is()`.

### Go 1.13+ (правильний спосіб)

```go
if err != nil {
    return fmt.Errorf("failed to process order %d: %w", orderID, err)
}
```

**Ключова різниця:** `%w` замість `%v`!

---

## Wrapping з %w

### Базовий приклад

```go
package main

import (
    "errors"
    "fmt"
)

var ErrDatabase = errors.New("database error")

func getUser(id int) error {
    err := queryDatabase(id)
    if err != nil {
        // Wrapping: додаємо контекст + зберігаємо оригінальну помилку
        return fmt.Errorf("failed to get user %d: %w", id, err)
    }
    return nil
}

func queryDatabase(id int) error {
    // Симуляція DB помилки
    return fmt.Errorf("query failed: %w", ErrDatabase)
}

func main() {
    err := getUser(42)
    if err != nil {
        fmt.Println("Error:", err)
        // Output: "failed to get user 42: query failed: database error"

        // Можемо перевірити оригінальну помилку!
        if errors.Is(err, ErrDatabase) {
            fmt.Println("✓ Database error detected in chain")
        }
    }
}
```

---

## Ланцюжок помилок (Error Chain)

```go
func serviceLayer() error {
    err := repositoryLayer()
    if err != nil {
        return fmt.Errorf("service: %w", err)
    }
    return nil
}

func repositoryLayer() error {
    err := databaseLayer()
    if err != nil {
        return fmt.Errorf("repository: %w", err)
    }
    return nil
}

func databaseLayer() error {
    return errors.New("connection timeout")
}

func main() {
    err := serviceLayer()
    fmt.Println(err)
    // Output: "service: repository: connection timeout"
}
```

**Візуалізація:**
```
serviceLayer
    ↓ wraps
repositoryLayer
    ↓ wraps
databaseLayer (original error)
```

---

## Unwrapping помилок

### errors.Unwrap()

```go
package main

import (
    "errors"
    "fmt"
)

func main() {
    err1 := errors.New("original error")
    err2 := fmt.Errorf("wrapped: %w", err1)
    err3 := fmt.Errorf("twice wrapped: %w", err2)

    // Unwrap один рівень
    unwrapped := errors.Unwrap(err3)
    fmt.Println(unwrapped)  // "wrapped: original error"

    // Unwrap ще раз
    unwrapped = errors.Unwrap(unwrapped)
    fmt.Println(unwrapped)  // "original error"

    // Unwrap nil, якщо більше немає
    unwrapped = errors.Unwrap(unwrapped)
    fmt.Println(unwrapped)  // <nil>
}
```

---

## Custom Error Type з Unwrap

```go
type QueryError struct {
    Query string
    Err   error
}

func (e *QueryError) Error() string {
    return fmt.Sprintf("query failed: %s", e.Query)
}

// Unwrap дозволяє errors.Is/As працювати
func (e *QueryError) Unwrap() error {
    return e.Err
}

func main() {
    original := errors.New("connection refused")
    wrapped := &QueryError{
        Query: "SELECT * FROM users",
        Err:   original,
    }

    // errors.Is знайде original в ланцюжку
    if errors.Is(wrapped, original) {
        fmt.Println("✓ Found original error")
    }
}
```

---

## Production Pattern

```go
package main

import (
    "errors"
    "fmt"
)

// Sentinel errors
var (
    ErrNotFound      = errors.New("not found")
    ErrUnauthorized  = errors.New("unauthorized")
    ErrInvalidInput  = errors.New("invalid input")
)

// Service layer
type UserService struct {
    repo *UserRepository
}

func (s *UserService) GetUser(id int) (*User, error) {
    user, err := s.repo.Find(id)
    if err != nil {
        // Wrapping з контекстом
        return nil, fmt.Errorf("UserService.GetUser(id=%d): %w", id, err)
    }
    return user, nil
}

// Repository layer
type UserRepository struct{}

func (r *UserRepository) Find(id int) (*User, error) {
    // Симуляція DB query
    if id == 0 {
        return nil, fmt.Errorf("invalid id: %w", ErrInvalidInput)
    }
    return nil, fmt.Errorf("user %d: %w", id, ErrNotFound)
}

type User struct {
    ID   int
    Name string
}

func main() {
    repo := &UserRepository{}
    service := &UserService{repo: repo}

    user, err := service.GetUser(42)
    if err != nil {
        // Повний контекст помилки
        fmt.Println("Error:", err)
        // Output: "UserService.GetUser(id=42): user 42: not found"

        // Але можемо перевірити оригінальну помилку!
        if errors.Is(err, ErrNotFound) {
            fmt.Println("✓ User not found (sentinel error detected)")
        }
    } else {
        fmt.Println("User:", user)
    }
}
```

---

## %v vs %w - Різниця

| Формат | Wrapping | errors.Is працює | Use Case |
|--------|----------|-----------------|----------|
| `%v` | ❌ Ні | ❌ Ні | Логування, коли не потрібно порівнювати |
| `%w` | ✅ Так | ✅ Так | Production код, API, коли потрібно порівнювати |

```go
original := errors.New("db error")

// %v - НЕ wrapping
err1 := fmt.Errorf("failed: %v", original)
errors.Is(err1, original)  // false ❌

// %w - wrapping
err2 := fmt.Errorf("failed: %w", original)
errors.Is(err2, original)  // true ✅
```

---

## Коли використовувати wrapping?

### ✅ Використовуйте %w коли:

1. **Повертаєте помилку вгору по стеку**
```go
func handler() error {
    err := service.Process()
    if err != nil {
        return fmt.Errorf("handler: %w", err)
    }
    return nil
}
```

2. **Потрібно перевірити конкретний тип помилки**
```go
if errors.Is(err, sql.ErrNoRows) {
    // specific handling
}
```

3. **Додаєте контекст до помилки**
```go
return fmt.Errorf("user %d: %w", userID, err)
```

### ❌ НЕ використовуйте %w коли:

1. **Логування (використовуйте %v)**
```go
log.Printf("error occurred: %v", err)
```

2. **Не хочете exposing internal errors**
```go
// Приховуємо internal деталі від користувача
if err != nil {
    return errors.New("internal server error")
}
```

---

## Best Practices

### 1. Завжди додавайте контекст

```go
// ❌ Погано - немає контексту
if err != nil {
    return err
}

// ✅ Добре - додали контекст
if err != nil {
    return fmt.Errorf("failed to load config file %s: %w", filename, err)
}
```

### 2. Не wrappingуйте двічі в одному місці

```go
// ❌ Погано
if err != nil {
    err = fmt.Errorf("step 1: %w", err)
    err = fmt.Errorf("step 2: %w", err)  // Зайвий wrap
    return err
}

// ✅ Добре
if err != nil {
    return fmt.Errorf("failed at step 1 and 2: %w", err)
}
```

### 3. Wrapping на межах абстракцій

```go
// Handler → Service → Repository → Database

// Handler wraps Service errors
func (h *Handler) Process() error {
    err := h.service.Execute()
    if err != nil {
        return fmt.Errorf("handler: %w", err)
    }
    return nil
}

// Service wraps Repository errors
func (s *Service) Execute() error {
    err := s.repo.Save()
    if err != nil {
        return fmt.Errorf("service: %w", err)
    }
    return nil
}
```

---

## Ключові моменти

1. ✅ Використовуйте `%w` для wrapping
2. ✅ Додавайте корисний контекст до кожної помилки
3. ✅ Wrapping зберігає оригінальну помилку для `errors.Is/As`
4. ✅ Створюйте error chains через layers
5. ✅ `%v` для логування, `%w` для повернення
6. ✅ Unwrap можна через `errors.Unwrap()`

---

## Домашнє завдання

1. Створіть 3-рівневу архітектуру (handler → service → repository)
2. Додайте wrapping на кожному рівні
3. Перевірте, чи `errors.Is()` працює з wrapped errors
4. Експериментуйте з `errors.Unwrap()`

---

**Наступна тема:** errors.Is / errors.As (`03_errors_is_as.md`)
