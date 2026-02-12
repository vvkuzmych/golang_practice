# Week 21 — os Package

**Ціль:** Повністю освоїти роботу з операційною системою через `os` package в Go.

---

## 📚 Теорія

### [01. File Operations](./theory/01_file_operations.md)
- Створення, читання, запис файлів
- os.Open, os.Create, os.OpenFile
- Читання/запис за допомогою os.File
- Закриття файлів (defer)
- Append, Truncate

### [02. Directory Operations](./theory/02_directories.md)
- os.Mkdir, os.MkdirAll
- os.ReadDir
- Обхід директорій (Walk)
- os.Remove, os.RemoveAll
- os.Rename

### [03. Environment & Process](./theory/03_environment_process.md)
- os.Getenv, os.Setenv, os.Environ
- os.Args (command-line arguments)
- os.Exit, os.Getpid
- os.Hostname
- Working directory (Getwd, Chdir)

### [04. File Info & Permissions](./theory/04_file_info_permissions.md)
- os.Stat, os.Lstat
- os.FileInfo interface
- File permissions (os.Chmod)
- File ownership (os.Chown)
- Symbolic links (os.Symlink, os.Readlink)
- Temporary files (os.CreateTemp)

---

## 🛠️ Практика

### [01. File Manager CLI](./practice/01_file_manager/)
- Команди: ls, cat, cp, mv, rm
- Використання os package

### [02. Log Rotator](./practice/02_log_rotator/)
- Ротація логів по розміру/даті
- os.Stat для перевірки розміру

### [03. Directory Sync](./practice/03_directory_sync/)
- Синхронізація двох директорій
- Порівняння файлів

### [04. Config Manager](./practice/04_config_manager/)
- Читання/запис конфігів
- Environment variables

---

## 📝 Exercises

### [Exercise 1: File Copy Tool](./exercises/exercise_1.md)
Створити утиліту для копіювання файлів з прогрес-баром.

### [Exercise 2: Directory Tree](./exercises/exercise_2.md)
Реалізувати `tree` command (показати структуру директорій).

### [Exercise 3: File Search](./exercises/exercise_3.md)
Пошук файлів за ім'ям/розміром/датою.

---

## 🎯 Learning Outcomes

Після цього тижня ви зможете:
- ✅ Створювати, читати, записувати файли
- ✅ Працювати з директоріями (створення, читання, видалення)
- ✅ Використовувати environment variables
- ✅ Отримувати інформацію про файли (розмір, дата, права)
- ✅ Працювати з file permissions
- ✅ Створювати CLI tools для роботи з файловою системою

---

## 📖 Key Concepts

### File Modes
```go
os.O_RDONLY  // Read-only
os.O_WRONLY  // Write-only
os.O_RDWR    // Read-write
os.O_APPEND  // Append to file
os.O_CREATE  // Create if doesn't exist
os.O_TRUNC   // Truncate file
```

### File Permissions (Unix)
```go
0644  // rw-r--r-- (owner: rw, group: r, others: r)
0755  // rwxr-xr-x (owner: rwx, group: rx, others: rx)
0600  // rw------- (owner: rw, others: none)
```

### Best Practices
1. Завжди закривай файли (`defer file.Close()`)
2. Перевіряй помилки після кожної операції
3. Використовуй `os.OpenFile` для точного контролю
4. Використовуй `filepath` package для кросплатформових шляхів
5. Використовуй `os.CreateTemp` для тимчасових файлів

---

## 📖 Additional Resources

- [Go os package documentation](https://pkg.go.dev/os)
- [Go filepath package](https://pkg.go.dev/path/filepath)
- [Working with Files in Go](https://gobyexample.com/reading-files)
- [File I/O in Go](https://yourbasic.org/golang/read-file-line-by-line/)

---

**Previous:** [Week 20 — System Design](../week_20/README.md)
