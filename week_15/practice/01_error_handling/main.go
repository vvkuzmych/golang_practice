package main

import (
	"errors"
	"fmt"
)

// ========================================
// Sentinel Errors
// ========================================

var (
	ErrNotFound      = errors.New("not found")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrInvalidInput  = errors.New("invalid input")
	ErrAlreadyExists = errors.New("already exists")
)

// ========================================
// Custom Error Types
// ========================================

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
}

type DatabaseError struct {
	Query string
	Err   error
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("database error: %s (query: %s)", e.Err, e.Query)
}

func (e *DatabaseError) Unwrap() error {
	return e.Err
}

// ========================================
// Functions that return errors
// ========================================

func GetUser(id int) (string, error) {
	if id < 0 {
		return "", ErrInvalidInput
	}
	if id == 0 {
		return "", ErrNotFound
	}
	if id == 999 {
		return "", ErrUnauthorized
	}
	return fmt.Sprintf("User-%d", id), nil
}

func ValidateEmail(email string) error {
	if email == "" {
		return &ValidationError{
			Field:   "email",
			Message: "cannot be empty",
		}
	}
	if !contains(email, "@") {
		return &ValidationError{
			Field:   "email",
			Message: "must contain @",
		}
	}
	return nil
}

func SaveUser(id int, email string) error {
	// Validate email first
	if err := ValidateEmail(email); err != nil {
		// Wrap error with context
		return fmt.Errorf("failed to save user %d: %w", id, err)
	}

	// Check if user doesn't exist first (to save it)
	_, err := GetUser(id)
	if err != nil {
		// Wrap the error (could be ErrNotFound, ErrInvalidInput, etc.)
		return fmt.Errorf("failed to verify user %d: %w", id, err)
	}

	// If GetUser succeeded, user already exists
	return fmt.Errorf("user %d: %w", id, ErrAlreadyExists)
}

// ========================================
// Helper
// ========================================

func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ========================================
// Examples
// ========================================

func example1_BasicErrors() {
	fmt.Println("1️⃣ Example 1: Basic Error Handling")
	fmt.Println("─────────────────────────────────────────")

	user, err := GetUser(1)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ User: %s\n", user)
	}

	user, err = GetUser(0)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	}
	fmt.Println()
}

func example2_SentinelErrors_OldWay() {
	fmt.Println("2️⃣ Example 2: Sentinel Errors (OLD WAY - ❌)")
	fmt.Println("─────────────────────────────────────────")

	_, err := GetUser(0)

	// ❌ OLD WAY: Using ==
	if err == ErrNotFound {
		fmt.Println("✅ Direct error: found with == (works)")
	}

	// This works for simple errors, but breaks with wrapping
	err = SaveUser(0, "test@mail.com")
	fmt.Printf("Wrapped error: %v\n", err)

	if err == ErrNotFound {
		fmt.Println("This won't print (error is wrapped)")
	} else {
		fmt.Println("❌ Can't find ErrNotFound with == in wrapped error")
	}
	fmt.Println()
}

func example3_ErrorsIs() {
	fmt.Println("3️⃣ Example 3: errors.Is() (NEW WAY - ✅)")
	fmt.Println("─────────────────────────────────────────")

	err := SaveUser(0, "test@mail.com")
	fmt.Printf("Wrapped error: %v\n", err)

	// ✅ NEW WAY: Using errors.Is()
	if errors.Is(err, ErrNotFound) {
		fmt.Println("✅ Found ErrNotFound in wrapped error (using errors.Is)")
	}

	// Check different sentinel errors
	_, err = GetUser(999)
	if errors.Is(err, ErrUnauthorized) {
		fmt.Println("✅ Found ErrUnauthorized (direct error)")
	}

	_, err = GetUser(-1)
	if errors.Is(err, ErrInvalidInput) {
		fmt.Println("✅ Found ErrInvalidInput (direct error)")
	}
	fmt.Println()
}

func example4_ErrorsAs() {
	fmt.Println("4️⃣ Example 4: errors.As() for Custom Types")
	fmt.Println("─────────────────────────────────────────")

	// Test validation error
	err := ValidateEmail("")

	var valErr *ValidationError
	if errors.As(err, &valErr) {
		fmt.Printf("✅ Validation Error:\n")
		fmt.Printf("   Field: %s\n", valErr.Field)
		fmt.Printf("   Message: %s\n", valErr.Message)
	}

	// Test wrapped validation error
	err = SaveUser(1, "invalid-email")

	if errors.As(err, &valErr) {
		fmt.Printf("✅ Wrapped Validation Error:\n")
		fmt.Printf("   Field: %s\n", valErr.Field)
		fmt.Printf("   Message: %s\n", valErr.Message)
		fmt.Printf("   Full error: %v\n", err)
	}
	fmt.Println()
}

func example5_DatabaseError() {
	fmt.Println("5️⃣ Example 5: Wrapped Errors with Unwrap")
	fmt.Println("─────────────────────────────────────────")

	// Create wrapped error
	err := SaveUser(0, "test@mail.com")
	fmt.Printf("Full error: %v\n", err)

	// Unwrap to get original error
	unwrapped := errors.Unwrap(err)
	fmt.Printf("Unwrapped once: %v\n", unwrapped)

	// Check if it's the sentinel error
	if unwrapped == ErrNotFound {
		fmt.Println("✅ Unwrapped error is ErrNotFound")
	}
	fmt.Println()
}

func example6_ErrorWrapping() {
	fmt.Println("6️⃣ Example 6: Error Wrapping Chain")
	fmt.Println("─────────────────────────────────────────")

	// Create error chain
	err := SaveUser(-1, "test@mail.com")

	fmt.Printf("Full error: %v\n\n", err)

	// Check if it contains ErrInvalidInput
	if errors.Is(err, ErrInvalidInput) {
		fmt.Println("✅ Contains ErrInvalidInput (wrapped)")
	}

	// Unwrap step by step
	fmt.Println("\nUnwrapping chain:")
	step := 1
	current := err
	for current != nil {
		fmt.Printf("  %d. %v\n", step, current)
		current = errors.Unwrap(current)
		step++
	}
	fmt.Println()
}

func example7_MultipleChecks() {
	fmt.Println("7️⃣ Example 7: Multiple Error Checks")
	fmt.Println("─────────────────────────────────────────")

	testCases := []struct {
		id    int
		email string
		desc  string
	}{
		{1, "test@mail.com", "User exists"},
		{0, "test@mail.com", "User not found"},
		{-1, "test@mail.com", "Invalid ID"},
		{999, "test@mail.com", "Unauthorized"},
		{1, "", "Empty email"},
		{1, "invalid", "Invalid email format"},
	}

	for _, tc := range testCases {
		fmt.Printf("\n📋 Test: %s (id=%d, email=%s)\n", tc.desc, tc.id, tc.email)
		err := SaveUser(tc.id, tc.email)

		if err == nil {
			fmt.Println("  ✅ Success!")
			continue
		}

		// Check sentinel errors
		switch {
		case errors.Is(err, ErrInvalidInput):
			fmt.Println("  ❌ Contains: ErrInvalidInput")
		case errors.Is(err, ErrNotFound):
			fmt.Println("  ❌ Contains: ErrNotFound")
		case errors.Is(err, ErrUnauthorized):
			fmt.Println("  ❌ Contains: ErrUnauthorized")
		case errors.Is(err, ErrAlreadyExists):
			fmt.Println("  ❌ Contains: ErrAlreadyExists")
		}

		// Check custom types
		var valErr *ValidationError
		if errors.As(err, &valErr) {
			fmt.Printf("  ⚠️  ValidationError → Field: %s, Message: %s\n",
				valErr.Field, valErr.Message)
		}
	}
	fmt.Println()
}

func example8_RealWorldAPI() {
	fmt.Println("8️⃣ Example 8: Real-World API Example")
	fmt.Println("─────────────────────────────────────────")

	// Simulate API handler
	handleGetUser := func(id int) (int, string) {
		user, err := GetUser(id)
		if err != nil {
			// Map errors to HTTP status codes
			switch {
			case errors.Is(err, ErrNotFound):
				return 404, "User not found"
			case errors.Is(err, ErrUnauthorized):
				return 401, "Unauthorized"
			case errors.Is(err, ErrInvalidInput):
				return 400, "Invalid input"
			default:
				return 500, "Internal server error"
			}
		}
		return 200, user
	}

	testIDs := []int{1, 0, -1, 999}
	for _, id := range testIDs {
		status, body := handleGetUser(id)
		fmt.Printf("GET /users/%d → %d %s\n", id, status, body)
	}
	fmt.Println()
}

// ========================================
// Main
// ========================================

func main() {
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║   Go Error Handling Examples          ║")
	fmt.Println("╚════════════════════════════════════════╝")
	fmt.Println()

	example1_BasicErrors()
	example2_SentinelErrors_OldWay()
	example3_ErrorsIs()
	example4_ErrorsAs()
	example5_DatabaseError()
	example6_ErrorWrapping()
	example7_MultipleChecks()
	example8_RealWorldAPI()

	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║   All Examples Completed! ✅           ║")
	fmt.Println("╚════════════════════════════════════════╝")
}
