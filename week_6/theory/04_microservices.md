# Мікросервісна Архітектура

Мікросервіси - це архітектурний стиль, де додаток складається з маленьких, незалежних сервісів.

---

## 📖 Зміст

1. [Монолітна vs Мікросервісна архітектура](#1-монолітна-vs-мікросервісна-архітектура)
2. [Principles of Microservices](#2-principles-of-microservices)
3. [API Gateway](#3-api-gateway)
4. [Service Discovery](#4-service-discovery)
5. [Inter-service Communication](#5-inter-service-communication)
6. [Data Management](#6-data-management)

---

## 1. Монолітна vs Мікросервісна архітектура

### Монолітна архітектура

```
┌─────────────────────────────────────┐
│        Monolithic Application       │
│                                     │
│  ┌────────────┐  ┌──────────────┐ │
│  │    UI      │  │   Business   │ │
│  │  Layer     │  │    Logic     │ │
│  └────────────┘  └──────────────┘ │
│                                     │
│  ┌──────────────────────────────┐  │
│  │      Data Access Layer        │  │
│  └──────────────────────────────┘  │
│                 │                   │
└─────────────────┼───────────────────┘
                  │
           ┌──────▼──────┐
           │   Database  │
           └─────────────┘
```

**Переваги:**
- ✅ Простота розробки
- ✅ Простота тестування
- ✅ Простота deployment
- ✅ Швидкий старт

**Недоліки:**
- ❌ Складно масштабувати
- ❌ Технологічна залежність
- ❌ Довгі релізні цикли
- ❌ Складність підтримки при зростанні

### Мікросервісна архітектура

```
┌─────────────────────────────────────────────────────┐
│                   API Gateway                        │
└────┬─────────────┬─────────────┬────────────┬───────┘
     │             │             │            │
┌────▼────┐  ┌────▼────┐  ┌────▼────┐  ┌───▼─────┐
│  User   │  │ Product │  │ Order   │  │ Payment │
│ Service │  │ Service │  │ Service │  │ Service │
└────┬────┘  └────┬────┘  └────┬────┘  └───┬─────┘
     │            │             │            │
┌────▼────┐  ┌───▼─────┐  ┌───▼─────┐  ┌───▼─────┐
│ User DB │  │Product  │  │ Order   │  │ Payment │
└─────────┘  │   DB    │  │   DB    │  │   DB    │
             └─────────┘  └─────────┘  └─────────┘
```

**Переваги:**
- ✅ Незалежне масштабування
- ✅ Технологічна свобода
- ✅ Ізоляція помилок
- ✅ Швидкі релізи
- ✅ Легше підтримувати

**Недоліки:**
- ❌ Складна інфраструктура
- ❌ Розподілена система (складність)
- ❌ Міжсервісна комунікація
- ❌ Тестування end-to-end
- ❌ Моніторинг та відладка

---

## 2. Principles of Microservices

### Single Responsibility

Кожен сервіс відповідає за одну бізнес-функцію.

```
✅ Good:
- User Service: управління користувачами
- Auth Service: автентифікація та авторизація
- Email Service: відправка email
- Notification Service: push notifications

❌ Bad:
- User Service: користувачі + auth + email + notifications
```

### Domain-Driven Design (DDD)

Сервіси організовані навколо бізнес-доменів.

```
E-commerce Domains:
┌──────────────────────────────────────────┐
│         Catalog (Domain)                 │
│  - Product Service                       │
│  - Category Service                      │
│  - Search Service                        │
└──────────────────────────────────────────┘

┌──────────────────────────────────────────┐
│         Ordering (Domain)                │
│  - Cart Service                          │
│  - Order Service                         │
│  - Shipping Service                      │
└──────────────────────────────────────────┘

┌──────────────────────────────────────────┐
│         Payment (Domain)                 │
│  - Payment Service                       │
│  - Billing Service                       │
└──────────────────────────────────────────┘
```

### Decentralized Data Management

Кожен сервіс має свою власну базу даних.

```go
// ❌ Shared Database (Anti-pattern)
// User Service → Shared DB ← Order Service

// ✅ Database per Service
// User Service → User DB
// Order Service → Order DB
```

### Communication Patterns

**Синхронна комунікація (HTTP/gRPC):**
```
User Service → [HTTP GET /products/123] → Product Service
             ← [JSON: {id:123, name:"Laptop"}] ←
```

**Асинхронна комунікація (Message Queue):**
```
Order Service → [Publish: OrderCreated] → Message Queue
                                              ↓
Email Service ← [Subscribe: OrderCreated] ←─┘
```

---

## 3. API Gateway

### Що таке API Gateway?

API Gateway - це єдина точка входу для всіх клієнтів.

```
┌──────────┐
│  Mobile  │
│   App    │
└─────┬────┘
      │
┌─────▼────┐        ┌─────────────────────────────┐
│   Web    │        │       API Gateway           │
│ Browser  │───────▶│                             │
└──────────┘        │ - Routing                   │
                    │ - Authentication            │
┌──────────┐        │ - Rate Limiting             │
│  IoT     │        │ - Load Balancing            │
│ Device   │───────▶│ - Request/Response          │
└──────────┘        │   Transformation            │
                    └───┬─────┬──────┬────────┬───┘
                        │     │      │        │
                   ┌────▼┐ ┌─▼───┐ ┌▼────┐ ┌─▼──────┐
                   │User │ │Auth │ │Order│ │Product │
                   │Svc  │ │Svc  │ │Svc  │ │  Svc   │
                   └─────┘ └─────┘ └─────┘ └────────┘
```

### Приклад API Gateway в Go

```go
package main

import (
    "encoding/json"
    "net/http"
    "net/http/httputil"
    "net/url"
)

type Gateway struct {
    userServiceURL    string
    productServiceURL string
    orderServiceURL   string
}

func NewGateway() *Gateway {
    return &Gateway{
        userServiceURL:    "http://localhost:8001",
        productServiceURL: "http://localhost:8002",
        orderServiceURL:   "http://localhost:8003",
    }
}

// Proxy до User Service
func (g *Gateway) UsersHandler(w http.ResponseWriter, r *http.Request) {
    target, _ := url.Parse(g.userServiceURL)
    proxy := httputil.NewSingleHostReverseProxy(target)
    proxy.ServeHTTP(w, r)
}

// Proxy до Product Service
func (g *Gateway) ProductsHandler(w http.ResponseWriter, r *http.Request) {
    target, _ := url.Parse(g.productServiceURL)
    proxy := httputil.NewSingleHostReverseProxy(target)
    proxy.ServeHTTP(w, r)
}

// Aggregation - об'єднання даних з кількох сервісів
func (g *Gateway) OrderDetailsHandler(w http.ResponseWriter, r *http.Request) {
    orderID := r.URL.Query().Get("id")
    
    // 1. Отримуємо order
    order := g.fetchOrder(orderID)
    
    // 2. Отримуємо user
    user := g.fetchUser(order.UserID)
    
    // 3. Отримуємо products
    products := g.fetchProducts(order.ProductIDs)
    
    // 4. Об'єднуємо
    response := map[string]interface{}{
        "order":    order,
        "user":     user,
        "products": products,
    }
    
    json.NewEncoder(w).Encode(response)
}

func main() {
    gateway := NewGateway()
    
    http.HandleFunc("/api/users/", gateway.UsersHandler)
    http.HandleFunc("/api/products/", gateway.ProductsHandler)
    http.HandleFunc("/api/orders/details", gateway.OrderDetailsHandler)
    
    http.ListenAndServe(":8080", nil)
}
```

---

## 4. Service Discovery

### Проблема

Сервіси динамічно змінюють адреси (IP/Port) при масштабуванні.

```
Order Service хоче викликати User Service
Але де він знаходиться? 192.168.1.5:8001? 192.168.1.6:8001?
```

### Рішення: Service Registry

```
┌─────────────────────────────────────┐
│       Service Registry              │
│  (Consul, Eureka, etcd, Zookeeper)  │
│                                     │
│  user-service:   192.168.1.5:8001  │
│  order-service:  192.168.1.6:8002  │
│  product-service: 192.168.1.7:8003 │
└───────▲──────────────┬──────────────┘
        │              │
   Registration    Discovery
        │              │
┌───────┴────────┐ ┌──▼────────────┐
│  User Service  │ │ Order Service │
└────────────────┘ └───────────────┘
```

### Client-Side Discovery

```go
package main

import (
    "fmt"
    "net/http"
)

type ServiceRegistry interface {
    Register(serviceName, address string) error
    Discover(serviceName string) (string, error)
}

// Простий in-memory registry
type InMemoryRegistry struct {
    services map[string][]string
    index    map[string]int
}

func NewRegistry() *InMemoryRegistry {
    return &InMemoryRegistry{
        services: make(map[string][]string),
        index:    make(map[string]int),
    }
}

func (r *InMemoryRegistry) Register(serviceName, address string) error {
    r.services[serviceName] = append(r.services[serviceName], address)
    return nil
}

// Round-robin load balancing
func (r *InMemoryRegistry) Discover(serviceName string) (string, error) {
    instances := r.services[serviceName]
    if len(instances) == 0 {
        return "", fmt.Errorf("no instances found for %s", serviceName)
    }
    
    idx := r.index[serviceName]
    address := instances[idx%len(instances)]
    r.index[serviceName]++
    
    return address, nil
}

// Використання
func main() {
    registry := NewRegistry()
    
    // Реєструємо сервіси
    registry.Register("user-service", "http://localhost:8001")
    registry.Register("user-service", "http://localhost:8002") // масштабування
    registry.Register("product-service", "http://localhost:8003")
    
    // Викликаємо сервіс
    address, _ := registry.Discover("user-service")
    resp, _ := http.Get(address + "/api/users/123")
    defer resp.Body.Close()
}
```

---

## 5. Inter-service Communication

### HTTP/REST

**Переваги:**
- ✅ Просто і зрозуміло
- ✅ Широка підтримка
- ✅ Debugging-friendly

**Недоліки:**
- ❌ Повільніше за gRPC
- ❌ Більше трафіку (JSON)

```go
// Order Service викликає User Service
func getUser(userID string) (*User, error) {
    resp, err := http.Get(fmt.Sprintf("http://user-service:8001/users/%s", userID))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var user User
    json.NewDecoder(resp.Body).Decode(&user)
    return &user, nil
}
```

### gRPC

**Переваги:**
- ✅ Швидкий (binary protocol)
- ✅ Strongly-typed (Protocol Buffers)
- ✅ Streaming підтримка

**Недоліки:**
- ❌ Складніше налаштувати
- ❌ Debugging складніший

```protobuf
// user.proto
syntax = "proto3";

service UserService {
    rpc GetUser(UserRequest) returns (UserResponse);
}

message UserRequest {
    string user_id = 1;
}

message UserResponse {
    string id = 1;
    string name = 2;
    string email = 3;
}
```

```go
// Order Service викликає User Service через gRPC
func getUser(userID string) (*UserResponse, error) {
    conn, _ := grpc.Dial("user-service:50051", grpc.WithInsecure())
    defer conn.Close()
    
    client := NewUserServiceClient(conn)
    user, err := client.GetUser(context.Background(), &UserRequest{
        UserId: userID,
    })
    return user, err
}
```

### Message Queue (Async)

**Переваги:**
- ✅ Decoupling
- ✅ Resilience
- ✅ Load leveling

**Недоліки:**
- ❌ Eventual consistency
- ❌ Складність

```go
// Order Service публікує подію
func createOrder(order Order) error {
    // Зберігаємо order
    saveOrder(order)
    
    // Публікуємо подію
    publishEvent("order.created", order)
    return nil
}

// Email Service підписується на подію
func subscribeToOrderEvents() {
    subscribe("order.created", func(order Order) {
        sendOrderConfirmationEmail(order)
    })
}
```

---

## 6. Data Management

### Database per Service

Кожен сервіс має свою БД - це ключовий принцип мікросервісів.

```
User Service → PostgreSQL (Users, Profiles)
Order Service → PostgreSQL (Orders, OrderItems)
Product Service → MongoDB (Products, Categories)
Analytics Service → Elasticsearch (Logs, Metrics)
```

### Saga Pattern

Для розподілених транзакцій використовуємо **Saga Pattern**.

**Приклад: Створення замовлення**

```
1. Order Service: CreateOrder
   → Success: Publish OrderCreated
   → Fail: End
   
2. Payment Service: ProcessPayment
   → Success: Publish PaymentProcessed
   → Fail: Publish PaymentFailed → Order Service: CancelOrder
   
3. Inventory Service: ReserveItems
   → Success: Publish ItemsReserved
   → Fail: Publish ReservationFailed → Payment Service: RefundPayment
                                     → Order Service: CancelOrder
4. Shipping Service: CreateShipment
   → Success: Order Completed!
   → Fail: Compensate all previous steps
```

```go
// Order Service
func CreateOrder(order Order) error {
    // 1. Створюємо order
    order.Status = "PENDING"
    db.Save(&order)
    
    // 2. Публікуємо подію
    publishEvent("order.created", order)
    
    // 3. Чекаємо на відповіді від інших сервісів
    return nil
}

// Compensation при помилці
func HandlePaymentFailed(event PaymentFailedEvent) {
    order := getOrder(event.OrderID)
    order.Status = "CANCELLED"
    db.Save(&order)
}
```

---

## ✅ Best Practices

1. **Start with Monolith** - не починайте з мікросервісів якщо не потрібно
2. **Domain-Driven Design** - організуйте сервіси навколо бізнес-доменів
3. **API Gateway** - єдина точка входу
4. **Service Discovery** - динамічне виявлення сервісів
5. **Circuit Breaker** - захист від каскадних помилок
6. **Distributed Tracing** - відстеження запитів (Jaeger, Zipkin)
7. **Centralized Logging** - централізовані логи (ELK, Loki)
8. **Health Checks** - моніторинг здоров'я сервісів
9. **API Versioning** - версіонування API

---

**Далі:** [05_databases.md](./05_databases.md)
