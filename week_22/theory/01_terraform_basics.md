# Terraform Basics

## Що таке Terraform?

**Terraform** - це Infrastructure as Code (IaC) tool від HashiCorp для створення, зміни та версіонування інфраструктури.

### Чому Terraform?

✅ **Declarative** - описуєш що хочеш, а не як це зробити  
✅ **Multi-cloud** - AWS, GCP, Azure, тощо  
✅ **Version Control** - інфраструктура в Git  
✅ **Idempotent** - можна запускати багато разів безпечно  
✅ **Plan before apply** - preview змін перед застосуванням  

---

## Основні концепції

### 1. Provider
**Інтеграція з cloud platform.**

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
```

### 2. Resource
**Інфраструктурний об'єкт (EC2, VPC, S3, etc.).**

```hcl
resource "aws_instance" "web" {
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t2.micro"
  
  tags = {
    Name = "WebServer"
  }
}
```

### 3. Data Source
**Читання існуючих даних.**

```hcl
data "aws_ami" "ubuntu" {
  most_recent = true
  
  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }
  
  owners = ["099720109477"]  # Canonical
}

resource "aws_instance" "web" {
  ami = data.aws_ami.ubuntu.id
  # ...
}
```

### 4. Variable
**Параметризація конфігурації.**

```hcl
# variables.tf
variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t2.micro"
}

variable "region" {
  type = string
}

variable "tags" {
  type = map(string)
  default = {
    Environment = "dev"
  }
}

# main.tf
resource "aws_instance" "web" {
  instance_type = var.instance_type
  tags          = var.tags
}
```

### 5. Output
**Експорт значень після apply.**

```hcl
output "instance_ip" {
  description = "Public IP of EC2 instance"
  value       = aws_instance.web.public_ip
}

output "instance_id" {
  value = aws_instance.web.id
}
```

---

## Terraform Workflow

```
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│   init   │ → │   plan   │ → │  apply   │ → │ destroy  │
└──────────┘    └──────────┘    └──────────┘    └──────────┘
```

### 1. terraform init
**Ініціалізація проекту (завантаження providers).**

```bash
terraform init

# Output:
# Initializing the backend...
# Initializing provider plugins...
# - Installing hashicorp/aws v5.0.0...
# Terraform has been successfully initialized!
```

### 2. terraform plan
**Preview змін (не застосовує).**

```bash
terraform plan

# Output:
# Terraform will perform the following actions:
# 
#   # aws_instance.web will be created
#   + resource "aws_instance" "web" {
#       + ami           = "ami-0c55b159cbfafe1f0"
#       + instance_type = "t2.micro"
#       ...
#     }
# 
# Plan: 1 to add, 0 to change, 0 to destroy.
```

### 3. terraform apply
**Застосувати зміни.**

```bash
terraform apply

# Preview + confirmation
# Do you want to perform these actions? yes

# Output:
# aws_instance.web: Creating...
# aws_instance.web: Creation complete after 30s [id=i-1234567890abcdef0]
# Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

### 4. terraform destroy
**Видалити всю інфраструктуру.**

```bash
terraform destroy

# Warning: will destroy all resources!
# Do you really want to destroy all resources? yes
```

---

## State Management

### terraform.tfstate
**Файл зі станом інфраструктури.**

```json
{
  "version": 4,
  "terraform_version": "1.5.0",
  "resources": [
    {
      "type": "aws_instance",
      "name": "web",
      "instances": [
        {
          "attributes": {
            "id": "i-1234567890abcdef0",
            "public_ip": "54.123.45.67"
          }
        }
      ]
    }
  ]
}
```

**ВАЖЛИВО:**
- ❌ НЕ комітити `.tfstate` в Git!
- ✅ Використовуй remote backend (S3, Terraform Cloud)
- ✅ State може містити секрети

### Remote Backend

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

---

## Практичний приклад

### main.tf

```hcl
terraform {
  required_version = ">= 1.0"
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.region
}

# VPC
resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  
  tags = {
    Name = "${var.project_name}-vpc"
  }
}

# Subnet
resource "aws_subnet" "public" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = true
  
  tags = {
    Name = "${var.project_name}-public-subnet"
  }
}

# Security Group
resource "aws_security_group" "web" {
  name        = "${var.project_name}-web-sg"
  description = "Allow HTTP and SSH"
  vpc_id      = aws_vpc.main.id
  
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# EC2 Instance
resource "aws_instance" "web" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = var.instance_type
  subnet_id     = aws_subnet.public.id
  
  vpc_security_group_ids = [aws_security_group.web.id]
  
  user_data = <<-EOF
    #!/bin/bash
    apt-get update
    apt-get install -y nginx
    echo "Hello from Terraform!" > /var/www/html/index.html
    systemctl start nginx
  EOF
  
  tags = {
    Name = "${var.project_name}-web-server"
  }
}

# Data source for AMI
data "aws_ami" "ubuntu" {
  most_recent = true
  
  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }
  
  owners = ["099720109477"]
}
```

### variables.tf

```hcl
variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "project_name" {
  description = "Project name for tagging"
  type        = string
  default     = "myapp"
}

variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t2.micro"
}
```

### outputs.tf

```hcl
output "instance_public_ip" {
  description = "Public IP address of EC2 instance"
  value       = aws_instance.web.public_ip
}

output "instance_id" {
  description = "ID of EC2 instance"
  value       = aws_instance.web.id
}

output "vpc_id" {
  description = "ID of VPC"
  value       = aws_vpc.main.id
}
```

### terraform.tfvars

```hcl
region        = "eu-west-1"
project_name  = "my-web-app"
instance_type = "t3.small"
```

### Використання

```bash
# Initialize
terraform init

# Format code
terraform fmt

# Validate
terraform validate

# Plan
terraform plan -out=tfplan

# Apply
terraform apply tfplan

# Output
terraform output instance_public_ip

# Show state
terraform show

# List resources
terraform state list

# Destroy
terraform destroy
```

---

## Корисні команди

```bash
# Format all .tf files
terraform fmt -recursive

# Validate configuration
terraform validate

# Show plan without asking for input
terraform plan -input=false

# Apply without confirmation (CI/CD)
terraform apply -auto-approve

# Destroy specific resource
terraform destroy -target=aws_instance.web

# Import existing resource
terraform import aws_instance.web i-1234567890abcdef0

# Refresh state
terraform refresh

# Show outputs
terraform output

# Unlock state (if locked)
terraform force-unlock <lock-id>

# Graph dependencies
terraform graph | dot -Tsvg > graph.svg
```

---

## Best Practices

### 1. Структура проекту

```
terraform/
├── environments/
│   ├── dev/
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── terraform.tfvars
│   │   └── backend.tf
│   ├── staging/
│   └── prod/
├── modules/
│   ├── vpc/
│   ├── ec2/
│   └── rds/
└── .gitignore
```

### 2. .gitignore

```gitignore
# Terraform
.terraform/
*.tfstate
*.tfstate.*
*.tfvars
.terraform.lock.hcl
```

### 3. Використовуй modules

```hcl
module "vpc" {
  source = "./modules/vpc"
  
  cidr_block   = "10.0.0.0/16"
  project_name = var.project_name
}

module "ec2" {
  source = "./modules/ec2"
  
  vpc_id       = module.vpc.vpc_id
  subnet_id    = module.vpc.public_subnet_id
  project_name = var.project_name
}
```

### 4. Версіонуй providers

```hcl
terraform {
  required_version = ">= 1.0"
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"  # >= 5.0, < 6.0
    }
  }
}
```

### 5. Використовуй workspaces

```bash
# Create workspace
terraform workspace new dev
terraform workspace new prod

# List workspaces
terraform workspace list

# Switch workspace
terraform workspace select prod

# Current workspace
terraform workspace show
```

---

## Підсумок

| Команда | Опис |
|---------|------|
| `terraform init` | Ініціалізація |
| `terraform plan` | Preview змін |
| `terraform apply` | Застосувати зміни |
| `terraform destroy` | Видалити все |
| `terraform fmt` | Форматувати код |
| `terraform validate` | Перевірити синтаксис |
| `terraform output` | Показати outputs |
| `terraform state list` | Список resources |

**Переваги Terraform:**
- ✅ Infrastructure as Code
- ✅ Multi-cloud
- ✅ Version control
- ✅ Preview changes
- ✅ Modules for reusability

**Недоліки:**
- ❌ Крива навчання
- ❌ State management complexity
- ❌ Потрібна обережність (може видалити prod!)
