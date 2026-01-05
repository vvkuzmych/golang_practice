package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ============= json.Marshaler & Unmarshaler =============

// type Marshaler interface {
//     MarshalJSON() ([]byte, error)
// }

// type Unmarshaler interface {
//     UnmarshalJSON([]byte) error
// }

// ============= Custom Time Format =============

type CustomTime struct {
	time.Time
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", ct.Format("02.01.2006 15:04"))
	return []byte(formatted), nil
}

func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	str := string(data)
	str = strings.Trim(str, "\"")

	t, err := time.Parse("02.01.2006 15:04", str)
	if err != nil {
		return err
	}

	ct.Time = t
	return nil
}

// ============= Sensitive Data (приховування) =============

type Password string

func (p Password) MarshalJSON() ([]byte, error) {
	return []byte("\"***hidden***\""), nil
}

// ============= Custom Format =============

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

func (p Person) MarshalJSON() ([]byte, error) {
	// Власний формат JSON
	custom := map[string]interface{}{
		"full_name": p.FirstName + " " + p.LastName,
		"age":       p.Age,
		"is_adult":  p.Age >= 18,
	}
	return json.Marshal(custom)
}

// ============= Money (форматування) =============

type Money struct {
	Amount   int64 // копійки
	Currency string
}

func (m Money) MarshalJSON() ([]byte, error) {
	// Перетворити в формат "100.50 UAH"
	whole := m.Amount / 100
	cents := m.Amount % 100
	str := fmt.Sprintf("\"%d.%02d %s\"", whole, cents, m.Currency)
	return []byte(str), nil
}

func (m *Money) UnmarshalJSON(data []byte) error {
	// Парсити формат "100.50 UAH"
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	var whole, cents int64
	var currency string
	_, err := fmt.Sscanf(str, "%d.%d %s", &whole, &cents, &currency)
	if err != nil {
		return err
	}

	m.Amount = whole*100 + cents
	m.Currency = currency
	return nil
}

// ============= Status (enum) =============

type Status int

const (
	StatusPending Status = iota
	StatusActive
	StatusCompleted
	StatusCancelled
)

var statusNames = map[Status]string{
	StatusPending:   "pending",
	StatusActive:    "active",
	StatusCompleted: "completed",
	StatusCancelled: "cancelled",
}

var statusValues = map[string]Status{
	"pending":   StatusPending,
	"active":    StatusActive,
	"completed": StatusCompleted,
	"cancelled": StatusCancelled,
}

func (s Status) MarshalJSON() ([]byte, error) {
	if name, ok := statusNames[s]; ok {
		return json.Marshal(name)
	}
	return json.Marshal("unknown")
}

func (s *Status) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}

	if value, ok := statusValues[name]; ok {
		*s = value
		return nil
	}

	return fmt.Errorf("unknown status: %s", name)
}

// ============= Conditional Fields =============

type User struct {
	ID       int
	Username string
	Email    string
	Password Password
	IsAdmin  bool
}

func (u User) MarshalJSON() ([]byte, error) {
	type Alias User // уникнути рекурсії

	// Якщо не admin - не показувати email
	if !u.IsAdmin {
		return json.Marshal(&struct {
			*Alias
			Email string `json:"email,omitempty"`
		}{
			Alias: (*Alias)(&u),
			Email: "", // приховати
		})
	}

	return json.Marshal((Alias)(u))
}

// ============= Array Format =============

type RGB struct {
	R, G, B uint8
}

func (c RGB) MarshalJSON() ([]byte, error) {
	// Зберегти як масив [R, G, B]
	return json.Marshal([]uint8{c.R, c.G, c.B})
}

func (c *RGB) UnmarshalJSON(data []byte) error {
	var arr []uint8
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}

	if len(arr) != 3 {
		return fmt.Errorf("expected 3 values, got %d", len(arr))
	}

	c.R, c.G, c.B = arr[0], arr[1], arr[2]
	return nil
}

// ============= Main =============

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║    json.Marshaler & Unmarshaler          ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	// ===== Custom Time Format =====
	fmt.Println("\n🔹 Custom Time Format")
	fmt.Println("─────────────────────────────────────────")

	type Event struct {
		Name string
		Time CustomTime
	}

	event := Event{
		Name: "Конференція",
		Time: CustomTime{time.Now()},
	}

	jsonData, _ := json.MarshalIndent(event, "", "  ")
	fmt.Printf("JSON:\n%s\n", jsonData)

	// Unmarshal назад
	var decoded Event
	json.Unmarshal(jsonData, &decoded)
	fmt.Printf("Decoded: %s at %s\n", decoded.Name, decoded.Time.Format("02.01.2006 15:04"))

	// ===== Password (приховування) =====
	fmt.Println("\n🔹 Password (приховування)")
	fmt.Println("─────────────────────────────────────────")

	type Account struct {
		Username string
		Password Password
	}

	acc := Account{
		Username: "john",
		Password: "super-secret-123",
	}

	jsonData, _ = json.MarshalIndent(acc, "", "  ")
	fmt.Printf("JSON (password hidden):\n%s\n", jsonData)

	// ===== Custom Person Format =====
	fmt.Println("\n🔹 Person (custom format)")
	fmt.Println("─────────────────────────────────────────")

	person := Person{
		FirstName: "Іван",
		LastName:  "Петренко",
		Age:       25,
	}

	jsonData, _ = json.MarshalIndent(person, "", "  ")
	fmt.Printf("JSON:\n%s\n", jsonData)

	// ===== Money =====
	fmt.Println("\n🔹 Money (форматування)")
	fmt.Println("─────────────────────────────────────────")

	type Product struct {
		Name  string
		Price Money
	}

	product := Product{
		Name:  "Ноутбук",
		Price: Money{Amount: 2500050, Currency: "UAH"},
	}

	jsonData, _ = json.MarshalIndent(product, "", "  ")
	fmt.Printf("JSON:\n%s\n", jsonData)

	// Unmarshal назад
	var decodedProduct Product
	json.Unmarshal(jsonData, &decodedProduct)
	fmt.Printf("Decoded: %s costs %d копійок\n",
		decodedProduct.Name, decodedProduct.Price.Amount)

	// ===== Status (enum) =====
	fmt.Println("\n🔹 Status (enum)")
	fmt.Println("─────────────────────────────────────────")

	type Task struct {
		Title  string
		Status Status
	}

	task := Task{
		Title:  "Написати код",
		Status: StatusActive,
	}

	jsonData, _ = json.MarshalIndent(task, "", "  ")
	fmt.Printf("JSON:\n%s\n", jsonData)

	// Unmarshal назад
	var decodedTask Task
	json.Unmarshal(jsonData, &decodedTask)
	fmt.Printf("Decoded: %s is %d\n", decodedTask.Title, decodedTask.Status)

	// ===== User (conditional fields) =====
	fmt.Println("\n🔹 User (conditional fields)")
	fmt.Println("─────────────────────────────────────────")

	regularUser := User{
		ID:       1,
		Username: "john",
		Email:    "john@example.com",
		Password: "secret",
		IsAdmin:  false,
	}

	adminUser := User{
		ID:       2,
		Username: "admin",
		Email:    "admin@example.com",
		Password: "admin-secret",
		IsAdmin:  true,
	}

	fmt.Println("Regular User (email hidden):")
	jsonData, _ = json.MarshalIndent(regularUser, "", "  ")
	fmt.Println(string(jsonData))

	fmt.Println("\nAdmin User (email visible):")
	jsonData, _ = json.MarshalIndent(adminUser, "", "  ")
	fmt.Println(string(jsonData))

	// ===== RGB Color =====
	fmt.Println("\n🔹 RGB (array format)")
	fmt.Println("─────────────────────────────────────────")

	type Image struct {
		Name  string
		Color RGB
	}

	img := Image{
		Name:  "Background",
		Color: RGB{R: 255, G: 100, B: 50},
	}

	jsonData, _ = json.MarshalIndent(img, "", "  ")
	fmt.Printf("JSON:\n%s\n", jsonData)

	// Unmarshal назад
	var decodedImg Image
	json.Unmarshal(jsonData, &decodedImg)
	fmt.Printf("Decoded: %s color RGB(%d, %d, %d)\n",
		decodedImg.Name, decodedImg.Color.R, decodedImg.Color.G, decodedImg.Color.B)

	// ===== Multiple Statuses =====
	fmt.Println("\n🔹 Collection with custom marshaling")
	fmt.Println("─────────────────────────────────────────")

	tasks := []Task{
		{"Завдання 1", StatusPending},
		{"Завдання 2", StatusActive},
		{"Завдання 3", StatusCompleted},
		{"Завдання 4", StatusCancelled},
	}

	jsonData, _ = json.MarshalIndent(tasks, "", "  ")
	fmt.Printf("JSON:\n%s\n", jsonData)

	// ===== Error Handling =====
	fmt.Println("\n🔹 Error Handling")
	fmt.Println("─────────────────────────────────────────")

	invalidJSON := `{"Status": "invalid_status"}`
	var t Task
	err := json.Unmarshal([]byte(invalidJSON), &t)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	}

	validJSON := `{"Title": "Test", "Status": "active"}`
	err = json.Unmarshal([]byte(validJSON), &t)
	if err == nil {
		fmt.Printf("✅ Успішно: %s (%d)\n", t.Title, t.Status)
	}

	// ===== Comparison =====
	fmt.Println("\n🔹 Порівняння: з і без custom marshaling")
	fmt.Println("─────────────────────────────────────────")

	type SimpleTask struct {
		Title  string
		Status int
	}

	simple := SimpleTask{"Завдання", 1}
	custom := Task{"Завдання", StatusActive}

	simpleJSON, _ := json.Marshal(simple)
	customJSON, _ := json.Marshal(custom)

	fmt.Printf("Без custom: %s\n", simpleJSON)
	fmt.Printf("З custom:   %s\n", customJSON)

	// ===== Summary =====
	fmt.Println("\n\n📝 ВИСНОВКИ")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("✅ json.Marshaler дозволяє:")
	fmt.Println("   • Контролювати формат JSON")
	fmt.Println("   • Приховувати чутливі дані")
	fmt.Println("   • Форматувати дати/числа")
	fmt.Println("   • Enum → string")
	fmt.Println()
	fmt.Println("✅ json.Unmarshaler дозволяє:")
	fmt.Println("   • Парсити custom формати")
	fmt.Println("   • Валідацію при десеріалізації")
	fmt.Println("   • string → enum")
	fmt.Println()
	fmt.Println("💡 Коли використовувати:")
	fmt.Println("   • Custom формат дат")
	fmt.Println("   • Приховування паролів/токенів")
	fmt.Println("   • Enum types")
	fmt.Println("   • Складні структури даних")
	fmt.Println("   • API compatibility")
	fmt.Println()
	fmt.Println("⚠️  Увага:")
	fmt.Println("   • Уникати рекурсії (type Alias)")
	fmt.Println("   • Перевіряти помилки")
	fmt.Println("   • Тестувати marshal/unmarshal разом")
}
