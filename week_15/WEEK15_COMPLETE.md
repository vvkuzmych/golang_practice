# Week 15 - Completion Report

## ✅ Module Complete: Recap & Advanced Topics

**Created:** 2026-01-28  
**Status:** ✅ Complete  
**Type:** Recap + Advanced  

---

## 📦 Структура

```
week_15/
├── README.md                          # Огляд модулю
├── QUICK_START.md                     # Швидкий старт
├── WEEK15_COMPLETE.md                 # Цей файл
├── ERRORS_IS_VS_AS_GUIDE.md          # ⭐ Quick reference guide
├── theory/
│   ├── 01_error_handling.md          # ✅ Error handling (NEW)
│   ├── 02_mock_generation.md         # ✅ Mock generation (NEW)
│   └── 03_indexes.md                 # ✅ Database indexes (NEW)
└── practice/
    ├── 01_error_handling/
    │   ├── main.go                    # ✅ 8 Examples
    │   ├── main_test.go               # ✅ Unit tests
    │   └── README.md                  # ✅ Docs
    ├── 02_mock_generation/            # (Ready for practice)
    └── 03_indexes/                    # (Ready for practice)
```

---

## 📚 Topics Covered

### Existing Topics (Cross-Referenced)

#### 1. Maps ✅ → [Week 1](../week_1/theory/01_types.md)
- Valid key types (comparable)
- Invalid key types (slice, map, function)
- Runtime errors (nil map assignment)
- Initialization (`make`, literal syntax)
- Operations (set, get, delete, check existence)

**Key Learning:**
```go
// ❌ PANIC: assignment to entry in nil map
var m map[string]int
m["key"] = 1

// ✅ OK
m = make(map[string]int)
m["key"] = 1
```

---

#### 2. Runes ✅ → [Week 3](../week_3_rune_bytes/)
- `rune` type (`int32` for Unicode code point)
- UTF-8 encoding
- `string` vs `[]rune` vs `[]byte`
- `len(string)` returns bytes, not characters
- Iterating over strings correctly

**Key Learning:**
```go
s := "Привіт"
fmt.Println(len(s))         // 12 (bytes)
fmt.Println(len([]rune(s))) // 6 (characters)
```

---

#### 3. Select Statement ✅ → [Week 6](../week_6/theory/07_goroutines_concurrency.md)
- Multiple channel operations
- Timeout pattern with `time.After()`
- Non-blocking with `default` case
- Random selection when multiple channels ready

**Key Learning:**
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

---

#### 4. Interfaces ✅ → [Week 2](../week_2/theory/)
- Implicit implementation (no `implements` keyword)
- Interface segregation principle
- Empty interface (`interface{}`, `any`)
- Type assertions
- Type switches

**Key Learning:**
```go
// No need to declare "implements Reader"
type File struct {}

func (f *File) Read(p []byte) (n int, err error) {
    return 0, nil
}

// File automatically satisfies io.Reader
```

---

#### 5. Database Normalization ✅ → [Week 14](../week_14/theory/04_normalization.md)
- **1NF:** Atomic values, no repeating groups
- **2NF:** 1NF + no partial dependencies
- **3NF:** 2NF + no transitive dependencies
- **BCNF:** Advanced normal form
- Data anomalies (insert, update, delete)
- When to denormalize

**Key Learning:**
```
"The key, the whole key, and nothing but the key, so help me Codd"
```

---

### New Topics (Week 15)

#### 6. Error Handling 🆕
**File:** `theory/01_error_handling.md`

**Topics:**
- `error` interface
- Sentinel errors (`var ErrNotFound = errors.New(...)`)
- Custom error types
- `errors.Is()` vs `==`
- `errors.As()` for type assertion
- Error wrapping with `fmt.Errorf("%w", err)`
- `errors.Unwrap()`
- Best practices

**Key Concepts:**

✅ **Sentinel Errors:**
```go
var ErrNotFound = errors.New("not found")

if errors.Is(err, ErrNotFound) {  // ✅ Works with wrapped
    // Handle
}
```

✅ **Custom Types:**
```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation: %s - %s", e.Field, e.Message)
}
```

✅ **errors.As():**
```go
var valErr *ValidationError
if errors.As(err, &valErr) {
    fmt.Println(valErr.Field)  // Extract field
}
```

✅ **Wrapping:**
```go
return fmt.Errorf("failed to save user: %w", err)

// Check wrapped error
if errors.Is(err, ErrNotFound) {  // ✅ Still works!
    // ...
}
```

**Practice:** `practice/01_error_handling/`
- ✅ 8 complete examples
- ✅ Unit tests with `errors.Is/As`
- ✅ Benchmarks
- ✅ Real-world API example

---

#### 7. Mock Generation 🆕
**File:** `theory/02_mock_generation.md`

**Topics:**
- What are mocks and why use them
- `gomock` + `mockgen` setup
- Generating mocks from interfaces
- Setting expectations
- Argument matchers
- Return values and side effects
- Best practices
- Common patterns

**Key Concepts:**

✅ **Generate Mocks:**
```bash
mockgen -source=interface.go -destination=mocks/mock.go -package=mocks
```

✅ **Write Tests:**
```go
ctrl := gomock.NewController(t)
defer ctrl.Finish()

mock := mocks.NewMockUserRepository(ctrl)
mock.EXPECT().GetByID(1).Return(user, nil)

// Use mock
svc := service.NewUserService(mock)
result, err := svc.GetUser(1)
```

✅ **Argument Matchers:**
```go
mock.EXPECT().Save(gomock.Any()).Return(nil)
mock.EXPECT().GetByID(gomock.Not(0)).Return(user, nil)
```

✅ **Call Order:**
```go
gomock.InOrder(
    mock.EXPECT().GetByID(1),
    mock.EXPECT().GetByID(2),
)
```

---

#### 8. Database Indexes 🆕
**File:** `theory/03_indexes.md`

**Topics:**
- What are indexes and how they work
- Types (B-Tree, Hash, Composite, Unique, Partial)
- Advantages (faster reads)
- **Downsides** (slower writes, space, maintenance)
- When to use indexes
- EXPLAIN query plans
- Best practices

**Key Concepts:**

✅ **Advantages:**
```sql
-- Without index: Full table scan (slow)
-- With index: Index lookup (fast)
-- Speedup: 500x for large tables
```

❌ **Downsides:**

1. **Slower Writes:**
```
INSERT without index: 1 operation
INSERT with 3 indexes: 4 operations (30-50% slower per index)
```

2. **Extra Disk Space:**
```
Table: 1 GB
+ 3 indexes: +450 MB
Total: 1.45 GB (45% overhead)
```

3. **Maintenance Overhead:**
- B-Tree rebalancing
- Fragmentation
- VACUUM/OPTIMIZE needed

4. **Query Planner Confusion:**
- Too many indexes → harder to choose
- Wrong choice → slow query

✅ **When to Use:**
```sql
-- Foreign keys ✅
CREATE INDEX idx_orders_user_id ON orders(user_id);

-- Frequent WHERE ✅
CREATE INDEX idx_users_email ON users(email);

-- Low cardinality ❌
-- (gender: only 3 values - not useful)
```

✅ **Measure Performance:**
```sql
EXPLAIN ANALYZE
SELECT * FROM users WHERE email = 'test@example.com';
```

---

## 🎯 Learning Objectives

### After Week 15, you should:

#### Existing Topics Review:
- [ ] Understand which types can be map keys
- [ ] Know the difference between `string`, `[]rune`, and `[]byte`
- [ ] Use `select` for multiple channel operations
- [ ] Understand implicit interface implementation
- [ ] Apply database normalization (1NF, 2NF, 3NF)

#### New Topics Mastery:
- [ ] Use `errors.Is()` instead of `==`
- [ ] Use `errors.As()` for type assertion
- [ ] Wrap errors with `fmt.Errorf("%w")`
- [ ] Create custom error types
- [ ] Generate mocks with `mockgen`
- [ ] Write tests with `gomock`
- [ ] Understand index trade-offs
- [ ] Know when to use indexes
- [ ] Measure query performance with EXPLAIN

---

## 💻 Practical Examples

### Error Handling

8 complete examples demonstrating:

1. ✅ Basic error handling
2. ✅ Sentinel errors (old way with `==`)
3. ✅ `errors.Is()` (new way - works with wrapped)
4. ✅ `errors.As()` for custom types
5. ✅ Wrapped errors with `Unwrap()`
6. ✅ Error wrapping chain
7. ✅ Multiple error checks
8. ✅ Real-world API example

**Run:**
```bash
cd practice/01_error_handling
go run main.go
go test -v
```

---

## 📊 Cross-Reference Matrix

| Topic | Week | Theory | Practice | Status |
|-------|------|--------|----------|--------|
| Maps | 1 | `week_1/theory/01_types.md` | `week_1/practice/` | ✅ |
| Runes | 3 | `week_3_rune_bytes/theory/` | `week_3_rune_bytes/practice/` | ✅ |
| Select | 6 | `week_6/theory/07_goroutines_concurrency.md` | `week_6/practice/06_goroutines/` | ✅ |
| Interfaces | 2 | `week_2/theory/` | `week_2/practice/` | ✅ |
| Normalization | 14 | `week_14/theory/04_normalization.md` | `week_14/practice/` | ✅ |
| **Error Handling** | **15** | **`theory/01_error_handling.md`** | **`practice/01_error_handling/`** | **✅** |
| **errors.Is vs As** | **15** | **`ERRORS_IS_VS_AS_GUIDE.md`** ⭐ | **`practice/01_error_handling/`** | **✅** |
| **Mock Generation** | **15** | **`theory/02_mock_generation.md`** | **`practice/02_mock_generation/`** | **📝 Theory** |
| **Indexes** | **15** | **`theory/03_indexes.md`** | **`practice/03_indexes/`** | **📝 Theory** |

---

## 🎓 Best Practices Summary

### Error Handling
1. ✅ Always use `errors.Is()` instead of `==`
2. ✅ Use `errors.As()` for type assertions
3. ✅ Wrap errors with `%w` for context
4. ✅ Export sentinel errors for callers
5. ✅ Implement `Unwrap()` for custom error chains

### Mock Generation
1. ✅ Mock interfaces, not concrete types
2. ✅ One mock per test
3. ✅ Always `defer ctrl.Finish()`
4. ✅ Test behavior, not implementation
5. ✅ Don't over-mock simple logic

### Database Indexes
1. ✅ Always index foreign keys
2. ✅ Index frequently queried WHERE columns
3. ✅ Use composite indexes for multi-column queries
4. ✅ Don't over-index (every index has a cost)
5. ✅ Monitor with EXPLAIN

---

## 📚 Additional Resources

### Error Handling
- ⭐ **Quick Guide:** [ERRORS_IS_VS_AS_GUIDE.md](./ERRORS_IS_VS_AS_GUIDE.md) - Complete comparison with examples
- Go Blog: [Error Handling and Go](https://go.dev/blog/error-handling-and-go)
- Go Blog: [Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors)

### Mock Generation
- GitHub: [golang/mock](https://github.com/golang/mock)
- Tutorial: [Mocking in Go with gomock](https://bmuschko.com/blog/go-testing-mocking/)

### Database Indexes
- PostgreSQL Docs: [Indexes](https://www.postgresql.org/docs/current/indexes.html)
- Use The Index, Luke: [Guide to Database Performance](https://use-the-index-luke.com/)

---

## 🎯 Recommended Learning Path

### Day 1: Review Existing Topics
1. Read Week 1 - Maps (30 min)
2. Read Week 3 - Runes (30 min)
3. Read Week 6 - Select (30 min)
4. Quick review of Interfaces & Normalization (30 min)

### Day 2: Error Handling
1. Read `theory/01_error_handling.md` (1 hour)
2. Run `practice/01_error_handling/main.go` (30 min)
3. Study tests `practice/01_error_handling/main_test.go` (30 min)
4. Implement own examples (1 hour)

### Day 3: Mock Generation
1. Read `theory/02_mock_generation.md` (1 hour)
2. Install `mockgen` (10 min)
3. Create own interface and generate mocks (1 hour)
4. Write tests with mocks (1 hour)

### Day 4: Database Indexes
1. Read `theory/03_indexes.md` (1 hour)
2. Practice EXPLAIN on real queries (1 hour)
3. Experiment with index creation (1 hour)
4. Measure before/after performance (1 hour)

### Day 5: Integration
1. Apply error handling to real project (2 hours)
2. Add mocks to existing tests (2 hours)
3. Optimize database queries with indexes (2 hours)

---

## ✅ Completion Checklist

### Documentation
- [x] README.md with overview and cross-references
- [x] theory/01_error_handling.md
- [x] ERRORS_IS_VS_AS_GUIDE.md ⭐ (Quick reference)
- [x] theory/02_mock_generation.md
- [x] theory/03_indexes.md
- [x] QUICK_START.md
- [x] WEEK15_COMPLETE.md

### Practice
- [x] practice/01_error_handling/main.go (8 examples)
- [x] practice/01_error_handling/main_test.go
- [x] practice/01_error_handling/README.md
- [ ] practice/02_mock_generation/ (theory ready, practice pending)
- [ ] practice/03_indexes/ (theory ready, practice pending)

### Cross-References
- [x] Maps → Week 1
- [x] Runes → Week 3
- [x] Select → Week 6
- [x] Interfaces → Week 2
- [x] Normalization → Week 14

---

## 🎊 Summary

**Week 15** успішно поєднує:
- ✅ **Повторення** 5 existing topics (Maps, Runes, Select, Interfaces, Normalization)
- ✅ **Додавання** 3 advanced topics (Error Handling, Mock Generation, Indexes)
- ✅ **Практика** з error handling (8 examples + tests)
- ✅ **Посилання** на всі попередні тижні

**Total Content:**
- 📄 3 нових теоретичних файлів
- ⭐ 1 quick reference guide (errors.Is vs errors.As)
- 💻 1 повний практичний проект (Error Handling)
- 🔗 5 cross-references до попередніх тижнів
- 📚 2 guide files (README, QUICK_START)
- 📊 1 completion report (цей файл)

**Week 15 Module: Complete!** ✅🎓🚀

---

**Created:** 2026-01-28  
**Status:** ✅ Complete  
**Next Steps:** Practice Mock Generation & Indexes

**Week 15: Recap + Advanced Topics Master!** 🔄⚡✨
