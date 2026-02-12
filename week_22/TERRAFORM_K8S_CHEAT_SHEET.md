# Terraform & Kubernetes — Cheat Sheet

## 🏗️ Terraform Commands

```bash
# Initialization
terraform init                      # Initialize working directory
terraform init -upgrade             # Upgrade providers

# Formatting & Validation
terraform fmt                       # Format code
terraform fmt -recursive            # Format all files
terraform validate                  # Validate configuration

# Planning
terraform plan                      # Show execution plan
terraform plan -out=tfplan          # Save plan to file
terraform plan -destroy             # Show destroy plan

# Applying
terraform apply                     # Apply changes
terraform apply tfplan              # Apply saved plan
terraform apply -auto-approve       # Skip confirmation
terraform apply -target=resource    # Apply specific resource

# Destroying
terraform destroy                   # Destroy all resources
terraform destroy -auto-approve     # Skip confirmation
terraform destroy -target=resource  # Destroy specific resource

# State Management
terraform state list                # List resources in state
terraform state show resource       # Show resource details
terraform state rm resource         # Remove from state
terraform state mv source dest      # Move/rename resource
terraform refresh                   # Refresh state
terraform force-unlock <id>         # Unlock state

# Outputs
terraform output                    # Show all outputs
terraform output <name>             # Show specific output

# Import
terraform import resource id        # Import existing resource

# Workspaces
terraform workspace list            # List workspaces
terraform workspace new <name>      # Create workspace
terraform workspace select <name>   # Switch workspace
terraform workspace delete <name>   # Delete workspace

# Other
terraform graph                     # Generate dependency graph
terraform show                      # Show current state
terraform version                   # Show version
```

---

## ☸️ Kubernetes (kubectl) Commands

### Cluster Info

```bash
kubectl cluster-info                # Cluster information
kubectl version                     # Client and server version
kubectl get nodes                   # List nodes
kubectl describe node <name>        # Node details
```

### Resources (CRUD)

```bash
# Get
kubectl get pods                    # List pods
kubectl get deployments             # List deployments
kubectl get services                # List services
kubectl get all                     # List all resources
kubectl get pods -o wide            # Detailed info
kubectl get pods -o yaml            # YAML output
kubectl get pods -o json            # JSON output
kubectl get pods -n <namespace>     # Specific namespace
kubectl get pods -A                 # All namespaces
kubectl get pods -w                 # Watch mode

# Describe
kubectl describe pod <name>         # Pod details + events
kubectl describe deployment <name>  # Deployment details
kubectl describe service <name>     # Service details

# Create/Apply
kubectl apply -f file.yaml          # Create/update from file
kubectl apply -f directory/         # Apply all files in directory
kubectl create -f file.yaml         # Create from file (fails if exists)

# Edit
kubectl edit pod <name>             # Edit resource
kubectl edit deployment <name>      # Edit deployment

# Delete
kubectl delete pod <name>           # Delete pod
kubectl delete -f file.yaml         # Delete from file
kubectl delete deployment <name>    # Delete deployment
kubectl delete all --all            # Delete all resources
```

### Pods

```bash
# Logs
kubectl logs <pod>                  # Show logs
kubectl logs <pod> -f               # Follow logs
kubectl logs <pod> -c <container>   # Specific container
kubectl logs <pod> --previous       # Previous container logs
kubectl logs <pod> --tail=100       # Last 100 lines

# Exec
kubectl exec <pod> -- <command>     # Execute command
kubectl exec -it <pod> -- sh        # Interactive shell
kubectl exec -it <pod> -- bash      # Interactive bash

# Port Forward
kubectl port-forward <pod> 8080:80  # Forward pod port
kubectl port-forward svc/<name> 8080:80  # Forward service port

# Copy Files
kubectl cp <pod>:/path /local/path  # Copy from pod
kubectl cp /local/path <pod>:/path  # Copy to pod
```

### Deployments

```bash
# Create
kubectl create deployment <name> --image=<image>  # Create deployment

# Scale
kubectl scale deployment <name> --replicas=5      # Scale to 5 replicas
kubectl autoscale deployment <name> --min=2 --max=10 --cpu-percent=80  # Autoscale

# Update
kubectl set image deployment/<name> <container>=<image>:<tag>  # Update image
kubectl rollout restart deployment/<name>         # Restart deployment

# Rollout
kubectl rollout status deployment/<name>          # Rollout status
kubectl rollout history deployment/<name>         # Rollout history
kubectl rollout undo deployment/<name>            # Rollback to previous
kubectl rollout undo deployment/<name> --to-revision=2  # Rollback to specific
```

### Services

```bash
# Expose
kubectl expose deployment <name> --port=80 --target-port=8080  # Create service
kubectl expose pod <name> --type=NodePort          # Expose as NodePort
kubectl expose deployment <name> --type=LoadBalancer  # Expose as LoadBalancer
```

### ConfigMaps & Secrets

```bash
# ConfigMap
kubectl create configmap <name> --from-file=config.yaml
kubectl create configmap <name> --from-literal=KEY=VALUE
kubectl get configmap <name> -o yaml

# Secret
kubectl create secret generic <name> --from-file=secret.txt
kubectl create secret generic <name> --from-literal=PASSWORD=pass123
kubectl get secret <name> -o yaml

# Base64 encode/decode
echo -n 'password' | base64         # Encode
echo 'cGFzc3dvcmQ=' | base64 -d     # Decode
```

### Namespaces

```bash
kubectl create namespace <name>     # Create namespace
kubectl get namespaces              # List namespaces
kubectl config set-context --current --namespace=<name>  # Set default namespace
kubectl delete namespace <name>     # Delete namespace
```

### Labels & Selectors

```bash
# Get by label
kubectl get pods -l app=myapp
kubectl get pods -l 'env in (prod,staging)'
kubectl get pods -l app,env

# Label
kubectl label pod <name> env=prod   # Add label
kubectl label pod <name> env-       # Remove label
```

### Context & Config

```bash
kubectl config view                 # View kubeconfig
kubectl config current-context      # Current context
kubectl config get-contexts         # List contexts
kubectl config use-context <name>   # Switch context
```

### Debugging

```bash
# Top (resource usage)
kubectl top nodes                   # Node resource usage
kubectl top pods                    # Pod resource usage
kubectl top pods -A                 # All pods

# Events
kubectl get events                  # List events
kubectl get events -w               # Watch events
kubectl get events --sort-by='.metadata.creationTimestamp'

# Drain & Cordon
kubectl cordon <node>               # Mark node unschedulable
kubectl uncordon <node>             # Mark node schedulable
kubectl drain <node>                # Drain node
```

---

## 📋 Terraform Syntax

### Resource

```hcl
resource "aws_instance" "web" {
  ami           = "ami-123456"
  instance_type = "t2.micro"
  
  tags = {
    Name = "WebServer"
  }
}
```

### Variable

```hcl
variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t2.micro"
}

# Use: var.instance_type
```

### Output

```hcl
output "instance_ip" {
  description = "Public IP"
  value       = aws_instance.web.public_ip
}
```

### Data Source

```hcl
data "aws_ami" "ubuntu" {
  most_recent = true
  
  filter {
    name   = "name"
    values = ["ubuntu-*"]
  }
}

# Use: data.aws_ami.ubuntu.id
```

### Module

```hcl
module "vpc" {
  source = "./modules/vpc"
  
  cidr_block = "10.0.0.0/16"
  name       = "my-vpc"
}

# Use: module.vpc.vpc_id
```

### Backend

```hcl
terraform {
  backend "s3" {
    bucket = "terraform-state"
    key    = "prod/terraform.tfstate"
    region = "us-east-1"
  }
}
```

---

## 📋 Kubernetes YAML Syntax

### Pod

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: myapp
  labels:
    app: myapp
spec:
  containers:
  - name: myapp
    image: myapp:1.0
    ports:
    - containerPort: 8080
    env:
    - name: ENV
      value: "prod"
```

### Deployment

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

### Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp-service
spec:
  type: LoadBalancer  # ClusterIP, NodePort, LoadBalancer
  selector:
    app: myapp
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
```

### ConfigMap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: myapp-config
data:
  APP_NAME: "MyApp"
  config.yaml: |
    server:
      port: 8080
```

### Secret

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: myapp-secret
type: Opaque
data:
  PASSWORD: cGFzc3dvcmQxMjM=  # base64
```

---

## 🔥 Common Patterns

### Terraform: Remote State

```hcl
terraform {
  backend "s3" {
    bucket         = "terraform-state"
    key            = "prod/terraform.tfstate"
    region         = "us-east-1"
    dynamodb_table = "terraform-locks"
    encrypt        = true
  }
}
```

### Terraform: Variables File

```bash
# terraform.tfvars
region        = "us-east-1"
instance_type = "t2.micro"
```

### K8s: Multi-container Pod

```yaml
spec:
  containers:
  - name: app
    image: myapp:1.0
  - name: sidecar
    image: logger:1.0
```

### K8s: Health Checks

```yaml
spec:
  containers:
  - name: myapp
    image: myapp:1.0
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

### K8s: ConfigMap in Pod

```yaml
spec:
  containers:
  - name: myapp
    image: myapp:1.0
    env:
    - name: APP_NAME
      valueFrom:
        configMapKeyRef:
          name: myapp-config
          key: APP_NAME
    envFrom:
    - configMapRef:
        name: myapp-config
```

### K8s: Secret in Pod

```yaml
spec:
  containers:
  - name: myapp
    image: myapp:1.0
    env:
    - name: DB_PASSWORD
      valueFrom:
        secretKeyRef:
          name: myapp-secret
          key: PASSWORD
```

---

## ⚠️ Common Errors

### Terraform

| Error | Solution |
|-------|----------|
| State locked | `terraform force-unlock <id>` |
| Provider not found | `terraform init` |
| Resource already exists | `terraform import` or delete manually |
| Invalid syntax | `terraform validate` |

### Kubernetes

| Error | Solution |
|-------|----------|
| ImagePullBackOff | Check image name, check registry auth |
| CrashLoopBackOff | Check logs: `kubectl logs <pod>` |
| Pending | Check resources: `kubectl describe pod <pod>` |
| Error: connection refused | Check service/port configuration |
| OOMKilled | Increase memory limits |
| Context deadline exceeded | Increase timeouts |

---

## 🚀 Quick Reference

| Task | Terraform | Kubernetes |
|------|-----------|------------|
| Initialize | `terraform init` | - |
| Create | `terraform apply` | `kubectl apply -f file.yaml` |
| Update | `terraform apply` | `kubectl apply -f file.yaml` |
| Delete | `terraform destroy` | `kubectl delete -f file.yaml` |
| List | `terraform state list` | `kubectl get all` |
| Details | `terraform show` | `kubectl describe <resource>` |
| Logs | - | `kubectl logs <pod>` |
| Shell | - | `kubectl exec -it <pod> -- sh` |

---

## 📚 Best Practices

### Terraform
- ✅ Use remote state (S3 + DynamoDB)
- ✅ Version providers
- ✅ Use modules for reusability
- ✅ `.gitignore`: `*.tfstate`, `*.tfvars`
- ✅ Use workspaces for environments

### Kubernetes
- ✅ Use namespaces for separation
- ✅ Set resource limits
- ✅ Use health checks (liveness, readiness)
- ✅ Use ConfigMaps/Secrets (not hardcode)
- ✅ Use labels for organization
- ✅ Version control manifests
