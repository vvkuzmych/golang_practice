# Week 16 - Quick Start

## 🚀 OSI Model Quick Reference

### 7 Layers

```
7. Application  → HTTP, DNS, SMTP        (Data)
6. Presentation → SSL/TLS                (Data)
5. Session      → NetBIOS                (Data)
4. Transport    → TCP, UDP               (Segment, Port)
3. Network      → IPv4, IPv6             (Packet, IP)
2. Data Link    → Ethernet, WiFi         (Frame, MAC)
1. Physical     → Cable, WiFi radio      (Bit)
```

---

## 📊 Key Concepts

### Addressing

```bash
Port (Layer 4):  443, 80, 22
IP (Layer 3):    192.168.1.1, 2001:0db8::1
MAC (Layer 2):   AA:BB:CC:DD:EE:FF
```

### Devices

```bash
Router (Layer 3):  Routes between networks
Switch (Layer 2):  Forwards frames by MAC
Hub (Layer 1):     Broadcasts to all ports
```

### TCP vs UDP

```bash
TCP:  Reliable, ordered, slow      (HTTP, SSH)
UDP:  Fast, unreliable, no order   (DNS, streaming)
```

---

## 🔧 Quick Commands

### Test Connectivity

```bash
# Layer 7 - Application
curl -v https://google.com
nslookup google.com

# Layer 4 - Transport
telnet google.com 80
nc -zv google.com 80

# Layer 3 - Network
ping google.com
traceroute google.com

# Layer 2 - Data Link
arp -a
ifconfig
```

### Check Ports

```bash
# List open ports
lsof -i -P -n | grep LISTEN
netstat -tuln

# Check specific port
lsof -i :8080
```

---

## 💻 Go Examples

### TCP Client

```go
conn, _ := net.Dial("tcp", "google.com:80")
defer conn.Close()
fmt.Fprintf(conn, "GET / HTTP/1.1\r\nHost: google.com\r\n\r\n")
```

### UDP Client

```go
conn, _ := net.Dial("udp", "8.8.8.8:53")
defer conn.Close()
conn.Write([]byte("query"))
```

### DNS Lookup

```go
ips, _ := net.LookupIP("google.com")
for _, ip := range ips {
    fmt.Println(ip)
}
```

---

## 🎯 Common Ports

```
20/21  → FTP
22     → SSH
25     → SMTP
53     → DNS
80     → HTTP
110    → POP3
143    → IMAP
443    → HTTPS
3306   → MySQL
5432   → PostgreSQL
6379   → Redis
```

---

## 📚 Files

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_16

# Quick references
cat OSI_CHEAT_SHEET.md
cat TCP_VS_UDP.md
cat PORT_NUMBERS.md
```

---

## ✅ Quick Checklist

- [ ] Назвати всі 7 рівнів OSI
- [ ] Знати PDU: bit, frame, packet, segment, data
- [ ] Пояснити TCP vs UDP
- [ ] Знати різницю IP vs MAC vs Port
- [ ] Знати різницю Router vs Switch vs Hub
- [ ] Використати ping, traceroute, telnet

---

**Week 16: OSI Model!** 🌐⚡
