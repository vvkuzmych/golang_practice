# Task 1: Reverse String

**Level:** Junior  
**Time:** 5 minutes  
**Language:** Ruby, Go, or JavaScript  

---

## 📝 Task

Напиши функцію, яка розвертає рядок.

---

## 📥 Input

```
string (String) - рядок для розвороту
```

---

## 📤 Output

```
String - розвернутий рядок
```

---

## 💡 Examples

```ruby
reverse_string("hello")
# => "olleh"

reverse_string("Ruby")
# => "ybuR"

reverse_string("12345")
# => "54321"

reverse_string("")
# => ""

reverse_string("a")
# => "a"
```

---

## ✅ Requirements

- НЕ використовуй вбудовані методи `.reverse()` (якщо мова)
- Підтримай пусті рядки
- Підтримай рядки з одним символом

---

## 🎯 Test Cases

```ruby
# Test 1: Normal string
input: "hello"
expected: "olleh"

# Test 2: Empty string
input: ""
expected: ""

# Test 3: Single character
input: "a"
expected: "a"

# Test 4: Numbers
input: "12345"
expected: "54321"

# Test 5: Special characters
input: "a!b@c#"
expected: "#c@b!a"
```

---

**Good luck!** 🚀
