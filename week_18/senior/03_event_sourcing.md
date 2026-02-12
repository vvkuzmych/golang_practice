# Task 3: Event Sourcing System

**Level:** Senior  
**Time:** 50 minutes  
**Language:** Ruby, Go, or JavaScript  

---

## 📝 Task

Імплементуй Event Sourcing систему для Bank Account.

**Concept:** Замість зберігання current state, зберігай всі events і rebuild state з events.

**Events:**
- AccountCreated
- MoneyDeposited
- MoneyWithdrawn
- AccountClosed

---

## 📥 API Design

```ruby
# Event Store
class EventStore
  def append(aggregate_id, event)
    # Зберегти event для aggregate
  end
  
  def get_events(aggregate_id)
    # Отримати всі events для aggregate
  end
end

# Aggregate (Bank Account)
class BankAccount
  attr_reader :id, :balance, :status
  
  def initialize(id)
    @id = id
    @balance = 0
    @status = :active
    @events = []
  end
  
  # Commands (generate events)
  def create_account(owner)
    # Apply AccountCreated event
  end
  
  def deposit(amount)
    # Apply MoneyDeposited event
  end
  
  def withdraw(amount)
    # Apply MoneyWithdrawn event
    # Validate: balance >= amount
  end
  
  def close
    # Apply AccountClosed event
  end
  
  # Rebuild state from events
  def self.from_events(events)
    # Create account and replay all events
  end
  
  private
  
  def apply_event(event)
    # Apply event to change state
  end
end
```

---

## 💡 Examples

```ruby
# Create account
account = BankAccount.new("acc-123")
account.create_account(owner: "John Doe")

# Deposit
account.deposit(1000)
account.deposit(500)

# Withdraw
account.withdraw(300)

# Check balance
account.balance  # => 1200

# Save events
store = EventStore.new
account.events.each do |event|
  store.append("acc-123", event)
end

# Later... rebuild from events
events = store.get_events("acc-123")
reconstructed = BankAccount.from_events(events)
reconstructed.balance  # => 1200
```

---

## ✅ Requirements

### Functional Requirements
- Support all 4 event types
- Validate commands (can't withdraw more than balance)
- Can't operate on closed account
- Rebuild state from events
- Event immutability (events never change)

### Event Structure
```ruby
{
  event_type: "MoneyDeposited",
  aggregate_id: "acc-123",
  timestamp: Time.now,
  data: { amount: 100 },
  version: 1  # Event version for aggregate
}
```

### Non-Functional Requirements
- Events are append-only (never update/delete)
- Order matters (events must be ordered)
- Idempotency (replaying same events = same state)
- Performance: rebuilding from 1000 events < 100ms

---

## 🎯 Design Considerations

### 1. Event Store
- Як зберігати events?
- Database table? File? In-memory?
- Indexing strategy?

### 2. Versioning
- Як handle event version conflicts?
- Optimistic locking?

### 3. Snapshots
- Що якщо 100K events для одного account?
- Snapshot strategy?
- Коли створювати snapshots?

### 4. Projections
- Як query current state швидко?
- Read models vs Event store?

### 5. CQRS
- Command side vs Query side?
- Eventual consistency?

---

## 🧪 Test Cases

```ruby
# Test 1: Create and deposit
account = BankAccount.new("acc-1")
account.create_account(owner: "John")
account.deposit(1000)
assert account.balance == 1000

# Test 2: Multiple operations
account.deposit(500)
account.withdraw(300)
assert account.balance == 1200

# Test 3: Can't withdraw more than balance
assert_raises(InsufficientFundsError) do
  account.withdraw(2000)
end

# Test 4: Rebuild from events
store = EventStore.new
account.events.each { |e| store.append("acc-1", e) }

events = store.get_events("acc-1")
new_account = BankAccount.from_events(events)
assert new_account.balance == account.balance

# Test 5: Can't operate on closed account
account.close
assert_raises(AccountClosedError) do
  account.deposit(100)
end

# Test 6: Event order matters
events = [
  { type: "AccountCreated", data: { owner: "John" } },
  { type: "MoneyDeposited", data: { amount: 1000 } },
  { type: "MoneyWithdrawn", data: { amount: 300 } }
]
account = BankAccount.from_events(events)
assert account.balance == 700

# Test 7: Snapshots (bonus)
# After 100 events, create snapshot
# Rebuild from snapshot + recent events
```

---

## 💡 Bonus Points

### Advanced Features
- Snapshots (rebuild from snapshot + recent events)
- Event versioning (handle schema changes)
- Projections (read models)
- CQRS implementation
- Event replay (replay events through new logic)
- Time travel (state at any point in time)

### Production Considerations
- Event store persistence (PostgreSQL, EventStoreDB)
- Event bus (publish events to subscribers)
- Saga pattern (distributed transactions)
- Event sourcing vs CRUD trade-offs

---

## 📊 Evaluation Criteria

- **Event Design** (25%) - proper event structure, immutability
- **State Rebuild** (25%) - correct state from events
- **Validation** (20%) - proper command validation
- **Design Decisions** (20%) - storage, versioning, snapshots
- **Code Quality** (10%) - clean, testable code

---

## 🔗 Reference

- Martin Fowler: Event Sourcing
- Greg Young: CQRS & Event Sourcing
- Event Store DB
- Trade-offs: Event Sourcing vs CRUD

---

**Good luck!** 🚀
