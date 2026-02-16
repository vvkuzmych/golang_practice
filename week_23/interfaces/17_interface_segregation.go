package main

import "fmt"

// 17. Interface Segregation - Prefer small, focused interfaces

type Worker interface {
	Work()
}

type Eater interface {
	Eat()
}

type Human struct{ Name string }

func (h Human) Work() { fmt.Printf("%s is working\n", h.Name) }
func (h Human) Eat()  { fmt.Printf("%s is eating\n", h.Name) }

type Robot struct{ ID string }

func (r Robot) Work() { fmt.Printf("Robot %s is working\n", r.ID) }

func StartShift(w Worker) {
	w.Work()
}

func main() {
	human := Human{Name: "Alice"}
	robot := Robot{ID: "R-001"}
	StartShift(human)
	StartShift(robot)
	human.Eat()
}
