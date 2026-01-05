package main

import (
	"fmt"
	"strings"
)

// ============= Product Interfaces =============

// Button інтерфейс
type Button interface {
	Click() string
	Render() string
}

// Checkbox інтерфейс
type Checkbox interface {
	Check() string
	Render() string
}

// ============= Abstract Factory =============

type GUIFactory interface {
	CreateButton() Button
	CreateCheckbox() Checkbox
	GetName() string
}

// ============= Windows Family =============

type WindowsButton struct{}

func (w *WindowsButton) Click() string {
	return "🖱️  Windows button clicked"
}

func (w *WindowsButton) Render() string {
	return "[Windows Button]"
}

type WindowsCheckbox struct{}

func (w *WindowsCheckbox) Check() string {
	return "☑️  Windows checkbox checked"
}

func (w *WindowsCheckbox) Render() string {
	return "[x] Windows Checkbox"
}

type WindowsFactory struct{}

func (w *WindowsFactory) CreateButton() Button {
	return &WindowsButton{}
}

func (w *WindowsFactory) CreateCheckbox() Checkbox {
	return &WindowsCheckbox{}
}

func (w *WindowsFactory) GetName() string {
	return "Windows"
}

// ============= Mac Family =============

type MacButton struct{}

func (m *MacButton) Click() string {
	return "🖱️  Mac button clicked"
}

func (m *MacButton) Render() string {
	return "( Mac Button )"
}

type MacCheckbox struct{}

func (m *MacCheckbox) Check() string {
	return "☑️  Mac checkbox checked"
}

func (m *MacCheckbox) Render() string {
	return "✓ Mac Checkbox"
}

type MacFactory struct{}

func (m *MacFactory) CreateButton() Button {
	return &MacButton{}
}

func (m *MacFactory) CreateCheckbox() Checkbox {
	return &MacCheckbox{}
}

func (m *MacFactory) GetName() string {
	return "Mac"
}

// ============= Linux Family =============

type LinuxButton struct{}

func (l *LinuxButton) Click() string {
	return "🖱️  Linux button clicked"
}

func (l *LinuxButton) Render() string {
	return "[ Linux Button ]"
}

type LinuxCheckbox struct{}

func (l *LinuxCheckbox) Check() string {
	return "☑️  Linux checkbox checked"
}

func (l *LinuxCheckbox) Render() string {
	return "[X] Linux Checkbox"
}

type LinuxFactory struct{}

func (l *LinuxFactory) CreateButton() Button {
	return &LinuxButton{}
}

func (l *LinuxFactory) CreateCheckbox() Checkbox {
	return &LinuxCheckbox{}
}

func (l *LinuxFactory) GetName() string {
	return "Linux"
}

// ============= Application =============

type Application struct {
	factory  GUIFactory
	button   Button
	checkbox Checkbox
}

func NewApplication(factory GUIFactory) *Application {
	return &Application{
		factory:  factory,
		button:   factory.CreateButton(),
		checkbox: factory.CreateCheckbox(),
	}
}

func (a *Application) Render() {
	fmt.Printf("\n🎨 Rendering %s UI:\n", a.factory.GetName())
	fmt.Printf("   %s\n", a.button.Render())
	fmt.Printf("   %s\n", a.checkbox.Render())
}

func (a *Application) Interact() {
	fmt.Println("\n👆 User interactions:")
	fmt.Printf("   %s\n", a.button.Click())
	fmt.Printf("   %s\n", a.checkbox.Check())
}

// ============= Factory Selector =============

func GetFactory(os string) GUIFactory {
	switch strings.ToLower(os) {
	case "windows", "win":
		return &WindowsFactory{}
	case "mac", "macos":
		return &MacFactory{}
	case "linux":
		return &LinuxFactory{}
	default:
		return &WindowsFactory{}
	}
}

// ============= Main =============

func main() {
	fmt.Println("╔════════════════════════════════════════════════╗")
	fmt.Println("║      Abstract Factory Pattern Demo            ║")
	fmt.Println("╚════════════════════════════════════════════════╝")

	// ===== Demo 1: Different OS UIs =====
	fmt.Println("\n🔹 Demo 1: Creating UI for Different OS")
	fmt.Println(strings.Repeat("─", 50))

	// Windows
	fmt.Println("\n💻 Creating Windows application:")
	winFactory := &WindowsFactory{}
	winApp := NewApplication(winFactory)
	winApp.Render()
	winApp.Interact()

	// Mac
	fmt.Println("\n🍎 Creating Mac application:")
	macFactory := &MacFactory{}
	macApp := NewApplication(macFactory)
	macApp.Render()
	macApp.Interact()

	// Linux
	fmt.Println("\n🐧 Creating Linux application:")
	linuxFactory := &LinuxFactory{}
	linuxApp := NewApplication(linuxFactory)
	linuxApp.Render()
	linuxApp.Interact()

	// ===== Demo 2: Runtime Selection =====
	fmt.Println("\n\n🔹 Demo 2: Runtime OS Detection")
	fmt.Println(strings.Repeat("─", 50))

	// Симулюємо визначення OS
	detectedOS := "mac" // можна було б використати runtime.GOOS

	fmt.Printf("\n🔍 Detected OS: %s\n", detectedOS)
	factory := GetFactory(detectedOS)
	app := NewApplication(factory)
	app.Render()
	app.Interact()

	// ===== Demo 3: Guaranteed Compatibility =====
	fmt.Println("\n\n🔹 Demo 3: Guaranteed Compatibility")
	fmt.Println(strings.Repeat("─", 50))

	fmt.Println("\n✅ All components from same family:")
	fmt.Println("   - Button and Checkbox match")
	fmt.Println("   - Consistent look and feel")
	fmt.Println("   - No mixing Windows button with Mac checkbox")

	// ===== Comparison =====
	fmt.Println("\n\n🔹 Visual Comparison")
	fmt.Println(strings.Repeat("─", 50))

	factories := []GUIFactory{
		&WindowsFactory{},
		&MacFactory{},
		&LinuxFactory{},
	}

	for _, f := range factories {
		btn := f.CreateButton()
		chk := f.CreateCheckbox()
		fmt.Printf("\n%-10s %s  |  %s\n",
			f.GetName()+":", btn.Render(), chk.Render())
	}

	// ===== Summary =====
	fmt.Println("\n\n📝 ВИСНОВКИ")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println("✅ Створює сімейства пов'язаних об'єктів")
	fmt.Println("✅ Гарантує сумісність компонентів")
	fmt.Println("✅ Легко додавати нові платформи")
	fmt.Println("✅ Ізоляція конкретних класів")

	fmt.Println("\n💡 ВИКОРИСТАННЯ:")
	fmt.Println("   - UI frameworks (Windows/Mac/Linux)")
	fmt.Println("   - Database drivers (MySQL/Postgres/MongoDB)")
	fmt.Println("   - Document generators (PDF/XML/JSON)")
	fmt.Println("   - Cloud providers (AWS/Azure/GCP)")

	fmt.Println("\n🔄 Відмінність від Factory Method:")
	fmt.Println("   Factory Method: створює ОДИН тип продукту")
	fmt.Println("   Abstract Factory: створює СІМЕЙСТВО продуктів")
}
