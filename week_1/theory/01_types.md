# Типи даних в Go

## Основні типи

### 1. Числові типи (int)

```go
package main

import "fmt"

func main() {
    // Цілі числа
    var age int = 25
    var temperature int32 = -10
    var population int64 = 1000000
    
    fmt.Printf("Вік: %d (тип: %T)\n", age, age)
    fmt.Printf("Температура: %d (тип: %T)\n", temperature, temperature)
    fmt.Printf("Населення: %d (тип: %T)\n", population, population)
}
```

**Типи int:**
- `int` - залежить від архітектури (32 або 64 біти)
- `int8` - від -128 до 127
- `int16` - від -32768 до 32767
- `int32` - від -2147483648 до 2147483647
- `int64` - від -9223372036854775808 до 9223372036854775807

**Беззнакові (unsigned):**
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`

---

### 2. Дійсні числа (float)

```go
package main

import "fmt"

func main() {
    var price float32 = 19.99
    var distance float64 = 123.456789
    
    fmt.Printf("Ціна: %.2f грн\n", price)
    fmt.Printf("Відстань: %.2f км\n", distance)
}
```

**Типи float:**
- `float32` - 32-бітне дійсне число
- `float64` - 64-бітне дійсне число (рекомендовано)

---

### 3. Рядки (string)

```go
package main

import "fmt"

func main() {
    var name string = "Іван"
    var greeting string = "Привіт, світ!"
    
    // Багаторядковий рядок
    var multiline string = `
    Це багаторядковий
    рядок в Go
    `
    
    fmt.Printf("Ім'я: %s\n", name)
    fmt.Printf("Привітання: %s\n", greeting)
    fmt.Printf("Багаторядковий: %s\n", multiline)
    
    // Довжина рядка
    fmt.Printf("Довжина імені: %d\n", len(name))
}
```

**Особливості string:**
- Незмінні (immutable)
- UTF-8 кодування
- Можна конкатенувати з `+`

---

### 4. Логічний тип (bool)

```go
package main

import "fmt"

func main() {
    var isStudent bool = true
    var hasJob bool = false
    
    fmt.Printf("Студент? %t\n", isStudent)
    fmt.Printf("Має роботу? %t\n", hasJob)
    
    // Логічні операції
    fmt.Printf("Студент І має роботу: %t\n", isStudent && hasJob)
    fmt.Printf("Студент АБО має роботу: %t\n", isStudent || hasJob)
    fmt.Printf("НЕ студент: %t\n", !isStudent)
}
```

**Значення bool:**
- `true` - істина
- `false` - хиба

---

### 5. Структури (struct)

```go
package main

import "fmt"

// Оголошення структури
type Person struct {
    FirstName string
    LastName  string
    Age       int
    IsStudent bool
}

func main() {
    // Створення екземпляра структури
    
    // Спосіб 1: з іменами полів
    person1 := Person{
        FirstName: "Іван",
        LastName:  "Петренко",
        Age:       25,
        IsStudent: true,
    }
    
    // Спосіб 2: за порядком полів
    person2 := Person{"Марія", "Іванова", 22, true}
    
    // Спосіб 3: частково
    person3 := Person{
        FirstName: "Петро",
        Age:       30,
    }
    
    fmt.Printf("Персона 1: %+v\n", person1)
    fmt.Printf("Персона 2: %+v\n", person2)
    fmt.Printf("Персона 3: %+v\n", person3)
    
    // Доступ до полів
    fmt.Printf("Ім'я: %s %s\n", person1.FirstName, person1.LastName)
    fmt.Printf("Вік: %d\n", person1.Age)
}
```

**Використання struct:**
- Групування пов'язаних даних
- Створення власних типів
- Методи на структурах

---

### 6. Масиви (array)

```go
package main

import "fmt"

func main() {
    // Фіксована довжина
    var numbers [5]int
    numbers[0] = 10
    numbers[1] = 20
    
    // Ініціалізація
    fruits := [3]string{"яблуко", "банан", "апельсин"}
    
    // Авто-довжина
    colors := [...]string{"червоний", "зелений", "синій"}
    
    fmt.Printf("Числа: %v\n", numbers)
    fmt.Printf("Фрукти: %v\n", fruits)
    fmt.Printf("Кольори: %v\n", colors)
    fmt.Printf("Довжина: %d\n", len(fruits))
}
```

---

### 7. Зрізи (slice)

```go
package main

import "fmt"

func main() {
    // Динамічна довжина
    var numbers []int
    numbers = append(numbers, 10, 20, 30)
    
    // Ініціалізація
    fruits := []string{"яблуко", "банан"}
    fruits = append(fruits, "апельсин")
    
    // Make
    colors := make([]string, 0, 5) // len=0, cap=5
    colors = append(colors, "червоний")
    
    fmt.Printf("Числа: %v\n", numbers)
    fmt.Printf("Фрукти: %v\n", fruits)
    fmt.Printf("Кольори: %v (len=%d, cap=%d)\n", colors, len(colors), cap(colors))
    
    // Slice операції
    slice := []int{1, 2, 3, 4, 5}
    fmt.Printf("Перші 3: %v\n", slice[:3])
    fmt.Printf("Останні 2: %v\n", slice[3:])
    fmt.Printf("Середина: %v\n", slice[1:4])
}
```

---

### 8. Мапи (map)

```go
package main

import "fmt"

func main() {
    // Створення мапи
    ages := make(map[string]int)
    ages["Іван"] = 25
    ages["Марія"] = 22
    
    // Літерал мапи
    prices := map[string]float64{
        "хліб":   15.50,
        "молоко": 28.00,
        "яйця":   45.00,
    }
    
    fmt.Printf("Вік Івана: %d\n", ages["Іван"])
    fmt.Printf("Ціна хліба: %.2f грн\n", prices["хліб"])
    
    // Перевірка наявності ключа
    age, exists := ages["Петро"]
    if exists {
        fmt.Printf("Вік Петра: %d\n", age)
    } else {
        fmt.Println("Петро не знайдено")
    }
    
    // Ітерація
    for name, age := range ages {
        fmt.Printf("%s: %d років\n", name, age)
    }
    
    // Видалення
    delete(ages, "Іван")
    fmt.Printf("Після видалення: %v\n", ages)
}
```

---

## Порівняння типів

| Тип | Розмір | Zero Value | Змінюваний |
|-----|--------|------------|-----------|
| `int` | 32/64 біт | `0` | ✅ |
| `float64` | 64 біт | `0.0` | ✅ |
| `string` | змінний | `""` | ❌ |
| `bool` | 1 біт | `false` | ✅ |
| `struct` | сума полів | zero для кожного поля | ✅ |
| `array` | фіксований | zero для кожного елемента | ✅ |
| `slice` | 24 байти | `nil` | ✅ |
| `map` | 8 байт | `nil` | ✅ |

---

## Конверсія типів

```go
package main

import (
    "fmt"
    "strconv"
)

func main() {
    // Числові конверсії
    var i int = 42
    var f float64 = float64(i)
    var u uint = uint(i)
    
    fmt.Printf("int: %d, float64: %.2f, uint: %d\n", i, f, u)
    
    // string <-> int
    str := "123"
    num, err := strconv.Atoi(str)
    if err == nil {
        fmt.Printf("Рядок '%s' → число %d\n", str, num)
    }
    
    numStr := strconv.Itoa(456)
    fmt.Printf("Число 456 → рядок '%s'\n", numStr)
    
    // string <-> float
    floatStr := "3.14"
    floatNum, _ := strconv.ParseFloat(floatStr, 64)
    fmt.Printf("Рядок '%s' → float %.2f\n", floatStr, floatNum)
}
```

---

## Завдання для практики

1. Створити struct `Book` з полями: назва, автор, рік, ціна
2. Створити slice з 3 книг
3. Вивести всі книги за допомогою range
4. Створити map з назвами книг як ключами
5. Знайти найдорожчу книгу

---

## Корисні посилання

- [Go Tour - Basic Types](https://go.dev/tour/basics/11)
- [Go Spec - Types](https://go.dev/ref/spec#Types)
- [Effective Go - Data](https://go.dev/doc/effective_go#data)

