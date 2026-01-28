# Week 10 - Quick Start ⚡

## 🎯 Мета

Писати **швидкий Go код** через розуміння allocations, GC, та sync.Pool.

---

## ⚡ 3-хвилинний старт

### 1. Check Escape Analysis

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_10

# Create test file
cat > test.go << 'GOEOF'
package main

type Data struct {
    value int
}

func stack() Data {
    return Data{value: 42}
}

func heap() *Data {
    return &Data{value: 42}
}

func main() {}
GOEOF

# Check what escapes
go build -gcflags="-m" test.go
```

**Output:**
```
./test.go:12:9: &Data{...} escapes to heap
```

### 2. Monitor GC

```bash
# Run with GC trace
GODEBUG=gctrace=1 go run test.go
```

### 3. Benchmark Example

```bash
# Create benchmark
cat > bench_test.go << 'GOEOF'
package main

import "testing"

func BenchmarkStack(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = stack()
    }
}

func BenchmarkHeap(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = heap()
    }
}
GOEOF

# Run benchmark
go test -bench=. -benchmem
```

---

## 📚 Читати теорію

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_10

# Allocations
cat theory/01_allocations.md

# GC
cat theory/02_gc_basics.md

# sync.Pool
cat theory/03_sync_pool.md
```

---

## 🎯 3 Головні Концепції

### 1. Stack vs Heap

```go
// Stack (fast ✅)
func stack() Data {
    x := Data{value: 42}
    return x
}

// Heap (slow ❌)
func heap() *Data {
    x := Data{value: 42}
    return &x  // Escapes to heap!
}
```

**Benchmark:**
```
BenchmarkStack-8    1000000000    0.25 ns/op    0 B/op    0 allocs/op
BenchmarkHeap-8      50000000      30 ns/op     8 B/op    1 allocs/op
```

### 2. GC Cycle

```
1. Mark Setup (STW ~50μs)
   ↓
2. Marking (Concurrent)
   ↓
3. Mark Termination (STW ~50μs)
   ↓
4. Sweep (Concurrent)
```

**Monitor:**
```bash
GODEBUG=gctrace=1 ./myapp
```

### 3. sync.Pool

```go
var pool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 1024)
    },
}

func process() {
    buf := pool.Get().([]byte)
    defer pool.Put(buf)
    // Use buf
}
```

**Improvement:** 5x faster, 0 allocations!

---

## 🔧 Essential Commands

### Escape Analysis

```bash
# What escapes to heap?
go build -gcflags="-m" main.go

# More verbose
go build -gcflags="-m -m" main.go
```

### Benchmarking

```bash
# Basic
go test -bench=.

# With memory stats
go test -bench=. -benchmem

# Memory profile
go test -bench=. -memprofile=mem.out
go tool pprof mem.out
```

### GC Monitoring

```bash
# GC trace
GODEBUG=gctrace=1 ./myapp

# Configure GC
GOGC=50 ./myapp   # More frequent
GOGC=200 ./myapp  # Less frequent
```

---

## 📊 Quick Wins

### Win 1: Pre-allocate Slices

```go
// ❌ SLOW
var results []int
for i := 0; i < 1000; i++ {
    results = append(results, i)
}

// ✅ FAST (10x faster!)
results := make([]int, 0, 1000)
for i := 0; i < 1000; i++ {
    results = append(results, i)
}
```

### Win 2: strings.Builder

```go
// ❌ SLOW
result := ""
for _, s := range items {
    result += s
}

// ✅ FAST (50x faster!)
var b strings.Builder
for _, s := range items {
    b.WriteString(s)
}
result := b.String()
```

### Win 3: sync.Pool

```go
// ❌ SLOW
func process() {
    buf := make([]byte, 1024)
    // use buf
}

// ✅ FAST (5x faster!)
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

---

## 🎯 Optimization Workflow

### Step 1: Measure

```bash
go test -bench=. -benchmem > before.txt
```

### Step 2: Optimize

- Pre-allocate slices
- Use sync.Pool
- Avoid interface{}
- Use strings.Builder

### Step 3: Verify

```bash
go test -bench=. -benchmem > after.txt
benchstat before.txt after.txt
```

---

## 🔍 Debug Commands

### Memory Profiling

```bash
go test -bench=. -memprofile=mem.out
go tool pprof mem.out
> top
> list functionName
```

### Allocation Trace

```bash
GODEBUG=allocfreetrace=1 ./myapp 2>&1 | grep myFunction
```

### Execution Trace

```bash
go test -bench=. -trace=trace.out
go tool trace trace.out
```

---

## ✅ Quick Checklist

### Before Coding

- [ ] Know if it's hot path
- [ ] Have benchmark
- [ ] Profiled if needed

### While Coding

- [ ] Pre-allocate known sizes
- [ ] Reuse buffers (sync.Pool)
- [ ] Return values, not pointers (small structs)
- [ ] Use strings.Builder

### After Coding

- [ ] Benchmark shows improvement
- [ ] Check escape analysis
- [ ] Monitor GC if relevant

---

## 📖 Структура

```
week_10/
├── README.md              # Повний опис
├── QUICK_START.md         # Цей файл
├── theory/
│   ├── 01_allocations.md  # Stack vs Heap
│   ├── 02_gc_basics.md    # Garbage Collector
│   └── 03_sync_pool.md    # Object Pooling
└── practice/
    ├── 01_benchmarks/     # Benchmark examples
    ├── 02_allocations/    # Optimization examples
    └── 03_sync_pool/      # Pool patterns
```

---

## 🎓 Learning Path

### День 1: Allocations
1. Читай `theory/01_allocations.md`
2. Розумій Stack vs Heap
3. Експериментуй з escape analysis
4. Write benchmarks

### День 2: GC
1. Читай `theory/02_gc_basics.md`
2. Моніторь GC з gctrace
3. Tune GOGC
4. Profile memory

### День 3: sync.Pool
1. Читай `theory/03_sync_pool.md`
2. Implement buffer pool
3. Benchmark improvement
4. Apply to real code

---

## 🎯 Golden Rules

### Rule 1: Measure First

```bash
# Always benchmark before optimizing
go test -bench=. -benchmem
```

### Rule 2: Reduce Allocations

```go
// Pre-allocate if size known
results := make([]T, 0, expectedSize)
```

### Rule 3: Reuse Objects

```go
// Use sync.Pool for temp objects
var pool = sync.Pool{New: ...}
obj := pool.Get()
defer pool.Put(obj)
```

### Rule 4: Check Escape

```bash
# See what allocates
go build -gcflags="-m" main.go
```

### Rule 5: Profile

```bash
# Find real bottlenecks
go test -bench=. -memprofile=mem.out
go tool pprof mem.out
```

---

## 📊 Expected Results

### Good Benchmark

```
BenchmarkOptimized-8    10000000    120 ns/op    64 B/op    1 allocs/op
```

### Bad Benchmark

```
BenchmarkSlow-8    100000    12000 ns/op    4096 B/op    50 allocs/op
```

**Target:** < 5 allocs/op для hot paths

---

## 🎉 Quick Test

```bash
# 1. Create simple benchmark
cat > quick_test.go << 'GOEOF'
package main

import (
    "strings"
    "testing"
)

func concat(items []string) string {
    result := ""
    for _, s := range items {
        result += s
    }
    return result
}

func builder(items []string) string {
    var b strings.Builder
    for _, s := range items {
        b.WriteString(s)
    }
    return b.String()
}

func BenchmarkConcat(b *testing.B) {
    items := []string{"a", "b", "c", "d", "e"}
    for i := 0; i < b.N; i++ {
        _ = concat(items)
    }
}

func BenchmarkBuilder(b *testing.B) {
    items := []string{"a", "b", "c", "d", "e"}
    for i := 0; i < b.N; i++ {
        _ = builder(items)
    }
}
GOEOF

# 2. Run
go test -bench=. -benchmem

# 3. See dramatic difference!
```

---

**"Fast code = Happy users!" ⚡**
