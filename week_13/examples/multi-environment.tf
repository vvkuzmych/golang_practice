# Example: Multi-Environment Setup

# Use workspaces or separate state files for different environments

# Option 1: Terraform Workspaces
# Commands:
# terraform workspace new dev
# terraform workspace new staging
# terraform workspace new prod
# terraform workspace select dev
# terraform workspace list

# Option 2: Separate directories (recommended)
# Directory structure:
# environments/
#   dev/
#     main.tf
#     terraform.tfvars
#   staging/
#     main.tf
#     terraform.tfvars
#   prod/
#     main.tf
#     terraform.tfvars

# Shared variables
variable "environment" {
  description = "Environment name"
  type        = string
  validation {
    condition     = contains(["dev", "staging", "prod"], var.environment)
    error_message = "Environment must be dev, staging, or prod."
  }
}

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

# Environment-specific config
locals {
  config = {
    dev = {
      lambda_memory = 512
      lambda_timeout = 30
      retention_days = 7
    }
    staging = {
      lambda_memory = 1024
      lambda_timeout = 60
      retention_days = 14
    }
    prod = {
      lambda_memory = 2048
      lambda_timeout = 300
      retention_days = 30
    }
  }

  environment_config = local.config[var.environment]
}

# Lambda with environment-specific settings
resource "aws_lambda_function" "processor" {
  function_name = "${var.environment}-processor"
  role          = aws_iam_role.lambda.arn
  runtime       = "provided.al2023"
  
  memory_size = local.environment_config.lambda_memory
  timeout     = local.environment_config.lambda_timeout
  
  environment {
    variables = {
      ENVIRONMENT = var.environment
    }
  }

  tags = {
    Environment = var.environment
    ManagedBy   = "terraform"
  }
}

# CloudWatch with environment-specific retention
resource "aws_cloudwatch_log_group" "lambda" {
  name              = "/aws/lambda/${var.environment}-processor"
  retention_in_days = local.environment_config.retention_days

  tags = {
    Environment = var.environment
    ManagedBy   = "terraform"
  }
}

# IAM role (same for all environments)
resource "aws_iam_role" "lambda" {
  name = "${var.environment}-lambda-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })

  tags = {
    Environment = var.environment
    ManagedBy   = "terraform"
  }
}

# Usage:
# 
# dev/terraform.tfvars:
# environment = "dev"
#
# staging/terraform.tfvars:
# environment = "staging"
#
# prod/terraform.tfvars:
# environment = "prod"
#
# Deploy:
# cd environments/dev && terraform apply
# cd environments/staging && terraform apply
# cd environments/prod && terraform apply
