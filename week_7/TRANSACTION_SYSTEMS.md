# 📊 Які бувають системи транзакцій?

## 🎯 8 Основних типів

```
1. ACID (Локальні)              ⭐⭐⭐⭐⭐
2. Two-Phase Commit (2PC)        ⭐
3. Saga Pattern                  ⭐⭐⭐⭐⭐
4. Event Sourcing                ⭐⭐⭐
5. CQRS                          ⭐⭐⭐
6. Outbox Pattern                ⭐⭐⭐⭐⭐
7. Try-Confirm/Cancel (TCC)      ⭐⭐
8. Reserve/Confirm               ⭐⭐⭐⭐
```

---

## 1️⃣ ACID (Локальні Транзакції)

**Одна база даних**

```go
tx, _ := db.Begin()
tx.Exec("UPDATE accounts SET balance = balance - 100")
tx.Exec("INSERT INTO transactions ...")
tx.Commit() // Або всі, або жодна!
```

✅ **Use case:** Монолітні системи, одна БД  
✅ **Consistency:** Immediate  
✅ **Complexity:** Low  

---

## 2️⃣ Two-Phase Commit (2PC)

**Кілька баз даних з координатором**

```
Phase 1: PREPARE всіх
Phase 2: COMMIT всіх або ROLLBACK всіх
```

⚠️ **Use case:** Legacy distributed databases  
✅ **Consistency:** Immediate  
❌ **Problem:** Blocking, single point of failure  

---

## 3️⃣ Saga Pattern ⭐ (Найпопулярніше для мікросервісів)

**Послідовність локальних транзакцій + compensations**

```
Service A: Debit ✅ → Event
         ↓
Service B: Process ✅ → Event
         ↓
Service C: Notify ✅

Якщо B fails:
└─> Compensate A (Refund)
```

✅ **Use case:** Мікросервіси, event-driven  
⚠️ **Consistency:** Eventual  
✅ **Scalability:** High  

---

## 4️⃣ Event Sourcing

**Зберігати всі зміни як події**

```
Event Store:
├─ Event 1: AccountCreated
├─ Event 2: MoneyDeposited $100
├─ Event 3: MoneyWithdrawn $50
└─ Current State = Replay events = $50
```

✅ **Use case:** Audit trail, фінанси, compliance  
✅ **Features:** Time travel, повна історія  
❌ **Complexity:** High  

---

## 5️⃣ CQRS

**Різні моделі для read і write**

```
Write Side: PostgreSQL (normalized)
         ↓
    Event Bus
         ↓
Read Side: ElasticSearch (denormalized)
```

✅ **Use case:** High read/write різниця  
✅ **Performance:** Optimize reads окремо  
❌ **Complexity:** High  

---

## 6️⃣ Outbox Pattern ⭐ (З message queues)

**Гарантія delivery через БД**

```
BEGIN TRANSACTION
├─ UPDATE accounts
├─ INSERT INTO outbox (event)
└─ COMMIT (atomic!)
         ↓
Background Worker
└─> Publish to queue ✅
```

✅ **Use case:** Microservices + Kafka/RabbitMQ  
✅ **Guarantee:** At-least-once delivery  
✅ **Reliability:** No lost messages  

---

## 7️⃣ Try-Confirm/Cancel (TCC)

**Трифазний протокол**

```
1. TRY (reserve)
2. CONFIRM (finalize) або CANCEL (compensate)
```

⚠️ **Use case:** Фінанси, booking systems  
✅ **Consistency:** Immediate-ish  
❌ **Complexity:** High (3 endpoints per service)  

---

## 8️⃣ Reserve/Confirm ⭐ (External Systems)

**Для ATM, payment gateway, тощо**

```
1. RESERVE (hold, не списувати)
2. TRY external system
3a. SUCCESS → CONFIRM (deduct)
3b. FAILURE → REFUND (release hold)
```

✅ **Use case:** ATM, payment gateway, shipping  
✅ **Safety:** No money lost  
✅ **Reconciliation:** Possible  

---

## 📊 Порівняльна таблиця

| Type | Consistency | Complexity | Use Case |
|------|-------------|------------|----------|
| ACID | Immediate ✅ | Low ⭐ | Single DB |
| 2PC | Immediate ✅ | High ⭐⭐⭐ | Legacy |
| **Saga** | Eventual ⚠️ | Medium ⭐⭐ | **Microservices** ⭐ |
| Event Sourcing | Eventual ⚠️ | High ⭐⭐⭐ | Audit |
| CQRS | Eventual ⚠️ | High ⭐⭐⭐ | Read/Write diff |
| **Outbox** | Eventual ⚠️ | Medium ⭐⭐ | **With queues** ⭐ |
| TCC | Immediate ✅ | High ⭐⭐⭐ | Finance |
| **Reserve/Confirm** | Eventual ⚠️ | Medium ⭐⭐ | **External** ⭐ |

---

## 🎯 Як вибрати?

### Одна БД?
```
✅ ACID (просто і надійно)
```

### Мікросервіси?
```
✅ Saga Pattern (choreography)
✅ Outbox Pattern (з Kafka/RabbitMQ)
```

### Зовнішні системи (ATM, API)?
```
✅ Reserve/Confirm Pattern
```

### Audit trail критичний?
```
✅ Event Sourcing
```

### High read/write різниця?
```
✅ CQRS
```

---

## 💡 Найпопулярніші в 2026

### 1. **Saga Pattern** ⭐⭐⭐⭐⭐
```
Мікросервіси, event-driven
Eventual consistency OK
Scalable, no blocking
```

### 2. **Outbox Pattern** ⭐⭐⭐⭐⭐
```
З Kafka, RabbitMQ, SQS
At-least-once delivery
Reliable messaging
```

### 3. **ACID** ⭐⭐⭐⭐
```
Монолітні системи
Immediate consistency
Simple & reliable
```

### 4. **Reserve/Confirm** ⭐⭐⭐⭐
```
External systems
ATM, payment gateway
Safe for hardware
```

---

## 🔥 Real-World Combinations

### E-commerce
```
✅ Saga (order workflow)
✅ Outbox (notifications)
✅ Reserve/Confirm (payment, warehouse)
✅ CQRS (product catalog)
```

### Banking
```
✅ Event Sourcing (audit)
✅ ACID (core transactions)
✅ Reserve/Confirm (ATM)
✅ TCC (межбанк)
```

### Social Network
```
✅ ACID (posts, comments)
✅ Saga (notifications)
✅ CQRS (news feed)
```

---

## 📖 Детальні файли

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_7/theory

# Огляд всіх типів
cat 16_transaction_systems_overview.md

# Конкретні типи:
cat 14_acid_transactions.md           # ACID
cat 12_distributed_transactions.md    # Saga, 2PC, Outbox
cat 15_external_systems_transactions.md # Reserve/Confirm
```

---

## 🎓 Висновок

**Не існує "найкращого" - залежить від:**

1. Архітектури (монолітна / мікросервіси)
2. Consistency needs (immediate / eventual)
3. Performance requirements
4. Business domain
5. Team expertise

**Розумійте всі, вибирайте правильний!** 🎯

---

**Файл:** `theory/16_transaction_systems_overview.md`  
**Обсяг:** Повний огляд + приклади коду для всіх 8 типів
