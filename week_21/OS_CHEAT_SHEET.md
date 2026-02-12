# os Package — Cheat Sheet

## 📁 File Operations

```go
// Create/Open
os.Create("file.txt")                                    // Create or truncate
os.Open("file.txt")                                      // Open read-only
os.OpenFile("file.txt", os.O_RDWR|os.O_CREATE, 0644)    // Custom flags

// Read
os.ReadFile("file.txt")                                  // Read entire file
file.Read(buffer)                                        // Read into buffer

// Write
os.WriteFile("file.txt", data, 0644)                    // Write entire file
file.Write([]byte("data"))                              // Write bytes
file.WriteString("data")                                // Write string

// Copy/Move/Delete
io.Copy(dst, src)                                       // Copy file
os.Rename("old.txt", "new.txt")                         // Rename/move
os.Remove("file.txt")                                   // Delete file

// Info
os.Stat("file.txt")                                     // File info
os.Truncate("file.txt", 0)                              // Truncate
```

---

## 📂 Directory Operations

```go
// Create
os.Mkdir("dir", 0755)                                   // Single directory
os.MkdirAll("path/to/dir", 0755)                        // Nested directories

// Read
os.ReadDir("dir")                                       // List entries
filepath.Walk("dir", func(path, info, err) {...})       // Recursive walk
filepath.WalkDir("dir", func(path, d, err) {...})       // Faster walk

// Delete
os.Remove("emptydir")                                   // Remove empty dir
os.RemoveAll("dir")                                     // Remove dir + contents

// Navigation
os.Getwd()                                              // Current directory
os.Chdir("/path")                                       // Change directory
```

---

## 🌍 Environment Variables

```go
// Get
os.Getenv("PORT")                                       // Get variable
os.LookupEnv("PORT")                                    // Check if exists
os.Environ()                                            // All variables
os.ExpandEnv("$HOME/projects")                          // Expand variables

// Set
os.Setenv("MY_VAR", "value")                            // Set variable
os.Unsetenv("MY_VAR")                                   // Remove variable
```

---

## ⚙️ Process Management

```go
// Process info
os.Getpid()                                             // Process ID
os.Getppid()                                            // Parent PID
os.Hostname()                                           // Hostname
os.Getuid()                                             // User ID (Unix)
os.Getgid()                                             // Group ID (Unix)

// Arguments
os.Args                                                 // Command-line args
flag.String("name", "default", "description")           // Parse flags

// Exit
os.Exit(0)                                              // Exit with code

// Directories
os.UserHomeDir()                                        // Home directory
os.UserCacheDir()                                       // Cache directory
os.UserConfigDir()                                      // Config directory
os.TempDir()                                            // Temp directory
```

---

## ℹ️ File Info & Permissions

```go
// File info
info, _ := os.Stat("file.txt")
info.Name()                                             // File name
info.Size()                                             // Size in bytes
info.Mode()                                             // File mode/permissions
info.ModTime()                                          // Modification time
info.IsDir()                                            // Is directory?

// Permissions
os.Chmod("file.txt", 0644)                              // Change permissions
os.Chown("file.txt", uid, gid)                          // Change owner (Unix)
os.Chtimes("file.txt", atime, mtime)                    // Change times

// Symlinks
os.Symlink("target", "link")                            // Create symlink
os.Readlink("link")                                     // Read symlink target
os.Lstat("link")                                        // Info about link itself

// Temp files
os.CreateTemp("", "prefix-*.txt")                       // Create temp file
os.MkdirTemp("", "prefix-")                             // Create temp directory
```

---

## 🔍 File Checks

```go
// Check existence
if _, err := os.Stat("file.txt"); err == nil {
    // File exists
}

// Error checks
os.IsNotExist(err)                                      // File doesn't exist
os.IsExist(err)                                         // File already exists
os.IsPermission(err)                                    // Permission denied
os.IsTimeout(err)                                       // Operation timed out
```

---

## 🎭 File Modes & Permissions

```go
// Open modes
os.O_RDONLY                                             // Read-only
os.O_WRONLY                                             // Write-only
os.O_RDWR                                               // Read-write
os.O_APPEND                                             // Append
os.O_CREATE                                             // Create if not exists
os.O_TRUNC                                              // Truncate
os.O_EXCL                                               // Exclusive (with O_CREATE)
os.O_SYNC                                               // Synchronous I/O

// Permissions (Unix)
0644                                                    // rw-r--r--
0755                                                    // rwxr-xr-x
0600                                                    // rw-------
0777                                                    // rwxrwxrwx

// File types
mode.IsDir()                                            // Directory
mode.IsRegular()                                        // Regular file
mode&os.ModeSymlink != 0                                // Symbolic link
mode&os.ModeNamedPipe != 0                              // Named pipe
```

---

## 🛠️ filepath Package

```go
// Path operations
filepath.Join("dir", "file.txt")                        // Join paths
filepath.Split("/path/to/file.txt")                     // Split to dir + file
filepath.Base("/path/to/file.txt")                      // Get filename
filepath.Dir("/path/to/file.txt")                       // Get directory
filepath.Ext("file.txt")                                // Get extension
filepath.Abs("relative/path")                           // Absolute path
filepath.Rel("/base", "/base/sub/file")                 // Relative path
filepath.Clean("path/./to/../file")                     // Clean path
```

---

## 📋 Common Patterns

### Always Close Files
```go
file, err := os.Open("file.txt")
if err != nil {
    return err
}
defer file.Close()  // Always use defer
```

### Check Errors
```go
if err := os.Remove("file.txt"); err != nil {
    if os.IsNotExist(err) {
        // File doesn't exist
    } else {
        return err
    }
}
```

### Atomic Write
```go
tmp, _ := os.CreateTemp("", "*.tmp")
defer os.Remove(tmp.Name())
tmp.Write(data)
tmp.Close()
os.Rename(tmp.Name(), "target.txt")
```

### Read Line by Line
```go
file, _ := os.Open("file.txt")
defer file.Close()
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    line := scanner.Text()
    // Process line
}
```

### Copy with Progress
```go
src, _ := os.Open("source.txt")
defer src.Close()
dst, _ := os.Create("dest.txt")
defer dst.Close()
io.Copy(dst, src)
```

---

## 🚨 Common Mistakes

### ❌ Forgot to Close
```go
file, _ := os.Open("file.txt")
// Memory leak!
```

### ❌ Ignored Errors
```go
os.WriteFile("file.txt", data, 0644)  // No error check
```

### ❌ Hardcoded Paths
```go
path := "dir/file.txt"  // Won't work on Windows
```

### ❌ defer in Loop
```go
for _, file := range files {
    f, _ := os.Open(file)
    defer f.Close()  // Won't close until function exits!
}
```

---

## ✅ Best Practices

1. **Always close files** with `defer`
2. **Check all errors** - don't ignore
3. **Use filepath.Join** for cross-platform paths
4. **Use os.CreateTemp** for temp files
5. **Validate file info** before operations
6. **Handle permissions** correctly (0644 for files, 0755 for dirs)
7. **Use atomic writes** for critical files
8. **Cleanup temp files** with `defer os.Remove()`

---

## 📚 Quick Reference

| Task | Function |
|------|----------|
| Read file | `os.ReadFile("file.txt")` |
| Write file | `os.WriteFile("file.txt", data, 0644)` |
| Create dir | `os.MkdirAll("path/to/dir", 0755)` |
| List dir | `os.ReadDir(".")` |
| Walk tree | `filepath.WalkDir(".", func...)` |
| File info | `os.Stat("file.txt")` |
| Copy file | `io.Copy(dst, src)` |
| Delete | `os.Remove("file.txt")` |
| ENV var | `os.Getenv("PORT")` |
| Temp file | `os.CreateTemp("", "*.tmp")` |
