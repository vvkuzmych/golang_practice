package main

import (
	"fmt"
	"strings"
)

// ============= Task Status Enum =============

type TaskStatus int

const (
	TaskPending TaskStatus = iota
	TaskInProgress
	TaskReview
	TaskCompleted
	TaskCancelled
)

// String повертає назву статусу
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

// IsTerminal перевіряє чи статус фінальний
func (ts TaskStatus) IsTerminal() bool {
	return ts == TaskCompleted || ts == TaskCancelled
}

// CanTransitionTo перевіряє чи можливий перехід
func (ts TaskStatus) CanTransitionTo(target TaskStatus) bool {
	if ts.IsTerminal() {
		return false
	}
	transitions := map[TaskStatus][]TaskStatus{
		TaskPending:    {TaskInProgress, TaskCancelled},
		TaskInProgress: {TaskReview, TaskCancelled},
		TaskReview:     {TaskCompleted, TaskInProgress, TaskCancelled},
	}
	allowed := transitions[ts]
	for _, s := range allowed {
		if s == target {
			return true
		}
	}
	return false
}

// ============= Permission Bit Flags =============

type Permission uint

const (
	PermRead Permission = 1 << iota
	PermWrite
	PermExecute
	PermDelete
	PermAdmin
)

// String повертає список всіх прав
func (p Permission) String() string {
	var perms []string

	if p&PermRead != 0 {
		perms = append(perms, "Read")
	}
	if p&PermWrite != 0 {
		perms = append(perms, "Write")
	}
	if p&PermExecute != 0 {
		perms = append(perms, "Execute")
	}
	if p&PermDelete != 0 {
		perms = append(perms, "Delete")
	}
	if p&PermAdmin != 0 {
		perms = append(perms, "Admin")
	}

	if len(perms) == 0 {
		return "[None]"
	}

	return fmt.Sprintf("[%s]", strings.Join(perms, ", "))
}

// Has перевіряє чи є право
func (p Permission) Has(perm Permission) bool {
	return p&perm != 0
}

// Add додає право
func (p Permission) Add(perm Permission) Permission {
	return p | perm
}

// Remove видаляє право
func (p Permission) Remove(perm Permission) Permission {
	return p &^ perm
}

// ============= Log Level =============

type LogLevel int

const (
	LogTrace LogLevel = iota
	LogDebug
	LogInfo
	LogWarning
	LogError
	LogFatal
)

// String повертає назву рівня
func (ll LogLevel) String() string {
	levels := []string{"TRACE", "DEBUG", "INFO", "WARNING", "ERROR", "FATAL"}
	if ll < 0 || int(ll) >= len(levels) {
		return "UNKNOWN"
	}
	return levels[ll]
}

// ShouldLog перевіряє чи логувати
func (ll LogLevel) ShouldLog(currentLevel LogLevel) bool {
	return ll >= currentLevel
}

// ============= Priority =============

type Priority int

const (
	PriorityLow      Priority = 1
	PriorityMedium   Priority = 5
	PriorityHigh     Priority = 10
	PriorityCritical Priority = 100
)

// String повертає назву пріоритету
func (p Priority) String() string {
	switch p {
	case PriorityLow:
		return "Low"
	case PriorityMedium:
		return "Medium"
	case PriorityHigh:
		return "High"
	case PriorityCritical:
		return "Critical"
	default:
		return "Unknown"
	}
}

// ============= Structures =============

type Task struct {
	ID       int
	Title    string
	Status   TaskStatus
	Priority Priority
}

// UpdateStatus оновлює статус задачі
func (t *Task) UpdateStatus(newStatus TaskStatus) error {
	if t.Status.IsTerminal() {
		return fmt.Errorf("cannot update terminal task")
	}
	if !t.Status.CanTransitionTo(newStatus) {
		return fmt.Errorf("invalid transition: %s -> %s", t.Status, newStatus)
	}
	t.Status = newStatus
	return nil
}

type User struct {
	Name  string
	Perms Permission
}

type Logger struct {
	level LogLevel
}

// NewLogger створює новий logger
func NewLogger(level LogLevel) *Logger {
	return &Logger{level: level}
}

// log виводить повідомлення якщо рівень дозволяє
func (l *Logger) log(level LogLevel, message string) {
	if level.ShouldLog(l.level) {
		fmt.Printf("[%s] %s\n", level, message)
	}
}

func (l *Logger) Trace(msg string)   { l.log(LogTrace, msg) }
func (l *Logger) Debug(msg string)   { l.log(LogDebug, msg) }
func (l *Logger) Info(msg string)    { l.log(LogInfo, msg) }
func (l *Logger) Warning(msg string) { l.log(LogWarning, msg) }
func (l *Logger) Error(msg string)   { l.log(LogError, msg) }
func (l *Logger) Fatal(msg string)   { l.log(LogFatal, msg) }

// ============= Helper Functions =============

func printHeader(title string) {
	fmt.Printf("\n🔹 %s\n", title)
	fmt.Println(strings.Repeat("─", 50))
}

// ============= Main =============

func main() {
	fmt.Println("╔════════════════════════════════════════════════╗")
	fmt.Println("║      Status Management System Solution         ║")
	fmt.Println("╚════════════════════════════════════════════════╝")

	// ===== Task Status System =====
	printHeader("Task Status System")

	task := Task{
		ID:       1,
		Title:    "Implement login",
		Status:   TaskPending,
		Priority: PriorityHigh,
	}

	fmt.Printf("Created task: \"%s\"\n", task.Title)
	fmt.Printf("Status: %s\n", task.Status)
	fmt.Printf("Priority: %s (%d)\n", task.Priority, task.Priority)
	fmt.Printf("Is terminal: %v\n", task.Status.IsTerminal())

	// Status transitions
	fmt.Println("\nStatus transitions:")

	transitions := []TaskStatus{
		TaskInProgress,
		TaskReview,
		TaskCompleted,
	}

	for _, newStatus := range transitions {
		err := task.UpdateStatus(newStatus)
		if err != nil {
			fmt.Printf("  ❌ %s -> %s: %v\n", task.Status, newStatus, err)
		} else {
			fmt.Printf("  ✅ %s\n", task.Status)
		}
	}

	fmt.Printf("\nTask is now terminal: %v\n", task.Status.IsTerminal())

	// Try to update completed task
	fmt.Println("\nTrying to update completed task...")
	err := task.UpdateStatus(TaskInProgress)
	if err != nil {
		fmt.Printf("❌ %v\n", err)
	}

	// ===== Permission System =====
	printHeader("Permission System")

	user := User{
		Name:  "John",
		Perms: PermRead | PermWrite,
	}

	fmt.Printf("User: %s\n", user.Name)
	fmt.Printf("Initial permissions: %s\n", user.Perms)

	// Check permissions
	fmt.Println("\nChecking permissions:")
	checks := []struct {
		perm Permission
		name string
	}{
		{PermRead, "Read"},
		{PermWrite, "Write"},
		{PermExecute, "Execute"},
		{PermDelete, "Delete"},
	}

	for _, check := range checks {
		if user.Perms.Has(check.perm) {
			fmt.Printf("  ✅ Has %s\n", check.name)
		} else {
			fmt.Printf("  ❌ No %s permission\n", check.name)
		}
	}

	// Add permission
	fmt.Println("\nAdding Execute permission...")
	user.Perms = user.Perms.Add(PermExecute)
	fmt.Printf("New permissions: %s\n", user.Perms)

	// Remove permission
	fmt.Println("\nRemoving Write permission...")
	user.Perms = user.Perms.Remove(PermWrite)
	fmt.Printf("New permissions: %s\n", user.Perms)

	// Admin user
	fmt.Println("\nAdmin user:")
	admin := User{
		Name:  "Admin",
		Perms: PermRead | PermWrite | PermExecute | PermDelete | PermAdmin,
	}
	fmt.Printf("Permissions: %s\n", admin.Perms)

	// ===== Permission Bit Values =====
	printHeader("Permission Bit Values")

	fmt.Printf("PermRead:    %08b (%d)\n", PermRead, PermRead)
	fmt.Printf("PermWrite:   %08b (%d)\n", PermWrite, PermWrite)
	fmt.Printf("PermExecute: %08b (%d)\n", PermExecute, PermExecute)
	fmt.Printf("PermDelete:  %08b (%d)\n", PermDelete, PermDelete)
	fmt.Printf("PermAdmin:   %08b (%d)\n", PermAdmin, PermAdmin)

	combined := PermRead | PermWrite | PermExecute
	fmt.Printf("\nCombined (R+W+X): %08b (%d)\n", combined, combined)

	// ===== Log Level System =====
	printHeader("Log Level System")

	fmt.Println("Current log level: Warning\n")
	logger := NewLogger(LogWarning)

	logger.Trace("Trace message")
	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warning("Warning message")
	logger.Error("Error occurred")
	logger.Fatal("Fatal error")

	fmt.Println("\nChanging log level to Debug:")
	logger = NewLogger(LogDebug)

	logger.Trace("Trace message")
	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warning("Warning message")

	// ===== Priority System =====
	printHeader("Priority System")

	tasks := []Task{
		{ID: 1, Title: "Bug fix", Priority: PriorityLow},
		{ID: 2, Title: "New feature", Priority: PriorityMedium},
		{ID: 3, Title: "Security patch", Priority: PriorityHigh},
		{ID: 4, Title: "System outage", Priority: PriorityCritical},
	}

	fmt.Println("Task priorities:")
	for _, t := range tasks {
		fmt.Printf("  %s: %s (%d)\n", t.Title, t.Priority, t.Priority)
	}

	// Sort by priority (simple bubble sort for demo)
	for i := 0; i < len(tasks); i++ {
		for j := i + 1; j < len(tasks); j++ {
			if tasks[i].Priority < tasks[j].Priority {
				tasks[i], tasks[j] = tasks[j], tasks[i]
			}
		}
	}

	fmt.Println("\nTasks sorted by priority:")
	for i, t := range tasks {
		fmt.Printf("  %d. [%s] %s (%d)\n", i+1, t.Priority, t.Title, t.Priority)
	}

	// ===== Enum Values =====
	printHeader("All Enum Values")

	fmt.Println("Task Statuses:")
	statuses := []TaskStatus{
		TaskPending, TaskInProgress, TaskReview, TaskCompleted, TaskCancelled,
	}
	for _, s := range statuses {
		fmt.Printf("  %d: %s (terminal: %v)\n", s, s, s.IsTerminal())
	}

	fmt.Println("\nLog Levels:")
	levels := []LogLevel{
		LogTrace, LogDebug, LogInfo, LogWarning, LogError, LogFatal,
	}
	for _, l := range levels {
		fmt.Printf("  %d: %s\n", l, l)
	}

	fmt.Println("\nPriorities:")
	priorities := []Priority{
		PriorityLow, PriorityMedium, PriorityHigh, PriorityCritical,
	}
	for _, p := range priorities {
		fmt.Printf("  %d: %s\n", p, p)
	}

	// ===== State Machine Example =====
	printHeader("State Machine Example")

	fmt.Println("Valid transitions from each state:")

	stateMap := map[TaskStatus][]TaskStatus{
		TaskPending:    {TaskInProgress, TaskCancelled},
		TaskInProgress: {TaskReview, TaskCancelled},
		TaskReview:     {TaskCompleted, TaskInProgress, TaskCancelled},
		TaskCompleted:  {},
		TaskCancelled:  {},
	}

	for state, transitions := range stateMap {
		fmt.Printf("\n%s can transition to:\n", state)
		if len(transitions) == 0 {
			fmt.Println("  (terminal state)")
		}
		for _, t := range transitions {
			fmt.Printf("  → %s\n", t)
		}
	}

	// ===== Summary =====
	fmt.Println("\n\n📝 ВИСНОВКИ")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println("✅ iota - auto-increment constants")
	fmt.Println("✅ Enum pattern для type-safety")
	fmt.Println("✅ Bit flags для permissions (1 << iota)")
	fmt.Println("✅ String() для красивого виводу")
	fmt.Println("✅ Методи для enum logic")
	fmt.Println("✅ State machines з validation")

	fmt.Println("\nБітові операції:")
	fmt.Println("  & (AND) - check if bit is set")
	fmt.Println("  | (OR) - set bit")
	fmt.Println("  &^ (AND NOT) - clear bit")
	fmt.Println("  ^ (XOR) - toggle bit")
}
