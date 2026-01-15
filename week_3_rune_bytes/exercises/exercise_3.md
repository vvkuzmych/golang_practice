# Вправа 3: Status Management System - Система управління статусами

## Ціль
Створити систему управління статусами з використанням iota та enum patterns.

---

## Завдання

Створіть програму `status_system.go`, яка:

1. Використовує iota для створення enum
2. Реалізує систему статусів задач
3. Підтримує permissions з bit flags
4. Має log levels з пріоритетами

---

## Вимоги

### 1. Task Status Enum

```go
type TaskStatus int

const (
    TaskPending TaskStatus = iota
    TaskInProgress
    TaskReview
    TaskCompleted
    TaskCancelled
)

// String() для красивого виводу
func (ts TaskStatus) String() string

// IsTerminal перевіряє чи статус фінальний
func (ts TaskStatus) IsTerminal() bool
```

### 2. Permission Bit Flags

```go
type Permission uint

const (
    PermRead Permission = 1 << iota
    PermWrite
    PermExecute
    PermDelete
    PermAdmin
)

// String() для виводу всіх прав
func (p Permission) String() string

// Has перевіряє чи є право
func (p Permission) Has(perm Permission) bool

// Add додає право
func (p Permission) Add(perm Permission) Permission

// Remove видаляє право
func (p Permission) Remove(perm Permission) Permission
```

### 3. Log Level

```go
type LogLevel int

const (
    LogTrace LogLevel = iota
    LogDebug
    LogInfo
    LogWarning
    LogError
    LogFatal
)

// String() для виводу
func (ll LogLevel) String() string

// ShouldLog перевіряє чи логувати
func (ll LogLevel) ShouldLog(currentLevel LogLevel) bool
```

### 4. Priority

```go
type Priority int

const (
    PriorityLow Priority = 1
    PriorityMedium Priority = 5
    PriorityHigh Priority = 10
    PriorityCritical Priority = 100
)
```

---

## Приклад використання

```go
func main() {
    // Task Status
    task := Task{
        ID: 1,
        Title: "Implement feature",
        Status: TaskPending,
    }
    
    task.UpdateStatus(TaskInProgress)
    fmt.Printf("Task status: %s\n", task.Status)
    
    if task.Status.IsTerminal() {
        fmt.Println("Task is done!")
    }
    
    // Permissions
    user := User{
        Name: "John",
        Perms: PermRead | PermWrite,
    }
    
    if user.Perms.Has(PermWrite) {
        fmt.Println("User can write")
    }
    
    user.Perms = user.Perms.Add(PermExecute)
    fmt.Printf("User permissions: %s\n", user.Perms)
    
    // Logging
    logger := NewLogger(LogWarning)
    logger.Trace("This won't show")
    logger.Warning("This will show")
    logger.Error("Error occurred!")
}
```

---

## Очікуваний вивід

```
=== Task Status System ===
Creating task: "Implement login"
Status: Pending

Updating status: Pending -> InProgress
✅ Status updated

Updating status: InProgress -> Review  
✅ Status updated

Updating status: Review -> Completed
✅ Status updated
Task is now terminal: true

Trying to update completed task...
❌ Cannot update terminal task

=== Permission System ===
User: John
Initial permissions: [Read, Write]

Adding Execute permission...
New permissions: [Read, Write, Execute]

Checking permissions:
  ✅ Has Read
  ✅ Has Write
  ✅ Has Execute
  ❌ No Delete permission

Admin user permissions: [Read, Write, Execute, Delete, Admin]

=== Log Level System ===
Current log level: Warning

[TRACE] Trace message         ← не показується
[DEBUG] Debug message         ← не показується  
[INFO] Info message           ← не показується
[WARNING] Warning message     ← показується
[ERROR] Error occurred        ← показується
[FATAL] Fatal error           ← показується

=== Priority System ===
Task priorities:
  Bug fix: Low (1)
  Feature: Medium (5)
  Security: High (10)
  Outage: Critical (100)

Tasks sorted by priority:
  1. [CRITICAL] System outage (100)
  2. [HIGH] Security patch (10)
  3. [MEDIUM] New feature (5)
  4. [LOW] Bug fix (1)
```

---

## Підказки

### 1. TaskStatus String()
```go
func (ts TaskStatus) String() string {
    statuses := []string{
        "Pending",
        "InProgress",
        "Review",
        "Completed",
        "Cancelled",
    }
    if ts < 0 || int(ts) >= len(statuses) {
        return "Unknown"
    }
    return statuses[ts]
}
```

### 2. Permission Has/Add/Remove
```go
func (p Permission) Has(perm Permission) bool {
    return p&perm != 0
}

func (p Permission) Add(perm Permission) Permission {
    return p | perm
}

func (p Permission) Remove(perm Permission) Permission {
    return p &^ perm
}
```

### 3. Permission String()
```go
func (p Permission) String() string {
    var perms []string
    if p&PermRead != 0 {
        perms = append(perms, "Read")
    }
    if p&PermWrite != 0 {
        perms = append(perms, "Write")
    }
    // ... інші
    return fmt.Sprintf("[%s]", strings.Join(perms, ", "))
}
```

### 4. Logger Implementation
```go
type Logger struct {
    level LogLevel
}

func (l *Logger) log(level LogLevel, message string) {
    if level >= l.level {
        fmt.Printf("[%s] %s\n", level, message)
    }
}

func (l *Logger) Error(msg string) {
    l.log(LogError, msg)
}
```

---

## Структури

```go
type Task struct {
    ID       int
    Title    string
    Status   TaskStatus
    Priority Priority
}

type User struct {
    Name  string
    Perms Permission
}

type Logger struct {
    level LogLevel
}
```

---

## Бонус завдання

1. **HTTP Status Codes**:
   ```go
   const (
       StatusOK HTTPStatus = 200 + iota
       StatusCreated
       StatusAccepted
   )
   ```

2. **File Mode Permissions**:
   ```go
   type FileMode uint
   const (
       ModeRead FileMode = 0400 + iota*0100
       ModeWrite
       ModeExecute
   )
   ```

3. **State Machine**:
   ```go
   func (ts TaskStatus) CanTransitionTo(target TaskStatus) bool
   ```
   Перевірка чи можливий перехід між статусами

4. **Size Units**:
   ```go
   const (
       KB Size = 1 << (10 * iota)
       MB
       GB
       TB
   )
   ```

5. **Color Codes**:
   ```go
   type Color int
   const (
       Red Color = iota
       Green
       Blue
   )
   func (c Color) RGB() (r, g, b int)
   ```

---

## Критерії оцінки

- ✅ TaskStatus enum реалізований
- ✅ Permission bit flags працюють
- ✅ LogLevel з пріоритетами
- ✅ String() методи для всіх enum
- ✅ Has/Add/Remove для permissions
- ✅ IsTerminal для статусів
- ✅ Код чистий і зрозумілий

---

## Тестування

```go
// Test permissions
perms := PermRead | PermWrite
assert(perms.Has(PermRead))
assert(perms.Has(PermWrite))
assert(!perms.Has(PermExecute))

perms = perms.Add(PermExecute)
assert(perms.Has(PermExecute))

perms = perms.Remove(PermWrite)
assert(!perms.Has(PermWrite))

// Test status transitions
status := TaskPending
status = TaskInProgress
status = TaskCompleted
assert(status.IsTerminal())

// Test logging
logger := NewLogger(LogWarning)
logger.Debug("Not shown")  // won't print
logger.Error("Shown")      // will print
```

---

## Важливі моменти

### Bit Operations:
```go
// AND (&) - check if bit is set
if perms & PermRead != 0 {  // has read?
}

// OR (|) - set bit
perms = perms | PermWrite  // add write

// XOR (^) - toggle bit
perms = perms ^ PermExecute  // toggle execute

// AND NOT (&^) - clear bit
perms = perms &^ PermDelete  // remove delete
```

### iota Reset:
```go
const (
    A = iota  // 0
    B         // 1
)

const (
    C = iota  // 0 (reset!)
    D         // 1
)
```

---

## Корисні пакети

- `fmt` - форматування
- `strings` - робота з strings
- `sort` - сортування за priority

---

## Рішення

Рішення знаходиться в `solutions/solution_3.go`.

Спробуйте виконати завдання самостійно перед тим, як дивитись рішення!

---

## Навчальні цілі

Після виконання цієї вправи ви будете:
- ✅ Розуміти iota та enum patterns
- ✅ Вміти використовувати bit flags
- ✅ Створювати type-safe enumerations
- ✅ Реалізовувати String() для enum
- ✅ Розуміти bit operations

---

## Реальне застосування

Подібні patterns використовуються в:
- **Status management** - task tracking, orders, etc.
- **Permission systems** - file permissions, user roles
- **Logging** - log levels, filtering
- **File modes** - Unix file permissions
- **HTTP codes** - status codes
- **State machines** - workflow management

