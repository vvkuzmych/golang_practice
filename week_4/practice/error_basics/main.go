package main

import (
	"errors"
	"fmt"
)

// ============= Sentinel Errors =============

var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidInput = errors.New("invalid input")
)

// ============= Custom Error Type =============

type MyError struct {
	Code    int
	Message string
}

func (e MyError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// ============= Examples =============

func example1_SimpleError() {
	fmt.Println("1️⃣ Проста помилка з errors.New()")
	fmt.Println("─────────────────────────────────────────")

	err := errors.New("something went wrong")
	fmt.Printf("Error: %v\n", err)
	fmt.Printf("Type: %T\n\n", err)
}

func example2_SentinelErrors() {
	fmt.Println("2️⃣ Sentinel Errors")
	fmt.Println("─────────────────────────────────────────")

	err := findUser(999)
	if err != nil {
		// Порівняння з sentinel error
		if err == ErrNotFound {
			fmt.Println("✓ User not found (detected with ==)")
		}
	}
	fmt.Println()
}

func findUser(id int) error {
	return ErrNotFound
}

func example3_CustomErrorType() {
	fmt.Println("3️⃣ Custom Error Type")
	fmt.Println("─────────────────────────────────────────")

	err := validateInput("")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("Type: %T\n", err)

		// Type assertion для отримання деталей
		if myErr, ok := err.(MyError); ok {
			fmt.Printf("Code: %d\n", myErr.Code)
			fmt.Printf("Message: %s\n", myErr.Message)
		}
	}
	fmt.Println()
}

func validateInput(input string) error {
	if input == "" {
		return MyError{
			Code:    400,
			Message: "input cannot be empty",
		}
	}
	return nil
}

func example4_ErrorFormatting() {
	fmt.Println("4️⃣ Error Formatting з fmt.Errorf()")
	fmt.Println("─────────────────────────────────────────")

	userID := 123
	err := fmt.Errorf("failed to process user %d", userID)
	fmt.Printf("Error: %v\n\n", err)
}

func example5_NilError() {
	fmt.Println("5️⃣ Nil Error (Success)")
	fmt.Println("─────────────────────────────────────────")

	err := successfulOperation()
	if err != nil {
		fmt.Println("Error occurred")
	} else {
		fmt.Println("✓ Operation successful (err == nil)")
	}
	fmt.Println()
}

func successfulOperation() error {
	// Все добре
	return nil
}

func example6_ErrorChecking() {
	fmt.Println("6️⃣ Правильна перевірка помилок")
	fmt.Println("─────────────────────────────────────────")

	result, err := divide(10, 0)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("✓ Result: %d\n", result)
	fmt.Println()
}

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

func example7_MultipleErrors() {
	fmt.Println("7️⃣ Множинні помилки")
	fmt.Println("─────────────────────────────────────────")

	errs := validateUser("", "invalid-email", -5)

	fmt.Printf("Found %d errors:\n", len(errs))
	for i, err := range errs {
		fmt.Printf("  %d. %v\n", i+1, err)
	}
	fmt.Println()
}

func validateUser(username, email string, age int) []error {
	var errs []error

	if username == "" {
		errs = append(errs, errors.New("username is required"))
	}

	if email == "" || !contains(email, "@") {
		errs = append(errs, errors.New("invalid email"))
	}

	if age < 0 {
		errs = append(errs, errors.New("age must be positive"))
	}

	return errs
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && s != substr
}

// ============= Main =============

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║         Error Basics Examples            ║")
	fmt.Println("╚══════════════════════════════════════════╝")
	fmt.Println()

	example1_SimpleError()
	example2_SentinelErrors()
	example3_CustomErrorType()
	example4_ErrorFormatting()
	example5_NilError()
	example6_ErrorChecking()
	example7_MultipleErrors()

	fmt.Println("📝 Висновки:")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("✅ error - це interface з методом Error() string")
	fmt.Println("✅ errors.New() для простих помилок")
	fmt.Println("✅ Sentinel errors для порівняння (==)")
	fmt.Println("✅ Custom types для складних помилок")
	fmt.Println("✅ Завжди перевіряйте err != nil")
	fmt.Println("✅ Повертайте nil при успіху")
}
