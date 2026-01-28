# Garbage Collector (GC) в Go

## 🎯 Що таке GC?

**Garbage Collector** - це автоматичне управління пам'яттю, що звільняє неіснуючі об'єкти з heap.

```
Allocate → Use → GC finds unused → Free memory
```

---

## 📊 Go GC: Concurrent Mark & Sweep

### Характеристики Go GC

- **Concurrent** (runs in parallel з application)
- **Non-generational** (no young/old separation)
- **Tri-color mark-and-sweep** algorithm
- **Low latency** (< 1ms pause for most apps)
- **Write barrier** (tracks pointer changes)

### GC Cycle

```
1. Mark Setup (STW)  ← ~10-50μs pause
2. Marking (Concurrent)
3. Mark Termination (STW)  ← ~10-50μs pause
4. Sweep (Concurrent)
```

**STW** = Stop The World (application pauses)

---

## 🔄 GC Phases

### Phase 1: Mark Setup (STW)

```
Application → PAUSE → Enable write barrier
                   → Scan stacks
                   → RESUME
```

**Duration:** 10-50 microseconds (дуже швидко)

### Phase 2: Concurrent Marking

```
Application runs ║ GC marks reachable objects
                 ║ (tri-color algorithm)
```

**No pause!** Application і GC працюють паралельно.

### Phase 3: Mark Termination (STW)

```
Application → PAUSE → Finalize marking
                   → Prepare sweep
                   → RESUME
```

**Duration:** 10-50 microseconds

### Phase 4: Concurrent Sweep

```
Application runs ║ GC frees unmarked objects
```

**No pause!** Пам'ять звільняється concurrently.

---

## 🎨 Tri-Color Algorithm

### Кольори об'єктів

```
White (не перевірено) → Gray (в черзі) → Black (перевірено)
```

**Algorithm:**

1. Start: All objects **White**
2. Roots (stacks, globals) → **Gray**
3. While gray exists:
   - Pick gray object
   - Mark it **Black**
   - Mark its children **Gray**
4. End: White objects = garbage

### Приклад

```
Step 1:
Root → [A]white → [B]white
              ↓
           [C]white

Step 2: Root gray
Root → [A]gray → [B]white
              ↓
           [C]white

Step 3: Process A (gray → black)
Root → [A]black → [B]gray
              ↓
           [C]gray

Step 4: Process B & C
Root → [A]black → [B]black
              ↓
           [C]black

All reachable = black
Unreachable = white (garbage)
```

---

## 🔧 GC Configuration

### GOGC Environment Variable

```bash
# Default: GOGC=100 (run GC when heap doubles)
GOGC=100 ./myapp

# More frequent GC (lower latency, higher CPU)
GOGC=50 ./myapp

# Less frequent GC (higher throughput, more memory)
GOGC=200 ./myapp

# Disable GC (testing only!)
GOGC=off ./myapp
```

**Formula:**
```
Target heap = Live heap * (1 + GOGC/100)

GOGC=100: Target = Live * 2    (double)
GOGC=50:  Target = Live * 1.5  (1.5x)
GOGC=200: Target = Live * 3    (triple)
```

### SetGCPercent Programmatically

```go
import "runtime/debug"

// Set GC target
old := debug.SetGCPercent(50)  // More frequent GC

// Disable GC temporarily
debug.SetGCPercent(-1)
// Critical section
debug.SetGCPercent(old)  // Re-enable
```

---

## 📊 Monitoring GC

### GC Stats

```go
import "runtime"

var m runtime.MemStats
runtime.ReadMemStats(&m)

fmt.Printf("Alloc: %v MB\n", m.Alloc/1024/1024)
fmt.Printf("TotalAlloc: %v MB\n", m.TotalAlloc/1024/1024)
fmt.Printf("Sys: %v MB\n", m.Sys/1024/1024)
fmt.Printf("NumGC: %v\n", m.NumGC)
fmt.Printf("PauseTotalNs: %v ms\n", m.PauseTotalNs/1000000)
fmt.Printf("LastGC: %v\n", time.Unix(0, int64(m.LastGC)))
```

### GC Trace

```bash
# Detailed GC trace
GODEBUG=gctrace=1 ./myapp

# Output:
gc 1 @0.001s 0%: 0.009+0.23+0.005 ms clock, 0.074+0.11/0.18/0.27+0.043 ms cpu
```

**Decode:**
```
gc 1          # GC cycle number
@0.001s       # Time since start
0%            # % CPU time in GC
0.009 ms      # STW sweep termination
0.23 ms       # Concurrent mark/scan
0.005 ms      # STW mark termination
```

---

## 🎯 Reducing GC Pressure

### Technique 1: Reduce Allocations

```go
// ❌ BAD: Many allocations
func bad() {
    for i := 0; i < 1000; i++ {
        data := make([]byte, 1024)  // 1000 allocations!
        process(data)
    }
}

// ✅ GOOD: Reuse buffer
func good() {
    data := make([]byte, 1024)  // 1 allocation
    for i := 0; i < 1000; i++ {
        process(data)
    }
}
```

### Technique 2: Pre-allocate

```go
// ❌ BAD: Grows dynamically
results := []Result{}
for _, item := range items {
    results = append(results, process(item))
}

// ✅ GOOD: Pre-allocate
results := make([]Result, 0, len(items))
for _, item := range items {
    results = append(results, process(item))
}
```

### Technique 3: Use sync.Pool

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 1024)
    },
}

func process() {
    buf := bufferPool.Get().([]byte)
    defer bufferPool.Put(buf)
    // Use buf
}
```

### Technique 4: Avoid Pointers (for small structs)

```go
// ❌ BAD: More GC work (pointers)
type Data struct {
    value *int
    name  *string
}

// ✅ GOOD: Less GC work (values)
type Data struct {
    value int
    name  string
}
```

**Why?** GC має перевіряти кожен pointer.

---

## 📈 GC Performance Tips

### Tip 1: Batch Operations

```go
// ❌ BAD: Many small allocations
for i := 0; i < 1000; i++ {
    save(Data{value: i})  // 1000 allocations
}

// ✅ GOOD: Batch
batch := make([]Data, 1000)
for i := 0; i < 1000; i++ {
    batch[i] = Data{value: i}
}
saveBatch(batch)  // 1 allocation
```

### Tip 2: Limit Heap Size

```go
// Set memory limit (Go 1.19+)
debug.SetMemoryLimit(2 * 1024 * 1024 * 1024)  // 2GB
```

### Tip 3: Manual GC (rarely needed)

```go
// Force GC (testing/profiling only!)
runtime.GC()

// After cleanup
items = nil
runtime.GC()  // Help GC reclaim memory
```

---

## 🔍 GC Debugging

### GODEBUG Options

```bash
# GC trace
GODEBUG=gctrace=1 ./myapp

# Allocations trace
GODEBUG=allocfreetrace=1 ./myapp

# Scheduler trace
GODEBUG=schedtrace=1000 ./myapp
```

### Memory Profiling

```bash
# Run with profiling
go test -memprofile=mem.out

# Analyze
go tool pprof mem.out
> top
> list main.expensive
> web  # Visualize
```

### GC Timeline Trace

```go
import (
    "os"
    "runtime/trace"
)

func main() {
    f, _ := os.Create("trace.out")
    defer f.Close()
    
    trace.Start(f)
    defer trace.Stop()
    
    // Your code
}
```

```bash
go tool trace trace.out
# Opens in browser with GC timeline!
```

---

## 🎯 When GC Runs

### Conditions

1. **Heap target reached** (GOGC threshold)
2. **2 minutes idle** (periodic check)
3. **Manual trigger** (`runtime.GC()`)

### Example

```
Live heap: 100MB
GOGC=100
Target: 100MB * 2 = 200MB

Allocate...
Heap: 150MB (no GC yet)
Heap: 200MB → GC triggers!

After GC:
Live heap: 120MB
New target: 120MB * 2 = 240MB
```

---

## 📊 GC Metrics

### Key Metrics

| Metric | Good | Bad | Fix |
|--------|------|-----|-----|
| **Pause time** | < 1ms | > 10ms | Reduce live heap |
| **GC frequency** | < 10/sec | > 100/sec | Reduce allocations |
| **Heap size** | Stable | Growing | Memory leak! |
| **CPU %** | < 5% | > 25% | Optimize allocations |

### Monitoring in Production

```go
import "expvar"

var (
    gcPauses = expvar.NewInt("gc_pauses_ms")
    gcRuns   = expvar.NewInt("gc_runs")
)

func monitorGC() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    gcPauses.Set(int64(m.PauseTotalNs / 1000000))
    gcRuns.Set(int64(m.NumGC))
}
```

**Expose:**
```bash
curl http://localhost:8080/debug/vars
```

---

## ⚠️ Common Issues

### Issue 1: Memory Leak

```go
// ❌ BAD: Global slice keeps growing
var cache []Data

func add(d Data) {
    cache = append(cache, d)  // Never cleaned!
}

// ✅ GOOD: Limit size
func add(d Data) {
    if len(cache) > 1000 {
        cache = cache[1:]  // Remove oldest
    }
    cache = append(cache, d)
}
```

### Issue 2: Too Many Pointers

```go
// ❌ BAD: Many pointers (more GC work)
type Node struct {
    value *int
    next  *Node
}

// ✅ BETTER: Values where possible
type Node struct {
    value int
    next  *Node
}
```

### Issue 3: Large Objects in Slices

```go
// ❌ BAD: Keeps entire slice alive
func first(items []LargeStruct) LargeStruct {
    return items[0]  // Keeps whole slice!
}

// ✅ GOOD: Copy value
func first(items []LargeStruct) LargeStruct {
    result := items[0]
    return result  // Slice can be freed
}
```

---

## 🎓 Висновок

### Go GC:

✅ **Concurrent** (low latency)  
✅ **Automatic** (no manual memory management)  
✅ **Efficient** (< 1ms pauses)  
✅ **Tunable** (GOGC, SetGCPercent)  

### Key Points:

1. GC runs when heap reaches target (GOGC)
2. Two short STW pauses (< 100μs each)
3. Most work concurrent with application
4. Reduce allocations → Reduce GC pressure
5. Monitor with gctrace & profiling

### Golden Rules:

1. **Reduce allocations** (fewer objects to collect)
2. **Pre-allocate** (avoid dynamic growth)
3. **Reuse buffers** (sync.Pool)
4. **Monitor GC** (gctrace, profiling)
5. **Tune GOGC** (if needed)

---

## 📖 Далі

- `03_sync_pool.md` - Object Pooling
- `practice/02_allocations/` - Optimization examples
- GC profiling & tuning

**"Happy GC = Happy Application!" 🗑️**
