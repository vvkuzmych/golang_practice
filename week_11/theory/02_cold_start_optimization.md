# Cold Start Optimization

## 🎯 Що таке Cold Start?

**Cold Start** - це час, необхідний AWS Lambda для ініціалізації нового execution environment.

```
Cold Start = Download + Runtime Init + Code Init + First Invocation

Go Lambda: ~100-500ms (typical)
```

---

## 📊 Cold Start Breakdown

### Типовий Cold Start (Go)

```
Phase 1: Download code        100ms  ████████
Phase 2: Start runtime         50ms  ████
Phase 3: Init code            200ms  ████████████████
Phase 4: Handler execution     50ms  ████
                              ----
Total:                        400ms
```

### Warm Start (для порівняння)

```
Handler execution only         50ms  ████
```

**Warm start 8x faster!**

---

## 🔧 Optimization Strategies

### Strategy 1: Reduce Package Size

#### Problem

```bash
# Large binary
-rw-r--r--  1 user  staff  45M Jan 28 main

# Slow download (200-500ms)
```

#### Solution: Build Flags

```bash
# Strip debug info and symbol table
GOOS=linux GOARCH=amd64 go build \
  -ldflags="-s -w" \
  -trimpath \
  -o bootstrap main.go

# Result
-rw-r--r--  1 user  staff  8M Jan 28 bootstrap
```

**Flags:**
- `-s`: Strip symbol table
- `-w`: Strip DWARF debug info
- `-trimpath`: Remove file system paths

**Improvement:** 45 MB → 8 MB (5.6x smaller, ~300ms saved)

#### Solution: UPX Compression

```bash
# Further compress
upx --best --lzma bootstrap

# Result
-rw-r--r--  1 user  staff  3M Jan 28 bootstrap
```

**⚠️ Warning:** UPX adds decompression time (~50ms), test first!

---

### Strategy 2: Optimize Dependencies

#### Problem: Too Many Dependencies

```go
// ❌ BAD: Import entire AWS SDK
import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/sqs"
    // ... entire SDK = 100+ MB
)
```

#### Solution: Import Only What You Need

```go
// ✅ GOOD: Import only S3
import (
    "github.com/aws/aws-sdk-go-v2/service/s3"
)
```

**Improvement:** Binary size: 45 MB → 12 MB

#### Solution: Use AWS SDK v2

```go
// ✅ Better: Modular SDK
import (
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)

// Smaller binary, faster init
```

---

### Strategy 3: Lazy Initialization

#### Problem: Load Everything in init()

```go
func init() {
    // ❌ BAD: Load even if not needed
    s3Client = createS3Client()
    dbClient = createDBClient()
    cacheClient = createCacheClient()
    // All loaded on cold start!
}
```

#### Solution: Lazy Load

```go
var (
    s3Client    *s3.Client
    s3Once      sync.Once
)

func getS3Client() *s3.Client {
    s3Once.Do(func() {
        s3Client = createS3Client()
    })
    return s3Client
}

func handler(ctx context.Context) error {
    // Only load if needed
    if needsS3 {
        client := getS3Client()
        // Use client
    }
    // ...
}
```

**Improvement:** Init time reduced by 50-70% (if not all services used)

---

### Strategy 4: Increase Memory (= More CPU)

#### Problem: Slow CPU

```yaml
# Configuration
Memory: 128 MB
# Gets: 0.08 vCPU (slow!)
```

**Cold start:** 800ms

#### Solution: Increase Memory

```yaml
# Configuration
Memory: 1024 MB
# Gets: 0.58 vCPU (5x faster!)
```

**Cold start:** 200ms

**Improvement:** 4x faster cold start

**Cost trade-off:**
```
128 MB:  $0.0000000021 per ms
1024 MB: $0.0000000167 per ms (8x more)

But: Executes 4x faster
Net: 2x cost for 4x performance ← Worth it!
```

---

### Strategy 5: Connection Pooling

#### Problem: Reconnect on Every Request

```go
func handler(ctx context.Context) error {
    // ❌ BAD: New connection per invocation
    db, err := sql.Open("postgres", dsn)
    defer db.Close()
    
    // Use db
    return nil
}
```

**Overhead:** 100-300ms per request (even warm starts!)

#### Solution: Reuse Connections

```go
var db *sql.DB

func init() {
    var err error
    db, err = sql.Open("postgres", os.Getenv("DB_DSN"))
    if err != nil {
        panic(err)
    }
    
    // Configure pool
    db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)
}

func handler(ctx context.Context) error {
    // ✅ GOOD: Reuse connection
    row := db.QueryRowContext(ctx, "SELECT ...")
    // ...
    return nil
}
```

**Improvement:**
- Cold start: +100ms (one-time connection)
- Warm start: 0ms overhead (reuse connection)

---

### Strategy 6: Avoid VPC (if possible)

#### Problem: VPC ENI Creation

```yaml
# Lambda in VPC
VpcConfig:
  SecurityGroupIds:
    - sg-xxx
  SubnetIds:
    - subnet-xxx
```

**Cold start penalty:** +600-700ms (ENI creation)

#### Solution: No VPC (if possible)

```yaml
# Lambda outside VPC
# (can still access public AWS services)
```

**Improvement:** 600-700ms saved

**Alternative:** Use VPC endpoints or NAT if VPC required

---

### Strategy 7: Provisioned Concurrency

#### Problem: Cold Starts on Traffic Spikes

```
Traffic spike → Need 100 containers → 100 cold starts!
```

#### Solution: Provisioned Concurrency

```yaml
ProvisionedConcurrencyConfig:
  ProvisionedConcurrentExecutions: 10
```

**How it works:**
- AWS keeps 10 containers always warm
- No cold start for first 10 concurrent requests
- Additional requests still cold start

**Cost:**
```
$0.000004166 per GB-second
For 1024 MB, 10 containers = $1.08/hour

Only use for production APIs!
```

---

## 📊 Benchmarking Cold Starts

### Measure Cold Start in Code

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-lambda-go/lambdacontext"
)

var (
    initStart time.Time
    coldStart = true
)

func init() {
    initStart = time.Now()
    
    // Heavy init work here
    time.Sleep(200 * time.Millisecond)
    
    fmt.Printf("Init took: %v\n", time.Since(initStart))
}

func handler(ctx context.Context, event map[string]interface{}) (map[string]interface{}, error) {
    lc, _ := lambdacontext.FromContext(ctx)
    
    if coldStart {
        fmt.Printf("COLD START - RequestID: %s\n", lc.RequestID)
        fmt.Printf("Total init time: %v\n", time.Since(initStart))
        coldStart = false
    } else {
        fmt.Printf("WARM START - RequestID: %s\n", lc.RequestID)
    }
    
    // Handler logic
    return map[string]interface{}{
        "statusCode": 200,
        "body":       "OK",
    }, nil
}

func main() {
    lambda.Start(handler)
}
```

### CloudWatch Logs Analysis

```sql
-- Find cold starts
fields @timestamp, @message, @duration
| filter @message like /COLD START/
| stats count() as coldStarts, avg(@duration) as avgDuration

-- Cold start percentage
fields @timestamp
| filter @type = "REPORT"
| stats 
    count(*) as total,
    sum(@initDuration > 0) as coldStarts
| extend coldStartPercentage = (coldStarts / total) * 100
```

---

## 🎯 Optimization Checklist

### Build Optimization

- [ ] Use `-ldflags="-s -w"` to strip debug info
- [ ] Use `-trimpath` to remove file paths
- [ ] Consider UPX compression (test first!)
- [ ] Remove unused dependencies
- [ ] Use AWS SDK v2 (modular)

### Code Optimization

- [ ] Move heavy init to `init()`
- [ ] Use lazy loading for optional services
- [ ] Reuse connections (DB, HTTP)
- [ ] Minimize init() complexity
- [ ] Cache data in `/tmp` if needed

### Configuration Optimization

- [ ] Increase memory (1024 MB sweet spot)
- [ ] Avoid VPC if possible
- [ ] Use Provisioned Concurrency for production APIs
- [ ] Set appropriate timeout (don't use 15 min default!)

### Monitoring

- [ ] Track cold start percentage
- [ ] Monitor init duration
- [ ] Set CloudWatch alarms for P99 latency
- [ ] Use X-Ray for detailed tracing

---

## 📈 Real-World Example

### Before Optimization

```go
// main.go
import (
    "github.com/aws/aws-sdk-go/aws"           // 100 MB
    "github.com/gin-gonic/gin"                 // 20 MB (overkill!)
    // ... many deps
)

func init() {
    // Load everything
    s3Client = createS3()
    dbClient = createDB()
    cacheClient = createCache()
}

func handler(ctx context.Context) error {
    // Sometimes uses S3, sometimes DB
    // ...
}
```

**Results:**
```
Binary size:  45 MB
Cold start:   1200ms
Warm start:   100ms
Memory:       128 MB
Cold start %: 15%
```

### After Optimization

```go
// main.go
import (
    "github.com/aws/aws-sdk-go-v2/service/s3"  // 10 MB
    "github.com/aws/aws-lambda-go/lambda"
)

var (
    s3Client *s3.Client
    s3Once   sync.Once
)

func getS3Client() *s3.Client {
    s3Once.Do(func() {
        s3Client = createS3()
    })
    return s3Client
}

func handler(ctx context.Context) error {
    // Lazy load
    if needsS3 {
        client := getS3Client()
    }
    // ...
}
```

**Build:**
```bash
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bootstrap
zip function.zip bootstrap
```

**Results:**
```
Binary size:  8 MB (5.6x smaller)
Cold start:   300ms (4x faster)
Warm start:   50ms (2x faster)
Memory:       1024 MB (8x more)
Cold start %: 5% (3x less)

Cost: +50% (but 4x faster!)
```

---

## ✅ Best Practices Summary

### Do's ✅

1. **Build with optimization flags**
   ```bash
   -ldflags="-s -w" -trimpath
   ```

2. **Minimize dependencies**
   - Use AWS SDK v2
   - Remove unused imports

3. **Reuse connections**
   - DB connections in `init()`
   - HTTP clients with connection pooling

4. **Increase memory for CPU**
   - 1024 MB sweet spot
   - Balance cost vs performance

5. **Use Provisioned Concurrency** (production)
   - For critical APIs
   - Eliminate cold starts

### Don'ts ❌

1. **Don't use VPC** (unless required)
   - +600ms cold start penalty

2. **Don't import entire SDK**
   - Import only what you need

3. **Don't init unused services**
   - Use lazy loading

4. **Don't use 128 MB** (default)
   - Too slow CPU

5. **Don't ignore cold starts**
   - Monitor and optimize

---

## 🎓 Висновок

### Cold Start Optimization:

✅ **Binary size:** Reduce with build flags (5x smaller)  
✅ **Dependencies:** Use AWS SDK v2 (modular)  
✅ **Memory:** Increase to 1024 MB (4x faster)  
✅ **Lazy loading:** Initialize only what's needed  
✅ **Provisioned Concurrency:** Eliminate cold starts (prod)  

### Typical Improvements:

- **Binary:** 45 MB → 8 MB (5.6x)
- **Cold start:** 1200ms → 300ms (4x)
- **Warm start:** 100ms → 50ms (2x)
- **Cold start %:** 15% → 5% (3x less)

### Golden Rule:

**"Optimize for warm start latency, tolerate cold start latency, or pay for Provisioned Concurrency!"**

---

## 📖 Далі

- `practice/01_basic_lambda/` - Basic Lambda implementation
- `practice/02_api_gateway/` - Lambda + API Gateway
- `practice/03_cloudwatch/` - Logging and monitoring

**"Every millisecond counts in serverless!" ⚡☁️**
