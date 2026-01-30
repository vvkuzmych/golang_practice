# Practice 03: Go with JOINs

## 🎯 Мета

Робота з SQL JOINs через Go database/sql пакет.

---

## 🚀 Швидкий старт

### 1. Prerequisites

```bash
# PostgreSQL повинен бути запущений
# База joins_practice повинна існувати (див. practice/01_basic_joins)

# Перевірити
psql -d joins_practice -c "SELECT COUNT(*) FROM users;"
```

### 2. Встановити dependencies

```bash
go mod init joins-practice
go get github.com/lib/pq
```

### 3. Налаштувати підключення

**Відредагувати `main.go`:**
```go
connStr := "host=localhost port=5432 user=postgres password=postgres dbname=joins_practice sslmode=disable"
```

### 4. Запустити

```bash
go run main.go
```

---

## 📊 Що демонструє

### 1. INNER JOIN

```go
func innerJoin(db *sql.DB) {
    query := `
        SELECT u.name, o.id, o.total, o.status
        FROM users u
        INNER JOIN orders o ON u.id = o.user_id
    `
    
    rows, _ := db.Query(query)
    defer rows.Close()
    
    for rows.Next() {
        var name string
        var orderID int
        var total float64
        var status string
        
        rows.Scan(&name, &orderID, &total, &status)
        fmt.Printf("%s | Order #%d | $%.2f\n", name, orderID, total)
    }
}
```

**Output:**
```
Jane Smith      | Order #3  | $375.00  | completed
John Doe        | Order #1  | $1300.00 | completed
John Doe        | Order #2  | $25.00   | pending
```

---

### 2. LEFT JOIN з NULL handling

```go
type UserWithOrders struct {
    UserName   string
    OrderID    *int     // NULL-safe
    OrderTotal *float64
    Status     *string
}

func leftJoin(db *sql.DB) {
    query := `
        SELECT u.name, o.id, o.total, o.status
        FROM users u
        LEFT JOIN orders o ON u.id = o.user_id
    `
    
    rows, _ := db.Query(query)
    defer rows.Close()
    
    for rows.Next() {
        var uo UserWithOrders
        rows.Scan(&uo.UserName, &uo.OrderID, &uo.OrderTotal, &uo.Status)
        
        if uo.OrderID == nil {
            fmt.Printf("%s | No orders\n", uo.UserName)
        } else {
            fmt.Printf("%s | Order #%d | $%.2f\n",
                uo.UserName, *uo.OrderID, *uo.OrderTotal)
        }
    }
}
```

**Output:**
```
Bob Wilson      | No orders
Jane Smith      | Order #3  | $375.00 | completed
John Doe        | Order #1  | $1300.00 | completed
John Doe        | Order #2  | $25.00 | pending
```

---

### 3. Знайти без замовлень

```go
query := `
    SELECT u.name
    FROM users u
    LEFT JOIN orders o ON u.id = o.user_id
    WHERE o.id IS NULL
`
```

**Output:**
```
❌ Bob Wilson - No orders
```

---

### 4. Multiple JOINs

```go
query := `
    SELECT 
        u.name,
        o.id,
        o.total,
        p.name,
        oi.quantity,
        oi.price
    FROM users u
    INNER JOIN orders o ON u.id = o.user_id
    INNER JOIN order_items oi ON o.id = oi.order_id
    INNER JOIN products p ON oi.product_id = p.id
`
```

**Output:**
```
📦 Order #1 - John Doe (Total: $1300.00)
   • Laptop x1 @ $1200.00 = $1200.00
   • Mouse x4 @ $25.00 = $100.00

📦 Order #2 - John Doe (Total: $25.00)
   • Mouse x1 @ $25.00 = $25.00

📦 Order #3 - Jane Smith (Total: $375.00)
   • Keyboard x1 @ $75.00 = $75.00
   • Monitor x1 @ $300.00 = $300.00
```

---

### 5. Aggregation

```go
query := `
    SELECT 
        u.name,
        COUNT(o.id) AS order_count,
        COALESCE(SUM(o.total), 0) AS total_spent
    FROM users u
    LEFT JOIN orders o ON u.id = o.user_id
    GROUP BY u.id, u.name
    ORDER BY total_spent DESC
`
```

**Output:**
```
User            | Orders | Total Spent
────────────────┼────────┼────────────
John Doe        | 2      | $1325.00
Jane Smith      | 1      | $375.00
Bob Wilson      | 0      | $0.00
```

---

## 🎯 Key Patterns

### NULL Handling

```go
// ❌ BAD - паніка якщо NULL
var orderID int
rows.Scan(&orderID)  // panic!

// ✅ GOOD - pointer для NULL
var orderID *int
rows.Scan(&orderID)

if orderID == nil {
    fmt.Println("No order")
} else {
    fmt.Printf("Order #%d\n", *orderID)
}
```

### Multiple Rows

```go
rows, err := db.Query(query)
if err != nil {
    log.Fatal(err)
}
defer rows.Close()  // ✅ ЗАВЖДИ!

for rows.Next() {
    // Scan і process
}

if err := rows.Err(); err != nil {
    log.Fatal(err)
}
```

### Single Row

```go
var name string
var total float64

err := db.QueryRow(`
    SELECT name, total FROM orders WHERE id = $1
`, 1).Scan(&name, &total)

if err == sql.ErrNoRows {
    fmt.Println("Not found")
} else if err != nil {
    log.Fatal(err)
}
```

---

## ✅ Best Practices

### 1. Defer Close

```go
rows, _ := db.Query(query)
defer rows.Close()  // ✅ ЗАВЖДИ!
```

### 2. Check errors

```go
if err := rows.Err(); err != nil {
    log.Fatal(err)
}
```

### 3. Pointers для NULL

```go
type User struct {
    ID    int
    Email *string  // NULL-safe
}
```

### 4. Prepared statements для повторних queries

```go
stmt, _ := db.Prepare("SELECT * FROM users WHERE id = $1")
defer stmt.Close()

for _, id := range userIDs {
    stmt.QueryRow(id).Scan(...)
}
```

---

## 🔧 Troubleshooting

### Connection refused

```bash
# Перевірити PostgreSQL
pg_isready

# Запустити
brew services start postgresql
```

### Database does not exist

```bash
# Створити
createdb joins_practice

# Виконати schema
psql -d joins_practice -f ../01_basic_joins/schema.sql
```

### Driver not found

```bash
go get github.com/lib/pq
```

---

## 📖 Далі

- Спробуй додати свої queries
- Додай error handling
- Додай connection pool config
- Спробуй GORM або sqlx

**Go + SQL = 💪**
