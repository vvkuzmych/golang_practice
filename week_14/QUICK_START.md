# Week 14 - Quick Start

## ⚡ Швидкий старт (3 хвилини)

### 1. Створити базу даних

```bash
# PostgreSQL
createdb joins_practice
```

### 2. Виконати schema

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_14

psql -d joins_practice -f practice/01_basic_joins/schema.sql
```

**Результат:**
```
CREATE TABLE
CREATE TABLE
CREATE TABLE
CREATE TABLE
INSERT 0 3  (users)
INSERT 0 4  (products)
INSERT 0 4  (orders)
INSERT 0 6  (order_items)
```

### 3. Спробувати JOINs

```bash
# INNER JOIN
psql -d joins_practice -c "
SELECT u.name, o.total
FROM users u
INNER JOIN orders o ON u.id = o.user_id;
"

# LEFT JOIN (всі користувачі)
psql -d joins_practice -c "
SELECT u.name, o.total
FROM users u
LEFT JOIN orders o ON u.id = o.user_id;
"

# Знайти без замовлень
psql -d joins_practice -c "
SELECT u.name
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
WHERE o.id IS NULL;
"
```

### 4. Go приклад

```bash
cd practice/03_go_joins

# Встановити driver
go mod init joins-practice
go get github.com/lib/pq

# Запустити
go run main.go
```

---

## 📊 Тестові дані

```
Users:
- John Doe (2 orders)
- Jane Smith (1 order)
- Bob Wilson (0 orders)

Products:
- Laptop ($1200)
- Mouse ($25)
- Keyboard ($75)
- Monitor ($300)

Orders:
#1: John - Laptop + Mouse ($1300)
#2: John - Mouse ($25)
#3: Jane - Monitor + Keyboard ($375)
#4: Guest - Keyboard ($75)
```

---

## 🎯 Основні запити

### INNER JOIN

```sql
SELECT u.name, o.total
FROM users u
INNER JOIN orders o ON u.id = o.user_id;
```

### LEFT JOIN

```sql
SELECT u.name, o.total
FROM users u
LEFT JOIN orders o ON u.id = o.user_id;
```

### Без зв'язків

```sql
SELECT u.name
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
WHERE o.id IS NULL;
```

### Aggregation

```sql
SELECT 
    u.name,
    COUNT(o.id) AS orders,
    SUM(o.total) AS total
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
GROUP BY u.id, u.name;
```

---

## 📖 Далі

- `theory/01_data_models.md` - Типи відношень
- `theory/02_sql_joins.md` - Всі JOIN types
- `practice/01_basic_joins/queries.sql` - 50+ прикладів

---

## 🌐 OSI Model Quick Reference

```
7. Application   → HTTP, DNS, SMTP (your code!)
4. Transport     → TCP/UDP, ports
3. Network       → IP addresses
1. Physical      → Cables, WiFi

TCP = reliable, UDP = fast
Port 443 = HTTPS, 5432 = PostgreSQL
```

**Debugging:**
```bash
ping 8.8.8.8           # Network layer
netstat -tuln          # Transport layer
curl https://example.com  # Application layer
```

---

## 📊 Normalization Quick Reference

```
1NF: One value per cell
     ❌ "Math, Physics"
     ✅ 2 rows: Math, Physics

2NF: No partial dependencies
     ❌ product_name depends on product_id (part of key)
     ✅ products table separate

3NF: No transitive dependencies
     ❌ dept_name depends on dept_id (not PK)
     ✅ departments table separate
```

**Golden Rule:** "The key, the whole key, and nothing but the key"

**Week 14: JOINs, OSI & Normalization!** 🗄️🌐📊
