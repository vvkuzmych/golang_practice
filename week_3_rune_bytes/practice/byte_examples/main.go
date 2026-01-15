package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║         Byte Examples                    ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	// ===== Basic byte =====
	fmt.Println("\n🔹 Основи byte")
	fmt.Println("─────────────────────────────────────────")

	var b byte = 65
	fmt.Printf("Byte значення: %d\n", b)
	fmt.Printf("Як символ: %c\n", b)
	fmt.Printf("Бінарне: %08b\n", b)
	fmt.Printf("Hex: %02x\n", b)

	// ===== String to []byte =====
	fmt.Println("\n🔹 String → []byte")
	fmt.Println("─────────────────────────────────────────")

	text := "Hello"
	byteSlice := []byte(text)

	fmt.Printf("String: %s\n", text)
	fmt.Printf("Bytes: %v\n", byteSlice)
	fmt.Printf("Довжина: %d\n", len(byteSlice))

	// ===== Модифікація bytes =====
	fmt.Println("\n🔹 Модифікація байтів")
	fmt.Println("─────────────────────────────────────────")

	data := []byte("Hello")
	fmt.Printf("Оригінал: %s\n", string(data))

	data[0] = 'h' // H -> h
	fmt.Printf("Після зміни: %s\n", string(data))

	// Uppercase -> Lowercase
	for i, b := range data {
		if b >= 'A' && b <= 'Z' {
			data[i] = b + 32
		}
	}
	fmt.Printf("Lowercase: %s\n", string(data))

	// ===== ASCII операції =====
	fmt.Println("\n🔹 ASCII перевірки")
	fmt.Println("─────────────────────────────────────────")

	testBytes := []byte{'A', 'z', '5', '@', ' '}

	for _, b := range testBytes {
		fmt.Printf("'%c' (%d): ", b, b)

		if b >= 'A' && b <= 'Z' {
			fmt.Print("UPPERCASE ")
		}
		if b >= 'a' && b <= 'z' {
			fmt.Print("lowercase ")
		}
		if b >= '0' && b <= '9' {
			fmt.Print("digit ")
		}
		if b == ' ' {
			fmt.Print("space ")
		}
		fmt.Println()
	}

	// ===== Hex Encoding =====
	fmt.Println("\n🔹 Hex кодування")
	fmt.Println("─────────────────────────────────────────")

	message := []byte("Go")
	hexStr := hex.EncodeToString(message)
	fmt.Printf("Original: %s\n", string(message))
	fmt.Printf("Hex: %s\n", hexStr)

	decoded, _ := hex.DecodeString(hexStr)
	fmt.Printf("Decoded: %s\n", string(decoded))

	// ===== Base64 Encoding =====
	fmt.Println("\n🔹 Base64 кодування")
	fmt.Println("─────────────────────────────────────────")

	secret := []byte("Secret Message")
	encoded := base64.StdEncoding.EncodeToString(secret)
	fmt.Printf("Original: %s\n", string(secret))
	fmt.Printf("Base64: %s\n", encoded)

	decodedB64, _ := base64.StdEncoding.DecodeString(encoded)
	fmt.Printf("Decoded: %s\n", string(decodedB64))

	// ===== bytes.Buffer =====
	fmt.Println("\n🔹 bytes.Buffer")
	fmt.Println("─────────────────────────────────────────")

	var buf bytes.Buffer
	buf.WriteString("Hello")
	buf.WriteByte(' ')
	buf.Write([]byte("World"))
	buf.WriteByte('!')

	fmt.Printf("Buffer: %s\n", buf.String())
	fmt.Printf("Length: %d bytes\n", buf.Len())

	// ===== bytes операції =====
	fmt.Println("\n🔹 bytes операції")
	fmt.Println("─────────────────────────────────────────")

	b1 := []byte("Hello")
	b2 := []byte("World")
	b3 := []byte("Hello")

	fmt.Printf("bytes.Equal(b1, b2): %t\n", bytes.Equal(b1, b2))
	fmt.Printf("bytes.Equal(b1, b3): %t\n", bytes.Equal(b1, b3))

	joined := bytes.Join([][]byte{b1, b2}, []byte(", "))
	fmt.Printf("Joined: %s\n", string(joined))

	fmt.Printf("Contains 'ell': %t\n", bytes.Contains(b1, []byte("ell")))

	replaced := bytes.Replace(b1, []byte("l"), []byte("L"), -1)
	fmt.Printf("Replaced: %s\n", string(replaced))

	// ===== XOR шифрування =====
	fmt.Println("\n🔹 XOR шифрування")
	fmt.Println("─────────────────────────────────────────")

	plaintext := []byte("Secret")
	key := byte(42)

	// Encrypt
	encrypted := xorCipher(plaintext, key)
	fmt.Printf("Plaintext: %s\n", string(plaintext))
	fmt.Printf("Encrypted: %v\n", encrypted)

	// Decrypt
	decrypted := xorCipher(encrypted, key)
	fmt.Printf("Decrypted: %s\n", string(decrypted))

	// ===== Підрахунок частоти =====
	fmt.Println("\n🔹 Підрахунок частоти байтів")
	fmt.Println("─────────────────────────────────────────")

	sample := []byte("Hello, World!")
	freq := countByteFrequency(sample)

	fmt.Printf("Текст: %s\n", string(sample))
	fmt.Println("Частота:")
	for b, count := range freq {
		if b >= 32 && b <= 126 { // printable
			fmt.Printf("  '%c': %d\n", b, count)
		}
	}

	// ===== Checksum =====
	fmt.Println("\n🔹 Simple Checksum")
	fmt.Println("─────────────────────────────────────────")

	data1 := []byte("Hello")
	data2 := []byte("World")

	fmt.Printf("Data1: %s, Checksum: %d\n", string(data1), simpleChecksum(data1))
	fmt.Printf("Data2: %s, Checksum: %d\n", string(data2), simpleChecksum(data2))

	// ===== Summary =====
	fmt.Println("\n\n📝 ВИСНОВКИ")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("✅ byte = uint8 (0-255)")
	fmt.Println("✅ []byte для роботи з binary data")
	fmt.Println("✅ String immutable, []byte mutable")
	fmt.Println("✅ Hex, Base64 для кодування")
	fmt.Println("✅ bytes package для операцій")
	fmt.Println("✅ Ідеально для ASCII і binary протоколів")
}

// ============= Helper Functions =============

func xorCipher(data []byte, key byte) []byte {
	result := make([]byte, len(data))
	for i, b := range data {
		result[i] = b ^ key
	}
	return result
}

func countByteFrequency(data []byte) map[byte]int {
	freq := make(map[byte]int)
	for _, b := range data {
		freq[b]++
	}
	return freq
}

func simpleChecksum(data []byte) byte {
	var sum byte
	for _, b := range data {
		sum += b
	}
	return sum
}
