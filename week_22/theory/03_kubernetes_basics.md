# Kubernetes Fundamentals

## Що таке Kubernetes?

**Kubernetes (K8s)** - система оркестрації контейнерів для автоматизації розгортання, масштабування та управління контейнеризованими застосунками.

### Чому Kubernetes?

✅ **Auto-scaling** - автоматичне масштабування  
✅ **Self-healing** - автоматичне відновлення  
✅ **Load balancing** - розподіл навантаження  
✅ **Rolling updates** - оновлення без downtime  
✅ **Secret management** - безпечне зберігання секретів  
✅ **Multi-cloud** - працює на AWS, GCP, Azure, bare metal  

---

## Архітектура Kubernetes

```
┌─────────────────────────────────────────┐
│          Control Plane                   │
│  (API Server, Scheduler, Controllers)    │
└─────────────────────────────────────────┘
                  │
    ┌─────────────┼─────────────┐
    │             │             │
┌───▼────┐   ┌───▼────┐   ┌───▼────┐
│ Node 1 │   │ Node 2 │   │ Node 3 │
│        │   │        │   │        │
│  Pods  │   │  Pods  │   │  Pods  │
└────────┘   └────────┘   └────────┘
```

---

## Основні об'єкти

### 1. Pod
**Найменша одиниця. Один або більше контейнерів.**

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: myapp-pod
  labels:
    app: myapp
spec:
  containers:
  - name: myapp-container
    image: myapp:1.0
    ports:
    - containerPort: 8080
    env:
    - name: ENV
      value: "production"
```

**Створення:**

```bash
kubectl apply -f pod.yaml

# Альтернативно (imperative)
kubectl run myapp --image=myapp:1.0 --port=8080
```

### 2. Deployment
**Управління replicas, rolling updates.**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp-deployment
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
        image: myapp:1.0
        ports:
        - containerPort: 8080
        resources:
          requests:
            memory: "64Mi"
            cpu: "250m"
          limits:
            memory: "128Mi"
            cpu: "500m"
```

**Команди:**

```bash
# Створити deployment
kubectl apply -f deployment.yaml

# Масштабувати
kubectl scale deployment myapp-deployment --replicas=5

# Rolling update
kubectl set image deployment/myapp-deployment myapp=myapp:2.0

# Rollback
kubectl rollout undo deployment/myapp-deployment

# Status
kubectl rollout status deployment/myapp-deployment

# History
kubectl rollout history deployment/myapp-deployment
```

### 3. Service
**Мережевий доступ до Pods.**

#### ClusterIP (default)
**Доступ тільки всередині кластера.**

```yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp-service
spec:
  type: ClusterIP
  selector:
    app: myapp
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
```

#### NodePort
**Доступ ззовні через port на node.**

```yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp-nodeport
spec:
  type: NodePort
  selector:
    app: myapp
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
    nodePort: 30080  # 30000-32767
```

#### LoadBalancer
**Cloud load balancer (AWS ELB, GCP LB).**

```yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp-loadbalancer
spec:
  type: LoadBalancer
  selector:
    app: myapp
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
```

**Команди:**

```bash
# Створити service
kubectl apply -f service.yaml

# Port forward (для тестування)
kubectl port-forward service/myapp-service 8080:80

# Відкрити в браузері (Minikube)
minikube service myapp-service
```

### 4. ConfigMap
**Конфігурація застосунку.**

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: myapp-config
data:
  APP_NAME: "MyApp"
  LOG_LEVEL: "info"
  DATABASE_HOST: "postgres.default.svc.cluster.local"
  config.yaml: |
    server:
      port: 8080
      timeout: 30s
```

**Використання в Pod:**

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: myapp-pod
spec:
  containers:
  - name: myapp
    image: myapp:1.0
    env:
    # Single env var
    - name: APP_NAME
      valueFrom:
        configMapKeyRef:
          name: myapp-config
          key: APP_NAME
    # All keys as env vars
    envFrom:
    - configMapRef:
        name: myapp-config
    # Mount as file
    volumeMounts:
    - name: config-volume
      mountPath: /etc/config
  volumes:
  - name: config-volume
    configMap:
      name: myapp-config
```

**Команди:**

```bash
# Створити з файлу
kubectl create configmap myapp-config --from-file=config.yaml

# Створити з literal values
kubectl create configmap myapp-config --from-literal=APP_NAME=MyApp

# Переглянути
kubectl get configmap myapp-config -o yaml
```

### 5. Secret
**Чутливі дані (паролі, токени, ключі).**

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: myapp-secret
type: Opaque
data:
  # Base64 encoded!
  DB_PASSWORD: cGFzc3dvcmQxMjM=
  API_KEY: c2VjcmV0a2V5MTIz
```

**Створення:**

```bash
# З файлу
kubectl create secret generic myapp-secret --from-file=./secret.txt

# З literal values
kubectl create secret generic myapp-secret \
  --from-literal=DB_PASSWORD=password123 \
  --from-literal=API_KEY=secretkey123

# Base64 encode
echo -n 'password123' | base64
# cGFzc3dvcmQxMjM=
```

**Використання:**

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: myapp-pod
spec:
  containers:
  - name: myapp
    image: myapp:1.0
    env:
    - name: DB_PASSWORD
      valueFrom:
        secretKeyRef:
          name: myapp-secret
          key: DB_PASSWORD
```

---

## Namespace
**Логічне розділення кластера.**

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: development
```

**Команди:**

```bash
# Створити namespace
kubectl create namespace development

# Список namespaces
kubectl get namespaces

# Встановити default namespace
kubectl config set-context --current --namespace=development

# Переглянути pods в namespace
kubectl get pods -n development

# Всі pods в всіх namespaces
kubectl get pods --all-namespaces
kubectl get pods -A  # short version
```

---

## kubectl Commands

### Basics

```bash
# Cluster info
kubectl cluster-info
kubectl version

# Nodes
kubectl get nodes
kubectl describe node <node-name>

# Pods
kubectl get pods
kubectl get pods -o wide  # more info
kubectl describe pod <pod-name>
kubectl logs <pod-name>
kubectl logs <pod-name> -f  # follow
kubectl logs <pod-name> -c <container-name>  # specific container

# Exec into pod
kubectl exec -it <pod-name> -- /bin/bash
kubectl exec -it <pod-name> -- /bin/sh
```

### CRUD Operations

```bash
# Create
kubectl apply -f manifest.yaml
kubectl create deployment myapp --image=myapp:1.0

# Read
kubectl get pods
kubectl get deployments
kubectl get services
kubectl get all  # all resources

# Update
kubectl apply -f manifest.yaml
kubectl edit deployment myapp
kubectl set image deployment/myapp myapp=myapp:2.0

# Delete
kubectl delete -f manifest.yaml
kubectl delete pod <pod-name>
kubectl delete deployment <deployment-name>
kubectl delete all --all  # delete all resources
```

### Debugging

```bash
# Logs
kubectl logs <pod-name>
kubectl logs <pod-name> --previous  # previous container logs

# Describe (detailed info + events)
kubectl describe pod <pod-name>

# Port forward
kubectl port-forward pod/<pod-name> 8080:8080
kubectl port-forward service/<service-name> 8080:80

# Exec
kubectl exec -it <pod-name> -- sh

# Copy files
kubectl cp <pod-name>:/path/to/file ./local-file
kubectl cp ./local-file <pod-name>:/path/to/file

# Top (resource usage)
kubectl top nodes
kubectl top pods

# Events
kubectl get events
kubectl get events --watch
```

---

## Labels & Selectors

### Labels

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: myapp-pod
  labels:
    app: myapp
    environment: production
    version: "1.0"
    tier: frontend
```

### Selectors

```bash
# Get pods by label
kubectl get pods -l app=myapp
kubectl get pods -l environment=production
kubectl get pods -l app=myapp,environment=production

# Not equal
kubectl get pods -l environment!=production

# In set
kubectl get pods -l 'environment in (production,staging)'

# Exists
kubectl get pods -l app

# Multiple labels in YAML
selector:
  matchLabels:
    app: myapp
  matchExpressions:
  - key: environment
    operator: In
    values:
    - production
    - staging
```

---

## Практичний приклад: Full Stack App

### 1. Namespace

```yaml
# namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: myapp
```

### 2. ConfigMap

```yaml
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: myapp-config
  namespace: myapp
data:
  DATABASE_HOST: "postgres-service"
  DATABASE_PORT: "5432"
  DATABASE_NAME: "myapp"
```

### 3. Secret

```yaml
# secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: myapp-secret
  namespace: myapp
type: Opaque
data:
  DATABASE_PASSWORD: cGFzc3dvcmQxMjM=
```

### 4. PostgreSQL Deployment

```yaml
# postgres-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
  namespace: myapp
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:15
        env:
        - name: POSTGRES_DB
          valueFrom:
            configMapKeyRef:
              name: myapp-config
              key: DATABASE_NAME
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: myapp-secret
              key: DATABASE_PASSWORD
        ports:
        - containerPort: 5432
```

### 5. PostgreSQL Service

```yaml
# postgres-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: postgres-service
  namespace: myapp
spec:
  selector:
    app: postgres
  ports:
  - protocol: TCP
    port: 5432
    targetPort: 5432
```

### 6. Go App Deployment

```yaml
# app-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
  namespace: myapp
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
        image: myapp:1.0
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
        - name: DATABASE_HOST
          valueFrom:
            configMapKeyRef:
              name: myapp-config
              key: DATABASE_HOST
        - name: DATABASE_PASSWORD
          valueFrom:
            secretKeyRef:
              name: myapp-secret
              key: DATABASE_PASSWORD
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

### 7. Go App Service

```yaml
# app-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp-service
  namespace: myapp
spec:
  type: LoadBalancer
  selector:
    app: myapp
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
```

### Розгортання

```bash
# Apply all
kubectl apply -f namespace.yaml
kubectl apply -f configmap.yaml
kubectl apply -f secret.yaml
kubectl apply -f postgres-deployment.yaml
kubectl apply -f postgres-service.yaml
kubectl apply -f app-deployment.yaml
kubectl apply -f app-service.yaml

# Альтернативно (apply всю директорію)
kubectl apply -f ./k8s/

# Перевірити
kubectl get all -n myapp

# Logs
kubectl logs -f deployment/myapp -n myapp

# Port forward
kubectl port-forward service/myapp-service 8080:80 -n myapp
```

---

## Підсумок

| Об'єкт | Опис |
|--------|------|
| **Pod** | Один або більше контейнерів |
| **Deployment** | Управління replicas, rolling updates |
| **Service** | Мережевий доступ до Pods |
| **ConfigMap** | Конфігурація |
| **Secret** | Чутливі дані |
| **Namespace** | Логічне розділення |

| Команда | Опис |
|---------|------|
| `kubectl get` | Список resources |
| `kubectl describe` | Детальна інформація |
| `kubectl logs` | Логи контейнера |
| `kubectl exec` | Виконати команду в контейнері |
| `kubectl apply` | Створити/оновити resource |
| `kubectl delete` | Видалити resource |

**Переваги Kubernetes:**
- ✅ Auto-scaling
- ✅ Self-healing
- ✅ Rolling updates
- ✅ Multi-cloud

**Недоліки:**
- ❌ Складність
- ❌ Крива навчання
- ❌ Overhead для малих проектів
