# Вправа 2: Shape Interface

## Ціль
Створити систему геометричних фігур з використанням інтерфейсів і polymorphism.

---

## Завдання

Створіть програму `shapes.go`, яка:

1. Має інтерфейс `Shape` з методами для обчислення площі та периметру
2. Реалізує кілька типів фігур
3. Використовує polymorphism для роботи з різними фігурами
4. Підтримує різні операції над колекцією фігур

---

## Вимоги

### Інтерфейс Shape
```go
type Shape interface {
    Area() float64
    Perimeter() float64
    Name() string
}
```

### Обов'язкові фігури:
1. **Rectangle** (прямокутник)
   - Поля: `Width`, `Height`
   
2. **Circle** (коло)
   - Поле: `Radius`
   
3. **Square** (квадрат)
   - Поле: `Side`

### Додаткові функції:
- `PrintShapeInfo(s Shape)` - вивести інформацію про фігуру
- `TotalArea(shapes []Shape) float64` - загальна площа всіх фігур
- `LargestShape(shapes []Shape) Shape` - знайти найбільшу фігуру за площею
- `FilterByMinArea(shapes []Shape, minArea float64) []Shape` - фільтрувати фігури

---

## Приклад використання

```go
func main() {
    rect := Rectangle{Width: 10, Height: 5}
    circle := Circle{Radius: 7}
    square := Square{Side: 6}
    
    shapes := []Shape{rect, circle, square}
    
    for _, shape := range shapes {
        PrintShapeInfo(shape)
    }
    
    fmt.Printf("Total area: %.2f\n", TotalArea(shapes))
    
    largest := LargestShape(shapes)
    fmt.Printf("Largest: %s\n", largest.Name())
}
```

---

## Підказки

### 1. Rectangle
```go
type Rectangle struct {
    Width  float64
    Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

func (r Rectangle) Name() string {
    return "Rectangle"
}
```

### 2. Circle
```go
import "math"

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}
```

### 3. Робота з slice через інтерфейс
```go
func TotalArea(shapes []Shape) float64 {
    total := 0.0
    for _, shape := range shapes {
        total += shape.Area()
    }
    return total
}
```

---

## Очікуваний вивід

```
=== Shape Information ===

Rectangle
  Width: 10.00, Height: 5.00
  Area: 50.00
  Perimeter: 30.00

Circle
  Radius: 7.00
  Area: 153.94
  Perimeter: 43.98

Square
  Side: 6.00
  Area: 36.00
  Perimeter: 24.00

=== Statistics ===
Total shapes: 3
Total area: 239.94
Average area: 79.98
Largest shape: Circle (153.94)

=== Filter: Area > 40 ===
Circle: 153.94
Rectangle: 50.00
```

---

## Бонус завдання

1. **Додаткові фігури**:
   - Triangle (трикутник)
   - Pentagon (п'ятикутник)
   - Ellipse (еліпс)

2. **Colorable Interface**:
   ```go
   type Colorable interface {
       SetColor(color string)
       GetColor() string
   }
   
   type ColoredCircle struct {
       Circle
       color string
   }
   ```

3. **Drawable Interface**:
   ```go
   type Drawable interface {
       Draw() string
   }
   
   func (r Rectangle) Draw() string {
       return "Drawing a rectangle"
   }
   ```

4. **Comparison**:
   ```go
   type Comparable interface {
       IsLargerThan(other Shape) bool
   }
   ```

5. **Sorting**:
   ```go
   func SortByArea(shapes []Shape) {
       // Відсортувати за площею
   }
   ```

---

## Критерії оцінки

- ✅ Інтерфейс `Shape` правильно оголошений
- ✅ Всі фігури реалізують інтерфейс (неявно)
- ✅ Функції працюють з будь-якою фігурою через інтерфейс
- ✅ Обчислення площі та периметру коректні
- ✅ Використовується polymorphism
- ✅ Код чистий і зрозумілий

---

## Рішення

Рішення знаходиться в `solutions/solution_2.go`.

Спробуйте виконати завдання самостійно перед тим, як дивитись рішення!

---

## Навчальні цілі

Після виконання цієї вправи ви будете вміти:
- Оголошувати та використовувати інтерфейси
- Реалізувати інтерфейси неявно
- Використовувати polymorphism
- Працювати зі slice інтерфейсів
- Фільтрувати та сортувати через інтерфейси

---

## Подальше вдосконалення

Подумайте як додати:
- 3D фігури (Cube, Sphere, Cylinder)
- Метод `Volume()` для 3D інтерфейсу
- Візуалізацію фігур (ASCII art)
- Трансформації (масштабування, поворот)
- Перетин фігур (collision detection)

---

## Математичні формули

### Коло
- Площа: π × r²
- Периметр: 2 × π × r

### Прямокутник
- Площа: width × height
- Периметр: 2 × (width + height)

### Квадрат
- Площа: side²
- Периметр: 4 × side

### Трикутник (Формула Герона)
- s = (a + b + c) / 2
- Площа: √(s × (s-a) × (s-b) × (s-c))
- Периметр: a + b + c

