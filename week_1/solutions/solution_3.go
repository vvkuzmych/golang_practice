package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const programName1 = "greet"

// Task представляє одне завдання
type Task struct {
	ID          int
	Description string
	Done        bool
	CreatedAt   string
}

// TodoList представляє список завдань
type TodoList struct {
	Tasks  []Task
	NextID int
}

func main() {
	// Ініціалізація списку завдань
	var todoList TodoList
	todoList.NextID = 1

	// Перевірка аргументів
	if len(os.Args) < 2 {
		printError("Команда не вказана")
		printHelp()
		os.Exit(1)
	}

	// Отримання команди
	command := strings.ToLower(os.Args[1])

	// Виконання команди
	switch command {
	case "add":
		if len(os.Args) < 3 {
			printError("Не вказано текст завдання")
			fmt.Println("Використання: todo add <текст завдання>")
			os.Exit(1)
		}
		description := strings.Join(os.Args[2:], " ")
		addTask(&todoList, description)

	case "list", "ls":
		listTasks(&todoList)

	case "done", "complete":
		if len(os.Args) < 3 {
			printError("Не вказано ID завдання")
			fmt.Println("Використання: todo done <id>")
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			printError(fmt.Sprintf("'%s' не є коректним ID", os.Args[2]))
			os.Exit(1)
		}
		markDone(&todoList, id)

	case "delete", "del", "rm":
		if len(os.Args) < 3 {
			printError("Не вказано ID завдання")
			fmt.Println("Використання: todo delete <id>")
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			printError(fmt.Sprintf("'%s' не є коректним ID", os.Args[2]))
			os.Exit(1)
		}
		deleteTask(&todoList, id)

	case "help", "-h", "--help":
		printHelp()

	default:
		printError(fmt.Sprintf("Невідома команда: '%s'", command))
		printHelp()
		os.Exit(1)
	}
}

// addTask додає нове завдання до списку
func addTask(list *TodoList, description string) {
	// Створення нового завдання
	task := Task{
		ID:          list.NextID,
		Description: description,
		Done:        false,
		CreatedAt:   time.Now().Format("2006-01-02 15:04"),
	}

	// Додавання до списку
	list.Tasks = append(list.Tasks, task)
	list.NextID++

	// Вивід підтвердження
	fmt.Printf("✅ Завдання додано: %q (ID: %d)\n", description, task.ID)
}

// listTasks виводить всі завдання
func listTasks(list *TodoList) {
	if len(list.Tasks) == 0 {
		fmt.Println("📝 Список завдань порожній")
		fmt.Println("\nДодайте перше завдання:")
		fmt.Println("  todo add \"Моє перше завдання\"")
		return
	}

	fmt.Println("\n=== TODO List ===\n")

	// Заголовок таблиці
	fmt.Printf("%-4s | %-6s | %-30s | %s\n", "ID", "Статус", "Завдання", "Створено")
	fmt.Println("-----|--------|--------------------------------|------------------")

	// Лічильники
	var completed, active int

	// Вивід завдань
	for _, task := range list.Tasks {
		status := "[ ]"
		if task.Done {
			status = "[✓]"
			completed++
		} else {
			active++
		}

		// Обрізати опис якщо занадто довгий
		description := task.Description
		if len(description) > 30 {
			description = description[:27] + "..."
		}

		fmt.Printf("%-4d | %-6s | %-30s | %s\n",
			task.ID, status, description, task.CreatedAt)
	}

	// Підсумок
	fmt.Printf("\nВсього: %d завдань (%d виконано, %d активних)\n",
		len(list.Tasks), completed, active)
}

// markDone позначає завдання як виконане
func markDone(list *TodoList, id int) {
	// Пошук завдання
	for i := range list.Tasks {
		if list.Tasks[i].ID == id {
			if list.Tasks[i].Done {
				fmt.Printf("⚠️  Завдання #%d вже виконано\n", id)
				return
			}
			list.Tasks[i].Done = true
			fmt.Printf("✅ Завдання #%d позначено як виконане\n", id)
			return
		}
	}

	// Якщо не знайдено
	printError(fmt.Sprintf("Завдання з ID %d не знайдено", id))
	os.Exit(1)
}

// deleteTask видаляє завдання зі списку
func deleteTask(list *TodoList, id int) {
	// Пошук завдання
	for i, task := range list.Tasks {
		if task.ID == id {
			// Видалення елемента зі slice
			list.Tasks = append(list.Tasks[:i], list.Tasks[i+1:]...)
			fmt.Printf("✅ Завдання #%d видалено\n", id)
			return
		}
	}

	// Якщо не знайдено
	printError(fmt.Sprintf("Завдання з ID %d не знайдено", id))
	os.Exit(1)
}

// printHelp виводить довідку
func printHelp() {
	fmt.Println("\n=== TODO Manager - Довідка ===\n")
	fmt.Printf("Використання: %s <команда> [аргументи]\n\n", programName1)

	fmt.Println("Доступні команди:")
	fmt.Println("  add <text>       - Додати нове завдання")
	fmt.Println("  list             - Показати всі завдання (або: ls)")
	fmt.Println("  done <id>        - Позначити завдання як виконане (або: complete)")
	fmt.Println("  delete <id>      - Видалити завдання (або: del, rm)")
	fmt.Println("  help             - Показати цю довідку (або: -h, --help)")

	fmt.Println("\nПриклади:")
	fmt.Printf("  %s add \"Вивчити Go\"\n", programName1)
	fmt.Printf("  %s add \"Написати тести\"\n", programName1)
	fmt.Printf("  %s list\n", programName1)
	fmt.Printf("  %s done 1\n", programName1)
	fmt.Printf("  %s delete 2\n", programName1)

	fmt.Println("\nПідказки:")
	fmt.Println("  💡 Використовуйте лапки для завдань з пробілами")
	fmt.Println("  💡 ID завдань можна побачити командою 'list'")
	fmt.Println("  💡 Завдання з [✓] вже виконані")
}

// printError виводить повідомлення про помилку
func printError(message string) {
	fmt.Printf("❌ Помилка: %s\n\n", message)
}

/*
ДЕМОНСТРАЦІЯ КОНЦЕПЦІЙ ТИЖНЯ 1:

1. Типи даних:
   - int (ID, лічильники)
   - string (Description, CreatedAt, команди)
   - bool (Done)
   - struct (Task, TodoList)
   - slice ([]Task)

2. Zero Values:
   - var todoList TodoList (struct з zero values)
   - var completed, active int (int = 0)
   - Task.Done = false (bool = false)

3. var vs :=
   - var todoList TodoList (zero value)
   - var completed, active int (zero values)
   - command := strings.ToLower(...) (коротке оголошення)
   - task := Task{...} (коротке оголошення)

4. Пакети та функції:
   - package main
   - func main()
   - Імпорти: fmt, os, strconv, strings, time
   - Власні функції: addTask, listTasks, markDone, deleteTask

5. fmt.Printf формати:
   - %s - string
   - %d - int
   - %q - quoted string
   - %v - default format
   - %T - type
   - %+v - struct з полями
   - Ширина: %-4s, %-30s

6. Робота з командним рядком:
   - os.Args
   - len(os.Args)
   - os.Exit()

7. Структури:
   - type Task struct
   - type TodoList struct
   - Створення екземплярів
   - Робота з полями

8. Slice операції:
   - append()
   - Видалення: append(s[:i], s[i+1:]...)
   - Ітерація: for range

9. Pointer:
   - func addTask(list *TodoList, ...)
   - Передача по посиланню для модифікації
*/
