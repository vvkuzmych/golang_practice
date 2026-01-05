package main

import (
	"fmt"
	"strings"
	"time"
)

// ============= fmt.Stringer Interface =============

// type Stringer interface {
//     String() string
// }

// ============= Person =============

type Person struct {
	FirstName string
	LastName  string
	Age       int
	Email     string
}

// Реалізація String()
func (p Person) String() string {
	return fmt.Sprintf("%s %s (%d років) <%s>",
		p.FirstName, p.LastName, p.Age, p.Email)
}

// ============= Money =============

type Money struct {
	Amount   int64 // копійки
	Currency string
}

func (m Money) String() string {
	whole := m.Amount / 100
	cents := m.Amount % 100
	return fmt.Sprintf("%d.%02d %s", whole, cents, m.Currency)
}

// ============= LogLevel =============

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func (l LogLevel) String() string {
	levels := []string{"DEBUG", "INFO", "WARNING", "ERROR", "FATAL"}
	if l < 0 || int(l) >= len(levels) {
		return "UNKNOWN"
	}
	return levels[l]
}

// ============= HTTPStatus =============

type HTTPStatus int

func (s HTTPStatus) String() string {
	statuses := map[int]string{
		200: "200 OK",
		201: "201 Created",
		400: "400 Bad Request",
		401: "401 Unauthorized",
		403: "403 Forbidden",
		404: "404 Not Found",
		500: "500 Internal Server Error",
	}

	if status, ok := statuses[int(s)]; ok {
		return status
	}
	return fmt.Sprintf("%d Unknown Status", s)
}

// ============= IPAddress =============

type IPAddress [4]byte

func (ip IPAddress) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

// ============= Duration (custom) =============

type Duration struct {
	Seconds int
}

func (d Duration) String() string {
	if d.Seconds < 60 {
		return fmt.Sprintf("%d секунд", d.Seconds)
	}
	minutes := d.Seconds / 60
	seconds := d.Seconds % 60
	if minutes < 60 {
		return fmt.Sprintf("%d хв %d сек", minutes, seconds)
	}
	hours := minutes / 60
	minutes = minutes % 60
	return fmt.Sprintf("%d год %d хв %d сек", hours, minutes, seconds)
}

// ============= Book =============

type Book struct {
	Title       string
	Author      string
	Year        int
	Pages       int
	ISBN        string
	IsAvailable bool
}

func (b Book) String() string {
	status := "доступна"
	if !b.IsAvailable {
		status = "видана"
	}
	return fmt.Sprintf("📚 \"%s\" by %s (%d) - %d стор., ISBN: %s [%s]",
		b.Title, b.Author, b.Year, b.Pages, b.ISBN, status)
}

// ============= Temperature =============

type Temperature struct {
	Value float64
	Unit  string // C або F
}

func (t Temperature) String() string {
	return fmt.Sprintf("%.1f°%s", t.Value, t.Unit)
}

func (t Temperature) ToCelsius() Temperature {
	if t.Unit == "F" {
		return Temperature{Value: (t.Value - 32) * 5 / 9, Unit: "C"}
	}
	return t
}

func (t Temperature) ToFahrenheit() Temperature {
	if t.Unit == "C" {
		return Temperature{Value: t.Value*9/5 + 32, Unit: "F"}
	}
	return t
}

// ============= List =============

type StringList []string

func (sl StringList) String() string {
	if len(sl) == 0 {
		return "[]"
	}
	return "[" + strings.Join(sl, ", ") + "]"
}

// ============= Card =============

type Card struct {
	Rank string
	Suit string
}

func (c Card) String() string {
	suits := map[string]string{
		"hearts":   "♥",
		"diamonds": "♦",
		"clubs":    "♣",
		"spades":   "♠",
	}

	suit := suits[c.Suit]
	if suit == "" {
		suit = c.Suit
	}

	return fmt.Sprintf("%s%s", c.Rank, suit)
}

// ============= Main =============

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║       fmt.Stringer Interface             ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	// ===== Person =====
	fmt.Println("\n🔹 Person")
	fmt.Println("─────────────────────────────────────────")

	person := Person{
		FirstName: "Іван",
		LastName:  "Петренко",
		Age:       30,
		Email:     "ivan@example.com",
	}

	// fmt.Println автоматично викликає String()
	fmt.Println(person)

	// Явний виклик
	fmt.Printf("Явний виклик: %s\n", person.String())

	// З форматом %v
	fmt.Printf("%%v format: %v\n", person)

	// ===== Money =====
	fmt.Println("\n🔹 Money")
	fmt.Println("─────────────────────────────────────────")

	prices := []Money{
		{Amount: 15050, Currency: "UAH"},
		{Amount: 10000, Currency: "USD"},
		{Amount: 50500, Currency: "EUR"},
	}

	for _, price := range prices {
		fmt.Printf("  Ціна: %s\n", price)
	}

	// ===== LogLevel =====
	fmt.Println("\n🔹 LogLevel")
	fmt.Println("─────────────────────────────────────────")

	levels := []LogLevel{DEBUG, INFO, WARNING, ERROR, FATAL}
	for _, level := range levels {
		fmt.Printf("[%s] Повідомлення рівня %s\n", level, level)
	}

	// ===== HTTPStatus =====
	fmt.Println("\n🔹 HTTPStatus")
	fmt.Println("─────────────────────────────────────────")

	statuses := []HTTPStatus{200, 404, 500, 201, 401}
	for _, status := range statuses {
		fmt.Printf("  Status: %s\n", status)
	}

	// ===== IPAddress =====
	fmt.Println("\n🔹 IPAddress")
	fmt.Println("─────────────────────────────────────────")

	ips := []IPAddress{
		{192, 168, 1, 1},
		{10, 0, 0, 1},
		{127, 0, 0, 1},
		{8, 8, 8, 8},
	}

	for _, ip := range ips {
		fmt.Printf("  IP: %s\n", ip)
	}

	// ===== Duration =====
	fmt.Println("\n🔹 Duration")
	fmt.Println("─────────────────────────────────────────")

	durations := []Duration{
		{30},
		{90},
		{3600},
		{7265},
	}

	for _, d := range durations {
		fmt.Printf("  %d секунд = %s\n", d.Seconds, d)
	}

	// ===== Book =====
	fmt.Println("\n🔹 Book")
	fmt.Println("─────────────────────────────────────────")

	books := []Book{
		{
			Title:       "The Go Programming Language",
			Author:      "Donovan & Kernighan",
			Year:        2015,
			Pages:       380,
			ISBN:        "978-0134190440",
			IsAvailable: true,
		},
		{
			Title:       "Go in Action",
			Author:      "William Kennedy",
			Year:        2015,
			Pages:       264,
			ISBN:        "978-1617291784",
			IsAvailable: false,
		},
	}

	for _, book := range books {
		fmt.Println(book)
	}

	// ===== Temperature =====
	fmt.Println("\n🔹 Temperature")
	fmt.Println("─────────────────────────────────────────")

	tempC := Temperature{Value: 25, Unit: "C"}
	tempF := tempC.ToFahrenheit()

	fmt.Printf("Температура: %s\n", tempC)
	fmt.Printf("У Фаренгейтах: %s\n", tempF)
	fmt.Printf("Назад в Цельсії: %s\n", tempF.ToCelsius())

	// ===== StringList =====
	fmt.Println("\n🔹 StringList")
	fmt.Println("─────────────────────────────────────────")

	languages := StringList{"Go", "Python", "JavaScript", "Rust"}
	emptyList := StringList{}

	fmt.Printf("Мови: %s\n", languages)
	fmt.Printf("Порожній список: %s\n", emptyList)

	// ===== Card =====
	fmt.Println("\n🔹 Playing Cards")
	fmt.Println("─────────────────────────────────────────")

	cards := []Card{
		{"A", "spades"},
		{"K", "hearts"},
		{"Q", "diamonds"},
		{"J", "clubs"},
		{"10", "spades"},
	}

	fmt.Print("Карти: ")
	for i, card := range cards {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(card)
	}
	fmt.Println()

	// ===== Використання в логуванні =====
	fmt.Println("\n🔹 Використання в логуванні")
	fmt.Println("─────────────────────────────────────────")

	type LogEntry struct {
		Time    time.Time
		Level   LogLevel
		Message string
		User    Person
	}

	entry := LogEntry{
		Time:    time.Now(),
		Level:   ERROR,
		Message: "Failed to connect",
		User:    person,
	}

	fmt.Printf("[%s] %s: %s (user: %s)\n",
		entry.Time.Format("15:04:05"),
		entry.Level,
		entry.Message,
		entry.User)

	// ===== Порівняння з і без String() =====
	fmt.Println("\n🔹 З і без String()")
	fmt.Println("─────────────────────────────────────────")

	type SimpleStruct struct {
		Name  string
		Value int
	}

	simple := SimpleStruct{Name: "Test", Value: 42}

	fmt.Printf("Без String():   %v\n", simple)
	fmt.Printf("Із String():    %s\n", person)

	// ===== Summary =====
	fmt.Println("\n\n📝 ВИСНОВКИ")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("✅ fmt.Stringer дозволяє:")
	fmt.Println("   • Контролювати вивід структур")
	fmt.Println("   • Красиве форматування")
	fmt.Println("   • Зручне логування")
	fmt.Println("   • Читабельний код")
	fmt.Println()
	fmt.Println("💡 Коли реалізувати:")
	fmt.Println("   • Складні структури даних")
	fmt.Println("   • Custom типи (enum, id, тощо)")
	fmt.Println("   • Об'єкти для логів")
	fmt.Println("   • API response objects")
	fmt.Println()
	fmt.Println("⚡ fmt.Println() автоматично викликає String()")
	fmt.Println("   Якщо String() є - використовує його")
	fmt.Println("   Якщо немає - використовує дефолтний формат")
}
