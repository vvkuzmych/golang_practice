# Week 6: Швидкий Старт

## 🚀 За 5 хвилин

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_6

# 1. Прочитайте теорію
cat theory/01_oop_principles.md

# 2. Запустіть практичні приклади
go run practice/01_oop/main.go
go run practice/02_http_server/main.go

# 3. Виконайте вправу
cat exercises/exercise_1.md
```

---

## 📚 Рекомендований порядок

### День 1: ООП і Патерни
```bash
# Теорія
cat theory/01_oop_principles.md
cat theory/02_design_patterns.md

# Практика
go run practice/01_oop/main.go

# Вправа
cat exercises/exercise_1.md
```

### День 2: HTTP і Сервери
```bash
# Теорія
cat theory/03_net_http.md

# Практика
go run practice/02_http_server/main.go

# Тестуємо API
curl http://localhost:8080/api/users
```

### День 3: Мікросервіси
```bash
# Теорія
cat theory/04_microservices.md

# Практика (запустіть в різних терміналах)
go run practice/03_microservices/service_a/main.go
go run practice/03_microservices/service_b/main.go
go run practice/03_microservices/gateway/main.go
```

### День 4: Бази даних
```bash
# Теорія
cat theory/05_databases.md

# Практика
go run practice/04_database/main.go
```

### День 5: Нетворкінг
```bash
# Теорія
cat theory/06_networking.md

# Практика
go run practice/05_networking/tcp_server.go  # Термінал 1
go run practice/05_networking/tcp_client.go  # Термінал 2
```

### День 6: Goroutines і Конкурентність
```bash
# Теорія
cat theory/07_goroutines_concurrency.md

# Практика
go run practice/06_goroutines/main.go

# Перевірка race conditions
go run -race practice/06_goroutines/main.go
```

---

## ✅ Перевірте своє розуміння

- [ ] Розумію 4 принципи ООП в Go
- [ ] Можу створити HTTP сервер з routing
- [ ] Знаю різницю між монолітом і мікросервісами
- [ ] Вмію працювати з PostgreSQL через GORM
- [ ] Розумію різницю між TCP і UDP
- [ ] Вмію налаштовувати timeouts і retries
- [ ] Використовую goroutines для конкурентності
- [ ] Розумію channels і select
- [ ] Знаю sync.Mutex і sync.WaitGroup

---

## 🔗 Корисні команди

```bash
# Встановити залежності
go get -u gorm.io/gorm
go get -u github.com/gorilla/mux
go get -u github.com/lib/pq

# Запустити тести
go test ./...

# Форматування
go fmt ./...

# Лінтер
golangci-lint run
```

---

**Успіхів у навчанні!** 🎉
