# 🌍 Real-World Use Cases: Goroutines & Channels

## Коли і де використовуються Goroutines та Channels в Production

Цей документ показує **20 реальних сценаріїв** використання goroutines та channels у виробничих застосунках.

---

## 📊 Категорії Use Cases

1. [HTTP Servers & APIs](#1-http-servers--apis) (5 cases)
2. [Data Processing & Pipelines](#2-data-processing--pipelines) (4 cases)
3. [Background Jobs & Workers](#3-background-jobs--workers) (3 cases)
4. [Real-time Systems](#4-real-time-systems) (3 cases)
5. [Infrastructure & DevOps](#5-infrastructure--devops) (3 cases)
6. [Distributed Systems](#6-distributed-systems) (2 cases)

---

## 1. HTTP Servers & APIs

### 🔹 Case 1: HTTP Server з Graceful Shutdown

**Проблема:** HTTP сервер повинен коректно завершитись при Ctrl+C, дочекавшись завершення всіх поточних requests.

**Рішення:**
```go
func main() {
    srv := &http.Server{Addr: ":8080"}
    
    // Goroutine для HTTP server
    go func() {
        if err := srv.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()
    
    // Signal handling
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan
    
    // Graceful shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    srv.Shutdown(ctx)
}
```

**Використання:**
- ✅ Всі web servers (Gin, Echo, Fiber, net/http)
- ✅ Microservices
- ✅ REST APIs

---

### 🔹 Case 2: Паралельні HTTP Requests (Fan-Out)

**Проблема:** Потрібно зробити 10 HTTP requests паралельно і зібрати результати.

**Рішення:**
```go
type Result struct {
    URL      string
    Response *http.Response
    Error    error
}

func fetchURLs(urls []string) []Result {
    results := make(chan Result, len(urls))
    var wg sync.WaitGroup
    
    for _, url := range urls {
        wg.Add(1)
        go func(u string) {
            defer wg.Done()
            resp, err := http.Get(u)
            results <- Result{URL: u, Response: resp, Error: err}
        }(url)
    }
    
    go func() {
        wg.Wait()
        close(results)
    }()
    
    var collected []Result
    for r := range results {
        collected = append(collected, r)
    }
    return collected
}
```

**Використання:**
- ✅ Aggregation APIs (збирають дані з кількох джерел)
- ✅ Health check systems (перевірка кількох endpoints)
- ✅ Web scrapers
- ✅ API gateways

**Реальні проекти:**
- Kubernetes health checks
- Prometheus scraping targets
- GraphQL resolvers (parallel field resolution)

---

### 🔹 Case 3: Rate Limiting для API

**Проблема:** Обмежити кількість requests до зовнішнього API (наприклад, 100 req/sec).

**Рішення:**
```go
type RateLimiter struct {
    ticker   *time.Ticker
    requests chan struct{}
}

func NewRateLimiter(rps int) *RateLimiter {
    rl := &RateLimiter{
        ticker:   time.NewTicker(time.Second / time.Duration(rps)),
        requests: make(chan struct{}, rps),
    }
    
    go func() {
        for range rl.ticker.C {
            select {
            case rl.requests <- struct{}{}:
            default:
            }
        }
    }()
    
    return rl
}

func (rl *RateLimiter) Wait() {
    <-rl.requests
}

// Використання:
limiter := NewRateLimiter(100)
for _, item := range items {
    limiter.Wait()
    go processItem(item)
}
```

**Використання:**
- ✅ API clients (Twitter, GitHub, Stripe APIs)
- ✅ Database connection pools
- ✅ External service integrations

**Реальні проекти:**
- GitHub API client (5000 req/hour limit)
- AWS SDK (service-specific limits)
- Redis rate limiters

---

### 🔹 Case 4: Request Timeout з Context

**Проблема:** HTTP request не повинен виконуватись довше 5 секунд.

**Рішення:**
```go
func makeRequest(url string) (*http.Response, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    return http.DefaultClient.Do(req)
}
```

**Використання:**
- ✅ HTTP clients з timeout
- ✅ Database queries з timeout
- ✅ gRPC calls з deadline
- ✅ Microservice communication

---

### 🔹 Case 5: WebSocket Broadcasting

**Проблема:** Відправити повідомлення всім підключеним WebSocket clients.

**Рішення:**
```go
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}

func (h *Hub) run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
            
        case client := <-h.unregister:
            delete(h.clients, client)
            close(client.send)
            
        case message := <-h.broadcast:
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
        }
    }
}
```

**Використання:**
- ✅ Chat applications (Slack, Discord)
- ✅ Real-time dashboards (Grafana, Datadog)
- ✅ Live notifications
- ✅ Multiplayer games

**Реальні проекти:**
- Gorilla WebSocket
- Centrifugo (real-time messaging)

---

## 2. Data Processing & Pipelines

### 🔹 Case 6: ETL Pipeline (Extract-Transform-Load)

**Проблема:** Обробити мільйони records: читати з DB → transform → записати в іншу DB.

**Рішення:**
```go
func ETLPipeline(ctx context.Context) {
    // Stage 1: Extract (read from DB)
    records := extract(ctx, db)
    
    // Stage 2: Transform (process data)
    transformed := transform(ctx, records)
    
    // Stage 3: Load (write to DB)
    load(ctx, transformedDB, transformed)
}

func extract(ctx context.Context, db *sql.DB) <-chan Record {
    out := make(chan Record)
    go func() {
        defer close(out)
        rows, _ := db.QueryContext(ctx, "SELECT * FROM source")
        defer rows.Close()
        
        for rows.Next() {
            var r Record
            rows.Scan(&r.ID, &r.Data)
            select {
            case out <- r:
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}

func transform(ctx context.Context, in <-chan Record) <-chan Record {
    out := make(chan Record)
    const numWorkers = 10
    var wg sync.WaitGroup
    
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for record := range in {
                // Transform logic
                record.Data = process(record.Data)
                select {
                case out <- record:
                case <-ctx.Done():
                    return
                }
            }
        }()
    }
    
    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}
```

**Використання:**
- ✅ Data warehouses (Snowflake, BigQuery)
- ✅ Data migration tools
- ✅ Log aggregation systems
- ✅ Analytics pipelines

**Реальні проекти:**
- Apache Kafka consumers
- Airflow tasks (написані на Go)
- Custom ETL tools

---

### 🔹 Case 7: Image Processing Pipeline

**Проблема:** Обробити 10,000 images: resize → compress → upload to S3.

**Рішення:**
```go
func processImages(images []string) {
    const numWorkers = 20
    jobs := make(chan string, len(images))
    results := make(chan Result, len(images))
    
    // Workers
    for w := 0; w < numWorkers; w++ {
        go func() {
            for imagePath := range jobs {
                // Resize
                img := resize(imagePath)
                // Compress
                compressed := compress(img)
                // Upload
                url := uploadToS3(compressed)
                results <- Result{Path: imagePath, URL: url}
            }
        }()
    }
    
    // Send jobs
    for _, img := range images {
        jobs <- img
    }
    close(jobs)
    
    // Collect results
    for i := 0; i < len(images); i++ {
        <-results
    }
}
```

**Використання:**
- ✅ Image CDNs (Cloudinary, Imgix)
- ✅ Video processing (thumbnail generation)
- ✅ PDF generation services
- ✅ Document conversion tools

---

### 🔹 Case 8: Log Processing з Buffer

**Проблема:** Зібрати 1000 log messages в batch перед відправкою в Elasticsearch.

**Рішення:**
```go
type LogBatcher struct {
    logs   chan LogEntry
    batch  []LogEntry
    size   int
    ticker *time.Ticker
}

func (lb *LogBatcher) Start() {
    go func() {
        for {
            select {
            case log := <-lb.logs:
                lb.batch = append(lb.batch, log)
                if len(lb.batch) >= lb.size {
                    lb.flush()
                }
            case <-lb.ticker.C:
                if len(lb.batch) > 0 {
                    lb.flush()
                }
            }
        }
    }()
}

func (lb *LogBatcher) flush() {
    sendToElasticsearch(lb.batch)
    lb.batch = lb.batch[:0]
}
```

**Використання:**
- ✅ Logging systems (Logrus, Zap)
- ✅ Metrics aggregation (Prometheus)
- ✅ Event streaming (Kafka producers)
- ✅ Analytics tracking

**Реальні проекти:**
- Fluentd/Fluent Bit (log forwarders)
- Beats (Elastic stack)
- Vector (observability data pipeline)

---

### 🔹 Case 9: CSV File Processing

**Проблема:** Обробити великий CSV файл (100GB) паралельно.

**Рішення:**
```go
func processCSV(filename string) {
    file, _ := os.Open(filename)
    defer file.Close()
    
    scanner := bufio.NewScanner(file)
    lines := make(chan string, 100)
    results := make(chan Result, 100)
    
    // Readers
    go func() {
        for scanner.Scan() {
            lines <- scanner.Text()
        }
        close(lines)
    }()
    
    // Workers
    var wg sync.WaitGroup
    for w := 0; w < 10; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for line := range lines {
                result := processLine(line)
                results <- result
            }
        }()
    }
    
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Writer
    for result := range results {
        writeToDatabase(result)
    }
}
```

**Використання:**
- ✅ Data import tools
- ✅ Batch processing systems
- ✅ Financial reporting
- ✅ Data analysis tools

---

## 3. Background Jobs & Workers

### 🔹 Case 10: Job Queue з Priority

**Проблема:** Обробляти jobs з різними пріоритетами (high, medium, low).

**Рішення:**
```go
type JobQueue struct {
    high   chan Job
    medium chan Job
    low    chan Job
}

func (jq *JobQueue) worker(id int) {
    for {
        select {
        case job := <-jq.high:
            fmt.Printf("Worker %d: processing HIGH priority job\n", id)
            process(job)
        case job := <-jq.medium:
            fmt.Printf("Worker %d: processing MEDIUM priority job\n", id)
            process(job)
        case job := <-jq.low:
            fmt.Printf("Worker %d: processing LOW priority job\n", id)
            process(job)
        }
    }
}
```

**Використання:**
- ✅ Background job processors (Sidekiq equivalent in Go)
- ✅ Task schedulers
- ✅ Email sending queues
- ✅ Notification systems

**Реальні проекти:**
- Asynq (distributed task queue)
- Machinery (async task queue)
- RiverQueue

---

### 🔹 Case 11: Scheduled Tasks (Cron Jobs)

**Проблема:** Виконувати tasks кожні 5 хвилин (cleanup, reports, backups).

**Рішення:**
```go
type Scheduler struct {
    tasks map[string]*ScheduledTask
    stop  chan bool
}

type ScheduledTask struct {
    Name     string
    Interval time.Duration
    Task     func()
}

func (s *Scheduler) Start() {
    for _, task := range s.tasks {
        go func(t *ScheduledTask) {
            ticker := time.NewTicker(t.Interval)
            defer ticker.Stop()
            
            for {
                select {
                case <-ticker.C:
                    t.Task()
                case <-s.stop:
                    return
                }
            }
        }(task)
    }
}

// Використання:
scheduler.AddTask("cleanup", 5*time.Minute, cleanupOldFiles)
scheduler.AddTask("backup", 1*time.Hour, backupDatabase)
scheduler.Start()
```

**Використання:**
- ✅ Database cleanup tasks
- ✅ Report generation
- ✅ Cache invalidation
- ✅ Health checks

**Реальні проекти:**
- Kubernetes CronJobs
- Nomad periodic jobs
- Custom schedulers

---

### 🔹 Case 12: Email Sending Queue

**Проблема:** Відправити 10,000 emails асинхронно з retry логікою.

**Рішення:**
```go
type EmailQueue struct {
    emails  chan Email
    results chan EmailResult
    workers int
}

func (eq *EmailQueue) Start() {
    for w := 0; w < eq.workers; w++ {
        go func(id int) {
            for email := range eq.emails {
                err := sendEmail(email)
                if err != nil && shouldRetry(err) {
                    // Retry logic
                    time.Sleep(1 * time.Second)
                    eq.emails <- email // Re-queue
                } else {
                    eq.results <- EmailResult{
                        Email: email,
                        Error: err,
                    }
                }
            }
        }(w)
    }
}
```

**Використання:**
- ✅ Transactional emails (order confirmations)
- ✅ Marketing campaigns
- ✅ Password reset emails
- ✅ Notification systems

**Реальні проекти:**
- Mailgun queue
- SendGrid batch processing
- AWS SES with SQS

---

## 4. Real-time Systems

### 🔹 Case 13: Server-Sent Events (SSE) Broadcasting

**Проблема:** Відправляти real-time updates клієнтам (live scores, stock prices).

**Рішення:**
```go
type SSEServer struct {
    clients  map[chan string]bool
    addCh    chan chan string
    removeCh chan chan string
    broadcast chan string
}

func (s *SSEServer) Run() {
    go func() {
        for {
            select {
            case client := <-s.addCh:
                s.clients[client] = true
            case client := <-s.removeCh:
                delete(s.clients, client)
                close(client)
            case msg := <-s.broadcast:
                for client := range s.clients {
                    select {
                    case client <- msg:
                    default:
                        // Client slow/disconnected
                        delete(s.clients, client)
                        close(client)
                    }
                }
            }
        }
    }()
}
```

**Використання:**
- ✅ Live sports scores
- ✅ Stock tickers
- ✅ Real-time dashboards
- ✅ Live auction systems

---

### 🔹 Case 14: Metrics Collection System

**Проблема:** Збирати metrics з 1000 servers кожні 10 секунд.

**Рішення:**
```go
type MetricsCollector struct {
    servers []Server
    metrics chan Metric
}

func (mc *MetricsCollector) Start() {
    ticker := time.NewTicker(10 * time.Second)
    
    go func() {
        for range ticker.C {
            var wg sync.WaitGroup
            for _, server := range mc.servers {
                wg.Add(1)
                go func(s Server) {
                    defer wg.Done()
                    metric := s.CollectMetrics()
                    mc.metrics <- metric
                }(server)
            }
        }
    }()
    
    // Processor
    go func() {
        for metric := range mc.metrics {
            storeInDB(metric)
            checkAlerts(metric)
        }
    }()
}
```

**Використання:**
- ✅ Monitoring systems (Prometheus, Datadog)
- ✅ APM tools (New Relic, AppDynamics)
- ✅ Infrastructure monitoring
- ✅ Application metrics

**Реальні проекти:**
- Prometheus exporters
- Telegraf plugins
- Custom monitoring agents

---

### 🔹 Case 15: Event Sourcing System

**Проблема:** Обробляти events асинхронно і апдейтити read models.

**Рішення:**
```go
type EventStore struct {
    events     chan Event
    projectors []Projector
}

func (es *EventStore) Start() {
    // Fan-out: один event → кілька projectors
    for _, projector := range es.projectors {
        go func(p Projector) {
            for event := range es.events {
                p.Project(event)
            }
        }(projector)
    }
}

// Використання:
eventStore := NewEventStore()
eventStore.AddProjector(&UserProjector{})
eventStore.AddProjector(&OrderProjector{})
eventStore.AddProjector(&AnalyticsProjector{})
eventStore.Start()

// Publish event
eventStore.Publish(OrderCreatedEvent{OrderID: 123})
```

**Використання:**
- ✅ CQRS systems
- ✅ Event-driven architectures
- ✅ Audit logging
- ✅ Microservices communication

---

## 5. Infrastructure & DevOps

### 🔹 Case 16: Health Check System

**Проблема:** Перевіряти health 50 microservices кожні 30 секунд.

**Рішення:**
```go
type HealthChecker struct {
    services []Service
    results  chan HealthResult
}

func (hc *HealthChecker) Start() {
    ticker := time.NewTicker(30 * time.Second)
    
    go func() {
        for range ticker.C {
            var wg sync.WaitGroup
            for _, service := range hc.services {
                wg.Add(1)
                go func(s Service) {
                    defer wg.Done()
                    
                    ctx, cancel := context.WithTimeout(
                        context.Background(), 
                        5*time.Second,
                    )
                    defer cancel()
                    
                    healthy := s.CheckHealth(ctx)
                    hc.results <- HealthResult{
                        Service: s.Name,
                        Healthy: healthy,
                        Time:    time.Now(),
                    }
                }(service)
            }
        }
    }()
}
```

**Використання:**
- ✅ Kubernetes liveness/readiness probes
- ✅ Load balancer health checks
- ✅ Service mesh (Istio, Linkerd)
- ✅ Monitoring dashboards

---

### 🔹 Case 17: Distributed Cache Warming

**Проблема:** Прогріти cache на 20 servers паралельно після deploy.

**Рішення:**
```go
func warmupCache(servers []string, data []CacheEntry) {
    var wg sync.WaitGroup
    
    for _, server := range servers {
        wg.Add(1)
        go func(s string) {
            defer wg.Done()
            
            client := redis.NewClient(&redis.Options{Addr: s})
            for _, entry := range data {
                client.Set(entry.Key, entry.Value, entry.TTL)
            }
        }(server)
    }
    
    wg.Wait()
    log.Println("Cache warmed up on all servers")
}
```

**Використання:**
- ✅ CDN cache warming
- ✅ Redis/Memcached warming
- ✅ Application cache initialization
- ✅ Deployment automation

---

### 🔹 Case 18: Log Aggregation System

**Проблема:** Збирати logs з 100 pods в Kubernetes cluster.

**Рішення:**
```go
type LogAggregator struct {
    pods    []Pod
    logs    chan LogEntry
    storage LogStorage
}

func (la *LogAggregator) Start() {
    // Tail logs from each pod
    for _, pod := range la.pods {
        go func(p Pod) {
            stream := p.StreamLogs()
            for log := range stream {
                la.logs <- LogEntry{
                    Pod:       p.Name,
                    Timestamp: time.Now(),
                    Message:   log,
                }
            }
        }(pod)
    }
    
    // Batch writer
    go func() {
        batch := make([]LogEntry, 0, 1000)
        ticker := time.NewTicker(5 * time.Second)
        
        for {
            select {
            case log := <-la.logs:
                batch = append(batch, log)
                if len(batch) >= 1000 {
                    la.storage.WriteBatch(batch)
                    batch = batch[:0]
                }
            case <-ticker.C:
                if len(batch) > 0 {
                    la.storage.WriteBatch(batch)
                    batch = batch[:0]
                }
            }
        }
    }()
}
```

**Використання:**
- ✅ ELK stack (Elasticsearch, Logstash, Kibana)
- ✅ Grafana Loki
- ✅ Splunk
- ✅ Custom log aggregators

---

## 6. Distributed Systems

### 🔹 Case 19: Distributed Task Execution

**Проблема:** Виконати task на 10 remote workers і зібрати results.

**Рішення:**
```go
type DistributedExecutor struct {
    workers []WorkerClient
}

func (de *DistributedExecutor) Execute(task Task) []Result {
    results := make(chan Result, len(de.workers))
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    var wg sync.WaitGroup
    for _, worker := range de.workers {
        wg.Add(1)
        go func(w WorkerClient) {
            defer wg.Done()
            result, err := w.ExecuteTask(ctx, task)
            if err != nil {
                results <- Result{Error: err, Worker: w.ID}
            } else {
                results <- Result{Data: result, Worker: w.ID}
            }
        }(worker)
    }
    
    go func() {
        wg.Wait()
        close(results)
    }()
    
    var collected []Result
    for r := range results {
        collected = append(collected, r)
    }
    return collected
}
```

**Використання:**
- ✅ MapReduce frameworks
- ✅ Distributed computing
- ✅ Grid computing
- ✅ Parallel test execution

**Реальні проекти:**
- Apache Spark (Go clients)
- Dask (distributed computing)
- Custom distributed systems

---

### 🔹 Case 20: Message Queue Consumer (Kafka/RabbitMQ)

**Проблема:** Споживати messages з Kafka topic паралельно.

**Рішення:**
```go
type KafkaConsumer struct {
    consumer    *kafka.Consumer
    workers     int
    messages    chan *kafka.Message
    stopCh      chan bool
}

func (kc *KafkaConsumer) Start() {
    // Fetcher: читає messages з Kafka
    go func() {
        for {
            select {
            case <-kc.stopCh:
                return
            default:
                msg, err := kc.consumer.ReadMessage(1 * time.Second)
                if err == nil {
                    kc.messages <- msg
                }
            }
        }
    }()
    
    // Workers: обробляють messages паралельно
    for w := 0; w < kc.workers; w++ {
        go func(id int) {
            for msg := range kc.messages {
                processMessage(msg)
                kc.consumer.CommitMessages(msg)
            }
        }(w)
    }
}

func (kc *KafkaConsumer) Stop() {
    close(kc.stopCh)
    close(kc.messages)
}
```

**Використання:**
- ✅ Event streaming platforms
- ✅ Message queue consumers
- ✅ Data pipelines
- ✅ Microservices communication

**Реальні проекти:**
- Kafka Go clients (Sarama, Confluent)
- RabbitMQ Go library (amqp091-go)
- NATS streaming
- Redis Streams consumers

---

## 📊 Статистика використання

### За типом застосунку:

| Тип | Use Cases | % |
|-----|-----------|---|
| **Web/API** | 5 | 25% |
| **Data Processing** | 4 | 20% |
| **Background Jobs** | 3 | 15% |
| **Real-time** | 3 | 15% |
| **Infrastructure** | 3 | 15% |
| **Distributed** | 2 | 10% |

### Топ-5 паттернів:

1. **Worker Pool** (Cases: 2, 6, 7, 9, 10, 12, 20) - **35%**
2. **Fan-Out/Fan-In** (Cases: 2, 14, 15, 19) - **20%**
3. **Pipeline** (Cases: 6, 7, 8, 9) - **20%**
4. **Broadcasting** (Cases: 5, 13, 15) - **15%**
5. **Graceful Shutdown** (Cases: 1, 11, 12, 20) - **20%**

---

## ✅ Коли використовувати Goroutines & Channels

### ✅ Використовуйте коли:

1. **I/O-bound operations** - HTTP requests, DB queries, file operations
2. **Parallel processing** - обробка великих datasets
3. **Real-time systems** - WebSockets, SSE, streaming
4. **Background tasks** - email sending, cleanup, reports
5. **Event-driven architecture** - event sourcing, pub/sub
6. **Distributed systems** - microservices, distributed computing

### ❌ НЕ використовуйте коли:

1. **CPU-bound без parallelism** - складні математичні обчислення (використайте GOMAXPROCS)
2. **Прості синхронні операції** - не потрібна concurrency
3. **Малий обсяг даних** - overhead може бути більшим за benefit
4. **Складна state synchronization** - можливо простіше single-threaded

---

## 📚 Корисні ресурси

### Open Source проекти з гарними прикладами:
- **Docker** - container orchestration
- **Kubernetes** - cluster management
- **Prometheus** - monitoring system
- **InfluxDB** - time-series database
- **CockroachDB** - distributed database
- **NATS** - messaging system
- **Caddy** - web server
- **Traefik** - reverse proxy

### Книги:
- "Concurrency in Go" by Katherine Cox-Buday
- "Go in Action" by William Kennedy
- "The Go Programming Language" by Alan Donovan

---

## 🎯 Висновок

**Goroutines та Channels - це фундамент Go для:**

1. 🚀 **Concurrency** - паралельна обробка
2. 📡 **Communication** - координація між goroutines
3. ⚡ **Performance** - ефективне використання ресурсів
4. 🎨 **Simplicity** - простий concurrent code

**Головне правило:**
> "Don't communicate by sharing memory; share memory by communicating."

**Практикуйтесь на реальних проектах!** 🎉

---

**Створено:** 2026-01-15  
**Week 5:** Goroutines & Channels
