# Вправа 1: ValidationError система

## Ціль
Створити систему валідації з custom error types та детальними повідомленнями про помилки.

---

## Завдання

Створіть програму `validator.go`, яка:

1. Має custom error type `ValidationError` з детальною інформацією
2. Реалізує валідатор для User struct
3. Повертає всі помилки валідації одразу (не тільки першу)
4. Використовує sentinel errors для типових помилок
5. Демонструє правильну обробку помилок через `errors.Is()` та `errors.As()`

---

## Вимоги

### Sentinel Errors
```go
var (
    ErrRequired      = errors.New("field is required")
    ErrInvalidFormat = errors.New("invalid format")
    ErrOutOfRange    = errors.New("value out of range")
)
```

### Custom Error Type
```go
type ValidationError struct {
    Field   string
    Value   interface{}
    Rule    string
    Err     error  // Wrapped sentinel error
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Rule)
}

func (e ValidationError) Unwrap() error {
    return e.Err
}
```

### User Struct
```go
type User struct {
    Username string
    Email    string
    Age      int
    Password string
}
```

### Обов'язкові функції:

- `ValidateUser(user User) []error` - повертає всі помилки валідації
- `ValidateUsername(username string) error` - валідація username
- `ValidateEmail(email string) error` - валідація email
- `ValidateAge(age int) error` - валідація віку
- `ValidatePassword(password string) error` - валідація пароля

---

## Правила валідації

### Username:
- ✅ Не порожнє
- ✅ Мінімум 3 символи
- ✅ Максимум 20 символів
- ✅ Тільки літери, цифри, '_'

### Email:
- ✅ Не порожнє
- ✅ Містить '@'
- ✅ Містить '.'

### Age:
- ✅ Більше 0
- ✅ Менше 150

### Password:
- ✅ Не порожнє
- ✅ Мінімум 8 символів
- ✅ Містить хоча б одну цифру

---

## Приклад використання

```go
func main() {
    // Валідний користувач
    validUser := User{
        Username: "alice_123",
        Email:    "alice@example.com",
        Age:      25,
        Password: "securePass123",
    }
    
    errs := ValidateUser(validUser)
    if len(errs) == 0 {
        fmt.Println("✓ User is valid")
    }
    
    // Невалідний користувач
    invalidUser := User{
        Username: "ab",           // Занадто короткий
        Email:    "invalid",      // Немає @
        Age:      200,            // Занадто старий
        Password: "short",        // Занадто короткий + немає цифр
    }
    
    errs = ValidateUser(invalidUser)
    for _, err := range errs {
        fmt.Printf("❌ %v\n", err)
        
        // Перевірка через errors.Is()
        if errors.Is(err, ErrRequired) {
            fmt.Println("   Type: Required field")
        }
        
        // Отримання деталей через errors.As()
        var valErr ValidationError
        if errors.As(err, &valErr) {
            fmt.Printf("   Field: %s\n", valErr.Field)
            fmt.Printf("   Value: %v\n", valErr.Value)
            fmt.Printf("   Rule: %s\n", valErr.Rule)
        }
    }
}
```

---

## Очікуваний вивід

```
╔══════════════════════════════════════════╗
║        User Validation System            ║
╚══════════════════════════════════════════╝

🔹 Valid User
─────────────────────────────────────────
Username: alice_123
Email: alice@example.com
Age: 25
Password: ********

✓ All validations passed!


🔹 Invalid User
─────────────────────────────────────────
Username: ab
Email: invalid
Age: 200
Password: short

❌ Validation errors:

1. validation failed for username: must be between 3 and 20 characters
   Field: username
   Value: ab
   Rule: must be between 3 and 20 characters

2. validation failed for email: must contain '@'
   Field: email
   Value: invalid
   Rule: must contain '@'

3. validation failed for age: must be between 0 and 150
   Field: age
   Value: 200
   Rule: must be between 0 and 150

4. validation failed for password: must be at least 8 characters
   Field: password
   Value: short
   Rule: must be at least 8 characters

5. validation failed for password: must contain at least one digit
   Field: password
   Value: short
   Rule: must contain at least one digit


🔹 Error Type Detection
─────────────────────────────────────────
Detecting ErrOutOfRange: ✓ Found
Detecting ErrInvalidFormat: ✓ Found
Extracting ValidationError details: ✓ Success
```

---

## Підказки

### 1. ValidationError Implementation

```go
type ValidationError struct {
    Field   string
    Value   interface{}
    Rule    string
    Err     error
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Rule)
}

func (e ValidationError) Unwrap() error {
    return e.Err
}

func NewValidationError(field string, value interface{}, rule string, err error) ValidationError {
    return ValidationError{
        Field: field,
        Value: value,
        Rule:  rule,
        Err:   err,
    }
}
```

### 2. String Contains Check

```go
if !strings.Contains(email, "@") {
    return NewValidationError("email", email, "must contain '@'", ErrInvalidFormat)
}
```

### 3. Collecting Multiple Errors

```go
func ValidateUser(user User) []error {
    var errs []error
    
    if err := ValidateUsername(user.Username); err != nil {
        errs = append(errs, err)
    }
    
    if err := ValidateEmail(user.Email); err != nil {
        errs = append(errs, err)
    }
    
    // ... more validations
    
    return errs
}
```

---

## Бонус завдання

### 1. MultiError Type

Створіть custom error type для множинних помилок:

```go
type MultiError struct {
    Errors []error
}

func (m MultiError) Error() string {
    return fmt.Sprintf("%d validation errors occurred", len(m.Errors))
}

func (m MultiError) Unwrap() []error {
    return m.Errors
}
```

### 2. Field-Specific Validators

```go
type FieldValidator func(value interface{}) error

var validators = map[string][]FieldValidator{
    "username": {ValidateNotEmpty, ValidateLength(3, 20), ValidateAlphanumeric},
    "email":    {ValidateNotEmpty, ValidateEmailFormat},
    "age":      {ValidateRange(0, 150)},
    "password": {ValidateNotEmpty, ValidateMinLength(8), ValidateContainsDigit},
}
```

### 3. Custom Validation Rules

```go
type Rule interface {
    Validate(value interface{}) error
}

type LengthRule struct {
    Min int
    Max int
}

func (r LengthRule) Validate(value interface{}) error {
    str := value.(string)
    if len(str) < r.Min || len(str) > r.Max {
        return fmt.Errorf("length must be between %d and %d", r.Min, r.Max)
    }
    return nil
}
```

### 4. JSON Error Response

```go
type ErrorResponse struct {
    Field   string `json:"field"`
    Message string `json:"message"`
    Code    string `json:"code"`
}

func (v ValidationError) ToJSON() ErrorResponse {
    return ErrorResponse{
        Field:   v.Field,
        Message: v.Rule,
        Code:    v.Code(),
    }
}
```

---

## Критерії оцінки

- ✅ ValidationError реалізує error interface
- ✅ ValidationError має метод Unwrap()
- ✅ Всі валідації працюють коректно
- ✅ ValidateUser() повертає всі помилки одразу
- ✅ errors.Is() працює з wrapped errors
- ✅ errors.As() витягує ValidationError
- ✅ Код чистий і зрозумілий

---

## Рішення

Рішення знаходиться в `solutions/solution_1.go`.

Спробуйте виконати завдання самостійно перед тим, як дивитись рішення!

---

## Навчальні цілі

Після виконання цієї вправи ви будете вміти:
- Створювати custom error types
- Реалізовувати Unwrap() метод
- Використовувати errors.Is() для перевірки типу
- Використовувати errors.As() для витягування деталей
- Збирати множинні помилки валідації
- Додавати контекст до помилок

---

## Подальше вдосконалення

Подумайте як додати:
- Локалізацію повідомлень про помилки
- Динамічні правила валідації
- Валідацію вкладених структур
- Custom error codes
- Severity levels (error, warning, info)
- Валідацію через struct tags
