package main

import "fmt"

type Memento struct {
	state string
}

type Editor struct {
	content string
}

func (e *Editor) Type(text string) {
	e.content += text
}

func (e *Editor) Save() *Memento {
	return &Memento{state: e.content}
}

func (e *Editor) Restore(m *Memento) {
	e.content = m.state
}

func (e *Editor) GetContent() string {
	return e.content
}

type History struct {
	mementos []*Memento
}

func (h *History) Push(m *Memento) {
	h.mementos = append(h.mementos, m)
}

func (h *History) Pop() *Memento {
	if len(h.mementos) == 0 {
		return nil
	}
	last := h.mementos[len(h.mementos)-1]
	h.mementos = h.mementos[:len(h.mementos)-1]
	return last
}

func main() {
	fmt.Println("=== Memento Pattern (Undo/Redo) ===\n")

	editor := &Editor{}
	history := &History{}

	editor.Type("Hello ")
	history.Push(editor.Save())

	editor.Type("World!")
	fmt.Printf("Current: %s\n", editor.GetContent())

	// Undo
	editor.Restore(history.Pop())
	fmt.Printf("After undo: %s\n", editor.GetContent())
}
