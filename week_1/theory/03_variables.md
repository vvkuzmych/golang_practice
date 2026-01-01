# Змінні в Go: var vs :=

## Два способи оголошення змінних

В Go є два основні способи оголосити змінну:
1. **`var`** - повне оголошення
2. **`:=`** - коротке оголошення (short declaration)

---

## 1. var - Повне оголошення

### Синтаксис

```go
var ім'я_змінної тип = значення
```

### Приклади

```go
package main

import "fmt"

func main() {
    // Повне оголошення з ініціалізацією
    var name string = "Іван"
    var age int = 25
    var isStudent bool = true
    
    fmt.Println(name, age, isStudent)
    
    // Оголошення без ініціалізації (zero value)
    var counter int
    var message string
    var active bool
    
    fmt.Printf("counter: %d, message: '%s', active: %t\n", 
        counter, message, active)
    // counter: 0, message: '', active: false
    
    // Автоматичне визначення типу
    var price = 19.99      // float64
    var product = "хліб"   // string
    
    fmt.Printf("price: %.2f (type: %T)\n", price, price)
    fmt.Printf("product: %s (type: %T)\n", product, product)
}
```

### Множинне оголошення

```go
package main

import "fmt"

func main() {
    // Кілька змінних одного типу
    var x, y, z int
    x = 1
    y = 2
    z = 3
    
    // Кілька змінних з ініціалізацією
    var a, b, c = 1, "два", true
    
    // Блок оголошень
    var (
        firstName string = "Іван"
        lastName  string = "Петренко"
        age       int    = 25
        height    float64 = 180.5
    )
    
    fmt.Println(x, y, z)
    fmt.Println(a, b, c)
    fmt.Println(firstName, lastName, age, height)
}
```

---

## 2. := - Коротке оголошення

### Синтаксис

```go
ім'я_змінної := значення
```

### Приклади

```go
package main

import "fmt"

func main() {
    // Коротке оголошення
    name := "Марія"
    age := 22
    isStudent := true
    
    fmt.Println(name, age, isStudent)
    
    // Тип визначається автоматично
    price := 29.99    // float64
    count := 10       // int
    
    fmt.Printf("price: %.2f (type: %T)\n", price, price)
    fmt.Printf("count: %d (type: %T)\n", count, count)
    
    // Множинне присвоєння
    x, y := 10, 20
    firstName, lastName := "Петро", "Сидоренко"
    
    fmt.Println(x, y)
    fmt.Println(firstName, lastName)
}
```

---

## Порівняння var vs :=

| Характеристика | `var` | `:=` |
|----------------|-------|------|
| **Де можна використовувати** | Всюди | Тільки в функції |
| **Вказівка типу** | Можна | Неможливо |
| **Zero value** | Можливо | Неможливо |
| **Пакетний рівень** | ✅ Так | ❌ Ні |
| **В функції** | ✅ Так | ✅ Так |
| **Блок оголошень** | ✅ Так | ❌ Ні |
| **Читабельність** | Більше коду | Менше коду |

---

## Коли використовувати var

### 1. Пакетний рівень (Package Level)

```go
package main

import "fmt"

// Можна тільки var
var globalCounter int = 0
var appName string = "MyApp"
var debug bool = false

// НЕ можна :=
// globalCounter := 0  // ❌ ПОМИЛКА!

func main() {
    fmt.Println(globalCounter, appName, debug)
}
```

### 2. Zero Values

```go
package main

import "fmt"

func main() {
    // Потрібен zero value
    var counter int      // 0
    var message string   // ""
    var active bool      // false
    
    // Не можна з :=
    // counter := ???  // якщо не знаємо значення
    
    for i := 0; i < 5; i++ {
        counter++
    }
    
    fmt.Println(counter)  // 5
}
```

### 3. Явна вказівка типу

```go
package main

import "fmt"

func main() {
    // Хочемо конкретний тип
    var age int32 = 25          // int32
    var price float32 = 19.99   // float32
    
    // З := було б
    // age := 25        // int (не int32!)
    // price := 19.99   // float64 (не float32!)
    
    fmt.Printf("age: %d (type: %T)\n", age, age)
    fmt.Printf("price: %.2f (type: %T)\n", price, price)
}
```

### 4. Блок оголошень

```go
package main

import "fmt"

func main() {
    var (
        name      string  = "Іван"
        age       int     = 25
        isStudent bool    = true
        gpa       float64 = 3.8
    )
    
    fmt.Println(name, age, isStudent, gpa)
}
```

---

## Коли використовувати :=

### 1. Локальні змінні в функціях

```go
package main

import "fmt"

func main() {
    // Швидке оголошення локальних змінних
    name := "Олена"
    age := 28
    city := "Київ"
    
    fmt.Println(name, age, city)
}
```

### 2. В циклах

```go
package main

import "fmt"

func main() {
    // for loop
    for i := 0; i < 5; i++ {
        fmt.Println(i)
    }
    
    // range
    numbers := []int{10, 20, 30, 40, 50}
    for index, value := range numbers {
        fmt.Printf("numbers[%d] = %d\n", index, value)
    }
}
```

### 3. Повернення з функцій

```go
package main

import (
    "fmt"
    "errors"
)

func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("ділення на нуль")
    }
    return a / b, nil
}

func main() {
    // Швидко отримати результат
    result, err := divide(10, 2)
    if err != nil {
        fmt.Println("Помилка:", err)
        return
    }
    fmt.Println("Результат:", result)
}
```

### 4. Короткий і читабельний код

```go
package main

import "fmt"

func main() {
    // Замість:
    var firstName string = "Андрій"
    var lastName string = "Коваленко"
    var fullName string = firstName + " " + lastName
    
    // Краще:
    first := "Андрій"
    last := "Коваленко"
    full := first + " " + last
    
    fmt.Println(full)
}
```

---

## Важливі деталі :=

### 1. Переприсвоєння (Re-declaration)

```go
package main

import "fmt"

func main() {
    // Оголошення
    x := 10
    fmt.Println("x =", x)
    
    // ❌ НЕ можна повторно оголосити
    // x := 20  // ПОМИЛКА: no new variables on left side
    
    // ✅ Але можна присвоїти нове значення
    x = 20
    fmt.Println("x =", x)
    
    // ✅ Можна, якщо хоча б одна змінна нова
    x, y := 30, 40  // y - нова
    fmt.Println("x =", x, ", y =", y)
}
```

### 2. Область видимості (Scope)

```go
package main

import "fmt"

func main() {
    x := 10
    fmt.Println("outer x:", x)
    
    if true {
        x := 20  // Нова змінна в блоці!
        fmt.Println("inner x:", x)  // 20
    }
    
    fmt.Println("outer x:", x)  // 10 (не змінилася!)
}
```

### 3. Не можна поза функцією

```go
package main

// ✅ Правильно
var globalVar = "hello"

// ❌ ПОМИЛКА
// shortVar := "world"  // syntax error: non-declaration statement outside function body

func main() {
    // ✅ Тут можна
    localVar := "привіт"
    println(localVar)
}
```

---

## Практичні приклади

### Приклад 1: Конфігурація

```go
package main

import "fmt"

// Пакетний рівень - var
var (
    appName    = "MyApp"
    appVersion = "1.0.0"
    debug      = false
)

func main() {
    // Локальний рівень - :=
    host := "localhost"
    port := 8080
    timeout := 30
    
    fmt.Printf("%s v%s\n", appName, appVersion)
    fmt.Printf("Server: %s:%d (timeout: %ds)\n", host, port, timeout)
}
```

### Приклад 2: Обробка помилок

```go
package main

import (
    "fmt"
    "strconv"
)

func main() {
    // := для результатів функцій
    num, err := strconv.Atoi("123")
    if err != nil {
        fmt.Println("Помилка:", err)
        return
    }
    
    fmt.Println("Число:", num)
    
    // var коли потрібен zero value
    var result int
    if num > 100 {
        result = num * 2
    } else {
        result = num
    }
    
    fmt.Println("Результат:", result)
}
```

### Приклад 3: Структури

```go
package main

import "fmt"

type Person struct {
    Name string
    Age  int
}

func main() {
    // := для створення екземплярів
    person1 := Person{
        Name: "Іван",
        Age:  25,
    }
    
    person2 := Person{"Марія", 22}
    
    // var коли потрібен zero value
    var person3 Person  // Name: "", Age: 0
    
    fmt.Printf("Person 1: %+v\n", person1)
    fmt.Printf("Person 2: %+v\n", person2)
    fmt.Printf("Person 3: %+v\n", person3)
}
```

---

## Стиль та конвенції

### ✅ Гарний стиль

```go
// 1. Використовуйте := для локальних змінних
func process() {
    result := calculate()
    data := fetch()
}

// 2. Використовуйте var для пакетного рівня
var defaultTimeout = 30

// 3. Використовуйте var для zero values
func reset() {
    var counter int  // явно показує zero value
    // ...
}

// 4. Групуйте пов'язані var
var (
    maxRetries = 3
    timeout    = 30
    debug      = false
)
```

### ❌ Поганий стиль

```go
// 1. Надмірне використання var в функції
func process() {
    var result int = calculate()     // краще result := calculate()
    var data string = fetch()        // краще data := fetch()
}

// 2. Довгі імена з := (важко читати)
func main() {
    thisIsAVeryLongVariableNameThatMakesCodeHardToRead := getValue()
}
```

---

## Резюме

### Використовуйте `var` коли:
- ✅ Пакетний рівень (package level)
- ✅ Потрібен zero value
- ✅ Потрібен конкретний тип (int32, float32)
- ✅ Блок оголошень пов'язаних змінних
- ✅ Експорт змінної (перша літера велика)

### Використовуйте `:=` коли:
- ✅ Локальні змінні в функції
- ✅ Цикли (for, range)
- ✅ Результати функцій
- ✅ Швидке присвоєння
- ✅ Коротший і читабельніший код

---

## Завдання для практики

1. Перепишіть код з `var` на `:=` де можливо
2. Створіть пакетні константи та змінні
3. Напишіть функцію з правильним використанням var vs :=
4. Поясніть, чому певний код використовує var замість :=

