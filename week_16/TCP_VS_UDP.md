# TCP vs UDP - Quick Comparison

## 📊 Overview

**TCP (Transmission Control Protocol)** та **UDP (User Datagram Protocol)** - два основних протоколи Transport Layer (Layer 4).

---

## 🎯 Side-by-Side Comparison

| Feature | TCP | UDP |
|---------|-----|-----|
| **Full Name** | Transmission Control Protocol | User Datagram Protocol |
| **Connection** | Connection-oriented (3-way handshake) | Connectionless |
| **Reliability** | ✅ Guaranteed delivery | ❌ Best effort (no guarantees) |
| **Order** | ✅ Ordered (packets arrive in order) | ❌ Unordered (packets may arrive out of order) |
| **Error Checking** | ✅ Extensive (checksums, retransmission) | ⚠️ Basic checksum only |
| **Speed** | ❌ Slower (more overhead) | ✅ Faster (minimal overhead) |
| **Header Size** | 20-60 bytes | 8 bytes |
| **Flow Control** | ✅ Yes (sliding window) | ❌ No |
| **Congestion Control** | ✅ Yes | ❌ No |
| **Use Cases** | HTTP, SMTP, SSH, FTP, databases | DNS, video streaming, gaming, VoIP |

---

## 🔧 TCP - Transmission Control Protocol

### Characteristics

✅ **Reliable** - guarantees delivery  
✅ **Ordered** - packets arrive in correct order  
✅ **Connection-oriented** - establishes connection before sending  
✅ **Error recovery** - automatic retransmission  
✅ **Flow control** - prevents overwhelming receiver  
✅ **Congestion control** - adjusts to network conditions  

### 3-Way Handshake

```
Client                    Server
   │                         │
   │──────── SYN ────────────>│  1. Client: "Want to connect"
   │                         │
   │<──── SYN-ACK ───────────│  2. Server: "OK, let's connect"
   │                         │
   │──────── ACK ────────────>│  3. Client: "Connection established"
   │                         │
   │═════ DATA TRANSFER ═════│
```

### TCP Header

```
0                   16                  31
+-------------------+-------------------+
|   Source Port     | Destination Port  |
+-------------------+-------------------+
|        Sequence Number                |
+---------------------------------------+
|     Acknowledgment Number             |
+---------------------------------------+
| Offset| Flags     |   Window Size     |
+---------------------------------------+
|   Checksum        |  Urgent Pointer   |
+---------------------------------------+
|            Options (optional)         |
+---------------------------------------+

Size: 20-60 bytes
```

### Use Cases

```
✅ HTTP/HTTPS  - Web browsing (need reliability)
✅ SMTP/POP3   - Email (need guaranteed delivery)
✅ SSH/FTP     - File transfer (need complete data)
✅ Databases   - SQL queries (need accuracy)
✅ APIs        - REST/GraphQL (need reliability)
```

---

## ⚡ UDP - User Datagram Protocol

### Characteristics

✅ **Fast** - minimal overhead  
✅ **Lightweight** - small header  
❌ **Unreliable** - no delivery guarantees  
❌ **Unordered** - packets may arrive out of order  
❌ **Connectionless** - no handshake  
❌ **No flow control** - sender doesn't know receiver state  

### UDP Header

```
0                   16                  31
+-------------------+-------------------+
|   Source Port     | Destination Port  |
+-------------------+-------------------+
|     Length        |     Checksum      |
+-------------------+-------------------+

Size: 8 bytes (fixed)
```

### Use Cases

```
✅ DNS          - Quick lookups (53 UDP)
✅ Video/Audio  - Streaming (ok to lose packets)
✅ Gaming       - Real-time (low latency critical)
✅ VoIP         - Voice calls (speed > reliability)
✅ Broadcasting - One-to-many (DHCP, TFTP)
```

---

## 🎯 When to Use Which?

### Use TCP when:

✅ Data integrity is critical (banking, databases)  
✅ You need guaranteed delivery (file transfer)  
✅ You need ordered packets (HTTP)  
✅ Connection state is important  

### Use UDP when:

✅ Speed is more important than reliability (streaming)  
✅ You can tolerate packet loss (video conferencing)  
✅ Low latency is critical (gaming)  
✅ Broadcasting to multiple recipients (DHCP)  
✅ Small, simple queries (DNS)  

---

## 💻 Go Examples

### TCP Server

```go
package main

import (
	"fmt"
	"net"
)

func main() {
	// Listen on TCP port 8080
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	fmt.Println("TCP server listening on :8080")

	for {
		// Accept connection (blocking)
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		// Handle connection
		go func(c net.Conn) {
			defer c.Close()
			buf := make([]byte, 1024)
			n, _ := c.Read(buf)
			fmt.Printf("Received: %s\n", buf[:n])
			c.Write([]byte("ACK\n"))
		}(conn)
	}
}
```

### TCP Client

```go
package main

import (
	"fmt"
	"net"
)

func main() {
	// Connect to TCP server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Send data
	conn.Write([]byte("Hello TCP\n"))

	// Receive response
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	fmt.Printf("Response: %s\n", buf[:n])
}
```

---

### UDP Server

```go
package main

import (
	"fmt"
	"net"
)

func main() {
	// Listen on UDP port 8080
	conn, err := net.ListenPacket("udp", ":8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("UDP server listening on :8080")

	buf := make([]byte, 1024)
	for {
		// Read datagram (blocking)
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			continue
		}

		fmt.Printf("Received from %s: %s\n", addr, buf[:n])

		// Send response
		conn.WriteTo([]byte("ACK\n"), addr)
	}
}
```

### UDP Client

```go
package main

import (
	"fmt"
	"net"
)

func main() {
	// Resolve UDP address
	addr, _ := net.ResolveUDPAddr("udp", "localhost:8080")

	// Create UDP connection
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Send datagram
	conn.Write([]byte("Hello UDP\n"))

	// Receive response
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	fmt.Printf("Response: %s\n", buf[:n])
}
```

---

## 📊 Performance Comparison

### TCP

```
Connection:     SYN → SYN-ACK → ACK (3 packets)
Send 1KB:       Data → ACK (2+ packets)
Close:          FIN → ACK (2 packets)
Total:          7+ packets for 1KB

Latency:        Higher (handshake + ACKs)
Throughput:     Good (with flow control)
Overhead:       20-60 bytes per packet
```

### UDP

```
Connection:     None (0 packets)
Send 1KB:       Data (1 packet, no ACK)
Close:          None (0 packets)
Total:          1 packet for 1KB

Latency:        Lower (no handshake)
Throughput:     Excellent (no ACKs)
Overhead:       8 bytes per packet
```

---

## 🎯 Real-World Examples

### DNS (UDP + TCP)

```
Primary: UDP port 53 (fast queries)
Fallback: TCP port 53 (large responses > 512 bytes)

Why UDP?
- Quick, simple queries
- Small responses
- Retry is acceptable (client resends if timeout)
```

### HTTP/HTTPS (TCP)

```
Always: TCP port 80/443

Why TCP?
- Need complete HTML/CSS/JS files
- Can't tolerate data loss
- Order matters (HTML before rendering)
```

### Video Streaming (UDP)

```
Typically: UDP or UDP-based (like QUIC)

Why UDP?
- Speed > reliability
- Missing frames = slight glitch (acceptable)
- Can't wait for retransmission (video is real-time)
```

### Online Gaming (UDP)

```
Game state updates: UDP
Chat messages: TCP

Why UDP for game state?
- Low latency critical (player position, shots)
- Old data is useless (outdated position)
- 60 updates/sec - can skip some frames

Why TCP for chat?
- Messages must arrive intact
- Order matters
- Not time-sensitive
```

---

## ✅ Summary

### TCP = Reliable Car Delivery

```
📦 → Truck picks up package
📦 → Drives to destination (planned route)
📦 → Delivers package (confirms receipt)
📦 → Returns with confirmation

Slow but guaranteed!
```

### UDP = Throwing Packages

```
📦 → Throw package over fence
📦 → Hope it lands safely
📦 → No confirmation
📦 → Maybe lost, maybe broken, but FAST!

Fast but risky!
```

---

**Week 16: TCP vs UDP!** 🔄⚡
