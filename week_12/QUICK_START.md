# Week 12 - Quick Start 📨

## 🎯 Мета

Навчитися створювати **async event-driven apps** з Amazon SQS + AWS Lambda.

---

## ⚡ 5-хвилинний старт

### 1. Create SQS Queue with DLQ

```bash
# Create main queue
aws sqs create-queue --queue-name orders-queue

# Create DLQ
aws sqs create-queue --queue-name orders-queue-dlq

# Get URLs
QUEUE_URL=$(aws sqs get-queue-url --queue-name orders-queue --query 'QueueUrl' --output text)
DLQ_URL=$(aws sqs get-queue-url --queue-name orders-queue-dlq --query 'QueueUrl' --output text)

# Get DLQ ARN
DLQ_ARN=$(aws sqs get-queue-attributes \
  --queue-url $DLQ_URL \
  --attribute-names QueueArn \
  --query 'Attributes.QueueArn' --output text)

# Configure DLQ (max 3 retries)
aws sqs set-queue-attributes \
  --queue-url $QUEUE_URL \
  --attributes "{\"RedrivePolicy\":\"{\\\"deadLetterTargetArn\\\":\\\"$DLQ_ARN\\\",\\\"maxReceiveCount\\\":\\\"3\\\"}\"}"
```

### 2. Send Test Message

```bash
# Send message
aws sqs send-message \
  --queue-url $QUEUE_URL \
  --message-body '{"order_id":"123","amount":100,"customer_id":"alice"}'

# Check queue depth
aws sqs get-queue-attributes \
  --queue-url $QUEUE_URL \
  --attribute-names ApproximateNumberOfMessagesVisible
```

### 3. Deploy Lambda Consumer

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_12/practice/01_lambda_consumer

# Build
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bootstrap main.go
zip function.zip bootstrap

# Deploy
aws lambda create-function \
  --function-name sqs-processor \
  --runtime provided.al2023 \
  --handler bootstrap \
  --zip-file fileb://function.zip \
  --role arn:aws:iam::ACCOUNT:role/lambda-role \
  --timeout 30

# Create event source mapping
QUEUE_ARN=$(aws sqs get-queue-attributes \
  --queue-url $QUEUE_URL \
  --attribute-names QueueArn \
  --query 'Attributes.QueueArn' --output text)

aws lambda create-event-source-mapping \
  --function-name sqs-processor \
  --event-source-arn $QUEUE_ARN \
  --batch-size 10 \
  --maximum-batching-window-in-seconds 5
```

---

## 📚 Читати теорію

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_12

# SQS & At-Least-Once
cat theory/01_sqs_at_least_once.md

# Dead Letter Queue
cat theory/02_dlq_dead_letter_queue.md
```

---

## 🎯 2 Головні Концепції

### 1. At-Least-Once Delivery

```
Send → SQS stores → Consumer receives
                 → May receive AGAIN!

Why duplicates?
1. Consumer didn't ACK in time
2. Visibility timeout expired
3. Network issue

Solution: Idempotency!
```

**Idempotency pattern:**
```go
func processMessage(messageID string, body string) error {
    // Check if already processed
    if alreadyProcessed(messageID) {
        return nil  // Skip duplicate
    }
    
    // Process
    err := doWork(body)
    if err != nil {
        return err
    }
    
    // Mark as processed
    markProcessed(messageID)
    return nil
}
```

### 2. Dead Letter Queue

```
Main Queue
  ↓ Process (attempt 1)
  ↓ Fail, retry
  ↓ Process (attempt 2)
  ↓ Fail, retry
  ↓ Process (attempt 3)
  ↓ Fail, maxReceiveCount reached
  ↓
DLQ (for investigation)
```

**Configuration:**
```bash
# maxReceiveCount = 3 (recommended)
aws sqs set-queue-attributes \
  --queue-url QUEUE_URL \
  --attributes '{
    "RedrivePolicy": "{\"deadLetterTargetArn\":\"DLQ_ARN\",\"maxReceiveCount\":\"3\"}"
  }'
```

---

## 🔧 Essential Commands

### Queue Operations

```bash
# Create queue
aws sqs create-queue --queue-name my-queue

# Get URL
aws sqs get-queue-url --queue-name my-queue

# Send message
aws sqs send-message \
  --queue-url QUEUE_URL \
  --message-body '{"key":"value"}'

# Receive (long polling)
aws sqs receive-message \
  --queue-url QUEUE_URL \
  --max-number-of-messages 10 \
  --wait-time-seconds 20

# Delete
aws sqs delete-message \
  --queue-url QUEUE_URL \
  --receipt-handle RECEIPT_HANDLE

# Get attributes
aws sqs get-queue-attributes \
  --queue-url QUEUE_URL \
  --attribute-names All
```

### Lambda Event Source

```bash
# Create mapping
aws lambda create-event-source-mapping \
  --function-name my-func \
  --event-source-arn QUEUE_ARN \
  --batch-size 10

# List mappings
aws lambda list-event-source-mappings \
  --function-name my-func

# Update mapping
aws lambda update-event-source-mapping \
  --uuid MAPPING_UUID \
  --batch-size 5

# Delete mapping
aws lambda delete-event-source-mapping \
  --uuid MAPPING_UUID
```

---

## 📊 Quick Patterns

### Pattern 1: Lambda Consumer

```go
package main

import (
    "context"
    "encoding/json"
    "log"
    
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
    for _, record := range sqsEvent.Records {
        log.Printf("Processing: %s", record.MessageId)
        
        var msg map[string]interface{}
        json.Unmarshal([]byte(record.Body), &msg)
        
        // Process message
        if err := process(msg); err != nil {
            return err  // Will retry
        }
    }
    
    return nil  // Success, all deleted
}

func main() {
    lambda.Start(handler)
}
```

### Pattern 2: Idempotency Check

```go
import "github.com/aws/aws-sdk-go-v2/service/dynamodb"

func processMessage(messageID string) error {
    // Check DynamoDB
    exists, _ := checkProcessed(messageID)
    if exists {
        return nil  // Already processed
    }
    
    // Do work
    err := doWork()
    if err != nil {
        return err
    }
    
    // Mark processed
    markProcessed(messageID)
    return nil
}
```

### Pattern 3: Error Types

```go
func processMessage(msg string) error {
    err := validate(msg)
    if err != nil {
        // Permanent error → will go to DLQ
        return err
    }
    
    err = callAPI(msg)
    if err != nil {
        if isTransient(err) {
            // Retry makes sense
            return err
        } else {
            // Permanent (404, 400) → DLQ
            return err
        }
    }
    
    return nil
}
```

---

## ✅ Quick Checklist

### Before Production

- [ ] Queue created with DLQ
- [ ] maxReceiveCount = 3-5
- [ ] Long polling enabled (WaitTimeSeconds: 20)
- [ ] Lambda timeout appropriate
- [ ] Idempotency implemented
- [ ] Error handling proper

### Monitoring

- [ ] CloudWatch alarm for queue depth
- [ ] CloudWatch alarm for DLQ depth
- [ ] CloudWatch alarm for Lambda errors
- [ ] CloudWatch Logs Insights queries
- [ ] Dashboard created

---

## 🎓 Golden Rules

### Rule 1: At-Least-Once = Idempotency

```go
// Always check if already processed
if alreadyProcessed(messageID) {
    return nil
}
```

### Rule 2: DLQ for Every Queue

```bash
# Every queue needs DLQ
# maxReceiveCount: 3-5
```

### Rule 3: Long Polling

```bash
# WaitTimeSeconds: 20 (max)
# Reduces empty responses, costs
```

### Rule 4: Monitor DLQ

```bash
# Alert if DLQ depth > 0
# Investigate immediately
```

### Rule 5: Batch Operations

```bash
# Send in batches (up to 10)
# Receive in batches (up to 10)
# Cheaper, faster
```

---

## 🚀 Quick Test

```bash
# 1. Create queue
aws sqs create-queue --queue-name test-queue
QUEUE_URL=$(aws sqs get-queue-url --queue-name test-queue --query 'QueueUrl' --output text)

# 2. Send message
aws sqs send-message \
  --queue-url $QUEUE_URL \
  --message-body '{"test":"message"}'

# 3. Receive
aws sqs receive-message \
  --queue-url $QUEUE_URL \
  --max-number-of-messages 1

# 4. Delete
aws sqs delete-message \
  --queue-url $QUEUE_URL \
  --receipt-handle "RECEIPT_HANDLE_FROM_RECEIVE"

# 5. Clean up
aws sqs delete-queue --queue-url $QUEUE_URL
```

---

**"Async = Decouple, Scale, Resilient!" 📨⚡**
