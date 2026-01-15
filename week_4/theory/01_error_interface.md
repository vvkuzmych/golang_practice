# Error Interface в Go

## Що таке error?

В Go `error` - це просто **interface**:

```go
type error interface {
    Error() string
}
```

Будь-який тип, що має метод `Error() string`, автоматично реалізує interface `error`.

---

## Базове використання

### 1. Повернення помилки

```go
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

func main() {
    result, err := divide(10, 0)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Result:", result)
}
```

**Важливо:** В Go помилка завжди повертається як **останнє** значення.

---

## Створення помилок

### 1. errors.New()

```go
import "errors"

var ErrNotFound = errors.New("not found")
var ErrUnauthorized = errors.New("unauthorized")

func getUser(id int) (*User, error) {
    if id == 0 {
        return nil, ErrNotFound
    }
    // ...
    return user, nil
}
```

### 2. fmt.Errorf()

```go
import "fmt"

func processFile(filename string) error {
    if filename == "" {
        return fmt.Errorf("filename cannot be empty")
    }
    // ...
    return nil
}
```

### 3. Custom Error Type

```go
type ValidationError struct {
    Field string
    Value interface{}
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation failed for field '%s': invalid value %v", 
        e.Field, e.Value)
}

// Використання
func validateAge(age int) error {
    if age < 0 {
        return ValidationError{
            Field: "age",
            Value: age,
        }
    }
    return nil
}
```

---

## Sentinel Errors

**Sentinel error** - це заздалегідь оголошена змінна з помилкою, яку можна порівнювати.

```go
var (
    ErrNotFound      = errors.New("resource not found")
    ErrAlreadyExists = errors.New("resource already exists")
    ErrInvalidInput  = errors.New("invalid input provided")
)

func findResource(id int) error {
    // ...
    return ErrNotFound
}

func main() {
    err := findResource(123)
    if err == ErrNotFound {
        fmt.Println("Resource not found!")
    }
}
```

### Переваги sentinel errors:
- ✅ Можна порівнювати напряму (`==`)
- ✅ Зрозуміло для користувачів API
- ✅ Документовано у коді
- ✅ Typed errors

---

## nil Error

Якщо функція виконалась успішно, вона повертає `nil` як error:

```go
func doSomething() error {
    // Все добре
    return nil
}

func main() {
    err := doSomething()
    if err != nil {
        // Це НЕ виконається
        panic(err)
    }
    fmt.Println("Success!")
}
```

**Важливо:** Завжди перевіряйте `if err != nil`, а не `if err == nil`.

---

## Приклад з реального життя

```go
package main

import (
    "errors"
    "fmt"
)

var (
    ErrUserNotFound    = errors.New("user not found")
    ErrInvalidPassword = errors.New("invalid password")
)

type User struct {
    ID       int
    Username string
    Password string
}

var users = map[int]User{
    1: {ID: 1, Username: "alice", Password: "secret123"},
}

func authenticate(userID int, password string) (*User, error) {
    user, exists := users[userID]
    if !exists {
        return nil, ErrUserNotFound
    }

    if user.Password != password {
        return nil, ErrInvalidPassword
    }

    return &user, nil
}

func main() {
    // Успішна автентифікація
    user, err := authenticate(1, "secret123")
    if err != nil {
        fmt.Println("Auth failed:", err)
        return
    }
    fmt.Printf("✓ Authenticated: %s\n", user.Username)

    // Невірний пароль
    _, err = authenticate(1, "wrong")
    if err == ErrInvalidPassword {
        fmt.Println("✗ Invalid password")
    }

    // Користувач не знайдений
    _, err = authenticate(999, "password")
    if err == ErrUserNotFound {
        fmt.Println("✗ User not found")
    }
}
```

---

## Best Practices

### ✅ DO:

```go
// Оголошуйте sentinel errors як package-level змінні
var ErrNotFound = errors.New("not found")

// Завжди перевіряйте помилки
result, err := doSomething()
if err != nil {
    return err
}

// Додавайте контекст до помилок
if err != nil {
    return fmt.Errorf("failed to process user %d: %v", userID, err)
}

// Створюйте кастомні типи для складних помилок
type DatabaseError struct {
    Query string
    Err   error
}
```

### ❌ DON'T:

```go
// ❌ Не ігноруйте помилки
data, _ := os.ReadFile("file.txt")

// ❌ Не використовуйте panic для звичайних помилок
if err != nil {
    panic(err)  // Тільки для критичних помилок!
}

// ❌ Не втрачайте оригінальну помилку
if err != nil {
    return errors.New("something failed")  // Де оригінальна помилка?
}

// ❌ Не порівнюйте помилки через string
if err.Error() == "not found" {  // Погано!
    // ...
}
```

---

## Коли використовувати що?

| Ситуація | Рішення |
|----------|---------|
| Проста помилка з одним повідомленням | `errors.New()` |
| Помилка з динамічним текстом | `fmt.Errorf()` |
| Помилка, яку потрібно перевіряти ззовні | Sentinel error (`var Err...`) |
| Помилка з додатковими полями | Custom error type |
| Wrapping іншої помилки | `fmt.Errorf("... %w", err)` |

---

## Ключові моменти

1. ✅ `error` - це інтерфейс з одним методом `Error() string`
2. ✅ Повертайте `nil` при успіху
3. ✅ Завжди перевіряйте помилки (`if err != nil`)
4. ✅ Використовуйте sentinel errors для API
5. ✅ Створюйте custom types для складних випадків
6. ✅ Додавайте контекст до помилок

---

## Домашнє завдання

1. Створіть функцію `validateEmail(email string) error`, яка перевіряє email
2. Створіть sentinel error `ErrInvalidEmail`
3. Створіть custom error type `ValidationError` з полями `Field` і `Reason`
4. Напишіть функцію, яка повертає різні типи помилок

---

**Наступна тема:** Error Wrapping (`02_error_wrapping.md`)
