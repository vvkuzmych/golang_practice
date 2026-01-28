# Week 6: Backend Fundamentals & Architecture

## 🎯 Мета тижня

Опанувати фундаментальні концепції backend розробки: ООП, патерни проєктування, HTTP сервери, мікросервіси, бази даних та нетворкінг.

---

## 📚 Теми тижня

### 1. ООП принципи в Go
- ✅ **Інкапсуляція** - приховування даних
- ✅ **Поліморфізм** - інтерфейси та duck typing
- ✅ **Абстракція** - спрощення складності
- ✅ **Композиція** замість успадкування

### 2. Патерни проєктування
- ✅ **Creational** - Singleton, Factory, Builder
- ✅ **Structural** - Adapter, Decorator, Facade
- ✅ **Behavioral** - Strategy, Observer, Command

### 3. Golang net/http
- ✅ HTTP Server & Client
- ✅ Routing & Middleware
- ✅ Request/Response handling
- ✅ Context & Timeouts

### 4. Мікросервісна архітектура
- ✅ Монолітна vs Мікросервісна
- ✅ API Gateway
- ✅ Service Discovery
- ✅ Inter-service communication

### 5. Бази даних
- ✅ PostgreSQL & SQL
- ✅ GORM (ORM)
- ✅ Міграції
- ✅ Transactions

### 6. Нетворкінг
- ✅ TCP/UDP
- ✅ HTTP Semantics
- ✅ TLS/SSL
- ✅ DNS
- ✅ Timeouts & Retries

### 7. Goroutines і Конкурентність
- ✅ Goroutines
- ✅ Channels (buffered/unbuffered)
- ✅ Select
- ✅ Sync Package (Mutex, WaitGroup, Once)
- ✅ Конкурентні Патерни
- ✅ Pipeline Pattern (Fan-Out/Fan-In)

---

## 📂 Структура

```
week_6/
├── README.md                          # Ви тут
├── QUICK_START.md                     # Швидкий старт
│
├── theory/                            # 📖 Теорія
│   ├── 01_oop_principles.md          # ООП в Go
│   ├── 02_design_patterns.md         # Патерни проєктування
│   ├── 03_net_http.md                # net/http
│   ├── 04_microservices.md           # Мікросервіси
│   ├── 05_databases.md               # Бази даних
│   ├── 06_networking.md              # Нетворкінг
│   ├── 07_goroutines_concurrency.md  # Goroutines
│   └── 08_pipeline_pattern.md        # Pipeline Pattern
│
├── practice/                          # 💻 Практика
│   ├── 01_oop/                       # ООП приклади
│   ├── 02_http_server/               # HTTP сервер
│   ├── 03_microservices/             # Мікросервіси
│   ├── 04_database/                  # БД
│   ├── 05_networking/                # Нетворкінг
│   └── 06_goroutines/                # Goroutines
│
├── exercises/                         # ✏️ Завдання
│   ├── exercise_1.md                 # ООП + Патерни
│   ├── exercise_2.md                 # HTTP Server
│   └── exercise_3.md                 # Full Stack задача
│
└── solutions/                         # ✅ Рішення
    └── solutions.md
```

---

## 🚀 Швидкий старт

### 1. Вивчити теорію
```bash
# Почніть з теорії
cat theory/01_oop_principles.md
cat theory/02_design_patterns.md
cat theory/03_net_http.md
```

### 2. Запустити практичні приклади
```bash
# ООП приклади
go run practice/01_oop/main.go

# HTTP сервер
go run practice/02_http_server/main.go

# Мікросервіси
go run practice/03_microservices/service_a/main.go
go run practice/03_microservices/service_b/main.go
```

### 3. Виконати вправи
```bash
# Прочитайте завдання
cat exercises/exercise_1.md

# Створіть своє рішення
mkdir my_solutions
cd my_solutions
```

---

## 📖 Рекомендований порядок вивчення

### День 1-2: ООП і Патерни
1. Прочитайте `theory/01_oop_principles.md`
2. Прочитайте `theory/02_design_patterns.md`
3. Запустіть приклади з `practice/01_oop/`
4. Виконайте `exercises/exercise_1.md`

### День 3-4: HTTP і Мікросервіси
1. Прочитайте `theory/03_net_http.md`
2. Прочитайте `theory/04_microservices.md`
3. Запустіть `practice/02_http_server/`
4. Запустіть `practice/03_microservices/`
5. Виконайте `exercises/exercise_2.md`

### День 5-6: Бази даних і Нетворкінг
1. Прочитайте `theory/05_databases.md`
2. Прочитайте `theory/06_networking.md`
3. Запустіть `practice/04_database/`
4. Запустіть `practice/05_networking/`
5. Виконайте `exercises/exercise_3.md`

### День 7: Проєкт
Створіть власний мікросервіс з:
- HTTP API
- PostgreSQL
- Proper error handling
- Timeouts & retries

---

## 🎓 Що ви вивчите

### ООП в Go
```go
// Інкапсуляція
type User struct {
    name  string // приватне
    email string
}

// Поліморфізм через інтерфейси
type PaymentProcessor interface {
    Process(amount float64) error
}
```

### HTTP Server
```go
http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(users)
})
http.ListenAndServe(":8080", nil)
```

### Мікросервіси
```
API Gateway → User Service → Auth Service
           ↘ Product Service → Inventory Service
```

### PostgreSQL
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### TCP/UDP
```go
// TCP Server
ln, _ := net.Listen("tcp", ":8080")
conn, _ := ln.Accept()
```

---

## 🔗 Корисні ресурси

### Документація
- [Go net/http](https://pkg.go.dev/net/http)
- [GORM](https://gorm.io/docs/)
- [PostgreSQL](https://www.postgresql.org/docs/)

### Книги
- "Design Patterns in Go" by Mario Zupan
- "Building Microservices" by Sam Newman
- "HTTP: The Definitive Guide" by David Gourley

### Статті
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Microservices.io](https://microservices.io/)

---

## ✅ Чеклист прогресу

- [ ] Розумію 4 принципи ООП
- [ ] Знаю основні патерни проєктування
- [ ] Вмію створювати HTTP сервери
- [ ] Розумію мікросервісну архітектуру
- [ ] Працював з PostgreSQL через GORM
- [ ] Розумію різницю між TCP і UDP
- [ ] Знаю що таке TLS/SSL
- [ ] Вмію правильно налаштовувати timeouts
- [ ] Виконав всі вправи
- [ ] Створив власний проєкт

---

## 💡 Поради

1. **Не поспішайте** - краще розібратись глибоко, ніж швидко пробігти
2. **Практикуйте** - кожна теорія має практичний приклад
3. **Експериментуйте** - змінюйте код, дивіться що станеться
4. **Питайте** - якщо щось незрозуміло, шукайте відповіді
5. **Створюйте** - найкращий спосіб вчитись це будувати

---

## 🎯 Наступні кроки

Після завершення Week 6 ви будете готові до:
- Створення production-ready API
- Розробки мікросервісів
- Роботи з базами даних
- Налаштування мережевої інфраструктури
- Проходження технічних інтерв'ю

**Успіхів у навчанні!** 🚀

---

**Автор:** Golang Practice Course  
**Версія:** 1.0  
**Дата:** 2026-01-27
