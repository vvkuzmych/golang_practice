# Environment Variables & Process Management

## Environment Variables

### os.Getenv
**Отримати значення змінної оточення.**

```go
// Get environment variable
home := os.Getenv("HOME")
fmt.Println("Home directory:", home)

// Якщо змінна не існує - повертає пустий string
path := os.Getenv("PATH")
if path == "" {
    fmt.Println("PATH not set")
}

// Check if variable is set
port := os.Getenv("PORT")
if port == "" {
    port = "8080"  // default value
}
```

### os.LookupEnv
**Перевірити чи змінна встановлена.**

```go
value, exists := os.LookupEnv("MY_VAR")
if !exists {
    fmt.Println("MY_VAR is not set")
} else {
    fmt.Println("MY_VAR =", value)
}
```

### os.Setenv
**Встановити змінну оточення.**

```go
err := os.Setenv("MY_VAR", "my_value")
if err != nil {
    panic(err)
}

// Перевірити
fmt.Println(os.Getenv("MY_VAR"))  // my_value
```

### os.Unsetenv
**Видалити змінну оточення.**

```go
err := os.Unsetenv("MY_VAR")
if err != nil {
    panic(err)
}
```

### os.Environ
**Отримати всі змінні оточення.**

```go
for _, env := range os.Environ() {
    fmt.Println(env)
}

// Output:
// HOME=/home/user
// PATH=/usr/local/bin:/usr/bin
// LANG=en_US.UTF-8
// ...
```

### os.ExpandEnv
**Розгорнути змінні в string.**

```go
path := os.ExpandEnv("$HOME/projects/$USER")
// path: "/home/alice/projects/alice"

// Або з ${} синтаксисом
path = os.ExpandEnv("${HOME}/projects")
```

---

## Практичні приклади з ENV

### 1. Конфігурація з ENV variables

```go
type Config struct {
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    AppPort    string
    LogLevel   string
}

func LoadConfig() *Config {
    return &Config{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: os.Getenv("DB_PASSWORD"),  // no default for security
        AppPort:    getEnv("APP_PORT", "8080"),
        LogLevel:   getEnv("LOG_LEVEL", "info"),
    }
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}

func main() {
    config := LoadConfig()
    fmt.Printf("Connecting to %s:%s\n", config.DBHost, config.DBPort)
}
```

### 2. .env file loader

```go
import (
    "bufio"
    "os"
    "strings"
)

func LoadEnvFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        
        // Skip comments and empty lines
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }
        
        // Parse KEY=VALUE
        parts := strings.SplitN(line, "=", 2)
        if len(parts) != 2 {
            continue
        }
        
        key := strings.TrimSpace(parts[0])
        value := strings.TrimSpace(parts[1])
        
        // Remove quotes
        value = strings.Trim(value, "\"'")
        
        os.Setenv(key, value)
    }
    
    return scanner.Err()
}

// .env file:
// DB_HOST=localhost
// DB_PORT=5432
// # This is a comment
// APP_NAME="My App"

// Usage:
LoadEnvFile(".env")
fmt.Println(os.Getenv("DB_HOST"))  // localhost
```

---

## Command-Line Arguments

### os.Args
**Отримати аргументи командного рядка.**

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // os.Args[0] - ім'я програми
    // os.Args[1:] - аргументи
    
    fmt.Println("Program:", os.Args[0])
    fmt.Println("Arguments:", os.Args[1:])
    
    if len(os.Args) < 2 {
        fmt.Println("Usage: program <arg1> <arg2>")
        os.Exit(1)
    }
    
    arg1 := os.Args[1]
    fmt.Println("First argument:", arg1)
}

// Run: ./program hello world
// Output:
// Program: ./program
// Arguments: [hello world]
// First argument: hello
```

### flag package (рекомендується)

```go
import "flag"

func main() {
    // Define flags
    name := flag.String("name", "World", "Name to greet")
    count := flag.Int("count", 1, "Number of times to greet")
    verbose := flag.Bool("verbose", false, "Verbose output")
    
    // Parse
    flag.Parse()
    
    // Use
    for i := 0; i < *count; i++ {
        fmt.Printf("Hello, %s!\n", *name)
    }
    
    if *verbose {
        fmt.Println("Done!")
    }
}

// Run: ./program -name=Alice -count=3 -verbose
// Output:
// Hello, Alice!
// Hello, Alice!
// Hello, Alice!
// Done!
```

---

## Process Information

### os.Getpid
**Отримати PID поточного процесу.**

```go
pid := os.Getpid()
fmt.Println("Process ID:", pid)
```

### os.Getppid
**Отримати PID батьківського процесу.**

```go
ppid := os.Getppid()
fmt.Println("Parent Process ID:", ppid)
```

### os.Getuid, os.Getgid
**Отримати user/group ID (Unix).**

```go
uid := os.Getuid()
gid := os.Getgid()
fmt.Printf("User ID: %d, Group ID: %d\n", uid, gid)
```

### os.Hostname
**Отримати ім'я хоста.**

```go
hostname, err := os.Hostname()
if err != nil {
    panic(err)
}
fmt.Println("Hostname:", hostname)
```

### os.UserCacheDir, os.UserConfigDir, os.UserHomeDir
**Отримати стандартні директорії користувача.**

```go
// Cache directory
cache, _ := os.UserCacheDir()
// macOS: /Users/alice/Library/Caches
// Linux: /home/alice/.cache
// Windows: C:\Users\alice\AppData\Local

// Config directory
config, _ := os.UserConfigDir()
// macOS: /Users/alice/Library/Application Support
// Linux: /home/alice/.config

// Home directory
home, _ := os.UserHomeDir()
// /Users/alice (macOS)
// /home/alice (Linux)
// C:\Users\alice (Windows)
```

---

## Process Exit

### os.Exit
**Вийти з програми з кодом.**

```go
func main() {
    if len(os.Args) < 2 {
        fmt.Println("Error: missing argument")
        os.Exit(1)  // exit with error code
    }
    
    // Success
    os.Exit(0)
}

// Exit codes:
// 0 - success
// 1 - general error
// 2 - misuse of command
// 126 - command cannot execute
// 127 - command not found
// 130 - terminated by Ctrl+C
```

### defer не виконується з os.Exit!

```go
func main() {
    defer fmt.Println("This won't be printed!")
    os.Exit(1)  // defer не виконається
}

// Замість цього:
func main() {
    code := run()
    os.Exit(code)
}

func run() int {
    defer fmt.Println("Cleanup done")  // виконається
    
    if error {
        return 1
    }
    return 0
}
```

---

## Практичний приклад: CLI Tool

```go
package main

import (
    "flag"
    "fmt"
    "os"
)

func main() {
    // Flags
    verbose := flag.Bool("v", false, "Verbose output")
    output := flag.String("o", "output.txt", "Output file")
    
    flag.Parse()
    
    // Remaining arguments
    args := flag.Args()
    if len(args) == 0 {
        fmt.Fprintf(os.Stderr, "Usage: %s [options] <input>\n", os.Args[0])
        flag.PrintDefaults()
        os.Exit(1)
    }
    
    input := args[0]
    
    // Verbose mode
    if *verbose {
        fmt.Printf("Processing: %s\n", input)
        fmt.Printf("Output: %s\n", *output)
        fmt.Printf("PID: %d\n", os.Getpid())
    }
    
    // Check environment
    logLevel := os.Getenv("LOG_LEVEL")
    if logLevel == "" {
        logLevel = "info"
    }
    
    if *verbose {
        fmt.Printf("Log level: %s\n", logLevel)
    }
    
    // Do work...
    fmt.Println("Done!")
}

// Run:
// LOG_LEVEL=debug ./program -v -o result.txt input.txt
```

---

## Signal Handling

```go
import (
    "os"
    "os/signal"
    "syscall"
)

func main() {
    // Channel для signals
    sigChan := make(chan os.Signal, 1)
    
    // Підписатися на SIGINT (Ctrl+C) та SIGTERM
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    
    // Чекати на signal
    fmt.Println("Waiting for signal...")
    sig := <-sigChan
    
    fmt.Printf("\nReceived signal: %v\n", sig)
    fmt.Println("Cleaning up...")
    
    // Cleanup
    // ...
    
    os.Exit(0)
}

// Press Ctrl+C:
// Received signal: interrupt
// Cleaning up...
```

### Graceful shutdown

```go
func main() {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    
    // Goroutine для обробки signal
    go func() {
        sig := <-sigChan
        fmt.Printf("\nReceived %v, shutting down...\n", sig)
        
        // Cleanup
        cleanupResources()
        
        os.Exit(0)
    }()
    
    // Main loop
    for {
        fmt.Println("Working...")
        time.Sleep(1 * time.Second)
    }
}

func cleanupResources() {
    fmt.Println("Closing database connections...")
    fmt.Println("Flushing logs...")
    fmt.Println("Goodbye!")
}
```

---

## Підсумок

| Операція | Функція | Опис |
|----------|---------|------|
| Get ENV | `os.Getenv` | Отримати змінну |
| Set ENV | `os.Setenv` | Встановити змінну |
| Check ENV | `os.LookupEnv` | Перевірити існування |
| All ENV | `os.Environ` | Всі змінні |
| CLI args | `os.Args` | Аргументи командного рядка |
| PID | `os.Getpid` | Process ID |
| Exit | `os.Exit` | Вийти з кодом |
| Hostname | `os.Hostname` | Ім'я хоста |
| Home dir | `os.UserHomeDir` | Домашня директорія |
| Signals | `signal.Notify` | Обробка сигналів |

**Best Practices:**
- Використовуй ENV variables для конфігурації
- Використовуй `flag` package для CLI args
- Завжди перевіряй наявність обов'язкових ENV variables
- Використовуй `os.Exit(0)` для success, `os.Exit(1)` для помилок
- Обробляй SIGINT/SIGTERM для graceful shutdown
