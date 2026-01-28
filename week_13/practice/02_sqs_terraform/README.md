# Practice 02: SQS with Terraform

## 🎯 Мета

Навчитися створювати SQS queue з Dead Letter Queue та CloudWatch alarms за допомогою Terraform.

---

## 📊 Що створюється

1. **Main SQS Queue** - З long polling та DLQ
2. **Dead Letter Queue (DLQ)** - Для poison messages
3. **CloudWatch Alarms** - Моніторинг (3 alarms)

---

## 🏗️ Архітектура

```
Producer → Main Queue (my-queue)
              ↓
        Consumer tries 3x
              ↓
         ✅ Success
              OR
         ❌ After 3 failures
              ↓
         Dead Letter Queue (my-queue-dlq)
              ↓
         CloudWatch Alarm 🚨
```

---

## 🚀 Швидкий старт

### 1. Initialize Terraform

```bash
terraform init
```

### 2. Plan

```bash
terraform plan
```

### 3. Apply

```bash
terraform apply
```

### 4. Send Test Message

```bash
# Get queue URL
QUEUE_URL=$(terraform output -raw queue_url)

# Send message
aws sqs send-message \
  --queue-url $QUEUE_URL \
  --message-body '{"order_id":"123","amount":100}'

# Send batch
aws sqs send-message-batch \
  --queue-url $QUEUE_URL \
  --entries file://messages.json
```

**messages.json:**
```json
[
  {
    "Id": "1",
    "MessageBody": "{\"order_id\":\"123\",\"amount\":100}"
  },
  {
    "Id": "2",
    "MessageBody": "{\"order_id\":\"456\",\"amount\":200}"
  }
]
```

### 5. Receive Messages

```bash
# Receive with long polling
aws sqs receive-message \
  --queue-url $QUEUE_URL \
  --max-number-of-messages 10 \
  --wait-time-seconds 20

# Response:
# {
#   "Messages": [
#     {
#       "MessageId": "abc123",
#       "ReceiptHandle": "xyz...",
#       "Body": "{\"order_id\":\"123\",\"amount\":100}"
#     }
#   ]
# }
```

### 6. Delete Message (ACK)

```bash
aws sqs delete-message \
  --queue-url $QUEUE_URL \
  --receipt-handle "RECEIPT_HANDLE_FROM_ABOVE"
```

### 7. Check DLQ

```bash
DLQ_URL=$(terraform output -raw dlq_url)

aws sqs receive-message \
  --queue-url $DLQ_URL \
  --max-number-of-messages 10
```

### 8. Monitor

```bash
# Check queue metrics
aws cloudwatch get-metric-statistics \
  --namespace AWS/SQS \
  --metric-name ApproximateNumberOfMessagesVisible \
  --dimensions Name=QueueName,Value=my-queue \
  --start-time $(date -u -d '1 hour ago' +%Y-%m-%dT%H:%M:%S) \
  --end-time $(date -u +%Y-%m-%dT%H:%M:%S) \
  --period 300 \
  --statistics Average

# Check alarms
aws cloudwatch describe-alarms --alarm-names my-queue-dlq-not-empty
```

### 9. Cleanup

```bash
terraform destroy
```

---

## 📝 Конфігурація

### Variables

**`terraform.tfvars`** (optional):
```hcl
aws_region              = "us-east-1"
queue_name              = "orders-queue"
visibility_timeout      = 300     # 5 minutes
message_retention_period = 345600  # 4 days
max_receive_count       = 3       # Retries before DLQ
```

### Queue Parameters

| Parameter | Default | Description |
|-----------|---------|-------------|
| `visibility_timeout` | 300s | How long message hidden after receive |
| `message_retention_seconds` | 345600s (4 days) | How long message stays in queue |
| `receive_wait_time_seconds` | 20s | Long polling wait time |
| `max_receive_count` | 3 | Retries before moving to DLQ |

### DLQ Parameters

| Parameter | Default | Description |
|-----------|---------|-------------|
| `message_retention_seconds` | 1209600s (14 days) | Longer retention for analysis |

---

## 📊 CloudWatch Alarms

### 1. DLQ Not Empty

```
Metric: ApproximateNumberOfMessagesVisible (DLQ)
Threshold: > 0
Action: Alert immediately
```

### 2. High Queue Depth

```
Metric: ApproximateNumberOfMessagesVisible (Main)
Threshold: > 1000
Action: Alert if sustained for 10 minutes
```

### 3. Old Messages

```
Metric: ApproximateAgeOfOldestMessage
Threshold: > 600 seconds (10 minutes)
Action: Alert if messages not processed
```

---

## 🧪 Testing

### Simulate Poison Message

```bash
# Send invalid message
aws sqs send-message \
  --queue-url $QUEUE_URL \
  --message-body "INVALID_JSON"

# Receive and DON'T delete (simulate failure)
for i in {1..3}; do
  echo "Attempt $i"
  aws sqs receive-message --queue-url $QUEUE_URL
  sleep 10  # Wait for visibility timeout
done

# Check DLQ
aws sqs receive-message --queue-url $DLQ_URL

# Should see the message moved to DLQ after 3 attempts
```

---

## ⚙️ Advanced

### FIFO Queue

```hcl
resource "aws_sqs_queue" "main" {
  name = "${var.queue_name}.fifo"
  
  fifo_queue                  = true
  content_based_deduplication = true
  deduplication_scope         = "messageGroup"
  fifo_throughput_limit       = "perMessageGroupId"
}
```

### Add IAM Policy

```hcl
resource "aws_sqs_queue_policy" "main" {
  queue_url = aws_sqs_queue.main.id
  
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
        Action   = "sqs:SendMessage"
        Resource = aws_sqs_queue.main.arn
      }
    ]
  })
}
```

### SNS to SQS Fan-Out

```hcl
resource "aws_sns_topic" "events" {
  name = "events-topic"
}

resource "aws_sns_topic_subscription" "queue" {
  topic_arn = aws_sns_topic.events.arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.main.arn
}

resource "aws_sqs_queue_policy" "sns" {
  queue_url = aws_sqs_queue.main.id
  
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "sns.amazonaws.com"
        }
        Action   = "sqs:SendMessage"
        Resource = aws_sqs_queue.main.arn
        Condition = {
          ArnEquals = {
            "aws:SourceArn" = aws_sns_topic.events.arn
          }
        }
      }
    ]
  })
}
```

---

## ✅ Перевірка

### 1. Queues Created

```bash
aws sqs list-queues

# Should show:
# - my-queue
# - my-queue-dlq
```

### 2. DLQ Configured

```bash
aws sqs get-queue-attributes \
  --queue-url $(terraform output -raw queue_url) \
  --attribute-names RedrivePolicy

# Should show maxReceiveCount=3
```

### 3. Alarms Created

```bash
aws cloudwatch describe-alarms | grep my-queue

# Should show 3 alarms
```

---

## 🎓 Висновок

Ти навчився:
- ✅ Створювати SQS queue з Terraform
- ✅ Налаштовувати Dead Letter Queue
- ✅ Включати long polling
- ✅ Створювати CloudWatch alarms
- ✅ Тестувати message flow

**Далі:** `03_iam_terraform/` - IAM roles та policies з Terraform
