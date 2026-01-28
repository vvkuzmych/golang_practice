# Testing in Go

## 📖 Зміст

1. [Unit Testing](#unit-testing)
2. [Table-Driven Tests](#table-driven-tests)
3. [Mocking](#mocking)
4. [Integration Testing](#integration-testing)
5. [Test Coverage](#test-coverage)

---

## Unit Testing

### Basic Test

```go
// math.go
package math

func Add(a, b int) int {
    return a + b
}

// math_test.go
package math

import "testing"

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    expected := 5
    
    if result != expected {
        t.Errorf("Add(2, 3) = %d; want %d", result, expected)
    }
}
```

### Using testify

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    assert.Equal(t, 5, result, "should add correctly")
}

func TestDivide(t *testing.T) {
    result, err := Divide(10, 2)
    require.NoError(t, err, "should not error")
    assert.Equal(t, 5.0, result)
}
```

---

## Table-Driven Tests

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive", 2, 3, 5},
        {"negative", -1, -1, -2},
        {"zero", 0, 5, 5},
        {"large", 1000, 2000, 3000},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d", 
                    tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

### With testify

```go
func TestUserService(t *testing.T) {
    tests := []struct {
        name    string
        userID  int
        wantErr bool
    }{
        {"valid user", 1, false},
        {"invalid user", -1, true},
        {"not found", 999, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            user, err := service.GetUser(tt.userID)
            
            if tt.wantErr {
                assert.Error(t, err)
                assert.Nil(t, user)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, user)
            }
        })
    }
}
```

---

## Mocking

### Manual Mocks

```go
// service.go
type UserRepository interface {
    GetUser(id int) (*User, error)
    CreateUser(user *User) error
}

type UserService struct {
    repo UserRepository
}

func (s *UserService) GetUser(id int) (*User, error) {
    return s.repo.GetUser(id)
}

// service_test.go
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) GetUser(id int) (*User, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) CreateUser(user *User) error {
    args := m.Called(user)
    return args.Error(0)
}

// Test
func TestUserService_GetUser(t *testing.T) {
    mockRepo := new(MockUserRepository)
    service := &UserService{repo: mockRepo}
    
    expectedUser := &User{ID: 1, Name: "John"}
    mockRepo.On("GetUser", 1).Return(expectedUser, nil)
    
    user, err := service.GetUser(1)
    
    assert.NoError(t, err)
    assert.Equal(t, expectedUser, user)
    mockRepo.AssertExpectations(t)
}
```

### Using testify/mock

```go
import "github.com/stretchr/testify/mock"

type MockDB struct {
    mock.Mock
}

func (m *MockDB) Query(sql string, args ...interface{}) (*sql.Rows, error) {
    mockArgs := m.Called(sql, args)
    return mockArgs.Get(0).(*sql.Rows), mockArgs.Error(1)
}

func TestService(t *testing.T) {
    mockDB := new(MockDB)
    service := NewService(mockDB)
    
    // Setup expectations
    mockDB.On("Query", "SELECT * FROM users WHERE id = ?", 1).
        Return(&sql.Rows{}, nil)
    
    // Run test
    _, err := service.GetUser(1)
    
    // Verify
    assert.NoError(t, err)
    mockDB.AssertExpectations(t)
}
```

### HTTP Mocking (httptest)

```go
func TestHTTPHandler(t *testing.T) {
    req := httptest.NewRequest("GET", "/users/1", nil)
    w := httptest.NewRecorder()
    
    handler := NewUserHandler(service)
    handler.GetUser(w, req)
    
    resp := w.Result()
    assert.Equal(t, http.StatusOK, resp.StatusCode)
    
    var user User
    json.NewDecoder(resp.Body).Decode(&user)
    assert.Equal(t, 1, user.ID)
}
```

---

## Integration Testing

### Database Tests

```go
func TestUserRepository_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    // Setup test database
    db, err := sql.Open("postgres", testDBURL)
    require.NoError(t, err)
    defer db.Close()
    
    // Cleanup
    defer db.Exec("TRUNCATE users CASCADE")
    
    repo := NewUserRepository(db)
    
    // Test Create
    user := &User{Name: "Test", Email: "test@example.com"}
    err = repo.Create(user)
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
    
    // Test Get
    found, err := repo.GetByID(user.ID)
    assert.NoError(t, err)
    assert.Equal(t, user.Name, found.Name)
}

// Run: go test -v (skips integration tests)
// Run: go test -v -short=false (runs all tests)
```

### Docker Test Containers

```go
import (
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
)

func setupPostgres(t *testing.T) (*sql.DB, func()) {
    ctx := context.Background()
    
    req := testcontainers.ContainerRequest{
        Image:        "postgres:14",
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_PASSWORD": "test",
            "POSTGRES_DB":       "testdb",
        },
        WaitingFor: wait.ForLog("database system is ready"),
    }
    
    container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    require.NoError(t, err)
    
    host, _ := container.Host(ctx)
    port, _ := container.MappedPort(ctx, "5432")
    
    dsn := fmt.Sprintf("postgres://postgres:test@%s:%s/testdb?sslmode=disable", 
        host, port.Port())
    db, err := sql.Open("postgres", dsn)
    require.NoError(t, err)
    
    cleanup := func() {
        db.Close()
        container.Terminate(ctx)
    }
    
    return db, cleanup
}

func TestWithDocker(t *testing.T) {
    db, cleanup := setupPostgres(t)
    defer cleanup()
    
    // Run tests with real database
    repo := NewRepository(db)
    // ... test
}
```

---

## Test Coverage

### Run Coverage

```bash
# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out

# Coverage per function
go tool cover -func=coverage.out
```

### Coverage in Code

```go
// Use build tags for test-only code
// +build !test

func ProductionOnlyCode() {
    // This won't be included in test coverage
}
```

---

## Benchmarking

```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(100, 200)
    }
}

func BenchmarkConcatStrings(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = "Hello" + " " + "World"
    }
}

// Run: go test -bench=. -benchmem
```

### Table-Driven Benchmarks

```go
func BenchmarkFibonacci(b *testing.B) {
    tests := []struct {
        name string
        n    int
    }{
        {"small", 10},
        {"medium", 20},
        {"large", 30},
    }
    
    for _, tt := range tests {
        b.Run(tt.name, func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                Fibonacci(tt.n)
            }
        })
    }
}
```

---

## Test Helpers

### Setup/Teardown

```go
func TestMain(m *testing.M) {
    // Setup
    setup()
    
    // Run tests
    code := m.Run()
    
    // Teardown
    teardown()
    
    os.Exit(code)
}

func setup() {
    // Initialize test resources
    db = connectTestDB()
}

func teardown() {
    // Cleanup
    db.Close()
}
```

### Test Fixtures

```go
func newTestUser(t *testing.T) *User {
    t.Helper() // Marks this as a helper function
    return &User{
        ID:    1,
        Name:  "Test User",
        Email: "test@example.com",
    }
}

func TestSomething(t *testing.T) {
    user := newTestUser(t)
    // Use user in test
}
```

---

## Best Practices

✅ **Test naming:** `Test<Function>_<Scenario>`
✅ **Use table-driven tests** for multiple cases
✅ **Test behavior, not implementation**
✅ **Mock external dependencies**
✅ **Write tests before fixing bugs** (TDD)
✅ **Keep tests fast** (< 100ms per unit test)
✅ **Use t.Helper()** for test utilities
✅ **Clean up resources** (defer cleanup())
✅ **Test edge cases** (nil, empty, negative)
✅ **Aim for 80%+ coverage** (but don't obsess)

---

## Common Patterns

### Testing Errors

```go
func TestDivide_ByZero(t *testing.T) {
    _, err := Divide(10, 0)
    assert.Error(t, err)
    assert.Equal(t, ErrDivisionByZero, err)
}
```

### Testing Panics

```go
func TestPanicFunction(t *testing.T) {
    assert.Panics(t, func() {
        PanicFunction()
    }, "should panic")
}
```

### Testing Timeouts

```go
func TestSlowOperation(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()
    
    err := SlowOperation(ctx)
    assert.NoError(t, err)
}
```

---

## CI Integration

```yaml
# .github/workflows/test.yml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.21
      
      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out ./...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v2
        with:
          file: ./coverage.out
```

---

**Testing is not optional - it's part of the code!** ✅
