# Pointer Receivers

## Value Receiver vs Pointer Receiver

В Go є два типи receivers:
1. **Value receiver** - `func (r Type) Method()`
2. **Pointer receiver** - `func (r *Type) Method()`

---

## Value Receiver

**Value receiver** отримує **копію** значення.

```go
package main

import "fmt"

type Counter struct {
    count int
}

// Value receiver - працює з копією
func (c Counter) Increment() {
    c.count++  // змінює копію!
    fmt.Printf("Всередині Increment: %d\n", c.count)
}

func (c Counter) Value() int {
    return c.count
}

func main() {
    counter := Counter{count: 5}
    
    fmt.Printf("До: %d\n", counter.Value())
    counter.Increment()
    fmt.Printf("Після: %d\n", counter.Value())  // все ще 5!
}
```

**Вивід:**
```
До: 5
Всередині Increment: 6
Після: 5
```

❌ **Зміни не зберігаються!**

---

## Pointer Receiver

**Pointer receiver** отримує **вказівник** на значення.

```go
package main

import "fmt"

type Counter struct {
    count int
}

// Pointer receiver - працює з оригіналом
func (c *Counter) Increment() {
    c.count++  // змінює оригінал!
    fmt.Printf("Всередині Increment: %d\n", c.count)
}

func (c *Counter) Value() int {
    return c.count
}

func main() {
    counter := Counter{count: 5}
    
    fmt.Printf("До: %d\n", counter.Value())
    counter.Increment()
    fmt.Printf("Після: %d\n", counter.Value())  // тепер 6!
}
```

**Вивід:**
```
До: 5
Всередині Increment: 6
Після: 6
```

✅ **Зміни зберігаються!**

---

## Порівняння

### Value Receiver
```go
type Rectangle struct {
    Width  int
    Height int
}

// Value receiver - НЕ змінює оригінал
func (r Rectangle) Scale(factor int) {
    r.Width *= factor
    r.Height *= factor
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}
    rect.Scale(2)
    fmt.Printf("%+v\n", rect)  // {Width:10 Height:5} - не змінилось!
}
```

### Pointer Receiver
```go
type Rectangle struct {
    Width  int
    Height int
}

// Pointer receiver - змінює оригінал
func (r *Rectangle) Scale(factor int) {
    r.Width *= factor
    r.Height *= factor
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}
    rect.Scale(2)
    fmt.Printf("%+v\n", rect)  // {Width:20 Height:10} - змінилось!
}
```

---

## Коли використовувати Pointer Receiver?

### ✅ Використовуйте Pointer Receiver коли:

#### 1. Метод змінює дані
```go
type BankAccount struct {
    balance float64
}

// Змінює balance - потрібен pointer
func (b *BankAccount) Deposit(amount float64) {
    b.balance += amount
}

func (b *BankAccount) Withdraw(amount float64) {
    b.balance -= amount
}
```

#### 2. Struct великий (економія пам'яті)
```go
type LargeStruct struct {
    data [1000000]int
    // багато даних...
}

// Pointer - не копіюємо 1000000 елементів!
func (l *LargeStruct) Process() {
    // обробка...
}
```

#### 3. Консистентність (якщо один метод pointer, всі pointer)
```go
type User struct {
    Name  string
    Email string
}

// Якщо є хоч один pointer receiver...
func (u *User) UpdateEmail(email string) {
    u.Email = email
}

// ...краще зробити всі pointer receivers
func (u *User) FullInfo() string {
    return u.Name + " <" + u.Email + ">"
}
```

---

## Коли використовувати Value Receiver?

### ✅ Використовуйте Value Receiver коли:

#### 1. Метод НЕ змінює дані
```go
type Point struct {
    X, Y int
}

// Тільки читає дані
func (p Point) String() string {
    return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

func (p Point) DistanceFromOrigin() float64 {
    return math.Sqrt(float64(p.X*p.X + p.Y*p.Y))
}
```

#### 2. Struct маленький
```go
type Color struct {
    R, G, B byte
}

// Struct маленький - value receiver OK
func (c Color) Hex() string {
    return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}
```

#### 3. Потрібна незмінність (immutability)
```go
type Money struct {
    amount int
}

// Value receiver - immutable
func (m Money) Add(other Money) Money {
    return Money{amount: m.amount + other.amount}
}

func main() {
    m1 := Money{amount: 100}
    m2 := Money{amount: 50}
    m3 := m1.Add(m2)
    
    // m1 і m2 не змінилися
    fmt.Println(m1.amount)  // 100
    fmt.Println(m2.amount)  // 50
    fmt.Println(m3.amount)  // 150
}
```

---

## Правила вибору

| Критерій | Value Receiver | Pointer Receiver |
|----------|----------------|------------------|
| Змінює дані? | ❌ Ні | ✅ Так |
| Struct маленький? | ✅ Так | ❌ Ні |
| Тільки читання? | ✅ Так | Необов'язково |
| Потрібна незмінність? | ✅ Так | ❌ Ні |
| Інші методи pointer? | ❌ | ✅ Так (консистентність) |

### 🎯 Загальне правило:

> **Якщо сумніваєтесь - використовуйте pointer receiver**

---

## Автоматична конверсія

Go автоматично конвертує між value і pointer при виклику методів.

```go
type Rectangle struct {
    Width, Height int
}

func (r *Rectangle) Scale(factor int) {
    r.Width *= factor
    r.Height *= factor
}

func main() {
    // Value
    rect1 := Rectangle{Width: 10, Height: 5}
    rect1.Scale(2)  // Go автоматично: (&rect1).Scale(2)
    
    // Pointer
    rect2 := &Rectangle{Width: 10, Height: 5}
    rect2.Scale(2)  // працює як є
    
    fmt.Printf("rect1: %+v\n", rect1)
    fmt.Printf("rect2: %+v\n", rect2)
}
```

---

## Практичний приклад: User Management

```go
package main

import (
    "fmt"
    "strings"
    "time"
)

type User struct {
    ID        int
    Username  string
    Email     string
    CreatedAt time.Time
    IsActive  bool
}

// Конструктор (функція, не метод)
func NewUser(username, email string) *User {
    return &User{
        ID:        generateID(),
        Username:  username,
        Email:     email,
        CreatedAt: time.Now(),
        IsActive:  true,
    }
}

// Pointer receiver - змінює дані
func (u *User) Activate() {
    u.IsActive = true
}

func (u *User) Deactivate() {
    u.IsActive = false
}

func (u *User) UpdateEmail(email string) error {
    if !strings.Contains(email, "@") {
        return fmt.Errorf("invalid email: %s", email)
    }
    u.Email = email
    return nil
}

// Value receiver - тільки читання
func (u User) FullInfo() string {
    status := "активний"
    if !u.IsActive {
        status = "неактивний"
    }
    
    return fmt.Sprintf(
        "ID: %d\nUsername: %s\nEmail: %s\nСтатус: %s\nСтворено: %s",
        u.ID,
        u.Username,
        u.Email,
        status,
        u.CreatedAt.Format("2006-01-02 15:04"),
    )
}

func (u User) IsValid() bool {
    return u.Username != "" && 
           u.Email != "" && 
           strings.Contains(u.Email, "@")
}

// Допоміжна функція
var nextID = 1
func generateID() int {
    id := nextID
    nextID++
    return id
}

func main() {
    // Створення користувача
    user := NewUser("ivan_petro", "ivan@example.com")
    
    fmt.Println("=== Початковий стан ===")
    fmt.Println(user.FullInfo())
    
    // Зміна даних (pointer receivers)
    user.UpdateEmail("new_email@example.com")
    user.Deactivate()
    
    fmt.Println("\n=== Після змін ===")
    fmt.Println(user.FullInfo())
    
    // Перевірка (value receiver)
    fmt.Printf("\nValid? %t\n", user.IsValid())
}
```

---

## Map, Slice, Chan - особливий випадок

❗ **Map, Slice, Channel** - це вже **reference types**.

Їх НЕ потрібно передавати через pointer!

```go
type UserList struct {
    users []User  // slice - вже reference type
}

// ✅ Добре - value receiver для slice
func (ul UserList) Add(user User) {
    ul.users = append(ul.users, user)  // append може змінити slice
}

// ❌ Не потрібно - зайвий pointer
func (ul *UserList) AddPointer(user User) {
    ul.users = append(ul.users, user)
}
```

**Але:** якщо потрібна консистентність з іншими методами, можна використати pointer.

---

## Інтерфейси і Pointer Receivers

⚠️ **Важливо:** Тип з pointer receiver НЕ задовольняє інтерфейс для value!

```go
type Writer interface {
    Write(data string)
}

type FileWriter struct {
    filename string
}

// Pointer receiver
func (f *FileWriter) Write(data string) {
    fmt.Printf("Writing to %s: %s\n", f.filename, data)
}

func SaveData(w Writer, data string) {
    w.Write(data)
}

func main() {
    // ✅ Працює - передаємо pointer
    fw := &FileWriter{filename: "data.txt"}
    SaveData(fw, "hello")
    
    // ❌ НЕ працює - value не реалізує Writer
    // fw2 := FileWriter{filename: "data.txt"}
    // SaveData(fw2, "hello")  // compilation error!
}
```

---

## Помилки початківців

### ❌ Помилка 1: Забули pointer receiver
```go
type Counter struct {
    count int
}

// ❌ Погано - не змінить оригінал
func (c Counter) Increment() {
    c.count++
}

// ✅ Добре
func (c *Counter) Increment() {
    c.count++
}
```

### ❌ Помилка 2: Pointer receiver для маленького read-only struct
```go
type Point struct {
    X, Y int
}

// ❌ Погано - не потрібен pointer
func (p *Point) String() string {
    return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

// ✅ Добре - value receiver
func (p Point) String() string {
    return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}
```

### ❌ Помилка 3: Мікс value і pointer receivers без причини
```go
// ❌ Погано - непослідовно
func (u User) GetName() string { ... }      // value
func (u *User) GetEmail() string { ... }    // pointer
func (u User) GetAge() int { ... }          // value

// ✅ Добре - послідовно
func (u User) GetName() string { ... }      // всі value
func (u User) GetEmail() string { ... }
func (u User) GetAge() int { ... }

// АБО всі pointer (якщо є методи що змінюють)
func (u *User) GetName() string { ... }
func (u *User) SetEmail(email string) { ... }
```

---

## Best Practices

### 1. Консистентність
```go
// ✅ Добре - всі pointer receivers
type User struct { ... }
func (u *User) SetName(name string) { u.Name = name }
func (u *User) GetName() string { return u.Name }
func (u *User) String() string { return u.Name }

// ✅ Добре - всі value receivers (immutable)
type Point struct { X, Y int }
func (p Point) Add(other Point) Point { ... }
func (p Point) String() string { ... }
func (p Point) Distance() float64 { ... }
```

### 2. Документуйте рішення
```go
// Point - immutable point in 2D space.
// All methods use value receivers to maintain immutability.
type Point struct {
    X, Y int
}
```

### 3. Думайте про інтерфейси
```go
// Якщо тип має реалізувати інтерфейс,
// подумайте про receiver type
type Writer interface {
    Write(data []byte) error
}

// Pointer receiver якщо потрібен state
func (f *FileWriter) Write(data []byte) error { ... }

// Value receiver якщо не потрібен state
func (f NoOpWriter) Write(data []byte) error { return nil }
```

---

## Резюме

| Концепція | Опис |
|-----------|------|
| **Value receiver** | `(r Type)` - отримує копію |
| **Pointer receiver** | `(r *Type)` - отримує вказівник |
| **Зміна даних** | Потрібен pointer receiver |
| **Великий struct** | Краще pointer receiver |
| **Маленький struct** | Value receiver OK |
| **Консистентність** | Якщо один pointer - краще всі pointer |

---

## Питання для самоперевірки

1. В чому різниця між value і pointer receiver?
2. Коли обов'язково потрібен pointer receiver?
3. Чи можна змінити дані через value receiver?
4. Що таке "автоматична конверсія" receivers?
5. Чому map/slice не потребують pointer receiver?

---

## Завдання

1. Створіть struct `BankAccount` з полем `balance`
2. Додайте методи `Deposit()` і `Withdraw()` (pointer receivers)
3. Додайте метод `Balance()` для читання (value receiver)
4. Протестуйте, що зміни зберігаються

---

## Корисні посилання

- [Go FAQ - Should I define methods on values or pointers?](https://go.dev/doc/faq#methods_on_values_or_pointers)
- [Effective Go - Pointers vs Values](https://go.dev/doc/effective_go#pointers_vs_values)
- [Go Tour - Pointer receivers](https://go.dev/tour/methods/4)

