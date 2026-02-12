# File Info & Permissions

## File Information

### os.Stat
**Отримати інформацію про файл/директорію.**

```go
info, err := os.Stat("file.txt")
if err != nil {
    if os.IsNotExist(err) {
        fmt.Println("File doesn't exist")
        return
    }
    panic(err)
}

// FileInfo methods
fmt.Println("Name:", info.Name())        // file.txt
fmt.Println("Size:", info.Size())        // bytes
fmt.Println("Mode:", info.Mode())        // permissions
fmt.Println("ModTime:", info.ModTime())  // last modified
fmt.Println("IsDir:", info.IsDir())      // true if directory
```

### os.Lstat
**Як Stat, але для symbolic links повертає інформацію про сам link, а не про target.**

```go
// file.txt -> /path/to/target.txt

info, _ := os.Stat("file.txt")    // Info про target.txt
info, _ := os.Lstat("file.txt")   // Info про сам link
```

---

## FileInfo Interface

```go
type FileInfo interface {
    Name() string       // base name of the file
    Size() int64        // length in bytes
    Mode() FileMode     // file mode bits
    ModTime() time.Time // modification time
    IsDir() bool        // abbreviation for Mode().IsDir()
    Sys() interface{}   // underlying data source (platform-specific)
}
```

### Приклад використання

```go
info, err := os.Stat("file.txt")
if err != nil {
    panic(err)
}

// Name
fmt.Println("Name:", info.Name())

// Size (в різних форматах)
bytes := info.Size()
kb := float64(bytes) / 1024
mb := kb / 1024
fmt.Printf("Size: %d bytes (%.2f KB, %.2f MB)\n", bytes, kb, mb)

// Modified time
modTime := info.ModTime()
fmt.Println("Modified:", modTime.Format("2006-01-02 15:04:05"))

// Age
age := time.Since(modTime)
fmt.Printf("Age: %v\n", age)

// Type
if info.IsDir() {
    fmt.Println("Type: Directory")
} else {
    fmt.Println("Type: File")
}
```

---

## File Permissions

### Unix File Permissions

```
Format: rwxrwxrwx
        ├─┤├─┤├─┤
         │  │  └─ Others
         │  └──── Group
         └─────── Owner

r = read    (4)
w = write   (2)
x = execute (1)

Examples:
0644 = rw-r--r--  (owner: rw, group: r, others: r)
0755 = rwxr-xr-x  (owner: rwx, group: rx, others: rx)
0600 = rw-------  (owner: rw, others: none)
0777 = rwxrwxrwx  (all permissions for all)
```

### os.Chmod
**Змінити права доступу.**

```go
// Set permissions to rw-r--r-- (0644)
err := os.Chmod("file.txt", 0644)
if err != nil {
    panic(err)
}

// Set permissions to rwxr-xr-x (0755)
err = os.Chmod("script.sh", 0755)

// Remove all permissions
err = os.Chmod("secret.txt", 0000)
```

### Перевірка permissions з FileMode

```go
info, _ := os.Stat("file.txt")
mode := info.Mode()

// Check if file
if mode.IsRegular() {
    fmt.Println("Regular file")
}

// Check if directory
if mode.IsDir() {
    fmt.Println("Directory")
}

// Check permissions
perm := mode.Perm()
fmt.Printf("Permissions: %o\n", perm)  // octal format

// Check specific permission bits
if perm&0400 != 0 {
    fmt.Println("Owner can read")
}
if perm&0200 != 0 {
    fmt.Println("Owner can write")
}
if perm&0100 != 0 {
    fmt.Println("Owner can execute")
}
```

---

## File Ownership

### os.Chown
**Змінити власника файлу (тільки root або owner).**

```go
// Change owner to UID 1000, GID 1000
err := os.Chown("file.txt", 1000, 1000)
if err != nil {
    panic(err)
}

// Don't change owner or group (use -1)
os.Chown("file.txt", -1, 1000)  // only change group
```

---

## Symbolic Links

### os.Symlink
**Створити symbolic link.**

```go
// Create link: mylink -> target.txt
err := os.Symlink("target.txt", "mylink")
if err != nil {
    panic(err)
}
```

### os.Readlink
**Читати target symbolic link.**

```go
target, err := os.Readlink("mylink")
if err != nil {
    panic(err)
}
fmt.Println("Link points to:", target)
```

### Перевірка чи файл - symbolic link

```go
info, err := os.Lstat("mylink")
if err != nil {
    panic(err)
}

if info.Mode()&os.ModeSymlink != 0 {
    fmt.Println("This is a symbolic link")
    target, _ := os.Readlink("mylink")
    fmt.Println("Points to:", target)
}
```

---

## Temporary Files

### os.CreateTemp
**Створити тимчасовий файл.**

```go
// Create temp file in default temp directory
tmpFile, err := os.CreateTemp("", "myapp-*.txt")
if err != nil {
    panic(err)
}
defer os.Remove(tmpFile.Name())  // cleanup
defer tmpFile.Close()

fmt.Println("Temp file:", tmpFile.Name())
// Output: /tmp/myapp-123456789.txt

// Write to temp file
tmpFile.WriteString("Temporary data")
```

### os.MkdirTemp
**Створити тимчасову директорію.**

```go
// Create temp directory
tmpDir, err := os.MkdirTemp("", "myapp-")
if err != nil {
    panic(err)
}
defer os.RemoveAll(tmpDir)  // cleanup

fmt.Println("Temp dir:", tmpDir)
// Output: /tmp/myapp-987654321
```

### os.TempDir
**Отримати шлях до системної temp директорії.**

```go
tmpDir := os.TempDir()
fmt.Println("System temp dir:", tmpDir)
// macOS/Linux: /tmp
// Windows: C:\Users\Username\AppData\Local\Temp
```

---

## Практичні приклади

### 1. Перевірка файлу перед обробкою

```go
func validateFile(path string) error {
    info, err := os.Stat(path)
    if err != nil {
        return fmt.Errorf("file error: %w", err)
    }
    
    // Check if it's a regular file
    if !info.Mode().IsRegular() {
        return errors.New("not a regular file")
    }
    
    // Check size
    if info.Size() == 0 {
        return errors.New("file is empty")
    }
    
    if info.Size() > 100*1024*1024 {  // 100 MB
        return errors.New("file too large")
    }
    
    // Check age
    age := time.Since(info.ModTime())
    if age > 24*time.Hour {
        return errors.New("file is too old")
    }
    
    return nil
}
```

### 2. Безпечний запис файлу (atomic write)

```go
func safeWriteFile(filename string, data []byte) error {
    // Create temp file in same directory
    dir := filepath.Dir(filename)
    tmpFile, err := os.CreateTemp(dir, ".tmp-")
    if err != nil {
        return err
    }
    tmpName := tmpFile.Name()
    
    // Cleanup on error
    defer func() {
        tmpFile.Close()
        os.Remove(tmpName)
    }()
    
    // Write to temp file
    if _, err := tmpFile.Write(data); err != nil {
        return err
    }
    
    // Sync to disk
    if err := tmpFile.Sync(); err != nil {
        return err
    }
    
    // Close temp file
    if err := tmpFile.Close(); err != nil {
        return err
    }
    
    // Atomic rename (on same filesystem)
    if err := os.Rename(tmpName, filename); err != nil {
        return err
    }
    
    return nil
}
```

### 3. Копіювання з permissions

```go
func copyFileWithPerms(src, dst string) error {
    // Get source info
    srcInfo, err := os.Stat(src)
    if err != nil {
        return err
    }
    
    // Read source
    data, err := os.ReadFile(src)
    if err != nil {
        return err
    }
    
    // Write destination with same permissions
    err = os.WriteFile(dst, data, srcInfo.Mode().Perm())
    if err != nil {
        return err
    }
    
    // Copy modification time
    return os.Chtimes(dst, srcInfo.ModTime(), srcInfo.ModTime())
}
```

### 4. Пошук найновіших файлів

```go
func findNewestFiles(dir string, n int) ([]string, error) {
    entries, err := os.ReadDir(dir)
    if err != nil {
        return nil, err
    }
    
    type fileInfo struct {
        path    string
        modTime time.Time
    }
    
    var files []fileInfo
    
    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }
        
        info, err := entry.Info()
        if err != nil {
            continue
        }
        
        files = append(files, fileInfo{
            path:    filepath.Join(dir, entry.Name()),
            modTime: info.ModTime(),
        })
    }
    
    // Sort by modification time (newest first)
    sort.Slice(files, func(i, j int) bool {
        return files[i].modTime.After(files[j].modTime)
    })
    
    // Return top N
    result := make([]string, 0, n)
    for i := 0; i < n && i < len(files); i++ {
        result = append(result, files[i].path)
    }
    
    return result, nil
}
```

### 5. Cleanup старих файлів

```go
func cleanupOldFiles(dir string, maxAge time.Duration) error {
    entries, err := os.ReadDir(dir)
    if err != nil {
        return err
    }
    
    now := time.Now()
    
    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }
        
        info, err := entry.Info()
        if err != nil {
            continue
        }
        
        age := now.Sub(info.ModTime())
        if age > maxAge {
            path := filepath.Join(dir, entry.Name())
            fmt.Printf("Removing old file: %s (age: %v)\n", path, age)
            os.Remove(path)
        }
    }
    
    return nil
}

// Usage: cleanup files older than 7 days
cleanupOldFiles("/var/log", 7*24*time.Hour)
```

---

## os.Chtimes
**Змінити access та modification time.**

```go
// Set modification time to now
now := time.Now()
err := os.Chtimes("file.txt", now, now)

// Set to specific time
customTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
err = os.Chtimes("file.txt", customTime, customTime)
```

---

## FileMode Type

```go
type FileMode uint32

const (
    ModeDir        FileMode = 1 << (32 - 1 - iota) // d: directory
    ModeAppend                                      // a: append-only
    ModeExclusive                                   // l: exclusive use
    ModeTemporary                                   // T: temporary file
    ModeSymlink                                     // L: symbolic link
    ModeDevice                                      // D: device file
    ModeNamedPipe                                   // p: named pipe (FIFO)
    ModeSocket                                      // S: Unix domain socket
    ModeSetuid                                      // u: setuid
    ModeSetgid                                      // g: setgid
    ModeCharDevice                                  // c: character device
    ModeSticky                                      // t: sticky
    ModeIrregular                                   // ?: non-regular file
)
```

### Приклад перевірки типу файлу

```go
info, _ := os.Stat("somefile")
mode := info.Mode()

if mode.IsDir() {
    fmt.Println("Directory")
} else if mode.IsRegular() {
    fmt.Println("Regular file")
} else if mode&os.ModeSymlink != 0 {
    fmt.Println("Symbolic link")
} else if mode&os.ModeNamedPipe != 0 {
    fmt.Println("Named pipe")
} else if mode&os.ModeSocket != 0 {
    fmt.Println("Socket")
}
```

---

## Підсумок

| Операція | Функція | Опис |
|----------|---------|------|
| File info | `os.Stat` | Інформація про файл |
| Link info | `os.Lstat` | Інформація про link |
| Permissions | `os.Chmod` | Змінити права |
| Ownership | `os.Chown` | Змінити власника |
| Create link | `os.Symlink` | Створити symlink |
| Read link | `os.Readlink` | Читати symlink |
| Temp file | `os.CreateTemp` | Тимчасовий файл |
| Temp dir | `os.MkdirTemp` | Тимчасова dir |
| Times | `os.Chtimes` | Змінити access/mod time |

**Best Practices:**
- Використовуй `os.Stat` перед операціями з файлом
- Використовуй `defer os.Remove()` для cleanup temp files
- Перевіряй permissions перед читанням/записом
- Використовуй atomic writes для критичних файлів
- Cleanup старі temp files періодично
