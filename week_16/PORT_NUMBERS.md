# Port Numbers Reference

## 📊 Port Ranges

| Range | Type | Description |
|-------|------|-------------|
| **0-1023** | Well-Known Ports | Reserved for common services (requires root/admin) |
| **1024-49151** | Registered Ports | Assigned by IANA for specific applications |
| **49152-65535** | Dynamic/Private Ports | Client-side ephemeral ports |

**Total ports:** 65,536 (0-65535) per IP address

---

## 🌐 Well-Known Ports (0-1023)

### File Transfer

| Port | Protocol | Service | Description |
|------|----------|---------|-------------|
| 20 | FTP Data | File Transfer Protocol | Data transfer |
| 21 | FTP Control | File Transfer Protocol | Control commands |
| 22 | SSH/SFTP | Secure Shell | Encrypted remote access & file transfer |
| 69 | TFTP | Trivial FTP | Simple file transfer (UDP) |

### Remote Access

| Port | Protocol | Service | Description |
|------|----------|---------|-------------|
| 22 | SSH | Secure Shell | Encrypted terminal |
| 23 | Telnet | Telnet | Unencrypted terminal (⚠️ insecure) |
| 3389 | RDP | Remote Desktop | Windows remote desktop |
| 5900 | VNC | Virtual Network Computing | Remote desktop |

### Email

| Port | Protocol | Service | Description |
|------|----------|---------|-------------|
| 25 | SMTP | Simple Mail Transfer Protocol | Send email |
| 110 | POP3 | Post Office Protocol v3 | Receive email (downloads) |
| 143 | IMAP | Internet Message Access Protocol | Receive email (keeps on server) |
| 465 | SMTPS | SMTP over SSL | Secure email sending |
| 587 | SMTP | SMTP Submission | Modern email sending |
| 993 | IMAPS | IMAP over SSL | Secure IMAP |
| 995 | POP3S | POP3 over SSL | Secure POP3 |

### Web

| Port | Protocol | Service | Description |
|------|----------|---------|-------------|
| 80 | HTTP | Hypertext Transfer Protocol | Web traffic |
| 443 | HTTPS | HTTP over SSL/TLS | Secure web traffic |
| 8080 | HTTP Alt | HTTP Alternative | Development/proxy |
| 8443 | HTTPS Alt | HTTPS Alternative | Alternative secure web |

### DNS & Network

| Port | Protocol | Service | Description |
|------|----------|---------|-------------|
| 53 | DNS | Domain Name System | Name resolution (UDP/TCP) |
| 67 | DHCP | DHCP Server | IP address assignment (UDP) |
| 68 | DHCP | DHCP Client | IP address assignment (UDP) |
| 123 | NTP | Network Time Protocol | Time synchronization (UDP) |
| 161 | SNMP | Simple Network Management | Network monitoring (UDP) |
| 162 | SNMP | SNMP Trap | Network alerts (UDP) |

### Other Services

| Port | Protocol | Service | Description |
|------|----------|---------|-------------|
| 389 | LDAP | Lightweight Directory Access | Directory services |
| 636 | LDAPS | LDAP over SSL | Secure directory |
| 445 | SMB | Server Message Block | Windows file sharing |
| 514 | Syslog | System Logging | Log collection (UDP) |

---

## 🗄️ Database Ports (Registered)

| Port | Service | Database |
|------|---------|----------|
| 1433 | MSSQL | Microsoft SQL Server |
| 1521 | Oracle | Oracle Database |
| 3306 | MySQL | MySQL/MariaDB |
| 5432 | PostgreSQL | PostgreSQL |
| 6379 | Redis | Redis Cache |
| 9042 | Cassandra | Apache Cassandra |
| 27017 | MongoDB | MongoDB |
| 28015 | RethinkDB | RethinkDB |

---

## 📡 Application Ports (Registered)

### Message Queues

| Port | Service | Description |
|------|---------|-------------|
| 5672 | AMQP | RabbitMQ |
| 9092 | Kafka | Apache Kafka |
| 61613 | STOMP | Message protocol |

### Application Servers

| Port | Service | Description |
|------|---------|-------------|
| 8000 | Django | Django dev server |
| 3000 | React/Node | React dev server |
| 4200 | Angular | Angular dev server |
| 5000 | Flask | Flask dev server |
| 8080 | Tomcat | Apache Tomcat |
| 9000 | PHP-FPM | PHP FastCGI |

### Monitoring & Tools

| Port | Service | Description |
|------|---------|-------------|
| 2379 | etcd | Distributed key-value store |
| 8086 | InfluxDB | Time-series database |
| 9090 | Prometheus | Metrics collection |
| 9200 | Elasticsearch | Search engine |
| 9300 | Elasticsearch | Node communication |
| 3000 | Grafana | Monitoring dashboard |

---

## 🎯 Common Port Usage Patterns

### Web Application Stack

```
443 (HTTPS)       → Nginx/Apache
    ↓
8080 (HTTP)       → Application server (Go, Node, etc.)
    ↓
5432 (PostgreSQL) → Database
    ↓
6379 (Redis)      → Cache
```

### Microservices

```
API Gateway:     8080
Auth Service:    8081
User Service:    8082
Order Service:   8083
Payment Service: 8084
```

---

## 🔧 Checking Ports

### List Open Ports

```bash
# Linux/Mac
lsof -i -P -n | grep LISTEN
netstat -tuln
ss -tuln

# Check specific port
lsof -i :8080
netstat -an | grep 8080
```

### Test Port Connection

```bash
# Telnet (Layer 4 test)
telnet google.com 80

# Netcat
nc -zv google.com 80

# Nmap (port scanner)
nmap google.com
```

### Go Example

```go
package main

import (
	"fmt"
	"net"
	"time"
)

func isPortOpen(host string, port string) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func main() {
	if isPortOpen("google.com", "80") {
		fmt.Println("Port 80 is open")
	} else {
		fmt.Println("Port 80 is closed")
	}
}
```

---

## 🔐 Security Considerations

### Dangerous Open Ports

```
⚠️  21 (FTP)      - Unencrypted file transfer
⚠️  23 (Telnet)   - Unencrypted terminal
⚠️  3389 (RDP)    - Often targeted by attacks
⚠️  445 (SMB)     - WannaCry ransomware vector
⚠️  1433 (MSSQL)  - Should not be exposed publicly
⚠️  3306 (MySQL)  - Should not be exposed publicly
```

### Best Practices

✅ **Close unused ports** - Firewall everything except necessary  
✅ **Use non-standard ports** - Don't use defaults for public services  
✅ **Enable SSL/TLS** - 443 instead of 80, 993 instead of 143  
✅ **Restrict by IP** - Whitelist known IPs  
✅ **Use VPN** - For internal services  

---

## 📊 Port Assignment Examples

### Development Environment

```
Frontend:       3000
Backend API:    8080
Database:       5432
Redis:          6379
```

### Production Environment

```
Load Balancer:  443 (public)
    ↓
App Servers:    8080, 8081, 8082 (internal)
    ↓
Database:       5432 (internal, VPN only)
Cache:          6379 (internal, VPN only)
```

---

## 🎓 Quick Quiz

1. **Q:** What port does HTTPS use?  
   **A:** 443

2. **Q:** What port does SSH use?  
   **A:** 22

3. **Q:** What port does MySQL use?  
   **A:** 3306

4. **Q:** What port does DNS use?  
   **A:** 53 (UDP/TCP)

5. **Q:** What port range is for ephemeral/client ports?  
   **A:** 49152-65535

6. **Q:** Should you expose database ports publicly?  
   **A:** No! Use VPN or internal network only

---

## 📚 IANA Reference

Full list: https://www.iana.org/assignments/service-names-port-numbers/

---

**Week 16: Port Numbers!** 🔢🌐
