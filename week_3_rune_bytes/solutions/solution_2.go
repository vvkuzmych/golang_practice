package main

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// ============= Character Counting =============

// CharCount повертає кількість символів (не байтів!)
func CharCount(text string) int {
	return utf8.RuneCountInString(text)
}

// UkrainianCount підраховує українські літери
func UkrainianCount(text string) int {
	count := 0
	for _, r := range text {
		if isUkrainian(r) {
			count++
		}
	}
	return count
}

// isUkrainian перевіряє чи символ українська літера
func isUkrainian(r rune) bool {
	return (r >= 'А' && r <= 'Я') || (r >= 'а' && r <= 'я') ||
		r == 'Є' || r == 'І' || r == 'Ї' || r == 'Ґ' ||
		r == 'є' || r == 'і' || r == 'ї' || r == 'ґ'
}

// EmojiCount підраховує emoji
func EmojiCount(text string) int {
	count := 0
	for _, r := range text {
		if isEmoji(r) {
			count++
		}
	}
	return count
}

// isEmoji перевіряє чи символ emoji
func isEmoji(r rune) bool {
	return (r >= 0x1F300 && r <= 0x1F9FF) || // основні emoji
		(r >= 0x2600 && r <= 0x26FF) || // різні символи
		(r >= 0x2700 && r <= 0x27BF) // dingbats
}

// ============= String Manipulation =============

// Reverse реверсує string (правильно для Unicode)
func Reverse(text string) string {
	runes := []rune(text)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Substring витягує підстроку за індексами символів (не байтів!)
func Substring(text string, start, end int) string {
	runes := []rune(text)
	if start < 0 {
		start = 0
	}
	if end > len(runes) {
		end = len(runes)
	}
	if start >= end {
		return ""
	}
	return string(runes[start:end])
}

// ============= Statistics =============

// TextStats повертає детальну статистику
func TextStats(text string) map[string]int {
	stats := map[string]int{
		"chars":       CharCount(text),
		"bytes":       len(text),
		"ukrainian":   0,
		"latin":       0,
		"digits":      0,
		"emoji":       0,
		"spaces":      0,
		"punctuation": 0,
		"unique":      0,
	}

	uniqueRunes := make(map[rune]bool)

	for _, r := range text {
		uniqueRunes[r] = true

		if isUkrainian(r) {
			stats["ukrainian"]++
		} else if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
			stats["latin"]++
		} else if unicode.IsDigit(r) {
			stats["digits"]++
		} else if isEmoji(r) {
			stats["emoji"]++
		} else if unicode.IsSpace(r) {
			stats["spaces"]++
		} else if unicode.IsPunct(r) {
			stats["punctuation"]++
		}
	}

	stats["unique"] = len(uniqueRunes)

	return stats
}

// UTF8Distribution показує розподіл за розмірами UTF-8
func UTF8Distribution(text string) map[int]int {
	dist := make(map[int]int)

	for _, r := range text {
		size := utf8.RuneLen(r)
		dist[size]++
	}

	return dist
}

// MostFrequentRune знаходить найчастішу руну
func MostFrequentRune(text string) (rune, int) {
	freq := make(map[rune]int)

	for _, r := range text {
		freq[r]++
	}

	var mostRune rune
	maxCount := 0

	for r, count := range freq {
		if count > maxCount {
			mostRune = r
			maxCount = count
		}
	}

	return mostRune, maxCount
}

// WordCount підраховує слова
func WordCount(text string) int {
	return len(strings.Fields(text))
}

// ============= Helper Functions =============

func printSeparator() {
	fmt.Println(strings.Repeat("─", 50))
}

func printHeader(title string) {
	fmt.Printf("\n🔹 %s\n", title)
	printSeparator()
}

// ============= Main =============

func main() {
	fmt.Println("╔════════════════════════════════════════════════╗")
	fmt.Println("║      Unicode Text Analyzer Solution            ║")
	fmt.Println("╚════════════════════════════════════════════════╝")

	// ===== Test Cases =====
	testCases := []string{
		"Привіт",
		"Hello, World!",
		"Привіт, World! 👋🎉",
		"你好世界",
		"مرحبا",
		"Слава Україні! 🇺🇦",
	}

	for _, text := range testCases {
		printHeader(fmt.Sprintf("Analyzing: %s", text))

		chars := CharCount(text)
		bytes := len(text)
		ratio := float64(bytes) / float64(chars)

		fmt.Printf("Символів: %d\n", chars)
		fmt.Printf("Байтів: %d\n", bytes)
		fmt.Printf("Співвідношення: %.2f байт/символ\n", ratio)

		// Detailed stats
		stats := TextStats(text)
		if stats["ukrainian"] > 0 {
			fmt.Printf("  Українських: %d\n", stats["ukrainian"])
		}
		if stats["latin"] > 0 {
			fmt.Printf("  Латинських: %d\n", stats["latin"])
		}
		if stats["emoji"] > 0 {
			fmt.Printf("  Emoji: %d\n", stats["emoji"])
		}
		if stats["digits"] > 0 {
			fmt.Printf("  Цифр: %d\n", stats["digits"])
		}
		if stats["spaces"] > 0 {
			fmt.Printf("  Пробілів: %d\n", stats["spaces"])
		}

		// UTF-8 distribution
		dist := UTF8Distribution(text)
		fmt.Printf("UTF-8 розподіл: ")
		for size := 1; size <= 4; size++ {
			if count, ok := dist[size]; ok {
				fmt.Printf("%db:%d ", size, count)
			}
		}
		fmt.Println()
	}

	// ===== Detailed Example =====
	printHeader("Детальний приклад")

	text := "Привіт, World! 👋🎉"
	fmt.Printf("Text: %s\n", text)
	fmt.Printf("Chars: %d, Bytes: %d\n", CharCount(text), len(text))

	// Character by character
	fmt.Println("\nПосимвольно:")
	for i, r := range text {
		size := utf8.RuneLen(r)
		fmt.Printf("  [%d] '%c' (U+%04X) - %d bytes\n", i, r, r, size)
	}

	// ===== Operations =====
	printHeader("Операції")

	fmt.Printf("Original: %s\n", text)

	reversed := Reverse(text)
	fmt.Printf("Reversed: %s\n", reversed)

	upper := strings.ToUpper(text)
	fmt.Printf("Uppercase: %s\n", upper)

	lower := strings.ToLower(text)
	fmt.Printf("Lowercase: %s\n", lower)

	// Substring
	sub := Substring(text, 0, 6)
	fmt.Printf("Substring(0,6): %s\n", sub)

	// ===== Statistics =====
	printHeader("Статистика")

	stats := TextStats(text)
	fmt.Printf("Всього символів: %d\n", stats["chars"])
	fmt.Printf("Унікальних символів: %d\n", stats["unique"])
	fmt.Printf("Українських літер: %d\n", stats["ukrainian"])
	fmt.Printf("Латинських літер: %d\n", stats["latin"])
	fmt.Printf("Emoji: %d\n", stats["emoji"])
	fmt.Printf("Пробілів: %d\n", stats["spaces"])
	fmt.Printf("Розділових знаків: %d\n", stats["punctuation"])

	mostRune, mostCount := MostFrequentRune(text)
	fmt.Printf("Найчастіший символ: '%c' (%d разів)\n", mostRune, mostCount)

	// ===== Ukrainian Text =====
	printHeader("Українські тексти")

	ukrainian := []string{
		"Слава Україні!",
		"Героям слава!",
		"Київ - столиця України",
		"Я люблю Україну 🇺🇦",
	}

	for _, ua := range ukrainian {
		ukrCount := UkrainianCount(ua)
		totalChars := CharCount(ua)
		fmt.Printf("%s\n", ua)
		fmt.Printf("  Українських літер: %d/%d\n", ukrCount, totalChars)
	}

	// ===== Word Count =====
	printHeader("Підрахунок слів")

	sentences := []string{
		"Привіт світ",
		"Hello World",
		"Це речення має п'ять слів",
		"This sentence has five words",
	}

	for _, sentence := range sentences {
		words := WordCount(sentence)
		fmt.Printf("%s → %d слів\n", sentence, words)
	}

	// ===== Common Mistakes =====
	printHeader("Поширені помилки")

	ukrText := "Привіт"

	fmt.Println("❌ Неправильно:")
	fmt.Printf("  len(\"%s\") = %d  ← байти, не символи!\n", ukrText, len(ukrText))
	fmt.Printf("  \"%s\"[0] = %d ← байт, не літера!\n", ukrText, ukrText[0])

	fmt.Println("\n✅ Правильно:")
	fmt.Printf("  utf8.RuneCountInString(\"%s\") = %d\n", ukrText, utf8.RuneCountInString(ukrText))
	runes := []rune(ukrText)
	fmt.Printf("  []rune(\"%s\")[0] = '%c'\n", ukrText, runes[0])

	// ===== Emoji Examples =====
	printHeader("Emoji приклади")

	emojis := []string{
		"👋",
		"🎉",
		"🚀",
		"🇺🇦",
		"👨‍👩‍👧‍👦", // складні emoji з ZWJ
	}

	for _, emoji := range emojis {
		chars := CharCount(emoji)
		bytes := len(emoji)
		emojiCount := EmojiCount(emoji)
		fmt.Printf("%s: chars=%d, bytes=%d, emoji=%d\n", emoji, chars, bytes, emojiCount)
	}

	// ===== Summary =====
	fmt.Println("\n\n📝 ВИСНОВКИ")
	printSeparator()
	fmt.Println("✅ len() повертає БАЙТИ, не символи")
	fmt.Println("✅ utf8.RuneCountInString() для підрахунку символів")
	fmt.Println("✅ []rune(text) для індексації по символах")
	fmt.Println("✅ range ітерує по рунах (символах)")
	fmt.Println("✅ Українські літери - 2 байти (UTF-8)")
	fmt.Println("✅ Emoji - 4 байти (іноді більше)")
	fmt.Println("✅ ASCII - 1 байт")
}
