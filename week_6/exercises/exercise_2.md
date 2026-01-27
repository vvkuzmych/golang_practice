# Exercise 2: REST API Server

## 🎯 Мета

Створити RESTful API для управління списком завдань (TODO list).

---

## 📝 Завдання

### Частина 1: Модель даних

Створіть структуру `Todo`:
```go
type Todo struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Completed   bool      `json:"completed"`
    CreatedAt   time.Time `json:"created_at"`
}
```

### Частина 2: HTTP Endpoints

Реалізуйте наступні endpoints:

| Method | Path | Опис |
|--------|------|------|
| GET | `/api/todos` | Отримати всі завдання |
| GET | `/api/todos/:id` | Отримати завдання по ID |
| POST | `/api/todos` | Створити нове завдання |
| PUT | `/api/todos/:id` | Оновити завдання |
| DELETE | `/api/todos/:id` | Видалити завдання |

### Частина 3: Middleware

Реалізуйте middleware для:
1. **Logging** - логувати всі запити
2. **CORS** - дозволити cross-origin requests
3. **Authentication** - перевірка API key в header

### Частина 4: Error Handling

Централізована обробка помилок:
```go
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
    Code    int    `json:"code"`
}
```

### Частина 5: Валідація

Перевіряйте:
- `Title` не пустий
- `Title` не довший за 100 символів
- `Description` не довший за 500 символів

---

## ✅ Критерії успіху

- [ ] Всі endpoints працюють
- [ ] JSON response для всіх endpoints
- [ ] Middleware застосовано коректно
- [ ] Помилки обробляються централізовано
- [ ] Валідація працює

---

## 🧪 Тестування

```bash
# Створити TODO
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Buy groceries","description":"Milk, bread, eggs"}'

# Отримати всі TODOs
curl http://localhost:8080/api/todos

# Оновити TODO
curl -X PUT http://localhost:8080/api/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"completed":true}'

# Видалити TODO
curl -X DELETE http://localhost:8080/api/todos/1
```

---

## 🚀 Бонус

Додайте:
- Фільтрація по `completed` status
- Пагінація (`?page=1&limit=10`)
- Сортування по `created_at`
- Rate limiting middleware

---

**Час виконання:** 3-4 години
