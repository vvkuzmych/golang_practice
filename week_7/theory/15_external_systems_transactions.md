# Транзакції з зовнішніми системами (External Systems)

## 🏧 Проблема: Банкомат не видав гроші

### Сценарій

```
User requests $100 from ATM
         ↓
BEGIN TRANSACTION
├─ Check balance ✅
├─ Deduct $100 from account ✅
└─ COMMIT ✅
         ↓
Command ATM to dispense $100
         ↓
   💥 ATM механізм застряв!
         ↓
❌ Гроші НЕ видані, але списані з рахунку!
```

### Чому не можна використати просту ACID транзакцію?

```
❌ НЕ МОЖНА:

BEGIN TRANSACTION
├─ Deduct $100 from database ✅
├─ Dispense cash from ATM   ❌ (не в БД!)
└─ COMMIT

Problem: ATM - це зовнішня hardware система!
         БД не може контролювати механіку!
```

---

## 🎯 Рішення: Compensating Transaction Pattern

### Архітектура

```
┌──────────────────────────────────────────────────┐
│  Database                ATM (Hardware)          │
│                                                  │
│  1. Reserve $100 ✅  →  2. Try dispense         │
│                              ↓                   │
│                         ┌────┴────┐              │
│                         ↓         ↓              │
│                      SUCCESS    FAIL             │
│                         ↓         ↓              │
│  3a. Confirm ✅         3b. Refund ✅            │
│  (finalize)             (compensate)             │
└──────────────────────────────────────────────────┘
```

### Крок за кроком

```
Step 1: RESERVE (не списувати!)
├─ Перевірити баланс
├─ "Зарезервувати" $100 (status = 'reserved')
└─ НЕ списувати остаточно!

Step 2: DISPENSE (спроба видати)
├─ Відправити команду до ATM
├─ Чекати на відповідь
└─ Timeout: 30 seconds

Step 3a: SUCCESS (якщо видано)
├─ Оновити status = 'completed'
└─ Списати остаточно

Step 3b: FAILURE (якщо НЕ видано)
├─ Оновити status = 'refunded'
└─ Розрезервувати (compensating transaction)
```

---

## 💻 Go Implementation

### Database Schema

```sql
CREATE TABLE atm_transactions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) NOT NULL, -- 'reserved', 'dispensing', 'completed', 'failed', 'refunded'
    atm_transaction_id VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP,
    error_message TEXT
);

CREATE TABLE account_holds (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    transaction_id INT REFERENCES atm_transactions(id),
    status VARCHAR(20) NOT NULL, -- 'active', 'released'
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Step 1: Reserve

```go
type ATMService struct {
    db  *sql.DB
    atm ATMClient
}

func (s *ATMService) WithdrawCash(ctx context.Context, userID int64, amount float64) error {
    // Step 1: Reserve money (не списувати!)
    txnID, err := s.reserveMoney(ctx, userID, amount)
    if err != nil {
        return fmt.Errorf("reserve failed: %w", err)
    }
    
    // Step 2: Try to dispense
    atmTxnID, err := s.tryDispenseCash(ctx, txnID, amount)
    if err != nil {
        // Step 3b: Compensate (refund)
        s.refundMoney(ctx, txnID)
        return fmt.Errorf("dispense failed: %w", err)
    }
    
    // Step 3a: Confirm (finalize)
    err = s.confirmWithdrawal(ctx, txnID, atmTxnID)
    if err != nil {
        // Тут складніше - гроші видані, але не confirmed
        // Потрібен manual reconciliation
        s.logCriticalError(ctx, txnID, "Money dispensed but not confirmed")
        return err
    }
    
    return nil
}

func (s *ATMService) reserveMoney(ctx context.Context, userID int64, amount float64) (int64, error) {
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return 0, err
    }
    defer tx.Rollback()
    
    // 1. Check balance (with lock)
    var balance float64
    err = tx.QueryRowContext(ctx,
        "SELECT balance FROM accounts WHERE user_id = $1 FOR UPDATE",
        userID,
    ).Scan(&balance)
    if err != nil {
        return 0, err
    }
    
    if balance < amount {
        return 0, errors.New("insufficient funds")
    }
    
    // 2. Create ATM transaction record
    var txnID int64
    err = tx.QueryRowContext(ctx,
        "INSERT INTO atm_transactions (user_id, amount, status) VALUES ($1, $2, 'reserved') RETURNING id",
        userID, amount,
    ).Scan(&txnID)
    if err != nil {
        return 0, err
    }
    
    // 3. Create hold (не списувати, а тримати)
    _, err = tx.ExecContext(ctx,
        "INSERT INTO account_holds (user_id, amount, transaction_id, status) VALUES ($1, $2, $3, 'active')",
        userID, amount, txnID,
    )
    if err != nil {
        return 0, err
    }
    
    // 4. Update available balance (virtual)
    _, err = tx.ExecContext(ctx,
        "UPDATE accounts SET available_balance = balance - $1 WHERE user_id = $2",
        amount, userID,
    )
    if err != nil {
        return 0, err
    }
    
    if err = tx.Commit(); err != nil {
        return 0, err
    }
    
    log.Printf("✅ Reserved $%.2f for user %d (txn %d)", amount, userID, txnID)
    return txnID, nil
}
```

### Step 2: Try Dispense

```go
type ATMClient interface {
    DispenseCash(ctx context.Context, amount float64) (string, error)
}

func (s *ATMService) tryDispenseCash(ctx context.Context, txnID int64, amount float64) (string, error) {
    // Update status to 'dispensing'
    _, err := s.db.ExecContext(ctx,
        "UPDATE atm_transactions SET status = 'dispensing', updated_at = NOW() WHERE id = $1",
        txnID,
    )
    if err != nil {
        return "", err
    }
    
    log.Printf("💰 Attempting to dispense $%.2f (txn %d)", amount, txnID)
    
    // Call ATM hardware (з timeout!)
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()
    
    atmTxnID, err := s.atm.DispenseCash(ctx, amount)
    if err != nil {
        // ATM failed!
        log.Printf("❌ ATM dispense failed: %v", err)
        
        // Update status
        s.db.ExecContext(context.Background(),
            "UPDATE atm_transactions SET status = 'failed', error_message = $1, updated_at = NOW() WHERE id = $2",
            err.Error(), txnID,
        )
        
        return "", fmt.Errorf("ATM error: %w", err)
    }
    
    log.Printf("✅ ATM dispensed cash (ATM txn: %s)", atmTxnID)
    return atmTxnID, nil
}
```

### Step 3a: Confirm (Success)

```go
func (s *ATMService) confirmWithdrawal(ctx context.Context, txnID int64, atmTxnID string) error {
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    // 1. Update transaction status
    _, err = tx.ExecContext(ctx,
        "UPDATE atm_transactions SET status = 'completed', atm_transaction_id = $1, completed_at = NOW() WHERE id = $2",
        atmTxnID, txnID,
    )
    if err != nil {
        return err
    }
    
    // 2. Actually deduct money (finalize)
    var userID int64
    var amount float64
    err = tx.QueryRowContext(ctx,
        "SELECT user_id, amount FROM atm_transactions WHERE id = $1",
        txnID,
    ).Scan(&userID, &amount)
    if err != nil {
        return err
    }
    
    _, err = tx.ExecContext(ctx,
        "UPDATE accounts SET balance = balance - $1 WHERE user_id = $2",
        amount, userID,
    )
    if err != nil {
        return err
    }
    
    // 3. Release hold
    _, err = tx.ExecContext(ctx,
        "UPDATE account_holds SET status = 'released' WHERE transaction_id = $1",
        txnID,
    )
    if err != nil {
        return err
    }
    
    if err = tx.Commit(); err != nil {
        return err
    }
    
    log.Printf("✅ Confirmed withdrawal (txn %d, ATM txn %s)", txnID, atmTxnID)
    return nil
}
```

### Step 3b: Refund (Failure - Compensating Transaction)

```go
func (s *ATMService) refundMoney(ctx context.Context, txnID int64) error {
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    // 1. Update transaction status
    _, err = tx.ExecContext(ctx,
        "UPDATE atm_transactions SET status = 'refunded', updated_at = NOW() WHERE id = $1",
        txnID,
    )
    if err != nil {
        return err
    }
    
    // 2. Release hold (compensating transaction)
    var userID int64
    var amount float64
    err = tx.QueryRowContext(ctx,
        "SELECT user_id, amount FROM atm_transactions WHERE id = $1",
        txnID,
    ).Scan(&userID, &amount)
    if err != nil {
        return err
    }
    
    _, err = tx.ExecContext(ctx,
        "UPDATE account_holds SET status = 'released' WHERE transaction_id = $1",
        txnID,
    )
    if err != nil {
        return err
    }
    
    // 3. Restore available balance
    _, err = tx.ExecContext(ctx,
        "UPDATE accounts SET available_balance = balance WHERE user_id = $1",
        userID,
    )
    if err != nil {
        return err
    }
    
    if err = tx.Commit(); err != nil {
        return err
    }
    
    log.Printf("✅ Refunded $%.2f to user %d (txn %d)", amount, userID, txnID)
    return nil
}
```

---

## 🔄 Flow Diagram

### Успішний сценарій

```
User → ATM Service
         ↓
    1. Reserve $100
    ├─ balance: $1000
    ├─ available: $900 (reserved)
    ├─ status: 'reserved'
    └─ hold: $100 ✅
         ↓
    2. Dispense cash
    ├─ ATM command sent
    ├─ ATM responds: SUCCESS ✅
    └─ status: 'dispensing'
         ↓
    3. Confirm
    ├─ balance: $900 (deducted)
    ├─ status: 'completed'
    └─ hold: released ✅
         ↓
    User receives cash ✅
```

### Сценарій з помилкою

```
User → ATM Service
         ↓
    1. Reserve $100
    ├─ balance: $1000
    ├─ available: $900
    └─ hold: $100 ✅
         ↓
    2. Dispense cash
    ├─ ATM command sent
    ├─ ATM responds: ERROR ❌
    └─ status: 'failed'
         ↓
    3. Refund (Compensating)
    ├─ balance: $1000 (unchanged)
    ├─ available: $1000 (restored)
    ├─ status: 'refunded'
    └─ hold: released ✅
         ↓
    User does NOT receive cash ✅
    Balance correct ✅
```

---

## ⚠️ Edge Cases

### Edge Case 1: ATM видав, але timeout

```
Problem:
├─ ATM механізм видав гроші ✅
├─ Але network timeout ❌
└─ Service думає що не видано

Solution:
├─ ATM має свій transaction ID
├─ Reconciliation process:
│   └─ Запитати ATM: "Чи видав ти txn X?"
│       ├─ YES → Confirm в DB
│       └─ NO → Refund в DB
```

```go
func (s *ATMService) reconcileTransaction(ctx context.Context, txnID int64) error {
    // Get transaction details
    var atmTxnID sql.NullString
    var status string
    err := s.db.QueryRowContext(ctx,
        "SELECT atm_transaction_id, status FROM atm_transactions WHERE id = $1",
        txnID,
    ).Scan(&atmTxnID, &status)
    if err != nil {
        return err
    }
    
    if status == "dispensing" && atmTxnID.Valid {
        // Check with ATM
        dispensed, err := s.atm.CheckTransactionStatus(ctx, atmTxnID.String)
        if err != nil {
            return err
        }
        
        if dispensed {
            // Money was dispensed, confirm it
            return s.confirmWithdrawal(ctx, txnID, atmTxnID.String)
        } else {
            // Money was NOT dispensed, refund it
            return s.refundMoney(ctx, txnID)
        }
    }
    
    return nil
}
```

### Edge Case 2: Duplicate request

```go
// Idempotency key
func (s *ATMService) WithdrawCashIdempotent(ctx context.Context, idempotencyKey string, userID int64, amount float64) error {
    // Check if already processed
    var existingTxnID sql.NullInt64
    err := s.db.QueryRowContext(ctx,
        "SELECT id FROM atm_transactions WHERE user_id = $1 AND idempotency_key = $2",
        userID, idempotencyKey,
    ).Scan(&existingTxnID)
    
    if err == nil && existingTxnID.Valid {
        // Already processed
        log.Printf("Duplicate request, returning existing txn %d", existingTxnID.Int64)
        return nil
    }
    
    // Process new request...
}
```

---

## 🔧 Reconciliation Process (Звірка)

### Background Job

```go
func (s *ATMService) ReconciliationJob(ctx context.Context) {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            s.reconcileStuckTransactions(ctx)
        case <-ctx.Done():
            return
        }
    }
}

func (s *ATMService) reconcileStuckTransactions(ctx context.Context) {
    // Find stuck transactions (dispensing > 5 minutes)
    rows, err := s.db.QueryContext(ctx, `
        SELECT id, atm_transaction_id 
        FROM atm_transactions 
        WHERE status = 'dispensing' 
          AND updated_at < NOW() - INTERVAL '5 minutes'
    `)
    if err != nil {
        log.Printf("Reconciliation query failed: %v", err)
        return
    }
    defer rows.Close()
    
    for rows.Next() {
        var txnID int64
        var atmTxnID sql.NullString
        rows.Scan(&txnID, &atmTxnID)
        
        log.Printf("🔍 Reconciling stuck transaction %d", txnID)
        
        if err := s.reconcileTransaction(ctx, txnID); err != nil {
            log.Printf("❌ Reconciliation failed for txn %d: %v", txnID, err)
        } else {
            log.Printf("✅ Reconciled transaction %d", txnID)
        }
    }
}
```

---

## 📊 State Machine

```
[RESERVED] ──────┐
     |           │
     ↓           │
[DISPENSING] ────┤
     |           │
     ├──> [COMPLETED] (success path)
     |           │
     └──> [FAILED] ──> [REFUNDED] (failure path)
```

### Valid Transitions

```go
var validTransitions = map[string][]string{
    "reserved":    {"dispensing", "refunded"},
    "dispensing":  {"completed", "failed"},
    "failed":      {"refunded"},
    "completed":   {}, // terminal
    "refunded":    {}, // terminal
}

func (s *ATMService) updateStatus(ctx context.Context, txnID int64, newStatus string) error {
    tx, _ := s.db.BeginTx(ctx, nil)
    defer tx.Rollback()
    
    // Get current status
    var currentStatus string
    tx.QueryRowContext(ctx,
        "SELECT status FROM atm_transactions WHERE id = $1 FOR UPDATE",
        txnID,
    ).Scan(&currentStatus)
    
    // Check if transition is valid
    validNext, ok := validTransitions[currentStatus]
    if !ok {
        return fmt.Errorf("invalid current status: %s", currentStatus)
    }
    
    allowed := false
    for _, status := range validNext {
        if status == newStatus {
            allowed = true
            break
        }
    }
    
    if !allowed {
        return fmt.Errorf("invalid transition from %s to %s", currentStatus, newStatus)
    }
    
    // Update status
    _, err := tx.ExecContext(ctx,
        "UPDATE atm_transactions SET status = $1, updated_at = NOW() WHERE id = $2",
        newStatus, txnID,
    )
    if err != nil {
        return err
    }
    
    return tx.Commit()
}
```

---

## 🎯 Best Practices

### 1. Always Reserve First

```
✅ GOOD: Reserve → Try → Confirm/Refund
❌ BAD:  Deduct → Try → Refund (якщо fail)
```

### 2. Idempotency

```go
// Use unique idempotency key
idempotencyKey := fmt.Sprintf("%d-%s-%s", userID, amount, requestID)
```

### 3. Timeouts

```go
// Always set timeout for external systems
ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
defer cancel()
```

### 4. Reconciliation

```go
// Background job для stuck transactions
go s.ReconciliationJob(ctx)
```

### 5. Monitoring & Alerts

```go
// Alert if too many failures
if failureRate > 0.05 { // 5%
    alert.Send("ATM failure rate too high")
}
```

---

## 📈 Metrics to Track

```go
type ATMMetrics struct {
    TotalRequests      int64
    SuccessfulDispense int64
    FailedDispense     int64
    RefundedAmount     float64
    StuckTransactions  int64
}

func (s *ATMService) GetMetrics() ATMMetrics {
    // Query metrics from DB
}
```

---

## 🎓 Висновок

**Транзакція з зовнішньою системою (банкомат):**

```
1. RESERVE (не списувати!)
   └─> Balance залишається, але "hold" створено

2. TRY EXTERNAL OPERATION
   └─> Спроба взаємодії з hardware

3a. SUCCESS → CONFIRM
    └─> Остаточно списати

3b. FAILURE → REFUND (Compensating)
    └─> Розрезервувати (rollback)

4. RECONCILIATION
   └─> Background job для stuck transactions
```

**Ключові принципи:**

✅ Reserve, не Deduct  
✅ Try external operation  
✅ Confirm або Compensate  
✅ Reconciliation для edge cases  
✅ Idempotency для retry  
✅ Timeouts для всього  
✅ State machine для статусів  

**Це НЕ проста ACID транзакція, тому що зовнішня система (ATM) не в контролі БД!** 🏧
