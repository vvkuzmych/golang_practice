# ✅ Week 14 - Завершено!

## 🎯 Що створено

**Week 14: Data Models, SQL JOINs, OSI & Normalization** - модуль про моделі даних (1:1, 1:N, N:M), всі типи SQL JOINs, 7 рівнів OSI моделі та Database Normalization (1NF, 2NF, 3NF).

---

## 📊 Статистика

### Створено файлів

**Теорія:** 4 файли (~2400 рядків)
- `theory/01_data_models.md` - Типи відношень, keys
- `theory/02_sql_joins.md` - Всі JOIN types з візуалізаціями
- `theory/03_osi_model.md` - 7 рівнів OSI моделі, протоколи, debugging
- `theory/04_normalization.md` - Normal Forms (1NF, 2NF, 3NF, BCNF), anomalies

**Практика:** 3 файли (~600 рядків SQL + Go)
- `practice/01_basic_joins/schema.sql` - Schema + test data
- `practice/01_basic_joins/queries.sql` - 50+ приклад queries
- `practice/03_go_joins/main.go` - Go код з database/sql

**Документація:** 5 файлів
- `README.md` - Повний guide
- `QUICK_START.md` - Швидкий старт
- `WEEK14_COMPLETE.md` - Цей звіт
- `JOINS_CHEAT_SHEET.md` - JOINs довідка
- `OSI_CHEAT_SHEET.md` - OSI довідка
- `NORMALIZATION_CHEAT_SHEET.md` - Normalization довідка

**Загалом:** 12 файлів, ~3800 рядків

---

## 📚 Що покрито

### 1. Data Models 🗄️

**Типи відношень:**

**One-to-One (1:1)**
```
User ←→ Profile
```
- UNIQUE foreign key
- Приклад: user.id ←→ profiles.user_id (UNIQUE)

**One-to-Many (1:N)**
```
User ───< Posts
```
- Regular foreign key
- Приклад: user.id ←→ posts.user_id

**Many-to-Many (N:M)**
```
Students >──< Courses
```
- Junction table з двома FK
- Приклад: enrollments(student_id, course_id)

**Keys:**
- **Primary Key (PK)** - унікальний ID
- **Foreign Key (FK)** - посилання на PK
- **Composite Key** - комбінація полів

---

### 2. SQL JOINs 🔗

**INNER JOIN** - тільки співпадіння
```sql
SELECT u.name, o.total
FROM users u
INNER JOIN orders o ON u.id = o.user_id;
```
```
users ∩ orders
```

**LEFT JOIN** - всі з лівої + співпадіння
```sql
SELECT u.name, o.total
FROM users u
LEFT JOIN orders o ON u.id = o.user_id;
```
```
users ∪ orders (all users)
```

**RIGHT JOIN** - всі з правої + співпадіння
```sql
SELECT u.name, o.total
FROM users u
RIGHT JOIN orders o ON u.id = o.user_id;
```
```
users ∪ orders (all orders)
```

**FULL OUTER JOIN** - всі з обох
```sql
SELECT u.name, o.total
FROM users u
FULL OUTER JOIN orders o ON u.id = o.user_id;
```
```
users ∪ orders (all)
```

**CROSS JOIN** - всі комбінації
```sql
SELECT u.name, p.name
FROM users u
CROSS JOIN products p;
```
```
users × products (Cartesian)
```

**SELF JOIN** - таблиця з собою
```sql
SELECT e.name, m.name AS manager
FROM employees e
LEFT JOIN employees m ON e.manager_id = m.id;
```

---

### 3. OSI Model 🌐

**7 рівнів мережевої взаємодії:**

```
7. Application   → HTTP, DNS, SMTP, FTP, SSH
6. Presentation  → SSL/TLS, encryption, JPEG, GZIP
5. Session       → Sessions, authentication
4. Transport     → TCP (reliable), UDP (fast), ports
3. Network       → IP addresses, routing, packets
2. Data Link     → MAC addresses, switches, frames
1. Physical      → Cables, hubs, bits, electrical signals
```

**Мнемоніка:** **P**lease **D**o **N**ot **T**hrow **S**ausage **P**izza **A**way

**TCP vs UDP:**
- **TCP:** Надійний, handshake, HTTP/HTTPS/SSH
- **UDP:** Швидкий, без підтвердження, DNS/games/streaming

**Common Ports:**
- 22: SSH, 80: HTTP, 443: HTTPS
- 3306: MySQL, 5432: PostgreSQL, 6379: Redis

---

### 4. Database Normalization 📊

**Normal Forms:**

```
Unnormalized → 1NF → 2NF → 3NF
(дублювання)          (оптимально)
```

**1NF (First Normal Form):**
- ✅ Atomic values (one value per cell)
- ❌ No "Math, Physics, CS" in one field
- ✅ Primary key exists

**2NF (Second Normal Form):**
- ✅ 1NF +
- ❌ No partial dependencies
- ✅ Fields depend on WHOLE key

**3NF (Third Normal Form):**
- ✅ 2NF +
- ❌ No transitive dependencies
- ✅ Non-key fields depend ONLY on PK

**Example: Unnormalized → 3NF**

```sql
-- Unnormalized (bad)
orders
──────────────────────────────────────────
id | customer | email         | product
1  | John     | john@mail.com | Laptop
2  | John     | john@mail.com | Mouse    ← Duplicate!

-- 3NF (good)
customers         orders
──────────────   ─────────────────
id | name|email  id | cust_id | prod_id
1  | John|j@m.c  1  | 1       | 1
                 2  | 1       | 2
```

**Anomalies:**
- **Insert:** Can't add customer without order (fixed in 3NF)
- **Update:** Email duplicated (fixed in 3NF)
- **Delete:** Delete order = lose customer (fixed in 3NF)

**Denormalization:**
- Для performance (after measuring!)
- Read-heavy tables, dashboards
- Materialized views > duplicate data

**Golden Rule:** "The key, the whole key, and nothing but the key"

---

## 🎯 Практичні приклади

### Schema

```sql
users (PK: id)
  ↓ 1:N
orders (FK: user_id)
  ↓ 1:N
order_items (FK: order_id, product_id)
  ↓ N:1
products (PK: id)
```

**4 таблиці, 3 FK, realistic e-commerce structure**

### Тестові дані

```
3 users (John, Jane, Bob)
4 products (Laptop, Mouse, Keyboard, Monitor)
4 orders (3 з user_id, 1 guest)
6 order_items (різні комбінації)
```

---

## 📊 Queries покрито

### Basic JOINs

1. ✅ INNER JOIN - користувачі з замовленнями
2. ✅ LEFT JOIN - всі користувачі + замовлення
3. ✅ RIGHT JOIN - всі замовлення + користувачі
4. ✅ FULL OUTER JOIN - всі з обох
5. ✅ CROSS JOIN - всі комбінації

### Advanced Patterns

6. ✅ LEFT + WHERE IS NULL - знайти без зв'язків
7. ✅ Multiple JOINs - 4 таблиці (users → orders → order_items → products)
8. ✅ SELF JOIN - employees з managers
9. ✅ Aggregations - COUNT, SUM, AVG з GROUP BY
10. ✅ Subqueries з JOINs
11. ✅ LATERAL JOIN - останнє замовлення на користувача

### Go Integration

12. ✅ database/sql з PostgreSQL
13. ✅ Handling NULL values (*int, *float64, *string)
14. ✅ Scanning rows з JOINs
15. ✅ Formatted output

---

## 📊 Visual Guide

### JOIN Types

```
INNER:         LEFT:          RIGHT:         FULL:
┌──┐  ┌──┐    ┌██┐  ┌──┐    ┌──┐  ┌██┐    ┌██┐  ┌██┐
│  ├──┤  │    │██├──┤  │    │  ├──┤██│    │██├──┤██│
└──┘  └──┘    └██┘  └──┘    └──┘  └██┘    └██┘  └██┘
 Only match    All left      All right     All both
```

### Data Models

```
1:1              1:N              N:M
───              ───              ───
User             User             Student
 │                │                 │
 │ UNIQUE FK      │ FK              ├─ FK1
 ▼                ▼                 │
Profile          Posts            Junction
                                   │
                                   ├─ FK2
                                   │
                                  Course
```

---

## 🎯 Use Cases

### INNER JOIN
```sql
-- Тільки активні зв'язки
SELECT u.name, o.total
FROM users u
INNER JOIN orders o ON u.id = o.user_id;
```
**Коли:** Потрібні тільки користувачі з замовленнями

### LEFT JOIN
```sql
-- Всі користувачі (навіть без замовлень)
SELECT u.name, COUNT(o.id) AS orders
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
GROUP BY u.id, u.name;
```
**Коли:** Потрібні ВСІ користувачі

### LEFT + WHERE IS NULL
```sql
-- Знайти БЕЗ зв'язків
SELECT u.name
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
WHERE o.id IS NULL;
```
**Коли:** Знайти inactive users, unused products, etc.

### Multiple JOINs
```sql
-- Order details (4 tables)
SELECT u.name, o.id, p.name, oi.quantity
FROM users u
JOIN orders o ON u.id = o.user_id
JOIN order_items oi ON o.id = oi.order_id
JOIN products p ON oi.product_id = p.id;
```
**Коли:** Complex reporting

---

## ✅ Best Practices

### 1. Завжди вказуй таблицю

```sql
-- ❌ BAD
SELECT name, total FROM users JOIN orders;

-- ✅ GOOD
SELECT users.name, orders.total FROM users JOIN orders;
```

### 2. Використовуй aliases

```sql
-- ✅ GOOD
SELECT u.name, o.total
FROM users u
JOIN orders o ON u.id = o.user_id;
```

### 3. Index на FK

```sql
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
```

### 4. ON DELETE/UPDATE

```sql
FOREIGN KEY (user_id) REFERENCES users(id)
  ON DELETE CASCADE      -- Видалити orders при видаленні user
  ON UPDATE CASCADE;     -- Оновити при зміні id
```

### 5. Naming Convention

```sql
-- Таблиці: plural, lowercase
users, orders, order_items

-- FK: singular_id
user_id, order_id, product_id

-- Junction: table1_table2
student_courses, post_tags
```

---

## 🎓 Go Integration

### NULL Handling

```go
type UserWithOrders struct {
    UserName   string
    OrderID    *int     // NULL-safe
    OrderTotal *float64
    Status     *string
}

rows.Scan(&uo.UserName, &uo.OrderID, &uo.OrderTotal, &uo.Status)

if uo.OrderID == nil {
    fmt.Println("No orders")
} else {
    fmt.Printf("Order #%d\n", *uo.OrderID)
}
```

### Multiple Rows

```go
rows, err := db.Query(`
    SELECT u.name, o.id, o.total
    FROM users u
    LEFT JOIN orders o ON u.id = o.user_id
`)
defer rows.Close()

for rows.Next() {
    var name string
    var orderID *int
    var total *float64
    
    rows.Scan(&name, &orderID, &total)
    // Process
}
```

---

## 📊 Query Patterns Cheat Sheet

### Знайти без зв'язків
```sql
LEFT JOIN ... WHERE right.id IS NULL
```

### TOP N
```sql
LEFT JOIN ... GROUP BY ... ORDER BY COUNT(*) DESC LIMIT N
```

### Aggregation
```sql
LEFT JOIN ... GROUP BY ... COUNT/SUM/AVG
```

### Latest per group
```sql
LEFT JOIN LATERAL (
    SELECT * FROM orders
    WHERE orders.user_id = users.id
    ORDER BY created_at DESC
    LIMIT 1
) AS latest ON true
```

### Hierarchical (tree)
```sql
SELF JOIN employees m ON e.manager_id = m.id
```

---

## 🔗 Real-World Example

### E-commerce Schema

```sql
-- Core entities
users (id, name, email)
products (id, name, price)

-- Transactional
orders (id, user_id, total, status)
order_items (id, order_id, product_id, quantity, price)

-- Relationships
users 1:N orders (user_id FK)
orders 1:N order_items (order_id FK)
products 1:N order_items (product_id FK)

-- Effective N:M
orders N:M products (через order_items)
```

### Common Queries

**Order summary:**
```sql
SELECT u.name, o.id, SUM(oi.quantity * oi.price) AS total
FROM users u
JOIN orders o ON u.id = o.user_id
JOIN order_items oi ON o.id = oi.order_id
GROUP BY u.id, u.name, o.id;
```

**Top products:**
```sql
SELECT p.name, SUM(oi.quantity) AS sold
FROM products p
JOIN order_items oi ON p.id = oi.product_id
GROUP BY p.id, p.name
ORDER BY sold DESC;
```

**Inactive users:**
```sql
SELECT u.name
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
WHERE o.id IS NULL;
```

---

## 🎯 Висновок

### Data Models:

✅ **1:1** - UNIQUE FK (User ↔ Profile)  
✅ **1:N** - FK (User → Posts)  
✅ **N:M** - Junction table (Students ↔ Courses)  

### JOINs:

✅ **INNER** - тільки співпадіння  
✅ **LEFT** - всі з лівої + співпадіння  
✅ **RIGHT** - всі з правої + співпадіння  
✅ **FULL** - всі з обох  
✅ **CROSS** - всі комбінації  
✅ **SELF** - таблиця з собою  

### OSI Model:

✅ **7 layers** - Physical → Application  
✅ **TCP/UDP** - Transport layer protocols  
✅ **IP** - Network layer addressing  
✅ **Ports** - Transport endpoints  
✅ **Debugging** - по рівнях (ping, netstat, curl)  

### Normalization:

✅ **1NF** - atomic values, no repeating groups  
✅ **2NF** - no partial dependencies  
✅ **3NF** - no transitive dependencies  
✅ **Anomalies** - insert, update, delete  
✅ **Denormalization** - for performance (justified!)  

### Patterns:

✅ **LEFT + WHERE IS NULL** - знайти без зв'язків  
✅ **Multiple JOINs** - complex queries  
✅ **Aggregations** - COUNT, SUM, AVG  
✅ **Go integration** - database/sql, NULL handling  

---

## ✅ Week 14 Complete!

```
Progress: 100% ✅

Theory:   ████████████ 4/4 (Data Models, JOINs, OSI, Normalization)
Practice: ████████████ 2/2 (SQL, Go)
Docs:     ████████████ 6/6 (README, Quick Start, Complete, 3x Cheat Sheets)
```

**Дата завершення:** 2026-01-28  
**Статус:** COMPLETE ✅  
**Локація:** `/Users/vkuzm/GolandProjects/golang_practice/week_14`

---

## 🎉 Вітаємо!

Тепер ти вмієш:
- ✅ Проектувати моделі даних (1:1, 1:N, N:M)
- ✅ Використовувати всі типи JOINs
- ✅ Знаходити дані без зв'язків
- ✅ Писати складні queries з multiple JOINs
- ✅ Aggregations з GROUP BY
- ✅ Інтеграція з Go (database/sql)
- ✅ Handling NULL values в Go
- ✅ Optimization з indexes
- ✅ Розуміти 7 рівнів OSI моделі
- ✅ Debugging мережевих проблем
- ✅ TCP vs UDP різниця
- ✅ Common ports та протоколи
- ✅ Normalization (1NF, 2NF, 3NF)
- ✅ Identify anomalies
- ✅ Denormalization trade-offs

**"Data + Networks + Normalization = Full Stack!"** 🗄️🌐📊

---

**Next Steps:**
- Week 15: Transactions & ACID
- Week 16: Query Optimization & EXPLAIN
- Week 17: Advanced SQL (Window Functions, CTEs)

**Week 14: COMPLETE!** 🎯🗄️🌐
