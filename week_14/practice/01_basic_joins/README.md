# Practice 01: Basic SQL JOINs

## 🎯 Мета

Практика всіх типів SQL JOINs з реальними даними.

---

## 📊 Schema

```
users
  ↓ 1:N
orders
  ↓ 1:N
order_items
  ↓ N:1
products
```

**4 таблиці, 3 foreign keys**

---

## 🚀 Швидкий старт

### 1. Створити базу

```bash
createdb joins_practice
```

### 2. Виконати schema

```bash
psql -d joins_practice -f schema.sql
```

**Результат:**
```
CREATE TABLE users
CREATE TABLE products
CREATE TABLE orders
CREATE TABLE order_items
INSERT 0 3  (users)
INSERT 0 4  (products)
INSERT 0 4  (orders)
INSERT 0 6  (order_items)
```

### 3. Спробувати queries

```bash
# Інтерактивно
psql -d joins_practice

# Або всі queries одразу
psql -d joins_practice -f queries.sql
```

---

## 📚 Queries покрито

### Basic JOINs

1. **INNER JOIN** - користувачі з замовленнями
2. **LEFT JOIN** - всі користувачі + замовлення
3. **RIGHT JOIN** - всі замовлення + користувачі
4. **FULL OUTER JOIN** - всі з обох
5. **CROSS JOIN** - всі комбінації

### Patterns

6. **LEFT + WHERE IS NULL** - знайти без замовлень
7. **Multiple JOINs** - 4 таблиці
8. **SELF JOIN** - employees з managers
9. **Aggregations** - COUNT, SUM, AVG
10. **Subqueries** - користувачі з total > $500

---

## 🎯 Тестові дані

```
Users:
1. John Doe   (john@example.com)
2. Jane Smith (jane@example.com)
3. Bob Wilson (bob@example.com)

Products:
1. Laptop  ($1200, stock: 10)
2. Mouse   ($25, stock: 50)
3. Keyboard ($75, stock: 30)
4. Monitor ($300, stock: 15)

Orders:
1. John - $1300 (completed)
2. John - $25 (pending)
3. Jane - $375 (completed)
4. Guest - $75 (pending)

Order Items:
Order #1: Laptop x1, Mouse x4
Order #2: Mouse x1
Order #3: Monitor x1, Keyboard x1
Order #4: Keyboard x1
```

---

## 📖 Приклади

### INNER JOIN

```sql
SELECT u.name, o.total
FROM users u
INNER JOIN orders o ON u.id = o.user_id;

-- Результат: 3 рядки (John x2, Jane x1)
-- Пропущено: Bob, Order #4 (guest)
```

### LEFT JOIN

```sql
SELECT u.name, o.total
FROM users u
LEFT JOIN orders o ON u.id = o.user_id;

-- Результат: 4 рядки (включаючи Bob з NULL)
```

### Без замовлень

```sql
SELECT u.name
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
WHERE o.id IS NULL;

-- Результат: Bob Wilson
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

-- John: 2 orders, $1325
-- Jane: 1 order, $375
-- Bob: 0 orders, $0
```

---

## ✅ Перевірка

```bash
# Всі користувачі
psql -d joins_practice -c "SELECT * FROM users;"

# Всі замовлення
psql -d joins_practice -c "SELECT * FROM orders;"

# INNER JOIN
psql -d joins_practice -c "
SELECT u.name, o.total 
FROM users u 
INNER JOIN orders o ON u.id = o.user_id;
"
```

---

**Готово!** Тепер спробуй всі queries з `queries.sql` 🚀
