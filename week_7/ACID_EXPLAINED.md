# 💡 Проста ACID Транзакція

## 🎯 Що таке ACID?

**ACID** - 4 властивості надійних транзакцій:

```
A - Atomicity      (Атомарність)     → All-or-nothing
C - Consistency    (Консистентність) → Дані валідні
I - Isolation      (Ізоляція)        → Не заважають
D - Durability     (Довговічність)   → Назавжди
```

---

## 🏦 Приклад: Переказ $100

```go
tx, _ := db.Begin()
defer tx.Rollback()

// 1. Зняти з рахунку A
tx.Exec("UPDATE accounts SET balance = balance - 100 WHERE id = 1")

// 2. Додати на рахунок B
tx.Exec("UPDATE accounts SET balance = balance + 100 WHERE id = 2")

tx.Commit()
```

Це **проста ACID транзакція** ✅

---

## 1️⃣ A - Atomicity (Атомарність)

**All-or-nothing**

```
BEGIN TRANSACTION
├─ Зняти $100 з A
├─ Додати $100 на B
└─ COMMIT
         ↓
   ┌─────┴─────┐
   ↓           ↓
SUCCESS      ERROR
   ↓           ↓
Обидві ✅    Обидві ❌
```

✅ **Гарантує:** Гроші не зникнуть і не подвояться

---

## 2️⃣ C - Consistency (Консистентність)

**Дані завжди валідні**

```
ДО:
├─ Account A: $1000
├─ Account B: $500
└─ Total: $1500 ✅

ПІСЛЯ:
├─ Account A: $900
├─ Account B: $600
└─ Total: $1500 ✅ (не змінився!)
```

✅ **Гарантує:** Всі constraints виконані

---

## 3️⃣ I - Isolation (Ізоляція)

**Транзакції не заважають**

```
User A: Withdraw $100
User B: Withdraw $50 (одночасно)
         ↓
Final: $850 ✅

NOT: $900 ❌ (втрата транзакції B)
NOT: $950 ❌ (втрата транзакції A)
```

✅ **Гарантує:** Кожна transaction бачить консистентний snapshot

---

## 4️⃣ D - Durability (Довговічність)

**COMMIT = назавжди**

```
COMMIT ✅
  |
  └─> Дані на диск
      └─> 💥 Server crash
          └─> Restart
              └─> Дані все ще там! ✅
```

✅ **Гарантує:** Write-Ahead Log захищає дані

---

## 💻 Повний приклад

```go
func TransferMoney(fromID, toID int64, amount float64) error {
    // BEGIN
    tx, _ := db.Begin()
    defer tx.Rollback() // Auto-rollback on error
    
    // ATOMICITY: обидві або жодна
    tx.Exec("UPDATE accounts SET balance = balance - $1 WHERE id = $2", 
        amount, fromID)
    tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", 
        amount, toID)
    
    // CONSISTENCY: перевірка constraints
    // ISOLATION: locks на рядки
    // DURABILITY: WAL flush
    
    // COMMIT
    return tx.Commit()
}
```

---

## 📊 ACID vs Non-ACID

### ❌ БЕЗ транзакції

```go
db.Exec("UPDATE accounts SET balance = balance - 100 WHERE id = 1")
// 💥 CRASH
db.Exec("UPDATE accounts SET balance = balance + 100 WHERE id = 2")

Problems:
❌ Гроші зникли
❌ Total неправильний
❌ Інші бачать проміжний стан
```

### ✅ З транзакцією

```go
tx, _ := db.Begin()
tx.Exec("UPDATE accounts SET balance = balance - 100 WHERE id = 1")
tx.Exec("UPDATE accounts SET balance = balance + 100 WHERE id = 2")
tx.Commit()

Guarantees:
✅ Atomicity
✅ Consistency
✅ Isolation
✅ Durability
```

---

## 🔬 Як це працює?

### Transaction Log

```
BEGIN TX 123
├─ UPDATE accounts: balance = 900
├─ UPDATE accounts: balance = 600
└─ COMMIT TX 123 ✅

Якщо помилка:
└─> ROLLBACK (undo всі операції)
```

### Write-Ahead Log (WAL)

```
1. Зміни → WAL (log file)
2. WAL → flush to disk
3. COMMIT ✅
4. Пізніше → data files

Crash після COMMIT?
└─> Recovery з WAL ✅
```

---

## 🎯 Коли використовувати?

### ✅ Обов'язково

- Фінансові операції
- E-commerce (замовлення, інвентар)
- Критичні бізнес-операції

### ⚠️ Необов'язково

- Логи і аналітика
- Cache updates
- Read-only запити

---

## 💡 Best Practices

### 1. Короткі транзакції

```go
// ✅ GOOD
tx, _ := db.Begin()
tx.Exec("UPDATE ...")
tx.Commit() // < 100ms

// ❌ BAD
tx, _ := db.Begin()
time.Sleep(5 * time.Second) // ❌
tx.Commit()
```

### 2. defer Rollback

```go
tx, _ := db.Begin()
defer tx.Rollback() // ✅ Завжди!

// Your code...

tx.Commit()
```

### 3. Обробляй помилки

```go
if err := tx.Exec(...); err != nil {
    return err // Auto-rollback
}

if err := tx.Commit(); err != nil {
    return err
}
```

---

## 🎓 Чому "проста"?

**Проста ACID транзакція:**

```
✅ Одна база даних
✅ Всі операції в одній транзакції
✅ BEGIN → операції → COMMIT
✅ ACID гарантії from DB
```

**Складніше (розподілені системи):**

```
⚠️ Кілька баз даних
⚠️ Мікросервіси (Saga Pattern)
⚠️ Message queues (Outbox Pattern)
```

---

## 📖 Читати повний файл

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_7

cat theory/14_acid_transactions.md
```

---

## 🎯 Key Takeaways

1. **ACID = 4 гарантії**
   - Atomicity (all-or-nothing)
   - Consistency (дані валідні)
   - Isolation (не заважають)
   - Durability (назавжди)

2. **BEGIN → операції → COMMIT**
   - Або всі ✅ або жодна ❌

3. **Завжди використовуй для критичних операцій**
   - Фінанси
   - Замовлення
   - Інвентар

4. **Тримай транзакції короткими**
   - < 100ms ideal
   - Не роби I/O всередині

---

**ACID = Foundation для надійних систем!** ✅🎯

**Файл:** `theory/14_acid_transactions.md`  
**Обсяг:** Детальні пояснення всіх 4 властивостей + код
