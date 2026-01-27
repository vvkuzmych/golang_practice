# Exercise 3: Мікросервіси + База Даних

## 🎯 Мета

Створити простий e-commerce backend з мікросервісами та PostgreSQL.

---

## 📝 Завдання

### Архітектура

```
API Gateway (:8080)
    ├── Product Service (:8001) → PostgreSQL (products DB)
    ├── Order Service (:8002) → PostgreSQL (orders DB)
    └── User Service (:8003) → PostgreSQL (users DB)
```

### Частина 1: Database Schema

**Product Service:**
```sql
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);
```

**Order Service:**
```sql
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL,
    total DECIMAL(10,2) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW()
);
```

**User Service:**
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Частина 2: Product Service

Endpoints:
- `GET /products` - список продуктів
- `GET /products/:id` - деталі продукту
- `POST /products` - створити продукт
- `PUT /products/:id/stock` - оновити кількість

Використайте GORM!

### Частина 3: Order Service

Endpoints:
- `POST /orders` - створити замовлення
  - Перевірити наявність stock (виклик Product Service)
  - Перевірити user (виклик User Service)
  - Зменшити stock
  - Створити order

- `GET /orders/:id` - деталі замовлення
- `GET /orders/user/:user_id` - замовлення користувача

### Частина 4: API Gateway

Роутинг:
- `/api/products/*` → Product Service
- `/api/orders/*` → Order Service
- `/api/users/*` → User Service

Aggregation endpoint:
- `GET /api/orders/:id/full` - order + product details + user details

### Частина 5: Service Discovery

Реалізуйте простий Service Registry:
```go
type ServiceRegistry struct {
    services map[string]string // "product-service" -> "http://localhost:8001"
}
```

---

## ✅ Критерії успіху

- [ ] 3 сервіси працюють незалежно
- [ ] API Gateway роутить запити
- [ ] GORM використовується для всіх сервісів
- [ ] Міжсервісна комунікація через HTTP
- [ ] Транзакції при створенні замовлення
- [ ] Graceful shutdown для всіх сервісів

---

## 🧪 Тестування

```bash
# Створити продукт
curl -X POST http://localhost:8080/api/products \
  -d '{"name":"Laptop","price":999.99,"stock":10}'

# Створити користувача
curl -X POST http://localhost:8080/api/users \
  -d '{"email":"user@example.com","name":"John Doe"}'

# Створити замовлення
curl -X POST http://localhost:8080/api/orders \
  -d '{"user_id":1,"product_id":1,"quantity":2}'

# Отримати full order details
curl http://localhost:8080/api/orders/1/full
```

---

## 🚀 Бонус

Додайте:
- Circuit Breaker для міжсервісних викликів
- Retry logic з exponential backoff
- Distributed tracing (IDs для requests)
- Health check endpoints для кожного сервісу
- Docker Compose для запуску всього стека

---

## 💡 Підказки

1. Кожен сервіс має свою БД - не використовуйте shared database
2. Використовуйте `context.WithTimeout` для HTTP викликів
3. Логуйте всі міжсервісні виклики
4. Обробляйте випадки, коли сервіс недоступний

---

**Час виконання:** 6-8 годин

**Це фінальна вправа Week 6!** 🎉
