# Zero Values в Go

## Що таке Zero Value?

**Zero Value** - це значення за замовчуванням, яке змінна отримує при оголошенні без явної ініціалізації.

В Go **немає** `undefined` або `null` для примітивних типів. Кожна змінна завжди має якесь значення.

---

## Zero Values для різних типів

```go
package main

import "fmt"

func main() {
    // Числові типи → 0
    var i int
    var f float64
    var b byte
    
    fmt.Printf("int: %d\n", i)         // 0
    fmt.Printf("float64: %.1f\n", f)   // 0.0
    fmt.Printf("byte: %d\n", b)        // 0
    
    // Рядок → ""
    var s string
    fmt.Printf("string: '%s' (len=%d)\n", s, len(s))  // '' (len=0)
    
    // Логічний → false
    var bo bool
    fmt.Printf("bool: %t\n", bo)       // false
    
    // Вказівник → nil
    var p *int
    fmt.Printf("pointer: %v\n", p)     // <nil>
    
    // Slice → nil
    var slice []int
    fmt.Printf("slice: %v (nil=%t)\n", slice, slice == nil)  // [] (nil=true)
    
    // Map → nil
    var m map[string]int
    fmt.Printf("map: %v (nil=%t)\n", m, m == nil)  // map[] (nil=true)
    
    // Функція → nil
    var fn func()
    fmt.Printf("func: %v\n", fn)       // <nil>
    
    // Інтерфейс → nil
    var i interface{}
    fmt.Printf("interface: %v (nil=%t)\n", i, i == nil)  // <nil> (nil=true)
}
```

---

## Таблиця Zero Values

| Тип | Zero Value | Можна використовувати? |
|-----|------------|------------------------|
| `int`, `int8`, `int16`, `int32`, `int64` | `0` | ✅ |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `0` | ✅ |
| `float32`, `float64` | `0.0` | ✅ |
| `string` | `""` (порожній рядок) | ✅ |
| `bool` | `false` | ✅ |
| `pointer` | `nil` | ⚠️ потрібна перевірка |
| `slice` | `nil` | ⚠️ можна append, НЕ можна індексувати |
| `map` | `nil` | ❌ потрібен make |
| `chan` | `nil` | ❌ потрібен make |
| `func` | `nil` | ⚠️ потрібна перевірка |
| `interface` | `nil` | ⚠️ потрібна перевірка |

---

## Структури та Zero Values

```go
package main

import "fmt"

type Person struct {
    Name    string
    Age     int
    IsAdmin bool
    Address *Address
}

type Address struct {
    City   string
    Street string
}

func main() {
    // Структура без ініціалізації
    var p Person
    
    fmt.Printf("Person: %+v\n", p)
    // Person: {Name: Age:0 IsAdmin:false Address:<nil>}
    
    // Кожне поле має zero value свого типу:
    fmt.Printf("Name: '%s'\n", p.Name)       // ''
    fmt.Printf("Age: %d\n", p.Age)           // 0
    fmt.Printf("IsAdmin: %t\n", p.IsAdmin)   // false
    fmt.Printf("Address: %v\n", p.Address)   // <nil>
}
```

---

## Практичні приклади

### Приклад 1: Безпечна робота з Zero Values

```go
package main

import "fmt"

func main() {
    // ✅ Безпечно: int zero value = 0
    var counter int
    counter++
    fmt.Println(counter)  // 1
    
    // ✅ Безпечно: string zero value = ""
    var message string
    message += "Привіт"
    fmt.Println(message)  // Привіт
    
    // ✅ Безпечно: slice zero value = nil (можна append)
    var numbers []int
    numbers = append(numbers, 1, 2, 3)
    fmt.Println(numbers)  // [1 2 3]
    
    // ❌ НЕБЕЗПЕЧНО: map zero value = nil (не можна записувати)
    var ages map[string]int
    // ages["Іван"] = 25  // panic: assignment to entry in nil map
    
    // ✅ Правильно: ініціалізувати map
    ages = make(map[string]int)
    ages["Іван"] = 25
    fmt.Println(ages)  // map[Іван:25]
}
```

### Приклад 2: Перевірка на nil

```go
package main

import "fmt"

type User struct {
    Name  string
    Email *string
}

func main() {
    user := User{Name: "Іван"}
    
    // Перевірка pointer на nil
    if user.Email != nil {
        fmt.Printf("Email: %s\n", *user.Email)
    } else {
        fmt.Println("Email не вказано")
    }
    
    // Встановлення значення
    email := "ivan@example.com"
    user.Email = &email
    
    if user.Email != nil {
        fmt.Printf("Email: %s\n", *user.Email)
    }
}
```

### Приклад 3: Функція з Zero Value

```go
package main

import "fmt"

// Zero value для повернення
func divide(a, b int) (int, error) {
    if b == 0 {
        // Повертаємо zero value для int (0) та error
        return 0, fmt.Errorf("ділення на нуль")
    }
    return a / b, nil
}

func main() {
    result, err := divide(10, 2)
    if err != nil {
        fmt.Println("Помилка:", err)
    } else {
        fmt.Println("Результат:", result)
    }
    
    result2, err2 := divide(10, 0)
    if err2 != nil {
        fmt.Println("Помилка:", err2)
        fmt.Println("Результат (zero value):", result2)  // 0
    }
}
```

---

## Коли важливі Zero Values?

### 1. Ініціалізація в циклах

```go
package main

import "fmt"

func main() {
    var sum int  // zero value = 0
    numbers := []int{1, 2, 3, 4, 5}
    
    for _, num := range numbers {
        sum += num  // працює бо sum = 0 на початку
    }
    
    fmt.Println("Сума:", sum)  // 15
}
```

### 2. Опціональні поля

```go
package main

import "fmt"

type Config struct {
    Host string  // обов'язкове
    Port int     // опціональне, zero value = 0
}

func NewConfig(host string, port int) Config {
    cfg := Config{Host: host}
    
    if port == 0 {
        cfg.Port = 8080  // default
    } else {
        cfg.Port = port
    }
    
    return cfg
}

func main() {
    cfg1 := NewConfig("localhost", 0)
    fmt.Printf("Config: %+v\n", cfg1)  // Port: 8080
    
    cfg2 := NewConfig("localhost", 3000)
    fmt.Printf("Config: %+v\n", cfg2)  // Port: 3000
}
```

### 3. Flags та опції

```go
package main

import "fmt"

type Options struct {
    Verbose bool    // zero value = false
    Timeout int     // zero value = 0
    Output  string  // zero value = ""
}

func Run(opts Options) {
    if opts.Verbose {
        fmt.Println("Verbose mode ON")
    }
    
    timeout := opts.Timeout
    if timeout == 0 {
        timeout = 30  // default
    }
    fmt.Printf("Timeout: %d seconds\n", timeout)
    
    output := opts.Output
    if output == "" {
        output = "output.txt"  // default
    }
    fmt.Printf("Output: %s\n", output)
}

func main() {
    // Всі поля мають zero values
    Run(Options{})
    
    // Деякі поля встановлені
    Run(Options{Verbose: true, Timeout: 60})
}
```

---

## Поширені помилки

### ❌ Помилка 1: Використання nil map

```go
var m map[string]int
// m["key"] = 1  // PANIC!

// Правильно:
m = make(map[string]int)
m["key"] = 1
```

### ❌ Помилка 2: Індексація nil slice

```go
var s []int
// fmt.Println(s[0])  // PANIC!

// Правильно:
s = append(s, 1)
fmt.Println(s[0])
```

### ❌ Помилка 3: Дереференс nil pointer

```go
var p *int
// fmt.Println(*p)  // PANIC!

// Правильно:
if p != nil {
    fmt.Println(*p)
}
```

---

## Корисні паттерни

### Паттерн 1: Lazy Initialization

```go
type Cache struct {
    data map[string]string
}

func (c *Cache) Get(key string) string {
    // Ініціалізація при першому використанні
    if c.data == nil {
        c.data = make(map[string]string)
    }
    return c.data[key]
}

func (c *Cache) Set(key, value string) {
    if c.data == nil {
        c.data = make(map[string]string)
    }
    c.data[key] = value
}
```

### Паттерн 2: Optional Fields з Pointers

```go
type User struct {
    Name     string   // обов'язкове
    Age      *int     // опціональне
    Email    *string  // опціональне
}

// nil означає "не вказано"
// pointer означає "вказано, значення X"
```

---

## Резюме

✅ **Zero Values - це безпечно**
- Кожна змінна завжди ініціалізована
- Немає "undefined" поведінки

⚠️ **Винятки потребують уваги**
- `map`: потрібен `make()`
- `pointer`: потрібна перевірка на `nil`
- `slice`: можна `append()`, але не індексувати пустий

🎯 **Використовуйте zero values**
- Для дефолтних значень
- Для опціональних параметрів
- Для прапорців (flags)

---

## Завдання

1. Створити функцію, яка приймає pointer і повертає zero value, якщо pointer = nil
2. Написати функцію з опціональними параметрами через struct
3. Реалізувати lazy initialization для map

