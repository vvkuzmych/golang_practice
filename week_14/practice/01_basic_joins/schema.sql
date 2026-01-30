-- Basic JOINs Practice - Schema Setup

-- Drop tables if exist
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS users;

-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Products table
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    stock INT DEFAULT 0
);

-- Orders table
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT,  -- NULL = guest order
    total DECIMAL(10,2) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Order Items table (N:M между Orders и Products)
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL,
    price DECIMAL(10,2) NOT NULL,  -- Price at order time
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT
);

-- Indexes for performance
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_product_id ON order_items(product_id);

-- Insert sample data
INSERT INTO users (name, email) VALUES
    ('John Doe', 'john@example.com'),
    ('Jane Smith', 'jane@example.com'),
    ('Bob Wilson', 'bob@example.com');

INSERT INTO products (name, price, stock) VALUES
    ('Laptop', 1200.00, 10),
    ('Mouse', 25.00, 50),
    ('Keyboard', 75.00, 30),
    ('Monitor', 300.00, 15);

INSERT INTO orders (user_id, total, status) VALUES
    (1, 1300.00, 'completed'),  -- John's order
    (1, 25.00, 'pending'),      -- John's order
    (2, 375.00, 'completed'),   -- Jane's order
    (NULL, 75.00, 'pending');   -- Guest order

INSERT INTO order_items (order_id, product_id, quantity, price) VALUES
    -- Order 1: John bought Laptop + Mouse
    (1, 1, 1, 1200.00),
    (1, 2, 4, 25.00),
    -- Order 2: John bought Mouse
    (2, 2, 1, 25.00),
    -- Order 3: Jane bought Monitor + Keyboard
    (3, 4, 1, 300.00),
    (3, 3, 1, 75.00),
    -- Order 4: Guest bought Keyboard
    (4, 3, 1, 75.00);
