# OSI Model Cheat Sheet

## 📊 7 Layers (знизу вгору)

```
┌─────────────────────────────────────────────────────┐
│ 7. APPLICATION  │ HTTP, HTTPS, DNS, SMTP, FTP, SSH  │ ← User programs
├─────────────────┼───────────────────────────────────┤
│ 6. PRESENTATION │ SSL/TLS, JPEG, GZIP, encryption   │ ← Format, encrypt
├─────────────────┼───────────────────────────────────┤
│ 5. SESSION      │ Sessions, authentication          │ ← Connections
├─────────────────┼───────────────────────────────────┤
│ 4. TRANSPORT    │ TCP, UDP, ports                   │ ← End-to-end
├─────────────────┼───────────────────────────────────┤
│ 3. NETWORK      │ IP, routing, packets              │ ← Addressing
├─────────────────┼───────────────────────────────────┤
│ 2. DATA LINK    │ MAC, switches, frames             │ ← Local network
├─────────────────┼───────────────────────────────────┤
│ 1. PHYSICAL     │ Cables, hubs, bits                │ ← Hardware
└─────────────────┴───────────────────────────────────┘
```

**Мнемоніка:** **P**lease **D**o **N**ot **T**hrow **S**ausage **P**izza **A**way

---

## 🎯 Quick Reference

| Layer | Unit | Address | Device | Protocol |
|-------|------|---------|--------|----------|
| 7. Application | Data | - | - | HTTP, DNS, SMTP |
| 6. Presentation | Data | - | - | SSL/TLS |
| 5. Session | Data | - | - | NetBIOS |
| 4. Transport | Segment | Port | - | TCP, UDP |
| 3. Network | Packet | IP | Router | IPv4, IPv6, ICMP |
| 2. Data Link | Frame | MAC | Switch | Ethernet, WiFi |
| 1. Physical | Bit | - | Hub, Cable | - |

---

## 🔧 Transport Layer (4)

### TCP (Transmission Control Protocol)

```
✅ Connection-oriented (handshake)
✅ Reliable (guaranteed delivery)
✅ Ordered (packets in sequence)
✅ Error checking
❌ Slower
```

**Use:** HTTP, HTTPS, SSH, FTP, Email

**3-Way Handshake:**
```
Client → Server: SYN
Server → Client: SYN-ACK
Client → Server: ACK
```

### UDP (User Datagram Protocol)

```
✅ Connectionless (no handshake)
✅ Fast
✅ Low overhead
❌ No reliability guarantee
❌ No ordering
```

**Use:** DNS, DHCP, streaming, games, VoIP

---

## 🌐 Network Layer (3)

### IPv4

```
Format: 192.168.1.1 (32-bit)
Classes: A, B, C
Private: 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16
```

### IPv6

```
Format: 2001:0db8:85a3::8a2e:0370:7334 (128-bit)
```

---

## 📡 Common Ports

### Well-Known (0-1023)

| Port | Protocol | Service |
|------|----------|---------|
| 20-21 | FTP | File Transfer |
| 22 | SSH | Secure Shell |
| 23 | Telnet | Remote login (insecure) |
| 25 | SMTP | Email send |
| 53 | DNS | Domain names |
| 80 | HTTP | Web |
| 110 | POP3 | Email receive |
| 143 | IMAP | Email |
| 443 | HTTPS | Secure web |

### Registered (1024-49151)

| Port | Protocol | Service |
|------|----------|---------|
| 3000 | HTTP | Dev server |
| 3306 | MySQL | Database |
| 5432 | PostgreSQL | Database |
| 6379 | Redis | Cache |
| 8080 | HTTP-alt | Web proxy |
| 27017 | MongoDB | Database |

---

## 🔍 Debugging Commands

### Layer 1: Physical

```bash
# Check interface status
ip link show
ethtool eth0

# Check cable
dmesg | grep eth0
```

### Layer 2: Data Link

```bash
# MAC address
ip link show
ifconfig

# ARP table (IP → MAC)
arp -a
ip neigh show
```

### Layer 3: Network

```bash
# Test connectivity
ping 8.8.8.8
ping google.com

# Trace route
traceroute google.com
tracepath google.com

# IP configuration
ip addr show
ifconfig
ip route show
```

### Layer 4: Transport

```bash
# Open ports
netstat -tuln
ss -tuln
lsof -i

# Active connections
netstat -an | grep ESTABLISHED
ss -tan

# Test port
telnet example.com 80
nc -zv example.com 80
```

### Layer 7: Application

```bash
# HTTP request
curl -v https://example.com
wget -S https://example.com

# DNS lookup
nslookup google.com
dig google.com
host google.com

# Test SMTP
telnet mail.example.com 25
```

---

## 🚨 Troubleshooting Flow

```
Problem: Can't access website
        ↓
1. Physical: Cable connected? WiFi on?
   → ip link show
        ↓
2. Data Link: MAC address OK?
   → arp -a
        ↓
3. Network: Can ping gateway? Can ping 8.8.8.8?
   → ping <gateway_ip>
   → ping 8.8.8.8
        ↓
4. Network: DNS working?
   → nslookup google.com
        ↓
5. Transport: Port open?
   → telnet example.com 80
        ↓
6. Application: HTTP working?
   → curl -v http://example.com
```

---

## 📊 Data Flow Example

### Sending HTTP Request

```
Application (7): Generate HTTP request
                 "GET / HTTP/1.1"
        ↓
Presentation (6): Encrypt with TLS (HTTPS)
        ↓
Session (5): Manage connection
        ↓
Transport (4): Add TCP header (port 443)
               Segment = Header + Data
        ↓
Network (3): Add IP header
             Packet = IP Header + Segment
        ↓
Data Link (2): Add MAC header + trailer
               Frame = MAC Header + Packet + FCS
        ↓
Physical (1): Convert to bits (01010101...)
              Send over cable/WiFi
```

### Receiving Response (reverse)

```
Physical (1): Receive bits
        ↓
Data Link (2): Remove MAC header, check FCS
        ↓
Network (3): Remove IP header, check dest IP
        ↓
Transport (4): Remove TCP header, reassemble
        ↓
Session (5): Associate with session
        ↓
Presentation (6): Decrypt TLS
        ↓
Application (7): Parse HTTP, display webpage
```

---

## 🔐 Security by Layer

| Layer | Attack | Defense |
|-------|--------|---------|
| 1. Physical | Wiretapping | Physical security |
| 2. Data Link | MAC spoofing, ARP poisoning | Port security, VLAN |
| 3. Network | IP spoofing, DoS | Firewall, IDS/IPS |
| 4. Transport | Port scan, SYN flood | Firewall rules, rate limiting |
| 5. Session | Session hijacking | Encryption, secure tokens |
| 6. Presentation | MITM | TLS/SSL, certificates |
| 7. Application | SQL injection, XSS | Input validation, WAF |

---

## ⚡ Quick Tips

### Remember Layers
```
All People Seem To Need Data Processing
(Application, Presentation, Session, Transport, Network, Data Link, Physical)
```

### TCP vs UDP
```
TCP = Phone call (connection)
UDP = Postcard (no connection)
```

### IP vs MAC
```
IP = Postal address (routable)
MAC = House number (local only)
```

### Ports
```
< 1024 = System (need root)
≥ 1024 = User
```

---

## 🎓 Study Tips

### Bottom-Up Approach
1. Start with Physical (cables, bits)
2. Build up to Application (HTTP, email)
3. Understand encapsulation at each layer

### Real-World Mapping
- **Physical:** Your WiFi router
- **Data Link:** Your MAC address
- **Network:** Your IP (192.168.x.x)
- **Transport:** Port 443 (HTTPS)
- **Application:** Web browser

### Practice
```bash
# Trace layers in action
tcpdump -i eth0
wireshark
```

---

**OSI Model = Network Stack Master!** 🌐📡
