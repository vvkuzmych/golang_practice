# Go Best Practices & Efficient Code

## 📖 Зміст

1. [Project Structure](#1-project-structure)
2. [Error Handling](#2-error-handling)
3. [Memory Management](#3-memory-management)
4. [Go Idioms](#4-go-idioms)
5. [Anti-Patterns](#5-anti-patterns)

---

## 1. Project Structure

### Standard Go Project Layout

```
myapp/
├── cmd/                    # Main applications
│   └── api/
│       └── main.go
├── internal/               # Private code
│   ├── handler/
│   ├── service/
│   ├── repository/
│   └── model/
├── pkg/                    # Public libraries
│   └── logger/
├── api/                    # API definitions (OpenAPI, protobuf)
├── web/                    # Web assets
├── configs/                # Configuration files
├── scripts/                # Build scripts
├── test/                   # Additional test data
├── docs/                   # Documentation
├── go.mod
├── go.sum
├── Makefile
├── Dockerfile
└── README.md
```

### Пояснення

- **cmd/** - точки входу (main packages)
- **internal/** - приватний код (не може бути імпортований ззовні)
- **pkg/** - публічні бібліотеки (можуть бути переиспользовані)
- **api/** - API contracts

---

## 2. Error Handling

### ✅ Правильно

```go
// Wrap errors з context
func ReadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("read config %s: %w", path, err)
    }
    
    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("parse config: %w", err)
    }
    
    return &cfg, nil
}

// Custom errors
var ErrNotFound = errors.New("not found")

func GetUser(id int) (*User, error) {
    user, ok := cache[id]
    if !ok {
        return nil, fmt.Errorf("user %d: %w", id, ErrNotFound)
    }
    return user, nil
}

// Error checking
if errors.Is(err, ErrNotFound) {
    // Handle not found
}
```

### ❌ Неправильно

```go
// Не ігноруйте errors
data, _ := os.ReadFile("config.json") // BAD!

// Не використовуйте panic для бізнес-логіки
if user == nil {
    panic("user not found") // BAD!
}

// Не втрачайте context
if err != nil {
    return errors.New("failed") // Втратили оригінальну помилку
}
```

---

## 3. Memory Management

### Slice pre-allocation

```go
// ❌ Неефективно
var users []User
for i := 0; i < 1000; i++ {
    users = append(users, User{ID: i})
}

// ✅ Ефективно
users := make([]User, 0, 1000) // pre-allocate capacity
for i := 0; i < 1000; i++ {
    users = append(users, User{ID: i})
}
```

### String concatenation

```go
// ❌ Неефективно (багато allocations)
result := ""
for i := 0; i < 1000; i++ {
    result += strconv.Itoa(i)
}

// ✅ Ефективно
var builder strings.Builder
builder.Grow(1000 * 4) // pre-allocate
for i := 0; i < 1000; i++ {
    builder.WriteString(strconv.Itoa(i))
}
result := builder.String()
```

### Pointer vs Value

```go
// Великі struct - передавайте pointer
type BigStruct struct {
    Data [1000]int
}

// ✅ Good
func ProcessBig(s *BigStruct) {
    // No copy
}

// ❌ Bad - копіює всю структуру
func ProcessBig(s BigStruct) {
    // Copies 8KB
}

// Малі struct - value OK
type Point struct {
    X, Y int
}

// ✅ OK
func Distance(p Point) float64 {
    return math.Sqrt(float64(p.X*p.X + p.Y*p.Y))
}
```

---

## 4. Go Idioms

### Accept interfaces, return structs

```go
// ✅ Good
type Reader interface {
    Read(p []byte) (n int, err error)
}

func ProcessData(r Reader) error {
    // Працює з будь-яким Reader
}

// Return concrete type
func NewFile(path string) *File {
    return &File{path: path}
}
```

### Keep interfaces small

```go
// ✅ Good - single method
type Stringer interface {
    String() string
}

// ❌ Bad - too many methods
type UserService interface {
    Create(user User) error
    Update(user User) error
    Delete(id int) error
    Get(id int) (User, error)
    List() ([]User, error)
    // ... 10 more methods
}
```

### Use context

```go
func Handler(ctx context.Context, req *Request) error {
    // Timeout
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    // Pass context
    return service.Process(ctx, req)
}
```

---

## 5. Anti-Patterns

### ❌ Goroutine leaks

```go
// BAD - goroutine never stops
func leak() {
    ch := make(chan int)
    go func() {
        val := <-ch // Blocks forever
        fmt.Println(val)
    }()
}

// ✅ GOOD - with context
func noLeak(ctx context.Context) {
    ch := make(chan int)
    go func() {
        select {
        case val := <-ch:
            fmt.Println(val)
        case <-ctx.Done():
            return
        }
    }()
}
```

### ❌ Не використовуйте init()

```go
// ❌ BAD - непередбачуваний порядок
func init() {
    db = connectDatabase()
}

// ✅ GOOD - explicit initialization
func main() {
    db, err := NewDatabase(config)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
}
```

---

## Best Practices Summary

✅ **DO:**
- Use `go fmt`, `go vet`, `golangci-lint`
- Write tests
- Handle all errors
- Use context for cancellation
- Document public APIs
- Pre-allocate slices/maps
- Close resources (defer)

❌ **DON'T:**
- Ignore errors
- Use panic for business logic
- Create goroutine leaks
- Over-engineer
- Mutate shared state without locks
