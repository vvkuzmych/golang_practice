# Exercise 1: Pipeline with Goroutines

## 🎯 Мета

Створити **pipeline** з goroutines для обробки даних через кілька етапів:

```
Generator → Processor → Consumer
```

---

## 📋 Завдання

Реалізуйте pipeline з трьох етапів:

1. **Generator** - генерує числа від 1 до 20
2. **Processor** - обробляє числа (множить на 2 і додає 1)
3. **Consumer** - друкує результати

### Вимоги:

- ✅ Кожен етап - окрема goroutine
- ✅ Використати **unbuffered channels** для комунікації
- ✅ Коректно закрити channels після завершення роботи
- ✅ Generator закриває свій output channel
- ✅ Processor закриває свій output channel
- ✅ Consumer читає до закриття channel

---

## 💡 Підказки

### Структура Pipeline:

```go
// Generator: генерує числа
func generator() <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)  // ✅ Закриваємо після відправки всіх даних
        for i := 1; i <= 20; i++ {
            out <- i
        }
    }()
    return out
}

// Processor: обробляє числа
func processor(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)  // ✅ Закриваємо після обробки всіх даних
        for n := range in {
            // Обробка: n * 2 + 1
            out <- n*2 + 1
        }
    }()
    return out
}

// Consumer: друкує результати
func consumer(in <-chan int) {
    for result := range in {
        fmt.Printf("Result: %d\n", result)
    }
}
```

### Використання:

```go
func main() {
    // Pipeline: generator → processor → consumer
    numbers := generator()
    processed := processor(numbers)
    consumer(processed)
}
```

---

## 🎓 Ключові концепції

1. **Unbuffered channels** - синхронізують goroutines
2. **Channel closure** - сигналізує "no more data"
3. **Range over channel** - автоматично завершується після close()
4. **Unidirectional channels** (`<-chan`, `chan<-`) - type safety

---

## ✅ Критерії успіху

- [ ] Generator генерує числа від 1 до 20
- [ ] Processor обробляє: `result = n * 2 + 1`
- [ ] Consumer друкує всі результати
- [ ] Channels коректно закриті (без deadlock!)
- [ ] Використано unbuffered channels
- [ ] Unidirectional channels для type safety

---

## 🚀 Очікуваний результат

```
Result: 3
Result: 5
Result: 7
Result: 9
...
Result: 41
```

(20 результатів, від 3 до 41)

---

## 🔥 Бонус (опціонально)

### Бонус 1: Додайте фільтр
Додайте етап **Filter** між Processor та Consumer, який пропускає тільки парні числа:

```
Generator → Processor → Filter → Consumer
```

### Бонус 2: Fan-Out Pattern
Додайте **2 Processors** які паралельно обробляють дані з Generator:

```
                  ┌─→ Processor 1 ─┐
Generator ─────→ │                  ├─→ Consumer
                  └─→ Processor 2 ─┘
```

**Підказка:** Використайте `merge()` функцію для об'єднання результатів.

### Бонус 3: Додайте Context
Додайте `context.Context` для можливості скасування pipeline:

```go
func generator(ctx context.Context) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for i := 1; i <= 20; i++ {
            select {
            case out <- i:
            case <-ctx.Done():
                return  // Скасовано!
            }
        }
    }()
    return out
}
```

---

## 📚 Корисні посилання

- Theory: `week_5/theory/02_channels.md` - про channels
- Practice: `week_5/practice/channel_patterns/main.go` - приклад pipeline
- Solution: `week_5/solutions/solution_1.go` (після виконання)

---

**Удачі! 🎉**

**Час виконання:** 30-45 хвилин
