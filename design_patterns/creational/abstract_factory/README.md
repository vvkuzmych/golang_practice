# Abstract Factory Pattern

## 📋 Опис

**Abstract Factory** - породжуючий патерн, який дозволяє створювати сімейства пов'язаних об'єктів без прив'язки до конкретних класів.

---

## 🎯 Проблема

- Потрібно створити **сімейство** пов'язаних об'єктів
- Різні варіанти одного сімейства (Windows UI vs Mac UI)
- Хочете гарантувати сумісність об'єктів
- Хочете приховати деталі створення

**Приклад:**  
UI елементи: Windows (Button, Checkbox) vs Mac (Button, Checkbox)

---

## ✅ Рішення

1. Створити інтерфейс Abstract Factory
2. Кожне сімейство має свою фабрику
3. Фабрика створює всі пов'язані об'єкти
4. Клієнт працює через інтерфейси

---

## 🔧 Реалізація в Go

```go
// Інтерфейси продуктів
type Button interface {
    Click() string
}

type Checkbox interface {
    Check() string
}

// Abstract Factory
type GUIFactory interface {
    CreateButton() Button
    CreateCheckbox() Checkbox
}

// Windows сімейство
type WindowsFactory struct{}

func (w *WindowsFactory) CreateButton() Button {
    return &WindowsButton{}
}

func (w *WindowsFactory) CreateCheckbox() Checkbox {
    return &WindowsCheckbox{}
}

// Mac сімейство
type MacFactory struct{}

func (m *MacFactory) CreateButton() Button {
    return &MacButton{}
}

func (m *MacFactory) CreateCheckbox() Checkbox {
    return &MacCheckbox{}
}
```

---

## ✅ Переваги

- ✅ Гарантує сумісність продуктів
- ✅ Легко додавати нові сімейства
- ✅ Ізолює конкретні класи

## ❌ Недоліки

- ❌ Складніше ніж Factory Method
- ❌ Важко додавати нові типи продуктів

---

## 🎓 Коли використовувати

✅ **Використовуйте коли:**
- Потрібні сімейства пов'язаних об'єктів
- Гарантія сумісності об'єктів
- Різні варіанти одного набору (themes, platforms)

❌ **Не використовуйте коли:**
- Всього один тип продукту
- Простий Factory Method достатньо

---

## 🌍 Реальні приклади

### В реальних проектах:
- UI frameworks (Windows/Mac/Linux themes)
- Database drivers (MySQL/PostgreSQL connections)
- Document exporters (PDF/XML/JSON generators)
- Cloud providers (AWS/Azure/GCP services)

---

## 💻 Запустити приклад

```bash
go run main.go
```

---

## 📚 Більше інформації

- [Refactoring.Guru - Abstract Factory](https://refactoring.guru/design-patterns/abstract-factory)

