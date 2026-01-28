# 🏧 Банкомат не видав гроші - Як транзакція?

## ❌ Проблема

```
User requests $100
         ↓
Database: Deduct $100 ✅
         ↓
ATM: Dispense cash
         ↓
   💥 Механізм застряв!
         ↓
❌ Гроші списані, але НЕ видані!
```

---

## ❓ Чому НЕ проста ACID транзакція?

```
BEGIN TRANSACTION
├─ UPDATE accounts SET balance = balance - 100 ✅
├─ ATM.dispenseCash() ❌ НЕ в БД!
└─ COMMIT

Problem: ATM - це зовнішня hardware система!
         БД не може контролювати механіку!
         Не можна зробити ROLLBACK механічної видачі!
```

---

## ✅ Рішення: Reserve → Try → Confirm/Refund

### Архітектура

```
1. RESERVE (не списувати!)
├─ balance: $1000
├─ available: $900 (hold $100)
└─ status: 'reserved'

2. TRY DISPENSE
├─ Команда до ATM
└─ Чекати на відповідь

3a. SUCCESS → CONFIRM
    ├─ balance: $900 (списати)
    └─ status: 'completed'

3b. FAILURE → REFUND
    ├─ balance: $1000 (не змінився)
    ├─ available: $1000 (розрезервувати)
    └─ status: 'refunded'
```

---

## 💻 Код (спрощено)

### Step 1: Reserve

```go
func reserveMoney(userID int64, amount float64) (txnID int64, err error) {
    tx, _ := db.Begin()
    defer tx.Rollback()
    
    // НЕ списувати, а тільки зарезервувати
    tx.Exec("INSERT INTO account_holds (user_id, amount, status) VALUES ($1, $2, 'active')",
        userID, amount)
    
    tx.Exec("UPDATE accounts SET available_balance = balance - $1 WHERE user_id = $2",
        amount, userID)
    
    tx.Exec("INSERT INTO atm_transactions (user_id, amount, status) VALUES ($1, $2, 'reserved') RETURNING id",
        userID, amount).Scan(&txnID)
    
    tx.Commit()
    return txnID, nil
}
```

### Step 2: Try Dispense

```go
func tryDispense(txnID int64, amount float64) (success bool, err error) {
    // Update status
    db.Exec("UPDATE atm_transactions SET status = 'dispensing' WHERE id = $1", txnID)
    
    // Try ATM (з timeout!)
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    atmTxnID, err := atm.DispenseCash(ctx, amount)
    if err != nil {
        // ATM failed!
        db.Exec("UPDATE atm_transactions SET status = 'failed' WHERE id = $1", txnID)
        return false, err
    }
    
    return true, nil
}
```

### Step 3a: Confirm (якщо SUCCESS)

```go
func confirmWithdrawal(txnID int64) error {
    tx, _ := db.Begin()
    defer tx.Rollback()
    
    // Остаточно списати
    tx.Exec("UPDATE accounts SET balance = balance - (SELECT amount FROM atm_transactions WHERE id = $1)", txnID)
    
    // Release hold
    tx.Exec("UPDATE account_holds SET status = 'released' WHERE transaction_id = $1", txnID)
    
    // Update status
    tx.Exec("UPDATE atm_transactions SET status = 'completed' WHERE id = $1", txnID)
    
    return tx.Commit()
}
```

### Step 3b: Refund (якщо FAILURE)

```go
func refundMoney(txnID int64) error {
    tx, _ := db.Begin()
    defer tx.Rollback()
    
    // НЕ списувати (balance залишається)
    
    // Release hold (розрезервувати)
    tx.Exec("UPDATE account_holds SET status = 'released' WHERE transaction_id = $1", txnID)
    
    // Restore available balance
    tx.Exec("UPDATE accounts SET available_balance = balance WHERE user_id = (SELECT user_id FROM atm_transactions WHERE id = $1)", txnID)
    
    // Update status
    tx.Exec("UPDATE atm_transactions SET status = 'refunded' WHERE id = $1", txnID)
    
    return tx.Commit()
}
```

### Повний Flow

```go
func WithdrawCash(userID int64, amount float64) error {
    // 1. Reserve
    txnID, err := reserveMoney(userID, amount)
    if err != nil {
        return err // Insufficient funds
    }
    
    // 2. Try dispense
    success, err := tryDispense(txnID, amount)
    if err != nil {
        // 3b. Refund (compensating transaction)
        refundMoney(txnID)
        return err // ATM failed
    }
    
    // 3a. Confirm
    return confirmWithdrawal(txnID)
}
```

---

## 🔄 Flow Diagram

### Успіх

```
Reserve $100
├─ balance: $1000
├─ available: $900
└─ hold: $100
         ↓
ATM dispenses ✅
         ↓
Confirm
├─ balance: $900 (deducted)
└─ hold released
         ↓
✅ User has cash
✅ Balance correct
```

### Помилка

```
Reserve $100
├─ balance: $1000
├─ available: $900
└─ hold: $100
         ↓
ATM fails ❌
         ↓
Refund (Compensating)
├─ balance: $1000 (unchanged!)
├─ available: $1000 (restored)
└─ hold released
         ↓
❌ User has NO cash
✅ Balance correct (not deducted!)
```

---

## ⚠️ Edge Case: Timeout

```
Problem:
├─ ATM видав гроші ✅
├─ Але network timeout ❌
└─ Service думає що НЕ видано

Solution: Reconciliation
├─ Background job перевіряє stuck transactions
├─ Запитує ATM: "Чи видав ти txn X?"
│   ├─ YES → Confirm в БД ✅
│   └─ NO → Refund в БД ✅
```

```go
// Reconciliation job
func reconcileStuckTransactions() {
    // Find stuck (status='dispensing' > 5 min)
    rows, _ := db.Query("SELECT id, atm_txn_id FROM atm_transactions WHERE status = 'dispensing' AND updated_at < NOW() - INTERVAL '5 minutes'")
    
    for rows.Next() {
        var txnID int64
        var atmTxnID string
        rows.Scan(&txnID, &atmTxnID)
        
        // Ask ATM
        dispensed, _ := atm.CheckTransactionStatus(atmTxnID)
        
        if dispensed {
            confirmWithdrawal(txnID) // Гроші видані
        } else {
            refundMoney(txnID) // Гроші НЕ видані
        }
    }
}
```

---

## 📊 State Machine

```
[RESERVED] → [DISPENSING] → [COMPLETED] ✅
                  ↓
              [FAILED] → [REFUNDED] ✅
```

---

## 🎯 Ключові принципи

### 1. Reserve, не Deduct
```
✅ Reserve → Try → Confirm
❌ Deduct → Try → Refund
```

### 2. Compensating Transaction
```
Якщо ATM fails:
└─> Refund (compensate = розрезервувати)
```

### 3. Idempotency
```go
idempotencyKey := fmt.Sprintf("%d-%f-%s", userID, amount, requestID)
```

### 4. Timeout
```go
ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
```

### 5. Reconciliation
```go
// Background job для stuck transactions
go reconciliationJob()
```

---

## 🎓 Чому це НЕ проста ACID?

| Проста ACID | ATM Transaction |
|-------------|-----------------|
| Одна БД ✅ | БД + Hardware ❌ |
| BEGIN-COMMIT ✅ | Reserve-Try-Confirm ⚠️ |
| Rollback в БД ✅ | Compensating Transaction ⚠️ |
| Immediate consistency ✅ | Eventual consistency ⚠️ |

---

## 📖 Читати повний файл

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_7

cat theory/15_external_systems_transactions.md
```

---

## 🎯 Висновок

**Коли є зовнішня система (ATM, payment gateway, shipping API):**

```
1. RESERVE (не змінюй незворотно!)
2. TRY (спробуй зовнішню операцію)
3a. SUCCESS → CONFIRM (finalize)
3b. FAILURE → REFUND (compensate)
4. RECONCILIATION (для edge cases)
```

**Це pattern для всіх зовнішніх систем, не тільки банкомат!** 🏧

**Файл:** `theory/15_external_systems_transactions.md`  
**Обсяг:** Повна реалізація + edge cases + reconciliation
