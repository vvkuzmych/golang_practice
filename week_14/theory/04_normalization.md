# Database Normalization (Нормалізація БД)

## 🎯 Що таке Normalization?

**Normalization (Нормалізація)** - процес організації даних для:
- ❌ Усунення дублювання
- ✅ Забезпечення цілісності
- ✅ Полегшення оновлень

---

## 📊 Normal Forms (Нормальні форми)

```
Unnormalized → 1NF → 2NF → 3NF → BCNF → 4NF → 5NF
(погано)                   (добре)         (ідеально)
```

**На практиці:** Зазвичай достатньо **3NF**

---

## 0️⃣ Unnormalized (Ненормалізована)

**Проблема:** Дублювання даних, складно оновлювати

```sql
orders
────────────────────────────────────────────────────────────────
id | customer_name | customer_email | products           | total
1  | John Doe      | john@mail.com  | Laptop, Mouse     | 1225
2  | John Doe      | john@mail.com  | Keyboard          | 75
3  | Jane Smith    | jane@mail.com  | Monitor, Keyboard | 375
```

**Проблеми:**
- ❌ John's email дублюється (Update Anomaly)
- ❌ Products в одному полі (Multiple values)
- ❌ Не можна додати customer без order (Insert Anomaly)
- ❌ Видалення order = втрата customer info (Delete Anomaly)

---

## 1️⃣ First Normal Form (1NF)

**Правило:** 
- ✅ Atomic values (кожна комірка = одне значення)
- ✅ No repeating groups
- ✅ Primary key

### ❌ Порушення 1NF

```sql
students
─────────────────────────────────
id | name  | courses
1  | John  | Math, Physics, CS     ← Multiple values!
2  | Jane  | Math, Chemistry
```

### ✅ 1NF

```sql
students
─────────────
id | name
1  | John
2  | Jane

student_courses
────────────────────
student_id | course
1          | Math
1          | Physics        ← Atomic values
1          | CS
2          | Math
2          | Chemistry
```

**Досягнуто:**
- ✅ Кожна комірка = 1 значення
- ✅ Primary key (student_id + course)

**Але залишилось:**
- ❌ Дублювання customer info

---

## 2️⃣ Second Normal Form (2NF)

**Правило:**
- ✅ Must be in 1NF
- ✅ No partial dependencies (всі non-key поля залежать від ВСЬОГО ключа)

### ❌ Порушення 2NF

```sql
order_items
──────────────────────────────────────────────────────────
order_id | product_id | product_name | product_price | quantity
1        | 101        | Laptop       | 1200          | 1
1        | 102        | Mouse        | 25            | 4
2        | 103        | Keyboard     | 75            | 1
         ↑            ↑
    Composite PK   Залежить тільки від product_id!
```

**Проблема:** `product_name`, `product_price` залежать тільки від `product_id`, а не від всього ключа `(order_id, product_id)`

### ✅ 2NF

```sql
-- Products окремо
products
─────────────────────────────
product_id | name     | price
101        | Laptop   | 1200
102        | Mouse    | 25
103        | Keyboard | 75

-- Order items без product info
order_items
──────────────────────────────
order_id | product_id | quantity
1        | 101        | 1
1        | 102        | 4
2        | 103        | 1
```

**Досягнуто:**
- ✅ No partial dependencies
- ✅ Product info тільки в products

**Але залишилось:**
- ❌ Transitive dependencies

---

## 3️⃣ Third Normal Form (3NF)

**Правило:**
- ✅ Must be in 2NF
- ✅ No transitive dependencies (non-key поля не залежать один від одного)

### ❌ Порушення 3NF

```sql
employees
───────────────────────────────────────────────────────
id | name  | dept_id | dept_name     | dept_location
1  | John  | 10      | Engineering   | Building A
2  | Jane  | 10      | Engineering   | Building A    ← Дублювання!
3  | Bob   | 20      | Sales         | Building B
              ↓           ↓
         dept_name залежить від dept_id, не від id!
```

**Проблема:** `dept_name`, `dept_location` залежать від `dept_id`, а не від `id` (transitive dependency)

### ✅ 3NF

```sql
-- Departments окремо
departments
─────────────────────────────────
dept_id | name        | location
10      | Engineering | Building A
20      | Sales       | Building B

-- Employees без dept info
employees
──────────────────────
id | name  | dept_id
1  | John  | 10
2  | Jane  | 10
3  | Bob   | 20
```

**Досягнуто:**
- ✅ No transitive dependencies
- ✅ Department info тільки в departments
- ✅ Легко оновити dept_location для всіх employees

---

## 📊 Comparison: Unnormalized vs 3NF

### Unnormalized

```sql
orders
────────────────────────────────────────────────────────────────────
id | customer | email         | product  | price | quantity | total
1  | John     | john@mail.com | Laptop   | 1200  | 1        | 1200
2  | John     | john@mail.com | Mouse    | 25    | 4        | 100
3  | John     | john@mail.com | Keyboard | 75    | 1        | 75
4  | Jane     | jane@mail.com | Monitor  | 300   | 1        | 300
```

**Проблеми:**
- John's email дублюється 3 рази
- Product info дублюється
- Update John's email = 3 updates

### 3NF

```sql
customers
─────────────────────────
id | name  | email
1  | John  | john@mail.com
2  | Jane  | jane@mail.com

products
─────────────────────────────
id  | name     | price
101 | Laptop   | 1200
102 | Mouse    | 25
103 | Keyboard | 75
104 | Monitor  | 300

orders
──────────────────────────
id | customer_id | total
1  | 1           | 1200
2  | 1           | 100
3  | 1           | 75
4  | 2           | 300

order_items
──────────────────────────────
order_id | product_id | quantity
1        | 101        | 1
2        | 102        | 4
3        | 103        | 1
4        | 104        | 1
```

**Переваги:**
- ✅ John's email в одному місці
- ✅ Update email = 1 update
- ✅ No duplicates

---

## 🎯 BCNF (Boyce-Codd Normal Form)

**Правило:**
- ✅ Must be in 3NF
- ✅ Every determinant must be a candidate key

**Коли потрібно:** Рідко, для складних випадків

### Приклад

```sql
-- Порушення BCNF
student_advisor
─────────────────────────────────
student_id | subject | advisor
1          | Math    | Prof. Smith
1          | Physics | Prof. Jones
2          | Math    | Prof. Smith

-- Якщо advisor визначається subject (кожен subject має 1 advisor)
-- То advisor → subject, але advisor не є candidate key

-- BCNF:
subjects
──────────────────────
subject | advisor
Math    | Prof. Smith
Physics | Prof. Jones

student_subjects
─────────────────────
student_id | subject
1          | Math
1          | Physics
2          | Math
```

---

## ⚖️ Denormalization (Денормалізація)

**Коли доцільно повернути дублювання:**
- 🚀 Performance (часті JOINs повільні)
- 📊 Reporting/Analytics
- 📱 Read-heavy applications

### Приклад: Денормалізація для швидкості

```sql
-- 3NF (повільний SELECT)
SELECT 
    o.id,
    c.name,
    c.email,
    SUM(oi.quantity * p.price) as total
FROM orders o
JOIN customers c ON o.customer_id = c.id
JOIN order_items oi ON o.id = oi.order_id
JOIN products p ON oi.product_id = p.id
GROUP BY o.id, c.name, c.email;

-- Денормалізована (швидкий SELECT)
orders_denormalized
───────────────────────────────────────────
id | customer_name | customer_email | total
1  | John Doe      | john@mail.com  | 1200
2  | John Doe      | john@mail.com  | 100

-- ✅ Швидше (no JOINs)
-- ❌ Дублювання customer info
-- ⚠️ Потрібні тригери для sync
```

**Trade-offs:**
```
Normalization:
✅ No duplicates
✅ Easy updates
❌ Slow reads (many JOINs)

Denormalization:
✅ Fast reads (no JOINs)
❌ Duplicates
❌ Complex updates
```

---

## 🎯 Practical Guidelines

### Завжди 3NF для:
- 📝 Transactional data (orders, users)
- 🔄 Frequent updates
- 🔐 Critical data integrity

### Денормалізація для:
- 📊 Analytics/reporting tables
- 📱 Read-heavy APIs
- 🚀 Performance-critical queries
- 📈 Dashboards

---

## 📊 Real-World Example: E-commerce

### 3NF Design

```sql
-- Core entities (3NF)
users (id, email, name)
products (id, name, price, stock)
orders (id, user_id, created_at, status)
order_items (order_id, product_id, quantity, price_at_order)

-- Queries потребують JOINs:
SELECT u.name, o.id, SUM(oi.quantity * oi.price)
FROM users u
JOIN orders o ON u.id = o.user_id
JOIN order_items oi ON o.id = oi.order_id
GROUP BY u.id, o.id;
```

### Hybrid (3NF + Denormalized views)

```sql
-- Core tables (3NF) - для writes
users, products, orders, order_items

-- Denormalized view - для reads
CREATE MATERIALIZED VIEW order_summary AS
SELECT 
    o.id as order_id,
    u.name as customer_name,
    u.email as customer_email,
    o.created_at,
    SUM(oi.quantity * oi.price) as total
FROM orders o
JOIN users u ON o.user_id = u.id
JOIN order_items oi ON o.id = oi.order_id
GROUP BY o.id, u.name, u.email, o.created_at;

-- ✅ Best of both worlds:
-- Writes → 3NF tables
-- Reads → Denormalized view
```

---

## ✅ Decision Tree

```
Нова таблиця?
├─ Є повторювані поля?
│  └─ ТАК → Split into 1NF
│
├─ Є partial dependencies?
│  └─ ТАК → Extract to 2NF
│
├─ Є transitive dependencies?
│  └─ ТАК → Extract to 3NF
│
└─ Queries повільні?
   ├─ ТАК → Consider indexes first
   └─ Still slow? → Denormalize specific queries
```

---

## 🎓 Cheat Sheet

### 1NF
```
❌ courses: "Math, Physics, CS"
✅ 3 рядки: Math, Physics, CS
```

### 2NF
```
❌ product_name залежить від product_id (частина composite key)
✅ products таблиця окремо
```

### 3NF
```
❌ dept_name залежить від dept_id (не від PK)
✅ departments таблиця окремо
```

### Denormalization
```
❌ Завжди нормалізувати
✅ Нормалізуй core, денормалізуй views/cache
```

---

## 📊 Anomalies (Аномалії)

### Insert Anomaly
```
-- Unnormalized
❌ Не можу додати customer без order

-- 3NF
✅ Можу додати customer окремо
```

### Update Anomaly
```
-- Unnormalized
❌ Update email = багато updates + можливі помилки

-- 3NF
✅ Update email = 1 update
```

### Delete Anomaly
```
-- Unnormalized
❌ Delete останній order = втрата customer info

-- 3NF
✅ Delete order ≠ delete customer
```

---

## ✅ Висновок

### Normal Forms:

✅ **1NF** - atomic values, no repeating groups  
✅ **2NF** - no partial dependencies  
✅ **3NF** - no transitive dependencies  
✅ **BCNF** - advanced (рідко потрібно)  

### In Practice:

**3NF = Sweet Spot** для більшості cases

**Denormalization:**
- Тільки після вимірювання performance
- Зазвичай для read-heavy tables
- Materialized views > duplicate data

### Golden Rule:

**"Normalize until it hurts, denormalize until it works"**

**Week 14: Normalization Master!** 🗄️📊
