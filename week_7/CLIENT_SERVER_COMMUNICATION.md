# 🎉 Створено: Client-Server Communication

## 📖 Що це?

Детальний файл про **всі способи, якими клієнт може питати сервер** та отримувати дані.

**Файл:** `theory/11_client_server_communication.md`

---

## 7 Способів комунікації

### 1. **HTTP Request-Response** (класичний)
```
Клієнт запитує → Сервер відповідає
```
- GET, POST, PUT, DELETE
- Найпростіший підхід
- Використовується скрізь

### 2. **Short Polling** (періодичні запити)
```
Клієнт запитує кожні 5 секунд
```
- Просто реалізувати
- Неефективно (багато пустих запитів)

### 3. **Long Polling** (тримає з'єднання)
```
Клієнт запитує → Сервер чекає до події → Відповідає
```
- Майже real-time
- Менше запитів

### 4. **WebSocket** (постійне з'єднання) ⭐
```
Клієнт ↔ Сервер (двостороння комунікація)
```
- **Real-time chat**
- Multiplayer games
- Live collaboration
- **Повний приклад Chat Room в файлі!**

### 5. **Server-Sent Events (SSE)**
```
Сервер → Клієнт (одностороння)
```
- Live notifications
- News feeds
- Stock tickers

### 6. **gRPC**
```
Клієнт → Сервер (binary, HTTP/2)
```
- Microservices
- High performance

### 7. **GraphQL**
```
Клієнт → Flexible query → Сервер
```
- Frontend flexibility

---

## 💻 Приклади коду для кожного способу

### WebSocket Chat Room (повний приклад!)

```go
// Server
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}

// Broadcast to all clients
for client := range h.clients {
    client.send <- message
}
```

### SSE (Server-Sent Events)

```go
func sseHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/event-stream")
    
    for {
        fmt.Fprintf(w, "data: {\"time\": \"%s\"}\n\n", time.Now())
        flusher.Flush()
        time.Sleep(2 * time.Second)
    }
}
```

### Long Polling

```go
func longPollHandler(w http.ResponseWriter, r *http.Request) {
    timeout := time.After(30 * time.Second)
    
    select {
    case event := <-eventChan:
        json.NewEncoder(w).Encode(event)
    case <-timeout:
        w.WriteHeader(http.StatusNoContent)
    }
}
```

---

## 📊 Порівняльна таблиця

| Спосіб         | Real-time | Складність | Use Case              |
|----------------|-----------|------------|-----------------------|
| HTTP REST      | ❌        | ⭐         | Standard APIs         |
| Short Polling  | ❌        | ⭐         | Simple updates        |
| Long Polling   | ⚠️        | ⭐⭐       | Near real-time        |
| **WebSocket**  | ✅        | ⭐⭐⭐     | **Chat, Games**       |
| SSE            | ✅        | ⭐⭐       | Notifications, Feeds  |
| gRPC           | ✅        | ⭐⭐       | Microservices         |

---

## 🎯 Коли що використовувати?

### HTTP REST ✅
- Standard CRUD
- Public API
- Need caching

### WebSocket ✅ (найпопулярніше для real-time)
- **Chat додатки**
- Multiplayer games
- Live collaboration
- Trading platforms

### SSE ✅
- Live notifications
- News feeds
- One-way updates

### gRPC ✅
- Microservices
- Internal APIs
- High performance

---

## 📖 Як читати файл

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_7

# Прочитайте файл
cat theory/11_client_server_communication.md

# Або відкрийте в редакторі
```

---

## 🚀 Що далі?

1. Прочитайте файл повністю
2. Спробуйте WebSocket приклад (повний Chat Room!)
3. Спробуйте SSE приклад
4. Порівняйте різні підходи

---

## ✨ Highlights

### Повний Chat Room на WebSocket

```go
// Hub manages all connected clients
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}

// Client reads & writes messages
type Client struct {
    conn *websocket.Conn
    send chan []byte
}

// Run in main
hub := newHub()
go hub.run()

http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
    serveWs(hub, w, r)
})
```

### SSE з автоматичним reconnect

```javascript
const eventSource = new EventSource('/events');

eventSource.onmessage = (event) => {
    console.log('New event:', JSON.parse(event.data));
};
```

---

**Тепер ви знаєте всі способи комунікації клієнт-сервер!** 🎉

**Файл:** `theory/11_client_server_communication.md`  
**Обсяг:** ~2,500 слів  
**Приклади:** 15+ робочих snippets  
**Статус:** ✅ Готовий до використання
