# Example: Using Data Sources

# Data sources allow you to query existing resources
# without creating or managing them

# Get current AWS account ID
data "aws_caller_identity" "current" {}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "caller_arn" {
  value = data.aws_caller_identity.current.arn
}

# Get current AWS region
data "aws_region" "current" {}

output "region" {
  value = data.aws_region.current.name
}

# Get existing VPC by tag
data "aws_vpc" "main" {
  tags = {
    Name = "main-vpc"
  }
}

output "vpc_id" {
  value = data.aws_vpc.main.id
}

# Get existing subnet by filter
data "aws_subnet" "private" {
  filter {
    name   = "tag:Name"
    values = ["private-subnet-1"]
  }
}

output "subnet_id" {
  value = data.aws_subnet.private.id
}

# Get existing SQS queue
data "aws_sqs_queue" "existing" {
  name = "existing-queue"
}

output "queue_arn" {
  value = data.aws_sqs_queue.existing.arn
}

# Get existing Lambda function
data "aws_lambda_function" "existing" {
  function_name = "existing-function"
}

output "lambda_arn" {
  value = data.aws_lambda_function.existing.arn
}

# Get existing IAM role
data "aws_iam_role" "existing" {
  name = "existing-role"
}

output "role_arn" {
  value = data.aws_iam_role.existing.arn
}

# Get latest Amazon Linux 2 AMI
data "aws_ami" "amazon_linux_2" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-*-x86_64-gp2"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

output "ami_id" {
  value = data.aws_ami.amazon_linux_2.id
}

# Get availability zones
data "aws_availability_zones" "available" {
  state = "available"
}

output "availability_zones" {
  value = data.aws_availability_zones.available.names
}

# Get secret from Secrets Manager
data "aws_secretsmanager_secret" "db_password" {
  name = "db-password"
}

data "aws_secretsmanager_secret_version" "db_password" {
  secret_id = data.aws_secretsmanager_secret.db_password.id
}

# Use secret in resource (not visible in output)
resource "aws_db_instance" "example" {
  # ... other config ...
  password = data.aws_secretsmanager_secret_version.db_password.secret_string
}

# Get SSM parameter
data "aws_ssm_parameter" "ami_id" {
  name = "/aws/service/ami-amazon-linux-latest/amzn2-ami-hvm-x86_64-gp2"
}

output "latest_ami" {
  value = data.aws_ssm_parameter.ami_id.value
}

# Example: Reference existing queue in new Lambda
data "aws_sqs_queue" "orders" {
  name = "orders-queue"
}

resource "aws_lambda_event_source_mapping" "sqs" {
  event_source_arn = data.aws_sqs_queue.orders.arn
  function_name    = aws_lambda_function.processor.arn
  batch_size       = 10
}

# Use data source to build resource names
locals {
  account_id = data.aws_caller_identity.current.account_id
  region     = data.aws_region.current.name
  
  bucket_name = "my-app-${local.account_id}-${local.region}"
}

resource "aws_s3_bucket" "app" {
  bucket = local.bucket_name
}

# Common use cases:
# 1. Get existing resources to reference
# 2. Query AWS-managed resources (AMIs, AZs)
# 3. Read secrets/parameters
# 4. Build dynamic names based on account/region
# 5. Import existing infrastructure
