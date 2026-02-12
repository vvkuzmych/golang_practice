# Week 20 — Quick Start

## 🎯 Мета тижня
Підготуватися до System Design інтерв'ю та розуміти масштабування систем.

---

## 📖 Швидке навчання (45 хв)

```bash
# 1. CAP Theorem
cat theory/01_cap_theorem.md

# 2. Scaling Strategies
cat theory/02_scaling.md
```

---

## 💡 Ключові концепції

### CAP Theorem
- **C + P** = Consistency over Availability (MongoDB, HBase)
- **A + P** = Availability over Consistency (Cassandra, DynamoDB)

### Scaling
- **Vertical** = більше CPU/RAM (простіше, але є ліміт)
- **Horizontal** = більше серверів (складніше, але нескінченно)

---

## 🎤 System Design Interview Framework

```
1. Уточнити requirements (5 хв)
   - Functional: що система має робити?
   - Non-functional: users, QPS, latency, storage

2. Back-of-the-envelope calculations (5 хв)
   - Users: 100M DAU
   - Requests: 1000 QPS
   - Storage: 1TB/day

3. High-level design (10 хв)
   - Схема компонентів
   - Client → LB → API → Cache → DB

4. Deep dive (15 хв)
   - Bottleneck аналіз
   - Scaling strategy
   - Database choice

5. Wrap up (5 хв)
   - Trade-offs
   - Improvements
```

---

## 📝 Приклади питань

- Design Twitter
- Design URL Shortener
- Design Instagram
- Design Rate Limiter
- Design Chat System

---

## ✅ Перевірка розуміння

- [ ] Можу пояснити CAP theorem
- [ ] Знаю різницю між vertical/horizontal scaling
- [ ] Розумію load balancing алгоритми
- [ ] Можу обрати caching strategy
- [ ] Розумію коли використовувати sharding

---

## 🚀 Наступний крок

Практикуй System Design інтерв'ю на:
- [Pramp](https://www.pramp.com/)
- [interviewing.io](https://interviewing.io/)
- [Exponent](https://www.tryexponent.com/)
