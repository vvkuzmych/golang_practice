# Week 12: SQS + Lambda (Async Processing)

## 🎯 Мета

Навчитися створювати **async event-driven applications** з Amazon SQS + AWS Lambda, розуміти at-least-once delivery, та працювати з Dead Letter Queue.

---

## 📚 Теорія

### 1. SQS & At-Least-Once Delivery
**Файл:** `theory/01_sqs_at_least_once.md`

- Amazon SQS basics
- Standard vs FIFO queues
- At-least-once delivery semantics
- Visibility timeout
- Idempotency patterns
- Long polling vs short polling

**Key concepts:**
```
At-Least-Once = Message delivered ≥ 1 time
May duplicate → Implement idempotency!
```

### 2. Dead Letter Queue (DLQ)
**Файл:** `theory/02_dlq_dead_letter_queue.md`

- What is DLQ?
- maxReceiveCount configuration
- Poison messages handling
- DLQ monitoring & analysis
- Reprocessing strategies
- Best practices

**DLQ Flow:**
```
Main Queue → Retry (1) → Retry (2) → Retry (3) → DLQ
             fail         fail         fail         (max reached)
```

---

## 💻 Практика

### 1. Lambda SQS Consumer
**Директорія:** `practice/01_lambda_consumer/`

**Features:**
- Lambda triggered by SQS
- Batch processing (up to 10 messages)
- Automatic retry on failure
- Idempotency with DynamoDB
- Structured logging
- Error handling

**Example code:**
```go
package main

import (
    "context"
    "encoding/json"
    "log"
    
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
)

type Order struct {
    OrderID    string  `json:"order_id"`
    Amount     float64 `json:"amount"`
    CustomerID string  `json:"customer_id"`
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
    log.Printf("Processing %d messages", len(sqsEvent.Records))
    
    for _, record := range sqsEvent.Records {
        if err := processMessage(ctx, record); err != nil {
            // Lambda will retry
            return err
        }
    }
    
    return nil  // Success, all deleted
}

func processMessage(ctx context.Context, record events.SQSMessage) error {
    var order Order
    if err := json.Unmarshal([]byte(record.Body), &order); err != nil {
        return err
    }
    
    // Check idempotency
    if alreadyProcessed(record.MessageId) {
        return nil
    }
    
    // Process order
    if err := createOrder(order); err != nil {
        return err
    }
    
    // Mark as processed
    markProcessed(record.MessageId)
    
    return nil
}

func main() {
    lambda.Start(handler)
}
```

### 2. Retry Handling with DLQ
**Директорія:** `practice/02_retry_handling/`

**Features:**
- Exponential backoff simulation
- Transient vs permanent error detection
- DLQ configuration
- CloudWatch alarms
- Reprocessing from DLQ

**Retry strategies:**
```go
// Transient errors: Retry
- Network timeout
- 5xx server errors
- Rate limiting

// Permanent errors: DLQ
- Invalid JSON format
- Missing required fields
- 404 Not Found
- 400 Bad Request
```

### 3. Batch Processing
**Директорія:** `practice/03_batch_processing/`

**Features:**
- Process multiple messages in batch
- Partial batch failure handling
- Batch size optimization
- Performance metrics

---

## 🔧 Основні команди

### SQS Queue Management

```bash
# Create queue
aws sqs create-queue --queue-name my-queue

# Create DLQ
aws sqs create-queue --queue-name my-queue-dlq

# Get queue URL
aws sqs get-queue-url --queue-name my-queue

# Configure DLQ
aws sqs set-queue-attributes \
  --queue-url QUEUE_URL \
  --attributes '{
    "RedrivePolicy": "{\"deadLetterTargetArn\":\"DLQ_ARN\",\"maxReceiveCount\":\"3\"}"
  }'
```

### Send Messages

```bash
# Send single message
aws sqs send-message \
  --queue-url QUEUE_URL \
  --message-body '{"order_id":"123","amount":100}'

# Send batch
aws sqs send-message-batch \
  --queue-url QUEUE_URL \
  --entries '[
    {"Id":"1","MessageBody":"{\"order_id\":\"123\"}"},
    {"Id":"2","MessageBody":"{\"order_id\":\"456\"}"}
  ]'
```

### Receive & Delete

```bash
# Receive (long polling)
aws sqs receive-message \
  --queue-url QUEUE_URL \
  --max-number-of-messages 10 \
  --wait-time-seconds 20

# Delete
aws sqs delete-message \
  --queue-url QUEUE_URL \
  --receipt-handle RECEIPT_HANDLE
```

### Lambda Event Source Mapping

```bash
# Create mapping
aws lambda create-event-source-mapping \
  --function-name my-processor \
  --event-source-arn arn:aws:sqs:region:account:my-queue \
  --batch-size 10 \
  --maximum-batching-window-in-seconds 5

# Update mapping
aws lambda update-event-source-mapping \
  --uuid MAPPING_UUID \
  --batch-size 5

# List mappings
aws lambda list-event-source-mappings \
  --function-name my-processor
```

---

## 📊 Architecture Patterns

### Pattern 1: Simple Processing

```
Producer → SQS Queue → Lambda → Database
                    ↓
                   DLQ
```

### Pattern 2: Fan-Out

```
SNS Topic
  ├─> SQS Queue 1 → Lambda 1 → Action A
  ├─> SQS Queue 2 → Lambda 2 → Action B
  └─> SQS Queue 3 → Lambda 3 → Action C
```

### Pattern 3: Chain Processing

```
Queue 1 → Lambda 1 → Queue 2 → Lambda 2 → Queue 3 → Lambda 3
   ↓                   ↓                     ↓
  DLQ 1               DLQ 2                 DLQ 3
```

### Pattern 4: Retry with Delay

```
Main Queue (fast retry)
   ↓ (after 3 failures)
Retry Queue (30s delay)
   ↓ (after 3 failures)
DLQ (manual intervention)
```

---

## 🎯 Idempotency Strategies

### Strategy 1: Redis Cache

```go
var redisClient *redis.Client

func init() {
    redisClient = redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
}

func processMessage(messageID string, body string) error {
    // Check if already processed
    exists, err := redisClient.Exists(ctx, messageID).Result()
    if err != nil {
        return err
    }
    
    if exists > 0 {
        log.Printf("Already processed: %s", messageID)
        return nil
    }
    
    // Process
    err = doWork(body)
    if err != nil {
        return err
    }
    
    // Mark as processed (24h TTL)
    redisClient.Set(ctx, messageID, "1", 24*time.Hour)
    
    return nil
}
```

### Strategy 2: DynamoDB

```go
func checkIdempotency(messageID string) (bool, error) {
    result, err := dynamoClient.GetItem(ctx, &dynamodb.GetItemInput{
        TableName: aws.String("processed_messages"),
        Key: map[string]types.AttributeValue{
            "message_id": &types.AttributeValueMemberS{Value: messageID},
        },
    })
    
    if err != nil {
        return false, err
    }
    
    return result.Item != nil, nil
}

func markProcessed(messageID string) error {
    _, err := dynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
        TableName: aws.String("processed_messages"),
        Item: map[string]types.AttributeValue{
            "message_id": &types.AttributeValueMemberS{Value: messageID},
            "processed_at": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", time.Now().Unix())},
            "ttl": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", time.Now().Add(24*time.Hour).Unix())},
        },
    })
    
    return err
}
```

### Strategy 3: Database Unique Constraint

```sql
CREATE TABLE orders (
    id UUID PRIMARY KEY,
    message_id VARCHAR(255) UNIQUE NOT NULL,
    order_id VARCHAR(255),
    amount DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT NOW()
);
```

```go
func createOrder(messageID, orderID string, amount float64) error {
    _, err := db.Exec(
        "INSERT INTO orders (id, message_id, order_id, amount) VALUES ($1, $2, $3, $4)",
        uuid.New(), messageID, orderID, amount,
    )
    
    if err != nil {
        if isDuplicateKeyError(err) {
            log.Printf("Order already exists: %s", messageID)
            return nil  // Idempotent
        }
        return err
    }
    
    return nil
}
```

---

## 📈 Monitoring & Observability

### CloudWatch Metrics

**SQS Metrics:**
```
ApproximateNumberOfMessagesVisible    # Queue depth
ApproximateNumberOfMessagesNotVisible # In-flight
ApproximateAgeOfOldestMessage         # Lag
NumberOfMessagesReceived              # Throughput
NumberOfMessagesSent                  # Producer rate
NumberOfMessagesDeleted               # Success rate
```

**Lambda Metrics:**
```
Invocations          # Total invocations
Errors               # Failed invocations
Duration             # Execution time
ConcurrentExecutions # Active lambdas
Throttles            # Rate limited
```

### CloudWatch Alarms

```bash
# High queue depth
aws cloudwatch put-metric-alarm \
  --alarm-name high-queue-depth \
  --metric-name ApproximateNumberOfMessagesVisible \
  --namespace AWS/SQS \
  --statistic Average \
  --threshold 1000

# Messages in DLQ
aws cloudwatch put-metric-alarm \
  --alarm-name dlq-not-empty \
  --metric-name ApproximateNumberOfMessagesVisible \
  --namespace AWS/SQS \
  --dimensions Name=QueueName,Value=my-queue-dlq \
  --threshold 0 \
  --comparison-operator GreaterThanThreshold

# Old messages
aws cloudwatch put-metric-alarm \
  --alarm-name old-messages \
  --metric-name ApproximateAgeOfOldestMessage \
  --namespace AWS/SQS \
  --threshold 300

# Lambda errors
aws cloudwatch put-metric-alarm \
  --alarm-name lambda-errors \
  --metric-name Errors \
  --namespace AWS/Lambda \
  --threshold 10
```

---

## ✅ Best Practices

### SQS

1. ✅ **Use long polling** (WaitTimeSeconds: 20)
2. ✅ **Set appropriate visibility timeout** (3x processing time)
3. ✅ **Batch operations** (send/receive/delete up to 10)
4. ✅ **Monitor queue depth** (alert if too high)
5. ✅ **Implement idempotency** (handle duplicates)

### DLQ

1. ✅ **Configure DLQ** for every queue
2. ✅ **Set maxReceiveCount** 3-5
3. ✅ **Monitor DLQ depth** (alert if > 0)
4. ✅ **Longer retention** for DLQ (14 days)
5. ✅ **Analyze failures** regularly

### Lambda

1. ✅ **Batch size** 10 for throughput
2. ✅ **Timeout** appropriate (not too long)
3. ✅ **Reserved concurrency** if needed
4. ✅ **Error handling** proper logging
5. ✅ **Idempotency** always implement

---

## 🚀 Quick Start

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_12

# Read theory
cat theory/01_sqs_at_least_once.md
cat theory/02_dlq_dead_letter_queue.md

# Create queue with DLQ
aws sqs create-queue --queue-name test-queue
aws sqs create-queue --queue-name test-queue-dlq

# Get ARNs
QUEUE_URL=$(aws sqs get-queue-url --queue-name test-queue --query 'QueueUrl' --output text)
DLQ_ARN=$(aws sqs get-queue-attributes --queue-url DLQ_URL --attribute-names QueueArn --query 'Attributes.QueueArn' --output text)

# Configure DLQ
aws sqs set-queue-attributes \
  --queue-url $QUEUE_URL \
  --attributes '{
    "RedrivePolicy": "{\"deadLetterTargetArn\":\"'$DLQ_ARN'\",\"maxReceiveCount\":\"3\"}"
  }'

# Deploy Lambda (practice examples)
cd practice/01_lambda_consumer
GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
zip function.zip bootstrap

aws lambda create-function \
  --function-name sqs-processor \
  --runtime provided.al2023 \
  --handler bootstrap \
  --zip-file fileb://function.zip \
  --role LAMBDA_ROLE_ARN

# Create event source mapping
aws lambda create-event-source-mapping \
  --function-name sqs-processor \
  --event-source-arn QUEUE_ARN \
  --batch-size 10
```

---

## 🎓 Learning Path

### День 1: SQS Basics

1. Читай `theory/01_sqs_at_least_once.md`
2. Create SQS queue
3. Send/receive messages via AWS CLI
4. Understand visibility timeout

### День 2: Lambda Consumer

1. Deploy `practice/01_lambda_consumer/`
2. Create event source mapping
3. Test message processing
4. Monitor CloudWatch logs

### День 3: DLQ & Retry

1. Читай `theory/02_dlq_dead_letter_queue.md`
2. Configure DLQ with maxReceiveCount
3. Test poison messages
4. Monitor DLQ depth

### День 4: Production Ready

1. Implement idempotency
2. Set up CloudWatch alarms
3. Test failure scenarios
4. Document runbooks

---

## 📖 Ресурси

- [Amazon SQS Documentation](https://docs.aws.amazon.com/sqs/)
- [Lambda + SQS Tutorial](https://docs.aws.amazon.com/lambda/latest/dg/with-sqs.html)
- [SQS Best Practices](https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-best-practices.html)
- Week 11: AWS Lambda basics

---

**"Async = Decoupled, Scalable, Resilient!" 📨⚡**

**Status:** Week 12 Materials Complete ✅  
**Created:** 2026-01-28
