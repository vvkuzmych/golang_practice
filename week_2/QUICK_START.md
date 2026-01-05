# Швидкий старт - Тиждень 2

## 🚀 Як почати

### 1. Перейти в папку week_2
```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_2
```

### 2. Прочитати README
```bash
cat README.md
```

---

## 📚 Порядок навчання

### День 1-2: Теорія

```bash
# 1. Methods vs Functions
cat theory/01_methods_vs_functions.md

# 2. Pointer Receivers
cat theory/02_pointer_receivers.md

# 3. Implicit Interfaces
cat theory/03_implicit_interfaces.md
```

### День 3-4: Практика

```bash
# 1. Methods Demo
cd practice/methods_demo
go run main.go
cat main.go

# 2. Interface Demo
cd ../interface_demo
go run main.go
cat main.go

# 3. UserService - головний приклад
cd ../user_service
go run main.go
cat main.go
```

### День 5-6: Вправи

```bash
# Прочитати завдання
cd ../../exercises
cat exercise_1.md  # Calculator з методами
cat exercise_2.md  # Shape інтерфейс
cat exercise_3.md  # Storage інтерфейс

# Створити файли і виконати вправи
# my_exercise_1.go
# my_exercise_2.go
# my_exercise_3.go

# Перевірити рішення
cd ../solutions
cat solution_1.go
cat solution_2.go
cat solution_3.go
```

### День 7: Контроль знань

Відповісти на питання:

#### 1. Methods vs Functions
**Q: В чому різниця?**
- Метод прив'язаний до типу через receiver
- Функція - самостійна

**Q: Приклад?**
```go
// Функція
func Area(r Rectangle) int { ... }

// Метод
func (r Rectangle) Area() int { ... }
```

#### 2. Pointer Receiver
**Q: Коли потрібен pointer receiver?**
- Коли треба змінити дані
- Коли struct великий (економія пам'яті)
- Коли інші методи використовують pointer

**Q: Приклад?**
```go
// Value receiver - не змінює
func (r Rectangle) Double() {
    r.Width *= 2  // змінює копію!
}

// Pointer receiver - змінює
func (r *Rectangle) Double() {
    r.Width *= 2  // змінює оригінал
}
```

#### 3. Implicit Interfaces
**Q: Чому неявна реалізація?**
- Гнучкість (не залежить від конкретного типу)
- Легше тестувати
- Можна додати інтерфейс до існуючого коду

**Q: Як працює?**
```go
type Writer interface {
    Write([]byte) error
}

type FileWriter struct{}

// Реалізує Write - автоматично Writer!
func (f FileWriter) Write(data []byte) error {
    return nil
}
```

---

## ⚡ Швидкі команди

### Запустити приклади
```bash
# Methods demo
cd practice/methods_demo && go run main.go

# Interface demo
cd ../interface_demo && go run main.go

# UserService
cd ../user_service && go run main.go
```

### Створити власний приклад
```bash
cat > my_example.go << 'EOF'
package main

import "fmt"

// Struct
type Rectangle struct {
    Width  int
    Height int
}

// Method
func (r Rectangle) Area() int {
    return r.Width * r.Height
}

// Pointer method
func (r *Rectangle) Scale(factor int) {
    r.Width *= factor
    r.Height *= factor
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}
    fmt.Printf("Area: %d\n", rect.Area())
    
    rect.Scale(2)
    fmt.Printf("After scale: %+v\n", rect)
}
EOF

go run my_example.go
rm my_example.go
```

---

## 🎯 Контрольний список

Після тижня 2 ви повинні:

### Теорія
- [ ] Розумію різницю між методами та функціями
- [ ] Знаю що таке receiver (value і pointer)
- [ ] Розумію коли використовувати pointer receiver
- [ ] Розумію неявну реалізацію інтерфейсів
- [ ] Можу пояснити переваги інтерфейсів

### Практика
- [ ] Створив методи на власному struct
- [ ] Використав value та pointer receivers
- [ ] Створив інтерфейс
- [ ] Написав 2+ реалізації інтерфейсу
- [ ] Передав різні реалізації через інтерфейс

### Код
- [ ] Написав Calculator з методами
- [ ] Написав Shape інтерфейс з фігурами
- [ ] Написав Storage інтерфейс (memory + file)
- [ ] Розібрав UserService приклад
- [ ] Можу пояснити свій код

---

## 💡 Підказки

### Коли використовувати pointer receiver?

✅ **Використовуйте pointer receiver коли:**
1. Метод змінює дані struct
2. Struct великий (економія пам'яті)
3. Для консистентності (якщо один метод pointer, всі pointer)

❌ **НЕ використовуйте pointer receiver коли:**
1. Struct маленький (наприклад, кілька int)
2. Struct незмінний (immutable)
3. Receiver - це map, slice, chan (вони вже reference types)

### Правила іменування інтерфейсів

```go
// ✅ Добре: -er суфікс
type Reader interface { ... }
type Writer interface { ... }
type Logger interface { ... }

// ✅ Добре: описова назва
type UserService interface { ... }
type DataStore interface { ... }

// ❌ Погано: занадто загально
type Data interface { ... }
type Manager interface { ... }
```

### Маленькі інтерфейси > Великі

```go
// ✅ Добре: маленькі, focused
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// ❌ Погано: один великий інтерфейс
type Storage interface {
    Read()
    Write()
    Delete()
    Update()
    List()
    // ... ще 10 методів
}
```

**Правило:** "The bigger the interface, the weaker the abstraction"

---

## 📚 Додаткові ресурси

### Офіційна документація
- [Go Tour - Methods](https://go.dev/tour/methods/)
- [Effective Go - Interfaces](https://go.dev/doc/effective_go#interfaces)
- [Go Spec - Interface Types](https://go.dev/ref/spec#Interface_types)

### Статті
- [How to use interfaces in Go](https://jordanorelli.com/post/32665860244/how-to-use-interfaces-in-go)
- [Accept interfaces, return structs](https://bryanftan.medium.com/accept-interfaces-return-structs-in-go-d4cab29a301b)

### Відео
- [justforfunc #19 - Understanding Interfaces](https://www.youtube.com/watch?v=F4wUrj6pmSI)
- [GopherCon 2015 - The Design of the Go Assembler](https://www.youtube.com/watch?v=KINIAgRpkDA)

---

## 🎓 Практичні патерни

### 1. Dependency Injection через інтерфейси

```go
type UserService struct {
    storage Storage  // інтерфейс, не конкретний тип!
}

// Можна передати будь-яку реалізацію
service := UserService{storage: &MemoryStorage{}}
service := UserService{storage: &FileStorage{}}
service := UserService{storage: &MockStorage{}}
```

### 2. Mock для тестів

```go
// Реальна реалізація
type RealEmailSender struct{}
func (r RealEmailSender) Send(to, msg string) error {
    // справжня відправка email
}

// Mock для тестів
type MockEmailSender struct{}
func (m MockEmailSender) Send(to, msg string) error {
    // просто логує
    fmt.Printf("Mock: sending to %s\n", to)
    return nil
}

// Обидва реалізують інтерфейс!
type EmailSender interface {
    Send(to, msg string) error
}
```

### 3. Композиція інтерфейсів

```go
// Маленькі інтерфейси
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Closer interface {
    Close() error
}

// Композиція
type ReadCloser interface {
    Reader
    Closer
}
```

---

## 🚧 Поширені помилки

### 1. Забувають pointer receiver

```go
// ❌ Погано - не змінить оригінал
func (r Rectangle) Scale(factor int) {
    r.Width *= factor
}

// ✅ Добре
func (r *Rectangle) Scale(factor int) {
    r.Width *= factor
}
```

### 2. Порожні інтерфейси всюди

```go
// ❌ Погано - втрачаємо type safety
func Process(data interface{}) { ... }

// ✅ Добре - конкретний тип або специфічний інтерфейс
func Process(data User) { ... }
func Process(data Processor) { ... }
```

### 3. Надто великі інтерфейси

```go
// ❌ Погано - важко реалізувати
type SuperService interface {
    Create()
    Read()
    Update()
    Delete()
    List()
    Search()
    Export()
    Import()
}

// ✅ Добре - маленькі інтерфейси
type Creator interface { Create() }
type Reader interface { Read() }
type Updater interface { Update() }
```

---

## ❓ Питання та відповіді

### Q: Чи можу я додати метод до чужого типу?
A: Ні, напряму не можна. Але можна створити свій тип на основі чужого:
```go
type MyInt int

func (m MyInt) Double() MyInt {
    return m * 2
}
```

### Q: Скільки інтерфейсів може реалізувати один тип?
A: Скільки завгодно! Якщо тип має всі необхідні методи.

### Q: Чи можна зберігати nil в інтерфейсі?
A: Так, але це може призвести до runtime panic. Будьте обережні!

### Q: Pointer receiver чи value receiver для маленького struct?
A: Для маленького struct (кілька простих полів) зазвичай value receiver. Але якщо треба змінювати - pointer.

---

## 🎉 Успіхів у навчанні!

**Пам'ятайте:**
- Methods роблять код виразнішим
- Pointer receivers для зміни даних
- Інтерфейси роблять код гнучким
- Маленькі інтерфейси > Великі
- Accept interfaces, return structs

---

**Happy coding! 🚀**

