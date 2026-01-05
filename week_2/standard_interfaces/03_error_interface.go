package main

import (
	"errors"
	"fmt"
	"time"
)

// ============= error Interface =============

// type error interface {
//     Error() string
// }

// ============= Custom Errors =============

// ValidationError - помилка валідації
type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation failed on field '%s' (value: %v): %s",
		e.Field, e.Value, e.Message)
}

// NotFoundError - ресурс не знайдено
type NotFoundError struct {
	Resource string
	ID       int
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %d not found", e.Resource, e.ID)
}

// AuthError - помилка автентифікації
type AuthError struct {
	Username string
	Reason   string
}

func (e AuthError) Error() string {
	return fmt.Sprintf("authentication failed for user '%s': %s", e.Username, e.Reason)
}

// DBError - помилка бази даних
type DBError struct {
	Operation string
	Table     string
	Err       error
}

func (e DBError) Error() string {
	return fmt.Sprintf("database error during %s on table '%s': %v",
		e.Operation, e.Table, e.Err)
}

// Unwrap дозволяє обгортати помилки
func (e DBError) Unwrap() error {
	return e.Err
}

// NetworkError - мережева помилка
type NetworkError struct {
	Host      string
	Port      int
	Operation string
	Timestamp time.Time
}

func (e NetworkError) Error() string {
	return fmt.Sprintf("[%s] network error: %s failed for %s:%d",
		e.Timestamp.Format("15:04:05"), e.Operation, e.Host, e.Port)
}

// ============= Error with code =============

type APIError struct {
	Code    int
	Message string
	Details map[string]string
}

func (e APIError) Error() string {
	return fmt.Sprintf("API Error %d: %s", e.Code, e.Message)
}

func (e APIError) IsClientError() bool {
	return e.Code >= 400 && e.Code < 500
}

func (e APIError) IsServerError() bool {
	return e.Code >= 500
}

// ============= Functions that return errors =============

func ValidateEmail(email string) error {
	if email == "" {
		return ValidationError{
			Field:   "email",
			Value:   email,
			Message: "email cannot be empty",
		}
	}

	if len(email) < 5 {
		return ValidationError{
			Field:   "email",
			Value:   email,
			Message: "email too short",
		}
	}

	// Проста перевірка на @
	hasAt := false
	for _, ch := range email {
		if ch == '@' {
			hasAt = true
			break
		}
	}

	if !hasAt {
		return ValidationError{
			Field:   "email",
			Value:   email,
			Message: "email must contain @",
		}
	}

	return nil
}

func GetUser(id int) (string, error) {
	// Симуляція бази даних
	users := map[int]string{
		1: "Alice",
		2: "Bob",
		3: "Charlie",
	}

	if name, ok := users[id]; ok {
		return name, nil
	}

	return "", NotFoundError{Resource: "User", ID: id}
}

func Authenticate(username, password string) error {
	if username == "" {
		return AuthError{Username: username, Reason: "username is empty"}
	}

	if password == "" {
		return AuthError{Username: username, Reason: "password is empty"}
	}

	if password == "wrong" {
		return AuthError{Username: username, Reason: "invalid credentials"}
	}

	return nil
}

func ConnectDB(table string) error {
	// Симуляція помилки БД
	baseErr := errors.New("connection timeout")
	return DBError{
		Operation: "connect",
		Table:     table,
		Err:       baseErr,
	}
}

// ============= Error wrapping =============

func ProcessData(data string) error {
	if data == "" {
		return fmt.Errorf("process failed: %w", errors.New("empty data"))
	}
	return nil
}

func SaveData(data string) error {
	if err := ProcessData(data); err != nil {
		return fmt.Errorf("save failed: %w", err)
	}
	return nil
}

// ============= Multiple errors =============

type MultiError struct {
	Errors []error
}

func (m MultiError) Error() string {
	if len(m.Errors) == 0 {
		return "no errors"
	}

	msg := fmt.Sprintf("multiple errors (%d):", len(m.Errors))
	for i, err := range m.Errors {
		msg += fmt.Sprintf("\n  %d. %v", i+1, err)
	}
	return msg
}

func ValidateUser(username, email, password string) error {
	var errs []error

	if username == "" {
		errs = append(errs, ValidationError{
			Field:   "username",
			Value:   username,
			Message: "username is required",
		})
	}

	if err := ValidateEmail(email); err != nil {
		errs = append(errs, err)
	}

	if len(password) < 8 {
		errs = append(errs, ValidationError{
			Field:   "password",
			Value:   "***",
			Message: "password must be at least 8 characters",
		})
	}

	if len(errs) > 0 {
		return MultiError{Errors: errs}
	}

	return nil
}

// ============= Main =============

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║         error Interface                  ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	// ===== ValidationError =====
	fmt.Println("\n🔹 ValidationError")
	fmt.Println("─────────────────────────────────────────")

	testEmails := []string{"", "abc", "test", "test@example.com"}
	for _, email := range testEmails {
		err := ValidateEmail(email)
		if err != nil {
			fmt.Printf("❌ %v\n", err)
		} else {
			fmt.Printf("✅ Email '%s' is valid\n", email)
		}
	}

	// ===== NotFoundError =====
	fmt.Println("\n🔹 NotFoundError")
	fmt.Println("─────────────────────────────────────────")

	userIDs := []int{1, 2, 5, 10}
	for _, id := range userIDs {
		name, err := GetUser(id)
		if err != nil {
			fmt.Printf("❌ %v\n", err)
		} else {
			fmt.Printf("✅ User %d: %s\n", id, name)
		}
	}

	// ===== Type assertion на помилку =====
	fmt.Println("\n🔹 Type assertion")
	fmt.Println("─────────────────────────────────────────")

	_, err := GetUser(999)
	if notFoundErr, ok := err.(NotFoundError); ok {
		fmt.Printf("NotFoundError detected!\n")
		fmt.Printf("  Resource: %s\n", notFoundErr.Resource)
		fmt.Printf("  ID: %d\n", notFoundErr.ID)
	}

	// ===== AuthError =====
	fmt.Println("\n🔹 AuthError")
	fmt.Println("─────────────────────────────────────────")

	authTests := []struct {
		user string
		pass string
	}{
		{"john", "secret123"},
		{"", "password"},
		{"alice", ""},
		{"bob", "wrong"},
	}

	for _, test := range authTests {
		err := Authenticate(test.user, test.pass)
		if err != nil {
			fmt.Printf("❌ %v\n", err)
		} else {
			fmt.Printf("✅ User '%s' authenticated\n", test.user)
		}
	}

	// ===== DBError with Unwrap =====
	fmt.Println("\n🔹 DBError (with wrapping)")
	fmt.Println("─────────────────────────────────────────")

	err = ConnectDB("users")
	fmt.Printf("Error: %v\n", err)

	// Unwrap базову помилку
	if dbErr, ok := err.(DBError); ok {
		fmt.Printf("Operation: %s\n", dbErr.Operation)
		fmt.Printf("Table: %s\n", dbErr.Table)
		fmt.Printf("Base error: %v\n", dbErr.Unwrap())
	}

	// ===== NetworkError =====
	fmt.Println("\n🔹 NetworkError")
	fmt.Println("─────────────────────────────────────────")

	netErr := NetworkError{
		Host:      "api.example.com",
		Port:      443,
		Operation: "GET",
		Timestamp: time.Now(),
	}
	fmt.Printf("%v\n", netErr)

	// ===== APIError =====
	fmt.Println("\n🔹 APIError")
	fmt.Println("─────────────────────────────────────────")

	apiErrors := []APIError{
		{Code: 400, Message: "Bad Request"},
		{Code: 404, Message: "Not Found"},
		{Code: 500, Message: "Internal Server Error"},
	}

	for _, err := range apiErrors {
		fmt.Printf("%v", err)
		if err.IsClientError() {
			fmt.Print(" (client error)")
		}
		if err.IsServerError() {
			fmt.Print(" (server error)")
		}
		fmt.Println()
	}

	// ===== Error wrapping з fmt.Errorf =====
	fmt.Println("\n🔹 Error wrapping")
	fmt.Println("─────────────────────────────────────────")

	err = SaveData("")
	fmt.Printf("Error: %v\n", err)

	// errors.Is - перевірка обгорнутої помилки
	if errors.Is(err, errors.New("empty data")) {
		fmt.Println("Contains 'empty data' error")
	}

	// ===== MultiError =====
	fmt.Println("\n🔹 MultiError (кілька помилок)")
	fmt.Println("─────────────────────────────────────────")

	err = ValidateUser("", "bad-email", "123")
	if err != nil {
		fmt.Printf("%v\n", err)

		// Обробка кожної помилки окремо
		if multiErr, ok := err.(MultiError); ok {
			fmt.Printf("\nВсього помилок: %d\n", len(multiErr.Errors))
		}
	}

	// ===== Стандартні помилки =====
	fmt.Println("\n🔹 Стандартні помилки")
	fmt.Println("─────────────────────────────────────────")

	// errors.New
	simpleErr := errors.New("something went wrong")
	fmt.Printf("Simple error: %v\n", simpleErr)

	// fmt.Errorf
	formattedErr := fmt.Errorf("failed to process item %d", 42)
	fmt.Printf("Formatted error: %v\n", formattedErr)

	// ===== Порівняння помилок =====
	fmt.Println("\n🔹 Порівняння помилок")
	fmt.Println("─────────────────────────────────────────")

	err1 := errors.New("test error")
	err2 := errors.New("test error")
	err3 := err1

	fmt.Printf("err1 == err2: %t (різні екземпляри)\n", err1 == err2)
	fmt.Printf("err1 == err3: %t (той самий екземпляр)\n", err1 == err3)

	// Sentinel errors
	var ErrNotFound = errors.New("not found")
	var ErrUnauthorized = errors.New("unauthorized")

	testErr := ErrNotFound
	fmt.Printf("\ntestErr == ErrNotFound: %t\n", testErr == ErrNotFound)
	fmt.Printf("testErr == ErrUnauthorized: %t\n", testErr == ErrUnauthorized)

	// ===== Best practices =====
	fmt.Println("\n🔹 Best Practices")
	fmt.Println("─────────────────────────────────────────")

	// 1. Перевірка на nil
	err = ValidateEmail("test@example.com")
	if err != nil {
		fmt.Println("Error occurred")
	} else {
		fmt.Println("✅ No error")
	}

	// 2. Sentinel errors для відомих помилок
	var (
		ErrInvalidInput = errors.New("invalid input")
		ErrTimeout      = errors.New("timeout")
	)

	fmt.Printf("Sentinel errors: %v, %v\n", ErrInvalidInput, ErrTimeout)

	// 3. Контекст в помилках
	contextErr := fmt.Errorf("failed to save user: %w",
		errors.New("database connection lost"))
	fmt.Printf("With context: %v\n", contextErr)

	// ===== Summary =====
	fmt.Println("\n\n📝 ВИСНОВКИ")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("✅ error - це інтерфейс з одним методом Error()")
	fmt.Println()
	fmt.Println("💡 Створення помилок:")
	fmt.Println("   • errors.New(\"message\")")
	fmt.Println("   • fmt.Errorf(\"format\", args...)")
	fmt.Println("   • Власні типи помилок")
	fmt.Println()
	fmt.Println("🔗 Wrapping:")
	fmt.Println("   • fmt.Errorf(\"context: %w\", err)")
	fmt.Println("   • errors.Is() - перевірка")
	fmt.Println("   • errors.As() - type assertion")
	fmt.Println()
	fmt.Println("⚡ Best practices:")
	fmt.Println("   • Завжди перевіряти err != nil")
	fmt.Println("   • Додавати контекст до помилок")
	fmt.Println("   • Використовувати sentinel errors")
	fmt.Println("   • Створювати власні типи для складних випадків")
}
