package main

import "fmt"

type Shape interface {
	Area() float64
}
type Circle struct{ Radius float64 }

func (c Circle) Area() float64 { return 3.14 * c.Radius * c.Radius }

type Rectangle struct{ Width, Height float64 }

func (r Rectangle) Area() float64 { return r.Width * r.Height }
func printArea(s Shape) {
	fmt.Printf("Area: %.2f\n", s.Area())
}
func main() {
	printArea(Circle{5})
	printArea(Rectangle{3, 4})
}
