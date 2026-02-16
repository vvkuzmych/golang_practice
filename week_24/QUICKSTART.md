# Week 24 - Quick Start Guide 🚀

## Швидкий старт

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_24

# Запустити будь-який приклад
go run 01_fan_in.go
```

---

## Всі паттерни

### 1. Fan-In (01_fan_in.go)
**Що робить**: Об'єднує кілька каналів в один.

```bash
go run 01_fan_in.go
```

**Коли використовувати**:
- Збір логів з кількох сервісів
- Агрегація метрик
- Об'єднання результатів

---

### 2. Fan-Out (02_fan_out.go)
**Що робить**: Розподіляє роботу між кількома workers.

```bash
go run 02_fan_out.go
```

**Коли використовувати**:
- Паралельна обробка великих датасетів
- Worker pools
- Batch операції

---

### 3. Tee (03_tee.go)
**Що робить**: Дублює дані в кілька каналів.

```bash
go run 03_tee.go
```

**Коли використовувати**:
- Логування (консоль + файл + мережа)
- Broadcasting подій
- Audit trails

---

### 4. Transform (04_transform.go)
**Що робить**: Перетворює дані в каналі.

```bash
go run 04_transform.go
```

**Коли використовувати**:
- ETL pipelines
- Нормалізація даних
- Конвертація форматів

---

### 5. Filter (05_filter.go)
**Що робить**: Фільтрує дані за умовою.

```bash
go run 05_filter.go
```

**Коли використовувати**:
- Валідація даних
- Фільтрація подій
- Stream processing

---

### 6. Pipeline (06_pipeline.go)
**Що робить**: Послідовна обробка даних через кілька стадій.

```bash
go run 06_pipeline.go
```

**Коли використовувати**:
- ETL (Extract, Transform, Load)
- Image/Video processing
- Request transformation

---

### 7. Semaphore (07_semaphore.go)
**Що робить**: Обмежує кількість одночасних операцій.

```bash
go run 07_semaphore.go
```

**Коли використовувати**:
- Connection pools
- Rate limiting
- Resource management

---

### 8. Barrier (08_barrier.go)
**Що робить**: Синхронізація: всі чекають одне одного.

```bash
go run 08_barrier.go
```

**Коли використовувати**:
- Multi-phase computations
- Distributed algorithms
- Game loop sync

---

### 9. Promise (09_promise.go)
**Що робить**: Контейнер для майбутнього результату.

```bash
go run 09_promise.go
```

**Коли використовувати**:
- Async API calls
- JavaScript-like pattern
- Future-based APIs

---

### 10. Future (10_future.go)
**Що робить**: Асинхронне обчислення з отриманням результату пізніше.

```bash
go run 10_future.go
```

**Коли використовувати**:
- Lazy evaluation
- Background tasks
- Parallel computations

---

### 11. Future + Promise (11_future_promise.go)
**Що робить**: Комбінація Future і Promise.

```bash
go run 11_future_promise.go
```

**Коли використовувати**:
- HTTP client libraries
- RPC calls
- Promise.all() pattern

---

### 12. Generator (12_generator.go)
**Що робить**: Нескінченний потік даних.

```bash
go run 12_generator.go
```

**Коли використовувати**:
- Infinite sequences
- Event streams
- ID generators

---

### 13. ErrGroup (13_errgroup.go)
**Що робить**: Групова обробка з автоматичним cancel при помилці.

```bash
go run 13_errgroup.go
```

**Коли використовувати**:
- Parallel API calls (stop on error)
- Database migrations
- Health checks

---

### 14. SingleFlight (14_singleflight.go)
**Що робить**: Дедуплікація одночасних запитів.

```bash
go run 14_singleflight.go
```

**Коли використовувати**:
- Cache warming (thundering herd)
- Дедуплікація DB queries
- Preventing duplicate work

---

### 15. Rate Limiter (15_rate_limiter.go)
**Що робить**: Обмеження швидкості виконання.

```bash
go run 15_rate_limiter.go
```

**Коли використовувати**:
- API rate limiting
- Database throttling
- DoS prevention

---

## Must-Know Patterns ⭐⭐⭐

```bash
# Топ 5 найважливіших
go run 01_fan_in.go         # Fan-In
go run 02_fan_out.go        # Fan-Out
go run 06_pipeline.go       # Pipeline
go run 13_errgroup.go       # ErrGroup
go run 15_rate_limiter.go   # Rate Limiter
```

---

## Тестування

```bash
# Перевірка компіляції всіх файлів
for f in *.go; do go build -o /dev/null "$f" && echo "✅ $f"; done
```

---

## Production Usage

Ці паттерни активно використовуються в:
- **Microservices** (fan-out, errgroup)
- **API Gateway** (rate limiter, singleflight)
- **Data Processing** (pipeline, transform, filter)
- **Worker Pools** (fan-in, fan-out, semaphore)
- **Caching** (future, promise, singleflight)

Успішного навчання! 🎉
