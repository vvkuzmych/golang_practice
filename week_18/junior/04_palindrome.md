# Task 4: Palindrome Check

**Level:** Junior  
**Time:** 5 minutes  
**Language:** Ruby, Go, or JavaScript  

---

## 📝 Task

Напиши функцію, яка перевіряє, чи є рядок palindrome (читається однаково з обох сторін).

---

## 📥 Input

```
string (String) - рядок для перевірки
```

---

## 📤 Output

```
Boolean - true якщо palindrome, false якщо ні
```

---

## 💡 Examples

```ruby
is_palindrome("racecar")
# => true

is_palindrome("hello")
# => false

is_palindrome("a")
# => true

is_palindrome("")
# => true

is_palindrome("A man a plan a canal Panama")
# => true (ignore spaces and case)

is_palindrome("race a car")
# => false
```

---

## ✅ Requirements

- Ignore spaces (пробіли)
- Ignore case (регістр)
- Ignore punctuation (розділові знаки)
- Пусті рядки вважаються palindrome
- Один символ - palindrome

---

## 🎯 Test Cases

```ruby
# Test 1: Simple palindrome
input: "racecar"
expected: true

# Test 2: Not palindrome
input: "hello"
expected: false

# Test 3: Single character
input: "a"
expected: true

# Test 4: Empty string
input: ""
expected: true

# Test 5: With spaces and capitals
input: "A man a plan a canal Panama"
expected: true

# Test 6: Not palindrome with spaces
input: "race a car"
expected: false

# Test 7: Numbers
input: "12321"
expected: true
```

---

**Good luck!** 🚀
