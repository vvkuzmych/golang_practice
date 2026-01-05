package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ============= Storage Interface =============

type Storage interface {
	Save(key, value string) error
	Load(key string) (string, error)
	Delete(key string) error
	Exists(key string) bool
	Keys() []string
	Clear() error
}

// ============= Memory Storage =============

type MemoryStorage struct {
	data map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string]string),
	}
}

func (m *MemoryStorage) Save(key, value string) error {
	m.data[key] = value
	return nil
}

func (m *MemoryStorage) Load(key string) (string, error) {
	value, exists := m.data[key]
	if !exists {
		return "", fmt.Errorf("key not found: %s", key)
	}
	return value, nil
}

func (m *MemoryStorage) Delete(key string) error {
	if !m.Exists(key) {
		return fmt.Errorf("key not found: %s", key)
	}
	delete(m.data, key)
	return nil
}

func (m *MemoryStorage) Exists(key string) bool {
	_, exists := m.data[key]
	return exists
}

func (m *MemoryStorage) Keys() []string {
	keys := make([]string, 0, len(m.data))
	for key := range m.data {
		keys = append(keys, key)
	}
	return keys
}

func (m *MemoryStorage) Clear() error {
	m.data = make(map[string]string)
	return nil
}

func (m *MemoryStorage) Size() int {
	return len(m.data)
}

// ============= File Storage =============

type FileStorage struct {
	filename string
	data     map[string]string
}

func NewFileStorage(filename string) *FileStorage {
	fs := &FileStorage{
		filename: filename,
		data:     make(map[string]string),
	}
	fs.loadFromFile()
	return fs
}

func (f *FileStorage) Save(key, value string) error {
	f.data[key] = value
	return f.saveToFile()
}

func (f *FileStorage) Load(key string) (string, error) {
	value, exists := f.data[key]
	if !exists {
		return "", fmt.Errorf("key not found: %s", key)
	}
	return value, nil
}

func (f *FileStorage) Delete(key string) error {
	if !f.Exists(key) {
		return fmt.Errorf("key not found: %s", key)
	}
	delete(f.data, key)
	return f.saveToFile()
}

func (f *FileStorage) Exists(key string) bool {
	_, exists := f.data[key]
	return exists
}

func (f *FileStorage) Keys() []string {
	keys := make([]string, 0, len(f.data))
	for key := range f.data {
		keys = append(keys, key)
	}
	return keys
}

func (f *FileStorage) Clear() error {
	f.data = make(map[string]string)
	return f.saveToFile()
}

func (f *FileStorage) saveToFile() error {
	file, err := os.Create(f.filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	for key, value := range f.data {
		_, err := fmt.Fprintf(file, "%s=%s\n", key, value)
		if err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	return nil
}

func (f *FileStorage) loadFromFile() error {
	file, err := os.Open(f.filename)
	if err != nil {
		// Файл не існує - це OK для нового storage
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			f.data[parts[0]] = parts[1]
		}
	}

	return scanner.Err()
}

// ============= Mock Storage =============

type MockStorage struct {
	data         map[string]string
	saveCalled   int
	loadCalled   int
	deleteCalled int
	shouldFail   bool
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		data: make(map[string]string),
	}
}

func (m *MockStorage) SetShouldFail(fail bool) {
	m.shouldFail = fail
}

func (m *MockStorage) Save(key, value string) error {
	m.saveCalled++
	if m.shouldFail {
		return fmt.Errorf("mock: save failed")
	}
	m.data[key] = value
	return nil
}

func (m *MockStorage) Load(key string) (string, error) {
	m.loadCalled++
	if m.shouldFail {
		return "", fmt.Errorf("mock: load failed")
	}
	value, exists := m.data[key]
	if !exists {
		return "", fmt.Errorf("key not found: %s", key)
	}
	return value, nil
}

func (m *MockStorage) Delete(key string) error {
	m.deleteCalled++
	if m.shouldFail {
		return fmt.Errorf("mock: delete failed")
	}
	delete(m.data, key)
	return nil
}

func (m *MockStorage) Exists(key string) bool {
	_, exists := m.data[key]
	return exists
}

func (m *MockStorage) Keys() []string {
	keys := make([]string, 0, len(m.data))
	for key := range m.data {
		keys = append(keys, key)
	}
	return keys
}

func (m *MockStorage) Clear() error {
	if m.shouldFail {
		return fmt.Errorf("mock: clear failed")
	}
	m.data = make(map[string]string)
	return nil
}

func (m *MockStorage) Stats() string {
	return fmt.Sprintf("Save: %d, Load: %d, Delete: %d",
		m.saveCalled, m.loadCalled, m.deleteCalled)
}

// ============= Data Manager =============

type DataManager struct {
	storage Storage
}

func NewDataManager(storage Storage) *DataManager {
	return &DataManager{storage: storage}
}

func (d *DataManager) Set(key, value string) error {
	return d.storage.Save(key, value)
}

func (d *DataManager) Get(key string) (string, error) {
	return d.storage.Load(key)
}

func (d *DataManager) Remove(key string) error {
	return d.storage.Delete(key)
}

func (d *DataManager) Has(key string) bool {
	return d.storage.Exists(key)
}

func (d *DataManager) AllKeys() []string {
	return d.storage.Keys()
}

func (d *DataManager) Reset() error {
	return d.storage.Clear()
}

func (d *DataManager) PrintAll() {
	keys := d.AllKeys()
	if len(keys) == 0 {
		fmt.Println("  (empty)")
		return
	}

	for _, key := range keys {
		value, _ := d.Get(key)
		fmt.Printf("  %s = %s\n", key, value)
	}
}

// ============= Main =============

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║        Storage Interface Solution        ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	// ===== Memory Storage =====
	fmt.Println("\n📦 MEMORY STORAGE")
	fmt.Println("─────────────────────────────────────────")

	memStorage := NewMemoryStorage()
	manager1 := NewDataManager(memStorage)

	fmt.Println("\n🔹 Збереження даних:")
	manager1.Set("name", "Іван")
	manager1.Set("age", "25")
	manager1.Set("city", "Київ")
	manager1.Set("country", "Україна")

	fmt.Println("✅ Дані збережені")

	fmt.Println("\n🔹 Читання даних:")
	manager1.PrintAll()

	fmt.Println("\n🔹 Отримання конкретного значення:")
	if name, err := manager1.Get("name"); err == nil {
		fmt.Printf("  name = %s\n", name)
	}

	fmt.Println("\n🔹 Перевірка існування:")
	fmt.Printf("  Exists 'city'? %t\n", manager1.Has("city"))
	fmt.Printf("  Exists 'email'? %t\n", manager1.Has("email"))

	fmt.Println("\n🔹 Всі ключі:")
	fmt.Printf("  %v\n", manager1.AllKeys())

	fmt.Println("\n🔹 Видалення:")
	manager1.Remove("age")
	fmt.Println("  Видалено 'age'")
	manager1.PrintAll()

	// ===== File Storage =====
	fmt.Println("\n\n💾 FILE STORAGE")
	fmt.Println("─────────────────────────────────────────")

	filename := "storage_test.txt"
	fileStorage := NewFileStorage(filename)
	manager2 := NewDataManager(fileStorage)

	fmt.Println("\n🔹 Збереження в файл:")
	manager2.Set("config", "production")
	manager2.Set("version", "1.0.0")
	manager2.Set("debug", "false")
	manager2.Set("port", "8080")

	fmt.Println("✅ Дані збережені в файл:", filename)

	fmt.Println("\n🔹 Дані:")
	manager2.PrintAll()

	// Прочитати файл для демонстрації
	fmt.Println("\n🔹 Вміст файлу:")
	content, err := os.ReadFile(filename)
	if err == nil {
		fmt.Println(string(content))
	}

	fmt.Println("🔹 Перезавантаження з файлу:")
	fileStorage2 := NewFileStorage(filename)
	manager3 := NewDataManager(fileStorage2)
	manager3.PrintAll()

	// ===== Mock Storage =====
	fmt.Println("\n\n🎭 MOCK STORAGE (For Testing)")
	fmt.Println("─────────────────────────────────────────")

	mockStorage := NewMockStorage()
	manager4 := NewDataManager(mockStorage)

	fmt.Println("\n🔹 Операції з Mock:")
	manager4.Set("test1", "value1")
	manager4.Set("test2", "value2")
	manager4.Get("test1")
	manager4.Get("test2")
	manager4.Remove("test1")

	fmt.Println("\n🔹 Mock Статистика:")
	fmt.Printf("  %s\n", mockStorage.Stats())

	fmt.Println("\n🔹 Тестування помилок:")
	mockStorage.SetShouldFail(true)
	err = manager4.Set("test3", "value3")
	if err != nil {
		fmt.Printf("  ❌ Expected error: %v\n", err)
	}

	// ===== Comparison =====
	fmt.Println("\n\n⚖️  ПОРІВНЯННЯ РЕАЛІЗАЦІЙ")
	fmt.Println("─────────────────────────────────────────")

	storages := []struct {
		name    string
		storage Storage
	}{
		{"Memory", memStorage},
		{"File", fileStorage},
		{"Mock", mockStorage},
	}

	for _, s := range storages {
		keys := s.storage.Keys()
		fmt.Printf("%s Storage: %d keys\n", s.name, len(keys))
	}

	// ===== Use Cases =====
	fmt.Println("\n\n💡 ПРИКЛАДИ ВИКОРИСТАННЯ")
	fmt.Println("─────────────────────────────────────────")

	fmt.Println("\n1️⃣  Конфігурація програми:")
	config := NewDataManager(NewMemoryStorage())
	config.Set("app_name", "MyApp")
	config.Set("environment", "production")
	config.Set("max_connections", "100")
	config.PrintAll()

	fmt.Println("\n2️⃣  Кеш даних:")
	cache := NewDataManager(NewMemoryStorage())
	cache.Set("user:1", "John Doe")
	cache.Set("user:2", "Jane Smith")
	fmt.Printf("Кешовано користувачів: %d\n", len(cache.AllKeys()))

	fmt.Println("\n3️⃣  Збереження налаштувань:")
	settings := NewDataManager(NewFileStorage("settings.txt"))
	settings.Set("theme", "dark")
	settings.Set("language", "uk")
	settings.Set("notifications", "enabled")
	fmt.Println("✅ Налаштування збережені у файл")

	// ===== Summary =====
	fmt.Println("\n\n📝 ВИСНОВКИ")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("✅ Реалізовано:")
	fmt.Println("   • Interface Storage")
	fmt.Println("   • MemoryStorage (швидке in-memory зберігання)")
	fmt.Println("   • FileStorage (персистентне зберігання)")
	fmt.Println("   • MockStorage (для тестування)")
	fmt.Println("   • DataManager (працює з будь-яким Storage)")
	fmt.Println()
	fmt.Println("💡 Переваги:")
	fmt.Println("   • Легко змінити реалізацію")
	fmt.Println("   • Легко тестувати через Mock")
	fmt.Println("   • Dependency Injection")
	fmt.Println("   • Один інтерфейс - багато реалізацій")

	// Cleanup
	fmt.Println("\n\n🧹 Очищення тестових файлів...")
	os.Remove(filename)
	os.Remove("settings.txt")
	fmt.Println("✅ Готово!")
}
