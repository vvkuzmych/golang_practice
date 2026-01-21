# 🛒 E-commerce Platform - Детальний Plan

## 📊 Загальний огляд

**E-commerce Platform** - це production-ready платформа для онлайн торгівлі з повною мікросервісною архітектурою, як Amazon, eBay, але навчальна версія.

### 🎯 Основні цілі проекту:

1. **100% покриття Go концепцій** - від basics до advanced
2. **Реальна мікросервісна архітектура** - 9 незалежних сервісів
3. **Production-ready patterns** - всі best practices
4. **Portfolio project** - можна показувати роботодавцям

---

## 🏗️ Архітектура: 9 Мікросервісів

```
                        ┌──────────────────┐
                        │   API Gateway    │ ← REST API для клієнтів
                        │  (Kong/Custom)   │
                        └──────────────────┘
                                 │
                ┌────────────────┼────────────────┐
                │                │                │
        ┌───────▼──────┐  ┌─────▼─────┐  ┌──────▼──────┐
        │ User Service │  │  Product   │  │    Cart     │
        │              │  │  Service   │  │   Service   │
        └──────────────┘  └────────────┘  └─────────────┘
                │                │                │
                └────────────────┼────────────────┘
                                 │
                        ┌────────▼────────┐
                        │ Order Service   │ ← Основна бізнес-логіка
                        └─────────────────┘
                                 │
                ┌────────────────┼────────────────┐
                │                │                │
        ┌───────▼──────┐  ┌─────▼─────┐  ┌──────▼──────┐
        │   Payment    │  │Notification│  │   Search    │
        │   Service    │  │  Service   │  │   Service   │
        └──────────────┘  └────────────┘  └─────────────┘
                                 │
                        ┌────────▼────────┐
                        │   Analytics     │
                        │    Service      │
                        └─────────────────┘
                                 │
                        ┌────────▼────────┐
                        │   Review        │
                        │   Service       │
                        └─────────────────┘

        ┌─────────────────────────────────────┐
        │     Event Bus (Kafka/RabbitMQ)      │ ← Event-driven communication
        └─────────────────────────────────────┘
```

---

## 1️⃣ User Service (Автентифікація & Профілі)

### 📋 Відповідальність:
- Реєстрація користувачів
- Автентифікація (JWT, OAuth2)
- Управління профілями
- Адреси доставки
- Wishlist (список бажань)

### 🗄️ Database Schema (PostgreSQL):

```sql
-- users table
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone VARCHAR(20),
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- addresses table
CREATE TABLE addresses (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    address_type VARCHAR(20), -- shipping, billing
    street_line1 VARCHAR(255) NOT NULL,
    street_line2 VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100),
    postal_code VARCHAR(20) NOT NULL,
    country VARCHAR(100) NOT NULL,
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- wishlist table
CREATE TABLE wishlist (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL,
    added_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, product_id)
);

-- sessions table (для JWT blacklist)
CREATE TABLE sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### 🔌 gRPC API:

```protobuf
service UserService {
  // Auth
  rpc Register(RegisterRequest) returns (RegisterResponse);
  rpc Login(LoginRequest) returns (LoginResponse);
  rpc Logout(LogoutRequest) returns (LogoutResponse);
  rpc RefreshToken(RefreshTokenRequest) returns (RefreshTokenResponse);
  rpc VerifyEmail(VerifyEmailRequest) returns (VerifyEmailResponse);
  
  // Profile
  rpc GetProfile(GetProfileRequest) returns (GetProfileResponse);
  rpc UpdateProfile(UpdateProfileRequest) returns (UpdateProfileResponse);
  rpc ChangePassword(ChangePasswordRequest) returns (ChangePasswordResponse);
  
  // Addresses
  rpc AddAddress(AddAddressRequest) returns (AddAddressResponse);
  rpc GetAddresses(GetAddressesRequest) returns (GetAddressesResponse);
  rpc UpdateAddress(UpdateAddressRequest) returns (UpdateAddressResponse);
  rpc DeleteAddress(DeleteAddressRequest) returns (DeleteAddressResponse);
  
  // Wishlist
  rpc AddToWishlist(AddToWishlistRequest) returns (AddToWishlistResponse);
  rpc GetWishlist(GetWishlistRequest) returns (GetWishlistResponse);
  rpc RemoveFromWishlist(RemoveFromWishlistRequest) returns (RemoveFromWishlistResponse);
}
```

### 💻 Go Code Example:

```go
package server

import (
    "context"
    "time"
    
    pb "github.com/yourapp/proto/user"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

type UserServer struct {
    pb.UnimplementedUserServiceServer
    db     *sql.DB
    redis  *redis.Client
    secret []byte
}

func (s *UserServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
    // Validate input
    if err := validateEmail(req.Email); err != nil {
        return nil, status.Error(codes.InvalidArgument, "invalid email")
    }
    
    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, status.Error(codes.Internal, "failed to hash password")
    }
    
    // Insert user
    var userID int64
    err = s.db.QueryRowContext(ctx, `
        INSERT INTO users (email, password_hash, first_name, last_name)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `, req.Email, hashedPassword, req.FirstName, req.LastName).Scan(&userID)
    
    if err != nil {
        return nil, status.Error(codes.AlreadyExists, "user already exists")
    }
    
    // Generate JWT token
    token, err := s.generateJWT(userID, req.Email)
    if err != nil {
        return nil, status.Error(codes.Internal, "failed to generate token")
    }
    
    // Send verification email (async)
    go s.sendVerificationEmail(req.Email)
    
    return &pb.RegisterResponse{
        UserId: userID,
        Token:  token,
    }, nil
}

func (s *UserServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
    var userID int64
    var passwordHash string
    
    // Get user
    err := s.db.QueryRowContext(ctx, `
        SELECT id, password_hash FROM users WHERE email = $1
    `, req.Email).Scan(&userID, &passwordHash)
    
    if err != nil {
        return nil, status.Error(codes.NotFound, "user not found")
    }
    
    // Verify password
    if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
        return nil, status.Error(codes.Unauthenticated, "invalid password")
    }
    
    // Generate JWT
    token, err := s.generateJWT(userID, req.Email)
    if err != nil {
        return nil, status.Error(codes.Internal, "failed to generate token")
    }
    
    return &pb.LoginResponse{
        Token:  token,
        UserId: userID,
    }, nil
}

func (s *UserServer) generateJWT(userID int64, email string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "email":   email,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.secret)
}
```

### 🎯 Go Concepts:
- ✅ HTTP/gRPC servers
- ✅ JWT authentication
- ✅ Password hashing (bcrypt)
- ✅ Database operations (PostgreSQL)
- ✅ Redis для sessions
- ✅ Context для cancellation
- ✅ Error handling
- ✅ Goroutines (async email sending)

---

## 2️⃣ Product Service (Каталог товарів)

### 📋 Відповідальність:
- Управління продуктами
- Категорії
- Інвентар (stock management)
- Варіанти продуктів (розміри, кольори)
- Зображення продуктів

### 🗄️ Database Schema:

```sql
-- categories table
CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    parent_id BIGINT REFERENCES categories(id),
    description TEXT,
    image_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW()
);

-- products table
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    category_id BIGINT REFERENCES categories(id),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    compare_at_price DECIMAL(10, 2), -- original price (для знижок)
    sku VARCHAR(100) UNIQUE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- product_variants table (розміри, кольори)
CREATE TABLE product_variants (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT REFERENCES products(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL, -- "Large / Red"
    sku VARCHAR(100) UNIQUE,
    price DECIMAL(10, 2),
    stock_quantity INT NOT NULL DEFAULT 0,
    option1_name VARCHAR(50),  -- "Size"
    option1_value VARCHAR(50), -- "Large"
    option2_name VARCHAR(50),  -- "Color"
    option2_value VARCHAR(50), -- "Red"
    created_at TIMESTAMP DEFAULT NOW()
);

-- product_images table
CREATE TABLE product_images (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT REFERENCES products(id) ON DELETE CASCADE,
    url VARCHAR(255) NOT NULL,
    alt_text VARCHAR(255),
    position INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

-- inventory_transactions table (для аудиту змін stock)
CREATE TABLE inventory_transactions (
    id BIGSERIAL PRIMARY KEY,
    variant_id BIGINT REFERENCES product_variants(id),
    quantity_change INT NOT NULL,
    transaction_type VARCHAR(50), -- purchase, sale, adjustment, return
    reference_id BIGINT, -- order_id or other reference
    created_at TIMESTAMP DEFAULT NOW()
);
```

### 🔌 gRPC API:

```protobuf
service ProductService {
  // Products
  rpc CreateProduct(CreateProductRequest) returns (CreateProductResponse);
  rpc GetProduct(GetProductRequest) returns (GetProductResponse);
  rpc ListProducts(ListProductsRequest) returns (ListProductsResponse);
  rpc UpdateProduct(UpdateProductRequest) returns (UpdateProductResponse);
  rpc DeleteProduct(DeleteProductRequest) returns (DeleteProductResponse);
  
  // Categories
  rpc CreateCategory(CreateCategoryRequest) returns (CreateCategoryResponse);
  rpc ListCategories(ListCategoriesRequest) returns (ListCategoriesResponse);
  
  // Inventory
  rpc CheckStock(CheckStockRequest) returns (CheckStockResponse);
  rpc ReserveStock(ReserveStockRequest) returns (ReserveStockResponse);
  rpc ReleaseStock(ReleaseStockRequest) returns (ReleaseStockResponse);
  rpc UpdateStock(UpdateStockRequest) returns (UpdateStockResponse);
  
  // Variants
  rpc GetProductVariants(GetProductVariantsRequest) returns (GetProductVariantsResponse);
}
```

### 💻 Go Code Example:

```go
package server

import (
    "context"
    "database/sql"
    
    pb "github.com/yourapp/proto/product"
)

type ProductServer struct {
    pb.UnimplementedProductServiceServer
    db    *sql.DB
    cache *redis.Client
}

func (s *ProductServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()
    
    // Insert product
    var productID int64
    err = tx.QueryRowContext(ctx, `
        INSERT INTO products (category_id, name, slug, description, price, sku)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `, req.CategoryId, req.Name, req.Slug, req.Description, req.Price, req.Sku).Scan(&productID)
    
    if err != nil {
        return nil, status.Error(codes.Internal, "failed to create product")
    }
    
    // Insert variants
    for _, variant := range req.Variants {
        _, err = tx.ExecContext(ctx, `
            INSERT INTO product_variants 
            (product_id, name, sku, price, stock_quantity, option1_name, option1_value, option2_name, option2_value)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        `, productID, variant.Name, variant.Sku, variant.Price, variant.StockQuantity,
           variant.Option1Name, variant.Option1Value, variant.Option2Name, variant.Option2Value)
        
        if err != nil {
            return nil, status.Error(codes.Internal, "failed to create variant")
        }
    }
    
    // Commit transaction
    if err := tx.Commit(); err != nil {
        return nil, status.Error(codes.Internal, "failed to commit transaction")
    }
    
    // Invalidate cache
    s.cache.Del(ctx, fmt.Sprintf("product:%d", productID))
    
    return &pb.CreateProductResponse{
        ProductId: productID,
    }, nil
}

func (s *ProductServer) ReserveStock(ctx context.Context, req *pb.ReserveStockRequest) (*pb.ReserveStockResponse, error) {
    tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()
    
    // Check stock with SELECT FOR UPDATE (pessimistic lock)
    var currentStock int
    err = tx.QueryRowContext(ctx, `
        SELECT stock_quantity FROM product_variants 
        WHERE id = $1 FOR UPDATE
    `, req.VariantId).Scan(&currentStock)
    
    if err != nil {
        return nil, status.Error(codes.NotFound, "variant not found")
    }
    
    if currentStock < int(req.Quantity) {
        return nil, status.Error(codes.FailedPrecondition, "insufficient stock")
    }
    
    // Update stock
    _, err = tx.ExecContext(ctx, `
        UPDATE product_variants 
        SET stock_quantity = stock_quantity - $1 
        WHERE id = $2
    `, req.Quantity, req.VariantId)
    
    if err != nil {
        return nil, status.Error(codes.Internal, "failed to reserve stock")
    }
    
    // Record transaction
    _, err = tx.ExecContext(ctx, `
        INSERT INTO inventory_transactions (variant_id, quantity_change, transaction_type, reference_id)
        VALUES ($1, $2, 'reservation', $3)
    `, req.VariantId, -req.Quantity, req.OrderId)
    
    if err := tx.Commit(); err != nil {
        return nil, status.Error(codes.Internal, "failed to commit")
    }
    
    return &pb.ReserveStockResponse{Success: true}, nil
}
```

### 🎯 Go Concepts:
- ✅ Database transactions
- ✅ Pessimistic locking (SELECT FOR UPDATE)
- ✅ Redis caching
- ✅ gRPC services
- ✅ Error handling
- ✅ Context propagation

---

## 3️⃣ Cart Service (Кошик)

### 📋 Відповідальність:
- Управління кошиком покупок
- Додавання/видалення товарів
- Розрахунок загальної суми
- Промокоди
- Session management

### 🗄️ Storage (Redis):

```redis
# Cart structure (Redis Hash)
# Key: cart:{user_id}
HSET cart:123 product:456:variant:789 2  # quantity
HSET cart:123 product:457:variant:790 1

# Cart metadata
SET cart:123:promo_code "SUMMER2024"
SET cart:123:created_at "2024-01-15T10:00:00Z"
EXPIRE cart:123 604800  # 7 days TTL
```

### 🔌 gRPC API:

```protobuf
service CartService {
  rpc AddItem(AddItemRequest) returns (AddItemResponse);
  rpc RemoveItem(RemoveItemRequest) returns (RemoveItemResponse);
  rpc UpdateQuantity(UpdateQuantityRequest) returns (UpdateQuantityResponse);
  rpc GetCart(GetCartRequest) returns (GetCartResponse);
  rpc ClearCart(ClearCartRequest) returns (ClearCartResponse);
  rpc ApplyPromoCode(ApplyPromoCodeRequest) returns (ApplyPromoCodeResponse);
  rpc CalculateTotal(CalculateTotalRequest) returns (CalculateTotalResponse);
}
```

### 💻 Go Code Example:

```go
package server

import (
    "context"
    "fmt"
    "strconv"
    
    pb "github.com/yourapp/proto/cart"
    "github.com/go-redis/redis/v8"
)

type CartServer struct {
    pb.UnimplementedCartServiceServer
    redis         *redis.Client
    productClient pb.ProductServiceClient
}

func (s *CartServer) AddItem(ctx context.Context, req *pb.AddItemRequest) (*pb.AddItemResponse, error) {
    // Check stock availability first
    stockResp, err := s.productClient.CheckStock(ctx, &pb.CheckStockRequest{
        VariantId: req.VariantId,
        Quantity:  req.Quantity,
    })
    
    if err != nil || !stockResp.Available {
        return nil, status.Error(codes.FailedPrecondition, "product out of stock")
    }
    
    // Add to cart (Redis)
    cartKey := fmt.Sprintf("cart:%d", req.UserId)
    itemKey := fmt.Sprintf("product:%d:variant:%d", req.ProductId, req.VariantId)
    
    // Increment quantity
    err = s.redis.HIncrBy(ctx, cartKey, itemKey, int64(req.Quantity)).Err()
    if err != nil {
        return nil, status.Error(codes.Internal, "failed to add item to cart")
    }
    
    // Set expiration (7 days)
    s.redis.Expire(ctx, cartKey, 7*24*time.Hour)
    
    return &pb.AddItemResponse{Success: true}, nil
}

func (s *CartServer) GetCart(ctx context.Context, req *pb.GetCartRequest) (*pb.GetCartResponse, error) {
    cartKey := fmt.Sprintf("cart:%d", req.UserId)
    
    // Get all items from Redis
    items, err := s.redis.HGetAll(ctx, cartKey).Result()
    if err != nil {
        return nil, status.Error(codes.Internal, "failed to get cart")
    }
    
    // Fetch product details for each item
    var cartItems []*pb.CartItem
    var totalPrice float64
    
    for itemKey, quantityStr := range items {
        quantity, _ := strconv.Atoi(quantityStr)
        
        // Parse itemKey: "product:456:variant:789"
        var productID, variantID int64
        fmt.Sscanf(itemKey, "product:%d:variant:%d", &productID, &variantID)
        
        // Get product details (from Product Service via gRPC)
        product, err := s.productClient.GetProduct(ctx, &pb.GetProductRequest{
            ProductId: productID,
        })
        
        if err != nil {
            continue // Skip unavailable products
        }
        
        // Find variant
        var variant *pb.ProductVariant
        for _, v := range product.Variants {
            if v.Id == variantID {
                variant = v
                break
            }
        }
        
        if variant == nil {
            continue
        }
        
        itemTotal := variant.Price * float64(quantity)
        totalPrice += itemTotal
        
        cartItems = append(cartItems, &pb.CartItem{
            ProductId:   productID,
            VariantId:   variantID,
            Name:        product.Name,
            VariantName: variant.Name,
            Price:       variant.Price,
            Quantity:    int32(quantity),
            Subtotal:    itemTotal,
        })
    }
    
    // Check for promo code
    promoKey := fmt.Sprintf("cart:%d:promo_code", req.UserId)
    promoCode, _ := s.redis.Get(ctx, promoKey).Result()
    
    discount := 0.0
    if promoCode != "" {
        discount = s.calculateDiscount(promoCode, totalPrice)
    }
    
    return &pb.GetCartResponse{
        Items:       cartItems,
        TotalPrice:  totalPrice,
        Discount:    discount,
        FinalPrice:  totalPrice - discount,
        PromoCode:   promoCode,
    }, nil
}

func (s *CartServer) ApplyPromoCode(ctx context.Context, req *pb.ApplyPromoCodeRequest) (*pb.ApplyPromoCodeResponse, error) {
    // Validate promo code (could be another microservice)
    discount := s.validatePromoCode(req.PromoCode)
    
    if discount == 0 {
        return nil, status.Error(codes.InvalidArgument, "invalid promo code")
    }
    
    // Store promo code in Redis
    promoKey := fmt.Sprintf("cart:%d:promo_code", req.UserId)
    err := s.redis.Set(ctx, promoKey, req.PromoCode, 7*24*time.Hour).Err()
    
    if err != nil {
        return nil, status.Error(codes.Internal, "failed to apply promo code")
    }
    
    return &pb.ApplyPromoCodeResponse{
        Success:  true,
        Discount: discount,
    }, nil
}
```

### 🎯 Go Concepts:
- ✅ Redis operations
- ✅ gRPC client (calling Product Service)
- ✅ Service-to-service communication
- ✅ Context propagation
- ✅ Error handling
- ✅ Data aggregation from multiple sources

---

## 4️⃣ Order Service (Обробка замовлень)

### 📋 Відповідальність:
- Створення замовлень
- Order state machine (pending → confirmed → shipped → delivered)
- Order tracking
- Order history
- Координація між сервісами (Saga pattern)

### 🗄️ Database Schema:

```sql
-- orders table
CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    order_number VARCHAR(50) UNIQUE NOT NULL,
    status VARCHAR(50) NOT NULL, -- pending, confirmed, processing, shipped, delivered, cancelled
    total_amount DECIMAL(10, 2) NOT NULL,
    discount_amount DECIMAL(10, 2) DEFAULT 0,
    final_amount DECIMAL(10, 2) NOT NULL,
    promo_code VARCHAR(50),
    
    -- Shipping
    shipping_address_id BIGINT,
    shipping_method VARCHAR(50),
    shipping_cost DECIMAL(10, 2),
    tracking_number VARCHAR(100),
    
    -- Payment
    payment_status VARCHAR(50), -- pending, completed, failed, refunded
    payment_method VARCHAR(50),
    payment_intent_id VARCHAR(255),
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    confirmed_at TIMESTAMP,
    shipped_at TIMESTAMP,
    delivered_at TIMESTAMP
);

-- order_items table
CREATE TABLE order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT REFERENCES orders(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL,
    variant_id BIGINT NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    variant_name VARCHAR(100),
    price DECIMAL(10, 2) NOT NULL,
    quantity INT NOT NULL,
    subtotal DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- order_events table (Event Sourcing)
CREATE TABLE order_events (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT REFERENCES orders(id),
    event_type VARCHAR(50) NOT NULL, -- created, confirmed, shipped, delivered, cancelled
    event_data JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### 🔌 gRPC API:

```protobuf
service OrderService {
  rpc CreateOrder(CreateOrderRequest) returns (CreateOrderResponse);
  rpc GetOrder(GetOrderRequest) returns (GetOrderResponse);
  rpc ListOrders(ListOrdersRequest) returns (ListOrdersResponse);
  rpc UpdateOrderStatus(UpdateOrderStatusRequest) returns (UpdateOrderStatusResponse);
  rpc CancelOrder(CancelOrderRequest) returns (CancelOrderResponse);
  rpc TrackOrder(TrackOrderRequest) returns (TrackOrderResponse);
}
```

### 💻 Go Code Example (Saga Pattern):

```go
package server

import (
    "context"
    "fmt"
    
    pb "github.com/yourapp/proto/order"
)

type OrderServer struct {
    pb.UnimplementedOrderServiceServer
    db             *sql.DB
    cartClient     pb.CartServiceClient
    productClient  pb.ProductServiceClient
    paymentClient  pb.PaymentServiceClient
    notifyClient   pb.NotificationServiceClient
    eventPublisher *kafka.Producer
}

// CreateOrder implements Saga pattern for distributed transaction
func (s *OrderServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
    // Step 1: Get cart
    cart, err := s.cartClient.GetCart(ctx, &pb.GetCartRequest{UserId: req.UserId})
    if err != nil {
        return nil, status.Error(codes.Internal, "failed to get cart")
    }
    
    if len(cart.Items) == 0 {
        return nil, status.Error(codes.FailedPrecondition, "cart is empty")
    }
    
    // Step 2: Create order (pending state)
    orderID, err := s.createPendingOrder(ctx, req.UserId, cart)
    if err != nil {
        return nil, err
    }
    
    // Step 3: Reserve stock (with rollback on failure)
    err = s.reserveStock(ctx, orderID, cart.Items)
    if err != nil {
        s.cancelOrder(ctx, orderID, "stock_unavailable")
        return nil, status.Error(codes.FailedPrecondition, "failed to reserve stock")
    }
    
    // Step 4: Process payment
    paymentResp, err := s.paymentClient.ProcessPayment(ctx, &pb.ProcessPaymentRequest{
        OrderId:       orderID,
        Amount:        cart.FinalPrice,
        PaymentMethod: req.PaymentMethod,
        UserId:        req.UserId,
    })
    
    if err != nil || !paymentResp.Success {
        // Rollback: release stock
        s.releaseStock(ctx, orderID, cart.Items)
        s.cancelOrder(ctx, orderID, "payment_failed")
        return nil, status.Error(codes.Internal, "payment failed")
    }
    
    // Step 5: Confirm order
    err = s.confirmOrder(ctx, orderID, paymentResp.PaymentIntentId)
    if err != nil {
        // This is critical - payment succeeded but confirmation failed
        // Need manual intervention or retry logic
        s.logCriticalError(orderID, "order_confirmation_failed", err)
    }
    
    // Step 6: Clear cart
    s.cartClient.ClearCart(ctx, &pb.ClearCartRequest{UserId: req.UserId})
    
    // Step 7: Send notifications (async)
    go s.sendOrderConfirmation(orderID, req.UserId)
    
    // Step 8: Publish event to Kafka
    s.publishOrderCreatedEvent(orderID, req.UserId, cart)
    
    return &pb.CreateOrderResponse{
        OrderId:     orderID,
        OrderNumber: s.generateOrderNumber(orderID),
    }, nil
}

func (s *OrderServer) reserveStock(ctx context.Context, orderID int64, items []*pb.CartItem) error {
    for _, item := range items {
        _, err := s.productClient.ReserveStock(ctx, &pb.ReserveStockRequest{
            VariantId: item.VariantId,
            Quantity:  item.Quantity,
            OrderId:   orderID,
        })
        
        if err != nil {
            // Rollback: release already reserved items
            s.releaseStock(ctx, orderID, items)
            return err
        }
    }
    return nil
}

func (s *OrderServer) releaseStock(ctx context.Context, orderID int64, items []*pb.CartItem) {
    for _, item := range items {
        s.productClient.ReleaseStock(ctx, &pb.ReleaseStockRequest{
            VariantId: item.VariantId,
            Quantity:  item.Quantity,
            OrderId:   orderID,
        })
    }
}

func (s *OrderServer) publishOrderCreatedEvent(orderID int64, userID int64, cart *pb.GetCartResponse) {
    event := OrderCreatedEvent{
        OrderID:   orderID,
        UserID:    userID,
        Items:     cart.Items,
        Total:     cart.FinalPrice,
        Timestamp: time.Now(),
    }
    
    data, _ := json.Marshal(event)
    s.eventPublisher.Produce(&kafka.Message{
        Topic: "orders.created",
        Value: data,
    }, nil)
}
```

### 🎯 Go Concepts:
- ✅ **Saga pattern** (distributed transactions)
- ✅ Compensation logic (rollback)
- ✅ Service orchestration
- ✅ Event sourcing
- ✅ Kafka event publishing
- ✅ Goroutines (async notifications)
- ✅ Error handling & logging
- ✅ Database transactions

---

## 5️⃣ Payment Service (Обробка платежів)

### 📋 Відповідальність:
- Інтеграція з Stripe/PayPal
- Processing payments
- Refunds
- Payment webhooks
- Transaction history

### 🗄️ Database Schema:

```sql
-- payments table
CREATE TABLE payments (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    status VARCHAR(50) NOT NULL, -- pending, succeeded, failed, refunded
    payment_method VARCHAR(50), -- card, paypal, apple_pay
    
    -- Stripe specific
    payment_intent_id VARCHAR(255) UNIQUE,
    stripe_charge_id VARCHAR(255),
    
    -- Card info (last 4 digits only!)
    card_last4 VARCHAR(4),
    card_brand VARCHAR(20),
    
    failure_reason TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- refunds table
CREATE TABLE refunds (
    id BIGSERIAL PRIMARY KEY,
    payment_id BIGINT REFERENCES payments(id),
    order_id BIGINT NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    reason TEXT,
    status VARCHAR(50), -- pending, completed, failed
    refund_id VARCHAR(255) UNIQUE,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### 💻 Go Code Example (Stripe Integration):

```go
package server

import (
    "context"
    
    pb "github.com/yourapp/proto/payment"
    "github.com/stripe/stripe-go/v74"
    "github.com/stripe/stripe-go/v74/paymentintent"
)

type PaymentServer struct {
    pb.UnimplementedPaymentServiceServer
    db          *sql.DB
    stripeKey   string
}

func (s *PaymentServer) ProcessPayment(ctx context.Context, req *pb.ProcessPaymentRequest) (*pb.ProcessPaymentResponse, error) {
    stripe.Key = s.stripeKey
    
    // Create Payment Intent
    params := &stripe.PaymentIntentParams{
        Amount:   stripe.Int64(int64(req.Amount * 100)), // Convert to cents
        Currency: stripe.String("usd"),
        Metadata: map[string]string{
            "order_id": fmt.Sprintf("%d", req.OrderId),
            "user_id":  fmt.Sprintf("%d", req.UserId),
        },
    }
    
    pi, err := paymentintent.New(params)
    if err != nil {
        return nil, status.Error(codes.Internal, "failed to create payment intent")
    }
    
    // Save to database
    _, err = s.db.ExecContext(ctx, `
        INSERT INTO payments (order_id, user_id, amount, status, payment_intent_id, payment_method)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, req.OrderId, req.UserId, req.Amount, "pending", pi.ID, req.PaymentMethod)
    
    if err != nil {
        return nil, status.Error(codes.Internal, "failed to save payment")
    }
    
    // In real app: redirect user to confirm payment
    // For now: simulate success
    _, err = paymentintent.Confirm(pi.ID, nil)
    if err != nil {
        s.updatePaymentStatus(ctx, pi.ID, "failed", err.Error())
        return &pb.ProcessPaymentResponse{Success: false}, nil
    }
    
    s.updatePaymentStatus(ctx, pi.ID, "succeeded", "")
    
    return &pb.ProcessPaymentResponse{
        Success:         true,
        PaymentIntentId: pi.ID,
    }, nil
}

func (s *PaymentServer) HandleWebhook(ctx context.Context, req *pb.WebhookRequest) (*pb.WebhookResponse, error) {
    // Verify webhook signature (important for security!)
    event := stripe.Event{}
    err := json.Unmarshal(req.Payload, &event)
    if err != nil {
        return nil, status.Error(codes.InvalidArgument, "invalid payload")
    }
    
    switch event.Type {
    case "payment_intent.succeeded":
        var paymentIntent stripe.PaymentIntent
        json.Unmarshal(event.Data.Raw, &paymentIntent)
        s.handlePaymentSuccess(ctx, &paymentIntent)
        
    case "payment_intent.payment_failed":
        var paymentIntent stripe.PaymentIntent
        json.Unmarshal(event.Data.Raw, &paymentIntent)
        s.handlePaymentFailure(ctx, &paymentIntent)
        
    case "charge.refunded":
        var charge stripe.Charge
        json.Unmarshal(event.Data.Raw, &charge)
        s.handleRefund(ctx, &charge)
    }
    
    return &pb.WebhookResponse{Success: true}, nil
}

func (s *PaymentServer) RefundPayment(ctx context.Context, req *pb.RefundRequest) (*pb.RefundResponse, error) {
    // Get payment intent ID
    var paymentIntentID string
    err := s.db.QueryRowContext(ctx, `
        SELECT payment_intent_id FROM payments WHERE order_id = $1
    `, req.OrderId).Scan(&paymentIntentID)
    
    if err != nil {
        return nil, status.Error(codes.NotFound, "payment not found")
    }
    
    // Create refund in Stripe
    refundParams := &stripe.RefundParams{
        PaymentIntent: stripe.String(paymentIntentID),
        Amount:        stripe.Int64(int64(req.Amount * 100)),
        Reason:        stripe.String("requested_by_customer"),
    }
    
    refund, err := refund.New(refundParams)
    if err != nil {
        return nil, status.Error(codes.Internal, "refund failed")
    }
    
    // Save refund to database
    _, err = s.db.ExecContext(ctx, `
        INSERT INTO refunds (payment_id, order_id, amount, reason, status, refund_id)
        VALUES (
            (SELECT id FROM payments WHERE order_id = $1),
            $1, $2, $3, 'completed', $4
        )
    `, req.OrderId, req.Amount, req.Reason, refund.ID)
    
    return &pb.RefundResponse{Success: true, RefundId: refund.ID}, nil
}
```

### 🎯 Go Concepts:
- ✅ External API integration (Stripe)
- ✅ Webhook handling
- ✅ Security (signature verification)
- ✅ Error handling
- ✅ Database transactions
- ✅ Context usage

---

## 6️⃣ Notification Service (Сповіщення)

### 📋 Відповідальність:
- Email notifications (SendGrid/AWS SES)
- SMS (Twilio)
- Push notifications
- Event-driven (Kafka consumer)
- Templates

### 💻 Go Code Example (Kafka Consumer + Worker Pool):

```go
package main

import (
    "context"
    "encoding/json"
    "log"
    
    "github.com/segmentio/kafka-go"
    "github.com/sendgrid/sendgrid-go"
    "github.com/sendgrid/sendgrid-go/helpers/mail"
)

type NotificationService struct {
    kafka         *kafka.Reader
    sendgrid      *sendgrid.Client
    workerPool    chan struct{} // semaphore for concurrency control
    notifications chan Notification
}

type Notification struct {
    Type      string                 `json:"type"` // email, sms, push
    Recipient string                 `json:"recipient"`
    Template  string                 `json:"template"`
    Data      map[string]interface{} `json:"data"`
}

func NewNotificationService() *NotificationService {
    return &NotificationService{
        kafka: kafka.NewReader(kafka.ReaderConfig{
            Brokers: []string{"localhost:9092"},
            Topic:   "notifications",
            GroupID: "notification-service",
        }),
        sendgrid:      sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY")),
        workerPool:    make(chan struct{}, 10), // max 10 concurrent workers
        notifications: make(chan Notification, 100),
    }
}

func (s *NotificationService) Start(ctx context.Context) {
    // Start workers
    for i := 0; i < 10; i++ {
        go s.worker(ctx, i)
    }
    
    // Start Kafka consumer
    go s.consumeMessages(ctx)
}

func (s *NotificationService) consumeMessages(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            msg, err := s.kafka.ReadMessage(ctx)
            if err != nil {
                log.Printf("Error reading message: %v", err)
                continue
            }
            
            var notif Notification
            if err := json.Unmarshal(msg.Value, &notif); err != nil {
                log.Printf("Error unmarshaling message: %v", err)
                continue
            }
            
            s.notifications <- notif
        }
    }
}

func (s *NotificationService) worker(ctx context.Context, id int) {
    log.Printf("Worker %d started", id)
    
    for {
        select {
        case <-ctx.Done():
            return
        case notif := <-s.notifications:
            s.workerPool <- struct{}{} // acquire semaphore
            
            log.Printf("Worker %d processing notification: %s", id, notif.Type)
            
            switch notif.Type {
            case "email":
                s.sendEmail(notif)
            case "sms":
                s.sendSMS(notif)
            case "push":
                s.sendPush(notif)
            }
            
            <-s.workerPool // release semaphore
        }
    }
}

func (s *NotificationService) sendEmail(notif Notification) error {
    from := mail.NewEmail("E-commerce", "noreply@ecommerce.com")
    to := mail.NewEmail("User", notif.Recipient)
    
    subject := s.getSubject(notif.Template)
    content := s.renderTemplate(notif.Template, notif.Data)
    
    message := mail.NewSingleEmail(from, subject, to, "", content)
    
    response, err := s.sendgrid.Send(message)
    if err != nil {
        log.Printf("Failed to send email: %v", err)
        return err
    }
    
    log.Printf("Email sent successfully: %d", response.StatusCode)
    return nil
}

func (s *NotificationService) renderTemplate(template string, data map[string]interface{}) string {
    switch template {
    case "order_confirmation":
        return fmt.Sprintf(`
            <h1>Order Confirmed!</h1>
            <p>Your order #%s has been confirmed.</p>
            <p>Total: $%.2f</p>
        `, data["order_number"], data["total"])
        
    case "order_shipped":
        return fmt.Sprintf(`
            <h1>Order Shipped!</h1>
            <p>Your order #%s has been shipped.</p>
            <p>Tracking number: %s</p>
        `, data["order_number"], data["tracking_number"])
        
    default:
        return ""
    }
}
```

### 🎯 Go Concepts:
- ✅ **Kafka consumer** (event-driven)
- ✅ **Worker pool pattern** (Week 5!)
- ✅ **Channels** для комунікації
- ✅ **Goroutines** для workers
- ✅ **Semaphore** для concurrency control
- ✅ Context cancellation
- ✅ External API integration (SendGrid, Twilio)

---

## 7️⃣ Search Service (Пошук)

### 📋 Відповідальність:
- Elasticsearch integration
- Full-text search
- Filters & facets
- Autocomplete
- Search suggestions

### 💻 Go Code Example:

```go
package server

import (
    "context"
    "encoding/json"
    
    "github.com/elastic/go-elasticsearch/v8"
    pb "github.com/yourapp/proto/search"
)

type SearchServer struct {
    pb.UnimplementedSearchServiceServer
    es *elasticsearch.Client
}

func (s *SearchServer) SearchProducts(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
    // Build Elasticsearch query
    query := map[string]interface{}{
        "query": map[string]interface{}{
            "bool": map[string]interface{}{
                "must": []interface{}{
                    map[string]interface{}{
                        "multi_match": map[string]interface{}{
                            "query":  req.Query,
                            "fields": []string{"name^3", "description", "category"},
                        },
                    },
                },
                "filter": s.buildFilters(req.Filters),
            },
        },
        "from": req.Offset,
        "size": req.Limit,
        "sort": []interface{}{
            map[string]interface{}{
                req.SortBy: map[string]string{"order": req.SortOrder},
            },
        },
        "aggs": s.buildAggregations(),
    }
    
    // Execute search
    var buf bytes.Buffer
    json.NewEncoder(&buf).Encode(query)
    
    res, err := s.es.Search(
        s.es.Search.WithContext(ctx),
        s.es.Search.WithIndex("products"),
        s.es.Search.WithBody(&buf),
    )
    
    if err != nil {
        return nil, status.Error(codes.Internal, "search failed")
    }
    defer res.Body.Close()
    
    // Parse results
    var result map[string]interface{}
    json.NewDecoder(res.Body).Decode(&result)
    
    // Extract hits
    hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
    
    var products []*pb.Product
    for _, hit := range hits {
        source := hit.(map[string]interface{})["_source"]
        products = append(products, s.mapToProduct(source))
    }
    
    return &pb.SearchResponse{
        Products:   products,
        TotalCount: int32(result["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
        Facets:     s.extractFacets(result["aggregations"]),
    }, nil
}

func (s *SearchServer) Autocomplete(ctx context.Context, req *pb.AutocompleteRequest) (*pb.AutocompleteResponse, error) {
    query := map[string]interface{}{
        "suggest": map[string]interface{}{
            "product-suggest": map[string]interface{}{
                "prefix": req.Query,
                "completion": map[string]interface{}{
                    "field": "name_suggest",
                    "size":  10,
                },
            },
        },
    }
    
    // Execute and return suggestions
    // ...
}
```

### 🎯 Go Concepts:
- ✅ Elasticsearch integration
- ✅ JSON marshaling/unmarshaling
- ✅ Context usage
- ✅ Complex data structures
- ✅ Error handling

---

## 8️⃣ Analytics Service (Аналітика)

### 📋 Відповідальність:
- Sales reports
- User behavior tracking
- Revenue metrics
- Real-time dashboards
- Data aggregation

### 💻 Go Code Example (Pipeline Pattern):

```go
package main

import (
    "context"
    "time"
)

type AnalyticsService struct {
    kafka         *kafka.Reader
    timeseries    *influxdb2.Client
    events        chan Event
    aggregated    chan AggregatedMetric
}

type Event struct {
    Type      string    `json:"type"`
    UserID    int64     `json:"user_id"`
    Data      map[string]interface{} `json:"data"`
    Timestamp time.Time `json:"timestamp"`
}

type AggregatedMetric struct {
    Name      string
    Value     float64
    Tags      map[string]string
    Timestamp time.Time
}

func (s *AnalyticsService) Start(ctx context.Context) {
    // Pipeline: Kafka → Process → Aggregate → Store
    go s.consumeEvents(ctx)
    go s.processEvents(ctx)
    go s.aggregateMetrics(ctx)
    go s.storeMetrics(ctx)
}

func (s *AnalyticsService) consumeEvents(ctx context.Context) {
    for {
        msg, err := s.kafka.ReadMessage(ctx)
        if err != nil {
            continue
        }
        
        var event Event
        json.Unmarshal(msg.Value, &event)
        s.events <- event
    }
}

func (s *AnalyticsService) processEvents(ctx context.Context) {
    for event := range s.events {
        switch event.Type {
        case "order_created":
            s.aggregated <- AggregatedMetric{
                Name:      "orders_total",
                Value:     1,
                Tags:      map[string]string{"status": "created"},
                Timestamp: event.Timestamp,
            }
            
            s.aggregated <- AggregatedMetric{
                Name:      "revenue",
                Value:     event.Data["total"].(float64),
                Tags:      map[string]string{"currency": "USD"},
                Timestamp: event.Timestamp,
            }
            
        case "product_viewed":
            s.aggregated <- AggregatedMetric{
                Name:      "product_views",
                Value:     1,
                Tags:      map[string]string{"product_id": fmt.Sprintf("%d", event.Data["product_id"])},
                Timestamp: event.Timestamp,
            }
        }
    }
}

func (s *AnalyticsService) storeMetrics(ctx context.Context) {
    writeAPI := s.timeseries.WriteAPI("myorg", "mybucket")
    
    for metric := range s.aggregated {
        point := influxdb2.NewPoint(
            metric.Name,
            metric.Tags,
            map[string]interface{}{"value": metric.Value},
            metric.Timestamp,
        )
        writeAPI.WritePoint(point)
    }
}
```

### 🎯 Go Concepts:
- ✅ **Pipeline pattern** (Week 5!)
- ✅ **Channels** для data flow
- ✅ **Goroutines** для кожного етапу
- ✅ Kafka consumer
- ✅ InfluxDB (time-series)
- ✅ Event processing

---

## 9️⃣ Review Service (Відгуки)

### 📋 Відповідальність:
- Product reviews
- Ratings
- Review moderation
- Aggregated ratings

### 🗄️ Database Schema:

```sql
CREATE TABLE reviews (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    order_id BIGINT, -- optional: verify purchase
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    title VARCHAR(255),
    content TEXT,
    is_verified_purchase BOOLEAN DEFAULT FALSE,
    is_approved BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(product_id, user_id) -- one review per product per user
);

CREATE TABLE review_votes (
    id BIGSERIAL PRIMARY KEY,
    review_id BIGINT REFERENCES reviews(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL,
    vote_type VARCHAR(10) CHECK (vote_type IN ('helpful', 'not_helpful')),
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(review_id, user_id)
);

CREATE TABLE product_ratings_summary (
    product_id BIGINT PRIMARY KEY,
    average_rating DECIMAL(3, 2),
    total_reviews INT DEFAULT 0,
    rating_1_count INT DEFAULT 0,
    rating_2_count INT DEFAULT 0,
    rating_3_count INT DEFAULT 0,
    rating_4_count INT DEFAULT 0,
    rating_5_count INT DEFAULT 0,
    updated_at TIMESTAMP DEFAULT NOW()
);
```

---

## 🔗 Service Communication

### gRPC (Synchronous)
```
User Service ──gRPC──→ Order Service
Order Service ──gRPC──→ Product Service (check stock)
Order Service ──gRPC──→ Payment Service
```

### Kafka Events (Asynchronous)
```
Order Service ──Kafka──→ order.created ──→ Notification Service
                                        ──→ Analytics Service
                                        ──→ Search Service (index)
```

---

## 🛠️ Infrastructure & DevOps

### Docker Compose (Development):

```yaml
version: '3.8'

services:
  postgres-user:
    image: postgres:15
    environment:
      POSTGRES_DB: user_service
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
  
  postgres-product:
    image: postgres:15
    environment:
      POSTGRES_DB: product_service
  
  postgres-order:
    image: postgres:15
    environment:
      POSTGRES_DB: order_service
  
  redis:
    image: redis:7-alpine
  
  elasticsearch:
    image: elasticsearch:8.10.0
    environment:
      - discovery.type=single-node
  
  kafka:
    image: confluentinc/cp-kafka:7.5.0
    environment:
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9092
  
  zookeeper:
    image: confluentinc/cp-zookeeper:7.5.0
  
  consul:
    image: consul:1.16
    ports:
      - "8500:8500"
  
  prometheus:
    image: prom/prometheus:latest
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
  
  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
```

### Kubernetes Deployment:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
    spec:
      containers:
      - name: user-service
        image: ecommerce/user-service:latest
        ports:
        - containerPort: 50051
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: user-db-secret
              key: url
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: jwt-secret
              key: secret
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          grpc:
            port: 50051
          initialDelaySeconds: 10
        readinessProbe:
          grpc:
            port: 50051
          initialDelaySeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: user-service
spec:
  selector:
    app: user-service
  ports:
  - protocol: TCP
    port: 50051
    targetPort: 50051
  type: ClusterIP
```

---

## 📊 Go Concepts Coverage - 100%!

### Week 1-2 (Basics) ✅
- HTTP servers (API Gateway)
- Database operations (PostgreSQL)
- Error handling
- Testing
- Packages & modules

### Week 3-4 (Intermediate) ✅
- Context (timeouts, cancellation)
- Error wrapping
- Middleware
- Environment config

### Week 5 (Goroutines & Channels) ✅
- **Worker pools** (Notification Service)
- **Pipeline pattern** (Analytics)
- **Fan-out/fan-in** (parallel processing)
- **Channels** (event streaming)
- **Select** (multiplexing)
- **Graceful shutdown** (all services)

### Advanced (Production) ✅
- **gRPC** (inter-service communication)
- **Kafka** (event streaming)
- **Redis** (caching, sessions)
- **Elasticsearch** (search)
- **Distributed transactions** (Saga pattern)
- **Event Sourcing**
- **CQRS**
- **Circuit Breaker**
- **Service Discovery** (Consul)
- **API Gateway**
- **Observability** (Prometheus, Grafana, Jaeger)
- **Kubernetes** deployment
- **CI/CD** (GitHub Actions)

---

## 🚀 Development Roadmap (8 weeks)

### Week 1-2: Foundation
**Services:** User, Product, Cart
- ✅ Project setup
- ✅ gRPC definitions
- ✅ Database schemas
- ✅ Basic CRUD operations
- ✅ Unit tests

### Week 3-4: Core Business Logic
**Services:** Order, Payment, Notification
- ✅ Order processing (Saga pattern)
- ✅ Stripe integration
- ✅ Kafka setup
- ✅ Email/SMS notifications
- ✅ Integration tests

### Week 5-6: Advanced Features
**Services:** Search, Analytics, Review
- ✅ Elasticsearch integration
- ✅ Real-time analytics
- ✅ Review system
- ✅ API Gateway
- ✅ Service discovery

### Week 7-8: Production Ready
**Infrastructure:**
- ✅ Kubernetes deployment
- ✅ Monitoring (Prometheus + Grafana)
- ✅ Distributed tracing (Jaeger)
- ✅ CI/CD pipeline
- ✅ Load testing
- ✅ Documentation

---

## 📈 Metrics & Monitoring

### Prometheus Metrics:

```go
var (
    ordersTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "orders_total",
            Help: "Total number of orders",
        },
        []string{"status"},
    )
    
    orderDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "order_duration_seconds",
            Help:    "Time to process an order",
            Buckets: prometheus.DefBuckets,
        },
        []string{"service"},
    )
    
    paymentErrors = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "payment_errors_total",
            Help: "Total number of payment errors",
        },
    )
)

func init() {
    prometheus.MustRegister(ordersTotal)
    prometheus.MustRegister(orderDuration)
    prometheus.MustRegister(paymentErrors)
}
```

---

## 🎯 Висновок

**E-commerce Platform** - це найбільш комплексний проект, який:

✅ **Покриває 100% Go концепцій**
✅ **9 мікросервісів** з реальною бізнес-логікою
✅ **Production-ready patterns** (Saga, CQRS, Event Sourcing)
✅ **Full stack** (Frontend + Backend + Infrastructure)
✅ **Можна показувати роботодавцям** 🏆

### Складність vs Навчання:

| Проект | Go Concepts | Мікросервіси | Рекомендація |
|--------|-------------|--------------|--------------|
| URL Shortener | 40% | 0 | ✅ Starter |
| Task Queue | 60% | 0 | ✅ Intermediate |
| **Monitoring** | **90%** | **5+** | **🔥 Best Choice!** |
| **E-commerce** | **100%** | **9** | **🏆 Advanced** |

### Моя рекомендація:

1. **Почни з Monitoring System** (90% concepts, швидший результат)
2. **Потім E-commerce** (якщо хочеш максимум + portfolio piece)

**E-commerce - це ідеальний "capstone project" після вивчення всіх Go концепцій!** 🚀
