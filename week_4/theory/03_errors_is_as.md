# errors.Is / errors.As в Go

## Проблема з простим порівнянням

```go
var ErrNotFound = errors.New("not found")

func findUser(id int) error {
    return fmt.Errorf("user %d: %w", id, ErrNotFound)
}

func main() {
    err := findUser(42)
    
    // ❌ Не працює!
    if err == ErrNotFound {
        fmt.Println("Not found")  // Не виконається!
    }
    
    // err це "user 42: not found", а не ErrNotFound
}
```

**Рішення:** `errors.Is()` та `errors.As()`

---

## errors.Is() - Перевірка типу помилки

### Базовий приклад

```go
package main

import (
    "errors"
    "fmt"
)

var ErrNotFound = errors.New("not found")

func findUser(id int) error {
    return fmt.Errorf("user %d: %w", id, ErrNotFound)
}

func main() {
    err := findUser(42)
    
    // ✅ Працює!
    if errors.Is(err, ErrNotFound) {
        fmt.Println("User not found")  // Виконається!
    }
}
```

### Як це працює?

`errors.Is()` розгортає (unwrap) помилки по ланцюжку:

```
errors.Is(err, target) {
    1. err == target? → return true
    2. err.Unwrap() → перевірити знову
    3. Повторювати до nil
    4. return false
}
```

---

## errors.As() - Type Assertion для помилок

### Базовий приклад

```go
package main

import (
    "errors"
    "fmt"
)

type ValidationError struct {
    Field string
    Value interface{}
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation failed: %s = %v", e.Field, e.Value)
}

func validateAge(age int) error {
    if age < 0 {
        return ValidationError{Field: "age", Value: age}
    }
    return nil
}

func main() {
    err := validateAge(-5)
    
    // Type assertion через errors.As
    var valErr ValidationError
    if errors.As(err, &valErr) {
        fmt.Printf("Field: %s, Value: %v\n", valErr.Field, valErr.Value)
        // Output: "Field: age, Value: -5"
    }
}
```

### Wrapped Error з errors.As

```go
func processUser(age int) error {
    err := validateAge(age)
    if err != nil {
        return fmt.Errorf("process failed: %w", err)
    }
    return nil
}

func main() {
    err := processUser(-5)
    
    // errors.As знайде ValidationError навіть після wrapping
    var valErr ValidationError
    if errors.As(err, &valErr) {
        fmt.Printf("Found validation error: %s\n", valErr.Field)
        // Працює! Навіть через wrapping
    }
}
```

---

## errors.Is vs == vs errors.As

| Метод | Використання | Wrapped errors | Type assertion |
|-------|--------------|----------------|----------------|
| `==` | Прямe порівняння | ❌ Не працює | ❌ Ні |
| `errors.Is()` | Порівняння sentinel errors | ✅ Працює | ❌ Ні |
| `errors.As()` | Отримати конкретний тип | ✅ Працює | ✅ Так |

### Приклади

```go
var ErrNotFound = errors.New("not found")

type CustomError struct {
    Code int
}

func (e CustomError) Error() string {
    return fmt.Sprintf("error code: %d", e.Code)
}

func main() {
    // Sentinel error (errors.Is)
    err1 := fmt.Errorf("wrapped: %w", ErrNotFound)
    if errors.Is(err1, ErrNotFound) {
        fmt.Println("✓ Found ErrNotFound")
    }
    
    // Custom type (errors.As)
    err2 := fmt.Errorf("wrapped: %w", CustomError{Code: 404})
    var customErr CustomError
    if errors.As(err2, &customErr) {
        fmt.Printf("✓ Found CustomError with code: %d\n", customErr.Code)
    }
}
```

---

## Production Patterns

### Pattern 1: HTTP Error Handling

```go
package main

import (
    "errors"
    "fmt"
    "net/http"
)

var (
    ErrNotFound      = errors.New("not found")
    ErrUnauthorized  = errors.New("unauthorized")
    ErrBadRequest    = errors.New("bad request")
)

type HTTPError struct {
    StatusCode int
    Message    string
    Err        error
}

func (e HTTPError) Error() string {
    return fmt.Sprintf("[%d] %s: %v", e.StatusCode, e.Message, e.Err)
}

func (e HTTPError) Unwrap() error {
    return e.Err
}

func handleError(w http.ResponseWriter, err error) {
    // Спробуємо отримати HTTPError
    var httpErr HTTPError
    if errors.As(err, &httpErr) {
        w.WriteHeader(httpErr.StatusCode)
        fmt.Fprintf(w, "Error: %s", httpErr.Message)
        return
    }
    
    // Перевіримо sentinel errors
    switch {
    case errors.Is(err, ErrNotFound):
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprint(w, "Not Found")
    case errors.Is(err, ErrUnauthorized):
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprint(w, "Unauthorized")
    default:
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprint(w, "Internal Server Error")
    }
}
```

### Pattern 2: Database Errors

```go
type DBError struct {
    Query string
    Err   error
}

func (e DBError) Error() string {
    return fmt.Sprintf("database error: %s", e.Query)
}

func (e DBError) Unwrap() error {
    return e.Err
}

func queryUser(id int) error {
    // Симуляція DB помилки
    return DBError{
        Query: "SELECT * FROM users WHERE id = ?",
        Err:   errors.New("connection refused"),
    }
}

func main() {
    err := queryUser(123)
    
    var dbErr DBError
    if errors.As(err, &dbErr) {
        fmt.Printf("Query failed: %s\n", dbErr.Query)
        fmt.Printf("Underlying error: %v\n", dbErr.Err)
    }
}
```

---

## Складні ланцюжки помилок

```go
package main

import (
    "errors"
    "fmt"
)

var (
    ErrDatabase = errors.New("database error")
    ErrNetwork  = errors.New("network error")
)

type ServiceError struct {
    Service string
    Err     error
}

func (e ServiceError) Error() string {
    return fmt.Sprintf("service %s failed", e.Service)
}

func (e ServiceError) Unwrap() error {
    return e.Err
}

func databaseLayer() error {
    return fmt.Errorf("connection timeout: %w", ErrDatabase)
}

func repositoryLayer() error {
    err := databaseLayer()
    if err != nil {
        return ServiceError{Service: "repository", Err: err}
    }
    return nil
}

func serviceLayer() error {
    err := repositoryLayer()
    if err != nil {
        return fmt.Errorf("service layer: %w", err)
    }
    return nil
}

func main() {
    err := serviceLayer()
    
    // errors.Is знайде ErrDatabase через весь ланцюжок
    if errors.Is(err, ErrDatabase) {
        fmt.Println("✓ Database error detected in deep chain")
    }
    
    // errors.As знайде ServiceError
    var serviceErr ServiceError
    if errors.As(err, &serviceErr) {
        fmt.Printf("✓ Service error: %s\n", serviceErr.Service)
    }
    
    // Повний error message
    fmt.Printf("\nFull error: %v\n", err)
}
```

**Output:**
```
✓ Database error detected in deep chain
✓ Service error: repository

Full error: service layer: service repository failed
```

---

## Коли використовувати що?

### errors.Is() - коли:
- ✅ Перевіряєте sentinel error
- ✅ Потрібно знати тип помилки, але не деталі
- ✅ Boolean check: "чи це та помилка?"

```go
if errors.Is(err, sql.ErrNoRows) {
    // handle "no rows"
}
```

### errors.As() - коли:
- ✅ Потрібні деталі з custom error type
- ✅ Хочете отримати додаткові поля
- ✅ Type assertion для structured errors

```go
var valErr ValidationError
if errors.As(err, &valErr) {
    fmt.Println("Field:", valErr.Field)
}
```

---

## Best Practices

### ✅ DO:

```go
// Використовуйте errors.Is для sentinel errors
if errors.Is(err, io.EOF) {
    // handle EOF
}

// Використовуйте errors.As для custom types
var netErr *net.OpError
if errors.As(err, &netErr) {
    // handle network error
}

// Створюйте Unwrap() для custom errors
func (e MyError) Unwrap() error {
    return e.OriginalErr
}
```

### ❌ DON'T:

```go
// ❌ Не використовуйте == для wrapped errors
if err == ErrNotFound {  // Не спрацює після wrapping
}

// ❌ Не використовуйте type assertion напряму
if _, ok := err.(ValidationError); ok {  // Не спрацює після wrapping
}

// ❌ Не порівнюйте через err.Error()
if err.Error() == "not found" {  // Крихке, легко ламається
}
```

---

## Ключові моменти

### 1. errors.Is() - Для перевірки sentinel errors

**Використання:** Коли потрібно просто **перевірити ТИП помилки** (є це та помилка чи ні?)

```go
var ErrNotFound = errors.New("not found")
var ErrTimeout = errors.New("timeout")

func process() error {
    return fmt.Errorf("database query: %w", ErrNotFound)
}

func main() {
    err := process()
    
    // errors.Is() - Boolean питання: "Це ErrNotFound?"
    if errors.Is(err, ErrNotFound) {
        fmt.Println("✓ Так, це ErrNotFound")
        // Але ми НЕ маємо доступу до деталей помилки
    }
    
    if errors.Is(err, ErrTimeout) {
        fmt.Println("Це timeout")  // Не виконається
    }
}
```

**Коли використовувати errors.Is():**
- ✅ Sentinel errors (заздалегідь оголошені як `var Err...`)
- ✅ Стандартні errors (`io.EOF`, `os.ErrNotExist`, `sql.ErrNoRows`)
- ✅ Коли потрібен простий Boolean check
- ✅ Коли деталі помилки не потрібні

**Аналогія:** Це як запитати "Чи є це яблуко?", не питаючи "Яке воно яблуко?"

---

### 2. errors.As() - Для type assertion custom errors

**Використання:** Коли потрібно **отримати ДЕТАЛІ** з custom error type

```go
type ValidationError struct {
    Field   string
    Value   interface{}
    Reason  string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation failed: %s", e.Field)
}

func validate(age int) error {
    if age < 0 {
        return ValidationError{
            Field:  "age",
            Value:  age,
            Reason: "must be positive",
        }
    }
    return nil
}

func main() {
    err := validate(-5)
    
    // errors.As() - Отримуємо ДЕТАЛІ помилки
    var valErr ValidationError
    if errors.As(err, &valErr) {
        fmt.Println("✓ Це ValidationError")
        // Маємо доступ до всіх полів!
        fmt.Printf("  Field: %s\n", valErr.Field)     // "age"
        fmt.Printf("  Value: %v\n", valErr.Value)     // -5
        fmt.Printf("  Reason: %s\n", valErr.Reason)   // "must be positive"
    }
}
```

**Коли використовувати errors.As():**
- ✅ Custom error types з додатковими полями
- ✅ Коли потрібні деталі помилки (fields, values, codes)
- ✅ Коли потрібно викликати методи на error type
- ✅ Structured errors для логування

**Аналогія:** Це як запитати "Чи є це яблуко? І якщо так, то якого сорту, кольору, розміру?"

---

### 3. Порівняння errors.Is() vs errors.As()

| Аспект | errors.Is() | errors.As() |
|--------|-------------|-------------|
| **Питання** | "Чи це та помилка?" | "Дай мені деталі цієї помилки" |
| **Результат** | `bool` | `bool` + заповнює target |
| **Use case** | Sentinel errors | Custom types з полями |
| **Доступ до полів** | ❌ Ні | ✅ Так |
| **Простота** | ✅ Простіше | Трохи складніше |

#### Візуальне порівняння:

```go
var ErrDatabase = errors.New("database error")

type DatabaseError struct {
    Query    string
    Duration time.Duration
    Err      error
}

func query() error {
    return DatabaseError{
        Query:    "SELECT * FROM users",
        Duration: 2 * time.Second,
        Err:      ErrDatabase,
    }
}

func main() {
    err := query()
    
    // errors.Is() - Перевірка типу (простий check)
    if errors.Is(err, ErrDatabase) {
        fmt.Println("✓ Database error detected")
        // Але Query, Duration недоступні!
    }
    
    // errors.As() - Отримання деталей (складніший check)
    var dbErr DatabaseError
    if errors.As(err, &dbErr) {
        fmt.Println("✓ Database error with details:")
        fmt.Printf("  Query: %s\n", dbErr.Query)         // Доступно!
        fmt.Printf("  Duration: %v\n", dbErr.Duration)   // Доступно!
        fmt.Printf("  Underlying: %v\n", dbErr.Err)      // Доступно!
    }
}
```

---

### 4. Практичні сценарії

#### Сценарій 1: HTTP Error Handling

```go
// Sentinel errors - використовуємо errors.Is()
var (
    ErrNotFound     = errors.New("not found")
    ErrUnauthorized = errors.New("unauthorized")
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
    err := processRequest(r)
    
    // Простий switch з errors.Is()
    switch {
    case errors.Is(err, ErrNotFound):
        w.WriteHeader(404)
        fmt.Fprint(w, "Not Found")
    case errors.Is(err, ErrUnauthorized):
        w.WriteHeader(401)
        fmt.Fprint(w, "Unauthorized")
    default:
        w.WriteHeader(500)
        fmt.Fprint(w, "Internal Server Error")
    }
}
```

#### Сценарій 2: Детальне логування

```go
// Custom error type - використовуємо errors.As()
type APIError struct {
    StatusCode int
    Method     string
    URL        string
    Err        error
}

func (e APIError) Error() string {
    return fmt.Sprintf("API error [%d]: %s %s", e.StatusCode, e.Method, e.URL)
}

func logError(err error) {
    var apiErr APIError
    if errors.As(err, &apiErr) {
        // Детальне логування з усіма полями
        log.Printf("API Error Details:")
        log.Printf("  Status: %d", apiErr.StatusCode)
        log.Printf("  Method: %s", apiErr.Method)
        log.Printf("  URL: %s", apiErr.URL)
        log.Printf("  Cause: %v", apiErr.Err)
    } else {
        // Звичайне логування
        log.Printf("Error: %v", err)
    }
}
```

#### Сценарій 3: Retry Logic

```go
var ErrTemporary = errors.New("temporary error")

type RetryableError struct {
    Attempt int
    Err     error
}

func (e RetryableError) Error() string {
    return fmt.Sprintf("retryable error (attempt %d): %v", e.Attempt, e.Err)
}

func processWithRetry(data string) error {
    for attempt := 1; attempt <= 3; attempt++ {
        err := process(data)
        
        // errors.Is() - перевірка чи temporary
        if errors.Is(err, ErrTemporary) {
            time.Sleep(time.Second * time.Duration(attempt))
            continue  // Retry
        }
        
        // errors.As() - отримання attempt number
        var retryErr RetryableError
        if errors.As(err, &retryErr) {
            log.Printf("Failed after %d attempts", retryErr.Attempt)
        }
        
        return err
    }
    return nil
}
```

---

### 5. Як вони працюють всередині (псевдокод)

#### errors.Is() алгоритм:

```go
func Is(err, target error) bool {
    for err != nil {
        if err == target {
            return true  // Знайшли!
        }
        // Unwrap і перевіряємо наступний в ланцюжку
        err = errors.Unwrap(err)
    }
    return false  // Не знайшли
}
```

**Візуалізація:**
```
err3: "service: %w"
  ↓ unwrap
err2: "repository: %w"
  ↓ unwrap
err1: ErrDatabase ← errors.Is() зупиняється ТУТ, якщо target == ErrDatabase
```

#### errors.As() алгоритм:

```go
func As(err error, target interface{}) bool {
    for err != nil {
        // Спробуємо type assertion
        if err, ok := err.(type of target); ok {
            *target = err  // Заповнюємо target
            return true
        }
        // Unwrap і перевіряємо наступний
        err = errors.Unwrap(err)
    }
    return false
}
```

**Візуалізація:**
```
err3: ServiceError{...}
  ↓ unwrap
err2: RepositoryError{...}
  ↓ unwrap
err1: DatabaseError{...} ← errors.As() КОПІЮЄ цей об'єкт в target
```

---

### 6. Комбінування errors.Is() та errors.As()

```go
func handleError(err error) {
    // Спочатку перевіримо чи це критична помилка
    if errors.Is(err, ErrCritical) {
        // Критична помилка - одразу panic
        panic(err)
    }
    
    // Тепер перевіримо чи є деталі
    var valErr ValidationError
    if errors.As(err, &valErr) {
        // Логуємо детально
        log.Printf("Validation failed:")
        log.Printf("  Field: %s", valErr.Field)
        log.Printf("  Value: %v", valErr.Value)
        return
    }
    
    // Стандартна обробка
    log.Printf("Error: %v", err)
}
```

---

### Підсумок

**errors.Is() = "Чи це та помилка?"**
- Boolean check
- Для sentinel errors
- Не дає доступу до полів

**errors.As() = "Дай мені цю помилку з деталями"**
- Type assertion + копіювання
- Для custom error types
- Повний доступ до всіх полів

**Обидва:**
- ✅ Працюють з wrapped errors
- ✅ Проходять весь error chain
- ✅ Не ламаються після fmt.Errorf("%w")
- ✅ Production-ready

---

## Домашнє завдання

1. Створіть custom error type з декількома полями
2. Заwrapте його кілька разів
3. Використайте errors.Is() для перевірки
4. Використайте errors.As() для отримання полів
5. Порівняйте з `==` (побачите що не працює)

---

**Наступна тема:** Context Basics (`04_context_basics.md`)
