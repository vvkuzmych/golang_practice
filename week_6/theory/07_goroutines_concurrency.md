# Goroutines і Конкурентність в Go

Go має вбудовану підтримку конкурентності через goroutines і channels.

---

## 📖 Зміст

1. [Goroutines](#1-goroutines)
2. [Channels](#2-channels)
3. [Select](#3-select)
4. [Sync Package](#4-sync-package)
5. [Конкурентні Патерни](#5-конкурентні-патерни)
6. [Common Pitfalls](#6-common-pitfalls)

---

## 1. Goroutines

### Що таке Goroutine?

**Goroutine** - це легковісний потік, керований Go runtime.

```go
// Звичайна функція - виконується синхронно
func sayHello() {
    fmt.Println("Hello")
}

func main() {
    sayHello()        // Блокує виконання до завершення
    fmt.Println("World")
}
// Output: Hello, World
```

```go
// Goroutine - виконується асинхронно
func main() {
    go sayHello()     // Запускається в окремій goroutine
    fmt.Println("World")
}
// Output: World (і можливо "Hello", якщо встигне)
```

### Створення Goroutines

```go
package main

import (
    "fmt"
    "time"
)

func task(id int) {
    fmt.Printf("Task %d started\n", id)
    time.Sleep(time.Second)
    fmt.Printf("Task %d finished\n", id)
}

func main() {
    // Запускаємо 5 goroutines
    for i := 1; i <= 5; i++ {
        go task(i)
    }
    
    // Чекаємо, щоб goroutines встигли завершитись
    time.Sleep(2 * time.Second)
    fmt.Println("All tasks completed")
}
```

### Anonymous Functions

```go
func main() {
    // Goroutine з анонімною функцією
    go func() {
        fmt.Println("Running in goroutine")
    }()
    
    // З параметрами
    for i := 0; i < 5; i++ {
        go func(n int) {
            fmt.Printf("Number: %d\n", n)
        }(i) // Передаємо i як аргумент
    }
    
    time.Sleep(time.Second)
}
```

### ⚠️ Common Mistake: Closure

```go
// ❌ Неправильно - всі goroutines бачать одне значення i
func main() {
    for i := 0; i < 5; i++ {
        go func() {
            fmt.Println(i) // Може вивести 5, 5, 5, 5, 5
        }()
    }
    time.Sleep(time.Second)
}

// ✅ Правильно - передаємо i як параметр
func main() {
    for i := 0; i < 5; i++ {
        go func(n int) {
            fmt.Println(n) // Виведе 0, 1, 2, 3, 4
        }(i)
    }
    time.Sleep(time.Second)
}
```

---

## 2. Channels

### Що таке Channel?

**Channel** - це типізований канал для комунікації між goroutines.

```go
// Створення channel
ch := make(chan int)        // unbuffered channel
ch := make(chan int, 10)    // buffered channel (buffer size 10)
```

### Відправка і Отримання

```go
ch := make(chan string)

// Відправка в channel
ch <- "Hello"

// Отримання з channel
message := <-ch

// Закриття channel
close(ch)
```

### Unbuffered Channel (синхронний)

```go
func main() {
    ch := make(chan int)
    
    go func() {
        ch <- 42 // Блокується доки хтось не прочитає
    }()
    
    value := <-ch // Блокується доки хтось не запише
    fmt.Println(value) // 42
}
```

### Buffered Channel (асинхронний)

```go
func main() {
    ch := make(chan int, 3) // buffer size 3
    
    // Можна записати 3 значення без блокування
    ch <- 1
    ch <- 2
    ch <- 3
    // ch <- 4  // Блокується, якщо buffer повний
    
    // Читаємо
    fmt.Println(<-ch) // 1
    fmt.Println(<-ch) // 2
    fmt.Println(<-ch) // 3
}
```

### Channel Direction

```go
// Send-only channel
func sender(ch chan<- int) {
    ch <- 42
    // val := <-ch  // Compile error! Не можна читати
}

// Receive-only channel
func receiver(ch <-chan int) {
    val := <-ch
    // ch <- 42  // Compile error! Не можна писати
}

func main() {
    ch := make(chan int)
    
    go sender(ch)
    go receiver(ch)
    
    time.Sleep(time.Second)
}
```

### Range over Channel

```go
func main() {
    ch := make(chan int, 5)
    
    // Відправляємо значення
    go func() {
        for i := 1; i <= 5; i++ {
            ch <- i
        }
        close(ch) // Важливо закрити!
    }()
    
    // Читаємо доки channel не закритий
    for value := range ch {
        fmt.Println(value)
    }
}
```

### Перевірка на закритий Channel

```go
func main() {
    ch := make(chan int, 2)
    ch <- 1
    ch <- 2
    close(ch)
    
    // Перевірка, чи channel закритий
    value, ok := <-ch
    fmt.Println(value, ok) // 1 true
    
    value, ok = <-ch
    fmt.Println(value, ok) // 2 true
    
    value, ok = <-ch
    fmt.Println(value, ok) // 0 false (channel закритий)
}
```

---

## 3. Select

### Що таке Select?

**Select** дозволяє чекати на кілька channel операцій одночасно.

```go
func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)
    
    go func() {
        time.Sleep(1 * time.Second)
        ch1 <- "from ch1"
    }()
    
    go func() {
        time.Sleep(2 * time.Second)
        ch2 <- "from ch2"
    }()
    
    // Select чекає на перший доступний channel
    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-ch1:
            fmt.Println(msg1)
        case msg2 := <-ch2:
            fmt.Println(msg2)
        }
    }
}
```

### Default Case (non-blocking)

```go
func main() {
    ch := make(chan string)
    
    select {
    case msg := <-ch:
        fmt.Println(msg)
    default:
        fmt.Println("No message received")
    }
}
```

### Timeout Pattern

```go
func main() {
    ch := make(chan string)
    
    go func() {
        time.Sleep(2 * time.Second)
        ch <- "result"
    }()
    
    select {
    case result := <-ch:
        fmt.Println("Received:", result)
    case <-time.After(1 * time.Second):
        fmt.Println("Timeout!")
    }
}
```

### Multiple Channels

```go
func worker(id int, jobs <-chan int, results chan<- int) {
    for job := range jobs {
        fmt.Printf("Worker %d processing job %d\n", id, job)
        time.Sleep(time.Second)
        results <- job * 2
    }
}

func main() {
    jobs := make(chan int, 5)
    results := make(chan int, 5)
    
    // Запускаємо 3 workers
    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }
    
    // Відправляємо 5 jobs
    for j := 1; j <= 5; j++ {
        jobs <- j
    }
    close(jobs)
    
    // Збираємо результати
    for r := 1; r <= 5; r++ {
        result := <-results
        fmt.Println("Result:", result)
    }
}
```

---

## 4. Sync Package

### sync.WaitGroup

Чекати завершення goroutines без `time.Sleep`.

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done() // Декремент counter при завершенні
    
    fmt.Printf("Worker %d starting\n", id)
    time.Sleep(time.Second)
    fmt.Printf("Worker %d done\n", id)
}

func main() {
    var wg sync.WaitGroup
    
    for i := 1; i <= 5; i++ {
        wg.Add(1) // Інкремент counter
        go worker(i, &wg)
    }
    
    wg.Wait() // Блокується доки counter != 0
    fmt.Println("All workers completed")
}
```

### sync.Mutex

Захист спільних даних від race conditions.

```go
package main

import (
    "fmt"
    "sync"
)

type SafeCounter struct {
    mu    sync.Mutex
    count int
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}

func (c *SafeCounter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.count
}

func main() {
    counter := &SafeCounter{}
    var wg sync.WaitGroup
    
    // 1000 goroutines інкрементують counter
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter.Increment()
        }()
    }
    
    wg.Wait()
    fmt.Println("Final count:", counter.Value()) // 1000
}
```

### sync.RWMutex

Дозволяє кілька читачів або одного писаря.

```go
type Cache struct {
    mu    sync.RWMutex
    data  map[string]string
}

func (c *Cache) Get(key string) (string, bool) {
    c.mu.RLock() // Читання - кілька goroutines одночасно
    defer c.mu.RUnlock()
    
    value, ok := c.data[key]
    return value, ok
}

func (c *Cache) Set(key, value string) {
    c.mu.Lock() // Запис - тільки одна goroutine
    defer c.mu.Unlock()
    
    c.data[key] = value
}
```

### sync.Once

Виконати функцію тільки один раз (thread-safe).

```go
var (
    instance *Database
    once     sync.Once
)

func GetDatabase() *Database {
    once.Do(func() {
        fmt.Println("Creating database instance")
        instance = &Database{}
    })
    return instance
}
```

### sync.Pool

Переиспользування об'єктів для зменшення GC pressure.

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func processData(data string) {
    // Отримуємо buffer з pool
    buf := bufferPool.Get().(*bytes.Buffer)
    defer bufferPool.Put(buf) // Повертаємо в pool
    
    buf.Reset()
    buf.WriteString(data)
    // ... обробка
}
```

---

## 5. Конкурентні Патерни

### Worker Pool

```go
func workerPool(jobs <-chan int, results chan<- int, numWorkers int) {
    var wg sync.WaitGroup
    
    // Запускаємо workers
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            for job := range jobs {
                fmt.Printf("Worker %d processing job %d\n", workerID, job)
                time.Sleep(100 * time.Millisecond)
                results <- job * 2
            }
        }(i)
    }
    
    wg.Wait()
    close(results)
}

func main() {
    jobs := make(chan int, 10)
    results := make(chan int, 10)
    
    go workerPool(jobs, results, 3)
    
    // Відправляємо jobs
    for i := 1; i <= 10; i++ {
        jobs <- i
    }
    close(jobs)
    
    // Збираємо результати
    for result := range results {
        fmt.Println("Result:", result)
    }
}
```

### Pipeline

```go
// Generator
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

// Square
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

// Filter
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

func main() {
    // Pipeline: generate → square → filterEven
    nums := generate(1, 2, 3, 4, 5)
    squared := square(nums)
    even := filterEven(squared)
    
    for result := range even {
        fmt.Println(result) // 4, 16
    }
}
```

### Fan-Out, Fan-In

```go
// Fan-Out: один input, кілька workers
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
        for n := range input {
            out <- n * 2
        }
        close(out)
    }()
    return out
}

// Fan-In: кілька inputs, один output
func fanIn(channels ...<-chan int) <-chan int {
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

### Context for Cancellation

```go
func worker(ctx context.Context, id int) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %d stopped\n", id)
            return
        default:
            fmt.Printf("Worker %d working\n", id)
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    
    // Запускаємо workers
    for i := 1; i <= 3; i++ {
        go worker(ctx, i)
    }
    
    // Працюємо 2 секунди, потім зупиняємо
    time.Sleep(2 * time.Second)
    cancel() // Зупиняє всі workers
    
    time.Sleep(time.Second)
}
```

---

## 6. Common Pitfalls

### 1. Race Conditions

```go
// ❌ Race condition
var counter int

func increment() {
    for i := 0; i < 1000; i++ {
        counter++ // NOT thread-safe!
    }
}

func main() {
    go increment()
    go increment()
    time.Sleep(time.Second)
    fmt.Println(counter) // Може бути < 2000
}

// ✅ Виправлено з Mutex
var (
    counter int
    mu      sync.Mutex
)

func increment() {
    for i := 0; i < 1000; i++ {
        mu.Lock()
        counter++
        mu.Unlock()
    }
}
```

**Перевірка Race Conditions:**
```bash
go run -race main.go
```

### 2. Goroutine Leaks

```go
// ❌ Goroutine leak - ніколи не завершиться
func leak() {
    ch := make(chan int)
    go func() {
        val := <-ch // Блокується назавжди
        fmt.Println(val)
    }()
    // Забули відправити значення в channel
}

// ✅ Виправлено з timeout
func noLeak() {
    ch := make(chan int)
    go func() {
        select {
        case val := <-ch:
            fmt.Println(val)
        case <-time.After(1 * time.Second):
            fmt.Println("Timeout")
            return
        }
    }()
}
```

### 3. Deadlock

```go
// ❌ Deadlock
func main() {
    ch := make(chan int)
    ch <- 42 // Блокується назавжди (unbuffered channel)
    // fatal error: all goroutines are asleep - deadlock!
}

// ✅ Виправлено
func main() {
    ch := make(chan int, 1) // buffered
    ch <- 42
    fmt.Println(<-ch)
}
```

### 4. Закриття закритого Channel

```go
// ❌ Panic
ch := make(chan int)
close(ch)
close(ch) // panic: close of closed channel

// ✅ Використовуйте sync.Once
var once sync.Once
once.Do(func() {
    close(ch)
})
```

---

## ✅ Best Practices

1. **Завжди закривайте channels** - коли більше не потрібні
2. **Використовуйте WaitGroup** - замість `time.Sleep`
3. **Mutex для shared state** - захищайте спільні дані
4. **Buffered channels** - якщо знаєте розмір
5. **Context для cancellation** - graceful shutdown
6. **Перевіряйте race conditions** - `go run -race`
7. **Не передавайте channels по channels** - занадто складно
8. **Обмежуйте кількість goroutines** - Worker Pool pattern
9. **Profile performance** - `go tool pprof`

---

## 📊 Порівняльна таблиця

| Концепція | Коли використовувати |
|-----------|---------------------|
| **Goroutine** | Асинхронні задачі |
| **Channel** | Комунікація між goroutines |
| **Buffered Channel** | Відомий розмір черги |
| **Select** | Множинні channel операції |
| **Mutex** | Захист shared state |
| **RWMutex** | Багато читачів, мало писарів |
| **WaitGroup** | Чекати завершення goroutines |
| **Once** | Виконати один раз (Singleton) |
| **Context** | Cancellation, timeouts |

---

## 🚀 Приклад: HTTP Server з Goroutines

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Кожен request обробляється в окремій goroutine автоматично
    fmt.Fprintf(w, "Request from: %s\n", r.RemoteAddr)
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
```

---

**Goroutines - це супер-сила Go!** 💪🚀

**Далі:** Практичні вправи з конкурентністю!
