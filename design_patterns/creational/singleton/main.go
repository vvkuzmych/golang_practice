package main

import (
	"fmt"
	"sync"
)

// ============= Singleton: Database Connection =============

// Database представляє підключення до бази даних
type Database struct {
	connection string
	queries    int
}

var (
	dbInstance *Database
	once       sync.Once
)

// GetDatabase повертає singleton екземпляр Database
func GetDatabase() *Database {
	once.Do(func() {
		fmt.Println("🔧 Creating database instance...")
		dbInstance = &Database{
			connection: "postgresql://localhost:5432/mydb",
			queries:    0,
		}
	})
	return dbInstance
}

// Query виконує запит до БД
func (db *Database) Query(sql string) {
	db.queries++
	fmt.Printf("📊 Executing query #%d: %s\n", db.queries, sql)
}

// GetStats повертає статистику
func (db *Database) GetStats() string {
	return fmt.Sprintf("Connection: %s, Queries: %d", db.connection, db.queries)
}

// ============= Singleton: Logger =============

type Logger struct {
	prefix string
}

var (
	loggerInstance *Logger
	loggerOnce     sync.Once
)

// GetLogger повертає singleton logger
func GetLogger() *Logger {
	loggerOnce.Do(func() {
		fmt.Println("📝 Creating logger instance...")
		loggerInstance = &Logger{
			prefix: "[APP]",
		}
	})
	return loggerInstance
}

// Info логує info повідомлення
func (l *Logger) Info(message string) {
	fmt.Printf("%s INFO: %s\n", l.prefix, message)
}

// Error логує error повідомлення
func (l *Logger) Error(message string) {
	fmt.Printf("%s ERROR: %s\n", l.prefix, message)
}

// ============= Singleton: Configuration =============

type Config struct {
	AppName string
	Port    int
	Debug   bool
}

var (
	configInstance *Config
	configOnce     sync.Once
)

// GetConfig повертає singleton конфігурацію
func GetConfig() *Config {
	configOnce.Do(func() {
		fmt.Println("⚙️  Loading configuration...")
		configInstance = &Config{
			AppName: "MyApp",
			Port:    8080,
			Debug:   true,
		}
	})
	return configInstance
}

// String повертає string представлення конфігурації
func (c *Config) String() string {
	return fmt.Sprintf("App: %s, Port: %d, Debug: %v", c.AppName, c.Port, c.Debug)
}

// ============= Demo Functions =============

func userService() {
	fmt.Println("\n👤 UserService working...")

	db := GetDatabase()
	db.Query("SELECT * FROM users")

	logger := GetLogger()
	logger.Info("UserService initialized")
}

func orderService() {
	fmt.Println("\n📦 OrderService working...")

	db := GetDatabase()
	db.Query("SELECT * FROM orders")

	logger := GetLogger()
	logger.Info("OrderService initialized")
}

func paymentService() {
	fmt.Println("\n💳 PaymentService working...")

	db := GetDatabase()
	db.Query("SELECT * FROM payments")

	logger := GetLogger()
	logger.Info("PaymentService initialized")
}

// ============= Thread-Safety Demo =============

func concurrentAccess() {
	fmt.Println("\n🔄 Testing concurrent access...")

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			db := GetDatabase()
			db.Query(fmt.Sprintf("Query from goroutine %d", id))
		}(i)
	}

	wg.Wait()
}

// ============= Main =============

func main() {
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║       Singleton Pattern Demo          ║")
	fmt.Println("╚════════════════════════════════════════╝")

	// ===== Demo 1: Basic Singleton =====
	fmt.Println("\n🔹 Demo 1: Basic Singleton Usage")
	fmt.Println("─────────────────────────────────────")

	fmt.Println("Getting database instance (1st time):")
	db1 := GetDatabase()
	fmt.Printf("Instance: %p\n", db1)

	fmt.Println("\nGetting database instance (2nd time):")
	db2 := GetDatabase()
	fmt.Printf("Instance: %p\n", db2)

	if db1 == db2 {
		fmt.Println("✅ Same instance! Singleton works!")
	}

	// ===== Demo 2: Multiple Services =====
	fmt.Println("\n🔹 Demo 2: Multiple Services Using Singletons")
	fmt.Println("─────────────────────────────────────")

	userService()
	orderService()
	paymentService()

	// Перевіряємо статистику
	fmt.Println("\n📊 Database Statistics:")
	fmt.Println(db1.GetStats())

	// ===== Demo 3: Logger Singleton =====
	fmt.Println("\n🔹 Demo 3: Logger Singleton")
	fmt.Println("─────────────────────────────────────")

	logger1 := GetLogger()
	logger1.Info("First log message")

	logger2 := GetLogger()
	logger2.Error("Error occurred")

	if logger1 == logger2 {
		fmt.Println("✅ Logger is singleton too!")
	}

	// ===== Demo 4: Config Singleton =====
	fmt.Println("\n🔹 Demo 4: Configuration Singleton")
	fmt.Println("─────────────────────────────────────")

	config1 := GetConfig()
	fmt.Printf("Config1: %s\n", config1)

	config2 := GetConfig()
	fmt.Printf("Config2: %s\n", config2)

	if config1 == config2 {
		fmt.Println("✅ Config is singleton!")
	}

	// ===== Demo 5: Thread Safety =====
	fmt.Println("\n🔹 Demo 5: Thread-Safe Singleton")
	fmt.Println("─────────────────────────────────────")

	concurrentAccess()

	fmt.Println("\n📊 Final Database Statistics:")
	fmt.Println(db1.GetStats())

	// ===== Summary =====
	fmt.Println("\n\n📝 ВИСНОВКИ")
	fmt.Println("─────────────────────────────────────")
	fmt.Println("✅ Один екземпляр на всю програму")
	fmt.Println("✅ Глобальний доступ")
	fmt.Println("✅ Lazy initialization (створюється при потребі)")
	fmt.Println("✅ Thread-safe (завдяки sync.Once)")
	fmt.Println("✅ Підходить для: DB, Logger, Config, Cache")

	fmt.Println("\n⚠️  ОБЕРЕЖНО:")
	fmt.Println("❌ Не зловживайте! Ускладнює тестування")
	fmt.Println("❌ Краще використовувати DI де можливо")
	fmt.Println("❌ Глобальний стан = складніший код")
}
