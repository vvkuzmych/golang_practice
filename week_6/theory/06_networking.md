# Нетворкінг: TCP/UDP, HTTP, TLS, DNS

---

## 📖 Зміст

1. [TCP vs UDP](#1-tcp-vs-udp)
2. [TCP Server & Client](#2-tcp-server--client)
3. [HTTP Semantics](#3-http-semantics)
4. [TLS/SSL](#4-tlsssl)
5. [DNS](#5-dns)
6. [Timeouts & Retries](#6-timeouts--retries)

---

## 1. TCP vs UDP

### TCP (Transmission Control Protocol)

**Характеристики:**
- ✅ **Надійний** - гарантує доставку
- ✅ **Ordered** - пакети в правильному порядку
- ✅ **Connection-oriented** - встановлює з'єднання
- ❌ **Повільніший** - через overhead
- ❌ **Більше ресурсів**

**Використання:**
- HTTP/HTTPS
- Email (SMTP, IMAP)
- File transfers (FTP, SSH)
- Databases

### UDP (User Datagram Protocol)

**Характеристики:**
- ✅ **Швидкий** - мінімальний overhead
- ✅ **Легкий** - менше ресурсів
- ❌ **Ненадійний** - може втратити пакети
- ❌ **No ordering** - пакети можуть прийти не по порядку
- ❌ **Connectionless** - немає встановлення з'єднання

**Використання:**
- DNS
- Video streaming
- Online gaming
- VoIP

---

## 2. TCP Server & Client

### TCP Server

```go
package main

import (
    "bufio"
    "fmt"
    "net"
    "strings"
)

func main() {
    // Слухаємо на порту 8080
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        panic(err)
    }
    defer listener.Close()
    
    fmt.Println("TCP Server listening on :8080")
    
    for {
        // Приймаємо з'єднання
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting:", err)
            continue
        }
        
        // Обробляємо кожне з'єднання в окремій goroutine
        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    
    fmt.Printf("Client connected: %s\n", conn.RemoteAddr())
    
    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        message := scanner.Text()
        fmt.Printf("Received: %s\n", message)
        
        // Echo server - відправляємо назад
        response := strings.ToUpper(message) + "\n"
        conn.Write([]byte(response))
    }
    
    fmt.Printf("Client disconnected: %s\n", conn.RemoteAddr())
}
```

### TCP Client

```go
package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
)

func main() {
    // Підключаємось до сервера
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        panic(err)
    }
    defer conn.Close()
    
    fmt.Println("Connected to server")
    
    // Відправляємо повідомлення
    for {
        reader := bufio.NewReader(os.Stdin)
        fmt.Print("Enter message: ")
        message, _ := reader.ReadString('\n')
        
        // Відправляємо
        conn.Write([]byte(message))
        
        // Читаємо відповідь
        response, _ := bufio.NewReader(conn).ReadString('\n')
        fmt.Print("Server response: " + response)
    }
}
```

### UDP Server & Client

**UDP Server:**
```go
func main() {
    addr, _ := net.ResolveUDPAddr("udp", ":8080")
    conn, _ := net.ListenUDP("udp", addr)
    defer conn.Close()
    
    buffer := make([]byte, 1024)
    
    for {
        n, clientAddr, _ := conn.ReadFromUDP(buffer)
        message := string(buffer[:n])
        fmt.Printf("Received from %s: %s\n", clientAddr, message)
        
        // Відповідь
        response := []byte("Received: " + message)
        conn.WriteToUDP(response, clientAddr)
    }
}
```

**UDP Client:**
```go
func main() {
    addr, _ := net.ResolveUDPAddr("udp", "localhost:8080")
    conn, _ := net.DialUDP("udp", nil, addr)
    defer conn.Close()
    
    message := []byte("Hello UDP!")
    conn.Write(message)
    
    buffer := make([]byte, 1024)
    n, _ := conn.Read(buffer)
    fmt.Println("Response:", string(buffer[:n]))
}
```

---

## 3. HTTP Semantics

### HTTP Methods

| Method | Призначення | Idempotent? | Safe? |
|--------|-------------|-------------|-------|
| **GET** | Отримати ресурс | ✅ Yes | ✅ Yes |
| **POST** | Створити ресурс | ❌ No | ❌ No |
| **PUT** | Оновити/замінити | ✅ Yes | ❌ No |
| **PATCH** | Часткове оновлення | ❌ No | ❌ No |
| **DELETE** | Видалити ресурс | ✅ Yes | ❌ No |
| **HEAD** | Метадані (без body) | ✅ Yes | ✅ Yes |
| **OPTIONS** | Підтримувані методи | ✅ Yes | ✅ Yes |

### Status Codes

**2xx Success:**
- `200 OK` - Успіх
- `201 Created` - Ресурс створено
- `204 No Content` - Успіх без body

**3xx Redirection:**
- `301 Moved Permanently` - Постійний redirect
- `302 Found` - Тимчасовий redirect
- `304 Not Modified` - Кеш валідний

**4xx Client Errors:**
- `400 Bad Request` - Невалідний запит
- `401 Unauthorized` - Потрібна автентифікація
- `403 Forbidden` - Немає доступу
- `404 Not Found` - Ресурс не знайдено
- `429 Too Many Requests` - Rate limit

**5xx Server Errors:**
- `500 Internal Server Error` - Помилка сервера
- `502 Bad Gateway` - Помилка проксі
- `503 Service Unavailable` - Сервіс недоступний

### Headers

**Request Headers:**
```
GET /api/users HTTP/1.1
Host: example.com
User-Agent: MyApp/1.0
Accept: application/json
Authorization: Bearer token123
Content-Type: application/json
Cache-Control: no-cache
```

**Response Headers:**
```
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 123
Cache-Control: max-age=3600
ETag: "abc123"
X-RateLimit-Remaining: 99
```

### Content Negotiation

```go
func handler(w http.ResponseWriter, r *http.Request) {
    accept := r.Header.Get("Accept")
    
    switch {
    case strings.Contains(accept, "application/json"):
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(data)
    case strings.Contains(accept, "application/xml"):
        w.Header().Set("Content-Type", "application/xml")
        xml.NewEncoder(w).Encode(data)
    default:
        w.Header().Set("Content-Type", "text/plain")
        fmt.Fprintf(w, "%v", data)
    }
}
```

---

## 4. TLS/SSL

### Що таке TLS?

**TLS (Transport Layer Security)** - протокол шифрування для безпечної комунікації.

```
HTTP  → Незахищений (port 80)
HTTPS → Захищений TLS (port 443)
```

### HTTPS Server

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Secure Hello!")
    })
    
    // HTTPS сервер
    err := http.ListenAndServeTLS(":443", 
        "server.crt",  // Certificate
        "server.key",  // Private key
        nil)
    
    if err != nil {
        panic(err)
    }
}
```

### Генерація self-signed сертифіката

```bash
openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes
```

### HTTPS Client з custom TLS config

```go
func main() {
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{
            InsecureSkipVerify: false, // В продакшн завжди false!
            MinVersion:         tls.VersionTLS12,
        },
    }
    
    client := &http.Client{Transport: tr}
    resp, err := client.Get("https://example.com")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
}
```

### Mutual TLS (mTLS)

```go
// Server вимагає client certificate
func main() {
    caCert, _ := os.ReadFile("ca.crt")
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)
    
    tlsConfig := &tls.Config{
        ClientCAs:  caCertPool,
        ClientAuth: tls.RequireAndVerifyClientCert,
    }
    
    server := &http.Server{
        Addr:      ":443",
        TLSConfig: tlsConfig,
    }
    
    server.ListenAndServeTLS("server.crt", "server.key")
}
```

---

## 5. DNS

### Що таке DNS?

DNS (Domain Name System) перетворює доменні імена на IP адреси.

```
example.com → 93.184.216.34
```

### DNS Lookup в Go

```go
package main

import (
    "fmt"
    "net"
)

func main() {
    // Простий lookup
    ips, err := net.LookupIP("google.com")
    if err != nil {
        panic(err)
    }
    
    for _, ip := range ips {
        fmt.Println(ip)
    }
    
    // Lookup host by IP (reverse DNS)
    names, _ := net.LookupAddr("8.8.8.8")
    fmt.Println(names) // [dns.google]
    
    // MX records (email servers)
    mx, _ := net.LookupMX("gmail.com")
    for _, record := range mx {
        fmt.Printf("%s (priority: %d)\n", record.Host, record.Pref)
    }
    
    // TXT records
    txt, _ := net.LookupTXT("google.com")
    fmt.Println(txt)
}
```

### Custom DNS Resolver

```go
func main() {
    resolver := &net.Resolver{
        PreferGo: true, // Використовувати Go resolver замість системного
        Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
            d := net.Dialer{
                Timeout: time.Second * 5,
            }
            // Використовуємо Google DNS
            return d.DialContext(ctx, network, "8.8.8.8:53")
        },
    }
    
    ips, _ := resolver.LookupIP(context.Background(), "network", "example.com")
    fmt.Println(ips)
}
```

---

## 6. Timeouts & Retries

### Типи Timeouts

```go
client := &http.Client{
    Timeout: 10 * time.Second, // Загальний timeout для запиту
    
    Transport: &http.Transport{
        DialContext: (&net.Dialer{
            Timeout:   5 * time.Second,  // Час на встановлення з'єднання
            KeepAlive: 30 * time.Second, // Keep-alive проби
        }).DialContext,
        
        TLSHandshakeTimeout:   5 * time.Second,  // TLS handshake
        ResponseHeaderTimeout: 10 * time.Second, // Час на отримання headers
        ExpectContinueTimeout: 1 * time.Second,  // Expect: 100-continue
        
        IdleConnTimeout:       90 * time.Second, // Час до закриття idle з'єднання
        MaxIdleConns:          100,              // Максимум idle з'єднань
        MaxIdleConnsPerHost:   10,               // На один host
    },
}
```

### Context Timeout

```go
func makeRequest(url string) error {
    // Timeout через context
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            return fmt.Errorf("request timed out")
        }
        return err
    }
    defer resp.Body.Close()
    
    return nil
}
```

### Retry Logic

```go
func makeRequestWithRetry(url string, maxRetries int) (*http.Response, error) {
    var resp *http.Response
    var err error
    
    for attempt := 0; attempt < maxRetries; attempt++ {
        resp, err = http.Get(url)
        
        if err == nil && resp.StatusCode < 500 {
            return resp, nil // Успіх
        }
        
        if resp != nil {
            resp.Body.Close()
        }
        
        // Exponential backoff
        waitTime := time.Duration(1<<uint(attempt)) * time.Second
        fmt.Printf("Attempt %d failed, retrying in %v...\n", attempt+1, waitTime)
        time.Sleep(waitTime)
    }
    
    return nil, fmt.Errorf("max retries reached: %w", err)
}
```

### Retry з exponential backoff + jitter

```go
func retry(fn func() error, maxRetries int) error {
    for attempt := 0; attempt < maxRetries; attempt++ {
        err := fn()
        if err == nil {
            return nil // Успіх
        }
        
        if attempt == maxRetries-1 {
            return err // Останя спроба
        }
        
        // Exponential backoff з jitter
        backoff := time.Duration(1<<uint(attempt)) * time.Second
        jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
        waitTime := backoff + jitter
        
        fmt.Printf("Attempt %d failed, waiting %v...\n", attempt+1, waitTime)
        time.Sleep(waitTime)
    }
    
    return fmt.Errorf("max retries exceeded")
}

// Використання
err := retry(func() error {
    _, err := http.Get("https://api.example.com/data")
    return err
}, 3)
```

### Circuit Breaker Pattern

```go
type CircuitBreaker struct {
    maxFailures  int
    resetTimeout time.Duration
    failures     int
    lastFailTime time.Time
    state        string // "closed", "open", "half-open"
    mu           sync.Mutex
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    cb.mu.Lock()
    defer cb.mu.Unlock()
    
    // Якщо circuit open і timeout не минув
    if cb.state == "open" {
        if time.Since(cb.lastFailTime) < cb.resetTimeout {
            return fmt.Errorf("circuit breaker is open")
        }
        // Переходимо в half-open
        cb.state = "half-open"
    }
    
    // Виконуємо функцію
    err := fn()
    
    if err != nil {
        cb.failures++
        cb.lastFailTime = time.Now()
        
        if cb.failures >= cb.maxFailures {
            cb.state = "open"
        }
        return err
    }
    
    // Успіх - скидаємо
    cb.failures = 0
    cb.state = "closed"
    return nil
}
```

---

## ✅ Best Practices

1. **Завжди встановлюйте timeouts** - на сервері та клієнті
2. **Використовуйте context** для cancellation
3. **Retry з exponential backoff** - не молотіть сервіс
4. **Circuit Breaker** - захист від каскадних помилок
5. **Connection pooling** - переиспользуйте з'єднання
6. **Graceful Shutdown** - коректне завершення
7. **TLS everywhere** - використовуйте HTTPS
8. **DNS caching** - кешуйте DNS lookup
9. **Health checks** - перевіряйте стан сервісів
10. **Monitoring & Logging** - слідкуйте за мережею

---

**Це завершує теоретичну частину Week 6!**
