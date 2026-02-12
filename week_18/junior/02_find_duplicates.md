# Task 2: Find Duplicates

**Level:** Junior  
**Time:** 7 minutes  
**Language:** Ruby, Go, or JavaScript  

---

## 📝 Task

Напиши функцію, яка знаходить всі дублікати в масиві.

---

## 📥 Input

```
array (Array of Integers) - масив чисел
```

---

## 📤 Output

```
Array of Integers - масив дублікатів (унікальні значення)
```

---

## 💡 Examples

```ruby
find_duplicates([1, 2, 3, 2, 4, 5, 1])
# => [1, 2]

find_duplicates([1, 2, 3, 4, 5])
# => []

find_duplicates([1, 1, 1, 1])
# => [1]

find_duplicates([])
# => []

find_duplicates([5, 5, 3, 3, 1, 1])
# => [1, 3, 5]  # або в будь-якому порядку
```

---

## ✅ Requirements

- Поверни тільки унікальні дублікати (без повторів)
- Порядок не важливий
- Підтримай пусті масиви
- Підтримай масиви без дублікатів

---

## 🎯 Test Cases

```ruby
# Test 1: Mix of duplicates and uniques
input: [1, 2, 3, 2, 4, 5, 1]
expected: [1, 2] (or [2, 1])

# Test 2: No duplicates
input: [1, 2, 3, 4, 5]
expected: []

# Test 3: All same
input: [1, 1, 1, 1]
expected: [1]

# Test 4: Empty array
input: []
expected: []

# Test 5: Multiple duplicates
input: [5, 5, 3, 3, 1, 1]
expected: [1, 3, 5] (any order)
```

---

**Good luck!** 🚀
