package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

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

// Helper функції для читання
func readString(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func readInt(reader *bufio.Reader, prompt string) (int, error) {
	input := readString(reader, prompt)
	return strconv.Atoi(input)
}

// Функції для роботи зі списком завдань
func addTask(list *TodoList, reader *bufio.Reader) {
	description := readString(reader, "\nОпис завдання: ")

	if description == "" {
		fmt.Println("❌ Опис не може бути порожнім")
		return
	}

	task := Task{
		ID:          list.NextID,
		Description: description,
		Done:        false,
		CreatedAt:   time.Now().Format("2006-01-02 15:04"),
	}

	list.Tasks = append(list.Tasks, task)
	list.NextID++

	fmt.Printf("✅ Завдання додано з ID: %d\n", task.ID)
}

func listTasks(list *TodoList) {
	if len(list.Tasks) == 0 {
		fmt.Println("\n📝 Список завдань порожній")
		return
	}

	fmt.Println("\n╔══════════════════════════════════════════════════════╗")
	fmt.Println("║                   TODO List                          ║")
	fmt.Println("╚══════════════════════════════════════════════════════╝\n")

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
	fmt.Printf("\n📊 Всього: %d | ✅ Виконано: %d | ⏳ Активних: %d\n",
		len(list.Tasks), completed, active)
}

func markDone(list *TodoList, reader *bufio.Reader) {
	listTasks(list)

	if len(list.Tasks) == 0 {
		return
	}

	id, err := readInt(reader, "\nВведіть ID завдання для позначки: ")
	if err != nil {
		fmt.Println("❌ Некоректний ID")
		return
	}

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

	fmt.Printf("❌ Завдання з ID %d не знайдено\n", id)
}

func deleteTask(list *TodoList, reader *bufio.Reader) {
	listTasks(list)

	if len(list.Tasks) == 0 {
		return
	}

	id, err := readInt(reader, "\nВведіть ID завдання для видалення: ")
	if err != nil {
		fmt.Println("❌ Некоректний ID")
		return
	}

	for i, task := range list.Tasks {
		if task.ID == id {
			list.Tasks = append(list.Tasks[:i], list.Tasks[i+1:]...)
			fmt.Printf("✅ Завдання #%d видалено\n", id)
			return
		}
	}

	fmt.Printf("❌ Завдання з ID %d не знайдено\n", id)
}

func printMenu() {
	fmt.Println("\n╔════════════════════════════════════╗")
	fmt.Println("║      TODO Manager - Меню           ║")
	fmt.Println("╠════════════════════════════════════╣")
	fmt.Println("║  1 - Додати завдання               ║")
	fmt.Println("║  2 - Показати всі завдання         ║")
	fmt.Println("║  3 - Позначити виконаним           ║")
	fmt.Println("║  4 - Видалити завдання             ║")
	fmt.Println("║  0 - Вихід                         ║")
	fmt.Println("╚════════════════════════════════════╝")
}

func main() {
	todoList := TodoList{
		Tasks:  []Task{},
		NextID: 1,
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n╔════════════════════════════════════════════════════╗")
	fmt.Println("║     Інтерактивний TODO Manager                     ║")
	fmt.Println("║     Ласкаво просимо!                               ║")
	fmt.Println("╚════════════════════════════════════════════════════╝")

	for {
		printMenu()

		choice, err := readInt(reader, "\nВаш вибір: ")

		if err != nil {
			fmt.Println("❌ Невірний вибір! Введіть число 0-4")
			continue
		}

		switch choice {
		case 1:
			addTask(&todoList, reader)
		case 2:
			listTasks(&todoList)
		case 3:
			markDone(&todoList, reader)
		case 4:
			deleteTask(&todoList, reader)
		case 0:
			fmt.Println("\n👋 До побачення!")
			fmt.Printf("📊 Ви виконали завдань: %d\n", countCompleted(&todoList))
			return
		default:
			fmt.Println("❌ Невірний вибір! Виберіть 0-4")
		}

		// Пауза перед наступною ітерацією
		fmt.Print("\nНатисніть Enter щоб продовжити...")
		reader.ReadString('\n')
	}
}

func countCompleted(list *TodoList) int {
	count := 0
	for _, task := range list.Tasks {
		if task.Done {
			count++
		}
	}
	return count
}
