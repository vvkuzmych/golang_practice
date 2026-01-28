# Scalable Backend Services

## 📖 Зміст

1. [Horizontal vs Vertical Scaling](#horizontal-vs-vertical-scaling)
2. [Load Balancing](#load-balancing)
3. [Caching Strategies](#caching-strategies)
4. [Message Queues](#message-queues)
5. [Database Scaling](#database-scaling)

---

## Horizontal vs Vertical Scaling

### Vertical Scaling (Scale Up)
- Збільшення ресурсів одного сервера (CPU, RAM)
- **Pros:** Простіше, без змін в коді
- **Cons:** Обмежений потенціал, single point of failure

### Horizontal Scaling (Scale Out)
- Додавання більше серверів
- **Pros:** Необмежене масштабування, high availability
- **Cons:** Складніша архітектура, потрібен load balancer

```go
// Stateless design для horizontal scaling
type UserService struct {
    db    *sql.DB
    cache *redis.Client
    // НЕ зберігаємо state в пам'яті!
}

func (s *UserService) GetUser(ctx context.Context, id int) (*User, error) {
    // Check cache first
    if user, err := s.getUserFromCache(ctx, id); err == nil {
        return user, nil
    }
    
    // Fallback to DB
    user, err := s.getUserFromDB(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Update cache
    s.cacheUser(ctx, user)
    
    return user, nil
}
```

---

## Load Balancing

### Nginx Config

```nginx
upstream backend {
    # Round-robin (default)
    server backend1:8080;
    server backend2:8080;
    server backend3:8080;
    
    # Health check
    keepalive 32;
}

server {
    listen 80;
    
    location / {
        proxy_pass http://backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### Load Balancing Algorithms

1. **Round Robin** - по черзі
2. **Least Connections** - найменше активних з'єднань
3. **IP Hash** - за IP клієнта (sticky sessions)
4. **Weighted** - з вагами (performance-based)

---

## Caching Strategies

### 1. Cache-Aside (Lazy Loading)

```go
func (s *Service) GetUser(ctx context.Context, id int) (*User, error) {
    // 1. Check cache
    key := fmt.Sprintf("user:%d", id)
    if data, err := s.cache.Get(ctx, key).Bytes(); err == nil {
        var user User
        json.Unmarshal(data, &user)
        return &user, nil
    }
    
    // 2. Load from DB
    user, err := s.db.GetUser(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // 3. Store in cache
    data, _ := json.Marshal(user)
    s.cache.Set(ctx, key, data, 1*time.Hour)
    
    return user, nil
}
```

### 2. Write-Through Cache

```go
func (s *Service) UpdateUser(ctx context.Context, user *User) error {
    // 1. Update DB
    if err := s.db.UpdateUser(ctx, user); err != nil {
        return err
    }
    
    // 2. Update cache immediately
    key := fmt.Sprintf("user:%d", user.ID)
    data, _ := json.Marshal(user)
    s.cache.Set(ctx, key, data, 1*time.Hour)
    
    return nil
}
```

### 3. Write-Behind (Write-Back) Cache

```go
func (s *Service) UpdateUser(ctx context.Context, user *User) error {
    // 1. Update cache first
    key := fmt.Sprintf("user:%d", user.ID)
    data, _ := json.Marshal(user)
    s.cache.Set(ctx, key, data, 1*time.Hour)
    
    // 2. Queue DB write (async)
    s.queue.Publish("user_updates", user)
    
    return nil
}
```

### Redis Example

```go
import "github.com/go-redis/redis/v8"

// Setup Redis
rdb := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
})

// Cache user
func cacheUser(ctx context.Context, user *User) error {
    data, err := json.Marshal(user)
    if err != nil {
        return err
    }
    
    key := fmt.Sprintf("user:%d", user.ID)
    return rdb.Set(ctx, key, data, 1*time.Hour).Err()
}

// Get cached user
func getCachedUser(ctx context.Context, id int) (*User, error) {
    key := fmt.Sprintf("user:%d", id)
    data, err := rdb.Get(ctx, key).Bytes()
    if err == redis.Nil {
        return nil, ErrNotFound
    }
    if err != nil {
        return nil, err
    }
    
    var user User
    if err := json.Unmarshal(data, &user); err != nil {
        return nil, err
    }
    
    return &user, nil
}

// Invalidate cache
func invalidateUser(ctx context.Context, id int) error {
    key := fmt.Sprintf("user:%d", id)
    return rdb.Del(ctx, key).Err()
}
```

---

## Message Queues

### RabbitMQ

```go
import "github.com/streadway/amqp"

// Producer
func publishMessage(ch *amqp.Channel, queueName, message string) error {
    return ch.Publish(
        "",        // exchange
        queueName, // routing key
        false,     // mandatory
        false,     // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:        []byte(message),
        },
    )
}

// Consumer
func consumeMessages(ch *amqp.Channel, queueName string) error {
    msgs, err := ch.Consume(
        queueName, // queue
        "",        // consumer
        false,     // auto-ack
        false,     // exclusive
        false,     // no-local
        false,     // no-wait
        nil,       // args
    )
    if err != nil {
        return err
    }
    
    for msg := range msgs {
        processMessage(msg.Body)
        msg.Ack(false) // Acknowledge
    }
    
    return nil
}
```

### AWS SQS

```go
import "github.com/aws/aws-sdk-go/service/sqs"

// Send message
func sendSQSMessage(svc *sqs.SQS, queueURL, message string) error {
    _, err := svc.SendMessage(&sqs.SendMessageInput{
        QueueUrl:    aws.String(queueURL),
        MessageBody: aws.String(message),
    })
    return err
}

// Receive messages
func receiveSQSMessages(svc *sqs.SQS, queueURL string) error {
    result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
        QueueUrl:            aws.String(queueURL),
        MaxNumberOfMessages: aws.Int64(10),
        WaitTimeSeconds:     aws.Int64(20), // Long polling
    })
    if err != nil {
        return err
    }
    
    for _, msg := range result.Messages {
        processMessage(*msg.Body)
        
        // Delete after processing
        svc.DeleteMessage(&sqs.DeleteMessageInput{
            QueueUrl:      aws.String(queueURL),
            ReceiptHandle: msg.ReceiptHandle,
        })
    }
    
    return nil
}
```

---

## Database Scaling

### Read Replicas

```go
type DB struct {
    master *sql.DB
    slaves []*sql.DB
}

func (db *DB) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
    // Read from random slave
    slave := db.slaves[rand.Intn(len(db.slaves))]
    return slave.QueryContext(ctx, query, args...)
}

func (db *DB) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
    // Write to master
    return db.master.ExecContext(ctx, query, args...)
}
```

### Sharding

```go
func getShardDB(userID int) *sql.DB {
    shardID := userID % numShards
    return shards[shardID]
}

func getUserByID(userID int) (*User, error) {
    db := getShardDB(userID)
    return db.QueryUser(userID)
}
```

### Connection Pooling

```go
db, err := sql.Open("postgres", connStr)

// Configure pool
db.SetMaxOpenConns(25)               // Maximum open connections
db.SetMaxIdleConns(5)                // Maximum idle connections
db.SetConnMaxLifetime(5 * time.Minute) // Connection lifetime
```

---

## Scalability Patterns

### 1. CQRS (Command Query Responsibility Segregation)

```go
// Write model
type CommandHandler struct {
    db *sql.DB
}

func (h *CommandHandler) CreateUser(user *User) error {
    // Write to master DB
    return h.db.Create(user)
}

// Read model
type QueryHandler struct {
    cache *redis.Client
    db    *sql.DB
}

func (h *QueryHandler) GetUser(id int) (*User, error) {
    // Try cache first, then DB
    return h.getUserOptimized(id)
}
```

### 2. Event Sourcing

```go
type Event struct {
    Type      string
    Aggregate string
    Data      json.RawMessage
    Timestamp time.Time
}

func (s *EventStore) Append(event Event) error {
    // Append to event log
    return s.db.Insert(event)
}

func (s *EventStore) GetEvents(aggregateID string) ([]Event, error) {
    // Replay events
    return s.db.Query("SELECT * FROM events WHERE aggregate = ?", aggregateID)
}
```

### 3. Circuit Breaker

```go
type CircuitBreaker struct {
    maxFailures int
    timeout     time.Duration
    failures    int
    lastFail    time.Time
    state       string
    mu          sync.Mutex
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    if cb.isOpen() {
        return errors.New("circuit breaker open")
    }
    
    err := fn()
    if err != nil {
        cb.recordFailure()
        return err
    }
    
    cb.recordSuccess()
    return nil
}
```

---

## Best Practices

✅ **Design for horizontal scaling**
✅ **Use caching strategically**
✅ **Implement health checks**
✅ **Monitor performance metrics**
✅ **Use connection pooling**
✅ **Implement rate limiting**
✅ **Plan for failure (circuit breaker)**
✅ **Use message queues for async tasks**
✅ **Implement graceful shutdown**
✅ **Log everything (structured logging)**

---

## Monitoring

```go
import "github.com/prometheus/client_golang/prometheus"

var (
    requestCount = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration",
        },
        []string{"method", "endpoint"},
    )
)

func middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        next.ServeHTTP(w, r)
        
        duration := time.Since(start).Seconds()
        requestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
        requestCount.WithLabelValues(r.Method, r.URL.Path, "200").Inc()
    })
}
```

---

**Scalability - це не про код, а про архітектуру!** 🚀
