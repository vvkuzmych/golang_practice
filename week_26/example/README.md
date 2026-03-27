# 🐛 Проблема в week_26/example/main.go

## ❌ Проблема

**Поточний код видає:**
```
not yet
```

**Замість очікуваних URLs!**

---

## 🔍 Що не так?

### Проблемний код (рядки 54-59):

```go
results := newSorter.SortMethod(urls)

select {
case url := <-results:
    fmt.Println(url)
default:                    // ← ПРОБЛЕМА!
    fmt.Println("not yet")
}
```

### Пояснення:

**`select` з `default` - це NON-BLOCKING операція!**

```
Timeline виконання:

1. Запускається goroutine
2. Main thread одразу йде до select
3. Channel ще порожній (goroutine не встигла надіслати)
4. Є default → виконується одразу
5. Друкує "not yet"
6. Програма завершується
7. Goroutine не встигла надіслати дані - вони втрачені! 😢
```

---

## ✅ Рішення

### Варіант 1: Видалити `default` (НАЙКРАЩЕ)

```go
results := newSorter.SortMethod(urls)

// Просто чекаємо результат
url := <-results
fmt.Println("Отримані URLs:", url)
```

**Результат:**
```
Отримані URLs: [https://www.baidu1.com.cn/ https://www.baidu2.com.cn/ https://www.baidu3.com.cn/]
```

✅ **Працює!**

---

### Варіант 2: Додати timeout

```go
results := newSorter.SortMethod(urls)

select {
case url := <-results:
    fmt.Println("Отримані URLs:", url)
case <-time.After(5 * time.Second):
    fmt.Println("Timeout! Занадто довго")
}
```

---

## 🚀 Запустити

### Проблемний код:
```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_26/example
go run main.go
# Output: not yet ❌
```

### Виправлений код:
```bash
go run main_fixed.go
# Output: Отримані URLs: [...] ✅
```

---

## 💡 Головний урок

| Select Type | Поведінка | Використання |
|-------------|-----------|--------------|
| `select { case ...; default ... }` | **NON-BLOCKING** - одразу виконує default якщо channel порожній | Polling, non-blocking checks |
| `select { case ...; case ... }` | **BLOCKING** - чекає поки хоч один case спрацює | Чекання на результат |
| `<-channel` | **BLOCKING** - чекає поки дані прийдуть | Найпростіший спосіб отримати дані |

---

## 📁 Файли

- `main.go` - Оригінальний (з проблемою)
- `main_fixed.go` - Виправлений (працює!)
- `PROBLEM_ANALYSIS.md` - Детальний аналіз з усіма варіантами рішень
- `README.md` - Це файл

---

## 🎓 Додатково

Детальний аналіз з 4 варіантами рішень дивись в `PROBLEM_ANALYSIS.md`
