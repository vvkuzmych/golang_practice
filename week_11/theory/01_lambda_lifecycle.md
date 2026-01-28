# AWS Lambda Lifecycle

## 🎯 Що таке AWS Lambda?

**AWS Lambda** - це serverless compute service, що запускає код у відповідь на events без управління серверами.

```
Event → Lambda Function → Response
```

**Key benefits:**
- ✅ No server management
- ✅ Auto-scaling
- ✅ Pay per invocation
- ✅ High availability

---

## 📊 Lambda Lifecycle

### Повний цикл Lambda execution

```
1. INIT (cold start)
   ├─ Download code
   ├─ Start runtime
   ├─ Run init code
   └─ Create execution environment

2. INVOKE
   ├─ Run handler function
   └─ Return response

3. SHUTDOWN (after idle)
   └─ Destroy environment
```

---

## 🔄 Execution Phases

### Phase 1: Init (Cold Start)

```
Event arrives → No warm container → INIT phase
```

**Steps:**
1. **Download code** from S3 (~100-200ms)
2. **Start runtime** (Go runtime ~50-100ms)
3. **Run init code** (outside handler)
4. **Create execution environment**

**Duration:** 100ms - 5 seconds (depending on code size)

### Phase 2: Invoke

```
Handler function executes
```

**Steps:**
1. Call handler function
2. Execute business logic
3. Return response

**Duration:** Your code execution time

### Phase 3: Freeze (After Response)

```
Response sent → Container frozen → Keep for reuse
```

**AWS keeps container warm for:**
- 5-15 minutes (typical)
- Can vary based on load

---

## ❄️ Cold Start vs Warm Start

### Cold Start

```
Event → No container → INIT + INVOKE
```

**Latency:** 100ms - 5 seconds

**Happens when:**
- First invocation
- After container expired (~15 min idle)
- Traffic spike (need more containers)
- Code deployment

### Warm Start

```
Event → Warm container exists → INVOKE only
```

**Latency:** < 10ms (just handler execution)

**Happens when:**
- Container already initialized
- Within 5-15 min of last invocation

---

## 🏗️ Go Lambda Structure

### Basic Lambda Handler

```go
package main

import (
    "context"
    "github.com/aws/aws-lambda-go/lambda"
)

// Init code (runs once during cold start)
var (
    dbConnection *DB
    config       *Config
)

func init() {
    // Runs during INIT phase (cold start only)
    dbConnection = connectDB()
    config = loadConfig()
}

// Handler (runs on every invocation)
func handler(ctx context.Context, event MyEvent) (MyResponse, error) {
    // Business logic
    // Uses dbConnection (already initialized)
    return MyResponse{Message: "Hello"}, nil
}

func main() {
    lambda.Start(handler)
}
```

### Execution Flow

```
Cold Start:
1. init() runs        ← 500ms
2. handler() runs     ← 50ms
Total: 550ms

Warm Start:
1. handler() runs     ← 50ms
Total: 50ms
```

---

## 📊 Lambda Lifecycle Diagram

```
┌─────────────────────────────────────────────────┐
│  COLD START (first invocation)                  │
├─────────────────────────────────────────────────┤
│                                                  │
│  1. Download Code (100-200ms)                   │
│     └─ Fetch .zip from S3                       │
│                                                  │
│  2. Start Runtime (50-100ms)                    │
│     └─ Boot Go runtime                          │
│                                                  │
│  3. Init Code (variable)                        │
│     ├─ Run init()                               │
│     ├─ Connect to DB                            │
│     └─ Load config                              │
│                                                  │
│  4. Invoke Handler                              │
│     └─ Run handler(event)                       │
│                                                  │
└─────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────┐
│  WARM START (subsequent invocations)            │
├─────────────────────────────────────────────────┤
│                                                  │
│  1. Invoke Handler (fast!)                      │
│     └─ Reuse existing container                 │
│                                                  │
└─────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────┐
│  SHUTDOWN (after idle ~15 min)                  │
├─────────────────────────────────────────────────┤
│                                                  │
│  1. Container destroyed                         │
│  2. Next invocation = cold start again          │
│                                                  │
└─────────────────────────────────────────────────┘
```

---

## 🎯 Optimizing Init Code

### What to Put in init()

```go
func init() {
    // ✅ DO: Heavy initialization (runs once)
    
    // Database connections
    db = connectToDB()
    
    // HTTP clients (reuse connections)
    httpClient = &http.Client{
        Timeout: 10 * time.Second,
        Transport: &http.Transport{
            MaxIdleConns: 100,
        },
    }
    
    // Config loading
    config = loadFromSSM()
    
    // AWS SDK clients
    s3Client = s3.New(session.Must(session.NewSession()))
}
```

### What NOT to Put in init()

```go
func init() {
    // ❌ DON'T: Request-specific logic
    
    // Current time (changes per request)
    // now := time.Now()
    
    // Random values (should be per-request)
    // requestID := uuid.New()
    
    // Request data
    // user := getCurrentUser()
}
```

---

## 📈 Cold Start Factors

### What Affects Cold Start Time?

1. **Code Package Size**
   ```
   10 MB  → ~100ms download
   50 MB  → ~500ms download
   250 MB → ~2s download
   ```

2. **Runtime**
   ```
   Go     → ~50-100ms  ✅ Fastest
   Node   → ~100-200ms
   Python → ~150-300ms
   Java   → ~3-10s     ❌ Slowest
   ```

3. **Memory Configuration**
   ```
   128 MB  → Slower CPU
   1024 MB → Faster CPU
   3008 MB → Fastest CPU
   
   More memory = More CPU = Faster init!
   ```

4. **VPC Configuration**
   ```
   No VPC → ~100ms
   With VPC → +600-700ms (ENI creation)
   ```

5. **Init Code Complexity**
   ```
   Simple init()  → ~10ms
   DB connection  → ~100-300ms
   Heavy compute  → 1-5s
   ```

---

## 🔧 Lambda Configuration

### Memory & CPU

```yaml
# AWS Lambda configuration
Memory: 1024 MB  # Also determines CPU!
Timeout: 30s     # Max execution time
```

**Memory → CPU relationship:**
```
128 MB   → 0.08 vCPU
512 MB   → 0.33 vCPU
1024 MB  → 0.58 vCPU
1792 MB  → 1.00 vCPU  ← Sweet spot
3008 MB  → 1.75 vCPU  ← Max
```

### Environment Variables

```go
func init() {
    // Read env vars during init
    dbHost := os.Getenv("DB_HOST")
    apiKey := os.Getenv("API_KEY")
}
```

---

## 📊 Container Reuse

### How Long Does AWS Keep Container?

```
Invocation 1 (t=0s)     → Cold start
Invocation 2 (t=10s)    → Warm ✅
Invocation 3 (t=60s)    → Warm ✅
Invocation 4 (t=300s)   → Warm ✅
Invocation 5 (t=900s)   → Warm ✅
Invocation 6 (t=1800s)  → Cold ❄️ (15+ min idle)
```

**Typical:** 5-15 minutes

### Multiple Concurrent Invocations

```
Request 1 → Container A (cold start)
Request 2 → Container B (cold start)  ← New container!
Request 3 → Container C (cold start)  ← New container!
Request 4 → Container A (warm)        ← Reuse!
```

**Each concurrent execution = separate container!**

---

## 🎯 Execution Context Reuse

### Global Variables Persist

```go
var counter int  // Persists across warm invocations!

func handler(ctx context.Context) (Response, error) {
    counter++  // 1, 2, 3... on warm starts
    
    fmt.Printf("Invocation count: %d\n", counter)
    
    return Response{Count: counter}, nil
}
```

**Output:**
```
Cold start: counter = 1
Warm:       counter = 2
Warm:       counter = 3
...
Cold start: counter = 1  (new container)
```

### /tmp Directory Persists

```go
func init() {
    // Download large file once
    downloadToTmp("/tmp/model.dat")
}

func handler(ctx context.Context) error {
    // Reuse file from /tmp (warm start)
    data, err := ioutil.ReadFile("/tmp/model.dat")
    // ...
}
```

**Important:**
- `/tmp` has 512 MB limit
- Persists only within same container
- Clean up if needed

---

## ⚡ Provisioned Concurrency

### Eliminate Cold Starts

```yaml
# Keep N containers always warm
ProvisionedConcurrency: 5
```

**How it works:**
```
AWS keeps 5 containers always initialized

Request arrives → Use pre-warmed container → No cold start!
```

**Cost:**
- Pay for all provisioned containers (running or not)
- Good for production APIs
- Expensive for low-traffic functions

---

## 🔍 Measuring Cold Starts

### CloudWatch Logs

```go
import "github.com/aws/aws-lambda-go/lambdacontext"

var coldStart = true

func handler(ctx context.Context, event Event) (Response, error) {
    lc, _ := lambdacontext.FromContext(ctx)
    
    if coldStart {
        log.Printf("COLD START - RequestID: %s", lc.RequestID)
        coldStart = false
    } else {
        log.Printf("WARM START - RequestID: %s", lc.RequestID)
    }
    
    // Business logic
    return Response{}, nil
}
```

### CloudWatch Insights Query

```sql
fields @timestamp, @message
| filter @message like /COLD START/
| stats count() as coldStarts by bin(5m)
```

---

## ✅ Best Practices

### 1. Minimize Package Size

```bash
# Strip debug info
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main

# Compress
zip function.zip main

# Result: 5 MB instead of 15 MB
```

### 2. Move Heavy Init to init()

```go
func init() {
    // ✅ Run once during cold start
    db = connectDB()
    cache = loadCache()
}
```

### 3. Reuse Connections

```go
var httpClient = &http.Client{
    Timeout: 10 * time.Second,
}

func handler(ctx context.Context) error {
    // Reuse client (connection pooling)
    resp, err := httpClient.Get(url)
    // ...
}
```

### 4. Lazy Loading (if appropriate)

```go
var cache *Cache

func getCache() *Cache {
    if cache == nil {
        cache = loadCache()  // Load on first use
    }
    return cache
}
```

### 5. Use Appropriate Memory

```yaml
# Balance: cost vs performance
Memory: 1024 MB  # Good for most Go functions
```

---

## 🎓 Висновок

### Lambda Lifecycle:

✅ **INIT** → Cold start (100ms - 5s)  
✅ **INVOKE** → Handler execution  
✅ **FREEZE** → Container kept warm (5-15 min)  
✅ **SHUTDOWN** → Destroy after idle  

### Key Points:

1. **Cold start** happens on first invocation or after idle
2. **Warm start** reuses existing container (much faster)
3. **init()** runs only during cold start
4. **Global variables** persist across warm invocations
5. **Optimize** package size and init code
6. **Provisioned Concurrency** eliminates cold starts (costs $)

### Golden Rule:

**"Keep init() heavy, handler() light!"**

---

## 📖 Далі

- `02_cold_start_optimization.md` - Reducing cold start time
- `practice/01_basic_lambda/` - Basic Lambda examples
- `practice/02_api_gateway/` - Lambda + API Gateway

**"Serverless = No servers to manage, but still code to optimize!" ☁️**
