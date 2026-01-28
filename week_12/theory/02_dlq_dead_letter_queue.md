# Dead Letter Queue (DLQ)

## 🎯 Що таке DLQ?

**Dead Letter Queue (DLQ)** - це окрема SQS queue для messages, що не вдалося обробити після максимальної кількості спроб.

```
Main Queue → Processing fails → Retry → Retry → Retry → DLQ
             (attempt 1)         (2)      (3)      (4)     (max reached)
```

**Мета:** Isolate poison messages, prevent blocking queue

---

## 📊 Why DLQ?

### Проблема без DLQ

```
Main Queue: [Msg1] [Msg2] [PoisonMsg] [Msg4] [Msg5]
                             ↑
                    Fails repeatedly!
                    Blocks processing!
```

**Issues:**
- Poison message блокує queue
- Wasted retries (costs money)
- Good messages delayed
- Hard to debug

### Рішення з DLQ

```
Main Queue: [Msg1] [Msg2] [Msg4] [Msg5]
                    ↓
            [PoisonMsg] moved to DLQ
                    ↓
            Investigate & fix
```

**Benefits:**
- Queue не блокується
- Failed messages isolated
- Can analyze failures
- Manual reprocessing

---

## 🔧 DLQ Configuration

### Create DLQ

```bash
# 1. Create DLQ
aws sqs create-queue --queue-name my-queue-dlq

# 2. Get DLQ ARN
DLQ_ARN=$(aws sqs get-queue-attributes \
  --queue-url https://sqs.region.amazonaws.com/account/my-queue-dlq \
  --attribute-names QueueArn \
  --query 'Attributes.QueueArn' --output text)

# 3. Configure main queue with DLQ
aws sqs set-queue-attributes \
  --queue-url https://sqs.region.amazonaws.com/account/my-queue \
  --attributes '{
    "RedrivePolicy": "{\"deadLetterTargetArn\":\"'$DLQ_ARN'\",\"maxReceiveCount\":\"3\"}"
  }'
```

### Redrive Policy

```json
{
  "deadLetterTargetArn": "arn:aws:sqs:region:account:my-queue-dlq",
  "maxReceiveCount": "3"
}
```

**maxReceiveCount:**
- Number of receive attempts before DLQ
- Recommended: 3-5
- Too low: False positives
- Too high: Wasted resources

---

## 📊 Message Flow

### Normal Flow (Success)

```
1. Message sent to Main Queue
2. Lambda receives message (attempt 1)
3. Processing succeeds
4. Lambda deletes message
5. Done ✅
```

### Retry Flow (Transient Error)

```
1. Message sent to Main Queue
2. Lambda receives message (attempt 1)
3. Processing fails (timeout)
4. Message returns to queue (visibility timeout)
5. Lambda receives message (attempt 2)
6. Processing succeeds
7. Lambda deletes message
8. Done ✅
```

### DLQ Flow (Permanent Error)

```
1. Message sent to Main Queue
2. Lambda receives message (attempt 1)
3. Processing fails (invalid format)
4. Message returns to queue
5. Lambda receives message (attempt 2)
6. Processing fails again
7. Message returns to queue
8. Lambda receives message (attempt 3)
9. Processing fails again
10. maxReceiveCount (3) reached
11. Message moved to DLQ automatically
12. Alert triggered 🚨
```

---

## 🎯 Lambda + SQS + DLQ

### Lambda Configuration

```yaml
# Lambda Event Source Mapping
FunctionName: my-processor
EventSourceArn: arn:aws:sqs:region:account:my-queue
BatchSize: 10
MaximumBatchingWindowInSeconds: 5

# Lambda will automatically:
# - Poll SQS
# - Delete on success
# - Retry on failure
# - Move to DLQ after maxReceiveCount
```

### Lambda Handler

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
)

type OrderMessage struct {
    OrderID        string  `json:"order_id"`
    Amount         float64 `json:"amount"`
    CustomerID     string  `json:"customer_id"`
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
    log.Printf("Processing %d messages", len(sqsEvent.Records))
    
    for _, record := range sqsEvent.Records {
        err := processMessage(ctx, record)
        if err != nil {
            // Lambda will retry automatically
            // After maxReceiveCount, moved to DLQ
            log.Printf("ERROR processing message %s: %v", record.MessageId, err)
            return err  // Fail entire batch
        }
    }
    
    return nil  // Success, Lambda deletes all messages
}

func processMessage(ctx context.Context, record events.SQSMessage) error {
    log.Printf("Processing message: %s", record.MessageId)
    
    // Parse message
    var msg OrderMessage
    err := json.Unmarshal([]byte(record.Body), &msg)
    if err != nil {
        // Invalid JSON = poison message
        log.Printf("Invalid JSON: %v", err)
        return err  // Will eventually go to DLQ
    }
    
    // Validate
    if msg.OrderID == "" || msg.Amount <= 0 {
        // Invalid data = poison message
        return fmt.Errorf("invalid order: %+v", msg)
    }
    
    // Process order
    err = createOrder(msg.OrderID, msg.Amount, msg.CustomerID)
    if err != nil {
        return err
    }
    
    log.Printf("Order created: %s", msg.OrderID)
    return nil
}

func createOrder(orderID string, amount float64, customerID string) error {
    // Business logic
    return nil
}

func main() {
    lambda.Start(handler)
}
```

---

## ⚠️ Poison Messages

### Types of Poison Messages

**1. Invalid Format**
```json
{
  "order_id": "invalid json  // ← Missing closing brace
}
```

**2. Missing Required Fields**
```json
{
  "amount": 100
  // Missing order_id!
}
```

**3. Invalid Data**
```json
{
  "order_id": "123",
  "amount": -50  // ← Negative amount!
}
```

**4. External Dependency Down**
```json
{
  "order_id": "123",
  "external_api_id": "nonexistent"  // ← 404 from external API
}
```

### Handling Strategies

#### Strategy 1: Fail Fast (Invalid Format)

```go
func processMessage(msg string) error {
    var data Order
    err := json.Unmarshal([]byte(msg), &data)
    if err != nil {
        // Invalid JSON = permanent error
        // Let it go to DLQ immediately
        return err
    }
    
    // Process...
}
```

#### Strategy 2: Validation Error (Invalid Data)

```go
func processMessage(msg Order) error {
    // Validate
    if msg.Amount <= 0 {
        // Invalid data = permanent error
        return fmt.Errorf("invalid amount: %f", msg.Amount)
    }
    
    // Process...
}
```

#### Strategy 3: Transient Error (Retry Worth It)

```go
func processMessage(msg Order) error {
    err := callExternalAPI(msg)
    if err != nil {
        if isTransientError(err) {
            // Retry makes sense (timeout, rate limit)
            return err
        } else {
            // Permanent error (404, 400)
            // Mark as unprocessable, skip retry
            return nil  // Or send to DLQ explicitly
        }
    }
    
    // Process...
}

func isTransientError(err error) bool {
    // Check for timeouts, 5xx errors, rate limits
    return isTimeout(err) || is5xxError(err) || isRateLimitError(err)
}
```

---

## 🔍 DLQ Monitoring & Analysis

### CloudWatch Alarms

```bash
# Alert when messages in DLQ
aws cloudwatch put-metric-alarm \
  --alarm-name dlq-messages \
  --metric-name ApproximateNumberOfMessagesVisible \
  --namespace AWS/SQS \
  --dimensions Name=QueueName,Value=my-queue-dlq \
  --statistic Sum \
  --period 60 \
  --threshold 1 \
  --comparison-operator GreaterThanThreshold \
  --evaluation-periods 1

# SNS notification
aws sns publish \
  --topic-arn arn:aws:sns:region:account:alerts \
  --message "Messages in DLQ detected!"
```

### Analyze DLQ Messages

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    
    "github.com/aws/aws-sdk-go-v2/service/sqs"
)

func analyzeDLQ(ctx context.Context, dlqURL string) error {
    // Receive messages from DLQ (don't delete)
    result, err := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
        QueueUrl:            aws.String(dlqURL),
        MaxNumberOfMessages: 10,
        WaitTimeSeconds:     5,
        AttributeNames:      []types.QueueAttributeName{"All"},
        MessageAttributeNames: []string{"All"},
    })
    
    if err != nil {
        return err
    }
    
    log.Printf("Found %d messages in DLQ", len(result.Messages))
    
    for _, msg := range result.Messages {
        // Analyze failure reason
        log.Printf("Message ID: %s", *msg.MessageId)
        log.Printf("Receive Count: %s", msg.Attributes["ApproximateReceiveCount"])
        log.Printf("First Received: %s", msg.Attributes["ApproximateFirstReceiveTimestamp"])
        log.Printf("Body: %s", *msg.Body)
        
        // Try to understand why it failed
        var data map[string]interface{}
        if err := json.Unmarshal([]byte(*msg.Body), &data); err != nil {
            log.Printf("  → Invalid JSON format")
        } else {
            log.Printf("  → Valid JSON, check business logic")
        }
        
        fmt.Println("---")
    }
    
    return nil
}
```

---

## 🔄 Reprocessing from DLQ

### Manual Reprocessing

```go
func reprocessDLQ(ctx context.Context, dlqURL, mainQueueURL string) error {
    for {
        // Receive from DLQ
        result, err := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
            QueueUrl:            aws.String(dlqURL),
            MaxNumberOfMessages: 10,
            WaitTimeSeconds:     5,
        })
        
        if err != nil {
            return err
        }
        
        if len(result.Messages) == 0 {
            break  // No more messages
        }
        
        for _, msg := range result.Messages {
            // Option 1: Fix and reprocess
            fixed, err := fixMessage(*msg.Body)
            if err != nil {
                log.Printf("Cannot fix message %s: %v", *msg.MessageId, err)
                continue
            }
            
            // Send back to main queue
            _, err = sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
                QueueUrl:    aws.String(mainQueueURL),
                MessageBody: aws.String(fixed),
            })
            
            if err != nil {
                return err
            }
            
            // Delete from DLQ
            _, err = sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
                QueueUrl:      aws.String(dlqURL),
                ReceiptHandle: msg.ReceiptHandle,
            })
            
            if err != nil {
                return err
            }
            
            log.Printf("Reprocessed message: %s", *msg.MessageId)
        }
    }
    
    return nil
}

func fixMessage(body string) (string, error) {
    // Fix common issues
    // e.g., add missing fields, fix format, etc.
    return body, nil
}
```

### Redrive to Source (AWS Feature)

```bash
# Move messages from DLQ back to main queue
aws sqs start-message-move-task \
  --source-arn arn:aws:sqs:region:account:my-queue-dlq \
  --destination-arn arn:aws:sqs:region:account:my-queue \
  --max-number-of-messages-per-second 10
```

---

## ✅ Best Practices

### 1. Set Appropriate maxReceiveCount

```yaml
# Too low (1-2): False positives
# Too high (10+): Wasted resources
# Recommended: 3-5
maxReceiveCount: 3
```

### 2. Different Retention for DLQ

```bash
# Main queue: 4 days
# DLQ: 14 days (max, for investigation)
aws sqs set-queue-attributes \
  --queue-url DLQ_URL \
  --attributes MessageRetentionPeriod=1209600
```

### 3. Monitor DLQ Depth

```bash
# Alert if > 0 messages
aws cloudwatch put-metric-alarm \
  --alarm-name dlq-not-empty \
  --metric-name ApproximateNumberOfMessagesVisible \
  --threshold 0 \
  --comparison-operator GreaterThanThreshold
```

### 4. Log Failure Reasons

```go
func processMessage(msg string) error {
    err := process(msg)
    if err != nil {
        // Log detailed error for DLQ investigation
        log.Printf("FAILURE: message=%s, error=%v, stack=%s", 
            msg, err, debug.Stack())
        return err
    }
    return nil
}
```

### 5. Separate DLQ per Queue

```
orders-queue → orders-dlq
payments-queue → payments-dlq
notifications-queue → notifications-dlq
```

---

## 🎓 Висновок

### Dead Letter Queue:

✅ **Isolates** poison messages  
✅ **Prevents** queue blocking  
✅ **Enables** failure analysis  
✅ **Supports** reprocessing  
✅ **Reduces** wasted resources  

### Key Points:

1. **DLQ** = окрема queue для failed messages
2. **maxReceiveCount** = кількість спроб (3-5)
3. **Monitor** DLQ depth (alert if > 0)
4. **Analyze** failures для fix
5. **Reprocess** після fix

### Golden Rule:

**"Every queue needs a DLQ, every DLQ needs monitoring!"**

---

## 📖 Далі

- `practice/01_lambda_consumer/` - Lambda SQS consumer
- `practice/02_retry_handling/` - Retry strategies with DLQ
- DLQ analysis tools

**"DLQ = Safety net for your messages!" 🛡️📨**
