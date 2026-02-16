package main

import (
	"bytes"
	"fmt"
	"io"
)

func main() {
	var b bytes.Buffer
	b.Write([]byte("Hello "))
	fmt.Fprintf(&b, "World!")
	io.Copy(io.Discard, &b)
	fmt.Println("Done")
}
