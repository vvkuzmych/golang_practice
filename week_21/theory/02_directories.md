# Directory Operations

## Створення директорій

### os.Mkdir
**Створює одну директорію.**

```go
// Створити директорію з permissions 0755
err := os.Mkdir("mydir", 0755)
if err != nil {
    if os.IsExist(err) {
        fmt.Println("Directory already exists")
    } else {
        panic(err)
    }
}
```

### os.MkdirAll
**Створює всі відсутні директорії в шляху (як `mkdir -p`).**

```go
// Створити вкладені директорії
err := os.MkdirAll("path/to/nested/dir", 0755)
if err != nil {
    panic(err)
}

// Якщо директорія вже існує - не помилка
os.MkdirAll("existing/dir", 0755)  // OK
```

---

## Читання директорій

### os.ReadDir
**Читає вміст директорії (Go 1.16+).**

```go
entries, err := os.ReadDir(".")
if err != nil {
    panic(err)
}

for _, entry := range entries {
    fmt.Printf("%s (dir: %v)\n", entry.Name(), entry.IsDir())
}
```

### DirEntry методи

```go
entry := entries[0]

// Ім'я
name := entry.Name()

// Чи це директорія?
isDir := entry.IsDir()

// FileInfo (додаткова інформація)
info, err := entry.Info()
if err == nil {
    fmt.Println("Size:", info.Size())
    fmt.Println("ModTime:", info.ModTime())
}

// Type (file type bits)
fileType := entry.Type()
```

---

## Обхід директорій (Walk)

### filepath.Walk

```go
import (
    "fmt"
    "os"
    "path/filepath"
)

func main() {
    root := "."
    
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        // Вивести шлях
        fmt.Println(path)
        
        // Фільтр: тільки .go файли
        if !info.IsDir() && filepath.Ext(path) == ".go" {
            fmt.Println("  Go file:", path)
        }
        
        return nil
    })
    
    if err != nil {
        panic(err)
    }
}
```

### filepath.WalkDir (швидше, Go 1.16+)

```go
err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
    if err != nil {
        return err
    }
    
    fmt.Printf("%s (dir: %v)\n", path, d.IsDir())
    
    // Skip .git directory
    if d.IsDir() && d.Name() == ".git" {
        return filepath.SkipDir
    }
    
    return nil
})
```

---

## Видалення директорій

### os.Remove
**Видаляє ПОРОЖНЮ директорію.**

```go
err := os.Remove("emptydir")
if err != nil {
    panic(err)
}
```

### os.RemoveAll
**Видаляє директорію та весь вміст (як `rm -rf`).**

```go
// ОБЕРЕЖНО: видалить все всередині!
err := os.RemoveAll("mydir")
if err != nil {
    panic(err)
}
```

---

## Working Directory

### os.Getwd
**Отримати поточну робочу директорію.**

```go
wd, err := os.Getwd()
if err != nil {
    panic(err)
}
fmt.Println("Current directory:", wd)
```

### os.Chdir
**Змінити робочу директорію.**

```go
err := os.Chdir("/tmp")
if err != nil {
    panic(err)
}

// Перевірити
wd, _ := os.Getwd()
fmt.Println("Now in:", wd)  // /tmp
```

---

## Практичні приклади

### 1. Рекурсивне копіювання директорії

```go
func copyDir(src, dst string) error {
    // Створити destination directory
    err := os.MkdirAll(dst, 0755)
    if err != nil {
        return err
    }
    
    // Читати source directory
    entries, err := os.ReadDir(src)
    if err != nil {
        return err
    }
    
    for _, entry := range entries {
        srcPath := filepath.Join(src, entry.Name())
        dstPath := filepath.Join(dst, entry.Name())
        
        if entry.IsDir() {
            // Рекурсивно копіювати sub-directory
            if err := copyDir(srcPath, dstPath); err != nil {
                return err
            }
        } else {
            // Копіювати файл
            if err := copyFile(srcPath, dstPath); err != nil {
                return err
            }
        }
    }
    
    return nil
}

func copyFile(src, dst string) error {
    sourceFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer sourceFile.Close()
    
    destFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer destFile.Close()
    
    _, err = io.Copy(destFile, sourceFile)
    return err
}
```

### 2. Пошук файлів за розширенням

```go
func findFiles(root, ext string) ([]string, error) {
    var files []string
    
    err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        
        if !d.IsDir() && filepath.Ext(path) == ext {
            files = append(files, path)
        }
        
        return nil
    })
    
    return files, err
}

// Використання
goFiles, err := findFiles(".", ".go")
if err != nil {
    panic(err)
}

for _, file := range goFiles {
    fmt.Println(file)
}
```

### 3. Обчислити розмір директорії

```go
func dirSize(root string) (int64, error) {
    var size int64
    
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if !info.IsDir() {
            size += info.Size()
        }
        
        return nil
    })
    
    return size, err
}

// Використання
size, err := dirSize(".")
if err != nil {
    panic(err)
}

fmt.Printf("Directory size: %.2f MB\n", float64(size)/1024/1024)
```

### 4. Створити структуру проекту

```go
func createProjectStructure(root string) error {
    dirs := []string{
        "cmd/api",
        "internal/domain",
        "internal/usecase",
        "internal/adapter/http",
        "internal/infrastructure/postgres",
        "pkg/logger",
        "migrations",
        "configs",
        "docs",
    }
    
    for _, dir := range dirs {
        path := filepath.Join(root, dir)
        if err := os.MkdirAll(path, 0755); err != nil {
            return err
        }
    }
    
    // Створити .gitkeep файли
    for _, dir := range dirs {
        gitkeep := filepath.Join(root, dir, ".gitkeep")
        if err := os.WriteFile(gitkeep, []byte(""), 0644); err != nil {
            return err
        }
    }
    
    return nil
}

// Використання
err := createProjectStructure("myproject")
if err != nil {
    panic(err)
}
```

### 5. Tree command (показати структуру)

```go
func printTree(root string, prefix string) error {
    entries, err := os.ReadDir(root)
    if err != nil {
        return err
    }
    
    for i, entry := range entries {
        isLast := i == len(entries)-1
        
        // Вибрати символи
        var branch string
        if isLast {
            branch = "└── "
        } else {
            branch = "├── "
        }
        
        // Вивести entry
        fmt.Println(prefix + branch + entry.Name())
        
        // Рекурсивно для sub-directories
        if entry.IsDir() {
            newPrefix := prefix
            if isLast {
                newPrefix += "    "
            } else {
                newPrefix += "│   "
            }
            
            subPath := filepath.Join(root, entry.Name())
            printTree(subPath, newPrefix)
        }
    }
    
    return nil
}

// Використання
fmt.Println(".")
printTree(".", "")

// Output:
// .
// ├── go.mod
// ├── go.sum
// ├── cmd
// │   └── api
// │       └── main.go
// └── internal
//     └── domain
//         └── user.go
```

---

## filepath package (допоміжні функції)

### filepath.Join
**Кросплатформове об'єднання шляхів.**

```go
// Windows: path\to\file.txt
// Unix:    path/to/file.txt
path := filepath.Join("path", "to", "file.txt")
```

### filepath.Split
**Розділити шлях на директорію і файл.**

```go
dir, file := filepath.Split("/path/to/file.txt")
// dir:  "/path/to/"
// file: "file.txt"
```

### filepath.Base
**Отримати ім'я файлу.**

```go
name := filepath.Base("/path/to/file.txt")  // "file.txt"
```

### filepath.Dir
**Отримати директорію.**

```go
dir := filepath.Dir("/path/to/file.txt")  // "/path/to"
```

### filepath.Ext
**Отримати розширення.**

```go
ext := filepath.Ext("file.txt")  // ".txt"
```

### filepath.Abs
**Отримати абсолютний шлях.**

```go
abs, err := filepath.Abs("relative/path")
// abs: "/home/user/project/relative/path"
```

### filepath.Rel
**Отримати відносний шлях.**

```go
rel, err := filepath.Rel("/home/user", "/home/user/project/main.go")
// rel: "project/main.go"
```

### filepath.Clean
**Очистити шлях (видалити .., .).**

```go
clean := filepath.Clean("path/./to/../file.txt")
// clean: "path/file.txt"
```

---

## Best Practices

### 1. Використовуй filepath.Join

```go
// ✅ Good (кросплатформове)
path := filepath.Join("dir", "subdir", "file.txt")

// ❌ Bad (працює тільки на Unix)
path := "dir/subdir/file.txt"
```

### 2. Перевіряй чи існує перед видаленням

```go
// ✅ Good
if _, err := os.Stat("dir"); err == nil {
    os.RemoveAll("dir")
}

// ❌ Bad (panic якщо не існує)
os.RemoveAll("dir")  // OK, але краще перевіряти
```

### 3. Skip .git, node_modules в Walk

```go
filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
    // Skip hidden directories
    if d.IsDir() && strings.HasPrefix(d.Name(), ".") {
        return filepath.SkipDir
    }
    
    // Skip node_modules
    if d.IsDir() && d.Name() == "node_modules" {
        return filepath.SkipDir
    }
    
    return nil
})
```

---

## Підсумок

| Операція | Функція | Опис |
|----------|---------|------|
| Створити dir | `os.Mkdir` | Одна директорія |
| Створити nested | `os.MkdirAll` | Вкладені директорії |
| Читати dir | `os.ReadDir` | Список entries |
| Обійти дерево | `filepath.Walk` | Рекурсивний обхід |
| Видалити порожню | `os.Remove` | Тільки порожню dir |
| Видалити все | `os.RemoveAll` | Dir + вміст |
| Поточна dir | `os.Getwd` | Working directory |
| Змінити dir | `os.Chdir` | Change directory |
| Об'єднати шлях | `filepath.Join` | Кросплатформово |
