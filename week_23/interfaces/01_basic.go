package main

import "fmt"

// 01. Basic Interface
type Speaker interface {
	Speak() string
}
type Dog struct{ Name string }

func (d Dog) Speak() string { return "Woof!" }

type Cat struct{ Name string }

func (c Cat) Speak() string { return "Meow!" }
func main() {
	var s Speaker
	s = Dog{Name: "Buddy"}
	fmt.Println(s.Speak())
	s = Cat{Name: "Whiskers"}
	fmt.Println(s.Speak())
}
