package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

// ============= io.Reader =============

// io.Reader interface:
// type Reader interface {
//     Read(p []byte) (n int, err error)
// }

// Власна реалізація Reader
type UppercaseReader struct {
	reader io.Reader
}

func (u *UppercaseReader) Read(p []byte) (int, error) {
	n, err := u.reader.Read(p)
	for i := 0; i < n; i++ {
		if p[i] >= 'a' && p[i] <= 'z' {
			p[i] = p[i] - 32 // перетворити на великі літери
		}
	}
	return n, err
}

// ============= io.Writer =============

// io.Writer interface:
// type Writer interface {
//     Write(p []byte) (n int, err error)
// }

// Власна реалізація Writer
type CountingWriter struct {
	writer io.Writer
	count  int
}

func (c *CountingWriter) Write(p []byte) (int, error) {
	n, err := c.writer.Write(p)
	c.count += n
	return n, err
}

func (c *CountingWriter) BytesWritten() int {
	return c.count
}

// ============= Префіксний Writer =============

type PrefixWriter struct {
	writer io.Writer
	prefix string
}

func (p *PrefixWriter) Write(data []byte) (int, error) {
	prefixed := []byte(p.prefix + string(data))
	return p.writer.Write(prefixed)
}

// ============= Helper Functions =============

// Універсальна функція - працює з будь-яким Reader
func ReadAll(r io.Reader) (string, error) {
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, r)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Універсальна функція - працює з будь-яким Writer
func WriteMessage(w io.Writer, message string) error {
	_, err := w.Write([]byte(message))
	return err
}

// Copy з прогресом
func CopyWithProgress(dst io.Writer, src io.Reader) (int64, error) {
	var written int64
	buf := make([]byte, 32*1024) // 32KB buffer

	for {
		nr, err := src.Read(buf)
		if nr > 0 {
			nw, err := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if err != nil {
				return written, err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return written, err
		}
	}
	return written, nil
}

// ============= Main =============

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║    io.Reader & io.Writer Interface      ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	// ===== strings.Reader =====
	fmt.Println("\n🔹 strings.Reader (реалізує io.Reader)")
	fmt.Println("─────────────────────────────────────────")

	strReader := strings.NewReader("Hello, Go!")
	data, _ := ReadAll(strReader)
	fmt.Printf("Прочитано: %s\n", data)

	// ===== bytes.Buffer =====
	fmt.Println("\n🔹 bytes.Buffer (реалізує Reader і Writer)")
	fmt.Println("─────────────────────────────────────────")

	var buf bytes.Buffer
	buf.WriteString("Привіт")
	buf.WriteString(" ")
	buf.WriteString("світ!")
	fmt.Printf("Buffer: %s\n", buf.String())

	// ===== Custom UppercaseReader =====
	fmt.Println("\n🔹 Custom UppercaseReader")
	fmt.Println("─────────────────────────────────────────")

	original := strings.NewReader("hello world")
	uppercase := &UppercaseReader{reader: original}
	result, _ := ReadAll(uppercase)
	fmt.Printf("Оригінал: hello world\n")
	fmt.Printf("Uppercase: %s\n", result)

	// ===== Custom CountingWriter =====
	fmt.Println("\n🔹 Custom CountingWriter")
	fmt.Println("─────────────────────────────────────────")

	var output bytes.Buffer
	counter := &CountingWriter{writer: &output}

	counter.Write([]byte("Перший рядок\n"))
	counter.Write([]byte("Другий рядок\n"))
	counter.Write([]byte("Третій рядок\n"))

	fmt.Printf("Записано байтів: %d\n", counter.BytesWritten())
	fmt.Printf("Вміст:\n%s", output.String())

	// ===== PrefixWriter =====
	fmt.Println("\n🔹 PrefixWriter")
	fmt.Println("─────────────────────────────────────────")

	var logBuf bytes.Buffer
	logger := &PrefixWriter{writer: &logBuf, prefix: "[LOG] "}

	logger.Write([]byte("Application started\n"))
	logger.Write([]byte("Processing data\n"))
	logger.Write([]byte("Done\n"))

	fmt.Print(logBuf.String())

	// ===== io.Copy =====
	fmt.Println("\n🔹 io.Copy (Reader → Writer)")
	fmt.Println("─────────────────────────────────────────")

	source := strings.NewReader("Копіюємо цей текст")
	var destination bytes.Buffer

	written, _ := io.Copy(&destination, source)
	fmt.Printf("Скопійовано %d байтів\n", written)
	fmt.Printf("Результат: %s\n", destination.String())

	// ===== MultiWriter =====
	fmt.Println("\n🔹 io.MultiWriter (запис в кілька місць)")
	fmt.Println("─────────────────────────────────────────")

	var buf1, buf2, buf3 bytes.Buffer
	multi := io.MultiWriter(&buf1, &buf2, &buf3)

	multi.Write([]byte("Цей текст іде в 3 місця!\n"))

	fmt.Printf("Buffer 1: %s", buf1.String())
	fmt.Printf("Buffer 2: %s", buf2.String())
	fmt.Printf("Buffer 3: %s", buf3.String())

	// ===== TeeReader =====
	fmt.Println("\n🔹 io.TeeReader (читання + копіювання)")
	fmt.Println("─────────────────────────────────────────")

	input := strings.NewReader("Original data")
	var copy bytes.Buffer
	tee := io.TeeReader(input, &copy)

	// Читаємо з tee
	result2, _ := ReadAll(tee)

	fmt.Printf("Прочитано: %s\n", result2)
	fmt.Printf("Копія: %s\n", copy.String())

	// ===== LimitReader =====
	fmt.Println("\n🔹 io.LimitReader (обмеження читання)")
	fmt.Println("─────────────────────────────────────────")

	longText := strings.NewReader("Це довгий текст, але ми прочитаємо тільки частину")
	limited := io.LimitReader(longText, 20)

	partial, _ := ReadAll(limited)
	fmt.Printf("Прочитано перші 20 байтів: %s\n", partial)

	// ===== Pipe =====
	fmt.Println("\n🔹 io.Pipe (синхронний канал)")
	fmt.Println("─────────────────────────────────────────")

	pr, pw := io.Pipe()

	// Goroutine для запису
	go func() {
		defer pw.Close()
		pw.Write([]byte("Дані через pipe\n"))
	}()

	// Читання
	pipeData, _ := ReadAll(pr)
	fmt.Printf("Отримано: %s", pipeData)

	// ===== Комбінація =====
	fmt.Println("\n🔹 Комбінація: Uppercase + Prefix + Count")
	fmt.Println("─────────────────────────────────────────")

	inputText := strings.NewReader("go is awesome")
	uppercaseReader := &UppercaseReader{reader: inputText}

	var finalOutput bytes.Buffer
	prefixWriter := &PrefixWriter{writer: &finalOutput, prefix: ">>> "}
	countWriter := &CountingWriter{writer: prefixWriter}

	io.Copy(countWriter, uppercaseReader)

	fmt.Printf("Результат: %s", finalOutput.String())
	fmt.Printf("Записано байтів: %d\n", countWriter.BytesWritten())

	// ===== Приклад з os.Stdout =====
	fmt.Println("\n🔹 Запис в os.Stdout (io.Writer)")
	fmt.Println("─────────────────────────────────────────")

	WriteMessage(os.Stdout, "Це пишеться напряму в stdout!\n")

	// ===== Summary =====
	fmt.Println("\n📝 ВИСНОВКИ")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("✅ io.Reader - універсальний інтерфейс для читання")
	fmt.Println("   • strings.Reader, bytes.Buffer, os.File")
	fmt.Println("   • http.Response.Body, net.Conn")
	fmt.Println()
	fmt.Println("✅ io.Writer - універсальний інтерфейс для запису")
	fmt.Println("   • bytes.Buffer, os.File, os.Stdout")
	fmt.Println("   • http.ResponseWriter")
	fmt.Println()
	fmt.Println("💡 Переваги:")
	fmt.Println("   • Одна функція - багато реалізацій")
	fmt.Println("   • Легко тестувати (bytes.Buffer)")
	fmt.Println("   • Композиція (io.MultiWriter, io.TeeReader)")
	fmt.Println("   • Стандарт в екосистемі Go")
}
