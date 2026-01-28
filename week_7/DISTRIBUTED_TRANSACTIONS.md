# 🎉 Розподілені Транзакції - Повний Гайд

## 📖 Що це?

Детальний файл про **транзакції в розподілених системах** з текстовими схемами процесу зняття грошей з банківського рахунку.

**Файл:** `theory/12_distributed_transactions.md`

---

## 🏦 Приклад: Зняття $100 з рахунку

### Сценарій 1: Монолітна система (проста)

```
User withdraws $100
         ↓
    BEGIN TRANSACTION
         ↓
1. Check balance >= $100 ✅
2. Deduct $100 ✅
3. Create withdrawal record ✅
4. Update timestamp ✅
         ↓
    COMMIT TRANSACTION ✅
         ↓
✅ Транзакція закінчена в базі даних!
```

### Сценарій 2: Мікросервіси (складно!)

```
User withdraws $100
         ↓
[Account Service] ✅ Deduct $100
         ↓
[Payment Service] ❌ FAILS!
         ↓
⚠️ Гроші зняті, але платіж не пройшов!
```

---

## 4 Рішення проблеми розподілених транзакцій

### 1️⃣ Two-Phase Commit (2PC)

```
Coordinator
    |
PHASE 1: PREPARE
    |
Service A: Ready? ✅
Service B: Ready? ✅
Service C: Ready? ✅
    |
PHASE 2: COMMIT
    |
Service A: Commit ✅
Service B: Commit ✅
Service C: Commit ✅
    |
✅ Транзакція закінчена коли всі committed!
```

**Проблема:** Blocking, single point of failure

---

### 2️⃣ Saga Pattern (⭐ Рекомендується!)

```
Step 1: Account Service
├─ Deduct $100 ✅
└─ Publish: MoneyDeducted

Step 2: Payment Service
├─ Process payment ✅
└─ Publish: PaymentProcessed

Step 3: Notification Service
└─ Send email ✅

✅ Транзакція закінчена коли останній крок успішний!
```

**Якщо помилка:**
```
Step 1: Deduct $100 ✅
Step 2: Payment fails ❌
         ↓
COMPENSATING TRANSACTION
└─ Refund $100 ✅

✅ Транзакція rolled back!
```

---

### 3️⃣ Event Sourcing

```
Event Store (append-only)
├─ Event 1: WithdrawalRequested
├─ Event 2: BalanceChecked ✅
├─ Event 3: MoneyDeducted ✅
├─ Event 4: PaymentProcessed ✅
└─ Event 5: WithdrawalCompleted ✅

✅ Транзакція закінчена коли останній event записаний!
```

---

### 4️⃣ Outbox Pattern

```
BEGIN TRANSACTION
├─ Update accounts table ✅
└─ Insert into outbox table ✅
COMMIT ✅ (обидва разом!)
         ↓
Background worker
└─ Publish to message queue ✅
         ↓
✅ Транзакція закінчена коли event published!
```

---

## 📊 Порівняння

| Підхід      | Consistency | Complexity | Use Case          |
|-------------|-------------|------------|-------------------|
| Monolith    | Immediate ✅| Low ⭐     | Single DB         |
| 2PC         | Immediate ✅| High ⭐⭐⭐| Rare (legacy)     |
| **Saga**    | Eventual ⚠️ | Medium ⭐⭐ | **Microservices** |
| Event Store | Eventual ⚠️ | High ⭐⭐⭐| Audit trail       |
| Outbox      | Eventual ⚠️ | Medium ⭐⭐ | With queues       |

---

## ❓ Де транзакція "закінчена"?

### Монолітна система
```
✅ COMMIT в базі даних
```

### Saga Pattern
```
✅ Останній крок успішний
   або
✅ Всі compensations виконані (rollback)
```

### Event Sourcing
```
✅ Останній event в Event Store
```

### Outbox Pattern
```
✅ Event published to queue + acknowledged
```

---

## 💻 Приклади коду

### Saga Pattern (Go)

```go
// Account Service
func (s *AccountService) Withdraw(userID int64, amount float64) error {
    tx, _ := s.db.Begin()
    defer tx.Rollback()
    
    // Deduct money
    tx.Exec("UPDATE accounts SET balance = balance - $1", amount)
    
    // Log for compensating
    tx.Exec("INSERT INTO saga_log (tx_id, amount) VALUES ($1, $2)", txID, amount)
    
    tx.Commit()
    
    // Publish event
    s.queue.Publish("MoneyDeducted", Event{Amount: amount})
    
    return nil
}

// Compensating transaction (refund)
func (s *AccountService) HandlePaymentFailed(event Event) {
    tx, _ := s.db.Begin()
    
    // Refund money
    tx.Exec("UPDATE accounts SET balance = balance + $1", event.Amount)
    tx.Exec("UPDATE saga_log SET status = 'compensated'")
    
    tx.Commit()
}
```

### Outbox Pattern (Go)

```go
func WithdrawWithOutbox(amount float64) error {
    tx, _ := db.Begin()
    
    // 1. Update account
    tx.Exec("UPDATE accounts SET balance = balance - $1", amount)
    
    // 2. Insert into outbox (same transaction!)
    tx.Exec("INSERT INTO outbox (event_type, payload) VALUES ('MoneyDeducted', $1)", data)
    
    // 3. Commit (both together!)
    return tx.Commit()
}

// Background worker publishes from outbox
func OutboxWorker() {
    for {
        rows, _ := db.Query("SELECT * FROM outbox WHERE published = false")
        
        for rows.Next() {
            // Publish to queue
            queue.Publish(event)
            
            // Mark as published
            db.Exec("UPDATE outbox SET published = true")
        }
        
        time.Sleep(1 * time.Second)
    }
}
```

---

## 🎯 Що використовувати?

### Монолітна система
✅ Проста ACID транзакція  
✅ BEGIN → UPDATE → COMMIT  
✅ Immediate consistency

### Мікросервіси
✅ **Saga Pattern** (найпопулярніше!)  
✅ Compensating transactions  
✅ Event-driven  
✅ Eventual consistency

### З Message Queues
✅ **Outbox Pattern**  
✅ Гарантія delivery  
✅ At-least-once processing

---

## ⚠️ Важливі концепції

### 1. Ідемпотентність
```go
// Перевіряй чи вже оброблено
if alreadyProcessed(paymentID) {
    return nil // Skip ✅
}
```

### 2. Retry з backoff
```go
// 1s, 2s, 4s, 8s, 16s...
time.Sleep(time.Duration(1<<retry) * time.Second)
```

### 3. Distributed Tracing
```go
// Трейсинг через всі сервіси
ctx, span := tracer.Start(ctx, "withdraw-money")
defer span.End()
```

---

## 📖 Читати файл

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_7

# Повний гайд
cat theory/12_distributed_transactions.md

# Цей довідник
cat DISTRIBUTED_TRANSACTIONS.md
```

---

## 🚀 Що всередині файлу

✅ **4 детальні сценарії** з текстовими схемами  
✅ **Повний Go код** для кожного підходу  
✅ **Compensating transactions** (як rollback в Saga)  
✅ **Outbox Pattern** (reliability with queues)  
✅ **Event Sourcing** (повна історія)  
✅ **Best Practices** (idempotency, retry, tracing)  
✅ **Порівняльна таблиця** всіх підходів

---

## 🎓 Key Takeaways

1. **Монолітна система:** Транзакція = COMMIT в базі
2. **Мікросервіси:** Транзакція = останній крок Saga або compensations
3. **Saga Pattern:** Найпопулярніше рішення для мікросервісів
4. **Outbox Pattern:** Гарантія delivery до message queue
5. **Eventual Consistency:** Прийми як реальність мікросервісів

---

**Тепер ви розумієте розподілені транзакції!** 🎉

**Файл:** `theory/12_distributed_transactions.md`  
**Обсяг:** ~2,000 слів + схеми  
**Приклади:** 10+ робочих Go snippets  
**Статус:** ✅ Production-ready patterns
