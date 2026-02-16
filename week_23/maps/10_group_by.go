package main

import "fmt"

// 10. Group By - Group slice elements by key

type Person struct {
	Name string
	Dept string
}

func main() {
	people := []Person{
		{"Alice", "Engineering"},
		{"Bob", "Sales"},
		{"Carol", "Engineering"},
		{"Dave", "Sales"},
	}

	// Group by department
	byDept := make(map[string][]Person)
	for _, p := range people {
		byDept[p.Dept] = append(byDept[p.Dept], p)
	}
	fmt.Println("By department:", byDept)

	// Group by first letter (use string for readable keys)
	words := []string{"apple", "banana", "apricot", "berry"}
	byLetter := make(map[string][]string)
	for _, w := range words {
		key := string(w[0])
		byLetter[key] = append(byLetter[key], w)
	}
	fmt.Println("By first letter:", byLetter)
}
