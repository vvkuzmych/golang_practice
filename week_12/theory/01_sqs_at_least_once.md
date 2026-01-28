# Amazon SQS & At-Least-Once Delivery

## 🎯 Що таке Amazon SQS?

**Amazon SQS (Simple Queue Service)** - це fully managed message queuing service для decoupling та scaling distributed systems.

```
Producer → [SQS Queue] → Consumer
```

**Key benefits:**
- ✅ Fully managed (no servers)
- ✅ Auto-scaling
- ✅ Durable (replicated across AZs)
- ✅ Pay per request

---

## 📊 SQS Queue Types

### Standard Queue

```
Throughput:     Unlimited
Ordering:       Best-effort (not guaranteed)
Delivery:       At-least-once (may duplicate)
Latency:        < 10ms
```

**Use case:** High throughput, order не критичний

### FIFO Queue

```
Throughput:     300 msgs/sec (3000 with batching)
Ordering:       Guaranteed (strict FIFO)
Delivery:       Exactly-once (no duplicates)
Latency:        ~20ms
Name suffix:    .fifo
```

**Use case:** Order критичний, no duplicates

---

## 🔄 At-Least-Once Delivery

### Що таке At-Least-Once?

**At-least-once delivery** означає, що message буде доставлено **принаймні один раз**, але може бути доставлено **кілька разів**.

```
Send message → SQS stores → Consumer receives
                         → Consumer may receive again!
```

### Чому відбуваються дублікати?

**Scenario 1: Consumer не встиг підтвердити**
```
1. Consumer receives message
2. Processing starts
3. Network glitch
4. SQS не отримав ACK
5. SQS re-delivers message
```

**Scenario 2: Visibility Timeout закінчився**
```
1. Consumer receives message (visibility timeout = 30s)
2. Processing takes 45s (too long!)
3. Message стає visible знову
4. Another consumer receives same message
```

---

## ⏱️ Visibility Timeout

### Як працює?

```
Message in queue (visible)
    ↓
Consumer receives message
    ↓
Message hidden (invisible to others)
    ↓
Visibility timeout running... (30s)
    ↓
Option 1: Consumer deletes → Message gone ✅
Option 2: Timeout expires → Message visible again ⚠️
```

### Configuration

```go
// Receive message with 30s visibility timeout
result, err := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
    QueueUrl:            aws.String(queueURL),
    MaxNumberOfMessages: 10,
    WaitTimeSeconds:     20,
    VisibilityTimeout:   30,  // 30 seconds
})
```

### Extend Visibility Timeout

```go
// If processing takes longer, extend timeout
_, err = sqsClient.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{
    QueueUrl:          aws.String(queueURL),
    ReceiptHandle:     message.ReceiptHandle,
    VisibilityTimeout: 60,  // Extend to 60s
})
```

---

## 🎯 Idempotency (Ідемпотентність)

### Проблема: Дублікати

```
Message 1: CreateOrder {orderId: 123, amount: 100}
Message 2: CreateOrder {orderId: 123, amount: 100}  ← Duplicate!

Without idempotency: 2 orders created ❌
With idempotency:    1 order created ✅
```

### Рішення 1: Idempotency Key

```go
func processOrder(ctx context.Context, msg OrderMessage) error {
    // Check if already processed
    exists, err := redis.Exists(ctx, msg.IdempotencyKey)
    if err != nil {
        return err
    }
    
    if exists {
        log.Printf("Message already processed: %s", msg.IdempotencyKey)
        return nil  // Skip duplicate
    }
    
    // Process order
    err = createOrder(msg.OrderID, msg.Amount)
    if err != nil {
        return err
    }
    
    // Mark as processed (TTL 24h)
    redis.Set(ctx, msg.IdempotencyKey, "processed", 24*time.Hour)
    
    return nil
}
```

### Рішення 2: Database Constraint

```sql
-- Unique constraint на idempotency key
CREATE TABLE orders (
    id UUID PRIMARY KEY,
    idempotency_key VARCHAR(255) UNIQUE NOT NULL,
    amount DECIMAL(10, 2),
    created_at TIMESTAMP DEFAULT NOW()
);
```

```go
func createOrder(idempotencyKey string, amount float64) error {
    // INSERT will fail if idempotency_key exists
    _, err := db.Exec(
        "INSERT INTO orders (id, idempotency_key, amount) VALUES ($1, $2, $3)",
        uuid.New(), idempotencyKey, amount,
    )
    
    if err != nil {
        if isDuplicateKeyError(err) {
            log.Printf("Duplicate order: %s", idempotencyKey)
            return nil  // Already processed
        }
        return err
    }
    
    return nil
}
```

---

## 📊 Message Lifecycle

### Full Flow

```
1. Producer sends message
   ↓
2. SQS stores (replicated 3x across AZs)
   ↓
3. Message visible in queue
   ↓
4. Consumer polls (long polling)
   ↓
5. Message received → Invisible (visibility timeout)
   ↓
6a. Consumer deletes → Message gone ✅
6b. Timeout expires → Back to step 3 ⚠️
```

### States

```
SENT      → Message in queue (visible)
RECEIVED  → Message being processed (invisible)
DELETED   → Message processed successfully
EXPIRED   → Visibility timeout expired (back to SENT)
```

---

## 🔧 SQS Configuration

### Queue Attributes

```yaml
# Standard Queue
QueueName: my-queue
VisibilityTimeout: 30           # 0-43200s (12h)
MessageRetentionPeriod: 345600  # 60s - 1209600s (14 days)
DelaySeconds: 0                 # 0-900s (15 min)
ReceiveMessageWaitTimeSeconds: 20  # Long polling (0-20s)
MaximumMessageSize: 262144      # 1KB - 256KB
```

### FIFO Queue

```yaml
QueueName: my-queue.fifo
ContentBasedDeduplication: true  # Auto deduplication
FifoThroughputLimit: perQueue    # or perMessageGroupId
DeduplicationScope: queue        # or messageGroup
```

---

## 📨 Sending Messages

### Send Single Message

```go
import (
    "github.com/aws/aws-sdk-go-v2/service/sqs"
    "github.com/aws/aws-sdk-go-v2/aws"
)

func sendMessage(ctx context.Context, queueURL string, body string) error {
    result, err := sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
        QueueUrl:    aws.String(queueURL),
        MessageBody: aws.String(body),
        MessageAttributes: map[string]types.MessageAttributeValue{
            "Source": {
                DataType:    aws.String("String"),
                StringValue: aws.String("order-service"),
            },
        },
    })
    
    if err != nil {
        return err
    }
    
    log.Printf("Message sent: %s", *result.MessageId)
    return nil
}
```

### Send Batch (up to 10)

```go
func sendBatch(ctx context.Context, queueURL string, messages []string) error {
    entries := make([]types.SendMessageBatchRequestEntry, len(messages))
    
    for i, msg := range messages {
        entries[i] = types.SendMessageBatchRequestEntry{
            Id:          aws.String(fmt.Sprintf("msg-%d", i)),
            MessageBody: aws.String(msg),
        }
    }
    
    result, err := sqsClient.SendMessageBatch(ctx, &sqs.SendMessageBatchInput{
        QueueUrl: aws.String(queueURL),
        Entries:  entries,
    })
    
    if err != nil {
        return err
    }
    
    log.Printf("Sent %d messages, %d failed", 
        len(result.Successful), len(result.Failed))
    
    return nil
}
```

---

## 📬 Receiving Messages

### Short Polling (❌ Inefficient)

```go
// ReceiveMessageWaitTimeSeconds = 0 (default)
result, err := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
    QueueUrl:            aws.String(queueURL),
    MaxNumberOfMessages: 10,
    WaitTimeSeconds:     0,  // Short polling (return immediately)
})

// Problem: Many empty responses, high costs!
```

### Long Polling (✅ Efficient)

```go
// ReceiveMessageWaitTimeSeconds = 20
result, err := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
    QueueUrl:            aws.String(queueURL),
    MaxNumberOfMessages: 10,
    WaitTimeSeconds:     20,  // Long polling (wait up to 20s)
})

// Benefits: Fewer empty responses, lower costs
```

---

## 🎯 Best Practices

### 1. Always Delete Messages

```go
func processMessage(ctx context.Context, msg types.Message) error {
    // Process
    err := doWork(msg)
    if err != nil {
        return err  // Don't delete, will retry
    }
    
    // Delete on success
    _, err = sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
        QueueUrl:      aws.String(queueURL),
        ReceiptHandle: msg.ReceiptHandle,
    })
    
    return err
}
```

### 2. Set Appropriate Visibility Timeout

```go
// Processing takes ~10s on average
// Set visibility timeout to 30s (3x buffer)
VisibilityTimeout: 30
```

### 3. Implement Idempotency

```go
// Use UUID or hash as idempotency key
type Message struct {
    IdempotencyKey string `json:"idempotency_key"`
    OrderID        string `json:"order_id"`
    Amount         float64 `json:"amount"`
}

// Check before processing
if alreadyProcessed(msg.IdempotencyKey) {
    return nil  // Skip
}
```

### 4. Use Long Polling

```go
// Set WaitTimeSeconds to 20 (max)
WaitTimeSeconds: 20  // ✅ Reduces empty responses
```

### 5. Batch Operations

```go
// Send in batches (up to 10)
SendMessageBatch()

// Receive in batches (up to 10)
ReceiveMessage() with MaxNumberOfMessages: 10

// Delete in batches (up to 10)
DeleteMessageBatch()
```

---

## 🔍 Monitoring

### CloudWatch Metrics

```
ApproximateNumberOfMessagesVisible  # Messages in queue
ApproximateNumberOfMessagesNotVisible  # Being processed
ApproximateAgeOfOldestMessage  # Oldest message age
NumberOfMessagesReceived  # Messages received
NumberOfMessagesSent  # Messages sent
NumberOfMessagesDeleted  # Messages processed
```

### Alarms

```bash
# Queue depth alarm
aws cloudwatch put-metric-alarm \
  --alarm-name high-queue-depth \
  --metric-name ApproximateNumberOfMessagesVisible \
  --namespace AWS/SQS \
  --statistic Average \
  --period 300 \
  --threshold 1000 \
  --comparison-operator GreaterThanThreshold

# Old message alarm
aws cloudwatch put-metric-alarm \
  --alarm-name old-messages \
  --metric-name ApproximateAgeOfOldestMessage \
  --namespace AWS/SQS \
  --statistic Maximum \
  --period 60 \
  --threshold 300 \
  --comparison-operator GreaterThanThreshold
```

---

## 🎓 Висновок

### At-Least-Once Delivery:

✅ **Message delivered** принаймні один раз  
⚠️ **May duplicate** (network issues, timeouts)  
✅ **High throughput** (Standard queue)  
✅ **Reliable** (replicated across AZs)  

### Key Points:

1. **At-least-once** = може бути дублікати
2. **Visibility timeout** = час на обробку
3. **Idempotency** = захист від дублікатів
4. **Long polling** = ефективніше
5. **Always delete** після успішної обробки

### Golden Rule:

**"Design for at-least-once, implement idempotency!"**

---

## 📖 Далі

- `02_dlq_dead_letter_queue.md` - Dead Letter Queue
- `practice/01_lambda_consumer/` - Lambda SQS consumer
- `practice/02_retry_handling/` - Retry strategies

**"At-least-once = Handle duplicates, ensure reliability!" 📨**
