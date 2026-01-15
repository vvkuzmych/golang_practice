package main

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║         Rune & Unicode Examples         ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	// ===== Що таке rune =====
	fmt.Println("\n🔹 Основи rune")
	fmt.Println("─────────────────────────────────────────")

	var r rune = 'П'
	fmt.Printf("Rune: %c\n", r)
	fmt.Printf("Code point: %d\n", r)
	fmt.Printf("Unicode: U+%04X\n", r)
	fmt.Printf("Тип: %T\n", r)

	// ===== len() vs RuneCount =====
	fmt.Println("\n🔹 len() vs RuneCount (ВАЖЛИВО!)")
	fmt.Println("─────────────────────────────────────────")

	text := "Привіт"
	fmt.Printf("String: %s\n", text)
	fmt.Printf("len(): %d байт\n", len(text))
	fmt.Printf("RuneCount: %d символів\n", utf8.RuneCountInString(text))

	// ===== Українська мова =====
	fmt.Println("\n🔹 Українська мова")
	fmt.Println("─────────────────────────────────────────")

	ukrainian := "Слава Україні!"

	fmt.Printf("Текст: %s\n", ukrainian)
	fmt.Printf("Байтів: %d\n", len(ukrainian))
	fmt.Printf("Символів: %d\n", utf8.RuneCountInString(ukrainian))

	fmt.Println("\nІтерація по рунах:")
	for i, r := range ukrainian {
		fmt.Printf("Позиція %2d: %c (U+%04X)\n", i, r, r)
	}

	// ===== Emoji =====
	fmt.Println("\n🔹 Emoji (4 байти)")
	fmt.Println("─────────────────────────────────────────")

	emoji := "👋🎉🚀"
	fmt.Printf("String: %s\n", emoji)
	fmt.Printf("Байтів: %d\n", len(emoji))
	fmt.Printf("Символів: %d\n", utf8.RuneCountInString(emoji))

	for i, r := range emoji {
		fmt.Printf("%d: %c (U+%04X) - %d байт\n",
			i, r, r, utf8.RuneLen(r))
	}

	// ===== Багатомовний текст =====
	fmt.Println("\n🔹 Багатомовний текст")
	fmt.Println("─────────────────────────────────────────")

	multilang := "Hello Привіт 你好 مرحبا"

	fmt.Printf("Текст: %s\n", multilang)
	fmt.Printf("Байтів: %d\n", len(multilang))
	fmt.Printf("Символів: %d\n", utf8.RuneCountInString(multilang))

	// ===== String indexing problem =====
	fmt.Println("\n🔹 Проблема індексації (ПОМИЛКА)")
	fmt.Println("─────────────────────────────────────────")

	str := "Київ"
	fmt.Printf("String: %s\n", str)

	// ❌ Неправильно - отримуємо byte!
	fmt.Printf("❌ str[0] = %c (%d) - це byte!\n", str[0], str[0])

	// ✅ Правильно - конвертуємо в []rune
	runes := []rune(str)
	fmt.Printf("✅ runes[0] = %c (%d) - це rune!\n", runes[0], runes[0])

	// ===== Конверсія string ↔ []rune =====
	fmt.Println("\n🔹 string ↔ []rune")
	fmt.Println("─────────────────────────────────────────")

	original := "Україна"
	runeSlice := []rune(original)

	fmt.Printf("String: %s\n", original)
	fmt.Printf("[]rune: %v\n", runeSlice)
	fmt.Printf("Кількість рун: %d\n", len(runeSlice))

	// Модифікація
	runeSlice[0] = 'У'
	runeSlice[6] = 'а'

	modified := string(runeSlice)
	fmt.Printf("Modified: %s\n", modified)

	// ===== Ітерація: for vs range =====
	fmt.Println("\n🔹 Ітерація: for vs range")
	fmt.Println("─────────────────────────────────────────")

	word := "Go!"

	fmt.Println("❌ Погано (for по індексу):")
	for i := 0; i < len(word); i++ {
		fmt.Printf("  %d: %c (byte)\n", i, word[i])
	}

	fmt.Println("\n✅ Добре (range):")
	for i, r := range word {
		fmt.Printf("  %d: %c (rune)\n", i, r)
	}

	// ===== UTF-8 encoding sizes =====
	fmt.Println("\n🔹 UTF-8: розміри символів")
	fmt.Println("─────────────────────────────────────────")

	chars := []rune{'A', 'П', '中', '🎉'}

	for _, r := range chars {
		size := utf8.RuneLen(r)
		fmt.Printf("%c (U+%04X): %d байт\n", r, r, size)
	}

	// ===== Unicode categories =====
	fmt.Println("\n🔹 Unicode категорії")
	fmt.Println("─────────────────────────────────────────")

	testRunes := []rune{'A', 'п', '5', ' ', '!', '©'}

	for _, r := range testRunes {
		fmt.Printf("'%c': ", r)

		if unicode.IsLetter(r) {
			fmt.Print("Letter ")
		}
		if unicode.IsUpper(r) {
			fmt.Print("Upper ")
		}
		if unicode.IsLower(r) {
			fmt.Print("Lower ")
		}
		if unicode.IsDigit(r) {
			fmt.Print("Digit ")
		}
		if unicode.IsSpace(r) {
			fmt.Print("Space ")
		}
		if unicode.IsSymbol(r) {
			fmt.Print("Symbol ")
		}
		fmt.Println()
	}

	// ===== ToUpper/ToLower =====
	fmt.Println("\n🔹 ToUpper / ToLower")
	fmt.Println("─────────────────────────────────────────")

	text2 := "Привіт, Світ!"

	fmt.Printf("Оригінал: %s\n", text2)
	fmt.Printf("Upper: %s\n", strings.ToUpper(text2))
	fmt.Printf("Lower: %s\n", strings.ToLower(text2))
	fmt.Printf("Title: %s\n", strings.Title(text2))

	// ===== Підрахунок українських літер =====
	fmt.Println("\n🔹 Підрахунок українських літер")
	fmt.Println("─────────────────────────────────────────")

	ukrText := "Слава Україні! Героям слава!"
	ukrCount := countUkrainianLetters(ukrText)
	totalCount := utf8.RuneCountInString(ukrText)

	fmt.Printf("Текст: %s\n", ukrText)
	fmt.Printf("Всього символів: %d\n", totalCount)
	fmt.Printf("Українських літер: %d\n", ukrCount)

	// ===== Reverse string =====
	fmt.Println("\n🔹 Реверс строки (правильно)")
	fmt.Println("─────────────────────────────────────────")

	toReverse := "Привіт 👋"
	reversed := reverseString(toReverse)

	fmt.Printf("Original: %s\n", toReverse)
	fmt.Printf("Reversed: %s\n", reversed)

	// ===== Substring з рунами =====
	fmt.Println("\n🔹 Substring (правильний спосіб)")
	fmt.Println("─────────────────────────────────────────")

	longText := "Україна - це Європа"
	sub := substring(longText, 0, 7)

	fmt.Printf("Original: %s\n", longText)
	fmt.Printf("Substring(0, 7): %s\n", sub)

	// ===== Summary =====
	fmt.Println("\n\n📝 ВИСНОВКИ")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("✅ rune = int32 (Unicode code point)")
	fmt.Println("✅ len() повертає БАЙТИ, не символи")
	fmt.Println("✅ Використовуйте utf8.RuneCountInString()")
	fmt.Println("✅ range по string ітерує по runes")
	fmt.Println("✅ []rune(s) для індексації символів")
	fmt.Println("✅ Українська: 2 байти/літера")
	fmt.Println("✅ Emoji: 4 байти")
	fmt.Println("✅ strings.ToUpper/ToLower працює з Unicode")
}

// ============= Helper Functions =============

func countUkrainianLetters(text string) int {
	count := 0
	for _, r := range text {
		if (r >= 'А' && r <= 'Я') || (r >= 'а' && r <= 'я') ||
			r == 'Є' || r == 'І' || r == 'Ї' || r == 'Ґ' ||
			r == 'є' || r == 'і' || r == 'ї' || r == 'ґ' {
			count++
		}
	}
	return count
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func substring(s string, start, end int) string {
	runes := []rune(s)
	if start < 0 {
		start = 0
	}
	if end > len(runes) {
		end = len(runes)
	}
	return string(runes[start:end])
}
