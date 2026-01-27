# Бази Даних: PostgreSQL, SQL, GORM

---

## 📖 Зміст

1. [PostgreSQL Basics](#1-postgresql-basics)
2. [SQL Queries](#2-sql-queries)
3. [Go database/sql](#3-go-databasesql)
4. [GORM (ORM)](#4-gorm-orm)
5. [Migrations](#5-migrations)
6. [Transactions](#6-transactions)

---

## 1. PostgreSQL Basics

### Установка (macOS)
```bash
brew install postgresql@15
brew services start postgresql@15

# Створити БД
createdb myapp_dev
```

### Підключення
```bash
psql -d myapp_dev
```

### Основні типи даних
```sql
INTEGER, BIGINT          -- Числа
VARCHAR(255), TEXT       -- Текст
BOOLEAN                  -- true/false
DATE, TIMESTAMP          -- Дати
JSON, JSONB              -- JSON дані
UUID                     -- Унікальні ID
```

---

## 2. SQL Queries

### CREATE TABLE
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    published BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_posts_user_id ON posts(user_id);
```

### INSERT
```sql
INSERT INTO users (username, email, password_hash)
VALUES ('john', 'john@example.com', '$2a$10$...');

INSERT INTO posts (user_id, title, content, published)
VALUES (1, 'My First Post', 'Hello World!', true);
```

### SELECT
```sql
-- Всі користувачі
SELECT * FROM users;

-- З умовою
SELECT * FROM users WHERE email = 'john@example.com';

-- JOIN
SELECT users.username, posts.title
FROM users
INNER JOIN posts ON users.id = posts.user_id
WHERE users.id = 1;

-- Агрегація
SELECT user_id, COUNT(*) as post_count
FROM posts
GROUP BY user_id;
```

### UPDATE
```sql
UPDATE users
SET email = 'newemail@example.com', updated_at = NOW()
WHERE id = 1;
```

### DELETE
```sql
DELETE FROM posts WHERE id = 5;
DELETE FROM users WHERE created_at < NOW() - INTERVAL '1 year';
```

---

## 3. Go database/sql

### Підключення
```go
package main

import (
    "database/sql"
    "fmt"
    
    _ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
    connStr := "postgresql://user:password@localhost/myapp_dev?sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }
    defer db.Close()
    
    // Перевірка підключення
    if err := db.Ping(); err != nil {
        panic(err)
    }
    
    fmt.Println("Connected to database!")
}
```

### Query (багато рядків)
```go
func getUsers(db *sql.DB) ([]User, error) {
    rows, err := db.Query("SELECT id, username, email FROM users")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var users []User
    for rows.Next() {
        var u User
        if err := rows.Scan(&u.ID, &u.Username, &u.Email); err != nil {
            return nil, err
        }
        users = append(users, u)
    }
    
    return users, rows.Err()
}
```

### QueryRow (один рядок)
```go
func getUserByID(db *sql.DB, id int) (*User, error) {
    var u User
    err := db.QueryRow("SELECT id, username, email FROM users WHERE id = $1", id).
        Scan(&u.ID, &u.Username, &u.Email)
    
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("user not found")
    }
    if err != nil {
        return nil, err
    }
    
    return &u, nil
}
```

### Exec (INSERT/UPDATE/DELETE)
```go
func createUser(db *sql.DB, username, email, password string) (int, error) {
    var id int
    err := db.QueryRow(
        "INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id",
        username, email, password,
    ).Scan(&id)
    
    return id, err
}

func updateUser(db *sql.DB, id int, email string) error {
    _, err := db.Exec("UPDATE users SET email = $1 WHERE id = $2", email, id)
    return err
}
```

### Prepared Statements
```go
func getUsersByIDs(db *sql.DB, ids []int) ([]User, error) {
    stmt, err := db.Prepare("SELECT id, username, email FROM users WHERE id = $1")
    if err != nil {
        return nil, err
    }
    defer stmt.Close()
    
    var users []User
    for _, id := range ids {
        var u User
        err := stmt.QueryRow(id).Scan(&u.ID, &u.Username, &u.Email)
        if err != nil && err != sql.ErrNoRows {
            return nil, err
        }
        users = append(users, u)
    }
    
    return users, nil
}
```

---

## 4. GORM (ORM)

### Установка
```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
```

### Підключення
```go
package main

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func main() {
    dsn := "host=localhost user=myuser password=mypass dbname=myapp_dev port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic(err)
    }
    
    // Auto Migration
    db.AutoMigrate(&User{}, &Post{})
}
```

### Models
```go
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Username  string    `gorm:"uniqueIndex;not null"`
    Email     string    `gorm:"uniqueIndex;not null"`
    Password  string    `gorm:"not null"`
    Posts     []Post    `gorm:"foreignKey:UserID"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Post struct {
    ID        uint   `gorm:"primaryKey"`
    UserID    uint   `gorm:"index;not null"`
    User      User   `gorm:"foreignKey:UserID"`
    Title     string `gorm:"not null"`
    Content   string `gorm:"type:text"`
    Published bool   `gorm:"default:false"`
    CreatedAt time.Time
}
```

### CRUD Operations

**Create:**
```go
// Створити користувача
user := User{
    Username: "john",
    Email:    "john@example.com",
    Password: "hashedpassword",
}
result := db.Create(&user)
fmt.Println("ID:", user.ID) // ID автоматично встановлюється
```

**Read:**
```go
// Знайти по ID
var user User
db.First(&user, 1) // SELECT * FROM users WHERE id = 1;

// Знайти по умові
db.Where("email = ?", "john@example.com").First(&user)

// Знайти всіх
var users []User
db.Find(&users) // SELECT * FROM users;

// З умовами
db.Where("created_at > ?", time.Now().AddDate(0, -1, 0)).Find(&users)
```

**Update:**
```go
// Оновити одне поле
db.Model(&user).Update("Email", "newemail@example.com")

// Оновити кілька полів
db.Model(&user).Updates(User{
    Email:    "newemail@example.com",
    Username: "john_updated",
})

// Або через map
db.Model(&user).Updates(map[string]interface{}{
    "email":    "newemail@example.com",
    "username": "john_updated",
})
```

**Delete:**
```go
// Soft delete (якщо є DeletedAt поле)
db.Delete(&user, 1)

// Permanent delete
db.Unscoped().Delete(&user, 1)
```

### Associations (Relationships)

**Preload (eager loading):**
```go
// Завантажити користувача з його постами
var user User
db.Preload("Posts").First(&user, 1)

fmt.Println(user.Username)
for _, post := range user.Posts {
    fmt.Println(post.Title)
}
```

**Create with associations:**
```go
user := User{
    Username: "jane",
    Email:    "jane@example.com",
    Posts: []Post{
        {Title: "First Post", Content: "Hello"},
        {Title: "Second Post", Content: "World"},
    },
}
db.Create(&user) // Створить user і 2 posts
```

### Advanced Queries

**Where conditions:**
```go
// AND
db.Where("username = ? AND email = ?", "john", "john@example.com").First(&user)

// OR
db.Where("username = ?", "john").Or("email = ?", "john@example.com").Find(&users)

// IN
db.Where("id IN ?", []int{1, 2, 3}).Find(&users)

// LIKE
db.Where("email LIKE ?", "%@gmail.com").Find(&users)
```

**Order, Limit, Offset:**
```go
db.Order("created_at DESC").Limit(10).Offset(20).Find(&users)
```

**Count:**
```go
var count int64
db.Model(&User{}).Where("created_at > ?", yesterday).Count(&count)
```

**Raw SQL:**
```go
db.Raw("SELECT * FROM users WHERE email = ?", "john@example.com").Scan(&user)
```

---

## 5. Migrations

### golang-migrate

**Установка:**
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

**Створення міграції:**
```bash
migrate create -ext sql -dir db/migrations -seq create_users_table
```

**Міграції:**
```sql
-- 000001_create_users_table.up.sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 000001_create_users_table.down.sql
DROP TABLE users;
```

**Запуск:**
```bash
migrate -path db/migrations -database "postgresql://localhost/myapp_dev?sslmode=disable" up
migrate -path db/migrations -database "postgresql://localhost/myapp_dev?sslmode=disable" down
```

---

## 6. Transactions

### database/sql
```go
func transferMoney(db *sql.DB, fromID, toID int, amount float64) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback() // Rollback якщо не зробили Commit
    
    // Зняти гроші
    _, err = tx.Exec("UPDATE accounts SET balance = balance - $1 WHERE id = $2", amount, fromID)
    if err != nil {
        return err
    }
    
    // Додати гроші
    _, err = tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, toID)
    if err != nil {
        return err
    }
    
    // Commit транзакції
    return tx.Commit()
}
```

### GORM
```go
func transferMoney(db *gorm.DB, fromID, toID uint, amount float64) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // Зняти гроші
        if err := tx.Model(&Account{}).Where("id = ?", fromID).
            Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
            return err
        }
        
        // Додати гроші
        if err := tx.Model(&Account{}).Where("id = ?", toID).
            Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
            return err
        }
        
        return nil // Commit
    })
}
```

---

## ✅ Best Practices

1. **Завжди використовуйте параметризовані запити** ($1, $2) - захист від SQL injection
2. **Використовуйте transactions** для критичних операцій
3. **Індекси** на часто використовувані колонки (WHERE, JOIN)
4. **Connection pooling** - налаштуйте правильно
5. **Migrations** - версіонуйте зміни БД
6. **Backups** - регулярні бекапи
7. **Monitoring** - слідкуйте за slow queries

```go
// Connection pool налаштування
sqlDB, _ := db.DB()
sqlDB.SetMaxOpenConns(25)
sqlDB.SetMaxIdleConns(5)
sqlDB.SetConnMaxLifetime(5 * time.Minute)
```

---

**Далі:** [06_networking.md](./06_networking.md)
