# Solutions - How to Run

## ⚠️ Important: Multiple main() Functions

Each solution file is a **separate, independent program** with its own `main()` function. You **cannot** compile or run them all together.

---

## ✅ How to Run Solutions

### Option 1: Using `go run` (Recommended)

Run each solution individually:

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_1/solutions

# Solution 1 - User Greeting
go run solution_1.go Іван 25
go run solution_1.go Марія
go run solution_1.go

# Solution 2 - Student Info
go run solution_2.go

# Solution 3 - TODO Manager
go run solution_3.go add "Вивчити Go"
go run solution_3.go list
go run solution_3.go done 1
go run solution_3.go delete 1
go run solution_3.go help
```

---

### Option 2: Build Each Separately

```bash
# Build each solution into its own executable
go build -o greet solution_1.go
go build -o student solution_2.go
go build -o todo solution_3.go

# Run the executables
./greet Петро 30
./student
./todo add "Моє завдання"
./todo list
```

---

## 🚫 What NOT to Do

```bash
# ❌ This will FAIL with "main redeclared" error
go build .
go run *.go
go build solution_1.go solution_2.go solution_3.go
```

**Why?** Go sees multiple `main()` functions and doesn't know which one to use as the entry point.

---

## 💡 In GoLand/IDE

### To run a specific solution:

1. **Open** the solution file (e.g., `solution_1.go`)
2. **Right-click** in the editor
3. **Select** "Run 'go build solution_1.go'" or "Run 'go run solution_1.go'"
4. **Add arguments** if needed (Run → Edit Configurations → Program arguments)

### Or use the green play button ▶️ next to `func main()`

---

## 📁 Alternative: Separate Directories (Optional)

If you want to build all solutions at once, restructure like this:

```
solutions/
├── README.md
├── solution_1/
│   └── main.go        (rename solution_1.go)
├── solution_2/
│   └── main.go        (rename solution_2.go)
└── solution_3/
    └── main.go        (rename solution_3.go)
```

Then you can build all:

```bash
go build ./solution_1
go build ./solution_2
go build ./solution_3
```

But for learning purposes, **keeping them in one directory is fine** - just run them individually!

---

## 📝 Quick Reference

| Solution | What it does | Example command |
|----------|--------------|-----------------|
| `solution_1.go` | User greeting with arguments | `go run solution_1.go Іван 25` |
| `solution_2.go` | Student info (structs demo) | `go run solution_2.go` |
| `solution_3.go` | TODO Manager CLI | `go run solution_3.go help` |

---

## 🎯 Expected Behavior

### Solution 1
```bash
$ go run solution_1.go Іван 25
Доброго вечора, Іван! 👋
Тобі 25 років.
Продуктивного дня! 💼
```

### Solution 2
```bash
$ go run solution_2.go
=== Інформація про студентів ===
[Shows student information with structs]
```

### Solution 3
```bash
$ go run solution_3.go add "Test task"
✅ Завдання додано: "Test task" (ID: 1)

$ go run solution_3.go list
=== TODO List ===
ID | Статус | Завдання | Створено
...
```

---

**Happy coding! 🚀**

