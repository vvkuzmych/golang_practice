# Task 3: Word Frequency Counter

**Level:** Middle  
**Time:** 15 minutes  
**Language:** Ruby, Go, or JavaScript  

---

## 📝 Task

Напиши функцію, яка приймає текст і повертає топ N найбільш вживаних слів.

---

## 📥 Input

```
text (String) - текст для аналізу
n (Integer) - кількість топ слів
```

---

## 📤 Output

```
Array of [word, count] pairs - топ N слів з їх частотою
```

---

## 💡 Examples

```ruby
text = "the quick brown fox jumps over the lazy dog the fox"

top_words(text, 3)
# => [["the", 3], ["fox", 2], ["quick", 1]]
# або [["the", 3], ["fox", 2], ["brown", 1]]  (порядок однакової частоти не важливий)

top_words("Hello world! Hello Ruby. Ruby is great.", 2)
# => [["hello", 2], ["ruby", 2]]

top_words("", 5)
# => []
```

---

## ✅ Requirements

- Case insensitive ("Hello" === "hello")
- Ignore punctuation (видали `,`, `.`, `!`, `?`, тощо)
- Ignore single character words
- Якщо слів з однаковою частотою, повернути будь-які
- Якщо в тексті менше N слів, повернути всі
- Порядок: від найбільшої частоти до найменшої

---

## 🎯 Test Cases

```ruby
# Test 1: Normal text
text = "the quick brown fox jumps over the lazy dog the fox"
top_words(text, 3)
# => [["the", 3], ["fox", 2], ...] (one more word with count 1)

# Test 2: Case insensitive
text = "Hello world! Hello Ruby. Ruby is great."
top_words(text, 2)
# => [["hello", 2], ["ruby", 2]]

# Test 3: Empty text
top_words("", 5)
# => []

# Test 4: More N than words
text = "hello world"
top_words(text, 10)
# => [["hello", 1], ["world", 1]]

# Test 5: With punctuation
text = "hello, hello! world. world?"
top_words(text, 2)
# => [["hello", 2], ["world", 2]]
```

---

## 💡 Hints

- Split text into words
- Clean punctuation
- Count frequency (HashMap)
- Sort by frequency
- Take top N

---

**Good luck!** 🚀
