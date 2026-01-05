# Standard Interfaces в Go

Ця папка містить приклади найпоширеніших інтерфейсів зі стандартної бібліотеки Go.

---

## 📚 Файли

| Файл | Інтерфейси | Опис |
|------|-----------|------|
| `01_io_reader_writer.go` | `io.Reader`, `io.Writer` | Читання і запис даних |
| `02_fmt_stringer.go` | `fmt.Stringer` | Текстове представлення |
| `03_error_interface.go` | `error` | Обробка помилок |
| `04_sort_interface.go` | `sort.Interface` | Сортування колекцій |
| `05_http_handler.go` | `http.Handler` | HTTP обробники |
| `06_json_marshaler.go` | `json.Marshaler/Unmarshaler` | JSON серіалізація |
| `07_io_closer.go` | `io.Closer`, `io.ReadCloser` | Закриття ресурсів |
| `08_context_usage.go` | `context.Context` | Контекст і скасування |

---

## 🚀 Як запускати

Кожен файл - це окрема програма з `main`:

```bash
cd standard_interfaces

# Запустити конкретний приклад
go run 01_io_reader_writer.go
go run 02_fmt_stringer.go
go run 03_error_interface.go
# і т.д.
```

---

## 📖 Огляд інтерфейсів

### 1. io.Reader і io.Writer

Найважливіші інтерфейси в Go для роботи з даними.

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}
```

**Використання:**
- Файли (`os.File`)
- Мережа (`net.Conn`)
- Буфери (`bytes.Buffer`)
- HTTP (`http.Request.Body`, `http.ResponseWriter`)
- Стандартний ввід/вивід (`os.Stdin`, `os.Stdout`)

---

### 2. fmt.Stringer

Визначає як тип перетворюється на рядок.

```go
type Stringer interface {
    String() string
}
```

**Використання:**
- Красивий вивід структур
- Логування
- Відлагодження
- `fmt.Println()` автоматично викликає `String()`

---

### 3. error

Стандартний інтерфейс для помилок.

```go
type error interface {
    Error() string
}
```

**Використання:**
- Повернення помилок з функцій
- Власні типи помилок
- Обгортання помилок (`fmt.Errorf`, `errors.Wrap`)

---

### 4. sort.Interface

Для сортування колекцій.

```go
type Interface interface {
    Len() int
    Less(i, j int) bool
    Swap(i, j int)
}
```

**Використання:**
- Сортування слайсів
- Власна логіка порівняння
- `sort.Sort()`, `sort.Stable()`

---

### 5. http.Handler

Обробка HTTP запитів.

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

**Використання:**
- Web сервери
- API endpoints
- Middleware
- Роутери

---

### 6. json.Marshaler / Unmarshaler

Контроль JSON серіалізації.

```go
type Marshaler interface {
    MarshalJSON() ([]byte, error)
}

type Unmarshaler interface {
    UnmarshalJSON([]byte) error
}
```

**Використання:**
- Власний формат JSON
- Приховування полів
- Трансформація даних

---

### 7. io.Closer

Закриття ресурсів.

```go
type Closer interface {
    Close() error
}
```

**Використання:**
- Файли
- Мережеві з'єднання
- Database connections
- Завжди з `defer`

---

### 8. context.Context

Контроль виконання і скасування.

```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key interface{}) interface{}
}
```

**Використання:**
- Таймаути
- Скасування операцій
- Передача метаданих
- HTTP requests

---

## 💡 Чому це важливо?

### 1. Композиція інтерфейсів

Go комбінує маленькі інтерфейси в більші:

```go
type ReadWriter interface {
    Reader
    Writer
}

type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}
```

### 2. Dependency Injection

Залежність від інтерфейсу, не від конкретного типу:

```go
func ProcessData(r io.Reader) error {
    // Працює з будь-яким Reader: файл, мережа, рядок
}
```

### 3. Тестування

Легко створити mock:

```go
type MockReader struct {
    data []byte
}

func (m *MockReader) Read(p []byte) (int, error) {
    // mock implementation
}
```

---

## 🎯 Практичні приклади

### Приклад 1: Універсальна функція

```go
// Працює з будь-яким Reader
func CountLines(r io.Reader) (int, error) {
    scanner := bufio.NewScanner(r)
    count := 0
    for scanner.Scan() {
        count++
    }
    return count, scanner.Err()
}

// Можна використати з:
CountLines(os.Stdin)                          // консоль
CountLines(strings.NewReader("a\nb\nc"))      // рядок
file, _ := os.Open("file.txt")
CountLines(file)                              // файл
```

### Приклад 2: Middleware Pattern

```go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Request: %s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}
```

### Приклад 3: Власний Error

```go
type ValidationError struct {
    Field string
    Error string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Error)
}
```

---

## 📊 Статистика використання

Найчастіше використовувані інтерфейси в Go проектах:

1. **io.Reader/Writer** - 90% проектів
2. **error** - 100% проектів
3. **fmt.Stringer** - 70% проектів
4. **http.Handler** - 80% web проектів
5. **context.Context** - 85% сучасних проектів
6. **json.Marshaler** - 60% API проектів
7. **io.Closer** - 75% проектів
8. **sort.Interface** - 40% проектів

---

## 🎓 Навчальні цілі

Після вивчення цих прикладів ви будете:

- ✅ Розуміти всі основні інтерфейси Go
- ✅ Вміти використовувати їх у власному коді
- ✅ Створювати власні реалізації
- ✅ Розуміти стандартну бібліотеку Go
- ✅ Писати ідіоматичний Go код
- ✅ Розуміти код інших Go розробників

---

## 📚 Додаткові ресурси

### Офіційна документація
- [io package](https://pkg.go.dev/io)
- [fmt package](https://pkg.go.dev/fmt)
- [errors package](https://pkg.go.dev/errors)
- [sort package](https://pkg.go.dev/sort)
- [net/http package](https://pkg.go.dev/net/http)
- [encoding/json package](https://pkg.go.dev/encoding/json)
- [context package](https://pkg.go.dev/context)

### Статті
- [Go Interfaces Explained](https://go.dev/doc/effective_go#interfaces)
- [The Laws of Reflection](https://go.dev/blog/laws-of-reflection)
- [Error handling in Go](https://go.dev/blog/error-handling-and-go)

---

## 🎯 Порядок вивчення

Рекомендуємо вивчати в такому порядку:

1. **io.Reader/Writer** - фундамент
2. **fmt.Stringer** - просто і корисно
3. **error** - критично важливо
4. **io.Closer** - управління ресурсами
5. **sort.Interface** - практичний приклад
6. **json.Marshaler** - для API
7. **http.Handler** - для web
8. **context.Context** - для конкурентності

---

## ⚡ Швидкий старт

```bash
# Запустити всі приклади підряд
for file in *.go; do
    echo "=== Running $file ==="
    go run "$file"
    echo ""
done
```

---

**Удачі у вивченні стандартних інтерфейсів Go! 🚀**

