package main

import "fmt"

// 20. Interface Best Practices - Summary of Go interface idioms

type Storer interface {
	Store(data string)
}

type FileStore struct{}

func (f *FileStore) Store(data string) {
	fmt.Println("Stored to file:", data)
}

func processStore(s Storer) {
	s.Store("data")
}

type UserFetcher interface {
	Fetch(id int) string
}

func handleRequest(f UserFetcher, id int) {
	fmt.Println("User:", f.Fetch(id))
}

type DBUserService struct{}

func (d *DBUserService) Fetch(id int) string {
	return fmt.Sprintf("User-%d", id)
}

func printAny[T any](v T) {
	fmt.Println(v)
}

func main() {
	processStore(&FileStore{})
	handleRequest(&DBUserService{}, 42)
	printAny("hello")
	printAny(42)
	printAny(3.14)
}
