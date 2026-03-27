# ⚡ Швидкий довідник: Без `time.After`

## 🎯 Питання: Як обійтись без `5*time.Second` в select?

**Відповідь: Є 7 способів! Ось найкращі 3:**

---

## 1️⃣ Найпростіший (БЕЗ select взагалі) ⭐

```go
results := newSorter.SortMethod(urls)

// Просто отримуємо!
url := <-results
fmt.Println(url)
```

**Для:** Швидкі операції, навчання

---

## 2️⃣ Production (Context з timeout) ⭐⭐⭐

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

results := newSorter.SortMethodWithContext(ctx, urls)

select {
case url := <-results:
    fmt.Println("Успіх:", url)
case <-ctx.Done():
    fmt.Println("Timeout:", ctx.Err())
}
```

**Для:** HTTP servers, production код, мережеві операції

---

## 3️⃣ UI з Cancel (Context з cancel) ⭐⭐

```go
ctx, cancel := context.WithCancel(context.Background())

// Користувач натиснув Cancel
go func() {
    <-cancelButton
    cancel()
}()

results := newSorter.SortMethodWithContext(ctx, urls)

select {
case url := <-results:
    fmt.Println("Успіх:", url)
case <-ctx.Done():
    fmt.Println("Скасовано")
}
```

**Для:** CLI/GUI додатки, довгі операції

---

## 📊 Швидке порівняння

| Варіант | Код | Timeout | Cancel | Складність |
|---------|-----|---------|--------|------------|
| No select | `url := <-ch` | ❌ | ❌ | ⭐ |
| Context timeout | `WithTimeout()` | ✅ | ✅ | ⭐⭐ |
| Context cancel | `WithCancel()` | ❌ | ✅ | ⭐⭐ |

---

## ✅ Що вибрати для твого коду?

Для `week_26/example/main.go`:

```go
func main() {
    urls := []string{"url1", "url2", "url3"}
    s := &DefaultSortingStruct{}
    newSorter := NewSorterService(s)

    results := newSorter.SortMethod(urls)

    // ✅ НАЙПРОСТІШЕ рішення!
    url := <-results
    fmt.Println("Отримані URLs:", url)
}
```

**Чому:** Операція швидка, гарантовано виконається, найчистіший код.

---

## 🚀 Спробувати

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_26/example

# Виправлений код
go run main_fixed.go

# Всі варіанти
go run alternatives.go
```

---

## 📚 Більше інформації

- `ALTERNATIVES.md` - Детальний опис всіх 7 варіантів
- `alternatives.go` - Робочий код всіх прикладів
- `PROBLEM_ANALYSIS.md` - Повний аналіз проблеми

---

## 💡 Правило

**"Використовуй найпростіше рішення що працює!"**

```
Простий код → Легко читати → Менше багів → Щасливий розробник! 😊
```
