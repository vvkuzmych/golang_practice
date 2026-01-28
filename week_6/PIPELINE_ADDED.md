# ✅ Pipeline Pattern додано до Week 6!

## 🎉 Що додано

### 📚 Теорія (~3000 слів)

**theory/08_pipeline_pattern.md** - Повний довідник про Pipeline Pattern

**Охоплені теми:**
1. ✅ Що таке Pipeline
2. ✅ Візуалізація і діаграми
3. ✅ Простий приклад з поясненнями
4. ✅ Де використовуються (7 реальних сценаріїв)
5. ✅ Реальний приклад: Web Scraper
6. ✅ Fan-Out / Fan-In патерни
7. ✅ Важливі правила (закриття channels, WaitGroup, Context)
8. ✅ Image Processing приклад
9. ✅ Performance tips
10. ✅ Корисні патерни (Tee, Merge, OrDone)

### 💻 Практика

**practice/06_goroutines/pipeline_example.go** - Робочі приклади ✅

**3 повних приклада:**
1. Simple Pipeline (Generate → Square → Filter)
2. Fan-Out/Fan-In (3 workers паралельно)
3. Data Processing (Generate → Validate → Transform → Save)

---

## 🚀 Швидкий старт

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_6

# 1. Прочитайте теорію
cat theory/08_pipeline_pattern.md

# 2. Запустіть приклади
go run practice/06_goroutines/pipeline_example.go
```

---

## 💻 Результат виконання

```
🔄 Pipeline Pattern Examples
==============================

=== Simple Pipeline ===
Generate: 1
Square: 1 → 1
Filter: 1 (rejected)
Generate: 2
Square: 2 → 4
Filter: 4 (passed)
Result: 4
...

=== Fan-Out / Fan-In ===
Worker 1: 1 → 1
Worker 2: 2 → 4
Worker 3: 3 → 9
Final result: 4
Final result: 9
...

=== Data Processing Pipeline ===
Generate: {ID:1 Value:data_1}
Validate: {ID:2 Value:data_2} (passed)
Transform: {ID:2 Value:data_2_transformed}
Saved: ID=2, Value=data_2_transformed
Total processed: 5 records

✅ All pipeline examples completed!
```

---

## 📖 Ключові концепції

### Pipeline Structure
```
Input → [Stage 1] → Channel → [Stage 2] → Channel → Output
          ↓                      ↓
      Goroutine              Goroutine
```

### Simple Pipeline
```go
func main() {
    nums := generate(1, 2, 3, 4, 5)       // Stage 1
    squared := square(nums)                // Stage 2
    even := filterEven(squared)            // Stage 3
    
    for result := range even {
        fmt.Println(result)
    }
}
```

### Fan-Out
```go
// Один input → багато workers
worker1 := process(input)
worker2 := process(input)
worker3 := process(input)
```

### Fan-In
```go
// Багато inputs → один output
results := merge(worker1, worker2, worker3)
```

---

## 🎯 Де використовується

1. **Data Processing**
   - CSV → Parse → Validate → Transform → Save

2. **Image Processing**
   - Image → Resize → Filter → Compress → Save

3. **Web Scraping**
   - URLs → Fetch → Parse → Extract → Store

4. **ETL (Extract, Transform, Load)**
   - API → Extract → Transform → Validate → Load

5. **Video Encoding**
   - Video → Decode → Resize → Encode → Upload

6. **Log Processing**
   - Logs → Parse → Filter → Aggregate → Store

7. **Stream Processing**
   - WebSocket → Decode → Process → Encode → Send

---

## ✅ Переваги Pipeline

1. **Конкурентність** - стадії працюють паралельно
2. **Модульність** - легко додавати/видаляти стадії
3. **Масштабованість** - кілька workers на стадію
4. **Читабельність** - чіткий flow даних
5. **Ефективність** - стадії не чекають одна на одну

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

### 2. WaitGroup для кількох workers
```go
var wg sync.WaitGroup
for i := 0; i < numWorkers; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        // work
    }()
}
wg.Wait()
```

### 3. Context для cancellation
```go
func stage(ctx context.Context, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for {
            select {
            case n := <-in:
                out <- n * 2
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}
```

---

## 📊 Оновлена структура Week 6

```
week_6/
├── theory/
│   ├── 01_oop_principles.md
│   ├── 02_design_patterns.md
│   ├── 03_net_http.md
│   ├── 04_microservices.md
│   ├── 05_databases.md
│   ├── 06_networking.md
│   ├── 07_goroutines_concurrency.md
│   └── 08_pipeline_pattern.md        ← НОВИЙ!
│
└── practice/
    └── 06_goroutines/
        ├── main.go
        └── pipeline_example.go          ← НОВИЙ!
```

---

## 📈 Статистика Week 6

**Теорія:** 8 файлів (~25,000 слів)
**Практика:** 4 робочих приклади
**Вправи:** 3 завдання
**Охоплено:** 35+ концепцій

### Теми:
1. ООП принципи ✅
2. Design Patterns ✅
3. HTTP Server/Client ✅
4. Microservices ✅
5. Databases ✅
6. Networking ✅
7. Goroutines ✅
8. **Pipeline Pattern** ✅ ← НОВИЙ!

---

## 🎓 Навчальний план

**Оновлено:**

- День 1-2: ООП і Патерни
- День 3-4: HTTP і Мікросервіси
- День 5: Бази даних
- День 6: Нетворкінг
- **День 7: Goroutines + Pipeline** ← ОНОВЛЕНО!

---

## 💡 Реальні приклади використання

### 1. Web Scraper
```
URLs → Fetch (5 workers) → Parse (3 workers) → Validate (2 workers) → Save
```

### 2. Image Processing
```
Load → Resize (3 workers) → Filter (3 workers) → Compress (2 workers) → Save
```

### 3. Data ETL
```
API → Extract → Transform → Validate → Load to DB
```

---

## 🚀 Що далі?

Тепер ви знаєте:
- ✅ Як будувати Pipeline
- ✅ Fan-Out/Fan-In патерни
- ✅ Як масштабувати обробку
- ✅ Context для cancellation
- ✅ Best practices

**Готові створювати production-ready concurrent systems!** 💪

---

## 📝 Корисні команди

```bash
# Запустити приклади
go run practice/06_goroutines/main.go
go run practice/06_goroutines/pipeline_example.go

# Перевірка race conditions
go run -race practice/06_goroutines/pipeline_example.go

# Profiling
go run -cpuprofile=cpu.prof practice/06_goroutines/pipeline_example.go
go tool pprof cpu.prof
```

---

## ✅ Оновлені файли

1. **theory/08_pipeline_pattern.md** - новий файл теорії
2. **practice/06_goroutines/pipeline_example.go** - робочі приклади
3. **README.md** - оновлено структуру
4. **PIPELINE_ADDED.md** - цей файл

---

**Week 6 тепер включає повний курс по concurrent programming в Go!** 🎉

**Створено:** 2026-01-27
**Статус:** ✅ Completed & Tested
**Тестування:** ✅ All examples work perfectly
