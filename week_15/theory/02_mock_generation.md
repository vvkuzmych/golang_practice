# Mock Generation in Go

## 🎯 Що таке Mocks?

**Mock** - це тестовий об'єкт, який імітує поведінку real dependencies.

### Навіщо потрібні?

✅ **Ізоляція** - тестуємо лише свій код  
✅ **Швидкість** - без DB, HTTP, file I/O  
✅ **Контроль** - можемо задавати будь-які responses  
✅ **Edge cases** - легко симулювати помилки  

---

## 📦 gomock + mockgen

**gomock** - бібліотека для mock objects  
**mockgen** - генератор mocks з interfaces

### Встановлення

```bash
go install github.com/golang/mock/mockgen@latest
go get github.com/golang/mock/gomock
```

---

## 🎯 Basic Example

### Step 1: Define Interface

```go
// user_service.go
package service

type User struct {
    ID    int
    Name  string
    Email string
}

type UserRepository interface {
    GetByID(id int) (*User, error)
    Save(user *User) error
    Delete(id int) error
}
```

### Step 2: Generate Mock

```bash
# From source file
mockgen -source=user_service.go -destination=mocks/user_repository_mock.go -package=mocks

# From package + interface
mockgen -destination=mocks/user_repository_mock.go \
        -package=mocks \
        github.com/yourname/project/service UserRepository
```

### Step 3: Write Test

```go
// user_service_test.go
package service_test

import (
    "errors"
    "testing"
    
    "github.com/golang/mock/gomock"
    "github.com/yourname/project/mocks"
    "github.com/yourname/project/service"
)

func TestGetUser(t *testing.T) {
    // 1. Create mock controller
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    // 2. Create mock
    mockRepo := mocks.NewMockUserRepository(ctrl)
    
    // 3. Set expectations
    mockRepo.EXPECT().
        GetByID(1).
        Return(&service.User{ID: 1, Name: "John"}, nil)
    
    // 4. Use mock in your code
    svc := service.NewUserService(mockRepo)
    user, err := svc.GetUser(1)
    
    // 5. Assertions
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if user.Name != "John" {
        t.Errorf("expected John, got %s", user.Name)
    }
}
```

---

## 🔧 mockgen Options

### Generate from Source

```bash
mockgen -source=interface.go \
        -destination=mocks/mock.go \
        -package=mocks
```

### Generate from Package

```bash
mockgen -destination=mocks/mock.go \
        -package=mocks \
        github.com/yourname/project/pkg InterfaceName
```

### Generate Multiple Interfaces

```bash
mockgen -destination=mocks/mock.go \
        -package=mocks \
        github.com/yourname/project/pkg Interface1,Interface2
```

### Add to go:generate

```go
//go:generate mockgen -source=user_service.go -destination=mocks/user_repository_mock.go -package=mocks
```

Then run:
```bash
go generate ./...
```

---

## 🎯 Setting Expectations

### Basic Expectation

```go
mock.EXPECT().
    MethodName(arg1, arg2).
    Return(result1, result2)
```

### Multiple Calls

```go
// Exactly once (default)
mock.EXPECT().GetByID(1).Return(user, nil)

// Exactly N times
mock.EXPECT().GetByID(1).Return(user, nil).Times(3)

// At least N times
mock.EXPECT().GetByID(1).Return(user, nil).MinTimes(2)

// At most N times
mock.EXPECT().GetByID(1).Return(user, nil).MaxTimes(5)

// Any number of times (including 0)
mock.EXPECT().GetByID(1).Return(user, nil).AnyTimes()
```

### Call Order

```go
gomock.InOrder(
    mock.EXPECT().GetByID(1).Return(user1, nil),
    mock.EXPECT().GetByID(2).Return(user2, nil),
)
```

---

## 🎯 Argument Matchers

### Exact Match

```go
mock.EXPECT().GetByID(1).Return(user, nil)
```

### Any

```go
mock.EXPECT().GetByID(gomock.Any()).Return(user, nil)
```

### Eq (Explicit equality)

```go
mock.EXPECT().GetByID(gomock.Eq(1)).Return(user, nil)
```

### Not

```go
mock.EXPECT().GetByID(gomock.Not(0)).Return(user, nil)
```

### Nil

```go
mock.EXPECT().Save(gomock.Nil()).Return(errors.New("nil user"))
```

### AssignableToTypeOf

```go
mock.EXPECT().
    Save(gomock.AssignableToTypeOf(&User{})).
    Return(nil)
```

### Custom Matcher

```go
// Define custom matcher
type userWithEmail string

func (u userWithEmail) Matches(x interface{}) bool {
    user, ok := x.(*User)
    return ok && user.Email == string(u)
}

func (u userWithEmail) String() string {
    return fmt.Sprintf("user with email %s", string(u))
}

// Use
mock.EXPECT().
    Save(userWithEmail("test@example.com")).
    Return(nil)
```

---

## 🎯 Return Values

### Simple Return

```go
mock.EXPECT().GetByID(1).Return(user, nil)
```

### Multiple Returns

```go
mock.EXPECT().GetByID(1).Return(user1, nil).Times(1)
mock.EXPECT().GetByID(2).Return(user2, nil).Times(1)
```

### Do() для Side Effects

```go
var savedUser *User

mock.EXPECT().
    Save(gomock.Any()).
    Do(func(u *User) {
        savedUser = u
        u.ID = 123  // Simulate DB auto-increment
    }).
    Return(nil)
```

### DoAndReturn()

```go
mock.EXPECT().
    GetByID(gomock.Any()).
    DoAndReturn(func(id int) (*User, error) {
        if id < 0 {
            return nil, errors.New("invalid id")
        }
        return &User{ID: id, Name: "Test"}, nil
    })
```

---

## 📊 Real-World Example

### Service Code

```go
// user_service.go
package service

import "errors"

type UserService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) GetUser(id int) (*User, error) {
    if id <= 0 {
        return nil, errors.New("invalid id")
    }
    return s.repo.GetByID(id)
}

func (s *UserService) CreateUser(name, email string) (*User, error) {
    if name == "" || email == "" {
        return nil, errors.New("name and email required")
    }
    
    user := &User{Name: name, Email: email}
    if err := s.repo.Save(user); err != nil {
        return nil, err
    }
    
    return user, nil
}
```

### Test Code

```go
// user_service_test.go
package service_test

import (
    "errors"
    "testing"
    
    "github.com/golang/mock/gomock"
    "github.com/yourname/project/mocks"
    "github.com/yourname/project/service"
)

func TestUserService_GetUser(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := mocks.NewMockUserRepository(ctrl)
    svc := service.NewUserService(mockRepo)
    
    t.Run("success", func(t *testing.T) {
        expectedUser := &service.User{ID: 1, Name: "John"}
        
        mockRepo.EXPECT().
            GetByID(1).
            Return(expectedUser, nil)
        
        user, err := svc.GetUser(1)
        
        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
        if user.Name != "John" {
            t.Errorf("expected John, got %s", user.Name)
        }
    })
    
    t.Run("invalid id", func(t *testing.T) {
        // No mock call expected (validation fails first)
        
        _, err := svc.GetUser(-1)
        
        if err == nil {
            t.Fatal("expected error for invalid id")
        }
    })
    
    t.Run("not found", func(t *testing.T) {
        mockRepo.EXPECT().
            GetByID(999).
            Return(nil, errors.New("not found"))
        
        _, err := svc.GetUser(999)
        
        if err == nil {
            t.Fatal("expected error")
        }
    })
}

func TestUserService_CreateUser(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := mocks.NewMockUserRepository(ctrl)
    svc := service.NewUserService(mockRepo)
    
    t.Run("success", func(t *testing.T) {
        mockRepo.EXPECT().
            Save(gomock.Any()).
            DoAndReturn(func(u *service.User) error {
                u.ID = 123  // Simulate DB
                return nil
            })
        
        user, err := svc.CreateUser("John", "john@example.com")
        
        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
        if user.ID != 123 {
            t.Errorf("expected ID 123, got %d", user.ID)
        }
    })
    
    t.Run("validation error", func(t *testing.T) {
        // No mock call expected
        
        _, err := svc.CreateUser("", "")
        
        if err == nil {
            t.Fatal("expected validation error")
        }
    })
    
    t.Run("repository error", func(t *testing.T) {
        mockRepo.EXPECT().
            Save(gomock.Any()).
            Return(errors.New("db error"))
        
        _, err := svc.CreateUser("John", "john@example.com")
        
        if err == nil {
            t.Fatal("expected error")
        }
    })
}
```

---

## ✅ Best Practices

### 1. One mock per test

```go
// ❌ BAD - shared mock
var mockRepo *mocks.MockUserRepository

func TestA(t *testing.T) {
    // Uses shared mock
}

// ✅ GOOD - new mock per test
func TestA(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockRepo := mocks.NewMockUserRepository(ctrl)
    // ...
}
```

### 2. Always defer ctrl.Finish()

```go
ctrl := gomock.NewController(t)
defer ctrl.Finish()  // ✅ Verifies expectations
```

### 3. Use table-driven tests

```go
tests := []struct {
    name    string
    id      int
    mockFn  func(*mocks.MockUserRepository)
    wantErr bool
}{
    {
        name: "success",
        id:   1,
        mockFn: func(m *mocks.MockUserRepository) {
            m.EXPECT().GetByID(1).Return(user, nil)
        },
        wantErr: false,
    },
    // ...
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        ctrl := gomock.NewController(t)
        defer ctrl.Finish()
        
        mockRepo := mocks.NewMockUserRepository(ctrl)
        tt.mockFn(mockRepo)
        
        // Test logic
    })
}
```

### 4. Don't over-mock

```go
// ❌ BAD - mocking simple types
type StringValidator interface {
    IsEmpty(s string) bool
}

// ✅ GOOD - use real implementation
func IsEmpty(s string) bool {
    return s == ""
}
```

### 5. Test behavior, not implementation

```go
// ❌ BAD - testing internal calls
mock.EXPECT().InternalMethod().Times(1)

// ✅ GOOD - testing public API
result := svc.PublicMethod()
// Assert result
```

---

## 🎯 Common Patterns

### Pattern 1: Factory Functions

```go
func newMockUserService(t *testing.T) (*service.UserService, *mocks.MockUserRepository) {
    ctrl := gomock.NewController(t)
    t.Cleanup(ctrl.Finish)
    
    mockRepo := mocks.NewMockUserRepository(ctrl)
    svc := service.NewUserService(mockRepo)
    
    return svc, mockRepo
}

func TestSomething(t *testing.T) {
    svc, mockRepo := newMockUserService(t)
    
    mockRepo.EXPECT().GetByID(1).Return(user, nil)
    // Test logic
}
```

### Pattern 2: Subtests

```go
func TestUserService(t *testing.T) {
    t.Run("GetUser", func(t *testing.T) {
        t.Run("success", func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()
            // ...
        })
        
        t.Run("not found", func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()
            // ...
        })
    })
}
```

---

## 🎓 Висновок

### Коли використовувати mocks:

✅ **External dependencies** - DB, HTTP, file I/O  
✅ **Slow operations** - network calls, heavy computation  
✅ **Non-deterministic** - time, random, external APIs  
✅ **Hard to reproduce** - edge cases, errors  

### Коли НЕ використовувати:

❌ **Simple logic** - pure functions, calculations  
❌ **Data structures** - structs, primitives  
❌ **Standard library** - usually don't mock  

### Golden Rules:

1. **Mock interfaces, not concrete types**
2. **One mock per test**
3. **Always defer ctrl.Finish()**
4. **Test behavior, not implementation**
5. **Keep mocks simple**

**Week 15: Mock Generation Master!** 🎭✅
