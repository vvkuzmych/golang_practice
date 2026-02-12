# Quick Start Guide 🚀

Швидкий старт для практики Go concurrency!

---

## 📂 Structure

```
practice/
├── README.md              # Повна документація
├── QUICK_START.md         # Цей файл
├── tasks/                 # Завдання (без рішень)
│   ├── TASK_01_parallel_sum.md
│   ├── TASK_02_url_checker.md
│   ├── TASK_03_worker_pool.md
│   ├── TASK_04_context_timeout.md
│   └── TASK_05_race_condition.md
└── solutions/             # Готові рішення
    ├── solution_01_parallel_sum.go
    ├── solution_02_url_checker.go
    ├── solution_03_worker_pool.go
    ├── solution_04_context_timeout.go
    └── solution_05_race_condition.go
```

---

## ⚡ Quick Commands

### Test Solution

```bash
# Run solution
cd /Users/vkuzm/GolandProjects/golang_practice/interviewtasks/practice
go run solutions/solution_01_parallel_sum.go

# Run with race detector (important!)
go run -race solutions/solution_01_parallel_sum.go
```

### View Task

```bash
# Read task
cat tasks/TASK_01_parallel_sum.md

# Or open in editor
code tasks/TASK_01_parallel_sum.md
```

---

## 🎯 Learning Path

### Beginner (Start Here)

**Task 1: Parallel Sum** - 10 minutes
```bash
cat tasks/TASK_01_parallel_sum.md
go run solutions/solution_01_parallel_sum.go
```

**Key concepts:**
- `sync.WaitGroup`
- `sync.Mutex`
- Goroutines

---

### Intermediate

**Task 2: URL Checker** - 15 minutes
```bash
cat tasks/TASK_02_url_checker.md
go run solutions/solution_02_url_checker.go
```

**Key concepts:**
- Channels
- Error handling
- Order preservation

---

**Task 3: Worker Pool** - 20 minutes
```bash
cat tasks/TASK_03_worker_pool.md
go run solutions/solution_03_worker_pool.go
```

**Key concepts:**
- Worker pool pattern
- Buffered channels
- Job queue

---

### Advanced

**Task 4: Context Timeout** - 15 minutes
```bash
cat tasks/TASK_04_context_timeout.md
go run solutions/solution_04_context_timeout.go
```

**Key concepts:**
- `context.Context`
- Timeout/cancellation
- Graceful shutdown

---

**Task 5: Race Condition** - 15 minutes
```bash
cat tasks/TASK_05_race_condition.md
go run -race solutions/solution_05_race_condition.go
```

**Key concepts:**
- Race detection
- `sync.Mutex` vs `sync.RWMutex`
- Thread safety

---

## 🧪 Test All Solutions

```bash
# Test all (if tests exist)
go test ./...

# With race detector
go test -race ./...

# Verbose
go test -v ./...
```

---

## 🔥 Challenge Yourself

### Step 1: Read Task (DON'T look at solution!)

```bash
cat tasks/TASK_01_parallel_sum.md
```

### Step 2: Create Your Solution

```bash
touch my_solution.go
```

### Step 3: Test Your Solution

```bash
go run my_solution.go
go run -race my_solution.go  # Check for race conditions
```

### Step 4: Compare with Official Solution

```bash
cat solutions/solution_01_parallel_sum.go
```

---

## 💡 Tips

### Always Use Race Detector

```bash
# This can save you hours of debugging!
go run -race your_code.go
```

### Common Mistakes to Avoid

1. **Forgetting `defer wg.Done()`**
2. **Not closing channels**
3. **Capturing loop variables in goroutines**
4. **Race conditions on shared variables**

---

## 📊 Your Progress

Track your progress:

- [ ] Task 1: Parallel Sum
- [ ] Task 2: URL Checker
- [ ] Task 3: Worker Pool
- [ ] Task 4: Context Timeout
- [ ] Task 5: Race Condition

---

## 🎓 After Completion

You'll be ready for:

- ✅ Technical interviews на позиції Go Developer
- ✅ Writing production-grade concurrent code
- ✅ Understanding Go concurrency patterns
- ✅ Debugging race conditions

---

## 🚀 Run Everything Now

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/interviewtasks/practice

# Run all solutions
for file in solutions/*.go; do
    echo "Running $file..."
    go run "$file"
    echo "---"
done
```

---

**Ready? Start with Task 1!** 💪

```bash
cat tasks/TASK_01_parallel_sum.md
```
