package main

import "fmt"

type Visitor interface {
	VisitCircle(*Circle)
	VisitRectangle(*Rectangle)
}

type Shape interface {
	Accept(Visitor)
}

type Circle struct {
	Radius float64
}

func (c *Circle) Accept(v Visitor) {
	v.VisitCircle(c)
}

type Rectangle struct {
	Width, Height float64
}

func (r *Rectangle) Accept(v Visitor) {
	v.VisitRectangle(r)
}

type AreaCalculator struct {
	Area float64
}

func (a *AreaCalculator) VisitCircle(c *Circle) {
	a.Area = 3.14 * c.Radius * c.Radius
	fmt.Printf("Circle area: %.2f\n", a.Area)
}

func (a *AreaCalculator) VisitRectangle(r *Rectangle) {
	a.Area = r.Width * r.Height
	fmt.Printf("Rectangle area: %.2f\n", a.Area)
}

func main() {
	fmt.Println("=== Visitor Pattern ===\n")

	shapes := []Shape{
		&Circle{Radius: 5},
		&Rectangle{Width: 10, Height: 20},
	}

	calculator := &AreaCalculator{}

	for _, shape := range shapes {
		shape.Accept(calculator)
	}
}
