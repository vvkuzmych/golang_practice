# Hardware-Software Integration: ATM State Machine

## 🎯 Питання: Як програма знає що гроші видані?

**Коротка відповідь:** Через **sensors (датчики)** і **hardware events**!

---

## 🏧 Архітектура ATM

### Повна схема

```
┌─────────────────────────────────────────────────────────────┐
│                    SOFTWARE LAYER                            │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Application (Banking Software)                       │   │
│  │  ├─ Transaction State Machine                        │   │
│  │  ├─ Database (Reserve/Confirm)                       │   │
│  │  └─ API calls to Core Banking                        │   │
│  └──────────────────┬───────────────────────────────────┘   │
│                     ↓                                         │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  ATM Middleware (XFS / NDC Protocol)                 │   │
│  │  ├─ Command Queue                                    │   │
│  │  ├─ Event Handler                                    │   │
│  │  └─ Hardware Abstraction Layer                       │   │
│  └──────────────────┬───────────────────────────────────┘   │
└────────────────────┼─────────────────────────────────────────┘
                     ↓
        ═════════════════════════
              HARDWARE BUS
        (USB / Serial / I2C)
        ═════════════════════════
                     ↓
┌─────────────────────────────────────────────────────────────┐
│                   HARDWARE LAYER                             │
│                                                              │
│  ┌─────────────────────────────────────────────┐            │
│  │  Cash Dispenser (CDU)                       │            │
│  │  ├─ Motors (витягування купюр)              │            │
│  │  ├─ Sensors:                                │            │
│  │  │  ├─ Note Counter Sensor ✅               │            │
│  │  │  ├─ Exit Sensor ✅                       │            │
│  │  │  ├─ Shutter Sensor ✅                    │            │
│  │  │  └─ Jam Sensor ⚠️                        │            │
│  │  └─ Firmware (мікроконтролер)               │            │
│  └─────────────────┬───────────────────────────┘            │
│                    ↓                                         │
│           Physical Cash Output                               │
│                 💵💵💵                                       │
└─────────────────────────────────────────────────────────────┘
```

---

## 🔌 Як це працює: Event-Driven State Machine

### Flow Diagram

```
SOFTWARE                    MIDDLEWARE              HARDWARE
   |                            |                       |
   | 1. Command: Dispense $100  |                       |
   |─────────────────────────>  |                       |
   |                            |                       |
   |                            | 2. CMD: Count 5 notes |
   |                            |───────────────────>   |
   |                            |                       |
   |                            |                 [Motors ON]
   |                            |                 [Counting...]
   |                            |                       |
   |                            | 3. EVENT: Note 1 ✅   |
   |                            |<──────────────────    |
   |                            | 4. EVENT: Note 2 ✅   |
   |                            |<──────────────────    |
   |                            | 5. EVENT: Note 3 ✅   |
   |                            |<──────────────────    |
   |                            | 6. EVENT: Note 4 ✅   |
   |                            |<──────────────────    |
   |                            | 7. EVENT: Note 5 ✅   |
   |                            |<──────────────────    |
   |                            |                       |
   |                            | 8. EVENT: Exit Sensor ✅
   |                            |<──────────────────    |
   |                            |                [Shutter Close]
   |                            |                       |
   | 9. CALLBACK: Success ✅     |                       |
   |<─────────────────────────  |                       |
   |                            |                       |
   | 10. Update DB: Confirmed   |                       |
   |                            |                       |
```

---

## 📡 Sensors (Датчики)

### 1. Note Counter Sensor

**Що робить:** Підраховує кожну купюру, що проходить

```
Hardware:
├─ Infrared sensor (інфрачервоний)
├─ Note passes через sensor
└─> Event: "NOTE_COUNTED"

Firmware код (псевдокод):
if (infrared_beam_broken) {
    note_count++
    send_event("NOTE_COUNTED", note_count)
}
```

### 2. Exit Sensor

**Що робить:** Детектує що купюри вийшли до клієнта

```
Hardware:
├─ Sensor в виході (shutter)
├─ Notes pass через exit
└─> Event: "NOTES_PRESENTED"

Firmware код:
if (exit_sensor_triggered && shutter_open) {
    send_event("NOTES_PRESENTED", total_count)
}
```

### 3. Shutter Sensor

**Що робить:** Контролює дверцю виходу

```
States:
├─ CLOSED (default)
├─ OPENING (в процесі)
├─ OPEN (купюри можна взяти)
└─ CLOSING (після взяття)

Events:
├─> "SHUTTER_OPENED"
├─> "SHUTTER_CLOSED"
└─> "CUSTOMER_TOOK_CASH" (за timeout)
```

### 4. Jam Sensor

**Що робить:** Детектує застрявання купюр

```
Hardware:
├─ Mechanical sensor
├─ Якщо купюра застрягла
└─> Event: "NOTE_JAM"

Firmware код:
if (motor_current_too_high || note_stuck) {
    motors_stop()
    send_event("NOTE_JAM", position)
}
```

---

## 🔄 State Machine з Hardware Events

### Go Implementation

```go
type ATMState string

const (
    StateIdle          ATMState = "IDLE"
    StateDispensingCmd ATMState = "DISPENSING_COMMAND_SENT"
    StateCounting      ATMState = "COUNTING_NOTES"
    StatePresenting    ATMState = "PRESENTING_CASH"
    StateCompleted     ATMState = "COMPLETED"
    StateError         ATMState = "ERROR"
)

type HardwareEvent struct {
    Type      string    // "NOTE_COUNTED", "NOTES_PRESENTED", "NOTE_JAM"
    Data      map[string]interface{}
    Timestamp time.Time
}

type ATMStateMachine struct {
    currentState   ATMState
    expectedNotes  int
    countedNotes   int
    txID           string
    eventChan      chan HardwareEvent
    db             *sql.DB
    mu             sync.Mutex
}

func NewATMStateMachine(txID string, expectedNotes int) *ATMStateMachine {
    return &ATMStateMachine{
        currentState:  StateIdle,
        expectedNotes: expectedNotes,
        countedNotes:  0,
        txID:          txID,
        eventChan:     make(chan HardwareEvent, 100),
    }
}

// Головний loop - слухає hardware events
func (sm *ATMStateMachine) Run(ctx context.Context) error {
    for {
        select {
        case event := <-sm.eventChan:
            if err := sm.handleEvent(event); err != nil {
                return err
            }
            
        case <-ctx.Done():
            return ctx.Err()
        }
    }
}

// Hardware event handler
func (sm *ATMStateMachine) handleEvent(event HardwareEvent) error {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    
    log.Printf("[%s] Received event: %s (state: %s)", sm.txID, event.Type, sm.currentState)
    
    switch sm.currentState {
    case StateDispensingCmd:
        return sm.handleDispensingState(event)
        
    case StateCounting:
        return sm.handleCountingState(event)
        
    case StatePresenting:
        return sm.handlePresentingState(event)
        
    default:
        log.Printf("Unexpected event %s in state %s", event.Type, sm.currentState)
    }
    
    return nil
}

func (sm *ATMStateMachine) handleCountingState(event HardwareEvent) error {
    switch event.Type {
    case "NOTE_COUNTED":
        sm.countedNotes++
        log.Printf("[%s] Note counted: %d / %d", sm.txID, sm.countedNotes, sm.expectedNotes)
        
        if sm.countedNotes == sm.expectedNotes {
            // Всі купюри пораховані, чекаємо presentation
            sm.currentState = StatePresenting
            log.Printf("[%s] All notes counted, waiting for presentation", sm.txID)
        }
        
    case "NOTE_JAM":
        log.Printf("[%s] ERROR: Note jam detected!", sm.txID)
        sm.currentState = StateError
        return sm.handleJam(event)
        
    default:
        log.Printf("[%s] Unexpected event %s in COUNTING state", sm.txID, event.Type)
    }
    
    return nil
}

func (sm *ATMStateMachine) handlePresentingState(event HardwareEvent) error {
    switch event.Type {
    case "NOTES_PRESENTED":
        // Hardware підтвердив що купюри вийшли!
        log.Printf("[%s] ✅ Cash presented to customer!", sm.txID)
        sm.currentState = StateCompleted
        
        // Update database: CONFIRM transaction
        return sm.confirmTransaction()
        
    case "CUSTOMER_TOOK_CASH":
        // Customer взяв гроші (shutter закрився)
        log.Printf("[%s] ✅ Customer took cash", sm.txID)
        
    case "PRESENTATION_TIMEOUT":
        // Customer НЕ взяв гроші - retract
        log.Printf("[%s] ⚠️ Customer didn't take cash, retracting", sm.txID)
        return sm.retractCash()
        
    default:
        log.Printf("[%s] Unexpected event %s in PRESENTING state", sm.txID, event.Type)
    }
    
    return nil
}

func (sm *ATMStateMachine) handleJam(event HardwareEvent) error {
    // Hardware jam - потрібен rollback
    log.Printf("[%s] Handling jam, initiating REFUND", sm.txID)
    
    // Update database: REFUND transaction
    tx, _ := sm.db.Begin()
    defer tx.Rollback()
    
    tx.Exec("UPDATE atm_transactions SET status = 'failed', error_message = 'Note jam' WHERE id = $1", sm.txID)
    tx.Exec("UPDATE account_holds SET status = 'released' WHERE transaction_id = $1", sm.txID)
    
    return tx.Commit()
}

func (sm *ATMStateMachine) confirmTransaction() error {
    // Гроші ФІЗИЧНО видані, тепер можна CONFIRM в БД
    log.Printf("[%s] CONFIRMING transaction in database", sm.txID)
    
    tx, _ := sm.db.Begin()
    defer tx.Rollback()
    
    // Списати гроші остаточно
    tx.Exec("UPDATE accounts SET balance = balance - (SELECT amount FROM atm_transactions WHERE id = $1)", sm.txID)
    
    // Update transaction status
    tx.Exec("UPDATE atm_transactions SET status = 'completed', completed_at = NOW() WHERE id = $1", sm.txID)
    
    // Release hold
    tx.Exec("UPDATE account_holds SET status = 'released' WHERE transaction_id = $1", sm.txID)
    
    return tx.Commit()
}

// Зовнішній API: send command до hardware
func (sm *ATMStateMachine) DispenseCash(amount float64, noteCount int) error {
    sm.mu.Lock()
    sm.currentState = StateDispensingCmd
    sm.mu.Unlock()
    
    log.Printf("[%s] Sending dispense command: $%.2f (%d notes)", sm.txID, amount, noteCount)
    
    // Відправити команду до hardware через middleware
    cmd := HardwareCommand{
        Type: "DISPENSE_CASH",
        Data: map[string]interface{}{
            "amount":     amount,
            "note_count": noteCount,
        },
    }
    
    // Після відправки команди, state machine чекає на events
    if err := SendCommandToHardware(cmd); err != nil {
        sm.currentState = StateError
        return err
    }
    
    sm.currentState = StateCounting
    return nil
}
```

---

## 📨 Hardware Communication Protocol

### XFS (eXtensions for Financial Services)

**Industry standard для ATM hardware**

```go
type XFSMiddleware struct {
    conn net.Conn
}

// Відправити команду до hardware
func (x *XFSMiddleware) SendCommand(cmd string, params map[string]interface{}) error {
    message := XFSMessage{
        Command: cmd,
        Params:  params,
    }
    
    data, _ := json.Marshal(message)
    _, err := x.conn.Write(data)
    return err
}

// Слухати events від hardware
func (x *XFSMiddleware) ListenForEvents(eventChan chan<- HardwareEvent) {
    scanner := bufio.NewScanner(x.conn)
    
    for scanner.Scan() {
        var event HardwareEvent
        json.Unmarshal(scanner.Bytes(), &event)
        
        // Відправити event до state machine
        eventChan <- event
    }
}
```

### NDC Protocol (NCR ATMs)

**Бінарний протокол для NCR банкоматів**

```
Message Format:
┌──────┬─────────┬──────────┬─────────┬─────┐
│ STX  │ Command │ Data     │ Checksum│ ETX │
│ 0x02 │ 2 bytes │ Variable │ 1 byte  │0x03 │
└──────┴─────────┴──────────┴─────────┴─────┘

Example: Dispense 5 notes
0x02 0x31 0x30 0x05 0x00 0x00 0x00 0x3A 0x03
     │    │    │                   │
     │    │    └─ Count: 5          └─ Checksum
     └────┴─ Command: 0x3130 (Dispense)

Response: Note counted
0x02 0x32 0x30 0x01 0x3B 0x03
     │    │    │    │
     │    │    │    └─ Checksum
     └────┴────┴─ Event: Note counted (1)
```

---

## 🎬 Повний Life Cycle

### Успішна транзакція

```
Step 1: Reserve in Database
├─ balance: $1000 (не змінюється)
├─ available: $900 (hold $100)
└─ status: 'reserved'
         ↓
Step 2: Send Command to Hardware
├─ Command: "DISPENSE $100 (5 x $20)"
└─ State: DISPENSING_COMMAND_SENT
         ↓
Step 3: Hardware Starts
├─ Motors ON
├─ State: COUNTING_NOTES
└─ Firmware → Software events:
         ↓
Step 4: Hardware Events (real-time)
├─ Event: NOTE_COUNTED (1/5) ✅
├─ Event: NOTE_COUNTED (2/5) ✅
├─ Event: NOTE_COUNTED (3/5) ✅
├─ Event: NOTE_COUNTED (4/5) ✅
├─ Event: NOTE_COUNTED (5/5) ✅
└─ State: PRESENTING_CASH
         ↓
Step 5: Presentation
├─ Event: NOTES_PRESENTED ✅
├─ Shutter OPEN
├─ Customer takes cash
└─ Event: CUSTOMER_TOOK_CASH ✅
         ↓
Step 6: Confirm in Database
├─ balance: $900 (списано!)
├─ status: 'completed'
└─ State: COMPLETED ✅
```

### Помилка: Jam (Застрявання)

```
Step 1-3: Same as above
         ↓
Step 4: Hardware Events
├─ Event: NOTE_COUNTED (1/5) ✅
├─ Event: NOTE_COUNTED (2/5) ✅
├─ Event: NOTE_JAM ❌
└─ Motors STOP
         ↓
Step 5: Error Handling
├─ State: ERROR
├─ Firmware: Retract notes (повернути назад)
└─ Software: Initiate REFUND
         ↓
Step 6: Refund in Database
├─ balance: $1000 (без змін!)
├─ available: $1000 (hold released)
├─ status: 'refunded'
└─ State: ERROR (manual intervention needed)
```

---

## 🧩 Інтеграція: Повний приклад

```go
package main

import (
    "context"
    "database/sql"
    "log"
    "time"
)

type ATMService struct {
    db         *sql.DB
    middleware *XFSMiddleware
}

func (s *ATMService) WithdrawCash(userID int64, amount float64) error {
    // 1. Reserve money in database (не списувати!)
    txID, err := s.reserveMoney(userID, amount)
    if err != nil {
        return err
    }
    
    // 2. Create state machine
    noteCount := int(amount / 20) // Assume $20 notes
    sm := NewATMStateMachine(txID, noteCount)
    sm.db = s.db
    
    // 3. Start state machine (слухати hardware events)
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()
    
    go sm.Run(ctx)
    
    // 4. Connect hardware events до state machine
    go s.middleware.ListenForEvents(sm.eventChan)
    
    // 5. Send command до hardware
    if err := sm.DispenseCash(amount, noteCount); err != nil {
        cancel()
        s.refundMoney(txID)
        return err
    }
    
    // 6. Wait for completion
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            // Timeout - щось пішло не так
            s.refundMoney(txID)
            return ctx.Err()
            
        case <-ticker.C:
            sm.mu.Lock()
            state := sm.currentState
            sm.mu.Unlock()
            
            if state == StateCompleted {
                log.Printf("[%s] ✅ Transaction completed successfully!", txID)
                return nil
            }
            
            if state == StateError {
                log.Printf("[%s] ❌ Transaction failed", txID)
                return fmt.Errorf("dispense failed")
            }
        }
    }
}

func main() {
    db, _ := sql.Open("postgres", "...")
    middleware := NewXFSMiddleware("COM1") // Serial port до hardware
    
    service := &ATMService{
        db:         db,
        middleware: middleware,
    }
    
    // User withdraws $100
    if err := service.WithdrawCash(12345, 100.0); err != nil {
        log.Fatal(err)
    }
}
```

---

## 🎯 Ключові моменти

### 1. **Hardware → Software через Events**

```
Sensor детектує купюру
         ↓
Firmware (мікроконтролер) генерує event
         ↓
Middleware отримує event (через USB/Serial)
         ↓
State Machine обробляє event
         ↓
Database update (якщо потрібно)
```

### 2. **Software → Hardware через Commands**

```
Application: "Dispense $100"
         ↓
State Machine: перехід до DISPENSING
         ↓
Middleware: відправляє binary command
         ↓
Firmware: включає motors
         ↓
Hardware: механіка працює
```

### 3. **State Machine синхронізує все**

```
Software State      Hardware State       Database State
──────────────      ──────────────       ──────────────
RESERVED     ────>  IDLE                 'reserved'
DISPENSING   ────>  MOTORS_ON            'dispensing'
COUNTING     <────  COUNTING_NOTES       'dispensing'
PRESENTING   <────  SHUTTER_OPEN         'dispensing'
COMPLETED    <────  SHUTTER_CLOSED       'completed' ✅
```

---

## 🔍 Debugging

### Logs приклад

```
[TXN-123] Reserve $100 for user 12345
[TXN-123] Sending dispense command: $100 (5 notes)
[TXN-123] State: IDLE → DISPENSING_COMMAND_SENT
[TXN-123] State: DISPENSING_COMMAND_SENT → COUNTING_NOTES
[TXN-123] Received event: NOTE_COUNTED (state: COUNTING_NOTES)
[TXN-123] Note counted: 1 / 5
[TXN-123] Received event: NOTE_COUNTED (state: COUNTING_NOTES)
[TXN-123] Note counted: 2 / 5
[TXN-123] Received event: NOTE_COUNTED (state: COUNTING_NOTES)
[TXN-123] Note counted: 3 / 5
[TXN-123] Received event: NOTE_COUNTED (state: COUNTING_NOTES)
[TXN-123] Note counted: 4 / 5
[TXN-123] Received event: NOTE_COUNTED (state: COUNTING_NOTES)
[TXN-123] Note counted: 5 / 5
[TXN-123] All notes counted, waiting for presentation
[TXN-123] State: COUNTING_NOTES → PRESENTING_CASH
[TXN-123] Received event: NOTES_PRESENTED (state: PRESENTING_CASH)
[TXN-123] ✅ Cash presented to customer!
[TXN-123] CONFIRMING transaction in database
[TXN-123] State: PRESENTING_CASH → COMPLETED
[TXN-123] ✅ Transaction completed successfully!
```

---

## 🎓 Висновок

**Як програма знає що гроші видані?**

1. **Hardware sensors** детектують фізичні події
2. **Firmware** в hardware генерує events
3. **Middleware** пересилає events до application
4. **State Machine** обробляє events і переходить між станами
5. **Database** update тільки коли hardware підтвердив успіх!

**Це event-driven architecture з hardware integration!** 🎯

```
Physical World ─(Sensors)→ Events ─(Middleware)→ State Machine ─(Logic)→ Database
    💵💵💵           ✅          USB/Serial           FSM              💾
```

**Ключове:** Database update (CONFIRM) відбувається **ПІСЛЯ** hardware event `NOTES_PRESENTED`!
