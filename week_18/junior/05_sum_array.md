# Task 5: Sum of Array

**Level:** Junior  
**Time:** 5 minutes  
**Language:** Ruby, Go, or JavaScript  

---

## 📝 Task

Напиши функцію, яка знаходить суму всіх чисел у масиві.

**Twist:** Якщо в масиві є вкладені масиви, суму всіх чисел з усіх рівнів.

---

## 📥 Input

```
array (Array) - масив чисел (можливо вкладений)
```

---

## 📤 Output

```
Integer - сума всіх чисел
```

---

## 💡 Examples

```ruby
sum_array([1, 2, 3, 4, 5])
# => 15

sum_array([1, [2, 3], 4, [5, 6]])
# => 21

sum_array([])
# => 0

sum_array([10])
# => 10

sum_array([1, [2, [3, [4, [5]]]]])
# => 15
```

---

## ✅ Requirements

- Підтримай вкладені масиви (nested arrays)
- Підтримай пусті масиви
- Підтримай масиви з одним елементом
- Ігноруй нечислові значення (strings, nil, etc.)

---

## 🎯 Test Cases

```ruby
# Test 1: Simple array
input: [1, 2, 3, 4, 5]
expected: 15

# Test 2: Nested array
input: [1, [2, 3], 4, [5, 6]]
expected: 21

# Test 3: Empty array
input: []
expected: 0

# Test 4: Single element
input: [10]
expected: 10

# Test 5: Deep nesting
input: [1, [2, [3, [4, [5]]]]]
expected: 15

# Test 6: With non-numbers
input: [1, "a", 2, nil, 3]
expected: 6
```

---

**Good luck!** 🚀
