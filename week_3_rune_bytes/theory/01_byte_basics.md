# Byte Basics - Основи роботи з байтами

## Що таке byte?

**byte** - це аліас для `uint8` (беззнаковий 8-бітний цілий тип)

```go
type byte = uint8
```

**Діапазон:** 0 до 255 (2^8 - 1)

---

## 🔄 byte vs rune - Ключова різниця

### Швидке порівняння

| Характеристика | byte | rune |
|----------------|------|------|
| **Тип** | `uint8` | `int32` |
| **Розмір** | 1 байт | 4 байти |
| **Діапазон** | 0-255 | 0-1,114,111 (Unicode) |
| **Використання** | ASCII, binary data | Unicode, багатомовні тексти |
| **Приклад** | `'A'`, `65` | `'П'`, `'你'`, `'👋'` |

### ⚠️ Критична різниця

```go
// ASCII текст - byte працює
text1 := "Hello"
fmt.Println(len(text1))                    // 5 байт = 5 символів ✅

// Українська - потрібен rune!
text2 := "Привіт"
fmt.Println(len(text2))                    // 12 байт ≠ 6 символів ❌
fmt.Println(utf8.RuneCountInString(text2)) // 6 символів ✅
```

### 🎯 Коли використовувати byte?

✅ **Використовуйте byte коли:**
- Робота з ASCII текстом (A-Z, 0-9)
- Binary data (файли, мережа)
- Кодування (hex, base64)
- Криптографія
- Один символ = один байт
- HTTP request/response bodies
- Файлові операції

**Приклад:**
```go
// ✅ ДОБРЕ для ASCII
data := []byte("GET /api/users HTTP/1.1")
password := []byte("secret123")
hexData := hex.EncodeToString([]byte{0xFF, 0xAB})
```

### 🎯 Коли використовувати rune?

✅ **Використовуйте rune коли:**
- Unicode текст (не тільки ASCII)
- Українська мова (Привіт, Київ)
- Інші мови (中文, العربية, 日本語)
- Emoji та спеціальні символи (👋, 🎉)
- Потрібна коректна робота з символами
- Підрахунок символів (не байтів!)
- Validація довжини імен користувачів

**Приклад:**
```go
// ✅ ДОБРЕ для Unicode
name := "Олександра"
for _, r := range name {  // ітерація по рунах
    fmt.Printf("%c ", r)
}

charCount := utf8.RuneCountInString(name)  // 10 символів
if charCount < 2 || charCount > 50 {
    // валідація довжини імені
}
```

### 📊 Реальний приклад: HTTP запит

```go
// HTTP використовує []byte
func handleUser(w http.ResponseWriter, r *http.Request) {
    // 1. Читаємо body як []byte
    bodyBytes, _ := io.ReadAll(r.Body)  // []byte
    
    // 2. Парсимо JSON
    var user User
    json.Unmarshal(bodyBytes, &user)
    
    // 3. Перевіряємо ім'я (використовуємо rune!)
    nameLength := utf8.RuneCountInString(user.Name)  // rune count
    if nameLength < 2 || nameLength > 50 {
        http.Error(w, "Invalid name length", http.StatusBadRequest)
        return
    }
    
    // 4. Відправляємо відповідь як []byte
    response := []byte(fmt.Sprintf("Привіт, %s!", user.Name))
    w.Write(response)  // []byte
}
```

### 💡 Правило великого пальця

```
📝 Якщо text може містити:
   - Українську ✓
   - Китайську ✓
   - Emoji ✓
   → використовуйте RUNE

📦 Якщо data це:
   - Файли ✓
   - Network packets ✓
   - Binary data ✓
   - HTTP bodies ✓
   → використовуйте BYTE
```

---

## Основні концепції

### 1. byte представляє один байт даних

```go
package main

import "fmt"

func main() {
    var b byte = 65
    
    fmt.Printf("Значення: %d\n", b)      // 65
    fmt.Printf("Символ: %c\n", b)        // A
    fmt.Printf("Бінарне: %b\n", b)       // 1000001
    fmt.Printf("Hex: %x\n", b)           // 41
}
```

---

## String ↔ []byte конверсія

### String to []byte

```go
package main

import "fmt"

func main() {
    s := "Hello"
    
    // Конвертація в байтовий масив
    bytes := []byte(s)
    
    fmt.Printf("String: %s\n", s)
    fmt.Printf("Bytes: %v\n", bytes)      // [72 101 108 108 111]
    fmt.Printf("Довжина: %d\n", len(bytes))
    
    // Зміна байтів
    bytes[0] = 72  // 'H'
    bytes[0] = 104 // 'h'
    
    // Назад в string
    modified := string(bytes)
    fmt.Printf("Modified: %s\n", modified)  // hello
}
```

### []byte to String

```go
bytes := []byte{72, 101, 108, 108, 111}
s := string(bytes)
fmt.Println(s)  // Hello
```

---

## Робота з ASCII

```go
package main

import "fmt"

func main() {
    // ASCII символи
    letters := []byte{'A', 'B', 'C'}
    
    for i, b := range letters {
        fmt.Printf("%d: %c (decimal: %d)\n", i, b, b)
    }
    
    // Перевірка діапазону
    var char byte = 'Z'
    if char >= 'A' && char <= 'Z' {
        fmt.Println("Велика літера")
    }
    
    // Перетворення uppercase → lowercase
    lowercase := char + 32  // 'Z' -> 'z'
    fmt.Printf("%c -> %c\n", char, lowercase)
}
```

---

## Binary Data

```go
package main

import (
    "encoding/binary"
    "fmt"
)

func main() {
    // Число в байти
    var num uint32 = 1000
    bytes := make([]byte, 4)
    binary.LittleEndian.PutUint32(bytes, num)
    
    fmt.Printf("Number: %d\n", num)
    fmt.Printf("Bytes: %v\n", bytes)
    
    // Байти в число
    decoded := binary.LittleEndian.Uint32(bytes)
    fmt.Printf("Decoded: %d\n", decoded)
}
```

---

## Hex Encoding

```go
package main

import (
    "encoding/hex"
    "fmt"
)

func main() {
    data := []byte("Hello")
    
    // Bytes -> Hex string
    hexString := hex.EncodeToString(data)
    fmt.Printf("Hex: %s\n", hexString)  // 48656c6c6f
    
    // Hex string -> Bytes
    decoded, _ := hex.DecodeString(hexString)
    fmt.Printf("Decoded: %s\n", string(decoded))  // Hello
}
```

---

## Читання файлів як байти

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // Записати байти у файл
    data := []byte("Hello, World!")
    os.WriteFile("test.txt", data, 0644)
    
    // Прочитати файл як байти
    content, err := os.ReadFile("test.txt")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("File content: %s\n", string(content))
    fmt.Printf("Bytes: %v\n", content)
    
    // Cleanup
    os.Remove("test.txt")
}
```

---

## Байтові операції

```go
package main

import (
    "bytes"
    "fmt"
)

func main() {
    b1 := []byte("Hello")
    b2 := []byte("World")
    
    // Порівняння
    if bytes.Equal(b1, b2) {
        fmt.Println("Однакові")
    } else {
        fmt.Println("Різні")
    }
    
    // Об'єднання
    result := bytes.Join([][]byte{b1, b2}, []byte(" "))
    fmt.Printf("Joined: %s\n", string(result))  // Hello World
    
    // Пошук
    if bytes.Contains(b1, []byte("ll")) {
        fmt.Println("Містить 'll'")
    }
    
    // Заміна
    replaced := bytes.Replace(b1, []byte("l"), []byte("L"), -1)
    fmt.Printf("Replaced: %s\n", string(replaced))  // HeLLo
}
```

---

## Буферизація

```go
package main

import (
    "bytes"
    "fmt"
)

func main() {
    var buf bytes.Buffer
    
    // Запис
    buf.Write([]byte("Hello"))
    buf.WriteString(" ")
    buf.WriteByte('W')
    buf.WriteString("orld")
    
    // Читання
    fmt.Printf("Buffer: %s\n", buf.String())
    fmt.Printf("Bytes: %v\n", buf.Bytes())
    fmt.Printf("Length: %d\n", buf.Len())
}
```

---

## Base64 Encoding

```go
package main

import (
    "encoding/base64"
    "fmt"
)

func main() {
    data := []byte("Hello, World!")
    
    // Encode
    encoded := base64.StdEncoding.EncodeToString(data)
    fmt.Printf("Base64: %s\n", encoded)
    
    // Decode
    decoded, _ := base64.StdEncoding.DecodeString(encoded)
    fmt.Printf("Decoded: %s\n", string(decoded))
}
```

---

## Практичні приклади

### Приклад 1: Підрахунок частоти байтів

```go
func countBytes(data []byte) map[byte]int {
    counts := make(map[byte]int)
    for _, b := range data {
        counts[b]++
    }
    return counts
}

func main() {
    text := "Hello, World!"
    counts := countBytes([]byte(text))
    
    for b, count := range counts {
        fmt.Printf("'%c' (%d): %d разів\n", b, b, count)
    }
}
```

### Приклад 2: XOR шифрування

```go
func xorEncrypt(data []byte, key byte) []byte {
    result := make([]byte, len(data))
    for i, b := range data {
        result[i] = b ^ key
    }
    return result
}

func main() {
    message := []byte("Secret")
    key := byte(123)
    
    // Encrypt
    encrypted := xorEncrypt(message, key)
    fmt.Printf("Encrypted: %v\n", encrypted)
    
    // Decrypt (XOR двічі повертає оригінал)
    decrypted := xorEncrypt(encrypted, key)
    fmt.Printf("Decrypted: %s\n", string(decrypted))
}
```

### Приклад 3: Checksum

```go
func simpleChecksum(data []byte) byte {
    var sum byte
    for _, b := range data {
        sum += b
    }
    return sum
}

func main() {
    data := []byte("Hello")
    checksum := simpleChecksum(data)
    fmt.Printf("Checksum: %d\n", checksum)
}
```

---

## 🎯 Практичне рішення: byte vs rune

### Дерево рішень

```
Потрібно обробити text/data?
│
├─ Це binary data? (файли, мережа, зображення)
│  └─ ✅ Використовуй []byte
│
├─ Це текст тільки ASCII? (A-Z, 0-9, без кирилиці)
│  └─ ✅ Можна []byte
│
├─ Текст може містити:
│  ├─ Українську? (Привіт, Київ)
│  ├─ Інші мови? (中文, العربية)
│  ├─ Emoji? (👋, 🎉)
│  └─ ✅ Використовуй rune / string з utf8.RuneCountInString()
│
└─ HTTP request/response?
   ├─ Body (читання/запис) → []byte
   └─ Валідація контенту → rune
```

### ✅ Використовуйте byte коли:

#### 1. Binary Data
```go
// Читання файлу
data, _ := os.ReadFile("image.png")  // []byte

// HTTP body
bodyBytes, _ := io.ReadAll(r.Body)  // []byte

// Network packet
packet := []byte{0xFF, 0xAB, 0xCD, 0xEF}
```

#### 2. Encoding/Decoding
```go
// Hex encoding
hexStr := hex.EncodeToString([]byte{255, 171, 205})

// Base64
b64 := base64.StdEncoding.EncodeToString([]byte("data"))

// JSON
jsonBytes, _ := json.Marshal(data)  // []byte
```

#### 3. Cryptography
```go
// Hashing
hash := sha256.Sum256([]byte("password"))

// Encryption
ciphertext := encrypt([]byte("secret"))
```

#### 4. ASCII Protocol
```go
// HTTP request line
request := []byte("GET /api HTTP/1.1\r\n")

// SMTP command
cmd := []byte("MAIL FROM:<user@example.com>\r\n")
```

#### 5. File Operations
```go
// Write to file
os.WriteFile("data.bin", []byte{1, 2, 3}, 0644)

// Read from file
data, _ := os.ReadFile("config.bin")
```

### ✅ Використовуйте rune коли:

#### 1. Підрахунок символів
```go
// ❌ WRONG - counts bytes
len("Привіт")  // 12

// ✅ CORRECT - counts characters
utf8.RuneCountInString("Привіт")  // 6
```

#### 2. Валідація довжини
```go
// Перевірка імені користувача
func validateName(name string) bool {
    charCount := utf8.RuneCountInString(name)
    return charCount >= 2 && charCount <= 50  // символи, не байти!
}
```

#### 3. Ітерація по символах
```go
// ✅ CORRECT - iterates by runes
for _, r := range "Привіт👋" {
    fmt.Printf("%c ", r)  // П р и в і т 👋
}
```

#### 4. Substring операції
```go
// ✅ CORRECT way
runes := []rune("Привіт")
first3 := string(runes[0:3])  // "При"
```

#### 5. Character manipulation
```go
// Перевірка українських літер
func isUkrainian(r rune) bool {
    return (r >= 'А' && r <= 'Я') || (r >= 'а' && r <= 'я') ||
        r == 'Є' || r == 'І' || r == 'Ї' || r == 'Ґ'
}
```

### ❌ НЕ використовуйте byte для:

1. **Unicode тексту**
   ```go
   // ❌ WRONG
   text := []byte("Привіт")
   firstChar := text[0]  // 208, не 'П'!
   ```

2. **Підрахунку символів**
   ```go
   // ❌ WRONG
   count := len([]byte("Слава Україні"))  // 25 байт, не 14 символів!
   ```

3. **Індексації багатобайтових символів**
   ```go
   // ❌ WRONG
   s := "Київ"
   char := s[0]  // byte, не 'К'!
   ```

### ❌ НЕ використовуйте rune для:

1. **Binary Data**
   ```go
   // ❌ WRONG - rune для text, не для binary
   data := []rune{0xFF, 0xAB}  // Неправильно!
   
   // ✅ CORRECT
   data := []byte{0xFF, 0xAB}
   ```

2. **HTTP Bodies**
   ```go
   // ❌ WRONG
   // io.ReadAll повертає []byte, не []rune
   
   // ✅ CORRECT
   bodyBytes, _ := io.ReadAll(r.Body)
   ```

3. **File I/O**
   ```go
   // ❌ WRONG
   os.WriteFile("data", []rune{...}, 0644)
   
   // ✅ CORRECT
   os.WriteFile("data", []byte{...}, 0644)
   ```

### 📋 Швидка довідка

| Задача | Тип | Приклад |
|--------|-----|---------|
| HTTP request body | `[]byte` | `io.ReadAll(r.Body)` |
| HTTP response write | `[]byte` | `w.Write([]byte("OK"))` |
| Валідація імені | `rune` | `utf8.RuneCountInString(name)` |
| File read/write | `[]byte` | `os.ReadFile()` |
| JSON marshal | `[]byte` | `json.Marshal()` |
| Підрахунок emoji | `rune` | `for _, r := range text` |
| Binary protocol | `[]byte` | `packet := []byte{0xFF}` |
| Unicode substring | `rune` | `[]rune(text)[0:3]` |
| Hex encoding | `[]byte` | `hex.EncodeToString()` |
| Ukrainian text | `rune` | `isUkrainian(r rune)` |

### 🎓 Запам'ятайте

```
byte → Binary, Files, Network, HTTP bodies, ASCII
rune → Unicode, Multilingual, Characters, User input

HTTP працює з BYTES
Користувачі думають СИМВОЛАМИ (runes)
```

---

## Резюме

| Операція | Код |
|----------|-----|
| String → []byte | `bytes := []byte(s)` |
| []byte → String | `s := string(bytes)` |
| ASCII перевірка | `if b >= 'A' && b <= 'Z'` |
| Hex encoding | `hex.EncodeToString(bytes)` |
| Base64 | `base64.StdEncoding.EncodeToString(bytes)` |
| Порівняння | `bytes.Equal(b1, b2)` |
| Об'єднання | `bytes.Join([][]byte{b1, b2}, sep)` |

---

## Корисні пакети

- `bytes` - операції з байтовими слайсами
- `encoding/hex` - hex кодування
- `encoding/base64` - base64 кодування
- `encoding/binary` - binary data
- `io` - читання/запис байтів

---

## Завдання для практики

1. Створити функцію для конвертації hex string в []byte
2. Написати простий XOR шифр
3. Реалізувати підрахунок CRC
4. Створити функцію для зчитування binary файлу
5. Написати encoder/decoder для власного формату

---

## Наступний крок

Тепер, коли ви розумієте byte, перейдіть до вивчення **rune** для роботи з Unicode!

