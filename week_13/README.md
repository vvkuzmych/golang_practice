# Week 13 - Infrastructure as Code (Terraform)

## 🎯 Мета

Навчитися створювати та управляти AWS інфраструктурою за допомогою **Terraform** (Infrastructure as Code).

---

## 📚 Теорія

### 1. [Terraform State](./theory/01_terraform_state.md)

**Що таке State:**
- Mapping між config та real resources
- Local vs Remote backends (S3)
- State locking (DynamoDB)
- State security (encryption, versioning)

**Key Concepts:**
```
Config (.tf) ←→ State File ←→ Real Resources (AWS)
```

### 2. [Terraform Plan & Apply](./theory/02_terraform_plan_apply.md)

**Terraform Workflow:**
```
Write → Init → Plan → Apply → Destroy
```

**Core Commands:**
- `terraform init` - Initialize, download providers
- `terraform plan` - Preview changes
- `terraform apply` - Execute changes
- `terraform destroy` - Delete resources

---

## 🎯 Практика

### [01: Lambda with Terraform](./practice/01_lambda_terraform/)

**Що створюється:**
- Lambda Function (Go runtime)
- IAM Role з policies
- CloudWatch Log Group
- Lambda Function URL

**Commands:**
```bash
cd practice/01_lambda_terraform
make deploy    # Build + deploy
make test      # Test function
make logs      # View logs
make destroy   # Cleanup
```

---

### [02: SQS with Terraform](./practice/02_sqs_terraform/)

**Що створюється:**
- SQS Main Queue (long polling)
- Dead Letter Queue (DLQ)
- CloudWatch Alarms (3 alarms)

**Commands:**
```bash
cd practice/02_sqs_terraform
make deploy        # Deploy queues
make send-batch    # Send 10 messages
make receive       # Receive messages
make check-dlq     # Check DLQ
make test-poison   # Test error handling
make destroy       # Cleanup
```

---

### [03: IAM with Terraform](./practice/03_iam_terraform/)

**Що створюється:**
- IAM Role (для Lambda)
- Custom IAM Policies (SQS, DynamoDB, Secrets)
- IAM User (для CI/CD)
- IAM Group (для команди)
- Access Keys

**Commands:**
```bash
cd practice/03_iam_terraform
make deploy       # Deploy IAM resources
make outputs      # Show all outputs
make destroy      # Cleanup
```

---

### [04: Full Stack](./practice/04_full_stack/) ⭐

**Повний async processing pipeline:**
- SQS Queue + DLQ
- Lambda Function (Go)
- IAM Role з permissions
- Event Source Mapping (SQS → Lambda)
- CloudWatch Logs + Alarms

**Architecture:**
```
Producer → SQS → Lambda → CloudWatch
                   ↓
                  DLQ → Alarm 🚨
```

**Commands:**
```bash
cd practice/04_full_stack
make deploy        # Deploy full stack
make send-batch    # Send 10 orders
make logs          # Tail logs
make monitor       # Show metrics
make test-load     # Load test (100 msgs)
make test-errors   # Error handling test
make destroy       # Cleanup
```

---

## 🔧 Prerequisites

### Install Terraform

```bash
# macOS
brew install terraform

# Verify
terraform version
```

### Configure AWS CLI

```bash
aws configure

# Set credentials
AWS Access Key ID: YOUR_KEY
AWS Secret Access Key: YOUR_SECRET
Default region: us-east-1
```

---

## 🎯 Learning Path

### Step 1: Understand State

```bash
cd practice/01_lambda_terraform
make deploy

# Check state
terraform state list
terraform state show aws_lambda_function.processor
terraform show
```

### Step 2: Practice Plan & Apply

```bash
# Modify main.tf (change memory_size)
vim main.tf

# See what will change
terraform plan

# Apply changes
terraform apply
```

### Step 3: Build Individual Components

```bash
# Lambda
cd practice/01_lambda_terraform && make deploy

# SQS
cd practice/02_sqs_terraform && make deploy

# IAM
cd practice/03_iam_terraform && make deploy
```

### Step 4: Build Full Stack

```bash
cd practice/04_full_stack
make deploy
make send-batch
make logs
make monitor
```

---

## 📊 Key Terraform Concepts

### 1. Resources

```hcl
resource "aws_lambda_function" "processor" {
  function_name = "my-function"
  runtime       = "provided.al2023"
  memory_size   = 1024
}
```

### 2. Variables

```hcl
variable "memory_size" {
  description = "Lambda memory in MB"
  type        = number
  default     = 1024
}
```

### 3. Outputs

```hcl
output "lambda_arn" {
  description = "Lambda ARN"
  value       = aws_lambda_function.processor.arn
}
```

### 4. Dependencies

```hcl
resource "aws_lambda_function" "processor" {
  depends_on = [
    aws_cloudwatch_log_group.lambda,
    aws_iam_role_policy_attachment.lambda_basic
  ]
}
```

### 5. Data Sources

```hcl
data "aws_caller_identity" "current" {}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}
```

---

## ✅ Best Practices

### 1. Use Remote Backend

```hcl
terraform {
  backend "s3" {
    bucket         = "my-terraform-state"
    key            = "prod/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "terraform-locks"
  }
}
```

### 2. Enable State Locking

```hcl
# Use DynamoDB for locking
terraform {
  backend "s3" {
    dynamodb_table = "terraform-locks"
  }
}
```

### 3. Separate Environments

```
s3://bucket/dev/terraform.tfstate
s3://bucket/staging/terraform.tfstate
s3://bucket/prod/terraform.tfstate
```

### 4. Use Variables

```hcl
variable "environment" {
  type = string
}

resource "aws_lambda_function" "processor" {
  function_name = "${var.environment}-processor"
}
```

### 5. Tag Everything

```hcl
tags = {
  Name        = var.function_name
  Environment = var.environment
  ManagedBy   = "terraform"
}
```

### 6. Never Commit State

```bash
# .gitignore
*.tfstate
*.tfstate.*
.terraform/
```

---

## 🔍 Common Commands

```bash
# Initialize
terraform init

# Validate syntax
terraform validate

# Format code
terraform fmt

# Plan changes
terraform plan

# Save plan
terraform plan -out=tfplan

# Apply
terraform apply

# Apply saved plan
terraform apply tfplan

# Destroy
terraform destroy

# Show state
terraform show

# List resources
terraform state list

# Show resource
terraform state show aws_lambda_function.processor

# Output values
terraform output

# Refresh state
terraform refresh

# Import existing resource
terraform import aws_lambda_function.processor function-name
```

---

## 📊 Typical Workflow

### Development

```bash
# 1. Write config
vim main.tf

# 2. Initialize
terraform init

# 3. Validate
terraform validate

# 4. Format
terraform fmt

# 5. Plan
terraform plan

# 6. Apply
terraform apply

# 7. Verify
aws lambda get-function --function-name my-function

# 8. Iterate
vim main.tf
terraform plan
terraform apply
```

### Production

```bash
# 1. Pull latest
git pull

# 2. Initialize
terraform init

# 3. Plan and save
terraform plan -out=prod.tfplan

# 4. Review
terraform show prod.tfplan

# 5. Get approval

# 6. Apply
terraform apply prod.tfplan

# 7. Verify
# Check AWS Console

# 8. Tag release
git tag v1.0.0
```

---

## 🎓 Що ти навчився

### Terraform Basics

✅ **State management** (local vs remote)  
✅ **Plan/Apply workflow** (preview → execute)  
✅ **Resources** (Lambda, SQS, IAM)  
✅ **Variables & Outputs** (parameterization)  
✅ **Dependencies** (depends_on, implicit)  

### AWS Resources

✅ **Lambda** (function, role, logs)  
✅ **SQS** (queue, DLQ, alarms)  
✅ **IAM** (roles, policies, users, groups)  
✅ **CloudWatch** (logs, alarms, metrics)  
✅ **Event Source Mapping** (SQS → Lambda)  

### Best Practices

✅ **Remote backend** (S3 + DynamoDB)  
✅ **State locking** (prevent conflicts)  
✅ **Encryption** (secure state)  
✅ **Versioning** (rollback capability)  
✅ **Least privilege** (IAM policies)  
✅ **Tagging** (resource organization)  

---

## 🚀 Next Steps

### Week 14: Advanced Terraform

- Modules (reusable components)
- Workspaces (multiple environments)
- Data sources (query existing resources)
- Provisioners (post-creation scripts)
- Terraform Cloud (team collaboration)

### Production Considerations

- Multi-region deployment
- Blue-green deployments
- Canary releases
- Cost optimization
- Security hardening
- Compliance (SOC2, HIPAA)

---

## 📖 Resources

- [Terraform Documentation](https://www.terraform.io/docs)
- [Terraform AWS Provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [AWS Lambda Best Practices](https://docs.aws.amazon.com/lambda/latest/dg/best-practices.html)
- [Terraform Best Practices](https://www.terraform-best-practices.com/)
- Week 11: AWS Lambda basics
- Week 12: SQS + Lambda integration

---

## 🎉 Вітаємо!

Тепер ти вмієш:
- ✅ Створювати інфраструктуру as code
- ✅ Управляти Terraform state
- ✅ Використовувати plan/apply workflow
- ✅ Створювати Lambda з Terraform
- ✅ Налаштовувати SQS з DLQ
- ✅ Управляти IAM permissions
- ✅ Будувати async processing pipelines
- ✅ Моніторити через CloudWatch

**"Infrastructure as Code = Reproducible, Versionable, Auditable!"** 🏗️

---

**Week 13: COMPLETE!** ✅

**Progress: Lambda → SQS → IaC → Production-Ready!** 📨☁️🏗️
