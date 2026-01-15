# Deadlock в Go

## 🎯 Що таке Deadlock?

**Deadlock** — це ситуація коли всі goroutines заблоковані і немає можливості розблокуватись. Програма не може продовжити виконання і зависає.

Go runtime автоматично виявляє deadlock і викидає помилку:
```
fatal error: all goroutines are asleep - deadlock!
```

---

## ⚠️ Коли виникає Deadlock?

Deadlock виникає коли **одночасно виконуються ВСІ** умови:

1. **Всі goroutines заблоковані** — жодна goroutine не може продовжити роботу
2. **Немає зовнішніх подій** — ніщо не може розблокувати goroutines
3. **Програма не може завершитись** — main goroutine теж заблокована

---

## 📋 Типові сценарії Deadlock

### 1️⃣ Unbuffered Channel без Receiver

**Найпоширеніша помилка!**

```go
package main

func main() {
    ch := make(chan int)  // Unbuffered channel
    ch <- 42              // ❌ DEADLOCK! Блокується - ніхто не читає
}
```

**Що відбувається:**
1. `ch <- 42` намагається відправити в unbuffered channel
2. Unbuffered channel блокує sender до receiver
3. Receiver НЕ існує → main goroutine назавжди заблокована
4. Go runtime: `fatal error: all goroutines are asleep - deadlock!`

**Виправлення:**

```go
// ✅ Варіант 1: Використати buffered channel
ch := make(chan int, 1)  // Capacity = 1
ch <- 42                 // OK! Не блокує

// ✅ Варіант 2: Receiver в goroutine
ch := make(chan int)
go func() {
    value := <-ch        // Receiver готовий
    fmt.Println(value)
}()
ch <- 42                 // OK! Sender не блокується
```

---

### 2️⃣ Забули `close()` в Range Loop

```go
package main

import "fmt"

func main() {
    ch := make(chan int, 3)
    
    // Відправляємо дані
    ch <- 1
    ch <- 2
    ch <- 3
    // ❌ Забули close(ch)!
    
    // Range чекає на close()
    for v := range ch {
        fmt.Println(v)
    }  // ❌ DEADLOCK! Range ніколи не завершиться
}
```

**Що відбувається:**
1. `range ch` читає 1, 2, 3
2. `range` чекає на наступне значення АБО `close(ch)`
3. `close(ch)` НЕ викликано → `range` чекає вічно
4. Main goroutine заблокована → deadlock

**Виправлення:**

```go
ch := make(chan int, 3)

ch <- 1
ch <- 2
ch <- 3
close(ch)  // ✅ Закриваємо channel!

for v := range ch {
    fmt.Println(v)  // Прочитає 1, 2, 3 і завершиться
}
// ✅ OK! Range завершився після close()
```

**Правило:** Завжди викликайте `close(ch)` після того як відправили ВСІ дані в channel, якщо використовуєте `range`.

---

### 3️⃣ Циклічне очікування (Circular Wait)

```go
package main

func main() {
    ch1 := make(chan int)
    ch2 := make(chan int)
    
    // Goroutine 1: чекає ch2, щоб відправити в ch1
    go func() {
        ch1 <- <-ch2  // Чекає на ch2
    }()
    
    // Goroutine 2: чекає ch1, щоб відправити в ch2
    go func() {
        ch2 <- <-ch1  // Чекає на ch1
    }()
    
    // ❌ DEADLOCK! Обидві goroutines чекають одна одну
    select {}  // Чекаємо вічно
}
```

**Що відбувається:**
1. Goroutine 1: `<-ch2` блокується (ніхто не відправляє в ch2)
2. Goroutine 2: `<-ch1` блокується (ніхто не відправляє в ch1)
3. Обидві goroutines чекають одна одну → циклічне очікування
4. Deadlock!

**Виправлення:**

```go
// ✅ Варіант 1: Використати buffered channels
ch1 := make(chan int, 1)
ch2 := make(chan int, 1)

go func() {
    ch1 <- 1          // Не блокується (buffered)
    value := <-ch2
    fmt.Println(value)
}()

go func() {
    ch2 <- 2          // Не блокується (buffered)
    value := <-ch1
    fmt.Println(value)
}()

time.Sleep(100 * time.Millisecond)  // Даємо час виконатись
```

---

### 4️⃣ WaitGroup без Done()

```go
package main

import "sync"

func main() {
    var wg sync.WaitGroup
    
    wg.Add(1)
    go func() {
        // Робимо щось
        // ❌ Забули wg.Done()!
    }()
    
    wg.Wait()  // ❌ DEADLOCK! Чекає вічно
}
```

**Виправлення:**

```go
var wg sync.WaitGroup

wg.Add(1)
go func() {
    defer wg.Done()  // ✅ Завжди викликаємо Done()!
    // Робимо щось
}()

wg.Wait()  // ✅ OK!
```

**Best practice:** Завжди використовуйте `defer wg.Done()` одразу після `wg.Add()`.

---

### 5️⃣ Читання з пустого unbuffered channel

```go
package main

import "fmt"

func main() {
    ch := make(chan int)  // Unbuffered
    
    value := <-ch  // ❌ DEADLOCK! Ніхто не відправляє
    fmt.Println(value)
}
```

**Виправлення:**

```go
ch := make(chan int)

go func() {
    ch <- 42  // ✅ Sender в goroutine
}()

value := <-ch  // ✅ OK! Receiver готовий
fmt.Println(value)
```

---

### 6️⃣ Select без Default (всі cases заблоковані)

```go
package main

func main() {
    ch1 := make(chan int)
    ch2 := make(chan int)
    
    select {
    case v := <-ch1:  // Заблоковано (ніхто не відправляє)
        fmt.Println(v)
    case v := <-ch2:  // Заблоковано (ніхто не відправляє)
        fmt.Println(v)
    }  // ❌ DEADLOCK! Всі cases заблоковані
}
```

**Виправлення:**

```go
// ✅ Варіант 1: Додати default
select {
case v := <-ch1:
    fmt.Println(v)
case v := <-ch2:
    fmt.Println(v)
default:
    fmt.Println("No data available")  // ✅ Non-blocking
}

// ✅ Варіант 2: Додати timeout
select {
case v := <-ch1:
    fmt.Println(v)
case <-time.After(1 * time.Second):
    fmt.Println("Timeout!")  // ✅ Не deadlock
}
```

---

## 🛡️ Як уникнути Deadlock?

### 1. Завжди забезпечте receiver для unbuffered channel

```go
// ❌ Погано
ch := make(chan int)
ch <- 42  // Deadlock!

// ✅ Добре
ch := make(chan int)
go func() {
    <-ch  // Receiver готовий
}()
ch <- 42  // OK!
```

---

### 2. Завжди викликайте `close()` якщо використовуєте `range`

```go
// ❌ Погано
for v := range ch {
    fmt.Println(v)
}  // Deadlock якщо ch не закритий!

// ✅ Добре
close(ch)  // Закриваємо після відправки всіх даних
for v := range ch {
    fmt.Println(v)
}  // OK!
```

---

### 3. Використовуйте `defer wg.Done()` ЗАВЖДИ

```go
// ❌ Погано
wg.Add(1)
go func() {
    // робота
    wg.Done()  // Може бути пропущено через panic або return!
}()

// ✅ Добре
wg.Add(1)
go func() {
    defer wg.Done()  // Виконається завжди!
    // робота
}()
```

---

### 4. Додавайте `default` або timeout в `select`

```go
// ❌ Погано (може deadlock)
select {
case v := <-ch:
    process(v)
}

// ✅ Добре (non-blocking)
select {
case v := <-ch:
    process(v)
default:
    // Альтернативна логіка
}

// ✅ Добре (з timeout)
select {
case v := <-ch:
    process(v)
case <-time.After(5 * time.Second):
    // Timeout handling
}
```

---

### 5. Використовуйте buffered channels для асинхронної роботи

```go
// ❌ Unbuffered може deadlock
ch := make(chan int)
ch <- 42  // Deadlock!

// ✅ Buffered не deadlock
ch := make(chan int, 1)
ch <- 42  // OK! (якщо buffer не повний)
```

**Але увага:** Buffered channels не завжди вирішують deadlock! Якщо buffer заповнений, sender все одно заблокується.

---

## 🔍 Як виявити Deadlock?

### 1. Go Runtime Detection

Go runtime **автоматично** виявляє deadlock в main goroutine:

```
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan send]:
main.main()
    /path/to/file.go:10 +0x59
```

**Коли виявляється:**
- Всі goroutines заблоковані на channel operations, WaitGroup, Mutex, etc.
- Програма не може продовжити виконання

**Коли НЕ виявляється:**
- Є активні goroutines (наприклад, `time.Sleep`)
- Є network I/O або інші блокуючі операції
- Deadlock в окремих goroutines (не main)

---

### 2. Race Detector

Хоча race detector не виявляє deadlock напряму, він допомагає знайти проблеми з concurrent code:

```bash
go run -race main.go
```

---

### 3. Timeout Pattern

Використовуйте timeout для виявлення потенційних deadlock:

```go
done := make(chan bool)

go func() {
    // Потенційно блокуючий код
    result := <-ch
    done <- true
}()

select {
case <-done:
    fmt.Println("Success!")
case <-time.After(5 * time.Second):
    fmt.Println("Possible deadlock detected!")
}
```

---

## 📊 Deadlock vs Livelock vs Starvation

| Проблема | Опис | Goroutines активні? |
|----------|------|---------------------|
| **Deadlock** | Всі goroutines заблоковані, не можуть продовжити | ❌ Ні (заблоковані) |
| **Livelock** | Goroutines активні, але не роблять прогресу | ✅ Так (активні, але циклюються) |
| **Starvation** | Деякі goroutines ніколи не отримують ресурси | ⚠️ Частково (одні працюють, інші чекають) |

### Deadlock Example:
```go
ch := make(chan int)
ch <- 42  // ❌ Заблоковано назавжди
```

### Livelock Example:
```go
// Дві goroutines постійно змінюють стан, але не роблять прогресу
for {
    select {
    case <-ch1:
        ch2 <- 1  // Відправляє назад
    case <-ch2:
        ch1 <- 1  // Відправляє назад
    }
}
// ✅ Goroutines активні, але прогресу немає!
```

---

## ✅ Checklist: Як уникнути Deadlock

- [ ] Unbuffered channel має receiver перед sender?
- [ ] Викликається `close(ch)` після відправки всіх даних (для `range`)?
- [ ] Використовується `defer wg.Done()` для WaitGroup?
- [ ] `select` має `default` або timeout для non-blocking?
- [ ] Немає циклічного очікування між goroutines?
- [ ] Buffered channels мають достатній capacity?
- [ ] Кожна goroutine має шлях завершення?

---

## 🎓 Висновок

**Deadlock виникає коли:**
1. Всі goroutines заблоковані
2. Немає можливості розблокуватись
3. Програма не може продовжити

**Як уникнути:**
- Завжди забезпечуйте receiver для sender
- Закривайте channels після відправки всіх даних
- Використовуйте `defer wg.Done()`
- Додавайте `default` або timeout в `select`
- Уникайте циклічного очікування

**Пам'ятайте:** Go runtime допомагає виявити deadlock, але краще його **уникати** ніж виправляти!

---

**Наступний файл:** `05_channel_vs_queue.md` — чому channel не queue
