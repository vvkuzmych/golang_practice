# Practice 04: Full Stack (Lambda + SQS + IAM)

## 🎯 Мета

Створити повний async processing stack: SQS → Lambda → CloudWatch з правильними IAM permissions.

---

## 📊 Що створюється

1. **SQS Main Queue** - З DLQ та long polling
2. **SQS Dead Letter Queue** - Для poison messages
3. **Lambda Function** - Go runtime, async processing
4. **IAM Role** - З SQS permissions
5. **Event Source Mapping** - SQS → Lambda інтеграція
6. **CloudWatch Logs** - Lambda logs
7. **CloudWatch Alarms** - DLQ та Lambda errors

---

## 🏗️ Архітектура

```
Producer → SQS Queue (order-processor-queue)
              ↓
        Lambda polls (Event Source Mapping)
              ↓
        Lambda Function (order-processor)
              ↓
        ✅ Success → Delete message
              OR
        ❌ Failure (3x) → DLQ
              ↓
        CloudWatch Alarm 🚨
              ↓
        CloudWatch Logs
```

---

## 🚀 Швидкий старт

### 1. Build Lambda

```bash
# Initialize Go module
go mod init order-processor
go get github.com/aws/aws-lambda-go/lambda
go get github.com/aws/aws-lambda-go/events
go mod tidy

# Build
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bootstrap main.go

# Create ZIP
zip function.zip bootstrap
```

### 2. Deploy

```bash
# Initialize Terraform
terraform init

# Plan
terraform plan

# Apply
terraform apply
```

### 3. Send Test Messages

```bash
# Get queue URL
QUEUE_URL=$(terraform output -raw queue_url)

# Send single message
aws sqs send-message \
  --queue-url $QUEUE_URL \
  --message-body '{
    "order_id": "ORDER-001",
    "amount": 99.99,
    "customer_id": "CUST-123",
    "timestamp": "2026-01-28T10:00:00Z"
  }'

# Send batch
for i in {1..10}; do
  aws sqs send-message \
    --queue-url $QUEUE_URL \
    --message-body "{
      \"order_id\": \"ORDER-$i\",
      \"amount\": $((100 * i)),
      \"customer_id\": \"CUST-$i\",
      \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
    }"
done
```

### 4. Monitor Lambda Execution

```bash
# Tail logs
LOG_GROUP=$(terraform output -raw log_group_name)
aws logs tail $LOG_GROUP --follow

# Expected output:
# Processing batch of 10 messages
# Processing message: abc123...
# Processing order: ID=ORDER-001 Customer=CUST-123 Amount=$99.99
# Order processed successfully: ORDER-001
```

### 5. Check Metrics

```bash
# Lambda invocations
aws cloudwatch get-metric-statistics \
  --namespace AWS/Lambda \
  --metric-name Invocations \
  --dimensions Name=FunctionName,Value=order-processor \
  --start-time $(date -u -d '1 hour ago' +%Y-%m-%dT%H:%M:%S) \
  --end-time $(date -u +%Y-%m-%dT%H:%M:%S) \
  --period 300 \
  --statistics Sum

# Lambda errors
aws cloudwatch get-metric-statistics \
  --namespace AWS/Lambda \
  --metric-name Errors \
  --dimensions Name=FunctionName,Value=order-processor \
  --start-time $(date -u -d '1 hour ago' +%Y-%m-%dT%H:%M:%S) \
  --end-time $(date -u +%Y-%m-%dT%H:%M:%S) \
  --period 300 \
  --statistics Sum
```

### 6. Test Error Handling

```bash
# Send invalid message (will trigger retry → DLQ)
aws sqs send-message \
  --queue-url $QUEUE_URL \
  --message-body 'INVALID_JSON'

# Wait 2 minutes (3 retries with visibility timeout)
sleep 120

# Check DLQ
DLQ_URL=$(terraform output -raw dlq_url)
aws sqs receive-message --queue-url $DLQ_URL

# Should show the poison message
```

### 7. Cleanup

```bash
terraform destroy
```

---

## 📝 Key Components

### 1. Event Source Mapping

**Connects SQS to Lambda:**

```hcl
resource "aws_lambda_event_source_mapping" "sqs" {
  event_source_arn = aws_sqs_queue.main.arn
  function_name    = aws_lambda_function.processor.arn
  
  batch_size                         = 10
  maximum_batching_window_in_seconds = 5
  function_response_types            = ["ReportBatchItemFailures"]
}
```

**Parameters:**
- `batch_size`: How many messages per Lambda invocation (1-10)
- `maximum_batching_window_in_seconds`: Wait time to fill batch (0-300)
- `function_response_types`: Partial batch failure support

**How it works:**
```
1. Lambda polls SQS (long polling)
2. Receives up to 10 messages
3. Invokes Lambda with batch
4. Lambda returns success/failure per message
5. Lambda deletes successful messages
6. Failed messages → retry (visibility timeout)
7. After maxReceiveCount → DLQ
```

### 2. Partial Batch Failure

**Lambda returns which messages failed:**

```go
func handler(ctx context.Context, sqsEvent events.SQSEvent) (events.SQSEventResponse, error) {
    var batchItemFailures []events.SQSBatchItemFailure
    
    for _, record := range sqsEvent.Records {
        if err := processMessage(record); err != nil {
            // Mark this specific message as failed
            batchItemFailures = append(batchItemFailures, events.SQSBatchItemFailure{
                ItemIdentifier: record.MessageId,
            })
        }
    }
    
    return events.SQSEventResponse{
        BatchItemFailures: batchItemFailures,
    }, nil
}
```

**Benefits:**
- Don't reprocess successful messages
- Only failed messages retry
- More efficient

### 3. Idempotency

**Prevent duplicate processing:**

```go
var processedMessages = make(map[string]bool)

func handler(ctx context.Context, sqsEvent events.SQSEvent) {
    for _, record := range sqsEvent.Records {
        // Check if already processed
        if processedMessages[record.MessageId] {
            log.Printf("Skipping duplicate: %s", record.MessageId)
            continue
        }
        
        processMessage(record)
        
        // Mark as processed
        processedMessages[record.MessageId] = true
    }
}
```

**In production:**
- Use DynamoDB or Redis for persistence
- Set TTL (24 hours)

---

## 📊 Monitoring

### CloudWatch Alarms

**1. DLQ Not Empty**
```
Triggers when: DLQ has messages
Action: Investigate poison messages
```

**2. Lambda Errors**
```
Triggers when: > 5 errors in 5 minutes
Action: Check logs, fix code
```

### Key Metrics

**SQS:**
- `ApproximateNumberOfMessagesVisible` - Queue depth
- `ApproximateAgeOfOldestMessage` - Message lag
- `NumberOfMessagesSent` - Producer rate
- `NumberOfMessagesDeleted` - Consumer success rate

**Lambda:**
- `Invocations` - Total calls
- `Errors` - Failed invocations
- `Duration` - Execution time
- `ConcurrentExecutions` - Active instances
- `Throttles` - Rate limited

---

## ⚙️ Configuration

### Environment-Specific Config

**`terraform.tfvars`:**
```hcl
aws_region   = "us-east-1"
project_name = "order-processor"
environment  = "dev"
```

**For production:**
```hcl
environment = "prod"
```

### Scaling

**Increase batch size:**
```hcl
resource "aws_lambda_event_source_mapping" "sqs" {
  batch_size = 10  # Max for SQS
}
```

**Concurrent Lambda executions:**
```hcl
resource "aws_lambda_function" "processor" {
  reserved_concurrent_executions = 100  # Max concurrent
}
```

**Queue parameters:**
```hcl
resource "aws_sqs_queue" "main" {
  visibility_timeout_seconds = 300  # 5 min
  receive_wait_time_seconds  = 20   # Long polling
}
```

---

## 🧪 Testing

### Load Test

```bash
# Send 1000 messages
for i in {1..1000}; do
  aws sqs send-message \
    --queue-url $(terraform output -raw queue_url) \
    --message-body "{\"order_id\":\"$i\",\"amount\":100}" \
    > /dev/null
done

# Monitor processing
aws logs tail $(terraform output -raw log_group_name) --follow
```

### Error Rate Test

```bash
# Send 50% valid, 50% invalid
for i in {1..100}; do
  if [ $((i % 2)) -eq 0 ]; then
    # Valid
    BODY="{\"order_id\":\"$i\",\"amount\":100}"
  else
    # Invalid
    BODY="INVALID_JSON_$i"
  fi
  
  aws sqs send-message \
    --queue-url $(terraform output -raw queue_url) \
    --message-body "$BODY"
done

# Check DLQ after processing
sleep 60
aws sqs receive-message --queue-url $(terraform output -raw dlq_url)
```

---

## ✅ Best Practices

### 1. Use Partial Batch Failures

```hcl
resource "aws_lambda_event_source_mapping" "sqs" {
  function_response_types = ["ReportBatchItemFailures"]
}
```

### 2. Implement Idempotency

```go
// Use message ID as idempotency key
if processed[record.MessageId] {
    continue
}
```

### 3. Set Appropriate Timeouts

```hcl
# Lambda timeout > visibility timeout
resource "aws_lambda_function" "processor" {
  timeout = 60  # 1 minute
}

resource "aws_sqs_queue" "main" {
  visibility_timeout_seconds = 300  # 5 minutes (6x Lambda)
}
```

### 4. Enable Long Polling

```hcl
resource "aws_sqs_queue" "main" {
  receive_wait_time_seconds = 20  # ✅ Reduces costs
}
```

### 5. Monitor DLQ

```hcl
resource "aws_cloudwatch_metric_alarm" "dlq_not_empty" {
  threshold = 0  # Alert immediately
}
```

---

## 🎓 Висновок

Ти створив production-ready async processing pipeline:

- ✅ SQS для message queuing
- ✅ Lambda для processing
- ✅ IAM для security
- ✅ DLQ для poison messages
- ✅ CloudWatch для monitoring
- ✅ Idempotency для at-least-once delivery
- ✅ Partial batch failures для efficiency

**Everything as code with Terraform!** 🏗️
