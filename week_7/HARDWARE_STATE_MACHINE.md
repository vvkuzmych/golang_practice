# 🏧 Як State Machine з'єднана з механікою?

## ❓ Питання: Як програма знає що гроші фізично видані?

**Відповідь: Через SENSORS (датчики) і HARDWARE EVENTS!**

---

## 📊 Текстова схема

### Повна архітектура

```
┌────────────────────────────────────────────────────────┐
│                   SOFTWARE                             │
│                                                        │
│  ┌──────────────────────────────────┐                 │
│  │  State Machine (Go code)         │                 │
│  │  ├─ IDLE                         │                 │
│  │  ├─ DISPENSING                   │                 │
│  │  ├─ COUNTING ←─ Events           │                 │
│  │  ├─ PRESENTING                   │                 │
│  │  └─ COMPLETED ✅                  │                 │
│  └──────────────┬───────────────────┘                 │
│                 ↓                                      │
│  ┌──────────────────────────────────┐                 │
│  │  Middleware (XFS / NDC Protocol) │                 │
│  │  ├─ Commands → Hardware          │                 │
│  │  └─ Events ← Hardware            │                 │
│  └──────────────┬───────────────────┘                 │
└─────────────────┼──────────────────────────────────────┘
                  ↓
        ═══════════════════
          USB / Serial
        ═══════════════════
                  ↓
┌─────────────────────────────────────────────────────────┐
│                   HARDWARE                              │
│                                                         │
│  ┌───────────────────────────────────────┐             │
│  │  Cash Dispenser (CDU)                 │             │
│  │  ├─ Motors (механіка витягування)    │             │
│  │  ├─ Sensors:                          │             │
│  │  │  ├─ Note Counter ✅ (рахує купюри)│             │
│  │  │  ├─ Exit Sensor ✅ (купюри вийшли)│             │
│  │  │  ├─ Shutter ✅ (дверцята)         │             │
│  │  │  └─ Jam Sensor ⚠️ (застрявання)  │             │
│  │  └─ Firmware (генерує events)        │             │
│  └───────────────────────────────────────┘             │
│                      ↓                                  │
│              Physical Cash 💵💵💵                       │
└─────────────────────────────────────────────────────────┘
```

---

## 🔄 Як працює: Event Flow

### Успішна видача грошей

```
SOFTWARE          MIDDLEWARE         HARDWARE
   |                  |                  |
   | 1. CMD: Dispense |                  |
   |───────────────>  |                  |
   |                  |                  |
   | State: DISPENSING|                  |
   |                  |                  |
   |                  | 2. CMD: Count 5  |
   |                  |──────────────>   |
   |                  |            [Motors ON]
   |                  |            [Counting...]
   |                  |                  |
   |                  | 3. EVENT: Note 1 ✅
   |                  |<─────────────    |
   | 4. Handle Event  |                  |
   |<─────────────    |                  |
   |   counted: 1/5   |                  |
   |                  |                  |
   |                  | 5. EVENT: Note 2 ✅
   |                  |<─────────────    |
   |   counted: 2/5   |                  |
   |                  |                  |
   |                  | 6. EVENT: Note 3 ✅
   |                  |<─────────────    |
   |   counted: 3/5   |                  |
   |                  |                  |
   |                  | 7. EVENT: Note 4 ✅
   |                  |<─────────────    |
   |   counted: 4/5   |                  |
   |                  |                  |
   |                  | 8. EVENT: Note 5 ✅
   |                  |<─────────────    |
   |   counted: 5/5   |                  |
   | State: PRESENTING|                  |
   |                  |                  |
   |                  | 9. EVENT: PRESENTED ✅
   |                  |<─────────────    |
   |                  |         [Shutter OPEN]
   | 10. Confirm DB ✅ |                  |
   | State: COMPLETED |                  |
   |                  |                  |
```

---

## 📡 Sensors (Датчики) - Як працюють

### 1. Note Counter Sensor

```
Physical:
├─ Infrared LED і receptor
├─ Купюра проходить між ними
└─> Промінь перериваєтьсяFirmware:
if (infrared_beam_broken) {
    note_count++
    send_event("NOTE_COUNTED", note_count)
}

Software отримує:
Event: { type: "NOTE_COUNTED", count: 3 }
```

### 2. Exit Sensor (Shutter)

```
Physical:
├─ Mechanical sensor в виході
├─ Купюри виходять через shutter
└─> Sensor спрацьовує

Firmware:
if (exit_sensor_triggered && shutter_open) {
    send_event("NOTES_PRESENTED")
}

Software отримує:
Event: { type: "NOTES_PRESENTED", total: 5 }
```

### 3. Jam Sensor

```
Physical:
├─ Motor current sensor
├─ Якщо купюра застрягла → current ↑
└─> Jam detected

Firmware:
if (motor_current_too_high) {
    motors_stop()
    send_event("NOTE_JAM", position)
}

Software отримує:
Event: { type: "NOTE_JAM", position: "exit_roller" }
```

---

## 🎬 State Machine + Hardware Events

### Go Code (спрощено)

```go
type ATMStateMachine struct {
    state         string
    countedNotes  int
    expectedNotes int
    eventChan     chan HardwareEvent
}

// Головний loop - слухає hardware events
func (sm *ATMStateMachine) Run() {
    for event := range sm.eventChan {
        sm.handleEvent(event)
    }
}

func (sm *ATMStateMachine) handleEvent(event HardwareEvent) {
    switch sm.state {
    case "COUNTING":
        if event.Type == "NOTE_COUNTED" {
            sm.countedNotes++
            log.Printf("Counted: %d / %d", sm.countedNotes, sm.expectedNotes)
            
            if sm.countedNotes == sm.expectedNotes {
                sm.state = "PRESENTING"
            }
        }
        
    case "PRESENTING":
        if event.Type == "NOTES_PRESENTED" {
            // ✅ Hardware підтвердив: гроші вийшли!
            log.Printf("✅ Cash dispensed!")
            sm.state = "COMPLETED"
            
            // Тепер можна CONFIRM в БД
            db.Exec("UPDATE atm_transactions SET status = 'completed'")
            db.Exec("UPDATE accounts SET balance = balance - 100")
        }
    }
}
```

---

## 🔌 Communication Protocol

### XFS (Industry Standard)

```
Software ─────> Middleware ─────> Hardware
         JSON            Binary/Serial
         
Command:
{
  "command": "DISPENSE",
  "amount": 100,
  "note_count": 5
}

Events (від hardware):
{
  "event": "NOTE_COUNTED",
  "count": 1,
  "timestamp": "2026-01-28T10:15:30Z"
}

{
  "event": "NOTES_PRESENTED",
  "total": 5,
  "timestamp": "2026-01-28T10:15:35Z"
}
```

---

## 🎯 Критичний момент

### Коли програма знає що гроші видані?

```
❌ НЕ коли:
├─ Команда відправлена
├─ Motors почали працювати
└─ Купюри пораховані

✅ ТАК коли:
└─> Event: "NOTES_PRESENTED" від hardware ✅
    └─> Exit Sensor спрацював
        └─> Купюри ФІЗИЧНО вийшли до клієнта

Тільки ПІСЛЯ цього:
└─> Database CONFIRM (списати гроші)
```

---

## ⚠️ Edge Case: Jam (Застрявання)

```
Software          Hardware
   |                  |
   | State: COUNTING  |
   |                  |
   |                  | Note 1 ✅
   |<─────────────    |
   |                  | Note 2 ✅
   |<─────────────    |
   |                  | Note 3... застрягла!
   |                  | EVENT: JAM ❌
   |<─────────────    |
   |                  [Motors STOP]
   |                  |
   | State: ERROR     |
   | REFUND в БД ✅   |
   |                  |
```

### Що відбувається

```go
if event.Type == "NOTE_JAM" {
    log.Printf("❌ Jam detected!")
    sm.state = "ERROR"
    
    // REFUND (compensating transaction)
    tx.Exec("UPDATE atm_transactions SET status = 'refunded'")
    tx.Exec("UPDATE account_holds SET status = 'released'")
    tx.Commit()
    
    // Гроші НЕ списані, balance без змін ✅
}
```

---

## 📊 State Transitions

```
State Machine синхронізується з Hardware:

Software State      Hardware Event       Action
──────────────      ──────────────       ──────
IDLE                -                    -
    ↓
DISPENSING          CMD sent             Motors start
    ↓
COUNTING            NOTE_COUNTED ✅      Count++
    ↓
PRESENTING          NOTES_PRESENTED ✅   Shutter open
    ↓
COMPLETED           -                    DB CONFIRM ✅

OR (якщо помилка):

COUNTING            NOTE_JAM ❌          Motors stop
    ↓
ERROR               -                    DB REFUND ✅
```

---

## 🎓 Висновок

### Як це працює?

```
1. Hardware має SENSORS (датчики)
   └─> Note Counter, Exit Sensor, Jam Sensor

2. Firmware генерує EVENTS з sensors
   └─> "NOTE_COUNTED", "NOTES_PRESENTED", "NOTE_JAM"

3. Middleware пересилає events через USB/Serial
   └─> Binary protocol → JSON/Events

4. State Machine обробляє events
   └─> Переходи між станами

5. Database update ТІЛЬКИ після hardware підтвердження
   └─> "NOTES_PRESENTED" → CONFIRM ✅
```

### Ключове правило

```
❌ НЕ довіряй команді (може не виконатися)
✅ Довіряй тільки EVENT від hardware!

Database CONFIRM = ПІСЛЯ "NOTES_PRESENTED" event ✅
```

---

## 📖 Детальний файл

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_7

# Повна теорія з кодом
cat theory/17_hardware_software_integration.md
```

**Обсяг:** Повна Go реалізація + protocols + debugging

---

**Event-Driven Architecture з Hardware Integration!** 🎯

**Physical World ─(Sensors)→ Events ─(Middleware)→ State Machine ─(Logic)→ Database**  
    💵💵💵           ✅          USB/Serial           FSM              💾
