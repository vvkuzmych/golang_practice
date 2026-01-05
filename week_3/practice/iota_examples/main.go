package main

import "fmt"

// ============= Basic Enum =============

type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

func (w Weekday) String() string {
	days := []string{
		"Sunday", "Monday", "Tuesday", "Wednesday",
		"Thursday", "Friday", "Saturday",
	}
	if w < 0 || int(w) >= len(days) {
		return "Unknown"
	}
	return days[w]
}

// ============= HTTP Status =============

type HTTPStatus int

const (
	StatusOK       HTTPStatus = 200 + iota
	StatusCreated             // 201
	StatusAccepted            // 202
)

const (
	StatusBadRequest   HTTPStatus = 400 + iota
	StatusUnauthorized            // 401
	StatusForbidden               // 402
	StatusNotFound                // 403
)

// ============= Bit Flags =============

type Permission uint

const (
	Read    Permission = 1 << iota // 1 << 0 = 1
	Write                          // 1 << 1 = 2
	Execute                        // 1 << 2 = 4
	Delete                         // 1 << 3 = 8
)

func (p Permission) String() string {
	var perms []string
	if p&Read != 0 {
		perms = append(perms, "Read")
	}
	if p&Write != 0 {
		perms = append(perms, "Write")
	}
	if p&Execute != 0 {
		perms = append(perms, "Execute")
	}
	if p&Delete != 0 {
		perms = append(perms, "Delete")
	}
	if len(perms) == 0 {
		return "None"
	}
	return fmt.Sprintf("%v", perms)
}

// ============= Size Units =============

type Size int64

const (
	_       = iota             // ignore first value
	KB Size = 1 << (10 * iota) // 1 << 10 = 1024
	MB                         // 1 << 20
	GB                         // 1 << 30
	TB                         // 1 << 40
)

// ============= Log Level =============

type LogLevel int

const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	WARNING
	ERROR
	FATAL
)

func (l LogLevel) String() string {
	levels := []string{"TRACE", "DEBUG", "INFO", "WARNING", "ERROR", "FATAL"}
	if l < 0 || int(l) >= len(levels) {
		return "UNKNOWN"
	}
	return levels[l]
}

// ============= Priority =============

type Priority int

const (
	Low    Priority = 1
	Medium Priority = 5
	High   Priority = 10
	Urgent Priority = 100
)

// ============= Status =============

type Status int

const (
	StatusPending Status = iota
	StatusActive
	StatusPaused
	StatusCompleted
	StatusCancelled
	StatusFailed
)

// ============= Color =============

type Color int

const (
	Red Color = iota
	Green
	Blue
	_ // пропускаємо значення
	Yellow
	Purple
)

// ============= Main =============

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║         iota & Enum Examples             ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	// ===== Basic iota =====
	fmt.Println("\n🔹 Базовий iota (Weekday)")
	fmt.Println("─────────────────────────────────────────")

	fmt.Printf("Sunday: %d (%s)\n", Sunday, Sunday)
	fmt.Printf("Monday: %d (%s)\n", Monday, Monday)
	fmt.Printf("Saturday: %d (%s)\n", Saturday, Saturday)

	today := Wednesday
	fmt.Printf("\nСьогодні: %s\n", today)

	if today == Wednesday {
		fmt.Println("Середа тижня!")
	}

	// ===== HTTP Status =====
	fmt.Println("\n🔹 HTTP Status коди")
	fmt.Println("─────────────────────────────────────────")

	fmt.Printf("OK: %d\n", StatusOK)
	fmt.Printf("Created: %d\n", StatusCreated)
	fmt.Printf("Bad Request: %d\n", StatusBadRequest)
	fmt.Printf("Unauthorized: %d\n", StatusUnauthorized)
	fmt.Printf("Not Found: %d\n", StatusNotFound)

	// ===== Bit Flags =====
	fmt.Println("\n🔹 Bit Flags (Permissions)")
	fmt.Println("─────────────────────────────────────────")

	fmt.Printf("Read: %d (%b)\n", Read, Read)
	fmt.Printf("Write: %d (%b)\n", Write, Write)
	fmt.Printf("Execute: %d (%b)\n", Execute, Execute)
	fmt.Printf("Delete: %d (%b)\n", Delete, Delete)

	// Комбінації прав
	fmt.Println("\nКомбінації:")

	readWrite := Read | Write
	fmt.Printf("Read+Write: %d (%b) - %s\n", readWrite, readWrite, readWrite)

	fullAccess := Read | Write | Execute | Delete
	fmt.Printf("Full Access: %d (%b) - %s\n", fullAccess, fullAccess, fullAccess)

	// Перевірка прав
	userPerms := Read | Write

	if userPerms&Read != 0 {
		fmt.Println("✅ Є право Read")
	}
	if userPerms&Execute != 0 {
		fmt.Println("✅ Є право Execute")
	} else {
		fmt.Println("❌ Немає права Execute")
	}

	// ===== Size Units =====
	fmt.Println("\n🔹 Size Units")
	fmt.Println("─────────────────────────────────────────")

	fmt.Printf("1 KB = %d bytes\n", KB)
	fmt.Printf("1 MB = %d bytes\n", MB)
	fmt.Printf("1 GB = %d bytes\n", GB)
	fmt.Printf("1 TB = %d bytes\n", TB)

	fileSize := 5 * GB
	fmt.Printf("\nРозмір файлу: %d bytes (%.2f GB)\n",
		fileSize, float64(fileSize)/float64(GB))

	// ===== Log Levels =====
	fmt.Println("\n🔹 Log Levels")
	fmt.Println("─────────────────────────────────────────")

	levels := []LogLevel{TRACE, DEBUG, INFO, WARNING, ERROR, FATAL}

	for _, level := range levels {
		fmt.Printf("[%s] %d\n", level, level)
	}

	// Фільтрація по рівню
	currentLevel := WARNING
	fmt.Printf("\nПоточний рівень: %s\n", currentLevel)
	fmt.Println("Показуємо тільки:")

	for _, level := range levels {
		if level >= currentLevel {
			fmt.Printf("  - %s\n", level)
		}
	}

	// ===== Priority =====
	fmt.Println("\n🔹 Priority (custom values)")
	fmt.Println("─────────────────────────────────────────")

	fmt.Printf("Low: %d\n", Low)
	fmt.Printf("Medium: %d\n", Medium)
	fmt.Printf("High: %d\n", High)
	fmt.Printf("Urgent: %d\n", Urgent)

	taskPriority := High
	if taskPriority >= High {
		fmt.Println("\n⚠️  Високий пріоритет!")
	}

	// ===== Status Flow =====
	fmt.Println("\n🔹 Status Transitions")
	fmt.Println("─────────────────────────────────────────")

	status := StatusPending
	fmt.Printf("1. Status: %d (Pending)\n", status)

	status = StatusActive
	fmt.Printf("2. Status: %d (Active)\n", status)

	status = StatusCompleted
	fmt.Printf("3. Status: %d (Completed)\n", status)

	// ===== Color (з пропуском) =====
	fmt.Println("\n🔹 Color (з пропуском значення)")
	fmt.Println("─────────────────────────────────────────")

	fmt.Printf("Red: %d\n", Red)
	fmt.Printf("Green: %d\n", Green)
	fmt.Printf("Blue: %d\n", Blue)
	fmt.Printf("Yellow: %d (пропустили 3)\n", Yellow)
	fmt.Printf("Purple: %d\n", Purple)

	// ===== Practical Example: File Modes =====
	fmt.Println("\n🔹 Практичний приклад: File Modes")
	fmt.Println("─────────────────────────────────────────")

	type FileMode Permission

	// rwx для різних ролей
	ownerPerms := Read | Write | Execute
	groupPerms := Read | Execute
	otherPerms := Read

	fmt.Printf("Owner: %s (%03b)\n", ownerPerms, ownerPerms)
	fmt.Printf("Group: %s (%03b)\n", groupPerms, groupPerms)
	fmt.Printf("Other: %s (%03b)\n", otherPerms, otherPerms)

	// ===== Enum Methods =====
	fmt.Println("\n🔹 Enum з методами")
	fmt.Println("─────────────────────────────────────────")

	day := Friday
	fmt.Printf("Day: %s\n", day)
	fmt.Printf("Is Weekend? %t\n", isWeekend(day))
	fmt.Printf("Is Workday? %t\n", isWorkday(day))

	// ===== Reset in new const block =====
	fmt.Println("\n🔹 iota Reset")
	fmt.Println("─────────────────────────────────────────")

	const (
		A = iota // 0
		B        // 1
		C        // 2
	)

	const (
		X = iota // 0 (reset!)
		Y        // 1
		Z        // 2
	)

	fmt.Printf("A=%d, B=%d, C=%d\n", A, B, C)
	fmt.Printf("X=%d, Y=%d, Z=%d\n", X, Y, Z)

	// ===== Complex iota expressions =====
	fmt.Println("\n🔹 Складні вирази з iota")
	fmt.Println("─────────────────────────────────────────")

	const (
		Val1 = iota * 10  // 0
		Val2              // 10
		Val3              // 20
		Val4 = iota * 100 // 300
		Val5              // 400
	)

	fmt.Printf("Val1=%d, Val2=%d, Val3=%d\n", Val1, Val2, Val3)
	fmt.Printf("Val4=%d, Val5=%d\n", Val4, Val5)

	// ===== Summary =====
	fmt.Println("\n\n📝 ВИСНОВКИ")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("✅ iota - авто-інкремент константа")
	fmt.Println("✅ Починається з 0 в кожному const блоці")
	fmt.Println("✅ Ідеально для enum")
	fmt.Println("✅ Bit flags: 1 << iota")
	fmt.Println("✅ Можна пропускати значення через _")
	fmt.Println("✅ Реалізуйте String() для красивого виводу")
	fmt.Println("✅ Використовуйте для status codes, priorities, тощо")
}

// ============= Helper Functions =============

func isWeekend(day Weekday) bool {
	return day == Saturday || day == Sunday
}

func isWorkday(day Weekday) bool {
	return !isWeekend(day)
}
