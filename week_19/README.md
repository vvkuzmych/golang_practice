# Week 19 — Clean Architecture

**Ціль:** Зрозуміти та впровадити зрілу архітектуру з чіткими layers та SOLID принципами.

---

## 📚 Теорія

### [01. SOLID Principles](./theory/01_solid_principles.md)
- Single Responsibility Principle (SRP)
- Open/Closed Principle (OCP)
- Liskov Substitution Principle (LSP)
- Interface Segregation Principle (ISP)
- Dependency Inversion Principle (DIP)

### [02. Dependency Inversion](./theory/02_dependency_inversion.md)
- Що таке Dependency Inversion
- Інверсія залежностей vs Dependency Injection
- Приклади в Go

### [03. Clean Architecture Layers](./theory/03_clean_architecture.md)
- Domain Layer (Entities, Business Logic)
- Use Case Layer (Application Logic)
- Interface Adapters (Controllers, Presenters)
- Infrastructure (Database, HTTP, External Services)

---

## 🛠️ Практика

### [01. Refactored API with Clean Architecture](./practice/01_refactored_api/)
- Чітке розділення на layers
- Dependency Injection
- Repository pattern
- Use Cases
- SOLID principles

### [02. Hexagonal Architecture Example](./practice/02_hexagonal/)
- Ports and Adapters
- Domain-driven design
- Тестування без залежностей

---

## 📝 Exercises

### [Exercise 1: Apply SOLID to Existing Code](./exercises/exercise_1.md)
Рефакторинг коду з порушеннями SOLID.

### [Exercise 2: Implement Clean Architecture](./exercises/exercise_2.md)
Створити REST API з чітким поділом на layers.

### [Exercise 3: Add New Feature](./exercises/exercise_3.md)
Додати нову функцію без зміни існуючого коду (OCP).

---

## 🎯 Learning Outcomes

Після цього тижня ви зможете:
- ✅ Застосовувати SOLID принципи в Go коді
- ✅ Розділяти код на чіткі архітектурні layers
- ✅ Використовувати Dependency Injection
- ✅ Писати тестований та підтримуваний код
- ✅ Розуміти Clean Architecture та Hexagonal Architecture

---

## 📖 Additional Resources

- [Clean Architecture (Robert C. Martin)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [SOLID in Go](https://dave.cheney.net/2016/08/20/solid-go-design)
- [Go Domain-Driven Design](https://github.com/marcusolsson/goddd)

---

**Next:** [Week 20 — System Design](../week_20/README.md)
