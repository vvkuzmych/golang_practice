# Practice 03: IAM with Terraform

## 🎯 Мета

Навчитися створювати IAM roles, policies, users, та groups за допомогою Terraform.

---

## 📊 Що створюється

1. **IAM Role** - Для Lambda з assume role policy
2. **IAM Policies** - Custom policies (SQS, DynamoDB, Secrets Manager)
3. **IAM User** - Для CI/CD deployment
4. **IAM Group** - Для команди розробників
5. **Access Keys** - Для programmatic access

---

## 🏗️ Архітектура

```
IAM Role (Lambda)
  ├─ AWSLambdaBasicExecutionRole (managed)
  ├─ Custom SQS Policy
  ├─ Custom DynamoDB Policy
  ├─ Custom Secrets Manager Policy
  └─ Inline S3 Policy

IAM User (Deployer)
  ├─ Access Key + Secret
  └─ Deployer Policy

IAM Group (Developers)
  ├─ Group Policy
  └─ Members: [Deployer]
```

---

## 🚀 Швидкий старт

### 1. Initialize Terraform

```bash
terraform init
```

### 2. Plan

```bash
terraform plan
```

### 3. Apply

```bash
terraform apply
```

### 4. Get Outputs

```bash
# All outputs
terraform output

# Specific output
terraform output lambda_role_arn

# Sensitive output (access key secret)
terraform output -raw deployer_secret_access_key
```

### 5. Test Role

```bash
# Assume role (requires AWS CLI configured)
aws sts assume-role \
  --role-arn $(terraform output -raw lambda_role_arn) \
  --role-session-name test-session
```

### 6. Test User Access

```bash
# Configure profile with deployer credentials
aws configure --profile deployer

# Set access keys from Terraform output
AWS_ACCESS_KEY_ID=$(terraform output -raw deployer_access_key_id)
AWS_SECRET_ACCESS_KEY=$(terraform output -raw deployer_secret_access_key)

# Test access
AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID \
AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY \
aws sts get-caller-identity
```

### 7. Cleanup

```bash
terraform destroy
```

---

## 📝 IAM Components

### 1. IAM Role

**Purpose:** Для AWS services (Lambda, EC2, etc.)

```hcl
resource "aws_iam_role" "lambda" {
  name = "my-lambda-role"
  
  assume_role_policy = jsonencode({
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })
}
```

**Trust policy:** Хто може assume цю роль

### 2. IAM Policy (Managed)

**Purpose:** Reusable permission set

```hcl
resource "aws_iam_policy" "sqs_access" {
  name = "sqs-access-policy"
  
  policy = jsonencode({
    Statement = [{
      Effect = "Allow"
      Action = [
        "sqs:ReceiveMessage",
        "sqs:DeleteMessage"
      ]
      Resource = "*"
    }]
  })
}
```

**Attach to role:**
```hcl
resource "aws_iam_role_policy_attachment" "lambda_sqs" {
  role       = aws_iam_role.lambda.name
  policy_arn = aws_iam_policy.sqs_access.arn
}
```

### 3. IAM Policy (Inline)

**Purpose:** Policy embedded directly in role

```hcl
resource "aws_iam_role_policy" "lambda_inline" {
  name = "inline-policy"
  role = aws_iam_role.lambda.id
  
  policy = jsonencode({
    Statement = [{
      Effect   = "Allow"
      Action   = ["s3:GetObject"]
      Resource = "arn:aws:s3:::my-bucket/*"
    }]
  })
}
```

**Inline vs Managed:**
- Inline: Deleted with role, single use
- Managed: Reusable, can attach to multiple roles

### 4. IAM User

**Purpose:** Programmatic or console access

```hcl
resource "aws_iam_user" "deployer" {
  name = "deployer"
}

resource "aws_iam_access_key" "deployer" {
  user = aws_iam_user.deployer.name
}
```

**Access keys:** For API/CLI access

### 5. IAM Group

**Purpose:** Manage permissions for multiple users

```hcl
resource "aws_iam_group" "developers" {
  name = "developers"
}

resource "aws_iam_group_policy" "developers" {
  group  = aws_iam_group.developers.name
  policy = jsonencode({ ... })
}

resource "aws_iam_user_group_membership" "deployer" {
  user   = aws_iam_user.deployer.name
  groups = [aws_iam_group.developers.name]
}
```

---

## 🔐 Policy Examples

### SQS Consumer

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "sqs:ReceiveMessage",
        "sqs:DeleteMessage",
        "sqs:GetQueueAttributes",
        "sqs:ChangeMessageVisibility"
      ],
      "Resource": "arn:aws:sqs:us-east-1:123456789012:my-queue"
    }
  ]
}
```

### DynamoDB Read/Write

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "dynamodb:GetItem",
        "dynamodb:PutItem",
        "dynamodb:UpdateItem",
        "dynamodb:Query",
        "dynamodb:Scan"
      ],
      "Resource": "arn:aws:dynamodb:us-east-1:123456789012:table/my-table"
    }
  ]
}
```

### Secrets Manager

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "secretsmanager:GetSecretValue"
      ],
      "Resource": "arn:aws:secretsmanager:us-east-1:123456789012:secret:my-secret-*"
    }
  ]
}
```

### S3 Read/Write

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:GetObject",
        "s3:PutObject"
      ],
      "Resource": "arn:aws:s3:::my-bucket/*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "s3:ListBucket"
      ],
      "Resource": "arn:aws:s3:::my-bucket"
    }
  ]
}
```

---

## ⚙️ Best Practices

### 1. Least Privilege

```hcl
# ❌ BAD: Too broad
Resource = "*"

# ✅ GOOD: Specific resources
Resource = [
  "arn:aws:sqs:us-east-1:123456789012:my-queue",
  "arn:aws:sqs:us-east-1:123456789012:my-queue-dlq"
]
```

### 2. Use Managed Policies for Common Patterns

```hcl
# ✅ GOOD: Use AWS managed policy
resource "aws_iam_role_policy_attachment" "lambda_basic" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}
```

### 3. Separate Policies by Service

```hcl
# ✅ GOOD: One policy per service
resource "aws_iam_policy" "sqs_access" { ... }
resource "aws_iam_policy" "dynamodb_access" { ... }
resource "aws_iam_policy" "secrets_access" { ... }
```

### 4. Tag Everything

```hcl
tags = {
  Name        = "${var.project_name}-role"
  Environment = "dev"
  ManagedBy   = "terraform"
}
```

### 5. Use Variables for ARNs

```hcl
variable "queue_arn" {
  description = "SQS queue ARN"
  type        = string
}

resource "aws_iam_policy" "sqs_access" {
  policy = jsonencode({
    Statement = [{
      Resource = var.queue_arn  # ✅ Parameterized
    }]
  })
}
```

### 6. Protect Sensitive Outputs

```hcl
output "deployer_secret_access_key" {
  value     = aws_iam_access_key.deployer.secret
  sensitive = true  # ✅ Hidden from console output
}
```

---

## 🧪 Testing

### Test Role Permissions

```bash
# Get role ARN
ROLE_ARN=$(terraform output -raw lambda_role_arn)

# Simulate policy (dry-run)
aws iam simulate-principal-policy \
  --policy-source-arn $ROLE_ARN \
  --action-names sqs:ReceiveMessage sqs:DeleteMessage
```

### Test User Permissions

```bash
# Get access keys
ACCESS_KEY=$(terraform output -raw deployer_access_key_id)
SECRET_KEY=$(terraform output -raw deployer_secret_access_key)

# Test CLI access
AWS_ACCESS_KEY_ID=$ACCESS_KEY \
AWS_SECRET_ACCESS_KEY=$SECRET_KEY \
aws sts get-caller-identity
```

### Check Attached Policies

```bash
# List policies attached to role
aws iam list-attached-role-policies \
  --role-name $(terraform output -raw lambda_role_name)

# List inline policies
aws iam list-role-policies \
  --role-name $(terraform output -raw lambda_role_name)
```

---

## ✅ Перевірка

### 1. Role Created

```bash
aws iam get-role --role-name sqs-processor-lambda-role
```

### 2. Policies Attached

```bash
aws iam list-attached-role-policies --role-name sqs-processor-lambda-role
```

### 3. User Created

```bash
aws iam get-user --user-name sqs-processor-deployer
```

### 4. Group Created

```bash
aws iam get-group --group-name sqs-processor-developers
```

---

## 🎓 Висновок

Ти навчився:
- ✅ Створювати IAM roles з Terraform
- ✅ Писати custom IAM policies
- ✅ Розуміти різницю між inline та managed policies
- ✅ Створювати IAM users з access keys
- ✅ Організовувати permissions через groups
- ✅ Застосовувати least privilege principle

**Далі:** `04_full_stack/` - Повний стек (Lambda + SQS + IAM)
