# Назва програми: os.Args[0] - Кращі практики

## Проблема з os.Args[0]

`os.Args[0]` містить шлях до виконуваного файлу, який може виглядати по-різному:

```bash
# При компіляції
$ go build -o greet main.go
$ ./greet
os.Args[0] = "./greet"  ✅ Нормально

# При go run
$ go run main.go
os.Args[0] = "/tmp/go-build123456/exe/main"  ❌ Незручно!

# При запуску з іншої директорії
$ /usr/local/bin/greet
os.Args[0] = "/usr/local/bin/greet"  ⚠️ Довгий шлях
```

---

## Рішення

### ✅ Рішення 1: Константа (Найпростіше)

```go
package main

import "fmt"

// Константа з назвою програми
const programName = "greet"

func printUsage() {
    fmt.Printf("Використання: %s <ім'я> [вік]\n", programName)
    fmt.Printf("Приклади:\n")
    fmt.Printf("  %s Іван\n", programName)
    fmt.Printf("  %s Марія 25\n", programName)
}

func main() {
    printUsage()
}
```

**Переваги:**
- ✅ Просто
- ✅ Зрозуміло
- ✅ Завжди однакове
- ✅ Легко змінити

**Недоліки:**
- ❌ Потрібно вручну оновлювати при зміні назви

---

### ✅ Рішення 2: filepath.Base() (Автоматично)

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
)

func getProgramName() string {
    // Отримати тільки назву файлу без шляху
    return filepath.Base(os.Args[0])
}

func printUsage() {
    progName := getProgramName()
    fmt.Printf("Використання: %s <ім'я> [вік]\n", progName)
    fmt.Printf("Приклади:\n")
    fmt.Printf("  %s Іван\n", progName)
    fmt.Printf("  %s Марія 25\n", progName)
}

func main() {
    printUsage()
}
```

**Результат:**
```bash
$ go run main.go
# progName = "main" (з go run отримаємо назву файлу)

$ go build -o greet main.go
$ ./greet
# progName = "greet" ✅
```

**Переваги:**
- ✅ Автоматично отримує назву
- ✅ Працює з різними шляхами

**Недоліки:**
- ⚠️ При `go run` показує "main" замість реальної назви

---

### ✅ Рішення 3: Змінна з можливістю перевизначення

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
)

var (
    // Назва програми за замовчуванням
    defaultName = "greet"
    
    // Фактична назва (можна змінити під час компіляції)
    programName string
)

func init() {
    if programName == "" {
        // Спробувати отримати з os.Args[0]
        name := filepath.Base(os.Args[0])
        
        // Якщо це go run, використати defaultName
        if name == "main" || filepath.Ext(name) != "" {
            programName = defaultName
        } else {
            programName = name
        }
    }
}

func printUsage() {
    fmt.Printf("Використання: %s <ім'я> [вік]\n", programName)
    fmt.Printf("Приклади:\n")
    fmt.Printf("  %s Іван\n", programName)
    fmt.Printf("  %s Марія 25\n", programName)
}

func main() {
    printUsage()
}
```

**Переваги:**
- ✅ Працює з `go run` і скомпільованим бінарником
- ✅ Можна перевизначити при компіляції
- ✅ Автоматично визначає правильну назву

---

### ✅ Рішення 4: Встановлення під час компіляції (Advanced)

```go
package main

import "fmt"

var (
    // Буде встановлено під час компіляції
    version   = "dev"
    buildTime = "unknown"
    appName   = "greet"
)

func printUsage() {
    fmt.Printf("%s v%s (build: %s)\n\n", appName, version, buildTime)
    fmt.Printf("Використання: %s <ім'я> [вік]\n", appName)
}

func main() {
    printUsage()
}
```

**Компіляція з параметрами:**
```bash
go build -ldflags "\
  -X main.version=1.0.0 \
  -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%S) \
  -X main.appName=greet" \
  -o greet main.go

./greet
# Вивід: greet v1.0.0 (build: 2024-01-15T10:30:00)
```

**Переваги:**
- ✅ Професійно
- ✅ Повний контроль
- ✅ Можна додати версію та інфо про білд

---

## Повний приклад з кращими практиками

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strconv"
)

// Конфігурація програми
const (
    appName        = "greet"
    appVersion     = "1.0.0"
    appDescription = "Привітання користувача"
)

func getProgramName() string {
    // Спробувати отримати з os.Args[0]
    if len(os.Args) > 0 {
        name := filepath.Base(os.Args[0])
        // Якщо це не "main" з go run, використати його
        if name != "main" && filepath.Ext(name) != ".go" {
            return name
        }
    }
    // Інакше використати константу
    return appName
}

func printVersion() {
    fmt.Printf("%s v%s\n", appName, appVersion)
    fmt.Printf("%s\n", appDescription)
}

func printUsage() {
    progName := getProgramName()
    
    fmt.Println("❌ Помилка: не вказано ім'я\n")
    printVersion()
    fmt.Println("\nВикористання:")
    fmt.Printf("  %s <ім'я> [вік]\n", progName)
    fmt.Printf("  %s --version\n", progName)
    fmt.Printf("  %s --help\n\n", progName)
    
    fmt.Println("Аргументи:")
    fmt.Println("  <ім'я>     Ваше ім'я (обов'язково)")
    fmt.Println("  [вік]      Ваш вік, число 0-120 (опційно)")
    
    fmt.Println("\nОпції:")
    fmt.Println("  --version  Показати версію програми")
    fmt.Println("  --help     Показати цю довідку")
    
    fmt.Println("\nПриклади:")
    fmt.Printf("  %s Іван\n", progName)
    fmt.Printf("  %s Марія 25\n", progName)
}

func main() {
    // Обробка спеціальних аргументів
    if len(os.Args) > 1 {
        switch os.Args[1] {
        case "--version", "-v":
            printVersion()
            return
        case "--help", "-h":
            printUsage()
            return
        }
    }
    
    // Перевірка аргументів
    if len(os.Args) < 2 {
        printUsage()
        os.Exit(1)
    }
    
    name := os.Args[1]
    
    // Обробка віку
    var age int
    var hasAge bool
    
    if len(os.Args) >= 3 {
        parsedAge, err := strconv.Atoi(os.Args[2])
        if err != nil {
            fmt.Printf("❌ Помилка: '%s' не є числом\n", os.Args[2])
            os.Exit(1)
        }
        age = parsedAge
        hasAge = true
    }
    
    // Вивід
    fmt.Printf("Привіт, %s! 👋\n", name)
    if hasAge {
        fmt.Printf("Тобі %d років.\n", age)
    }
}
```

---

## Порівняння підходів

| Підхід | Простота | Гнучкість | Для go run | Для build |
|--------|----------|-----------|------------|-----------|
| **Константа** | ⭐⭐⭐⭐⭐ | ⭐⭐ | ✅ | ✅ |
| **filepath.Base()** | ⭐⭐⭐⭐ | ⭐⭐⭐ | ⚠️ | ✅ |
| **Змінна + init()** | ⭐⭐⭐ | ⭐⭐⭐⭐ | ✅ | ✅ |
| **Build-time flags** | ⭐⭐ | ⭐⭐⭐⭐⭐ | ⚠️ | ✅ |

---

## Рекомендації

### Для навчання (Week 1):
```go
const programName = "greet"
fmt.Printf("Використання: %s <args>\n", programName)
```
✅ **Просто і зрозуміло**

### Для невеликих проектів:
```go
func getProgramName() string {
    name := filepath.Base(os.Args[0])
    if name == "main" {
        return "myapp"
    }
    return name
}
```
✅ **Баланс між простотою і функціональністю**

### Для великих проектів:
```go
var (
    appName    = "myapp"
    appVersion = "1.0.0"
)
// + build-time flags
```
✅ **Професійно з версіонуванням**

---

## Практичні приклади

### Приклад 1: CLI калькулятор

```go
package main

import "fmt"

const programName = "calc"

func printUsage() {
    fmt.Printf("Використання: %s <число1> <операція> <число2>\n", programName)
    fmt.Println("\nОперації: +, -, *, /")
    fmt.Println("\nПриклади:")
    fmt.Printf("  %s 10 + 5\n", programName)
    fmt.Printf("  %s 20 - 7\n", programName)
    fmt.Printf("  %s 6 mul 3\n", programName)
}

func main() {
    printUsage()
}
```

### Приклад 2: TODO Manager

```go
package main

import "fmt"

const (
    appName = "todo"
    appVersion = "1.0.0"
)

func printHelp() {
    fmt.Printf("\n%s v%s - Менеджер завдань\n\n", appName, appVersion)
    fmt.Printf("Використання: %s <команда> [аргументи]\n\n", appName)
    
    fmt.Println("Команди:")
    fmt.Printf("  %s add <текст>      Додати завдання\n", appName)
    fmt.Printf("  %s list             Показати всі\n", appName)
    fmt.Printf("  %s done <id>        Позначити виконаним\n", appName)
    fmt.Printf("  %s delete <id>      Видалити\n", appName)
    fmt.Printf("  %s help             Ця довідка\n", appName)
}

func main() {
    printHelp()
}
```

---

## Резюме

✅ **Використовуйте константу** для простих програм
✅ **Додайте версію** для більш серйозних проектів
✅ **Уникайте прямого використання** `os.Args[0]` у виводі
✅ **Зробіть код читабельним** і зручним для користувачів

---

## Завдання

1. Оновити solution_1.go з константою programName
2. Додати функцію getProgramName() у solution_3.go
3. Додати --version та --help у ваші програми

