package main

import "fmt"

// 07. Iteration - for loops and range over slices

func main() {
	fruits := []string{"apple", "banana", "cherry"}

	// Classic for loop with index
	fmt.Print("Index loop: ")
	for i := 0; i < len(fruits); i++ {
		fmt.Print(fruits[i], " ")
	}
	fmt.Println()

	// Range - index and value
	fmt.Print("Range index+value: ")
	for i, f := range fruits {
		fmt.Printf("%d:%s ", i, f)
	}
	fmt.Println()

	// Range - value only
	fmt.Print("Range value only: ")
	for _, f := range fruits {
		fmt.Print(f, " ")
	}
	fmt.Println()

	// Range - index only
	fmt.Print("Range index only: ")
	for i := range fruits {
		fmt.Print(i, " ")
	}
	fmt.Println()

	// Modify copy - range gives value copy
	for i, v := range fruits {
		v = "x" // doesn't affect fruits
		_ = v
		fruits[i] = fruits[i] + "!" // modifies original
	}
	fmt.Println("Modified:", fruits)
}
