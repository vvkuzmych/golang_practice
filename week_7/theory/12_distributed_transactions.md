# Distributed Transactions (Розподілені Транзакції)

## 📖 Приклад: Зняття грошей з банківського рахунку

---

## Сценарій 1: Монолітна система (проста транзакція)

### Схема процесу

```
User clicks "Withdraw $100"
         ↓
    [API Server]
         ↓
    BEGIN TRANSACTION
         ↓
1. Check balance >= $100 ────────┐
         ↓                        │
2. Deduct $100 from account      │ Database
         ↓                        │ Transaction
3. Create withdrawal record      │
         ↓                        │
4. Update account.updated_at ────┘
         ↓
    COMMIT TRANSACTION
         ↓
    Return success ✅
```

### Go Code

```go
func WithdrawMoney(userID int64, amount float64) error {
    // Start transaction
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback() // Rollback if not committed
    
    // 1. Check balance
    var balance float64
    err = tx.QueryRow(
        "SELECT balance FROM accounts WHERE user_id = $1 FOR UPDATE",
        userID,
    ).Scan(&balance)
    if err != nil {
        return err
    }
    
    if balance < amount {
        return errors.New("insufficient funds")
    }
    
    // 2. Deduct money
    _, err = tx.Exec(
        "UPDATE accounts SET balance = balance - $1 WHERE user_id = $2",
        amount, userID,
    )
    if err != nil {
        return err
    }
    
    // 3. Create withdrawal record
    _, err = tx.Exec(
        "INSERT INTO transactions (user_id, type, amount) VALUES ($1, 'withdrawal', $2)",
        userID, amount,
    )
    if err != nil {
        return err
    }
    
    // 4. Commit
    if err = tx.Commit(); err != nil {
        return err
    }
    
    return nil // ✅ Transaction complete!
}
```

### ✅ Транзакція закінчена: В одній базі даних
- ACID гарантії
- Rollback якщо помилка
- Просто і надійно

---

## Сценарій 2: Мікросервіси (розподілена транзакція)

### Архітектура

```
                    API Gateway
                         |
        ┌────────────────┼────────────────┐
        ↓                ↓                ↓
  Account Service   Payment Service   Notification Service
  (PostgreSQL)      (PostgreSQL)      (MongoDB)
```

### Проблема: Транзакція через кілька сервісів!

```
User withdraws $100
         ↓
[Account Service]
├─ Check balance ✅
├─ Deduct $100 ✅
└─ Save transaction ✅
         ↓
[Payment Service]
├─ Process payment... ❌ FAILS!
└─ Timeout / Network error
         ↓
[Notification Service]
└─ Never reached

⚠️ PROBLEM: Гроші зняті з рахунку, але платіж не пройшов!
```

---

## Рішення 1: Two-Phase Commit (2PC)

### Схема процесу

```
Coordinator (Orchestrator)
         |
    PHASE 1: PREPARE
         |
    ┌────┴────┬────────┬────────┐
    ↓         ↓        ↓        ↓
Service A  Service B  Service C  Service D
    |         |        |        |
"Ready?"  "Ready?"  "Ready?" "Ready?"
    |         |        |        |
   Yes ✅    Yes ✅   Yes ✅   Yes ✅
    |         |        |        |
    └────┬────┴────────┴────────┘
         ↓
    PHASE 2: COMMIT
         |
    ┌────┴────┬────────┬────────┐
    ↓         ↓        ↓        ↓
Commit ✅  Commit ✅  Commit ✅  Commit ✅
```

### Go Code (Simplified)

```go
type Coordinator struct {
    services []TransactionalService
}

type TransactionalService interface {
    Prepare(ctx context.Context, txID string) error
    Commit(ctx context.Context, txID string) error
    Rollback(ctx context.Context, txID string) error
}

func (c *Coordinator) ExecuteTransaction(ctx context.Context) error {
    txID := generateTxID()
    
    // PHASE 1: PREPARE
    for _, service := range c.services {
        if err := service.Prepare(ctx, txID); err != nil {
            // Rollback all
            c.rollbackAll(ctx, txID)
            return fmt.Errorf("prepare failed: %w", err)
        }
    }
    
    // PHASE 2: COMMIT
    for _, service := range c.services {
        if err := service.Commit(ctx, txID); err != nil {
            // ⚠️ Point of no return!
            log.Error("commit failed, but can't rollback")
            return err
        }
    }
    
    return nil // ✅ Transaction complete!
}

func (c *Coordinator) rollbackAll(ctx context.Context, txID string) {
    for _, service := range c.services {
        service.Rollback(ctx, txID)
    }
}
```

### Account Service Implementation

```go
type AccountService struct {
    db *sql.DB
}

func (s *AccountService) Prepare(ctx context.Context, txID string) error {
    // Start local transaction but don't commit
    tx, _ := s.db.Begin()
    
    // Save transaction handle
    transactions[txID] = tx
    
    // Do work
    _, err := tx.Exec("UPDATE accounts SET balance = balance - 100")
    if err != nil {
        tx.Rollback()
        return err
    }
    
    // Don't commit yet!
    return nil // ✅ Ready to commit
}

func (s *AccountService) Commit(ctx context.Context, txID string) error {
    tx := transactions[txID]
    return tx.Commit() // ✅ Final commit
}

func (s *AccountService) Rollback(ctx context.Context, txID string) error {
    tx := transactions[txID]
    return tx.Rollback()
}
```

### ✅ Переваги 2PC
- Строга consistency
- All-or-nothing гарантії

### ❌ Недоліки 2PC
- Blocking (якщо coordinator fails)
- Performance overhead
- Single point of failure

---

## Рішення 2: Saga Pattern (Рекомендується!)

### Схема: Choreography-based Saga

```
User withdraws $100
         ↓
┌─────────────────────────────────────────────┐
│ Step 1: Account Service                     │
│ ├─ Deduct $100 ✅                           │
│ └─ Publish event: MoneyDeducted             │
└─────────────────────────────────────────────┘
         ↓
┌─────────────────────────────────────────────┐
│ Step 2: Payment Service                     │
│ ├─ Listen: MoneyDeducted                    │
│ ├─ Process payment ✅                       │
│ └─ Publish event: PaymentProcessed          │
└─────────────────────────────────────────────┘
         ↓
┌─────────────────────────────────────────────┐
│ Step 3: Notification Service                │
│ ├─ Listen: PaymentProcessed                 │
│ └─ Send email ✅                            │
└─────────────────────────────────────────────┘

✅ Transaction complete!
```

### Що якщо помилка?

```
User withdraws $100
         ↓
Step 1: Deduct $100 ✅
         ↓
Step 2: Payment fails ❌
         ↓
┌─────────────────────────────────────────────┐
│ COMPENSATING TRANSACTION                    │
│ ├─ Publish event: PaymentFailed             │
│ └─ Account Service listens                  │
│    └─ Refund $100 ✅ (compensate)           │
└─────────────────────────────────────────────┘

✅ Transaction rolled back!
```

### Go Code: Saga Implementation

```go
// Event types
type Event struct {
    Type      string
    TxID      string
    UserID    int64
    Amount    float64
    Timestamp time.Time
}

// Account Service
type AccountService struct {
    db    *sql.DB
    queue MessageQueue
}

func (s *AccountService) Withdraw(userID int64, amount float64) error {
    txID := generateTxID()
    
    // Local transaction
    tx, _ := s.db.Begin()
    defer tx.Rollback()
    
    // Deduct money
    _, err := tx.Exec(
        "UPDATE accounts SET balance = balance - $1 WHERE user_id = $2",
        amount, userID,
    )
    if err != nil {
        return err
    }
    
    // Save transaction log (for compensating)
    _, err = tx.Exec(
        "INSERT INTO saga_log (tx_id, user_id, amount, status) VALUES ($1, $2, $3, 'pending')",
        txID, userID, amount,
    )
    if err != nil {
        return err
    }
    
    if err = tx.Commit(); err != nil {
        return err
    }
    
    // Publish event
    event := Event{
        Type:   "MoneyDeducted",
        TxID:   txID,
        UserID: userID,
        Amount: amount,
    }
    s.queue.Publish("account-events", event)
    
    return nil
}

// Compensating transaction (refund)
func (s *AccountService) HandlePaymentFailed(event Event) {
    tx, _ := s.db.Begin()
    defer tx.Rollback()
    
    // Refund money
    tx.Exec(
        "UPDATE accounts SET balance = balance + $1 WHERE user_id = $2",
        event.Amount, event.UserID,
    )
    
    // Update saga log
    tx.Exec(
        "UPDATE saga_log SET status = 'compensated' WHERE tx_id = $1",
        event.TxID,
    )
    
    tx.Commit()
    
    log.Printf("✅ Refunded %v to user %d", event.Amount, event.UserID)
}

// Payment Service
type PaymentService struct {
    db    *sql.DB
    queue MessageQueue
}

func (s *PaymentService) HandleMoneyDeducted(event Event) {
    err := s.processPayment(event)
    
    if err != nil {
        // Publish failure event (triggers compensating transaction)
        s.queue.Publish("payment-events", Event{
            Type:   "PaymentFailed",
            TxID:   event.TxID,
            UserID: event.UserID,
            Amount: event.Amount,
        })
        return
    }
    
    // Publish success event
    s.queue.Publish("payment-events", Event{
        Type:   "PaymentProcessed",
        TxID:   event.TxID,
        UserID: event.UserID,
        Amount: event.Amount,
    })
}
```

### ✅ Переваги Saga
- No blocking
- Eventual consistency
- Good for microservices
- Scalable

### ❌ Недоліки Saga
- Складніше реалізувати
- Eventual consistency (не immediate)
- Потребує compensating transactions

---

## Рішення 3: Event Sourcing + CQRS

### Схема

```
User withdraws $100
         ↓
┌─────────────────────────────────────────────┐
│ Event Store (append-only log)              │
│ ├─ Event 1: WithdrawalRequested            │
│ ├─ Event 2: BalanceChecked ✅              │
│ ├─ Event 3: MoneyDeducted ✅               │
│ ├─ Event 4: PaymentProcessed ✅            │
│ └─ Event 5: WithdrawalCompleted ✅         │
└─────────────────────────────────────────────┘
         ↓
    Projections (Read Models)
         ↓
┌─────────────────────────────────────────────┐
│ Account Balance View: $900                  │
│ Transaction History View: [..., -$100]      │
└─────────────────────────────────────────────┘
```

### Go Code

```go
type Event struct {
    ID        string
    Type      string
    Aggregate string // e.g., "account:123"
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

func (es *EventStore) GetEvents(aggregateID string) ([]Event, error) {
    rows, err := es.db.Query(
        "SELECT id, type, aggregate, data, timestamp FROM events WHERE aggregate = $1 ORDER BY timestamp",
        aggregateID,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var events []Event
    for rows.Next() {
        var e Event
        rows.Scan(&e.ID, &e.Type, &e.Aggregate, &e.Data, &e.Timestamp)
        events = append(events, e)
    }
    
    return events, nil
}

// Rebuild account state from events
func RebuildAccountState(eventStore *EventStore, accountID string) (*Account, error) {
    events, err := eventStore.GetEvents(fmt.Sprintf("account:%s", accountID))
    if err != nil {
        return nil, err
    }
    
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

### ✅ Переваги Event Sourcing
- Повна історія змін
- Audit trail
- Time travel (replay events)
- Projections для різних use cases

### ❌ Недоліки Event Sourcing
- Складно реалізувати
- Query складніше
- Storage overhead

---

## Рішення 4: Outbox Pattern

### Проблема: Що якщо база даних committed, але message queue failed?

```
❌ BAD:
1. Update DB ✅
2. Send message to queue ❌ (network error)

Result: Inconsistency!
```

### Рішення: Outbox Pattern

```
✅ GOOD:

1. Begin transaction
2. Update accounts table
3. Insert into outbox table ─┐
4. Commit transaction ───────┘ (atomic!)
         ↓
Background worker polls outbox
         ↓
Publish to message queue ✅
         ↓
Mark as sent in outbox
```

### Go Code

```go
func WithdrawWithOutbox(userID int64, amount float64) error {
    tx, _ := db.Begin()
    defer tx.Rollback()
    
    // 1. Deduct money
    _, err := tx.Exec(
        "UPDATE accounts SET balance = balance - $1 WHERE user_id = $2",
        amount, userID,
    )
    if err != nil {
        return err
    }
    
    // 2. Insert into outbox (same transaction!)
    event := Event{
        Type:   "MoneyDeducted",
        UserID: userID,
        Amount: amount,
    }
    eventJSON, _ := json.Marshal(event)
    
    _, err = tx.Exec(
        "INSERT INTO outbox (event_type, payload, created_at) VALUES ($1, $2, NOW())",
        "MoneyDeducted", eventJSON,
    )
    if err != nil {
        return err
    }
    
    // 3. Commit (both happen together!)
    return tx.Commit()
}

// Background worker
func OutboxWorker(db *sql.DB, queue MessageQueue) {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        // Get unpublished events
        rows, _ := db.Query(
            "SELECT id, event_type, payload FROM outbox WHERE published = false LIMIT 100",
        )
        defer rows.Close()
        
        for rows.Next() {
            var id int64
            var eventType string
            var payload []byte
            rows.Scan(&id, &eventType, &payload)
            
            // Publish to queue
            err := queue.Publish(eventType, payload)
            if err != nil {
                continue // Retry later
            }
            
            // Mark as published
            db.Exec("UPDATE outbox SET published = true WHERE id = $1", id)
        }
    }
}
```

---

## Порівняння підходів

| Підхід           | Consistency     | Complexity | Performance | Use Case                    |
|------------------|-----------------|------------|-------------|-----------------------------|
| Monolith ACID    | Immediate ✅    | Low ⭐     | High ✅     | Single database             |
| 2PC              | Immediate ✅    | High ⭐⭐⭐ | Low ❌      | Rare (legacy systems)       |
| **Saga**         | Eventual ⚠️     | Medium ⭐⭐ | High ✅     | **Microservices (popular)** |
| Event Sourcing   | Eventual ⚠️     | High ⭐⭐⭐ | Medium      | Audit trail needed          |
| Outbox Pattern   | Eventual ⚠️     | Medium ⭐⭐ | High ✅     | **With message queues**     |

---

## Де транзакція "закінчена"?

### Монолітна система
```
✅ COMMIT в базі даних = транзакція завершена
```

### Мікросервіси (Saga)
```
✅ Останній крок Saga успішний = транзакція завершена
   (або всі compensating transactions виконані = rollback)
```

### Event Sourcing
```
✅ Останній event в Event Store = транзакція завершена
```

### Outbox Pattern
```
✅ Event published to queue + acknowledged = транзакція завершена
```

---

## Best Practices

### 1. Ідемпотентність (Idempotency)

```go
func ProcessPayment(paymentID string) error {
    // Check if already processed
    var exists bool
    db.QueryRow(
        "SELECT EXISTS(SELECT 1 FROM payments WHERE id = $1)",
        paymentID,
    ).Scan(&exists)
    
    if exists {
        return nil // Already processed ✅
    }
    
    // Process payment...
}
```

### 2. Retry з експоненціальним backoff

```go
func RetryWithBackoff(fn func() error, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        err := fn()
        if err == nil {
            return nil
        }
        
        // Exponential backoff: 1s, 2s, 4s, 8s...
        time.Sleep(time.Duration(1<<uint(i)) * time.Second)
    }
    return errors.New("max retries exceeded")
}
```

### 3. Distributed Tracing

```go
import "go.opentelemetry.io/otel/trace"

func WithdrawMoney(ctx context.Context, userID int64) error {
    ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "withdraw-money")
    defer span.End()
    
    // Your code...
    // Trace propagates through all services
}
```

---

## Висновок

### ✅ Монолітна система
- Проста ACID транзакція
- Одна база даних
- Immediate consistency

### ✅ Мікросервіси
- **Saga Pattern** (найпопулярніше!)
- Outbox Pattern для reliability
- Eventual consistency
- Compensating transactions

**Транзакція "закінчена" коли всі кроки виконані або всі compensations виконані!** 🎯
