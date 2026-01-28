# Proxy Pattern (Проксі / Замісник)

## 📋 Опис

**Proxy** - структурний патерн, що надає об'єкт-замінник для контролю доступу до іншого об'єкта.

## 🎯 Проблема

- Потрібно контролювати доступ до об'єкта
- Об'єкт "дорогий" у створенні (lazy loading)
- Потрібна додаткова функціональність (кешування, логування, перевірка прав)

## ✅ Рішення

Створити proxy-об'єкт з тим самим інтерфейсом, що і реальний об'єкт.

## 🏗️ Структура

```
Client → Proxy → RealSubject
         │
         └─> Додаткова логіка:
             - Кешування
             - Lazy loading
             - Access control
             - Logging
```

## 💻 Типи Proxy

### 1. Virtual Proxy (Віртуальний)
- Lazy initialization
- Створює об'єкт тільки коли потрібно
- **Use case:** Завантаження великих зображень

### 2. Protection Proxy (Захисний)
- Контроль прав доступу
- Перевіряє права перед викликом
- **Use case:** Доступ до sensitive data

### 3. Caching Proxy (Кешуючий)
- Кешує результати
- Повертає з кешу якщо можливо
- **Use case:** API requests, database queries

### 4. Remote Proxy (Віддалений)
- Представляє об'єкт в іншому address space
- **Use case:** RPC, gRPC clients

### 5. Smart Proxy (Розумний)
- Додаткова функціональність
- Logging, metrics, retry logic
- **Use case:** Production services

## ✅ Переваги

- Контроль доступу без зміни реального об'єкта
- Lazy initialization
- Додаткова функціональність (кеш, логи)
- Прозорість для клієнта

## ❌ Недоліки

- Додатковий рівень абстракції
- Може уповільнити код (overhead)
- Складніше тестувати

## 🎯 Коли використовувати

✅ Lazy initialization (дорогі об'єкти)  
✅ Access control (права доступу)  
✅ Caching (оптимізація)  
✅ Logging/Monitoring  
✅ Remote objects (RPC)  

## 📊 Real-World приклади

### Go Standard Library

```go
// net/http: ReverseProxy
proxy := &httputil.ReverseProxy{
    Director: func(req *http.Request) {
        req.URL.Host = "backend:8080"
    },
}

// database/sql: Connection pooling
db.Query() // Proxy до реального connection
```

### Production Use Cases

1. **API Gateway** - Proxy до мікросервісів
2. **CDN** - Caching Proxy для статики
3. **Database Connection Pool** - Proxy до DB connections
4. **gRPC Interceptors** - Logging/Auth Proxy

## 🔄 Порівняння з іншими патернами

### Proxy vs Decorator
- **Proxy:** Контроль доступу
- **Decorator:** Додавання функціональності

### Proxy vs Adapter
- **Proxy:** Той самий інтерфейс
- **Adapter:** Різні інтерфейси

### Proxy vs Facade
- **Proxy:** 1-to-1 (один об'єкт)
- **Facade:** 1-to-many (багато об'єктів)

## 🚀 Запуск

```bash
cd design_patterns/structural/proxy
go run main.go
```

## 📖 Output

```
=== 1. Caching Proxy ===

First request:
Proxy: Creating RealSubject (lazy init)
Creating RealSubject: API Service
Proxy: Logging request: data1
RealSubject handling request: data1
Proxy: Request completed
Got: Result from API Service: data1

Second request (same data):
Proxy: Returning cached result for: data1
Got: Result from API Service: data1

=== 2. Protection Proxy ===

Admin request:
Creating RealSubject: Secure Service
RealSubject handling request: sensitive data
Result from Secure Service: sensitive data

Regular user request:
ProtectionProxy: Access denied! Admin role required.

=== 3. Virtual Proxy (Lazy Loading) ===
Proxies created (images not loaded yet)

Displaying image1 for the first time:
ImageProxy: Loading image for the first time...
Loading image from disk: photo1.jpg
Displaying image: photo1.jpg

Displaying image1 again:
Displaying image: photo1.jpg
```
