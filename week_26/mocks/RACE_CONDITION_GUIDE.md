# 🔍 Як перевірити Race Condition в Go

## 📚 Швидкий старт

Race Condition - це коли дві або більше goroutines одночасно намагаються отримати доступ до однієї змінної, і хоча б одна з них робить запис.

### Go має вбудований **Race Detector**! 🎉

---

## 🚀 Команди для перевірки

### 1️⃣ Запустити програму з race detector:
```bash
go run -race main.go
```

### 2️⃣ Запустити тести з race detector:
```bash
go test -race
go test -race -v              # детальний вивід
go test -race -v -run TestName  # конкретний тест
```

### 3️⃣ Перевірити весь проект:
```bash
go test -race ./...
```

### 4️⃣ Benchmark з race detector:
```bash
go test -race -bench=.
```

### 5️⃣ Збілдити з race detector (для manual testing):
```bash
go build -race myprogram.go
./myprogram
```

---

## 📋 Приклади в цьому проекті

### Запустити демонстрацію:
```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_26/mocks

# Звичайний запуск (побачиш втрачені інкременти)
go run race_condition_demo.go

# З race detector (побачиш де саме проблема)
go run -race race_condition_demo.go
```

### Запустити тести:
```bash
# Спочатку перемістимо demo файли
mkdir -p examples
mv demo_goroutine.go race_condition_demo.go examples/

# Тепер можна запустити тести
go test -race -v

# Повернути файли назад
mv examples/*.go .
rmdir examples
```

### Запустити тести для основного коду:
```bash
# Тестуємо URLService на race conditions
go test -race -v -run TestProcessAsync
go test -race -v -run TestProcessWithCallback
```

---

## 🎯 Що шукає Race Detector?

Race Detector виявляє:

- ❌ **Одночасний запис** в одну змінну з різних goroutines
- ❌ **Одночасне читання/запис** без синхронізації
- ❌ **Небезпечний доступ до maps**
- ❌ **Небезпечний доступ до slices**
- ❌ **Неправильне використання channels**

---

## 📊 Приклад виводу Race Detector

Коли race detector знаходить проблему, ви побачите:

```
==================
WARNING: DATA RACE
Write at 0x00c000014098 by goroutine 7:
  main.(*UnsafeCounter).Increment()
      /Users/vkuzm/GolandProjects/golang_practice/week_26/mocks/race_condition_demo.go:15 +0x38
  main.demoRaceCondition.func1()
      /Users/vkuzm/GolandProjects/golang_practice/week_26/mocks/race_condition_demo.go:28 +0x40

Previous write at 0x00c000014098 by goroutine 6:
  main.(*UnsafeCounter).Increment()
      /Users/vkuzm/GolandProjects/golang_practice/week_26/mocks/race_condition_demo.go:15 +0x38
  main.demoRaceCondition.func1()
      /Users/vkuzm/GolandProjects/golang_practice/week_26/mocks/race_condition_demo.go:28 +0x40
==================
```

### Що означає цей вивід:

1. **`WARNING: DATA RACE`** - знайдено race condition!
2. **`Write at 0x...`** - адреса пам'яті де відбувся запис
3. **`by goroutine 7`** - goroutine яка робила запис
4. **Стек викликів** - де саме в коді проблема
5. **`Previous write`** - попередній доступ з іншої goroutine

---

## ❌ Приклади небезпечного коду

### 1. Небезпечний лічильник:
```go
type UnsafeCounter struct {
    count int  // ❌ Небезпечно!
}

func (c *UnsafeCounter) Increment() {
    c.count++  // ❌ RACE CONDITION!
}

// Використання
counter := &UnsafeCounter{}
for i := 0; i < 100; i++ {
    go counter.Increment()  // ❌ Багато goroutines пишуть одночасно!
}
```

### 2. Небезпечна мапа:
```go
m := make(map[string]int)  // ❌ Звичайна map не thread-safe!

for i := 0; i < 10; i++ {
    go func(id int) {
        m[fmt.Sprintf("key-%d", id)] = id  // ❌ RACE CONDITION!
    }(i)
}
```

### 3. Небезпечний slice:
```go
var results []string  // ❌ Небезпечно!

for i := 0; i < 10; i++ {
    go func(id int) {
        results = append(results, fmt.Sprintf("result-%d", id))  // ❌ RACE!
    }(i)
}
```

---

## ✅ Як виправити Race Conditions

### 1. Використовуй `sync.Mutex`:
```go
type SafeCounter struct {
    mu    sync.Mutex
    count int
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()         // ✅ Блокуємо
    defer c.mu.Unlock() // ✅ Розблокуємо
    c.count++
}

func (c *SafeCounter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.count
}
```

### 2. Використовуй `sync.RWMutex` (якщо більше читань):
```go
type SafeCache struct {
    mu   sync.RWMutex
    data map[string]string
}

func (c *SafeCache) Get(key string) string {
    c.mu.RLock()         // ✅ Read lock
    defer c.mu.RUnlock()
    return c.data[key]
}

func (c *SafeCache) Set(key, value string) {
    c.mu.Lock()          // ✅ Write lock
    defer c.mu.Unlock()
    c.data[key] = value
}
```

### 3. Використовуй `channels`:
```go
type SafeCounterChannel struct {
    ch chan int
}

func NewSafeCounterChannel() *SafeCounterChannel {
    c := &SafeCounterChannel{ch: make(chan int)}
    go c.run()  // Одна goroutine обробляє всі запити
    return c
}

func (c *SafeCounterChannel) run() {
    count := 0
    for increment := range c.ch {
        count += increment  // ✅ Тільки одна goroutine має доступ!
    }
}

func (c *SafeCounterChannel) Increment() {
    c.ch <- 1  // ✅ Channels thread-safe!
}
```

### 4. Використовуй `sync.Map` (для concurrent maps):
```go
var safeMap sync.Map  // ✅ Thread-safe map!

for i := 0; i < 10; i++ {
    go func(id int) {
        safeMap.Store(fmt.Sprintf("key-%d", id), id)  // ✅ Безпечно!
    }(i)
}

// Читання
value, ok := safeMap.Load("key-1")

// Ітерація
safeMap.Range(func(key, value interface{}) bool {
    fmt.Printf("%s = %v\n", key, value)
    return true
})
```

### 5. Використовуй `sync/atomic` (для простих операцій):
```go
type AtomicCounter struct {
    count int64
}

func (c *AtomicCounter) Increment() {
    atomic.AddInt64(&c.count, 1)  // ✅ Атомарна операція!
}

func (c *AtomicCounter) Value() int64 {
    return atomic.LoadInt64(&c.count)  // ✅ Атомарне читання!
}
```

---

## 🎓 Порівняльна таблиця рішень

| Метод | Складність | Продуктивність | Коли використовувати |
|-------|-----------|----------------|---------------------|
| `sync.Mutex` | ⭐ Проста | ⭐⭐⭐ Швидко | Загальний випадок |
| `sync.RWMutex` | ⭐⭐ Середня | ⭐⭐⭐⭐ Дуже швидко | Багато читань, мало записів |
| `channels` | ⭐⭐⭐ Складна | ⭐⭐ Повільніше | Комунікація між goroutines |
| `sync.Map` | ⭐ Проста | ⭐⭐⭐ Швидко | Concurrent map access |
| `sync/atomic` | ⭐ Проста | ⭐⭐⭐⭐⭐ Найшвидше | Прості операції (int, bool) |

---

## ⚠️ Важливі нюанси

### 1. Performance Impact:
```
Race detector збільшує:
- ⏱️  Час виконання: ~5-10x
- 💾 Використання пам'яті: ~5-10x
```

### 2. Використовуй тільки для розробки:
```bash
# ✅ ДОБРЕ - для тестів
go test -race

# ❌ ПОГАНО - для production
go build -race && ./app  # НЕ РОБИ ТАК В ПРОДІ!
```

### 3. Race detector НЕ знаходить всі race conditions:
- Знаходить тільки ті, які **відбулись** під час виконання
- Тести повинні покривати concurrent scenarios
- Потрібно запускати з різними входами

---

## 🔧 Best Practices

### 1. ✅ Завжди запускай тести з `-race`:
```bash
# В CI/CD pipeline
go test -race -coverprofile=coverage.out ./...
```

### 2. ✅ Тестуй concurrent scenarios:
```go
func TestConcurrent(t *testing.T) {
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            // твій код тут
        }()
    }
    wg.Wait()
}
```

### 3. ✅ Використовуй `go vet`:
```bash
go vet ./...  # Знаходить деякі concurrent проблеми
```

### 4. ✅ Code review checklist:
- [ ] Чи є shared state між goroutines?
- [ ] Чи використовується синхронізація (mutex/channels)?
- [ ] Чи є доступ до maps/slices з різних goroutines?
- [ ] Чи тести покривають concurrent scenarios?
- [ ] Чи запускались тести з `-race`?

---

## 📁 Файли в цьому проекті

| Файл | Опис | Як запустити |
|------|------|--------------|
| `race_condition_demo.go` | Демонстрація race conditions | `go run race_condition_demo.go` |
| `race_test.go` | Тести з race conditions | `go test -race -v` |
| `main.go` | Основний код URLService | `go run main.go mock_url_sorter.go` |
| `RACE_CONDITION_GUIDE.md` | Цей guide | - |

---

## 🎯 Практичні вправи

### Вправа 1: Знайди race condition
```bash
# Запусти демо БЕЗ race detector
go run race_condition_demo.go

# Що побачив? Скільки інкрементів втрачено?
```

### Вправа 2: Виправ race condition
```bash
# Запусти демо З race detector
go run -race race_condition_demo.go

# Race detector покаже де проблема!
```

### Вправа 3: Перевір свій код
```bash
# Запусти тести для URLService
go test -race -v -run TestProcessAsync

# Чи є race conditions?
```

---

## 📚 Додаткові матеріали

### Офіційна документація:
- [Go Race Detector](https://go.dev/doc/articles/race_detector)
- [Go Memory Model](https://go.dev/ref/mem)
- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)

### Корисні статті:
- [Understanding Race Conditions](https://www.ardanlabs.com/blog/2013/09/detecting-race-conditions-with-go.html)
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)

---

## 💡 Коротке резюме

### Як перевірити race condition:
```bash
go test -race        # Для тестів
go run -race main.go # Для програми
go build -race       # Для білда (тільки розробка!)
```

### Як виправити:
1. ✅ `sync.Mutex` - для shared state
2. ✅ `channels` - для комунікації
3. ✅ `sync.Map` - для concurrent maps
4. ✅ `sync/atomic` - для простих операцій
5. ✅ Уникай shared state (якщо можливо)

### Пам'ятай:
- ⚠️ Race detector знаходить race conditions тільки під час виконання
- ⚠️ Не використовуй `-race` в production
- ⚠️ Завжди тестуй concurrent code
- ⚠️ Code review обов'язковий для concurrent code

---

## 🚀 Готовий перевірити свій код?

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_26/mocks
go test -race -v
```

Успіхів! 🎉
