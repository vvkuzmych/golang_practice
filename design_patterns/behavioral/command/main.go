package main

import (
	"fmt"
	"strings"
)

// ============= Command Interface =============

type Command interface {
	Execute()
	Undo()
	GetName() string
}

// ============= Receiver: Light =============

type Light struct {
	isOn bool
}

func (l *Light) On() {
	l.isOn = true
	fmt.Println("   💡 Light is ON")
}

func (l *Light) Off() {
	l.isOn = false
	fmt.Println("   💡 Light is OFF")
}

// Commands
type LightOnCommand struct {
	light *Light
}

func (c *LightOnCommand) Execute() {
	c.light.On()
}

func (c *LightOnCommand) Undo() {
	c.light.Off()
}

func (c *LightOnCommand) GetName() string {
	return "Light ON"
}

type LightOffCommand struct {
	light *Light
}

func (c *LightOffCommand) Execute() {
	c.light.Off()
}

func (c *LightOffCommand) Undo() {
	c.light.On()
}

func (c *LightOffCommand) GetName() string {
	return "Light OFF"
}

// ============= Receiver: Document =============

type Document struct {
	content string
}

type WriteCommand struct {
	doc      *Document
	text     string
	prevText string
}

func (c *WriteCommand) Execute() {
	c.prevText = c.doc.content
	c.doc.content += c.text
	fmt.Printf("   ✏️  Wrote: '%s' (total: '%s')\n", c.text, c.doc.content)
}

func (c *WriteCommand) Undo() {
	c.doc.content = c.prevText
	fmt.Printf("   ↩️  Undo write (reverted to: '%s')\n", c.doc.content)
}

func (c *WriteCommand) GetName() string {
	return fmt.Sprintf("Write '%s'", c.text)
}

// ============= Invoker: Remote Control =============

type RemoteControl struct {
	history []Command
	current int
}

func NewRemoteControl() *RemoteControl {
	return &RemoteControl{
		history: []Command{},
		current: -1,
	}
}

func (r *RemoteControl) Execute(cmd Command) {
	// Remove commands after current position
	r.history = r.history[:r.current+1]

	cmd.Execute()
	r.history = append(r.history, cmd)
	r.current++
}

func (r *RemoteControl) Undo() {
	if r.current >= 0 {
		cmd := r.history[r.current]
		fmt.Printf("   ⬅️  Undoing: %s\n", cmd.GetName())
		cmd.Undo()
		r.current--
	} else {
		fmt.Println("   ❌ Nothing to undo")
	}
}

func (r *RemoteControl) Redo() {
	if r.current < len(r.history)-1 {
		r.current++
		cmd := r.history[r.current]
		fmt.Printf("   ➡️  Redoing: %s\n", cmd.GetName())
		cmd.Execute()
	} else {
		fmt.Println("   ❌ Nothing to redo")
	}
}

func (r *RemoteControl) ShowHistory() {
	fmt.Println("\n📜 Command History:")
	for i, cmd := range r.history {
		marker := " "
		if i == r.current {
			marker = "→"
		}
		fmt.Printf("   %s %d. %s\n", marker, i+1, cmd.GetName())
	}
}

// ============= Main =============

func main() {
	fmt.Println("╔════════════════════════════════════════════════╗")
	fmt.Println("║          Command Pattern Demo                  ║")
	fmt.Println("╚════════════════════════════════════════════════╝")

	// ===== Example 1: Light Control =====
	fmt.Println("\n🔹 Example 1: Light Control with Undo/Redo")
	fmt.Println(strings.Repeat("─", 50))

	light := &Light{}
	remote := NewRemoteControl()

	fmt.Println("\n💡 Turning light on/off:")
	remote.Execute(&LightOnCommand{light: light})
	remote.Execute(&LightOffCommand{light: light})
	remote.Execute(&LightOnCommand{light: light})

	fmt.Println("\n↩️  Undoing commands:")
	remote.Undo()
	remote.Undo()

	fmt.Println("\n↪️  Redoing commands:")
	remote.Redo()

	remote.ShowHistory()

	// ===== Example 2: Text Editor =====
	fmt.Println("\n\n🔹 Example 2: Text Editor with Undo/Redo")
	fmt.Println(strings.Repeat("─", 50))

	doc := &Document{content: ""}
	editor := NewRemoteControl()

	fmt.Println("\n✏️  Writing text:")
	editor.Execute(&WriteCommand{doc: doc, text: "Hello"})
	editor.Execute(&WriteCommand{doc: doc, text: " World"})
	editor.Execute(&WriteCommand{doc: doc, text: "!"})

	fmt.Printf("\n📄 Current document: '%s'\n", doc.content)

	fmt.Println("\n↩️  Undo last 2 operations:")
	editor.Undo()
	editor.Undo()
	fmt.Printf("📄 Document now: '%s'\n", doc.content)

	fmt.Println("\n↪️  Redo 1 operation:")
	editor.Redo()
	fmt.Printf("📄 Document now: '%s'\n", doc.content)

	editor.ShowHistory()

	// ===== Example 3: Macro (Multiple Commands) =====
	fmt.Println("\n\n🔹 Example 3: Macro Recording")
	fmt.Println(strings.Repeat("─", 50))

	fmt.Println("\n🎬 Recording macro:")
	macro := []Command{
		&LightOnCommand{light: light},
		&WriteCommand{doc: doc, text: " Macro"},
		&LightOffCommand{light: light},
	}

	macroControl := NewRemoteControl()
	for _, cmd := range macro {
		macroControl.Execute(cmd)
	}

	fmt.Println("\n↩️  Undo entire macro:")
	for range macro {
		macroControl.Undo()
	}

	// ===== Example 4: Command Queue =====
	fmt.Println("\n\n🔹 Example 4: Command Queue")
	fmt.Println(strings.Repeat("─", 50))

	queue := []Command{
		&LightOnCommand{light: light},
		&WriteCommand{doc: doc, text: " Queue1"},
		&WriteCommand{doc: doc, text: " Queue2"},
		&LightOffCommand{light: light},
	}

	fmt.Println("\n⏳ Executing queued commands:")
	queueControl := NewRemoteControl()
	for i, cmd := range queue {
		fmt.Printf("\n▶️  Executing command %d:\n", i+1)
		queueControl.Execute(cmd)
	}

	// ===== Summary =====
	fmt.Println("\n\n📝 ВИСНОВКИ")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println("✅ Інкапсулює запити як об'єкти")
	fmt.Println("✅ Undo/Redo підтримка")
	fmt.Println("✅ Історія команд")
	fmt.Println("✅ Черги команд")
	fmt.Println("✅ Macro recording")

	fmt.Println("\n💡 ВИКОРИСТАННЯ:")
	fmt.Println("   - Text editors (Undo/Redo)")
	fmt.Println("   - Task queues")
	fmt.Println("   - Transaction systems")
	fmt.Println("   - UI actions")
	fmt.Println("   - Macro automation")

	fmt.Println("\n🎯 Ключові переваги:")
	fmt.Println("   - Відокремлює відправника від отримувача")
	fmt.Println("   - Легко додавати нові команди")
	fmt.Println("   - Можна комбінувати команди")
}
