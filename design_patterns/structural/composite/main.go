package main

import (
	"fmt"
	"strings"
)

// Component - базовий інтерфейс для всіх компонентів
type Component interface {
	Operation() string
	Add(Component) error
	Remove(Component) error
	GetChild(int) (Component, error)
}

// Leaf - листок дерева (не має дітей)
type Leaf struct {
	name string
}

func NewLeaf(name string) *Leaf {
	return &Leaf{name: name}
}

func (l *Leaf) Operation() string {
	return l.name
}

func (l *Leaf) Add(c Component) error {
	return fmt.Errorf("cannot add to a leaf")
}

func (l *Leaf) Remove(c Component) error {
	return fmt.Errorf("cannot remove from a leaf")
}

func (l *Leaf) GetChild(index int) (Component, error) {
	return nil, fmt.Errorf("leaf has no children")
}

// Composite - контейнер (може мати дітей)
type Composite struct {
	name     string
	children []Component
}

func NewComposite(name string) *Composite {
	return &Composite{
		name:     name,
		children: make([]Component, 0),
	}
}

func (c *Composite) Operation() string {
	var builder strings.Builder
	builder.WriteString(c.name)
	builder.WriteString(" [")

	for i, child := range c.children {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(child.Operation())
	}

	builder.WriteString("]")
	return builder.String()
}

func (c *Composite) Add(component Component) error {
	c.children = append(c.children, component)
	return nil
}

func (c *Composite) Remove(component Component) error {
	for i, child := range c.children {
		if child == component {
			c.children = append(c.children[:i], c.children[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("component not found")
}

func (c *Composite) GetChild(index int) (Component, error) {
	if index < 0 || index >= len(c.children) {
		return nil, fmt.Errorf("index out of range")
	}
	return c.children[index], nil
}

// FileSystem Example - більш практичний приклад

type FileSystemNode interface {
	GetName() string
	GetSize() int
	Print(indent string)
}

// File - листок
type File struct {
	name string
	size int
}

func NewFile(name string, size int) *File {
	return &File{name: name, size: size}
}

func (f *File) GetName() string {
	return f.name
}

func (f *File) GetSize() int {
	return f.size
}

func (f *File) Print(indent string) {
	fmt.Printf("%s📄 %s (%d KB)\n", indent, f.name, f.size)
}

// Directory - composite
type Directory struct {
	name  string
	nodes []FileSystemNode
}

func NewDirectory(name string) *Directory {
	return &Directory{
		name:  name,
		nodes: make([]FileSystemNode, 0),
	}
}

func (d *Directory) GetName() string {
	return d.name
}

func (d *Directory) GetSize() int {
	total := 0
	for _, node := range d.nodes {
		total += node.GetSize()
	}
	return total
}

func (d *Directory) Add(node FileSystemNode) {
	d.nodes = append(d.nodes, node)
}

func (d *Directory) Print(indent string) {
	fmt.Printf("%s📁 %s (%d KB)\n", indent, d.name, d.GetSize())
	for _, node := range d.nodes {
		node.Print(indent + "  ")
	}
}

// Organization Example - компанія з відділами

type Employee interface {
	GetName() string
	GetSalary() int
	Print(indent string)
}

// Developer - листок
type Developer struct {
	name   string
	salary int
}

func NewDeveloper(name string, salary int) *Developer {
	return &Developer{name: name, salary: salary}
}

func (d *Developer) GetName() string {
	return d.name
}

func (d *Developer) GetSalary() int {
	return d.salary
}

func (d *Developer) Print(indent string) {
	fmt.Printf("%s👨‍💻 %s (Developer, $%d)\n", indent, d.name, d.salary)
}

// Manager - composite
type Manager struct {
	name         string
	salary       int
	subordinates []Employee
}

func NewManager(name string, salary int) *Manager {
	return &Manager{
		name:         name,
		salary:       salary,
		subordinates: make([]Employee, 0),
	}
}

func (m *Manager) GetName() string {
	return m.name
}

func (m *Manager) GetSalary() int {
	total := m.salary
	for _, sub := range m.subordinates {
		total += sub.GetSalary()
	}
	return total
}

func (m *Manager) Add(employee Employee) {
	m.subordinates = append(m.subordinates, employee)
}

func (m *Manager) Print(indent string) {
	fmt.Printf("%s👔 %s (Manager, $%d, manages %d people)\n",
		indent, m.name, m.salary, len(m.subordinates))
	for _, sub := range m.subordinates {
		sub.Print(indent + "  ")
	}
}

func main() {
	fmt.Println("=== 1. Basic Composite Pattern ===")

	// Створення листків
	leaf1 := NewLeaf("Leaf 1")
	leaf2 := NewLeaf("Leaf 2")
	leaf3 := NewLeaf("Leaf 3")

	// Створення composite
	comp1 := NewComposite("Composite 1")
	comp1.Add(leaf1)
	comp1.Add(leaf2)

	comp2 := NewComposite("Composite 2")
	comp2.Add(leaf3)
	comp2.Add(comp1)

	// Виконання операції
	fmt.Println(comp2.Operation())

	fmt.Println("\n=== 2. File System Example ===")

	// Створення файлової системи
	root := NewDirectory("root")

	// Documents директорія
	documents := NewDirectory("documents")
	documents.Add(NewFile("report.pdf", 150))
	documents.Add(NewFile("presentation.pptx", 500))

	// Photos директорія
	photos := NewDirectory("photos")
	photos.Add(NewFile("vacation1.jpg", 1200))
	photos.Add(NewFile("vacation2.jpg", 1300))

	// Projects директорія з вкладеними
	projects := NewDirectory("projects")
	projectA := NewDirectory("project-a")
	projectA.Add(NewFile("main.go", 10))
	projectA.Add(NewFile("README.md", 5))
	projects.Add(projectA)

	// Додаємо все до root
	root.Add(documents)
	root.Add(photos)
	root.Add(projects)
	root.Add(NewFile("config.yaml", 2))

	// Виводимо структуру
	root.Print("")

	fmt.Printf("\nTotal size: %d KB\n", root.GetSize())

	fmt.Println("\n=== 3. Organization Example ===")

	// CEO
	ceo := NewManager("John Smith", 10000)

	// CTO і його команда
	cto := NewManager("Jane Doe", 8000)
	cto.Add(NewDeveloper("Bob Wilson", 5000))
	cto.Add(NewDeveloper("Alice Brown", 5500))
	cto.Add(NewDeveloper("Charlie Davis", 4800))

	// CFO і його команда
	cfo := NewManager("Mike Johnson", 7500)
	cfo.Add(NewDeveloper("Emily White", 4500))
	cfo.Add(NewDeveloper("David Lee", 4300))

	// CEO керує CTO і CFO
	ceo.Add(cto)
	ceo.Add(cfo)

	// Виводимо структуру
	ceo.Print("")

	fmt.Printf("\nTotal company salary: $%d\n", ceo.GetSalary())
}
