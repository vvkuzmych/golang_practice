package main

import "fmt"

func main() {
	var s []int
	fmt.Println(s, len(s), cap(s))
	s = []int{1, 2, 3}
	fmt.Println(s)
	s = make([]int, 5)
	fmt.Println(s)
}
