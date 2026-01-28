# Terraform Examples

Корисні Terraform patterns та приклади.

---

## 📚 Файли

### 1. [backend-s3.tf](./backend-s3.tf)

**Remote State Backend з S3 та DynamoDB**

**Що створює:**
- S3 bucket для state файлів
- Versioning (для recovery)
- Encryption (AES256)
- Public access block
- DynamoDB table для locking

**Використання:**
```bash
# 1. Deploy backend infrastructure
terraform apply

# 2. Get outputs
terraform output s3_bucket_name
terraform output dynamodb_table_name

# 3. Configure backend in main.tf
terraform {
  backend "s3" {
    bucket         = "OUTPUT_FROM_STEP_2"
    key            = "prod/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "terraform-locks"
  }
}

# 4. Initialize with new backend
terraform init -reconfigure
```

---

### 2. [multi-environment.tf](./multi-environment.tf)

**Multi-Environment Setup (dev/staging/prod)**

**Підхід 1: Workspaces**
```bash
terraform workspace new dev
terraform workspace new staging
terraform workspace new prod

terraform workspace select dev
terraform apply
```

**Підхід 2: Separate Directories (рекомендується)**
```
environments/
  dev/
    main.tf
    terraform.tfvars
  staging/
    main.tf
    terraform.tfvars
  prod/
    main.tf
    terraform.tfvars
```

**Environment-Specific Config:**
```hcl
locals {
  config = {
    dev = {
      lambda_memory = 512
      retention_days = 7
    }
    prod = {
      lambda_memory = 2048
      retention_days = 30
    }
  }
}
```

---

### 3. [data-sources.tf](./data-sources.tf)

**Data Sources - Query Existing Resources**

**Приклади:**

**Get AWS account info:**
```hcl
data "aws_caller_identity" "current" {}
data "aws_region" "current" {}
```

**Get existing resources:**
```hcl
data "aws_vpc" "main" {
  tags = { Name = "main-vpc" }
}

data "aws_sqs_queue" "existing" {
  name = "existing-queue"
}

data "aws_lambda_function" "existing" {
  function_name = "existing-function"
}
```

**Get secrets:**
```hcl
data "aws_secretsmanager_secret_version" "db_password" {
  secret_id = "db-password"
}

resource "aws_db_instance" "example" {
  password = data.aws_secretsmanager_secret_version.db_password.secret_string
}
```

**Get latest AMI:**
```hcl
data "aws_ami" "amazon_linux_2" {
  most_recent = true
  owners      = ["amazon"]
  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-*"]
  }
}
```

---

## 🎯 Common Patterns

### Pattern 1: Remote State Reference

```hcl
# In project A (creates VPC)
output "vpc_id" {
  value = aws_vpc.main.id
}

# In project B (uses VPC)
data "terraform_remote_state" "network" {
  backend = "s3"
  config = {
    bucket = "terraform-state"
    key    = "network/terraform.tfstate"
    region = "us-east-1"
  }
}

resource "aws_subnet" "app" {
  vpc_id = data.terraform_remote_state.network.outputs.vpc_id
}
```

### Pattern 2: Count for Multiple Resources

```hcl
variable "subnet_count" {
  default = 3
}

resource "aws_subnet" "private" {
  count             = var.subnet_count
  vpc_id            = aws_vpc.main.id
  cidr_block        = cidrsubnet(aws_vpc.main.cidr_block, 8, count.index)
  availability_zone = data.aws_availability_zones.available.names[count.index]
  
  tags = {
    Name = "private-subnet-${count.index + 1}"
  }
}
```

### Pattern 3: For Each for Map Resources

```hcl
variable "queues" {
  type = map(object({
    visibility_timeout = number
    retention_seconds  = number
  }))
  default = {
    orders = {
      visibility_timeout = 300
      retention_seconds  = 345600
    }
    notifications = {
      visibility_timeout = 60
      retention_seconds  = 86400
    }
  }
}

resource "aws_sqs_queue" "queues" {
  for_each = var.queues
  
  name                       = each.key
  visibility_timeout_seconds = each.value.visibility_timeout
  message_retention_seconds  = each.value.retention_seconds
}
```

### Pattern 4: Conditional Resources

```hcl
variable "enable_monitoring" {
  type    = bool
  default = true
}

resource "aws_cloudwatch_metric_alarm" "lambda_errors" {
  count = var.enable_monitoring ? 1 : 0
  
  alarm_name = "lambda-errors"
  # ... other config
}
```

### Pattern 5: Dynamic Blocks

```hcl
variable "ingress_rules" {
  type = list(object({
    from_port   = number
    to_port     = number
    protocol    = string
    cidr_blocks = list(string)
  }))
  default = [
    {
      from_port   = 80
      to_port     = 80
      protocol    = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
    },
    {
      from_port   = 443
      to_port     = 443
      protocol    = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
    }
  ]
}

resource "aws_security_group" "web" {
  name = "web-sg"
  
  dynamic "ingress" {
    for_each = var.ingress_rules
    content {
      from_port   = ingress.value.from_port
      to_port     = ingress.value.to_port
      protocol    = ingress.value.protocol
      cidr_blocks = ingress.value.cidr_blocks
    }
  }
}
```

---

## 🔧 Best Practices

### 1. Use Variables

```hcl
# ❌ BAD
resource "aws_lambda_function" "processor" {
  function_name = "my-processor"
  memory_size   = 1024
}

# ✅ GOOD
variable "function_name" { type = string }
variable "memory_size" { type = number, default = 1024 }

resource "aws_lambda_function" "processor" {
  function_name = var.function_name
  memory_size   = var.memory_size
}
```

### 2. Use Locals for Computed Values

```hcl
locals {
  common_tags = {
    Environment = var.environment
    ManagedBy   = "terraform"
    Project     = var.project_name
  }
  
  resource_prefix = "${var.environment}-${var.project_name}"
}

resource "aws_lambda_function" "processor" {
  function_name = "${local.resource_prefix}-processor"
  tags          = local.common_tags
}
```

### 3. Use Data Sources to Reference Existing Resources

```hcl
# ✅ GOOD - Reference existing VPC
data "aws_vpc" "main" {
  tags = { Name = "main" }
}

resource "aws_subnet" "app" {
  vpc_id = data.aws_vpc.main.id
}
```

### 4. Separate State per Environment

```
s3://bucket/dev/terraform.tfstate
s3://bucket/staging/terraform.tfstate
s3://bucket/prod/terraform.tfstate
```

### 5. Use Remote State

```hcl
# ✅ Always use remote backend for teams
terraform {
  backend "s3" {
    bucket         = "terraform-state"
    key            = "prod/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "terraform-locks"
  }
}
```

---

## 📖 Resources

- [Terraform Documentation](https://www.terraform.io/docs)
- [Terraform AWS Provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [Terraform Best Practices](https://www.terraform-best-practices.com/)

---

**Week 13: Infrastructure as Code!** 🏗️
