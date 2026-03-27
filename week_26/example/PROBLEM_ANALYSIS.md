# 🐛 Проблема в week_26/example/main.go

## ❌ Що не так?

Коли запускаєш програму, отримуєш:
```
not yet
```

Замість очікуваних URLs!

---

## 🔍 Аналіз проблеми

### Поточний код (рядки 52-59):
```go
results := newSorter.SortMethod(urls)

select {
case url := <-results:
    fmt.Println(url)
default:                    // ← ПРОБЛЕМА ТУТ!
    fmt.Println("not yet")
}
```

### Що відбувається:

```
Timeline:

Main Goroutine                    Worker Goroutine
──────────────                    ────────────────
1. Виклик SortMethod()
2. Створено channel
3. Запущено goroutine        →   4. Почав виконання
4. Повернуто channel               5. Виконує Swap()...
5. Одразу select{}                 
6. Channel ще ПУСТИЙ! ⚠️          
7. Є default → виконує його        6. Надсилає в channel
8. Друкує "not yet"                
9. ВИХІД з main()              →   7. Goroutine все ще працює...
10. Програма завершилась           8. Данні втрачені! 😢
```

**Проблема:** `select` з `default` - це **NON-BLOCKING** операція!
- Якщо channel порожній, одразу виконується `default`
- Goroutine не встигає надіслати дані
- Програма завершується до отримання результату

---

## ✅ 3 способи виправити

### Рішення 1: Видалити `default` (НАЙПРОСТІШЕ)

```go
func main() {
    urls := []string{
        "https://www.baidu1.com.cn/",
        "https://www.baidu2.com.cn/",
        "https://www.baidu3.com.cn/",
    }

    s := &DefaultSortingStruct{}
    newSorter := NewSorterService(s)

    results := newSorter.SortMethod(urls)

    // ✅ Блокуючий receive - чекає поки дані прийдуть
    url := <-results
    fmt.Println(url)
}
```

**Переваги:**
- ✅ Найпростіше рішення
- ✅ Гарантовано отримаємо результат
- ✅ Працює правильно

**Timeline:**
```
Main Goroutine                    Worker Goroutine
──────────────                    ────────────────
1. Виклик SortMethod()
2. results := <-results       
3. БЛОКУЄТЬСЯ тут...         →    4. Виконує Swap()
                                   5. Надсилає в channel  →
6. Отримує дані! ✅           ←    6. Закриває channel
7. Друкує результат
8. Завершення
```

---

### Рішення 2: Використати `time.Sleep` (ПОГАНЕ рішення)

```go
func main() {
    urls := []string{
        "https://www.baidu1.com.cn/",
        "https://www.baidu2.com.cn/",
        "https://www.baidu3.com.cn/",
    }

    s := &DefaultSortingStruct{}
    newSorter := NewSorterService(s)

    results := newSorter.SortMethod(urls)

    // ⚠️ Даємо час goroutine
    time.Sleep(100 * time.Millisecond)

    select {
    case url := <-results:
        fmt.Println(url)
    default:
        fmt.Println("not yet")
    }
}
```

**Недоліки:**
- ❌ Hardcoded затримка
- ❌ Може не вистачити часу для повільних операцій
- ❌ Марнує час якщо швидко виконується
- ❌ Не професійний підхід

---

### Рішення 3: `select` з `time.After` (якщо потрібен timeout)

```go
func main() {
    urls := []string{
        "https://www.baidu1.com.cn/",
        "https://www.baidu2.com.cn/",
        "https://www.baidu3.com.cn/",
    }

    s := &DefaultSortingStruct{}
    newSorter := NewSorterService(s)

    results := newSorter.SortMethod(urls)

    // ✅ Чекаємо результат АБО timeout
    select {
    case url := <-results:
        fmt.Println("Отримали:", url)
    case <-time.After(5 * time.Second):
        fmt.Println("Timeout! Занадто довго чекали")
    }
}
```

**Переваги:**
- ✅ Чекає на результат
- ✅ Має захист від зависання
- ✅ Професійний підхід

**Коли використовувати:**
- Коли операція може зависнути
- Коли потрібен timeout
- Для production коду

---

### Рішення 4: Цикл з `select` (для обробки багатьох значень)

Якщо channel може відправити багато значень:

```go
func main() {
    urls := []string{
        "https://www.baidu1.com.cn/",
        "https://www.baidu2.com.cn/",
        "https://www.baidu3.com.cn/",
    }

    s := &DefaultSortingStruct{}
    newSorter := NewSorterService(s)

    results := newSorter.SortMethod(urls)

    // ✅ Обробляємо всі результати
    for {
        select {
        case url, ok := <-results:
            if !ok {
                // Channel закрито
                fmt.Println("Всі результати отримано")
                return
            }
            fmt.Println("Отримали:", url)
        case <-time.After(5 * time.Second):
            fmt.Println("Timeout!")
            return
        }
    }
}
```

---

## 📊 Порівняння рішень

| Рішення | Складність | Надійність | Коли використовувати |
|---------|-----------|------------|---------------------|
| Блокуючий receive | ⭐ Проста | ⭐⭐⭐ Висока | Одне значення, швидка операція |
| `time.Sleep` | ⭐ Проста | ❌ Низька | НІКОЛИ (поганий код) |
| `time.After` | ⭐⭐ Середня | ⭐⭐⭐ Висока | Потрібен timeout |
| Цикл з select | ⭐⭐⭐ Складна | ⭐⭐⭐ Висока | Багато значень, потрібен timeout |

---

## 🎯 Рекомендоване рішення

**Для цього коду - Рішення 1 (видалити `default`):**

```go
func main() {
    urls := []string{
        "https://www.baidu1.com.cn/",
        "https://www.baidu2.com.cn/",
        "https://www.baidu3.com.cn/",
    }

    s := &DefaultSortingStruct{}
    newSorter := NewSorterService(s)

    results := newSorter.SortMethod(urls)

    // Просто чекаємо результат
    url := <-results
    fmt.Println(url)
}
```

**Очікуваний output:**
```
[https://www.baidu1.com.cn/ https://www.baidu2.com.cn/ https://www.baidu3.com.cn/]
```

---

## 💡 Додаткові покращення

### 1. Додати контекст для cancellation:

```go
func (newSorter *SorterService) SortMethod(ctx context.Context, urls []string) <-chan []string {
    results := make(chan []string)

    go func() {
        defer close(results)

        // Перевіряємо чи контекст не скасовано
        select {
        case <-ctx.Done():
            return
        default:
        }

        res := newSorter.sorter.Swap(urls)
        
        // Надсилаємо з можливістю cancellation
        select {
        case results <- res:
        case <-ctx.Done():
        }
    }()
    return results
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    urls := []string{"url1", "url2", "url3"}
    s := &DefaultSortingStruct{}
    newSorter := NewSorterService(s)

    results := newSorter.SortMethod(ctx, urls)

    select {
    case url := <-results:
        fmt.Println(url)
    case <-ctx.Done():
        fmt.Println("Timeout або cancellation")
    }
}
```

### 2. Додати error handling:

```go
type Result struct {
    URLs []string
    Err  error
}

func (newSorter *SorterService) SortMethod(urls []string) <-chan Result {
    results := make(chan Result, 1)

    go func() {
        defer close(results)

        res := newSorter.sorter.Swap(urls)
        results <- Result{URLs: res, Err: nil}
    }()
    return results
}

func main() {
    urls := []string{"url1", "url2", "url3"}
    s := &DefaultSortingStruct{}
    newSorter := NewSorterService(s)

    resultCh := newSorter.SortMethod(urls)
    result := <-resultCh

    if result.Err != nil {
        fmt.Printf("Error: %v\n", result.Err)
        return
    }

    fmt.Println("Success:", result.URLs)
}
```

---

## 🎓 Головний урок

**`select` with `default` = NON-BLOCKING**
```go
select {
case data := <-ch:
    // Отримав дані
default:
    // Одразу виконується якщо channel порожній!
}
```

**`select` without `default` = BLOCKING**
```go
select {
case data := <-ch:
    // Чекає поки дані прийдуть
case <-time.After(5 * time.Second):
    // АБО timeout
}
```

---

## ✅ Виправлений код

Готовий до використання код в `main_fixed.go`!
