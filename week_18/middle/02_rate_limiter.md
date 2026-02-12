# Task 2: Rate Limiter

**Level:** Middle  
**Time:** 20 minutes  
**Language:** Ruby, Go, or JavaScript  

---

## 📝 Task

Імплементуй Rate Limiter, який обмежує кількість requests від користувача.

**Rules:**
- Max N requests per M seconds
- Якщо ліміт перевищено, повертай `false`
- Якщо в межах ліміту, повертай `true`

---

## 📥 API

```ruby
class RateLimiter
  def initialize(max_requests, window_seconds)
    # max_requests - максимальна кількість requests
    # window_seconds - вікно часу в секундах
  end
  
  def allow_request(user_id)
    # Повертає true якщо request дозволено
    # Повертає false якщо ліміт перевищено
  end
end
```

---

## 💡 Examples

```ruby
# Max 3 requests per 10 seconds
limiter = RateLimiter.new(3, 10)

limiter.allow_request("user1")  # => true (1/3)
limiter.allow_request("user1")  # => true (2/3)
limiter.allow_request("user1")  # => true (3/3)
limiter.allow_request("user1")  # => false (exceeded)

sleep(10)  # Wait for window to reset

limiter.allow_request("user1")  # => true (1/3)
```

---

## ✅ Requirements

- Підтримуй multiple users (кожен user має свій ліміт)
- Sliding window approach (не fixed window)
- Старі requests автоматично "забуваються" після window_seconds
- Ефективність: O(1) або O(log N) per request

---

## 🎯 Test Cases

```ruby
# Test 1: Within limit
limiter = RateLimiter.new(3, 10)
limiter.allow_request("user1") # => true
limiter.allow_request("user1") # => true
limiter.allow_request("user1") # => true

# Test 2: Exceed limit
limiter.allow_request("user1") # => false

# Test 3: Different users
limiter = RateLimiter.new(2, 10)
limiter.allow_request("user1") # => true
limiter.allow_request("user2") # => true
limiter.allow_request("user1") # => true
limiter.allow_request("user2") # => true
limiter.allow_request("user1") # => false
limiter.allow_request("user2") # => false

# Test 4: Window reset
limiter = RateLimiter.new(1, 1)
limiter.allow_request("user1") # => true
limiter.allow_request("user1") # => false
sleep(1.1)
limiter.allow_request("user1") # => true
```

---

## 💡 Hints

- Зберігай timestamps requests для кожного user
- Видаляй старі timestamps (за межами window)
- Можна використати Queue або Array з timestamps

---

**Good luck!** 🚀
