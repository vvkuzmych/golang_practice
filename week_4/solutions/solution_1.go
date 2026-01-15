package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// ============= Sentinel Errors =============

var (
	ErrRequired      = errors.New("field is required")
	ErrInvalidFormat = errors.New("invalid format")
	ErrOutOfRange    = errors.New("value out of range")
)

// ============= ValidationError =============

type ValidationError struct {
	Field string
	Value interface{}
	Rule  string
	Err   error
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

// ============= ValidationUser =============

// ValidationUser представляє користувача для валідації
type ValidationUser struct {
	Username string
	Email    string
	Age      int
	Password string
}

// ============= Validation Functions =============

func ValidateUsername(username string) error {
	if username == "" {
		return NewValidationError("username", username, "cannot be empty", ErrRequired)
	}

	if len(username) < 3 {
		return NewValidationError("username", username, "must be at least 3 characters", ErrInvalidFormat)
	}

	if len(username) > 20 {
		return NewValidationError("username", username, "must be at most 20 characters", ErrInvalidFormat)
	}

	// Перевірка: тільки літери, цифри, '_'
	for _, r := range username {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return NewValidationError("username", username,
				"must contain only letters, digits, and underscore", ErrInvalidFormat)
		}
	}

	return nil
}

func ValidateEmail(email string) error {
	if email == "" {
		return NewValidationError("email", email, "cannot be empty", ErrRequired)
	}

	if !strings.Contains(email, "@") {
		return NewValidationError("email", email, "must contain '@'", ErrInvalidFormat)
	}

	if !strings.Contains(email, ".") {
		return NewValidationError("email", email, "must contain '.'", ErrInvalidFormat)
	}

	return nil
}

func ValidateAge(age int) error {
	if age <= 0 {
		return NewValidationError("age", age, "must be greater than 0", ErrOutOfRange)
	}

	if age >= 150 {
		return NewValidationError("age", age, "must be less than 150", ErrOutOfRange)
	}

	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return NewValidationError("password", password, "cannot be empty", ErrRequired)
	}

	if len(password) < 8 {
		return NewValidationError("password", password, "must be at least 8 characters", ErrInvalidFormat)
	}

	// Перевірка: містить хоча б одну цифру
	hasDigit := false
	for _, r := range password {
		if unicode.IsDigit(r) {
			hasDigit = true
			break
		}
	}

	if !hasDigit {
		return NewValidationError("password", password, "must contain at least one digit", ErrInvalidFormat)
	}

	return nil
}

func ValidateUser(user ValidationUser) []error {
	var errs []error

	if err := ValidateUsername(user.Username); err != nil {
		errs = append(errs, err)
	}

	if err := ValidateEmail(user.Email); err != nil {
		errs = append(errs, err)
	}

	if err := ValidateAge(user.Age); err != nil {
		errs = append(errs, err)
	}

	if err := ValidatePassword(user.Password); err != nil {
		errs = append(errs, err)
	}

	return errs
}

// ============= Main =============

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║        User Validation System            ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	// ===== Valid User =====
	fmt.Println("\n🔹 Valid User")
	fmt.Println("─────────────────────────────────────────")

	validUser := ValidationUser{
		Username: "alice_123",
		Email:    "alice@example.com",
		Age:      25,
		Password: "securePass123",
	}

	printUser(validUser)

	errs := ValidateUser(validUser)
	if len(errs) == 0 {
		fmt.Println("\n✓ All validations passed!")
	} else {
		fmt.Println("\n❌ Validation failed:")
		printErrors(errs)
	}

	// ===== Invalid User =====
	fmt.Println("\n\n🔹 Invalid User")
	fmt.Println("─────────────────────────────────────────")

	invalidUser := ValidationUser{
		Username: "ab",
		Email:    "invalid",
		Age:      200,
		Password: "short",
	}

	printUser(invalidUser)

	errs = ValidateUser(invalidUser)
	if len(errs) > 0 {
		fmt.Println("\n❌ Validation errors:\n")
		printErrors(errs)
	}

	// ===== Error Type Detection =====
	fmt.Println("\n\n🔹 Error Type Detection")
	fmt.Println("─────────────────────────────────────────")

	testUser := ValidationUser{
		Username: "ab",
		Email:    "test@mail",
		Age:      200,
		Password: "",
	}

	errs = ValidateUser(testUser)

	// Перевірка errors.Is()
	foundOutOfRange := false
	foundRequired := false

	for _, err := range errs {
		if errors.Is(err, ErrOutOfRange) {
			foundOutOfRange = true
		}
		if errors.Is(err, ErrRequired) {
			foundRequired = true
		}
	}

	if foundOutOfRange {
		fmt.Println("Detecting ErrOutOfRange: ✓ Found")
	}
	if foundRequired {
		fmt.Println("Detecting ErrRequired: ✓ Found")
	}

	// Перевірка errors.As()
	var valErr ValidationError
	if errors.As(errs[0], &valErr) {
		fmt.Println("Extracting ValidationError details: ✓ Success")
		fmt.Printf("  Field: %s\n", valErr.Field)
		fmt.Printf("  Value: %v\n", valErr.Value)
		fmt.Printf("  Rule: %s\n", valErr.Rule)
	}
}

func printUser(user ValidationUser) {
	fmt.Printf("Username: %s\n", user.Username)
	fmt.Printf("Email: %s\n", user.Email)
	fmt.Printf("Age: %d\n", user.Age)
	fmt.Printf("Password: %s\n", maskPassword(user.Password))
}

func printErrors(errs []error) {
	for i, err := range errs {
		fmt.Printf("%d. %v\n", i+1, err)

		var valErr ValidationError
		if errors.As(err, &valErr) {
			fmt.Printf("   Field: %s\n", valErr.Field)
			fmt.Printf("   Value: %v\n", valErr.Value)
			fmt.Printf("   Rule: %s\n", valErr.Rule)
		}
		fmt.Println()
	}
}

func maskPassword(password string) string {
	if len(password) == 0 {
		return "(empty)"
	}
	return strings.Repeat("*", len(password))
}
