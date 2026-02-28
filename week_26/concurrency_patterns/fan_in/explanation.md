# Чому не можна прибрати `_` у fan_in паттерні?

## ❌ Проблемний код (БАГ!)

```go
func MergeChannels[T any](channels ...<-chan T) <-chan T {
    var wg sync.WaitGroup
    wg.Add(len(channels))
    outputCh := make(chan T)
    
    // ❌ НЕПРАВИЛЬНО - без _
    for channel := range channels {
        go func() {
            defer wg.Done()
            for value := range channel {  // БАГ: channel захоплюється неправильно!
                outputCh <- value
            }
        }()
    }
    
    // ...
}
```

## 🐛 Що станеться?

**Проблема: Loop Variable Capture (Захоплення змінної циклу)**

### Покрокове пояснення:

1. **Цикл запускає 3 горутини:**
   ```go
   channels = [channel1, channel2, channel3]
   
   Ітерація 1: channel = channel1 → запускаємо goroutine #1
   Ітерація 2: channel = channel2 → запускаємо goroutine #2
   Ітерація 3: channel = channel3 → запускаємо goroutine #3
   ```

2. **Змінна `channel` - це ОДНА змінна, яка перевикористовується:**
   - Горутини НЕ копіюють значення `channel`
   - Всі горутини посилаються на **ту саму змінну** `channel`
   - Цикл працює **швидше**, ніж запускаються горутини

3. **Що бачать горутини:**
   ```
   Цикл закінчується → channel = channel3 (останнє значення)
   
   Goroutine #1: читає channel → бачить channel3 ❌
   Goroutine #2: читає channel → бачить channel3 ❌
   Goroutine #3: читає channel → бачить channel3 ✅
   ```

4. **Результат:**
   - Всі 3 горутини читають з **channel3**
   - channel1 і channel2 ніхто не читає → **deadlock** або втрачені дані
   - Програма зависає або працює неправильно

## ✅ Правильні рішення

### Варіант 1: Використати `_` (ігнорувати індекс)

```go
for _, channel := range channels {
    go func() {
        defer wg.Done()
        for value := range channel {  // ✅ channel = параметр range, не змінна циклу
            outputCh <- value
        }
    }()
}
```

**Чому це працює:**
- `for _, channel := range channels` → `channel` - це **змінна значення** з range
- Кожна ітерація отримує **копію** наступного каналу
- Горутина захоплює правильне значення

**ВАЖЛИВО:** Це все одно НЕ 100% безпечно в Go < 1.22!

---

### Варіант 2: Передати як параметр функції (Найбезпечніше! ✅)

```go
for _, channel := range channels {
    go func(ch <-chan T) {  // ✅ ch - це параметр, кожна горутина має свою копію
        defer wg.Done()
        for value := range ch {
            outputCh <- value
        }
    }(channel)  // Передаємо значення явно
}
```

**Чому це найкраще:**
- ✅ Явна передача значення при виклику
- ✅ Кожна горутина отримує **свій параметр**
- ✅ Працює в усіх версіях Go
- ✅ Рекомендований спосіб

---

### Варіант 3: Створити локальну змінну

```go
for _, channel := range channels {
    channel := channel  // ✅ Створюємо нову змінну в scope ітерації
    go func() {
        defer wg.Done()
        for value := range channel {
            outputCh <- value
        }
    }()
}
```

**Чому працює:**
- `channel := channel` створює **нову змінну** в межах ітерації
- Кожна горутина захоплює окрему змінну

---

## 📊 Порівняння

| Варіант | Go < 1.22 | Go >= 1.22 | Безпека | Читабельність |
|---------|-----------|------------|---------|---------------|
| `for channel := range` | ❌ Баг | ✅ OK | ⚠️ Залежить від версії | ⭐⭐⭐ |
| `for _, channel := range` | ⚠️ Майже OK | ✅ OK | ⭐⭐⭐ | ⭐⭐⭐ |
| `go func(ch)` | ✅ OK | ✅ OK | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| `channel := channel` | ✅ OK | ✅ OK | ⭐⭐⭐⭐ | ⭐⭐ |

---

## 🔍 Візуалізація проблеми

### ❌ Неправильно (for channel := range):
```
Цикл:
[channel1] ← змінна channel
[channel2] ← змінна channel (перезаписує)
[channel3] ← змінна channel (перезаписує)
    ↓
Усі 3 горутини:
goroutine1 → читає channel → channel3 ❌
goroutine2 → читає channel → channel3 ❌
goroutine3 → читає channel → channel3 ❌
```

### ✅ Правильно (go func(ch) {} (channel)):
```
Цикл:
[channel1] → func(ch = channel1) → goroutine1 ✅
[channel2] → func(ch = channel2) → goroutine2 ✅
[channel3] → func(ch = channel3) → goroutine3 ✅
```

---

## 🎯 Висновок

**Чому не можна прибрати `_`:**
- Без `_` використовується індекс: `for i, channel := range`
- Або просто: `for channel := range` (що те саме, що `for i := range`)
- Це означає, що `channel` - це **одна змінна**, яка перезаписується
- Всі горутини бачать **останнє значення** цієї змінної

**Найкраще рішення:**
```go
for _, channel := range channels {
    go func(ch <-chan T) {
        defer wg.Done()
        for value := range ch {
            outputCh <- value
        }
    }(channel)  // ✅ Передаємо явно
}
```

**Це класична Go пастка!** Завжди передавайте змінні як параметри при створенні горутин у циклах! 🎯
