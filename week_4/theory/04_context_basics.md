# Context в Go

## Що таке Context?

`context.Context` - це механізм для:
1. **Cancellation** - скасування операцій
2. **Timeouts** - обмеження часу виконання  
3. **Deadlines** - жорсткі дедлайни
4. **Values** - передача request-scoped даних

```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key any) any
}
```

---

## Створення Context

### 1. context.Background()

**Кореневий context**, який ніколи не скасовується.

```go
ctx := context.Background()
```

**Використання:**
- Main function
- Початок HTTP request
- Top-level горутини

### 2. context.TODO()

Placeholder, коли ще не знаєте який context використати.

```go
ctx := context.TODO()
```

---

## WithCancel - Manual Cancellation

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    
    go worker(ctx, "Worker-1")
    
    time.Sleep(2 * time.Second)
    fmt.Println("Cancelling...")
    cancel()  // Скасовуємо context
    
    time.Sleep(1 * time.Second)
}

func worker(ctx context.Context, name string) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("%s: stopped (%v)\n", name, ctx.Err())
            return
        default:
            fmt.Printf("%s: working...\n", name)
            time.Sleep(500 * time.Millisecond)
        }
    }
}
```

**Output:**
```
Worker-1: working...
Worker-1: working...
Worker-1: working...
Cancelling...
Worker-1: stopped (context canceled)
```

---

## WithTimeout - Automatic Cancellation

```go
func main() {
    // Автоматично скасується через 2 секунди
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()  // Завжди викликайте cancel()!
    
    result := make(chan string, 1)
    go slowOperation(ctx, result)
    
    select {
    case <-ctx.Done():
        fmt.Println("Timeout:", ctx.Err())
    case res := <-result:
        fmt.Println("Result:", res)
    }
}

func slowOperation(ctx context.Context, result chan<- string) {
    time.Sleep(5 * time.Second)  // Занадто повільно!
    result <- "completed"
}
```

---

## WithDeadline - Fixed Time

```go
func main() {
    deadline := time.Now().Add(3 * time.Second)
    ctx, cancel := context.WithDeadline(context.Background(), deadline)
    defer cancel()
    
    doWork(ctx)
}

func doWork(ctx context.Context) {
    for i := 0; i < 10; i++ {
        select {
        case <-ctx.Done():
            fmt.Println("Deadline exceeded!")
            return
        default:
            fmt.Printf("Step %d...\n", i+1)
            time.Sleep(500 * time.Millisecond)
        }
    }
}
```

---

## WithValue - Request-Scoped Data

```go
func main() {
    ctx := context.WithValue(context.Background(), "userID", 123)
    ctx = context.WithValue(ctx, "requestID", "abc-def")
    
    handleRequest(ctx)
}

func handleRequest(ctx context.Context) {
    userID := ctx.Value("userID").(int)
    requestID := ctx.Value("requestID").(string)
    
    fmt.Printf("User: %d, Request: %s\n", userID, requestID)
}
```

**⚠️ WARNING:** Використовуйте WithValue обережно!
- Тільки для request-scoped даних
- Не для передачі обов'язкових параметрів
- Краще передавати як explicit parameters

---

## ⚠️ Чому НЕ зберігати Context в Struct?

### ❌ ПОГАНА ПРАКТИКА:

```go
type Service struct {
    ctx context.Context  // НІ! НІ! НІ!
    db  *sql.DB
}

func NewService(db *sql.DB) *Service {
    return &Service{
        ctx: context.Background(),  // Який саме context?
        db:  db,
    }
}

func (s *Service) ProcessUser(userID int) error {
    // Який context використовувати? Той що в struct? Або новий?
    return s.db.QueryContext(s.ctx, "SELECT ...")
}
```

### ✅ ПРАВИЛЬНА ПРАКТИКА:

```go
type Service struct {
    db *sql.DB  // Тільки залежності, БЕЗ context!
}

func NewService(db *sql.DB) *Service {
    return &Service{db: db}
}

func (s *Service) ProcessUser(ctx context.Context, userID int) error {
    // Context приходить як параметр!
    return s.db.QueryContext(ctx, "SELECT ...")
}
```

---

## Чому Context НЕ в Struct? - Детальне пояснення

### Причина 1: Lifetime

```go
// ❌ Проблема: Context має lifetime request
type Handler struct {
    ctx context.Context  // Від якого request?
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // r.Context() - context ЦЬОго request
    // h.ctx - context якогось ІНШОГО request? 🤔
}
```

**Рішення:**
```go
// ✅ Context прив'язаний до операції, не до об'єкта
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()  // Context від request
    h.process(ctx)  // Передаємо як параметр
}
```

### Причина 2: Memory Leaks

```go
// ❌ Context може жити довше ніж потрібно
service := &UserService{
    ctx: ctx,  // ctx може бути вже cancelled!
}

// Service живе довго, ctx cancelled - але ми не знаємо!
```

### Причина 3: Race Conditions

```go
// ❌ Множинні горутини, один context в struct
service := &Service{ctx: ctx}

go service.Process()  // Використовує service.ctx
go service.Process()  // Використовує service.ctx
go service.Process()  // Використовує service.ctx

// Що станеться якщо ctx скасується посередині?
```

### Причина 4: Тестування

```go
// ❌ Важко тестувати з різними contexts
service := &Service{ctx: context.Background()}

// Як протестувати з timeout?
// Як протестувати з cancellation?
// Доведеться створювати новий Service для кожного тесту!
```

```go
// ✅ Легко тестувати
service := &Service{db: mockDB}

// Test 1: з timeout
ctx1, _ := context.WithTimeout(context.Background(), 1*time.Second)
service.Process(ctx1, data)

// Test 2: з cancellation
ctx2, cancel := context.WithCancel(context.Background())
cancel()
service.Process(ctx2, data)
```

---

## Context Best Practices

### 1. Завжди перший параметр

```go
// ✅ Правильно
func DoWork(ctx context.Context, data string) error

// ❌ Неправильно
func DoWork(data string, ctx context.Context) error
func DoWork(data string) error  // Де context?
```

### 2. Завжди defer cancel()

```go
// ✅ Правильно
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
defer cancel()  // Гарантує cleanup

// ❌ Неправильно
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
// Забули cancel() - memory leak!
```

### 3. Не ігноруйте ctx.Done()

```go
// ✅ Правильно
func LongOperation(ctx context.Context) error {
    for i := 0; i < 1000; i++ {
        select {
        case <-ctx.Done():
            return ctx.Err()  // Швидко виходимо
        default:
            doWork(i)
        }
    }
}

// ❌ Неправильно
func LongOperation(ctx context.Context) error {
    for i := 0; i < 1000; i++ {
        doWork(i)  // Ігноруємо cancellation!
    }
}
```

### 4. Propagate Context вниз

```go
// ✅ Context йде вниз по call stack
func HandleRequest(ctx context.Context) error {
    data, err := FetchData(ctx)  // Передаємо ctx
    if err != nil {
        return err
    }
    return ProcessData(ctx, data)  // Передаємо ctx
}
```

---

## Production Pattern: HTTP Server

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "time"
)

type Service struct {
    db *Database  // БЕЗ context!
}

func (s *Service) GetUser(ctx context.Context, id int) (*User, error) {
    // Timeout для DB query
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()
    
    // Використовуємо ctx від request
    return s.db.QueryUser(ctx, id)
}

func (s *Service) HandleHTTP(w http.ResponseWriter, r *http.Request) {
    // Context від HTTP request
    ctx := r.Context()
    
    user, err := s.GetUser(ctx, 123)
    if err != nil {
        if err == context.DeadlineExceeded {
            http.Error(w, "Request timeout", http.StatusGatewayTimeout)
            return
        }
        http.Error(w, "Internal error", http.StatusInternalServerError)
        return
    }
    
    fmt.Fprintf(w, "User: %+v", user)
}

type User struct {
    ID   int
    Name string
}

type Database struct{}

func (db *Database) QueryUser(ctx context.Context, id int) (*User, error) {
    // Перевіряємо cancellation перед тяжкою операцією
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }
    
    // Симуляція DB query
    time.Sleep(100 * time.Millisecond)
    return &User{ID: id, Name: "Alice"}, nil
}

func main() {
    db := &Database{}
    service := &Service{db: db}
    
    http.HandleFunc("/user", service.HandleHTTP)
    fmt.Println("Server running on :8080")
    http.ListenAndServe(":8080", nil)
}
```

---

## Context Hierarchy

```go
func main() {
    // 1. Root context
    root := context.Background()
    
    // 2. Request context (з timeout)
    reqCtx, cancel1 := context.WithTimeout(root, 10*time.Second)
    defer cancel1()
    
    // 3. Database context (з timeout)
    dbCtx, cancel2 := context.WithTimeout(reqCtx, 2*time.Second)
    defer cancel2()
    
    // 4. Query context (з cancellation)
    queryCtx, cancel3 := context.WithCancel(dbCtx)
    defer cancel3()
    
    // Якщо reqCtx скасується, всі дочірні теж скасуються!
}
```

**Візуалізація:**
```
root (Background)
  ↓
reqCtx (10s timeout)
  ↓
dbCtx (2s timeout)
  ↓
queryCtx (manual cancel)
```

---

## Поширені помилки

### ❌ Помилка 1: Забули defer cancel()

```go
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
// Забули defer cancel() - memory leak!
```

### ❌ Помилка 2: Context в struct

```go
type Service struct {
    ctx context.Context  // НІ!
}
```

### ❌ Помилка 3: Ігнорування ctx.Done()

```go
func Work(ctx context.Context) {
    for {
        doHeavyWork()  // Не перевіряємо ctx.Done()!
    }
}
```

### ❌ Помилка 4: Використання cancelled context

```go
ctx, cancel := context.WithCancel(context.Background())
cancel()  // Скасували

// Спроба використати cancelled context
DoWork(ctx)  // Одразу fail!
```

---

## Ключові моменти

1. ✅ Context завжди перший параметр (`ctx context.Context`)
2. ✅ Завжди `defer cancel()`
3. ✅ Перевіряйте `ctx.Done()` в loops
4. ✅ **НІКОЛИ** не зберігайте context в struct
5. ✅ Context прив'язаний до операції, не до об'єкта
6. ✅ Propagate context вниз по call stack
7. ✅ Використовуйте context для cancellation/timeout, не для передачі даних

---

## Домашнє завдання

1. Створіть HTTP server з timeout для кожного request
2. Напишіть функцію, яка поважає ctx.Done()
3. Створіть ланцюжок з 3 contexts (root → request → db)
4. Поясніть товаришу чому context НЕ в struct

---

## Додаткове читання

- [Go Blog: Context](https://go.dev/blog/context)
- [Context package docs](https://pkg.go.dev/context)
- [Context and structs](https://go.dev/blog/context-and-structs)

---

**Вітаємо! Ви завершили теорію Тижня 4! 🎉**

**Наступний крок:** Practice Examples (`practice/`)
