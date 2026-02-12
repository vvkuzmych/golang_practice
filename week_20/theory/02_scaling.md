# Scaling Strategies

## Vertical vs Horizontal Scaling

### Vertical Scaling (Scale Up) ⬆️
**Збільшити потужність одного сервера.**

```
Before: 4 CPU, 8GB RAM
After:  16 CPU, 64GB RAM
```

**Переваги:**
- ✅ Простіше реалізувати
- ✅ Немає складності розподіленої системи
- ✅ Consistency легше підтримувати

**Недоліки:**
- ❌ Є ліміт (максимальна конфігурація сервера)
- ❌ Single point of failure
- ❌ Дорого (high-end сервери)
- ❌ Downtime при upgrade

**Коли використовувати:**
- Монолітні застосунки
- До ~10K користувачів
- Прототипи, MVP

### Horizontal Scaling (Scale Out) ➡️
**Додати більше серверів.**

```
Before: 1 server
After:  10 servers + Load Balancer
```

**Переваги:**
- ✅ Немає верхньої межі
- ✅ High availability (один сервер падає - інші працюють)
- ✅ Дешевше (commodity hardware)
- ✅ Zero downtime deployments

**Недоліки:**
- ❌ Складніша архітектура
- ❌ Data consistency challenges
- ❌ Distributed tracing потрібен

**Коли використовувати:**
- Мікросервіси
- Від ~10K користувачів
- Production-grade системи

---

## Load Balancing

**Розподіл трафіку між серверами.**

```
       User Requests
            ↓
      Load Balancer
     /      |      \
Server 1  Server 2  Server 3
```

### Алгоритми

#### 1. Round Robin
```
Request 1 → Server 1
Request 2 → Server 2
Request 3 → Server 3
Request 4 → Server 1 (по колу)
```

#### 2. Least Connections
```
Server 1: 10 connections
Server 2: 5 connections  ← вибираємо
Server 3: 8 connections
```

#### 3. IP Hash
```
IP: 192.168.1.100 → hash → Server 2 (завжди той самий сервер для IP)
```

#### 4. Weighted Round Robin
```
Server 1 (weight 3): 60% requests
Server 2 (weight 2): 40% requests
```

### Приклад (Nginx)

```nginx
upstream backend {
    least_conn;  # algorithm
    
    server backend1.example.com weight=3;
    server backend2.example.com weight=2;
    server backend3.example.com backup;  # fallback
}

server {
    location / {
        proxy_pass http://backend;
    }
}
```

---

## Caching Strategies

### 1. Client-Side Caching
```
Browser Cache (HTTP headers)
Cache-Control: max-age=3600
ETag: "abc123"
```

### 2. CDN Caching
```
User (Ukraine) → CDN (Poland) → Origin (USA)
Static assets: images, CSS, JS
```

### 3. Application Caching (Redis, Memcached)

```go
// Cache aside pattern
func GetUser(id int) (*User, error) {
    // 1. Спробувати з кешу
    cached, err := redis.Get(fmt.Sprintf("user:%d", id))
    if err == nil {
        return unmarshal(cached), nil
    }
    
    // 2. Якщо немає - з БД
    user, err := db.Query("SELECT * FROM users WHERE id = ?", id)
    if err != nil {
        return nil, err
    }
    
    // 3. Записати в кеш
    redis.Set(fmt.Sprintf("user:%d", id), marshal(user), 1*time.Hour)
    
    return user, nil
}
```

### Cache Invalidation Strategies

#### 1. Time-based (TTL)
```
SET user:123 "Alice" EX 3600  # expires in 1 hour
```

#### 2. Event-based
```go
// When user updated
func UpdateUser(user *User) error {
    db.Update(user)
    redis.Del(fmt.Sprintf("user:%d", user.ID))  // invalidate cache
}
```

#### 3. Write-through
```go
// Write to cache and DB at the same time
func UpdateUser(user *User) error {
    redis.Set(fmt.Sprintf("user:%d", user.ID), marshal(user))
    return db.Update(user)
}
```

---

## Database Scaling

### 1. Read Replicas
```
          Primary (Write)
         /       |       \
Replica 1   Replica 2   Replica 3
(Read)      (Read)      (Read)
```

**Коли:**
- Read-heavy workload (90% reads, 10% writes)
- Analytics queries

### 2. Sharding (Horizontal Partitioning)
```
Shard 1: users 1-1000000
Shard 2: users 1000001-2000000
Shard 3: users 2000001-3000000
```

**Sharding strategies:**

#### Hash-based
```go
shardID := hash(userID) % totalShards
```

#### Range-based
```
Shard 1: A-M
Shard 2: N-Z
```

#### Geography-based
```
Shard 1: Europe
Shard 2: USA
Shard 3: Asia
```

**Challenges:**
- Cross-shard queries складні
- Rebalancing shards важко
- Joins між shards неможливі

### 3. Vertical Partitioning
```
Table: users (id, name, email, bio, photo, settings...)

Split into:
users_basic (id, name, email)
users_profile (id, bio, photo)
users_settings (id, settings)
```

---

## Microservices Scaling

```
API Gateway
    ↓
┌─────────────┬─────────────┬─────────────┐
│   User      │   Order     │  Payment    │
│  Service    │  Service    │  Service    │
│ (3 instances)│ (5 instances)│ (2 instances)│
└─────────────┴─────────────┴─────────────┘
```

**Переваги:**
- Кожен сервіс масштабується окремо
- Технологічна незалежність
- Fault isolation

**Недоліки:**
- Network latency
- Distributed tracing потрібен
- Data consistency складніша

---

## Auto-Scaling

### Horizontal Pod Autoscaler (Kubernetes)

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: myapp-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: myapp
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

**Тригери:**
- CPU usage > 70% → scale up
- Memory usage > 80% → scale up
- Request rate > 1000/s → scale up
- Custom metrics (queue length, response time)

---

## Bottleneck Analysis

### 1. CPU-bound
**Symptom:** High CPU usage  
**Solution:** 
- Horizontal scaling
- Code optimization
- Caching

### 2. Memory-bound
**Symptom:** High memory usage, OOM kills  
**Solution:** 
- Vertical scaling (more RAM)
- Memory profiling & optimization
- Pagination

### 3. I/O-bound
**Symptom:** Slow disk/network  
**Solution:** 
- SSD instead of HDD
- CDN for static assets
- Connection pooling

### 4. Database-bound
**Symptom:** Slow queries  
**Solution:** 
- Indexes
- Query optimization
- Read replicas
- Caching
- Sharding

---

## Real-world Example: Twitter

```
Users: 330M active
Tweets: 500M per day
Reads: 5B per day (10,000:1 read/write ratio)

Architecture:
- Load Balancers (HAProxy)
- Application Servers (Ruby on Rails → Scala)
- Caching (Redis, Memcached)
- Databases:
  - MySQL (sharded by user_id)
  - Manhattan (distributed key-value)
- Message Queues (Kafka)
- CDN (Akamai)
- Timeline generation (Fan-out on write)
```

---

## Підсумок

| Strategy | Pros | Cons | Use Case |
|----------|------|------|----------|
| **Vertical Scaling** | Simple | Limited, expensive | MVP, small apps |
| **Horizontal Scaling** | Unlimited, HA | Complex | Production systems |
| **Load Balancing** | HA, performance | SPOF (LB itself) | All apps |
| **Caching** | Fast reads | Stale data | Read-heavy |
| **Sharding** | Unlimited DB scale | Complex queries | Huge datasets |
| **Microservices** | Independent scaling | Network overhead | Large teams |

**Золоте правило:**
1. Почни з vertical scaling
2. Додай caching
3. Перейди на horizontal scaling
4. Додай read replicas
5. Шардуй БД (якщо потрібно)
