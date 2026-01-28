# 🔄 Pipeline Pattern - Стисло

## Що таке Pipeline?

**Pipeline** - це послідовність стадій обробки даних, де:
- Кожна стадія виконується в окремій goroutine
- Стадії з'єднані через channels
- Дані течуть від одної стадії до іншої

```
Input → Stage 1 → Stage 2 → Stage 3 → Output
```

---

## 📊 Візуалізація

```
[1,2,3,4,5] → [Generate] → [Square] → [Filter] → [Sum]
              channel 1    channel 2   channel 3
```

---

## 💻 Простий приклад

```go
// Stage 1: Generate numbers
func generate(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

// Stage 2: Square numbers
func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

// Stage 3: Filter even
func filterEven(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            if n%2 == 0 {
                out <- n
            }
        }
        close(out)
    }()
    return out
}

// Використання
func main() {
    // Створюємо pipeline
    nums := generate(1, 2, 3, 4, 5)       // → 1, 2, 3, 4, 5
    squared := square(nums)                // → 1, 4, 9, 16, 25
    even := filterEven(squared)            // → 4, 16
    
    // Отримуємо результат
    for result := range even {
        fmt.Println(result) // Output: 4, 16
    }
}
```

---

## 🎯 Де використовуються?

### 1. **Data Processing**
```go
CSV → Parse → Validate → Transform → Save to DB
```

### 2. **Image Processing**
```go
Image → Resize → Apply Filter → Compress → Save
```

### 3. **Log Processing**
```go
Logs → Parse → Filter → Aggregate → Store
```

### 4. **ETL (Extract, Transform, Load)**
```go
API → Extract → Transform → Validate → Load to DB
```

### 5. **Stream Processing**
```go
WebSocket → Decode → Process → Encode → Send
```

### 6. **Video Encoding**
```go
Video → Decode → Resize → Encode → Upload
```

### 7. **Web Scraping**
```go
URLs → Fetch → Parse → Extract → Store
```

---

## ✅ Переваги

1. **Конкурентність** - кожна стадія працює паралельно
2. **Модульність** - легко додати/видалити стадії
3. **Масштабованість** - можна запустити кілька workers на стадію
4. **Читабельність** - чіткий flow даних
5. **Ефективність** - стадії не чекають одна на одну

---

## 📈 Реальний приклад: Web Scraper

```go
package main

import (
    "fmt"
    "net/http"
    "sync"
)

func main() {
    // Pipeline: URLs → Fetch → Parse → Validate → Save
    urls := generateURLs("https://example.com", 100)
    pages := fetchPages(urls, 5)        // 5 concurrent fetchers
    parsed := parseContent(pages, 3)    // 3 parsers
    validated := validate(parsed, 2)    // 2 validators
    
    // Save to DB
    for item := range validated {
        saveToDatabase(item)
    }
}

// Stage 1: Generate URLs
func generateURLs(baseURL string, count int) <-chan string {
    out := make(chan string)
    go func() {
        defer close(out)
        for i := 1; i <= count; i++ {
            out <- fmt.Sprintf("%s/page/%d", baseURL, i)
        }
    }()
    return out
}

// Stage 2: Fetch pages (with multiple workers)
func fetchPages(urls <-chan string, numWorkers int) <-chan string {
    out := make(chan string)
    var wg sync.WaitGroup
    
    // Запускаємо N workers
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            for url := range urls {
                content := httpGet(url)
                fmt.Printf("Worker %d fetched: %s\n", workerID, url)
                out <- content
            }
        }(i)
    }
    
    // Закриваємо channel після завершення всіх workers
    go func() {
        wg.Wait()
        close(out)
    }()
    
    return out
}

// Stage 3: Parse content
func parseContent(pages <-chan string, numWorkers int) <-chan ParsedData {
    out := make(chan ParsedData)
    var wg sync.WaitGroup
    
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for page := range pages {
                data := parse(page)
                out <- data
            }
        }()
    }
    
    go func() {
        wg.Wait()
        close(out)
    }()
    
    return out
}

func httpGet(url string) string {
    resp, _ := http.Get(url)
    defer resp.Body.Close()
    // ... read body
    return "page content"
}
```

---

## 🔥 Fan-Out / Fan-In

### Fan-Out (один input → багато workers)

```
         → [Worker 1] →
Input →  → [Worker 2] → results
         → [Worker 3] →
```

```go
func fanOut(input <-chan int, numWorkers int) []<-chan int {
    channels := make([]<-chan int, numWorkers)
    
    for i := 0; i < numWorkers; i++ {
        channels[i] = worker(input)
    }
    
    return channels
}

func worker(input <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range input {
            // Heavy computation
            result := n * n
            out <- result
        }
    }()
    return out
}
```

### Fan-In (багато inputs → один output)

```
[Worker 1] →
[Worker 2] → merge → results
[Worker 3] →
```

```go
func fanIn(channels ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    
    // Для кожного input channel
    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for n := range c {
                out <- n
            }
        }(ch)
    }
    
    // Закриваємо output після всіх inputs
    go func() {
        wg.Wait()
        close(out)
    }()
    
    return out
}
```

### Повний приклад Fan-Out/Fan-In

```go
func main() {
    // Input
    input := generate(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
    
    // Fan-Out: 3 workers обробляють паралельно
    worker1 := square(input)
    worker2 := square(input)
    worker3 := square(input)
    
    // Fan-In: об'єднуємо результати
    results := fanIn(worker1, worker2, worker3)
    
    // Output
    for result := range results {
        fmt.Println(result)
    }
}
```

---

## ⚠️ Важливі правила

### 1. Завжди закривайте channels

```go
func stage(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out) // ✅ Обов'язково!
        for n := range in {
            out <- n * 2
        }
    }()
    return out
}
```

### 2. Використовуйте WaitGroup

```go
func multiWorker(in <-chan int, numWorkers int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for n := range in {
                out <- n * 2
            }
        }()
    }
    
    go func() {
        wg.Wait()
        close(out) // Закриваємо після всіх workers
    }()
    
    return out
}
```

### 3. Context для Cancellation

```go
func stage(ctx context.Context, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for {
            select {
            case n, ok := <-in:
                if !ok {
                    return
                }
                out <- n * 2
            case <-ctx.Done():
                fmt.Println("Stage cancelled")
                return
            }
        }
    }()
    return out
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    nums := generate(1, 2, 3, 4, 5)
    squared := stage(ctx, nums)
    
    // Cancel після 1 секунди
    go func() {
        time.Sleep(1 * time.Second)
        cancel()
    }()
    
    for result := range squared {
        fmt.Println(result)
    }
}
```

### 4. Buffered Channels для performance

```go
// ❌ Unbuffered - повільно
func slow(in <-chan int) <-chan int {
    out := make(chan int) // size 0
    // ...
}

// ✅ Buffered - швидше
func fast(in <-chan int) <-chan int {
    out := make(chan int, 100) // buffer size 100
    // ...
}
```

---

## 📊 Реальний приклад: Image Processing

```go
package main

import (
    "fmt"
    "image"
    "sync"
)

func main() {
    // Pipeline: Load → Resize → Filter → Compress → Save
    images := loadImages("./photos", 100)
    resized := resize(images, 3)        // 3 workers
    filtered := applyFilter(resized, 3) // 3 workers
    compressed := compress(filtered, 2) // 2 workers
    
    for img := range compressed {
        save(img)
    }
}

func loadImages(dir string, count int) <-chan image.Image {
    out := make(chan image.Image)
    go func() {
        defer close(out)
        for i := 0; i < count; i++ {
            img := loadImage(fmt.Sprintf("%s/img%d.jpg", dir, i))
            out <- img
        }
    }()
    return out
}

func resize(images <-chan image.Image, numWorkers int) <-chan image.Image {
    out := make(chan image.Image, 10)
    var wg sync.WaitGroup
    
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            for img := range images {
                resized := resizeImage(img, 800, 600)
                fmt.Printf("Worker %d resized image\n", workerID)
                out <- resized
            }
        }(i)
    }
    
    go func() {
        wg.Wait()
        close(out)
    }()
    
    return out
}
```

---

## 🎯 Коли використовувати Pipeline?

### ✅ **Використовуйте, коли:**

- Послідовна обробка даних
- Різні стадії з різною швидкістю
- Потрібна конкурентність
- Stream processing
- Кожна стадія незалежна
- Великі обсяги даних

### ❌ **Не використовуйте, коли:**

- Проста синхронна логіка
- Стадії залежні від результату всіх попередніх
- Overhead більший за користь
- Малі обсяги даних
- Debugging критичний (складніше)

---

## 📝 Порівняння підходів

### Без Pipeline (послідовно)

```go
func processData(items []int) []int {
    // Stage 1
    for i := range items {
        items[i] = items[i] * 2
    }
    
    // Stage 2
    var filtered []int
    for _, item := range items {
        if item > 10 {
            filtered = append(filtered, item)
        }
    }
    
    // Stage 3
    for i := range filtered {
        filtered[i] = filtered[i] + 1
    }
    
    return filtered
}
```

### З Pipeline (конкурентно)

```go
func processData(items []int) <-chan int {
    input := generate(items...)
    doubled := double(input)
    filtered := filter(doubled, 10)
    incremented := increment(filtered)
    return incremented
}
```

---

## 🚀 Performance Tips

1. **Buffered channels** - зменшують блокування
2. **Правильна кількість workers** - залежить від CPU cores
3. **Profiling** - використовуйте `pprof` для оптимізації
4. **Batch processing** - обробляйте дані пачками

```go
// Batch processing в pipeline
func batchProcess(in <-chan int, batchSize int) <-chan []int {
    out := make(chan []int)
    go func() {
        defer close(out)
        batch := make([]int, 0, batchSize)
        
        for n := range in {
            batch = append(batch, n)
            if len(batch) >= batchSize {
                out <- batch
                batch = make([]int, 0, batchSize)
            }
        }
        
        if len(batch) > 0 {
            out <- batch
        }
    }()
    return out
}
```

---

## 📚 Корисні патерни

### 1. Tee (розгалуження)

```go
func tee(in <-chan int) (<-chan int, <-chan int) {
    out1 := make(chan int)
    out2 := make(chan int)
    
    go func() {
        defer close(out1)
        defer close(out2)
        
        for n := range in {
            out1 <- n
            out2 <- n
        }
    }()
    
    return out1, out2
}
```

### 2. Merge (об'єднання)

```go
func merge(channels ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    
    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for n := range c {
                out <- n
            }
        }(ch)
    }
    
    go func() {
        wg.Wait()
        close(out)
    }()
    
    return out
}
```

### 3. OrDone (з Context)

```go
func orDone(ctx context.Context, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for {
            select {
            case <-ctx.Done():
                return
            case v, ok := <-in:
                if !ok {
                    return
                }
                select {
                case out <- v:
                case <-ctx.Done():
                    return
                }
            }
        }
    }()
    return out
}
```

---

## 🎓 Висновок

**Pipeline Pattern** - це потужний інструмент для:
- Concurrent data processing
- Stream processing
- ETL tasks
- Real-time data transformation

**Ключові переваги:**
- Модульність коду
- Легке масштабування
- Ефективне використання CPU
- Чистий і читабельний код

**Пам'ятайте:**
- Завжди закривайте channels
- Використовуйте WaitGroup
- Context для cancellation
- Profile перед оптимізацією

---

**Pipeline робить ваш concurrent код модульним, масштабованим і ефективним!** 🚀

**Далі:** Практикуйте з реальними задачами!
