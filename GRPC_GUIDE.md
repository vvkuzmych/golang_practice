# 🔌 gRPC Communication Guide

## Що таке gRPC?

**gRPC** (gRPC Remote Procedure Call) - це сучасний, високопродуктивний RPC framework від Google для комунікації між сервісами.

### 📝 Простими словами:

**gRPC дозволяє викликати функції на іншому сервері так, ніби вони локальні.**

```go
// Замість HTTP REST:
resp, err := http.Get("http://user-service/users/123")
data := parseJSON(resp.Body)

// З gRPC просто викликаєш функцію:
user, err := client.GetUser(ctx, &pb.GetUserRequest{Id: 123})
```

---

## 🆚 gRPC vs REST API

| Аспект | REST API | gRPC |
|--------|----------|------|
| **Протокол** | HTTP/1.1 (text) | HTTP/2 (binary) |
| **Формат даних** | JSON (text) | Protocol Buffers (binary) |
| **Швидкість** | Повільніше | **7-10x швидше!** |
| **Розмір даних** | Більший (JSON) | **Менший на 30-50%** |
| **Streaming** | ❌ Складно | ✅ Вбудований |
| **Type safety** | ❌ Runtime | ✅ **Compile-time** |
| **Code generation** | ❌ Ручний | ✅ **Автоматичний** |
| **Browser support** | ✅ Native | ⚠️ Потребує gRPC-Web |
| **Human-readable** | ✅ Так (JSON) | ❌ Ні (binary) |
| **Best for** | Client-facing APIs | **Microservices** |

### 🎯 Коли використовувати gRPC:

✅ **Використовуй gRPC:**
- Комунікація між мікросервісами (backend-to-backend)
- Потрібна висока продуктивність
- Real-time streaming (logs, metrics)
- Polyglot environments (різні мови)
- Type safety критична

❌ **Використовуй REST:**
- Public API для фронтенду (browsers)
- Простота важливіша за performance
- Debugging має бути легким (curl, Postman)

---

## 🏗️ Як працює gRPC?

### Архітектура:

```
┌─────────────┐                           ┌─────────────┐
│   Client    │                           │   Server    │
│             │                           │             │
│  Go code    │                           │  Go code    │
│     ↓       │                           │     ↑       │
│  gRPC stub  │─────── HTTP/2 ────────→  │ gRPC server │
│     ↓       │   (Protocol Buffers)      │     ↑       │
│  Network    │                           │  Network    │
└─────────────┘                           └─────────────┘
```

### Кроки:

1. **Define** service у `.proto` файлі
2. **Generate** Go code з protobuf compiler
3. **Implement** server (реалізація методів)
4. **Call** методи з client

---

## 📦 Protocol Buffers (Protobuf)

**Protobuf** - це мова для опису структури даних (як JSON schema, але краще).

### Приклад `.proto` файлу:

```protobuf
syntax = "proto3";

package user;

option go_package = "github.com/yourapp/proto/user";

// Service definition
service UserService {
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
  rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
}

// Messages (data structures)
message User {
  int64 id = 1;
  string username = 2;
  string email = 3;
  int32 age = 4;
  bool is_active = 5;
}

message GetUserRequest {
  int64 id = 1;
}

message GetUserResponse {
  User user = 1;
  string error = 2;
}

message ListUsersRequest {
  int32 page = 1;
  int32 page_size = 2;
}

message ListUsersResponse {
  repeated User users = 1;  // repeated = array
  int32 total = 2;
}

message CreateUserRequest {
  string username = 1;
  string email = 2;
  int32 age = 3;
}

message CreateUserResponse {
  User user = 1;
  string error = 2;
}

message UpdateUserRequest {
  int64 id = 1;
  string username = 2;
  string email = 3;
  int32 age = 4;
}

message UpdateUserResponse {
  User user = 1;
  string error = 2;
}

message DeleteUserRequest {
  int64 id = 1;
}

message DeleteUserResponse {
  bool success = 1;
  string error = 2;
}
```

### Що означають числа (1, 2, 3...)?

Це **field numbers** - унікальні ідентифікатори полів у binary форматі. **НЕ змінюй їх після release!**

---

## 🔧 Налаштування Go Project

### 1. Встановити інструменти:

```bash
# Protocol Buffers compiler
brew install protobuf  # macOS
# or
sudo apt install protobuf-compiler  # Linux

# Go plugins для protoc
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Перевірити
protoc --version  # libprotoc 3.21.0 or higher
```

### 2. Структура проекту:

```
myapp/
├── proto/
│   └── user/
│       └── user.proto          # Protobuf definitions
├── pkg/
│   └── pb/                     # Generated code (pb = protobuf)
│       └── user/
│           ├── user.pb.go      # Generated message types
│           └── user_grpc.pb.go # Generated service code
├── services/
│   └── user-service/
│       ├── server/
│       │   └── server.go       # gRPC server implementation
│       └── main.go
└── client/
    └── main.go                 # gRPC client example
```

### 3. Generate Go code:

```bash
# З кореня проекту
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/user/user.proto

# Або Makefile:
.PHONY: proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       proto/**/*.proto
```

---

## 🖥️ Server Implementation

### server.go:

```go
package server

import (
	"context"
	"fmt"
	
	pb "github.com/yourapp/pkg/pb/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserServer реалізує UserServiceServer interface (згенерований)
type UserServer struct {
	pb.UnimplementedUserServiceServer  // Для forward compatibility
	users map[int64]*pb.User            // In-memory storage (для прикладу)
	nextID int64
}

func NewUserServer() *UserServer {
	return &UserServer{
		users: make(map[int64]*pb.User),
		nextID: 1,
	}
}

// GetUser - отримати користувача за ID
func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// Validation
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %d", req.Id)
	}
	
	// Find user
	user, exists := s.users[req.Id]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "user %d not found", req.Id)
	}
	
	return &pb.GetUserResponse{
		User: user,
	}, nil
}

// ListUsers - список користувачів (з пагінацією)
func (s *UserServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	// Default pagination
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	
	// Collect all users
	var users []*pb.User
	for _, user := range s.users {
		users = append(users, user)
	}
	
	// Pagination
	start := int((page - 1) * pageSize)
	end := int(page * pageSize)
	if start >= len(users) {
		return &pb.ListUsersResponse{Users: []*pb.User{}, Total: int32(len(users))}, nil
	}
	if end > len(users) {
		end = len(users)
	}
	
	return &pb.ListUsersResponse{
		Users: users[start:end],
		Total: int32(len(users)),
	}, nil
}

// CreateUser - створити користувача
func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// Validation
	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	
	// Create user
	user := &pb.User{
		Id:       s.nextID,
		Username: req.Username,
		Email:    req.Email,
		Age:      req.Age,
		IsActive: true,
	}
	
	s.users[s.nextID] = user
	s.nextID++
	
	return &pb.CreateUserResponse{
		User: user,
	}, nil
}

// UpdateUser - оновити користувача
func (s *UserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	user, exists := s.users[req.Id]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "user %d not found", req.Id)
	}
	
	// Update fields
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Age > 0 {
		user.Age = req.Age
	}
	
	return &pb.UpdateUserResponse{
		User: user,
	}, nil
}

// DeleteUser - видалити користувача
func (s *UserServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	_, exists := s.users[req.Id]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "user %d not found", req.Id)
	}
	
	delete(s.users, req.Id)
	
	return &pb.DeleteUserResponse{
		Success: true,
	}, nil
}
```

### main.go (server):

```go
package main

import (
	"fmt"
	"log"
	"net"
	
	pb "github.com/yourapp/pkg/pb/user"
	"github.com/yourapp/services/user-service/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Create TCP listener
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	
	// Create gRPC server
	grpcServer := grpc.NewServer()
	
	// Register service
	userServer := server.NewUserServer()
	pb.RegisterUserServiceServer(grpcServer, userServer)
	
	// Register reflection service (для grpcurl, grpcui)
	reflection.Register(grpcServer)
	
	fmt.Println("gRPC server listening on :50051")
	
	// Start server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
```

---

## 📱 Client Implementation

### client.go:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"
	
	pb "github.com/yourapp/pkg/pb/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to server
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()
	
	// Create client
	client := pb.NewUserServiceClient(conn)
	
	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// === Create User ===
	fmt.Println("Creating user...")
	createResp, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Username: "john_doe",
		Email:    "john@example.com",
		Age:      30,
	})
	if err != nil {
		log.Fatalf("CreateUser failed: %v", err)
	}
	fmt.Printf("Created user: ID=%d, Username=%s\n", 
		createResp.User.Id, createResp.User.Username)
	
	userID := createResp.User.Id
	
	// === Get User ===
	fmt.Println("\nGetting user...")
	getResp, err := client.GetUser(ctx, &pb.GetUserRequest{
		Id: userID,
	})
	if err != nil {
		log.Fatalf("GetUser failed: %v", err)
	}
	fmt.Printf("Got user: %+v\n", getResp.User)
	
	// === Update User ===
	fmt.Println("\nUpdating user...")
	updateResp, err := client.UpdateUser(ctx, &pb.UpdateUserRequest{
		Id:       userID,
		Username: "john_updated",
		Age:      31,
	})
	if err != nil {
		log.Fatalf("UpdateUser failed: %v", err)
	}
	fmt.Printf("Updated user: %+v\n", updateResp.User)
	
	// === List Users ===
	fmt.Println("\nListing users...")
	listResp, err := client.ListUsers(ctx, &pb.ListUsersRequest{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		log.Fatalf("ListUsers failed: %v", err)
	}
	fmt.Printf("Found %d users (total: %d)\n", len(listResp.Users), listResp.Total)
	for _, user := range listResp.Users {
		fmt.Printf("  - %d: %s (%s)\n", user.Id, user.Username, user.Email)
	}
	
	// === Delete User ===
	fmt.Println("\nDeleting user...")
	deleteResp, err := client.DeleteUser(ctx, &pb.DeleteUserRequest{
		Id: userID,
	})
	if err != nil {
		log.Fatalf("DeleteUser failed: %v", err)
	}
	fmt.Printf("Deleted user: success=%v\n", deleteResp.Success)
}
```

---

## 🚀 Запуск

### Terminal 1 (Server):

```bash
cd services/user-service
go run main.go

# Output:
# gRPC server listening on :50051
```

### Terminal 2 (Client):

```bash
cd client
go run main.go

# Output:
# Creating user...
# Created user: ID=1, Username=john_doe
# 
# Getting user...
# Got user: id:1 username:"john_doe" email:"john@example.com" age:30 is_active:true
# 
# Updating user...
# Updated user: id:1 username:"john_updated" email:"john@example.com" age:31 is_active:true
# 
# Listing users...
# Found 1 users (total: 1)
#   - 1: john_updated (john@example.com)
# 
# Deleting user...
# Deleted user: success=true
```

---

## 🔥 Advanced Features

### 1️⃣ Streaming (Real-time Data)

**4 типи streaming:**

#### Unary RPC (звичайний)
```protobuf
rpc GetUser(GetUserRequest) returns (GetUserResponse);
```

#### Server Streaming (server → client stream)
```protobuf
rpc StreamUsers(StreamUsersRequest) returns (stream User);
```

**Використання:** Real-time updates, log streaming

```go
// Server:
func (s *UserServer) StreamUsers(req *pb.StreamUsersRequest, stream pb.UserService_StreamUsersServer) error {
	for _, user := range s.users {
		if err := stream.Send(user); err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond) // Simulate delay
	}
	return nil
}

// Client:
stream, err := client.StreamUsers(ctx, &pb.StreamUsersRequest{})
for {
	user, err := stream.Recv()
	if err == io.EOF {
		break
	}
	fmt.Printf("Received: %+v\n", user)
}
```

#### Client Streaming (client → server stream)
```protobuf
rpc CreateBatchUsers(stream CreateUserRequest) returns (CreateBatchResponse);
```

**Використання:** File upload, batch processing

#### Bidirectional Streaming (обидва напрямки)
```protobuf
rpc Chat(stream ChatMessage) returns (stream ChatMessage);
```

**Використання:** Chat, real-time collaboration

---

### 2️⃣ Interceptors (Middleware)

```go
// Logging interceptor
func loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	
	// Before
	log.Printf("→ Method: %s, Request: %+v", info.FullMethod, req)
	
	// Call handler
	resp, err := handler(ctx, req)
	
	// After
	log.Printf("← Method: %s, Duration: %v, Error: %v", 
		info.FullMethod, time.Since(start), err)
	
	return resp, err
}

// Register interceptor
grpcServer := grpc.NewServer(
	grpc.UnaryInterceptor(loggingInterceptor),
)
```

**Use cases:**
- Authentication/Authorization
- Logging
- Metrics
- Rate limiting
- Error handling

---

### 3️⃣ Error Handling

```go
import "google.golang.org/grpc/codes"
import "google.golang.org/grpc/status"

// Server:
if user == nil {
	return nil, status.Errorf(codes.NotFound, "user %d not found", id)
}

// Client:
_, err := client.GetUser(ctx, req)
if err != nil {
	st, ok := status.FromError(err)
	if ok {
		switch st.Code() {
		case codes.NotFound:
			fmt.Println("User not found")
		case codes.InvalidArgument:
			fmt.Println("Invalid request")
		default:
			fmt.Printf("Error: %v\n", st.Message())
		}
	}
}
```

**gRPC Status Codes:**
- `OK` - success
- `InvalidArgument` - invalid data
- `NotFound` - resource not found
- `PermissionDenied` - auth failed
- `Unavailable` - service down
- `Internal` - server error
- [Full list](https://grpc.github.io/grpc/core/md_doc_statuscodes.html)

---

### 4️⃣ Metadata (Headers)

```go
import "google.golang.org/grpc/metadata"

// Client: Send metadata
md := metadata.Pairs(
	"authorization", "Bearer token123",
	"user-agent", "my-app/1.0",
)
ctx := metadata.NewOutgoingContext(context.Background(), md)
resp, err := client.GetUser(ctx, req)

// Server: Read metadata
func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		authToken := md.Get("authorization")
		fmt.Printf("Auth token: %v\n", authToken)
	}
	// ...
}
```

---

## 📊 Real-World Example: Microservices

### Monitoring System (з PROJECT_IDEAS.md):

```
┌─────────────┐     gRPC      ┌─────────────┐
│   Agent     │─────────────→ │  Collector  │
│ (Go/Python) │   Metrics     │   Service   │
└─────────────┘               └─────────────┘
                                      │
                                      │ gRPC
                                      ↓
                              ┌─────────────┐
                              │   Storage   │
                              │   Service   │
                              └─────────────┘
                                      │
                                      │ gRPC
                                      ↓
                              ┌─────────────┐
                              │    Query    │
                              │   Service   │
                              └─────────────┘
```

**Чому gRPC тут ідеальний:**
- ✅ Висока пропускна здатність (millions of metrics/sec)
- ✅ Низька латентність (<1ms between services)
- ✅ Server streaming (metrics flow)
- ✅ Type safety (метрики мають чітку структуру)
- ✅ Polyglot (agents можуть бути на Python, Go, Java)

---

## 🛠️ Інструменти для роботи з gRPC

### 1. grpcurl (curl для gRPC)

```bash
# Install
brew install grpcurl

# List services
grpcurl -plaintext localhost:50051 list

# List methods
grpcurl -plaintext localhost:50051 list user.UserService

# Call method
grpcurl -plaintext -d '{"id": 1}' \
  localhost:50051 user.UserService/GetUser
```

### 2. grpcui (Web UI для gRPC)

```bash
# Install
go install github.com/fullstorydev/grpcui/cmd/grpcui@latest

# Run
grpcui -plaintext localhost:50051

# Opens browser with GUI!
```

### 3. BloomRPC (GUI Client)

[bloomrpc.github.io](https://github.com/bloomrpc/bloomrpc) - як Postman для gRPC

---

## ✅ Best Practices

### 1. Versioning
```protobuf
// v1/user.proto
package user.v1;

// v2/user.proto
package user.v2;
```

### 2. Backward Compatibility
- ❌ НЕ видаляй field numbers
- ❌ НЕ змінюй тип поля
- ✅ Додавай нові поля з новими numbers

### 3. Error Handling
- Використовуй правильні status codes
- Додавай details до помилок
- Log помилки на сервері

### 4. Context
- Завжди передавай context
- Встановлюй timeout
- Обробляй cancellation

### 5. Testing
```go
// Mock server для тестів
import "google.golang.org/grpc/test/bufconn"

func setupTestServer(t *testing.T) (*grpc.Server, *bufconn.Listener) {
	lis := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, &UserServer{})
	go server.Serve(lis)
	return server, lis
}
```

---

## 📚 Ресурси

### Офіційна документація:
- [grpc.io](https://grpc.io)
- [Protocol Buffers](https://protobuf.dev)
- [gRPC-Go](https://github.com/grpc/grpc-go)

### Tutorials:
- [gRPC Basics - Go](https://grpc.io/docs/languages/go/basics/)
- [gRPC Masterclass](https://www.udemy.com/course/grpc-golang/)

### Книги:
- "gRPC: Up and Running" by Kasun Indrasiri

---

## 🎯 Висновок

### gRPC - це:

✅ **Швидкий** - 7-10x швидше за REST
✅ **Компактний** - binary формат
✅ **Type-safe** - compile-time перевірки
✅ **Streaming** - вбудована підтримка
✅ **Polyglot** - працює з усіма мовами
✅ **Ідеальний для мікросервісів**

### Коли використовувати:

**gRPC:** Backend-to-backend комунікація (мікросервіси)
**REST:** Public APIs для браузерів

### З чого почати:

1. Встановити `protoc` та Go plugins
2. Написати `.proto` файл
3. Згенерувати Go code
4. Реалізувати server
5. Написати client
6. Profit! 🚀

---

**gRPC - це must-have для сучасних мікросервісів на Go!** 💪

*Готовий створити свій перший gRPC сервіс?* 🔥
