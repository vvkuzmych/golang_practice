package main

import (
	"fmt"
	"sort"
)

// 15. Sort Interface - Implementing sort.Interface for custom types

type Person struct {
	Name string
	Age  int
}

type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ByName []Person

func (a ByName) Len() int           { return len(a) }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func main() {
	people := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 35},
	}
	fmt.Println("Original:", people)
	sort.Sort(ByAge(people))
	fmt.Println("Sorted by age:", people)
	sort.Sort(ByName(people))
	fmt.Println("Sorted by name:", people)

	numbers := []int{3, 1, 4, 1, 5}
	sort.Slice(numbers, func(i, j int) bool { return numbers[i] < numbers[j] })
	fmt.Println("Sorted numbers:", numbers)
}
