# Task 2: Concurrent URL Checker

**Level:** Intermediate  
**Time:** 15 minutes  
**Topics:** Goroutines, Channels, Error Handling

---

## 📝 Task

Напиши функцію, яка перевіряє доступність (availability) списку URLs **паралельно** і повертає результати.

---

## 📥 Function Signature

```go
type URLStatus struct {
    URL        string
    StatusCode int
    Error      error
}

func CheckURLs(urls []string) []URLStatus
```

**Parameters:**
- `urls` - slice URLs для перевірки

**Returns:**
- `[]URLStatus` - результати для кожного URL в **тому самому порядку**

---

## 💡 Examples

```go
urls := []string{
    "https://google.com",
    "https://github.com",
    "https://invalid-url-that-does-not-exist.com",
}

results := CheckURLs(urls)

// results[0] = URLStatus{URL: "https://google.com", StatusCode: 200, Error: nil}
// results[1] = URLStatus{URL: "https://github.com", StatusCode: 200, Error: nil}
// results[2] = URLStatus{URL: "https://invalid...", StatusCode: 0, Error: <error>}
```

---

## ✅ Requirements

- Використай goroutines для паралельних запитів
- Використай channel для збору результатів
- Збережи порядок результатів (result для `urls[0]` має бути в `results[0]`)
- Якщо URL недоступний, StatusCode = 0 та Error != nil
- Використай `http.Get()` для запитів
- Встанови timeout 5 секунд для кожного запиту
- Всі URLs мають перевірятись, навіть якщо деякі фейляться

---

## 🧪 Test Cases

```go
// Test 1: All URLs valid
urls := []string{"https://google.com", "https://github.com"}
results := CheckURLs(urls)
assert.Equal(t, 2, len(results))
assert.Nil(t, results[0].Error)
assert.Nil(t, results[1].Error)

// Test 2: Mix of valid and invalid
urls := []string{
    "https://google.com",
    "https://this-url-definitely-does-not-exist-12345.com",
}
results := CheckURLs(urls)
assert.Nil(t, results[0].Error)
assert.NotNil(t, results[1].Error)

// Test 3: Order preservation
urls := []string{"url1", "url2", "url3"}
results := CheckURLs(urls)
assert.Equal(t, "url1", results[0].URL)
assert.Equal(t, "url2", results[1].URL)
assert.Equal(t, "url3", results[2].URL)

// Test 4: Empty slice
urls := []string{}
results := CheckURLs(urls)
assert.Equal(t, 0, len(results))

// Test 5: Timeout handling
// Simulate slow server that takes >5 seconds
urls := []string{"https://httpbin.org/delay/10"}
results := CheckURLs(urls)
assert.NotNil(t, results[0].Error)  // Should timeout
```

---

## 💡 Hints

1. Створи channel для результатів: `ch := make(chan URLStatus)`
2. Запусти goroutine для кожного URL
3. В кожній goroutine зроби HTTP GET і надішли результат в channel
4. В main goroutine збери результати з channel
5. Використай `http.Client` з timeout:
```go
client := &http.Client{
    Timeout: 5 * time.Second,
}
resp, err := client.Get(url)
```
6. Використай слайс для збереження порядку результатів

---

**Рішення:** `solutions/solution_02_url_checker.go`

**Good luck!** 🚀
