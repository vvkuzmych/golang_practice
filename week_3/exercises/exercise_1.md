# Вправа 1: Byte Encoder - Кодувальник байтів

## Ціль
Створити систему кодування/декодування даних з використанням byte операцій.

---

## Завдання

Створіть програму `byte_encoder.go`, яка:

1. Реалізує кілька способів кодування даних
2. Підтримує Hex, Base64, та власний ROT13
3. Може кодувати та декодувати текст
4. Виводить статистику про байти

---

## Вимоги

### Обов'язкові функції:

```go
// HexEncode кодує []byte в hex string
func HexEncode(data []byte) string

// HexDecode декодує hex string в []byte
func HexDecode(hexString string) ([]byte, error)

// Base64Encode кодує []byte в base64
func Base64Encode(data []byte) string

// Base64Decode декодує base64 в []byte
func Base64Decode(encoded string) ([]byte, error)

// ROT13 застосовує ROT13 шифр (для ASCII A-Z, a-z)
func ROT13(data []byte) []byte

// ByteStats повертає статистику про байти
func ByteStats(data []byte) map[string]int
```

### ROT13 Пояснення:
ROT13 - це простий шифр зсуву на 13 позицій в алфавіті:
- A → N, B → O, C → P, ..., M → Z, N → A, ...
- Тільки для літер A-Z та a-z
- Інші символи залишаються без змін

---

## Приклад використання

```go
func main() {
    message := "Hello, World!"
    
    // Hex
    hex := HexEncode([]byte(message))
    fmt.Printf("Hex: %s\n", hex)
    decoded, _ := HexDecode(hex)
    fmt.Printf("Decoded: %s\n", decoded)
    
    // Base64
    b64 := Base64Encode([]byte(message))
    fmt.Printf("Base64: %s\n", b64)
    
    // ROT13
    rot := ROT13([]byte(message))
    fmt.Printf("ROT13: %s\n", rot)
    fmt.Printf("ROT13 again: %s\n", ROT13(rot))  // назад!
    
    // Stats
    stats := ByteStats([]byte(message))
    fmt.Printf("Stats: %v\n", stats)
}
```

---

## Очікуваний вивід

```
=== Hex Encoding ===
Original: Hello, World!
Hex: 48656c6c6f2c20576f726c6421
Decoded: Hello, World!

=== Base64 Encoding ===
Original: Hello, World!
Base64: SGVsbG8sIFdvcmxkIQ==
Decoded: Hello, World!

=== ROT13 Cipher ===
Original: Hello, World!
ROT13: Uryyb, Jbeyq!
ROT13 twice: Hello, World!

=== Byte Statistics ===
Text: Hello, World!
Total bytes: 13
Unique bytes: 10
Most frequent: 'l' (3 times)
Letter count: 10
Digit count: 0
Space count: 1
```

---

## Підказки

### 1. Hex Encoding
```go
import "encoding/hex"

func HexEncode(data []byte) string {
    return hex.EncodeToString(data)
}
```

### 2. ROT13 Implementation
```go
func ROT13(data []byte) []byte {
    result := make([]byte, len(data))
    for i, b := range data {
        if b >= 'A' && b <= 'Z' {
            result[i] = 'A' + (b-'A'+13)%26
        } else if b >= 'a' && b <= 'z' {
            result[i] = 'a' + (b-'a'+13)%26
        } else {
            result[i] = b
        }
    }
    return result
}
```

### 3. Byte Statistics
```go
func ByteStats(data []byte) map[string]int {
    stats := map[string]int{
        "total": len(data),
        "unique": 0,
        "letters": 0,
        "digits": 0,
        "spaces": 0,
    }
    
    freq := make(map[byte]int)
    for _, b := range data {
        freq[b]++
        if (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') {
            stats["letters"]++
        }
        if b >= '0' && b <= '9' {
            stats["digits"]++
        }
        if b == ' ' {
            stats["spaces"]++
        }
    }
    stats["unique"] = len(freq)
    
    return stats
}
```

---

## Бонус завдання

1. **Caesar Cipher**:
   ```go
   func Caesar(data []byte, shift int) []byte
   ```
   Зсув на будь-яку кількість позицій

2. **XOR Cipher**:
   ```go
   func XORCipher(data []byte, key byte) []byte
   ```
   XOR кожного байту з ключем

3. **Byte Frequency Analysis**:
   ```go
   func FrequencyAnalysis(data []byte) map[byte]int
   ```
   Підрахунок частоти кожного байту

4. **Binary Representation**:
   ```go
   func ToBinary(data []byte) string
   ```
   Конвертація в binary string (01010101...)

5. **Checksum/Hash**:
   ```go
   func SimpleChecksum(data []byte) byte
   func CRC8(data []byte) byte
   ```

---

## Критерії оцінки

- ✅ Всі функції реалізовані
- ✅ Hex encode/decode працює
- ✅ Base64 encode/decode працює
- ✅ ROT13 коректно шифрує та дешифрує
- ✅ ByteStats повертає правильну статистику
- ✅ Обробка помилок (invalid hex, base64)
- ✅ Код чистий і зрозумілий

---

## Тестування

```go
// Test ROT13
input := "Hello"
encrypted := ROT13([]byte(input))
decrypted := ROT13(encrypted)
if string(decrypted) != input {
    fmt.Println("ROT13 failed!")
}

// Test Hex
data := []byte("Test")
hex := HexEncode(data)
decoded, _ := HexDecode(hex)
if string(decoded) != "Test" {
    fmt.Println("Hex failed!")
}
```

---

## Корисні пакети

- `encoding/hex` - hex кодування
- `encoding/base64` - base64 кодування
- `fmt` - форматування виводу
- `strings` - string операції

---

## Рішення

Рішення знаходиться в `solutions/solution_1.go`.

Спробуйте виконати завдання самостійно перед тим, як дивитись рішення!

---

## Навчальні цілі

Після виконання цієї вправи ви будете:
- ✅ Розуміти роботу з []byte
- ✅ Вміти кодувати/декодувати дані
- ✅ Знати різні методи кодування
- ✅ Вміти аналізувати байтові дані
- ✅ Розуміти ASCII маніпуляції

---

## Реальне застосування

Подібні техніки використовуються в:
- **Криптографія** - шифрування даних
- **Networking** - кодування пакетів
- **Storage** - збереження binary даних
- **APIs** - передача даних (base64)
- **File formats** - парсинг binary файлів

