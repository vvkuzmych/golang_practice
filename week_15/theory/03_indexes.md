# Database Indexes

## 🎯 Що таке Index?

**Index** - це структура даних, яка прискорює пошук записів у таблиці.

**Аналогія:** Як алфавітний покажчик в кінці книги - швидко знайти сторінку за темою.

---

## 📊 How Indexes Work

### Without Index

```sql
SELECT * FROM users WHERE email = 'john@example.com';
```

**Full Table Scan:** перевіряє кожен рядок (O(n))

```
Row 1: alice@example.com   ❌
Row 2: bob@example.com     ❌
Row 3: john@example.com    ✅ FOUND
Row 4: mary@example.com    ❌
...
Row 1,000,000: ...         ❌
```

### With Index

```sql
CREATE INDEX idx_users_email ON users(email);
```

**Index Lookup:** знаходить за деревом (O(log n))

```
Index Tree:
          [m]
         /   \
      [d]     [t]
     /  \     /  \
   [a] [j]  [o] [z]
        |
   john@example.com → Row 3 ✅
```

---

## ✅ Advantages of Indexes

### 1. Faster SELECT with WHERE

```sql
-- Without index: Full table scan (slow)
SELECT * FROM users WHERE email = 'john@example.com';

-- With index: Index lookup (fast)
-- Speedup: 1000x for large tables
```

### 2. Faster JOIN

```sql
-- Index on foreign keys speeds up joins
SELECT o.id, u.name
FROM orders o
JOIN users u ON o.user_id = u.id;
```

### 3. Faster ORDER BY

```sql
-- Index on created_at speeds up sorting
SELECT * FROM posts ORDER BY created_at DESC LIMIT 10;
```

### 4. Faster Aggregations

```sql
-- Index on status speeds up GROUP BY
SELECT status, COUNT(*) 
FROM orders 
GROUP BY status;
```

---

## ❌ Downsides of Indexes

### 1. Slower Writes (INSERT, UPDATE, DELETE)

**Problem:** Index must be updated on every write.

```sql
INSERT INTO users (name, email) VALUES ('John', 'john@example.com');
```

**Without index:**
```
1. Write row to table
   Done! ✅ (1 operation)
```

**With 3 indexes (id, email, name):**
```
1. Write row to table
2. Update index on id
3. Update index on email
4. Update index on name
   Done! ✅ (4 operations)
```

**Impact:**
- INSERT: ~30-50% slower per index
- UPDATE: Depends on which columns change
- DELETE: All indexes must be updated

---

### 2. Extra Disk Space

**Table size:** 1 GB  
**Index on email:** +200 MB  
**Index on name:** +150 MB  
**Index on created_at:** +100 MB  

**Total:** 1 GB → 1.45 GB (45% overhead)

**Rule of thumb:** Each index adds ~10-30% of table size.

---

### 3. Index Maintenance Overhead

**B-Tree rebalancing:**
- After many INSERTs, tree becomes unbalanced
- DB needs to rebalance → expensive operation
- Can cause temporary slowdowns

**Fragmentation:**
- Over time, index pages become fragmented
- Requires VACUUM (PostgreSQL) or OPTIMIZE (MySQL)

---

### 4. Query Planner Confusion

**Too many indexes = harder to choose.**

```sql
-- Table with 10 indexes on different columns
SELECT * FROM users 
WHERE email = 'john@example.com' 
  AND country = 'USA' 
  AND age > 18;
```

**Query planner must decide:**
- Use index on email?
- Use index on country?
- Use index on age?
- Use composite index?
- Use multiple indexes?

**Problem:** Wrong choice → slow query.

---

## 🎯 Types of Indexes

### 1. B-Tree Index (Default)

```sql
CREATE INDEX idx_users_email ON users(email);
```

**Best for:**
- Equality: `WHERE email = 'john@example.com'`
- Range: `WHERE age BETWEEN 18 AND 30`
- Sorting: `ORDER BY created_at`
- Prefix: `WHERE name LIKE 'John%'`

**Structure:** Balanced tree (O(log n))

---

### 2. Hash Index

```sql
CREATE INDEX idx_users_email ON users USING HASH (email);
```

**Best for:**
- Equality only: `WHERE email = 'john@example.com'`

**Cannot do:**
- Range queries: `WHERE age > 18`
- Sorting: `ORDER BY`
- Prefix: `WHERE name LIKE 'John%'`

**Performance:** O(1) for equality, but limited use cases.

---

### 3. Composite Index

```sql
CREATE INDEX idx_users_country_age ON users(country, age);
```

**Can use for:**
- `WHERE country = 'USA'` ✅
- `WHERE country = 'USA' AND age > 18` ✅

**Cannot use for:**
- `WHERE age > 18` ❌ (age is second column)

**Rule:** Index works left-to-right.

---

### 4. Unique Index

```sql
CREATE UNIQUE INDEX idx_users_email ON users(email);
```

**Benefits:**
- Enforces uniqueness
- Same speed as regular index

---

### 5. Partial Index

```sql
CREATE INDEX idx_active_users ON users(email) WHERE active = true;
```

**Benefits:**
- Smaller index (only active users)
- Faster updates (only updates when active changes)

---

## 📊 When to Use Indexes

### ✅ Good Candidates

1. **Foreign keys**
```sql
CREATE INDEX idx_orders_user_id ON orders(user_id);
```

2. **Columns in WHERE**
```sql
-- Often queried
SELECT * FROM users WHERE email = ?;
CREATE INDEX idx_users_email ON users(email);
```

3. **Columns in JOIN**
```sql
SELECT * FROM orders o JOIN users u ON o.user_id = u.id;
CREATE INDEX idx_orders_user_id ON orders(user_id);
```

4. **Columns in ORDER BY**
```sql
SELECT * FROM posts ORDER BY created_at DESC;
CREATE INDEX idx_posts_created_at ON posts(created_at);
```

5. **High cardinality columns**
```sql
-- email: unique for each user ✅
CREATE INDEX idx_users_email ON users(email);

-- gender: only 2-3 values ❌ (not useful)
```

---

### ❌ Bad Candidates

1. **Low cardinality**
```sql
-- gender: male/female/other (only 3 values)
-- Index не допоможе
```

2. **Small tables**
```sql
-- Table with 100 rows
-- Full scan is faster than index lookup
```

3. **Frequently updated columns**
```sql
-- last_login updated every request
-- Index would slow down every UPDATE
```

4. **Columns rarely used in WHERE**
```sql
-- description: never used in WHERE
-- No need for index
```

---

## 🔍 Measuring Index Performance

### EXPLAIN Query Plan

```sql
-- PostgreSQL
EXPLAIN ANALYZE
SELECT * FROM users WHERE email = 'john@example.com';
```

**Without index:**
```
Seq Scan on users  (cost=0.00..180.00 rows=1 width=100) (actual time=50.123..50.123 rows=1 loops=1)
  Filter: (email = 'john@example.com'::text)
  Rows Removed by Filter: 9999
Planning Time: 0.123 ms
Execution Time: 50.456 ms
```

**With index:**
```
Index Scan using idx_users_email on users  (cost=0.29..8.30 rows=1 width=100) (actual time=0.045..0.046 rows=1 loops=1)
  Index Cond: (email = 'john@example.com'::text)
Planning Time: 0.234 ms
Execution Time: 0.098 ms
```

**Speedup:** 50ms → 0.1ms (500x faster!)

---

## 🎯 Index Strategies

### Strategy 1: Start with Foreign Keys

```sql
-- Always index foreign keys
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_product_id ON order_items(product_id);
```

---

### Strategy 2: Index Frequent WHERE Columns

```sql
-- Analyze slow queries
SELECT * FROM users WHERE email = ?;  -- Frequent
SELECT * FROM users WHERE country = ?;  -- Rare

-- Index only frequent
CREATE INDEX idx_users_email ON users(email);
```

---

### Strategy 3: Composite for Multiple Columns

```sql
-- Query uses both columns
SELECT * FROM orders WHERE user_id = ? AND status = ?;

-- Composite index
CREATE INDEX idx_orders_user_status ON orders(user_id, status);
```

---

### Strategy 4: Partial for Subsets

```sql
-- Most queries filter by active users
SELECT * FROM users WHERE active = true AND email = ?;

-- Partial index
CREATE INDEX idx_active_users_email ON users(email) WHERE active = true;
```

---

## ✅ Best Practices

### 1. Don't Over-Index

```sql
-- ❌ BAD - 10 indexes on one table
CREATE INDEX idx1 ON users(email);
CREATE INDEX idx2 ON users(name);
CREATE INDEX idx3 ON users(country);
CREATE INDEX idx4 ON users(age);
CREATE INDEX idx5 ON users(gender);
-- ...

-- ✅ GOOD - Only necessary indexes
CREATE INDEX idx_users_email ON users(email);  -- Frequent WHERE
CREATE INDEX idx_users_country_age ON users(country, age);  -- Composite
```

### 2. Monitor Query Performance

```sql
-- Use EXPLAIN to check if index is used
EXPLAIN SELECT * FROM users WHERE email = ?;
```

### 3. Consider Write/Read Ratio

```
High reads, low writes → More indexes OK ✅
High writes, low reads → Fewer indexes ✅
```

### 4. Composite Index Order Matters

```sql
-- Query: WHERE country = 'USA' AND age > 18
CREATE INDEX idx_users_country_age ON users(country, age);  -- ✅

-- Query: WHERE age > 18 AND country = 'USA'
-- Same index works! ✅ (equality comes first)
```

### 5. Unique Constraints are Indexes

```sql
-- This creates an index automatically
ALTER TABLE users ADD CONSTRAINT unique_email UNIQUE (email);
-- No need for: CREATE INDEX idx_users_email ON users(email);
```

---

## 🎓 Trade-offs Summary

| Aspect | Without Index | With Index |
|--------|---------------|------------|
| **SELECT speed** | Slow (O(n)) | Fast (O(log n)) |
| **INSERT speed** | Fast | Slower (~30-50% per index) |
| **UPDATE speed** | Fast | Slower (if indexed column changes) |
| **DELETE speed** | Fast | Slower (~30-50% per index) |
| **Disk space** | Less | More (+10-30% per index) |
| **Maintenance** | None | Rebalancing, vacuuming |

---

## 🎯 Висновок

### Indexes are:

✅ **Essential** for performance  
✅ **Fast** reads (SELECT, JOIN, ORDER BY)  
❌ **Slow** writes (INSERT, UPDATE, DELETE)  
❌ **Extra** disk space  
❌ **Maintenance** overhead  

### Golden Rules:

1. **Index foreign keys** - always
2. **Index WHERE columns** - if frequently queried
3. **Composite indexes** - for multi-column queries
4. **Don't over-index** - every index has a cost
5. **Monitor performance** - use EXPLAIN

### Decision Process:

```
1. Identify slow queries (EXPLAIN)
2. Check if column is in WHERE/JOIN/ORDER BY
3. Consider write frequency
4. Add index if reads >> writes
5. Measure impact (before/after)
```

**Week 15: Database Indexes Master!** 🗂️⚡
