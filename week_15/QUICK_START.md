# Week 15 - Quick Start

## 🚀 Швидкий старт

### Error Handling Practice

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_15/practice/01_error_handling

# Запустити всі приклади
go run main.go

# Запустити тести
go test -v

# Benchmark
go test -bench=.
```

### Огляд тем

#### 1. Maps (Week 1)
```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_1/theory
cat 01_types.md  # Секція "Мапи (map)"
```

**Key points:**
- Valid key types: int, string, bool, pointer, comparable struct
- Invalid key types: slice, map, function
- Runtime error: assignment to nil map

#### 2. Runes (Week 3)
```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_3_rune_bytes
ls theory/
```

**Key points:**
- `rune` = `int32` (Unicode code point)
- `len(string)` = bytes, not characters
- Use `[]rune` for character iteration

#### 3. Select Statement (Week 6)
```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_6/theory
cat 07_goroutines_concurrency.md  # Секція "Select"
```

**Key points:**
- Multiple channel operations
- Timeout pattern
- Non-blocking with `default`

#### 4. Error Handling (NEW)
```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_15/theory
cat 01_error_handling.md

# Quick reference guide
cat ../ERRORS_IS_VS_AS_GUIDE.md
```

**Key points:**
- `errors.Is()` for sentinel errors (checks WHAT)
- `errors.As()` for custom types (checks TYPE)
- Error wrapping with `fmt.Errorf("%w")`

**Quick Guide:** See [ERRORS_IS_VS_AS_GUIDE.md](../ERRORS_IS_VS_AS_GUIDE.md) ⭐

#### 5. Mock Generation (NEW)
```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_15/theory
cat 02_mock_generation.md
```

**Key points:**
- Install: `go install github.com/golang/mock/mockgen@latest`
- Generate: `mockgen -source=file.go -destination=mocks/mock.go`
- Use: `ctrl := gomock.NewController(t)`

#### 6. Interfaces (Week 2)
```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_2/theory
ls *.md
```

**Key points:**
- Implicit implementation
- No `implements` keyword
- Empty interface: `interface{}` or `any`

#### 7. Database Normalization (Week 14)
```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_14
cat theory/04_normalization.md
cat NORMALIZATION_CHEAT_SHEET.md
```

**Key points:**
- 1NF: Atomic values
- 2NF: No partial dependencies
- 3NF: No transitive dependencies

#### 8. Database Indexes (NEW)
```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_15/theory
cat 03_indexes.md
```

**Key points:**
- ✅ Faster reads (SELECT, JOIN, ORDER BY)
- ❌ Slower writes (INSERT, UPDATE, DELETE)
- ❌ Extra disk space
- ❌ Maintenance overhead

---

## 📊 Quick Test

### Test Your Knowledge

```bash
# 1. Error Handling
cd week_15/practice/01_error_handling
go run main.go

# 2. Review Maps
cd ../../week_1/theory
grep -A 20 "### 8. Мапи" 01_types.md

# 3. Review Runes
cd ../../week_3_rune_bytes
cat theory/*.md

# 4. Review Select
cd ../week_6/theory
grep -A 30 "### 5. Select" 07_goroutines_concurrency.md

# 5. Review Interfaces
cd ../../week_2/theory
ls *.md

# 6. Review Normalization
cd ../../week_14
cat NORMALIZATION_CHEAT_SHEET.md

# 7. Read about Indexes
cd ../week_15/theory
cat 03_indexes.md
```

---

## 🎯 Чеклист

### Existing Topics (Review)

- [ ] Maps - What types can be keys? (Week 1)
- [ ] Runes - What's the difference from bytes? (Week 3)
- [ ] Select - How to timeout a channel? (Week 6)
- [ ] Interfaces - Implicit vs explicit? (Week 2)
- [ ] Normalization - What's 3NF? (Week 14)

### New Topics (Learn)

- [ ] errors.Is() vs == 
- [ ] errors.As() usage
- [ ] Error wrapping with %w
- [ ] mockgen installation
- [ ] Mock generation from interface
- [ ] Setting expectations
- [ ] Index advantages
- [ ] Index downsides (writes, space, maintenance)
- [ ] When to use indexes
- [ ] EXPLAIN query plans

---

## 📚 Cross-References

| Topic | Week | Quick Access |
|-------|------|--------------|
| Maps | 1 | `cat week_1/theory/01_types.md` |
| Runes | 3 | `cd week_3_rune_bytes` |
| Select | 6 | `cat week_6/theory/07_goroutines_concurrency.md` |
| Interfaces | 2 | `cd week_2/theory` |
| Normalization | 14 | `cat week_14/NORMALIZATION_CHEAT_SHEET.md` |
| Error Handling | 15 | `cat week_15/theory/01_error_handling.md` |
| **errors.Is vs As** | **15** | **`cat week_15/ERRORS_IS_VS_AS_GUIDE.md`** ⭐ |
| Mock Generation | 15 | `cat week_15/theory/02_mock_generation.md` |
| Indexes | 15 | `cat week_15/theory/03_indexes.md` |

---

## ⚡ Quick Examples

### errors.Is() - Check WHAT (value)
```go
if errors.Is(err, ErrNotFound) {
    // Handle not found
}
```

### errors.As() - Check TYPE (and extract)
```go
var valErr *ValidationError
if errors.As(err, &valErr) {
    fmt.Println(valErr.Field)  // Access fields
}
```

### Complete Guide
See [ERRORS_IS_VS_AS_GUIDE.md](./ERRORS_IS_VS_AS_GUIDE.md) for detailed comparison and examples! ⭐

### mockgen
```bash
mockgen -source=service.go -destination=mocks/service_mock.go -package=mocks
```

### EXPLAIN
```sql
EXPLAIN ANALYZE SELECT * FROM users WHERE email = 'test@example.com';
```

---

**Week 15: Recap + Advanced!** 🔄🚀
