# STDIN - Динамічний ввід під час запуску

## Що таке STDIN?

**STDIN** (Standard Input) - стандартний потік вводу в Unix/Linux системах.

### 3 способи роботи з вводом:

| Спосіб | Як працює | Приклад |
|--------|-----------|---------|
| **os.Args** | Аргументи при запуску | `program arg1 arg2` |
| **Інтерактивний** | Діалог з користувачем | Програма запитує → Відповідь |
| **STDIN** | Потік даних | `echo "data" \| program` |

---

## STDIN - Найгнучкіший підхід!

### Переваги ✅

1. **Працює з pipe**
   ```bash
   echo "Іван" | program
   ```

2. **Працює з файлами**
   ```bash
   program < input.txt
   cat input.txt | program
   ```

3. **Працює інтерактивно**
   ```bash
   program
   # Вводите вручну
   ```

4. **Працює в скриптах**
   ```bash
   for name in Іван Марія Петро; do
       echo $name | program
   done
   ```

5. **Працює з heredoc**
   ```bash
   program << EOF
   Іван
   25
   EOF
   ```

---

## Як читати з STDIN в Go

### Метод 1: bufio.Scanner (Рекомендовано) ⭐⭐⭐⭐⭐

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func main() {
    // Створюємо scanner для STDIN
    scanner := bufio.NewScanner(os.Stdin)
    
    // Читаємо рядок за рядком
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        fmt.Printf("Прочитано: %s\n", line)
    }
    
    // Перевірка на помилки
    if err := scanner.Err(); err != nil {
        fmt.Fprintf(os.Stderr, "Помилка: %v\n", err)
    }
}
```

**Особливості:**
- ✅ Найпростіший API
- ✅ Автоматично обробляє \n
- ✅ Працює з будь-яким джерелом (pipe, файл, клавіатура)
- ✅ Можна читати по рядках

---

### Метод 2: bufio.Reader

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
    
    // Читаємо до '\n'
    line, err := reader.ReadString('\n')
    if err != nil {
        fmt.Fprintf(os.Stderr, "Помилка: %v\n", err)
        return
    }
    
    line = strings.TrimSpace(line)
    fmt.Printf("Прочитано: %s\n", line)
}
```

---

### Метод 3: io.ReadAll (Читає все одразу)

```go
package main

import (
    "fmt"
    "io"
    "os"
)

func main() {
    // Читаємо весь STDIN одразу
    data, err := io.ReadAll(os.Stdin)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Помилка: %v\n", err)
        return
    }
    
    fmt.Printf("Прочитано %d байт\n", len(data))
    fmt.Printf("Дані: %s\n", string(data))
}
```

---

## STDOUT vs STDERR

**Важливо!** При роботі з STDIN потрібно розрізняти:

| Потік | Призначення | Використання в Go |
|-------|-------------|-------------------|
| **STDIN** | Ввід даних | `os.Stdin` |
| **STDOUT** | Результат роботи | `fmt.Println()`, `os.Stdout` |
| **STDERR** | Повідомлення/помилки | `fmt.Fprintln(os.Stderr, ...)` |

### Чому це важливо?

```go
// ❌ ПОГАНО - підказки в STDOUT
fmt.Println("Введіть ім'я:")
// Якщо вивід перенаправлений, підказка потрапить в результат!

// ✅ ДОБРЕ - підказки в STDERR
fmt.Fprintln(os.Stderr, "Введіть ім'я:")
// Підказка йде в STDERR, результат в STDOUT
```

**Приклад проблеми:**
```bash
# ПОГАНО
echo "Іван" | program > output.txt
# В output.txt буде: "Введіть ім'я:\nІван"

# ДОБРЕ
echo "Іван" | program > output.txt
# В output.txt буде тільки: "Іван"
# "Введіть ім'я:" пішло в термінал (STDERR)
```

---

## Повний приклад з правильним використанням

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
    scanner := bufio.NewScanner(os.Stdin)
    
    // Підказки в STDERR (не потраплять в pipe/файл)
    fmt.Fprintln(os.Stderr, "Калькулятор (через STDIN)")
    fmt.Fprintln(os.Stderr, "Введіть перше число:")
    
    if !scanner.Scan() {
        fmt.Fprintln(os.Stderr, "Помилка читання")
        os.Exit(1)
    }
    
    num1, err := strconv.ParseFloat(strings.TrimSpace(scanner.Text()), 64)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Помилка: не є числом\n")
        os.Exit(1)
    }
    
    fmt.Fprintln(os.Stderr, "Введіть операцію (+, -, *, /):")
    if !scanner.Scan() {
        fmt.Fprintln(os.Stderr, "Помилка читання")
        os.Exit(1)
    }
    
    op := strings.TrimSpace(scanner.Text())
    
    fmt.Fprintln(os.Stderr, "Введіть друге число:")
    if !scanner.Scan() {
        fmt.Fprintln(os.Stderr, "Помилка читання")
        os.Exit(1)
    }
    
    num2, err := strconv.ParseFloat(strings.TrimSpace(scanner.Text()), 64)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Помилка: не є числом\n")
        os.Exit(1)
    }
    
    // Обчислення
    var result float64
    switch op {
    case "+":
        result = num1 + num2
    case "-":
        result = num1 - num2
    case "*":
        result = num1 * num2
    case "/":
        result = num1 / num2
    }
    
    // Результат в STDOUT (чистий вивід для pipe)
    fmt.Printf("%.2f\n", result)
    
    // Або детальний результат
    // fmt.Printf("%.2f %s %.2f = %.2f\n", num1, op, num2, result)
}
```

---

## Способи використання STDIN

### 1. Інтерактивно (з клавіатури)

```bash
go run program.go
# Вводите дані вручну
```

### 2. Через echo

```bash
echo "Іван" | go run program.go
echo -e "Іван\n25" | go run program.go
```

### 3. Через printf

```bash
printf "Марія\n22\n" | go run program.go
```

### 4. Heredoc

```bash
go run program.go << EOF
Петро
30
EOF
```

### 5. З файлу (redirect)

```bash
go run program.go < input.txt
```

### 6. Через cat

```bash
cat input.txt | go run program.go
```

### 7. З іншої програми

```bash
ls -la | grep ".go" | go run program.go
```

### 8. В скриптах

```bash
#!/bin/bash
for name in Іван Марія Петро; do
    echo $name | go run greet.go
done
```

---

## Приклади з реального життя

### Приклад 1: Фільтр (як grep)

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Fprintln(os.Stderr, "Використання: filter <шаблон>")
        os.Exit(1)
    }
    
    pattern := os.Args[1]
    scanner := bufio.NewScanner(os.Stdin)
    
    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, pattern) {
            fmt.Println(line)
        }
    }
}
```

**Використання:**
```bash
cat file.txt | go run filter.go "error"
ls -la | go run filter.go ".go"
```

---

### Приклад 2: Підрахунок рядків (як wc -l)

```go
package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    count := 0
    
    for scanner.Scan() {
        count++
    }
    
    fmt.Println(count)
}
```

**Використання:**
```bash
cat file.txt | go run count.go
echo -e "line1\nline2\nline3" | go run count.go  # → 3
```

---

### Приклад 3: Перетворення (uppercase)

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    
    for scanner.Scan() {
        line := scanner.Text()
        fmt.Println(strings.ToUpper(line))
    }
}
```

**Використання:**
```bash
echo "hello world" | go run uppercase.go  # → HELLO WORLD
```

---

### Приклад 4: JSON обробка

```go
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "os"
)

type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func main() {
    data, _ := io.ReadAll(os.Stdin)
    
    var user User
    if err := json.Unmarshal(data, &user); err != nil {
        fmt.Fprintf(os.Stderr, "Помилка JSON: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Printf("Привіт, %s! Вік: %d\n", user.Name, user.Age)
}
```

**Використання:**
```bash
echo '{"name":"Іван","age":25}' | go run json_reader.go
cat user.json | go run json_reader.go
```

---

## Поєднання підходів

### Універсальна програма (os.Args + STDIN)

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func main() {
    var name string
    
    // Якщо є аргументи - використовуємо CLI
    if len(os.Args) >= 2 {
        name = os.Args[1]
        fmt.Fprintln(os.Stderr, "📝 Режим: CLI")
    } else {
        // Інакше - читаємо з STDIN
        fmt.Fprintln(os.Stderr, "📝 Режим: STDIN")
        fmt.Fprintln(os.Stderr, "Введіть ім'я:")
        
        scanner := bufio.NewScanner(os.Stdin)
        if scanner.Scan() {
            name = strings.TrimSpace(scanner.Text())
        }
    }
    
    // Результат
    fmt.Printf("Привіт, %s! 👋\n", name)
}
```

**Використання:**
```bash
# CLI
go run program.go Іван

# STDIN (інтерактивно)
go run program.go
# Введіть ім'я: Марія

# STDIN (pipe)
echo "Петро" | go run program.go
```

---

## Відмінності: Інтерактивний vs STDIN

| Характеристика | Інтерактивний | STDIN |
|----------------|---------------|-------|
| **Діалог** | ✅ Так | ⚠️ Опційно |
| **Pipe** | ❌ Складно | ✅ Легко |
| **Файли** | ❌ Ні | ✅ Так |
| **Скрипти** | ❌ Складно | ✅ Легко |
| **Підказки** | В STDOUT | В STDERR |
| **Цикл** | В коді | Зовні (bash) |

---

## Best Practices

### ✅ DO

1. **Підказки в STDERR**
   ```go
   fmt.Fprintln(os.Stderr, "Введіть дані:")
   ```

2. **Результати в STDOUT**
   ```go
   fmt.Println(result)
   ```

3. **Перевірка помилок**
   ```go
   if err := scanner.Err(); err != nil {
       fmt.Fprintf(os.Stderr, "Помилка: %v\n", err)
       os.Exit(1)
   }
   ```

4. **TrimSpace для вводу**
   ```go
   input := strings.TrimSpace(scanner.Text())
   ```

### ❌ DON'T

1. **Не плутати STDOUT і STDERR**
   ```go
   // ❌ ПОГАНО
   fmt.Println("Введіть ім'я:")  // В pipe потрапить!
   
   // ✅ ДОБРЕ
   fmt.Fprintln(os.Stderr, "Введіть ім'я:")
   ```

2. **Не забувати про помилки**
   ```go
   // ❌ ПОГАНО
   scanner.Scan()
   data := scanner.Text()
   
   // ✅ ДОБРЕ
   if scanner.Scan() {
       data := scanner.Text()
   }
   if err := scanner.Err(); err != nil {
       // обробка
   }
   ```

---

## Резюме

| Метод | Коли використовувати |
|-------|---------------------|
| **os.Args** | CLI утиліти для розробників |
| **Інтерактивний** | Програми з складним діалогом |
| **STDIN** | ✅ **Універсальні утиліти** (pipe, файли, скрипти) |

**Рекомендація:** Для максимальної гнучкості використовуйте **STDIN** з підтримкою інтерактивного режиму!

---

## Завдання

1. Переробити solution_1.go на STDIN
2. Створити фільтр для тексту (як grep)
3. Створити програму для підрахунку слів
4. Створити конвертер JSON → CSV через STDIN

