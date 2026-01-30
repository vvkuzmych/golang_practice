# errors.Is vs errors.As - Quick Reference Guide

## 🎯 Overview

In Go, **`errors.Is`** is used to check for equality with a specific error value (sentinel error), while **`errors.As`** is used to check if an error is of a specific type and extract its underlying value. 

**Both functions traverse the entire chain of wrapped errors.**

---

## 📊 Comparison Table

| Feature | errors.Is | errors.As |
|---------|-----------|-----------|
| **Purpose** | Check for equality with a target error value. | Check for a specific type and extract the value. |
| **Target Parameter** | Takes the error and a target error value (e.g., `os.ErrNotExist`). | Takes the error and a pointer to a variable of the target custom error type (e.g., `*MyCustomError`). |
| **Return Value** | Returns `true` if a match is found in the error chain, `false` otherwise. | Returns `true` if a match is found and assigns the error to the target pointer, `false` otherwise. |
| **When to use** | Use with predefined sentinel errors (e.g., `io.EOF`, `os.ErrPermission`). | Use with custom error types when you need to access specific data fields or methods of that type. |
| **Analogy** | Like checking if a wrapped gift is a specific, known item. | Like checking if a wrapped gift is a specific kind of item (e.g., a book) and then using it as such. |

---

## 🎯 In Short

- **`errors.Is`** checks **what** the error is (by value equality)
- **`errors.As`** checks **which type** of error it is (by type assertion) and lets you work with that specific type

---

## 💻 Example Usage

```go
package main

import (
	"errors"
	"fmt"
)

// Define a custom error type with extra fields
type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field %s: %s", e.Field, e.Reason)
}

func operation() error {
	// Wrap the custom error to add context
	return fmt.Errorf("operation failed: %w", &ValidationError{
		Field:  "email",
		Reason: "invalid format",
	})
}

func main() {
	err := operation()

	// Use errors.Is to check for a specific sentinel value (not applicable here, but for illustration)
	// Example: if errors.Is(err, os.ErrNotExist) {...}

	// Use errors.As to check for the custom type and extract its value
	var validationErr *ValidationError
	if errors.As(err, &validationErr) {
		// Can access fields of the custom error type
		fmt.Printf("Handled specific validation error for field: %s, reason: %s\n",
			validationErr.Field, validationErr.Reason)
	} else {
		fmt.Printf("Handled general error: %v\n", err)
	}
}
```

---

## 🔍 Detailed Examples

### Example 1: Using errors.Is

```go
package main

import (
	"errors"
	"fmt"
	"os"
)

var ErrNotFound = errors.New("not found")

func getUser(id int) error {
	if id == 0 {
		return fmt.Errorf("failed to get user: %w", ErrNotFound)
	}
	return nil
}

func main() {
	err := getUser(0)

	// Check for sentinel error
	if errors.Is(err, ErrNotFound) {
		fmt.Println("✅ User not found (using errors.Is)")
	}

	// Check for os.ErrNotExist
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("This won't print")
	}
}
```

**Output:**
```
✅ User not found (using errors.Is)
```

---

### Example 2: Using errors.As

```go
package main

import (
	"errors"
	"fmt"
)

type NetworkError struct {
	Host string
	Port int
	Err  error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("network error: %s:%d - %v", e.Host, e.Port, e.Err)
}

func connectToDatabase() error {
	return fmt.Errorf("database connection failed: %w", &NetworkError{
		Host: "localhost",
		Port: 5432,
		Err:  errors.New("connection refused"),
	})
}

func main() {
	err := connectToDatabase()

	// Extract custom error type
	var netErr *NetworkError
	if errors.As(err, &netErr) {
		fmt.Printf("✅ Network error detected!\n")
		fmt.Printf("   Host: %s\n", netErr.Host)
		fmt.Printf("   Port: %d\n", netErr.Port)
		fmt.Printf("   Error: %v\n", netErr.Err)
	}
}
```

**Output:**
```
✅ Network error detected!
   Host: localhost
   Port: 5432
   Error: connection refused
```

---

### Example 3: Combining Both

```go
package main

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = errors.New("not found")
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation: %s - %s", e.Field, e.Message)
}

func processUser(email string) error {
	if email == "" {
		return fmt.Errorf("process failed: %w", &ValidationError{
			Field:   "email",
			Message: "cannot be empty",
		})
	}
	if email == "admin@example.com" {
		return fmt.Errorf("process failed: %w", ErrNotFound)
	}
	return nil
}

func main() {
	// Test 1: ValidationError
	err1 := processUser("")

	// Use errors.As to extract custom type
	var valErr *ValidationError
	if errors.As(err1, &valErr) {
		fmt.Printf("✅ Validation error: %s - %s\n", valErr.Field, valErr.Message)
	}

	// Test 2: Sentinel error
	err2 := processUser("admin@example.com")

	// Use errors.Is to check sentinel
	if errors.Is(err2, ErrNotFound) {
		fmt.Println("✅ User not found")
	}
}
```

**Output:**
```
✅ Validation error: email - cannot be empty
✅ User not found
```

---

## 🎓 When to Use Which

### Use `errors.Is` when:

✅ Checking for **sentinel errors** (predefined constants)
```go
var ErrNotFound = errors.New("not found")

if errors.Is(err, ErrNotFound) {
    // Handle not found
}
```

✅ Checking for **standard library errors**
```go
if errors.Is(err, io.EOF) {
    // Handle end of file
}

if errors.Is(err, os.ErrNotExist) {
    // Handle file not found
}
```

---

### Use `errors.As` when:

✅ Need to **access fields** of custom error
```go
var valErr *ValidationError
if errors.As(err, &valErr) {
    fmt.Println(valErr.Field)  // Access field
    fmt.Println(valErr.Message) // Access message
}
```

✅ Need to **call methods** on custom error
```go
type TimeoutError struct {
    Duration time.Duration
}

func (e *TimeoutError) IsRetryable() bool {
    return e.Duration < 5*time.Second
}

var timeoutErr *TimeoutError
if errors.As(err, &timeoutErr) {
    if timeoutErr.IsRetryable() {
        // Retry operation
    }
}
```

---

## ⚠️ Common Mistakes

### ❌ Mistake 1: Using == instead of errors.Is

```go
// ❌ BAD - doesn't work with wrapped errors
if err == ErrNotFound {
    // Won't match if error is wrapped
}

// ✅ GOOD - works with wrapped errors
if errors.Is(err, ErrNotFound) {
    // Will match even if wrapped
}
```

---

### ❌ Mistake 2: Using type assertion instead of errors.As

```go
// ❌ BAD - doesn't work with wrapped errors
if valErr, ok := err.(*ValidationError); ok {
    // Won't match if error is wrapped
}

// ✅ GOOD - works with wrapped errors
var valErr *ValidationError
if errors.As(err, &valErr) {
    // Will match even if wrapped
}
```

---

### ❌ Mistake 3: Forgetting to pass pointer to errors.As

```go
// ❌ BAD - valErr is not a pointer
var valErr ValidationError
if errors.As(err, &valErr) {  // Won't compile!
    // ...
}

// ✅ GOOD - valErr is a pointer
var valErr *ValidationError
if errors.As(err, &valErr) {
    // ...
}
```

---

## 📊 Decision Tree

```
Do you need to check if error is a specific value?
│
├─ Yes → Use errors.Is()
│        Example: errors.Is(err, ErrNotFound)
│
└─ No → Do you need to extract custom error type?
         │
         ├─ Yes → Use errors.As()
         │        Example: errors.As(err, &customErr)
         │
         └─ No → Just check if err != nil
```

---

## ✅ Best Practices

1. **Always use `errors.Is()` instead of `==`** for sentinel errors
2. **Always use `errors.As()` instead of type assertion** for custom types
3. **Use pointer types** with `errors.As()` (e.g., `*ValidationError`)
4. **Export sentinel errors** so callers can check them
5. **Implement `Unwrap()`** on custom errors if they wrap other errors

---

## 🔗 Related Topics

- [Error Handling Theory](./theory/01_error_handling.md)
- [Error Handling Practice](./practice/01_error_handling/)
- [Week 15 README](./README.md)

---

**Week 15: errors.Is vs errors.As Master!** ⚠️✅

**Quick Rule:** `errors.Is` for WHAT (value), `errors.As` for TYPE (and extract)!
