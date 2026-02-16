package main

import "fmt"

// 11. Dependency Injection - Ін'єкція залежностей

type Database interface {
	Query(sql string) string
}

type PostgresDB struct{}

func (p PostgresDB) Query(sql string) string {
	return "Postgres result: " + sql
}

type MockDB struct{}

func (m MockDB) Query(sql string) string {
	return "Mock result: " + sql
}

type UserService struct {
	db Database
}

func NewUserService(db Database) *UserService {
	return &UserService{db: db}
}

func (u *UserService) GetUser(id int) string {
	return u.db.Query(fmt.Sprintf("SELECT * FROM users WHERE id=%d", id))
}

func main() {
	// Production
	prodService := NewUserService(PostgresDB{})
	fmt.Println(prodService.GetUser(1))

	// Testing
	testService := NewUserService(MockDB{})
	fmt.Println(testService.GetUser(1))
}
