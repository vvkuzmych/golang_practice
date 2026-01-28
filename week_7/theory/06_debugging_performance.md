# Debugging & Performance Optimization

## CPU Profiling

```bash
# Run with CPU profiling
go test -cpuprofile=cpu.prof -bench=.

# Analyze
go tool pprof cpu.prof
# Commands: top, list, web

# HTTP profiling
import _ "net/http/pprof"

go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()

# Access: http://localhost:6060/debug/pprof/
```

## Memory Profiling

```bash
# Memory profile
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof

# Heap profile
curl http://localhost:6060/debug/pprof/heap > heap.prof
go tool pprof heap.prof
```

## Race Detector

```bash
go test -race ./...
go run -race main.go
```

## Benchmarking

```go
func BenchmarkMyFunction(b *testing.B) {
    for i := 0; i < b.N; i++ {
        MyFunction()
    }
}

// Run: go test -bench=. -benchmem
```

## OpenTelemetry Tracing

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
)

func handleRequest(ctx context.Context) {
    ctx, span := otel.Tracer("myapp").Start(ctx, "handleRequest")
    defer span.End()
    
    // Your code
    result, err := doWork(ctx)
    if err != nil {
        span.RecordError(err)
    }
}
```

## Performance Tips

✅ Pre-allocate slices: `make([]int, 0, expectedSize)`
✅ Use `strings.Builder` for concatenation
✅ Reuse buffers with `sync.Pool`
✅ Avoid unnecessary allocations
✅ Use pointer receivers for large structs
✅ Profile before optimizing
