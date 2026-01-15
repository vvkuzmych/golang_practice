package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"unicode/utf8"
)

// ============= HTTP with Bytes =============

// Example 1: Reading Request Body (returns []byte)
func handleRequestBody(w http.ResponseWriter, r *http.Request) {
	// io.ReadAll повертає []byte
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fmt.Printf("📦 Request body (bytes): %v\n", bodyBytes)
	fmt.Printf("📝 Request body (string): %s\n", string(bodyBytes))
	fmt.Printf("📊 Body size: %d bytes\n", len(bodyBytes))

	// Відправляємо відповідь (також []byte)
	response := []byte("Received: " + string(bodyBytes))
	w.Write(response)
}

// Example 2: Checking Content-Type and Reading Binary
func handleFileUpload(w http.ResponseWriter, r *http.Request) {
	// Читаємо binary файл як []byte
	fileBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Перевіряємо magic bytes для визначення типу файлу
	fileType := detectFileType(fileBytes)

	fmt.Printf("📁 File uploaded: %d bytes\n", len(fileBytes))
	fmt.Printf("🔍 Detected type: %s\n", fileType)
	fmt.Printf("🔢 First 16 bytes (hex): %x\n", fileBytes[:min(16, len(fileBytes))])

	w.Write([]byte(fmt.Sprintf("File received: %s (%d bytes)", fileType, len(fileBytes))))
}

// detectFileType визначає тип файлу по magic bytes
func detectFileType(data []byte) string {
	if len(data) < 4 {
		return "unknown"
	}

	// PNG: 89 50 4E 47
	if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		return "PNG image"
	}

	// JPEG: FF D8 FF
	if data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return "JPEG image"
	}

	// PDF: 25 50 44 46 (%PDF)
	if data[0] == 0x25 && data[1] == 0x50 && data[2] == 0x44 && data[3] == 0x46 {
		return "PDF document"
	}

	// ZIP: 50 4B 03 04
	if data[0] == 0x50 && data[1] == 0x4B && data[2] == 0x03 && data[3] == 0x04 {
		return "ZIP archive"
	}

	return "unknown"
}

// Example 3: bytes.Buffer for building response
func handleWithBuffer(w http.ResponseWriter, r *http.Request) {
	// Використовуємо bytes.Buffer для побудови відповіді
	var buf bytes.Buffer

	buf.WriteString("Request Info:\n")
	buf.WriteString(fmt.Sprintf("Method: %s\n", r.Method))
	buf.WriteString(fmt.Sprintf("URL: %s\n", r.URL.Path))
	buf.WriteString(fmt.Sprintf("Content-Length: %d bytes\n", r.ContentLength))

	// Конвертуємо в []byte і відправляємо
	responseBytes := buf.Bytes()
	w.Write(responseBytes)
}

// ============= HTTP with Runes (Unicode) =============

// Example 4: Handling Unicode in JSON
type User struct {
	Name    string `json:"name"`
	City    string `json:"city"`
	Message string `json:"message"`
}

func handleUnicodeJSON(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	var user User
	err := json.Unmarshal(bodyBytes, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Аналізуємо Unicode в даних
	fmt.Printf("\n🌍 Unicode JSON Analysis:\n")
	fmt.Printf("Name: %s\n", user.Name)
	fmt.Printf("  - Characters: %d\n", utf8.RuneCountInString(user.Name))
	fmt.Printf("  - Bytes: %d\n", len(user.Name))

	fmt.Printf("City: %s\n", user.City)
	fmt.Printf("  - Characters: %d\n", utf8.RuneCountInString(user.City))
	fmt.Printf("  - Bytes: %d\n", len(user.City))

	fmt.Printf("Message: %s\n", user.Message)
	fmt.Printf("  - Characters: %d\n", utf8.RuneCountInString(user.Message))
	fmt.Printf("  - Bytes: %d\n", len(user.Message))

	// Відповідь
	response := fmt.Sprintf("✅ Received: %s from %s", user.Name, user.City)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(response))
}

// Example 5: URL Encoding with Unicode
func handleUnicodeURL(w http.ResponseWriter, r *http.Request) {
	// Query parameters можуть містити Unicode
	query := r.URL.Query()

	name := query.Get("name")
	city := query.Get("city")

	fmt.Printf("\n🔗 URL Query Parameters:\n")
	fmt.Printf("Raw URL: %s\n", r.URL.String())
	fmt.Printf("Name: %s (%d chars, %d bytes)\n", name, utf8.RuneCountInString(name), len(name))
	fmt.Printf("City: %s (%d chars, %d bytes)\n", city, utf8.RuneCountInString(city), len(city))

	// Перевіряємо валідність UTF-8
	if !utf8.ValidString(name) {
		http.Error(w, "Invalid UTF-8 in name", http.StatusBadRequest)
		return
	}

	response := fmt.Sprintf("Привіт, %s з міста %s!", name, city)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(response))
}

// Example 6: Form Data with Unicode
func handleUnicodeForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	comment := r.FormValue("comment")

	fmt.Printf("\n📝 Form Data Analysis:\n")

	// Аналіз імені
	nameChars := utf8.RuneCountInString(name)
	nameBytes := len(name)
	fmt.Printf("Name: %s\n", name)
	fmt.Printf("  Characters: %d, Bytes: %d, Ratio: %.2f\n",
		nameChars, nameBytes, float64(nameBytes)/float64(nameChars))

	// Аналіз коментаря
	if comment != "" {
		fmt.Printf("Comment: %s\n", comment)
		fmt.Printf("  Characters: %d, Bytes: %d\n",
			utf8.RuneCountInString(comment), len(comment))
	}

	// Перевіряємо чи містить українські літери
	ukrainianCount := 0
	for _, r := range name {
		if isUkrainian(r) {
			ukrainianCount++
		}
	}
	if ukrainianCount > 0 {
		fmt.Printf("  🇺🇦 Ukrainian letters: %d\n", ukrainianCount)
	}

	w.Write([]byte(fmt.Sprintf("Form received from %s", name)))
}

func isUkrainian(r rune) bool {
	return (r >= 'А' && r <= 'Я') || (r >= 'а' && r <= 'я') ||
		r == 'Є' || r == 'І' || r == 'Ї' || r == 'Ґ' ||
		r == 'є' || r == 'і' || r == 'ї' || r == 'ґ'
}

// Example 7: Content-Length Validation (bytes vs chars)
func handleContentValidation(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	bodyString := string(bodyBytes)

	fmt.Printf("\n✅ Content Validation:\n")
	fmt.Printf("Content-Length header: %d\n", r.ContentLength)
	fmt.Printf("Actual bytes read: %d\n", len(bodyBytes))
	fmt.Printf("Characters (runes): %d\n", utf8.RuneCountInString(bodyString))

	// ВАЖЛИВО: Content-Length завжди в байтах, не символах!
	if int64(len(bodyBytes)) != r.ContentLength {
		fmt.Printf("⚠️  Mismatch detected!\n")
	}

	w.Write([]byte("Validation complete"))
}

// ============= Practical Examples =============

// Example 8: API Response with Mixed Content
func handleMixedContent(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":  "success",
		"message": "Вітаємо! Welcome! 👋",
		"user": map[string]string{
			"name": "Олександр",
			"city": "Київ",
		},
		"stats": map[string]int{
			"bytes": len("Вітаємо! Welcome! 👋"),
			"chars": utf8.RuneCountInString("Вітаємо! Welcome! 👋"),
		},
	}

	// JSON marshal повертає []byte
	jsonBytes, _ := json.MarshalIndent(response, "", "  ")

	fmt.Printf("\n📤 API Response:\n")
	fmt.Printf("Response size: %d bytes\n", len(jsonBytes))
	fmt.Printf("JSON:\n%s\n", string(jsonBytes))

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(jsonBytes)
}

// ============= Helper Functions =============

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func printSeparator(title string) {
	fmt.Printf("\n%s\n", strings.Repeat("═", 60))
	fmt.Printf("  %s\n", title)
	fmt.Printf("%s\n", strings.Repeat("═", 60))
}

// ============= Main Demo =============

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║      HTTP Examples: byte & rune in Action               ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")

	// Demo 1: Request Body ([]byte)
	printSeparator("1. Reading Request Body ([]byte)")

	req1 := httptest.NewRequest("POST", "/api/data",
		strings.NewReader("Hello, World!"))
	w1 := httptest.NewRecorder()
	handleRequestBody(w1, req1)

	// Demo 2: Binary File Upload
	printSeparator("2. Binary File Upload (magic bytes)")

	// Симулюємо PNG файл (magic bytes)
	pngData := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	req2 := httptest.NewRequest("POST", "/upload", bytes.NewReader(pngData))
	w2 := httptest.NewRecorder()
	handleFileUpload(w2, req2)

	// Demo 3: bytes.Buffer
	printSeparator("3. Using bytes.Buffer")

	req3 := httptest.NewRequest("GET", "/info", nil)
	w3 := httptest.NewRecorder()
	handleWithBuffer(w3, req3)
	fmt.Printf("Response:\n%s\n", w3.Body.String())

	// Demo 4: Unicode in JSON
	printSeparator("4. Unicode in JSON (runes)")

	jsonData := `{
		"name": "Олександр",
		"city": "Київ",
		"message": "Привіт! Hello! 👋"
	}`
	req4 := httptest.NewRequest("POST", "/api/user", strings.NewReader(jsonData))
	w4 := httptest.NewRecorder()
	handleUnicodeJSON(w4, req4)

	// Demo 5: URL Encoding with Unicode
	printSeparator("5. URL Query Parameters with Unicode")

	// URL-encoded Ukrainian text
	encodedName := url.QueryEscape("Олена")
	encodedCity := url.QueryEscape("Львів")
	urlStr := fmt.Sprintf("/api/greet?name=%s&city=%s", encodedName, encodedCity)

	req5 := httptest.NewRequest("GET", urlStr, nil)
	w5 := httptest.NewRecorder()
	handleUnicodeURL(w5, req5)
	fmt.Printf("Response: %s\n", w5.Body.String())

	// Demo 6: Form Data
	printSeparator("6. Form Data with Unicode")

	formData := url.Values{}
	formData.Set("name", "Марія")
	formData.Set("comment", "Дуже добре! 👍")

	req6 := httptest.NewRequest("POST", "/api/form",
		strings.NewReader(formData.Encode()))
	req6.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w6 := httptest.NewRecorder()
	handleUnicodeForm(w6, req6)

	// Demo 7: Content-Length
	printSeparator("7. Content-Length (bytes, not chars!)")

	ukrainianText := "Привіт"
	req7 := httptest.NewRequest("POST", "/api/validate",
		strings.NewReader(ukrainianText))
	req7.ContentLength = int64(len(ukrainianText))
	w7 := httptest.NewRecorder()
	handleContentValidation(w7, req7)

	// Demo 8: Mixed Content API
	printSeparator("8. API Response with Mixed Content")

	req8 := httptest.NewRequest("GET", "/api/info", nil)
	w8 := httptest.NewRecorder()
	handleMixedContent(w8, req8)

	// ===== Practical Tips =====
	printSeparator("💡 Key Takeaways")

	fmt.Println("\n🔹 BYTES in HTTP:")
	fmt.Println("  ✅ io.ReadAll() returns []byte")
	fmt.Println("  ✅ http.ResponseWriter.Write() accepts []byte")
	fmt.Println("  ✅ Content-Length is in BYTES")
	fmt.Println("  ✅ Binary files are []byte")
	fmt.Println("  ✅ JSON Marshal/Unmarshal uses []byte")
	fmt.Println("  ✅ Magic bytes for file type detection")

	fmt.Println("\n🔹 RUNES in HTTP:")
	fmt.Println("  ✅ URL parameters can contain Unicode")
	fmt.Println("  ✅ JSON strings support Unicode")
	fmt.Println("  ✅ Form data can have Ukrainian text")
	fmt.Println("  ✅ Content validation needs rune count")
	fmt.Println("  ✅ Always validate with utf8.ValidString()")

	fmt.Println("\n⚠️  IMPORTANT:")
	fmt.Println("  • Content-Length = BYTES, not characters")
	fmt.Println("  • Ukrainian 'Привіт' = 6 chars, 12 bytes")
	fmt.Println("  • Always set charset=utf-8 in headers")
	fmt.Println("  • URL-encode Unicode: url.QueryEscape()")

	// ===== Real-world Example =====
	printSeparator("🌍 Real-world Scenario")

	fmt.Println("\nScenario: User registration form with Ukrainian name")

	// Клієнт відправляє
	name := "Олександра"
	nameBytes := []byte(name)

	fmt.Printf("\n📤 Client sends:\n")
	fmt.Printf("   Name: %s\n", name)
	fmt.Printf("   As bytes: %v\n", nameBytes)
	fmt.Printf("   Content-Length: %d (bytes!)\n", len(nameBytes))

	// Сервер отримує
	fmt.Printf("\n📥 Server receives:\n")
	receivedBytes := nameBytes // з io.ReadAll()
	receivedString := string(receivedBytes)

	fmt.Printf("   Bytes received: %d\n", len(receivedBytes))
	fmt.Printf("   String: %s\n", receivedString)
	fmt.Printf("   Characters: %d\n", utf8.RuneCountInString(receivedString))

	// Валідація
	fmt.Printf("\n✅ Validation:\n")
	if utf8.ValidString(receivedString) {
		fmt.Printf("   ✓ Valid UTF-8\n")
	}

	charCount := utf8.RuneCountInString(receivedString)
	if charCount >= 2 && charCount <= 50 {
		fmt.Printf("   ✓ Name length OK (%d characters)\n", charCount)
	}

	// Відповідь
	response := fmt.Sprintf(`{"status":"success","message":"Вітаємо, %s!"}`, name)
	responseBytes := []byte(response)

	fmt.Printf("\n📤 Server responds:\n")
	fmt.Printf("   JSON: %s\n", response)
	fmt.Printf("   Content-Length: %d bytes\n", len(responseBytes))
	fmt.Printf("   Content-Type: application/json; charset=utf-8\n")

	// ===== Summary =====
	fmt.Println("\n\n" + strings.Repeat("═", 60))
	fmt.Println("📚 SUMMARY")
	fmt.Println(strings.Repeat("═", 60))

	fmt.Println("\nWhen to use []byte:")
	fmt.Println("  1. Reading HTTP request bodies")
	fmt.Println("  2. Writing HTTP responses")
	fmt.Println("  3. Binary file uploads/downloads")
	fmt.Println("  4. JSON encoding/decoding")
	fmt.Println("  5. Content-Length calculations")

	fmt.Println("\nWhen to use rune:")
	fmt.Println("  1. Counting characters in user input")
	fmt.Println("  2. Validating text length (not byte length!)")
	fmt.Println("  3. Processing Ukrainian/multilingual text")
	fmt.Println("  4. Substring operations on Unicode")
	fmt.Println("  5. Character-based validation rules")

	fmt.Println("\n🎯 Remember:")
	fmt.Println("  HTTP works with BYTES ([]byte)")
	fmt.Println("  Users think in CHARACTERS (runes)")
	fmt.Println("  Always convert appropriately!")
}
