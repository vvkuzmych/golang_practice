# sync.Pool - Object Pooling

## 🎯 Що таке sync.Pool?

**sync.Pool** - це cache для reusable objects, що зменшує кількість allocations.

```
Get() → Use object → Put() → Reuse
```

**Мета:** Reduce GC pressure шляхом reusing objects.

---

## 📊 Проблема без Pool

```go
func process() {
    buf := make([]byte, 1024)  // ❌ Allocation!
    // Use buf
    // buf freed by GC
}

// Call 1000 times = 1000 allocations!
```

**Issues:**
- 1000 allocations
- 1000 GC cycles
- High memory pressure

---

## ✅ Рішення: sync.Pool

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 1024)
    },
}

func process() {
    buf := bufferPool.Get().([]byte)  // ✅ Reuse!
    defer bufferPool.Put(buf)
    
    // Use buf
    // buf returned to pool
}

// Call 1000 times ≈ few allocations!
```

**Benefits:**
- Reuse objects
- Fewer allocations
- Less GC pressure

---

## 🏗️ Базове використання

### Create Pool

```go
var pool = sync.Pool{
    New: func() interface{} {
        // Create new object
        return &MyObject{}
    },
}
```

### Get & Put

```go
// Get from pool
obj := pool.Get().(*MyObject)

// Use object
obj.DoWork()

// Return to pool
pool.Put(obj)
```

### With defer

```go
func process() {
    obj := pool.Get().(*MyObject)
    defer pool.Put(obj)  // ✅ Always return
    
    obj.DoWork()
}
```

---

## 🎯 Real-World Examples

### Example 1: Buffer Pool

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func process(data []byte) []byte {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()  // ⚠️ Important!
        bufferPool.Put(buf)
    }()
    
    buf.Write(data)
    buf.WriteString(" processed")
    
    return buf.Bytes()
}
```

**⚠️ Important:** Always `Reset()` before `Put()`!

### Example 2: JSON Encoder Pool

```go
var encoderPool = sync.Pool{
    New: func() interface{} {
        return json.NewEncoder(nil)
    },
}

func encodeJSON(w io.Writer, v interface{}) error {
    enc := encoderPool.Get().(*json.Encoder)
    defer encoderPool.Put(enc)
    
    enc.SetEscapeHTML(false)
    return enc.Encode(v)
}
```

### Example 3: Slice Pool

```go
var slicePool = sync.Pool{
    New: func() interface{} {
        return make([]int, 0, 100)
    },
}

func process(data []int) []int {
    slice := slicePool.Get().([]int)
    defer func() {
        slice = slice[:0]  // ⚠️ Reset length!
        slicePool.Put(slice)
    }()
    
    for _, v := range data {
        slice = append(slice, v*2)
    }
    
    return slice
}
```

---

## ⚠️ Important Rules

### Rule 1: Always Reset Objects

```go
// ❌ BAD: Don't reset
defer pool.Put(buf)

// ✅ GOOD: Reset before Put
defer func() {
    buf.Reset()
    pool.Put(buf)
}()
```

**Why?** Prevent data leakage між uses.

### Rule 2: Don't Hold References

```go
// ❌ BAD: Держить reference
func bad() *bytes.Buffer {
    buf := pool.Get().(*bytes.Buffer)
    defer pool.Put(buf)
    return buf  // ❌ Reference after Put!
}

// ✅ GOOD: Copy data
func good() []byte {
    buf := pool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        pool.Put(buf)
    }()
    
    result := make([]byte, buf.Len())
    copy(result, buf.Bytes())
    return result  // ✅ Safe copy
}
```

### Rule 3: Objects Can Be Dropped

```go
// Pool is NOT a cache!
// Objects can be dropped by GC at any time
```

**⚠️ Warning:** Pool може бути очищений GC! Use for temporary objects only.

---

## 📊 Benchmarking Pool

### Without Pool

```go
func BenchmarkWithoutPool(b *testing.B) {
    for i := 0; i < b.N; i++ {
        buf := make([]byte, 1024)
        _ = buf
    }
}
```

### With Pool

```go
var pool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 1024)
    },
}

func BenchmarkWithPool(b *testing.B) {
    for i := 0; i < b.N; i++ {
        buf := pool.Get().([]byte)
        pool.Put(buf)
    }
}
```

### Results

```bash
$ go test -bench=. -benchmem
BenchmarkWithoutPool-8    50000000    30 ns/op    1024 B/op    1 allocs/op
BenchmarkWithPool-8      200000000     6 ns/op       0 B/op    0 allocs/op
```

**Improvement:**
- 5x faster
- 0 allocations!

---

## 🎯 Advanced Patterns

### Pattern 1: Typed Pool

```go
type BufferPool struct {
    pool sync.Pool
}

func NewBufferPool() *BufferPool {
    return &BufferPool{
        pool: sync.Pool{
            New: func() interface{} {
                return new(bytes.Buffer)
            },
        },
    }
}

func (p *BufferPool) Get() *bytes.Buffer {
    return p.pool.Get().(*bytes.Buffer)
}

func (p *BufferPool) Put(buf *bytes.Buffer) {
    buf.Reset()
    p.pool.Put(buf)
}

// Usage
var bufPool = NewBufferPool()

func process() {
    buf := bufPool.Get()
    defer bufPool.Put(buf)
    // Use buf
}
```

### Pattern 2: Size-Based Pools

```go
var (
    smallPool = sync.Pool{
        New: func() interface{} {
            return make([]byte, 1024)  // 1KB
        },
    }
    
    largePool = sync.Pool{
        New: func() interface{} {
            return make([]byte, 64*1024)  // 64KB
        },
    }
)

func getBuffer(size int) []byte {
    if size <= 1024 {
        return smallPool.Get().([]byte)
    }
    return largePool.Get().([]byte)
}

func putBuffer(buf []byte) {
    if cap(buf) <= 1024 {
        smallPool.Put(buf[:1024])
    } else {
        largePool.Put(buf[:64*1024])
    }
}
```

### Pattern 3: Pool with Cleanup

```go
type ResourcePool struct {
    pool    sync.Pool
    cleanup func(interface{})
}

func NewResourcePool(
    newFunc func() interface{},
    cleanupFunc func(interface{}),
) *ResourcePool {
    return &ResourcePool{
        pool: sync.Pool{New: newFunc},
        cleanup: cleanupFunc,
    }
}

func (p *ResourcePool) Get() interface{} {
    return p.pool.Get()
}

func (p *ResourcePool) Put(obj interface{}) {
    if p.cleanup != nil {
        p.cleanup(obj)
    }
    p.pool.Put(obj)
}

// Usage
var connPool = NewResourcePool(
    func() interface{} {
        return &Connection{/* ... */}
    },
    func(obj interface{}) {
        conn := obj.(*Connection)
        conn.Reset()  // Cleanup
    },
)
```

---

## 🔍 When to Use Pool

### ✅ Good Use Cases

1. **Temporary buffers** (bytes.Buffer, []byte)
2. **JSON encoders/decoders**
3. **HTTP request/response objects**
4. **Database connection wrappers**
5. **Parser states**
6. **Worker goroutine contexts**

### ❌ Bad Use Cases

1. **Long-lived objects** (use explicit cache)
2. **Objects with complex lifecycle**
3. **Objects requiring cleanup** (use defer/finalizers)
4. **Very small objects** (< 100 bytes - not worth it)

---

## 📈 Real-World Examples

### HTTP Handler

```go
var responsePool = sync.Pool{
    New: func() interface{} {
        return &Response{
            Headers: make(map[string]string),
            Body:    new(bytes.Buffer),
        }
    },
}

func handler(w http.ResponseWriter, r *http.Request) {
    resp := responsePool.Get().(*Response)
    defer func() {
        resp.Reset()
        responsePool.Put(resp)
    }()
    
    // Build response
    resp.Status = 200
    resp.Body.WriteString("Hello")
    
    // Write response
    w.WriteHeader(resp.Status)
    w.Write(resp.Body.Bytes())
}
```

### JSON Processing

```go
var (
    encoderPool = sync.Pool{
        New: func() interface{} {
            return json.NewEncoder(nil)
        },
    }
    
    bufferPool = sync.Pool{
        New: func() interface{} {
            return new(bytes.Buffer)
        },
    }
)

func marshalJSON(v interface{}) ([]byte, error) {
    buf := bufferPool.Get().(*bytes.Buffer)
    enc := encoderPool.Get().(*json.Encoder)
    
    defer func() {
        buf.Reset()
        bufferPool.Put(buf)
        encoderPool.Put(enc)
    }()
    
    enc.Reset(buf)
    err := enc.Encode(v)
    
    result := make([]byte, buf.Len())
    copy(result, buf.Bytes())
    
    return result, err
}
```

---

## ⚠️ Common Mistakes

### Mistake 1: Not Resetting

```go
// ❌ BAD
defer pool.Put(buf)

// ✅ GOOD
defer func() {
    buf.Reset()
    pool.Put(buf)
}()
```

### Mistake 2: Holding Reference

```go
// ❌ BAD
func bad() []byte {
    buf := pool.Get().([]byte)
    defer pool.Put(buf)
    return buf  // ❌ Dangerous!
}

// ✅ GOOD
func good() []byte {
    buf := pool.Get().([]byte)
    defer pool.Put(buf)
    result := make([]byte, len(buf))
    copy(result, buf)
    return result
}
```

### Mistake 3: Pool as Cache

```go
// ❌ BAD: Expect objects to persist
obj := pool.Get()
// ... later (may be GC'd!)
pool.Put(obj)

// ✅ GOOD: Use for short-lived operations only
obj := pool.Get()
defer pool.Put(obj)
// Use immediately
```

---

## 📊 Measuring Impact

### Before

```go
func process(data []byte) []byte {
    buf := new(bytes.Buffer)  // Allocation
    buf.Write(data)
    return buf.Bytes()
}
```

```bash
BenchmarkBefore-8    1000000    1200 ns/op    2048 B/op    2 allocs/op
```

### After

```go
var pool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func process(data []byte) []byte {
    buf := pool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        pool.Put(buf)
    }()
    
    buf.Write(data)
    result := make([]byte, buf.Len())
    copy(result, buf.Bytes())
    return result
}
```

```bash
BenchmarkAfter-8    5000000    250 ns/op    1024 B/op    1 allocs/op
```

**Improvement:**
- 4.8x faster
- 50% fewer allocations

---

## 🎓 Висновок

### sync.Pool:

✅ **Reuse objects** (reduce allocations)  
✅ **Automatic** (GC-aware)  
✅ **Thread-safe** (built-in locking)  
✅ **Zero configuration** (just New func)  

### Key Points:

1. Use for **temporary objects**
2. Always **Reset()** before **Put()**
3. Don't **hold references** after Put
4. Pool can be **cleared by GC**
5. Benchmark to **measure impact**

### Golden Rules:

1. **Reset objects** before returning
2. **Copy data** if needed after Put
3. **Use for temp objects** (not long-lived)
4. **Benchmark** to verify improvement
5. **Don't rely** on objects persisting

---

## 📖 Далі

- `practice/03_sync_pool/` - Pool examples
- `practice/02_allocations/` - Optimization benchmarks
- Real-world pool patterns

**"Reuse, don't reallocate!" 🔄**
