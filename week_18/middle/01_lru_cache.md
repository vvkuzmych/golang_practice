# Task 1: LRU Cache

**Level:** Middle  
**Time:** 25 minutes  
**Language:** Ruby, Go, or JavaScript  

---

## 📝 Task

Імплементуй LRU (Least Recently Used) Cache з обмеженою capacity.

Cache повинен підтримувати:
- `get(key)` - отримати значення за ключем (O(1))
- `put(key, value)` - додати/оновити значення (O(1))

Коли capacity перевищено, видаляй найменш недавно використаний елемент.

---

## 📥 API

```ruby
class LRUCache
  def initialize(capacity)
    # ...
  end
  
  def get(key)
    # Повертає value або nil якщо не знайдено
  end
  
  def put(key, value)
    # Додає або оновлює значення
  end
end
```

---

## 💡 Examples

```ruby
cache = LRUCache.new(2)

cache.put(1, "one")
cache.put(2, "two")
cache.get(1)          # => "one"

cache.put(3, "three") # Видаляє key=2 (least recently used)

cache.get(2)          # => nil (видалено)
cache.get(3)          # => "three"

cache.put(4, "four")  # Видаляє key=1

cache.get(1)          # => nil
cache.get(3)          # => "three"
cache.get(4)          # => "four"
```

---

## ✅ Requirements

- `get(key)` повинен бути O(1)
- `put(key, value)` повинен бути O(1)
- Коли capacity перевищено, видаляй LRU елемент
- `get` оновлює "recently used" статус
- `put` теж оновлює "recently used" статус

---

## 🎯 Test Cases

```ruby
# Test 1: Basic usage
cache = LRUCache.new(2)
cache.put(1, "one")
cache.put(2, "two")
cache.get(1) # => "one"
cache.put(3, "three") # removes 2
cache.get(2) # => nil

# Test 2: Update existing key
cache = LRUCache.new(2)
cache.put(1, "one")
cache.put(2, "two")
cache.put(1, "ONE") # update
cache.get(1) # => "ONE"

# Test 3: Get updates recency
cache = LRUCache.new(2)
cache.put(1, "one")
cache.put(2, "two")
cache.get(1) # makes 1 recently used
cache.put(3, "three") # removes 2 (not 1)
cache.get(2) # => nil
cache.get(1) # => "one"
```

---

## 💡 Hints

- Можеш використати HashMap + Doubly Linked List
- HashMap для O(1) lookup
- Linked List для LRU tracking
- Або знайди інший підхід!

---

**Good luck!** 🚀
