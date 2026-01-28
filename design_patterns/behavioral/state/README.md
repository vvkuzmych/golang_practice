# State Pattern (Стан)

## 📋 Опис

**State** - поведінковий патерн, що дозволяє об'єкту змінювати свою поведінку залежно від внутрішнього стану.

## 🎯 Проблема

- Об'єкт поводиться по-різному в різних станах
- Багато if/switch statements для різних станів
- Складно додавати нові стани

## ✅ Рішення

Винести кожен стан в окремий клас з єдиним інтерфейсом.

## 🏗️ Структура

```
Context (ATM)
    ↓
State (interface)
    ├─> IdleState
    ├─> CardInsertedState  
    ├─> AuthorizedState
    └─> DispensingState
```

## 💡 Зв'язок з Week 7!

**Це той самий State Machine pattern для ATM!** 🏧

```
IDLE → CARD_INSERTED → AUTHORIZED → DISPENSING → IDLE
```

Така сама логіка, як у файлі:
`week_7/theory/17_hardware_software_integration.md`

## ✅ Переваги

- Чистий код (без if/switch)
- Single Responsibility
- Open/Closed Principle
- Легко додавати нові стани

## 🎯 Коли використовувати

✅ **ATM transactions** (як у Week 7!)  
✅ Document workflow (draft → review → published)  
✅ Order processing (pending → paid → shipped)  
✅ Game character states  
✅ Connection states (disconnected → connecting → connected)  

## 📊 Real-World приклади

- **ATM:** IDLE → DISPENSING → COMPLETED
- **TCP Connection:** CLOSED → SYN_SENT → ESTABLISHED
- **Order:** PENDING → PROCESSING → SHIPPED → DELIVERED

## 🚀 Запуск

```bash
cd design_patterns/behavioral/state
go run main.go
```

## 🔗 Пов'язані файли

- `week_7/HARDWARE_STATE_MACHINE.md` - ATM State Machine
- `week_7/theory/17_hardware_software_integration.md` - Hardware + State Machine
