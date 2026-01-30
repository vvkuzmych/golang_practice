# Week 14 - Data Models, SQL JOINs, OSI & Normalization

## 🎯 Мета

Розуміння моделей даних (1:1, 1:N, N:M), SQL JOINs (INNER, LEFT, RIGHT, FULL, CROSS, SELF), OSI Model (7 рівнів) та Database Normalization (1NF, 2NF, 3NF).

---

## 📚 Теорія

### [01: Data Models](./theory/01_data_models.md)

**Типи відношень:**
- **1:1** (One-to-One) - User ↔ Profile
- **1:N** (One-to-Many) - User → Posts
- **N:M** (Many-to-Many) - Students ↔ Courses

**Keys:**
- **Primary Key (PK)** - унікальний ідентифікатор
- **Foreign Key (FK)** - посилання на інший PK
- **Composite Key** - комбінація полів

---

### [02: SQL JOINs](./theory/02_sql_joins.md)

**JOIN Types:**
- **INNER** - тільки співпадіння
- **LEFT** - всі з лівої + співпадіння
- **RIGHT** - всі з правої + співпадіння
- **FULL** - всі з обох
- **CROSS** - всі комбінації
- **SELF** - таблиця з собою

---

### [03: OSI Model](./theory/03_osi_model.md)

**7 рівнів мережевої взаємодії:**
1. **Physical** - cables, bits, signals
2. **Data Link** - MAC, switches, frames
3. **Network** - IP, routing, packets
4. **Transport** - TCP/UDP, ports
5. **Session** - sessions, connections
6. **Presentation** - encryption, compression
7. **Application** - HTTP, DNS, SMTP

**Мнемоніка:** Please Do Not Throw Sausage Pizza Away

---

### [04: Database Normalization](./theory/04_normalization.md)

**Normal Forms (Нормальні форми):**
- **1NF** - atomic values (no "Math, Physics")
- **2NF** - no partial dependencies
- **3NF** - no transitive dependencies
- **BCNF** - advanced

**Anomalies:**
- Insert, Update, Delete anomalies
- Denormalization для performance

---

## 🎯 Практика

### [01: Basic JOINs (SQL)](./practice/01_basic_joins/)

**Файли:**
- `schema.sql` - створення таблиць + тестові дані
- `queries.sql` - приклади всіх JOIN types

**Запуск:**
```bash
# Створити базу
createdb joins_practice

# Виконати schema
psql -d joins_practice -f practice/01_basic_joins/schema.sql

# Спробувати queries
psql -d joins_practice -f practice/01_basic_joins/queries.sql
```

---

### [03: Go with JOINs](./practice/03_go_joins/)

**Go код для роботи з JOINs через database/sql**

**Запуск:**
```bash
cd practice/03_go_joins

# Встановити PostgreSQL driver
go get github.com/lib/pq

# Запустити
go run main.go
```

**Демонструє:**
- INNER JOIN - користувачі з замовленнями
- LEFT JOIN - всі користувачі + замовлення
- WHERE IS NULL - знайти без замовлень
- Multiple JOINs - деталі замовлення
- Aggregation - COUNT, SUM з GROUP BY

---

## 📊 Швидка довідка

### Data Models

```
1:1  User ──── Profile     (UNIQUE FK)
1:N  User ───< Posts       (FK)
N:M  Student >─< Courses   (Junction table)
```

### JOINs Cheat Sheet

```sql
-- Тільки співпадіння
INNER JOIN

-- Всі з лівої + співпадіння
LEFT JOIN

-- Всі з правої + співпадіння  
RIGHT JOIN

-- Всі з обох
FULL OUTER JOIN

-- Всі комбінації
CROSS JOIN

-- Таблиця з собою
SELF JOIN
```

### OSI Model (7 layers)

```
7. Application   ← HTTP, DNS, SMTP
6. Presentation  ← SSL/TLS, encryption
5. Session       ← Sessions
4. Transport     ← TCP, UDP, ports
3. Network       ← IP, routing
2. Data Link     ← MAC, switches
1. Physical      ← Cables, bits

Мнемоніка: Please Do Not Throw Sausage Pizza Away
```

### Normalization

```
Unnormalized → 1NF → 2NF → 3NF
(дублювання)           (оптимально)

1NF: Atomic values (no "Math, Physics")
2NF: No partial dependencies
3NF: No transitive dependencies

3NF = Sweet spot для більшості проектів
```

### Visual Guide

```
INNER:  ╔═══╗
        ║ A ║═════╗
        ╚═══╝     ║
                ╔═╩═╗
                ║ B ║
                ╚═══╝

LEFT:   ╔═══╗
        ║ A ║═════╗
        ║ A ║     ║
        ╚═══╝   ╔═╩═╗
                ║ B ║
                ╚═══╝

RIGHT:  ╔═══╗
        ║ A ║═════╗
        ╚═══╝     ║
                ╔═╩═╗
                ║ B ║
                ║ B ║
                ╚═══╝

FULL:   ╔═══╗
        ║ A ║═════╗
        ║ A ║     ║
        ╚═══╝   ╔═╩═╗
                ║ B ║
                ║ B ║
                ╚═══╝
```

---

## ✅ Приклади

### Знайти без зв'язків

```sql
-- Користувачі БЕЗ замовлень
SELECT u.name
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
WHERE o.id IS NULL;
```

### TOP користувачів

```sql
SELECT 
    u.name,
    COUNT(o.id) AS orders,
    SUM(o.total) AS spent
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
GROUP BY u.id, u.name
ORDER BY spent DESC;
```

### Multiple JOINs

```sql
SELECT 
    u.name,
    o.id,
    p.name,
    oi.quantity
FROM users u
JOIN orders o ON u.id = o.user_id
JOIN order_items oi ON o.id = oi.order_id
JOIN products p ON oi.product_id = p.id;
```

---

## 🎓 Key Points

### Models

✅ **1:1** → UNIQUE FK  
✅ **1:N** → FK  
✅ **N:M** → Junction table  

### JOINs

✅ **INNER** → співпадіння  
✅ **LEFT** → всі з лівої  
✅ **LEFT + WHERE NULL** → без зв'язків  

### Normalization

✅ **1NF** → atomic values  
✅ **2NF** → no partial dependencies  
✅ **3NF** → no transitive dependencies  
✅ **3NF = Sweet spot** для більшості проектів  

### Best Practices

✅ Index на FK  
✅ Завжди вказуй таблицю (users.id)  
✅ Використовуй aliases (u, o, p)  
✅ ON DELETE CASCADE для dependencies  
✅ Normalize до 3NF, denormalize for performance  

**Week 14: Data Models, JOINs, OSI & Normalization!** 🗄️🌐📊
