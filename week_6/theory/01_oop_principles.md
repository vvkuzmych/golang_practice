# ООП Принципи в Go

Go - не класичн ООП мова, але підтримує всі основні принципи через свої механізми.

---

## 📖 Зміст

1. [Інкапсуляція](#1-інкапсуляція-encapsulation)
2. [Поліморфізм](#2-поліморфізм-polymorphism)
3. [Абстракція](#3-абстракція-abstraction)
4. [Композиція замість успадкування](#4-композиція-замість-успадкування)

---

## 1. Інкапсуляція (Encapsulation)

### Що це?
**Інкапсуляція** - приховування внутрішньої реалізації та надання доступу через публічні методи.

### В Go
У Go немає ключових слів `private`/`public`. Видимість визначається **регістром першої літери**:
- **Велика літера** = public (експортовано)
- **Маленька літера** = private (не експортовано)

### Приклад

```go
package user

import "fmt"

// User - публічна структура (експортована)
type User struct {
    ID       int    // публічне поле
    Username string // публічне поле
    password string // приватне поле ❌ не доступне ззовні
    email    string // приватне поле
}

// NewUser - публічний конструктор
func NewUser(username, password, email string) *User {
    return &User{
        Username: username,
        password: password, // встановлюємо приватне поле
        email:    email,
    }
}

// GetEmail - публічний getter
func (u *User) GetEmail() string {
    return u.email
}

// SetEmail - публічний setter з валідацією
func (u *User) SetEmail(email string) error {
    if len(email) < 5 {
        return fmt.Errorf("email too short")
    }
    u.email = email
    return nil
}

// ValidatePassword - публічний метод
func (u *User) ValidatePassword(password string) bool {
    return u.password == password // доступ до приватного поля
}

// hashPassword - приватний метод ❌ не доступний ззовні
func (u *User) hashPassword(password string) string {
    // логіка хешування
    return "hashed_" + password
}
```

### Використання

```go
package main

import "myapp/user"

func main() {
    // ✅ Створюємо користувача через конструктор
    u := user.NewUser("john", "secret123", "john@example.com")
    
    // ✅ Доступ до публічних полів
    fmt.Println(u.Username) // "john"
    
    // ❌ Немає доступу до приватних полів
    // fmt.Println(u.password) // compile error
    
    // ✅ Використовуємо getter
    fmt.Println(u.GetEmail()) // "john@example.com"
    
    // ✅ Використовуємо setter з валідацією
    u.SetEmail("newemail@example.com")
    
    // ✅ Публічний метод
    if u.ValidatePassword("secret123") {
        fmt.Println("Password is correct")
    }
    
    // ❌ Немає доступу до приватних методів
    // u.hashPassword("test") // compile error
}
```

### Переваги інкапсуляції

1. **Контроль доступу** - зовнішній код не може змінити критичні дані напряму
2. **Валідація** - можна перевіряти дані перед зміною
3. **Гнучкість** - можна змінити внутрішню реалізацію без зміни API
4. **Безпека** - приватні дані захищені від несанкціонованого доступу

---

## 2. Поліморфізм (Polymorphism)

### Що це?
**Поліморфізм** - можливість об'єктів різних типів відповідати на однакові повідомлення (методи).

### В Go
Поліморфізм реалізується через **інтерфейси** та **duck typing**.

### Приклад: Система платежів

```go
package main

import "fmt"

// PaymentProcessor - інтерфейс для обробки платежів
type PaymentProcessor interface {
    Process(amount float64) error
    GetName() string
}

// ===== Різні реалізації =====

// CreditCardProcessor - оплата кредиткою
type CreditCardProcessor struct {
    CardNumber string
}

func (c *CreditCardProcessor) Process(amount float64) error {
    fmt.Printf("Processing $%.2f via Credit Card %s\n", amount, c.CardNumber)
    return nil
}

func (c *CreditCardProcessor) GetName() string {
    return "Credit Card"
}

// PayPalProcessor - оплата через PayPal
type PayPalProcessor struct {
    Email string
}

func (p *PayPalProcessor) Process(amount float64) error {
    fmt.Printf("Processing $%.2f via PayPal %s\n", amount, p.Email)
    return nil
}

func (p *PayPalProcessor) GetName() string {
    return "PayPal"
}

// CryptoProcessor - оплата криптовалютою
type CryptoProcessor struct {
    WalletAddress string
}

func (c *CryptoProcessor) Process(amount float64) error {
    fmt.Printf("Processing $%.2f via Crypto Wallet %s\n", amount, c.WalletAddress)
    return nil
}

func (c *CryptoProcessor) GetName() string {
    return "Cryptocurrency"
}

// ===== Використання поліморфізму =====

// ProcessPayment - приймає будь-який PaymentProcessor
func ProcessPayment(processor PaymentProcessor, amount float64) error {
    fmt.Printf("Using %s processor\n", processor.GetName())
    return processor.Process(amount)
}

func main() {
    // Створюємо різні процесори
    creditCard := &CreditCardProcessor{CardNumber: "**** 1234"}
    paypal := &PayPalProcessor{Email: "user@example.com"}
    crypto := &CryptoProcessor{WalletAddress: "0x123..."}
    
    // ✅ Викликаємо одну функцію з різними типами
    ProcessPayment(creditCard, 100.50) // Credit Card
    ProcessPayment(paypal, 75.00)      // PayPal
    ProcessPayment(crypto, 200.00)     // Crypto
    
    // ✅ Список різних процесорів
    processors := []PaymentProcessor{creditCard, paypal, crypto}
    for _, p := range processors {
        ProcessPayment(p, 50.00)
    }
}
```

### Приклад: Duck Typing

```go
// Інтерфейс з одним методом
type Swimmer interface {
    Swim() string
}

type Duck struct{}
func (d Duck) Swim() string { return "Duck is swimming" }

type Fish struct{}
func (f Fish) Swim() string { return "Fish is swimming" }

type Human struct{}
func (h Human) Swim() string { return "Human is swimming" }

// Функція приймає будь-що, що вміє плавати
func MakeSwim(s Swimmer) {
    fmt.Println(s.Swim())
}

func main() {
    MakeSwim(Duck{})   // "Duck is swimming"
    MakeSwim(Fish{})   // "Fish is swimming"
    MakeSwim(Human{})  // "Human is swimming"
}
```

### Порожній інтерфейс `interface{}`

```go
// interface{} (або any в Go 1.18+) приймає будь-який тип
func Print(value interface{}) {
    fmt.Println(value)
}

Print(42)          // int
Print("hello")     // string
Print([]int{1,2})  // slice
Print(struct{}{})  // struct
```

### Переваги поліморфізму

1. **Гнучкість** - один інтерфейс, багато реалізацій
2. **Розширюваність** - легко додати нову реалізацію
3. **Тестування** - легко створити mock об'єкти
4. **Чистий код** - менше дублювання

---

## 3. Абстракція (Abstraction)

### Що це?
**Абстракція** - приховування складної реалізації за простим інтерфейсом.

### Приклад: Database Abstraction

```go
package main

import "fmt"

// Database - абстрактний інтерфейс для роботи з БД
type Database interface {
    Connect() error
    Query(sql string) ([]Row, error)
    Close() error
}

type Row map[string]interface{}

// ===== PostgreSQL реалізація =====
type PostgresDB struct {
    host     string
    port     int
    username string
    password string
}

func (p *PostgresDB) Connect() error {
    fmt.Printf("Connecting to PostgreSQL at %s:%d\n", p.host, p.port)
    // складна логіка підключення
    return nil
}

func (p *PostgresDB) Query(sql string) ([]Row, error) {
    fmt.Printf("Executing PostgreSQL query: %s\n", sql)
    // складна логіка запиту
    return []Row{{"id": 1, "name": "John"}}, nil
}

func (p *PostgresDB) Close() error {
    fmt.Println("Closing PostgreSQL connection")
    return nil
}

// ===== MongoDB реалізація =====
type MongoDB struct {
    connectionString string
}

func (m *MongoDB) Connect() error {
    fmt.Printf("Connecting to MongoDB: %s\n", m.connectionString)
    return nil
}

func (m *MongoDB) Query(sql string) ([]Row, error) {
    fmt.Printf("Executing MongoDB query: %s\n", sql)
    return []Row{{"_id": "abc", "name": "Jane"}}, nil
}

func (m *MongoDB) Close() error {
    fmt.Println("Closing MongoDB connection")
    return nil
}

// ===== Використання абстракції =====

// UserRepository - працює з будь-якою БД
type UserRepository struct {
    db Database // абстракція, не конкретна реалізація
}

func (r *UserRepository) GetAllUsers() ([]Row, error) {
    // Не знаємо з якою БД працюємо - це абстраговано
    return r.db.Query("SELECT * FROM users")
}

func main() {
    // ✅ Можемо використовувати PostgreSQL
    postgres := &PostgresDB{host: "localhost", port: 5432}
    userRepo1 := &UserRepository{db: postgres}
    postgres.Connect()
    users1, _ := userRepo1.GetAllUsers()
    fmt.Println(users1)
    postgres.Close()
    
    fmt.Println("---")
    
    // ✅ Або MongoDB - код не змінюється!
    mongo := &MongoDB{connectionString: "mongodb://localhost"}
    userRepo2 := &UserRepository{db: mongo}
    mongo.Connect()
    users2, _ := userRepo2.GetAllUsers()
    fmt.Println(users2)
    mongo.Close()
}
```

### Приклад: Logger Abstraction

```go
// Logger - абстракція для логування
type Logger interface {
    Info(message string)
    Error(message string)
    Debug(message string)
}

// ConsoleLogger - логує в консоль
type ConsoleLogger struct{}

func (c *ConsoleLogger) Info(message string) {
    fmt.Printf("[INFO] %s\n", message)
}

func (c *ConsoleLogger) Error(message string) {
    fmt.Printf("[ERROR] %s\n", message)
}

func (c *ConsoleLogger) Debug(message string) {
    fmt.Printf("[DEBUG] %s\n", message)
}

// FileLogger - логує у файл
type FileLogger struct {
    filename string
}

func (f *FileLogger) Info(message string) {
    // запис у файл
    fmt.Printf("Writing to %s: [INFO] %s\n", f.filename, message)
}

func (f *FileLogger) Error(message string) {
    fmt.Printf("Writing to %s: [ERROR] %s\n", f.filename, message)
}

func (f *FileLogger) Debug(message string) {
    fmt.Printf("Writing to %s: [DEBUG] %s\n", f.filename, message)
}

// Application - використовує абстракцію Logger
type Application struct {
    logger Logger // не знаємо який саме logger
}

func (a *Application) Start() {
    a.logger.Info("Application started")
    // складна логіка
    a.logger.Debug("Processing...")
}

func main() {
    // Легко міняємо реалізацію
    app1 := &Application{logger: &ConsoleLogger{}}
    app1.Start()
    
    app2 := &Application{logger: &FileLogger{filename: "app.log"}}
    app2.Start()
}
```

### Переваги абстракції

1. **Простота** - складна логіка прихована
2. **Модульність** - легко міняти реалізацію
3. **Тестування** - легко створити mock
4. **Підтримка** - зміни в одному місці

---

## 4. Композиція замість успадкування

### Проблема успадкування
У Go немає класичного успадкування (extends), і це **добре**! Успадкування створює жорсткі зв'язки між класами.

### Рішення: Композиція
**Композиція** - включення одного об'єкта в інший.

### Приклад: Embedding (вбудовування)

```go
package main

import "fmt"

// ===== Базові компоненти =====

// Engine - двигун
type Engine struct {
    Horsepower int
}

func (e *Engine) Start() {
    fmt.Println("Engine started")
}

func (e *Engine) Stop() {
    fmt.Println("Engine stopped")
}

// Wheels - колеса
type Wheels struct {
    Count int
}

func (w *Wheels) Rotate() {
    fmt.Println("Wheels are rotating")
}

// GPS - навігація
type GPS struct {
    Model string
}

func (g *GPS) Navigate(destination string) {
    fmt.Printf("Navigating to %s using %s\n", destination, g.Model)
}

// ===== Композиція =====

// Car - машина, що МІСТИТЬ інші компоненти
type Car struct {
    Brand  string
    Engine Engine  // композиція
    Wheels Wheels  // композиція
    GPS    *GPS    // опціональна композиція (pointer)
}

func (c *Car) Drive() {
    fmt.Printf("Driving %s\n", c.Brand)
    c.Engine.Start()
    c.Wheels.Rotate()
    if c.GPS != nil {
        c.GPS.Navigate("Home")
    }
}

// ===== Embedding (вбудовування) =====

// ElectricCar - електромобіль
type ElectricCar struct {
    Car          // вбудовуємо Car (anonymous field)
    BatteryLevel int
}

// Методи Car доступні напряму через ElectricCar
func (e *ElectricCar) Charge() {
    fmt.Println("Charging battery...")
    e.BatteryLevel = 100
}

func main() {
    // Звичайна машина
    car := Car{
        Brand:  "Toyota",
        Engine: Engine{Horsepower: 150},
        Wheels: Wheels{Count: 4},
        GPS:    &GPS{Model: "Garmin"},
    }
    car.Drive()
    
    fmt.Println("---")
    
    // Електромобіль
    tesla := ElectricCar{
        Car: Car{
            Brand:  "Tesla",
            Engine: Engine{Horsepower: 500},
            Wheels: Wheels{Count: 4},
        },
        BatteryLevel: 80,
    }
    
    // ✅ Методи Car доступні напряму
    tesla.Drive()
    
    // ✅ Доступ до полів Car
    fmt.Printf("Brand: %s\n", tesla.Brand)
    fmt.Printf("HP: %d\n", tesla.Engine.Horsepower)
    
    // ✅ Власні методи ElectricCar
    tesla.Charge()
}
```

### Приклад: Композиція vs Успадкування

```go
// ❌ Погано (якби було успадкування)
// class Manager extends Employee {
//     manages []Employee
// }
// class CEO extends Manager { // глибока ієрархія!
// }

// ✅ Добре (композиція в Go)
type Employee struct {
    Name   string
    Salary float64
}

type Manager struct {
    Employee        // вбудовуємо Employee
    Manages  []Employee
}

type CEO struct {
    Manager         // вбудовуємо Manager
    Company string
}

func main() {
    ceo := CEO{
        Manager: Manager{
            Employee: Employee{
                Name:   "John Doe",
                Salary: 500000,
            },
            Manages: []Employee{
                {Name: "Alice", Salary: 100000},
                {Name: "Bob", Salary: 100000},
            },
        },
        Company: "TechCorp",
    }
    
    // Доступ до всіх полів
    fmt.Println(ceo.Name)    // з Employee
    fmt.Println(ceo.Company) // з CEO
    fmt.Println(len(ceo.Manages)) // з Manager
}
```

### Interface Composition

```go
// Reader - читання
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Writer - запис
type Writer interface {
    Write(p []byte) (n int, err error)
}

// ReadWriter - композиція інтерфейсів
type ReadWriter interface {
    Reader  // вбудовуємо Reader
    Writer  // вбудовуємо Writer
}

// File реалізує ReadWriter
type File struct {
    name string
}

func (f *File) Read(p []byte) (n int, err error) {
    fmt.Println("Reading from file")
    return 0, nil
}

func (f *File) Write(p []byte) (n int, err error) {
    fmt.Println("Writing to file")
    return 0, nil
}

func ProcessData(rw ReadWriter) {
    rw.Read(nil)
    rw.Write(nil)
}

func main() {
    f := &File{name: "data.txt"}
    ProcessData(f) // File реалізує ReadWriter
}
```

### Переваги композиції

1. **Гнучкість** - легко комбінувати компоненти
2. **Простота** - немає глибоких ієрархій
3. **Повторне використання** - компоненти незалежні
4. **Тестування** - легко замінити компоненти на mocks

---

## 📊 Порівняльна таблиця

| Принцип | В Go | Приклад |
|---------|------|---------|
| **Інкапсуляція** | Регістр літери (великий/маленький) | `user.password` (private) |
| **Поліморфізм** | Інтерфейси | `PaymentProcessor` інтерфейс |
| **Абстракція** | Інтерфейси + реалізація | `Database` інтерфейс |
| **Успадкування** | ❌ Немає | — |
| **Композиція** | ✅ Embedding | `Car` містить `Engine` |

---

## ✅ Висновки

1. **Go - не класична ООП мова**, але підтримує всі принципи
2. **Інкапсуляція** через регістр літер (просто й елегантно)
3. **Поліморфізм** через інтерфейси (duck typing)
4. **Абстракція** через інтерфейси
5. **Композиція > Успадкування** - гнучкість і простота

---

## 📚 Додаткове читання

- [Effective Go](https://go.dev/doc/effective_go)
- [Go by Example: Interfaces](https://gobyexample.com/interfaces)
- [Composition vs Inheritance in Go](https://yourbasic.org/golang/inheritance-object-oriented/)

---

**Далі:** [02_design_patterns.md](./02_design_patterns.md)
