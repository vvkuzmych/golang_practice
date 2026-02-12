# Week 21 — Quick Start

## 🎯 Мета тижня
Повністю освоїти роботу з операційною системою через `os` package в Go.

---

## 📖 Швидке навчання (60 хв)

```bash
# 1. File Operations
cat theory/01_file_operations.md

# 2. Directory Operations  
cat theory/02_directories.md

# 3. Environment & Process
cat theory/03_environment_process.md

# 4. File Info & Permissions
cat theory/04_file_info_permissions.md
```

---

## 💡 Ключові функції

### File Operations
```go
// Create/Open/Write
file, _ := os.Create("file.txt")
file.WriteString("Hello")
file.Close()

// Read
data, _ := os.ReadFile("file.txt")
fmt.Println(string(data))

// Copy
io.Copy(dst, src)

// Delete
os.Remove("file.txt")
```

### Directory Operations
```go
// Create
os.MkdirAll("path/to/dir", 0755)

// Read
entries, _ := os.ReadDir(".")
for _, e := range entries {
    fmt.Println(e.Name())
}

// Walk
filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
    fmt.Println(path)
    return nil
})
```

### Environment Variables
```go
// Get
port := os.Getenv("PORT")
if port == "" {
    port = "8080"  // default
}

// Set
os.Setenv("MY_VAR", "value")

// All
for _, env := range os.Environ() {
    fmt.Println(env)
}
```

### File Info
```go
// Stat
info, _ := os.Stat("file.txt")
fmt.Println("Size:", info.Size())
fmt.Println("ModTime:", info.ModTime())
fmt.Println("Permissions:", info.Mode())

// Chmod
os.Chmod("file.txt", 0644)

// Symlink
os.Symlink("target.txt", "link.txt")
```

---

## 🚀 Практичні патерни

### 1. Safe File Write
```go
func safeWrite(filename string, data []byte) error {
    tmp, _ := os.CreateTemp("", "*.tmp")
    defer os.Remove(tmp.Name())
    tmp.Write(data)
    tmp.Close()
    return os.Rename(tmp.Name(), filename)
}
```

### 2. File Exists Check
```go
func fileExists(path string) bool {
    _, err := os.Stat(path)
    return err == nil
}
```

### 3. Find Files
```go
func findFiles(root, ext string) []string {
    var files []string
    filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
        if !d.IsDir() && filepath.Ext(path) == ext {
            files = append(files, path)
        }
        return nil
    })
    return files
}
```

---

## ⚠️ Common Mistakes

### ❌ Забути закрити файл
```go
// BAD
file, _ := os.Open("file.txt")
// forgot to close - memory leak

// GOOD
file, _ := os.Open("file.txt")
defer file.Close()
```

### ❌ Ігнорувати помилки
```go
// BAD
data, _ := os.ReadFile("file.txt")

// GOOD
data, err := os.ReadFile("file.txt")
if err != nil {
    return err
}
```

### ❌ Hardcoded paths
```go
// BAD (only works on Unix)
path := "dir/subdir/file.txt"

// GOOD (cross-platform)
path := filepath.Join("dir", "subdir", "file.txt")
```

---

## 📝 Mini Project Ideas

1. **File Manager CLI** - ls, cat, cp, mv, rm commands
2. **Log Rotator** - rotate logs by size/date
3. **Duplicate Finder** - find duplicate files by hash
4. **Config Manager** - read/write configs with ENV vars
5. **Backup Tool** - incremental backup with timestamps

---

## ✅ Перевірка розуміння

- [ ] Можу створювати, читати, записувати файли
- [ ] Розумію різницю між Open, Create, OpenFile
- [ ] Можу працювати з директоріями (Mkdir, ReadDir, Walk)
- [ ] Знаю як використовувати ENV variables
- [ ] Розумію file permissions (0644, 0755)
- [ ] Можу отримувати file info (Stat)
- [ ] Знаю як створювати temp files
- [ ] Розумію symlinks

---

## 🔗 Корисні команди

```bash
# Run Go file
go run main.go

# Build executable
go build -o myapp main.go

# Check file permissions
ls -la file.txt

# Create test file
echo "test" > test.txt

# Check environment variables
printenv

# Run with env var
PORT=8080 go run main.go
```

---

## 🚀 Наступний крок

Практикуй написання CLI tools:
- File utilities
- System monitors
- Backup scripts
- Config managers
