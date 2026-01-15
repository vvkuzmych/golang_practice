package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ============= Models =============

// APIUser представляє користувача в HTTP API
type APIUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Post struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Title  string `json:"title"`
}

// ============= APIDatabase Layer =============

type APIDatabase struct {
	// БЕЗ context! Context передається як параметр
}

func NewAPIDatabase() *APIDatabase {
	return &APIDatabase{}
}

func (db *APIDatabase) QueryUser(ctx context.Context, id int) (*APIUser, error) {
	// Симуляція повільного DB query
	resultChan := make(chan *APIUser, 1)
	errChan := make(chan error, 1)

	go func() {
		// Різна затримка для різних ID
		var delay time.Duration
		switch id {
		case 1:
			delay = 500 * time.Millisecond // Швидкий
		case 2:
			delay = 4 * time.Second // Повільний (timeout!)
		case 3:
			delay = 2 * time.Second // Середній
		default:
			delay = 1 * time.Second
		}

		time.Sleep(delay)

		// Перевірка cancellation перед поверненням
		select {
		case <-ctx.Done():
			errChan <- ctx.Err()
			return
		default:
			resultChan <- &APIUser{
				ID:       id,
				Username: fmt.Sprintf("user%d", id),
				Email:    fmt.Sprintf("user%d@example.com", id),
			}
		}
	}()

	// Чекаємо результат або cancellation
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("query cancelled: %w", ctx.Err())
	case err := <-errChan:
		return nil, err
	case user := <-resultChan:
		log.Printf("✓ Database query completed for user %d", id)
		return user, nil
	}
}

// ============= External API Layer =============

type ExternalAPI struct {
	baseURL string
}

func NewExternalAPI() *ExternalAPI {
	return &ExternalAPI{baseURL: "https://api.example.com"}
}

func (api *ExternalAPI) FetchUserPosts(ctx context.Context, userID int) ([]Post, error) {
	// Симуляція API call з затримкою
	resultChan := make(chan []Post, 1)

	go func() {
		time.Sleep(300 * time.Millisecond)

		select {
		case <-ctx.Done():
			return
		default:
			posts := []Post{
				{ID: 1, UserID: userID, Title: "Post 1"},
				{ID: 2, UserID: userID, Title: "Post 2"},
			}
			resultChan <- posts
		}
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("API call cancelled: %w", ctx.Err())
	case posts := <-resultChan:
		log.Printf("✓ External API call completed for user %d", userID)
		return posts, nil
	}
}

// ============= Service Layer (HTTP Context) =============

type HTTPUserService struct {
	db  *APIDatabase // БЕЗ context!
	api *ExternalAPI
}

func NewHTTPUserService(db *APIDatabase, api *ExternalAPI) *HTTPUserService {
	return &HTTPUserService{
		db:  db,
		api: api,
	}
}

// ✅ Context як параметр!
func (s *HTTPUserService) GetUser(ctx context.Context, id int) (*APIUser, error) {
	// Перевірка перед тяжкою операцією
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	log.Printf("→ Service: Getting user %d", id)
	start := time.Now()

	user, err := s.db.QueryUser(ctx, id)
	if err != nil {
		log.Printf("✗ Service: Failed to get user %d (%v)", id, time.Since(start))
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	log.Printf("✓ Service: Got user %d (%v)", id, time.Since(start))
	return user, nil
}

func (s *HTTPUserService) GetUserWithPosts(ctx context.Context, id int) (*APIUser, []Post, error) {
	user, err := s.GetUser(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	posts, err := s.api.FetchUserPosts(ctx, id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch posts: %w", err)
	}

	return user, posts, nil
}

// ============= HTTP Handlers =============

type Handler struct {
	service *HTTPUserService
}

func NewHandler(service *HTTPUserService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	// 1. Context від request
	ctx := r.Context()

	// 2. Додаємо timeout
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// 3. Витягуємо ID
	id, err := extractID(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	log.Printf("→ HTTP: GET /users/%d (timeout: 3s)", id)
	start := time.Now()

	// 4. Викликаємо service з context
	user, err := h.service.GetUser(ctx, id)
	duration := time.Since(start)

	if err != nil {
		h.handleError(w, err, duration)
		return
	}

	// 5. Успішна відповідь
	log.Printf("✓ HTTP: Request completed (%v)", duration)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) HandleGetUserPosts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	id, err := extractID(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	log.Printf("→ HTTP: GET /users/%d/posts (timeout: 2s)", id)

	user, posts, err := h.service.GetUserWithPosts(ctx, id)
	if err != nil {
		h.handleError(w, err, 0)
		return
	}

	response := map[string]interface{}{
		"user":  user,
		"posts": posts,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) handleError(w http.ResponseWriter, err error, duration time.Duration) {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		log.Printf("✗ HTTP: Request timeout (%v)", duration)
		http.Error(w, "Request timeout", http.StatusGatewayTimeout)
	case errors.Is(err, context.Canceled):
		log.Printf("✗ HTTP: Request cancelled (%v)", duration)
		http.Error(w, "Request cancelled", http.StatusRequestTimeout)
	default:
		log.Printf("✗ HTTP: Error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// ============= Helper Functions =============

func extractID(path string) (int, error) {
	// /users/123 or /users/123/posts
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		return 0, errors.New("invalid path")
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}

	return id, nil
}

// ============= Main =============

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║   HTTP Service with Context Timeout      ║")
	fmt.Println("╚══════════════════════════════════════════╝\n")

	// Setup
	db := NewAPIDatabase()
	api := NewExternalAPI()
	service := NewHTTPUserService(db, api)
	handler := NewHandler(service)

	// Routes
	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/posts") {
			handler.HandleGetUserPosts(w, r)
		} else {
			handler.HandleGetUser(w, r)
		}
	})

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	// Demo mode: симуляція різних scenarios
	fmt.Println("🚀 Server starting in DEMO mode")
	fmt.Println("   Demonstrating different timeout scenarios...")
	fmt.Println()

	go startServer()

	// Дамо серверу час запуститись
	time.Sleep(100 * time.Millisecond)

	// Симуляція клієнтських запитів
	runDemoRequests()

	// Чекаємо трохи перед завершенням
	time.Sleep(1 * time.Second)

	fmt.Println("\n✅ Demo completed!")
	fmt.Println("\nTo run as real server:")
	fmt.Println("  1. Comment out demo code in main()")
	fmt.Println("  2. Uncomment: select {}")
	fmt.Println("  3. Run: go run solution_3.go")
	fmt.Println("  4. Test: curl http://localhost:8080/users/1")
}

func startServer() {
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func runDemoRequests() {
	// Scenario 1: Fast request (успіх)
	fmt.Println("🔹 Scenario 1: Fast Request (Success)")
	fmt.Println("─────────────────────────────────────────")
	makeRequest("http://localhost:8080/users/1")

	time.Sleep(500 * time.Millisecond)

	// Scenario 2: Slow request (timeout)
	fmt.Println("\n🔹 Scenario 2: Slow Request (Timeout)")
	fmt.Println("─────────────────────────────────────────")
	makeRequest("http://localhost:8080/users/2")

	time.Sleep(500 * time.Millisecond)

	// Scenario 3: Medium request (успіх)
	fmt.Println("\n🔹 Scenario 3: Medium Request (Success)")
	fmt.Println("─────────────────────────────────────────")
	makeRequest("http://localhost:8080/users/3")
}

func makeRequest(url string) {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("❌ Client error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("📊 Response: %d %s\n", resp.StatusCode, resp.Status)

	if resp.StatusCode == http.StatusOK {
		var user APIUser
		json.NewDecoder(resp.Body).Decode(&user)
		fmt.Printf("   User: %+v\n", user)
	}
}
