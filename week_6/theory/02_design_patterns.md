# Патерни Проєктування в Go

Патерни проєктування - це перевірені рішення типових проблем у розробці ПЗ.

---

## 📖 Зміст

1. [Creational Patterns (Породжуючі)](#creational-patterns-породжуючі)
2. [Structural Patterns (Структурні)](#structural-patterns-структурні)
3. [Behavioral Patterns (Поведінкові)](#behavioral-patterns-поведінкові)

---

## Creational Patterns (Породжуючі)

Ці патерни відповідають за створення об'єктів.

### 1. Singleton

**Мета:** Гарантувати, що клас має тільки один екземпляр.

```go
package main

import (
    "fmt"
    "sync"
)

// Database - singleton
type Database struct {
    connections int
}

var (
    instance *Database
    once     sync.Once
)

// GetInstance - єдина точка доступу
func GetInstance() *Database {
    once.Do(func() {
        fmt.Println("Creating database instance...")
        instance = &Database{connections: 0}
    })
    return instance
}

func (db *Database) Connect() {
    db.connections++
    fmt.Printf("Connected. Total connections: %d\n", db.connections)
}

func main() {
    // Всі отримують той самий екземпляр
    db1 := GetInstance()
    db1.Connect() // "Creating database instance..." + "Connected. Total connections: 1"
    
    db2 := GetInstance()
    db2.Connect() // "Connected. Total connections: 2" (без створення)
    
    fmt.Println(db1 == db2) // true - той самий об'єкт
}
```

**Коли використовувати:**
- Конфігурація додатку
- Логер
- Підключення до БД
- Кеш

---

### 2. Factory

**Мета:** Створення об'єктів без вказування конкретного класу.

```go
package main

import "fmt"

// Transport - інтерфейс
type Transport interface {
    Deliver() string
}

// Truck - вантажівка
type Truck struct{}

func (t *Truck) Deliver() string {
    return "Delivering by truck 🚚"
}

// Ship - корабель
type Ship struct{}

func (s *Ship) Deliver() string {
    return "Delivering by ship 🚢"
}

// Plane - літак
type Plane struct{}

func (p *Plane) Deliver() string {
    return "Delivering by plane ✈️"
}

// TransportFactory - фабрика
func TransportFactory(transportType string) Transport {
    switch transportType {
    case "truck":
        return &Truck{}
    case "ship":
        return &Ship{}
    case "plane":
        return &Plane{}
    default:
        return &Truck{} // default
    }
}

func main() {
    // Не знаємо конкретний тип, працюємо через інтерфейс
    transport1 := TransportFactory("truck")
    fmt.Println(transport1.Deliver()) // "Delivering by truck 🚚"
    
    transport2 := TransportFactory("ship")
    fmt.Println(transport2.Deliver()) // "Delivering by ship 🚢"
    
    transport3 := TransportFactory("plane")
    fmt.Println(transport3.Deliver()) // "Delivering by plane ✈️"
}
```

**Коли використовувати:**
- Створення різних типів об'єктів залежно від умов
- Ізоляція логіки створення

---

### 3. Builder

**Мета:** Покрокове створення складних об'єктів.

```go
package main

import "fmt"

// Computer - складний об'єкт
type Computer struct {
    CPU     string
    RAM     int
    Storage int
    GPU     string
    OS      string
}

// ComputerBuilder - будівельник
type ComputerBuilder struct {
    computer *Computer
}

// NewComputerBuilder - створення будівельника
func NewComputerBuilder() *ComputerBuilder {
    return &ComputerBuilder{
        computer: &Computer{},
    }
}

// SetCPU - встановлює CPU
func (b *ComputerBuilder) SetCPU(cpu string) *ComputerBuilder {
    b.computer.CPU = cpu
    return b // fluent interface
}

// SetRAM - встановлює RAM
func (b *ComputerBuilder) SetRAM(ram int) *ComputerBuilder {
    b.computer.RAM = ram
    return b
}

// SetStorage - встановлює Storage
func (b *ComputerBuilder) SetStorage(storage int) *ComputerBuilder {
    b.computer.Storage = storage
    return b
}

// SetGPU - встановлює GPU
func (b *ComputerBuilder) SetGPU(gpu string) *ComputerBuilder {
    b.computer.GPU = gpu
    return b
}

// SetOS - встановлює OS
func (b *ComputerBuilder) SetOS(os string) *ComputerBuilder {
    b.computer.OS = os
    return b
}

// Build - повертає готовий об'єкт
func (b *ComputerBuilder) Build() *Computer {
    return b.computer
}

func main() {
    // ✅ Fluent interface - зручний синтаксис
    gamingPC := NewComputerBuilder().
        SetCPU("Intel i9").
        SetRAM(32).
        SetStorage(1000).
        SetGPU("NVIDIA RTX 4090").
        SetOS("Windows 11").
        Build()
    
    fmt.Printf("Gaming PC: %+v\n", gamingPC)
    
    // ✅ Можна створити мінімальну конфігурацію
    officePC := NewComputerBuilder().
        SetCPU("Intel i3").
        SetRAM(8).
        SetOS("Ubuntu").
        Build()
    
    fmt.Printf("Office PC: %+v\n", officePC)
}
```

**Коли використовувати:**
- Складні об'єкти з багатьма параметрами
- Потрібна валідація під час створення
- Різні конфігурації одного об'єкта

---

## Structural Patterns (Структурні)

Ці патерни описують, як об'єднувати об'єкти в більші структури.

### 4. Adapter

**Мета:** Адаптування інтерфейсу одного класу під інтерфейс іншого.

```go
package main

import "fmt"

// ===== Target інтерфейс (те, що очікує клієнт) =====
type PaymentGateway interface {
    Pay(amount float64) string
}

// ===== Новий сервіс (сумісний з інтерфейсом) =====
type StripeService struct{}

func (s *StripeService) Pay(amount float64) string {
    return fmt.Sprintf("Paid $%.2f via Stripe", amount)
}

// ===== Старий сервіс (НЕсумісний з інтерфейсом) =====
type OldPayPalService struct{}

// ❌ Метод називається інакше та має інші параметри
func (o *OldPayPalService) SendPayment(dollars int, cents int) string {
    return fmt.Sprintf("Sent payment of $%d.%02d via PayPal", dollars, cents)
}

// ===== Adapter для старого сервісу =====
type PayPalAdapter struct {
    paypal *OldPayPalService
}

// ✅ Адаптуємо до PaymentGateway інтерфейсу
func (a *PayPalAdapter) Pay(amount float64) string {
    dollars := int(amount)
    cents := int((amount - float64(dollars)) * 100)
    return a.paypal.SendPayment(dollars, cents)
}

// ===== Клієнт працює тільки з PaymentGateway =====
func ProcessPayment(gateway PaymentGateway, amount float64) {
    result := gateway.Pay(amount)
    fmt.Println(result)
}

func main() {
    // ✅ Новий сервіс - працює напряму
    stripe := &StripeService{}
    ProcessPayment(stripe, 99.99)
    
    // ✅ Старий сервіс - через адаптер
    oldPayPal := &OldPayPalService{}
    paypalAdapter := &PayPalAdapter{paypal: oldPayPal}
    ProcessPayment(paypalAdapter, 49.95)
}
```

**Коли використовувати:**
- Інтеграція зі сторонніми бібліотеками
- Робота з legacy кодом
- Об'єднання несумісних інтерфейсів

---

### 5. Decorator

**Мета:** Динамічне додавання нової функціональності об'єкту.

```go
package main

import "fmt"

// Coffee - базовий інтерфейс
type Coffee interface {
    GetDescription() string
    GetCost() float64
}

// SimpleCoffee - базова кава
type SimpleCoffee struct{}

func (c *SimpleCoffee) GetDescription() string {
    return "Simple Coffee"
}

func (c *SimpleCoffee) GetCost() float64 {
    return 2.00
}

// ===== Decorators =====

// MilkDecorator - додає молоко
type MilkDecorator struct {
    coffee Coffee
}

func (m *MilkDecorator) GetDescription() string {
    return m.coffee.GetDescription() + ", Milk"
}

func (m *MilkDecorator) GetCost() float64 {
    return m.coffee.GetCost() + 0.50
}

// SugarDecorator - додає цукор
type SugarDecorator struct {
    coffee Coffee
}

func (s *SugarDecorator) GetDescription() string {
    return s.coffee.GetDescription() + ", Sugar"
}

func (s *SugarDecorator) GetCost() float64 {
    return s.coffee.GetCost() + 0.20
}

// VanillaDecorator - додає ваніль
type VanillaDecorator struct {
    coffee Coffee
}

func (v *VanillaDecorator) GetDescription() string {
    return v.coffee.GetDescription() + ", Vanilla"
}

func (v *VanillaDecorator) GetCost() float64 {
    return v.coffee.GetCost() + 0.80
}

func PrintCoffee(c Coffee) {
    fmt.Printf("%s | $%.2f\n", c.GetDescription(), c.GetCost())
}

func main() {
    // Базова кава
    coffee := &SimpleCoffee{}
    PrintCoffee(coffee) // "Simple Coffee | $2.00"
    
    // Кава з молоком
    milkCoffee := &MilkDecorator{coffee: coffee}
    PrintCoffee(milkCoffee) // "Simple Coffee, Milk | $2.50"
    
    // ✅ Кава з молоком, цукром та ваніллю
    fancyCoffee := &VanillaDecorator{
        coffee: &SugarDecorator{
            coffee: &MilkDecorator{
                coffee: coffee,
            },
        },
    }
    PrintCoffee(fancyCoffee) // "Simple Coffee, Milk, Sugar, Vanilla | $3.50"
}
```

**Коли використовувати:**
- Додавання функціональності без зміни коду
- Потрібна гнучка комбінація поведінок
- Middleware в HTTP серверах

---

### 6. Facade

**Мета:** Простий інтерфейс до складної системи.

```go
package main

import "fmt"

// ===== Складна підсистема =====

type CPU struct{}

func (c *CPU) Freeze() {
    fmt.Println("CPU: Freezing...")
}

func (c *CPU) Execute() {
    fmt.Println("CPU: Executing...")
}

type Memory struct{}

func (m *Memory) Load() {
    fmt.Println("Memory: Loading...")
}

type HardDrive struct{}

func (h *HardDrive) Read() {
    fmt.Println("HardDrive: Reading...")
}

// ===== Facade - простий інтерфейс =====

type ComputerFacade struct {
    cpu       *CPU
    memory    *Memory
    hardDrive *HardDrive
}

func NewComputerFacade() *ComputerFacade {
    return &ComputerFacade{
        cpu:       &CPU{},
        memory:    &Memory{},
        hardDrive: &HardDrive{},
    }
}

// ✅ Один простий метод замість багатьох викликів
func (c *ComputerFacade) Start() {
    fmt.Println("=== Starting Computer ===")
    c.cpu.Freeze()
    c.memory.Load()
    c.hardDrive.Read()
    c.cpu.Execute()
    fmt.Println("=== Computer Started ===")
}

func main() {
    // ❌ Без Facade - багато кроків
    // cpu := &CPU{}
    // memory := &Memory{}
    // hardDrive := &HardDrive{}
    // cpu.Freeze()
    // memory.Load()
    // hardDrive.Read()
    // cpu.Execute()
    
    // ✅ З Facade - один виклик
    computer := NewComputerFacade()
    computer.Start()
}
```

**Коли використовувати:**
- Спрощення роботи зі складною системою
- API для зовнішніх користувачів
- Приховання реалізації

---

## Behavioral Patterns (Поведінкові)

Ці патерни визначають взаємодію між об'єктами.

### 7. Strategy

**Мета:** Визначення сімейства алгоритмів та можливість їх заміни.

```go
package main

import "fmt"

// Strategy - інтерфейс стратегії
type PaymentStrategy interface {
    Pay(amount float64) string
}

// ===== Конкретні стратегії =====

type CreditCardStrategy struct {
    cardNumber string
}

func (c *CreditCardStrategy) Pay(amount float64) string {
    return fmt.Sprintf("Paid $%.2f with Credit Card %s", amount, c.cardNumber)
}

type PayPalStrategy struct {
    email string
}

func (p *PayPalStrategy) Pay(amount float64) string {
    return fmt.Sprintf("Paid $%.2f with PayPal %s", amount, p.email)
}

type BitcoinStrategy struct {
    wallet string
}

func (b *BitcoinStrategy) Pay(amount float64) string {
    return fmt.Sprintf("Paid $%.2f with Bitcoin %s", amount, b.wallet)
}

// ===== Context =====

type PaymentProcessor struct {
    strategy PaymentStrategy
}

// SetStrategy - зміна стратегії в runtime
func (p *PaymentProcessor) SetStrategy(strategy PaymentStrategy) {
    p.strategy = strategy
}

func (p *PaymentProcessor) ProcessPayment(amount float64) string {
    return p.strategy.Pay(amount)
}

func main() {
    processor := &PaymentProcessor{}
    
    // ✅ Використовуємо Credit Card стратегію
    processor.SetStrategy(&CreditCardStrategy{cardNumber: "**** 1234"})
    fmt.Println(processor.ProcessPayment(100.50))
    
    // ✅ Міняємо на PayPal стратегію
    processor.SetStrategy(&PayPalStrategy{email: "user@example.com"})
    fmt.Println(processor.ProcessPayment(75.00))
    
    // ✅ Міняємо на Bitcoin стратегію
    processor.SetStrategy(&BitcoinStrategy{wallet: "1A2b3C..."})
    fmt.Println(processor.ProcessPayment(200.00))
}
```

**Коли використовувати:**
- Різні варіанти алгоритму
- Потрібна зміна поведінки в runtime
- Уникнення великих if/switch блоків

---

### 8. Observer

**Мета:** Сповіщення залежних об'єктів про зміни стану.

```go
package main

import "fmt"

// Observer - інтерфейс спостерігача
type Observer interface {
    Update(message string)
}

// Subject - об'єкт, який спостерігають
type Subject struct {
    observers []Observer
}

// Attach - додати спостерігача
func (s *Subject) Attach(observer Observer) {
    s.observers = append(s.observers, observer)
}

// Notify - сповістити всіх спостерігачів
func (s *Subject) Notify(message string) {
    for _, observer := range s.observers {
        observer.Update(message)
    }
}

// ===== Конкретні спостерігачі =====

type EmailNotifier struct {
    email string
}

func (e *EmailNotifier) Update(message string) {
    fmt.Printf("📧 Email to %s: %s\n", e.email, message)
}

type SMSNotifier struct {
    phone string
}

func (s *SMSNotifier) Update(message string) {
    fmt.Printf("📱 SMS to %s: %s\n", s.phone, message)
}

type PushNotifier struct {
    deviceID string
}

func (p *PushNotifier) Update(message string) {
    fmt.Printf("🔔 Push to %s: %s\n", p.deviceID, message)
}

// ===== Publisher =====

type NewsPublisher struct {
    Subject
}

func (n *NewsPublisher) PublishNews(news string) {
    fmt.Printf("\n📰 Publishing news: %s\n", news)
    n.Notify(news)
}

func main() {
    publisher := &NewsPublisher{}
    
    // Підписуємо спостерігачів
    publisher.Attach(&EmailNotifier{email: "user@example.com"})
    publisher.Attach(&SMSNotifier{phone: "+1234567890"})
    publisher.Attach(&PushNotifier{deviceID: "device123"})
    
    // Публікуємо новини - всі отримають сповіщення
    publisher.PublishNews("Breaking: Go 1.23 Released!")
    publisher.PublishNews("New design patterns tutorial available")
}
```

**Коли використовувати:**
- Система подій
- Реактивне програмування
- Pub/Sub системи

---

### 9. Command

**Мета:** Інкапсуляція запиту як об'єкт.

```go
package main

import "fmt"

// Command - інтерфейс команди
type Command interface {
    Execute()
    Undo()
}

// ===== Receiver =====

type Light struct {
    isOn bool
}

func (l *Light) TurnOn() {
    l.isOn = true
    fmt.Println("💡 Light is ON")
}

func (l *Light) TurnOff() {
    l.isOn = false
    fmt.Println("💡 Light is OFF")
}

// ===== Конкретні команди =====

type LightOnCommand struct {
    light *Light
}

func (c *LightOnCommand) Execute() {
    c.light.TurnOn()
}

func (c *LightOnCommand) Undo() {
    c.light.TurnOff()
}

type LightOffCommand struct {
    light *Light
}

func (c *LightOffCommand) Execute() {
    c.light.TurnOff()
}

func (c *LightOffCommand) Undo() {
    c.light.TurnOn()
}

// ===== Invoker =====

type RemoteControl struct {
    command Command
    history []Command
}

func (r *RemoteControl) SetCommand(command Command) {
    r.command = command
}

func (r *RemoteControl) PressButton() {
    r.command.Execute()
    r.history = append(r.history, r.command)
}

func (r *RemoteControl) PressUndo() {
    if len(r.history) == 0 {
        fmt.Println("Nothing to undo")
        return
    }
    
    lastCommand := r.history[len(r.history)-1]
    lastCommand.Undo()
    r.history = r.history[:len(r.history)-1]
}

func main() {
    light := &Light{}
    remote := &RemoteControl{}
    
    onCommand := &LightOnCommand{light: light}
    offCommand := &LightOffCommand{light: light}
    
    // Увімкнути світло
    remote.SetCommand(onCommand)
    remote.PressButton() // "💡 Light is ON"
    
    // Вимкнути світло
    remote.SetCommand(offCommand)
    remote.PressButton() // "💡 Light is OFF"
    
    // ✅ Undo - скасувати останню команду
    remote.PressUndo() // "💡 Light is ON"
    remote.PressUndo() // "💡 Light is OFF"
}
```

**Коли використовувати:**
- Потрібна історія команд (Undo/Redo)
- Черга завдань
- Логування операцій

---

## 📊 Коли використовувати які патерни

| Патерн | Проблема | Приклад |
|--------|----------|---------|
| **Singleton** | Потрібен один екземпляр | Config, Logger, DB Connection |
| **Factory** | Не знаємо заздалегідь який тип створювати | Transport Factory (Truck/Ship/Plane) |
| **Builder** | Складний об'єкт з багатьома параметрами | HTTP Request Builder, Computer Builder |
| **Adapter** | Несумісні інтерфейси | Legacy API integration |
| **Decorator** | Додаткова функціональність | HTTP Middleware, Coffee additions |
| **Facade** | Спрощення складної системи | SDK, API Wrapper |
| **Strategy** | Різні алгоритми | Payment methods, Sorting algorithms |
| **Observer** | Реакція на події | Event system, Notifications |
| **Command** | Undo/Redo, черги | Text editor, Task queue |

---

## ✅ Best Practices

1. **Не переусердкуйте** - використовуйте патерни тільки коли вони потрібні
2. **Спочатку простота** - почніть з простого рішення, додайте патерн якщо потрібно
3. **Знайте компроміси** - кожен патерн додає складності
4. **Go way** - використовуйте ідіоматичний Go код
5. **Інтерфейси** - маленькі інтерфейси краще великих

---

## 📚 Додаткове читання

- [Go Patterns](https://github.com/tmrts/go-patterns)
- [Design Patterns in Golang](https://refactoring.guru/design-patterns/go)
- "Head First Design Patterns" by Eric Freeman

---

**Далі:** [03_net_http.md](./03_net_http.md)
