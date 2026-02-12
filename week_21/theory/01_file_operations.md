# File Operations

## Створення і відкриття файлів

### os.Create
**Створює новий файл або truncate існуючий.**

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    file, err := os.Create("test.txt")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer file.Close()  // ВАЖЛИВО: завжди закривай файл
    
    fmt.Println("File created:", file.Name())
}
```

**Еквівалентно:**
```go
os.OpenFile("test.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
```

### os.Open
**Відкриває файл тільки для читання.**

```go
file, err := os.Open("test.txt")
if err != nil {
    if os.IsNotExist(err) {
        fmt.Println("File doesn't exist")
    }
    return
}
defer file.Close()
```

### os.OpenFile
**Найбільш гнучкий спосіб відкриття файлу.**

```go
// Flags
const (
    O_RDONLY int = syscall.O_RDONLY // read-only
    O_WRONLY int = syscall.O_WRONLY // write-only
    O_RDWR   int = syscall.O_RDWR   // read-write
    O_APPEND int = syscall.O_APPEND // append when writing
    O_CREATE int = syscall.O_CREAT  // create if not exists
    O_EXCL   int = syscall.O_EXCL   // used with O_CREATE, fail if exists
    O_SYNC   int = syscall.O_SYNC   // sync writes
    O_TRUNC  int = syscall.O_TRUNC  // truncate file
)

// Приклад: відкрити для запису, створити якщо немає
file, err := os.OpenFile("data.txt", os.O_WRONLY|os.O_CREATE, 0644)
if err != nil {
    panic(err)
}
defer file.Close()
```

---

## Читання файлів

### 1. Читання всього файлу

```go
// os.ReadFile (простіший спосіб з Go 1.16+)
data, err := os.ReadFile("test.txt")
if err != nil {
    panic(err)
}
fmt.Println(string(data))
```

### 2. Читання через os.File

```go
file, err := os.Open("test.txt")
if err != nil {
    panic(err)
}
defer file.Close()

// Читання в buffer
buffer := make([]byte, 1024)
n, err := file.Read(buffer)
if err != nil && err != io.EOF {
    panic(err)
}
fmt.Printf("Read %d bytes: %s\n", n, buffer[:n])
```

### 3. Читання по рядках

```go
import (
    "bufio"
    "os"
)

file, err := os.Open("test.txt")
if err != nil {
    panic(err)
}
defer file.Close()

scanner := bufio.NewScanner(file)
for scanner.Scan() {
    line := scanner.Text()
    fmt.Println(line)
}

if err := scanner.Err(); err != nil {
    panic(err)
}
```

### 4. Читання з позиції (Seek)

```go
file, _ := os.Open("test.txt")
defer file.Close()

// Перемістити курсор на 10 байт від початку
file.Seek(10, io.SeekStart)

// Читати з позиції
buffer := make([]byte, 5)
file.Read(buffer)
fmt.Println(string(buffer))

// Seek modes:
// io.SeekStart   - від початку файлу
// io.SeekCurrent - від поточної позиції
// io.SeekEnd     - від кінця файлу
```

---

## Запис у файли

### 1. Простий запис

```go
// os.WriteFile (простіший спосіб з Go 1.16+)
data := []byte("Hello, World!\n")
err := os.WriteFile("output.txt", data, 0644)
if err != nil {
    panic(err)
}
```

### 2. Запис через os.File

```go
file, err := os.Create("output.txt")
if err != nil {
    panic(err)
}
defer file.Close()

// Записати bytes
n, err := file.Write([]byte("Hello, "))
if err != nil {
    panic(err)
}
fmt.Printf("Wrote %d bytes\n", n)

// Записати string
n, err = file.WriteString("World!\n")
if err != nil {
    panic(err)
}
fmt.Printf("Wrote %d bytes\n", n)
```

### 3. Append (додавання в кінець)

```go
file, err := os.OpenFile("log.txt", 
    os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
if err != nil {
    panic(err)
}
defer file.Close()

file.WriteString("New log entry\n")
```

### 4. Буферизований запис

```go
import "bufio"

file, _ := os.Create("output.txt")
defer file.Close()

writer := bufio.NewWriter(file)
writer.WriteString("Line 1\n")
writer.WriteString("Line 2\n")
writer.Flush()  // ВАЖЛИВО: flush buffer
```

---

## Копіювання файлів

```go
func copyFile(src, dst string) error {
    // Відкрити source
    sourceFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer sourceFile.Close()
    
    // Створити destination
    destFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer destFile.Close()
    
    // Копіювати
    _, err = io.Copy(destFile, sourceFile)
    return err
}

// Використання
err := copyFile("source.txt", "destination.txt")
if err != nil {
    panic(err)
}
```

---

## Видалення файлів

```go
// Видалити файл
err := os.Remove("test.txt")
if err != nil {
    if os.IsNotExist(err) {
        fmt.Println("File doesn't exist")
    } else {
        panic(err)
    }
}

// Видалити директорію та весь вміст
err = os.RemoveAll("mydir")
if err != nil {
    panic(err)
}
```

---

## Перейменування/переміщення

```go
// Перейменувати або перемістити файл
err := os.Rename("old.txt", "new.txt")
if err != nil {
    panic(err)
}

// Перемістити в іншу директорію
err = os.Rename("file.txt", "backup/file.txt")
```

---

## Truncate (обрізати файл)

```go
// Обрізати до 0 bytes (очистити файл)
err := os.Truncate("test.txt", 0)

// Обрізати до 100 bytes
err = os.Truncate("test.txt", 100)
```

---

## Практичний приклад: Лог-файл

```go
package main

import (
    "fmt"
    "os"
    "time"
)

type Logger struct {
    file *os.File
}

func NewLogger(filename string) (*Logger, error) {
    file, err := os.OpenFile(filename, 
        os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }
    return &Logger{file: file}, nil
}

func (l *Logger) Log(message string) error {
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    logEntry := fmt.Sprintf("[%s] %s\n", timestamp, message)
    _, err := l.file.WriteString(logEntry)
    return err
}

func (l *Logger) Close() error {
    return l.file.Close()
}

func main() {
    logger, err := NewLogger("app.log")
    if err != nil {
        panic(err)
    }
    defer logger.Close()
    
    logger.Log("Application started")
    logger.Log("Processing data...")
    logger.Log("Application finished")
}
```

---

## Перевірка помилок

```go
file, err := os.Open("test.txt")
if err != nil {
    // Специфічні помилки
    if os.IsNotExist(err) {
        fmt.Println("File doesn't exist")
    } else if os.IsPermission(err) {
        fmt.Println("Permission denied")
    } else if os.IsExist(err) {
        fmt.Println("File already exists")
    } else if os.IsTimeout(err) {
        fmt.Println("Operation timed out")
    } else {
        fmt.Println("Unknown error:", err)
    }
    return
}
defer file.Close()
```

---

## Best Practices

### 1. Завжди закривай файли

```go
// ✅ Good
file, err := os.Open("test.txt")
if err != nil {
    return err
}
defer file.Close()  // Закриється автоматично

// ❌ Bad
file, _ := os.Open("test.txt")
// forgot to close - memory leak!
```

### 2. Перевіряй помилки

```go
// ✅ Good
n, err := file.Write(data)
if err != nil {
    return err
}
if n != len(data) {
    return errors.New("incomplete write")
}

// ❌ Bad
file.Write(data)  // ignoring errors
```

### 3. Використовуй defer для cleanup

```go
file, err := os.Create("temp.txt")
if err != nil {
    return err
}
defer os.Remove("temp.txt")  // cleanup temp file
defer file.Close()

// do work...
```

### 4. Використовуй io.Copy для великих файлів

```go
// ✅ Good (memory efficient)
io.Copy(dst, src)

// ❌ Bad (loads entire file into memory)
data, _ := os.ReadFile("large.txt")
os.WriteFile("copy.txt", data, 0644)
```

---

## Підсумок

| Операція | Функція | Опис |
|----------|---------|------|
| Створити | `os.Create` | Створює новий файл |
| Відкрити (read) | `os.Open` | Відкрити для читання |
| Відкрити (custom) | `os.OpenFile` | Відкрити з flags |
| Читати все | `os.ReadFile` | Читати весь файл |
| Записати все | `os.WriteFile` | Записати весь файл |
| Копіювати | `io.Copy` | Копіювати file → file |
| Видалити | `os.Remove` | Видалити файл |
| Перейменувати | `os.Rename` | Перейменувати/перемістити |
| Truncate | `os.Truncate` | Обрізати файл |
