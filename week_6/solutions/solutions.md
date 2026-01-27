# Рішення Вправ - Week 6

## Exercise 1: ООП і Патерни

### Інкапсуляція

```go
package main

type Book struct {
    title     string
    author    string
    isbn      string
    available bool
}

func NewBook(title, author, isbn string) *Book {
    return &Book{
        title:     title,
        author:    author,
        isbn:      isbn,
        available: true,
    }
}

func (b *Book) GetTitle() string { return b.title }
func (b *Book) GetAuthor() string { return b.author }
func (b *Book) GetISBN() string { return b.isbn }
func (b *Book) IsAvailable() bool { return b.available }

func (b *Book) Borrow() error {
    if !b.available {
        return fmt.Errorf("book is not available")
    }
    b.available = false
    return nil
}

func (b *Book) Return() {
    b.available = true
}
```

### Поліморфізм

```go
type LibraryItem interface {
    GetTitle() string
    GetType() string
    CanBeBorrowed() bool
}

func (b *Book) GetType() string { return "Book" }
func (b *Book) CanBeBorrowed() bool { return true }

type Magazine struct {
    title string
}

func (m *Magazine) GetTitle() string { return m.title }
func (m *Magazine) GetType() string { return "Magazine" }
func (m *Magazine) CanBeBorrowed() bool { return false }

type ReferenceBook struct {
    title string
}

func (r *ReferenceBook) GetTitle() string { return r.title }
func (r *ReferenceBook) GetType() string { return "Reference Book" }
func (r *ReferenceBook) CanBeBorrowed() bool { return false }
```

### Композиція

```go
type Library struct {
    name  string
    items []LibraryItem
}

func (l *Library) AddItem(item LibraryItem) {
    l.items = append(l.items, item)
}

func (l *Library) FindByTitle(title string) LibraryItem {
    for _, item := range l.items {
        if item.GetTitle() == title {
            return item
        }
    }
    return nil
}

func (l *Library) ListAvailableItems() []LibraryItem {
    var available []LibraryItem
    for _, item := range l.items {
        if item.CanBeBorrowed() {
            available = append(available, item)
        }
    }
    return available
}
```

### Singleton

```go
type LibraryManager struct {
    libraries map[string]*Library
}

var (
    instance *LibraryManager
    once     sync.Once
)

func GetLibraryManager() *LibraryManager {
    once.Do(func() {
        instance = &LibraryManager{
            libraries: make(map[string]*Library),
        }
    })
    return instance
}

func (m *LibraryManager) RegisterLibrary(name string) {
    m.libraries[name] = &Library{name: name}
}
```

### Factory

```go
func CreateLibraryItem(itemType, title, author string) LibraryItem {
    switch itemType {
    case "book":
        return NewBook(title, author, "ISBN-XXX")
    case "magazine":
        return &Magazine{title: title}
    case "reference":
        return &ReferenceBook{title: title}
    default:
        return nil
    }
}
```

---

## Exercise 2: REST API Server

### Повне рішення

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "strings"
    "sync"
    "time"
)

type Todo struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Completed   bool      `json:"completed"`
    CreatedAt   time.Time `json:"created_at"`
}

type TodoStore struct {
    mu     sync.RWMutex
    todos  map[int]Todo
    nextID int
}

type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
    Code    int    `json:"code"`
}

// Logging Middleware
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
        next.ServeHTTP(w, r)
        log.Printf("Took: %v", time.Since(start))
    })
}

// CORS Middleware
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-Key")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

// Auth Middleware
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        apiKey := r.Header.Get("X-API-Key")
        
        if apiKey != "secret123" {
            writeError(w, http.StatusUnauthorized, "Invalid API key")
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

func writeError(w http.ResponseWriter, code int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(ErrorResponse{
        Error:   http.StatusText(code),
        Message: message,
        Code:    code,
    })
}

func validateTodo(todo *Todo) error {
    if strings.TrimSpace(todo.Title) == "" {
        return fmt.Errorf("title is required")
    }
    if len(todo.Title) > 100 {
        return fmt.Errorf("title too long (max 100 chars)")
    }
    if len(todo.Description) > 500 {
        return fmt.Errorf("description too long (max 500 chars)")
    }
    return nil
}

func main() {
    store := &TodoStore{
        todos: make(map[int]Todo),
        nextID: 1,
    }
    
    mux := http.NewServeMux()
    
    // Routes
    mux.HandleFunc("/api/todos", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        switch r.Method {
        case http.MethodGet:
            store.mu.RLock()
            todos := make([]Todo, 0, len(store.todos))
            for _, todo := range store.todos {
                todos = append(todos, todo)
            }
            store.mu.RUnlock()
            json.NewEncoder(w).Encode(todos)
            
        case http.MethodPost:
            var todo Todo
            if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
                writeError(w, http.StatusBadRequest, "Invalid JSON")
                return
            }
            
            if err := validateTodo(&todo); err != nil {
                writeError(w, http.StatusBadRequest, err.Error())
                return
            }
            
            store.mu.Lock()
            todo.ID = store.nextID
            todo.CreatedAt = time.Now()
            store.todos[store.nextID] = todo
            store.nextID++
            store.mu.Unlock()
            
            w.WriteHeader(http.StatusCreated)
            json.NewEncoder(w).Encode(todo)
            
        default:
            writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
        }
    })
    
    // Chain middleware
    handler := loggingMiddleware(
        corsMiddleware(
            authMiddleware(mux),
        ),
    )
    
    log.Println("Server started at :8080")
    log.Fatal(http.ListenAndServe(":8080", handler))
}
```

---

## Exercise 3: Мікросервіси

### Product Service (основні частини)

```go
package main

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "net/http"
)

type Product struct {
    ID        uint    `gorm:"primaryKey" json:"id"`
    Name      string  `json:"name"`
    Price     float64 `json:"price"`
    Stock     int     `json:"stock"`
    CreatedAt time.Time `json:"created_at"`
}

func main() {
    dsn := "host=localhost user=postgres password=postgres dbname=products_db port=5432"
    db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    db.AutoMigrate(&Product{})
    
    http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "GET" {
            var products []Product
            db.Find(&products)
            json.NewEncoder(w).Encode(products)
        } else if r.Method == "POST" {
            var product Product
            json.NewDecoder(r.Body).Decode(&product)
            db.Create(&product)
            json.NewEncoder(w).Encode(product)
        }
    })
    
    log.Fatal(http.ListenAndServe(":8001", nil))
}
```

### API Gateway (основні частини)

```go
func (g *Gateway) OrderDetailsHandler(w http.ResponseWriter, r *http.Request) {
    orderID := r.URL.Query().Get("id")
    
    // 1. Get order
    order := g.fetchOrder(orderID)
    
    // 2. Get product details
    product := g.fetchProduct(order.ProductID)
    
    // 3. Get user details
    user := g.fetchUser(order.UserID)
    
    // 4. Combine
    response := map[string]interface{}{
        "order":   order,
        "product": product,
        "user":    user,
    }
    
    json.NewEncoder(w).Encode(response)
}
```

---

## 💡 Підказки для самостійного вивчення

1. Використовуйте gorilla/mux або chi для зручного роутингу
2. Логуйте всі помилки
3. Тестуйте з `curl` або Postman
4. Використовуйте `.env` файли для конфігурації
5. Пишіть unit tests для критичної логіки

---

**Це рішення для довідки. Спробуйте спочатку самостійно!** 💪
