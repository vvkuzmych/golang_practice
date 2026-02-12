# Task 1: Parallel Sum Calculator

**Level:** Beginner  
**Time:** 10 minutes  
**Topics:** Goroutines, WaitGroup, Mutex

---

## 📝 Task

Напиши функцію, яка розраховує суму великого slice чисел **паралельно**.

Розділи slice на N частин і обробляй кожну частину в окремій goroutine.

---

## 📥 Function Signature

```go
func ParallelSum(numbers []int, workers int) int
```

**Parameters:**
- `numbers` - slice чисел для підсумовування
- `workers` - кількість goroutines (наприклад, 4)

**Returns:**
- `int` - сума всіх чисел

---

## 💡 Examples

```go
numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
result := ParallelSum(numbers, 2)
// => 55 (1+2+3+4+5+6+7+8+9+10)

numbers := []int{10, 20, 30, 40}
result := ParallelSum(numbers, 4)
// => 100

numbers := []int{1, 2, 3}
result := ParallelSum(numbers, 10)  // більше workers ніж елементів
// => 6
```

---

## ✅ Requirements

- Використай `sync.WaitGroup` для очікування всіх goroutines
- Використай `sync.Mutex` для безпечного додавання до загальної суми
- Розділи slice на приблизно рівні частини
- Кожен worker обробляє свою частину
- Підтримай випадок коли `workers > len(numbers)`

---

## 🧪 Test Cases

```go
// Test 1: Simple sum
numbers := []int{1, 2, 3, 4, 5}
result := ParallelSum(numbers, 2)
assert.Equal(t, 15, result)

// Test 2: Single worker
numbers := []int{10, 20, 30}
result := ParallelSum(numbers, 1)
assert.Equal(t, 60, result)

// Test 3: More workers than elements
numbers := []int{1, 2, 3}
result := ParallelSum(numbers, 10)
assert.Equal(t, 6, result)

// Test 4: Large array
numbers := make([]int, 1000000)
for i := range numbers {
    numbers[i] = 1
}
result := ParallelSum(numbers, 4)
assert.Equal(t, 1000000, result)

// Test 5: Empty array
numbers := []int{}
result := ParallelSum(numbers, 2)
assert.Equal(t, 0, result)
```

---

## 💡 Hints

1. Розрахуй розмір chunk для кожного worker: `chunkSize = len(numbers) / workers`
2. Останній worker може мати трохи більше елементів
3. Кожен worker додає свою частину до загальної суми (потрібен mutex!)
4. Використай `wg.Add(workers)` та `wg.Done()` для синхронізації

---

**Рішення:** `solutions/solution_01_parallel_sum.go`

**Good luck!** 🚀
