# CAP Theorem

## Що таке CAP?

**CAP Theorem** (Brewer's theorem): В розподіленій системі можна гарантувати **максимум 2 з 3** властивостей одночасно.

```
       C (Consistency)
      / \
     /   \
    /     \
   /       \
  /    CA   \
 /   (RDBMS)\
/____________\
A            P
(Availability) (Partition Tolerance)

  AP                 CP
(Cassandra,      (MongoDB,
 DynamoDB)        HBase)
```

---

## 3 Властивості

### 1. Consistency (Узгодженість)
**Всі ноди бачать одні й ті ж дані в один і той же час.**

```
User → Write "X=1" → Node A
                   → Node B
                   → Node C
                   
User → Read from Node A → "X=1" ✅
User → Read from Node B → "X=1" ✅
User → Read from Node C → "X=1" ✅
```

**Приклад:** Банківський рахунок - баланс має бути однаковий на всіх серверах.

### 2. Availability (Доступність)
**Кожен запит отримує відповідь (успішну або помилку), навіть якщо деякі ноди не працюють.**

```
User → Request → Node A (down) ❌
              → Node B ✅ → Response
```

**Приклад:** Social media feed - краще показати трохи застарілі дані, ніж нічого.

### 3. Partition Tolerance (Толерантність до розділення)
**Система продовжує працювати, навіть якщо мережа між нодами розірвана.**

```
Node A ←--X--→ Node B
(Network partition)

Система все одно працює ✅
```

**Приклад:** Мікросервіси в різних датацентрах - мережа може розірватися, але сервіси мають працювати.

---

## Trade-offs

### CA (Consistency + Availability)
**НЕ ПРАЦЮЄ в розподілених системах!**

Якщо є network partition, треба вибрати: C або A.

**Приклади:** Single-server RDBMS (PostgreSQL, MySQL на одному сервері)

### CP (Consistency + Partition Tolerance)
**Жертвуємо доступністю заради узгодженості.**

```
User → Write "X=1" → Node A ✅
                   → Node B ❌ (network partition)
                   
Node A блокує запис, поки Node B не синхронізується
User → Read → "Service Unavailable" (waiting for sync)
```

**Коли використовувати:**
- Фінансові транзакції
- Інвентар товарів
- Критичні бізнес-дані

**Приклади:** MongoDB (з strong consistency), HBase, Redis (single master)

### AP (Availability + Partition Tolerance)
**Жертвуємо узгодженістю заради доступності.**

```
User → Write "X=1" → Node A ✅
                   → Node B ❌ (network partition)
                   
User → Read from Node A → "X=1" ✅
User → Read from Node B → "X=0" (old value) ❌

Eventual consistency: Node B синхронізується пізніше
```

**Коли використовувати:**
- Social media (лайки, коментарі)
- Analytics, logging
- Shopping cart (не критично, якщо дані трохи застарілі)

**Приклади:** Cassandra, DynamoDB, Riak, CouchDB

---

## Реальні Приклади

### PostgreSQL (CP)
```sql
-- Strong consistency
BEGIN;
UPDATE accounts SET balance = balance - 100 WHERE id = 1;
UPDATE accounts SET balance = balance + 100 WHERE id = 2;
COMMIT;

-- Якщо є network partition між репліками:
-- - Primary блокує write до sync з replica
-- - Read може бути недоступний
```

**Trade-off:** Availability ↓, Consistency ↑

### Cassandra (AP)
```cql
-- Eventual consistency
INSERT INTO users (id, name) VALUES (1, 'Alice');

-- Read може бути з різних нод з різними даними
SELECT * FROM users WHERE id = 1;
-- Node A: "Alice"
-- Node B: null (ще не синхронізувалося)

-- Через ~100ms всі ноди матимуть "Alice"
```

**Trade-off:** Availability ↑, Consistency ↓

### MongoDB (CP або AP - налаштовується)
```javascript
// Write Concern: "majority" (CP - wait for majority to acknowledge)
db.users.insertOne(
  { name: "Alice" },
  { writeConcern: { w: "majority" } }
);

// Read Preference: "secondary" (AP - може читати stale data)
db.users.find().readPref("secondary");
```

---

## Eventual Consistency

**Модель для AP систем:** Дані зрештою стануть узгодженими, але не одразу.

```
Timeline:
t=0s   User writes "X=1" to Node A
t=1s   User reads from Node B → "X=0" (stale)
t=2s   User reads from Node B → "X=0" (still stale)
t=5s   Node B syncs with Node A
t=6s   User reads from Node B → "X=1" ✅ (consistent)
```

**Приклади:**
- DNS (зміни DNS поширюються годинами)
- Amazon S3 (eventual consistency для overwrites)
- Social media likes/views count

---

## Як обрати?

### Обирай CP якщо:
- Фінансові дані
- Критичні транзакції
- Інвентар (stock management)
- Medical records

**Приклад:** Банк - краще показати "Service Unavailable", ніж неправильний баланс.

### Обирай AP якщо:
- Social media
- Analytics, metrics
- Logging
- Shopping cart
- User profiles

**Приклад:** Twitter - краще показати feed з трохи застарілими даними, ніж нічого.

---

## Приклад з життя: DNS

**DNS є AP система**

```
1. Ви змінюєте DNS запис: example.com → 1.2.3.4
2. Зміна поширюється через DNS сервери (propagation)
3. Деякі користувачі бачать old IP (1.2.3.3)
4. Інші користувачі бачать new IP (1.2.3.4)
5. Через 24-48 годин всі бачать new IP

Availability: ✅ (DNS завжди доступний)
Consistency: ❌ (різні користувачі бачать різні IP)
Partition Tolerance: ✅ (DNS сервери в різних країнах)
```

---

## Підсумок

| Тип | Consistency | Availability | Partition Tolerance | Приклад |
|-----|-------------|--------------|---------------------|---------|
| **CP** | ✅ | ❌ | ✅ | MongoDB (majority), HBase |
| **AP** | ❌ | ✅ | ✅ | Cassandra, DynamoDB |
| **CA** | ✅ | ✅ | ❌ | Single-server RDBMS |

**Правило:**
- Якщо система **розподілена** → Partition Tolerance обов'язково
- Залишається вибір: **C або A**

**В реальності:**
- Більшість систем - **AP з Eventual Consistency**
- Критичні системи - **CP з Strong Consistency**
- Hybrid підходи - різні consistency levels для різних даних
