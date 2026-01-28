# Типи Систем Транзакцій - Повний Огляд

## 📊 Класифікація

```
Системи Транзакцій
│
├─ 1. Локальні (Single Database)
│   └─ ACID транзакції
│
├─ 2. Розподілені (Multiple Databases)
│   ├─ Two-Phase Commit (2PC)
│   ├─ Three-Phase Commit (3PC)
│   └─ XA Transactions
│
├─ 3. Eventual Consistency
│   ├─ Saga Pattern
│   ├─ Event Sourcing
│   └─ CQRS
│
├─ 4. Hybrid
│   ├─ Outbox Pattern
│   ├─ Inbox Pattern
│   └─ Transactional Outbox/Inbox
│
└─ 5. External Systems
    ├─ Reserve/Confirm Pattern
    ├─ Try-Confirm/Cancel (TCC)
    └─ Compensating Transactions
```

---

## 1️⃣ Локальні Транзакції (ACID)

### Що це?

**Одна база даних, всі операції в одній транзакції**

### Схема

```
Application
     ↓
Single Database
     ↓
BEGIN TRANSACTION
├─ Operation 1 ✅
├─ Operation 2 ✅
├─ Operation 3 ✅
└─ COMMIT ✅
```

### Приклад

```go
func TransferMoney(fromID, toID int64, amount float64) error {
    tx, _ := db.Begin()
    defer tx.Rollback()
    
    tx.Exec("UPDATE accounts SET balance = balance - $1 WHERE id = $2", amount, fromID)
    tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, toID)
    
    return tx.Commit()
}
```

### Властивості

- ✅ **A**tomicity: все або нічого
- ✅ **C**onsistency: дані валідні
- ✅ **I**solation: не заважають
- ✅ **D**urability: назавжди

### Коли використовувати

✅ Монолітні застосунки  
✅ Одна база даних  
✅ Immediate consistency потрібна  

### Обмеження

❌ Не працює з кількома БД  
❌ Не працює з зовнішніми API  
❌ Не scalable для мікросервісів  

---

## 2️⃣ Two-Phase Commit (2PC)

### Що це?

**Координатор керує транзакцією через кілька БД**

### Схема

```
         Coordinator
              |
    ┌─────────┼─────────┐
    ↓         ↓         ↓
  DB 1      DB 2      DB 3
    |         |         |
PHASE 1: PREPARE
    |         |         |
 Ready? ✅  Ready? ✅  Ready? ✅
    |         |         |
PHASE 2: COMMIT
    |         |         |
 Commit ✅  Commit ✅  Commit ✅
```

### Приклад

```go
type Coordinator struct {
    databases []*sql.DB
}

func (c *Coordinator) ExecuteTransaction() error {
    txns := make([]*sql.Tx, len(c.databases))
    
    // PHASE 1: PREPARE
    for i, db := range c.databases {
        tx, _ := db.Begin()
        txns[i] = tx
        
        // Do work in transaction
        tx.Exec("UPDATE ...")
        
        // Don't commit yet!
    }
    
    // PHASE 2: COMMIT (або ROLLBACK всіх)
    for _, tx := range txns {
        if err := tx.Commit(); err != nil {
            // Rollback all
            for _, t := range txns {
                t.Rollback()
            }
            return err
        }
    }
    
    return nil
}
```

### Властивості

✅ Strong consistency  
✅ ACID гарантії  
✅ All-or-nothing  

### Проблеми

❌ **Blocking** - якщо coordinator fails  
❌ **Single point of failure**  
❌ **Performance** overhead  
❌ **Deadlocks** можливі  

### Коли використовувати

⚠️ Рідко в сучасних системах  
⚠️ Legacy distributed databases  
⚠️ Коли ACID критичний  

---

## 3️⃣ Saga Pattern (Choreography)

### Що це?

**Послідовність локальних транзакцій з compensating actions**

### Схема

```
Service A          Service B          Service C
    |                  |                  |
1. Debit account ✅    |                  |
   └─> Event ────────> |                  |
                  2. Process ✅           |
                     └─> Event ─────────> |
                                    3. Send ✅
                                       └─> Done
```

### Якщо помилка

```
Service A          Service B          Service C
    |                  |                  |
1. Debit ✅           |                  |
    |                  |                  |
    |            2. Process ❌           |
    |               └─> Event ──────────>|
    |                  |                  |
    |<─── Compensate ─┘                 |
    |                                     |
1b. Refund ✅ (compensating)            |
```

### Приклад

```go
// Service A: Account Service
func (s *AccountService) Withdraw(userID int64, amount float64) error {
    tx, _ := s.db.Begin()
    defer tx.Rollback()
    
    // Local transaction
    tx.Exec("UPDATE accounts SET balance = balance - $1", amount)
    tx.Exec("INSERT INTO saga_log (tx_id, status) VALUES ($1, 'pending')", txID)
    
    tx.Commit()
    
    // Publish event (не в транзакції!)
    s.queue.Publish("MoneyDeducted", Event{TxID: txID, Amount: amount})
    
    return nil
}

// Compensating transaction
func (s *AccountService) HandlePaymentFailed(event Event) {
    tx, _ := s.db.Begin()
    defer tx.Rollback()
    
    // Refund
    tx.Exec("UPDATE accounts SET balance = balance + $1", event.Amount)
    tx.Exec("UPDATE saga_log SET status = 'compensated' WHERE tx_id = $1", event.TxID)
    
    tx.Commit()
}
```

### Властивості

✅ No blocking  
✅ High availability  
✅ Scalable  
⚠️ Eventual consistency  

### Проблеми

❌ Складніше реалізувати  
❌ Потрібні compensating transactions  
❌ Eventual consistency (не immediate)  

### Коли використовувати

✅ **Мікросервіси** (найпопулярніше!)  
✅ Event-driven architecture  
✅ High scalability потрібна  

---

## 4️⃣ Event Sourcing

### Що це?

**Зберігати всі зміни як події (append-only log)**

### Схема

```
Event Store (immutable log)
├─ Event 1: AccountCreated
├─ Event 2: MoneyDeposited $100
├─ Event 3: MoneyWithdrawn $50
└─ Event 4: MoneyDeposited $200

Current State = Replay всіх events
└─> Balance = 0 + 100 - 50 + 200 = $250
```

### Приклад

```go
type Event struct {
    ID        string
    Type      string
    Aggregate string
    Data      json.RawMessage
    Timestamp time.Time
}

type EventStore struct {
    db *sql.DB
}

func (es *EventStore) Append(event Event) error {
    _, err := es.db.Exec(
        "INSERT INTO events (id, type, aggregate, data, timestamp) VALUES ($1, $2, $3, $4, $5)",
        event.ID, event.Type, event.Aggregate, event.Data, event.Timestamp,
    )
    return err
}

// Rebuild state
func RebuildAccount(es *EventStore, accountID string) (*Account, error) {
    events, _ := es.GetEvents(fmt.Sprintf("account:%s", accountID))
    
    account := &Account{ID: accountID, Balance: 0}
    
    for _, event := range events {
        switch event.Type {
        case "MoneyDeposited":
            var data struct{ Amount float64 }
            json.Unmarshal(event.Data, &data)
            account.Balance += data.Amount
            
        case "MoneyWithdrawn":
            var data struct{ Amount float64 }
            json.Unmarshal(event.Data, &data)
            account.Balance -= data.Amount
        }
    }
    
    return account, nil
}
```

### Властивості

✅ Повна історія змін  
✅ Audit trail  
✅ Time travel (replay to any point)  
✅ Projections для різних use cases  

### Проблеми

❌ Складна реалізація  
❌ Query складніше  
❌ Storage overhead  
❌ Event schema evolution  

### Коли використовувати

✅ Audit trail критичний  
✅ Фінансові системи  
✅ Compliance (GDPR, HIPAA)  
✅ Undo/Redo functionality  

---

## 5️⃣ CQRS (Command Query Responsibility Segregation)

### Що це?

**Різні моделі для write і read**

### Схема

```
           Command Side              Query Side
                |                        |
        ┌───────┴───────┐                |
        ↓               ↓                ↓
   Write Model     Event Bus      Read Models
   (normalized)         |         (denormalized)
        ↓               |               ↓
   PostgreSQL           └────────> ElasticSearch
                                        ↓
                                   Redis Cache
```

### Приклад

```go
// Write Side (Command)
type CommandHandler struct {
    db        *sql.DB
    eventBus  EventBus
}

func (h *CommandHandler) CreateOrder(order *Order) error {
    tx, _ := h.db.Begin()
    defer tx.Rollback()
    
    // Write to normalized schema
    tx.Exec("INSERT INTO orders (id, user_id, total) VALUES ($1, $2, $3)",
        order.ID, order.UserID, order.Total)
    
    for _, item := range order.Items {
        tx.Exec("INSERT INTO order_items (order_id, product_id, qty) VALUES ($1, $2, $3)",
            order.ID, item.ProductID, item.Quantity)
    }
    
    tx.Commit()
    
    // Publish event
    h.eventBus.Publish("OrderCreated", order)
    
    return nil
}

// Read Side (Query)
type QueryHandler struct {
    cache *redis.Client
    es    *elasticsearch.Client
}

func (h *QueryHandler) GetOrderDetails(orderID string) (*OrderView, error) {
    // Try cache first
    if data, err := h.cache.Get(orderID).Bytes(); err == nil {
        var view OrderView
        json.Unmarshal(data, &view)
        return &view, nil
    }
    
    // Fallback to ElasticSearch (denormalized)
    view, _ := h.es.GetOrderView(orderID)
    
    // Cache it
    data, _ := json.Marshal(view)
    h.cache.Set(orderID, data, 1*time.Hour)
    
    return view, nil
}

// Event handler (updates read model)
func (h *QueryHandler) HandleOrderCreated(event OrderCreatedEvent) {
    // Update ElasticSearch
    h.es.IndexOrder(event.OrderID, OrderView{
        ID:       event.OrderID,
        UserName: event.UserName,
        Items:    event.Items,
        Total:    event.Total,
    })
}
```

### Властивості

✅ Optimize reads незалежно від writes  
✅ Different storage для read/write  
✅ Scalability  
✅ Complex queries без impact на writes  

### Проблеми

❌ Eventual consistency  
❌ Складна архітектура  
❌ Sync між read/write models  

### Коли використовувати

✅ High read/write різниця  
✅ Complex queries потрібні  
✅ Different scaling для read/write  

---

## 6️⃣ Outbox Pattern

### Що це?

**Гарантія delivery to message queue через БД транзакцію**

### Схема

```
BEGIN TRANSACTION
├─ UPDATE accounts SET balance = balance - 100
├─ INSERT INTO outbox (event_type, payload)
└─ COMMIT ✅ (обидва разом!)
         ↓
Background Worker
├─ SELECT * FROM outbox WHERE published = false
├─ Publish to message queue ✅
└─ UPDATE outbox SET published = true
```

### Приклад

```go
func WithdrawWithOutbox(userID int64, amount float64) error {
    tx, _ := db.Begin()
    defer tx.Rollback()
    
    // 1. Business logic
    tx.Exec("UPDATE accounts SET balance = balance - $1 WHERE user_id = $2", 
        amount, userID)
    
    // 2. Insert to outbox (same transaction!)
    event := Event{Type: "MoneyWithdrawn", UserID: userID, Amount: amount}
    eventJSON, _ := json.Marshal(event)
    
    tx.Exec("INSERT INTO outbox (event_type, payload, created_at) VALUES ($1, $2, NOW())",
        event.Type, eventJSON)
    
    // 3. Commit (atomic!)
    return tx.Commit()
}

// Background worker
func OutboxWorker(db *sql.DB, queue MessageQueue) {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        rows, _ := db.Query("SELECT id, event_type, payload FROM outbox WHERE published = false LIMIT 100")
        defer rows.Close()
        
        for rows.Next() {
            var id int64
            var eventType string
            var payload []byte
            rows.Scan(&id, &eventType, &payload)
            
            // Publish
            if err := queue.Publish(eventType, payload); err != nil {
                continue // Retry later
            }
            
            // Mark as published
            db.Exec("UPDATE outbox SET published = true WHERE id = $1", id)
        }
    }
}
```

### Властивості

✅ **At-least-once delivery** гарантія  
✅ No lost messages  
✅ ACID + messaging разом  

### Проблеми

❌ Eventual consistency  
❌ Потрібен background worker  
❌ Idempotency на consumer side  

### Коли використовувати

✅ **З message queues** (RabbitMQ, Kafka, SQS)  
✅ Microservices communication  
✅ Event-driven architecture  

---

## 7️⃣ Try-Confirm/Cancel (TCC)

### Що це?

**Трифазний протокол для distributed transactions**

### Схема

```
Phase 1: TRY (reserve resources)
├─ Service A: Try ✅
├─ Service B: Try ✅
└─ Service C: Try ✅
         ↓
Phase 2: CONFIRM або CANCEL
         ↓
   ┌─────┴─────┐
   ↓           ↓
CONFIRM      CANCEL
   ↓           ↓
Finalize   Compensate
```

### Приклад

```go
type TCCService interface {
    Try(ctx context.Context, txID string) error
    Confirm(ctx context.Context, txID string) error
    Cancel(ctx context.Context, txID string) error
}

// Account Service
type AccountService struct {
    db *sql.DB
}

func (s *AccountService) Try(ctx context.Context, txID string) error {
    tx, _ := s.db.Begin()
    defer tx.Rollback()
    
    // Reserve (не списувати!)
    tx.Exec("INSERT INTO account_holds (tx_id, amount, status) VALUES ($1, $2, 'try')",
        txID, amount)
    
    return tx.Commit()
}

func (s *AccountService) Confirm(ctx context.Context, txID string) error {
    tx, _ := s.db.Begin()
    defer tx.Rollback()
    
    // Finalize
    tx.Exec("UPDATE accounts SET balance = balance - (SELECT amount FROM account_holds WHERE tx_id = $1)", txID)
    tx.Exec("UPDATE account_holds SET status = 'confirmed' WHERE tx_id = $1", txID)
    
    return tx.Commit()
}

func (s *AccountService) Cancel(ctx context.Context, txID string) error {
    tx, _ := s.db.Begin()
    defer tx.Rollback()
    
    // Compensate
    tx.Exec("UPDATE account_holds SET status = 'cancelled' WHERE tx_id = $1", txID)
    
    return tx.Commit()
}

// Coordinator
func ExecuteTCC(services []TCCService) error {
    txID := generateTxID()
    
    // Phase 1: TRY
    for _, svc := range services {
        if err := svc.Try(context.Background(), txID); err != nil {
            // Cancel all
            for _, s := range services {
                s.Cancel(context.Background(), txID)
            }
            return err
        }
    }
    
    // Phase 2: CONFIRM
    for _, svc := range services {
        if err := svc.Confirm(context.Background(), txID); err != nil {
            // Partial failure - need manual intervention
            log.Error("TCC confirm failed", err)
            return err
        }
    }
    
    return nil
}
```

### Властивості

✅ No blocking (порівняно з 2PC)  
✅ Explicit reserve/confirm  
✅ Compensating transactions  

### Проблеми

❌ Складна реалізація  
❌ Потребує 3 endpoints на service  
❌ Partial failures складні  

### Коли використовувати

⚠️ Фінансові системи  
⚠️ Booking systems (hotels, flights)  
⚠️ Коли reserve/confirm pattern природний  

---

## 8️⃣ External Systems (Reserve/Confirm)

### Що це?

**Для систем поза контролем БД (ATM, payment gateway)**

### Схема

```
1. RESERVE (в БД)
├─ balance: $1000 (не змінюється)
└─ available: $900 (hold $100)
         ↓
2. TRY EXTERNAL (ATM, API)
         ↓
   ┌─────┴─────┐
   ↓           ↓
SUCCESS      FAILURE
   ↓           ↓
3a. CONFIRM  3b. REFUND
(deduct)     (release hold)
```

### Приклад (детально в файлі 15)

```go
func WithdrawCash(userID int64, amount float64) error {
    // 1. Reserve
    txnID, err := reserveMoney(userID, amount)
    if err != nil {
        return err
    }
    
    // 2. Try external
    success, err := atmClient.DispenseCash(txnID, amount)
    if err != nil {
        // 3b. Refund
        refundMoney(txnID)
        return err
    }
    
    // 3a. Confirm
    return confirmWithdrawal(txnID)
}
```

### Властивості

✅ Safe для external systems  
✅ No money lost  
✅ Reconciliation можливий  

### Коли використовувати

✅ ATM transactions  
✅ Payment gateways  
✅ Shipping APIs  
✅ Будь-які external hardware/services  

---

## 📊 Порівняльна таблиця

| Pattern | Consistency | Complexity | Performance | Use Case |
|---------|-------------|------------|-------------|----------|
| **ACID** | Immediate ✅ | Low ⭐ | High ✅ | Single DB |
| **2PC** | Immediate ✅ | High ⭐⭐⭐ | Low ❌ | Legacy distributed |
| **Saga** | Eventual ⚠️ | Medium ⭐⭐ | High ✅ | **Microservices** ⭐ |
| **Event Sourcing** | Eventual ⚠️ | High ⭐⭐⭐ | Medium | Audit trail |
| **CQRS** | Eventual ⚠️ | High ⭐⭐⭐ | High ✅ | High read/write diff |
| **Outbox** | Eventual ⚠️ | Medium ⭐⭐ | High ✅ | **With queues** ⭐ |
| **TCC** | Immediate ✅ | High ⭐⭐⭐ | Medium | Financial systems |
| **Reserve/Confirm** | Eventual ⚠️ | Medium ⭐⭐ | Medium | **External systems** ⭐ |

---

## 🎯 Вибір правильного підходу

### Flowchart

```
START: Потрібна транзакція
         ↓
    Одна БД?
    ├─ YES → ACID ✅
    └─ NO ↓
         ↓
    Кілька БД в одній компанії?
    ├─ YES → 2PC (рідко) або Saga
    └─ NO ↓
         ↓
    Мікросервіси?
    ├─ YES → Saga + Outbox ⭐
    └─ NO ↓
         ↓
    Зовнішні системи?
    ├─ YES → Reserve/Confirm ⭐
    └─ NO ↓
         ↓
    Audit trail критичний?
    ├─ YES → Event Sourcing
    └─ NO ↓
         ↓
    High read/write різниця?
    └─ YES → CQRS
```

---

## 💡 Рекомендації по архітектурі

### Монолітна система

```
✅ ACID транзакції
✅ Проста, надійна
✅ Immediate consistency
```

### Мікросервіси (найпопулярніше)

```
✅ Saga Pattern (choreography)
✅ Outbox Pattern (з Kafka/RabbitMQ)
✅ Event-driven architecture
⚠️ Eventual consistency
```

### Фінансові системи

```
✅ Event Sourcing (audit trail)
✅ CQRS (read/write optimization)
✅ TCC (для critical operations)
✅ Saga (для workflow)
```

### E-commerce

```
✅ Saga Pattern (order → payment → shipping)
✅ Outbox Pattern (notifications)
✅ Reserve/Confirm (payment gateway, warehouse)
✅ CQRS (product catalog)
```

---

## 🎓 Висновок

**Не існує "найкращого" підходу - вибір залежить від:**

1. **Архітектури** (монолітна / мікросервіси)
2. **Consistency requirements** (immediate / eventual)
3. **Performance needs** (latency / throughput)
4. **Business domain** (фінанси / e-commerce / social)
5. **Team expertise** (складність реалізації)

**Найпопулярніші в 2026:**

1. **Saga Pattern** ⭐⭐⭐⭐⭐ (мікросервіси)
2. **Outbox Pattern** ⭐⭐⭐⭐⭐ (з message queues)
3. **ACID** ⭐⭐⭐⭐ (монолітні системи)
4. **Reserve/Confirm** ⭐⭐⭐⭐ (external systems)
5. **Event Sourcing** ⭐⭐⭐ (audit trail)

**Рідко використовуються:**

- 2PC (legacy, blocking problems)
- 3PC (складно, рідко підтримується)

---

**Розумійте всі підходи, вибирайте правильний для вашого use case!** 🎯
