package main

import "fmt"

type Iterator interface {
	HasNext() bool
	Next() interface{}
}

type Collection interface {
	CreateIterator() Iterator
}

// Book collection
type Book struct {
	Title  string
	Author string
}

type BookCollection struct {
	books []Book
}

func (bc *BookCollection) CreateIterator() Iterator {
	return &BookIterator{
		collection: bc,
		index:      0,
	}
}

type BookIterator struct {
	collection *BookCollection
	index      int
}

func (bi *BookIterator) HasNext() bool {
	return bi.index < len(bi.collection.books)
}

func (bi *BookIterator) Next() interface{} {
	if bi.HasNext() {
		book := bi.collection.books[bi.index]
		bi.index++
		return book
	}
	return nil
}

func main() {
	fmt.Println("=== Iterator Pattern ===\n")

	collection := &BookCollection{
		books: []Book{
			{"The Go Programming Language", "Donovan & Kernighan"},
			{"Clean Code", "Robert Martin"},
			{"Design Patterns", "Gang of Four"},
		},
	}

	iterator := collection.CreateIterator()

	for iterator.HasNext() {
		book := iterator.Next().(Book)
		fmt.Printf("📚 %s by %s\n", book.Title, book.Author)
	}
}
