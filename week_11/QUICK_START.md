# Week 11 - Quick Start ☁️

## 🎯 Мета

Навчитися створювати **AWS Lambda functions** на Go з оптимізацією cold starts.

---

## ⚡ 5-хвилинний старт

### 1. Create Basic Lambda

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_11/practice/01_basic_lambda

# Build for Linux (Lambda runtime)
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bootstrap main.go

# Package
zip function.zip bootstrap

# Check size
ls -lh function.zip
```

### 2. Test Locally (Optional)

```bash
# Install AWS SAM CLI
brew install aws-sam-cli

# Create test event
cat > test-event.json << 'JSON'
{
  "name": "Alice",
  "message": "Hello Lambda"
}
JSON

# Test locally
sam local invoke -e test-event.json
```

### 3. Deploy to AWS

```bash
# Deploy (replace YOUR_ROLE_ARN)
aws lambda create-function \
  --function-name my-first-lambda \
  --runtime provided.al2023 \
  --handler bootstrap \
  --zip-file fileb://function.zip \
  --role arn:aws:iam::ACCOUNT:role/lambda-role \
  --memory-size 1024 \
  --timeout 30

# Invoke
aws lambda invoke \
  --function-name my-first-lambda \
  --payload '{"name":"Alice","message":"Test"}' \
  response.json

cat response.json
```

---

## 📚 Читати теорію

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_11

# Lambda Lifecycle
cat theory/01_lambda_lifecycle.md

# Cold Start Optimization
cat theory/02_cold_start_optimization.md
```

---

## 🎯 2 Головні Концепції

### 1. Lambda Lifecycle

```
┌───────────────────────────────────┐
│  COLD START (first invocation)   │
├───────────────────────────────────┤
│  1. Download Code     100ms       │
│  2. Start Runtime      50ms       │
│  3. Init Code         200ms       │
│  4. Handler            50ms       │
│                                   │
│  Total: ~400ms                    │
└───────────────────────────────────┘

┌───────────────────────────────────┐
│  WARM START (reuse container)    │
├───────────────────────────────────┤
│  1. Handler            50ms       │
│                                   │
│  Total: ~50ms (8x faster!)        │
└───────────────────────────────────┘
```

### 2. Cold Start Optimization

```go
// ✅ DO: Heavy init once
var (
    db     *sql.DB
    config *Config
)

func init() {
    // Runs once during cold start
    db = connectDB()
    config = loadConfig()
}

// ✅ DO: Light handler
func handler(ctx context.Context) error {
    // Reuse db (already connected)
    result := db.Query(...)
    return nil
}
```

---

## 🔧 Essential Commands

### Build & Package

```bash
# Optimized build
GOOS=linux GOARCH=amd64 go build \
  -ldflags="-s -w" \
  -trimpath \
  -o bootstrap main.go

# Package
zip function.zip bootstrap
```

### Deploy & Update

```bash
# Create
aws lambda create-function \
  --function-name my-func \
  --runtime provided.al2023 \
  --handler bootstrap \
  --zip-file fileb://function.zip \
  --role ROLE_ARN

# Update code
aws lambda update-function-code \
  --function-name my-func \
  --zip-file fileb://function.zip

# Update config
aws lambda update-function-configuration \
  --function-name my-func \
  --memory-size 1024
```

### Invoke & Logs

```bash
# Invoke
aws lambda invoke \
  --function-name my-func \
  --payload '{"key":"value"}' \
  response.json

# View logs
aws logs tail /aws/lambda/my-func --follow

# Get recent logs
aws logs tail /aws/lambda/my-func --since 10m
```

---

## 📊 Optimization Wins

### Win 1: Build Flags

```bash
# Before
go build -o main main.go
# Size: 15 MB

# After
go build -ldflags="-s -w" -o bootstrap main.go
# Size: 5 MB

Improvement: 3x smaller!
```

### Win 2: Memory Config

```yaml
# Before
Memory: 128 MB
CPU: 0.08 vCPU
Cold start: 800ms

# After
Memory: 1024 MB
CPU: 0.58 vCPU (8x more)
Cold start: 200ms

Improvement: 4x faster!
Cost: Only 2x more (worth it!)
```

### Win 3: Connection Reuse

```go
// ❌ Before: New connection each time
func handler(ctx context.Context) error {
    db := connectDB()  // 200ms
    defer db.Close()
    // ...
}

// ✅ After: Reuse connection
var db *sql.DB

func init() {
    db = connectDB()  // Once!
}

func handler(ctx context.Context) error {
    // Reuse db (0ms overhead)
    // ...
}

Improvement: 200ms saved per warm invocation!
```

---

## 🎯 Quick Patterns

### Pattern 1: Basic Handler

```go
package main

import (
    "context"
    "github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
    Name string `json:"name"`
}

type Response struct {
    Message string `json:"message"`
}

func handler(ctx context.Context, event Event) (Response, error) {
    return Response{
        Message: "Hello " + event.Name,
    }, nil
}

func main() {
    lambda.Start(handler)
}
```

### Pattern 2: API Gateway

```go
import (
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    return events.APIGatewayProxyResponse{
        StatusCode: 200,
        Body:       `{"message":"Hello"}`,
    }, nil
}

func main() {
    lambda.Start(handler)
}
```

### Pattern 3: Cold Start Detection

```go
var coldStart = true

func handler(ctx context.Context) error {
    if coldStart {
        log.Println("COLD START")
        coldStart = false
    } else {
        log.Println("WARM START")
    }
    // ...
}
```

---

## 📖 Структура

```
week_11/
├── README.md                    # Повний опис
├── QUICK_START.md              # Цей файл
├── theory/
│   ├── 01_lambda_lifecycle.md  # Lifecycle & phases
│   └── 02_cold_start_optimization.md
└── practice/
    ├── 01_basic_lambda/        # Basic example
    ├── 02_api_gateway/         # REST API
    └── 03_cloudwatch/          # Structured logging
```

---

## ✅ Quick Checklist

### Before Deploying

- [ ] Build with optimization flags
- [ ] Package size < 10 MB
- [ ] Set memory to 1024 MB (or test)
- [ ] Init heavy resources in `init()`
- [ ] Add structured logging

### After Deploying

- [ ] Test cold start time
- [ ] Monitor CloudWatch logs
- [ ] Check error rate
- [ ] Set up alarms
- [ ] Monitor costs

---

## 🎓 Learning Path

### День 1: Basics
1. Read `theory/01_lambda_lifecycle.md`
2. Deploy `practice/01_basic_lambda/`
3. Test cold vs warm start

### День 2: Optimization
1. Read `theory/02_cold_start_optimization.md`
2. Optimize binary size
3. Test different memory configs
4. Measure improvements

### День 3: API Gateway
1. Deploy `practice/02_api_gateway/`
2. Create REST API
3. Test all endpoints

### День 4: Monitoring
1. Implement structured logging
2. Query logs with CloudWatch Insights
3. Set up alarms

---

## 🎯 Golden Rules

### Rule 1: Build Optimized

```bash
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bootstrap
```

### Rule 2: Init Heavy, Handler Light

```go
// Heavy init (once)
func init() {
    db = connectDB()
    config = loadConfig()
}

// Light handler (per request)
func handler(ctx context.Context) error {
    // Reuse db
}
```

### Rule 3: Use 1024 MB Memory

```yaml
Memory: 1024 MB  # Sweet spot
```

### Rule 4: Structured Logging

```go
log.Printf(`{"level":"INFO","message":"...","requestId":"%s"}`, id)
```

### Rule 5: Monitor Cold Starts

```sql
-- CloudWatch Insights
fields @timestamp
| filter @message like /COLD START/
| stats count() as coldStarts
```

---

## 🚀 Quick Test

```bash
# 1. Navigate
cd /Users/vkuzm/GolandProjects/golang_practice/week_11/practice/01_basic_lambda

# 2. Build
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bootstrap main.go

# 3. Package
zip function.zip bootstrap

# 4. Deploy (replace ROLE_ARN)
aws lambda create-function \
  --function-name quick-test \
  --runtime provided.al2023 \
  --handler bootstrap \
  --zip-file fileb://function.zip \
  --role arn:aws:iam::ACCOUNT:role/lambda-role

# 5. Invoke
aws lambda invoke \
  --function-name quick-test \
  --payload '{"name":"Test","message":"Hello"}' \
  out.json && cat out.json

# 6. Logs
aws logs tail /aws/lambda/quick-test --follow
```

---

**"Cold start is normal, warm start is fast!" ❄️⚡☁️**
