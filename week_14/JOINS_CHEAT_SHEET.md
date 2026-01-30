# SQL JOINs Cheat Sheet

## 📊 Visual Guide

```
INNER JOIN               LEFT JOIN               RIGHT JOIN              FULL OUTER JOIN
───────────             ───────────             ───────────             ───────────────
    A    B                  A    B                  A    B                  A    B
  ┌───┐ ┌───┐            ┌███┐ ┌───┐            ┌───┐ ┌███┐            ┌███┐ ┌███┐
  │   ├─┤   │            │███├─┤   │            │   ├─┤███│            │███├─┤███│
  └───┘ └───┘            └███┘ └───┘            └───┘ └███┘            └███┘ └███┘
Only overlap            All A + match          All B + match          All A + All B


CROSS JOIN              SELF JOIN
───────────            ───────────
Every A × Every B       Table with itself

A ∞ B                     employees
                            ↓
1,1  1,2  1,3           employees (managers)
2,1  2,2  2,3
3,1  3,2  3,3
```

---

## 🎯 Quick Reference

| JOIN Type | SQL | Результат | Use Case |
|-----------|-----|-----------|----------|
| **INNER** | `INNER JOIN` | Тільки співпадіння | Активні зв'язки |
| **LEFT** | `LEFT JOIN` | Всі з лівої + співпадіння | Всі основні + зв'язані |
| **RIGHT** | `RIGHT JOIN` | Всі з правої + співпадіння | Рідко (краще LEFT) |
| **FULL** | `FULL OUTER JOIN` | Всі з обох | Аудит зв'язків |
| **CROSS** | `CROSS JOIN` | A × B комбінації | Генерація комбінацій |
| **SELF** | Table with itself | Ієрархії | Employees, Categories |

---

## 💡 Common Patterns

### Pattern 1: Знайти без зв'язків

```sql
-- Користувачі БЕЗ замовлень
SELECT u.name
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
WHERE o.id IS NULL;
```

**Ключ:** `LEFT JOIN` + `WHERE right.id IS NULL`

---

### Pattern 2: TOP N з aggregation

```sql
-- TOP-5 користувачів за витратами
SELECT 
    u.name,
    COUNT(o.id) AS orders,
    SUM(o.total) AS spent
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
GROUP BY u.id, u.name
ORDER BY spent DESC
LIMIT 5;
```

**Ключ:** `LEFT JOIN` + `GROUP BY` + `ORDER BY` + `LIMIT`

---

### Pattern 3: Multiple JOINs

```sql
-- Деталі замовлення (4 таблиці)
SELECT 
    u.name AS customer,
    o.id AS order_id,
    p.name AS product,
    oi.quantity
FROM users u
INNER JOIN orders o ON u.id = o.user_id
INNER JOIN order_items oi ON o.id = oi.order_id
INNER JOIN products p ON oi.product_id = p.id;
```

**Ключ:** Ланцюг JOINs по FK

---

### Pattern 4: SELF JOIN (ієрархія)

```sql
-- Співробітники з їх менеджерами
SELECT 
    e.name AS employee,
    m.name AS manager
FROM employees e
LEFT JOIN employees m ON e.manager_id = m.id;
```

**Ключ:** Aliases (e, m) для розрізнення

---

### Pattern 5: LATERAL (останній на групу)

```sql
-- Останнє замовлення кожного користувача
SELECT 
    u.name,
    last_order.id,
    last_order.created_at
FROM users u
LEFT JOIN LATERAL (
    SELECT * FROM orders
    WHERE user_id = u.id
    ORDER BY created_at DESC
    LIMIT 1
) last_order ON true;
```

**Ключ:** `LATERAL` для correlated subquery

---

## 🔍 Decision Tree

```
Потрібні всі з основної таблиці?
├─ НІ  → INNER JOIN
└─ ТАК → LEFT JOIN
    │
    ├─ Знайти БЕЗ зв'язків?
    │  └─ ТАК → + WHERE right.id IS NULL
    │
    └─ Підрахунок (COUNT, SUM)?
       └─ ТАК → + GROUP BY
```

---

## ⚡ Performance Tips

### 1. Index на FK

```sql
CREATE INDEX idx_orders_user_id ON orders(user_id);
```

### 2. INNER замість LEFT (якщо можливо)

```sql
-- ✅ Швидше
INNER JOIN  -- Тільки співпадіння

-- ❌ Повільніше
LEFT JOIN   -- Всі + співпадіння
```

### 3. Фільтр в WHERE, не в JOIN

```sql
-- ❌ ПОВІЛЬНО
FROM users u
LEFT JOIN orders o ON u.id = o.user_id AND o.status = 'completed'

-- ✅ ШВИДКО
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
WHERE o.status = 'completed'
```

---

## 🎯 Go Code Patterns

### NULL Handling

```go
// Pointer для NULL
var orderID *int
rows.Scan(&orderID)

if orderID == nil {
    fmt.Println("No order")
} else {
    fmt.Printf("Order #%d\n", *orderID)
}
```

### Aggregation

```go
type UserStats struct {
    Name       string
    OrderCount int
    TotalSpent float64
}

query := `
    SELECT u.name, COUNT(o.id), COALESCE(SUM(o.total), 0)
    FROM users u
    LEFT JOIN orders o ON u.id = o.user_id
    GROUP BY u.id, u.name
`

rows, _ := db.Query(query)
defer rows.Close()

for rows.Next() {
    var stats UserStats
    rows.Scan(&stats.Name, &stats.OrderCount, &stats.TotalSpent)
    // Process
}
```

---

## 📝 Syntax Quick Copy

### INNER JOIN

```sql
SELECT a.*, b.*
FROM table_a a
INNER JOIN table_b b ON a.id = b.a_id;
```

### LEFT JOIN

```sql
SELECT a.*, b.*
FROM table_a a
LEFT JOIN table_b b ON a.id = b.a_id;
```

### Multiple JOINs

```sql
SELECT a.*, b.*, c.*, d.*
FROM table_a a
INNER JOIN table_b b ON a.id = b.a_id
INNER JOIN table_c c ON b.id = c.b_id
INNER JOIN table_d d ON c.id = d.c_id;
```

### SELF JOIN

```sql
SELECT 
    child.name AS child,
    parent.name AS parent
FROM table child
LEFT JOIN table parent ON child.parent_id = parent.id;
```

---

## ✅ Checklist

Перед написанням JOIN:

- [ ] Визначити основну таблицю (ліва)
- [ ] Визначити тип JOIN (INNER/LEFT/etc.)
- [ ] Перевірити FK для JOIN condition
- [ ] Index на FK колонках
- [ ] NULL handling в коді (якщо LEFT)
- [ ] Aliases для читабельності
- [ ] Specify table.column для disambiguation

---

**Week 14: JOINs Master!** 🗄️💪
