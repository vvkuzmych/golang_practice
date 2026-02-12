# Task 3: FizzBuzz

**Level:** Junior  
**Time:** 5 minutes  
**Language:** Ruby, Go, or JavaScript  

---

## 📝 Task

Класичне завдання FizzBuzz.

Напиши функцію, яка для чисел від 1 до N:
- Повертає "Fizz" якщо число ділиться на 3
- Повертає "Buzz" якщо число ділиться на 5
- Повертає "FizzBuzz" якщо число ділиться на 3 І 5
- Інакше повертає саме число

---

## 📥 Input

```
n (Integer) - максимальне число
```

---

## 📤 Output

```
Array - масив рядків/чисел
```

---

## 💡 Examples

```ruby
fizzbuzz(15)
# => [1, 2, "Fizz", 4, "Buzz", "Fizz", 7, 8, "Fizz", "Buzz", 11, "Fizz", 13, 14, "FizzBuzz"]

fizzbuzz(5)
# => [1, 2, "Fizz", 4, "Buzz"]

fizzbuzz(1)
# => [1]

fizzbuzz(0)
# => []
```

---

## ✅ Requirements

- Числа від 1 до N (inclusive)
- Якщо N = 0, повернути пустий масив
- Порядок має бути правильним

---

## 🎯 Test Cases

```ruby
# Test 1: Normal case
input: 15
expected: [1, 2, "Fizz", 4, "Buzz", "Fizz", 7, 8, "Fizz", "Buzz", 11, "Fizz", 13, 14, "FizzBuzz"]

# Test 2: Only 5
input: 5
expected: [1, 2, "Fizz", 4, "Buzz"]

# Test 3: Only 1
input: 1
expected: [1]

# Test 4: Zero
input: 0
expected: []

# Test 5: Up to 3
input: 3
expected: [1, 2, "Fizz"]
```

---

**Good luck!** 🚀
