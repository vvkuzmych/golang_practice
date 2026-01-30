-- Basic JOINs Practice - Queries

-- ========================================
-- 1. INNER JOIN
-- ========================================

-- Користувачі з їх замовленнями (тільки хто має замовлення)
SELECT 
    u.name AS customer,
    o.id AS order_id,
    o.total,
    o.status
FROM users u
INNER JOIN orders o ON u.id = o.user_id
ORDER BY u.name, o.id;

/*
 customer   | order_id | total   | status
------------+----------+---------+-----------
 Jane Smith | 3        | 375.00  | completed
 John Doe   | 1        | 1300.00 | completed
 John Doe   | 2        | 25.00   | pending

 Пропущено: Bob (немає замовлень), Order 4 (guest)
*/

-- ========================================
-- 2. LEFT JOIN
-- ========================================

-- ВСІ користувачі + їх замовлення (якщо є)
SELECT 
    u.name AS customer,
    o.id AS order_id,
    o.total,
    o.status
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
ORDER BY u.name, o.id;

/*
 customer   | order_id | total   | status
------------+----------+---------+-----------
 Bob Wilson | NULL     | NULL    | NULL        ← Немає замовлень
 Jane Smith | 3        | 375.00  | completed
 John Doe   | 1        | 1300.00 | completed
 John Doe   | 2        | 25.00   | pending
*/

-- Знайти користувачів БЕЗ замовлень
SELECT u.name
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
WHERE o.id IS NULL;

/*
 name
------------
 Bob Wilson
*/

-- ========================================
-- 3. RIGHT JOIN
-- ========================================

-- ВСІ замовлення + користувачі (якщо є)
SELECT 
    u.name AS customer,
    o.id AS order_id,
    o.total,
    o.status
FROM users u
RIGHT JOIN orders o ON u.id = o.user_id
ORDER BY o.id;

/*
 customer   | order_id | total   | status
------------+----------+---------+---------
 John Doe   | 1        | 1300.00 | completed
 John Doe   | 2        | 25.00   | pending
 Jane Smith | 3        | 375.00  | completed
 NULL       | 4        | 75.00   | pending    ← Guest order
*/

-- Знайти гостьові замовлення (без user_id)
SELECT o.*
FROM users u
RIGHT JOIN orders o ON u.id = o.user_id
WHERE u.id IS NULL;

/*
 id | user_id | total | status  | created_at
----+---------+-------+---------+-------------
 4  | NULL    | 75.00 | pending | ...
*/

-- ========================================
-- 4. FULL OUTER JOIN
-- ========================================

-- ВСІ користувачі + ВСІ замовлення
SELECT 
    u.name AS customer,
    o.id AS order_id,
    o.total
FROM users u
FULL OUTER JOIN orders o ON u.id = o.user_id
ORDER BY u.name, o.id;

/*
 customer   | order_id | total
------------+----------+---------
 NULL       | 4        | 75.00    ← Guest order
 Bob Wilson | NULL     | NULL     ← No orders
 Jane Smith | 3        | 375.00
 John Doe   | 1        | 1300.00
 John Doe   | 2        | 25.00
*/

-- ========================================
-- 5. Multiple JOINs
-- ========================================

-- Користувачі → Замовлення → Товари в замовленні
SELECT 
    u.name AS customer,
    o.id AS order_id,
    o.total AS order_total,
    p.name AS product,
    oi.quantity,
    oi.price AS item_price
FROM users u
INNER JOIN orders o ON u.id = o.user_id
INNER JOIN order_items oi ON o.id = oi.order_id
INNER JOIN products p ON oi.product_id = p.id
ORDER BY u.name, o.id, p.name;

/*
 customer   | order_id | order_total | product  | quantity | item_price
------------+----------+-------------+----------+----------+------------
 Jane Smith | 3        | 375.00      | Keyboard | 1        | 75.00
 Jane Smith | 3        | 375.00      | Monitor  | 1        | 300.00
 John Doe   | 1        | 1300.00     | Laptop   | 1        | 1200.00
 John Doe   | 1        | 1300.00     | Mouse    | 4        | 25.00
 John Doe   | 2        | 25.00       | Mouse    | 1        | 25.00
*/

-- ========================================
-- 6. Aggregations with JOINs
-- ========================================

-- Кількість замовлень на користувача
SELECT 
    u.name,
    COUNT(o.id) AS order_count,
    COALESCE(SUM(o.total), 0) AS total_spent
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
GROUP BY u.id, u.name
ORDER BY total_spent DESC;

/*
 name        | order_count | total_spent
-------------+-------------+-------------
 John Doe    | 2           | 1325.00
 Jane Smith  | 1           | 375.00
 Bob Wilson  | 0           | 0.00
*/

-- TOP-3 товари за кількістю продажів
SELECT 
    p.name AS product,
    SUM(oi.quantity) AS total_sold,
    SUM(oi.quantity * oi.price) AS revenue
FROM products p
INNER JOIN order_items oi ON p.id = oi.product_id
GROUP BY p.id, p.name
ORDER BY total_sold DESC
LIMIT 3;

/*
 product  | total_sold | revenue
----------+------------+---------
 Mouse    | 5          | 125.00
 Keyboard | 2          | 150.00
 Monitor  | 1          | 300.00
*/

-- ========================================
-- 7. Subqueries with JOINs
-- ========================================

-- Користувачі, які купили більше ніж на $500
SELECT DISTINCT u.name
FROM users u
INNER JOIN (
    SELECT user_id, SUM(total) AS total_spent
    FROM orders
    GROUP BY user_id
    HAVING SUM(total) > 500
) AS big_spenders ON u.id = big_spenders.user_id;

/*
 name
----------
 John Doe
*/

-- ========================================
-- 8. CROSS JOIN (Cartesian Product)
-- ========================================

-- Всі комбінації користувачів і товарів (рідко потрібно)
SELECT 
    u.name AS user,
    p.name AS product
FROM users u
CROSS JOIN products p
ORDER BY u.name, p.name
LIMIT 10;

/*
 user        | product
-------------+---------
 Bob Wilson  | Keyboard
 Bob Wilson  | Laptop
 Bob Wilson  | Monitor
 Bob Wilson  | Mouse
 Jane Smith  | Keyboard
 Jane Smith  | Laptop
 Jane Smith  | Monitor
 Jane Smith  | Mouse
 John Doe    | Keyboard
 John Doe    | Laptop
*/

-- ========================================
-- 9. SELF JOIN
-- ========================================

-- Створимо таблицю співробітників з менеджерами
CREATE TEMP TABLE employees (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    manager_id INT
);

INSERT INTO employees (name, manager_id) VALUES
    ('Alice', NULL),  -- CEO
    ('Bob', 1),       -- Manager: Alice
    ('Carol', 1),     -- Manager: Alice
    ('Dave', 2),      -- Manager: Bob
    ('Eve', 2);       -- Manager: Bob

-- Співробітники з їх менеджерами
SELECT 
    e.name AS employee,
    m.name AS manager
FROM employees e
LEFT JOIN employees m ON e.manager_id = m.id
ORDER BY e.id;

/*
 employee | manager
----------+---------
 Alice    | NULL     ← CEO (no manager)
 Bob      | Alice
 Carol    | Alice
 Dave     | Bob
 Eve      | Bob
*/

-- ========================================
-- 10. Advanced: Multiple LEFT JOINs
-- ========================================

-- Всі користувачі + їх останнє замовлення + товари в ньому
SELECT 
    u.name AS customer,
    last_order.id AS last_order_id,
    last_order.created_at,
    p.name AS product
FROM users u
LEFT JOIN LATERAL (
    SELECT * FROM orders o
    WHERE o.user_id = u.id
    ORDER BY o.created_at DESC
    LIMIT 1
) last_order ON true
LEFT JOIN order_items oi ON last_order.id = oi.order_id
LEFT JOIN products p ON oi.product_id = p.id
ORDER BY u.name;

/*
 customer   | last_order_id | created_at | product
------------+---------------+------------+----------
 Bob Wilson | NULL          | NULL       | NULL
 Jane Smith | 3             | ...        | Monitor
 Jane Smith | 3             | ...        | Keyboard
 John Doe   | 2             | ...        | Mouse
*/
