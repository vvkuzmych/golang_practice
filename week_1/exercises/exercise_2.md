# Вправа 2: Інформація про студента

## Ціль
Створити програму для роботи зі структурами та демонстрації zero values.

---

## Завдання

Створіть програму `student_info.go`, яка:

1. Визначає структуру `Student`
2. Створює кілька студентів різними способами
3. Демонструє zero values
4. Виводить інформацію у різних форматах

---

## Структура Student

```go
type Student struct {
    FirstName  string
    LastName   string
    Age        int
    GPA        float64
    IsActive   bool
    University *string  // pointer для демонстрації nil
}
```

---

## Завдання

### Частина 1: Створення студентів

Створіть 4 студенти різними способами:

1. **student1**: Повна ініціалізація з усіма полями
2. **student2**: Частко ініціалізація (тільки ім'я та прізвище)
3. **student3**: Zero value (без ініціалізації)
4. **student4**: З указкою університету (pointer)

### Частина 2: Вивід

Виведіть кожного студента з використанням:
- `%v` - default format
- `%+v` - з іменами полів
- `%#v` - Go syntax
- `%T` - тип

### Частина 3: Аналіз Zero Values

Створіть таблицю з zero values для всіх полів структури.

---

## Приклад виводу

```
=== Студенти ===

Студент 1 (повна ініціалізація):
  %v:  {Іван Петренко 20 3.8 true <nil>}
  %+v: {FirstName:Іван LastName:Петренко Age:20 GPA:3.8 IsActive:true University:<nil>}
  %T:  main.Student

Студент 2 (часткова ініціалізація):
  %v:  {Марія Коваленко 0 0 false <nil>}
  %+v: {FirstName:Марія LastName:Коваленко Age:0 GPA:0 IsActive:false University:<nil>}
  
Студент 3 (zero value):
  %v:  { 0 0 false <nil>}
  %+v: {FirstName: LastName: Age:0 GPA:0 IsActive:false University:<nil>}

Студент 4 (з університетом):
  FirstName:  Петро
  LastName:   Сидоренко
  University: КНУ імені Тараса Шевченка

=== Таблиця Zero Values ===

Поле         | Тип     | Zero Value
-------------|---------|------------
FirstName    | string  | ""
LastName     | string  | ""
Age          | int     | 0
GPA          | float64 | 0.0
IsActive     | bool    | false
University   | *string | nil
```

---

## Вимоги

1. ✅ Визначити структуру Student
2. ✅ Створити 4 екземпляри різними способами
3. ✅ Використати var та :=
4. ✅ Продемонструвати роботу з pointer полем
5. ✅ Використати всі формати виводу (%v, %+v, %#v, %T)
6. ✅ Створити таблицю zero values
7. ✅ Додати коментарі до коду

---

## Підказки

### 1. Структура з pointer
```go
type Student struct {
    University *string
}

// Використання
uni := "КНУ"
student := Student{
    University: &uni,
}

// Перевірка на nil
if student.University != nil {
    fmt.Println(*student.University)
}
```

### 2. Zero value структури
```go
var student Student  // всі поля мають zero values
```

### 3. Часткова ініціалізація
```go
student := Student{
    FirstName: "Іван",
    LastName:  "Петренко",
    // інші поля отримують zero values
}
```

---

## Бонус завдання

1. **Методи на структурі**:
   ```go
   func (s Student) FullName() string {
       return s.FirstName + " " + s.LastName
   }
   
   func (s Student) IsExcellent() bool {
       return s.GPA >= 4.0
   }
   ```

2. **Функція створення**:
   ```go
   func NewStudent(firstName, lastName string, age int) Student {
       return Student{
           FirstName: firstName,
           LastName:  lastName,
           Age:       age,
           IsActive:  true,
       }
   }
   ```

3. **Порівняння студентів**: Функція для порівняння GPA

4. **Валідація**: Функція для перевірки коректності даних

---

## Критерії оцінки

- ✅ Структура визначена правильно
- ✅ 4 студенти створені різними способами
- ✅ Всі формати виводу використані
- ✅ Zero values продемонстровані
- ✅ Робота з pointer полем
- ✅ Код чистий з коментарями

---

## Рішення

Рішення знаходиться в `solutions/solution_2.go`.

Спробуйте виконати завдання самостійно!

