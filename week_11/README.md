# Week 11: AWS Lambda & Serverless

## 🎯 Мета

Навчитися створювати **serverless applications** на Go з AWS Lambda, розуміти lifecycle, оптимізувати cold starts, та інтегруватися з API Gateway.

---

## 📚 Теорія

### 1. Lambda Lifecycle
**Файл:** `theory/01_lambda_lifecycle.md`

- Що таке AWS Lambda?
- Execution phases (INIT, INVOKE, SHUTDOWN)
- Cold start vs Warm start
- Container reuse
- Execution context
- Global variables persistence

**Key concepts:**
```
Cold Start = INIT + INVOKE (~100-500ms)
Warm Start = INVOKE only (~10-50ms)
```

### 2. Cold Start Optimization
**Файл:** `theory/02_cold_start_optimization.md`

- Cold start breakdown
- 7 optimization strategies
- Binary size reduction
- Lazy initialization
- Memory/CPU relationship
- Provisioned Concurrency
- Real-world benchmarks

**Typical improvements:**
```
Before: 1200ms cold start, 45 MB binary
After:   300ms cold start,  8 MB binary
4x faster! 5.6x smaller!
```

---

## 💻 Практика

### 1. Basic Lambda
**Директорія:** `practice/01_basic_lambda/`

**Файл:** `main.go`

**Features:**
- Basic Lambda handler
- Cold/warm start detection
- Lambda context usage
- Structured response
- Logging to CloudWatch

**How to deploy:**
```bash
cd practice/01_basic_lambda

# Build
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bootstrap main.go

# Package
zip function.zip bootstrap

# Deploy (using AWS CLI)
aws lambda create-function \
  --function-name basic-lambda \
  --runtime provided.al2023 \
  --handler bootstrap \
  --zip-file fileb://function.zip \
  --role arn:aws:iam::YOUR_ACCOUNT:role/lambda-role

# Invoke
aws lambda invoke \
  --function-name basic-lambda \
  --payload '{"name":"Alice","message":"Hello"}' \
  response.json

cat response.json
```

### 2. API Gateway Integration
**Директорія:** `practice/02_api_gateway/`

**Файл:** `main.go`

**Features:**
- API Gateway proxy integration
- HTTP routing (GET, POST)
- Path parameters
- Query parameters
- JSON request/response
- CORS headers

**API Routes:**
```
GET  /health        - Health check
GET  /users         - List users
GET  /users/{id}    - Get user by ID
POST /users         - Create user
```

**How to deploy:**
```bash
cd practice/02_api_gateway

# Build
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bootstrap main.go
zip function.zip bootstrap

# Deploy Lambda
aws lambda create-function \
  --function-name api-lambda \
  --runtime provided.al2023 \
  --handler bootstrap \
  --zip-file fileb://function.zip \
  --role arn:aws:iam::YOUR_ACCOUNT:role/lambda-role \
  --environment Variables="{STAGE=dev}"

# Create API Gateway (REST API)
# Use AWS Console or AWS SAM/Terraform
```

**Test locally:**
```bash
# Install SAM CLI
brew install aws-sam-cli

# Test locally
sam local start-api

# Call API
curl http://localhost:3000/health
curl http://localhost:3000/users
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com"}'
```

### 3. CloudWatch Logging
**Директорія:** `practice/03_cloudwatch/`

**Файл:** `main.go`

**Features:**
- Structured JSON logging
- Log levels (DEBUG, INFO, WARN, ERROR)
- Request tracking (requestId)
- Duration metrics
- Error logging with context

**Log format:**
```json
{
  "level": "INFO",
  "message": "Request completed successfully",
  "timestamp": "2026-01-28T12:00:00Z",
  "requestId": "abc-123",
  "function": "my-lambda",
  "extra": {
    "statusCode": 200,
    "duration": 150
  }
}
```

**CloudWatch Insights queries:**
```sql
-- Find cold starts
fields @timestamp, extra.requestId
| filter message = "Cold start detected"
| stats count() as coldStarts by bin(5m)

-- Average duration
fields @timestamp, extra.duration
| filter message = "Request completed successfully"
| stats avg(extra.duration) as avgDuration

-- Error rate
fields @timestamp, level
| filter level = "ERROR"
| stats count() as errors by bin(1h)

-- P99 latency
fields @timestamp, extra.duration
| filter message = "Request completed successfully"
| stats percentile(extra.duration, 99) as p99
```

---

## 🔧 Основні команди

### Build Lambda

```bash
# For AWS Lambda (Linux)
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bootstrap main.go

# Package
zip function.zip bootstrap

# Check size
ls -lh function.zip
```

### Deploy Lambda

```bash
# Create function
aws lambda create-function \
  --function-name my-function \
  --runtime provided.al2023 \
  --handler bootstrap \
  --zip-file fileb://function.zip \
  --role arn:aws:iam::ACCOUNT:role/lambda-role

# Update function code
aws lambda update-function-code \
  --function-name my-function \
  --zip-file fileb://function.zip

# Update configuration
aws lambda update-function-configuration \
  --function-name my-function \
  --memory-size 1024 \
  --timeout 30
```

### Invoke Lambda

```bash
# Sync invoke
aws lambda invoke \
  --function-name my-function \
  --payload '{"key":"value"}' \
  response.json

# Async invoke
aws lambda invoke \
  --function-name my-function \
  --invocation-type Event \
  --payload '{"key":"value"}' \
  response.json
```

### View Logs

```bash
# Get log groups
aws logs describe-log-groups \
  --log-group-name-prefix /aws/lambda/my-function

# Tail logs
aws logs tail /aws/lambda/my-function --follow

# Get recent logs
aws logs tail /aws/lambda/my-function \
  --since 10m \
  --format short
```

---

## 📊 Lambda Configuration

### Memory & Timeout

```yaml
# lambda-config.yaml
Memory: 1024      # 128 MB - 10240 MB
Timeout: 30       # Max 900s (15 min)
```

**Memory → CPU:**
```
128 MB   → 0.08 vCPU
1024 MB  → 0.58 vCPU  ← Recommended
1792 MB  → 1.00 vCPU  ← Sweet spot
3008 MB  → 1.75 vCPU
```

### Environment Variables

```bash
aws lambda update-function-configuration \
  --function-name my-function \
  --environment Variables="{DB_HOST=localhost,API_KEY=secret}"
```

**In code:**
```go
func init() {
    dbHost := os.Getenv("DB_HOST")
    apiKey := os.Getenv("API_KEY")
}
```

---

## 🎯 Best Practices

### Build Optimization

```bash
# Strip debug info
-ldflags="-s -w"

# Remove file paths
-trimpath

# Full command
GOOS=linux GOARCH=amd64 go build \
  -ldflags="-s -w" \
  -trimpath \
  -o bootstrap main.go
```

### Code Optimization

```go
// ✅ DO: Init heavy resources once
var (
    db     *sql.DB
    config *Config
)

func init() {
    db = connectDB()
    config = loadConfig()
}

// ✅ DO: Reuse connections
var httpClient = &http.Client{
    Timeout: 10 * time.Second,
}

// ❌ DON'T: Init per request
func handler(ctx context.Context) error {
    db := connectDB()  // ❌ Slow!
    defer db.Close()
    // ...
}
```

### Logging Best Practices

```go
// ✅ DO: Structured logging
log.Printf(`{"level":"INFO","message":"Request processed","requestId":"%s"}`, requestID)

// ✅ DO: Include context
logInfo("Processing request", map[string]interface{}{
    "requestId": requestID,
    "path":      path,
    "method":    method,
})

// ❌ DON'T: Plain text logs
log.Println("Processing request")  // Hard to query
```

---

## 📈 Monitoring & Debugging

### CloudWatch Metrics

**Default metrics:**
- **Invocations** - Total calls
- **Duration** - Execution time
- **Errors** - Failed invocations
- **Throttles** - Rate limit exceeded
- **ConcurrentExecutions** - Active containers

**Custom metrics:**
```go
import "github.com/aws/aws-sdk-go-v2/service/cloudwatch"

// Put custom metric
client.PutMetricData(ctx, &cloudwatch.PutMetricDataInput{
    Namespace: aws.String("MyApp"),
    MetricData: []types.MetricDatum{
        {
            MetricName: aws.String("ProcessingTime"),
            Value:      aws.Float64(duration),
            Unit:       types.StandardUnitMilliseconds,
        },
    },
})
```

### CloudWatch Alarms

```bash
# P99 latency alarm
aws cloudwatch put-metric-alarm \
  --alarm-name high-latency \
  --metric-name Duration \
  --namespace AWS/Lambda \
  --statistic Average \
  --period 300 \
  --threshold 1000 \
  --comparison-operator GreaterThanThreshold

# Error rate alarm
aws cloudwatch put-metric-alarm \
  --alarm-name high-errors \
  --metric-name Errors \
  --namespace AWS/Lambda \
  --statistic Sum \
  --period 60 \
  --threshold 10 \
  --comparison-operator GreaterThanThreshold
```

### X-Ray Tracing

```go
import "github.com/aws/aws-xray-sdk-go/xray"

func handler(ctx context.Context) error {
    // Create subsegment
    ctx, seg := xray.BeginSubsegment(ctx, "ProcessData")
    defer seg.Close(nil)
    
    // Your code
    result := processData()
    
    return nil
}
```

---

## 🚀 Quick Start

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_11

# Read theory
cat theory/01_lambda_lifecycle.md
cat theory/02_cold_start_optimization.md

# Try basic example
cd practice/01_basic_lambda

# Build
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bootstrap main.go
zip function.zip bootstrap

# Test locally (requires AWS SAM)
sam local invoke -e test-event.json

# Deploy to AWS
aws lambda create-function \
  --function-name test-lambda \
  --runtime provided.al2023 \
  --handler bootstrap \
  --zip-file fileb://function.zip \
  --role YOUR_LAMBDA_ROLE_ARN
```

---

## 🎓 Learning Path

### День 1: Lambda Basics

1. Читай `theory/01_lambda_lifecycle.md`
2. Розумій INIT → INVOKE → SHUTDOWN
3. Cold start vs warm start
4. Deploy `practice/01_basic_lambda/`

### День 2: Cold Start Optimization

1. Читай `theory/02_cold_start_optimization.md`
2. Optimize binary size
3. Test different memory configs
4. Measure cold start time

### День 3: API Gateway

1. Deploy `practice/02_api_gateway/`
2. Create API Gateway REST API
3. Test all endpoints
4. Add authentication

### День 4: CloudWatch

1. Implement structured logging
2. Create CloudWatch dashboards
3. Set up alarms
4. Query logs with Insights

---

## 📖 Ресурси

### AWS Documentation

- [AWS Lambda Go](https://docs.aws.amazon.com/lambda/latest/dg/lambda-golang.html)
- [Lambda Best Practices](https://docs.aws.amazon.com/lambda/latest/dg/best-practices.html)
- [CloudWatch Logs Insights](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/AnalyzingLogData.html)

### Libraries

- `github.com/aws/aws-lambda-go` - Official Lambda SDK
- `github.com/aws/aws-sdk-go-v2` - AWS SDK v2
- `github.com/aws/aws-xray-sdk-go` - X-Ray tracing

### Tools

- [AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html) - Local testing
- [Serverless Framework](https://www.serverless.com/) - Deployment framework
- [AWS CDK](https://aws.amazon.com/cdk/) - Infrastructure as Code

---

**"Serverless = Focus on code, not infrastructure!" ☁️**

**Status:** Week 11 Materials Complete ✅  
**Created:** 2026-01-28
