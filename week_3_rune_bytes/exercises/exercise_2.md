# Вправа 2: Unicode Text Analyzer - Аналізатор Unicode тексту

## Ціль
Створити програму для аналізу та обробки Unicode тексту, включаючи українську мову.

---

## Завдання

Створіть програму `unicode_analyzer.go`, яка:

1. Коректно обробляє українські тексти
2. Працює з emoji та спеціальними символами
3. Аналізує Unicode текст
4. Виконує різні маніпуляції з рунами

---

## Вимоги

### Обов'язкові функції:

```go
// CharCount повертає кількість символів (не байтів!)
func CharCount(text string) int

// UkrainianCount підраховує українські літери
func UkrainianCount(text string) int

// EmojiCount підраховує emoji
func EmojiCount(text string) int

// Reverse реверсує string (правильно для Unicode)
func Reverse(text string) string

// Substring витягує підстроку за індексами символів
func Substring(text string, start, end int) string

// RemoveAccents видаляє діакритичні знаки
func RemoveAccents(text string) string

// TextStats повертає детальну статистику
func TextStats(text string) map[string]int
```

---

## Приклад використання

```go
func main() {
    text := "Привіт, World! 👋🎉"
    
    // Підрахунок
    fmt.Printf("Символів: %d\n", CharCount(text))
    fmt.Printf("Байтів: %d\n", len(text))
    
    // Аналіз
    ukr := UkrainianCount(text)
    emoji := EmojiCount(text)
    fmt.Printf("Українських літер: %d\n", ukr)
    fmt.Printf("Emoji: %d\n", emoji)
    
    // Маніпуляції
    reversed := Reverse(text)
    fmt.Printf("Reversed: %s\n", reversed)
    
    // Статистика
    stats := TextStats(text)
    for k, v := range stats {
        fmt.Printf("%s: %d\n", k, v)
    }
}
```

---

## Очікуваний вивід

```
=== Базовий аналіз ===
Text: Привіт, World! 👋🎉
Символів: 18
Байтів: 30
Співвідношення: 1.67 байт/символ

=== Підрахунок ===
Українських літер: 6
Латинських літер: 5
Цифр: 0
Emoji: 2
Пробілів: 2
Розділових знаків: 2

=== Категорії Unicode ===
Letters: 11
Symbols: 2
Punctuation: 2
Spaces: 2

=== Операції ===
Original: Привіт, World! 👋🎉
Reversed: 🎉👋 !dlroW ,тівирП
Uppercase: ПРИВІТ, WORLD! 👋🎉
Lowercase: привіт, world! 👋🎉
Substring(0, 6): Привіт

=== Детальна статистика ===
Всього символів: 18
Унікальних символів: 17
Найчастіший символ: 'і' (2 рази)
UTF-8 розподіл:
  1 байт: 10 символів
  2 байти: 6 символів
  4 байти: 2 символів
```

---

## Підказки

### 1. Правильний підрахунок символів
```go
import "unicode/utf8"

func CharCount(text string) int {
    return utf8.RuneCountInString(text)
}
```

### 2. Українські літери
```go
func UkrainianCount(text string) int {
    count := 0
    for _, r := range text {
        if isUkrainian(r) {
            count++
        }
    }
    return count
}

func isUkrainian(r rune) bool {
    return (r >= 'А' && r <= 'Я') || (r >= 'а' && r <= 'я') ||
        r == 'Є' || r == 'І' || r == 'Ї' || r == 'Ґ' ||
        r == 'є' || r == 'і' || r == 'ї' || r == 'ґ'
}
```

### 3. Emoji Detection
```go
func EmojiCount(text string) int {
    count := 0
    for _, r := range text {
        if r >= 0x1F300 && r <= 0x1F9FF {  // emoji range
            count++
        }
    }
    return count
}
```

### 4. Правильний Reverse
```go
func Reverse(text string) string {
    runes := []rune(text)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}
```

### 5. Substring з символами
```go
func Substring(text string, start, end int) string {
    runes := []rune(text)
    if start < 0 {
        start = 0
    }
    if end > len(runes) {
        end = len(runes)
    }
    if start >= end {
        return ""
    }
    return string(runes[start:end])
}
```

---

## Бонус завдання

1. **Word Count** (з Unicode):
   ```go
   func WordCount(text string) int
   ```
   Підрахунок слів (складно з Unicode!)

2. **Language Detection**:
   ```go
   func DetectLanguage(text string) string
   ```
   Визначити мову тексту

3. **Transliteration**:
   ```go
   func UkrainianToLatin(text string) string
   ```
   Привіт → Pryvit

4. **UTF-8 Size Analysis**:
   ```go
   func UTF8Distribution(text string) map[int]int
   ```
   Скільки символів займає 1,2,3,4 байти

5. **Normalize Unicode**:
   ```go
   func NormalizeUnicode(text string) string
   ```
   NFD, NFC, NFKD, NFKC normalization

---

## Критерії оцінки

- ✅ CharCount правильно підраховує символи
- ✅ Коректно обробляє українську мову
- ✅ Правильно визначає emoji
- ✅ Reverse працює з Unicode
- ✅ Substring працює по символах, не байтах
- ✅ TextStats повертає всю статистику
- ✅ Код працює з багатомовними текстами

---

## Тестові дані

```go
testCases := []string{
    "Привіт",                    // українська
    "Hello",                      // англійська
    "你好",                       // китайська
    "مرحبا",                     // арабська
    "👋🎉🚀",                     // emoji
    "Привіт, World! 你好 👋",    // змішаний
}
```

---

## Важливі моменти

### ❌ Неправильно:
```go
// НЕ використовуйте len() для символів!
count := len("Привіт")  // 12 (байти), не 6!

// НЕ індексуйте string напряму!
first := "Привіт"[0]  // byte, не 'П'!

// НЕ використовуйте for по індексу!
for i := 0; i < len(text); i++ {
    char := text[i]  // byte, не rune!
}
```

### ✅ Правильно:
```go
// Використовуйте utf8.RuneCountInString
count := utf8.RuneCountInString("Привіт")  // 6

// Конвертуйте в []rune для індексації
runes := []rune("Привіт")
first := runes[0]  // 'П'

// Використовуйте range
for i, r := range text {
    // i - byte position, r - rune
}
```

---

## Корисні пакети

- `unicode/utf8` - UTF-8 utilities
- `unicode` - Unicode properties
- `strings` - String functions
- `golang.org/x/text/unicode/norm` - Unicode normalization

---

## Рішення

Рішення знаходиться в `solutions/solution_2.go`.

Спробуйте виконати завдання самостійно перед тим, як дивитись рішення!

---

## Навчальні цілі

Після виконання цієї вправи ви будете:
- ✅ Розуміти різницю між byte і rune
- ✅ Коректно обробляти Unicode тексти
- ✅ Знати як працювати з українською мовою
- ✅ Вміти аналізувати багатомовні тексти
- ✅ Розуміти UTF-8 encoding

---

## Реальне застосування

Подібні функції використовуються в:
- **Text editors** - правильна навігація по тексту
- **Search engines** - індексація багатомовних текстів
- **Social media** - обробка emoji та Unicode
- **Translation apps** - аналіз мов
- **Content moderation** - фільтрація тексту

