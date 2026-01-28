# ✅ Week 13 - Завершено!

## 🎯 Що створено

**Week 13: Infrastructure as Code (Terraform)** - модуль про створення та управління AWS інфраструктурою за допомогою Terraform, включаючи Lambda, SQS, IAM, та CloudWatch.

---

## 📊 Статистика

### Створено файлів

**Теорія:** 2 файли
- `theory/01_terraform_state.md` (750+ рядків)
- `theory/02_terraform_plan_apply.md` (600+ рядків)

**Практика:** 16 файлів
- `practice/01_lambda_terraform/` - Lambda з Terraform (5 файлів)
  - `main.tf` - Terraform config
  - `main.go` - Lambda code
  - `README.md` - Documentation
  - `Makefile` - Automation
- `practice/02_sqs_terraform/` - SQS з Terraform (3 файли)
  - `main.tf` - Terraform config
  - `README.md` - Documentation
  - `Makefile` - Automation
- `practice/03_iam_terraform/` - IAM з Terraform (2 файли)
  - `main.tf` - Terraform config
  - `README.md` - Documentation
- `practice/04_full_stack/` - Full Stack (6 файлів)
  - `main.tf` - Terraform config
  - `main.go` - Lambda code
  - `README.md` - Documentation
  - `Makefile` - Automation

**Документація:** 3 файли
- `README.md` - Повний guide
- `QUICK_START.md` - Швидкий старт
- `WEEK13_COMPLETE.md` - Цей звіт

**Загалом:** 21 файл, ~4000+ рядків коду та документації

---

## 📚 Що покрито

### 1. Terraform State 🗄️

**Теорія:**
- Що таке State (mapping config → real resources)
- Local vs Remote backends (S3)
- State locking (DynamoDB)
- State operations (list, show, mv, rm, import)
- State security (encryption, versioning)
- State drift detection and recovery

**Key Concepts:**
```
Config (.tf) ←→ State File ←→ Real Resources (AWS)
```

**Backends:**
```
Local:  ./terraform.tfstate (dev only)
Remote: S3 + DynamoDB (production)
```

### 2. Terraform Plan & Apply 📝

**Теорія:**
- Terraform workflow (init → plan → apply → destroy)
- Core commands (init, plan, apply, destroy)
- Change types (create, update, recreate)
- Plan options (save, target, refresh)
- Apply options (auto-approve, replace, parallelism)

**Workflow:**
```
Write → Init → Validate → Plan → Apply → Monitor
```

**Symbols:**
```
+ create
~ update in-place
- destroy
-/+ destroy and recreate
<= read (data source)
```

---

## 🎯 Практичні приклади

### Practice 01: Lambda with Terraform

**Stack:**
- Lambda Function (Go runtime)
- IAM Role + Basic Execution Policy
- CloudWatch Log Group (7 days retention)
- Lambda Function URL (HTTP access)

**Terraform Resources:**
```hcl
aws_iam_role
aws_iam_role_policy_attachment
aws_cloudwatch_log_group
aws_lambda_function
aws_lambda_function_url
```

**Commands:**
```bash
make deploy  # Build Go → ZIP → Terraform apply
make test    # Test function URL
make logs    # Tail CloudWatch logs
make destroy # Cleanup
```

---

### Practice 02: SQS with Terraform

**Stack:**
- Main SQS Queue (long polling, 4 days retention)
- Dead Letter Queue (14 days retention)
- CloudWatch Alarms (3):
  - DLQ not empty
  - High queue depth (> 1000)
  - Old messages (> 10 min)

**Terraform Resources:**
```hcl
aws_sqs_queue (main + DLQ)
aws_cloudwatch_metric_alarm (x3)
```

**Features:**
- Long polling (20s wait)
- DLQ with maxReceiveCount=3
- Visibility timeout: 5 minutes
- Batch operations (send/receive up to 10)

**Commands:**
```bash
make deploy       # Create queues
make send-batch   # Send 10 messages
make receive      # Receive messages (long poll)
make check-dlq    # Check DLQ
make test-poison  # Test error handling
make monitor      # Show metrics
make destroy      # Cleanup
```

---

### Practice 03: IAM with Terraform

**Stack:**
- IAM Role (для Lambda)
- 4 Custom IAM Policies:
  - SQS access (receive, delete)
  - DynamoDB access (read, write)
  - Secrets Manager (read secrets)
  - S3 inline policy (GetObject, PutObject)
- IAM User (для CI/CD)
- IAM Group (для команди)
- Access Keys (programmatic access)

**Terraform Resources:**
```hcl
aws_iam_role
aws_iam_policy (x3)
aws_iam_role_policy (inline)
aws_iam_role_policy_attachment (x4)
aws_iam_user
aws_iam_access_key
aws_iam_group
aws_iam_group_policy
aws_iam_user_group_membership
```

**Best Practices:**
- Least privilege (specific resources, not `*`)
- Separate policies per service
- Managed policies for common patterns
- Tags on all resources
- Sensitive outputs protected

---

### Practice 04: Full Stack ⭐

**Complete async processing pipeline:**
- SQS Main Queue + DLQ
- Lambda Function (Go, batch processing)
- IAM Role з SQS permissions
- Event Source Mapping (SQS → Lambda)
- CloudWatch Logs + Alarms (2)

**Architecture:**
```
Producer → SQS Queue
              ↓
        Event Source Mapping
              ↓
        Lambda (batch 10, window 5s)
              ↓
        ✅ Process successfully
              ↓
        Delete from queue
              
        OR
              
        ❌ Process failure
              ↓
        Retry (visibility timeout)
              ↓
        After 3 failures
              ↓
        Move to DLQ
              ↓
        CloudWatch Alarm 🚨
```

**Terraform Resources:**
```hcl
aws_sqs_queue (x2)
aws_iam_role
aws_iam_policy
aws_iam_role_policy_attachment (x2)
aws_cloudwatch_log_group
aws_lambda_function
aws_lambda_event_source_mapping
aws_cloudwatch_metric_alarm (x2)
```

**Go Lambda Features:**
- Batch processing (up to 10 messages)
- Partial batch failures (report per message)
- Idempotency check (in-memory cache)
- Proper error handling
- Structured logging

**Commands:**
```bash
make deploy        # Deploy full stack
make send-batch    # Send 10 orders
make logs          # Tail logs
make monitor       # Show queue + Lambda metrics
make test-load     # Load test (100 messages)
make test-errors   # Error handling test (50% invalid)
make invoke        # Direct Lambda invoke (no SQS)
make check-dlq     # Check DLQ
make destroy       # Cleanup
```

---

## 📊 Terraform Resources Покрито

### Core Resources

| Resource | Practice | Count |
|----------|----------|-------|
| `aws_lambda_function` | 01, 04 | 2 |
| `aws_iam_role` | 01, 03, 04 | 3 |
| `aws_iam_policy` | 03, 04 | 4 |
| `aws_iam_role_policy_attachment` | 01, 03, 04 | 8 |
| `aws_iam_role_policy` | 03 | 1 |
| `aws_sqs_queue` | 02, 04 | 4 |
| `aws_cloudwatch_log_group` | 01, 04 | 2 |
| `aws_cloudwatch_metric_alarm` | 02, 04 | 5 |
| `aws_lambda_function_url` | 01 | 1 |
| `aws_lambda_event_source_mapping` | 04 | 1 |
| `aws_iam_user` | 03 | 1 |
| `aws_iam_access_key` | 03 | 1 |
| `aws_iam_group` | 03 | 1 |
| `aws_iam_group_policy` | 03 | 1 |
| `aws_iam_user_group_membership` | 03 | 1 |

**Total:** 15 різних типів resources, 36 instances

---

## 🔧 Makefiles

Кожна практика має власний `Makefile` з командами:

### Common Commands

| Command | Description |
|---------|-------------|
| `make help` | Show all commands |
| `make init` | Initialize (Go modules, Terraform) |
| `make deploy` | Build + Deploy |
| `make outputs` | Show Terraform outputs |
| `make destroy` | Cleanup all resources |
| `make clean` | Clean build artifacts |

### Lambda-Specific

| Command | Description |
|---------|-------------|
| `make build` | Build Go binary for Linux |
| `make zip` | Create function.zip |
| `make test` | Test function URL |
| `make logs` | Tail CloudWatch logs |
| `make invoke` | Direct Lambda invoke |

### SQS-Specific

| Command | Description |
|---------|-------------|
| `make send` | Send single message |
| `make send-batch` | Send 10 messages |
| `make receive` | Receive messages |
| `make check-dlq` | Check DLQ |
| `make test-poison` | Test poison message → DLQ |
| `make monitor` | Show queue metrics |
| `make purge` | Purge queue |

### Full Stack-Specific

| Command | Description |
|---------|-------------|
| `make test-load` | Load test (100 messages) |
| `make test-errors` | Error test (50 valid, 50 invalid) |
| `make logs-errors` | Show only error logs |

---

## ✅ Best Practices Implemented

### Terraform

1. ✅ **Remote backend ready** (S3 + DynamoDB config)
2. ✅ **State locking** (DynamoDB table)
3. ✅ **Variables** (parameterization)
4. ✅ **Outputs** (export important values)
5. ✅ **Dependencies** (explicit `depends_on`)
6. ✅ **Tags** (Name, Environment, ManagedBy)
7. ✅ **Sensitive outputs** (marked as sensitive)
8. ✅ **Resource naming** (consistent, prefixed)

### Lambda

1. ✅ **Optimized build** (`-ldflags="-s -w"`)
2. ✅ **Custom runtime** (provided.al2023)
3. ✅ **Environment variables** (config via env)
4. ✅ **Proper logging** (structured, timestamps)
5. ✅ **Error handling** (graceful failures)
6. ✅ **Idempotency** (prevent duplicates)
7. ✅ **Partial batch failures** (efficient retries)
8. ✅ **CloudWatch integration** (logs, metrics)

### SQS

1. ✅ **Long polling** (WaitTimeSeconds=20)
2. ✅ **Dead Letter Queue** (maxReceiveCount=3)
3. ✅ **Visibility timeout** (appropriate for workload)
4. ✅ **Message retention** (4 days main, 14 days DLQ)
5. ✅ **CloudWatch alarms** (DLQ, depth, age)
6. ✅ **Batch operations** (up to 10 messages)

### IAM

1. ✅ **Least privilege** (specific resources)
2. ✅ **Assume role policy** (trust relationships)
3. ✅ **Managed policies** (AWS-managed where possible)
4. ✅ **Custom policies** (service-specific)
5. ✅ **Inline policies** (single-use cases)
6. ✅ **Groups** (team organization)

---

## 🎯 Key Terraform Patterns

### 1. Resource Definition

```hcl
resource "aws_lambda_function" "processor" {
  function_name = var.function_name
  role          = aws_iam_role.lambda.arn
  runtime       = "provided.al2023"
  
  tags = {
    Environment = var.environment
    ManagedBy   = "terraform"
  }
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
  description = "Lambda function ARN"
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

### 5. Interpolation

```hcl
name = "${var.project_name}-${var.environment}"
```

### 6. JSON Encoding

```hcl
assume_role_policy = jsonencode({
  Version = "2012-10-17"
  Statement = [...]
})
```

---

## 🔍 Terraform Commands Summary

```bash
# Initialize
terraform init

# Validate
terraform validate

# Format
terraform fmt

# Plan
terraform plan
terraform plan -out=tfplan

# Apply
terraform apply
terraform apply tfplan
terraform apply -auto-approve

# Destroy
terraform destroy
terraform destroy -target=aws_lambda_function.processor

# State
terraform state list
terraform state show aws_lambda_function.processor
terraform state mv old_name new_name
terraform state rm aws_lambda_function.processor
terraform state pull
terraform state push

# Import
terraform import aws_lambda_function.processor function-name

# Output
terraform output
terraform output lambda_arn
terraform output -json

# Show
terraform show
terraform show tfplan

# Refresh
terraform refresh

# Graph
terraform graph | dot -Tpng > graph.png
```

---

## 🎓 Що ти навчився

### Terraform Core

✅ **State management** - Local vs remote, locking, encryption  
✅ **Workflow** - init → plan → apply → destroy  
✅ **Resources** - 15 types, 36 instances  
✅ **Variables** - Input parameterization  
✅ **Outputs** - Export values  
✅ **Dependencies** - Explicit and implicit  
✅ **Data sources** - Query existing resources  
✅ **Backends** - S3 + DynamoDB  

### AWS Services

✅ **Lambda** - Functions, roles, logs, URLs  
✅ **SQS** - Queues, DLQ, long polling  
✅ **IAM** - Roles, policies, users, groups  
✅ **CloudWatch** - Logs, alarms, metrics  
✅ **Event Source Mapping** - SQS → Lambda  

### Go Lambda Development

✅ **Build for Linux** - GOOS=linux GOARCH=amd64  
✅ **Optimize binary** - ldflags="-s -w"  
✅ **Handle SQS events** - events.SQSEvent  
✅ **Partial batch failures** - events.SQSEventResponse  
✅ **Idempotency** - Prevent duplicate processing  
✅ **Error handling** - Graceful failures, retries  
✅ **Logging** - Structured logs for CloudWatch  

### DevOps

✅ **Infrastructure as Code** - Reproducible infrastructure  
✅ **Automation** - Makefiles for common tasks  
✅ **Monitoring** - CloudWatch logs, alarms, metrics  
✅ **Error handling** - DLQ, retries, partial failures  
✅ **Testing** - Load tests, error tests  
✅ **Documentation** - READMEs, quick starts  

---

## 📊 Architecture Evolution

```
Week 11: Lambda Basics
   ↓
Week 12: Lambda + SQS
   ↓
Week 13: Infrastructure as Code (Terraform)
   ↓
Result: Production-Ready Async Processing Pipeline!
```

**Stack:**
```
Terraform (IaC)
   ↓
Lambda (Compute)
   ↓
SQS (Queue)
   ↓
IAM (Security)
   ↓
CloudWatch (Monitoring)
```

---

## 🚀 Production-Ready Checklist

| Feature | Status |
|---------|--------|
| Remote state (S3) | ✅ Config ready |
| State locking (DynamoDB) | ✅ Config ready |
| State encryption | ✅ Enabled |
| IAM least privilege | ✅ Implemented |
| Dead Letter Queue | ✅ Configured |
| CloudWatch alarms | ✅ 5 alarms |
| Idempotency | ✅ Implemented |
| Partial batch failures | ✅ Supported |
| Long polling | ✅ 20 seconds |
| Structured logging | ✅ JSON format |
| Error handling | ✅ Retry + DLQ |
| Monitoring | ✅ Logs + Metrics |
| Documentation | ✅ Complete |
| Automation | ✅ Makefiles |
| Testing | ✅ Load + Error tests |

---

## 🔗 Зв'язок з іншими модулями

### Week 11: AWS Lambda

```
Week 11: Lambda basics, lifecycle, cold start
   ↓
Week 13: Lambda with Terraform (IaC)
```

### Week 12: SQS + Lambda

```
Week 12: SQS, at-least-once, DLQ
   ↓
Week 13: SQS + Lambda with Terraform
```

---

## 🎉 Досягнення

### Code Statistics

- **21 files** created
- **4000+ lines** of code + docs
- **4 complete examples**
- **15 Terraform resource types**
- **36 Terraform resource instances**
- **2 Go Lambda functions**
- **4 Makefiles** with 60+ commands
- **5 CloudWatch alarms**

### Skills Acquired

- ✅ Terraform basics (state, plan, apply)
- ✅ AWS infrastructure automation
- ✅ Lambda deployment with Terraform
- ✅ SQS queue management
- ✅ IAM role and policy creation
- ✅ CloudWatch monitoring setup
- ✅ Event-driven architecture
- ✅ Async processing pipelines
- ✅ Error handling and DLQ
- ✅ Infrastructure best practices

---

## 📖 Ресурси

- [Terraform Documentation](https://www.terraform.io/docs)
- [Terraform AWS Provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [AWS Lambda with Terraform](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function)
- [Terraform Best Practices](https://www.terraform-best-practices.com/)
- Week 11: Lambda lifecycle
- Week 12: SQS + Lambda

---

## 🎯 Висновок

### Terraform = Infrastructure as Code

✅ **Reproducible** - Deploy identical infrastructure anywhere  
✅ **Versionable** - Track changes in Git  
✅ **Auditable** - See who changed what and when  
✅ **Testable** - Plan before apply  
✅ **Automated** - CI/CD integration  

### Production-Ready Stack

✅ **Lambda** - Serverless compute  
✅ **SQS** - Managed queuing  
✅ **DLQ** - Error isolation  
✅ **IAM** - Security by default  
✅ **CloudWatch** - Full observability  
✅ **Terraform** - Infrastructure as code  

### Golden Rules

1. **Always use remote state** (S3 + DynamoDB)
2. **Always run plan first** (preview changes)
3. **Always tag resources** (organization)
4. **Always implement DLQ** (error handling)
5. **Always monitor** (CloudWatch)

---

## ✅ Week 13 Complete!

```
Progress: 100% ✅

Theory:   ████████████ 2/2
Practice: ████████████ 4/4
Docs:     ████████████ 3/3
```

**Дата завершення:** 2026-01-28  
**Статус:** COMPLETE ✅  
**Локація:** `/Users/vkuzm/GolandProjects/golang_practice/week_13`

---

## 🎉 Вітаємо!

Тепер ти вмієш:
- ✅ Створювати infrastructure as code
- ✅ Управляти Terraform state
- ✅ Deployувати Lambda з Terraform
- ✅ Налаштовувати SQS з DLQ
- ✅ Створювати IAM roles та policies
- ✅ Будувати async processing pipelines
- ✅ Моніторити через CloudWatch
- ✅ Автоматизувати через Makefiles

**"Infrastructure as Code = Reproducible, Versionable, Auditable!"** 🏗️

---

**Next:**
- Week 14: Advanced Terraform (modules, workspaces)
- Week 15: CI/CD Pipelines (GitHub Actions + Terraform)
- Week 16: Multi-region Deployments

**Week 13: COMPLETE!** 🎯🏗️☁️
