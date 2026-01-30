# Week 16 - OSI Model

## 🎯 Мета

Розуміння 7 рівнів OSI моделі: PDU, адресація, пристрої та протоколи.

---

## 📊 OSI Model - 7 Layers

| Layer | Unit | Address | Device | Protocol |
|-------|------|---------|--------|----------|
| **7. Application** | Data | - | - | HTTP, DNS, SMTP |
| **6. Presentation** | Data | - | - | SSL/TLS |
| **5. Session** | Data | - | - | NetBIOS |
| **4. Transport** | Segment | Port | - | TCP, UDP |
| **3. Network** | Packet | IP | Router | IPv4, IPv6, ICMP |
| **2. Data Link** | Frame | MAC | Switch | Ethernet, WiFi |
| **1. Physical** | Bit | - | Hub, Cable | - |

---

## 📚 Layer Descriptions

### 7. Application Layer

**Unit:** Data  
**Address:** None  
**Device:** None  
**Protocols:** HTTP, HTTPS, DNS, SMTP, FTP, SSH  

**Роль:** Інтерфейс між користувачем та мережею.

**Приклади:**
- HTTP: Веб-браузери, API
- DNS: Розв'язання доменів (google.com → 142.250.185.46)
- SMTP: Відправка email

---

### 6. Presentation Layer

**Unit:** Data  
**Address:** None  
**Device:** None  
**Protocols:** SSL/TLS, JPEG, ASCII, UTF-8  

**Роль:** Шифрування, кодування, стиснення даних.

**Приклади:**
- SSL/TLS: HTTPS шифрування
- Data encoding: JSON, XML, Base64
- Compression: gzip

---

### 5. Session Layer

**Unit:** Data  
**Address:** None  
**Device:** None  
**Protocols:** NetBIOS, RPC, SQL  

**Роль:** Управління сесіями (встановлення, підтримка, завершення).

**Приклади:**
- HTTP sessions з cookies
- Database connections
- RPC calls

---

### 4. Transport Layer

**Unit:** Segment  
**Address:** Port (0-65535)  
**Device:** None  
**Protocols:** TCP, UDP  

**Роль:** End-to-end комунікація, надійність доставки.

**TCP vs UDP:**
```
TCP (Transmission Control Protocol):
✅ Reliable (guaranteed delivery)
✅ Ordered (packets arrive in order)
✅ Connection-oriented (3-way handshake)
❌ Slower (overhead)
Use: HTTP, SMTP, SSH, databases

UDP (User Datagram Protocol):
✅ Fast (no overhead)
✅ Connectionless
❌ Unreliable (no guarantees)
❌ Unordered
Use: DNS, video streaming, gaming
```

**Port Examples:**
- 80: HTTP
- 443: HTTPS
- 22: SSH
- 25: SMTP
- 53: DNS

---

### 3. Network Layer

**Unit:** Packet  
**Address:** IP address  
**Device:** Router  
**Protocols:** IPv4, IPv6, ICMP  

**Роль:** Routing між мережами, logical addressing.

**IP Addresses:**
```
IPv4: 192.168.1.1 (32-bit, 4 billion addresses)
IPv6: 2001:0db8::1 (128-bit, 340 undecillion addresses)
```

**ICMP:**
- `ping` - перевірка доступності
- `traceroute` - шлях до хосту

**Router:** Пересилає пакети між різними мережами.

---

### 2. Data Link Layer

**Unit:** Frame  
**Address:** MAC address  
**Device:** Switch  
**Protocols:** Ethernet, WiFi, ARP  

**Роль:** Point-to-point комунікація, MAC addressing.

**MAC Address:**
```
Format: AA:BB:CC:DD:EE:FF (48-bit, 6 bytes)
Example: 00:1A:2B:3C:4D:5E
Unique per network interface card (NIC)
```

**ARP (Address Resolution Protocol):**
```
IP → MAC mapping
Example: 192.168.1.1 → AA:BB:CC:DD:EE:FF
```

**Switch:** Пересилає frames за MAC адресами в локальній мережі.

---

### 1. Physical Layer

**Unit:** Bit (0s and 1s)  
**Address:** None  
**Device:** Hub, Cable, NIC  
**Protocols:** None (hardware)  

**Роль:** Передача raw bits через фізичне середовище.

**Media Types:**
- Ethernet cable (Cat5e, Cat6)
- Fiber optic
- WiFi (radio waves)
- Bluetooth

**Hub:** Broadcast пристрій (посилає сигнал всім портам).

---

## 🎯 Data Flow Example

### Sending HTTP Request

```
[7] Application    → "GET / HTTP/1.1"
[6] Presentation   → Add SSL/TLS encryption
[5] Session        → Maintain session
[4] Transport      → Add TCP header [Port: 443]
[3] Network        → Add IP header [IP: 142.250.185.46]
[2] Data Link      → Add Ethernet header [MAC: AA:BB:CC:DD:EE:FF]
[1] Physical       → Convert to electrical signals → Send
```

### Receiving Response

```
[1] Physical       → Receive signals → Convert to bits
[2] Data Link      → Read MAC, verify, strip header
[3] Network        → Read IP, verify, strip header
[4] Transport      → Read Port, reassemble, strip header
[5] Session        → Continue session
[6] Presentation   → Decrypt SSL/TLS
[7] Application    → "HTTP/1.1 200 OK"
```

---

## 🔧 Common Ports

| Port | Protocol | Service |
|------|----------|---------|
| 20/21 | FTP | File Transfer |
| 22 | SSH | Secure Shell |
| 25 | SMTP | Email (send) |
| 53 | DNS | Domain resolution |
| 80 | HTTP | Web |
| 110 | POP3 | Email (receive) |
| 143 | IMAP | Email (receive) |
| 443 | HTTPS | Secure Web |
| 3306 | MySQL | Database |
| 5432 | PostgreSQL | Database |
| 6379 | Redis | Cache |
| 27017 | MongoDB | Database |

---

## 🛠️ Troubleshooting Tools

| Layer | Tool | Command | Purpose |
|-------|------|---------|---------|
| 7 | curl | `curl -v https://google.com` | Test HTTP |
| 7 | nslookup | `nslookup google.com` | DNS lookup |
| 4 | telnet | `telnet google.com 80` | Test TCP port |
| 3 | ping | `ping google.com` | Test reachability |
| 3 | traceroute | `traceroute google.com` | Trace route |
| 2 | arp | `arp -a` | View ARP table |
| 2 | ifconfig | `ifconfig` | View MAC address |
| 1 | ethtool | `ethtool eth0` | Check cable |

---

## 💻 Go Examples

### TCP Client (Layer 4)

```go
package main

import (
	"fmt"
	"net"
)

func main() {
	// Connect to TCP port
	conn, _ := net.Dial("tcp", "google.com:80")
	defer conn.Close()

	// Send HTTP request
	fmt.Fprintf(conn, "GET / HTTP/1.1\r\nHost: google.com\r\n\r\n")

	// Read response
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	fmt.Println(string(buf[:n]))
}
```

### DNS Lookup (Layer 7)

```go
package main

import (
	"fmt"
	"net"
)

func main() {
	// Lookup IP addresses
	ips, _ := net.LookupIP("google.com")
	for _, ip := range ips {
		fmt.Println(ip)
	}
}
```

### UDP Client (Layer 4)

```go
package main

import (
	"fmt"
	"net"
)

func main() {
	// UDP connection
	conn, _ := net.Dial("udp", "8.8.8.8:53")
	defer conn.Close()

	// Send DNS query (simplified)
	conn.Write([]byte("query"))

	// Read response
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	fmt.Println(string(buf[:n]))
}
```

---

## 📚 Quick References

- [OSI Cheat Sheet](./OSI_CHEAT_SHEET.md) - візуальний довідник
- [TCP vs UDP](./TCP_VS_UDP.md) - порівняння
- [Port Numbers](./PORT_NUMBERS.md) - список портів

---

## ✅ Learning Checklist

- [ ] Назвати всі 7 рівнів OSI моделі
- [ ] Знати PDU для кожного рівня (bit, frame, packet, segment, data)
- [ ] Пояснити різницю між IP та MAC адресами
- [ ] Знати різницю між TCP та UDP
- [ ] Розуміти роль Router (Layer 3) vs Switch (Layer 2)
- [ ] Знати основні порти (80, 443, 22, 25, 53)
- [ ] Використовувати ping, traceroute, telnet для troubleshooting

---

**Week 16: OSI Model!** 🌐📡
