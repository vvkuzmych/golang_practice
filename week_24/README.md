# Week 24 - Advanced Channel Patterns 🚀

Професійні паттерни роботи з каналами в Go.

## 📁 Структура

```
week_24/
├── 01_fan_in.go              Fan-In: об'єднання каналів
├── 02_fan_out.go             Fan-Out: розподіл роботи
├── 03_tee.go                 Tee: дублювання даних
├── 04_transform.go           Трансформація даних
├── 05_filter.go              Фільтрація даних
├── 06_pipeline.go            Pipeline: ланцюг обробки
├── 07_semaphore.go           Семафор: обмеження ресурсів
├── 08_barrier.go             Бар'єр: синхронізація
├── 09_promise.go             Promise: відкладений результат
├── 10_future.go              Future: async обчислення
├── 11_future_promise.go      Future + Promise разом
├── 12_generator.go           Генератор: нескінченний потік
├── 13_errgroup.go            ErrGroup: обробка помилок
├── 14_singleflight.go        SingleFlight: дедуплікація
└── 15_rate_limiter.go        Rate Limiter: throttling
```

---

## 🚀 Запуск

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_24

# Запустити будь-який приклад
go run 01_fan_in.go
go run 06_pipeline.go
go run 13_errgroup.go
```

---

## 📚 Теми

### Basic Patterns
1. **Fan-In** - Об'єднання кількох каналів в один
2. **Fan-Out** - Розподіл роботи між workers
3. **Tee** - Дублювання даних в кілька каналів

### Data Processing
4. **Transform** - Перетворення даних в каналі
5. **Filter** - Фільтрація даних за умовою
6. **Pipeline** - Послідовна обробка даних

### Synchronization
7. **Semaphore** - Обмеження кількості ресурсів
8. **Barrier** - Синхронізація горутин

### Async Patterns
9. **Promise** - Відкладений результат
10. **Future** - Асинхронне обчислення
11. **Future + Promise** - Комбінований паттерн

### Advanced
12. **Generator** - Нескінченний генератор даних
13. **ErrGroup** - Групова обробка з помилками
14. **SingleFlight** - Дедуплікація запитів
15. **Rate Limiter** - Обмеження швидкості

---

## 🎯 Найважливіші

```bash
# Must-know паттерни:
go run 01_fan_in.go          # ⭐⭐⭐
go run 02_fan_out.go         # ⭐⭐⭐
go run 06_pipeline.go        # ⭐⭐⭐
go run 13_errgroup.go        # ⭐⭐⭐
go run 15_rate_limiter.go    # ⭐⭐⭐
```

---

## 💡 Real-World Usage

Ці паттерни використовуються в:
- Microservices (fan-out, errgroup)
- API Gateway (rate limiter, singleflight)
- Data Processing (pipeline, transform, filter)
- Worker Pools (fan-in, fan-out, semaphore)
- Caching (future, promise, singleflight)

---

**15 Production-Ready Channel Patterns!** 🎉
