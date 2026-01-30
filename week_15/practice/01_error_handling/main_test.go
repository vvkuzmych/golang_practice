package main

import (
	"errors"
	"testing"
)

// ========================================
// Test GetUser
// ========================================

func TestGetUser(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		wantUser  string
		wantErr   error
		checkType bool
	}{
		{
			name:     "valid user",
			id:       1,
			wantUser: "User-1",
			wantErr:  nil,
		},
		{
			name:     "not found",
			id:       0,
			wantUser: "",
			wantErr:  ErrNotFound,
		},
		{
			name:     "invalid input",
			id:       -1,
			wantUser: "",
			wantErr:  ErrInvalidInput,
		},
		{
			name:     "unauthorized",
			id:       999,
			wantUser: "",
			wantErr:  ErrUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := GetUser(tt.id)

			// Check error
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			// Check success
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if user != tt.wantUser {
				t.Errorf("Expected user %s, got %s", tt.wantUser, user)
			}
		})
	}
}

// ========================================
// Test ValidateEmail
// ========================================

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		wantField string
		wantMsg   string
		shouldErr bool
	}{
		{
			name:      "valid email",
			email:     "test@example.com",
			shouldErr: false,
		},
		{
			name:      "empty email",
			email:     "",
			wantField: "email",
			wantMsg:   "cannot be empty",
			shouldErr: true,
		},
		{
			name:      "missing @",
			email:     "invalid",
			wantField: "email",
			wantMsg:   "must contain @",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)

			if !tt.shouldErr {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				return
			}

			// Check error exists
			if err == nil {
				t.Errorf("Expected error, got nil")
				return
			}

			// Use errors.As to extract ValidationError
			var valErr *ValidationError
			if !errors.As(err, &valErr) {
				t.Errorf("Expected ValidationError, got %T", err)
				return
			}

			// Check fields
			if valErr.Field != tt.wantField {
				t.Errorf("Expected field %s, got %s", tt.wantField, valErr.Field)
			}
			if valErr.Message != tt.wantMsg {
				t.Errorf("Expected message %s, got %s", tt.wantMsg, valErr.Message)
			}
		})
	}
}

// ========================================
// Test SaveUser
// ========================================

func TestSaveUser(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		email   string
		wantErr error
	}{
		{
			name:    "success",
			id:      5,
			email:   "test@example.com",
			wantErr: nil,
		},
		{
			name:    "invalid email",
			id:      1,
			email:   "",
			wantErr: nil, // Will check with errors.As
		},
		{
			name:    "user not found",
			id:      0,
			email:   "test@example.com",
			wantErr: ErrNotFound,
		},
		{
			name:    "database error",
			id:      777,
			email:   "test@example.com",
			wantErr: nil, // Will check with errors.As
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SaveUser(tt.id, tt.email)

			// Test success case
			if tt.name == "success" {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				return
			}

			// Test error cases
			if err == nil {
				t.Errorf("Expected error, got nil")
				return
			}

			// Special checks for different error types
			switch tt.name {
			case "invalid email":
				var valErr *ValidationError
				if !errors.As(err, &valErr) {
					t.Errorf("Expected ValidationError, got %T: %v", err, err)
				}

			case "user not found":
				if !errors.Is(err, ErrNotFound) {
					t.Errorf("Expected ErrNotFound, got %v", err)
				}

			case "database error":
				var dbErr *DatabaseError
				if !errors.As(err, &dbErr) {
					t.Errorf("Expected DatabaseError, got %T: %v", err, err)
				}
			}
		})
	}
}

// ========================================
// Test Error Wrapping
// ========================================

func TestErrorWrapping(t *testing.T) {
	t.Run("wrapped error preserves original", func(t *testing.T) {
		err := SaveUser(0, "test@example.com")

		// Check that original error is preserved
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("Expected wrapped error to contain ErrNotFound")
		}

		// Check that error message contains context
		errMsg := err.Error()
		expectedContains := "user 0"
		if !contains(errMsg, expectedContains) {
			t.Errorf("Expected error message to contain '%s', got: %s",
				expectedContains, errMsg)
		}
	})

	t.Run("unwrap returns original error", func(t *testing.T) {
		originalErr := errors.New("test error")
		dbErr := &DatabaseError{
			Query: "SELECT *",
			Err:   originalErr,
		}

		unwrapped := errors.Unwrap(dbErr)
		if unwrapped != originalErr {
			t.Errorf("Expected unwrapped error to be original error")
		}
	})
}

// ========================================
// Test errors.Is vs ==
// ========================================

func TestErrorsIsVsEqual(t *testing.T) {
	t.Run("== fails with wrapped errors", func(t *testing.T) {
		err := SaveUser(0, "test@example.com")

		// ❌ Using == fails
		if err == ErrNotFound {
			t.Errorf("== should not work with wrapped errors, but it did!")
		}
	})

	t.Run("errors.Is works with wrapped errors", func(t *testing.T) {
		err := SaveUser(0, "test@example.com")

		// ✅ Using errors.Is works
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("errors.Is should work with wrapped errors")
		}
	})

	t.Run("== works with direct errors", func(t *testing.T) {
		_, err := GetUser(0)

		// ✅ Both work for direct errors
		if err != ErrNotFound {
			t.Errorf("== should work with direct errors")
		}
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("errors.Is should also work with direct errors")
		}
	})
}

// ========================================
// Test errors.As extraction
// ========================================

func TestErrorsAs(t *testing.T) {
	t.Run("extract ValidationError", func(t *testing.T) {
		err := ValidateEmail("")

		var valErr *ValidationError
		if !errors.As(err, &valErr) {
			t.Fatalf("errors.As should extract ValidationError")
		}

		if valErr.Field != "email" {
			t.Errorf("Expected field 'email', got '%s'", valErr.Field)
		}
	})

	t.Run("extract wrapped ValidationError", func(t *testing.T) {
		err := SaveUser(1, "")

		var valErr *ValidationError
		if !errors.As(err, &valErr) {
			t.Fatalf("errors.As should extract wrapped ValidationError")
		}

		if valErr.Field != "email" {
			t.Errorf("Expected field 'email', got '%s'", valErr.Field)
		}
	})

	t.Run("fail to extract wrong type", func(t *testing.T) {
		err := ValidateEmail("")

		var dbErr *DatabaseError
		if errors.As(err, &dbErr) {
			t.Errorf("errors.As should not extract DatabaseError from ValidationError")
		}
	})
}

// ========================================
// Benchmark
// ========================================

func BenchmarkErrorsIs(b *testing.B) {
	err := SaveUser(0, "test@example.com")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = errors.Is(err, ErrNotFound)
	}
}

func BenchmarkErrorsAs(b *testing.B) {
	err := ValidateEmail("")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var valErr *ValidationError
		_ = errors.As(err, &valErr)
	}
}
