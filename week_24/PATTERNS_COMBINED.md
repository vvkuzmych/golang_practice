# Комбіновані паттерни 🔥

## Real-World приклади комбінування патернів

---

## 1. Worker Pool з Rate Limiting

**Паттерни**: Fan-Out + Semaphore + Rate Limiter

```go
// Producer → Rate Limiter → Fan-Out Workers → Results
```

**Use case**: API scraper з обмеженням запитів

```bash
# Див:
go run 02_fan_out.go   # Fan-Out
go run 07_semaphore.go # Semaphore
go run 15_rate_limiter.go # Rate Limiter
```

---

## 2. ETL Pipeline з Error Handling

**Паттерни**: Pipeline + Transform + Filter + ErrGroup

```go
// Extract → Transform → Filter → Load (with error handling)
```

**Use case**: Data processing з обробкою помилок

```bash
go run 06_pipeline.go  # Pipeline
go run 04_transform.go # Transform
go run 05_filter.go    # Filter
go run 13_errgroup.go  # ErrGroup
```

---

## 3. Cache Warming

**Паттерни**: SingleFlight + Future + Rate Limiter

```go
// Multiple requests → SingleFlight → Future → Cache
```

**Use case**: Prevent thundering herd при старті

```bash
go run 14_singleflight.go # SingleFlight
go run 10_future.go       # Future
go run 15_rate_limiter.go # Rate Limiter
```

---

## 4. Distributed Task Processing

**Паттерни**: Fan-In + Fan-Out + Barrier

```go
// Tasks → Fan-Out → Workers → Barrier → Fan-In → Results
```

**Use case**: Map-Reduce операції

```bash
go run 01_fan_in.go  # Fan-In
go run 02_fan_out.go # Fan-Out
go run 08_barrier.go # Barrier
```

---

## 5. Real-time Data Stream

**Паттерни**: Generator + Tee + Transform + Filter

```go
// Generator → Tee → [Console, File, Network] → Transform → Filter
```

**Use case**: Monitoring system

```bash
go run 12_generator.go # Generator
go run 03_tee.go       # Tee
go run 04_transform.go # Transform
go run 05_filter.go    # Filter
```

---

## 6. Async API Client

**Паттерни**: Promise + Future + ErrGroup + SingleFlight

```go
// Requests → Promise → Future → ErrGroup (parallel) → Results
```

**Use case**: HTTP client library

```bash
go run 09_promise.go       # Promise
go run 11_future_promise.go # Future+Promise
go run 13_errgroup.go      # ErrGroup
go run 14_singleflight.go  # SingleFlight
```

---

## Patterns Matrix

| Pattern | Combines Well With | Use Case |
|---------|-------------------|----------|
| Fan-In | Fan-Out, Pipeline | Aggregation |
| Fan-Out | Semaphore, Rate Limiter | Worker pools |
| Pipeline | Transform, Filter, ErrGroup | ETL |
| Semaphore | Fan-Out, Rate Limiter | Resource control |
| ErrGroup | Pipeline, Fan-Out | Error handling |
| SingleFlight | Future, Rate Limiter | Deduplication |
| Rate Limiter | Fan-Out, Semaphore | Throttling |
| Generator | Tee, Transform, Filter | Data streams |
| Promise/Future | ErrGroup, SingleFlight | Async ops |

---

## Best Practices

### 1. Завжди закривайте канали
```go
defer close(ch)
```

### 2. Використовуйте context для cancellation
```go
select {
case <-ctx.Done():
    return ctx.Err()
case result := <-ch:
    // process
}
```

### 3. Обмежуйте concurrency
```go
// Замість необмеженої кількості горутин
sem := make(chan struct{}, maxConcurrent)
```

### 4. Обробляйте помилки
```go
// Використовуйте ErrGroup замість WaitGroup коли потрібна обробка помилок
```

### 5. Уникайте deadlock
```go
// Завжди перевіряйте чи канал може бути прочитаний
```

---

## Production Checklist

- [ ] Context cancellation
- [ ] Error handling (ErrGroup)
- [ ] Resource limits (Semaphore)
- [ ] Rate limiting
- [ ] Graceful shutdown
- [ ] Monitoring/metrics
- [ ] Testing with race detector
- [ ] Deadlock prevention

---

## Testing

```bash
# Run with race detector
go run -race 01_fan_in.go
go run -race 02_fan_out.go
go run -race 13_errgroup.go

# Benchmark
go test -bench=. -benchmem
```

---

**Комбінуйте паттерни для production-ready рішень!** 🚀
