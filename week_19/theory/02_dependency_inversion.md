# Dependency Inversion & Dependency Injection

## Dependency Inversion Principle (DIP)

**High-level modules should not depend on low-level modules. Both should depend on abstractions.**

### Без DIP (погана практика)

```go
// Low-level module
type MySQLDatabase struct {
    conn *sql.DB
}

func (db *MySQLDatabase) SaveUser(name string) error {
    _, err := db.conn.Exec("INSERT INTO users (name) VALUES (?)", name)
    return err
}

// High-level module залежить від low-level
type UserService struct {
    db *MySQLDatabase // Пряма залежність!
}

func (s *UserService) CreateUser(name string) error {
    return s.db.SaveUser(name)
}

// Проблеми:
// 1. Неможливо змінити БД без зміни UserService
// 2. Складно тестувати (потрібна реальна MySQL)
// 3. Tight coupling
```

### З DIP (хороша практика)

```go
// Abstraction (інтерфейс)
type UserRepository interface {
    SaveUser(name string) error
}

// Low-level module імплементує абстракцію
type MySQLUserRepository struct {
    conn *sql.DB
}

func (db *MySQLUserRepository) SaveUser(name string) error {
    _, err := db.conn.Exec("INSERT INTO users (name) VALUES (?)", name)
    return err
}

// High-level module залежить від абстракції
type UserService struct {
    repo UserRepository // Залежність від інтерфейсу
}

func (s *UserService) CreateUser(name string) error {
    return s.repo.SaveUser(name)
}

// Переваги:
// 1. Легко змінити БД (передати іншу реалізацію)
// 2. Легко тестувати (передати mock)
// 3. Loose coupling
```

---

## Dependency Injection (DI)

**Техніка передачі залежностей ззовні, а не створення всередині.**

### 3 види DI

#### 1. Constructor Injection (найпопулярніший в Go)

```go
type UserService struct {
    repo   UserRepository
    logger *log.Logger
}

// Залежності передаються через конструктор
func NewUserService(repo UserRepository, logger *log.Logger) *UserService {
    return &UserService{
        repo:   repo,
        logger: logger,
    }
}

func main() {
    db, _ := sql.Open("mysql", "...")
    repo := &MySQLUserRepository{conn: db}
    logger := log.New(os.Stdout, "USER: ", log.LstdFlags)
    
    service := NewUserService(repo, logger)
    service.CreateUser("Alice")
}
```

#### 2. Method Injection

```go
type UserService struct{}

// Залежність передається в кожен метод
func (s *UserService) CreateUser(repo UserRepository, name string) error {
    return repo.SaveUser(name)
}

func main() {
    service := &UserService{}
    repo := &MySQLUserRepository{}
    
    service.CreateUser(repo, "Alice")
}
```

#### 3. Interface Injection (рідко в Go)

```go
type ServiceConfigurator interface {
    SetRepository(repo UserRepository)
}

type UserService struct {
    repo UserRepository
}

func (s *UserService) SetRepository(repo UserRepository) {
    s.repo = repo
}
```

---

## DI Container (для великих проектів)

### Без DI Container

```go
func main() {
    // Вручну створюємо всі залежності
    db, _ := sql.Open("mysql", "...")
    repo := NewUserRepository(db)
    logger := log.New(os.Stdout, "", log.LstdFlags)
    cache := NewRedisCache("localhost:6379")
    emailSender := NewEmailSender("smtp.gmail.com")
    
    userService := NewUserService(repo, logger, cache)
    orderService := NewOrderService(repo, logger, emailSender)
    paymentService := NewPaymentService(logger, cache, emailSender)
    
    // ...багато boilerplate коду
}
```

### З DI Container (wire, dig, fx)

```go
// Using google/wire
// wire.go
//go:build wireinject

func InitializeUserService() (*UserService, error) {
    wire.Build(
        NewDB,
        NewUserRepository,
        NewLogger,
        NewUserService,
    )
    return &UserService{}, nil
}

// Generated code автоматично створює залежності
func main() {
    service, err := InitializeUserService()
    if err != nil {
        log.Fatal(err)
    }
    
    service.CreateUser("Alice")
}
```

---

## Практичний приклад: Refactoring з DIP

### Before (без DIP)

```go
type OrderService struct{}

func (s *OrderService) PlaceOrder(userID int, productID int) error {
    // Пряма залежність від MySQL
    db, _ := sql.Open("mysql", "...")
    _, err := db.Exec("INSERT INTO orders...", userID, productID)
    if err != nil {
        return err
    }
    
    // Пряма залежність від SMTP
    smtp.SendMail("smtp.gmail.com:587", nil, "from@example.com", 
        []string{"to@example.com"}, []byte("Order placed"))
    
    // Пряма залежність від Stripe
    stripe.Key = "sk_test_..."
    _, err = charge.New(&stripe.ChargeParams{Amount: 1000})
    
    return err
}
```

### After (з DIP)

```go
// Abstractions
type OrderRepository interface {
    SaveOrder(userID, productID int) error
}

type EmailSender interface {
    SendOrderConfirmation(userID int) error
}

type PaymentProcessor interface {
    Charge(amount int) error
}

// Service залежить від абстракцій
type OrderService struct {
    repo     OrderRepository
    email    EmailSender
    payment  PaymentProcessor
}

func NewOrderService(repo OrderRepository, email EmailSender, payment PaymentProcessor) *OrderService {
    return &OrderService{
        repo:    repo,
        email:   email,
        payment: payment,
    }
}

func (s *OrderService) PlaceOrder(userID int, productID int) error {
    if err := s.repo.SaveOrder(userID, productID); err != nil {
        return err
    }
    
    if err := s.payment.Charge(1000); err != nil {
        return err
    }
    
    return s.email.SendOrderConfirmation(userID)
}

// Implementations
type MySQLOrderRepository struct {
    db *sql.DB
}

func (r *MySQLOrderRepository) SaveOrder(userID, productID int) error {
    _, err := r.db.Exec("INSERT INTO orders...", userID, productID)
    return err
}

type SMTPEmailSender struct {
    host string
}

func (e *SMTPEmailSender) SendOrderConfirmation(userID int) error {
    return smtp.SendMail(e.host, nil, "from@example.com", 
        []string{"to@example.com"}, []byte("Order placed"))
}

type StripePaymentProcessor struct {
    apiKey string
}

func (p *StripePaymentProcessor) Charge(amount int) error {
    stripe.Key = p.apiKey
    _, err := charge.New(&stripe.ChargeParams{Amount: amount})
    return err
}

// Usage
func main() {
    db, _ := sql.Open("mysql", "...")
    repo := &MySQLOrderRepository{db: db}
    email := &SMTPEmailSender{host: "smtp.gmail.com:587"}
    payment := &StripePaymentProcessor{apiKey: "sk_test_..."}
    
    service := NewOrderService(repo, email, payment)
    service.PlaceOrder(1, 100)
}
```

---

## Тестування з DIP

### Mock Implementation

```go
// Mock для тестування
type MockOrderRepository struct {
    SaveOrderFunc func(userID, productID int) error
}

func (m *MockOrderRepository) SaveOrder(userID, productID int) error {
    if m.SaveOrderFunc != nil {
        return m.SaveOrderFunc(userID, productID)
    }
    return nil
}

// Test
func TestPlaceOrder(t *testing.T) {
    // Arrange
    mockRepo := &MockOrderRepository{
        SaveOrderFunc: func(userID, productID int) error {
            assert.Equal(t, 1, userID)
            assert.Equal(t, 100, productID)
            return nil
        },
    }
    mockEmail := &MockEmailSender{}
    mockPayment := &MockPaymentProcessor{}
    
    service := NewOrderService(mockRepo, mockEmail, mockPayment)
    
    // Act
    err := service.PlaceOrder(1, 100)
    
    // Assert
    assert.NoError(t, err)
}
```

---

## Коли використовувати DIP?

### ✅ Використовуй DIP коли:
- Потрібна гнучкість (зміна БД, external services)
- Потрібне тестування (моки, стаби)
- Складна бізнес-логіка (багато залежностей)
- Мікросервісна архітектура

### ❌ Не переусложнюй коли:
- Прості CRUD операції
- Скрипти, CLI tools
- Прототипи, MVP
- Залежність ніколи не зміниться

---

## Підсумок

| Концепція | Опис |
|-----------|------|
| **DIP** | Принцип: залежати від абстракцій |
| **DI** | Техніка: передавати залежності ззовні |
| **IoC** | Контейнер: керує створенням залежностей |

**Переваги:**
- ✅ Легко тестувати (моки)
- ✅ Легко замінювати реалізації
- ✅ Loose coupling
- ✅ Гнучка архітектура

**Недоліки:**
- ❌ Більше коду (інтерфейси, конструктори)
- ❌ Складніша структура проекту
- ❌ Може бути overkill для простих проектів
