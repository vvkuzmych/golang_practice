package main

import (
	"fmt"
	"strings"
)

// ============= Example 1: House Builder =============

type House struct {
	Walls     int
	Doors     int
	Windows   int
	HasRoof   bool
	HasGarage bool
	HasGarden bool
	Floors    int
}

type HouseBuilder struct {
	house House
}

func NewHouseBuilder() *HouseBuilder {
	return &HouseBuilder{
		house: House{Floors: 1}, // default
	}
}

func (b *HouseBuilder) SetWalls(n int) *HouseBuilder {
	b.house.Walls = n
	return b
}

func (b *HouseBuilder) SetDoors(n int) *HouseBuilder {
	b.house.Doors = n
	return b
}

func (b *HouseBuilder) SetWindows(n int) *HouseBuilder {
	b.house.Windows = n
	return b
}

func (b *HouseBuilder) WithRoof() *HouseBuilder {
	b.house.HasRoof = true
	return b
}

func (b *HouseBuilder) WithGarage() *HouseBuilder {
	b.house.HasGarage = true
	return b
}

func (b *HouseBuilder) WithGarden() *HouseBuilder {
	b.house.HasGarden = true
	return b
}

func (b *HouseBuilder) SetFloors(n int) *HouseBuilder {
	b.house.Floors = n
	return b
}

func (b *HouseBuilder) Build() House {
	return b.house
}

func (h House) String() string {
	return fmt.Sprintf("🏠 House: %d walls, %d doors, %d windows, %d floors"+
		"\n   Roof: %v, Garage: %v, Garden: %v",
		h.Walls, h.Doors, h.Windows, h.Floors,
		h.HasRoof, h.HasGarage, h.HasGarden)
}

// ============= Example 2: HTTP Request Builder =============

type HTTPRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    string
	Timeout int
}

type RequestBuilder struct {
	request HTTPRequest
}

func NewRequestBuilder(url string) *RequestBuilder {
	return &RequestBuilder{
		request: HTTPRequest{
			URL:     url,
			Method:  "GET",
			Headers: make(map[string]string),
			Timeout: 30,
		},
	}
}

func (b *RequestBuilder) Method(method string) *RequestBuilder {
	b.request.Method = method
	return b
}

func (b *RequestBuilder) Header(key, value string) *RequestBuilder {
	b.request.Headers[key] = value
	return b
}

func (b *RequestBuilder) Body(body string) *RequestBuilder {
	b.request.Body = body
	return b
}

func (b *RequestBuilder) Timeout(seconds int) *RequestBuilder {
	b.request.Timeout = seconds
	return b
}

func (b *RequestBuilder) Build() HTTPRequest {
	return b.request
}

func (r HTTPRequest) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("🌐 %s %s\n", r.Method, r.URL))
	for k, v := range r.Headers {
		sb.WriteString(fmt.Sprintf("   %s: %s\n", k, v))
	}
	if r.Body != "" {
		sb.WriteString(fmt.Sprintf("   Body: %s\n", r.Body))
	}
	sb.WriteString(fmt.Sprintf("   Timeout: %ds", r.Timeout))
	return sb.String()
}

// ============= Example 3: SQL Query Builder =============

type SQLQuery struct {
	Table   string
	Columns []string
	Where   []string
	OrderBy string
	Limit   int
}

type QueryBuilder struct {
	query SQLQuery
}

func NewQueryBuilder(table string) *QueryBuilder {
	return &QueryBuilder{
		query: SQLQuery{
			Table:   table,
			Columns: []string{"*"},
		},
	}
}

func (b *QueryBuilder) Select(columns ...string) *QueryBuilder {
	b.query.Columns = columns
	return b
}

func (b *QueryBuilder) Where(condition string) *QueryBuilder {
	b.query.Where = append(b.query.Where, condition)
	return b
}

func (b *QueryBuilder) OrderBy(column string) *QueryBuilder {
	b.query.OrderBy = column
	return b
}

func (b *QueryBuilder) Limit(n int) *QueryBuilder {
	b.query.Limit = n
	return b
}

func (b *QueryBuilder) Build() string {
	var sb strings.Builder

	sb.WriteString("SELECT ")
	sb.WriteString(strings.Join(b.query.Columns, ", "))
	sb.WriteString(" FROM ")
	sb.WriteString(b.query.Table)

	if len(b.query.Where) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(b.query.Where, " AND "))
	}

	if b.query.OrderBy != "" {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(b.query.OrderBy)
	}

	if b.query.Limit > 0 {
		sb.WriteString(fmt.Sprintf(" LIMIT %d", b.query.Limit))
	}

	return sb.String()
}

// ============= Main =============

func main() {
	fmt.Println("╔════════════════════════════════════════════════╗")
	fmt.Println("║          Builder Pattern Demo                  ║")
	fmt.Println("╚════════════════════════════════════════════════╝")

	// ===== Example 1: House Builder =====
	fmt.Println("\n🔹 Example 1: Building Houses")
	fmt.Println(strings.Repeat("─", 50))

	// Маленький будинок
	smallHouse := NewHouseBuilder().
		SetWalls(4).
		SetDoors(1).
		SetWindows(2).
		WithRoof().
		Build()

	fmt.Println("\nSmall house:")
	fmt.Println(smallHouse)

	// Великий будинок
	bigHouse := NewHouseBuilder().
		SetWalls(8).
		SetDoors(3).
		SetWindows(12).
		SetFloors(2).
		WithRoof().
		WithGarage().
		WithGarden().
		Build()

	fmt.Println("\nBig house:")
	fmt.Println(bigHouse)

	// ===== Example 2: HTTP Request Builder =====
	fmt.Println("\n\n🔹 Example 2: Building HTTP Requests")
	fmt.Println(strings.Repeat("─", 50))

	// GET request
	getReq := NewRequestBuilder("https://api.example.com/users").
		Header("Authorization", "Bearer token123").
		Header("Accept", "application/json").
		Build()

	fmt.Println("\nGET Request:")
	fmt.Println(getReq)

	// POST request
	postReq := NewRequestBuilder("https://api.example.com/users").
		Method("POST").
		Header("Content-Type", "application/json").
		Header("Authorization", "Bearer token123").
		Body(`{"name":"John","email":"john@example.com"}`).
		Timeout(60).
		Build()

	fmt.Println("\n\nPOST Request:")
	fmt.Println(postReq)

	// ===== Example 3: SQL Query Builder =====
	fmt.Println("\n\n🔹 Example 3: Building SQL Queries")
	fmt.Println(strings.Repeat("─", 50))

	// Простий SELECT
	query1 := NewQueryBuilder("users").Build()
	fmt.Printf("\nQuery 1: %s\n", query1)

	// SELECT з умовами
	query2 := NewQueryBuilder("users").
		Select("name", "email", "age").
		Where("age > 18").
		Where("active = true").
		OrderBy("name").
		Limit(10).
		Build()
	fmt.Printf("\nQuery 2: %s\n", query2)

	// Складний запит
	query3 := NewQueryBuilder("orders").
		Select("id", "customer_name", "total_amount", "created_at").
		Where("status = 'completed'").
		Where("total_amount > 100").
		Where("created_at > '2024-01-01'").
		OrderBy("created_at DESC").
		Limit(50).
		Build()
	fmt.Printf("\nQuery 3: %s\n", query3)

	// ===== Comparison: With vs Without Builder =====
	fmt.Println("\n\n🔹 Comparison: Constructor vs Builder")
	fmt.Println(strings.Repeat("─", 50))

	fmt.Println("\n❌ Without Builder (ugly):")
	fmt.Println("house := House{")
	fmt.Println("    Walls: 4,")
	fmt.Println("    Doors: 2,")
	fmt.Println("    Windows: 6,")
	fmt.Println("    HasRoof: true,")
	fmt.Println("    HasGarage: false,  // чи це optional?")
	fmt.Println("    HasGarden: true,")
	fmt.Println("    Floors: 2,")
	fmt.Println("}")

	fmt.Println("\n✅ With Builder (clean):")
	fmt.Println("house := NewHouseBuilder().")
	fmt.Println("    SetWalls(4).")
	fmt.Println("    SetDoors(2).")
	fmt.Println("    SetWindows(6).")
	fmt.Println("    SetFloors(2).")
	fmt.Println("    WithRoof().")
	fmt.Println("    WithGarden().")
	fmt.Println("    Build()")

	// ===== Summary =====
	fmt.Println("\n\n📝 ВИСНОВКИ")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println("✅ Зручно для об'єктів з багатьма параметрами")
	fmt.Println("✅ Fluent interface (method chaining)")
	fmt.Println("✅ Чіткий код, легко читати")
	fmt.Println("✅ Optional параметри явно видно")
	fmt.Println("✅ Immutable результат")

	fmt.Println("\n💡 ВИКОРИСТАННЯ:")
	fmt.Println("   - HTTP requests/responses")
	fmt.Println("   - SQL query builders")
	fmt.Println("   - Configuration objects")
	fmt.Println("   - Complex data structures")

	fmt.Println("\n📚 Go stdlib приклади:")
	fmt.Println("   - strings.Builder")
	fmt.Println("   - bytes.Buffer")
}
