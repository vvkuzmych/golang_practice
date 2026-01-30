# OSI Model - Cheat Sheet

## 📊 Complete Table

| # | Layer | Unit | Address | Device | Protocol | Example |
|---|-------|------|---------|--------|----------|---------|
| **7** | Application | Data | - | - | HTTP, DNS, SMTP | Web browser, Email |
| **6** | Presentation | Data | - | - | SSL/TLS | HTTPS encryption |
| **5** | Session | Data | - | - | NetBIOS | Login session |
| **4** | Transport | Segment | Port | - | TCP, UDP | Port 443, 80 |
| **3** | Network | Packet | IP | Router | IPv4, IPv6 | 192.168.1.1 |
| **2** | Data Link | Frame | MAC | Switch | Ethernet | AA:BB:CC:DD:EE:FF |
| **1** | Physical | Bit | - | Hub, Cable | - | Cat6 cable, WiFi |

---

## 🎯 Remember with Mnemonic

### Top-Down (7 → 1)

**"All People Seem To Need Data Processing"**

- **A**ll = Application
- **P**eople = Presentation
- **S**eem = Session
- **T**o = Transport
- **N**eed = Network
- **D**ata = Data Link
- **P**rocessing = Physical

### Bottom-Up (1 → 7)

**"Please Do Not Throw Sausage Pizza Away"**

- **P**lease = Physical
- **D**o = Data Link
- **N**ot = Network
- **T**hrow = Transport
- **S**ausage = Session
- **P**izza = Presentation
- **A**way = Application

---

## 📦 PDU (Protocol Data Unit)

```
Layer 7-5: Data
Layer 4:   Segment (TCP) / Datagram (UDP)
Layer 3:   Packet
Layer 2:   Frame
Layer 1:   Bit
```

---

## 🔢 Addressing Summary

| Layer | Addressing | Format | Example |
|-------|------------|--------|---------|
| 4 - Transport | Port | 16-bit number | 443, 80, 22 |
| 3 - Network | IP address | IPv4 (32-bit) / IPv6 (128-bit) | 192.168.1.1 / 2001:0db8::1 |
| 2 - Data Link | MAC address | 48-bit (6 bytes) | AA:BB:CC:DD:EE:FF |

---

## 🌐 Protocol Mapping

### Layer 7 - Application
```
HTTP/HTTPS   → Web (80/443)
DNS          → Domain resolution (53)
SMTP         → Send email (25)
POP3/IMAP    → Receive email (110/143)
FTP          → File transfer (20/21)
SSH          → Secure shell (22)
```

### Layer 6 - Presentation
```
SSL/TLS      → Encryption
JPEG/PNG     → Image formats
ASCII/UTF-8  → Text encoding
gzip         → Compression
```

### Layer 5 - Session
```
NetBIOS      → Windows networking
RPC          → Remote procedure calls
SQL          → Database sessions
```

### Layer 4 - Transport
```
TCP          → Reliable, ordered (HTTP, SSH, SMTP)
UDP          → Fast, unreliable (DNS, streaming)
```

### Layer 3 - Network
```
IPv4         → 192.168.1.1
IPv6         → 2001:0db8::1
ICMP         → ping, traceroute
```

### Layer 2 - Data Link
```
Ethernet     → Wired LAN
WiFi         → Wireless LAN (802.11)
ARP          → IP → MAC mapping
```

### Layer 1 - Physical
```
Ethernet cable (Cat5e, Cat6)
Fiber optic
WiFi radio waves
```

---

## 🔧 Devices by Layer

| Layer | Device | Function |
|-------|--------|----------|
| 7-5 | Application | Software (browser, email client) |
| 4 | - | OS handles (no dedicated device) |
| 3 | **Router** | Routes packets between networks |
| 2 | **Switch** | Forwards frames by MAC address |
| 1 | **Hub** | Broadcasts bits to all ports |
| 1 | **Cable** | Physical medium |

---

## 🎯 TCP vs UDP Quick Comparison

| Feature | TCP | UDP |
|---------|-----|-----|
| **Reliability** | ✅ Guaranteed | ❌ Best effort |
| **Order** | ✅ Ordered | ❌ Unordered |
| **Connection** | ✅ Connection-oriented | ❌ Connectionless |
| **Speed** | ❌ Slower (overhead) | ✅ Faster |
| **Use Cases** | HTTP, SMTP, SSH, FTP | DNS, streaming, gaming |

---

## 📍 Port Numbers

### Well-Known (0-1023)

```
20/21  → FTP
22     → SSH
23     → Telnet
25     → SMTP
53     → DNS
80     → HTTP
110    → POP3
143    → IMAP
443    → HTTPS
```

### Common Application Ports (1024-65535)

```
3306   → MySQL
5432   → PostgreSQL
6379   → Redis
8080   → Alternative HTTP
27017  → MongoDB
```

---

## 🛠️ Troubleshooting Tools

### Test Connectivity by Layer

```bash
# Layer 7 - Application
curl -v https://google.com         # HTTP test
nslookup google.com                # DNS lookup
dig google.com                     # DNS query

# Layer 4 - Transport
telnet google.com 80               # TCP port test
nc -zv google.com 80               # Netcat port scan

# Layer 3 - Network
ping google.com                    # ICMP test
traceroute google.com              # Route trace
ip route                           # Routing table

# Layer 2 - Data Link
arp -a                             # ARP table
ifconfig / ip addr                 # MAC address
tcpdump                            # Packet capture

# Layer 1 - Physical
ethtool eth0                       # Cable status
iwconfig                           # WiFi status
```

---

## 🎯 Data Encapsulation

### Sending Data (Top → Down)

```
[7] Application    → HTTP Request
[6] Presentation   → Encrypt with TLS
[5] Session        → Add session ID
[4] Transport      → [TCP Header | Data]
[3] Network        → [IP Header | TCP Header | Data]
[2] Data Link      → [Ethernet Header | IP | TCP | Data | CRC]
[1] Physical       → 010101010101... (electrical signals)
```

### Receiving Data (Bottom → Up)

```
[1] Physical       → 010101010101... → Bits
[2] Data Link      → Remove Ethernet header
[3] Network        → Remove IP header
[4] Transport      → Remove TCP header, reassemble
[5] Session        → Verify session
[6] Presentation   → Decrypt TLS
[7] Application    → HTTP Response
```

---

## 📊 Layer Responsibilities

| Layer | Responsibility |
|-------|----------------|
| 7 | User interface, application services |
| 6 | Data format, encryption, compression |
| 5 | Session management |
| 4 | End-to-end communication, ports |
| 3 | Routing, IP addressing |
| 2 | MAC addressing, frame switching |
| 1 | Physical transmission of bits |

---

## 🎓 Quick Quiz

1. **Q:** Which layer uses MAC addresses?  
   **A:** Layer 2 (Data Link)

2. **Q:** Which layer uses IP addresses?  
   **A:** Layer 3 (Network)

3. **Q:** Which layer uses port numbers?  
   **A:** Layer 4 (Transport)

4. **Q:** What device works at Layer 3?  
   **A:** Router

5. **Q:** What device works at Layer 2?  
   **A:** Switch

6. **Q:** TCP or UDP for HTTP?  
   **A:** TCP (reliable, ordered)

7. **Q:** TCP or UDP for DNS?  
   **A:** UDP (fast, simple queries)

8. **Q:** What port does HTTPS use?  
   **A:** 443

---

## 🔍 Common Scenarios

### Scenario 1: Can't access website

```
1. ping google.com
   ✅ Works → Layer 3 OK
   ❌ Fails → Check Layer 1-3 (cable, IP, routing)

2. telnet google.com 80
   ✅ Works → Layer 4 OK
   ❌ Fails → Check firewall, port

3. curl -v https://google.com
   ✅ Works → Layer 7 OK
   ❌ Fails → Check DNS, SSL, application
```

### Scenario 2: Slow connection

```
Layer 1: Check cable, WiFi signal
Layer 2: Check for MAC conflicts
Layer 3: Check routing, packet loss (ping)
Layer 4: Check TCP congestion, retransmissions
```

---

**Week 16: OSI Cheat Sheet!** 🌐✅
