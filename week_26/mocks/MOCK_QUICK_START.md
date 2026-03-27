# 🚀 MockSortingStruct - Quick Start

## ⚡ TL;DR

```go
// 1. Створюємо mock з кастомною логікою
mock := &MockURLSorter{
    SortFunc: func(ctx context.Context, urls []string) ([]string, error) {
        // Твоя кастомна логіка тут
        return urls, nil
    },
}

// 2. Використовуємо як звичайний URLSorter
svc := NewURLService(mock)
resultCh := svc.ProcessAsync(ctx, []string{"a.com", "b.com"})
result := <-resultCh
```

---

## ❌ Твоя помилка

```go
type MockSortingStruct struct {
    Swap func(ctx context.Context, urls []string) []string  // ❌
}

func (m *MockSortingStruct) SortMethod(...) <-chan []string { // ❌
    // ❌ Не відповідає interface!
}
```

**Проблеми:**
- ❌ Метод має бути `Sort`, а не `SortMethod`
- ❌ Має повертати `([]string, error)`, а не channel
- ❌ Не відповідає `URLSorter` interface

---

## ✅ Правильна реалізація

```go
type MockSortingStruct struct {
    SortFunc func(ctx context.Context, urls []string) ([]string, error)
}

// Sort реалізує URLSorter interface
func (m *MockSortingStruct) Sort(ctx context.Context, urls []string) ([]string, error) {
    if m.SortFunc == nil {
        return urls, nil  // default
    }
    return m.SortFunc(ctx, urls)
}
```

---

## 🎯 Приклади використання

### 1. Reverse logic:
```go
mock := &MockURLSorter{
    SortFunc: func(ctx context.Context, urls []string) ([]string, error) {
        reversed := make([]string, len(urls))
        for i, url := range urls {
            reversed[len(urls)-1-i] = url
        }
        return reversed, nil
    },
}

svc := NewURLService(mock)
resultCh := svc.ProcessAsync(ctx, []string{"a.com", "b.com", "c.com"})
result := <-resultCh
fmt.Println(result.URLs)  // [c.com b.com a.com]
```

### 2. Return error:
```go
mock := &MockURLSorter{
    SortFunc: func(ctx context.Context, urls []string) ([]string, error) {
        return nil, fmt.Errorf("mock error")
    },
}

svc := NewURLService(mock)
resultCh := svc.ProcessAsync(ctx, []string{"test.com"})
result := <-resultCh
fmt.Println(result.Err)  // mock error
```

### 3. Default behavior (no function):
```go
mock := &MockURLSorter{
    // Не задаємо SortFunc
}

svc := NewURLService(mock)
resultCh := svc.ProcessAsync(ctx, []string{"x.com", "y.com"})
result := <-resultCh
fmt.Println(result.URLs)  // [x.com y.com] (без змін)
```

---

## 📊 Interface відповідність

```go
// Interface каже:
type URLSorter interface {
    Sort(ctx context.Context, urls []string) ([]string, error)
}

// Твій mock ПОВИНЕН мати:
func (m *MockSortingStruct) Sort(ctx context.Context, urls []string) ([]string, error) {
    //                       ^^^^                                  ^^^^^^^^^^^^^^^^^^^^
    //                    Назва методу                            Повертає це
    return m.SortFunc(ctx, urls)
}
```

**Все має співпадати! ✅**

---

## 🔄 Як це працює

```
1. Створюємо mock:
   mock := &MockURLSorter{
       SortFunc: func(...) { твоя логіка }
   }

2. Передаємо в service:
   svc := NewURLService(mock)
                         ^^^^
                    Приймає URLSorter

3. Service викликає Sort():
   res, err := s.sorter.Sort(ctx, urls)
                        ^^^^
                   Викликає твій mock

4. Mock викликає SortFunc:
   return m.SortFunc(ctx, urls)
          ^^^^^^^^^^
       Твоя кастомна логіка!
```

---

## 🚀 Запустити

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_26/mocks

# Main з mock
go run main.go mock_url_sorter.go

# Всі приклади
go run mock_examples.go main.go mock_url_sorter.go
```

**Output:**
```
1. Default Sorter:
Sorted URLs: [amazon.com github.com google.com stackoverflow.com]

2. Mock Sorter:
Mock result: [c.com b.com a.com]  ← Reverse order!
```

---

## 📁 Файли

- `HOW_TO_USE_MOCK.md` - Детальний guide
- `mock_examples.go` - 7 прикладів
- `MOCK_QUICK_START.md` - Це швидкий старт
- `main.go` - Working example

---

## 💡 Запам'ятай

1. **Mock повинен реалізувати interface** ✅
2. **Метод має бути `Sort`, не `SortMethod`** ✅
3. **Повертає `([]string, error)`, не channel** ✅
4. **Можна inject будь-яку логіку через `SortFunc`** ✅
