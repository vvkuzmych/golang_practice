# Week 13 - Quick Start

## ⚡ Швидкий старт (5 хвилин)

### Prerequisites

```bash
# Install Terraform
brew install terraform

# Verify
terraform version

# Configure AWS
aws configure
```

---

## 🚀 Deploy Full Stack

### Option 1: Full Stack (Recommended)

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_13/practice/04_full_stack

# Build and deploy everything
make init
make deploy

# Send test messages
make send-batch

# Monitor logs
make logs

# Check metrics
make monitor

# Cleanup
make destroy
```

### Option 2: Step by Step

**1. Lambda Only**
```bash
cd practice/01_lambda_terraform
make init
make deploy
make test
make destroy
```

**2. SQS Only**
```bash
cd practice/02_sqs_terraform
make deploy
make send-batch
make receive
make destroy
```

**3. IAM Only**
```bash
cd practice/03_iam_terraform
make deploy
make outputs
make destroy
```

---

## 📊 Quick Test

### Deploy + Test + Destroy

```bash
cd practice/04_full_stack

# 1. Deploy
make deploy

# 2. Send 10 orders
make send-batch

# 3. Watch processing (in separate terminal)
make logs

# 4. Check metrics
make monitor

# 5. Test error handling
make send-invalid
sleep 120
make check-dlq

# 6. Load test (100 messages)
make test-load

# 7. Cleanup
make destroy
```

---

## 🔍 Key Commands

### Terraform

```bash
terraform init      # Initialize
terraform plan      # Preview changes
terraform apply     # Apply changes
terraform destroy   # Delete resources
terraform output    # Show outputs
```

### Makefiles

```bash
make help       # Show all commands
make deploy     # Build + deploy
make test       # Test functionality
make logs       # View logs
make monitor    # Show metrics
make destroy    # Cleanup
```

---

## 📖 Learn More

- [README.md](./README.md) - Full documentation
- [theory/](./theory/) - Terraform theory
- [practice/](./practice/) - Hands-on examples

---

## ✅ Expected Results

After `make deploy`:
```
✅ Full stack deployed!

Queue URL: https://sqs.us-east-1.amazonaws.com/123/order-processor-queue
Lambda: order-processor

Try: make send
```

After `make send-batch`:
```
Sending 10 orders...
Sent 1/10
Sent 2/10
...
✅ 10 messages sent! Check logs with: make logs
```

After `make logs`:
```
Processing batch of 10 messages
Processing message: abc123...
Processing order: ID=ORDER-001 Customer=CUST-123 Amount=$99.99
Order processed successfully: ORDER-001
...
```

---

## 🎓 Next Steps

1. ✅ Complete `practice/04_full_stack`
2. ✅ Modify Lambda code and redeploy
3. ✅ Change SQS parameters (batch size, timeout)
4. ✅ Add custom CloudWatch alarms
5. ✅ Implement DynamoDB for idempotency

**Week 13: Infrastructure as Code!** 🏗️
