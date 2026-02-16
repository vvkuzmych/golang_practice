package main

import "fmt"

type Reader interface{ Read() string }
type Writer interface{ Write(string) }
type ReadWriter interface {
	Reader
	Writer
}
type File struct{ data string }

func (f *File) Read() string   { return f.data }
func (f *File) Write(s string) { f.data = s }
func main() {
	var rw ReadWriter = &File{}
	rw.Write("Hello")
	fmt.Println(rw.Read())
}
