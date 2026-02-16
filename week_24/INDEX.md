# Week 24 - Complete Index 📚

## Всі 15 паттернів каналів

| # | Файл | Паттерн | Складність | Use Case |
|---|------|---------|------------|----------|
| 1 | `01_fan_in.go` | Fan-In | ⭐⭐ | Aggregating logs from multiple services |
| 2 | `02_fan_out.go` | Fan-Out | ⭐⭐ | Processing large datasets |
| 3 | `03_tee.go` | Tee | ⭐ | Broadcasting events |
| 4 | `04_transform.go` | Transform | ⭐ | ETL pipelines |
| 5 | `05_filter.go` | Filter | ⭐ | Data validation |
| 6 | `06_pipeline.go` | Pipeline | ⭐⭐⭐ | Request/Response transformation |
| 7 | `07_semaphore.go` | Semaphore | ⭐⭐ | Database connection pool |
| 8 | `08_barrier.go` | Barrier | ⭐⭐⭐ | Multi-phase computations |
| 9 | `09_promise.go` | Promise | ⭐⭐ | Async API calls |
| 10 | `10_future.go` | Future | ⭐⭐ | Lazy evaluation |
| 11 | `11_future_promise.go` | Future + Promise | ⭐⭐⭐ | HTTP client libraries |
| 12 | `12_generator.go` | Generator | ⭐ | Infinite sequences |
| 13 | `13_errgroup.go` | ErrGroup | ⭐⭐⭐ | Parallel API calls with error handling |
| 14 | `14_singleflight.go` | SingleFlight | ⭐⭐⭐ | Cache warming (thundering herd) |
| 15 | `15_rate_limiter.go` | Rate Limiter | ⭐⭐⭐ | API rate limiting |

---

## За категоріями

### Basic Patterns (3)
- **Fan-In** - Об'єднання каналів
- **Fan-Out** - Розподіл роботи
- **Tee** - Дублювання даних

### Data Processing (3)
- **Transform** - Перетворення даних
- **Filter** - Фільтрація
- **Pipeline** - Послідовна обробка

### Synchronization (2)
- **Semaphore** - Обмеження ресурсів
- **Barrier** - Синхронізація точки зустрічі

### Async Patterns (3)
- **Promise** - Відкладений результат
- **Future** - Async обчислення
- **Future + Promise** - Комбінований

### Advanced (4)
- **Generator** - Нескінченний потік
- **ErrGroup** - Групова обробка з помилками
- **SingleFlight** - Дедуплікація запитів
- **Rate Limiter** - Throttling

---

## За складністю

### Початковий рівень ⭐
```
03_tee.go
04_transform.go
05_filter.go
12_generator.go
```

### Середній рівень ⭐⭐
```
01_fan_in.go
02_fan_out.go
07_semaphore.go
09_promise.go
10_future.go
```

### Просунутий рівень ⭐⭐⭐
```
06_pipeline.go
08_barrier.go
11_future_promise.go
13_errgroup.go
14_singleflight.go
15_rate_limiter.go
```

---

## Швидкий запуск

```bash
# Всі паттерни
cd /Users/vkuzm/GolandProjects/golang_practice/week_24

# Basic
go run 01_fan_in.go
go run 02_fan_out.go
go run 03_tee.go

# Data Processing
go run 04_transform.go
go run 05_filter.go
go run 06_pipeline.go

# Synchronization
go run 07_semaphore.go
go run 08_barrier.go

# Async
go run 09_promise.go
go run 10_future.go
go run 11_future_promise.go

# Advanced
go run 12_generator.go
go run 13_errgroup.go
go run 14_singleflight.go
go run 15_rate_limiter.go
```

---

## Топ 5 Must-Know

```bash
go run 02_fan_out.go        # Worker pools
go run 06_pipeline.go       # ETL
go run 13_errgroup.go       # Error handling
go run 14_singleflight.go   # Deduplication
go run 15_rate_limiter.go   # Rate limiting
```

---

## Real-World Examples

### Microservices
- Fan-Out: Parallel service calls
- ErrGroup: Health checks
- Rate Limiter: API protection

### API Gateway
- Fan-In: Response aggregation
- SingleFlight: Cache warming
- Rate Limiter: Throttling

### Data Processing
- Pipeline: ETL stages
- Transform: Normalization
- Filter: Validation

### Worker Pools
- Fan-Out: Task distribution
- Semaphore: Resource limits
- Barrier: Phase sync

---

## Тестування

```bash
# Compile all
for f in *.go; do go build -o /dev/null "$f" && echo "✅ $f"; done

# Run all (slow)
for f in *.go; do echo "Running $f..." && go run "$f"; done
```

---

**15 Production-Ready Channel Patterns!** 🎉

Детальніше: `README.md` | `QUICKSTART.md`
