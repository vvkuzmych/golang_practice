# Week 15 - Recap & Advanced Topics

## 🎯 Мета

Повторення ключових тем + додаткові просунуті topics (error handling, mock generation, indexes).

---

## 📚 Topics Covered

### 1. Maps ✅ → [Week 1](../week_1/theory/01_types.md#8-мапи-map)

**Що вже покрито:**
- Створення maps (`make`, literal)
- Операції (set, get, delete)
- Перевірка існування (`value, ok`)
- Zero value (`nil map`)

**Key types:**
```go
// ✅ Valid key types (comparable):
int, string, bool, pointer, struct (with comparable fields), array

// ❌ Invalid key types (not comparable):
slice, map, function
```

**Runtime errors:**
```go
var m map[string]int  // nil map
m["key"] = 1          // PANIC: assignment to entry in nil map

// ✅ Fix:
m = make(map[string]int)
m["key"] = 1
```

➡️ **Детальніше:** [Week 1 - Types & Maps](../week_1/theory/01_types.md)

---

### 2. Runes ✅ → [Week 3](../week_3_rune_bytes/)

**Що вже покрито:**
- Що таке `rune` (`int32` для Unicode code point)
- `string` vs `[]rune` vs `[]byte`
- UTF-8 encoding
- `len(string)` (bytes) vs `len([]rune)` (characters)

**Приклад:**
```go
s := "Привіт"
fmt.Println(len(s))         // 12 (bytes)
fmt.Println(len([]rune(s))) // 6 (characters)
```

➡️ **Детальніше:** [Week 3 - Runes & Bytes](../week_3_rune_bytes/)

---

### 3. Select Statement ✅ → [Week 6](../week_6/theory/07_goroutines_concurrency.md#5-select)

**Що вже покрито:**
- Multiple channel operations
- Timeout pattern
- Non-blocking operations
- `default` case

**Приклад:**
```go
select {
case msg := <-ch1:
    fmt.Println("Ch1:", msg)
case msg := <-ch2:
    fmt.Println("Ch2:", msg)
case <-time.After(1 * time.Second):
    fmt.Println("Timeout")
default:
    fmt.Println("Non-blocking")
}
```

➡️ **Детальніше:** [Week 6 - Goroutines & Select](../week_6/theory/07_goroutines_concurrency.md)

---

### 4. Error Handling 🆕 → [New Material](./theory/01_error_handling.md)

**Нові теми:**
- `errors.Is()` vs `==`
- `errors.As()` для type assertion
- `errors.Unwrap()`
- Custom error types
- Wrapping errors (`fmt.Errorf("%w", err)`)
- Sentinel errors

**Приклад:**
```go
import "errors"

var ErrNotFound = errors.New("not found")

func Get(id int) error {
    if id < 0 {
        return fmt.Errorf("invalid id: %w", ErrNotFound)
    }
    return nil
}

// Check:
if errors.Is(err, ErrNotFound) {
    // Handle not found
}
```

➡️ **Детальніше:** [Week 15 - Error Handling](./theory/01_error_handling.md)

---

### 5. Mock Generation 🆕 → [New Material](./theory/02_mock_generation.md)

**Нові теми:**
- `gomock` + `mockgen`
- Generating mocks from interfaces
- Setting expectations
- Verifying calls
- Argument matchers

**Приклад:**
```bash
# Generate mocks
mockgen -source=service.go -destination=mocks/service_mock.go

# Use in tests
mock := NewMockService(ctrl)
mock.EXPECT().GetUser(1).Return(&User{}, nil)
```

➡️ **Детальніше:** [Week 15 - Mock Generation](./theory/02_mock_generation.md)

---

### 6. Interfaces ✅ → [Week 2](../week_2/theory/)

**Що вже покрито:**
- Implicit implementation (no `implements` keyword)
- Interface segregation
- Empty interface (`interface{}`, `any`)
- Type assertions (`value.(Type)`)
- Type switches

**Приклад:**
```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

// File implements Reader implicitly
type File struct {}

func (f *File) Read(p []byte) (n int, err error) {
    return 0, nil
}
```

➡️ **Детальніше:** [Week 2 - Interfaces](../week_2/theory/)

---

### 7. Database Normalization ✅ → [Week 14](../week_14/theory/04_normalization.md)

**Що вже покрито:**
- **1NF:** Atomic values
- **2NF:** No partial dependencies
- **3NF:** No transitive dependencies
- Anomalies (insert, update, delete)
- Denormalization for performance

**Golden Rule:** "The key, the whole key, and nothing but the key"

➡️ **Детальніше:** [Week 14 - Normalization](../week_14/theory/04_normalization.md)

---

### 8. Database Indexes 🆕 → [New Material](./theory/03_indexes.md)

**Нові теми:**
- Що таке indexes
- B-Tree vs Hash indexes
- **Downsides:**
  - Slower writes (INSERT, UPDATE, DELETE)
  - Extra storage space
  - Index maintenance overhead
  - Too many indexes = query planner confusion

**Trade-offs:**
```
✅ Faster reads (SELECT with WHERE, JOIN, ORDER BY)
❌ Slower writes (INSERT, UPDATE, DELETE)
❌ Extra disk space
❌ Maintenance overhead
```

➡️ **Детальніше:** [Week 15 - Database Indexes](./theory/03_indexes.md)

---

## 🎯 Практика

### [01: Error Handling](./practice/01_error_handling/)

- Custom error types
- `errors.Is` / `errors.As`
- Error wrapping
- Sentinel errors

### [02: Mock Generation](./practice/02_mock_generation/)

- `mockgen` setup
- Generate mocks
- Write tests with mocks
- Verify expectations

### [03: Database Indexes](./practice/03_indexes/)

- Create indexes
- Measure performance
- EXPLAIN query plans
- Index downsides demo

---

## 📊 Quick Reference

### Topics Already Covered

| Topic | Week | Path |
|-------|------|------|
| Maps | Week 1 | [theory/01_types.md](../week_1/theory/01_types.md) |
| Runes | Week 3 | [week_3_rune_bytes/](../week_3_rune_bytes/) |
| Select | Week 6 | [theory/07_goroutines_concurrency.md](../week_6/theory/07_goroutines_concurrency.md) |
| Interfaces | Week 2 | [theory/](../week_2/theory/) |
| Normalization | Week 14 | [theory/04_normalization.md](../week_14/theory/04_normalization.md) |

### New Topics (Week 15)

| Topic | File |
|-------|------|
| Error Handling | [theory/01_error_handling.md](./theory/01_error_handling.md) |
| **errors.Is vs errors.As** | **[ERRORS_IS_VS_AS_GUIDE.md](./ERRORS_IS_VS_AS_GUIDE.md)** ⭐ |
| Mock Generation | [theory/02_mock_generation.md](./theory/02_mock_generation.md) |
| Database Indexes | [theory/03_indexes.md](./theory/03_indexes.md) |

---

## ✅ Checklist

- [ ] Review Maps (Week 1)
- [ ] Review Runes (Week 3)
- [ ] Review Select (Week 6)
- [ ] Learn Error Handling (NEW)
- [ ] Learn Mock Generation (NEW)
- [ ] Review Interfaces (Week 2)
- [ ] Review Normalization (Week 14)
- [ ] Learn Indexes Downsides (NEW)

---

**Week 15: Recap + Advanced Topics!** 🔄🚀
