package main

import "fmt"

func main() {
	src := []int{1, 2, 3}
	dst := make([]int, len(src))
	copy(dst, src)
	fmt.Println(dst)
}
