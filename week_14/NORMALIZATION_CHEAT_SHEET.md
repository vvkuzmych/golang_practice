# Database Normalization Cheat Sheet

## 📊 Normal Forms Quick Guide

```
Unnormalized → 1NF → 2NF → 3NF → BCNF
(дублювання)          (оптимально) (рідко)
```

---

## 1️⃣ First Normal Form (1NF)

**Правило:** Atomic values + Primary key

### ❌ Before

```sql
students
─────────────────────────────────
id | name  | courses
1  | John  | Math, Physics, CS     ← Multiple values!
```

### ✅ After

```sql
students              student_courses
─────────────        ────────────────────
id | name            student_id | course
1  | John            1          | Math
2  | Jane            1          | Physics
                     1          | CS
```

**Remember:** 1 cell = 1 value

---

## 2️⃣ Second Normal Form (2NF)

**Правило:** 1NF + No partial dependencies

### ❌ Before

```sql
order_items
──────────────────────────────────────────
order_id | product_id | product_name | price | qty
1        | 101        | Laptop       | 1200  | 1
         ↑            ↑
    Composite PK   Depends only on product_id!
```

### ✅ After

```sql
products              order_items
────────────────     ─────────────────────
id  | name  | price  order_id | product_id | qty
101 | Laptop| 1200   1        | 101        | 1
```

**Remember:** No field depends on PART of the key

---

## 3️⃣ Third Normal Form (3NF)

**Правило:** 2NF + No transitive dependencies

### ❌ Before

```sql
employees
────────────────────────────────────────────
id | name | dept_id | dept_name   | dept_location
1  | John | 10      | Engineering | Building A
2  | Jane | 10      | Engineering | Building A   ← Duplicate!
                ↓          ↓
           dept_name depends on dept_id, not id!
```

### ✅ After

```sql
departments           employees
─────────────────    ────────────────
id | name      | loc  id | name | dept_id
10 | Eng       | A    1  | John | 10
20 | Sales     | B    2  | Jane | 10
```

**Remember:** Non-key fields depend ONLY on PK

---

## 🎯 Quick Decision Tree

```
┌─ Multiple values in one cell? (Math, Physics)
│  └─ YES → 1NF (split to rows)
│
┌─ Field depends on PART of composite key?
│  └─ YES → 2NF (extract to separate table)
│
┌─ Non-key field depends on another non-key field?
│  └─ YES → 3NF (extract to separate table)
│
└─ Done! (3NF достатньо)
```

---

## 📊 Anomalies

### Insert Anomaly
```
❌ Can't add customer without order
✅ 3NF: customers table separate
```

### Update Anomaly
```
❌ Update email = multiple updates + errors
✅ 3NF: update email in 1 place
```

### Delete Anomaly
```
❌ Delete last order = lose customer info
✅ 3NF: customer exists independently
```

---

## ⚖️ Normalization vs Denormalization

### Normalize (3NF)

```sql
-- Good for writes
✅ No duplicates
✅ Easy updates
❌ Slow reads (many JOINs)

Use: transactional data, frequent updates
```

### Denormalize

```sql
-- Good for reads
✅ Fast reads (no JOINs)
❌ Duplicates
❌ Complex updates

Use: analytics, dashboards, read-heavy
```

---

## 🎯 Real Example

### Unnormalized

```sql
orders
─────────────────────────────────────────────────
id | customer | email         | product | price
1  | John     | john@mail.com | Laptop  | 1200
2  | John     | john@mail.com | Mouse   | 25    ← Duplicate!
```

**Problems:**
- Email duplicated
- Update email = 2 updates
- Can't add customer without order

### 3NF

```sql
customers            products           orders
───────────────     ────────────────   ─────────────────
id | name | email   id | name | price  id | cust_id | prod_id
1  | John | j@m.c   1  | Laptop| 1200  1  | 1       | 1
                    2  | Mouse | 25    2  | 1       | 2
```

**Benefits:**
- ✅ No duplicates
- ✅ Update email once
- ✅ Can add customer without order

---

## 💡 Practical Tips

### Always 3NF for:
- 📝 Core tables (users, products, orders)
- 🔄 Frequent updates
- 🔐 Data integrity critical

### Denormalize for:
- 📊 Reporting tables
- 🚀 Performance bottlenecks (after measuring!)
- 📈 Analytics/dashboards

### Hybrid Approach (Best)

```sql
-- Core (3NF) - for writes
CREATE TABLE orders (...);
CREATE TABLE products (...);

-- Denormalized view - for reads
CREATE MATERIALIZED VIEW order_summary AS
SELECT 
    o.id,
    c.name,
    SUM(oi.price * oi.qty) as total
FROM orders o
JOIN customers c ...
JOIN order_items oi ...;

-- ✅ Best of both worlds!
```

---

## 🎓 Remember

### 1NF
```
One cell = One value
```

### 2NF
```
All fields depend on WHOLE key
```

### 3NF
```
Non-key fields depend ONLY on key
```

### Golden Rule
```
"The key, the whole key, and nothing but the key"
```

---

## ✅ Checklist

Before deploying database:

- [ ] All tables in at least 1NF? (atomic values)
- [ ] No partial dependencies? (2NF)
- [ ] No transitive dependencies? (3NF)
- [ ] Indexes on foreign keys?
- [ ] Denormalization justified by metrics?

---

## 🔍 Common Mistakes

### ❌ Over-normalization
```sql
-- TOO MUCH (5NF, 6NF)
user_first_names
user_last_names
user_middle_names
```

### ❌ Premature denormalization
```sql
-- Before measuring performance
orders_with_everything  -- Just use 3NF + indexes first!
```

### ✅ Right approach
```sql
-- Start with 3NF
-- Add indexes
-- Measure performance
-- Denormalize ONLY if needed
```

---

**Normalization = Clean Data Design!** 🗄️📊
