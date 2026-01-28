# Memory Allocations в Go

## 🎯 Що таке Allocation?

**Allocation** - це виділення пам'яті для змінних, структур, слайсів, maps тощо.

```
Stack (швидко) ← Escape Analysis → Heap (повільно)
```

---

## 📊 Stack vs Heap

### Stack Allocation (швидко)

```go
func example() {
    x := 42  // Stack allocation (швидко)
    // x автоматично "очищується" при виході з функції
}
```

**Характеристики:**
- ✅ Дуже швидко (just pointer increment)
- ✅ Автоматичне очищення
- ✅ Predictable
- ❌ Обмежений розмір (~few MB)

### Heap Allocation (повільно)

```go
func example() *int {
    x := 42
    return &x  // ❌ Escapes to heap!
}
```

**Характеристики:**
- ❌ Повільніше (malloc)
- ❌ Потребує GC для очищення
- ❌ Фрагментація пам'яті
- ✅ Необмежений розмір

---

## 🔍 Escape Analysis

Go compiler автоматично вирішує: stack чи heap?

### Що Escapes to Heap?

```go
// 1. Return pointer
func escape1() *int {
    x := 42
    return &x  // ❌ Escapes
}

// 2. Store pointer in global
var global *int

func escape2() {
    x := 42
    global = &x  // ❌ Escapes
}

// 3. Send pointer to channel
func escape3(ch chan *int) {
    x := 42
    ch <- &x  // ❌ Escapes
}

// 4. Large struct
func escape4() {
    var big [10000]int  // ❌ Too big for stack
    _ = big
}

// 5. Interface
func escape5() {
    x := 42
    var i interface{} = x  // ❌ Escapes (boxing)
    _ = i
}

// 6. Closure capturing pointer
func escape6() func() int {
    x := 42
    return func() int {
        return x  // ❌ x escapes
    }
}
```

### Що НЕ Escapes?

```go
// 1. Local variables
func noEscape1() {
    x := 42  // ✅ Stack
    y := x * 2
    _ = y
}

// 2. Return value (not pointer)
func noEscape2() int {
    x := 42
    return x  // ✅ Stack
}

// 3. Small arrays
func noEscape3() {
    arr := [10]int{}  // ✅ Stack
    _ = arr
}
```

---

## 🔍 Як перевірити Escape Analysis?

### Build з `-gcflags`

```bash
go build -gcflags="-m -l" main.go
```

**Output:**
```
./main.go:5:6: can inline noEscape
./main.go:10:9: &x escapes to heap
./main.go:9:2: moved to heap: x
```

**Flags:**
- `-m`: Show optimization decisions
- `-m -m`: More verbose
- `-l`: Disable inlining (для чистого аналізу)

### Example

```go
package main

type Data struct {
    value int
}

func stackAlloc() Data {
    return Data{value: 42}  // ✅ Stack
}

func heapAlloc() *Data {
    return &Data{value: 42}  // ❌ Heap
}

func main() {
    _ = stackAlloc()
    _ = heapAlloc()
}
```

```bash
$ go build -gcflags="-m" main.go
./main.go:8:6: can inline stackAlloc
./main.go:12:9: &Data{...} escapes to heap
```

---

## 📊 Benchmarking Allocations

```go
func BenchmarkStackAlloc(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = stackAlloc()
    }
}

func BenchmarkHeapAlloc(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = heapAlloc()
    }
}
```

```bash
$ go test -bench=. -benchmem
BenchmarkStackAlloc-8    1000000000    0.25 ns/op    0 B/op    0 allocs/op
BenchmarkHeapAlloc-8     50000000      30 ns/op      8 B/op    1 allocs/op
```

**Аналіз:**
- Stack: 0.25 ns/op, 0 allocs ✅
- Heap: 30 ns/op, 1 alloc ❌ (120x повільніше!)

---

## 🎯 Common Allocation Sources

### 1. Slices

```go
// ❌ BAD: Multiple allocations
func bad() []int {
    var result []int
    for i := 0; i < 1000; i++ {
        result = append(result, i)  // Grows & reallocates
    }
    return result
}

// ✅ GOOD: Single allocation
func good() []int {
    result := make([]int, 0, 1000)  // Pre-allocate
    for i := 0; i < 1000; i++ {
        result = append(result, i)
    }
    return result
}
```

### 2. Maps

```go
// ❌ BAD: Default capacity
m := make(map[string]int)

// ✅ GOOD: Pre-allocate
m := make(map[string]int, 1000)
```

### 3. String Concatenation

```go
// ❌ BAD: Multiple allocations
func bad(items []string) string {
    result := ""
    for _, item := range items {
        result += item  // New string each iteration!
    }
    return result
}

// ✅ GOOD: Single allocation
func good(items []string) string {
    var b strings.Builder
    b.Grow(len(items) * 10)  // Estimate size
    for _, item := range items {
        b.WriteString(item)
    }
    return b.String()
}
```

### 4. []byte ↔ string Conversions

```go
// ❌ BAD: Allocation
s := string(bytes)  // Copies!

// ✅ GOOD: Zero-copy (unsafe)
import "unsafe"

func bytesToString(b []byte) string {
    return *(*string)(unsafe.Pointer(&b))
}
```

**⚠️ Warning:** Unsafe! Only if you know bytes won't change.

### 5. Interface Boxing

```go
// ❌ BAD: Boxing allocation
func bad(x int) {
    var i interface{} = x  // Allocation!
    process(i)
}

// ✅ GOOD: Avoid interface
func good(x int) {
    process(x)
}
```

---

## 🔧 Reducing Allocations

### Technique 1: Pre-allocate Slices

```go
// ❌ Before
var results []Result
for _, item := range items {
    results = append(results, process(item))
}

// ✅ After
results := make([]Result, 0, len(items))
for _, item := range items {
    results = append(results, process(item))
}
```

### Technique 2: Reuse Buffers

```go
// ❌ Before
func process(data []byte) {
    buf := make([]byte, 1024)  // Allocation!
    // Use buf
}

// ✅ After
var bufPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 1024)
    },
}

func process(data []byte) {
    buf := bufPool.Get().([]byte)
    defer bufPool.Put(buf)
    // Use buf
}
```

### Technique 3: Avoid String Concatenation

```go
// ❌ Before
func format(name, age string) string {
    return "Name: " + name + ", Age: " + age
}

// ✅ After (fmt - but still allocates)
func format(name, age string) string {
    return fmt.Sprintf("Name: %s, Age: %s", name, age)
}

// ✅ Better: strings.Builder
func format(name, age string) string {
    var b strings.Builder
    b.Grow(50)
    b.WriteString("Name: ")
    b.WriteString(name)
    b.WriteString(", Age: ")
    b.WriteString(age)
    return b.String()
}
```

### Technique 4: Return Values, Not Pointers

```go
// ❌ Before
func NewData() *Data {
    return &Data{value: 42}  // Heap allocation
}

// ✅ After
func NewData() Data {
    return Data{value: 42}  // Stack allocation
}
```

**Note:** Only if struct is small (<= few hundred bytes)

---

## 📊 Measuring Allocations

### In Tests

```go
func TestAllocations(t *testing.T) {
    var m1, m2 runtime.MemStats
    
    runtime.GC()
    runtime.ReadMemStats(&m1)
    
    // Code to measure
    result := expensive()
    
    runtime.ReadMemStats(&m2)
    
    allocs := m2.TotalAlloc - m1.TotalAlloc
    t.Logf("Allocations: %d bytes", allocs)
    
    if allocs > 1000 {
        t.Errorf("Too many allocations: %d", allocs)
    }
}
```

### With Benchmarks

```bash
go test -bench=. -benchmem -memprofile=mem.out

# Analyze
go tool pprof mem.out
> top
> list functionName
```

---

## 🎯 Allocation Hotspots

### 1. JSON Encoding/Decoding

```go
// ❌ Allocation-heavy
json.Marshal(data)
json.Unmarshal(bytes, &v)

// ✅ Better: Reuse encoder/decoder
encoder := json.NewEncoder(writer)
encoder.Encode(data)
```

### 2. Regex

```go
// ❌ Compile each time
re := regexp.MustCompile(`\d+`)

// ✅ Compile once
var digitRe = regexp.MustCompile(`\d+`)
```

### 3. Time Formatting

```go
// ❌ Allocates
s := time.Now().Format(time.RFC3339)

// ✅ Better: Append to buffer
var b []byte
b = time.Now().AppendFormat(b, time.RFC3339)
```

---

## 🔍 Tools

### 1. Escape Analysis

```bash
go build -gcflags="-m -m" main.go
```

### 2. Benchmarks

```bash
go test -bench=. -benchmem
```

### 3. Memory Profiler

```bash
go test -memprofile=mem.out
go tool pprof mem.out
```

### 4. Allocation Tracer

```bash
GODEBUG=allocfreetrace=1 ./myapp 2>&1 | grep myFunction
```

---

## ✅ Best Practices

1. **Pre-allocate slices** з відомою capacity
2. **Reuse buffers** з sync.Pool
3. **Return values** замість pointers (для малих structs)
4. **Avoid string concatenation** (use strings.Builder)
5. **Minimize interface{}** usage
6. **Profile before optimizing** (measure!)

---

## 🎓 Висновок

### Stack vs Heap:

✅ **Stack:** Fast, automatic cleanup  
❌ **Heap:** Slow, needs GC  

### Key Points:

1. Escape Analysis вирішує stack vs heap
2. Pointer return → heap allocation
3. Pre-allocate slices/maps
4. Reuse buffers (sync.Pool)
5. Benchmark to measure (`-benchmem`)

### Golden Rule:

**"Don't allocate if you can avoid it!"**

---

## 📖 Далі

- `02_gc_basics.md` - Garbage Collector
- `03_sync_pool.md` - Object Pooling
- `practice/02_allocations/` - Optimization examples

**"Every allocation avoided is 100x faster!" 🚀**
