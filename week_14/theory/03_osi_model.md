# OSI Model (Модель OSI)

## 🎯 Що таке OSI Model?

**OSI (Open Systems Interconnection)** - еталонна модель мережевої взаємодії з **7 рівнів**.

---

## 📊 7 Рівнів OSI

```
┌─────────────────────────┐
│ 7. Application Layer    │ ← HTTP, FTP, DNS, SMTP
├─────────────────────────┤
│ 6. Presentation Layer   │ ← SSL/TLS, encryption, JPEG, MP3
├─────────────────────────┤
│ 5. Session Layer        │ ← Sessions, authentication
├─────────────────────────┤
│ 4. Transport Layer      │ ← TCP, UDP, ports
├─────────────────────────┤
│ 3. Network Layer        │ ← IP, routing, packets
├─────────────────────────┤
│ 2. Data Link Layer      │ ← MAC, switches, frames
├─────────────────────────┤
│ 1. Physical Layer       │ ← Cables, bits, signals
└─────────────────────────┘
```

**Мнемоніка (знизу вгору):**
**P**lease **D**o **N**ot **T**hrow **S**ausage **P**izza **A**way

---

## 1️⃣ Physical Layer (Фізичний)

**Що робить:** Передача **бітів** (0 та 1) через фізичне середовище

**Компоненти:**
- Cables (Ethernet, fiber optic)
- Hubs
- Repeaters
- Radio waves (WiFi)

**Приклад:**
```
Біти: 10110010
       ↓
Electrical signal: ─┐ ┌─┐ ┌┐ ┌─
                    └─┘ └┘└─┘
```

**Проблеми:**
- Cable unplugged
- Signal interference
- Distance limits

---

## 2️⃣ Data Link Layer (Канальний)

**Що робить:** Передача **кадрів (frames)** між сусідніми вузлами, MAC addresses

**Компоненти:**
- MAC addresses (48-bit, e.g., `AA:BB:CC:DD:EE:FF`)
- Switches
- Bridges
- Ethernet, WiFi

**Frame структура:**
```
┌──────────┬──────────┬──────┬─────┐
│ Dest MAC │ Src MAC  │ Data │ CRC │
└──────────┴──────────┴──────┴─────┘
```

**Приклад:**
```
Source: AA:BB:CC:DD:EE:FF (твій комп'ютер)
Dest:   11:22:33:44:55:66 (роутер)
```

**Проблеми:**
- MAC address conflicts
- Switch failures

---

## 3️⃣ Network Layer (Мережевий)

**Що робить:** Маршрутизація **пакетів** між мережами, IP addresses

**Компоненти:**
- IP addresses (IPv4: `192.168.1.1`, IPv6: `2001:db8::1`)
- Routers
- IP protocol
- ICMP (ping)

**Packet структура:**
```
┌─────────────┬─────────────┬──────────┐
│ Source IP   │ Dest IP     │ Payload  │
└─────────────┴─────────────┴──────────┘
```

**Приклад:**
```
From: 192.168.1.100 (твій комп'ютер)
To:   8.8.8.8 (Google DNS)

Router визначає: "Через який шлях відправити?"
```

**Проблеми:**
- IP conflicts
- Routing loops
- Unreachable networks

---

## 4️⃣ Transport Layer (Транспортний)

**Що робить:** Надійна доставка даних між **процесами**, ports

**Протоколи:**

### TCP (Transmission Control Protocol)
```
✅ Надійний (гарантує доставку)
✅ З'єднання (handshake)
✅ Порядок пакетів
❌ Повільніше

Port: 80 (HTTP), 443 (HTTPS), 22 (SSH)
```

### UDP (User Datagram Protocol)
```
❌ Ненадійний (може загубити)
✅ Без з'єднання
✅ Швидкий
✅ Live streaming, games

Port: 53 (DNS), 123 (NTP), game servers
```

**Приклад TCP handshake:**
```
Client → Server: SYN (Хочу з'єднатись)
Server → Client: SYN-ACK (OK, я готовий)
Client → Server: ACK (Починаємо!)
```

**Segment структура:**
```
┌────────────┬────────────┬──────────┐
│ Source Port│ Dest Port  │ Data     │
└────────────┴────────────┴──────────┘
```

---

## 5️⃣ Session Layer (Сеансовий)

**Що робить:** Управління **сесіями** між додатками

**Функції:**
- Встановлення сесії
- Підтримка сесії
- Завершення сесії
- Checkpointing (відновлення після розриву)

**Приклад:**
```
1. Login → Create session
2. Keep-alive → Maintain session
3. Logout → Close session
```

**Протоколи:**
- NetBIOS
- RPC
- PPTP

**В реальності:** Часто комбінується з Application Layer

---

## 6️⃣ Presentation Layer (Представлення)

**Що робить:** Форматування, шифрування, стиснення даних

**Функції:**
- **Encoding:** ASCII, UTF-8, EBCDIC
- **Encryption:** SSL/TLS, AES
- **Compression:** GZIP, JPEG, MP3
- **Serialization:** JSON, XML, Protocol Buffers

**Приклад:**
```
Application: "Hello"
    ↓
Presentation: Encrypt with TLS → 0x8f3a2b...
    ↓
Session: Send encrypted data
```

**В реальності:** Часто частина Application Layer (наприклад, HTTPS = HTTP + TLS)

---

## 7️⃣ Application Layer (Прикладний)

**Що робить:** Взаємодія з **користувачем** та додатками

**Протоколи:**

### HTTP/HTTPS (Web)
```
GET /api/users HTTP/1.1
Host: example.com
```

### DNS (Domain Name System)
```
example.com → 93.184.216.34
```

### SMTP (Email sending)
```
From: user@gmail.com
To: friend@yahoo.com
Subject: Hello
```

### FTP (File Transfer)
```
PUT file.txt → Server
```

### SSH (Secure Shell)
```
ssh user@server
```

---

## 🌐 Real-World Example: Відкрити сайт

### Browser → example.com

```
7. Application:
   HTTP GET request: "GET / HTTP/1.1"

6. Presentation:
   Encrypt with TLS (HTTPS)
   
5. Session:
   Create TCP session

4. Transport:
   TCP segment with port 443 (HTTPS)
   SYN → SYN-ACK → ACK

3. Network:
   Add IP header
   Source: 192.168.1.100
   Dest: 93.184.216.34 (example.com)

2. Data Link:
   Add MAC header
   Source MAC: AA:BB:CC:DD:EE:FF
   Dest MAC: Router's MAC

1. Physical:
   Convert to electrical signals
   Send through Ethernet cable
```

### Відповідь від сервера (зворотній шлях):

```
1. Physical: Receive bits
2. Data Link: Parse frame, check MAC
3. Network: Parse packet, check IP
4. Transport: Parse segment, check port
5. Session: Associate with existing session
6. Presentation: Decrypt TLS
7. Application: Parse HTTP response, render HTML
```

---

## 🎯 TCP/IP Model vs OSI

**TCP/IP (4 layers)** - практична модель, що використовується:

```
OSI                    TCP/IP
─────────────────     ─────────────────
7. Application   ┐
6. Presentation  ├──► Application
5. Session       ┘
4. Transport     ───► Transport (TCP/UDP)
3. Network       ───► Internet (IP)
2. Data Link     ┐
1. Physical      ├──► Link (Ethernet, WiFi)
                 ┘
```

---

## 📊 Port Numbers (Transport Layer)

### Well-Known Ports (0-1023)

| Port | Protocol | Service |
|------|----------|---------|
| 20/21 | FTP | File Transfer |
| 22 | SSH | Secure Shell |
| 23 | Telnet | Remote login |
| 25 | SMTP | Email sending |
| 53 | DNS | Domain names |
| 80 | HTTP | Web |
| 110 | POP3 | Email receiving |
| 143 | IMAP | Email |
| 443 | HTTPS | Secure web |
| 3306 | MySQL | Database |
| 5432 | PostgreSQL | Database |
| 6379 | Redis | Cache |
| 8080 | HTTP-alt | Web (dev) |

---

## 🔍 Debugging по рівнях

### Physical (1)
```bash
# Check cable
ip link show

# Ethernet connected?
ethtool eth0
```

### Data Link (2)
```bash
# MAC address
ip link show

# ARP table (MAC ↔ IP)
arp -a
```

### Network (3)
```bash
# Ping (ICMP)
ping 8.8.8.8

# Trace route
traceroute google.com

# IP config
ip addr show
ifconfig
```

### Transport (4)
```bash
# Open ports
netstat -tuln
ss -tuln

# TCP connections
netstat -an | grep ESTABLISHED
```

### Application (7)
```bash
# HTTP request
curl -v https://example.com

# DNS lookup
nslookup example.com
dig example.com

# Test port
telnet example.com 80
nc -zv example.com 80
```

---

## 🎯 Common Problems

### "Cannot connect to server"

```
7. Application: ❓
6. Presentation: ❓
5. Session: ❓
4. Transport: ❓ Port closed? Firewall?
3. Network: ❓ Can ping server?
2. Data Link: ❓ MAC address resolved?
1. Physical: ❓ Cable connected?
```

**Debugging bottom-up:**
```bash
# 1. Physical
ip link show  # Interface UP?

# 2. Data Link
arp -a  # MAC addresses visible?

# 3. Network
ping <server_ip>  # Reachable?

# 4. Transport
telnet <server_ip> <port>  # Port open?

# 7. Application
curl http://<server_ip>  # HTTP works?
```

---

## 🔐 Security по рівнях

| Layer | Attacks | Defense |
|-------|---------|---------|
| 1. Physical | Wire tapping | Physical security |
| 2. Data Link | MAC spoofing, ARP poisoning | Port security, VLAN |
| 3. Network | IP spoofing, DoS | Firewall, IDS/IPS |
| 4. Transport | Port scanning, SYN flood | Firewall rules |
| 5. Session | Session hijacking | Encryption, tokens |
| 6. Presentation | Man-in-the-middle | TLS/SSL, certificates |
| 7. Application | SQL injection, XSS | Input validation, WAF |

---

## ✅ Висновок

### OSI Model:

✅ **7 layers** - від physical до application  
✅ **Each layer** - specific responsibility  
✅ **Encapsulation** - кожен рівень додає header  
✅ **De-encapsulation** - кожен рівень видаляє header  

### Key Layers:

**Physical (1)** - bits, cables  
**Data Link (2)** - frames, MAC  
**Network (3)** - packets, IP  
**Transport (4)** - segments, TCP/UDP, ports  
**Application (7)** - HTTP, DNS, SMTP  

### Golden Rule:

**"Please Do Not Throw Sausage Pizza Away"** = 7 layers! 📡

---

## 📖 Go Example

```go
// Application Layer - HTTP client
resp, _ := http.Get("https://example.com")
// ↓ Go handles:
// - TLS (Presentation)
// - TCP connection (Transport)
// - IP routing (Network)
// - Ethernet (Data Link)
// - Physical transmission

// You only work at Application layer!
```

**Week 14: OSI Model + JOINs!** 🌐🗄️
