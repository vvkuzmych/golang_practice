# Golang net/http

Пакет `net/http` надає все необхідне для створення HTTP клієнтів та серверів.

---

## 📖 Зміст

1. [HTTP Server](#1-http-server)
2. [HTTP Client](#2-http-client)
3. [Routing](#3-routing)
4. [Middleware](#4-middleware)
5. [Context & Timeouts](#5-context--timeouts)
6. [Error Handling](#6-error-handling)

---

## 1. HTTP Server

### Простий сервер

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    // Простий handler
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, World!")
    })
    
    // Запуск сервера
    fmt.Println("Server started at :8080")
    http.ListenAndServe(":8080", nil)
}
```

### Handler Interface

```go
// Handler - це інтерфейс
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}

// Власний handler
type MyHandler struct {
    message string
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, h.message)
}

func main() {
    handler := &MyHandler{message: "Custom Handler"}
    http.Handle("/", handler)
    http.ListenAndServe(":8080", nil)
}
```

### Request & Response

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // ===== Request =====
    
    // HTTP метод
    fmt.Println("Method:", r.Method) // GET, POST, PUT, DELETE, etc.
    
    // URL
    fmt.Println("Path:", r.URL.Path)   // /api/users
    fmt.Println("Query:", r.URL.Query()) // map[key:[value]]
    
    // Headers
    fmt.Println("User-Agent:", r.Header.Get("User-Agent"))
    fmt.Println("Content-Type:", r.Header.Get("Content-Type"))
    
    // Body
    body, _ := io.ReadAll(r.Body)
    defer r.Body.Close()
    fmt.Println("Body:", string(body))
    
    // Query параметри
    name := r.URL.Query().Get("name")
    age := r.URL.Query().Get("age")
    
    // Form дані
    r.ParseForm()
    username := r.FormValue("username")
    
    // ===== Response =====
    
    // Status code
    w.WriteHeader(http.StatusOK) // 200
    
    // Headers
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("X-Custom-Header", "value")
    
    // Body
    w.Write([]byte(`{"message": "success"}`))
    // або
    fmt.Fprintf(w, `{"message": "success"}`)
}
```

### JSON API

```go
package main

import (
    "encoding/json"
    "net/http"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func main() {
    // GET /api/users
    http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
        }
        
        users := []User{
            {ID: 1, Name: "John", Email: "john@example.com"},
            {ID: 2, Name: "Jane", Email: "jane@example.com"},
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(users)
    })
    
    // POST /api/users
    http.HandleFunc("/api/users/create", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
        }
        
        var user User
        if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{
                "error": "Invalid JSON",
            })
            return
        }
        defer r.Body.Close()
        
        // Логіка створення...
        user.ID = 3
        
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated) // 201
        json.NewEncoder(w).Encode(user)
    })
    
    http.ListenAndServe(":8080", nil)
}
```

### ServeMux (Router)

```go
func main() {
    // Створюємо власний mux
    mux := http.NewServeMux()
    
    mux.HandleFunc("/", homeHandler)
    mux.HandleFunc("/about", aboutHandler)
    mux.HandleFunc("/api/users", usersHandler)
    
    // Serve static files
    fs := http.FileServer(http.Dir("./static"))
    mux.Handle("/static/", http.StripPrefix("/static/", fs))
    
    http.ListenAndServe(":8080", mux)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    fmt.Fprintf(w, "Home Page")
}
```

---

## 2. HTTP Client

### Простий GET запит

```go
package main

import (
    "fmt"
    "io"
    "net/http"
)

func main() {
    // Простий GET
    resp, err := http.Get("https://api.github.com/users/golang")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    // Status code
    fmt.Println("Status:", resp.Status) // "200 OK"
    fmt.Println("Status Code:", resp.StatusCode) // 200
    
    // Headers
    fmt.Println("Content-Type:", resp.Header.Get("Content-Type"))
    
    // Body
    body, _ := io.ReadAll(resp.Body)
    fmt.Println("Body:", string(body))
}
```

### POST запит з JSON

```go
func main() {
    user := map[string]string{
        "name":  "John",
        "email": "john@example.com",
    }
    
    jsonData, _ := json.Marshal(user)
    
    resp, err := http.Post(
        "http://localhost:8080/api/users",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)
    fmt.Println(result)
}
```

### Кастомний запит

```go
func main() {
    // Створюємо запит
    req, _ := http.NewRequest("PUT", "http://localhost:8080/api/users/1", bytes.NewBuffer(jsonData))
    
    // Додаємо headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer token123")
    req.Header.Set("User-Agent", "MyApp/1.0")
    
    // Виконуємо запит
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

### HTTP Client з налаштуваннями

```go
func main() {
    // Кастомний client з timeouts
    client := &http.Client{
        Timeout: 10 * time.Second,
        Transport: &http.Transport{
            MaxIdleConns:        100,
            MaxIdleConnsPerHost: 10,
            IdleConnTimeout:     30 * time.Second,
        },
    }
    
    resp, err := client.Get("https://api.example.com/data")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
}
```

---

## 3. Routing

### Проблема стандартного роутера

```go
// ❌ Стандартний http.ServeMux не підтримує:
// - Параметри в URL (/users/:id)
// - HTTP методи (GET, POST, etc.)
// - Middleware chains
// - Групи роутів
```

### Рішення: gorilla/mux

```go
package main

import (
    "encoding/json"
    "net/http"
    
    "github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()
    
    // ===== Route Parameters =====
    r.HandleFunc("/users/{id}", getUserHandler).Methods("GET")
    r.HandleFunc("/users/{id}", updateUserHandler).Methods("PUT")
    r.HandleFunc("/users/{id}", deleteUserHandler).Methods("DELETE")
    
    // ===== Query Parameters =====
    r.HandleFunc("/search", searchHandler).Methods("GET").
        Queries("q", "{query}", "page", "{page}")
    
    // ===== Subrouters (групи) =====
    api := r.PathPrefix("/api/v1").Subrouter()
    api.HandleFunc("/products", getProductsHandler).Methods("GET")
    api.HandleFunc("/products", createProductHandler).Methods("POST")
    
    // ===== Static Files =====
    r.PathPrefix("/static/").Handler(
        http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))),
    )
    
    http.ListenAndServe(":8080", r)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
    // Отримуємо параметр з URL
    vars := mux.Vars(r)
    userID := vars["id"]
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "id":   userID,
        "name": "John Doe",
    })
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    query := vars["query"]
    page := vars["page"]
    
    fmt.Fprintf(w, "Searching for: %s, Page: %s", query, page)
}
```

### Chi Router (альтернатива)

```go
package main

import (
    "net/http"
    
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func main() {
    r := chi.NewRouter()
    
    // Middleware
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    
    // Routes
    r.Get("/", homeHandler)
    r.Route("/api", func(r chi.Router) {
        r.Get("/users", listUsersHandler)
        r.Post("/users", createUserHandler)
        
        r.Route("/users/{userID}", func(r chi.Router) {
            r.Get("/", getUserHandler)
            r.Put("/", updateUserHandler)
            r.Delete("/", deleteUserHandler)
        })
    })
    
    http.ListenAndServe(":8080", r)
}
```

---

## 4. Middleware

### Що таке Middleware?

Middleware - це функція, яка обробляє запит **перед** або **після** handler'а.

```
Request → Middleware 1 → Middleware 2 → Handler → Middleware 2 → Middleware 1 → Response
```

### Простий Middleware

```go
// Logger middleware
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // До handler'а
        fmt.Printf("[%s] %s %s\n", r.Method, r.URL.Path, r.RemoteAddr)
        
        // Викликаємо наступний handler
        next.ServeHTTP(w, r)
        
        // Після handler'а
        fmt.Printf("Request took: %v\n", time.Since(start))
    })
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", homeHandler)
    
    // Обгортаємо в middleware
    http.ListenAndServe(":8080", loggingMiddleware(mux))
}
```

### Chain of Middleware

```go
// Auth middleware
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        
        if token != "Bearer secret123" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

// Recovery middleware
func recoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                fmt.Printf("Panic recovered: %v\n", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }()
        
        next.ServeHTTP(w, r)
    })
}

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", homeHandler)
    
    // ✅ Chain middleware
    handler := recoveryMiddleware(
        loggingMiddleware(
            corsMiddleware(
                authMiddleware(mux),
            ),
        ),
    )
    
    http.ListenAndServe(":8080", handler)
}
```

### Middleware для конкретних роутів

```go
func main() {
    r := mux.NewRouter()
    
    // Public routes (без auth)
    r.HandleFunc("/", homeHandler)
    r.HandleFunc("/login", loginHandler)
    
    // Protected routes (з auth middleware)
    protected := r.PathPrefix("/api").Subrouter()
    protected.Use(authMiddleware)
    protected.HandleFunc("/users", getUsersHandler)
    protected.HandleFunc("/products", getProductsHandler)
    
    http.ListenAndServe(":8080", r)
}
```

---

## 5. Context & Timeouts

### Context в HTTP

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Кожен request має context
    ctx := r.Context()
    
    // Context скасовується коли:
    // 1. Client закриває з'єднання
    // 2. Таймаут спрацьовує
    // 3. Викликається cancel()
    
    select {
    case <-time.After(5 * time.Second):
        fmt.Fprintf(w, "Work completed")
    case <-ctx.Done():
        // Request скасовано
        fmt.Println("Request cancelled:", ctx.Err())
        http.Error(w, "Request cancelled", 499)
        return
    }
}
```

### Передача даних через Context

```go
type contextKey string

const userKey contextKey = "user"

// Middleware додає дані в context
func userMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        userID := r.Header.Get("X-User-ID")
        
        // Додаємо userID в context
        ctx := context.WithValue(r.Context(), userKey, userID)
        
        // Створюємо новий request з оновленим context
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Handler отримує дані з context
func handler(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(userKey).(string)
    fmt.Fprintf(w, "User ID: %s", userID)
}
```

### Timeouts

```go
func main() {
    // Server timeouts
    srv := &http.Server{
        Addr:         ":8080",
        Handler:      handler,
        ReadTimeout:  5 * time.Second,  // час на читання request
        WriteTimeout: 10 * time.Second, // час на запис response
        IdleTimeout:  120 * time.Second, // час до закриття keep-alive з'єднання
    }
    
    srv.ListenAndServe()
}

// Client timeouts
func makeRequest() {
    client := &http.Client{
        Timeout: 10 * time.Second, // загальний timeout
    }
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    req, _ := http.NewRequestWithContext(ctx, "GET", "http://example.com", nil)
    resp, err := client.Do(req)
    if err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            fmt.Println("Request timed out")
        }
        return
    }
    defer resp.Body.Close()
}
```

---

## 6. Error Handling

### Централізована обробка помилок

```go
// ErrorResponse - стандартна структура помилки
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
    Code    int    `json:"code"`
}

// HTTPError - кастомна помилка
type HTTPError struct {
    Code    int
    Message string
}

func (e *HTTPError) Error() string {
    return e.Message
}

// WriteError - запис помилки в response
func WriteError(w http.ResponseWriter, err error) {
    var httpErr *HTTPError
    
    if errors.As(err, &httpErr) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(httpErr.Code)
        json.NewEncoder(w).Encode(ErrorResponse{
            Error:   http.StatusText(httpErr.Code),
            Message: httpErr.Message,
            Code:    httpErr.Code,
        })
    } else {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(ErrorResponse{
            Error:   "Internal Server Error",
            Message: "An unexpected error occurred",
            Code:    500,
        })
    }
}

func handler(w http.ResponseWriter, r *http.Request) {
    user, err := getUser(123)
    if err != nil {
        WriteError(w, err)
        return
    }
    
    json.NewEncoder(w).Encode(user)
}

func getUser(id int) (*User, error) {
    // Якщо користувача не знайдено
    return nil, &HTTPError{
        Code:    http.StatusNotFound,
        Message: "User not found",
    }
}
```

---

## ✅ Best Practices

1. **Завжди закривайте Body**: `defer resp.Body.Close()`
2. **Використовуйте Context**: для cancellation та timeouts
3. **Встановлюйте timeouts**: на сервері та клієнті
4. **Structured Logging**: використовуйте structured логи
5. **Error Handling**: централізована обробка помилок
6. **Middleware**: для cross-cutting concerns
7. **Graceful Shutdown**: коректне завершення сервера

```go
// Graceful Shutdown
func main() {
    srv := &http.Server{Addr: ":8080", Handler: handler}
    
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()
    
    // Чекаємо на сигнал
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    // Gracefully shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }
    
    log.Println("Server exited")
}
```

---

**Далі:** [04_microservices.md](./04_microservices.md)
