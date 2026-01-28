# 💸 Що якщо купюри НЕ взяті? (Cash Not Taken)

## 🎯 Сценарій

```
ATM: Купюри вийшли, shutter OPEN 💵💵💵
User: ... не бере гроші (забув, відвернувся, пішов)
ATM: ⏱️ Чекає... 30 секунд
ATM: ⚠️ Timeout! Що робити?
```

---

## 🔄 Повний Flow з Timeout

### Схема

```
Step 1: Dispense команда
├─ State: DISPENSING
└─ Motors ON

Step 2: Counting
├─ Event: NOTE_COUNTED (1/5) ✅
├─ Event: NOTE_COUNTED (2/5) ✅
├─ Event: NOTE_COUNTED (3/5) ✅
├─ Event: NOTE_COUNTED (4/5) ✅
└─ Event: NOTE_COUNTED (5/5) ✅

Step 3: Presenting
├─ Event: NOTES_PRESENTED ✅
├─ Shutter OPEN 🚪
├─ Купюри доступні 💵💵💵
└─ Start TIMEOUT timer ⏱️ (30 секунд)

Step 4a: User TAKES cash (нормальний сценарій)
    ├─ Shutter sensor: Cash removed ✅
    ├─ Shutter CLOSE
    ├─ State: COMPLETED
    └─ Database: CONFIRM (списати $100) ✅

Step 4b: User DOES NOT take cash (edge case)
    ├─ ⏱️ 30 seconds elapsed
    ├─ Event: PRESENTATION_TIMEOUT ⚠️
    ├─ Command: RETRACT (втягнути назад)
    ├─ Motors REVERSE
    ├─ Event: NOTES_RETRACTED ✅
    ├─ State: RETRACTED
    └─ Database: REFUND (НЕ списувати!) ✅
```

---

## 🎬 State Machine з Timeout

### States

```
PRESENTING
    ↓
    ├─> CUSTOMER_TOOK_CASH ✅ → COMPLETED
    │   └─> Database: CONFIRM
    │
    └─> TIMEOUT ⚠️ → RETRACTING
        └─> RETRACTED
            └─> Database: REFUND
```

### Go Implementation

```go
type ATMStateMachine struct {
    state           string
    txID            string
    presentedAt     time.Time
    presentTimeout  time.Duration // 30 seconds
    eventChan       chan HardwareEvent
    db              *sql.DB
}

func (sm *ATMStateMachine) handlePresentingState(event HardwareEvent) error {
    switch event.Type {
    case "NOTES_PRESENTED":
        // Купюри вийшли!
        log.Printf("[%s] 💵 Cash presented, waiting for customer...", sm.txID)
        sm.state = "AWAITING_CUSTOMER"
        sm.presentedAt = time.Now()
        
        // Start timeout timer
        go sm.waitForCustomerOrTimeout()
        
    case "CUSTOMER_TOOK_CASH":
        // Customer взяв гроші ✅
        log.Printf("[%s] ✅ Customer took cash!", sm.txID)
        sm.state = "COMPLETED"
        return sm.confirmTransaction()
        
    default:
        log.Printf("[%s] Unexpected event: %s", sm.txID, event.Type)
    }
    
    return nil
}

func (sm *ATMStateMachine) waitForCustomerOrTimeout() {
    timeout := time.NewTimer(sm.presentTimeout) // 30 seconds
    defer timeout.Stop()
    
    <-timeout.C
    
    // Timeout elapsed!
    sm.mu.Lock()
    if sm.state == "AWAITING_CUSTOMER" {
        log.Printf("[%s] ⚠️ TIMEOUT: Customer didn't take cash!", sm.txID)
        
        // Trigger retract
        sm.eventChan <- HardwareEvent{
            Type: "PRESENTATION_TIMEOUT",
        }
    }
    sm.mu.Unlock()
}

func (sm *ATMStateMachine) handleTimeout(event HardwareEvent) error {
    if event.Type == "PRESENTATION_TIMEOUT" {
        log.Printf("[%s] Initiating retract...", sm.txID)
        sm.state = "RETRACTING"
        
        // Send command to hardware: retract cash
        cmd := HardwareCommand{
            Type: "RETRACT_CASH",
        }
        
        if err := SendCommandToHardware(cmd); err != nil {
            log.Printf("[%s] ERROR: Failed to retract: %v", sm.txID, err)
            return err
        }
        
        // Wait for hardware event
        // Hardware буде слати: NOTES_RETRACTED
    }
    
    return nil
}

func (sm *ATMStateMachine) handleRetractingState(event HardwareEvent) error {
    switch event.Type {
    case "NOTES_RETRACTED":
        // Купюри втягнуті назад успішно ✅
        log.Printf("[%s] ✅ Cash retracted successfully", sm.txID)
        sm.state = "RETRACTED"
        
        // REFUND transaction (НЕ списувати гроші!)
        return sm.refundTransaction()
        
    case "RETRACT_FAILED":
        // Помилка втягування ❌
        log.Printf("[%s] ❌ ERROR: Failed to retract cash!", sm.txID)
        sm.state = "ERROR_RETRACT_FAILED"
        
        // Це серйозна проблема - потрібна manual intervention
        return sm.handleRetractFailure()
        
    default:
        log.Printf("[%s] Unexpected event: %s", sm.txID, event.Type)
    }
    
    return nil
}

func (sm *ATMStateMachine) refundTransaction() error {
    // Купюри втягнуті, гроші НЕ списувати!
    log.Printf("[%s] REFUNDING transaction (cash not taken)", sm.txID)
    
    tx, _ := sm.db.Begin()
    defer tx.Rollback()
    
    // Update transaction status
    tx.Exec(`
        UPDATE atm_transactions 
        SET status = 'refunded_not_taken', 
            error_message = 'Customer did not take cash, notes retracted',
            updated_at = NOW() 
        WHERE id = $1
    `, sm.txID)
    
    // Release hold (гроші залишаються на рахунку)
    tx.Exec(`
        UPDATE account_holds 
        SET status = 'released' 
        WHERE transaction_id = $1
    `, sm.txID)
    
    // Balance НЕ змінюється! (гроші не списані)
    
    return tx.Commit()
}
```

---

## 📊 Detailed Flow Diagram

```
┌──────────────────────────────────────────────────────┐
│  NOTES_PRESENTED                                     │
│  ├─ Shutter OPEN 🚪                                  │
│  ├─ Купюри доступні 💵💵💵                            │
│  └─ Start Timer ⏱️ (30 sec)                          │
└──────────────┬───────────────────────────────────────┘
               ↓
         ┌─────┴─────┐
         ↓           ↓
    TOOK CASH    TIMEOUT
         ↓           ↓
         │           │
    ┌────┴────┐  ┌───┴────────────┐
    │ SUCCESS │  │ RETRACT        │
    │         │  │ ├─ Motors REV  │
    │ Confirm │  │ ├─ Pull back   │
    │ in DB ✅ │  │ └─ Close       │
    └─────────┘  └───┬────────────┘
                     ↓
               ┌─────┴─────┐
               ↓           ↓
         RETRACTED    RETRACT_FAIL
               ↓           ↓
          ┌────┴────┐  ┌───┴──────────┐
          │ REFUND  │  │ ERROR        │
          │ (no $)  │  │ Manual fix ⚠️ │
          │      ✅ │  │              │
          └─────────┘  └──────────────┘
```

---

## ⏱️ Timeout Configuration

### Typical Values

```go
const (
    // Час, скільки купюри доступні для взяття
    PresentationTimeout = 30 * time.Second
    
    // Час для retract operation
    RetractTimeout = 15 * time.Second
    
    // Максимальний час для всієї транзакції
    TotalTransactionTimeout = 60 * time.Second
)
```

### Configurable per Bank

```
Bank A (швидкий):
├─ Presentation: 20 sec
└─ Rationale: Busy city center, long queues

Bank B (звичайний):
├─ Presentation: 30 sec
└─ Rationale: Standard

Bank C (доступний):
├─ Presentation: 45 sec
└─ Rationale: Elderly customers, accessibility
```

---

## 🔧 Hardware: Retract Mechanism

### Як працює Retract

```
Physical Process:

1. Shutter OPEN, notes visible 💵💵💵
         ↓
2. Timeout detected ⏱️
         ↓
3. Motors REVERSE direction
         ↓
4. Vacuum/Rollers pull notes back
         ↓
5. Notes go to RETRACT BIN (not dispenser cassette!)
         ↓
6. Shutter CLOSE
         ↓
7. Event: NOTES_RETRACTED ✅
```

### Retract Bin

```
ATM має 2 місця для купюр:

1. DISPENSER CASSETTES
   ├─ Нові купюри для видачі
   └─ Регулярно поповнюються

2. RETRACT BIN
   ├─ Купюри, що не були взяті
   ├─ Купюри, що застрягли (jam)
   └─ Потім перевіряються вручну

Чому не назад в cassette?
└─> Security: Купюра могла бути підміненою!
```

---

## 📝 Database States

### Transaction Status Values

```sql
CREATE TYPE transaction_status AS ENUM (
    'reserved',              -- Гроші зарезервовані
    'dispensing',            -- Механіка працює
    'completed',             -- ✅ Гроші взяті клієнтом
    'refunded_not_taken',    -- ✅ Купюри не взяті, retracted
    'refunded_jam',          -- ✅ Застрявання
    'failed',                -- ❌ Інша помилка
    'error_retract_failed'   -- ❌ Не вдалося втягнути (manual fix!)
);
```

### Example Records

```sql
-- Успішна транзакція
id: 123
status: 'completed'
amount: 100
notes_presented: true
notes_taken: true
completed_at: '2026-01-28 10:15:35'

-- Купюри не взяті
id: 124
status: 'refunded_not_taken'
amount: 100
notes_presented: true
notes_taken: false
notes_retracted: true
error_message: 'Customer timeout, notes retracted'
```

---

## 🚨 Critical Edge Case: Retract Failed

### Що якщо retract НЕ вдався?

```
Problem:
├─ Timeout ⏱️
├─ Command: RETRACT
├─ Hardware: ❌ Failed to retract
│   └─> Купюри застрягли в shutter
│   └─> Або customer взяв ПІСЛЯ timeout
└─> Що робити?

Solution: Manual Intervention Required!
```

### Implementation

```go
func (sm *ATMStateMachine) handleRetractFailure() error {
    log.Printf("[%s] ❌ CRITICAL: Retract failed!", sm.txID)
    
    tx, _ := sm.db.Begin()
    defer tx.Rollback()
    
    // Mark as critical error
    tx.Exec(`
        UPDATE atm_transactions 
        SET status = 'error_retract_failed',
            error_message = 'Failed to retract notes - MANUAL INSPECTION REQUIRED',
            requires_manual_review = true,
            updated_at = NOW() 
        WHERE id = $1
    `, sm.txID)
    
    // НЕ release hold (чекаємо manual review)
    
    tx.Commit()
    
    // Alert support
    alertSupport(AlertCritical, sm.txID, "Retract failed")
    
    // Lock ATM (безпека!)
    lockATM("Retract failure, manual inspection required")
    
    return fmt.Errorf("retract failed, ATM locked")
}

func alertSupport(level AlertLevel, txID string, message string) {
    // Send SMS/Email/PagerDuty до техпідтримки
    log.Printf("🚨 ALERT [%s] %s: %s", level, txID, message)
    
    // Real implementation:
    // - SMS до on-call engineer
    // - Email до support team
    // - PagerDuty incident
    // - Log to monitoring system
}

func lockATM(reason string) {
    // Lock ATM until manual inspection
    log.Printf("🔒 ATM LOCKED: %s", reason)
    
    // Real implementation:
    // - Display "Out of Service" на екрані
    // - Disable card reader
    // - Відправити status до Central Monitoring
}
```

---

## 📊 Statistics & Monitoring

### Metrics to Track

```go
type ATMMetrics struct {
    TotalDispenses      int64
    SuccessfulTaken     int64   // Customer взяв
    NotTakenRetracted   int64   // Timeout → Retracted
    RetractFailures     int64   // Critical!
    
    AverageTakeTime     time.Duration
    TimeoutRate         float64  // %
}

// Example query
SELECT 
    COUNT(*) as total,
    COUNT(*) FILTER (WHERE status = 'completed') as taken,
    COUNT(*) FILTER (WHERE status = 'refunded_not_taken') as not_taken,
    COUNT(*) FILTER (WHERE status = 'error_retract_failed') as critical
FROM atm_transactions
WHERE created_at > NOW() - INTERVAL '1 day';
```

---

## 🎯 Best Practices

### 1. Generous Timeout

```
✅ GOOD: 30 seconds
⚠️ RISKY: 15 seconds (старі люди не встигнуть)
❌ BAD: 10 seconds (занадто швидко)
```

### 2. Visual/Audio Warnings

```
20 seconds: 🔊 "Please take your cash"
25 seconds: 🔊🔊 "Cash will be retracted in 5 seconds!"
30 seconds: RETRACT
```

### 3. Log Everything

```go
log.Printf("[%s] Notes presented at %s", txID, time.Now())
log.Printf("[%s] Waiting for customer... (timeout: %ds)", txID, timeout)
log.Printf("[%s] Timeout! Initiating retract", txID)
log.Printf("[%s] Notes retracted successfully", txID)
```

### 4. Monitoring

```
Alert if:
├─ Timeout rate > 5% (щось не так з ATM?)
├─ Retract failure > 0 (critical!)
└─ Average take time > 20 sec (slow customers?)
```

---

## 🎓 Висновок

### Що відбувається якщо купюри не взяті?

```
1. ⏱️ Timeout (30 sec)
2. 🔄 RETRACT (втягнути назад)
3. 📦 Купюри → Retract Bin
4. 💾 Database: REFUND (гроші не списані!)
5. 👤 Customer: Balance без змін ✅
```

### Ключові моменти

```
✅ Купюри НЕ залишаються в shutter (security)
✅ Гроші НЕ списуються (customer не втрачає)
✅ Retract bin перевіряється щоденно (manual)
⚠️ Якщо retract fails → ATM lock (critical!)
```

---

## 📖 Детальний файл

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_7

# Повна теорія
cat theory/17_hardware_software_integration.md

# Hardware State Machine
cat HARDWARE_STATE_MACHINE.md
```

---

**Safety First: Customer never loses money!** 🎯

**Flow:**  
💵 Presented → ⏱️ Timeout → 🔄 Retract → 💾 Refund → ✅ Balance OK
