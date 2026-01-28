package main

import "fmt"

// Handler - інтерфейс обробника
type Handler interface {
	SetNext(Handler) Handler
	Handle(request string) string
}

// BaseHandler - базова реалізація
type BaseHandler struct {
	next Handler
}

func (h *BaseHandler) SetNext(handler Handler) Handler {
	h.next = handler
	return handler
}

func (h *BaseHandler) Handle(request string) string {
	if h.next != nil {
		return h.next.Handle(request)
	}
	return ""
}

// AuthHandler - перевірка автентифікації
type AuthHandler struct {
	BaseHandler
}

func (h *AuthHandler) Handle(request string) string {
	if request == "auth-fail" {
		return "AuthHandler: Authentication failed!"
	}
	fmt.Println("AuthHandler: User authenticated ✓")
	return h.BaseHandler.Handle(request)
}

// LoggingHandler - логування
type LoggingHandler struct {
	BaseHandler
}

func (h *LoggingHandler) Handle(request string) string {
	fmt.Printf("LoggingHandler: Logging request '%s' ✓\n", request)
	return h.BaseHandler.Handle(request)
}

// ValidationHandler - валідація
type ValidationHandler struct {
	BaseHandler
}

func (h *ValidationHandler) Handle(request string) string {
	if request == "invalid" {
		return "ValidationHandler: Invalid request!"
	}
	fmt.Println("ValidationHandler: Request validated ✓")
	return h.BaseHandler.Handle(request)
}

// ProcessingHandler - обробка
type ProcessingHandler struct {
	BaseHandler
}

func (h *ProcessingHandler) Handle(request string) string {
	fmt.Println("ProcessingHandler: Processing request ✓")
	return "Success: Request processed!"
}

// HTTP Middleware Example

type HTTPRequest struct {
	Method string
	Path   string
	User   string
}

type HTTPHandler interface {
	SetNext(HTTPHandler) HTTPHandler
	Handle(*HTTPRequest) string
}

type BaseHTTPHandler struct {
	next HTTPHandler
}

func (h *BaseHTTPHandler) SetNext(handler HTTPHandler) HTTPHandler {
	h.next = handler
	return handler
}

func (h *BaseHTTPHandler) Handle(req *HTTPRequest) string {
	if h.next != nil {
		return h.next.Handle(req)
	}
	return ""
}

type RateLimitHandler struct {
	BaseHTTPHandler
	requests map[string]int
	limit    int
}

func NewRateLimitHandler(limit int) *RateLimitHandler {
	return &RateLimitHandler{
		requests: make(map[string]int),
		limit:    limit,
	}
}

func (h *RateLimitHandler) Handle(req *HTTPRequest) string {
	h.requests[req.User]++
	if h.requests[req.User] > h.limit {
		return "RateLimitHandler: Too many requests!"
	}
	fmt.Printf("RateLimitHandler: Request %d/%d ✓\n", h.requests[req.User], h.limit)
	return h.BaseHTTPHandler.Handle(req)
}

type CORSHandler struct {
	BaseHTTPHandler
}

func (h *CORSHandler) Handle(req *HTTPRequest) string {
	fmt.Println("CORSHandler: CORS headers added ✓")
	return h.BaseHTTPHandler.Handle(req)
}

type RouterHandler struct {
	BaseHTTPHandler
}

func (h *RouterHandler) Handle(req *HTTPRequest) string {
	fmt.Printf("RouterHandler: Routing %s %s ✓\n", req.Method, req.Path)
	return "Success: Request handled!"
}

func main() {
	fmt.Println("=== 1. Basic Chain of Responsibility ===\n")

	// Створення ланцюга
	auth := &AuthHandler{}
	logging := &LoggingHandler{}
	validation := &ValidationHandler{}
	processing := &ProcessingHandler{}

	auth.SetNext(logging).SetNext(validation).SetNext(processing)

	// Тест 1: Успішний запит
	fmt.Println("Request 1: valid data")
	result1 := auth.Handle("valid-data")
	fmt.Printf("Result: %s\n\n", result1)

	// Тест 2: Невалідний запит
	fmt.Println("Request 2: invalid data")
	result2 := auth.Handle("invalid")
	fmt.Printf("Result: %s\n\n", result2)

	// Тест 3: Не авторизований
	fmt.Println("Request 3: auth-fail")
	result3 := auth.Handle("auth-fail")
	fmt.Printf("Result: %s\n\n", result3)

	fmt.Println("=== 2. HTTP Middleware Chain ===\n")

	// Створення middleware ланцюга
	rateLimit := NewRateLimitHandler(3)
	cors := &CORSHandler{}
	router := &RouterHandler{}

	rateLimit.SetNext(cors).SetNext(router)

	// Тест: кілька запитів
	requests := []*HTTPRequest{
		{Method: "GET", Path: "/api/users", User: "user1"},
		{Method: "POST", Path: "/api/users", User: "user1"},
		{Method: "GET", Path: "/api/products", User: "user1"},
		{Method: "DELETE", Path: "/api/users/1", User: "user1"}, // Має бути відхилено
	}

	for i, req := range requests {
		fmt.Printf("Request %d: %s %s (User: %s)\n", i+1, req.Method, req.Path, req.User)
		result := rateLimit.Handle(req)
		fmt.Printf("Result: %s\n\n", result)
	}
}
