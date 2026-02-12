# Task 4: Merge Overlapping Intervals

**Level:** Middle  
**Time:** 20 minutes  
**Language:** Ruby, Go, or JavaScript  

---

## 📝 Task

Дано масив intervals (початок, кінець). Зmerge всі overlapping intervals.

---

## 📥 Input

```
intervals (Array of Arrays) - [[start, end], [start, end], ...]
```

---

## 📤 Output

```
Array of Arrays - merged intervals
```

---

## 💡 Examples

```ruby
merge_intervals([[1, 3], [2, 6], [8, 10], [15, 18]])
# => [[1, 6], [8, 10], [15, 18]]

merge_intervals([[1, 4], [4, 5]])
# => [[1, 5]]

merge_intervals([[1, 4], [0, 4]])
# => [[0, 4]]

merge_intervals([[1, 4]])
# => [[1, 4]]

merge_intervals([])
# => []
```

---

## ✅ Requirements

- Intervals які торкаються (touch) теж merge [[1, 3], [3, 5]] => [[1, 5]]
- Результат має бути відсортованим
- Підтримай пусті масиви
- Підтримай масив з одним interval

---

## 🎯 Test Cases

```ruby
# Test 1: Overlapping intervals
input: [[1, 3], [2, 6], [8, 10], [15, 18]]
expected: [[1, 6], [8, 10], [15, 18]]

# Test 2: Touching intervals
input: [[1, 4], [4, 5]]
expected: [[1, 5]]

# Test 3: Completely overlapping
input: [[1, 4], [2, 3]]
expected: [[1, 4]]

# Test 4: No overlap
input: [[1, 2], [3, 4], [5, 6]]
expected: [[1, 2], [3, 4], [5, 6]]

# Test 5: Unsorted input
input: [[8, 10], [1, 3], [2, 6]]
expected: [[1, 6], [8, 10]]

# Test 6: All merge into one
input: [[1, 4], [2, 5], [3, 6]]
expected: [[1, 6]]
```

---

## 💡 Hints

- Sort intervals by start time
- Iterate and merge when overlap detected
- Two intervals overlap if: start2 <= end1

---

**Good luck!** 🚀
