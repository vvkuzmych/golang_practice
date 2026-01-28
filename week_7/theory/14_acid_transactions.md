# ACID Транзакції - Просто про складне

## 🎯 Що таке ACID?

**ACID** - це 4 властивості, які гарантують надійність транзакцій в базі даних:

```
A - Atomicity      (Атомарність)
C - Consistency    (Консистентність)
I - Isolation      (Ізоляція)
D - Durability     (Довговічність)
```

---

## 📖 Проста ACID транзакція - Приклад

### Сценарій: Переказ $100 з рахунку A на рахунок B

```go
func TransferMoney(fromID, toID int64, amount float64) error {
    // BEGIN TRANSACTION
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    // 1. Зняти з рахунку A
    _, err = tx.Exec(
        "UPDATE accounts SET balance = balance - $1 WHERE id = $2",
        amount, fromID,
    )
    if err != nil {
        return err // Auto-rollback
    }
    
    // 2. Додати на рахунок B
    _, err = tx.Exec(
        "UPDATE accounts SET balance = balance + $1 WHERE id = $2",
        amount, toID,
    )
    if err != nil {
        return err // Auto-rollback
    }
    
    // COMMIT TRANSACTION
    return tx.Commit()
}
```

Це **проста ACID транзакція** тому що:
- ✅ Одна база даних
- ✅ Всі операції в одній транзакції
- ✅ BEGIN → операції → COMMIT
- ✅ Гарантії ACID

---

## 1️⃣ A - Atomicity (Атомарність)

### Що це?

**All-or-nothing** - або всі операції виконані, або жодна.

### Візуалізація

```
BEGIN TRANSACTION
├─ Зняти $100 з рахунку A
├─ Додати $100 на рахунок B
└─ COMMIT
         ↓
   ┌─────┴─────┐
   ↓           ↓
SUCCESS      ERROR
   ↓           ↓
Обидві ✅    Обидві ❌
виконані     скасовані
```

### Приклад

```go
tx, _ := db.Begin()
defer tx.Rollback()

// Операція 1
tx.Exec("UPDATE accounts SET balance = balance - 100 WHERE id = 1")

// 💥 Якщо тут помилка
tx.Exec("UPDATE accounts SET balance = balance + 100 WHERE id = 2")

// Обидві операції або виконані ✅ або скасовані ❌
tx.Commit()
```

### ✅ Що гарантує?

- Гроші не зникнуть
- Гроші не подвояться
- Дані завжди консистентні

---

## 2️⃣ C - Consistency (Консистентність)

### Що це?

База даних **завжди в валідному стані**. Всі constraints і rules виконані.

### Візуалізація

```
Стан ДО транзакції:
├─ Account A: balance = $1000 ✅
├─ Account B: balance = $500  ✅
└─ Total: $1500 ✅

    ↓ TRANSACTION ↓

Стан ПІСЛЯ транзакції:
├─ Account A: balance = $900 ✅
├─ Account B: balance = $600 ✅
└─ Total: $1500 ✅ (не змінився!)

Consistency rules виконані! ✅
```

### Приклад: Database Constraints

```sql
-- Constraint: balance не може бути негативним
CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    balance DECIMAL NOT NULL CHECK (balance >= 0)
);

-- Спроба зняти більше ніж є
BEGIN;
UPDATE accounts SET balance = balance - 1500 WHERE id = 1;
-- ❌ ERROR: violates check constraint "balance >= 0"
ROLLBACK; -- Транзакція скасована
```

### ✅ Що гарантує?

- Всі constraints виконані
- Business rules не порушені
- Дані валідні до і після транзакції

---

## 3️⃣ I - Isolation (Ізоляція)

### Що це?

**Транзакції не заважають одна одній**. Кожна транзакція виконується ніби вона єдина.

### Візуалізація

```
User A                        User B
  |                             |
BEGIN TRANSACTION          BEGIN TRANSACTION
  |                             |
Read: balance = $1000      Read: balance = $1000
  |                             |
Withdraw $100              Withdraw $50
  |                             |
Write: balance = $900      Write: balance = $950
  |                             |
COMMIT ✅                  COMMIT ✅
  |                             |
  
Final balance: $850 ✅ (обидві транзакції враховані)

NOT: $900 ❌ (втрачена транзакція User B)
NOT: $950 ❌ (втрачена транзакція User A)
```

### Приклад: Lost Update Problem

#### ❌ БЕЗ Isolation

```go
// User A
balance := getBalance(1) // $1000
balance -= 100            // $900
updateBalance(1, balance)

// User B (одночасно)
balance := getBalance(1) // $1000 (стара!)
balance -= 50             // $950
updateBalance(1, balance) // ❌ Перезаписав A!

// Result: $950 (втрачено $100 від User A!)
```

#### ✅ З Isolation

```go
// User A
tx, _ := db.Begin()
tx.Exec("UPDATE accounts SET balance = balance - 100 WHERE id = 1")
// Transaction locks row
tx.Commit()

// User B
tx, _ := db.Begin()
tx.Exec("UPDATE accounts SET balance = balance - 50 WHERE id = 1")
// Waits for User A to finish
tx.Commit()

// Result: $850 ✅ (обидві враховані)
```

### Isolation Levels

```
Serializable       🔒🔒🔒 Найсуворіше (найповільніше)
    ↑
Repeatable Read    🔒🔒   Середнє
    ↑
Read Committed     🔒     Default (найшвидше)
    ↑
Read Uncommitted   (без isolation)
```

### ✅ Що гарантує?

- Concurrent transactions не конфліктують
- Кожна transaction бачить консистентний snapshot
- Немає "lost updates"

---

## 4️⃣ D - Durability (Довговічність)

### Що це?

**Після COMMIT дані збережені назавжди**, навіть якщо сервер crash.

### Візуалізація

```
COMMIT ✅
  |
  ├─> Дані записані на диск (permanent storage)
  |
  ├─> Write-Ahead Log (WAL)
  |
  └─> 💥 Server crash
      └─> Restart
          └─> Дані все ще там! ✅
```

### Приклад

```go
tx, _ := db.Begin()
tx.Exec("UPDATE accounts SET balance = 900")
tx.Commit() // ✅ COMMITTED

// 💥 Server crashes after COMMIT

// После перезагрузки:
balance := getBalance(1) // $900 ✅ (збережено!)
```

### Як це працює?

```
Write-Ahead Logging (WAL):

1. Transaction модифікує дані
2. Зміни спочатку пишуться в WAL (log file)
3. WAL flush на диск
4. COMMIT
5. Пізніше дані записуються в основні файли

Якщо crash після COMMIT:
└─> Recovery з WAL ✅
```

### ✅ Що гарантує?

- COMMIT = permanent
- Crash не втрачає committed data
- Write-Ahead Log захищає дані

---

## 🔄 Життєвий цикл транзакції

```
┌─────────────────────────────────────────┐
│  1. BEGIN TRANSACTION                   │
│     ├─ Створюється transaction log      │
│     └─ Locks можуть бути взяті          │
├─────────────────────────────────────────┤
│  2. OPERATIONS (Read/Write)             │
│     ├─ Зміни в transaction buffer       │
│     ├─ Visibility: тільки для цієї tx   │
│     └─ Constraints перевіряються        │
├─────────────────────────────────────────┤
│  3. COMMIT або ROLLBACK                 │
│     │                                   │
│     ├─> COMMIT:                         │
│     │   ├─ Apply changes to DB          │
│     │   ├─ Release locks                │
│     │   ├─ WAL flush                    │
│     │   └─ Changes visible to all ✅    │
│     │                                   │
│     └─> ROLLBACK:                       │
│         ├─ Discard changes              │
│         ├─ Release locks                │
│         └─ No changes applied ❌        │
└─────────────────────────────────────────┘
```

---

## 💻 Повний приклад з усіма ACID властивостями

```go
package main

import (
    "database/sql"
    "fmt"
)

// Переказ грошей з демонстрацією ACID
func TransferMoneyWithACID(db *sql.DB, fromID, toID int64, amount float64) error {
    // BEGIN TRANSACTION
    tx, err := db.Begin()
    if err != nil {
        return fmt.Errorf("begin transaction: %w", err)
    }
    
    // ATOMICITY: або всі операції, або жодна
    defer func() {
        if err != nil {
            tx.Rollback() // Скасувати при помилці
        }
    }()
    
    // ISOLATION: використовуємо FOR UPDATE для lock
    var fromBalance float64
    err = tx.QueryRow(
        "SELECT balance FROM accounts WHERE id = $1 FOR UPDATE",
        fromID,
    ).Scan(&fromBalance)
    if err != nil {
        return fmt.Errorf("get from balance: %w", err)
    }
    
    // CONSISTENCY: перевірка business rule
    if fromBalance < amount {
        return fmt.Errorf("insufficient funds: have %.2f, need %.2f", 
            fromBalance, amount)
    }
    
    var toBalance float64
    err = tx.QueryRow(
        "SELECT balance FROM accounts WHERE id = $1 FOR UPDATE",
        toID,
    ).Scan(&toBalance)
    if err != nil {
        return fmt.Errorf("get to balance: %w", err)
    }
    
    // Операція 1: зняти гроші
    _, err = tx.Exec(
        "UPDATE accounts SET balance = balance - $1 WHERE id = $2",
        amount, fromID,
    )
    if err != nil {
        return fmt.Errorf("deduct from account: %w", err)
    }
    
    // Операція 2: додати гроші
    _, err = tx.Exec(
        "UPDATE accounts SET balance = balance + $1 WHERE id = $2",
        amount, toID,
    )
    if err != nil {
        return fmt.Errorf("add to account: %w", err)
    }
    
    // Операція 3: створити audit log
    _, err = tx.Exec(
        "INSERT INTO transfers (from_id, to_id, amount, created_at) VALUES ($1, $2, $3, NOW())",
        fromID, toID, amount,
    )
    if err != nil {
        return fmt.Errorf("create audit log: %w", err)
    }
    
    // COMMIT TRANSACTION
    // DURABILITY: після commit дані збережені назавжди
    if err = tx.Commit(); err != nil {
        return fmt.Errorf("commit transaction: %w", err)
    }
    
    fmt.Printf("✅ Transferred %.2f from account %d to %d\n", amount, fromID, toID)
    return nil
}

func main() {
    db, _ := sql.Open("postgres", "...")
    defer db.Close()
    
    // Переказ $100
    err := TransferMoneyWithACID(db, 1, 2, 100.0)
    if err != nil {
        fmt.Printf("❌ Transfer failed: %v\n", err)
        return
    }
    
    fmt.Println("✅ Transfer successful with ACID guarantees!")
}
```

---

## 📊 ACID vs Non-ACID

### Non-ACID (небезпечно)

```go
// ❌ БЕЗ транзакції
db.Exec("UPDATE accounts SET balance = balance - 100 WHERE id = 1")
// 💥 CRASH тут
db.Exec("UPDATE accounts SET balance = balance + 100 WHERE id = 2")

Problems:
❌ Atomicity: можливо тільки перша виконана
❌ Consistency: total balance неправильний
❌ Isolation: інші бачать проміжний стан
❌ Durability: може втратитись при crash
```

### ACID (безпечно)

```go
// ✅ З транзакцією
tx, _ := db.Begin()
tx.Exec("UPDATE accounts SET balance = balance - 100 WHERE id = 1")
tx.Exec("UPDATE accounts SET balance = balance + 100 WHERE id = 2")
tx.Commit()

Guarantees:
✅ Atomicity: обидві або жодна
✅ Consistency: total balance правильний
✅ Isolation: інші не бачать проміжний стан
✅ Durability: після commit збережено
```

---

## ⚙️ Як база даних реалізує ACID?

### Atomicity

```
Transaction Log:
├─ BEGIN TX 123
├─ UPDATE accounts SET balance = 900 WHERE id = 1
├─ UPDATE accounts SET balance = 600 WHERE id = 2
└─ COMMIT TX 123 ✅

Якщо ROLLBACK:
└─ Undo всі операції з log
```

### Consistency

```
Constraints:
├─ CHECK (balance >= 0)
├─ FOREIGN KEY (user_id) REFERENCES users(id)
└─ Перевіряються перед COMMIT

Якщо порушення:
└─> ROLLBACK автоматично
```

### Isolation

```
Locks:
├─ Shared Lock (read)
├─ Exclusive Lock (write)
└─ Row-level / Table-level

MVCC (Multi-Version Concurrency Control):
├─ Кожна transaction бачить snapshot
└─ Різні версії рядка для різних transactions
```

### Durability

```
Write-Ahead Log (WAL):
├─ Зміни спочатку в WAL
├─ fsync() - flush to disk
├─ COMMIT
└─ Recovery з WAL після crash

Checkpoints:
└─ Періодично WAL → data files
```

---

## 🎯 Коли використовувати ACID транзакції?

### ✅ Обов'язково

1. **Фінансові операції**
   ```
   ├─ Перекази грошей
   ├─ Платежі
   └─ Зняття/поповнення
   ```

2. **E-commerce**
   ```
   ├─ Створення замовлення
   ├─ Оновлення інвентаря
   └─ Резервації
   ```

3. **Критичні бізнес-операції**
   ```
   ├─ Реєстрація користувача
   ├─ Зміна прав доступу
   └─ Аудит логи
   ```

### ⚠️ Необов'язково

1. **Логи і аналітика**
   - Eventual consistency OK
   - Performance важливіше

2. **Cache updates**
   - Можна втратити без наслідків

3. **Read-only запити**
   - Не потребують транзакцій

---

## 💡 Best Practices

### 1. Тримай транзакції короткими

```go
// ✅ GOOD - швидка транзакція
tx, _ := db.Begin()
tx.Exec("UPDATE ...")
tx.Exec("INSERT ...")
tx.Commit() // < 100ms

// ❌ BAD - довга транзакція
tx, _ := db.Begin()
tx.Exec("UPDATE ...")
time.Sleep(5 * time.Second) // ❌ Lock тримається!
tx.Commit()
```

### 2. Використовуй defer для безпеки

```go
tx, _ := db.Begin()
defer tx.Rollback() // ✅ Auto-rollback якщо panic

// Your code...

tx.Commit() // Явний commit якщо все OK
```

### 3. Правильний Isolation Level

```go
// Default: Read Committed (достатньо для 90%)
tx, _ := db.Begin()

// Для критичних операцій: Serializable
tx, _ := db.BeginTx(ctx, &sql.TxOptions{
    Isolation: sql.LevelSerializable,
})
```

### 4. Обробляй помилки

```go
tx, err := db.Begin()
if err != nil {
    return fmt.Errorf("begin: %w", err)
}
defer tx.Rollback()

if err := tx.Exec(...); err != nil {
    return fmt.Errorf("update: %w", err) // Auto-rollback
}

if err := tx.Commit(); err != nil {
    return fmt.Errorf("commit: %w", err)
}
```

---

## 📚 Порівняння баз даних

| База даних  | ACID Support | Notes                        |
|-------------|--------------|------------------------------|
| PostgreSQL  | ✅ Full      | MVCC, Serializable isolation |
| MySQL       | ✅ Full      | InnoDB engine                |
| SQLite      | ✅ Full      | Single-file DB               |
| MongoDB     | ✅ Partial   | ACID з версії 4.0            |
| Redis       | ⚠️ Limited   | MULTI/EXEC (pseudo-tx)       |
| Cassandra   | ❌ No        | Eventual consistency         |

---

## 🎓 Висновок

**Проста ACID транзакція** - це:

```
┌──────────────────────────────────────┐
│  BEGIN TRANSACTION                   │
│                                      │
│  A - Atomicity:                      │
│      Або всі операції, або жодна     │
│                                      │
│  C - Consistency:                    │
│      Дані завжди валідні             │
│                                      │
│  I - Isolation:                      │
│      Транзакції не заважають         │
│                                      │
│  D - Durability:                     │
│      COMMIT = назавжди               │
│                                      │
│  COMMIT або ROLLBACK                 │
└──────────────────────────────────────┘
```

### Чому "проста"?

- ✅ Одна база даних
- ✅ Всі операції в одній транзакції
- ✅ BEGIN → операції → COMMIT
- ✅ Гарантії ACID from DB

### Коли складніше?

- ⚠️ Розподілені транзакції (кілька БД)
- ⚠️ Мікросервіси (Saga Pattern)
- ⚠️ Message queues (Outbox Pattern)

---

**ACID = Foundation для надійних додатків!** ✅🎯

**Запам'ятай:** Завжди використовуй транзакції для критичних операцій!
