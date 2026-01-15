package main

import "fmt"

// ============= Quick Demos =============
//
//func main() {
//	fmt.Println("╔════════════════════════════════════════════╗")
//	fmt.Println("║    Week 3: byte, rune & iota - Demo       ║")
//	fmt.Println("╚════════════════════════════════════════════╝")
//
//	// ===== byte demo =====
//	fmt.Println("\n🔹 byte (uint8)")
//	fmt.Println("─────────────────────────────────────────")
//
//	var b byte = 65
//	fmt.Printf("byte: %d, char: %c\n", b, b)
//
//	text := "Hello"
//	bytes := []byte(text)
//	fmt.Printf("String: %s\n", text)
//	fmt.Printf("Bytes: %v\n", bytes)
//
//	// ===== rune demo =====
//	fmt.Println("\n🔹 rune (int32) - Unicode")
//	fmt.Println("─────────────────────────────────────────")
//
//	ukrainian := "Привіт"
//	fmt.Printf("Text: %s\n", ukrainian)
//	fmt.Printf("len(): %d bytes\n", len(ukrainian))
//	fmt.Printf("RuneCount: %d chars\n", utf8.RuneCountInString(ukrainian))
//
//	fmt.Println("\nІтерація:")
//	for i, r := range ukrainian {
//		fmt.Printf("  %d: %c\n", i, r)
//	}
//
//	// ===== iota demo =====
//	fmt.Println("\n🔹 iota - Auto-increment константи")
//	fmt.Println("─────────────────────────────────────────")
//
//	const (
//		Monday    = iota // 0
//		Tuesday          // 1
//		Wednesday        // 2
//	)
//
//	fmt.Printf("Monday: %d\n", Monday)
//	fmt.Printf("Tuesday: %d\n", Tuesday)
//	fmt.Printf("Wednesday: %d\n", Wednesday)
//
//	// Bit flags
//	const (
//		Read    = 1 << iota // 1
//		Write               // 2
//		Execute             // 4
//	)
//
//	fmt.Println("\nBit flags:")
//	fmt.Printf("Read: %d (%03b)\n", Read, Read)
//	fmt.Printf("Write: %d (%03b)\n", Write, Write)
//	fmt.Printf("Execute: %d (%03b)\n", Execute, Execute)
//
//	perms := Read | Write
//	fmt.Printf("Read+Write: %d (%03b)\n", perms, perms)
//
//	// ===== Key Differences =====
//	fmt.Println("\n" + strings.Repeat("═", 44))
//	fmt.Println("📚 Ключові відмінності:")
//	fmt.Println(strings.Repeat("═", 44))
//
//	fmt.Println("\nbyte:")
//	fmt.Println("  • uint8 (0-255)")
//	fmt.Println("  • 1 байт")
//	fmt.Println("  • Для ASCII і binary data")
//
//	fmt.Println("\nrune:")
//	fmt.Println("  • int32 (Unicode code point)")
//	fmt.Println("  • 4 байти")
//	fmt.Println("  • Для Unicode символів")
//
//	fmt.Println("\niota:")
//	fmt.Println("  • Auto-increment константа")
//	fmt.Println("  • Для enum patterns")
//	fmt.Println("  • Reset в кожному const блоці")
//
//	// ===== Practical Examples =====
//	fmt.Println("\n" + strings.Repeat("═", 44))
//	fmt.Println("💡 Практичні приклади:")
//	fmt.Println(strings.Repeat("═", 44))
//
//	// byte: Simple cipher
//	fmt.Println("\n1️⃣  byte: Simple XOR cipher")
//	message := []byte("Secret")
//	key := byte(42)
//	encrypted := xorBytes(message, key)
//	fmt.Printf("Original: %s\n", string(message))
//	fmt.Printf("Encrypted: %v\n", encrypted)
//	fmt.Printf("Decrypted: %s\n", string(xorBytes(encrypted, key)))
//
//	// rune: Ukrainian text processing
//	fmt.Println("\n2️⃣  rune: Підрахунок українських літер")
//	ukrText := "Слава Україні!"
//	count := countUkrainianLetters(ukrText)
//	fmt.Printf("Text: %s\n", ukrText)
//	fmt.Printf("Українських літер: %d\n", count)
//
//	// iota: Status system
//	fmt.Println("\n3️⃣  iota: Status system")
//	type Status int
//	const (
//		Pending Status = iota
//		Active
//		Completed
//	)
//
//	status := Active
//	fmt.Printf("Current status: %d\n", status)
//	if status == Active {
//		fmt.Println("Status is Active")
//	}
//
//	// ===== Instructions =====
//	fmt.Println("\n" + strings.Repeat("═", 44))
//	fmt.Println("🚀 Навчальні матеріали:")
//	fmt.Println(strings.Repeat("═", 44))
//
//	fmt.Println("\n1️⃣  Теорія:")
//	fmt.Println("   cd theory")
//	fmt.Println("   cat 01_byte_basics.md")
//	fmt.Println("   cat 02_rune_unicode.md")
//	fmt.Println("   cat 03_iota_enums.md")
//
//	fmt.Println("\n2️⃣  Практика:")
//	fmt.Println("   cd practice/byte_examples && go run main.go")
//	fmt.Println("   cd practice/rune_examples && go run main.go")
//	fmt.Println("   cd practice/iota_examples && go run main.go")
//
//	fmt.Println("\n3️⃣  Вправи:")
//	fmt.Println("   cd exercises")
//	fmt.Println("   cat exercise_1.md  # Byte Encoder")
//	fmt.Println("   cat exercise_2.md  # Unicode Counter")
//	fmt.Println("   cat exercise_3.md  # Status System")
//
//	fmt.Println("\n4️⃣  Рішення:")
//	fmt.Println("   cd solutions")
//	fmt.Println("   go run solution_1.go")
//	fmt.Println("   go run solution_2.go")
//	fmt.Println("   go run solution_3.go")
//
//	fmt.Println("\n" + strings.Repeat("═", 44))
//	fmt.Println("📖 Почніть з:")
//	fmt.Println("   cat README.md")
//	fmt.Println("   cat QUICK_START.md")
//	fmt.Println(strings.Repeat("═", 44))
//
//	fmt.Println("\nУдачі у вивченні! 🎉")
//}
//
//// ============= Helper Functions =============
//
//func xorBytes(data []byte, key byte) []byte {
//	result := make([]byte, len(data))
//	for i, b := range data {
//		result[i] = b ^ key
//	}
//	return result
//}
//
//func countUkrainianLetters(text string) int {
//	count := 0
//	for _, r := range text {
//		if (r >= 'А' && r <= 'Я') || (r >= 'а' && r <= 'я') ||
//			r == 'Є' || r == 'І' || r == 'Ї' || r == 'Ґ' ||
//			r == 'є' || r == 'і' || r == 'ї' || r == 'ґ' {
//			count++
//		}
//	}
//	return count
//}

func main() {
	s := [][]int{{0, 6, 0},
		{5, 8, 7},
		{0, 9, 0}}

	g := getMaximumGold(s)
	fmt.Println("--------", g)
}

func getMaximumGold(grid [][]int) int {
	new_grid := []int{}

	for _, num := range grid {
		for _, other := range num {
			new_grid = append(new_grid, other)
		}
	}

	s := bubbleSort(new_grid)

	sum := sumLastThree(s)

	return sum
}

func bubbleSort(nums []int) []int {
	n := len(nums)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if nums[j] > nums[j+1] {
				nums[j], nums[j+1] = nums[j+1], nums[j]
			}
		}
	}

	return nums
}

func sumLastThree(nums []int) int {
	if len(nums) < 3 {
		return 0 // or handle error
	}

	sum := 0
	for _, v := range nums[len(nums)-3:] {
		sum += v
	}
	return sum
}
