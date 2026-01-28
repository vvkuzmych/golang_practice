# Atomic Transactions (Атомарні Транзакції)

## 🎯 Що таке "atomic"?

**Atomic** = "неподільний", "all-or-nothing"

Транзакція **або повністю виконується, або не виконується взагалі**. Не може бути "наполовину виконана".

---

## 📖 Візуалізація: Чому це важливо

### ❌ БЕЗ транзакції (небезпечно!)

```
1. UPDATE accounts SET balance = balance - 100  ✅ Виконано
         ↓
   💥 CRASH / NETWORK ERROR / POWER OUTAGE
         ↓
2. INSERT INTO outbox (event)                   ❌ НЕ виконано

Result: Гроші зняті, але event не записаний!
        Інші сервіси не знають про зняття!
```

### ✅ З транзакцією (безпечно!)

```
BEGIN TRANSACTION
         |
         ├─ 1. UPDATE accounts SET balance = balance - 100
         |
         ├─ 2. INSERT INTO outbox (event)
         |
         └─ COMMIT ─────┐
                        │
                        ├─> Якщо SUCCESS: обидві операції збережені ✅
                        │
                        └─> Якщо ERROR: обидві операції скасовані ❌
                            (ROLLBACK)

Result: Або обидві виконані, або жодна!
        Consistency гарантована! 🎯
```

---

## 💻 Приклад в коді

### ❌ БЕЗ atomic transaction (небезпечно)

```go
func WithdrawMoney(userID int64, amount float64) error {
    // Step 1: Deduct money
    _, err := db.Exec(
        "UPDATE accounts SET balance = balance - $1 WHERE user_id = $2",
        amount, userID,
    )
    if err != nil {
        return err
    }
    
    // 💥 Якщо тут crash - гроші зняті, але event не записаний!
    
    // Step 2: Insert event
    _, err = db.Exec(
        "INSERT INTO outbox (event_type, user_id, amount) VALUES ('withdrawal', $1, $2)",
        userID, amount,
    )
    if err != nil {
        // ⚠️ Помилка! Гроші вже зняті, але event не записаний!
        // Дані неконсистентні!
        return err
    }
    
    return nil
}
```

### ✅ З atomic transaction (безпечно)

```go
func WithdrawMoney(userID int64, amount float64) error {
    // BEGIN TRANSACTION
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback() // Auto-rollback if not committed
    
    // Step 1: Deduct money (в рамках транзакції)
    _, err = tx.Exec(
        "UPDATE accounts SET balance = balance - $1 WHERE user_id = $2",
        amount, userID,
    )
    if err != nil {
        return err // Rollback автоматично
    }
    
    // 💥 Якщо тут crash - транзакція rollback, нічого не збережено ✅
    
    // Step 2: Insert event (в тій же транзакції)
    _, err = tx.Exec(
        "INSERT INTO outbox (event_type, user_id, amount) VALUES ('withdrawal', $1, $2)",
        userID, amount,
    )
    if err != nil {
        return err // Rollback автоматично, обидві операції скасовані ✅
    }
    
    // COMMIT - обидві операції або збережені разом, або rollback разом!
    if err = tx.Commit(); err != nil {
        return err
    }
    
    return nil // ✅ Обидві операції збережені атомарно!
}
```

---

## 🔬 Деталі: Що відбувається всередині

### Під час транзакції (до COMMIT)

```
BEGIN TRANSACTION
         |
Database створює "тимчасову версію" даних
         |
         ├─ UPDATE accounts: balance = 900 (was 1000)
         |  └─ Збережено в transaction log (не видно іншим)
         |
         ├─ INSERT INTO outbox: event created
         |  └─ Збережено в transaction log (не видно іншим)
         |
         ├─ Інші користувачі бачать balance = 1000 (старе значення)
         |  └─ Isolation: зміни не видно до COMMIT
         |
         └─ COMMIT або ROLLBACK?
```

### COMMIT (успіх)

```
COMMIT
  |
  ├─ Database застосовує всі зміни з transaction log
  |  ├─ balance = 900 ✅
  |  └─ outbox event created ✅
  |
  ├─ Зміни стають видимими іншим користувачам
  |
  └─ Неможливо відкотити (permanent)

✅ Atomic: обидві операції збережені разом!
```

### ROLLBACK (помилка)

```
ERROR або CRASH
  |
  ├─ Database скасовує всі зміни з transaction log
  |  ├─ balance залишається 1000 ✅
  |  └─ outbox event не створений ✅
  |
  ├─ Як ніби нічого не відбулося
  |
  └─ Дані консистентні

✅ Atomic: обидві операції скасовані разом!
```

---

## 🎓 ACID Properties

**Atomic** - це "A" в ACID:

### **A**tomicity (Атомарність)
```
All or nothing
├─ Або всі операції виконані ✅
└─ Або жодна не виконана ❌
```

### **C**onsistency (Консистентність)
```
Дані завжди в валідному стані
├─ Не може бути balance зменшено без outbox event
└─ Constraints завжди виконані
```

### **I**solation (Ізоляція)
```
Транзакції не заважають одна одній
├─ User A не бачить незавершені зміни User B
└─ Рівні ізоляції: Read Committed, Repeatable Read, Serializable
```

### **D**urability (Довговічність)
```
Після COMMIT дані збережені назавжди
├─ Навіть якщо сервер crash ✅
└─ Гарантія persistence
```

---

## 🔍 Приклади в життєвих ситуаціях

### Приклад 1: Переказ грошей між рахунками

```go
func TransferMoney(fromUserID, toUserID int64, amount float64) error {
    tx, _ := db.Begin()
    defer tx.Rollback()
    
    // Операція 1: Зняти з рахунку A
    tx.Exec("UPDATE accounts SET balance = balance - $1 WHERE user_id = $2", 
        amount, fromUserID)
    
    // Операція 2: Додати на рахунок B
    tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE user_id = $2", 
        amount, toUserID)
    
    // ATOMIC: або обидві виконані, або жодна!
    return tx.Commit()
}
```

**Без atomic:**
```
❌ Гроші зняті з A, але не додані на B = гроші зникли!
```

**З atomic:**
```
✅ Або обидві операції успішні, або обидві скасовані
```

---

### Приклад 2: Створення замовлення

```go
func CreateOrder(userID int64, items []Item) error {
    tx, _ := db.Begin()
    defer tx.Rollback()
    
    // Операція 1: Створити order
    var orderID int64
    tx.QueryRow(
        "INSERT INTO orders (user_id, total) VALUES ($1, $2) RETURNING id",
        userID, calculateTotal(items),
    ).Scan(&orderID)
    
    // Операція 2-N: Додати order items
    for _, item := range items {
        tx.Exec(
            "INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1, $2, $3)",
            orderID, item.ProductID, item.Quantity,
        )
    }
    
    // Операція N+1: Зменшити inventory
    for _, item := range items {
        tx.Exec(
            "UPDATE inventory SET quantity = quantity - $1 WHERE product_id = $2",
            item.Quantity, item.ProductID,
        )
    }
    
    // ATOMIC: або всі 20 операцій виконані, або жодна!
    return tx.Commit()
}
```

**Без atomic:**
```
❌ Order створений, але inventory не зменшений = overselling!
❌ Order створений, але items не додані = broken data!
```

**З atomic:**
```
✅ Або order + items + inventory update всі успішні
✅ Або якщо помилка - нічого не збережено
```

---

## ⚠️ Поширені помилки

### Помилка 1: Забув використати tx

```go
func BadExample(userID int64) error {
    tx, _ := db.Begin()
    defer tx.Rollback()
    
    // ❌ Використав db замість tx - операція поза транзакцією!
    db.Exec("UPDATE accounts SET balance = balance - 100")
    
    // ✅ Правильно
    tx.Exec("INSERT INTO outbox ...")
    
    return tx.Commit()
}
```

### Помилка 2: Забув defer tx.Rollback()

```go
func BadExample2(userID int64) error {
    tx, _ := db.Begin()
    // ❌ Немає defer tx.Rollback()
    
    tx.Exec("UPDATE accounts SET balance = balance - 100")
    
    // Якщо тут panic - транзакція залишиться відкритою!
    // Connection leak!
    
    return tx.Commit()
}

// ✅ Завжди додавай defer
func GoodExample(userID int64) error {
    tx, _ := db.Begin()
    defer tx.Rollback() // ✅ Безпечно
    
    tx.Exec("UPDATE accounts SET balance = balance - 100")
    
    return tx.Commit()
}
```

### Помилка 3: Довгі транзакції

```go
// ❌ BAD - транзакція тримається довго
func BadExample3() error {
    tx, _ := db.Begin()
    defer tx.Rollback()
    
    tx.Exec("UPDATE accounts ...")
    
    // ❌ External API call inside transaction!
    time.Sleep(5 * time.Second) // Симулює довгий запит
    resp, _ := http.Get("https://api.example.com/verify")
    
    return tx.Commit()
}

// ✅ GOOD - транзакція коротка
func GoodExample3() error {
    // Спочатку зробити external calls
    resp, _ := http.Get("https://api.example.com/verify")
    
    // Потім швидка транзакція
    tx, _ := db.Begin()
    defer tx.Rollback()
    
    tx.Exec("UPDATE accounts ...")
    
    return tx.Commit() // Швидко!
}
```

---

## 📊 Performance: Atomic vs Non-Atomic

### Non-Atomic (швидше, але небезпечно)

```
Operation 1: 10ms ✅
Operation 2: 10ms ✅
Total: 20ms

⚠️ Але якщо Op2 fails - дані неконсистентні!
```

### Atomic (трохи повільніше, але безпечно)

```
BEGIN: 1ms
Operation 1: 10ms
Operation 2: 10ms
COMMIT: 5ms
Total: 26ms (на 30% повільніше)

✅ Але consistency гарантована!
```

**Trade-off:** +20-30% часу за consistency

---

## 🎯 Коли використовувати atomic transactions

### ✅ Обов'язково використовуй

1. **Фінансові операції**
   - Переказ грошей
   - Зняття/поповнення
   - Платежі

2. **Критичні бізнес-операції**
   - Створення замовлення
   - Інвентар
   - Резервації

3. **Linked data**
   - Order + Order Items
   - User + Profile
   - Parent + Children records

### ⚠️ Можна без atomic

1. **Logs/Analytics**
   - Не критично якщо втратиться 1 запис

2. **Cache updates**
   - Eventual consistency OK

3. **Read-only операції**
   - SELECT не потребує транзакції

---

## 💡 Best Practices

### 1. Тримай транзакції короткими

```go
// ✅ GOOD - швидка транзакція
tx, _ := db.Begin()
defer tx.Rollback()

tx.Exec("UPDATE ...")
tx.Exec("INSERT ...")

tx.Commit() // < 100ms
```

### 2. Не роби I/O в транзакції

```go
// ❌ BAD
tx, _ := db.Begin()
tx.Exec("UPDATE ...")
sendEmail() // ❌ External I/O в транзакції!
tx.Commit()

// ✅ GOOD
tx, _ := db.Begin()
tx.Exec("UPDATE ...")
tx.Exec("INSERT INTO outbox ...") // Event для email
tx.Commit()
// Окремий worker відправить email
```

### 3. Використовуй правильний Isolation Level

```go
// Default: Read Committed (достатньо для більшості)
tx, _ := db.Begin()

// Для критичних операцій: Serializable
tx, _ := db.BeginTx(ctx, &sql.TxOptions{
    Isolation: sql.LevelSerializable,
})
```

---

## 🎓 Висновок

**Atomic transaction** означає:

```
┌─────────────────────────────────────┐
│  BEGIN TRANSACTION                  │
│  ├─ Operation 1                     │
│  ├─ Operation 2                     │
│  ├─ Operation N                     │
│  └─ COMMIT                          │
│     ↓                               │
│  Або ВСІ разом ✅                   │
│  Або ЖОДНА ❌                       │
│                                     │
│  Неможливо "наполовину"!            │
└─────────────────────────────────────┘
```

**Це foundation для data consistency в будь-якій системі!** 🎯

---

**Atomic = All-or-Nothing = Consistency Guaranteed!** ✅
