package main

import "fmt"

// 19. Generic Interface - Interfaces with type parameters (Go 1.18+)

type Comparable interface {
	~int | ~string | ~float64
}

type Container[T any] interface {
	Add(item T)
	Get(index int) T
	Len() int
}

type SliceContainer[T any] struct {
	items []T
}

func (c *SliceContainer[T]) Add(item T) {
	c.items = append(c.items, item)
}

func (c *SliceContainer[T]) Get(index int) T {
	return c.items[index]
}

func (c *SliceContainer[T]) Len() int {
	return len(c.items)
}

type Maxer[T Comparable] interface {
	Max(other T) T
}

type OrderedInt int

func (a OrderedInt) Max(b OrderedInt) OrderedInt {
	if a > b {
		return a
	}
	return b
}

func main() {
	strContainer := &SliceContainer[string]{}
	strContainer.Add("hello")
	strContainer.Add("world")
	fmt.Println("Strings:", strContainer.Get(0), strContainer.Get(1))

	intContainer := &SliceContainer[int]{}
	intContainer.Add(10)
	intContainer.Add(20)
	fmt.Println("Ints:", intContainer.Get(0), intContainer.Get(1))

	var m Maxer[OrderedInt] = OrderedInt(5)
	fmt.Println("Max(5, 3):", m.Max(3))
}
