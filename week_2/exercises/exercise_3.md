# Вправа 3: Storage Interface

## Ціль
Створити абстракцію для зберігання даних з різними реалізаціями (memory та file).

---

## Завдання

Створіть програму `storage.go`, яка:

1. Має інтерфейс `Storage` для зберігання key-value даних
2. Реалізує `MemoryStorage` (зберігання в пам'яті)
3. Реалізує `FileStorage` (зберігання у файлі)
4. Демонструє можливість зміни реалізації без зміни коду

---

## Вимоги

### Інтерфейс Storage
```go
type Storage interface {
    Save(key, value string) error
    Load(key string) (string, error)
    Delete(key string) error
    Exists(key string) bool
    Keys() []string
}
```

### Обов'язкові реалізації:
1. **MemoryStorage** - зберігання в map
2. **FileStorage** - зберігання у текстовому файлі (можна спрощений варіант)

### Application Layer:
```go
type DataManager struct {
    storage Storage  // залежність від інтерфейсу
}
```

---

## Приклад використання

```go
func main() {
    // Memory Storage
    memStorage := NewMemoryStorage()
    manager1 := NewDataManager(memStorage)
    
    manager1.Set("name", "Іван")
    manager1.Set("age", "25")
    
    value, _ := manager1.Get("name")
    fmt.Println(value)  // Іван
    
    // File Storage
    fileStorage := NewFileStorage("data.txt")
    manager2 := NewDataManager(fileStorage)
    
    manager2.Set("config", "value1")
    // ... працює так само!
}
```

---

## Підказки

### 1. MemoryStorage
```go
type MemoryStorage struct {
    data map[string]string
}

func NewMemoryStorage() *MemoryStorage {
    return &MemoryStorage{
        data: make(map[string]string),
    }
}

func (m *MemoryStorage) Save(key, value string) error {
    m.data[key] = value
    return nil
}

func (m *MemoryStorage) Load(key string) (string, error) {
    value, exists := m.data[key]
    if !exists {
        return "", fmt.Errorf("key not found: %s", key)
    }
    return value, nil
}
```

### 2. FileStorage (спрощений)
```go
type FileStorage struct {
    filename string
    data     map[string]string
}

func NewFileStorage(filename string) *FileStorage {
    fs := &FileStorage{
        filename: filename,
        data:     make(map[string]string),
    }
    fs.loadFromFile()
    return fs
}

func (f *FileStorage) Save(key, value string) error {
    f.data[key] = value
    return f.saveToFile()
}

func (f *FileStorage) saveToFile() error {
    // Записати map у файл
    // Формат: key=value (кожен рядок)
}
```

### 3. DataManager
```go
type DataManager struct {
    storage Storage
}

func NewDataManager(s Storage) *DataManager {
    return &DataManager{storage: s}
}

func (d *DataManager) Set(key, value string) error {
    return d.storage.Save(key, value)
}

func (d *DataManager) Get(key string) (string, error) {
    return d.storage.Load(key)
}
```

---

## Очікуваний вивід

```
=== Memory Storage ===
✅ Saved: name=Іван
✅ Saved: age=25
✅ Saved: city=Київ

📖 Loading data:
  name: Іван
  age: 25
  city: Київ

🔑 All keys: [name age city]

🗑️  Deleted: age

📖 After deletion:
  name: Іван
  city: Київ

=== File Storage ===
✅ Saved to file: config=production
✅ Saved to file: version=1.0
✅ Saved to file: debug=false

📁 File content:
config=production
version=1.0
debug=false

=== Using Same Interface ===
Memory Stats: 2 keys
File Stats: 3 keys

💡 Обидві реалізації працюють через один інтерфейс!
```

---

## Бонус завдання

1. **JSON Storage**:
   ```go
   type JSONStorage struct {
       filename string
       data     map[string]string
   }
   
   func (j *JSONStorage) saveToFile() error {
       return json.Marshal(j.data)
   }
   ```

2. **Cache Layer**:
   ```go
   type CachedStorage struct {
       storage Storage
       cache   map[string]string
   }
   
   // Кешує Load операції
   func (c *CachedStorage) Load(key string) (string, error) {
       if value, ok := c.cache[key]; ok {
           return value, nil
       }
       value, err := c.storage.Load(key)
       if err == nil {
           c.cache[key] = value
       }
       return value, err
   }
   ```

3. **Mock Storage**:
   ```go
   type MockStorage struct {
       saveCalled  int
       loadCalled  int
       shouldFail  bool
   }
   
   // Для тестування
   ```

4. **Encrypted Storage**:
   ```go
   type EncryptedStorage struct {
       storage Storage
       key     []byte
   }
   
   func (e *EncryptedStorage) Save(key, value string) error {
       encrypted := encrypt(value, e.key)
       return e.storage.Save(key, encrypted)
   }
   ```

5. **Transaction Support**:
   ```go
   type TransactionalStorage interface {
       Storage
       Begin() Transaction
   }
   
   type Transaction interface {
       Commit() error
       Rollback() error
   }
   ```

---

## Критерії оцінки

- ✅ Інтерфейс `Storage` правильно оголошений
- ✅ MemoryStorage працює коректно
- ✅ FileStorage зберігає дані у файл
- ✅ DataManager працює з обома через інтерфейс
- ✅ Обробка помилок (ключ не знайдено)
- ✅ Код чистий і зрозумілий

---

## Рішення

Рішення знаходиться в `solutions/solution_3.go`.

Спробуйте виконати завдання самостійно перед тим, як дивитись рішення!

---

## Навчальні цілі

Після виконання цієї вправи ви будете вміти:
- Створювати абстракції через інтерфейси
- Реалізувати різні backend'и
- Використовувати Dependency Injection
- Працювати з файлами
- Тестувати через Mock

---

## Подальше вдосконалення

Подумайте як додати:
- Database Storage (SQLite, PostgreSQL)
- Redis Storage
- S3 Storage
- Компресію даних
- Шифрування
- Версіонування значень
- TTL (Time To Live) для ключів
- Bulk операції (SaveMany, LoadMany)

---

## Архітектурні патерни

Цей приклад демонструє:
- **Repository Pattern**: Storage як репозиторій
- **Strategy Pattern**: різні стратегії зберігання
- **Dependency Injection**: DataManager залежить від Storage
- **Adapter Pattern**: різні backend'и через один інтерфейс

---

## Формат файлу (простий варіант)

```
key1=value1
key2=value2
key3=value3
```

Кожен рядок: `key=value`

---

## Робота з файлами

### Запис
```go
func (f *FileStorage) saveToFile() error {
    file, err := os.Create(f.filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    for key, value := range f.data {
        fmt.Fprintf(file, "%s=%s\n", key, value)
    }
    return nil
}
```

### Читання
```go
func (f *FileStorage) loadFromFile() error {
    file, err := os.Open(f.filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        parts := strings.Split(line, "=")
        if len(parts) == 2 {
            f.data[parts[0]] = parts[1]
        }
    }
    return scanner.Err()
}
```

---

## Real-world Applications

Подібні абстракції використовуються в:
- **Web frameworks**: database abstraction
- **Cloud SDKs**: storage abstraction (S3, Azure Blob, GCS)
- **Caching**: Redis, Memcached, in-memory
- **Configuration**: файли, env vars, remote config
- **Logging**: console, file, remote

