# Configure Terraform and AWS provider
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
  region = var.aws_region
}

# Variables
variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "project_name" {
  description = "Project name"
  type        = string
  default     = "sqs-processor"
}

# ========================================
# IAM Role for Lambda
# ========================================

resource "aws_iam_role" "lambda" {
  name = "${var.project_name}-lambda-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })

  tags = {
    Name        = "${var.project_name}-lambda-role"
    Environment = "dev"
    ManagedBy   = "terraform"
  }
}

# ========================================
# IAM Policies
# ========================================

# Policy 1: CloudWatch Logs (Basic Lambda Execution)
resource "aws_iam_role_policy_attachment" "lambda_basic" {
  role       = aws_iam_role.lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Policy 2: Custom SQS Access
resource "aws_iam_policy" "sqs_access" {
  name        = "${var.project_name}-sqs-policy"
  description = "Allow Lambda to consume from SQS"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "sqs:ReceiveMessage",
          "sqs:DeleteMessage",
          "sqs:GetQueueAttributes",
          "sqs:ChangeMessageVisibility"
        ]
        Resource = "*"  # In production, specify exact queue ARN
      }
    ]
  })

  tags = {
    Name        = "${var.project_name}-sqs-policy"
    Environment = "dev"
    ManagedBy   = "terraform"
  }
}

resource "aws_iam_role_policy_attachment" "lambda_sqs" {
  role       = aws_iam_role.lambda.name
  policy_arn = aws_iam_policy.sqs_access.arn
}

# Policy 3: DynamoDB Access (for idempotency table)
resource "aws_iam_policy" "dynamodb_access" {
  name        = "${var.project_name}-dynamodb-policy"
  description = "Allow Lambda to access DynamoDB"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:Query"
        ]
        Resource = "*"  # In production, specify exact table ARN
      }
    ]
  })

  tags = {
    Name        = "${var.project_name}-dynamodb-policy"
    Environment = "dev"
    ManagedBy   = "terraform"
  }
}

resource "aws_iam_role_policy_attachment" "lambda_dynamodb" {
  role       = aws_iam_role.lambda.name
  policy_arn = aws_iam_policy.dynamodb_access.arn
}

# Policy 4: Secrets Manager (for sensitive data)
resource "aws_iam_policy" "secrets_access" {
  name        = "${var.project_name}-secrets-policy"
  description = "Allow Lambda to read secrets"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "secretsmanager:GetSecretValue"
        ]
        Resource = "*"  # In production, specify exact secret ARN
      }
    ]
  })

  tags = {
    Name        = "${var.project_name}-secrets-policy"
    Environment = "dev"
    ManagedBy   = "terraform"
  }
}

resource "aws_iam_role_policy_attachment" "lambda_secrets" {
  role       = aws_iam_role.lambda.name
  policy_arn = aws_iam_policy.secrets_access.arn
}

# ========================================
# Inline Policy Example
# ========================================

resource "aws_iam_role_policy" "lambda_inline" {
  name = "${var.project_name}-inline-policy"
  role = aws_iam_role.lambda.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject"
        ]
        Resource = "arn:aws:s3:::my-bucket/*"
      }
    ]
  })
}

# ========================================
# IAM User (for CI/CD or local testing)
# ========================================

resource "aws_iam_user" "deployer" {
  name = "${var.project_name}-deployer"

  tags = {
    Name        = "${var.project_name}-deployer"
    Environment = "dev"
    ManagedBy   = "terraform"
  }
}

# Access keys for the user
resource "aws_iam_access_key" "deployer" {
  user = aws_iam_user.deployer.name
}

# Policy for deployer user
resource "aws_iam_user_policy" "deployer" {
  name = "${var.project_name}-deployer-policy"
  user = aws_iam_user.deployer.name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "lambda:UpdateFunctionCode",
          "lambda:UpdateFunctionConfiguration",
          "lambda:PublishVersion"
        ]
        Resource = "*"
      },
      {
        Effect = "Allow"
        Action = [
          "sqs:SendMessage",
          "sqs:GetQueueUrl"
        ]
        Resource = "*"
      }
    ]
  })
}

# ========================================
# IAM Group (for team access)
# ========================================

resource "aws_iam_group" "developers" {
  name = "${var.project_name}-developers"
}

resource "aws_iam_group_policy" "developers" {
  name  = "${var.project_name}-developers-policy"
  group = aws_iam_group.developers.name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "lambda:InvokeFunction",
          "lambda:GetFunction",
          "logs:GetLogEvents",
          "logs:FilterLogEvents"
        ]
        Resource = "*"
      }
    ]
  })
}

# Add user to group
resource "aws_iam_user_group_membership" "deployer" {
  user = aws_iam_user.deployer.name

  groups = [
    aws_iam_group.developers.name
  ]
}

# ========================================
# Outputs
# ========================================

output "lambda_role_arn" {
  description = "Lambda IAM role ARN"
  value       = aws_iam_role.lambda.arn
}

output "lambda_role_name" {
  description = "Lambda IAM role name"
  value       = aws_iam_role.lambda.name
}

output "deployer_user_name" {
  description = "Deployer IAM user name"
  value       = aws_iam_user.deployer.name
}

output "deployer_access_key_id" {
  description = "Deployer access key ID"
  value       = aws_iam_access_key.deployer.id
}

output "deployer_secret_access_key" {
  description = "Deployer secret access key"
  value       = aws_iam_access_key.deployer.secret
  sensitive   = true
}

output "developers_group_name" {
  description = "Developers IAM group name"
  value       = aws_iam_group.developers.name
}

output "policies_created" {
  description = "Custom IAM policies created"
  value = [
    aws_iam_policy.sqs_access.name,
    aws_iam_policy.dynamodb_access.name,
    aws_iam_policy.secrets_access.name
  ]
}
