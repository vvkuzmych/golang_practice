package main

import "fmt"

// 16. Mock Testing - Interface enables dependency injection for testing

type UserRepository interface {
	GetByID(id int) (string, bool)
	Save(id int, name string)
}

type RealUserRepo struct {
	users map[int]string
}

func NewRealUserRepo() *RealUserRepo {
	return &RealUserRepo{users: map[int]string{1: "Alice", 2: "Bob"}}
}

func (r *RealUserRepo) GetByID(id int) (string, bool) {
	name, ok := r.users[id]
	return name, ok
}

func (r *RealUserRepo) Save(id int, name string) {
	r.users[id] = name
}

type MockUserRepo struct {
	GetByIDFunc func(id int) (string, bool)
}

func (m *MockUserRepo) GetByID(id int) (string, bool) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return "", false
}

func (m *MockUserRepo) Save(id int, name string) {}

type UserService struct {
	repo UserRepository
}

func (s *UserService) GetUserName(id int) string {
	name, ok := s.repo.GetByID(id)
	if !ok {
		return "Unknown"
	}
	return name
}

func main() {
	realRepo := NewRealUserRepo()
	svc := &UserService{repo: realRepo}
	fmt.Println("Real repo - User 1:", svc.GetUserName(1))

	mockRepo := &MockUserRepo{
		GetByIDFunc: func(id int) (string, bool) {
			return "MockUser", true
		},
	}
	svc.repo = mockRepo
	fmt.Println("Mock repo - User 99:", svc.GetUserName(99))
}
