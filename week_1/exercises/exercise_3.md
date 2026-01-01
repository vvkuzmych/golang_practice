# Вправа 3: Менеджер завдань (TODO CLI)

## Ціль
Створити повноцінний CLI додаток, який демонструє всі концепції тижня 1.

---

## Завдання

Створіть програму `todo.go` - простий менеджер завдань з командним рядком.

---

## Функціонал

### Команди:

1. **add** - Додати завдання
   ```bash
   go run todo.go add "Купити молоко"
   ```

2. **list** - Показати всі завдання
   ```bash
   go run todo.go list
   ```

3. **done** - Позначити завдання як виконане
   ```bash
   go run todo.go done 1
   ```

4. **delete** - Видалити завдання
   ```bash
   go run todo.go delete 1
   ```

5. **help** - Показати допомогу
   ```bash
   go run todo.go help
   ```

---

## Структури даних

```go
type Task struct {
    ID          int
    Description string
    Done        bool
    CreatedAt   string
}

type TodoList struct {
    Tasks []Task
}
```

---

## Приклад використання

```bash
# Додати завдання
$ go run todo.go add "Вивчити Go"
✅ Завдання додано: "Вивчити Go" (ID: 1)

$ go run todo.go add "Написати код"
✅ Завдання додано: "Написати код" (ID: 2)

$ go run todo.go add "Протестувати"
✅ Завдання додано: "Протестувати" (ID: 3)

# Показати список
$ go run todo.go list
=== TODO List ===

ID | Статус | Завдання           | Створено
---|--------|--------------------|-----------------
1  | [ ]    | Вивчити Go         | 2024-01-15 10:30
2  | [ ]    | Написати код       | 2024-01-15 10:31
3  | [ ]    | Протестувати       | 2024-01-15 10:32

Всього: 3 завдань (0 виконано, 3 активних)

# Позначити як виконане
$ go run todo.go done 1
✅ Завдання #1 позначено як виконане

$ go run todo.go list
ID | Статус | Завдання           | Створено
---|--------|--------------------|-----------------
1  | [✓]    | Вивчити Go         | 2024-01-15 10:30
2  | [ ]    | Написати код       | 2024-01-15 10:31
3  | [ ]    | Протестувати       | 2024-01-15 10:32

# Видалити завдання
$ go run todo.go delete 2
✅ Завдання #2 видалено

# Без аргументів
$ go run todo.go
Помилка: Команда не вказана

Використання: todo <команда> [аргументи]

Доступні команди:
  add <text>   - Додати завдання
  list         - Показати всі завдання
  done <id>    - Позначити як виконане
  delete <id>  - Видалити завдання
  help         - Показати цю допомогу
```

---

## Вимоги

### Обов'язкові:

1. ✅ Використати структури `Task` та `TodoList`
2. ✅ Використати os.Args для команд
3. ✅ Використати fmt.Printf для форматованого виводу
4. ✅ Обробити всі команди (add, list, done, delete, help)
5. ✅ Обробка помилок (невірна команда, невірний ID)
6. ✅ Використати як var, так і :=
7. ✅ Додати коментарі до коду

### Опціональні:

1. ⭐ Зберігати дані в файл (JSON)
2. ⭐ Додати пріоритети (high, medium, low)
3. ⭐ Додати дедлайни
4. ⭐ Фільтрувати (активні/виконані)
5. ⭐ Кольоровий вивід
6. ⭐ Пошук завдань

---

## Структура програми

```go
package main

import (
    "fmt"
    "os"
    "strconv"
    "time"
)

type Task struct {
    ID          int
    Description string
    Done        bool
    CreatedAt   string
}

type TodoList struct {
    Tasks []Task
}

func main() {
    // 1. Перевірка аргументів
    // 2. Отримання команди
    // 3. Виконання команди
}

func addTask(list *TodoList, description string) {
    // Додати завдання
}

func listTasks(list TodoList) {
    // Показати всі завдання
}

func markDone(list *TodoList, id int) error {
    // Позначити як виконане
}

func deleteTask(list *TodoList, id int) error {
    // Видалити завдання
}

func printHelp() {
    // Показати допомогу
}
```

---

## Підказки

### 1. Робота з slice
```go
// Додати елемент
tasks = append(tasks, newTask)

// Видалити елемент за індексом
tasks = append(tasks[:index], tasks[index+1:]...)

// Знайти елемент
for i, task := range tasks {
    if task.ID == targetID {
        // знайдено
    }
}
```

### 2. Поточний час
```go
import "time"

now := time.Now()
formatted := now.Format("2006-01-02 15:04")
```

### 3. Форматована таблиця
```go
fmt.Printf("%-4s | %-6s | %-20s | %s\n", 
    "ID", "Статус", "Завдання", "Створено")
fmt.Println("-----|--------|----------------------|------------------")

for _, task := range tasks {
    status := "[ ]"
    if task.Done {
        status = "[✓]"
    }
    fmt.Printf("%-4d | %-6s | %-20s | %s\n",
        task.ID, status, task.Description, task.CreatedAt)
}
```

### 4. Конвертація string → int
```go
import "strconv"

id, err := strconv.Atoi(idStr)
if err != nil {
    fmt.Println("Помилка: ID має бути числом")
    return
}
```

---

## Бонус: Збереження в файл

```go
import (
    "encoding/json"
    "io/ioutil"
)

func saveTasks(list TodoList) error {
    data, err := json.MarshalIndent(list.Tasks, "", "  ")
    if err != nil {
        return err
    }
    return ioutil.WriteFile("tasks.json", data, 0644)
}

func loadTasks() (TodoList, error) {
    data, err := ioutil.ReadFile("tasks.json")
    if err != nil {
        return TodoList{}, err
    }
    
    var tasks []Task
    err = json.Unmarshal(data, &tasks)
    return TodoList{Tasks: tasks}, err
}
```

---

## Критерії оцінки

- ✅ Всі команди працюють
- ✅ Обробка помилок
- ✅ Форматований вивід
- ✅ Структури використані правильно
- ✅ Код чистий і зрозумілий
- ✅ Коментарі присутні
- ✅ var та := використані доречно

---

## Приклад повного виводу

```
$ go run todo.go add "Вивчити типи в Go"
✅ Завдання додано: "Вивчити типи в Go" (ID: 1)

$ go run todo.go add "Зрозуміти zero values"
✅ Завдання додано: "Зрозуміти zero values" (ID: 2)

$ go run todo.go add "Написати CLI програму"
✅ Завдання додано: "Написати CLI програму" (ID: 3)

$ go run todo.go list
=== TODO List ===

ID   | Статус | Завдання                  | Створено
-----|--------|---------------------------|------------------
1    | [ ]    | Вивчити типи в Go        | 2024-01-15 10:30
2    | [ ]    | Зрозуміти zero values    | 2024-01-15 10:31
3    | [ ]    | Написати CLI програму    | 2024-01-15 10:32

Всього: 3 завдань (0 виконано, 3 активних)

$ go run todo.go done 1
✅ Завдання #1 позначено як виконане

$ go run todo.go done 2
✅ Завдання #2 позначено як виконане

$ go run todo.go list
=== TODO List ===

ID   | Статус | Завдання                  | Створено
-----|--------|---------------------------|------------------
1    | [✓]    | Вивчити типи в Go        | 2024-01-15 10:30
2    | [✓]    | Зрозуміти zero values    | 2024-01-15 10:31
3    | [ ]    | Написати CLI програму    | 2024-01-15 10:32

Всього: 3 завдань (2 виконано, 1 активних)
```

---

## Рішення

Рішення знаходиться в `solutions/solution_3.go`.

Це найскладніше завдання тижня - спробуйте виконати його самостійно!

