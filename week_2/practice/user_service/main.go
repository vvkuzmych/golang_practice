package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// ============= Domain Model =============

type User struct {
	ID        int
	Username  string
	Email     string
	CreatedAt time.Time
	IsActive  bool
}

func (u User) String() string {
	status := "активний"
	if !u.IsActive {
		status = "неактивний"
	}
	return fmt.Sprintf("[%d] %s <%s> - %s", u.ID, u.Username, u.Email, status)
}

// ============= UserService Interface =============

// UserService визначає операції над користувачами
type UserService interface {
	Create(username, email string) (*User, error)
	GetByID(id int) (*User, error)
	GetAll() ([]*User, error)
	Update(id int, username, email string) error
	Delete(id int) error
	Activate(id int) error
	Deactivate(id int) error
}

// ============= In-Memory Implementation (Real) =============

type InMemoryUserService struct {
	users  map[int]*User
	nextID int
}

func NewInMemoryUserService() *InMemoryUserService {
	return &InMemoryUserService{
		users:  make(map[int]*User),
		nextID: 1,
	}
}

func (s *InMemoryUserService) Create(username, email string) (*User, error) {
	// Валідація
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	if !strings.Contains(email, "@") {
		return nil, errors.New("invalid email format")
	}

	// Перевірка на дублікати
	for _, user := range s.users {
		if user.Username == username {
			return nil, errors.New("username already exists")
		}
		if user.Email == email {
			return nil, errors.New("email already exists")
		}
	}

	user := &User{
		ID:        s.nextID,
		Username:  username,
		Email:     email,
		CreatedAt: time.Now(),
		IsActive:  true,
	}

	s.users[s.nextID] = user
	s.nextID++

	return user, nil
}

func (s *InMemoryUserService) GetByID(id int) (*User, error) {
	user, exists := s.users[id]
	if !exists {
		return nil, fmt.Errorf("user with id %d not found", id)
	}
	return user, nil
}

func (s *InMemoryUserService) GetAll() ([]*User, error) {
	users := make([]*User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users, nil
}

func (s *InMemoryUserService) Update(id int, username, email string) error {
	user, exists := s.users[id]
	if !exists {
		return fmt.Errorf("user with id %d not found", id)
	}

	if username != "" {
		user.Username = username
	}
	if email != "" && strings.Contains(email, "@") {
		user.Email = email
	}

	return nil
}

func (s *InMemoryUserService) Delete(id int) error {
	if _, exists := s.users[id]; !exists {
		return fmt.Errorf("user with id %d not found", id)
	}
	delete(s.users, id)
	return nil
}

func (s *InMemoryUserService) Activate(id int) error {
	user, exists := s.users[id]
	if !exists {
		return fmt.Errorf("user with id %d not found", id)
	}
	user.IsActive = true
	return nil
}

func (s *InMemoryUserService) Deactivate(id int) error {
	user, exists := s.users[id]
	if !exists {
		return fmt.Errorf("user with id %d not found", id)
	}
	user.IsActive = false
	return nil
}

// ============= Mock Implementation (For Testing) =============

type MockUserService struct {
	users        []*User
	createCalled int
	getCalled    int
	getAllCalled int
	updateCalled int
	deleteCalled int
	shouldFail   bool
}

func NewMockUserService() *MockUserService {
	return &MockUserService{
		users: []*User{},
	}
}

func (m *MockUserService) SetShouldFail(fail bool) {
	m.shouldFail = fail
}

func (m *MockUserService) Create(username, email string) (*User, error) {
	m.createCalled++

	if m.shouldFail {
		return nil, errors.New("mock: create failed")
	}

	user := &User{
		ID:        len(m.users) + 1,
		Username:  username,
		Email:     email,
		CreatedAt: time.Now(),
		IsActive:  true,
	}
	m.users = append(m.users, user)
	return user, nil
}

func (m *MockUserService) GetByID(id int) (*User, error) {
	m.getCalled++

	if m.shouldFail {
		return nil, errors.New("mock: get failed")
	}

	return &User{
		ID:        id,
		Username:  "mock_user",
		Email:     "mock@example.com",
		CreatedAt: time.Now(),
		IsActive:  true,
	}, nil
}

func (m *MockUserService) GetAll() ([]*User, error) {
	m.getAllCalled++

	if m.shouldFail {
		return nil, errors.New("mock: getall failed")
	}

	return m.users, nil
}

func (m *MockUserService) Update(id int, username, email string) error {
	m.updateCalled++

	if m.shouldFail {
		return errors.New("mock: update failed")
	}

	return nil
}

func (m *MockUserService) Delete(id int) error {
	m.deleteCalled++

	if m.shouldFail {
		return errors.New("mock: delete failed")
	}

	return nil
}

func (m *MockUserService) Activate(id int) error {
	if m.shouldFail {
		return errors.New("mock: activate failed")
	}
	return nil
}

func (m *MockUserService) Deactivate(id int) error {
	if m.shouldFail {
		return errors.New("mock: deactivate failed")
	}
	return nil
}

func (m *MockUserService) Stats() string {
	return fmt.Sprintf("Mock Stats: Create=%d, Get=%d, GetAll=%d, Update=%d, Delete=%d",
		m.createCalled, m.getCalled, m.getAllCalled, m.updateCalled, m.deleteCalled)
}

// ============= Application Layer =============

type Application struct {
	userService UserService // залежність від інтерфейсу!
}

func NewApplication(service UserService) *Application {
	return &Application{
		userService: service,
	}
}

func (app *Application) RegisterUser(username, email string) error {
	user, err := app.userService.Create(username, email)
	if err != nil {
		return fmt.Errorf("registration failed: %w", err)
	}
	fmt.Printf("✅ User registered: %s\n", user)
	return nil
}

func (app *Application) ShowUser(id int) error {
	user, err := app.userService.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	fmt.Printf("👤 %s\n", user)
	return nil
}

func (app *Application) ListUsers() error {
	users, err := app.userService.GetAll()
	if err != nil {
		return fmt.Errorf("failed to list users: %w", err)
	}

	fmt.Println("👥 All Users:")
	if len(users) == 0 {
		fmt.Println("   (no users)")
		return nil
	}

	for _, user := range users {
		fmt.Printf("   • %s\n", user)
	}
	return nil
}

func (app *Application) RemoveUser(id int) error {
	err := app.userService.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	fmt.Printf("🗑️  User %d deleted\n", id)
	return nil
}

func (app *Application) ToggleUserStatus(id int, active bool) error {
	var err error
	if active {
		err = app.userService.Activate(id)
	} else {
		err = app.userService.Deactivate(id)
	}

	if err != nil {
		return fmt.Errorf("failed to toggle status: %w", err)
	}

	status := "активовано"
	if !active {
		status = "деактивовано"
	}
	fmt.Printf("🔄 User %d %s\n", id, status)
	return nil
}

// ============= Main =============

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║  UserService: Interface Demo            ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	// ===== In-Memory Implementation =====
	fmt.Println("\n📦 IN-MEMORY IMPLEMENTATION")
	fmt.Println("─────────────────────────────────────────")

	inMemoryService := NewInMemoryUserService()
	app1 := NewApplication(inMemoryService)

	fmt.Println("\n🔹 Реєстрація користувачів:")
	app1.RegisterUser("ivan_petro", "ivan@example.com")
	app1.RegisterUser("maria_ivanova", "maria@example.com")
	app1.RegisterUser("petro_sydarenko", "petro@example.com")

	fmt.Println("\n🔹 Список користувачів:")
	app1.ListUsers()

	fmt.Println("\n🔹 Отримання користувача #2:")
	app1.ShowUser(2)

	fmt.Println("\n🔹 Деактивація користувача #2:")
	app1.ToggleUserStatus(2, false)
	app1.ShowUser(2)

	fmt.Println("\n🔹 Видалення користувача #1:")
	app1.RemoveUser(1)
	app1.ListUsers()

	// ===== Mock Implementation =====
	fmt.Println("\n\n🎭 MOCK IMPLEMENTATION (For Testing)")
	fmt.Println("─────────────────────────────────────────")

	mockService := NewMockUserService()
	app2 := NewApplication(mockService)

	fmt.Println("\n🔹 Тестування з Mock:")
	app2.RegisterUser("test_user1", "test1@example.com")
	app2.RegisterUser("test_user2", "test2@example.com")
	app2.ListUsers()

	fmt.Println("\n🔹 Mock Statistics:")
	fmt.Println("   " + mockService.Stats())

	// ===== Error Handling with Mock =====
	fmt.Println("\n\n❌ ERROR HANDLING (Mock with failures)")
	fmt.Println("─────────────────────────────────────────")

	mockService.SetShouldFail(true)

	fmt.Println("\n🔹 Спроба створення (має провалитись):")
	err := app2.RegisterUser("fail_user", "fail@example.com")
	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	}

	fmt.Println("\n🔹 Спроба отримання (має провалитись):")
	err = app2.ShowUser(999)
	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	}

	// ===== Comparison =====
	fmt.Println("\n\n⚖️  ПОРІВНЯННЯ РЕАЛІЗАЦІЙ")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("In-Memory (Real):")
	fmt.Println("  ✅ Реальне зберігання даних")
	fmt.Println("  ✅ Валідація та бізнес-логіка")
	fmt.Println("  ✅ Використовується в продакшені")
	fmt.Println()
	fmt.Println("Mock:")
	fmt.Println("  ✅ Швидке тестування")
	fmt.Println("  ✅ Контроль над поведінкою")
	fmt.Println("  ✅ Статистика викликів")
	fmt.Println("  ✅ Симуляція помилок")

	// ===== Summary =====
	fmt.Println("\n\n📝 ВИСНОВКИ")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("✅ Interface дозволяє:")
	fmt.Println("   • Змінювати реалізацію без зміни коду")
	fmt.Println("   • Тестувати через Mock")
	fmt.Println("   • Dependency Injection")
	fmt.Println("   • Loose coupling")
	fmt.Println()
	fmt.Println("💡 Неявна реалізація:")
	fmt.Println("   • InMemoryUserService реалізує UserService")
	fmt.Println("   • MockUserService реалізує UserService")
	fmt.Println("   • Application працює з обома через інтерфейс!")
	fmt.Println()
	fmt.Println("🎯 Переваги:")
	fmt.Println("   • Гнучка архітектура")
	fmt.Println("   • Легко тестувати")
	fmt.Println("   • Легко додати нові реалізації")
	fmt.Println("   • Не залежимо від конкретної реалізації")
}
