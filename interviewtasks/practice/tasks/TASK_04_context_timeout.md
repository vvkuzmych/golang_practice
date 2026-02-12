# Task 4: Context with Timeout

**Level:** Advanced  
**Time:** 15 minutes  
**Topics:** Context, Timeout, Cancellation

---

## 📝 Task

Напиши функцію, яка виконує HTTP запит з можливістю **cancellation** через context.

Якщо context cancelled (timeout або manual cancel), запит має зупинитись негайно.

---

## 📥 Function Signature

```go
func FetchWithContext(ctx context.Context, url string) (string, error)
```

**Parameters:**
- `ctx` - context для cancellation
- `url` - URL для запиту

**Returns:**
- `string` - response body
- `error` - помилка (включаючи context.DeadlineExceeded)

---

## 💡 Examples

```go
// Example 1: Normal request (no timeout)
ctx := context.Background()
body, err := FetchWithContext(ctx, "https://google.com")
// => HTML content, nil

// Example 2: Request with timeout
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
defer cancel()

body, err := FetchWithContext(ctx, "https://httpbin.org/delay/5")  // 5 sec delay
// => "", context.DeadlineExceeded

// Example 3: Manual cancellation
ctx, cancel := context.WithCancel(context.Background())

go func() {
    time.Sleep(500 * time.Millisecond)
    cancel()  // Cancel після 500ms
}()

body, err := FetchWithContext(ctx, "https://google.com")
// => "", context.Canceled
```

---

## ✅ Requirements

- Використай `http.NewRequestWithContext()` для HTTP запиту з context
- Перевір `ctx.Done()` перед запитом
- Повертай `ctx.Err()` якщо context cancelled
- Підтримуй різні типи cancellation:
  - `context.WithTimeout` - timeout
  - `context.WithDeadline` - absolute deadline
  - `context.WithCancel` - manual cancellation

---

## 🧪 Test Cases

```go
// Test 1: Successful request
ctx := context.Background()
body, err := FetchWithContext(ctx, "https://google.com")
assert.Nil(t, err)
assert.NotEmpty(t, body)

// Test 2: Timeout
ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
defer cancel()
_, err := FetchWithContext(ctx, "https://httpbin.org/delay/1")
assert.Equal(t, context.DeadlineExceeded, err)

// Test 3: Already cancelled context
ctx, cancel := context.WithCancel(context.Background())
cancel()  // Cancel одразу
_, err := FetchWithContext(ctx, "https://google.com")
assert.Equal(t, context.Canceled, err)

// Test 4: Context cancelled during request
ctx, cancel := context.WithCancel(context.Background())
go func() {
    time.Sleep(50 * time.Millisecond)
    cancel()
}()
_, err := FetchWithContext(ctx, "https://httpbin.org/delay/10")
assert.NotNil(t, err)

// Test 5: Deadline
deadline := time.Now().Add(500 * time.Millisecond)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()
_, err := FetchWithContext(ctx, "https://httpbin.org/delay/2")
assert.Equal(t, context.DeadlineExceeded, err)
```

---

## 💡 Hints

1. Використай `http.NewRequestWithContext()` замість `http.Get()`:
```go
req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
if err != nil {
    return "", err
}

resp, err := http.DefaultClient.Do(req)
```

2. Перевір context перед запитом:
```go
select {
case <-ctx.Done():
    return "", ctx.Err()
default:
}
```

3. Context автоматично cancel HTTP request якщо timeout

---

## 🎯 Real-World Use Case

```go
// Microservice з timeout для external API
func GetUserData(userID int) (*User, error) {
    // 3 second timeout для external API
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    
    url := fmt.Sprintf("https://api.external.com/users/%d", userID)
    body, err := FetchWithContext(ctx, url)
    if err != nil {
        if err == context.DeadlineExceeded {
            return nil, fmt.Errorf("external API timeout")
        }
        return nil, err
    }
    
    // Parse response...
    return user, nil
}
```

---

**Рішення:** `solutions/solution_04_context_timeout.go`

**Good luck!** 🚀
