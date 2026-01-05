package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// ============= http.Handler Interface =============

// type Handler interface {
//     ServeHTTP(ResponseWriter, *Request)
// }

// ============= Simple Handler =============

type HelloHandler struct{}

func (h HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, World!\n")
}

// ============= Handler з даними =============

type GreetHandler struct {
	Name string
}

func (g GreetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!\n", g.Name)
}

// ============= Counter Handler =============

type CounterHandler struct {
	count int
}

func (c *CounterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.count++
	fmt.Fprintf(w, "Request count: %d\n", c.count)
}

// ============= JSON Handler =============

type StatusHandler struct{}

func (s StatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
}

// ============= Middleware Pattern =============

// LoggingMiddleware логує кожен запит
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("Completed in %v", time.Since(start))
	})
}

// AuthMiddleware перевіряє авторизацію
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if token != "Bearer secret-token" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RecoveryMiddleware ловить паніки
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// ============= Router Pattern =============

type Router struct {
	routes map[string]http.Handler
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]http.Handler),
	}
}

func (r *Router) Handle(path string, handler http.Handler) {
	r.routes[path] = handler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if handler, ok := r.routes[req.URL.Path]; ok {
		handler.ServeHTTP(w, req)
		return
	}

	http.NotFound(w, req)
}

// ============= Method Handler =============

type MethodHandler struct {
	GET    http.Handler
	POST   http.Handler
	PUT    http.Handler
	DELETE http.Handler
}

func (m MethodHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if m.GET != nil {
			m.GET.ServeHTTP(w, r)
			return
		}
	case http.MethodPost:
		if m.POST != nil {
			m.POST.ServeHTTP(w, r)
			return
		}
	case http.MethodPut:
		if m.PUT != nil {
			m.PUT.ServeHTTP(w, r)
			return
		}
	case http.MethodDelete:
		if m.DELETE != nil {
			m.DELETE.ServeHTTP(w, r)
			return
		}
	}

	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// ============= API Examples =============

// UserHandler обробляє операції з користувачами
type UserHandler struct{}

func (u UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `[{"id":1,"name":"John"},{"id":2,"name":"Jane"}]`)
	case http.MethodPost:
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{"id":3,"name":"New User"}`)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// ============= Main =============

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║       http.Handler Interface             ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	fmt.Println("\nНУВАГА: Цей приклад показує структуру,")
	fmt.Println("але не запускає сервер для демонстрації.")
	fmt.Println()
	fmt.Println("Для запуску справжнього сервера розкоментуйте")
	fmt.Println("останні рядки та запустіть програму.")

	// ===== Simple Handler =====
	fmt.Println("\n🔹 Simple Handler")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("type HelloHandler struct{}")
	fmt.Println("")
	fmt.Println("func (h HelloHandler) ServeHTTP(w, r) {")
	fmt.Println("    fmt.Fprintf(w, \"Hello, World!\")")
	fmt.Println("}")

	// ===== http.HandlerFunc =====
	fmt.Println("\n🔹 http.HandlerFunc (функція → Handler)")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("// Звичайна функція")
	fmt.Println("func homeHandler(w ResponseWriter, r *Request) {")
	fmt.Println("    fmt.Fprintf(w, \"Home Page\")")
	fmt.Println("}")
	fmt.Println("")
	fmt.Println("// Перетворення в Handler")
	fmt.Println("http.Handle(\"/\", http.HandlerFunc(homeHandler))")

	// ===== Middleware Pattern =====
	fmt.Println("\n🔹 Middleware Pattern")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("func LoggingMiddleware(next Handler) Handler {")
	fmt.Println("    return HandlerFunc(func(w, r) {")
	fmt.Println("        log.Printf(\"Request: %s\", r.URL)")
	fmt.Println("        next.ServeHTTP(w, r)")
	fmt.Println("    })")
	fmt.Println("}")
	fmt.Println("")
	fmt.Println("Використання:")
	fmt.Println("handler := LoggingMiddleware(myHandler)")

	// ===== Router Example =====
	fmt.Println("\n🔹 Router (кастомний)")
	fmt.Println("─────────────────────────────────────────")

	router := NewRouter()
	router.Handle("/", HelloHandler{})
	router.Handle("/greet", GreetHandler{Name: "Іван"})
	router.Handle("/status", StatusHandler{})

	fmt.Println("Router створено з маршрутами:")
	fmt.Println("  GET / → HelloHandler")
	fmt.Println("  GET /greet → GreetHandler")
	fmt.Println("  GET /status → StatusHandler")

	// ===== Method Handler =====
	fmt.Println("\n🔹 Method Handler (різні HTTP методи)")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("methodHandler := MethodHandler{")
	fmt.Println("    GET:  http.HandlerFunc(getHandler),")
	fmt.Println("    POST: http.HandlerFunc(postHandler),")
	fmt.Println("}")

	// ===== Middleware Chain =====
	fmt.Println("\n🔹 Middleware Chain")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("Ланцюг middleware:")
	fmt.Println("  RecoveryMiddleware")
	fmt.Println("    ↓")
	fmt.Println("  LoggingMiddleware")
	fmt.Println("    ↓")
	fmt.Println("  AuthMiddleware")
	fmt.Println("    ↓")
	fmt.Println("  YourHandler")

	// ===== Real Example Setup =====
	fmt.Println("\n🔹 Приклад реального сервера")
	fmt.Println("─────────────────────────────────────────")

	// Створення handlers
	hello := HelloHandler{}
	counter := &CounterHandler{}
	greet := GreetHandler{Name: "Go Developer"}
	status := StatusHandler{}
	users := UserHandler{}

	// Middleware wrapper
	protectedHandler := AuthMiddleware(
		LoggingMiddleware(users),
	)

	fmt.Println("\nМаршрути:")
	fmt.Println("  /              → Hello World")
	fmt.Println("  /counter       → Counter (з state)")
	fmt.Println("  /greet         → Greeting")
	fmt.Println("  /status        → JSON status")
	fmt.Println("  /api/users     → Users API (захищено)")

	fmt.Println("\nMiddleware:")
	fmt.Println("  Recovery   → ловить паніки")
	fmt.Println("  Logging    → логує запити")
	fmt.Println("  Auth       → перевіряє токен")

	// ===== Code Example =====
	fmt.Println("\n📝 Код для запуску сервера:")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println(`
// Реєстрація handlers
http.Handle("/", RecoveryMiddleware(hello))
http.Handle("/counter", counter)
http.Handle("/greet", greet)
http.Handle("/status", status)
http.Handle("/api/users", protectedHandler)

// Запуск сервера
log.Println("Server starting on :8080")
log.Fatal(http.ListenAndServe(":8080", nil))
	`)

	// ===== Testing Example =====
	fmt.Println("\n🔹 Тестування (без запуску сервера)")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println(`
import (
    "net/http/httptest"
    "testing"
)

func TestHandler(t *testing.T) {
    req := httptest.NewRequest("GET", "/", nil)
    w := httptest.NewRecorder()
    
    handler := HelloHandler{}
    handler.ServeHTTP(w, req)
    
    if w.Code != http.StatusOK {
        t.Errorf("expected 200, got %d", w.Code)
    }
}
	`)

	// ===== Summary =====
	fmt.Println("\n📝 ВИСНОВКИ")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("✅ http.Handler - основа HTTP в Go")
	fmt.Println("   • Один метод: ServeHTTP(w, r)")
	fmt.Println()
	fmt.Println("💡 Patterns:")
	fmt.Println("   • Handler struct - state full handlers")
	fmt.Println("   • HandlerFunc - звичайні функції")
	fmt.Println("   • Middleware - обгортання handlers")
	fmt.Println("   • Router - маршрутизація")
	fmt.Println()
	fmt.Println("🔗 Middleware Chain:")
	fmt.Println("   func(next Handler) Handler")
	fmt.Println("   • Logging")
	fmt.Println("   • Authentication")
	fmt.Println("   • Recovery")
	fmt.Println("   • CORS")
	fmt.Println()
	fmt.Println("⚡ Переваги:")
	fmt.Println("   • Простий інтерфейс")
	fmt.Println("   • Композиція через middleware")
	fmt.Println("   • Легко тестувати")
	fmt.Println("   • Стандарт екосистеми")

	fmt.Println("\n\n════════════════════════════════════════════")
	fmt.Println("Для запуску сервера розкоментуйте:")
	fmt.Println("════════════════════════════════════════════")

	// Uncomment to run actual server:
	/*
		http.Handle("/", RecoveryMiddleware(hello))
		http.Handle("/counter", counter)
		http.Handle("/greet", greet)
		http.Handle("/status", status)
		http.Handle("/api/users", protectedHandler)

		fmt.Println("\nServer starting on http://localhost:8080")
		fmt.Println("Try:")
		fmt.Println("  curl http://localhost:8080/")
		fmt.Println("  curl http://localhost:8080/counter")
		fmt.Println("  curl http://localhost:8080/status")
		fmt.Println("  curl -H 'Authorization: Bearer secret-token' http://localhost:8080/api/users")
		fmt.Println()

		log.Fatal(http.ListenAndServe(":8080", nil))
	*/

	_ = hello
	_ = counter
	_ = greet
	_ = status
	_ = users
	_ = protectedHandler
}
