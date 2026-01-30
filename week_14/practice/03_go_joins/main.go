package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Models
type User struct {
	ID    int
	Name  string
	Email string
}

type Order struct {
	ID     int
	UserID *int // NULL для гостьових замовлень
	Total  float64
	Status string
}

type UserWithOrders struct {
	UserName   string
	OrderID    *int     // NULL якщо немає замовлень
	OrderTotal *float64 // NULL якщо немає замовлень
	Status     *string  // NULL якщо немає замовлень
}

type OrderDetails struct {
	CustomerName string
	OrderID      int
	OrderTotal   float64
	ProductName  string
	Quantity     int
	ItemPrice    float64
}

func main() {
	// Підключення до PostgreSQL
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=joins_practice sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ Connected to PostgreSQL\n")

	// 1. INNER JOIN - Користувачі з замовленнями
	innerJoin(db)

	// 2. LEFT JOIN - Всі користувачі + замовлення
	leftJoin(db)

	// 3. Find users without orders
	usersWithoutOrders(db)

	// 4. Multiple JOINs - Order details
	orderDetails(db)

	// 5. Aggregation - Order count per user
	orderCountPerUser(db)
}

// 1. INNER JOIN - Тільки користувачі з замовленнями
func innerJoin(db *sql.DB) {
	fmt.Println("1️⃣ INNER JOIN - Users with orders")
	fmt.Println("─────────────────────────────────────")

	query := `
		SELECT 
			u.name,
			o.id AS order_id,
			o.total,
			o.status
		FROM users u
		INNER JOIN orders o ON u.id = o.user_id
		ORDER BY u.name, o.id
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var orderID int
		var total float64
		var status string

		if err := rows.Scan(&name, &orderID, &total, &status); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%-15s | Order #%-2d | $%-7.2f | %s\n", name, orderID, total, status)
	}
	fmt.Println()
}

// 2. LEFT JOIN - Всі користувачі (навіть без замовлень)
func leftJoin(db *sql.DB) {
	fmt.Println("2️⃣ LEFT JOIN - All users + their orders")
	fmt.Println("─────────────────────────────────────")

	query := `
		SELECT 
			u.name,
			o.id AS order_id,
			o.total,
			o.status
		FROM users u
		LEFT JOIN orders o ON u.id = o.user_id
		ORDER BY u.name, o.id
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var uo UserWithOrders

		if err := rows.Scan(&uo.UserName, &uo.OrderID, &uo.OrderTotal, &uo.Status); err != nil {
			log.Fatal(err)
		}

		if uo.OrderID == nil {
			fmt.Printf("%-15s | No orders\n", uo.UserName)
		} else {
			fmt.Printf("%-15s | Order #%-2d | $%-7.2f | %s\n",
				uo.UserName, *uo.OrderID, *uo.OrderTotal, *uo.Status)
		}
	}
	fmt.Println()
}

// 3. Знайти користувачів БЕЗ замовлень
func usersWithoutOrders(db *sql.DB) {
	fmt.Println("3️⃣ Users WITHOUT orders (LEFT JOIN + WHERE IS NULL)")
	fmt.Println("─────────────────────────────────────")

	query := `
		SELECT u.name
		FROM users u
		LEFT JOIN orders o ON u.id = o.user_id
		WHERE o.id IS NULL
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("❌ %s - No orders\n", name)
	}
	fmt.Println()
}

// 4. Multiple JOINs - Деталі замовлення
func orderDetails(db *sql.DB) {
	fmt.Println("4️⃣ Multiple JOINs - Order details")
	fmt.Println("─────────────────────────────────────")

	query := `
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
		ORDER BY u.name, o.id, p.name
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var currentOrder int
	for rows.Next() {
		var od OrderDetails

		if err := rows.Scan(
			&od.CustomerName,
			&od.OrderID,
			&od.OrderTotal,
			&od.ProductName,
			&od.Quantity,
			&od.ItemPrice,
		); err != nil {
			log.Fatal(err)
		}

		if od.OrderID != currentOrder {
			if currentOrder != 0 {
				fmt.Println()
			}
			fmt.Printf("📦 Order #%d - %s (Total: $%.2f)\n",
				od.OrderID, od.CustomerName, od.OrderTotal)
			currentOrder = od.OrderID
		}

		fmt.Printf("   • %s x%d @ $%.2f = $%.2f\n",
			od.ProductName, od.Quantity, od.ItemPrice, float64(od.Quantity)*od.ItemPrice)
	}
	fmt.Println()
}

// 5. Aggregation - Кількість замовлень на користувача
func orderCountPerUser(db *sql.DB) {
	fmt.Println("5️⃣ Aggregation - Order count per user")
	fmt.Println("─────────────────────────────────────")

	query := `
		SELECT 
			u.name,
			COUNT(o.id) AS order_count,
			COALESCE(SUM(o.total), 0) AS total_spent
		FROM users u
		LEFT JOIN orders o ON u.id = o.user_id
		GROUP BY u.id, u.name
		ORDER BY total_spent DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Printf("%-15s | Orders | Total Spent\n", "User")
	fmt.Println("────────────────┼────────┼────────────")

	for rows.Next() {
		var name string
		var orderCount int
		var totalSpent float64

		if err := rows.Scan(&name, &orderCount, &totalSpent); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%-15s | %-6d | $%.2f\n", name, orderCount, totalSpent)
	}
	fmt.Println()
}
