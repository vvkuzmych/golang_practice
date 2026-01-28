# RESTful APIs Best Practices

## API Design Principles

### 1. Use Proper HTTP Methods

```go
// ✅ GOOD
GET    /api/v1/users           // List
GET    /api/v1/users/:id       // Get one
POST   /api/v1/users           // Create
PUT    /api/v1/users/:id       // Update (full)
PATCH  /api/v1/users/:id       // Update (partial)
DELETE /api/v1/users/:id       // Delete
```

### 2. API Versioning

```go
// URL versioning
router.Group("/api/v1")
router.Group("/api/v2")

// Header versioning
r.Header.Get("API-Version")
```

### 3. Proper Status Codes

```go
// Success
200 OK                    // GET, PUT, PATCH
201 Created              // POST
204 No Content           // DELETE

// Client errors
400 Bad Request          // Invalid input
401 Unauthorized         // No auth
403 Forbidden            // Insufficient permissions
404 Not Found
409 Conflict            // Duplicate resource

// Server errors
500 Internal Server Error
503 Service Unavailable
```

### 4. Response Format

```go
type Response struct {
    Data    interface{} `json:"data,omitempty"`
    Error   *ErrorInfo  `json:"error,omitempty"`
    Meta    *Meta       `json:"meta,omitempty"`
}

type ErrorInfo struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

type Meta struct {
    Page      int `json:"page"`
    PerPage   int `json:"per_page"`
    Total     int `json:"total"`
}
```

## Pagination

```go
func ListUsers(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
    
    offset := (page - 1) * perPage
    
    users, total, err := service.List(offset, perPage)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, Response{
        Data: users,
        Meta: &Meta{
            Page:    page,
            PerPage: perPage,
            Total:   total,
        },
    })
}
```

## Filtering & Sorting

```go
// GET /api/v1/users?status=active&sort=created_at:desc

type Query struct {
    Status  string `form:"status"`
    Sort    string `form:"sort"`
    Page    int    `form:"page"`
    PerPage int    `form:"per_page"`
}

func ListUsers(c *gin.Context) {
    var q Query
    if err := c.ShouldBindQuery(&q); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    users := service.List(q)
    c.JSON(200, users)
}
```

## OpenAPI / Swagger

```go
// @Summary      List users
// @Description  Get all users with pagination
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        page     query  int  false  "Page number"
// @Param        per_page query  int  false  "Items per page"
// @Success      200  {array}  User
// @Failure      500  {object}  ErrorResponse
// @Router       /users [get]
func ListUsers(c *gin.Context) {
    // ...
}
```
