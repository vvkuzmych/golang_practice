# Design Patterns in Go

## 📊 Статус: 23/23 (100%) ✅

Повна колекція класичних Design Patterns реалізованих на Go!

---

## 🎨 Creational Patterns (5/5) ✅

Патерни створення об'єктів

| Pattern | Опис | Use Case |
|---------|------|----------|
| [Abstract Factory](creational/abstract_factory/) | Створення сімей пов'язаних об'єктів | UI themes, DB drivers |
| [Builder](creational/builder/) | Покрокове створення складних об'єктів | Query builders, Config |
| [Factory](creational/factory/) | Створення об'єктів через фабричний метод | Loggers, Parsers |
| [Prototype](creational/prototype/) | Клонування об'єктів | Deep copy, Caching |
| [Singleton](creational/singleton/) | Один екземпляр класу | Config, DB connection |

---

## 🏗️ Structural Patterns (7/7) ✅

Патерни компонування об'єктів

| Pattern | Опис | Use Case |
|---------|------|----------|
| [Adapter](structural/adapter/) | Адаптація інтерфейсів | Legacy code, 3rd party libs |
| [Bridge](structural/bridge/) | Розділення абстракції і реалізації | Cross-platform, Devices |
| [Composite](structural/composite/) | Деревоподібні структури | File system, UI, Org chart |
| [Decorator](structural/decorator/) | Динамічне додавання функціональності | Middleware, Wrappers |
| [Facade](structural/facade/) | Спрощений інтерфейс до складної системи | API Gateway, Subsystems |
| [Flyweight](structural/flyweight/) | Зменшення використання пам'яті | Object pooling, Particles |
| [Proxy](structural/proxy/) | Контроль доступу до об'єкта | Caching, Lazy loading, Auth |

---

## 🎭 Behavioral Patterns (11/11) ✅

Патерни взаємодії між об'єктами

| Pattern | Опис | Use Case |
|---------|------|----------|
| [Chain of Responsibility](behavioral/chain_of_responsibility/) | Ланцюжок обробників | Middleware, Validation pipeline |
| [Command](behavioral/command/) | Інкапсуляція запиту як об'єкта | Undo/Redo, Task queue |
| [Iterator](behavioral/iterator/) | Послідовний доступ до елементів | Collections traversal |
| [Mediator](behavioral/mediator/) | Зменшення зв'язаності компонентів | Chat rooms, UI coordination |
| [Memento](behavioral/memento/) | Збереження стану об'єкта | Undo/Redo, Snapshots |
| [Observer](behavioral/observer/) | Підписка на події | Event systems, Pub/Sub |
| [State](behavioral/state/) | Зміна поведінки залежно від стану | **ATM**, Order workflow |
| [Strategy](behavioral/strategy/) | Вибір алгоритму в runtime | Sorting, Payment methods |
| [Template Method](behavioral/template_method/) | Скелет алгоритму в базовому класі | Data processing, Testing |
| [Visitor](behavioral/visitor/) | Додавання операцій без зміни класів | AST traversal, Reporting |

---

## 🔗 Зв'язок з Week 7

### State Pattern = ATM State Machine! 🏧

**State Pattern** (`behavioral/state/`) - це той самий паттерн, який використовується для ATM у Week 7!

```
IDLE → CARD_INSERTED → AUTHORIZED → DISPENSING → COMPLETED
```

Детальніше:
- `week_7/HARDWARE_STATE_MACHINE.md` - ATM з hardware events
- `week_7/theory/17_hardware_software_integration.md` - Повна інтеграція
- `design_patterns/behavioral/state/` - Класичний State Pattern

---

## 📚 Як використовувати

### Структура кожного патерну

```
pattern_name/
├── main.go      # Робочий приклад з main()
└── README.md    # Опис, use cases, коли використовувати
```

### Запуск будь-якого патерну

```bash
cd design_patterns/<category>/<pattern_name>
go run main.go
```

### Приклад

```bash
cd design_patterns/behavioral/state
go run main.go
```

---

## 🎯 Коли використовувати який патерн?

### Потрібно створити об'єкт?
→ **Creational** (Factory, Builder, Singleton)

### Потрібно організувати структуру?
→ **Structural** (Composite, Proxy, Decorator)

### Потрібна взаємодія між об'єктами?
→ **Behavioral** (State, Strategy, Observer)

---

## 💡 Найпопулярніші в Go

### Top 10 (за частотою використання)

1. **Factory** ⭐⭐⭐⭐⭐ - Створення об'єктів
2. **Builder** ⭐⭐⭐⭐⭐ - Конфігурація
3. **Singleton** ⭐⭐⭐⭐⭐ - Глобальні об'єкти
4. **Decorator** ⭐⭐⭐⭐⭐ - HTTP middleware
5. **Strategy** ⭐⭐⭐⭐⭐ - Різні алгоритми
6. **Observer** ⭐⭐⭐⭐ - Event systems
7. **Proxy** ⭐⭐⭐⭐ - Caching, logging
8. **Chain of Responsibility** ⭐⭐⭐⭐ - Middleware
9. **State** ⭐⭐⭐ - State machines (ATM!)
10. **Composite** ⭐⭐⭐ - Tree structures

---

## 📖 Додаткові матеріали

### Теорія

- `week_6/theory/02_design_patterns.md` - Детальна теорія всіх патернів

### Practical Examples

- Week 6: OOP + Design Patterns
- Week 7: State Machine для ATM
- Sneakers Marketplace: Factory, Builder, Strategy, Observer

---

## 🎓 Рекомендована послідовність вивчення

### Рівень 1: Basics (почни з цих)
1. Singleton
2. Factory
3. Strategy
4. Observer

### Рівень 2: Intermediate
5. Builder
6. Decorator
7. Proxy
8. State

### Рівень 3: Advanced
9. Abstract Factory
10. Composite
11. Chain of Responsibility
12. Template Method

### Рівень 4: Specialized (рідко використовуються)
13. Flyweight
14. Visitor
15. Mediator
16. Memento

---

## ✅ Завдання для практики

### 1. HTTP Server з Middleware (Chain + Decorator)
Створи HTTP server з:
- Logging middleware
- Auth middleware
- Rate limiting middleware

### 2. Document Editor (Memento + Command)
Створи редактор з:
- Undo/Redo
- Command history

### 3. Notification System (Observer + Strategy)
Створи систему нотифікацій з:
- Email, SMS, Push subscribers
- Різні стратегії доставки

### 4. Game AI (State + Strategy)
Створи game character з:
- States: Idle, Attacking, Defending
- Strategies: Aggressive, Defensive, Balanced

---

## 🔍 Пошук патерну по проблемі

| Проблема | Патерн |
|----------|--------|
| Потрібен один екземпляр | Singleton |
| Складна конфігурація | Builder |
| Різні реалізації одного інтерфейсу | Factory, Strategy |
| Додати функціональність без зміни коду | Decorator, Proxy |
| Обхід колекції | Iterator |
| Pub/Sub система | Observer |
| Middleware pipeline | Chain of Responsibility |
| State machine (ATM) | State |
| Деревоподібна структура | Composite |
| Undo/Redo | Memento, Command |

---

## 📊 Progress

- ✅ Creational: 5/5 (100%)
- ✅ Structural: 7/7 (100%)
- ✅ Behavioral: 11/11 (100%)

**Total: 23/23 (100%) COMPLETE!** 🎉

---

## 🚀 Наступні кроки

1. ✅ ~~Створити всі патерни~~ **ГОТОВО!**
2. Вивчити кожен патерн (запустити приклади)
3. Виконати практичні завдання
4. Використати в реальних проектах

---

**Created:** 2026-01-28  
**Status:** Complete ✅  
**Author:** Week 6 & Week 7 Integration
