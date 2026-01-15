# Вправа 3: HTTP Service з Context Timeout

## Ціль
Створити HTTP сервіс з правильним використанням context для timeouts та cancellation.

---

## Завдання

Створіть програму `api_service.go`, яка:

1. HTTP сервер з endpoint `/users/{id}`
2. Симулює повільні операції (DB query, external API)
3. Використовує context для timeout control
4. Правильно обробляє cancellation
5. Логує час виконання та причини fail
6. **НЕ зберігає context в struct!**

---

## Вимоги

### API Endpoints

```go
GET  /users/{id}        - отримати користувача (3s timeout)
POST /users             - створити користувача (5s timeout)
GET  /users/{id}/posts  - отримати пости користувача (2s timeout)
```

### Service Structure

```go
type UserService struct {
    db  *Database      // БЕЗ context в struct!
    api *ExternalAPI
}

// ✅ Context як параметр
func (s *UserService) GetUser(ctx context.Context, id int) (*User, error)

// ❌ НІ! Context в struct
type BadService struct {
    ctx context.Context  // НІКОЛИ!
}
```

### Database Layer

```go
type Database struct {
    // Симуляція різних затримок
}

func (db *Database) QueryUser(ctx context.Context, id int) (*User, error) {
    // Симулює DB query з затримкою
    // Має перевіряти ctx.Done()
}
```

### External API Layer

```go
type ExternalAPI struct {
    baseURL string
}

func (api *ExternalAPI) FetchUserPosts(ctx context.Context, userID int) ([]Post, error) {
    // Симулює HTTP request
    // Має поважати context timeout
}
```

---

## Приклад використання

```go
func main() {
    db := &Database{}
    api := &ExternalAPI{baseURL: "https://api.example.com"}
    
    service := &UserService{
        db:  db,
        api: api,
    }
    
    // Створюємо HTTP handlers
    http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
        // Context від HTTP request
        ctx := r.Context()
        
        // Додаємо timeout
        ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
        defer cancel()
        
        // Передаємо context в service
        user, err := service.GetUser(ctx, userID)
        if err != nil {
            if errors.Is(err, context.DeadlineExceeded) {
                http.Error(w, "Request timeout", http.StatusGatewayTimeout)
                return
            }
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        json.NewEncoder(w).Encode(user)
    })
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## Очікуваний вивід

```
╔══════════════════════════════════════════╗
║   HTTP Service with Context Timeout      ║
╚══════════════════════════════════════════╝

🚀 Server starting on :8080

🔹 Test 1: Fast Request (Success)
─────────────────────────────────────────
→ GET /users/1
Database query: 500ms
External API call: 300ms
✓ Request completed in 800ms
Response: 200 OK


🔹 Test 2: Slow Request (Timeout)
─────────────────────────────────────────
→ GET /users/2
Database query: 2000ms
External API call: started...
✗ Request cancelled: context deadline exceeded
   Total time: 3000ms
   Reason: timeout reached
Response: 504 Gateway Timeout


🔹 Test 3: Client Cancellation
─────────────────────────────────────────
→ GET /users/3
Database query: started...
Client disconnected after 1s
✗ Request cancelled: context canceled
   Total time: 1000ms
   Reason: client cancelled
Response: (connection closed)


🔹 Test 4: Concurrent Requests
─────────────────────────────────────────
→ GET /users/1 (timeout: 5s)
→ GET /users/2 (timeout: 5s)
→ GET /users/3 (timeout: 5s)

Results:
  User 1: ✓ 800ms
  User 2: ✓ 1200ms
  User 3: ✗ timeout after 5000ms


🔹 Test 5: Graceful Shutdown
─────────────────────────────────────────
Received SIGINT, shutting down...

Waiting for active requests to complete...
   Active requests: 2
   Request 1: completing... ✓
   Request 2: completing... ✓

✓ All requests completed
✓ Server shutdown gracefully


📊 Statistics:
─────────────────────────────────────────
Total requests: 8
Successful: 6
Timeouts: 1
Cancelled: 1
Average response time: 1200ms
```

---

## Підказки

### 1. HTTP Handler з Context

```go
func (s *UserService) HandleGetUser(w http.ResponseWriter, r *http.Request) {
    // 1. Отримати context від request
    ctx := r.Context()
    
    // 2. Додати timeout
    ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()
    
    // 3. Витягнути ID з URL
    id := extractID(r.URL.Path)
    
    // 4. Викликати service з context
    user, err := s.GetUser(ctx, id)
    if err != nil {
        s.handleError(w, err)
        return
    }
    
    // 5. Відправити response
    json.NewEncoder(w).Encode(user)
}
```

### 2. Database Query з Context

```go
func (db *Database) QueryUser(ctx context.Context, id int) (*User, error) {
    // Симулюємо повільний query
    resultChan := make(chan *User, 1)
    errChan := make(chan error, 1)
    
    go func() {
        // Повільна операція
        time.Sleep(2 * time.Second)
        
        // Перевірка перед поверненням результату
        select {
        case <-ctx.Done():
            errChan <- ctx.Err()
        default:
            resultChan <- &User{ID: id, Name: "Alice"}
        }
    }()
    
    // Чекаємо результат або cancellation
    select {
    case <-ctx.Done():
        return nil, fmt.Errorf("query cancelled: %w", ctx.Err())
    case err := <-errChan:
        return nil, err
    case user := <-resultChan:
        return user, nil
    }
}
```

### 3. Перевірка Timeout

```go
func (s *UserService) GetUser(ctx context.Context, id int) (*User, error) {
    // Перевірка перед тяжкою операцією
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }
    
    user, err := s.db.QueryUser(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to query user: %w", err)
    }
    
    return user, nil
}
```

### 4. Error Handling

```go
func handleError(w http.ResponseWriter, err error) {
    switch {
    case errors.Is(err, context.DeadlineExceeded):
        http.Error(w, "Request timeout", http.StatusGatewayTimeout)
    case errors.Is(err, context.Canceled):
        http.Error(w, "Request cancelled", http.StatusRequestTimeout)
    default:
        http.Error(w, "Internal server error", http.StatusInternalServerError)
    }
}
```

---

## Бонус завдання

### 1. Request ID для Tracing

```go
func withRequestID(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        requestID := uuid.New().String()
        ctx := context.WithValue(r.Context(), "requestID", requestID)
        next(w, r.WithContext(ctx))
    }
}
```

### 2. Middleware для Logging

```go
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        next(w, r)
        
        duration := time.Since(start)
        log.Printf("%s %s - %v", r.Method, r.URL.Path, duration)
    }
}
```

### 3. Graceful Shutdown

```go
func main() {
    server := &http.Server{Addr: ":8080"}
    
    go func() {
        if err := server.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()
    
    // Чекаємо SIGINT
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt)
    <-stop
    
    // Graceful shutdown з timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := server.Shutdown(ctx); err != nil {
        log.Fatal(err)
    }
    
    log.Println("Server stopped gracefully")
}
```

### 4. Circuit Breaker

```go
type CircuitBreaker struct {
    maxFailures int
    timeout     time.Duration
    failures    int
    lastFailure time.Time
}

func (cb *CircuitBreaker) Call(ctx context.Context, fn func() error) error {
    if cb.isOpen() {
        return errors.New("circuit breaker open")
    }
    
    err := fn()
    if err != nil {
        cb.recordFailure()
        return err
    }
    
    cb.reset()
    return nil
}
```

---

## Критерії оцінки

- ✅ Context передається як параметр, НЕ в struct
- ✅ Кожен handler має timeout
- ✅ ctx.Done() перевіряється в loops
- ✅ Правильна обробка DeadlineExceeded
- ✅ Правильна обробка Canceled
- ✅ defer cancel() після WithTimeout
- ✅ HTTP responses відповідають статусам помилок

---

## Рішення

Рішення знаходиться в `solutions/solution_3.go`.

Спробуйте виконати завдання самостійно перед тим, як дивитись рішення!

---

## Навчальні цілі

Після виконання цієї вправи ви будете вміти:
- Використовувати context в HTTP handlers
- Додавати timeouts до операцій
- Обробляти cancellation gracefully
- Розуміти чому context НЕ в struct
- Створювати production-ready HTTP сервіси
- Логувати context-related errors

---

## Подальше вдосконалення

Подумайте як додати:
- Rate limiting per user
- Request prioritization
- Distributed tracing (OpenTelemetry)
- Metrics collection (Prometheus)
- Health check endpoint
- Readiness probe
- Load testing з різними timeouts
