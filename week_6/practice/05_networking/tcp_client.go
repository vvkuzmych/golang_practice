package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Connected to server at localhost:9000")
	fmt.Println("Type messages (Ctrl+C to exit):")

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		message, _ := reader.ReadString('\n')

		conn.Write([]byte(message))

		response, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Printf("Server: %s", response)
	}
}
