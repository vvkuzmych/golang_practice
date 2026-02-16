package main

import "fmt"

// 09. Interface Composition - Композиція інтерфейсів

type Reader interface {
	Read() string
}

type Writer interface {
	Write(string)
}

// Композиція: ReadWriter = Reader + Writer
type ReadWriter interface {
	Reader
	Writer
}

type File struct {
	content string
}

func (f *File) Read() string {
	return f.content
}

func (f *File) Write(s string) {
	f.content = s
}

// Функція приймає ReadWriter
func process(rw ReadWriter) {
	rw.Write("Hello from composition")
	fmt.Println("Read:", rw.Read())
}

func main() {
	file := &File{}

	// File реалізує обидва інтерфейси
	var r Reader = file
	var w Writer = file
	var rw ReadWriter = file

	w.Write("test data")
	fmt.Println("Reader:", r.Read())

	process(rw)
}
