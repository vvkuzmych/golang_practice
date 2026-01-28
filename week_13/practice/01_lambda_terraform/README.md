# Practice 01: Lambda with Terraform

## 🎯 Мета

Навчитися створювати AWS Lambda функцію за допомогою Terraform.

---

## 📊 Що створюється

1. **Lambda Function** - Go runtime (provided.al2023)
2. **IAM Role** - З правами для виконання Lambda
3. **CloudWatch Log Group** - Для логів
4. **Lambda Function URL** - HTTP endpoint

---

## 🏗️ Архітектура

```
Internet
   ↓
Lambda Function URL (HTTPS)
   ↓
Lambda Function (Go)
   ↓
CloudWatch Logs
```

---

## 🚀 Швидкий старт

### 1. Build Lambda

```bash
# Initialize Go module
go mod init lambda-example
go get github.com/aws/aws-lambda-go/lambda
go get github.com/aws/aws-lambda-go/events

# Build for Linux (Lambda runtime)
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bootstrap main.go

# Create ZIP
zip function.zip bootstrap
```

### 2. Initialize Terraform

```bash
terraform init
```

### 3. Plan

```bash
terraform plan
```

### 4. Apply

```bash
terraform apply
```

### 5. Test Lambda

```bash
# Get Lambda URL from output
LAMBDA_URL=$(terraform output -raw lambda_url)

# Test with curl
curl -X POST $LAMBDA_URL \
  -H "Content-Type: application/json" \
  -d '{"name":"John","message":"Hello from Terraform!"}'

# Expected response:
# {"message":"Hello John! Your message: Hello from Terraform!","stage":"dev"}
```

### 6. Check Logs

```bash
# Get log group name
LOG_GROUP=$(terraform output -raw log_group_name)

# View logs
aws logs tail $LOG_GROUP --follow
```

### 7. Cleanup

```bash
terraform destroy
```

---

## 📝 Конфігурація

### Variables

**`terraform.tfvars`** (optional):
```hcl
aws_region    = "us-east-1"
function_name = "my-processor"
memory_size   = 1024
timeout       = 30
```

### Outputs

After `terraform apply`:
```
lambda_arn       = "arn:aws:lambda:us-east-1:123:function:my-processor"
lambda_name      = "my-processor"
lambda_url       = "https://abc123.lambda-url.us-east-1.on.aws/"
lambda_role_arn  = "arn:aws:iam::123:role/my-processor-role"
log_group_name   = "/aws/lambda/my-processor"
```

---

## 🔧 Налаштування

### Increase Memory

```hcl
# main.tf
variable "memory_size" {
  default = 2048  # 2 GB
}
```

### Increase Timeout

```hcl
variable "timeout" {
  default = 60  # 60 seconds
}
```

### Add Environment Variables

```hcl
# main.tf
resource "aws_lambda_function" "processor" {
  environment {
    variables = {
      STAGE       = "prod"
      LOG_LEVEL   = "DEBUG"
      API_KEY     = var.api_key
    }
  }
}
```

### Enable IAM Auth for Function URL

```hcl
# main.tf
resource "aws_lambda_function_url" "processor" {
  authorization_type = "AWS_IAM"  # Requires signed requests
}
```

---

## 📊 Terraform State

After apply, check state:

```bash
# List resources
terraform state list

# Show Lambda details
terraform state show aws_lambda_function.processor

# Show outputs
terraform output
```

---

## ⚙️ Advanced

### Use Remote Backend

```hcl
# main.tf
terraform {
  backend "s3" {
    bucket         = "my-terraform-state"
    key            = "lambda/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "terraform-locks"
  }
}
```

### Add Alias

```hcl
resource "aws_lambda_alias" "prod" {
  name             = "prod"
  function_name    = aws_lambda_function.processor.function_name
  function_version = aws_lambda_function.processor.version
}
```

### Add Dead Letter Queue

```hcl
resource "aws_sqs_queue" "dlq" {
  name = "${var.function_name}-dlq"
}

resource "aws_lambda_function" "processor" {
  dead_letter_config {
    target_arn = aws_sqs_queue.dlq.arn
  }
}
```

---

## ✅ Перевірка

### 1. Lambda Created

```bash
aws lambda get-function --function-name my-processor
```

### 2. Function URL Works

```bash
curl -X POST $(terraform output -raw lambda_url) \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","message":"Hello"}'
```

### 3. Logs in CloudWatch

```bash
aws logs tail /aws/lambda/my-processor --follow
```

---

## 🎓 Висновок

Ти навчився:
- ✅ Створювати Lambda з Terraform
- ✅ Налаштовувати IAM роль
- ✅ Додавати CloudWatch логи
- ✅ Створювати Function URL
- ✅ Використовувати outputs

**Далі:** `02_sqs_terraform/` - SQS з Terraform
