package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// ============= io.Closer Interface =============

// type Closer interface {
//     Close() error
// }

// ============= Custom Closer =============

type LogFile struct {
	filename string
	file     *os.File
	writer   *bufio.Writer
}

func NewLogFile(filename string) (*LogFile, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	return &LogFile{
		filename: filename,
		file:     file,
		writer:   bufio.NewWriter(file),
	}, nil
}

func (l *LogFile) Write(message string) error {
	_, err := l.writer.WriteString(message + "\n")
	return err
}

func (l *LogFile) Close() error {
	// Flush буфер перед закриттям
	if err := l.writer.Flush(); err != nil {
		return err
	}
	return l.file.Close()
}

// ============= Resource Manager =============

type ResourceManager struct {
	name   string
	isOpen bool
	data   []string
}

func NewResourceManager(name string) *ResourceManager {
	fmt.Printf("📂 Opening resource: %s\n", name)
	return &ResourceManager{
		name:   name,
		isOpen: true,
		data:   []string{},
	}
}

func (r *ResourceManager) Add(item string) error {
	if !r.isOpen {
		return fmt.Errorf("resource %s is closed", r.name)
	}
	r.data = append(r.data, item)
	return nil
}

func (r *ResourceManager) Close() error {
	if !r.isOpen {
		return fmt.Errorf("resource %s already closed", r.name)
	}

	fmt.Printf("🔒 Closing resource: %s (items: %d)\n", r.name, len(r.data))
	r.isOpen = false
	r.data = nil
	return nil
}

// ============= Database Connection (mock) =============

type DBConnection struct {
	host      string
	connected bool
}

func Connect(host string) (*DBConnection, error) {
	fmt.Printf("📡 Connecting to database: %s\n", host)
	return &DBConnection{
		host:      host,
		connected: true,
	}, nil
}

func (db *DBConnection) Query(sql string) ([]string, error) {
	if !db.connected {
		return nil, fmt.Errorf("not connected")
	}

	// Mock результат
	return []string{"row1", "row2", "row3"}, nil
}

func (db *DBConnection) Close() error {
	if !db.connected {
		return fmt.Errorf("already closed")
	}

	fmt.Printf("🔌 Closing database connection: %s\n", db.host)
	db.connected = false
	return nil
}

// ============= MultiCloser =============

type MultiCloser struct {
	closers []io.Closer
}

func NewMultiCloser(closers ...io.Closer) *MultiCloser {
	return &MultiCloser{closers: closers}
}

func (m *MultiCloser) Close() error {
	var firstErr error

	for _, closer := range m.closers {
		if err := closer.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}

// ============= SafeCloser (wrapper) =============

type SafeCloser struct {
	closer io.Closer
	closed bool
}

func NewSafeCloser(closer io.Closer) *SafeCloser {
	return &SafeCloser{closer: closer}
}

func (s *SafeCloser) Close() error {
	if s.closed {
		return nil // вже закрито - OK
	}

	s.closed = true
	return s.closer.Close()
}

// ============= Helper Functions =====

// CloseQuietly закриває ресурс і ігнорує помилки
func CloseQuietly(closer io.Closer) {
	if closer != nil {
		_ = closer.Close()
	}
}

// CloseWithLog закриває ресурс і логує помилки
func CloseWithLog(closer io.Closer, name string) {
	if closer != nil {
		if err := closer.Close(); err != nil {
			fmt.Printf("⚠️  Error closing %s: %v\n", name, err)
		}
	}
}

// ============= Main =============

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║         io.Closer Interface              ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	// ===== Basic Example =====
	fmt.Println("\n🔹 Базовий приклад (файл)")
	fmt.Println("─────────────────────────────────────────")

	// Створення файлу
	file, err := os.Create("test.txt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Запис
	file.WriteString("Hello, World!\n")

	// Закриття
	err = file.Close()
	if err != nil {
		fmt.Printf("Error closing file: %v\n", err)
	} else {
		fmt.Println("✅ File closed successfully")
	}

	// ===== defer Pattern =====
	fmt.Println("\n🔹 defer Pattern (рекомендовано)")
	fmt.Println("─────────────────────────────────────────")

	func() {
		file, err := os.Create("defer_test.txt")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		defer file.Close() // автоматично закриється

		file.WriteString("Using defer\n")
		fmt.Println("✅ File will be closed automatically")
	}()

	// ===== Custom LogFile =====
	fmt.Println("\n🔹 Custom LogFile")
	fmt.Println("─────────────────────────────────────────")

	log, err := NewLogFile("app.log")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer log.Close()

	log.Write("Application started")
	log.Write("Processing data...")
	log.Write("Done")

	fmt.Println("✅ Log file created and will be closed")

	// ===== ResourceManager =====
	fmt.Println("\n🔹 ResourceManager")
	fmt.Println("─────────────────────────────────────────")

	resource := NewResourceManager("Cache")
	defer resource.Close()

	resource.Add("item1")
	resource.Add("item2")
	resource.Add("item3")

	fmt.Println("✅ Resource will be cleaned up")

	// ===== Database Connection =====
	fmt.Println("\n🔹 Database Connection")
	fmt.Println("─────────────────────────────────────────")

	db, err := Connect("localhost:5432")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer db.Close()

	results, _ := db.Query("SELECT * FROM users")
	fmt.Printf("Query results: %v\n", results)
	fmt.Println("✅ Connection will be closed")

	// ===== MultiCloser =====
	fmt.Println("\n🔹 MultiCloser (кілька ресурсів)")
	fmt.Println("─────────────────────────────────────────")

	r1 := NewResourceManager("Resource1")
	r2 := NewResourceManager("Resource2")
	r3 := NewResourceManager("Resource3")

	multi := NewMultiCloser(r1, r2, r3)
	defer multi.Close()

	fmt.Println("✅ All resources will be closed together")

	// ===== SafeCloser =====
	fmt.Println("\n🔹 SafeCloser (запобігання подвійного закриття)")
	fmt.Println("─────────────────────────────────────────")

	resource2 := NewResourceManager("SafeResource")
	safe := NewSafeCloser(resource2)

	safe.Close() // перше закриття
	fmt.Println("First close: OK")

	safe.Close() // друге закриття - безпечно
	fmt.Println("Second close: OK (ignored)")

	// ===== Error Handling =====
	fmt.Println("\n🔹 Error Handling")
	fmt.Println("─────────────────────────────────────────")

	// Приклад з помилкою
	func() {
		file, err := os.Create("/invalid/path/file.txt")
		if err != nil {
			fmt.Printf("❌ Cannot create file: %v\n", err)
			return
		}
		defer file.Close()

		// Цей код не виконається
		file.WriteString("This won't work")
	}()

	// ===== Named Return with defer =====
	fmt.Println("\n🔹 Named Return + defer")
	fmt.Println("─────────────────────────────────────────")

	readFileContent := func(filename string) (content string, err error) {
		file, err := os.Open(filename)
		if err != nil {
			return "", err
		}
		defer func() {
			closeErr := file.Close()
			if err == nil {
				err = closeErr // якщо немає іншої помилки, повернути помилку закриття
			}
		}()

		// Читання
		var buf strings.Builder
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			buf.WriteString(scanner.Text())
			buf.WriteString("\n")
		}

		return buf.String(), scanner.Err()
	}

	content, err := readFileContent("test.txt")
	if err == nil {
		fmt.Printf("✅ Read content: %s", content)
	}

	// ===== io.ReadCloser Example =====
	fmt.Println("\n🔹 io.ReadCloser (композиція)")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("type ReadCloser interface {")
	fmt.Println("    Reader")
	fmt.Println("    Closer")
	fmt.Println("}")
	fmt.Println()
	fmt.Println("Приклади:")
	fmt.Println("  • os.File")
	fmt.Println("  • http.Response.Body")
	fmt.Println("  • gzip.Reader")

	// ===== Multiple defer Order =====
	fmt.Println("\n🔹 Порядок виконання defer")
	fmt.Println("─────────────────────────────────────────")

	func() {
		r1 := NewResourceManager("First")
		defer r1.Close()

		r2 := NewResourceManager("Second")
		defer r2.Close()

		r3 := NewResourceManager("Third")
		defer r3.Close()

		fmt.Println("All resources opened")
	}()

	fmt.Println("(закриваються у зворотному порядку: LIFO)")

	// ===== Best Practices =====
	fmt.Println("\n🔹 Best Practices")
	fmt.Println("─────────────────────────────────────────")

	fmt.Println(`
✅ Добре:
func ProcessFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()  // ← одразу після відкриття
    
    // робота з файлом
    return nil
}

❌ Погано:
func ProcessFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    
    // багато коду...
    
    file.Close()  // ← можна забути або пропустити
    return nil
}
	`)

	// Cleanup test files
	os.Remove("test.txt")
	os.Remove("defer_test.txt")
	os.Remove("app.log")

	// ===== Summary =====
	fmt.Println("\n📝 ВИСНОВКИ")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("✅ io.Closer - для очищення ресурсів")
	fmt.Println("   • Файли")
	fmt.Println("   • Мережеві з'єднання")
	fmt.Println("   • Database connections")
	fmt.Println("   • Будь-які ресурси що потребують cleanup")
	fmt.Println()
	fmt.Println("💡 defer pattern:")
	fmt.Println("   • Завжди використовувати defer")
	fmt.Println("   • defer одразу після створення ресурсу")
	fmt.Println("   • Порядок: LIFO (останній відкритий - перший закритий)")
	fmt.Println()
	fmt.Println("🔗 Композиція:")
	fmt.Println("   • io.ReadCloser = Reader + Closer")
	fmt.Println("   • io.WriteCloser = Writer + Closer")
	fmt.Println("   • io.ReadWriteCloser = Reader + Writer + Closer")
	fmt.Println()
	fmt.Println("⚠️  Важливо:")
	fmt.Println("   • Завжди перевіряти помилки Close()")
	fmt.Println("   • Не закривати двічі (або використати SafeCloser)")
	fmt.Println("   • Використовувати defer для гарантії закриття")
}
