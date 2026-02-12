# Clean Architecture

## Концепція

**Clean Architecture** (Robert C. Martin) - архітектурний патерн з чіткими layers та правилом залежностей.

```
┌─────────────────────────────────────────┐
│         External Interfaces              │  UI, Web, CLI, External APIs
│  (Frameworks, Drivers, Web, DB, UI)      │
├─────────────────────────────────────────┤
│      Interface Adapters                  │  Controllers, Presenters,
│  (Controllers, Gateways, Presenters)     │  Gateways
├─────────────────────────────────────────┤
│         Use Cases                        │  Application Business Rules
│  (Application Business Rules)            │
├─────────────────────────────────────────┤
│         Entities                         │  Enterprise Business Rules
│  (Enterprise Business Rules, Domain)     │  (Domain Models)
└─────────────────────────────────────────┘

Dependency Rule: →→→→ (внутрішні layers не знають про зовнішні)
```

---

## 4 Layers

### 1. Domain Layer (Entities)
**Найвнутрішній layer. Бізнес-логіка підприємства.**

```go
// domain/user.go
package domain

import "errors"

// Entity (Domain Model)
type User struct {
    ID       int
    Name     string
    Email    string
    Password string
}

// Business Rules
func (u *User) Validate() error {
    if u.Email == "" {
        return errors.New("email is required")
    }
    if len(u.Password) < 8 {
        return errors.New("password must be at least 8 characters")
    }
    return nil
}

func (u *User) IsActive() bool {
    return u.ID > 0
}
```

### 2. Use Case Layer (Application Business Rules)
**Орк

естрація бізнес-логіки застосунку.**

```go
// usecase/user_usecase.go
package usecase

import "domain"

// Port (Interface for repository)
type UserRepository interface {
    Save(user *domain.User) error
    FindByID(id int) (*domain.User, error)
    FindByEmail(email string) (*domain.User, error)
}

// Port (Interface for password hasher)
type PasswordHasher interface {
    Hash(password string) (string, error)
    Compare(hashed, plain string) error
}

// Use Case
type UserUseCase struct {
    repo   UserRepository
    hasher PasswordHasher
}

func NewUserUseCase(repo UserRepository, hasher PasswordHasher) *UserUseCase {
    return &UserUseCase{
        repo:   repo,
        hasher: hasher,
    }
}

func (uc *UserUseCase) RegisterUser(name, email, password string) error {
    // Перевірка існуючого користувача
    existing, _ := uc.repo.FindByEmail(email)
    if existing != nil {
        return errors.New("user already exists")
    }
    
    // Створення domain entity
    user := &domain.User{
        Name:  name,
        Email: email,
    }
    
    // Валідація domain rules
    if err := user.Validate(); err != nil {
        return err
    }
    
    // Хешування пароля
    hashedPassword, err := uc.hasher.Hash(password)
    if err != nil {
        return err
    }
    user.Password = hashedPassword
    
    // Збереження
    return uc.repo.Save(user)
}
```

### 3. Interface Adapters (Controllers, Presenters, Gateways)
**Адаптація даних між use cases та зовнішнім світом.**

```go
// adapter/http/user_controller.go
package http

import (
    "encoding/json"
    "net/http"
    "usecase"
)

type UserController struct {
    useCase *usecase.UserUseCase
}

func NewUserController(uc *usecase.UserUseCase) *UserController {
    return &UserController{useCase: uc}
}

// HTTP Handler
func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Name     string `json:"name"`
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Виклик use case
    if err := c.useCase.RegisterUser(req.Name, req.Email, req.Password); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
```

### 4. Infrastructure (Frameworks, Drivers)
**Зовнішній layer. БД, Web frameworks, зовнішні API.**

```go
// infrastructure/postgres/user_repository.go
package postgres

import (
    "database/sql"
    "domain"
)

type PostgresUserRepository struct {
    db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
    return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Save(user *domain.User) error {
    query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`
    return r.db.QueryRow(query, user.Name, user.Email, user.Password).Scan(&user.ID)
}

func (r *PostgresUserRepository) FindByID(id int) (*domain.User, error) {
    var user domain.User
    query := `SELECT id, name, email, password FROM users WHERE id = $1`
    err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
    return &user, err
}

func (r *PostgresUserRepository) FindByEmail(email string) (*domain.User, error) {
    var user domain.User
    query := `SELECT id, name, email, password FROM users WHERE email = $1`
    err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    return &user, err
}
```

---

## Dependency Rule

**Залежності завжди йдуть ВНУТРІШНЬО:**

```
Infrastructure → Interface Adapters → Use Cases → Domain
```

**НЕ ДОЗВОЛЕНО:**
- Domain не може знати про Use Cases
- Use Cases не можуть знати про Controllers
- Controllers не можуть знати про Infrastructure напряму

---

## Структура проекту

```
myapp/
├── cmd/
│   └── api/
│       └── main.go              # Entry point
├── internal/
│   ├── domain/                  # Layer 1: Entities
│   │   ├── user.go
│   │   └── order.go
│   ├── usecase/                 # Layer 2: Use Cases
│   │   ├── user_usecase.go
│   │   └── order_usecase.go
│   ├── adapter/                 # Layer 3: Interface Adapters
│   │   ├── http/
│   │   │   ├── user_controller.go
│   │   │   └── order_controller.go
│   │   └── presenter/
│   │       └── user_presenter.go
│   └── infrastructure/          # Layer 4: Infrastructure
│       ├── postgres/
│       │   ├── user_repository.go
│       │   └── order_repository.go
│       ├── bcrypt/
│       │   └── password_hasher.go
│       └── smtp/
│           └── email_sender.go
└── go.mod
```

---

## Приклад: Повний flow

```go
// 1. main.go (Infrastructure setup)
package main

import (
    "database/sql"
    "log"
    "net/http"
    
    "myapp/internal/adapter/http"
    "myapp/internal/infrastructure/bcrypt"
    "myapp/internal/infrastructure/postgres"
    "myapp/internal/usecase"
)

func main() {
    // Infrastructure
    db, _ := sql.Open("postgres", "...")
    defer db.Close()
    
    // Repositories (Infrastructure)
    userRepo := postgres.NewPostgresUserRepository(db)
    
    // Services (Infrastructure)
    hasher := bcrypt.NewBCryptHasher()
    
    // Use Cases (Application)
    userUseCase := usecase.NewUserUseCase(userRepo, hasher)
    
    // Controllers (Adapter)
    userController := adapter.NewUserController(userUseCase)
    
    // HTTP Routes
    http.HandleFunc("/register", userController.Register)
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## Переваги Clean Architecture

✅ **Testability**: Кожен layer можна тестувати окремо  
✅ **Independence**: Frameworks, DB, UI можна легко замінити  
✅ **Maintainability**: Зміни локалізовані  
✅ **Flexibility**: Легко додавати нові features  
✅ **Business Logic First**: Domain rules в центрі

---

## Недоліки

❌ Більше boilerplate коду  
❌ Крива навчання  
❌ Overkill для простих CRUD  
❌ Більше abstractions (інтерфейси)

---

## Коли використовувати?

### ✅ Використовуй коли:
- Складна бізнес-логіка
- Довгострокові проекти
- Потрібна гнучкість (зміна БД, UI)
- Команда > 3 розробників

### ❌ Не використовуй коли:
- Простий CRUD
- Prototype, MVP
- Solo проект, скрипт
- Deadline < 1 місяць
