# 🔄 URL Sorting Service - Technical Flow

## Execution Flow when `main.go` is Launched

---

## 📋 Main Execution Flow

### **Step 1: Program Start**
```go
func main() {
    fmt.Println("=== URL Sorting Service Demo ===")
```

**What happens:**
- `main()` function is invoked
- Prints header to console

---

### **Step 2: Create Default Sorter**
```go
defaultSorter := &DefaultURLSorter{}
svc := NewURLService(defaultSorter)
```

**What happens:**
- Instantiate `DefaultURLSorter` - concrete implementation that sorts URLs by hostname
- Create `URLService` with the sorter injected

**🎯 Pattern: Dependency Injection**
- `URLService` accepts `URLSorter` interface
- Allows different implementations to be injected
- Enables testing with mocks

---

### **Step 3: Example 1 - Default Sorter (Async Channel)**
```go
ctx := context.Background()
resultCh := svc.ProcessAsync(ctx, []string{
    "github.com",
    "stackoverflow.com",
    "google.com",
    "amazon.com",
})

result := <-resultCh  // Block until result received
if result.Err != nil {
    fmt.Printf("Error: %v\n", result.Err)
} else {
    fmt.Printf("Sorted URLs: %v\n", result.URLs)
}
```

**What happens:**
1. Call `ProcessAsync()` which returns a channel
2. Sorting happens in a **goroutine** (asynchronous)
3. Main thread **blocks** on channel receive (`<-resultCh`)
4. When sorting completes, result is sent to channel
5. Main thread receives result and prints it

**⚡ Pattern: Channel-based Async**
- Goroutine runs `Sort()` in background
- Result is sent via channel
- Main thread can do other work or wait

**Expected Output:**
```
Sorted URLs: [amazon.com github.com google.com stackoverflow.com]
```

---

### **Step 4: Example 2 - Mock Sorter**
```go
mockSorter := &MockURLSorter{
    SortFunc: func(ctx context.Context, urls []string) ([]string, error) {
        // Custom mock behavior - reverse order
        reversed := make([]string, len(urls))
        for i, url := range urls {
            reversed[len(urls)-1-i] = url
        }
        return reversed, nil
    },
}

svcMock := NewURLService(mockSorter)
resultCh2 := svcMock.ProcessAsync(ctx, []string{"a.com", "b.com", "c.com"})
result2 := <-resultCh2
fmt.Printf("Mock result: %v\n", result2.URLs)
```

**What happens:**
1. Create `MockURLSorter` with **custom behavior**
2. Mock reverses the array instead of sorting
3. Same async pattern as Example 1
4. Demonstrates how mocks can inject different behavior

**🧪 Pattern: Test Double / Mock**
- `MockURLSorter` implements same `URLSorter` interface
- Can inject any custom behavior via `SortFunc`
- Used in tests to avoid real sorting logic

**Expected Output:**
```
Mock result: [c.com b.com a.com]
```

---

### **Step 5: Example 3 - Context Cancellation**
```go
ctxCancel, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
defer cancel()

time.Sleep(10 * time.Millisecond)  // Ensure timeout occurs
resultCh3 := svc.ProcessAsync(ctxCancel, []string{"x.com", "y.com"})
result3 := <-resultCh3
if result3.Err != nil {
    fmt.Printf("Expected error: %v\n", result3.Err)
}
```

**What happens:**
1. Create context with **1ms timeout**
2. Sleep for **10ms** to ensure timeout occurs
3. Call `ProcessAsync()` with expired context
4. `Sort()` checks context and returns error
5. Error is `context.DeadlineExceeded`

**⏱️ Pattern: Context Cancellation**
- Demonstrates graceful handling of timeouts
- `Sort()` checks `ctx.Done()` before processing
- Returns error if context is cancelled/timeout

**Expected Output:**
```
Expected error: context deadline exceeded
```

---

### **Step 6: Example 4 - Callback Pattern**
```go
done := make(chan bool)
svc.ProcessWithCallback(ctx, []string{"example.com", "test.com"}, 
    func(urls []string, err error) {
        if err != nil {
            fmt.Printf("Callback error: %v\n", err)
        } else {
            fmt.Printf("Callback result: %v\n", urls)
        }
        done <- true
    })
<-done  // Wait for callback to complete
```

**What happens:**
1. Create a `done` channel to wait for completion
2. Call `ProcessWithCallback()` with a **callback function**
3. Sorting happens in goroutine
4. When complete, callback is invoked with results
5. Callback sends signal to `done` channel
6. Main thread waits on `done` channel

**📞 Pattern: Callback-based Async**
- Alternative async pattern to channels
- Caller provides function to be called on completion
- Common in JavaScript, less common in Go

**Expected Output:**
```
Callback result: [example.com test.com]
```

---

### **Step 7: Program End**
```go
fmt.Println("=== Demo Complete ===")
```

**What happens:**
- Print completion message
- Program exits

---

## 🔍 Deep Dive: ProcessAsync() Goroutine Flow

### Detailed Execution Steps:

```go
func (s *URLService) ProcessAsync(ctx context.Context, urls []string) <-chan Result {
    // Step A: Create buffered channel
    resultCh := make(chan Result, 1)  // Buffer size 1

    // Step B: Spawn goroutine
    go func() {
        defer close(resultCh)  // Step F: Ensure channel is closed

        // Step C: Execute Sort() in goroutine
        res, err := s.sorter.Sort(ctx, urls)
        
        // Step D: Send result to channel
        resultCh <- Result{
            URLs: res,
            Err:  err,
        }
    }()

    // Step E: Return channel immediately (non-blocking)
    return resultCh
}
```

### Timeline:

```
Main Thread                          Goroutine Thread
───────────                          ────────────────

1. Call ProcessAsync()
2. Create channel
3. Spawn goroutine               →   4. Start executing
4. Return channel immediately
5. Continue execution
6. Block on <-resultCh                5. Call sorter.Sort()
                                      6. Wait for sorting...
                                      7. Send result to channel  →
7. Receive result                ←    8. Close channel (defer)
8. Continue with result               9. Goroutine exits
```

---

## 🎨 Design Patterns Used

### 1. **Dependency Injection**
```go
type URLService struct {
    sorter URLSorter  // Interface, not concrete type
}
```
- Service depends on interface, not implementation
- Enables testing with mocks
- Loose coupling

### 2. **Interface Segregation**
```go
type URLSorter interface {
    Sort(ctx context.Context, urls []string) ([]string, error)
}
```
- Small, focused interface
- Easy to implement
- Easy to mock

### 3. **Test Double / Mock**
```go
type MockURLSorter struct {
    SortFunc func(ctx context.Context, urls []string) ([]string, error)
}
```
- Implements same interface as real sorter
- Allows custom behavior injection
- Used for testing without real implementation

### 4. **Channel-based Concurrency**
```go
resultCh := make(chan Result, 1)
go func() {
    // async work
    resultCh <- result
}()
return resultCh
```
- Goroutine for async execution
- Channel for communication
- Non-blocking return

### 5. **Context Propagation**
```go
func (s *DefaultURLSorter) Sort(ctx context.Context, urls []string) ([]string, error) {
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }
    // ...
}
```
- Propagate cancellation/timeout through call stack
- Check context before expensive operations
- Graceful error handling

---

## 📊 Data Flow Diagram

```
┌─────────────┐
│   main()    │
└──────┬──────┘
       │
       │ 1. Create DefaultURLSorter
       ├─────────────────────────────────┐
       │                                 │
       │ 2. Inject into URLService       │
       ├────────────────────────────┐    │
       │                            ▼    ▼
       │                    ┌────────────────────┐
       │                    │   URLService       │
       │                    │  ┌──────────────┐  │
       │                    │  │ sorter field │  │
       │                    │  └──────────────┘  │
       │                    └────────┬───────────┘
       │                             │
       │ 3. Call ProcessAsync()      │
       ├─────────────────────────────┤
       │                             │
       │                             │ Create goroutine
       │                             ├──────────────────┐
       │                             │                  │
       │                             │              ┌───▼─────┐
       │ 4. Return channel           │              │ Sort()  │
       │◄────────────────────────────┤              │         │
       │                             │              └───┬─────┘
       │                             │                  │
       │ 5. Block on <-channel       │                  │
       │────────────────────X        │                  │
       │                    ║        │              ┌───▼─────────┐
       │                    ║        │              │ Send Result │
       │                    ║        │              └───┬─────────┘
       │                    ║        │                  │
       │ 6. Receive result  ║        │◄─────────────────┘
       │◄───────────────────X        │
       │                             │
       ▼                             ▼
   Continue                      Close channel
```

---

## ❓ ЧОМУ ВИКОРИСТОВУЄТЬСЯ goroutine В ProcessWithCallback?

### Запитання:
```go
func (s *URLService) ProcessWithCallback(
    ctx context.Context,
    urls []string,
    callback func([]string, error),
) {
    go func() {  // ← ЧОМУ ТУТ goroutine?
        res, err := s.sorter.Sort(ctx, urls)
        callback(res, err)
    }()
}
```

### Відповідь: Щоб зробити функцію АСИНХРОННОЮ (non-blocking)

---

### 🔴 Що буде БЕЗ goroutine:

```go
func ProcessWithCallbackSync(ctx context.Context, urls []string, callback func([]string, error)) {
    // ❌ БЕЗ goroutine
    res, err := s.sorter.Sort(ctx, urls)  // ← БЛОКУЄ основний потік!
    callback(res, err)
    // Функція повертається ПІСЛЯ всієї роботи
}
```

**Timeline:**
```
Main Thread:
  │
  ├─ Виклик ProcessWithCallbackSync()
  │
  ├─ Виконання sorter.Sort() ⏳⏳⏳  ← БЛОКУЄТЬСЯ ТУТ!
  │  (чекаємо 2 секунди...)
  │
  ├─ Виклик callback()
  │
  └─ Повернення з функції  ← Тільки через 2 секунди!

❌ ПРОБЛЕМА: Не можемо робити іншу роботу під час сортування!
```

---

### 🟢 Що буде З goroutine:

```go
func ProcessWithCallbackAsync(ctx context.Context, urls []string, callback func([]string, error)) {
    // ✅ З goroutine
    go func() {
        res, err := s.sorter.Sort(ctx, urls)
        callback(res, err)
    }()
    // Функція повертається МИТТЄВО!
}
```

**Timeline:**
```
Main Thread                      Goroutine Thread
────────────                     ────────────────
  │
  ├─ Виклик ProcessWithCallbackAsync()
  │
  ├─ Створення goroutine      →  │
  │                                ├─ Виконання sorter.Sort() ⏳⏳⏳
  └─ Повернення МИТТЄВО! ✓        │  (сортування в фоні)
  │                                │
  ├─ Можу робити іншу роботу ✓    │
  ├─ Ще одна задача ✓              │
  ├─ І ще одна ✓                   │
  │                                ├─ Виклик callback()
  │                                └─ Goroutine завершена
  └─ Продовжую роботу...

✅ ПЕРЕВАГА: Паралельне виконання! Основний потік не блокується!
```

---

### 📊 Порівняння:

| Аспект | БЕЗ goroutine | З goroutine |
|--------|---------------|-------------|
| **Повернення функції** | Через 2 сек | Миттєво (~0мс) |
| **Основний потік** | Блокується ❌ | Вільний ✅ |
| **Паралельність** | Ні | Так |
| **Інша робота** | Неможлива | Можлива |
| **Тип** | Синхронний | Асинхронний |

---

### 🎯 5 Причин використовувати goroutine:

#### 1. ⚡ **Асинхронність (Non-blocking)**
Функція повертається миттєво, не чекаючи завершення роботи.
```go
start := time.Now()
ProcessWithCallbackAsync(ctx, urls, callback)
fmt.Printf("Повернулось за: %v\n", time.Since(start))  // ~0ms
```

#### 2. 🔄 **Паралельність**
Основний потік може продовжувати роботу, поки callback виконується в фоні.
```go
ProcessWithCallbackAsync(ctx, urls, callback)
// Можемо одразу робити щось інше!
doOtherWork()
handleNextRequest()
```

#### 3. 📱 **Відповідність патерну "Callback"**
Callback pattern означає "**викличи мене пізніше**", а не "**зачекай тут**".
- JavaScript: `setTimeout(callback, 1000)` - асинхронний
- Go callback також має бути асинхронним

#### 4. 🎯 **Узгодженість з ProcessAsync()**
Обидва методи працюють однаково - асинхронно:
```go
// ProcessAsync - повертає channel, робота в goroutine
resultCh := service.ProcessAsync(ctx, urls)  // миттєво

// ProcessWithCallback - викликає callback, робота в goroutine
service.ProcessWithCallback(ctx, urls, callback)  // миттєво
```

#### 5. 💪 **Кращий User Experience**
- **UI**: Інтерфейс не зависає під час обробки
- **Server**: Може обробляти багато запитів паралельно
- **CLI**: Можна показувати прогрес або виконувати інші команди

---

### 🌐 Реальний приклад: Web Server

#### ❌ БЕЗ goroutine:
```go
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    service.ProcessWithCallbackSync(ctx, urls, func(result []string, err error) {
        json.NewEncoder(w).Encode(result)
    })
    // ← Тут БЛОКУЄТЬСЯ вся goroutine веб-сервера!
    // ← Інші запити чекають своєї черги!
}
```
**Результат:** Сервер може обробити **1 запит за раз** 😢

#### ✅ З goroutine:
```go
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    service.ProcessWithCallbackAsync(ctx, urls, func(result []string, err error) {
        json.NewEncoder(w).Encode(result)
    })
    // ← Одразу повернулись!
    // ← Можемо обробляти наступний запит!
}
```
**Результат:** Сервер може обробити **ТИСЯЧІ запитів паралельно** 🚀

---

### 💡 Висновок:

**Goroutine в `ProcessWithCallback` робить функцію асинхронною:**
- Повертається миттєво
- Не блокує основний потік
- Дозволяє паралельне виконання
- Відповідає семантиці callback pattern
- Покращує продуктивність і UX

**Без goroutine** - це була б **синхронна** функція, що блокує виконання!

---

## 🎯 Key Takeaways

1. **Dependency Injection** enables testing and flexibility
2. **Goroutines** provide true concurrency in Go
3. **Channels** are the primary communication mechanism
4. **Context** enables cancellation and timeouts
5. **Interfaces** allow multiple implementations
6. **Mocks** enable testing without real dependencies
7. **Goroutines in callbacks** make functions asynchronous and non-blocking

---

## 🚀 Running the Program

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_26/mocks
go run main.go mock_url_sorter.go
```

**Expected Output:**
```
=== URL Sorting Service Demo ===

1. Default Sorter:
Sorted URLs: [amazon.com github.com google.com stackoverflow.com]

2. Mock Sorter:
Mock result: [c.com b.com a.com]

3. Context Cancellation:
Expected error: context deadline exceeded

4. Callback Pattern:
Callback result: [example.com test.com]

=== Demo Complete ===
```
