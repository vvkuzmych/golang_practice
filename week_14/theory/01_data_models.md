# Data Models (Моделі даних)

## 🎯 Що таке Data Model?

**Data Model** - структура організації даних та відношень між ними.

---

## 📊 Типи відношень

### 1. One-to-One (1:1)

```
User ←→ Profile
```

**Приклад:**
```
users                  profiles
─────────────         ────────────────
id | name             id | user_id | bio
1  | John             1  | 1       | "Developer"
2  | Jane             2  | 2       | "Designer"
```

**SQL:**
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100)
);

CREATE TABLE profiles (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE,  -- UNIQUE забезпечує 1:1
    bio TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

---

### 2. One-to-Many (1:N)

```
User ──< Posts
(один користувач → багато постів)
```

**Приклад:**
```
users                  posts
─────────────         ────────────────────────
id | name             id | user_id | title
1  | John             1  | 1       | "Post 1"
2  | Jane             2  | 1       | "Post 2"
                       3  | 2       | "Post 3"
```

**SQL:**
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100)
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,  -- Без UNIQUE = Many
    title VARCHAR(255),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

---

### 3. Many-to-Many (N:M)

```
Students >──< Courses
(багато студентів → багато курсів)
```

**Приклад з Junction Table:**
```
students              enrollments           courses
────────────         ─────────────────     ────────────
id | name            student_id | course_id  id | name
1  | John            1          | 1          1  | Math
2  | Jane            1          | 2          2  | Physics
                     2          | 1          3  | Chemistry
                     2          | 3
```

**SQL:**
```sql
CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100)
);

CREATE TABLE courses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100)
);

-- Junction table (зв'язкова таблиця)
CREATE TABLE enrollments (
    student_id INT,
    course_id INT,
    enrolled_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (student_id, course_id),  -- Composite PK
    FOREIGN KEY (student_id) REFERENCES students(id),
    FOREIGN KEY (course_id) REFERENCES courses(id)
);
```

---

## 🗄️ Реальний приклад: E-commerce

```
┌──────────┐
│  users   │
└────┬─────┘
     │ 1:N
     ▼
┌──────────┐     N:M     ┌──────────┐
│  orders  │◄───────────►│ products │
└────┬─────┘  order_items└──────────┘
     │ 1:N
     ▼
┌──────────┐
│ order_   │
│ items    │
└──────────┘
```

**Schema:**
```sql
-- Users
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE,
    name VARCHAR(100)
);

-- Products
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    price DECIMAL(10,2)
);

-- Orders
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    total DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Order Items (Junction для Orders + Products)
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL,
    price DECIMAL(10,2),  -- Зберігаємо ціну на момент замовлення
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);
```

---

## 🎯 Keys (Ключі)

### Primary Key (PK)

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,  -- Унікальний ідентифікатор
    email VARCHAR(255)
);
```

**Правила:**
- Унікальний
- NOT NULL
- Один на таблицю

### Foreign Key (FK)

```sql
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

**Правила:**
- Посилається на PK іншої таблиці
- Забезпечує referential integrity

### Composite Key

```sql
CREATE TABLE enrollments (
    student_id INT,
    course_id INT,
    PRIMARY KEY (student_id, course_id)  -- Обидва разом = PK
);
```

---

## 📊 Normalization (Нормалізація)

### Denormalized (❌ погано)

```sql
orders
─────────────────────────────────────────────
id | user_name | user_email | product_name | price
1  | John      | j@g.com    | Laptop       | 1000
2  | John      | j@g.com    | Mouse        | 20    ← Дублювання!
```

### Normalized (✅ добре)

```sql
users                orders               products
───────────────     ──────────────      ──────────────
id | name | email   id | user_id       id | name | price
1  | John | j@g.com 1  | 1             1  | Laptop | 1000
                     2  | 1             2  | Mouse  | 20
```

**Переваги:**
- Немає дублювання
- Легше оновлювати
- Економія пам'яті

---

## 🎯 Схема відношень

### Text Diagram

```
users (1) ──< (N) posts
  │
  │ (1)
  ▼
  (1) profiles

posts (N) >──< (M) tags
         (через post_tags)
```

### SQL з відношеннями

```sql
-- 1:1 (User → Profile)
users.id ←→ profiles.user_id (UNIQUE)

-- 1:N (User → Posts)
users.id ←→ posts.user_id

-- N:M (Posts → Tags через post_tags)
posts.id ←→ post_tags.post_id
tags.id  ←→ post_tags.tag_id
```

---

## ✅ Best Practices

### 1. Завжди використовуй FK

```sql
-- ✅ GOOD
FOREIGN KEY (user_id) REFERENCES users(id)
```

### 2. Index на FK

```sql
CREATE INDEX idx_posts_user_id ON posts(user_id);
```

### 3. ON DELETE/UPDATE

```sql
FOREIGN KEY (user_id) REFERENCES users(id)
    ON DELETE CASCADE      -- Видалити пости при видаленні юзера
    ON UPDATE CASCADE      -- Оновити при зміні id
```

**Опції:**
- `CASCADE` - каскадне видалення/оновлення
- `SET NULL` - встановити NULL
- `RESTRICT` - заборонити (default)
- `NO ACTION` - нічого не робити

### 4. Naming Convention

```sql
-- Таблиці: множина, lowercase
users, posts, order_items

-- FK: singular_id
user_id, post_id, order_id

-- Junction: table1_table2
post_tags, student_courses
```

---

## 🎓 Висновок

### Відношення:

✅ **1:1** - UNIQUE FK  
✅ **1:N** - FK без UNIQUE  
✅ **N:M** - Junction table з двома FK  

### Ключі:

✅ **PK** - унікальний ідентифікатор  
✅ **FK** - посилання на інший PK  
✅ **Composite** - комбінація полів  

**Далі:** `02_sql_joins.md` - JOINs для з'єднання таблиць
