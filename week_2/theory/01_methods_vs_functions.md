# Methods vs Functions

## Що таке метод?

**Метод** - це функція з особливим **receiver** аргументом.

**Функція** - самостійна операція, не прив'язана до типу.

---

## Синтаксис

### Функція
```go
package main

import "fmt"

type Rectangle struct {
    Width  int
    Height int
}

// Функція - приймає Rectangle як звичайний параметр
func CalculateArea(r Rectangle) int {
    return r.Width * r.Height
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}
    area := CalculateArea(rect)
    fmt.Printf("Area (function): %d\n", area)
}
```

### Метод
```go
package main

import "fmt"

type Rectangle struct {
    Width  int
    Height int
}

// Метод - має receiver (r Rectangle)
func (r Rectangle) Area() int {
    return r.Width * r.Height
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}
    area := rect.Area()  // викликається на об'єкті!
    fmt.Printf("Area (method): %d\n", area)
}
```

---

## Основні відмінності

| Аспект | Функція | Метод |
|--------|---------|-------|
| **Синтаксис** | `func Name(params) result` | `func (receiver Type) Name(params) result` |
| **Виклик** | `FunctionName(arg)` | `object.MethodName()` |
| **Прив'язка** | Не прив'язана до типу | Прив'язана до типу через receiver |
| **Доступ до полів** | Через параметри | Через receiver напряму |

---

## Receiver - що це?

**Receiver** - це спеціальний параметр, який з'являється перед назвою методу.

```go
func (r Rectangle) Area() int {
    // r - це receiver
    // Rectangle - тип receiver
    return r.Width * r.Height
}
```

### Receiver може бути:
1. **Value receiver** - `func (r Rectangle) Method()`
2. **Pointer receiver** - `func (r *Rectangle) Method()`

---

## Приклад: Person struct

```go
package main

import "fmt"

type Person struct {
    FirstName string
    LastName  string
    Age       int
}

// Функція
func GetFullNameFunc(p Person) string {
    return p.FirstName + " " + p.LastName
}

// Метод
func (p Person) FullName() string {
    return p.FirstName + " " + p.LastName
}

// Ще методи
func (p Person) Greet() string {
    return fmt.Sprintf("Привіт, мене звати %s", p.FullName())
}

func (p Person) IsAdult() bool {
    return p.Age >= 18
}

func (p Person) Info() string {
    adult := "неповнолітній"
    if p.IsAdult() {
        adult = "повнолітній"
    }
    return fmt.Sprintf("%s, %d років (%s)", p.FullName(), p.Age, adult)
}

func main() {
    person := Person{
        FirstName: "Іван",
        LastName:  "Петренко",
        Age:       25,
    }
    
    // Функція
    fmt.Println("Функція:", GetFullNameFunc(person))
    
    // Методи
    fmt.Println("Метод:", person.FullName())
    fmt.Println(person.Greet())
    fmt.Println(person.Info())
}
```

**Вивід:**
```
Функція: Іван Петренко
Метод: Іван Петренко
Привіт, мене звати Іван Петренко
Іван Петренко, 25 років (повнолітній)
```

---

## Переваги методів

### 1. Виразніший код
```go
// Функції - багато параметрів
area := CalculateArea(rect)
perimeter := CalculatePerimeter(rect)
diagonal := CalculateDiagonal(rect)

// Методи - читабельніше
area := rect.Area()
perimeter := rect.Perimeter()
diagonal := rect.Diagonal()
```

### 2. Інкапсуляція
```go
type BankAccount struct {
    balance float64  // приватне поле
}

func (b *BankAccount) Deposit(amount float64) {
    if amount > 0 {
        b.balance += amount
    }
}

func (b BankAccount) Balance() float64 {
    return b.balance
}

// Не можна напряму змінити balance
// account.balance = 1000000  // не компілюється (якщо в іншому пакеті)
```

### 3. Логічне групування
```go
// Всі операції над Rectangle разом
func (r Rectangle) Area() int { ... }
func (r Rectangle) Perimeter() int { ... }
func (r *Rectangle) Scale(factor int) { ... }
func (r *Rectangle) Move(dx, dy int) { ... }

// В IDE легко знайти всі методи типу
```

---

## Value Receiver

**Value receiver** отримує копію значення.

```go
package main

import "fmt"

type Counter struct {
    count int
}

// Value receiver - отримує копію
func (c Counter) Increment() {
    c.count++  // змінює КОПІЮ, не оригінал!
}

func (c Counter) Value() int {
    return c.count
}

func main() {
    counter := Counter{count: 0}
    
    fmt.Printf("До: %d\n", counter.Value())
    counter.Increment()
    fmt.Printf("Після: %d\n", counter.Value())  // все ще 0!
}
```

**Вивід:**
```
До: 0
Після: 0
```

---

## Коли використовувати методи?

### ✅ Використовуйте методи коли:

1. **Операція логічно належить типу**
   ```go
   rect.Area()        // ✅ Добре
   user.FullName()    // ✅ Добре
   order.Total()      // ✅ Добре
   ```

2. **Потрібна інкапсуляція**
   ```go
   type Account struct {
       balance float64
   }
   
   func (a *Account) Deposit(amount float64) {
       // контрольована зміна balance
   }
   ```

3. **Робота з даними структури**
   ```go
   func (p Person) IsAdult() bool {
       return p.Age >= 18
   }
   ```

### ✅ Використовуйте функції коли:

1. **Операція над кількома типами**
   ```go
   func Max(a, b int) int { ... }
   func CopyFile(src, dst string) error { ... }
   ```

2. **Утиліти і хелпери**
   ```go
   func ParseDate(s string) (time.Time, error) { ... }
   func FormatJSON(v interface{}) string { ... }
   ```

3. **Конструктори**
   ```go
   func NewUser(name string, age int) *User {
       return &User{Name: name, Age: age}
   }
   ```

---

## Методи на різних типах

### На struct
```go
type Person struct {
    Name string
}

func (p Person) Greet() string {
    return "Привіт, " + p.Name
}
```

### На власному типі
```go
type MyInt int

func (m MyInt) Double() MyInt {
    return m * 2
}

func main() {
    var x MyInt = 5
    fmt.Println(x.Double())  // 10
}
```

### На slice type
```go
type IntSlice []int

func (s IntSlice) Sum() int {
    total := 0
    for _, v := range s {
        total += v
    }
    return total
}

func main() {
    nums := IntSlice{1, 2, 3, 4, 5}
    fmt.Println(nums.Sum())  // 15
}
```

---

## Обмеження

### ❌ Не можна додати метод до чужого типу

```go
// ❌ Це не працює!
// func (i int) Double() int {
//     return i * 2
// }

// ✅ Але можна обгорнути в свій тип
type MyInt int

func (i MyInt) Double() MyInt {
    return i * 2
}
```

### ❌ Receiver має бути в тому ж пакеті

```go
// Метод має бути в тому ж пакеті, що й тип
// Не можна додати метод до типу з іншого пакету
```

---

## Практичний приклад: Blog Post

```go
package main

import (
    "fmt"
    "strings"
    "time"
)

type BlogPost struct {
    Title     string
    Content   string
    Author    string
    CreatedAt time.Time
    Tags      []string
}

// Конструктор (функція!)
func NewBlogPost(title, content, author string) *BlogPost {
    return &BlogPost{
        Title:     title,
        Content:   content,
        Author:    author,
        CreatedAt: time.Now(),
        Tags:      []string{},
    }
}

// Методи
func (b *BlogPost) AddTag(tag string) {
    b.Tags = append(b.Tags, tag)
}

func (b BlogPost) HasTag(tag string) bool {
    for _, t := range b.Tags {
        if t == tag {
            return true
        }
    }
    return false
}

func (b BlogPost) Summary(maxLength int) string {
    if len(b.Content) <= maxLength {
        return b.Content
    }
    return b.Content[:maxLength] + "..."
}

func (b BlogPost) Display() string {
    return fmt.Sprintf(
        "📝 %s\nАвтор: %s\nДата: %s\nТеги: %s\n\n%s",
        b.Title,
        b.Author,
        b.CreatedAt.Format("2006-01-02"),
        strings.Join(b.Tags, ", "),
        b.Content,
    )
}

func main() {
    post := NewBlogPost(
        "Go Methods",
        "Methods в Go - це функції з receiver параметром. Вони роблять код більш виразним та організованим.",
        "Іван Петренко",
    )
    
    post.AddTag("go")
    post.AddTag("programming")
    post.AddTag("tutorial")
    
    fmt.Println(post.Display())
    fmt.Println("\n--- Короткий опис ---")
    fmt.Println(post.Summary(50))
    fmt.Printf("\nЄ тег 'go'? %t\n", post.HasTag("go"))
}
```

---

## Резюме

| Концепція | Опис |
|-----------|------|
| **Метод** | Функція з receiver, прив'язана до типу |
| **Receiver** | Спеціальний параметр перед назвою методу |
| **Value receiver** | `(r Type)` - отримує копію |
| **Pointer receiver** | `(r *Type)` - отримує вказівник |
| **Переваги** | Виразніший код, інкапсуляція, групування |

---

## Питання для самоперевірки

1. В чому різниця між методом і функцією?
2. Що таке receiver?
3. Чи можна додати метод до типу `int`?
4. Коли краще використовувати метод, а коли функцію?
5. Що станеться, якщо змінити дані у value receiver?

---

## Завдання

1. Створіть struct `Book` з полями: title, author, pages
2. Додайте методи:
   - `Info()` - повна інформація
   - `IsLong()` - чи книга довга (>300 сторінок)
3. Створіть кілька книг і викличте методи

---

## Корисні посилання

- [Go Tour - Methods](https://go.dev/tour/methods/1)
- [Effective Go - Methods](https://go.dev/doc/effective_go#methods)
- [Go by Example - Methods](https://gobyexample.com/methods)

