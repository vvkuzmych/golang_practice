# Пакети, main та init в Go

## Що таке пакет?

**Пакет (package)** - це спосіб організації та повторного використання коду в Go.

- Кожен `.go` файл належить до одного пакету
- Пакет = колекція `.go` файлів в одній директорії
- Пакет = модуль, бібліотека, компонент

---

## Оголошення пакету

### Синтаксис

```go
package назва_пакету
```

### Перший рядок кожного файлу

```go
package main  // Виконуваний пакет

package utils  // Бібліотечний пакет

package models  // Пакет моделей
```

---

## Package main

**`package main`** - спеціальний пакет для виконуваних програм.

### Характеристики:
- Точка входу в програму
- Містить функцію `main()`
- Компілюється в executable файл
- Може бути тільки один `main` пакет в програмі

### Приклад

```go
package main

import "fmt"

func main() {
    fmt.Println("Привіт, світ!")
}
```

### Компіляція та запуск

```bash
# Запустити без компіляції
go run main.go

# Скомпілювати в executable
go build main.go
./main

# Скомпілювати з назвою
go build -o myapp main.go
./myapp
```

---

## Функція main()

**`main()`** - точка входу в програму.

### Правила:
- Має бути в пакеті `main`
- Немає параметрів
- Нічого не повертає
- Викликається автоматично при запуску

### Приклад

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    fmt.Println("Програма запущена!")
    
    // Аргументи командного рядка
    fmt.Println("Аргументи:", os.Args)
    
    // Вихід з програми
    if len(os.Args) < 2 {
        fmt.Println("Недостатньо аргументів")
        os.Exit(1)
    }
    
    fmt.Println("Перший аргумент:", os.Args[1])
}
```

---

## Функція init()

**`init()`** - спеціальна функція для ініціалізації пакету.

### Характеристики:
- Викликається автоматично перед `main()`
- Немає параметрів та return значення
- Може бути кілька `init()` в одному пакеті
- Виконується один раз на пакет

### Порядок виконання

```go
package main

import "fmt"

// 1. Ініціалізація змінних пакету
var message = initMessage()

// 2. init() функція
func init() {
    fmt.Println("init() викликана")
}

// 3. main() функція
func main() {
    fmt.Println("main() викликана")
    fmt.Println(message)
}

func initMessage() string {
    fmt.Println("initMessage() викликана")
    return "Привіт з init"
}

// Вивід:
// initMessage() викликана
// init() викликана
// main() викликана
// Привіт з init
```

### Множинні init()

```go
package main

import "fmt"

func init() {
    fmt.Println("init 1")
}

func init() {
    fmt.Println("init 2")
}

func init() {
    fmt.Println("init 3")
}

func main() {
    fmt.Println("main")
}

// Вивід:
// init 1
// init 2
// init 3
// main
```

---

## Структура проєкту

### Прост проєкт

```
myproject/
├── go.mod
├── main.go
└── README.md
```

```go
// main.go
package main

import "fmt"

func main() {
    fmt.Println("Hello!")
}
```

### Проєкт з пакетами

```
myproject/
├── go.mod
├── main.go
├── utils/
│   ├── math.go
│   └── string.go
└── models/
    └── user.go
```

```go
// main.go
package main

import (
    "fmt"
    "myproject/utils"
    "myproject/models"
)

func main() {
    result := utils.Add(5, 3)
    fmt.Println("Result:", result)
    
    user := models.User{Name: "Іван"}
    fmt.Println("User:", user)
}

// utils/math.go
package utils

func Add(a, b int) int {
    return a + b
}

// models/user.go
package models

type User struct {
    Name string
    Age  int
}
```

---

## Імпорти

### Стандартні імпорти

```go
package main

import "fmt"          // Один пакет
import "os"
import "strings"

// Або блок імпортів (рекомендовано)
import (
    "fmt"
    "os"
    "strings"
)
```

### Алізи (Aliases)

```go
package main

import (
    "fmt"
    f "fmt"           // Аліас f
    str "strings"     // Аліас str
    . "math"          // Імпорт в поточний namespace
    _ "database/sql"  // Імпорт тільки для init()
)

func main() {
    fmt.Println("fmt")
    f.Println("f")
    str.ToUpper("test")
    Sqrt(4)           // Без math. завдяки .
}
```

### Порядок імпортів

```go
package main

import (
    // 1. Стандартна бібліотека
    "fmt"
    "os"
    "strings"
    
    // 2. Зовнішні пакети
    "github.com/user/package"
    
    // 3. Внутрішні пакети
    "myproject/utils"
    "myproject/models"
)
```

---

## Експорт та імпорт

### Правило капіталізації

- **Велика літера** = Експортовано (Public)
- **Маленька літера** = Не експортовано (Private)

```go
// utils/math.go
package utils

// Експортовано - доступно з інших пакетів
func Add(a, b int) int {
    return a + b
}

// Не експортовано - доступно тільки в пакеті utils
func subtract(a, b int) int {
    return a - b
}

// Експортована константа
const MaxValue = 100

// Не експортована константа
const minValue = 0

// Експортована структура
type User struct {
    Name string   // Експортоване поле
    age  int      // Не експортоване поле
}
```

```go
// main.go
package main

import (
    "fmt"
    "myproject/utils"
)

func main() {
    // ✅ Працює
    result := utils.Add(5, 3)
    fmt.Println(result)
    
    user := utils.User{Name: "Іван"}
    fmt.Println(user.Name)
    
    // ❌ Помилка - не експортовано
    // result := utils.subtract(5, 3)
    // user.age = 25
}
```

---

## Приклад повного проєкту

### Структура

```
calculator/
├── go.mod
├── main.go
├── operations/
│   ├── math.go
│   └── string.go
└── config/
    └── config.go
```

### go.mod

```go
module calculator

go 1.21
```

### main.go

```go
package main

import (
    "calculator/config"
    "calculator/operations"
    "fmt"
    "os"
)

func init() {
    fmt.Println("Ініціалізація програми...")
    config.Load()
}

func main() {
    if len(os.Args) < 4 {
        fmt.Println("Використання: calculator <число1> <операція> <число2>")
        os.Exit(1)
    }
    
    num1 := parseNumber(os.Args[1])
    op := os.Args[2]
    num2 := parseNumber(os.Args[3])
    
    var result float64
    switch op {
    case "+":
        result = operations.Add(num1, num2)
    case "-":
        result = operations.Subtract(num1, num2)
    case "*":
        result = operations.Multiply(num1, num2)
    case "/":
        result = operations.Divide(num1, num2)
    default:
        fmt.Println("Невідома операція:", op)
        os.Exit(1)
    }
    
    fmt.Printf("Результат: %.2f\n", result)
}

func parseNumber(s string) float64 {
    // Реалізація парсингу
    return 0
}
```

### operations/math.go

```go
package operations

import "fmt"

func init() {
    fmt.Println("Ініціалізація пакету operations...")
}

// Експортовані функції
func Add(a, b float64) float64 {
    return a + b
}

func Subtract(a, b float64) float64 {
    return a - b
}

func Multiply(a, b float64) float64 {
    return a * b
}

func Divide(a, b float64) float64 {
    if b == 0 {
        panic("ділення на нуль")
    }
    return a / b
}
```

### config/config.go

```go
package config

import "fmt"

var (
    AppName    = "Calculator"
    AppVersion = "1.0.0"
    Debug      = false
)

func init() {
    fmt.Println("Ініціалізація пакету config...")
}

func Load() {
    fmt.Printf("Завантажено конфігурацію: %s v%s\n", AppName, AppVersion)
}
```

---

## Порядок ініціалізації

### 1. Імпорт пакетів (dependency order)

```
main імпортує → operations та config
                   ↓
            init() operations
                   ↓
            init() config
                   ↓
            init() main
                   ↓
            main() main
```

### 2. В межах пакету

```
1. Константи
2. Змінні (в порядку залежності)
3. init() функції (в порядку появи)
4. main() (якщо є)
```

### Приклад

```go
package main

import "fmt"

// 1. Константи
const version = "1.0.0"

// 2. Змінні
var name = getName()

// 3. init()
func init() {
    fmt.Println("Запуск init 1")
}

func init() {
    fmt.Println("Запуск init 2")
}

// 4. main()
func main() {
    fmt.Println("Запуск main")
    fmt.Println("Version:", version)
    fmt.Println("Name:", name)
}

func getName() string {
    fmt.Println("Виклик getName()")
    return "MyApp"
}

// Вивід:
// Виклик getName()
// Запуск init 1
// Запуск init 2
// Запуск main
// Version: 1.0.0
// Name: MyApp
```

---

## Blank Import

**Blank import (`_`)** - імпорт пакету тільки для виклику `init()`.

```go
package main

import (
    "fmt"
    _ "myproject/drivers"  // Викликає тільки init()
)

func main() {
    fmt.Println("main")
}
```

```go
// drivers/postgres.go
package drivers

import "fmt"

func init() {
    fmt.Println("Реєстрація PostgreSQL драйвера")
    // RegisterDriver("postgres", &PostgresDriver{})
}
```

### Використання

```go
import (
    "database/sql"
    _ "github.com/lib/pq"  // PostgreSQL драйвер
)
```

---

## Практичні поради

### ✅ Гарна практика

```go
// 1. Одне package оголошення на файл
package main

// 2. Групуйте імпорти
import (
    // стандартна бібліотека
    "fmt"
    "os"
    
    // зовнішні пакети
    "github.com/pkg/errors"
    
    // внутрішні пакети
    "myapp/models"
)

// 3. Використовуйте значущі назви пакетів
package userservice  // ✅
package utils        // ⚠️ надто загально

// 4. init() тільки для ініціалізації
func init() {
    // Налаштування логера
    // Реєстрація драйверів
    // Валідація конфігурації
}
```

### ❌ Погана практика

```go
// 1. Використання . import (плутанина)
import . "fmt"

// 2. Занадто багато в init()
func init() {
    // Багато складної логіки
    // Мережеві запити
    // Тривала ініціалізація
}

// 3. Циклічні імпорти
// package A імпортує B
// package B імпортує A
// ❌ COMPILATION ERROR
```

---

## Резюме

### Package main
- ✅ Виконувана програма
- ✅ Містить `main()`
- ✅ Компілюється в executable

### Функція main()
- ✅ Точка входу
- ✅ Викликається автоматично
- ✅ Немає параметрів/return

### Функція init()
- ✅ Автоматична ініціалізація
- ✅ Викликається перед main()
- ✅ Може бути кілька в пакеті

### Експорт
- ✅ Велика літера = Public
- ✅ Маленька літера = Private

---

## Завдання

1. Створити проєкт з 3 пакетами (main, utils, models)
2. Використати init() для ініціалізації
3. Експортувати та імпортувати функції
4. Дослідити порядок виклику init()

