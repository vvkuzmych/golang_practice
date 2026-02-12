# SOLID Principles in Go

## 1. Single Responsibility Principle (SRP)

**Клас повинен мати тільки одну причину для зміни.**

### ❌ Порушення SRP

```go
type User struct {
    Name  string
    Email string
}

// Занадто багато відповідальностей
func (u *User) Save() error {
    // 1. Валідація
    if u.Email == "" {
        return errors.New("email is required")
    }
    
    // 2. З'єднання з БД
    db, _ := sql.Open("postgres", "...")
    
    // 3. SQL запит
    _, err := db.Exec("INSERT INTO users...")
    
    // 4. Логування
    log.Printf("User %s saved", u.Name)
    
    return err
}
```

### ✅ Дотримання SRP

```go
// Domain entity
type User struct {
    ID    int
    Name  string
    Email string
}

// Validator (одна відповідальність)
type UserValidator struct{}

func (v *UserValidator) Validate(u *User) error {
    if u.Email == "" {
        return errors.New("email is required")
    }
    return nil
}

// Repository (одна відповідальність)
type UserRepository struct {
    db *sql.DB
}

func (r *UserRepository) Save(u *User) error {
    _, err := r.db.Exec("INSERT INTO users...", u.Name, u.Email)
    return err
}

// Logger (одна відповідальність)
type UserLogger struct {
    logger *log.Logger
}

func (l *UserLogger) LogSaved(u *User) {
    l.logger.Printf("User %s saved", u.Name)
}
```

---

## 2. Open/Closed Principle (OCP)

**Відкритий для розширення, закритий для модифікації.**

### ❌ Порушення OCP

```go
type PaymentProcessor struct{}

func (p *PaymentProcessor) ProcessPayment(method string, amount float64) error {
    if method == "credit_card" {
        // Credit card logic
    } else if method == "paypal" {
        // PayPal logic
    } else if method == "bitcoin" {
        // Bitcoin logic
    }
    // Додавання нового методу вимагає зміни цього коду!
    return nil
}
```

### ✅ Дотримання OCP

```go
// Інтерфейс для розширення
type PaymentMethod interface {
    Process(amount float64) error
}

// Implementations
type CreditCardPayment struct{}

func (c *CreditCardPayment) Process(amount float64) error {
    // Credit card logic
    return nil
}

type PayPalPayment struct{}

func (p *PayPalPayment) Process(amount float64) error {
    // PayPal logic
    return nil
}

// Processor закритий для модифікації
type PaymentProcessor struct{}

func (p *PaymentProcessor) ProcessPayment(method PaymentMethod, amount float64) error {
    return method.Process(amount)
}

// Додавання нового методу БЕЗ зміни існуючого коду
type BitcoinPayment struct{}

func (b *BitcoinPayment) Process(amount float64) error {
    // Bitcoin logic
    return nil
}
```

---

## 3. Liskov Substitution Principle (LSP)

**Підкласи повинні замінювати базові класи без порушення функціональності.**

### ❌ Порушення LSP

```go
type Bird interface {
    Fly() error
}

type Sparrow struct{}

func (s *Sparrow) Fly() error {
    fmt.Println("Sparrow flying")
    return nil
}

type Penguin struct{}

func (p *Penguin) Fly() error {
    return errors.New("penguins can't fly") // Порушення контракту!
}
```

### ✅ Дотримання LSP

```go
type Bird interface {
    Move() error
}

type FlyingBird interface {
    Bird
    Fly() error
}

type Sparrow struct{}

func (s *Sparrow) Move() error {
    return s.Fly()
}

func (s *Sparrow) Fly() error {
    fmt.Println("Sparrow flying")
    return nil
}

type Penguin struct{}

func (p *Penguin) Move() error {
    fmt.Println("Penguin swimming")
    return nil
}
```

---

## 4. Interface Segregation Principle (ISP)

**Клієнти не повинні залежати від інтерфейсів, які вони не використовують.**

### ❌ Порушення ISP

```go
// Занадто великий інтерфейс
type Worker interface {
    Work() error
    Eat() error
    Sleep() error
}

type Robot struct{}

func (r *Robot) Work() error {
    return nil
}

func (r *Robot) Eat() error {
    return errors.New("robots don't eat") // Непотрібний метод!
}

func (r *Robot) Sleep() error {
    return errors.New("robots don't sleep") // Непотрібний метод!
}
```

### ✅ Дотримання ISP

```go
// Малі, специфічні інтерфейси
type Worker interface {
    Work() error
}

type Eater interface {
    Eat() error
}

type Sleeper interface {
    Sleep() error
}

// Robot імплементує тільки те, що потрібно
type Robot struct{}

func (r *Robot) Work() error {
    return nil
}

// Human імплементує всі три
type Human struct{}

func (h *Human) Work() error  { return nil }
func (h *Human) Eat() error   { return nil }
func (h *Human) Sleep() error { return nil }
```

---

## 5. Dependency Inversion Principle (DIP)

**Залежати від абстракцій, а не від конкретних реалізацій.**

### ❌ Порушення DIP

```go
// Висока залежність від конкретної реалізації
type MySQLDatabase struct{}

func (db *MySQLDatabase) Save(data string) error {
    fmt.Println("Saving to MySQL")
    return nil
}

type UserService struct {
    db *MySQLDatabase // Пряма залежність від MySQL
}

func (s *UserService) CreateUser(name string) error {
    return s.db.Save(name)
}
```

### ✅ Дотримання DIP

```go
// Абстракція
type Database interface {
    Save(data string) error
}

// Реалізація 1
type MySQLDatabase struct{}

func (db *MySQLDatabase) Save(data string) error {
    fmt.Println("Saving to MySQL")
    return nil
}

// Реалізація 2
type PostgresDatabase struct{}

func (db *PostgresDatabase) Save(data string) error {
    fmt.Println("Saving to Postgres")
    return nil
}

// Сервіс залежить від абстракції
type UserService struct {
    db Database // Залежність від інтерфейсу
}

func (s *UserService) CreateUser(name string) error {
    return s.db.Save(name)
}

// Легко змінити БД без зміни UserService
func main() {
    service1 := &UserService{db: &MySQLDatabase{}}
    service2 := &UserService{db: &PostgresDatabase{}}
    
    service1.CreateUser("Alice")
    service2.CreateUser("Bob")
}
```

---

## Практичний приклад: Всі SOLID разом

```go
// 1. SRP: Окремі структури для різних відповідальностей
type User struct {
    ID    int
    Name  string
    Email string
}

// 2. DIP: Залежність від інтерфейсів
type UserRepository interface {
    Save(u *User) error
    Find(id int) (*User, error)
}

type UserNotifier interface {
    Notify(u *User, message string) error
}

// 3. ISP: Малі, специфічні інтерфейси
type EmailNotifier struct{}

func (n *EmailNotifier) Notify(u *User, message string) error {
    fmt.Printf("Sending email to %s: %s\n", u.Email, message)
    return nil
}

// 4. OCP: Відкритий для розширення через інтерфейси
type UserService struct {
    repo     UserRepository
    notifier UserNotifier
}

func (s *UserService) RegisterUser(name, email string) error {
    user := &User{Name: name, Email: email}
    
    if err := s.repo.Save(user); err != nil {
        return err
    }
    
    return s.notifier.Notify(user, "Welcome!")
}

// 5. LSP: Різні реалізації можна замінити
type PostgresUserRepository struct {
    db *sql.DB
}

func (r *PostgresUserRepository) Save(u *User) error {
    _, err := r.db.Exec("INSERT INTO users...", u.Name, u.Email)
    return err
}

func (r *PostgresUserRepository) Find(id int) (*User, error) {
    var u User
    err := r.db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Email)
    return &u, err
}
```

---

## Підсумок

| Принцип | Ключова ідея |
|---------|--------------|
| **SRP** | Один клас = одна відповідальність |
| **OCP** | Відкритий для розширення, закритий для модифікації |
| **LSP** | Підкласи замінюють базові класи без порушень |
| **ISP** | Малі, специфічні інтерфейси |
| **DIP** | Залежати від абстракцій, а не конкретних реалізацій |

**Переваги SOLID:**
- ✅ Легше тестувати (моки, стаби)
- ✅ Легше підтримувати (зміни локалізовані)
- ✅ Легше розширювати (нові features без ламання старих)
- ✅ Менше coupling, більше cohesion
