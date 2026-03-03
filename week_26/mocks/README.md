# Mocking and Concurrency Patterns

This directory demonstrates **dependency injection, mocking, and concurrency patterns** in Go.

---

## 📁 Directory Structure

```
week_26/mocks/
├── main.go                 - Async patterns example (channels, callbacks)
├── main_test.go            - Tests for async patterns
├── mock_url_sorter.go      - Mock for URLSorter interface
└── fanout/
    ├── main.go             - Fan-Out pattern (simplest example)
    └── main_test.go        - Tests for fan-out pattern
```

---

## 📚 Examples

### **1. Async Patterns** (`main.go`)

Demonstrates asynchronous processing patterns with dependency injection.

**Key Concepts:**
- ✅ Interface-based design (`URLSorter`)
- ✅ Channel-based async (`ProcessAsync`)
- ✅ Callback-based async (`ProcessWithCallback`)
- ✅ Context cancellation
- ✅ Mocking for tests

**Run:**
```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_26/mocks
go run main.go mock_url_sorter.go
```

**Test:**
```bash
go test -v
```

**What it shows:**
- Using real implementation vs mock
- Channel-based async communication
- Callback pattern
- Context timeout handling

---

### **2. Fan-Out Pattern** (`fanout/`)

Demonstrates **Fan-Out concurrency pattern** - the simplest possible example.

**Key Concepts:**
- ✅ Fan-Out: Multiple workers reading from single channel
- ✅ No structs, no interfaces - just functions
- ✅ `ProcessSequential()` vs `ProcessFanOut()` comparison
- ✅ Pure function-based approach

**Run:**
```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_26/mocks/fanout
go run main.go
```

**Test:**
```bash
go test -v
```

**What it shows:**
- **Sequential**: 0.8 seconds for 8 numbers
- **Fan-Out (4 workers)**: 0.2 seconds for 8 numbers
- **Speedup**: ~4x faster! 🚀

---

## 🎯 What is Fan-Out?

### **Pattern Explanation:**

```
        Input Channel
             │
      ┌──────┴──────┬──────┬──────┐
      ↓              ↓      ↓      ↓
   Worker 1      Worker 2  Worker 3  Worker 4
      │              │      │      │
   Process        Process  Process  Process
      │              │      │      │
      └──────┬───────┴──────┴──────┘
             ↓
       Results Channel
```

**How it works:**
1. **One input channel** feeds work to all workers
2. **Multiple workers** (goroutines) read from the SAME channel
3. **Go's channel guarantees** each item goes to only ONE worker
4. **Parallel processing** - workers operate concurrently
5. **Results collected** from all workers

---

## 🔍 Key Differences Between Examples

### **Async Patterns (`main.go`)**
- **Focus**: Asynchronous communication patterns
- **Pattern**: Single goroutine per operation
- **Use case**: Non-blocking operations
- **Speedup**: Doesn't block caller

### **Fan-Out (`fanout/main.go`)**
- **Focus**: Parallel processing
- **Pattern**: Multiple workers from one source
- **Use case**: Batch processing, parallel work
- **Speedup**: ~N× faster (N = number of workers)

---

## 🧪 Test Coverage

### **Async Patterns Tests:**
- ✅ Success case
- ✅ Error handling
- ✅ Context cancellation
- ✅ Callback pattern
- ✅ Nil function handling
- 📊 Benchmarks

### **Fan-Out Tests:**
- ✅ Basic functionality
- ✅ Performance improvement (vs sequential)
- ✅ Error handling
- ✅ Context cancellation
- ✅ Different worker counts (1, 2, 4, 8)
- ✅ Order independence
- 📊 Benchmarks (sequential vs 2/4 workers)

---

## 🚀 Running Everything

### **Run both examples:**
```bash
# Async patterns
go run main.go mock_url_sorter.go

# Fan-out pattern
go run fanout/main.go
```

### **Run all tests:**
```bash
# Test async patterns
go test -v

# Test fan-out
go test -v ./fanout/

# All tests with coverage
go test -cover ./...
```

### **Benchmarks:**
```bash
# Compare fan-out performance
cd fanout && go test -bench=BenchmarkFanOut -benchmem
```

---

## 💡 Key Takeaways

### **Why Separate Mock Files?**
- ✅ **Reusability** - mock can be used in multiple test files
- ✅ **Clarity** - clear separation of concerns
- ✅ **Organization** - easy to find mocks
- ✅ **Same package** - tests can still access it

### **Why Fan-Out is Powerful?**
- ⚡ **Parallel processing** - utilize multiple CPU cores
- 🚀 **Scalable** - add more workers for more speed
- 📦 **Simple pattern** - easy to implement and understand
- 🎯 **Real-world use** - batch jobs, API calls, file processing

### **When to Use Each:**

| Use Case | Pattern |
|----------|---------|
| Non-blocking single operation | Async (channels/callbacks) |
| Batch processing many items | Fan-Out |
| I/O-bound work (network, disk) | Fan-Out |
| Need to aggregate results | Fan-Out + Fan-In |

---

## 📖 Further Reading

- **Worker Pool**: Similar to fan-out but with fixed workers processing from a queue
- **Fan-In**: Opposite of fan-out - merge multiple channels into one
- **Pipeline**: Chain fan-out stages together

---

**All examples include mocks separated in their own files for clean testing!** ✨
