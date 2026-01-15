# Швидкий старт - Тиждень 3

## 🚀 Як почати

### 1. Перейти в папку week_3
```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_3_rune_bytes
```

### 2. Прочитати README
```bash
cat README.md
```

---

## 📚 Порядок навчання

### День 1-2: Теорія

```bash
# 1. Byte Basics
cat theory/01_byte_basics.md

# 2. Rune & Unicode
cat theory/02_rune_unicode.md

# 3. iota & Enums
cat theory/03_iota_enums.md
```

### День 3-4: Практика

```bash
# 1. Byte Examples
cd practice/byte_examples
go run main.go

# 2. Rune Examples
cd ../rune_examples
go run main.go

# 3. Iota Examples
cd ../iota_examples
go run main.go

# 4. HTTP Examples (byte & rune в реальному застосуванні)
cd ../http_examples
go run main.go
```

### День 5-6: Вправи

```bash
# Прочитати завдання
cd ../../exercises
cat exercise_1.md  # Byte Encoder
cat exercise_2.md  # Unicode Counter
cat exercise_3.md  # Status System

# Перевірити рішення
cd ../solutions
go run solution_1.go
go run solution_2.go
go run solution_3.go
```

### День 7: Контроль знань

Відповісти на питання:

#### 1. byte vs rune
**Q: В чому різниця?**
- byte = uint8 (1 байт, 0-255)
- rune = int32 (4 байти, Unicode code point)

**Q: Приклад?**
```go
var b byte = 65     // 'A'
var r rune = 'П'    // Українська літера
```

#### 2. UTF-8 Encoding
**Q: Скільки байт займає "Привіт"?**
- 12 байт (кожна кирилична літера = 2 байти)
- len("Привіт") == 12
- utf8.RuneCountInString("Привіт") == 6

**Q: Як ітерувати?**
```go
// ❌ Погано - по байтах
for i := 0; i < len(s); i++ {
    b := s[i]  // byte
}

// ✅ Добре - по рунах
for i, r := range s {
    // r - це rune
}
```

#### 3. iota
**Q: Що таке iota?**
- Авто-інкремент константа
- Починається з 0
- Збільшується на 1 в кожному const блоці

**Q: Приклад?**
```go
const (
    Sunday = iota     // 0
    Monday            // 1
    Tuesday           // 2
)
```

---

## ⚡ Швидкі команди

### Експерименти з byte
```bash
cat > test_byte.go << 'EOF'
package main
import "fmt"

func main() {
    // byte - це uint8
    var b byte = 65
    fmt.Printf("byte: %d, char: %c\n", b, b)
    
    // string -> []byte
    s := "Hello"
    bytes := []byte(s)
    fmt.Printf("bytes: %v\n", bytes)
}
EOF

go run test_byte.go
rm test_byte.go
```

### Експерименти з rune
```bash
cat > test_rune.go << 'EOF'
package main
import (
    "fmt"
    "unicode/utf8"
)

func main() {
    s := "Привіт"
    
    // Довжина в байтах vs символах
    fmt.Printf("len: %d байт\n", len(s))
    fmt.Printf("runes: %d символів\n", utf8.RuneCountInString(s))
    
    // Ітерація по рунах
    for i, r := range s {
        fmt.Printf("%d: %c\n", i, r)
    }
}
EOF

go run test_rune.go
rm test_rune.go
```

### Експерименти з iota
```bash
cat > test_iota.go << 'EOF'
package main
import "fmt"

const (
    Monday = iota
    Tuesday
    Wednesday
)

const (
    Read = 1 << iota
    Write
    Execute
)

func main() {
    fmt.Printf("Days: %d %d %d\n", Monday, Tuesday, Wednesday)
    fmt.Printf("Permissions: %b %b %b\n", Read, Write, Execute)
}
EOF

go run test_iota.go
rm test_iota.go
```

---

## 🎯 Контрольний список

Після тижня 3 ви повинні:

### Теорія
- [ ] Розумію що таке byte (uint8)
- [ ] Розумію що таке rune (int32)
- [ ] Знаю різницю між byte і rune
- [ ] Розумію UTF-8 encoding
- [ ] Знаю що таке iota
- [ ] Можу створювати enum

### Практика
- [ ] Конвертую string ↔ []byte
- [ ] Правильно ітерую по Unicode string
- [ ] Обробляю українські тексти
- [ ] Створюю enum з iota
- [ ] Використовую bit flags
- [ ] Розумію len() vs RuneCount()

### Код
- [ ] Написав Byte Encoder
- [ ] Написав Unicode Counter
- [ ] Написав Status System з enum
- [ ] Можу пояснити свій код

---

## 💡 Підказки

### Коли використовувати byte?

✅ **Використовуйте byte коли:**
- Робота з ASCII текстом
- Binary data (файли, мережа)
- Один байт = один символ
- Потрібна швидкість і економія пам'яті

### Коли використовувати rune?

✅ **Використовуйте rune коли:**
- Unicode текст (не тільки ASCII)
- Українська, китайська, арабська мови
- Emoji і спеціальні символи
- Потрібна коректна робота з символами

### Коли використовувати iota?

✅ **Використовуйте iota коли:**
- Створюєте enum (послідовність констант)
- Bit flags (Read, Write, Execute)
- Status codes
- Послідовні значення (Monday, Tuesday...)

---

## 🌍 Приклади з українською

### Правильна обробка українського тексту
```go
package main

import (
    "fmt"
    "unicode/utf8"
)

func main() {
    text := "Слава Україні!"
    
    // ❌ Неправильно
    fmt.Println("Байтів:", len(text))  // 25 байт
    
    // ✅ Правильно
    count := utf8.RuneCountInString(text)
    fmt.Println("Символів:", count)  // 14 символів
    
    // ✅ Правильна ітерація
    for i, r := range text {
        fmt.Printf("Позиція %d: %c\n", i, r)
    }
}
```

---

## 📚 Додаткові ресурси

### Офіційна документація
- [Strings, bytes, runes and characters in Go](https://go.dev/blog/strings)
- [unicode/utf8 package](https://pkg.go.dev/unicode/utf8)
- [Constants in Go](https://go.dev/blog/constants)

### Корисні пакети
- `unicode/utf8` - UTF-8 utilities
- `unicode` - Unicode character properties
- `strings` - String manipulation
- `bytes` - Byte slice utilities

---

## 🚧 Поширені помилки та рішення

### Помилка 1: Неправильний підрахунок символів

```go
// ❌ Погано
s := "Київ"
length := len(s)  // 8 байт, не 4 символи!

// ✅ Добре
import "unicode/utf8"
length := utf8.RuneCountInString(s)  // 4 символи
```

### Помилка 2: Індексація string

```go
s := "Україна"

// ❌ Погано - отримуємо byte, не символ
first := s[0]  // byte, не 'У'!

// ✅ Добре - конвертуємо в []rune
runes := []rune(s)
first := runes[0]  // 'У'
```

### Помилка 3: Модифікація string через індекс

```go
s := "Hello"

// ❌ Погано - string immutable!
// s[0] = 'h'  // compilation error

// ✅ Добре - через []byte або []rune
bytes := []byte(s)
bytes[0] = 'h'
s = string(bytes)  // "hello"
```

### Помилка 4: iota в різних const блоках

```go
// ❌ Неправильне розуміння
const A = iota  // 0
const B = iota  // 0 (новий блок!)

// ✅ Правильно - один блок
const (
    A = iota  // 0
    B         // 1
)
```

---

## 🎓 Практичні сценарії

### Сценарій 1: Валідація email (ASCII)
```go
func isValidEmail(email string) bool {
    // Email зазвичай ASCII - можна використати byte
    bytes := []byte(email)
    hasAt := false
    for _, b := range bytes {
        if b == '@' {
            hasAt = true
        }
    }
    return hasAt
}
```

### Сценарій 2: Підрахунок українських літер
```go
func countUkrainianLetters(text string) int {
    count := 0
    for _, r := range text {
        if (r >= 'А' && r <= 'Я') || (r >= 'а' && r <= 'я') || 
           r == 'Є' || r == 'І' || r == 'Ї' || r == 'Ґ' ||
           r == 'є' || r == 'і' || r == 'ї' || r == 'ґ' {
            count++
        }
    }
    return count
}
```

### Сценарій 3: HTTP Status Codes з iota
```go
const (
    StatusOK = 200 + iota
    StatusCreated
    StatusAccepted
)

const (
    StatusBadRequest = 400 + iota
    StatusUnauthorized
    StatusForbidden
    StatusNotFound
)
```

### Сценарій 4: HTTP Request з українським ім'ям (byte & rune)
```go
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "unicode/utf8"
)

type User struct {
    Name string `json:"name"`
}

func handler(w http.ResponseWriter, r *http.Request) {
    // HTTP body - це []byte
    bodyBytes, _ := io.ReadAll(r.Body)
    defer r.Body.Close()
    
    var user User
    json.Unmarshal(bodyBytes, &user)
    
    // Перевірка довжини імені (в СИМВОЛАХ, не байтах!)
    nameLength := utf8.RuneCountInString(user.Name)
    if nameLength < 2 || nameLength > 50 {
        http.Error(w, "Name too short/long", http.StatusBadRequest)
        return
    }
    
    // ВАЖЛИВО: Content-Length в БАЙТАХ
    fmt.Printf("Name: %s\n", user.Name)
    fmt.Printf("  Characters: %d\n", nameLength)
    fmt.Printf("  Bytes: %d\n", len(user.Name))
    
    // Відповідь (також []byte)
    response := []byte(fmt.Sprintf("Привіт, %s!", user.Name))
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.Write(response)
}

// Приклад запиту:
// {"name": "Олександра"}  
//   - 10 символів (characters)
//   - 20 байт (bytes)
//   - Content-Length: 20
```

---

## ❓ Питання та відповіді

### Q: Чому len("Привіт") повертає 12, а не 6?
A: `len()` повертає кількість **байтів**, не символів. Кожна кирилична літера в UTF-8 займає 2 байти. 6 літер × 2 = 12 байтів.

### Q: Як отримати n-ий символ Unicode string?
A: Конвертуйте в `[]rune`: `runes := []rune(s); char := runes[n]`

### Q: Чи можна змінювати string?
A: Ні, strings immutable. Конвертуйте в `[]byte` або `[]rune`, змініть, потім назад в string.

### Q: Коли iota скидається?
A: В кожному новому `const` блоці iota починається з 0.

### Q: Як зберігати emoji?
A: Як звичайний string - Go підтримує UTF-8, emoji працюють автоматично.

### Q: Чи використовуються byte і rune в HTTP?
A: **Так!** 
- **byte**: HTTP request/response bodies завжди `[]byte`. `io.ReadAll()` повертає `[]byte`. Content-Length вимірюється в байтах.
- **rune**: Для валідації довжини імені користувача (не байтів!), підрахунку символів у формах, обробки Unicode в JSON/URL parameters.

**Приклад:**
```go
// Клієнт відправляє JSON з українським ім'ям
{"name": "Олена"}  // 5 символів (runes), 10 байт (bytes)

// Сервер читає body як []byte
bodyBytes, _ := io.ReadAll(r.Body)  // []byte

// Але перевіряє довжину в СИМВОЛАХ
nameLength := utf8.RuneCountInString(user.Name)  // 5, не 10!
```

Детальні приклади: `practice/http_examples/main.go`

---

## 🎉 Успіхів у навчанні!

**Пам'ятайте:**
- byte для ASCII і binary data
- rune для Unicode (українська, emoji)
- iota для enum і послідовних констант
- len() повертає байти, не символи
- range по string ітерує по runes

---

**Happy coding! 🚀**

