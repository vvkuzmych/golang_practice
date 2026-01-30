# Practice 01: Error Handling

## 🎯 Мета

Практика error handling в Go з `errors.Is()`, `errors.As()`, wrapping та custom error types.

---

## 🚀 Швидкий старт

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_15/practice/01_error_handling

# Запустити всі приклади
go run main.go

# Запустити тести
go test -v
```

---

## 📊 Що демонструється

### 1️⃣ Basic Error Handling
```go
user, err := GetUser(1)
if err != nil {
    fmt.Println("Error:", err)
}
```

### 2️⃣ Sentinel Errors (OLD vs NEW)
```go
// ❌ OLD WAY - breaks with wrapped errors
if err == ErrNotFound {
    // ...
}

// ✅ NEW WAY - works with wrapped errors
if errors.Is(err, ErrNotFound) {
    // ...
}
```

### 3️⃣ errors.Is() - Check Error Type
```go
_, err := SaveUser(0, "test@mail.com")

if errors.Is(err, ErrNotFound) {
    fmt.Println("User not found")
}
```

### 4️⃣ errors.As() - Extract Custom Type
```go
var valErr *ValidationError
if errors.As(err, &valErr) {
    fmt.Printf("Field: %s, Message: %s\n", 
        valErr.Field, valErr.Message)
}
```

### 5️⃣ Database Errors with Unwrap
```go
type DatabaseError struct {
    Query string
    Err   error
}

func (e *DatabaseError) Unwrap() error {
    return e.Err
}
```

### 6️⃣ Error Wrapping Chain
```go
return fmt.Errorf("failed to save user %d: %w", id, err)

// Check:
if errors.Is(err, ErrNotFound) {  // Works!
    // ...
}
```

### 7️⃣ Multiple Error Checks
```go
switch {
case errors.Is(err, ErrInvalidInput):
    return 400
case errors.Is(err, ErrNotFound):
    return 404
case errors.Is(err, ErrUnauthorized):
    return 401
}
```

### 8️⃣ Real-World API Example
```go
handleGetUser := func(id int) (int, string) {
    user, err := GetUser(id)
    if err != nil {
        switch {
        case errors.Is(err, ErrNotFound):
            return 404, "Not found"
        case errors.Is(err, ErrUnauthorized):
            return 401, "Unauthorized"
        default:
            return 500, "Internal error"
        }
    }
    return 200, user
}
```

---

## 📦 Структура

```
01_error_handling/
├── main.go              # Всі приклади
├── main_test.go         # Тести з errors.Is/As
└── README.md            # Цей файл
```

---

## 🎯 Key Concepts

### Sentinel Errors
```go
var (
    ErrNotFound      = errors.New("not found")
    ErrUnauthorized  = errors.New("unauthorized")
    ErrInvalidInput  = errors.New("invalid input")
)
```

### Custom Error Types
```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error: %s - %s", 
        e.Field, e.Message)
}
```

### Error Wrapping
```go
// Wrap with %w
return fmt.Errorf("operation failed: %w", originalErr)

// Check wrapped error
if errors.Is(err, originalErr) {  // ✅ Works!
    // ...
}
```

### Unwrap for Error Chains
```go
type MyError struct {
    Err error
}

func (e *MyError) Unwrap() error {
    return e.Err
}
```

---

## 📖 Output Example

```
╔════════════════════════════════════════╗
║   Go Error Handling Examples          ║
╚════════════════════════════════════════╝

1️⃣ Example 1: Basic Error Handling
─────────────────────────────────────────
✅ User: User-1
❌ Error: not found

2️⃣ Example 2: Sentinel Errors (OLD WAY - ❌)
─────────────────────────────────────────
✅ User not found (using ==)
❌ Can't check wrapped error with ==: user 0: not found

3️⃣ Example 3: errors.Is() (NEW WAY - ✅)
─────────────────────────────────────────
✅ User not found (using errors.Is)
✅ Unauthorized

4️⃣ Example 4: errors.As() for Custom Types
─────────────────────────────────────────
✅ Validation Error:
   Field: email
   Message: cannot be empty
✅ Wrapped Validation Error:
   Field: email
   Message: must contain @
   Full error: failed to save user 1: validation error: email - must contain @

...
```

---

## ✅ Висновок

### Що ти навчився:

- ✅ Використовувати `errors.Is()` замість `==`
- ✅ Використовувати `errors.As()` для type assertion
- ✅ Створювати custom error types
- ✅ Wrapping errors з `fmt.Errorf("%w")`
- ✅ Імплементувати `Unwrap()` method
- ✅ Handling errors в real-world scenarios

### Golden Rules:

1. **errors.Is()** для sentinel errors
2. **errors.As()** для custom types
3. **Wrap** errors з контекстом (`%w`)
4. **Implement Unwrap()** для error chains

**Week 15: Error Handling!** ⚠️✅
