# Client-Server Communication Patterns

## 📖 Зміст

1. [HTTP Request-Response](#1-http-request-response)
2. [Polling](#2-polling)
3. [Long Polling](#3-long-polling)
4. [WebSockets](#4-websockets)
5. [Server-Sent Events (SSE)](#5-server-sent-events-sse)
6. [gRPC](#6-grpc)
7. [GraphQL](#7-graphql)

---

## 1. HTTP Request-Response

### Класичний підхід

**Клієнт ініціює запит → Сервер відповідає**

```
Client                    Server
  |                          |
  |-------- GET /users ----->|
  |                          |
  |<----- 200 OK + Data -----|
  |                          |
```

### Go Client Example

```go
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

func main() {
    // 1. Simple GET request
    resp, err := http.Get("https://api.example.com/users/1")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))
    
    // 2. GET with parameters
    req, _ := http.NewRequest("GET", "https://api.example.com/users", nil)
    q := req.URL.Query()
    q.Add("page", "1")
    q.Add("per_page", "10")
    req.URL.RawQuery = q.Encode()
    
    client := &http.Client{}
    resp, _ = client.Do(req)
    defer resp.Body.Close()
    
    // 3. POST request
    user := User{Name: "John"}
    jsonData, _ := json.Marshal(user)
    
    resp, _ = http.Post(
        "https://api.example.com/users",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    defer resp.Body.Close()
}
```

### Go Server Example

```go
package main

import (
    "encoding/json"
    "net/http"
)

func getUser(w http.ResponseWriter, r *http.Request) {
    user := User{ID: 1, Name: "John"}
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(user)
}

func createUser(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Save user...
    user.ID = 123
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func main() {
    http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
            getUser(w, r)
        case "POST":
            createUser(w, r)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })
    
    http.ListenAndServe(":8080", nil)
}
```

### ✅ Переваги
- Простий
- Stateless
- Кешується
- Широко підтримується

### ❌ Недоліки
- Клієнт завжди ініціює
- Не real-time
- Overhead на кожен запит

---

## 2. Polling

### Short Polling

**Клієнт періодично запитує сервер (наприклад, кожні 5 секунд)**

```
Client                    Server
  |                          |
  |-------- GET /status ---->|
  |<----- 200 OK ------------|
  |                          |
  | (wait 5 seconds)         |
  |                          |
  |-------- GET /status ---->|
  |<----- 200 OK ------------|
  |                          |
```

### Go Client Example

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type Status struct {
    Online bool   `json:"online"`
    Users  int    `json:"users"`
}

func pollServer() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            resp, err := http.Get("https://api.example.com/status")
            if err != nil {
                fmt.Println("Error:", err)
                continue
            }
            
            var status Status
            json.NewDecoder(resp.Body).Decode(&status)
            resp.Body.Close()
            
            fmt.Printf("Status: Online=%v, Users=%d\n", status.Online, status.Users)
        }
    }
}

func main() {
    pollServer()
}
```

### ✅ Переваги
- Простий в реалізації
- Працює з будь-яким HTTP сервером

### ❌ Недоліки
- Неефективний (багато пустих запитів)
- Затримка (залежить від polling interval)
- Навантаження на сервер

---

## 3. Long Polling

### Як працює

**Клієнт робить запит → Сервер тримає з'єднання відкритим до появи нових даних**

```
Client                    Server
  |                          |
  |-------- GET /events ---->|
  |         (waiting...)     | (holds connection)
  |                          | (new data arrives)
  |<----- 200 OK + Data -----|
  |                          |
  |-------- GET /events ---->| (immediately reconnect)
  |         (waiting...)     |
```

### Go Server Example

```go
package main

import (
    "encoding/json"
    "net/http"
    "time"
)

type Event struct {
    Type string `json:"type"`
    Data string `json:"data"`
}

var eventChan = make(chan Event, 10)

func longPollHandler(w http.ResponseWriter, r *http.Request) {
    // Set timeout (e.g., 30 seconds)
    timeout := time.After(30 * time.Second)
    
    select {
    case event := <-eventChan:
        // Got new event, send it
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(event)
        
    case <-timeout:
        // Timeout, send empty response
        w.WriteHeader(http.StatusNoContent)
        
    case <-r.Context().Done():
        // Client disconnected
        return
    }
}

func main() {
    http.HandleFunc("/events", longPollHandler)
    
    // Simulate events
    go func() {
        for {
            time.Sleep(10 * time.Second)
            eventChan <- Event{Type: "update", Data: "New data"}
        }
    }()
    
    http.ListenAndServe(":8080", nil)
}
```

### Go Client Example

```go
func longPollClient() {
    for {
        resp, err := http.Get("http://localhost:8080/events")
        if err != nil {
            fmt.Println("Error:", err)
            time.Sleep(1 * time.Second)
            continue
        }
        
        if resp.StatusCode == http.StatusOK {
            var event Event
            json.NewDecoder(resp.Body).Decode(&event)
            fmt.Printf("Event: %+v\n", event)
        }
        
        resp.Body.Close()
        
        // Immediately reconnect
    }
}
```

### ✅ Переваги
- Майже real-time
- Менше запитів ніж polling
- Працює через HTTP

### ❌ Недоліки
- Тримає з'єднання відкритим
- Складніший ніж short polling
- Scaling issues

---

## 4. WebSockets

### Full-Duplex Communication

**Постійне двостороннє з'єднання між клієнтом і сервером**

```
Client                    Server
  |                          |
  |------ WS Handshake ----->|
  |<---- WS Accept ----------|
  |                          |
  |====== Connected =========|
  |                          |
  |<----- Message 1 ---------|
  |------ Message 2 -------->|
  |<----- Message 3 ---------|
  |                          |
```

### Go Server Example (gorilla/websocket)

```go
package main

import (
    "fmt"
    "net/http"
    
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Allow all origins
    },
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
    // Upgrade HTTP to WebSocket
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println("Upgrade error:", err)
        return
    }
    defer conn.Close()
    
    fmt.Println("Client connected")
    
    // Read messages
    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            fmt.Println("Read error:", err)
            break
        }
        
        fmt.Printf("Received: %s\n", message)
        
        // Echo back
        err = conn.WriteMessage(messageType, message)
        if err != nil {
            fmt.Println("Write error:", err)
            break
        }
    }
    
    fmt.Println("Client disconnected")
}

func main() {
    http.HandleFunc("/ws", wsHandler)
    
    fmt.Println("WebSocket server on :8080")
    http.ListenAndServe(":8080", nil)
}
```

### Go Client Example

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/gorilla/websocket"
)

func main() {
    // Connect to WebSocket
    conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
    if err != nil {
        log.Fatal("Dial error:", err)
    }
    defer conn.Close()
    
    fmt.Println("Connected to server")
    
    // Send messages
    go func() {
        for i := 0; i < 5; i++ {
            message := fmt.Sprintf("Message %d", i)
            err := conn.WriteMessage(websocket.TextMessage, []byte(message))
            if err != nil {
                fmt.Println("Write error:", err)
                return
            }
            fmt.Printf("Sent: %s\n", message)
            time.Sleep(2 * time.Second)
        }
    }()
    
    // Read messages
    for {
        _, message, err := conn.ReadMessage()
        if err != nil {
            fmt.Println("Read error:", err)
            break
        }
        fmt.Printf("Received: %s\n", message)
    }
}
```

### Real-World Example: Chat Room

```go
package main

import (
    "github.com/gorilla/websocket"
    "net/http"
    "sync"
)

type Client struct {
    conn *websocket.Conn
    send chan []byte
}

type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    mu         sync.Mutex
}

func newHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client] = true
            h.mu.Unlock()
            
        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
            h.mu.Unlock()
            
        case message := <-h.broadcast:
            h.mu.Lock()
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
            h.mu.Unlock()
        }
    }
}

func (c *Client) readPump(hub *Hub) {
    defer func() {
        hub.unregister <- c
        c.conn.Close()
    }()
    
    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            break
        }
        hub.broadcast <- message
    }
}

func (c *Client) writePump() {
    defer c.conn.Close()
    
    for message := range c.send {
        err := c.conn.WriteMessage(websocket.TextMessage, message)
        if err != nil {
            break
        }
    }
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    
    client := &Client{conn: conn, send: make(chan []byte, 256)}
    hub.register <- client
    
    go client.writePump()
    go client.readPump(hub)
}

func main() {
    hub := newHub()
    go hub.run()
    
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        serveWs(hub, w, r)
    })
    
    http.ListenAndServe(":8080", nil)
}
```

### ✅ Переваги
- **Real-time** двостороння комунікація
- Низька latency
- Ефективний (одне з'єднання)
- Підтримка binary data

### ❌ Недоліки
- Складніший в реалізації
- Не кешується
- Потребує підтримки proxy/load balancer
- Не працює через деякі firewalls

---

## 5. Server-Sent Events (SSE)

### One-Way Server → Client

**Сервер надсилає події клієнту через відкрите HTTP з'єднання**

```
Client                    Server
  |                          |
  |------ GET /events ------>|
  |<----- Headers -----------|
  |                          |
  |<===== Event Stream ======|
  |<----- Event 1 -----------|
  |<----- Event 2 -----------|
  |<----- Event 3 -----------|
  |                          |
```

### Go Server Example

```go
package main

import (
    "fmt"
    "net/http"
    "time"
)

func sseHandler(w http.ResponseWriter, r *http.Request) {
    // Set SSE headers
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    
    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
        return
    }
    
    // Send events
    ticker := time.NewTicker(2 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            // Send event
            fmt.Fprintf(w, "data: {\"time\": \"%s\"}\n\n", time.Now().Format(time.RFC3339))
            flusher.Flush()
            
        case <-r.Context().Done():
            // Client disconnected
            return
        }
    }
}

func main() {
    http.HandleFunc("/events", sseHandler)
    fmt.Println("SSE server on :8080")
    http.ListenAndServe(":8080", nil)
}
```

### JavaScript Client Example

```javascript
const eventSource = new EventSource('http://localhost:8080/events');

eventSource.onmessage = (event) => {
    const data = JSON.parse(event.data);
    console.log('Received:', data);
};

eventSource.onerror = (error) => {
    console.error('SSE error:', error);
};

// Close connection
// eventSource.close();
```

### Go Client Example

```go
package main

import (
    "bufio"
    "fmt"
    "net/http"
    "strings"
)

func sseClient() {
    resp, err := http.Get("http://localhost:8080/events")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    reader := bufio.NewReader(resp.Body)
    
    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            break
        }
        
        if strings.HasPrefix(line, "data: ") {
            data := strings.TrimPrefix(line, "data: ")
            fmt.Printf("Received: %s", data)
        }
    }
}

func main() {
    sseClient()
}
```

### ✅ Переваги
- Простіше ніж WebSocket
- Автоматичне reconnect
- Працює через HTTP
- Event IDs для resumed connections

### ❌ Недоліки
- Тільки Server → Client
- Text only (no binary)
- Обмеження на кількість з'єднань (browser limit)

---

## 6. gRPC

### HTTP/2 Based RPC

```
Client                    Server
  |                          |
  |------ gRPC Call -------->|
  |  (binary protobuf)       |
  |<----- Response ----------|
  |  (binary protobuf)       |
```

### Example в Week 6

Детальніше в `/Users/vkuzm/GolandProjects/golang_practice/week_6/theory/07_goroutines_concurrency.md`

---

## 7. GraphQL

### Flexible Queries

```graphql
query {
  user(id: 1) {
    name
    email
    posts {
      title
    }
  }
}
```

---

## Порівняння підходів

| Pattern       | Use Case                  | Latency | Complexity | Scalability |
|---------------|---------------------------|---------|------------|-------------|
| HTTP REST     | Standard APIs             | Medium  | Low        | High        |
| Short Polling | Simple updates            | High    | Low        | Medium      |
| Long Polling  | Near real-time            | Medium  | Medium     | Medium      |
| WebSocket     | Real-time chat, games     | Low     | High       | Medium      |
| SSE           | Live feeds, notifications | Low     | Medium     | High        |
| gRPC          | Microservices             | Very Low| Medium     | High        |

---

## Вибір підходу

### Use HTTP REST коли:
✅ Standard CRUD operations
✅ Public API
✅ Need caching
✅ Simple requirements

### Use WebSocket коли:
✅ Real-time chat
✅ Multiplayer games
✅ Live collaboration (Google Docs)
✅ Trading platforms

### Use SSE коли:
✅ Live notifications
✅ News feeds
✅ Stock tickers
✅ One-way updates

### Use Long Polling коли:
✅ Can't use WebSocket
✅ Need fallback
✅ Simple real-time updates

### Use gRPC коли:
✅ Microservices communication
✅ Need high performance
✅ Strong typing required
✅ Internal APIs

---

**Вибір залежить від вимог!** 🚀
