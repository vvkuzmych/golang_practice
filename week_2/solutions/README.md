# Рішення вправ - Тиждень 2

Тут знаходяться рішення всіх вправ тижня 2.

---

## 📁 Файли

- `solution_1.go` - Calculator з методами
- `solution_2.go` - Shape Interface з polymorphism
- `solution_3.go` - Storage Interface з різними реалізаціями

---

## 🚀 Як запустити

### Solution 1: Calculator
```bash
cd solutions
go run solution_1.go
```

**Що демонструє:**
- Методи на struct
- Pointer receivers для зміни стану
- Value receivers для читання
- Обробка помилок в методах
- Chainable методи
- Історія операцій

### Solution 2: Shape Interface
```bash
go run solution_2.go
```

**Що демонструє:**
- Оголошення інтерфейсу
- Неявна реалізація інтерфейсу
- Polymorphism (різні типи через один інтерфейс)
- Робота зі slice інтерфейсів
- Type assertions і type switch
- Фільтрація та сортування

### Solution 3: Storage Interface
```bash
go run solution_3.go
```

**Що демонструє:**
- Абстракція через інтерфейси
- Memory storage
- File storage (персистентність)
- Mock storage (для тестів)
- Dependency Injection
- Легка зміна реалізації

---

## 📊 Порівняння вправ

| Вправа | Складність | Концепції | Час |
|--------|-----------|-----------|-----|
| Solution 1 | ⭐⭐ | Methods, Receivers | 30-45 хв |
| Solution 2 | ⭐⭐⭐ | Interfaces, Polymorphism | 45-60 хв |
| Solution 3 | ⭐⭐⭐⭐ | Abstraction, DI, File I/O | 60-90 хв |

---

## 💡 Ключові моменти

### Solution 1: Calculator
```go
// Pointer receiver для зміни
func (c *Calculator) Add(value float64) *Calculator {
    c.result += value
    return c  // для chaining
}

// Value receiver для читання
func (c Calculator) Result() float64 {
    return c.result
}
```

### Solution 2: Shape Interface
```go
// Інтерфейс
type Shape interface {
    Area() float64
    Perimeter() float64
    Name() string
}

// Неявна реалізація
type Circle struct { Radius float64 }
func (c Circle) Area() float64 { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }
func (c Circle) Name() string { return "Circle" }

// Circle автоматично реалізує Shape!
```

### Solution 3: Storage Interface
```go
// Інтерфейс
type Storage interface {
    Save(key, value string) error
    Load(key string) (string, error)
    // ...
}

// DataManager не знає про конкретну реалізацію
type DataManager struct {
    storage Storage  // залежність від інтерфейсу!
}

// Можна передати будь-яку реалізацію:
manager1 := NewDataManager(NewMemoryStorage())
manager2 := NewDataManager(NewFileStorage("data.txt"))
manager3 := NewDataManager(NewMockStorage())
```

---

## 🎯 Що ви навчились

### Methods (Solution 1)
- ✅ Різниця між методами та функціями
- ✅ Value vs Pointer receivers
- ✅ Коли використовувати кожен тип
- ✅ Chainable методи
- ✅ Обробка помилок

### Interfaces (Solution 2)
- ✅ Оголошення інтерфейсів
- ✅ Неявна реалізація
- ✅ Polymorphism
- ✅ Type assertions
- ✅ Type switch
- ✅ Робота зі slice інтерфейсів

### Architecture (Solution 3)
- ✅ Абстракція через інтерфейси
- ✅ Dependency Injection
- ✅ Multiple implementations
- ✅ Mock для тестів
- ✅ Loose coupling
- ✅ File I/O

---

## 🔍 Важливі патерни

### 1. Pointer Receiver Pattern
```go
// ✅ Використовуйте pointer коли:
// 1. Змінюєте дані
func (c *Calculator) Add(v float64) { c.result += v }

// 2. Struct великий
func (b *BigStruct) Process() { /* ... */ }

// 3. Для консистентності
type User struct { /* ... */ }
func (u *User) Method1() { /* ... */ }
func (u *User) Method2() { /* ... */ }  // всі pointer
```

### 2. Interface Pattern
```go
// ✅ Маленькі, focused інтерфейси
type Reader interface {
    Read(p []byte) (n int, err error)
}

// ❌ Не робіть великі інтерфейси
// type SuperInterface interface {
//     Method1()
//     Method2()
//     // ... 20 методів
// }
```

### 3. Dependency Injection Pattern
```go
// ✅ Залежність від інтерфейсу
type Service struct {
    storage Storage  // interface!
}

// Можна підставити будь-яку реалізацію
service := Service{storage: MemoryStorage{}}
service := Service{storage: FileStorage{}}
```

---

## 📚 Додаткові ресурси

### Про Methods
- [Go Tour - Methods](https://go.dev/tour/methods/)
- [Effective Go - Methods](https://go.dev/doc/effective_go#methods)

### Про Interfaces
- [Go Tour - Interfaces](https://go.dev/tour/methods/9)
- [How to use interfaces in Go](https://jordanorelli.com/post/32665860244/how-to-use-interfaces-in-go)

### Про Architecture
- [Accept interfaces, return structs](https://bryanftan.medium.com/accept-interfaces-return-structs-in-go-d4cab29a301b)
- [Practical Go: Real world advice](https://dave.cheney.net/practical-go/presentations/qcon-china.html)

---

## 🎓 Наступні кроки

Після розбору рішень:

1. **Модифікуйте код:**
   - Додайте нові методи в Calculator
   - Створіть нові фігури в Shape
   - Додайте Database Storage

2. **Експериментуйте:**
   - Змініть pointer на value receivers і навпаки
   - Створіть власні інтерфейси
   - Комбінуйте різні реалізації

3. **Покращуйте:**
   - Додайте більше валідації
   - Покращіть обробку помилок
   - Додайте логування

4. **Тестуйте:**
   - Напишіть unit tests
   - Використайте Mock для тестування
   - Перевірте edge cases

---

## ❓ FAQ

### Q: Чому в Solution 1 використовується pointer receiver?
A: Бо методи змінюють поле `result`. Без pointer змінилася б тільки копія.

### Q: Чому в Solution 2 не треба явно вказувати `implements`?
A: Go використовує неявну реалізацію. Якщо тип має всі методи інтерфейсу - він його реалізує автоматично.

### Q: Навіщо Mock Storage в Solution 3?
A: Для тестування. Mock дозволяє контролювати поведінку та перевіряти виклики без реальних операцій.

### Q: Коли використовувати File Storage vs Memory Storage?
A: Memory - швидкий але не зберігає дані після перезапуску. File - повільніший але персистентний.

---

## 🎉 Вітаємо!

Ви пройшли всі вправи тижня 2!

Тепер ви розумієте:
- ✅ Methods і Receivers
- ✅ Interfaces і Polymorphism
- ✅ Dependency Injection
- ✅ Abstraction через інтерфейси

**Готові до наступного тижня!** 🚀

