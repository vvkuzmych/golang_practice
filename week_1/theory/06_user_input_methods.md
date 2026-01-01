# Методи вводу даних користувача в Go

## os.Args vs Інтерактивний ввід

### 🔀 Два підходи:

| Підхід | Коли використовувати | Приклад |
|--------|---------------------|---------|
| **os.Args** | CLI утиліти, скрипти, автоматизація | `program arg1 arg2` |
| **Інтерактивний ввід** | Програми з діалогом, меню, форми | Програма запитує → Користувач відповідає |

---

## 1. os.Args - Аргументи командного рядка

### Переваги ✅
- Швидко для досвідчених користувачів
- Легко автоматизувати (скрипти)
- Не потребує взаємодії
- Можна передати багато параметрів одразу

### Недоліки ❌
- Складно для початківців
- Треба пам'ятати синтаксис
- Важко виправити помилки
- Не інтуїтивно

### Приклад
```go
package main

import (
    "fmt"
    "os"
)

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Використання: program <ім'я> <вік>")
        return
    }
    
    name := os.Args[1]
    age := os.Args[2]
    
    fmt.Printf("Привіт, %s! Вік: %s\n", name, age)
}
```

**Використання:**
```bash
$ go run main.go Іван 25
Привіт, Іван! Вік: 25
```

---

## 2. Інтерактивний ввід - Діалог з користувачем

### Переваги ✅
- Інтуїтивно зрозуміло
- Можна виправити помилки
- Підказки для користувача
- Краще для складних даних

### Недоліки ❌
- Повільніше
- Важко автоматизувати
- Потребує взаємодії

---

## Методи інтерактивного вводу

### Метод 1: fmt.Scan() - Найпростіший

```go
package main

import "fmt"

func main() {
    var name string
    var age int
    
    fmt.Print("Введіть ім'я: ")
    fmt.Scan(&name)  // Зчитує до пробілу
    
    fmt.Print("Введіть вік: ")
    fmt.Scan(&age)
    
    fmt.Printf("Привіт, %s! Вік: %d\n", name, age)
}
```

**Особливості:**
- ✅ Просто
- ❌ Зчитує тільки до пробілу/Enter
- ❌ Важко обробляти помилки

**Приклад:**
```
Введіть ім'я: Іван
Введіть вік: 25
Привіт, Іван! Вік: 25
```

---

### Метод 2: fmt.Scanln() - З новим рядком

```go
package main

import "fmt"

func main() {
    var name string
    var age int
    
    fmt.Print("Введіть ім'я: ")
    fmt.Scanln(&name)  // Зчитує до Enter
    
    fmt.Print("Введіть вік: ")
    fmt.Scanln(&age)
    
    fmt.Printf("Привіт, %s! Вік: %d\n", name, age)
}
```

**Особливості:**
- ✅ Простіше за fmt.Scan()
- ❌ Все ще проблеми з пробілами
- ❌ Важко обробляти помилки

---

### Метод 3: bufio.Reader - Найкращий ✅

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
)

func main() {
    // Створюємо reader для читання з консолі
    reader := bufio.NewReader(os.Stdin)
    
    // Читаємо рядок
    fmt.Print("Введіть ваше ім'я: ")
    name, _ := reader.ReadString('\n')
    name = strings.TrimSpace(name)  // Видаляємо \n
    
    // Читаємо число
    fmt.Print("Введіть ваш вік: ")
    ageStr, _ := reader.ReadString('\n')
    ageStr = strings.TrimSpace(ageStr)
    age, err := strconv.Atoi(ageStr)
    
    if err != nil {
        fmt.Println("Помилка: вік має бути числом")
        return
    }
    
    fmt.Printf("Привіт, %s! Вік: %d років\n", name, age)
}
```

**Особливості:**
- ✅ Читає все, включно з пробілами
- ✅ Можна обробляти помилки
- ✅ Гнучко
- ✅ **РЕКОМЕНДОВАНО!**

---

### Метод 4: bufio.Scanner - Альтернатива

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
)

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    
    fmt.Print("Введіть ім'я: ")
    scanner.Scan()
    name := scanner.Text()
    
    fmt.Print("Введіть вік: ")
    scanner.Scan()
    ageStr := scanner.Text()
    age, _ := strconv.Atoi(ageStr)
    
    fmt.Printf("Привіт, %s! Вік: %d\n", name, age)
}
```

**Особливості:**
- ✅ Зручний API
- ✅ Автоматично видаляє \n
- ✅ Хороша обробка помилок

---

## Порівняння методів

| Метод | Складність | Пробіли | Помилки | Рекомендація |
|-------|-----------|---------|---------|--------------|
| `fmt.Scan()` | ⭐ | ❌ | ❌ | Тільки для простих випадків |
| `fmt.Scanln()` | ⭐ | ❌ | ❌ | Тільки для простих випадків |
| `bufio.Reader` | ⭐⭐⭐ | ✅ | ✅ | ✅ **РЕКОМЕНДОВАНО** |
| `bufio.Scanner` | ⭐⭐ | ✅ | ✅ | ✅ Добра альтернатива |

---

## Повні приклади

### Приклад 1: Простий діалог

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    
    fmt.Println("=== Анкета ===\n")
    
    // Ім'я
    fmt.Print("Ваше ім'я: ")
    name, _ := reader.ReadString('\n')
    name = strings.TrimSpace(name)
    
    // Місто
    fmt.Print("Ваше місто: ")
    city, _ := reader.ReadString('\n')
    city = strings.TrimSpace(city)
    
    // Хобі
    fmt.Print("Ваше хобі: ")
    hobby, _ := reader.ReadString('\n')
    hobby = strings.TrimSpace(hobby)
    
    // Вивід
    fmt.Println("\n--- Ваші дані ---")
    fmt.Printf("Ім'я: %s\n", name)
    fmt.Printf("Місто: %s\n", city)
    fmt.Printf("Хобі: %s\n", hobby)
}
```

---

### Приклад 2: Калькулятор з меню

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    
    for {
        fmt.Println("\n=== Калькулятор ===")
        fmt.Println("1 - Додавання")
        fmt.Println("2 - Віднімання")
        fmt.Println("3 - Множення")
        fmt.Println("4 - Ділення")
        fmt.Println("0 - Вихід")
        fmt.Print("\nВибір: ")
        
        choiceStr, _ := reader.ReadString('\n')
        choice, _ := strconv.Atoi(strings.TrimSpace(choiceStr))
        
        if choice == 0 {
            fmt.Println("До побачення!")
            break
        }
        
        if choice < 1 || choice > 4 {
            fmt.Println("Невірний вибір!")
            continue
        }
        
        // Введення чисел
        fmt.Print("Перше число: ")
        num1Str, _ := reader.ReadString('\n')
        num1, _ := strconv.ParseFloat(strings.TrimSpace(num1Str), 64)
        
        fmt.Print("Друге число: ")
        num2Str, _ := reader.ReadString('\n')
        num2, _ := strconv.ParseFloat(strings.TrimSpace(num2Str), 64)
        
        // Обчислення
        var result float64
        switch choice {
        case 1:
            result = num1 + num2
            fmt.Printf("%.2f + %.2f = %.2f\n", num1, num2, result)
        case 2:
            result = num1 - num2
            fmt.Printf("%.2f - %.2f = %.2f\n", num1, num2, result)
        case 3:
            result = num1 * num2
            fmt.Printf("%.2f × %.2f = %.2f\n", num1, num2, result)
        case 4:
            if num2 != 0 {
                result = num1 / num2
                fmt.Printf("%.2f ÷ %.2f = %.2f\n", num1, num2, result)
            } else {
                fmt.Println("Помилка: ділення на нуль!")
            }
        }
    }
}
```

---

### Приклад 3: Валідація вводу

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
)

func readString(reader *bufio.Reader, prompt string) string {
    fmt.Print(prompt)
    input, _ := reader.ReadString('\n')
    return strings.TrimSpace(input)
}

func readInt(reader *bufio.Reader, prompt string) (int, error) {
    input := readString(reader, prompt)
    return strconv.Atoi(input)
}

func readFloat(reader *bufio.Reader, prompt string) (float64, error) {
    input := readString(reader, prompt)
    return strconv.ParseFloat(input, 64)
}

func readBool(reader *bufio.Reader, prompt string) bool {
    input := strings.ToLower(readString(reader, prompt))
    return input == "так" || input == "yes" || input == "y" || input == "true"
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    
    // Читання з валідацією
    name := readString(reader, "Ім'я: ")
    
    var age int
    for {
        var err error
        age, err = readInt(reader, "Вік (0-120): ")
        if err != nil {
            fmt.Println("❌ Помилка: введіть число")
            continue
        }
        if age < 0 || age > 120 {
            fmt.Println("❌ Помилка: вік має бути 0-120")
            continue
        }
        break
    }
    
    isStudent := readBool(reader, "Студент? (так/ні): ")
    
    fmt.Println("\n--- Результат ---")
    fmt.Printf("Ім'я: %s\n", name)
    fmt.Printf("Вік: %d\n", age)
    fmt.Printf("Студент: %t\n", isStudent)
}
```

---

## Коли використовувати що?

### Використовуйте os.Args коли:
- ✅ Пишете CLI утиліту
- ✅ Потрібна автоматизація
- ✅ Користувачі технічно підковані
- ✅ Швидкість важлива

**Приклади:**
- `git commit -m "message"`
- `grep "pattern" file.txt`
- `docker run -p 8080:80 nginx`

---

### Використовуйте інтерактивний ввід коли:
- ✅ Програма для звичайних користувачів
- ✅ Багато параметрів
- ✅ Потрібна валідація
- ✅ Меню або вибір опцій

**Приклади:**
- Анкети
- Калькулятори
- Ігри
- Установщики

---

## Комбінований підхід (Найкраще!)

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    
    var name string
    var age int
    
    // Якщо є аргументи - використовуємо їх
    if len(os.Args) >= 3 {
        name = os.Args[1]
        age, _ = strconv.Atoi(os.Args[2])
        fmt.Println("📝 Використано аргументи командного рядка")
    } else {
        // Інакше - запитуємо інтерактивно
        fmt.Println("📝 Інтерактивний режим")
        fmt.Print("Ім'я: ")
        name, _ = reader.ReadString('\n')
        name = strings.TrimSpace(name)
        
        fmt.Print("Вік: ")
        ageStr, _ := reader.ReadString('\n')
        age, _ = strconv.Atoi(strings.TrimSpace(ageStr))
    }
    
    fmt.Printf("\nПривіт, %s! Вік: %d\n", name, age)
}
```

**Використання:**
```bash
# CLI режим
$ go run main.go Іван 25
Привіт, Іван! Вік: 25

# Інтерактивний режим
$ go run main.go
Ім'я: Марія
Вік: 22
Привіт, Марія! Вік: 22
```

---

## Резюме

| Ситуація | Рекомендація |
|----------|--------------|
| **CLI утиліта** | `os.Args` |
| **Програма з діалогом** | `bufio.Reader` |
| **Просте читання** | `fmt.Scan()` |
| **Складна валідація** | `bufio.Reader` + helper функції |
| **Універсальна програма** | Комбінований підхід |

---

## Корисні функції-помічники

```go
// Читання рядка з валідацією
func readNonEmptyString(reader *bufio.Reader, prompt string) string {
    for {
        fmt.Print(prompt)
        input, _ := reader.ReadString('\n')
        input = strings.TrimSpace(input)
        if input != "" {
            return input
        }
        fmt.Println("❌ Поле не може бути порожнім")
    }
}

// Читання числа в діапазоні
func readIntInRange(reader *bufio.Reader, prompt string, min, max int) int {
    for {
        fmt.Print(prompt)
        input, _ := reader.ReadString('\n')
        num, err := strconv.Atoi(strings.TrimSpace(input))
        if err != nil {
            fmt.Println("❌ Введіть число")
            continue
        }
        if num < min || num > max {
            fmt.Printf("❌ Число має бути між %d і %d\n", min, max)
            continue
        }
        return num
    }
}

// Підтвердження дії
func confirm(reader *bufio.Reader, prompt string) bool {
    fmt.Print(prompt + " (так/ні): ")
    input, _ := reader.ReadString('\n')
    input = strings.ToLower(strings.TrimSpace(input))
    return input == "так" || input == "yes" || input == "y"
}
```

---

## Завдання для практики

1. Переписати solution_1.go з інтерактивним вводом
2. Створити TODO менеджер з меню (без os.Args)
3. Додати валідацію для всіх полів
4. Створити програму з комбінованим підходом

---

**Рекомендація:** Для навчання почніть з `bufio.Reader` - це найуніверсальніший підхід! 🚀

