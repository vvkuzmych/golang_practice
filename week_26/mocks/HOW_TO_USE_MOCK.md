# 🎭 Як використовувати MockSortingStruct

## ❌ Твій код (не працює)

```go
type MockSortingStruct struct {
    Swap func(ctx context.Context, urls []string) []string
}

func (m *MockSortingStruct) SortMethod(ctx context.Context, urls []string) <-chan []string {
    result := make(chan []string)
    if m.Swap == nil {
        return result
    }
    return result
}
```

**Проблеми:**
1. ❌ Інтерфейс `URLSorter` має метод `Sort()`, а не `SortMethod()`
2. ❌ `Sort()` повинен повертати `([]string, error)`, а не `<-chan []string`
3. ❌ Метод повинен називатись `Sort`, щоб відповідати інтерфейсу

---

## ✅ Правильний код

### 1. Визначення Mock:

```go
type MockSortingStruct struct {
    // Функція яку ми inject
    SortFunc func(ctx context.Context, urls []string) ([]string, error)
    
    // Опціонально - для тестування
    CallCount int
    LastURLs  []string
    LastCtx   context.Context
}

// Sort реалізує URLSorter interface
func (m *MockSortingStruct) Sort(ctx context.Context, urls []string) ([]string, error) {
    // Tracking (опціонально)
    m.CallCount++
    m.LastURLs = urls
    m.LastCtx = ctx
    
    // Якщо не задано функцію - default behavior
    if m.SortFunc == nil {
        return urls, nil
    }
    
    // Викликаємо injected функцію
    return m.SortFunc(ctx, urls)
}
```

### 2. Використання в main.go:

```go
func main() {
    ctx := context.Background()
    
    // ════════════════════════════════════════════════════
    // Приклад 1: Mock з reverse logic
    // ════════════════════════════════════════════════════
    mockReverse := &MockSortingStruct{
        SortFunc: func(ctx context.Context, urls []string) ([]string, error) {
            // Кастомна логіка
            reversed := make([]string, len(urls))
            for i, url := range urls {
                reversed[len(urls)-1-i] = url
            }
            return reversed, nil
        },
    }
    
    svc := NewURLService(mockReverse)
    resultCh := svc.ProcessAsync(ctx, []string{"a.com", "b.com", "c.com"})
    result := <-resultCh
    
    fmt.Printf("Mock result: %v\n", result.URLs)
    // Output: [c.com b.com a.com]
    
    // ════════════════════════════════════════════════════
    // Приклад 2: Mock з помилкою
    // ════════════════════════════════════════════════════
    mockError := &MockSortingStruct{
        SortFunc: func(ctx context.Context, urls []string) ([]string, error) {
            return nil, fmt.Errorf("mock error")
        },
    }
    
    svc2 := NewURLService(mockError)
    resultCh2 := svc2.ProcessAsync(ctx, []string{"test.com"})
    result2 := <-resultCh2
    
    if result2.Err != nil {
        fmt.Printf("Error: %v\n", result2.Err)
    }
    
    // ════════════════════════════════════════════════════
    // Приклад 3: Mock без функції (default)
    // ════════════════════════════════════════════════════
    mockDefault := &MockSortingStruct{
        // Не задаємо SortFunc - використається default
    }
    
    svc3 := NewURLService(mockDefault)
    resultCh3 := svc3.ProcessAsync(ctx, []string{"x.com", "y.com"})
    result3 := <-resultCh3
    
    fmt.Printf("Default: %v\n", result3.URLs)
    // Output: [x.com y.com] (без змін)
    
    // ════════════════════════════════════════════════════
    // Приклад 4: Mock з tracking (для тестів)
    // ════════════════════════════════════════════════════
    mockTracking := &MockSortingStruct{
        SortFunc: func(ctx context.Context, urls []string) ([]string, error) {
            return urls, nil
        },
    }
    
    svc4 := NewURLService(mockTracking)
    
    // Викликаємо 3 рази
    for i := 1; i <= 3; i++ {
        resultCh := svc4.ProcessAsync(ctx, []string{"url.com"})
        <-resultCh
    }
    
    fmt.Printf("Викликано %d разів\n", mockTracking.CallCount)
    // Output: Викликано 3 разів
}
```

---

## 🎯 Чому це потрібно?

### URLSorter Interface:
```go
type URLSorter interface {
    Sort(ctx context.Context, urls []string) ([]string, error)
    //   ^                                      ^
    //   |                                      |
    //   Метод Sort                        Повертає []string і error
}
```

### Твій Mock ПОВИНЕН реалізувати цей interface:
```go
func (m *MockSortingStruct) Sort(ctx context.Context, urls []string) ([]string, error) {
    //                       ^                                      ^
    //                       |                                      |
    //                  Метод Sort                         Повертає []string і error
    
    return m.SortFunc(ctx, urls)
}
```

---

## 📊 Порівняння

| Аспект | DefaultURLSorter | MockSortingStruct |
|--------|------------------|-------------------|
| **Реальна логіка** | ✅ Справжнє сортування | ❌ Кастомна логіка |
| **Для production** | ✅ Так | ❌ Тільки для тестів |
| **Можна змінити логіку** | ❌ Ні | ✅ Так (через inject) |
| **Tracking викликів** | ❌ Ні | ✅ Так |
| **Можна симулювати помилки** | ❌ Ні | ✅ Так |

---

## 🎓 Приклади use cases

### 1. Тестування помилок:
```go
mockError := &MockSortingStruct{
    SortFunc: func(ctx context.Context, urls []string) ([]string, error) {
        return nil, fmt.Errorf("network timeout")
    },
}
// Тестуємо як додаток обробляє помилки
```

### 2. Тестування повільних операцій:
```go
mockSlow := &MockSortingStruct{
    SortFunc: func(ctx context.Context, urls []string) ([]string, error) {
        time.Sleep(5 * time.Second)
        return urls, nil
    },
}
// Тестуємо timeouts
```

### 3. Тестування cancellation:
```go
mockCancel := &MockSortingStruct{
    SortFunc: func(ctx context.Context, urls []string) ([]string, error) {
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        default:
            return urls, nil
        }
    },
}
// Тестуємо context cancellation
```

### 4. Tracking кількості викликів:
```go
mock := &MockSortingStruct{
    SortFunc: func(ctx context.Context, urls []string) ([]string, error) {
        return urls, nil
    },
}

// ... виконуємо операції ...

if mock.CallCount != 3 {
    t.Errorf("Expected 3 calls, got %d", mock.CallCount)
}
```

---

## 🚀 Запустити

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_26/mocks

# Оновлений main.go з mock
go run main.go mock_url_sorter.go

# Всі приклади з mock
go run mock_examples.go main.go mock_url_sorter.go
```

---

## 💡 Головне правило

**Mock ПОВИНЕН реалізувати той самий interface що і справжня реалізація!**

```
URLSorter interface
       ↓
   ┌───────────────────┐
   │                   │
   ↓                   ↓
DefaultURLSorter  MockSortingStruct
(production)         (testing)
```

Обидва реалізують `Sort(ctx, urls) ([]string, error)` ✅

---

## 📁 Файли

- `mock_examples.go` - 7 прикладів використання mock
- `main.go` - Оновлений з working mock
- `mock_url_sorter.go` - Правильна реалізація mock
- `HOW_TO_USE_MOCK.md` - Цей guide

---

## ✅ Checklist

Перевір що твій mock:
- [ ] Реалізує правильний interface (`URLSorter`)
- [ ] Має правильну сигнатуру методу `Sort()`
- [ ] Повертає правильний тип `([]string, error)`
- [ ] Можна inject кастомну логіку через `SortFunc`
- [ ] Має default behavior якщо `SortFunc == nil`
