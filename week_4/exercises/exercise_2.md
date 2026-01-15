# Вправа 2: Error Wrapping Chain

## Ціль
Створити багаторівневу архітектуру з правильним error wrapping на кожному рівні.

---

## Завдання

Створіть програму `user_service.go` з 3-рівневою архітектурою:

1. **Database Layer** - симуляція роботи з БД
2. **Repository Layer** - абстракція над БД
3. **Service Layer** - бізнес-логіка

Кожен рівень має:
- Додавати контекст до помилок через wrapping (%w)
- Мати власні sentinel errors
- Правильно передавати помилки вгору

---

## Вимоги

### Architecture

```
HTTP Handler
    ↓
Service Layer (бізнес-логіка)
    ↓
Repository Layer (дані)
    ↓
Database Layer (низький рівень)
```

### Sentinel Errors

```go
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
    ErrInvalidUser = errors.New("invalid user")
    ErrUnauthorized = errors.New("unauthorized")
)
```

### User Struct

```go
type User struct {
    ID       int
    Username string
    Email    string
    IsActive bool
}
```

### Layers

#### 1. Database Layer

```go
type Database struct {
    connected bool
}

func (db *Database) Query(query string) (map[string]interface{}, error)
func (db *Database) Execute(query string) error
```

#### 2. Repository Layer

```go
type UserRepository struct {
    db *Database
}

func (r *UserRepository) FindByID(id int) (*User, error)
func (r *UserRepository) Create(user User) error
func (r *UserRepository) Delete(id int) error
```

#### 3. Service Layer

```go
type UserService struct {
    repo *UserRepository
}

func (s *UserService) GetUser(id int) (*User, error)
func (s *UserService) CreateUser(user User) error
func (s *UserService) DeleteUser(id int) error
```

---

## Приклад використання

```go
func main() {
    // Setup
    db := &Database{connected: true}
    repo := &UserRepository{db: db}
    service := &UserService{repo: repo}
    
    // Спроба отримати неіснуючого користувача
    user, err := service.GetUser(999)
    if err != nil {
        fmt.Println("Error:", err)
        // Output: "service: failed to get user 999: repository: user not found: record not found"
        
        // Перевірка оригінальної помилки
        if errors.Is(err, ErrNotFound) {
            fmt.Println("✓ Original ErrNotFound detected through chain")
        }
        
        // Перевірка repository помилки
        if errors.Is(err, ErrUserNotFound) {
            fmt.Println("✓ Repository ErrUserNotFound detected")
        }
    }
    
    // Створення користувача
    newUser := User{
        Username: "alice",
        Email:    "alice@example.com",
    }
    
    err = service.CreateUser(newUser)
    if err != nil {
        fmt.Printf("Create failed: %v\n", err)
    } else {
        fmt.Println("✓ User created successfully")
    }
}
```

---

## Очікуваний вивід

```
╔══════════════════════════════════════════╗
║    Multi-Layer Error Wrapping Demo      ║
╚══════════════════════════════════════════╝

🔹 Scenario 1: User Not Found
─────────────────────────────────────────
Attempting to get user ID: 999

❌ Error occurred:
   service: failed to get user 999: repository: user query failed: database: query execution failed: record not found

Error chain analysis:
   ✓ ErrNotFound detected (database level)
   ✓ ErrUserNotFound detected (repository level)
   
Error traversal:
   Level 4: record not found
   Level 3: database: query execution failed: record not found
   Level 2: repository: user query failed: database: query execution failed: record not found
   Level 1: service: failed to get user 999: repository: user query failed: database: query execution failed: record not found


🔹 Scenario 2: Database Connection Error
─────────────────────────────────────────
Simulating connection failure...

❌ Error occurred:
   service: failed to create user: repository: failed to insert user: database: connection error

Error chain analysis:
   ✓ ErrConnection detected


🔹 Scenario 3: Successful Operation
─────────────────────────────────────────
Creating user: alice

✓ User created successfully
User ID: 1


🔹 Scenario 4: Duplicate Key Error
─────────────────────────────────────────
Attempting to create duplicate user...

❌ Error occurred:
   service: failed to create user: repository: user already exists: duplicate key

Error chain analysis:
   ✓ ErrDuplicateKey detected


🔹 Error Wrapping Best Practices Demonstrated:
─────────────────────────────────────────
✓ Each layer adds meaningful context
✓ Original errors preserved through %w
✓ errors.Is() works across all layers
✓ Error chain shows complete execution path
✓ Sentinel errors enable specific error handling
```

---

## Підказки

### 1. Database Layer Wrapping

```go
func (db *Database) Query(query string) (map[string]interface{}, error) {
    if !db.connected {
        return nil, fmt.Errorf("database: connection failed: %w", ErrConnection)
    }
    
    // Симуляція: користувач не знайдений
    if query == "SELECT * FROM users WHERE id = 999" {
        return nil, fmt.Errorf("database: query execution failed: %w", ErrNotFound)
    }
    
    return map[string]interface{}{"id": 1, "username": "alice"}, nil
}
```

### 2. Repository Layer Wrapping

```go
func (r *UserRepository) FindByID(id int) (*User, error) {
    query := fmt.Sprintf("SELECT * FROM users WHERE id = %d", id)
    result, err := r.db.Query(query)
    if err != nil {
        if errors.Is(err, ErrNotFound) {
            return nil, fmt.Errorf("repository: user query failed: %w", err)
        }
        return nil, fmt.Errorf("repository: database error: %w", err)
    }
    
    // Parse result...
    user := &User{ID: id}
    return user, nil
}
```

### 3. Service Layer Wrapping

```go
func (s *UserService) GetUser(id int) (*User, error) {
    user, err := s.repo.FindByID(id)
    if err != nil {
        return nil, fmt.Errorf("service: failed to get user %d: %w", id, err)
    }
    
    if !user.IsActive {
        return nil, fmt.Errorf("service: user %d is inactive: %w", id, ErrUnauthorized)
    }
    
    return user, nil
}
```

### 4. Error Chain Visualization

```go
func printErrorChain(err error) {
    fmt.Println("Error chain:")
    level := 1
    for err != nil {
        fmt.Printf("  Level %d: %v\n", level, err)
        err = errors.Unwrap(err)
        level++
    }
}
```

---

## Бонус завдання

### 1. Custom Error Types per Layer

```go
type DatabaseError struct {
    Query string
    Err   error
}

type RepositoryError struct {
    Operation string
    Entity    string
    Err       error
}

type ServiceError struct {
    Action string
    UserID int
    Err    error
}
```

### 2. Error Metrics

```go
type ErrorMetrics struct {
    TotalErrors      int
    ErrorsByType     map[string]int
    ErrorsByLayer    map[string]int
}

func (m *ErrorMetrics) RecordError(err error)
```

### 3. Retry Logic

```go
func (s *UserService) GetUserWithRetry(id int, maxRetries int) (*User, error) {
    for attempt := 1; attempt <= maxRetries; attempt++ {
        user, err := s.GetUser(id)
        if err == nil {
            return user, nil
        }
        
        // Retry only on temporary errors
        if errors.Is(err, ErrTimeout) {
            time.Sleep(time.Second * time.Duration(attempt))
            continue
        }
        
        return nil, err
    }
    return nil, errors.New("max retries exceeded")
}
```

### 4. Structured Logging

```go
func logError(err error) {
    var dbErr DatabaseError
    if errors.As(err, &dbErr) {
        log.Printf("DB_ERROR query=%s err=%v", dbErr.Query, dbErr.Err)
        return
    }
    
    var repoErr RepositoryError
    if errors.As(err, &repoErr) {
        log.Printf("REPO_ERROR op=%s entity=%s err=%v", 
            repoErr.Operation, repoErr.Entity, repoErr.Err)
        return
    }
    
    log.Printf("ERROR: %v", err)
}
```

---

## Критерії оцінки

- ✅ Кожен рівень додає контекст до помилки
- ✅ Використовується %w для wrapping
- ✅ errors.Is() працює через всі рівні
- ✅ Оригінальна помилка доступна в кінці ланцюжка
- ✅ Sentinel errors визначені для кожного рівня
- ✅ Код організований по layers
- ✅ Error messages зрозумілі і корисні

---

## Рішення

Рішення знаходиться в `solutions/solution_2.go`.

Спробуйте виконати завдання самостійно перед тим, як дивитись рішення!

---

## Навчальні цілі

Після виконання цієї вправи ви будете вміти:
- Створювати багаторівневу архітектуру
- Правильно wrapping помилок на кожному рівні
- Додавати корисний контекст до помилок
- Використовувати errors.Is() для перевірки в ланцюжку
- Проектувати sentinel errors для кожного layer
- Дебажити error chains

---

## Подальше вдосконалення

Подумайте як додати:
- Transaction support в Repository
- Caching layer між Service та Repository
- Event publishing при помилках
- Circuit breaker pattern
- Error translation для різних клієнтів (HTTP, gRPC, CLI)
- Distributed tracing для error chains
