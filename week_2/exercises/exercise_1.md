# Вправа 1: Calculator з методами

## Ціль
Створити калькулятор з використанням struct та methods (замість простих функцій).

---

## Завдання

Створіть програму `calculator.go`, яка:

1. Має struct `Calculator` з полем для зберігання результату
2. Використовує методи замість функцій
3. Підтримує базові операції: `+`, `-`, `*`, `/`
4. Зберігає результат попередньої операції
5. Має метод для скидання результату

---

## Вимоги

### Struct
```go
type Calculator struct {
    result float64
}
```

### Обов'язкові методи:
- `Add(value float64)` - додавання
- `Subtract(value float64)` - віднімання
- `Multiply(value float64)` - множення
- `Divide(value float64) error` - ділення (з перевіркою на 0)
- `Result() float64` - отримати результат
- `Reset()` - скинути результат до 0
- `String() string` - текстове представлення

---

## Приклад використання

```go
func main() {
    calc := Calculator{}
    
    calc.Add(10)
    fmt.Println(calc.Result())  // 10
    
    calc.Multiply(5)
    fmt.Println(calc.Result())  // 50
    
    calc.Subtract(20)
    fmt.Println(calc.Result())  // 30
    
    calc.Divide(3)
    fmt.Println(calc.Result())  // 10
    
    calc.Reset()
    fmt.Println(calc.Result())  // 0
}
```

---

## Підказки

### 1. Pointer Receiver
```go
// Використовуйте pointer receiver для зміни result
func (c *Calculator) Add(value float64) {
    c.result += value
}
```

### 2. Обробка помилок
```go
func (c *Calculator) Divide(value float64) error {
    if value == 0 {
        return errors.New("division by zero")
    }
    c.result /= value
    return nil
}
```

### 3. String метод
```go
func (c Calculator) String() string {
    return fmt.Sprintf("Calculator: %.2f", c.result)
}
```

---

## Очікуваний вивід

```
Початкове значення: 0.00

10 + 5 = 15.00
15 * 2 = 30.00
30 - 10 = 20.00
20 / 4 = 5.00

Результат: 5.00

Спробаділення на 0:
Error: division by zero
Результат після помилки: 5.00

Скидання: 0.00
```

---

## Бонус завдання

1. **History**: Зберігати історію операцій
   ```go
   type Calculator struct {
       result  float64
       history []string
   }
   
   func (c *Calculator) History() []string {
       return c.history
   }
   ```

2. **Square Root**: Додати метод `Sqrt()`
   ```go
   func (c *Calculator) Sqrt() error {
       if c.result < 0 {
           return errors.New("cannot take square root of negative number")
       }
       c.result = math.Sqrt(c.result)
       return nil
   }
   ```

3. **Power**: Додати метод `Power(exp float64)`
   ```go
   func (c *Calculator) Power(exp float64) {
       c.result = math.Pow(c.result, exp)
   }
   ```

4. **Chainable**: Зробити методи chainable
   ```go
   func (c *Calculator) Add(value float64) *Calculator {
       c.result += value
       return c
   }
   
   // Використання:
   calc.Add(10).Multiply(2).Subtract(5)
   ```

---

## Критерії оцінки

- ✅ Програма компілюється без помилок
- ✅ Використані pointer receivers для методів, що змінюють дані
- ✅ Коректно обробляється ділення на нуль
- ✅ Всі базові операції працюють правильно
- ✅ Метод `Reset()` скидає результат
- ✅ Код чистий і зрозумілий

---

## Рішення

Рішення знаходиться в `solutions/solution_1.go`.

Спробуйте виконати завдання самостійно перед тим, як дивитись рішення!

---

## Навчальні цілі

Після виконання цієї вправи ви будете вміти:
- Створювати methods на struct
- Використовувати pointer receivers
- Обробляти помилки в методах
- Реалізовувати `String()` метод
- Зберігати стан в struct

---

## Подальше вдосконалення

Подумайте як додати:
- Підтримку різних систем числення (binary, hex)
- Тригонометричні функції
- Константи (π, e)
- Memory операції (M+, M-, MR, MC)
- Undo/Redo функціонал

