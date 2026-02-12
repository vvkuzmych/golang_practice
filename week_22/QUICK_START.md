# Week 22 — Quick Start

## 🎯 Мета тижня
Освоїти Terraform та Kubernetes для production-ready розгортання Go застосунків.

---

## 🛠️ Installation

```bash
# Terraform
brew install terraform

# kubectl
brew install kubectl

# Minikube (local Kubernetes)
brew install minikube

# Start Minikube
minikube start

# Verify
kubectl get nodes
```

---

## 📖 Швидке навчання (90 хв)

```bash
# 1. Terraform Basics
cat theory/01_terraform_basics.md

# 2. Kubernetes Fundamentals
cat theory/03_kubernetes_basics.md
```

---

## ⚡ Terraform Quick Start

### 1. Create main.tf

```hcl
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

resource "aws_instance" "web" {
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t2.micro"
  
  tags = {
    Name = "WebServer"
  }
}

output "instance_ip" {
  value = aws_instance.web.public_ip
}
```

### 2. Run Terraform

```bash
# Initialize
terraform init

# Plan
terraform plan

# Apply
terraform apply

# Output
terraform output instance_ip

# Destroy
terraform destroy
```

---

## ⚡ Kubernetes Quick Start

### 1. Create deployment.yaml

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
      - name: myapp
        image: nginx:latest
        ports:
        - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: myapp-service
spec:
  type: NodePort
  selector:
    app: myapp
  ports:
  - protocol: TCP
    port: 80
    targetPort: 80
    nodePort: 30080
```

### 2. Deploy to Kubernetes

```bash
# Apply
kubectl apply -f deployment.yaml

# Check status
kubectl get pods
kubectl get services

# Port forward
kubectl port-forward service/myapp-service 8080:80

# Open in browser (Minikube)
minikube service myapp-service

# Logs
kubectl logs -f deployment/myapp

# Scale
kubectl scale deployment myapp --replicas=5

# Update image
kubectl set image deployment/myapp myapp=nginx:alpine

# Delete
kubectl delete -f deployment.yaml
```

---

## 🔑 Key Commands

### Terraform

```bash
terraform init          # Initialize
terraform fmt           # Format code
terraform validate      # Validate syntax
terraform plan          # Preview changes
terraform apply         # Apply changes
terraform destroy       # Destroy infrastructure
terraform output        # Show outputs
terraform state list    # List resources
```

### Kubernetes

```bash
kubectl get pods                    # List pods
kubectl get deployments             # List deployments
kubectl get services                # List services
kubectl describe pod <name>         # Details
kubectl logs <pod-name>             # Logs
kubectl exec -it <pod> -- sh        # Shell into pod
kubectl apply -f <file>             # Create/update
kubectl delete -f <file>            # Delete
kubectl scale deployment <name> --replicas=5  # Scale
kubectl port-forward service/<name> 8080:80   # Port forward
```

---

## 💡 Common Patterns

### Terraform: Remote State

```hcl
terraform {
  backend "s3" {
    bucket         = "my-terraform-state"
    key            = "prod/terraform.tfstate"
    region         = "us-east-1"
    dynamodb_table = "terraform-locks"
    encrypt        = true
  }
}
```

### Kubernetes: ConfigMap & Secret

```yaml
# ConfigMap
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  APP_NAME: "MyApp"
  LOG_LEVEL: "info"
---
# Secret
apiVersion: v1
kind: Secret
metadata:
  name: app-secret
type: Opaque
data:
  DB_PASSWORD: cGFzc3dvcmQxMjM=  # base64 encoded
---
# Use in Deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  template:
    spec:
      containers:
      - name: myapp
        image: myapp:1.0
        env:
        - name: APP_NAME
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: APP_NAME
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: app-secret
              key: DB_PASSWORD
```

---

## 🚀 Deploying Go App to Kubernetes

### 1. Dockerize Go App

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

```bash
# Build
docker build -t myapp:1.0 .

# Push to registry
docker tag myapp:1.0 myregistry/myapp:1.0
docker push myregistry/myapp:1.0
```

### 2. Kubernetes Manifest

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
      - name: myapp
        image: myregistry/myapp:1.0
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
        resources:
          requests:
            memory: "64Mi"
            cpu: "250m"
          limits:
            memory: "128Mi"
            cpu: "500m"
---
apiVersion: v1
kind: Service
metadata:
  name: myapp-service
spec:
  type: LoadBalancer
  selector:
    app: myapp
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
```

### 3. Deploy

```bash
kubectl apply -f k8s/deployment.yaml
kubectl get pods
kubectl logs -f deployment/myapp
```

---

## ⚠️ Common Issues

### Terraform: State Lock

```bash
# Force unlock (if stuck)
terraform force-unlock <lock-id>
```

### Kubernetes: Pod Not Starting

```bash
# Describe pod (check events)
kubectl describe pod <pod-name>

# Check logs
kubectl logs <pod-name>

# Check previous logs (if crashed)
kubectl logs <pod-name> --previous

# Common issues:
# - ImagePullBackOff: image not found
# - CrashLoopBackOff: app crashes on start
# - Pending: insufficient resources
```

---

## 📝 Mini Projects

1. **Terraform AWS Infrastructure** - VPC, EC2, RDS
2. **K8s Go API** - Deploy Go REST API with PostgreSQL
3. **Multi-environment** - dev/staging/prod with Terraform workspaces
4. **CI/CD Pipeline** - GitHub Actions → Docker → Kubernetes
5. **Monitoring** - Prometheus + Grafana on Kubernetes

---

## ✅ Перевірка розуміння

### Terraform
- [ ] Розумію resource, provider, variable, output
- [ ] Можу написати basic Terraform config
- [ ] Знаю Terraform workflow (init, plan, apply, destroy)
- [ ] Розумію state management
- [ ] Можу використовувати modules

### Kubernetes
- [ ] Розумію Pod, Deployment, Service
- [ ] Можу створити Deployment з ConfigMap/Secret
- [ ] Знаю kubectl основні команди
- [ ] Розумію різницю між ClusterIP, NodePort, LoadBalancer
- [ ] Можу налаштувати health checks

---

## 🔗 Корисні посилання

- [Terraform AWS Provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [Kubernetes kubectl Cheat Sheet](https://kubernetes.io/docs/reference/kubectl/cheatsheet/)
- [Docker Hub](https://hub.docker.com/)
- [Minikube](https://minikube.sigs.k8s.io/docs/)

---

## 🚀 Наступний крок

**Production-ready practices:**
- Terraform modules for reusability
- Kubernetes Helm charts
- CI/CD pipelines
- Monitoring & Logging (Prometheus, Grafana, ELK)
- Security (RBAC, Network Policies)
