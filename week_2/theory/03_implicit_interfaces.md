# Implicit Interfaces (Неявні інтерфейси)

## Що таке інтерфейс?

**Інтерфейс** - це набір сигнатур методів. Тип реалізує інтерфейс, якщо має всі ці методи.

**Ключова особливість Go:** реалізація інтерфейсів **неявна** (implicit).

---

## Explicit vs Implicit

### Explicit (Java, C#, TypeScript)
```java
// Java - явна реалізація
interface Writer {
    void write(String data);
}

class FileWriter implements Writer {  // ← явно вказуємо
    public void write(String data) {
        // ...
    }
}
```

### Implicit (Go)
```go
// Go - неявна реалізація
type Writer interface {
    Write(data string)
}

type FileWriter struct{}

// Просто маємо метод Write - автоматично реалізує Writer!
func (f FileWriter) Write(data string) {
    // ...
}

// FileWriter реалізує Writer БЕЗ явного оголошення!
```

---

## Базовий приклад

```go
package main

import "fmt"

// Оголошення інтерфейсу
type Greeter interface {
    Greet() string
}

// Тип 1
type Person struct {
    Name string
}

func (p Person) Greet() string {
    return "Привіт, я " + p.Name
}

// Тип 2
type Dog struct {
    Name string
}

func (d Dog) Greet() string {
    return "Гав! Я " + d.Name
}

// Функція приймає інтерфейс
func SayHello(g Greeter) {
    fmt.Println(g.Greet())
}

func main() {
    person := Person{Name: "Іван"}
    dog := Dog{Name: "Рекс"}
    
    // Обидва типи реалізують Greeter!
    SayHello(person)  // Привіт, я Іван
    SayHello(dog)     // Гав! Я Рекс
}
```

---

## Переваги неявної реалізації

### 1. Гнучкість

Можна створити інтерфейс для існуючого коду:

```go
// Існуючий код в іншому пакеті
type HTTPClient struct{}
func (h HTTPClient) Get(url string) ([]byte, error) { ... }

// Ваш код - створюєте інтерфейс
type Getter interface {
    Get(url string) ([]byte, error)
}

// HTTPClient автоматично реалізує Getter!
func FetchData(g Getter, url string) ([]byte, error) {
    return g.Get(url)
}
```

### 2. Dependency Inversion

```go
type Database interface {
    Save(data string) error
    Load(id int) (string, error)
}

// Реальна база даних
type PostgresDB struct{}
func (p PostgresDB) Save(data string) error { ... }
func (p PostgresDB) Load(id int) (string, error) { ... }

// Mock для тестів
type MockDB struct{}
func (m MockDB) Save(data string) error { return nil }
func (m MockDB) Load(id int) (string, error) { return "test", nil }

// Сервіс залежить від інтерфейсу, не від конкретного типу
type UserService struct {
    db Database  // інтерфейс!
}

func main() {
    // В продакшені
    service := UserService{db: PostgresDB{}}
    
    // В тестах
    service = UserService{db: MockDB{}}
}
```

### 3. Decoupling

Пакети не залежать один від одного:

```go
// package storage
type Storage interface {
    Save(key, value string) error
}

// package redis
type RedisClient struct{}
func (r RedisClient) Save(key, value string) error { ... }

// package mysql
type MySQLClient struct{}
func (m MySQLClient) Save(key, value string) error { ... }

// Redis і MySQL не знають про Storage інтерфейс!
// Але обидва його реалізують.
```

---

## Polymorphism (Поліморфізм)

```go
package main

import (
    "fmt"
    "math"
)

// Інтерфейс
type Shape interface {
    Area() float64
    Perimeter() float64
}

// Прямокутник
type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// Коло
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}

// Трикутник
type Triangle struct {
    A, B, C float64
}

func (t Triangle) Area() float64 {
    s := (t.A + t.B + t.C) / 2
    return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C))
}

func (t Triangle) Perimeter() float64 {
    return t.A + t.B + t.C
}

// Функція працює з будь-якою фігурою
func PrintShapeInfo(s Shape) {
    fmt.Printf("Area: %.2f\n", s.Area())
    fmt.Printf("Perimeter: %.2f\n", s.Perimeter())
    fmt.Println()
}

// Загальна площа всіх фігур
func TotalArea(shapes []Shape) float64 {
    total := 0.0
    for _, shape := range shapes {
        total += shape.Area()
    }
    return total
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}
    circle := Circle{Radius: 7}
    triangle := Triangle{A: 3, B: 4, C: 5}
    
    fmt.Println("=== Rectangle ===")
    PrintShapeInfo(rect)
    
    fmt.Println("=== Circle ===")
    PrintShapeInfo(circle)
    
    fmt.Println("=== Triangle ===")
    PrintShapeInfo(triangle)
    
    // Slice різних фігур
    shapes := []Shape{rect, circle, triangle}
    fmt.Printf("Total area: %.2f\n", TotalArea(shapes))
}
```

---

## Empty Interface `interface{}`

`interface{}` (або `any` в Go 1.18+) - інтерфейс без методів.

**Будь-який** тип реалізує порожній інтерфейс!

```go
package main

import "fmt"

func PrintAnything(v interface{}) {
    fmt.Printf("Value: %v, Type: %T\n", v, v)
}

func main() {
    PrintAnything(42)           // int
    PrintAnything("hello")      // string
    PrintAnything(true)         // bool
    PrintAnything([]int{1, 2})  // []int
    PrintAnything(struct{ Name string }{"John"})  // struct
}
```

**Вивід:**
```
Value: 42, Type: int
Value: hello, Type: string
Value: true, Type: bool
Value: [1 2], Type: []int
Value: {John}, Type: struct { Name string }
```

---

## Type Assertions

**Type assertion** - отримання конкретного типу з інтерфейсу.

### Синтаксис 1: Panic on failure
```go
value := interfaceVar.(ConcreteType)
```

### Синтаксис 2: Safe (з перевіркою)
```go
value, ok := interfaceVar.(ConcreteType)
if ok {
    // value має тип ConcreteType
}
```

### Приклад
```go
package main

import "fmt"

func Describe(i interface{}) {
    // Type assertion з перевіркою
    if s, ok := i.(string); ok {
        fmt.Printf("String: %s (length: %d)\n", s, len(s))
        return
    }
    
    if n, ok := i.(int); ok {
        fmt.Printf("Integer: %d (double: %d)\n", n, n*2)
        return
    }
    
    fmt.Printf("Unknown type: %T\n", i)
}

func main() {
    Describe("hello")
    Describe(42)
    Describe(true)
}
```

**Вивід:**
```
String: hello (length: 5)
Integer: 42 (double: 84)
Unknown type: bool
```

---

## Type Switch

Елегантний спосіб обробки різних типів:

```go
package main

import "fmt"

func ProcessValue(v interface{}) {
    switch value := v.(type) {
    case string:
        fmt.Printf("String: %s\n", value)
    case int:
        fmt.Printf("Integer: %d\n", value)
    case bool:
        fmt.Printf("Boolean: %t\n", value)
    case []int:
        fmt.Printf("Int slice: %v (sum: %d)\n", value, sum(value))
    default:
        fmt.Printf("Unknown type: %T\n", value)
    }
}

func sum(nums []int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

func main() {
    ProcessValue("Go")
    ProcessValue(2024)
    ProcessValue(true)
    ProcessValue([]int{1, 2, 3, 4, 5})
    ProcessValue(3.14)
}
```

---

## Стандартні інтерфейси Go

### io.Reader
```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Багато типів реалізують Reader:
// - os.File
// - strings.Reader
// - bytes.Buffer
// - http.Response.Body
```

### io.Writer
```go
type Writer interface {
    Write(p []byte) (n int, err error)
}

// Багато типів реалізують Writer:
// - os.File
// - os.Stdout
// - bytes.Buffer
// - http.ResponseWriter
```

### fmt.Stringer
```go
type Stringer interface {
    String() string
}

type Person struct {
    Name string
    Age  int
}

func (p Person) String() string {
    return fmt.Sprintf("%s (%d років)", p.Name, p.Age)
}

func main() {
    person := Person{Name: "Іван", Age: 25}
    fmt.Println(person)  // автоматично викликає String()
}
```

### error
```go
type error interface {
    Error() string
}

type MyError struct {
    Code    int
    Message string
}

func (e MyError) Error() string {
    return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}
```

---

## Практичний приклад: UserService

```go
package main

import (
    "fmt"
    "errors"
)

// ============= Domain Model =============

type User struct {
    ID    int
    Name  string
    Email string
}

// ============= Interface =============

type UserService interface {
    Create(name, email string) (*User, error)
    GetByID(id int) (*User, error)
    GetAll() ([]*User, error)
    Delete(id int) error
}

// ============= In-Memory Implementation =============

type InMemoryUserService struct {
    users  []*User
    nextID int
}

func NewInMemoryUserService() *InMemoryUserService {
    return &InMemoryUserService{
        users:  []*User{},
        nextID: 1,
    }
}

func (s *InMemoryUserService) Create(name, email string) (*User, error) {
    user := &User{
        ID:    s.nextID,
        Name:  name,
        Email: email,
    }
    s.nextID++
    s.users = append(s.users, user)
    return user, nil
}

func (s *InMemoryUserService) GetByID(id int) (*User, error) {
    for _, user := range s.users {
        if user.ID == id {
            return user, nil
        }
    }
    return nil, errors.New("user not found")
}

func (s *InMemoryUserService) GetAll() ([]*User, error) {
    return s.users, nil
}

func (s *InMemoryUserService) Delete(id int) error {
    for i, user := range s.users {
        if user.ID == id {
            s.users = append(s.users[:i], s.users[i+1:]...)
            return nil
        }
    }
    return errors.New("user not found")
}

// ============= Mock Implementation =============

type MockUserService struct {
    CreateCalled bool
    GetCalled    bool
}

func NewMockUserService() *MockUserService {
    return &MockUserService{}
}

func (m *MockUserService) Create(name, email string) (*User, error) {
    m.CreateCalled = true
    return &User{ID: 999, Name: name, Email: email}, nil
}

func (m *MockUserService) GetByID(id int) (*User, error) {
    m.GetCalled = true
    return &User{ID: id, Name: "Mock User", Email: "mock@example.com"}, nil
}

func (m *MockUserService) GetAll() ([]*User, error) {
    return []*User{
        {ID: 1, Name: "Mock 1", Email: "mock1@example.com"},
        {ID: 2, Name: "Mock 2", Email: "mock2@example.com"},
    }, nil
}

func (m *MockUserService) Delete(id int) error {
    return nil
}

// ============= Application Layer =============

type Application struct {
    userService UserService  // залежність від інтерфейсу!
}

func (app *Application) RegisterUser(name, email string) {
    user, err := app.userService.Create(name, email)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("User registered: %+v\n", user)
}

func (app *Application) ListUsers() {
    users, err := app.userService.GetAll()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Println("=== All Users ===")
    for _, user := range users {
        fmt.Printf("  [%d] %s <%s>\n", user.ID, user.Name, user.Email)
    }
}

// ============= Main =============

func main() {
    fmt.Println("=== Using In-Memory Service ===")
    inMemory := NewInMemoryUserService()
    app1 := Application{userService: inMemory}
    
    app1.RegisterUser("Іван", "ivan@example.com")
    app1.RegisterUser("Марія", "maria@example.com")
    app1.ListUsers()
    
    fmt.Println("\n=== Using Mock Service ===")
    mock := NewMockUserService()
    app2 := Application{userService: mock}
    
    app2.RegisterUser("Test", "test@example.com")
    app2.ListUsers()
    
    fmt.Printf("\nMock statistics: Create called = %t, Get called = %t\n", 
        mock.CreateCalled, mock.GetCalled)
}
```

---

## Interface Composition

Інтерфейси можна комбінувати:

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type Closer interface {
    Close() error
}

// Композиція інтерфейсів
type ReadWriter interface {
    Reader
    Writer
}

type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}
```

---

## Best Practices

### 1. Accept interfaces, return structs

```go
// ✅ Добре
func NewUserService(db Database) *UserService {
    return &UserService{db: db}  // повертаємо struct
}

func ProcessData(r io.Reader) error {  // приймаємо interface
    // ...
}

// ❌ Погано
func NewUserService(db Database) Database {  // не повертайте interface
    return &UserService{db: db}
}
```

### 2. Keep interfaces small

```go
// ✅ Добре - маленькі, focused
type Reader interface {
    Read(p []byte) (n int, err error)
}

// ❌ Погано - занадто великий
type SuperInterface interface {
    Method1()
    Method2()
    Method3()
    Method4()
    Method5()
    // ...
}
```

### 3. Interface in consumer package

```go
// ✅ Добре - інтерфейс там, де використовується
// package handler
type UserGetter interface {
    GetUser(id int) (*User, error)
}

func HandleRequest(ug UserGetter) {
    // ...
}

// ❌ Погано - інтерфейс в producer package
// package database
type UserGetter interface {  // не тут!
    GetUser(id int) (*User, error)
}
```

---

## Перевірка реалізації інтерфейсу

```go
type Writer interface {
    Write(data []byte) error
}

type FileWriter struct{}

func (f FileWriter) Write(data []byte) error {
    return nil
}

// Compile-time перевірка
var _ Writer = FileWriter{}        // якщо не реалізує - compilation error
var _ Writer = (*FileWriter)(nil)  // для pointer receiver
```

---

## Порожній інтерфейс - коли використовувати?

### ✅ Добре: Generic контейнери (до Go 1.18)
```go
func PrintSlice(items []interface{}) {
    for _, item := range items {
        fmt.Println(item)
    }
}
```

### ✅ Добре: JSON unmarshaling
```go
var data interface{}
json.Unmarshal(jsonBytes, &data)
```

### ❌ Погано: Замість конкретного типу
```go
// ❌ Погано
func ProcessUser(u interface{}) { ... }

// ✅ Добре
func ProcessUser(u User) { ... }
func ProcessUser(u UserProvider) { ... }  // специфічний interface
```

---

## Резюме

| Концепція | Опис |
|-----------|------|
| **Implicit implementation** | Не потрібно явно вказувати `implements` |
| **Duck typing** | "Якщо виглядає як качка і крякає як качка..." |
| **Polymorphism** | Різні типи через один інтерфейс |
| **Type assertion** | `value, ok := i.(Type)` |
| **Type switch** | `switch v := i.(type)` |
| **Empty interface** | `interface{}` або `any` |

---

## Питання для самоперевірки

1. Чому інтерфейси не реалізуються явно в Go?
2. Які переваги неявної реалізації?
3. Що таке "duck typing"?
4. Як перевірити, чи тип реалізує інтерфейс?
5. Коли використовувати empty interface?
6. В чому різниця між type assertion і type switch?

---

## Завдання

1. Створіть інтерфейс `Storage` з методами `Save()` та `Load()`
2. Зробіть дві реалізації:
   - `MemoryStorage` (в пам'яті)
   - `MockStorage` (для тестів)
3. Створіть функцію, яка приймає `Storage` і працює з обома реалізаціями

---

## Корисні посилання

- [Go Tour - Interfaces](https://go.dev/tour/methods/9)
- [Effective Go - Interfaces](https://go.dev/doc/effective_go#interfaces)
- [How to use interfaces in Go](https://jordanorelli.com/post/32665860244/how-to-use-interfaces-in-go)

