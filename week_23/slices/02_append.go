package main

import "fmt"

func main() {
	s := []int{1, 2, 3}
	s = append(s, 4, 5)
	fmt.Println(s)
}
