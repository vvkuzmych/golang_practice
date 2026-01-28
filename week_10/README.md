# Week 10: Performance Optimization

## 🎯 Мета

Навчитися **писати швидкий Go код** через розуміння allocations, GC, та використання `sync.Pool`.

---

## 📚 Теорія

### 1. Memory Allocations
**Файл:** `theory/01_allocations.md`

- Stack vs Heap
- Escape Analysis
- Allocation hotspots
- Reducing allocations
- Benchmarking techniques

**Key concepts:**
```
Stack (fast) ← Escape Analysis → Heap (slow + GC)
```

### 2. Garbage Collector (GC)
**Файл:** `theory/02_gc_basics.md`

- Go GC algorithm (Concurrent Mark & Sweep)
- GC phases (Mark Setup, Marking, Mark Termination, Sweep)
- Tri-color marking
- GOGC configuration
- Monitoring GC

**GC Cycle:**
```
1. Mark Setup (STW ~50μs)
2. Marking (Concurrent)
3. Mark Termination (STW ~50μs)
4. Sweep (Concurrent)
```

### 3. sync.Pool
**Файл:** `theory/03_sync_pool.md`

- Object pooling pattern
- Reusing objects
- Reducing allocations
- When to use Pool
- Best practices

**Pattern:**
```
Get() → Use → Reset() → Put() → Reuse
```

---

## 💻 Практика

### 1. Benchmarks
**Директорія:** `practice/01_benchmarks/`

**Файли:**
- `basic_bench_test.go` - Базові benchmarks
- `alloc_bench_test.go` - Benchmarking allocations
- `README.md` - Інструкції

**Команди:**
```bash
go test -bench=.
go test -bench=. -benchmem
go test -bench=. -benchmem -memprofile=mem.out
```

### 2. Allocation Optimization
**Директорія:** `practice/02_allocations/`

**Файли:**
- `before_after_test.go` - Before/After приклади
- `escape_analysis.go` - Escape Analysis examples
- `README.md` - Пояснення

**Benchmarks:**
- String concatenation optimization
- Slice pre-allocation
- Buffer reuse

### 3. sync.Pool Examples
**Директорія:** `practice/03_sync_pool/`

**Файли:**
- `buffer_pool_test.go` - Buffer pooling
- `json_pool_test.go` - JSON encoder pooling
- `README.md` - Use cases

---

## 🔧 Основні команди

### Benchmarking

```bash
# Run benchmarks
go test -bench=.

# With memory stats
go test -bench=. -benchmem

# Specific benchmark
go test -bench=BenchmarkName

# CPU profiling
go test -bench=. -cpuprofile=cpu.out

# Memory profiling
go test -bench=. -memprofile=mem.out
```

### Escape Analysis

```bash
# Check what escapes to heap
go build -gcflags="-m" main.go

# More verbose
go build -gcflags="-m -m" main.go

# Disable inlining for clean analysis
go build -gcflags="-m -l" main.go
```

### GC Monitoring

```bash
# GC trace
GODEBUG=gctrace=1 ./myapp

# Allocation trace
GODEBUG=allocfreetrace=1 ./myapp

# Memory profiling
go tool pprof mem.out
```

---

## 📊 Performance Patterns

### Pattern 1: Pre-allocate Slices

```go
// ❌ SLOW: Multiple allocations
var results []Result
for _, item := range items {
    results = append(results, process(item))
}

// ✅ FAST: Single allocation
results := make([]Result, 0, len(items))
for _, item := range items {
    results = append(results, process(item))
}
```

### Pattern 2: Reuse Buffers

```go
// ❌ SLOW: New allocation each time
func process() {
    buf := make([]byte, 1024)
    // use buf
}

// ✅ FAST: Reuse from pool
var pool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 1024)
    },
}

func process() {
    buf := pool.Get().([]byte)
    defer pool.Put(buf)
    // use buf
}
```

### Pattern 3: Return Values, Not Pointers

```go
// ❌ SLOW: Heap allocation
func NewData() *Data {
    return &Data{value: 42}
}

// ✅ FAST: Stack allocation
func NewData() Data {
    return Data{value: 42}
}
```

### Pattern 4: strings.Builder

```go
// ❌ SLOW: Many allocations
result := ""
for _, s := range items {
    result += s  // New string each time!
}

// ✅ FAST: Single allocation
var b strings.Builder
b.Grow(len(items) * 10)
for _, s := range items {
    b.WriteString(s)
}
result := b.String()
```

---

## 🎯 Optimization Workflow

### Step 1: Benchmark

```bash
go test -bench=. -benchmem
```

**Example output:**
```
BenchmarkSlow-8    1000    1200 ns/op    2048 B/op    15 allocs/op
```

### Step 2: Profile

```bash
go test -bench=. -memprofile=mem.out
go tool pprof mem.out
> top
> list slowFunction
```

### Step 3: Optimize

- Reduce allocations
- Pre-allocate slices
- Use sync.Pool
- Avoid interface{}
- Use strings.Builder

### Step 4: Verify

```bash
go test -bench=. -benchmem
```

**Expected:**
```
BenchmarkFast-8    10000    120 ns/op    512 B/op    2 allocs/op
```

**Improvement:** 10x faster, 7x fewer allocs!

---

## 📈 Benchmarking Tips

### Tip 1: Use -benchmem

```bash
# Always include memory stats
go test -bench=. -benchmem
```

### Tip 2: Benchtime для точності

```bash
# Run longer for accurate results
go test -bench=. -benchtime=10s
```

### Tip 3: Compare з benchstat

```bash
# Run before
go test -bench=. -benchmem > old.txt

# Make changes

# Run after
go test -bench=. -benchmem > new.txt

# Compare
benchstat old.txt new.txt
```

### Tip 4: Prevent Compiler Optimizations

```go
var result int

func BenchmarkSomething(b *testing.B) {
    var r int
    for i := 0; i < b.N; i++ {
        r = compute()  // Don't let compiler optimize away
    }
    result = r  // Store to global
}
```

---

## 🔍 Profiling Tools

### CPU Profiling

```bash
go test -bench=. -cpuprofile=cpu.out
go tool pprof cpu.out
> top
> web
```

### Memory Profiling

```bash
go test -bench=. -memprofile=mem.out
go tool pprof mem.out
> top
> list functionName
```

### Execution Trace

```bash
go test -bench=. -trace=trace.out
go tool trace trace.out
```

### Escape Analysis

```bash
go build -gcflags="-m" main.go
```

---

## ✅ Best Practices

### Allocations

1. ✅ **Pre-allocate** slices with known capacity
2. ✅ **Reuse buffers** з sync.Pool
3. ✅ **Return values** замість pointers (small structs)
4. ✅ **Use strings.Builder** for string concatenation
5. ✅ **Benchmark** before і after optimization

### GC

1. ✅ **Reduce allocations** (fewer objects to collect)
2. ✅ **Avoid pointers** where possible (less GC work)
3. ✅ **Monitor GC** з gctrace
4. ✅ **Tune GOGC** if needed (default 100 is usually good)
5. ✅ **Profile** to find hotspots

### sync.Pool

1. ✅ **Reset objects** before Put()
2. ✅ **Don't hold references** after Put()
3. ✅ **Use for temp objects** (not long-lived)
4. ✅ **Benchmark** to verify improvement
5. ✅ **Copy data** if needed after Put()

---

## 📊 Performance Checklist

### Before Optimizing

- [ ] Have benchmarks
- [ ] Know current performance
- [ ] Identified bottleneck
- [ ] Profiled the code

### Optimization

- [ ] Pre-allocate slices/maps
- [ ] Use sync.Pool for reusable objects
- [ ] Minimize interface{} usage
- [ ] Avoid unnecessary conversions
- [ ] Use strings.Builder

### After Optimizing

- [ ] Benchmarks show improvement
- [ ] No new bugs introduced
- [ ] Code still readable
- [ ] Documented trade-offs

---

## 🎯 Common Hotspots

### 1. JSON Encoding/Decoding

```go
// ❌ SLOW
json.Marshal(data)

// ✅ FASTER: Reuse encoder
var encoderPool = sync.Pool{
    New: func() interface{} {
        return json.NewEncoder(nil)
    },
}
```

### 2. String Operations

```go
// ❌ SLOW
result := ""
for _, s := range items {
    result += s
}

// ✅ FAST
var b strings.Builder
for _, s := range items {
    b.WriteString(s)
}
```

### 3. Slice Growth

```go
// ❌ SLOW
var results []Result
for range items {
    results = append(results, ...)
}

// ✅ FAST
results := make([]Result, 0, len(items))
```

---

## 🚀 Quick Start

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_10

# Read theory
cat theory/01_allocations.md
cat theory/02_gc_basics.md
cat theory/03_sync_pool.md

# Run benchmarks (when created)
cd practice/01_benchmarks
go test -bench=. -benchmem

cd ../02_allocations
go test -bench=. -benchmem

cd ../03_sync_pool
go test -bench=. -benchmem
```

---

## 🎓 Learning Path

### День 1: Allocations

1. Читай `theory/01_allocations.md`
2. Розумій Stack vs Heap
3. Вивчи Escape Analysis
4. Запускай escape analysis на своєму коді

### День 2: GC

1. Читай `theory/02_gc_basics.md`
2. Розумій GC phases
3. Експериментуй з GOGC
4. Моніторь GC з gctrace

### День 3: sync.Pool

1. Читай `theory/03_sync_pool.md`
2. Зрозумій pooling pattern
3. Implement власний pool
4. Benchmark improvement

### День 4: Practice

1. Write benchmarks
2. Profile your code
3. Optimize allocations
4. Use sync.Pool
5. Measure improvement

---

## 📖 Ресурси

### Documentation

- [Go Performance Workshop](https://dave.cheney.net/high-performance-go-workshop/gopherchina-2019.html)
- [Profiling Go Programs](https://go.dev/blog/pprof)
- Week 8: Race detector & debugging

### Books

- "Efficient Go" by Bartłomiej Płotka
- "High Performance Go Workshop" by Dave Cheney

---

**"Fast code is good code, but correct code is better!" ⚡**

**Status:** Week 10 Materials Complete ✅  
**Created:** 2026-01-28
