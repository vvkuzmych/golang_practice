# Week 22 — Terraform & Kubernetes

**Ціль:** Освоїти Infrastructure as Code (Terraform) та оркестрацію контейнерів (Kubernetes) для production-ready розгортання Go застосунків.

---

## 📚 Теорія

### [01. Terraform Basics](./theory/01_terraform_basics.md)
- Що таке Infrastructure as Code (IaC)
- Terraform providers
- Resources, Data sources
- Variables, Outputs
- State management
- Terraform workflow (init, plan, apply, destroy)

### [02. Terraform for AWS](./theory/02_terraform_aws.md)
- EC2, VPC, Security Groups
- RDS, S3, Lambda
- IAM roles and policies
- Terraform modules
- Remote state (S3 + DynamoDB)

### [03. Kubernetes Fundamentals](./theory/03_kubernetes_basics.md)
- Pods, Deployments, Services
- ConfigMaps, Secrets
- Namespaces
- Labels and Selectors
- kubectl commands

### [04. Kubernetes for Go Apps](./theory/04_kubernetes_go_apps.md)
- Containerizing Go applications
- Kubernetes manifests
- Health checks (liveness, readiness)
- Resource limits
- Horizontal Pod Autoscaler (HPA)
- Ingress

---

## 🛠️ Практика

### [01. Deploy Go API with Terraform](./practice/01_terraform_go_api/)
- EC2 instance
- Security groups
- Deploy Go binary
- Terraform state

### [02. Kubernetes Local Setup](./practice/02_k8s_local/)
- Minikube or Kind
- Deploy Go app to local k8s
- Service exposure
- ConfigMaps for configuration

### [03. Full Stack Deployment](./practice/03_full_stack/)
- Go API + PostgreSQL на Kubernetes
- Persistent volumes
- Secrets management
- Ingress configuration

### [04. CI/CD Pipeline](./practice/04_cicd/)
- GitHub Actions
- Build Docker image
- Push to registry
- Deploy to Kubernetes
- Rolling updates

---

## 📝 Exercises

### [Exercise 1: Terraform AWS Infrastructure](./exercises/exercise_1.md)
Створити повну AWS інфраструктуру: VPC, EC2, RDS, S3.

### [Exercise 2: Deploy to Kubernetes](./exercises/exercise_2.md)
Розгорнути Go microservices на Kubernetes з service mesh.

### [Exercise 3: Auto-scaling](./exercises/exercise_3.md)
Налаштувати HPA та Cluster Autoscaler.

---

## 🎯 Learning Outcomes

Після цього тижня ви зможете:
- ✅ Писати Terraform код для AWS/GCP/Azure
- ✅ Керувати infrastructure as code
- ✅ Розгортати Go застосунки на Kubernetes
- ✅ Налаштовувати auto-scaling
- ✅ Використовувати ConfigMaps та Secrets
- ✅ Створювати CI/CD pipelines
- ✅ Моніторити та дебажити k8s pods
- ✅ Розуміти різницю між Deployment, StatefulSet, DaemonSet

---

## 🔧 Tools to Install

```bash
# Terraform
brew install terraform

# kubectl
brew install kubectl

# Minikube (local k8s)
brew install minikube

# Docker
brew install docker

# Helm (optional)
brew install helm

# k9s (optional, terminal UI for k8s)
brew install k9s
```

---

## 📖 Key Concepts

### Terraform
- **Resource** - інфраструктурний об'єкт (EC2, VPC, etc.)
- **Provider** - інтеграція з cloud provider (AWS, GCP, Azure)
- **State** - поточний стан інфраструктури
- **Module** - reusable Terraform код
- **Backend** - де зберігається state (S3, Terraform Cloud)

### Kubernetes
- **Pod** - найменша одиниця, один або більше контейнерів
- **Deployment** - управління replicas
- **Service** - мережевий доступ до pods
- **ConfigMap** - конфігурація
- **Secret** - чутливі дані
- **Ingress** - HTTP(S) routing

---

## 📖 Additional Resources

### Terraform
- [Terraform Documentation](https://developer.hashicorp.com/terraform/docs)
- [Terraform AWS Examples](https://github.com/terraform-aws-modules)
- [Learn Terraform](https://learn.hashicorp.com/terraform)

### Kubernetes
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Kubernetes Patterns](https://k8spatterns.io/)
- [kubectl Cheat Sheet](https://kubernetes.io/docs/reference/kubectl/cheatsheet/)
- [Kubernetes the Hard Way](https://github.com/kelseyhightower/kubernetes-the-hard-way)

---

**Previous:** [Week 21 — os Package](../week_21/README.md)
