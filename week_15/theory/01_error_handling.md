# Error Handling in Go

## 🎯 Що таке error?

**`error`** - це built-in interface:

```go
type error interface {
    Error() string
}
```

Будь-який тип з методом `Error() string` є error!

---

## 📊 Basic Error Handling

### Creating Errors

```go
import "errors"

// Simple error
err1 := errors.New("something went wrong")

// Formatted error
err2 := fmt.Errorf("user %d not found", 123)
```

### Checking Errors

```go
func Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

result, err := Divide(10, 0)
if err != nil {
    fmt.Println("Error:", err)
    return
}
fmt.Println("Result:", result)
```

---

## 🎯 Sentinel Errors

**Sentinel error** - заздалегідь визначена error constant.

```go
package mypackage

import "errors"

// Sentinel errors (exported)
var (
    ErrNotFound      = errors.New("not found")
    ErrUnauthorized  = errors.New("unauthorized")
    ErrInvalidInput  = errors.New("invalid input")
)

func GetUser(id int) (*User, error) {
    if id < 0 {
        return nil, ErrInvalidInput
    }
    
    user := db.Find(id)
    if user == nil {
        return nil, ErrNotFound
    }
    
    return user, nil
}
```

### Checking Sentinel Errors (Old Way)

```go
user, err := GetUser(123)
if err == ErrNotFound {  // ❌ Breaks with wrapped errors
    fmt.Println("User not found")
}
```

---

## ⚡ errors.Is() - Modern Way

**`errors.Is()`** перевіряє error в ланцюжку (wrapped errors).

```go
import "errors"

// Check error
if errors.Is(err, ErrNotFound) {  // ✅ Works with wrapped errors
    fmt.Println("User not found")
}
```

### Приклад з Wrapping

```go
func GetUserProfile(id int) error {
    user, err := GetUser(id)
    if err != nil {
        // Wrap error with context
        return fmt.Errorf("failed to get profile: %w", err)
    }
    return nil
}

// Use
err := GetUserProfile(123)

// ✅ errors.Is works через wrapping
if errors.Is(err, ErrNotFound) {
    fmt.Println("User not found")
}

// ❌ == не працює
if err == ErrNotFound {  // false!
    fmt.Println("Never prints")
}
```

---

## 🔍 Error Wrapping

### fmt.Errorf() з %w

```go
func ProcessOrder(id int) error {
    order, err := GetOrder(id)
    if err != nil {
        // %w wraps the error
        return fmt.Errorf("process order %d: %w", id, err)
    }
    return nil
}

// Error chain:
// "process order 123: order not found"
//                      ^^^^^^^^^^^^^^^
//                      Original error
```

### errors.Unwrap()

```go
err := ProcessOrder(123)

// Unwrap to get original error
unwrapped := errors.Unwrap(err)
fmt.Println(unwrapped)  // "order not found"
```

---

## 🎯 Custom Error Types

### Struct Error

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
}

func ValidateUser(u *User) error {
    if u.Email == "" {
        return &ValidationError{
            Field:   "email",
            Message: "cannot be empty",
        }
    }
    return nil
}
```

### errors.As() для Type Assertion

```go
err := ValidateUser(&User{})

// Type assertion with errors.As
var validationErr *ValidationError
if errors.As(err, &validationErr) {
    fmt.Printf("Field: %s\n", validationErr.Field)  // "email"
    fmt.Printf("Message: %s\n", validationErr.Message)  // "cannot be empty"
}
```

**Переваги `errors.As()`:**
- Працює з wrapped errors
- Type-safe
- Витягує конкретний тип error

---

## 📊 errors.Is() vs errors.As()

| | errors.Is() | errors.As() |
|---|-------------|-------------|
| **Use** | Check if error IS specific sentinel | Extract specific error TYPE |
| **Return** | bool | bool + populates target |
| **Example** | `errors.Is(err, ErrNotFound)` | `errors.As(err, &validationErr)` |

```go
// errors.Is - для sentinel errors
if errors.Is(err, ErrNotFound) {
    // Handle not found
}

// errors.As - для custom error types
var valErr *ValidationError
if errors.As(err, &valErr) {
    fmt.Println("Field:", valErr.Field)
}
```

---

## 🎯 Practical Examples

### Example 1: HTTP API Errors

```go
var (
    ErrNotFound      = errors.New("not found")
    ErrUnauthorized  = errors.New("unauthorized")
    ErrBadRequest    = errors.New("bad request")
)

type APIError struct {
    StatusCode int
    Message    string
    Err        error
}

func (e *APIError) Error() string {
    return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

func (e *APIError) Unwrap() error {
    return e.Err
}

func GetUser(id int) (*User, error) {
    if id < 0 {
        return nil, &APIError{
            StatusCode: 400,
            Message:    "invalid user id",
            Err:        ErrBadRequest,
        }
    }
    
    user := db.Find(id)
    if user == nil {
        return nil, &APIError{
            StatusCode: 404,
            Message:    "user not found",
            Err:        ErrNotFound,
        }
    }
    
    return user, nil
}

// Use
user, err := GetUser(-1)
if err != nil {
    var apiErr *APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("HTTP %d: %s\n", apiErr.StatusCode, apiErr.Message)
    }
    
    if errors.Is(err, ErrBadRequest) {
        // Handle bad request
    }
}
```

---

### Example 2: Database Errors

```go
type DBError struct {
    Query string
    Err   error
}

func (e *DBError) Error() string {
    return fmt.Sprintf("db error: %s (query: %s)", e.Err, e.Query)
}

func (e *DBError) Unwrap() error {
    return e.Err
}

func GetUserByEmail(email string) (*User, error) {
    query := "SELECT * FROM users WHERE email = ?"
    
    user := &User{}
    err := db.QueryRow(query, email).Scan(&user.ID, &user.Email)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, &DBError{
                Query: query,
                Err:   ErrNotFound,
            }
        }
        return nil, &DBError{
            Query: query,
            Err:   err,
        }
    }
    
    return user, nil
}

// Use
user, err := GetUserByEmail("test@example.com")
if err != nil {
    var dbErr *DBError
    if errors.As(err, &dbErr) {
        log.Printf("Query failed: %s", dbErr.Query)
    }
    
    if errors.Is(err, ErrNotFound) {
        // User not found
    }
}
```

---

### Example 3: Multi-Error

```go
type MultiError struct {
    Errors []error
}

func (m *MultiError) Error() string {
    var msgs []string
    for _, err := range m.Errors {
        msgs = append(msgs, err.Error())
    }
    return strings.Join(msgs, "; ")
}

func (m *MultiError) Add(err error) {
    if err != nil {
        m.Errors = append(m.Errors, err)
    }
}

func (m *MultiError) HasErrors() bool {
    return len(m.Errors) > 0
}

// Use
func ValidateUser(u *User) error {
    var errs MultiError
    
    if u.Email == "" {
        errs.Add(errors.New("email required"))
    }
    if u.Age < 18 {
        errs.Add(errors.New("must be 18+"))
    }
    if u.Name == "" {
        errs.Add(errors.New("name required"))
    }
    
    if errs.HasErrors() {
        return &errs
    }
    return nil
}
```

---

## ✅ Best Practices

### 1. Використовуй errors.Is() замість ==

```go
// ❌ BAD
if err == ErrNotFound {

// ✅ GOOD
if errors.Is(err, ErrNotFound) {
```

### 2. Використовуй errors.As() для type assertion

```go
// ❌ BAD
if valErr, ok := err.(*ValidationError); ok {

// ✅ GOOD
var valErr *ValidationError
if errors.As(err, &valErr) {
```

### 3. Wrap errors з контекстом

```go
// ❌ BAD
return err

// ✅ GOOD
return fmt.Errorf("failed to save user %d: %w", id, err)
```

### 4. Експортуй sentinel errors

```go
// ✅ GOOD - exported, can be checked by callers
var ErrNotFound = errors.New("not found")

// ❌ BAD - unexported, can't be checked
var errNotFound = errors.New("not found")
```

### 5. Custom errors з Unwrap()

```go
type MyError struct {
    Err error
}

func (e *MyError) Error() string {
    return fmt.Sprintf("my error: %v", e.Err)
}

// ✅ Implement Unwrap for errors.Is/As to work
func (e *MyError) Unwrap() error {
    return e.Err
}
```

---

## 🎓 Висновок

### Error Types:

✅ **Sentinel errors** - `var ErrNotFound = errors.New(...)`  
✅ **Custom types** - `type MyError struct {...}`  
✅ **Wrapped errors** - `fmt.Errorf("...: %w", err)`  

### Checking:

✅ **errors.Is()** - перевірка sentinel errors (працює з wrapped)  
✅ **errors.As()** - витягнути конкретний тип (працює з wrapped)  
✅ **errors.Unwrap()** - отримати оригінальний error  

### Golden Rule:

**"Always use errors.Is() and errors.As(), never == for errors!"**

**Week 15: Error Handling Master!** ⚠️🎯
