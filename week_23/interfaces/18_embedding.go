package main

import "fmt"

// 18. Interface Embedding - Compose interfaces from smaller ones

type Reader interface {
	Read() string
}

type Writer interface {
	Write(s string)
}

type ReadWriter interface {
	Reader
	Writer
}

type BufferedRW struct {
	buffer string
}

func (b *BufferedRW) Read() string {
	return b.buffer
}

func (b *BufferedRW) Write(s string) {
	b.buffer = s
}

func main() {
	buf := &BufferedRW{}
	buf.Write("Hello, World!")

	var rw ReadWriter = buf
	fmt.Println("Read:", rw.Read())

	var r Reader = buf
	fmt.Println("Read via Reader:", r.Read())
}
