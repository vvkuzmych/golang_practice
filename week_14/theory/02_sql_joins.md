# SQL JOINs

## 🎯 Що таке JOIN?

**JOIN** - оператор для об'єднання даних з двох або більше таблиць.

---

## 📊 Тестові дані

```sql
-- Users
users
─────────────
id | name
1  | John
2  | Jane
3  | Bob

-- Orders
orders
─────────────────
id | user_id | total
1  | 1       | 100
2  | 1       | 200
3  | 2       | 150
4  | NULL    | 50   ← Гостьове замовлення
```

---

## 1. INNER JOIN

**Повертає тільки співпадіння з обох таблиць**

```
users ∩ orders
```

### Візуалізація

```
users          orders
┌──┐          ┌──┐
│  │          │  │
│  ├──────────┤  │  ← Тільки перетин
│  │          │  │
└──┘          └──┘
```

### SQL

```sql
SELECT users.name, orders.total
FROM users
INNER JOIN orders ON users.id = orders.user_id;
```

### Результат

```
name  | total
──────┼──────
John  | 100
John  | 200
Jane  | 150
```

**Пропущено:**
- Bob (немає замовлень)
- Order #4 (немає user_id)

---

## 2. LEFT JOIN (LEFT OUTER JOIN)

**Всі з лівої таблиці + співпадіння з правої**

```
users ∪ orders (all users)
```

### Візуалізація

```
users          orders
┌──┐          ┌──┐
│██│          │  │
│██├──────────┤  │  ← Вся ліва + перетин
│██│          │  │
└──┘          └──┘
```

### SQL

```sql
SELECT users.name, orders.total
FROM users
LEFT JOIN orders ON users.id = orders.user_id;
```

### Результат

```
name  | total
──────┼──────
John  | 100
John  | 200
Jane  | 150
Bob   | NULL   ← Немає замовлень, але є в результаті
```

**Використання:** Знайти користувачів БЕЗ замовлень

```sql
SELECT users.name
FROM users
LEFT JOIN orders ON users.id = orders.user_id
WHERE orders.id IS NULL;

-- Результат: Bob
```

---

## 3. RIGHT JOIN (RIGHT OUTER JOIN)

**Всі з правої таблиці + співпадіння з лівої**

```
users ∪ orders (all orders)
```

### Візуалізація

```
users          orders
┌──┐          ┌──┐
│  │          │██│
│  ├──────────┤██│  ← Вся права + перетин
│  │          │██│
└──┘          └──┘
```

### SQL

```sql
SELECT users.name, orders.total
FROM users
RIGHT JOIN orders ON users.id = orders.user_id;
```

### Результат

```
name  | total
──────┼──────
John  | 100
John  | 200
Jane  | 150
NULL  | 50    ← Гостьове замовлення
```

---

## 4. FULL OUTER JOIN

**Всі з обох таблиць**

```
users ∪ orders (all)
```

### Візуалізація

```
users          orders
┌──┐          ┌──┐
│██│          │██│
│██├──────────┤██│  ← Все з обох
│██│          │██│
└──┘          └──┘
```

### SQL

```sql
SELECT users.name, orders.total
FROM users
FULL OUTER JOIN orders ON users.id = orders.user_id;
```

### Результат

```
name  | total
──────┼──────
John  | 100
John  | 200
Jane  | 150
Bob   | NULL   ← Без замовлень
NULL  | 50     ← Гостьове замовлення
```

---

## 5. CROSS JOIN (Cartesian Product)

**Кожен рядок з першої × кожен рядок з другої**

```
users × orders
```

### SQL

```sql
SELECT users.name, orders.total
FROM users
CROSS JOIN orders;
```

### Результат

```
name  | total
──────┼──────
John  | 100
John  | 200
John  | 150
John  | 50
Jane  | 100
Jane  | 200
Jane  | 150
Jane  | 50
Bob   | 100
Bob   | 200
Bob   | 150
Bob   | 50
```

**Всього:** 3 users × 4 orders = 12 рядків

**Використання:** Генерація комбінацій (рідко)

---

## 6. SELF JOIN

**Таблиця з'єднується сама з собою**

### Приклад: Employees з manager_id

```sql
employees
─────────────────────────
id | name   | manager_id
1  | Alice  | NULL       ← CEO
2  | Bob    | 1          ← Manager: Alice
3  | Carol  | 1          ← Manager: Alice
4  | Dave   | 2          ← Manager: Bob
```

### SQL

```sql
SELECT 
    e.name AS employee,
    m.name AS manager
FROM employees e
LEFT JOIN employees m ON e.manager_id = m.id;
```

### Результат

```
employee | manager
─────────┼─────────
Alice    | NULL
Bob      | Alice
Carol    | Alice
Dave     | Bob
```

---

## 📊 JOIN Summary

| JOIN | Лівий NULL | Правий NULL | Використання |
|------|------------|-------------|--------------|
| INNER | ❌ | ❌ | Тільки співпадіння |
| LEFT | ✅ | ❌ | Всі з лівої + співпадіння |
| RIGHT | ❌ | ✅ | Всі з правої + співпадіння |
| FULL | ✅ | ✅ | Всі з обох |
| CROSS | - | - | Всі комбінації |

---

## 🎯 Multiple JOINs

```sql
SELECT 
    users.name,
    orders.id AS order_id,
    products.name AS product_name
FROM users
INNER JOIN orders ON users.id = orders.user_id
INNER JOIN order_items ON orders.id = order_items.order_id
INNER JOIN products ON order_items.product_id = products.id;
```

**Ланцюг:**
```
users → orders → order_items → products
```

---

## 🔍 JOIN Conditions

### ON (standard)

```sql
FROM users
JOIN orders ON users.id = orders.user_id
```

### USING (якщо однакова назва колонки)

```sql
FROM users
JOIN orders USING (user_id)
```

**Працює якщо:**
- Обидві таблиці мають колонку `user_id`
- Однакова назва

### NATURAL JOIN (автоматично по всіх однакових колонках)

```sql
FROM users
NATURAL JOIN orders
```

**⚠️ Небезпечно!** Може з'єднати не по тих колонках.

---

## ✅ Best Practices

### 1. Завжди вказуй таблицю

```sql
-- ❌ BAD (не ясно звідки колонка)
SELECT name, total
FROM users
JOIN orders ON users.id = orders.user_id;

-- ✅ GOOD
SELECT users.name, orders.total
FROM users
JOIN orders ON users.id = orders.user_id;
```

### 2. Використовуй aliases

```sql
-- ✅ GOOD
SELECT u.name, o.total
FROM users u
JOIN orders o ON u.id = o.user_id;
```

### 3. Index на JOIN колонках

```sql
CREATE INDEX idx_orders_user_id ON orders(user_id);
```

### 4. INNER vs LEFT - залежить від задачі

```sql
-- Тільки користувачі з замовленнями
INNER JOIN

-- Всі користувачі (навіть без замовлень)
LEFT JOIN
```

---

## 🎯 Практичні приклади

### 1. Знайти користувачів БЕЗ замовлень

```sql
SELECT u.name
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
WHERE o.id IS NULL;
```

### 2. Знайти замовлення БЕЗ користувача (гостьові)

```sql
SELECT o.*
FROM orders o
LEFT JOIN users u ON o.user_id = u.id
WHERE u.id IS NULL;
```

### 3. TOP користувачів за кількістю замовлень

```sql
SELECT 
    u.name,
    COUNT(o.id) AS order_count
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
GROUP BY u.id, u.name
ORDER BY order_count DESC;
```

### 4. Середня сума замовлення на користувача

```sql
SELECT 
    u.name,
    AVG(o.total) AS avg_order
FROM users u
INNER JOIN orders o ON u.id = o.user_id
GROUP BY u.id, u.name;
```

---

## 🎓 Висновок

### JOIN Types:

✅ **INNER** - тільки співпадіння  
✅ **LEFT** - всі з лівої + співпадіння  
✅ **RIGHT** - всі з правої + співпадіння  
✅ **FULL** - всі з обох  
✅ **CROSS** - всі комбінації  
✅ **SELF** - таблиця з собою  

### Коли використовувати:

| Задача | JOIN |
|--------|------|
| Тільки зв'язані дані | INNER |
| Всі з основної + зв'язані | LEFT |
| Знайти незв'язані | LEFT + WHERE IS NULL |
| Ієрархія (employees) | SELF |

**Далі:** Практика з реальними прикладами!
