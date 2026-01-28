package main

import "fmt"

// TreeType - flyweight (shared state)
type TreeType struct {
	Name    string
	Color   string
	Texture string
}

// TreeFactory - flyweight factory
type TreeFactory struct {
	types map[string]*TreeType
}

func NewTreeFactory() *TreeFactory {
	return &TreeFactory{types: make(map[string]*TreeType)}
}

func (tf *TreeFactory) GetTreeType(name, color, texture string) *TreeType {
	key := name + "_" + color + "_" + texture
	if _, ok := tf.types[key]; !ok {
		fmt.Printf("Creating new TreeType: %s\n", key)
		tf.types[key] = &TreeType{Name: name, Color: color, Texture: texture}
	}
	return tf.types[key]
}

// Tree - uses flyweight (unique state: x, y)
type Tree struct {
	X, Y int
	Type *TreeType
}

func (t *Tree) Draw() {
	fmt.Printf("Drawing %s tree (color: %s) at (%d, %d)\n",
		t.Type.Name, t.Type.Color, t.X, t.Y)
}

// Forest
type Forest struct {
	trees   []*Tree
	factory *TreeFactory
}

func NewForest() *Forest {
	return &Forest{
		trees:   make([]*Tree, 0),
		factory: NewTreeFactory(),
	}
}

func (f *Forest) PlantTree(x, y int, name, color, texture string) {
	treeType := f.factory.GetTreeType(name, color, texture)
	tree := &Tree{X: x, Y: y, Type: treeType}
	f.trees = append(f.trees, tree)
}

func (f *Forest) Draw() {
	for _, tree := range f.trees {
		tree.Draw()
	}
}

func main() {
	fmt.Println("=== Flyweight Pattern ===\n")

	forest := NewForest()

	// Садимо багато дерев (багато об'єктів використовують мало flyweight-ів)
	forest.PlantTree(1, 2, "Oak", "Green", "Rough")
	forest.PlantTree(5, 3, "Oak", "Green", "Rough") // Reuses flyweight!
	forest.PlantTree(10, 8, "Pine", "Dark", "Smooth")
	forest.PlantTree(3, 7, "Oak", "Green", "Rough")    // Reuses flyweight!
	forest.PlantTree(15, 12, "Pine", "Dark", "Smooth") // Reuses flyweight!

	fmt.Println("\nDrawing forest:")
	forest.Draw()

	fmt.Printf("\nTotal tree types created: %d\n", len(forest.factory.types))
	fmt.Printf("Total trees planted: %d\n", len(forest.trees))
}
