# Task 2: Distributed Lock Manager

**Level:** Senior  
**Time:** 45 minutes  
**Language:** Ruby, Go, or JavaScript  

---

## 📝 Task

Імплементуй Distributed Lock Manager для координації між multiple processes/servers.

**Scenario:** У тебе кілька серверів, які хочуть виконати критичну операцію (наприклад, обробити payment). Тільки один сервер повинен її виконати одночасно.

---

## 📥 API Design

```ruby
class DistributedLock
  def initialize(backend)
    # backend - Redis, Database, або in-memory (для тестів)
  end
  
  # Acquire lock
  def acquire(resource_id, ttl_seconds: 30)
    # Повертає true якщо lock отримано
    # Повертає false якщо resource вже locked
  end
  
  # Release lock
  def release(resource_id)
    # Повертає true якщо lock released
    # Повертає false якщо lock не існував або expired
  end
  
  # Check if locked
  def locked?(resource_id)
    # Повертає true/false
  end
  
  # Execute with lock
  def with_lock(resource_id, ttl: 30, &block)
    # Acquire lock, execute block, release lock
    # Raise error якщо не вдалося отримати lock
  end
end
```

---

## 💡 Examples

```ruby
lock = DistributedLock.new(backend)

# Acquire lock
lock.acquire("payment:123", ttl_seconds: 30)  # => true
lock.acquire("payment:123")  # => false (already locked)

# Release lock
lock.release("payment:123")  # => true
lock.acquire("payment:123")  # => true (now available)

# With block
lock.with_lock("payment:123", ttl: 30) do
  # Critical section - only one process executes this
  process_payment(123)
end

# TTL expiration
lock.acquire("payment:123", ttl_seconds: 1)
sleep(1.1)
lock.acquire("payment:123")  # => true (previous lock expired)
```

---

## ✅ Requirements

### Functional Requirements
- Support multiple resources (different resource_ids)
- TTL (Time To Live) - автоматичний release після N секунд
- Atomic operations (acquire має бути atomic)
- Prevent deadlocks
- Support re-entrancy (optional, але обговори trade-offs)

### Non-Functional Requirements
- Race condition safe
- Fast acquire/release (< 10ms)
- Low overhead
- Network failure tolerant

---

## 🎯 Design Considerations

### 1. Storage Backend
- In-memory (для тестів)
- Redis (SET NX EX)
- Database (row-level locks)
- Trade-offs кожного підходу?

### 2. Lock Ownership
- Як ідентифікувати хто володіє lock?
- UUID? Process ID? Server ID?

### 3. TTL Implementation
- Як expired locks видаляються?
- Background cleanup?
- Lazy deletion?

### 4. Edge Cases
- Process crashes while holding lock?
- Network partition?
- Clock skew між серверами?

### 5. Deadlock Prevention
- Як запобігти deadlocks?
- Timeout strategies?

---

## 🧪 Test Cases

```ruby
# Test 1: Basic acquire/release
lock.acquire("res1") # => true
lock.acquire("res1") # => false
lock.release("res1") # => true
lock.acquire("res1") # => true

# Test 2: Multiple resources
lock.acquire("res1") # => true
lock.acquire("res2") # => true (different resource)

# Test 3: TTL expiration
lock.acquire("res1", ttl_seconds: 1)
lock.locked?("res1") # => true
sleep(1.1)
lock.locked?("res1") # => false

# Test 4: with_lock helper
result = lock.with_lock("res1") do
  "success"
end
assert result == "success"

# Test 5: with_lock when locked
lock.acquire("res1")
assert_raises(LockError) do
  lock.with_lock("res1") { "fail" }
end

# Test 6: Concurrent access (simulate)
threads = []
counter = 0

10.times do
  threads << Thread.new do
    lock.with_lock("counter") do
      temp = counter
      sleep(0.01) # Simulate work
      counter = temp + 1
    end
  end
end

threads.each(&:join)
assert counter == 10  # No race condition!
```

---

## 💡 Bonus Points

### Advanced Features
- Blocking acquire (wait до отримання lock)
- Lock renewal (extend TTL)
- Fairness (FIFO queue)
- Watch/notify when lock released
- Distributed consensus (Raft/Paxos)

### Monitoring
- Lock metrics (hold time, wait time)
- Deadlock detection
- Alert on long-held locks

---

## 📊 Evaluation Criteria

- **Correctness** (35%) - race condition safe, atomic operations
- **Design** (30%) - architecture, backend choice, trade-offs
- **Edge Cases** (20%) - crashes, network failures, TTL
- **Performance** (15%) - overhead, scalability

---

## 🔗 Reference

- Redis SETNX: https://redis.io/commands/setnx
- Redlock algorithm: https://redis.io/topics/distlock
- Database row-level locks
- Compare approaches!

---

**Good luck!** 🚀
