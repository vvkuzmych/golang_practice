package main

import (
	"errors"
	"fmt"
)

// ============= Sentinel Errors =============

var (
	ErrDatabase = errors.New("database error")
	ErrNetwork  = errors.New("network error")
)

// ============= Examples =============

func example1_WithoutWrapping() {
	fmt.Println("1️⃣ БЕЗ Wrapping (%v) - втрачаємо оригінал")
	fmt.Println("─────────────────────────────────────────")

	err := operationWithoutWrapping()
	fmt.Printf("Error: %v\n", err)

	// ❌ Не можемо перевірити оригінальну помилку!
	if errors.Is(err, ErrDatabase) {
		fmt.Println("Database error detected")
	} else {
		fmt.Println("❌ Can't detect ErrDatabase (wrapping lost!)")
	}
	fmt.Println()
}

func operationWithoutWrapping() error {
	err := ErrDatabase
	// %v - НЕ зберігає оригінал
	return fmt.Errorf("operation failed: %v", err)
}

func example2_WithWrapping() {
	fmt.Println("2️⃣ З Wrapping (%w) - зберігаємо оригінал")
	fmt.Println("─────────────────────────────────────────")

	err := operationWithWrapping()
	fmt.Printf("Error: %v\n", err)

	// ✅ Можемо перевірити оригінальну помилку!
	if errors.Is(err, ErrDatabase) {
		fmt.Println("✓ Database error detected (wrapping works!)")
	}
	fmt.Println()
}

func operationWithWrapping() error {
	err := ErrDatabase
	// %w - ЗБЕРІГАЄ оригінал
	return fmt.Errorf("operation failed: %w", err)
}

func example3_ErrorChain() {
	fmt.Println("3️⃣ Ланцюжок помилок")
	fmt.Println("─────────────────────────────────────────")

	err := level3()
	fmt.Printf("Error: %v\n", err)

	// Проходимо ланцюжок
	fmt.Println("\nError chain:")
	current := err
	level := 1
	for current != nil {
		fmt.Printf("  Level %d: %v\n", level, current)
		current = errors.Unwrap(current)
		level++
	}
	fmt.Println()
}

func level1() error {
	return ErrDatabase
}

func level2() error {
	err := level1()
	if err != nil {
		return fmt.Errorf("level2: %w", err)
	}
	return nil
}

func level3() error {
	err := level2()
	if err != nil {
		return fmt.Errorf("level3: %w", err)
	}
	return nil
}

func example4_Unwrap() {
	fmt.Println("4️⃣ errors.Unwrap()")
	fmt.Println("─────────────────────────────────────────")

	err1 := errors.New("original error")
	err2 := fmt.Errorf("wrapped: %w", err1)
	err3 := fmt.Errorf("twice wrapped: %w", err2)

	fmt.Printf("err3: %v\n", err3)

	// Unwrap один раз
	unwrapped := errors.Unwrap(err3)
	fmt.Printf("After 1 unwrap: %v\n", unwrapped)

	// Unwrap ще раз
	unwrapped = errors.Unwrap(unwrapped)
	fmt.Printf("After 2 unwraps: %v\n", unwrapped)

	// Unwrap останній раз
	unwrapped = errors.Unwrap(unwrapped)
	fmt.Printf("After 3 unwraps: %v\n\n", unwrapped) // nil
}

func example5_AddingContext() {
	fmt.Println("5️⃣ Додавання контексту")
	fmt.Println("─────────────────────────────────────────")

	userID := 42
	filename := "config.json"

	err := processFile(userID, filename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		// Output: "failed to process file for user 42: failed to open config.json: file not found"
	}
	fmt.Println()
}

func processFile(userID int, filename string) error {
	err := openFile(filename)
	if err != nil {
		// Додаємо контекст: хто викликав
		return fmt.Errorf("failed to process file for user %d: %w", userID, err)
	}
	return nil
}

func openFile(filename string) error {
	// Симуляція помилки відкриття файлу
	err := errors.New("file not found")
	// Додаємо контекст: який файл
	return fmt.Errorf("failed to open %s: %w", filename, err)
}

func example6_Comparison() {
	fmt.Println("6️⃣ %v vs %w - Порівняння")
	fmt.Println("─────────────────────────────────────────")

	original := errors.New("original error")

	// З %v
	errV := fmt.Errorf("wrapped with %%v: %v", original)
	fmt.Printf("With %%v: errors.Is() = %v ❌\n", errors.Is(errV, original))

	// З %w
	errW := fmt.Errorf("wrapped with %%w: %w", original)
	fmt.Printf("With %%w: errors.Is() = %v ✅\n\n", errors.Is(errW, original))
}

func example7_RealWorld() {
	fmt.Println("7️⃣ Real-World приклад")
	fmt.Println("─────────────────────────────────────────")

	err := saveUser(User{ID: 123, Name: "Alice"})
	if err != nil {
		fmt.Printf("❌ Failed to save user:\n   %v\n", err)

		// Специфічна обробка для різних типів помилок
		if errors.Is(err, ErrDatabase) {
			fmt.Println("   → Database issue - retry later")
		} else if errors.Is(err, ErrNetwork) {
			fmt.Println("   → Network issue - check connection")
		}
	}
	fmt.Println()
}

type User struct {
	ID   int
	Name string
}

func saveUser(user User) error {
	err := connectDB()
	if err != nil {
		return fmt.Errorf("failed to save user %d: %w", user.ID, err)
	}
	return nil
}

func connectDB() error {
	// Симуляція DB помилки
	return fmt.Errorf("connection failed: %w", ErrDatabase)
}

// ============= Main =============

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║       Error Wrapping Examples            ║")
	fmt.Println("╚══════════════════════════════════════════╝")
	fmt.Println()

	example1_WithoutWrapping()
	example2_WithWrapping()
	example3_ErrorChain()
	example4_Unwrap()
	example5_AddingContext()
	example6_Comparison()
	example7_RealWorld()

	fmt.Println("📝 Висновки:")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("✅ Використовуйте %w для wrapping")
	fmt.Println("✅ %v НЕ зберігає оригінальну помилку")
	fmt.Println("✅ Додавайте контекст на кожному рівні")
	fmt.Println("✅ errors.Is() працює через wrapping")
	fmt.Println("✅ errors.Unwrap() розгортає ланцюжок")
	fmt.Println("✅ Кожен layer додає свій контекст")
}
