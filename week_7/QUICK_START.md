# Week 7: Швидкий Старт

## 🚀 За 10 хвилин

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_7

# Прочитайте README
cat README.md

# Почніть з best practices
cat theory/01_go_best_practices.md

# Запустіть приклад
go run practice/01_advanced_api/main.go
```

---

## 📖 Рекомендований порядок

### День 1-2: Advanced Go (6-8 годин)
```bash
cat theory/01_go_best_practices.md
cat theory/02_advanced_concurrency.md
go run practice/01_advanced_api/main.go
```

### День 3-4: APIs & Cloud (8-10 годин)
```bash
cat theory/03_restful_apis.md
cat theory/04_aws_cloud.md
cat theory/05_scalable_backend.md
go run practice/02_aws_integration/main.go
```

### День 5-6: Performance & Testing (6-8 годин)
```bash
cat theory/06_debugging_performance.md
cat theory/07_testing.md
go test practice/04_testing/...
```

### День 7: DevOps (8-10 годин)
```bash
cat theory/08_cicd_docker_k8s.md
cat theory/09_technical_english.md
cat theory/10_security_compliance.md

# Docker
cd practice/05_docker
docker build -t myapp .
docker run -p 8080:8080 myapp

# Kubernetes
kubectl apply -f practice/06_k8s/
```

---

## ✅ Швидка перевірка знань

```bash
# 1. Go best practices
go fmt ./...
go vet ./...
golangci-lint run

# 2. Testing
go test -v -cover ./...

# 3. Profiling
go test -cpuprofile=cpu.prof
go tool pprof cpu.prof

# 4. Docker
docker build -t test .
docker run test

# 5. Kubernetes
kubectl get pods
kubectl logs <pod-name>
```

---

## 🎯 Мінімальні вимоги

Перед Week 7 ви маєте знати:
- ✅ Week 1-6 completed
- ✅ Go basics
- ✅ HTTP servers
- ✅ Goroutines & channels
- ✅ Basic SQL

---

## 💡 Підказки

1. **AWS Free Tier** - використовуйте для практики
2. **Docker Desktop** - встановіть локально
3. **minikube** - для локального Kubernetes
4. **Postman/Insomnia** - для тестування APIs
5. **GitHub Actions** - безкоштовний CI/CD

---

**Успіхів!** 🚀
