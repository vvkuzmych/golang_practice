package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// ============= Encoding Functions =============

// HexEncode кодує []byte в hex string
func HexEncode(data []byte) string {
	return hex.EncodeToString(data)
}

// HexDecode декодує hex string в []byte
func HexDecode(hexString string) ([]byte, error) {
	return hex.DecodeString(hexString)
}

// Base64Encode кодує []byte в base64
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base64Decode декодує base64 в []byte
func Base64Decode(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}

// ROT13 застосовує ROT13 шифр (для ASCII A-Z, a-z)
func ROT13(data []byte) []byte {
	result := make([]byte, len(data))
	for i, b := range data {
		if b >= 'A' && b <= 'Z' {
			result[i] = 'A' + (b-'A'+13)%26
		} else if b >= 'a' && b <= 'z' {
			result[i] = 'a' + (b-'a'+13)%26
		} else {
			result[i] = b
		}
	}
	return result
}

// Caesar зсуває символи на shift позицій
func Caesar(data []byte, shift int) []byte {
	result := make([]byte, len(data))
	shift = shift % 26 // normalize shift
	if shift < 0 {
		shift += 26
	}

	for i, b := range data {
		if b >= 'A' && b <= 'Z' {
			result[i] = byte('A' + (int(b-'A')+shift)%26)
		} else if b >= 'a' && b <= 'z' {
			result[i] = byte('a' + (int(b-'a')+shift)%26)
		} else {
			result[i] = b
		}
	}
	return result
}

// XORCipher застосовує XOR з ключем
func XORCipher(data []byte, key byte) []byte {
	result := make([]byte, len(data))
	for i, b := range data {
		result[i] = b ^ key
	}
	return result
}

// ============= Statistics Functions =============

// ByteStats повертає статистику про байти
func ByteStats(data []byte) map[string]int {
	stats := map[string]int{
		"total":   len(data),
		"unique":  0,
		"letters": 0,
		"digits":  0,
		"spaces":  0,
		"upper":   0,
		"lower":   0,
	}

	freq := make(map[byte]int)
	for _, b := range data {
		freq[b]++

		if (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') {
			stats["letters"]++
		}
		if b >= 'A' && b <= 'Z' {
			stats["upper"]++
		}
		if b >= 'a' && b <= 'z' {
			stats["lower"]++
		}
		if b >= '0' && b <= '9' {
			stats["digits"]++
		}
		if b == ' ' {
			stats["spaces"]++
		}
	}
	stats["unique"] = len(freq)

	return stats
}

// FrequencyAnalysis повертає частоту кожного байту
func FrequencyAnalysis(data []byte) map[byte]int {
	freq := make(map[byte]int)
	for _, b := range data {
		freq[b]++
	}
	return freq
}

// MostFrequent знаходить найчастіший байт
func MostFrequent(data []byte) (byte, int) {
	freq := FrequencyAnalysis(data)

	var mostByte byte
	maxCount := 0

	for b, count := range freq {
		if count > maxCount {
			mostByte = b
			maxCount = count
		}
	}

	return mostByte, maxCount
}

// SimpleChecksum обчислює просту контрольну суму
func SimpleChecksum(data []byte) byte {
	var sum byte
	for _, b := range data {
		sum += b
	}
	return sum
}

// ToBinary конвертує в binary string
func ToBinary(data []byte) string {
	result := ""
	for _, b := range data {
		result += fmt.Sprintf("%08b ", b)
	}
	return result
}

// ============= Main =============

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║      Byte Encoder Solution               ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	message := "Hello, World!"

	// ===== Hex Encoding =====
	fmt.Println("\n🔹 Hex Encoding")
	fmt.Println("─────────────────────────────────────────")

	fmt.Printf("Original: %s\n", message)

	hexEncoded := HexEncode([]byte(message))
	fmt.Printf("Hex: %s\n", hexEncoded)

	hexDecoded, _ := HexDecode(hexEncoded)
	fmt.Printf("Decoded: %s\n", string(hexDecoded))

	// ===== Base64 Encoding =====
	fmt.Println("\n🔹 Base64 Encoding")
	fmt.Println("─────────────────────────────────────────")

	fmt.Printf("Original: %s\n", message)

	b64Encoded := Base64Encode([]byte(message))
	fmt.Printf("Base64: %s\n", b64Encoded)

	b64Decoded, _ := Base64Decode(b64Encoded)
	fmt.Printf("Decoded: %s\n", string(b64Decoded))

	// ===== ROT13 Cipher =====
	fmt.Println("\n🔹 ROT13 Cipher")
	fmt.Println("─────────────────────────────────────────")

	fmt.Printf("Original: %s\n", message)

	rot13Encrypted := ROT13([]byte(message))
	fmt.Printf("ROT13: %s\n", string(rot13Encrypted))

	rot13Decrypted := ROT13(rot13Encrypted)
	fmt.Printf("ROT13 twice: %s\n", string(rot13Decrypted))

	// ===== Caesar Cipher =====
	fmt.Println("\n🔹 Caesar Cipher (shift = 3)")
	fmt.Println("─────────────────────────────────────────")

	fmt.Printf("Original: %s\n", message)

	caesarEncrypted := Caesar([]byte(message), 3)
	fmt.Printf("Encrypted: %s\n", string(caesarEncrypted))

	caesarDecrypted := Caesar(caesarEncrypted, -3)
	fmt.Printf("Decrypted: %s\n", string(caesarDecrypted))

	// ===== XOR Cipher =====
	fmt.Println("\n🔹 XOR Cipher (key = 42)")
	fmt.Println("─────────────────────────────────────────")

	fmt.Printf("Original: %s\n", message)

	xorEncrypted := XORCipher([]byte(message), 42)
	fmt.Printf("Encrypted: %v\n", xorEncrypted)

	xorDecrypted := XORCipher(xorEncrypted, 42)
	fmt.Printf("Decrypted: %s\n", string(xorDecrypted))

	// ===== Byte Statistics =====
	fmt.Println("\n🔹 Byte Statistics")
	fmt.Println("─────────────────────────────────────────")

	fmt.Printf("Text: %s\n", message)

	stats := ByteStats([]byte(message))
	fmt.Printf("Total bytes: %d\n", stats["total"])
	fmt.Printf("Unique bytes: %d\n", stats["unique"])
	fmt.Printf("Letters: %d\n", stats["letters"])
	fmt.Printf("  Uppercase: %d\n", stats["upper"])
	fmt.Printf("  Lowercase: %d\n", stats["lower"])
	fmt.Printf("Digits: %d\n", stats["digits"])
	fmt.Printf("Spaces: %d\n", stats["spaces"])

	// ===== Frequency Analysis =====
	fmt.Println("\n🔹 Frequency Analysis")
	fmt.Println("─────────────────────────────────────────")

	mostByte, mostCount := MostFrequent([]byte(message))
	fmt.Printf("Most frequent: '%c' (%d times)\n", mostByte, mostCount)

	fmt.Println("\nAll frequencies:")
	freq := FrequencyAnalysis([]byte(message))
	for b, count := range freq {
		if b >= 32 && b <= 126 { // printable
			fmt.Printf("  '%c': %d\n", b, count)
		}
	}

	// ===== Checksum =====
	fmt.Println("\n🔹 Checksums")
	fmt.Println("─────────────────────────────────────────")

	checksum1 := SimpleChecksum([]byte(message))
	checksum2 := SimpleChecksum([]byte("Hello, world!"))

	fmt.Printf("%s → checksum: %d\n", message, checksum1)
	fmt.Printf("Hello, world! → checksum: %d\n", checksum2)

	// ===== Binary Representation =====
	fmt.Println("\n🔹 Binary Representation")
	fmt.Println("─────────────────────────────────────────")

	short := "Hi"
	fmt.Printf("Text: %s\n", short)
	fmt.Printf("Binary: %s\n", ToBinary([]byte(short)))

	// ===== Error Handling =====
	fmt.Println("\n🔹 Error Handling")
	fmt.Println("─────────────────────────────────────────")

	_, err := HexDecode("invalid hex!")
	if err != nil {
		fmt.Printf("❌ Hex decode error: %v\n", err)
	}

	_, err = Base64Decode("invalid base64!")
	if err != nil {
		fmt.Printf("❌ Base64 decode error: %v\n", err)
	}

	// Valid decoding
	validHex := "48656c6c6f"
	decoded, err := HexDecode(validHex)
	if err == nil {
		fmt.Printf("✅ Valid hex decoded: %s\n", string(decoded))
	}

	// ===== Practical Example =====
	fmt.Println("\n🔹 Практичний приклад: Секретне повідомлення")
	fmt.Println("─────────────────────────────────────────")

	secret := "Secret Message"
	fmt.Printf("1. Original: %s\n", secret)

	// XOR encrypt
	encrypted := XORCipher([]byte(secret), 123)
	fmt.Printf("2. XOR encrypted: %v\n", encrypted)

	// Base64 encode
	encoded := Base64Encode(encrypted)
	fmt.Printf("3. Base64 encoded: %s\n", encoded)

	// Для декодування:
	decodedB64, _ := Base64Decode(encoded)
	decrypted := XORCipher(decodedB64, 123)
	fmt.Printf("4. Decrypted: %s\n", string(decrypted))

	// ===== Summary =====
	fmt.Println("\n\n📝 ВИСНОВКИ")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("✅ Hex encoding для binary data")
	fmt.Println("✅ Base64 для передачі binary в text")
	fmt.Println("✅ ROT13/Caesar для простого шифрування")
	fmt.Println("✅ XOR для швидкого шифрування")
	fmt.Println("✅ Checksum для перевірки цілісності")
	fmt.Println("✅ Frequency analysis для криптоаналізу")
}
