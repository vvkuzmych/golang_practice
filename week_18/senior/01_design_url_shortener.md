# Task 1: Design URL Shortener

**Level:** Senior  
**Time:** 40 minutes  
**Language:** Ruby, Go, or JavaScript  

---

## 📝 Task

Design та імплементуй URL Shortener service (як bit.ly).

**Функціональність:**
1. Створити короткий URL з довгого
2. Redirect з короткого на довгий
3. Трекінг статистики (кількість кліків)
4. Expiration time (optional)

---

## 📥 API Design

```ruby
class URLShortener
  # Створити короткий URL
  def shorten(long_url, expire_at: nil)
    # Повертає короткий URL (string)
  end
  
  # Отримати довгий URL за коротким
  def expand(short_url)
    # Повертає довгий URL або nil якщо не знайдено/expired
  end
  
  # Отримати статистику
  def stats(short_url)
    # Повертає { clicks: N, created_at: Time, expire_at: Time }
  end
end
```

---

## 💡 Examples

```ruby
shortener = URLShortener.new

# Створити короткий URL
short = shortener.shorten("https://example.com/very/long/url")
# => "abc123"

# Expand (і збільшити counter)
shortener.expand("abc123")
# => "https://example.com/very/long/url"

# Статистика
shortener.stats("abc123")
# => { clicks: 1, created_at: ..., expire_at: nil }

# З expiration
short = shortener.shorten("https://example.com", expire_at: Time.now + 3600)
shortener.expand(short)  # => URL (якщо не expired)
sleep(3601)
shortener.expand(short)  # => nil (expired)
```

---

## ✅ Requirements

### Functional Requirements
- Генеруй унікальні короткі URLs
- Довгий URL може мати кілька коротких
- Короткий URL має бути коротким (6-8 символів)
- Підтримка expiration
- Трекінг кількості кліків

### Non-Functional Requirements
- `shorten()` - швидко (< 100ms)
- `expand()` - дуже швидко (< 10ms)
- Scale to 1M URLs
- Consider collision handling

---

## 🎯 Design Considerations

### 1. URL Generation
- Як генерувати короткі URLs?
- Hash? Random? Sequential?
- Collision handling?

### 2. Storage
- Як зберігати mapping (short → long)?
- Як зберігати stats?
- In-memory? Database? Cache?

### 3. Performance
- Як оптимізувати `expand()`?
- Індекси? Кешування?

### 4. Scalability
- Що якщо 100M URLs?
- Sharding strategy?
- Database design?

---

## 🧪 Test Cases

```ruby
# Test 1: Basic shortening
short = shortener.shorten("https://example.com/long")
expanded = shortener.expand(short)
assert expanded == "https://example.com/long"

# Test 2: Different URLs get different shorts
short1 = shortener.shorten("https://example.com/1")
short2 = shortener.shorten("https://example.com/2")
assert short1 != short2

# Test 3: Same URL can get different shorts
short1 = shortener.shorten("https://example.com")
short2 = shortener.shorten("https://example.com")
# Can be different (design choice)

# Test 4: Stats tracking
short = shortener.shorten("https://example.com")
shortener.expand(short)
shortener.expand(short)
stats = shortener.stats(short)
assert stats[:clicks] == 2

# Test 5: Expiration
short = shortener.shorten("https://example.com", 
                          expire_at: Time.now + 1)
sleep(1.1)
assert shortener.expand(short) == nil

# Test 6: Non-existent short URL
assert shortener.expand("notexist") == nil
```

---

## 💡 Bonus Points

- Custom short URLs (vanity URLs)
- QR code generation
- Analytics (clicks per day, geo location)
- Rate limiting
- A/B testing support

---

## 📊 Evaluation Criteria

- **Code Quality** (30%) - clean, readable, maintainable
- **Algorithm** (25%) - URL generation strategy
- **Design** (25%) - scalability, performance considerations
- **Edge Cases** (20%) - collision, expiration, validation

---

**Good luck!** 🚀
