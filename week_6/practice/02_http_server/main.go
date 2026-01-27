package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserStore struct {
	mu     sync.RWMutex
	users  map[int]User
	nextID int
}

func NewUserStore() *UserStore {
	return &UserStore{
		users:  make(map[int]User),
		nextID: 1,
	}
}

func (s *UserStore) Create(name, email string) User {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := User{
		ID:    s.nextID,
		Name:  name,
		Email: email,
	}
	s.users[s.nextID] = user
	s.nextID++
	return user
}

func (s *UserStore) GetAll() []User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
}

func (s *UserStore) GetByID(id int) (User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	return user, ok
}

type Server struct {
	store *UserStore
}

func NewServer() *Server {
	return &Server{
		store: NewUserStore(),
	}
}

func (s *Server) handleUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		users := s.store.GetAll()
		json.NewEncoder(w).Encode(users)

	case http.MethodPost:
		var req struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user := s.store.Create(req.Name, req.Email)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	user, ok := s.store.GetByID(id)
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// Middleware для логування
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		next(w, r)
	}
}

func main() {
	server := NewServer()

	// Seed data
	server.store.Create("John Doe", "john@example.com")
	server.store.Create("Jane Smith", "jane@example.com")

	// Routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to User API!\n")
		fmt.Fprintf(w, "\nEndpoints:\n")
		fmt.Fprintf(w, "GET    /api/users       - Get all users\n")
		fmt.Fprintf(w, "POST   /api/users       - Create user\n")
		fmt.Fprintf(w, "GET    /api/users/{id}  - Get user by ID\n")
	})

	http.HandleFunc("/api/users", loggingMiddleware(server.handleUsers))
	http.HandleFunc("/api/users/", loggingMiddleware(server.handleUserByID))

	fmt.Println("🚀 Server started at http://localhost:8080")
	fmt.Println("Try: curl http://localhost:8080/api/users")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
