package main

import (
	"errors"
	"fmt"
	"time"
)

// ============= Sentinel Errors =============

// Database errors
var (
	ErrConnection = errors.New("database connection error")
	ErrTimeout    = errors.New("database timeout")
	ErrNotFound   = errors.New("record not found")
)

// Repository errors
var (
	ErrUserNotFound = errors.New("user not found")
	ErrDuplicateKey = errors.New("duplicate key")
)

// Service errors
var (
	ErrInvalidUser  = errors.New("invalid user")
	ErrUnauthorized = errors.New("unauthorized")
)

// ============= DBUser (Database User) =============

type DBUser struct {
	ID       int
	Username string
	Email    string
	IsActive bool
}

// ============= MockDatabase Layer =============

type MockDatabase struct {
	connected bool
}

func NewMockDatabase() *MockDatabase {
	return &MockDatabase{connected: true}
}

func (db *MockDatabase) Query(query string) (map[string]interface{}, error) {
	// Симуляція: не підключено
	if !db.connected {
		return nil, fmt.Errorf("database: connection failed: %w", ErrConnection)
	}

	// Симуляція: користувач не знайдений
	if query == "SELECT * FROM users WHERE id = 999" {
		return nil, fmt.Errorf("database: query execution failed: %w", ErrNotFound)
	}

	// Симуляція: дублікат ключа
	if query == "INSERT INTO users (username) VALUES ('alice')" {
		return nil, fmt.Errorf("database: insert failed: %w", ErrDuplicateKey)
	}

	// Успішний запит
	return map[string]interface{}{
		"id":       1,
		"username": "alice",
		"email":    "alice@example.com",
		"active":   true,
	}, nil
}

func (db *MockDatabase) Execute(query string) error {
	if !db.connected {
		return fmt.Errorf("database: execution failed: %w", ErrConnection)
	}
	return nil
}

// ============= Repository Layer =============

type UserRepository struct {
	db *MockDatabase
}

func NewUserRepository(db *MockDatabase) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(id int) (*DBUser, error) {
	query := fmt.Sprintf("SELECT * FROM users WHERE id = %d", id)
	result, err := r.db.Query(query)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, fmt.Errorf("repository: user query failed: %w", err)
		}
		return nil, fmt.Errorf("repository: database error: %w", err)
	}

	// Парсинг результату
	user := &DBUser{
		ID:       result["id"].(int),
		Username: result["username"].(string),
		Email:    result["email"].(string),
		IsActive: result["active"].(bool),
	}

	return user, nil
}

func (r *UserRepository) Create(user DBUser) error {
	query := fmt.Sprintf("INSERT INTO users (username) VALUES ('%s')", user.Username)
	result, err := r.db.Query(query)
	if err != nil {
		if errors.Is(err, ErrDuplicateKey) {
			return fmt.Errorf("repository: user already exists: %w", err)
		}
		return fmt.Errorf("repository: failed to insert user: %w", err)
	}

	// Встановлюємо ID
	if id, ok := result["id"].(int); ok {
		user.ID = id
	}

	return nil
}

func (r *UserRepository) Delete(id int) error {
	query := fmt.Sprintf("DELETE FROM users WHERE id = %d", id)
	err := r.db.Execute(query)
	if err != nil {
		return fmt.Errorf("repository: failed to delete user: %w", err)
	}
	return nil
}

// ============= Service Layer (Error Wrapping Demo) =============

type ErrorWrappingService struct {
	repo *UserRepository
}

func NewErrorWrappingService(repo *UserRepository) *ErrorWrappingService {
	return &ErrorWrappingService{repo: repo}
}

func (s *ErrorWrappingService) GetUser(id int) (*DBUser, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get user %d: %w", id, err)
	}

	if !user.IsActive {
		return nil, fmt.Errorf("service: user %d is inactive: %w", id, ErrUnauthorized)
	}

	return user, nil
}

func (s *ErrorWrappingService) CreateUser(user DBUser) error {
	if user.Username == "" {
		return fmt.Errorf("service: invalid user data: %w", ErrInvalidUser)
	}

	err := s.repo.Create(user)
	if err != nil {
		return fmt.Errorf("service: failed to create user: %w", err)
	}

	return nil
}

func (s *ErrorWrappingService) DeleteUser(id int) error {
	err := s.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("service: failed to delete user: %w", err)
	}
	return nil
}

// ============= Helper Functions =============

func printErrorChain(err error) {
	fmt.Println("\nError chain traversal:")
	level := 1
	for err != nil {
		fmt.Printf("  Level %d: %v\n", level, err)
		err = errors.Unwrap(err)
		level++
	}
}

func analyzeError(err error) {
	fmt.Println("\nError chain analysis:")

	if errors.Is(err, ErrNotFound) {
		fmt.Println("   ✓ ErrNotFound detected (database level)")
	}

	if errors.Is(err, ErrUserNotFound) {
		fmt.Println("   ✓ ErrUserNotFound detected (repository level)")
	}

	if errors.Is(err, ErrConnection) {
		fmt.Println("   ✓ ErrConnection detected")
	}

	if errors.Is(err, ErrDuplicateKey) {
		fmt.Println("   ✓ ErrDuplicateKey detected")
	}
}

// ============= Main =============

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║    Multi-Layer Error Wrapping Demo      ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	// Setup
	db := NewMockDatabase()
	repo := NewUserRepository(db)
	service := NewErrorWrappingService(repo)

	// ===== Scenario 1: User Not Found =====
	fmt.Println("\n🔹 Scenario 1: User Not Found")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("Attempting to get user ID: 999")

	_, err := service.GetUser(999)
	if err != nil {
		fmt.Printf("\n❌ Error occurred:\n   %v\n", err)
		analyzeError(err)
		printErrorChain(err)
	}

	// ===== Scenario 2: Database Connection Error =====
	fmt.Println("\n\n🔹 Scenario 2: Database Connection Error")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("Simulating connection failure...")

	db.connected = false

	newUser := DBUser{Username: "bob"}
	err = service.CreateUser(newUser)
	if err != nil {
		fmt.Printf("\n❌ Error occurred:\n   %v\n", err)
		analyzeError(err)
	}

	db.connected = true // Відновлюємо з'єднання

	// ===== Scenario 3: Successful Operation =====
	fmt.Println("\n\n🔹 Scenario 3: Successful Operation")
	fmt.Println("─────────────────────────────────────────")

	successUser := DBUser{
		Username: "charlie",
		Email:    "charlie@example.com",
	}

	fmt.Printf("Creating user: %s\n", successUser.Username)

	err = service.CreateUser(successUser)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Println("\n✓ User created successfully")
		fmt.Printf("User ID: %d\n", successUser.ID)
	}

	// ===== Scenario 4: Duplicate Key Error =====
	fmt.Println("\n\n🔹 Scenario 4: Duplicate Key Error")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("Attempting to create duplicate user...")

	duplicateUser := DBUser{Username: "alice"}
	err = service.CreateUser(duplicateUser)
	if err != nil {
		fmt.Printf("\n❌ Error occurred:\n   %v\n", err)
		analyzeError(err)
	}

	// ===== Summary =====
	fmt.Println("\n\n🔹 Error Wrapping Best Practices Demonstrated:")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("✓ Each layer adds meaningful context")
	fmt.Println("✓ Original errors preserved through %w")
	fmt.Println("✓ errors.Is() works across all layers")
	fmt.Println("✓ Error chain shows complete execution path")
	fmt.Println("✓ Sentinel errors enable specific error handling")

	// ===== Performance Demo =====
	fmt.Println("\n\n🔹 Multiple Operations Performance")
	fmt.Println("─────────────────────────────────────────")

	operations := []struct {
		name string
		fn   func() error
	}{
		{"GetUser(1)", func() error {
			_, err := service.GetUser(1)
			return err
		}},
		{"CreateUser", func() error {
			return service.CreateUser(DBUser{Username: "dave"})
		}},
		{"GetUser(999)", func() error {
			_, err := service.GetUser(999)
			return err
		}},
	}

	for i, op := range operations {
		start := time.Now()
		err := op.fn()
		duration := time.Since(start)

		fmt.Printf("%d. %s - ", i+1, op.name)
		if err != nil {
			fmt.Printf("❌ failed (%v)\n", duration)
		} else {
			fmt.Printf("✓ success (%v)\n", duration)
		}
	}
}
