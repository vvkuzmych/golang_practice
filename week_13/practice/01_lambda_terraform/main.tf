# Configure Terraform and AWS provider
terraform {
  required_version = ">= 1.0"
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
  
  # Remote backend (uncomment for production)
  # backend "s3" {
  #   bucket         = "my-terraform-state"
  #   key            = "lambda/terraform.tfstate"
  #   region         = "us-east-1"
  #   encrypt        = true
  #   dynamodb_table = "terraform-locks"
  # }
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

variable "function_name" {
  description = "Lambda function name"
  type        = string
  default     = "my-processor"
}

variable "memory_size" {
  description = "Lambda memory in MB"
  type        = number
  default     = 1024
}

variable "timeout" {
  description = "Lambda timeout in seconds"
  type        = number
  default     = 30
}

# IAM Role for Lambda
resource "aws_iam_role" "lambda" {
  name = "${var.function_name}-role"

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
    Name        = "${var.function_name}-role"
    Environment = "dev"
    ManagedBy   = "terraform"
  }
}

# Attach basic Lambda execution policy
resource "aws_iam_role_policy_attachment" "lambda_basic" {
  role       = aws_iam_role.lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# CloudWatch Log Group
resource "aws_cloudwatch_log_group" "lambda" {
  name              = "/aws/lambda/${var.function_name}"
  retention_in_days = 7

  tags = {
    Name        = "${var.function_name}-logs"
    Environment = "dev"
    ManagedBy   = "terraform"
  }
}

# Lambda function
resource "aws_lambda_function" "processor" {
  filename      = "function.zip"  # You need to create this
  function_name = var.function_name
  role          = aws_iam_role.lambda.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  
  source_code_hash = filebase64sha256("function.zip")
  
  memory_size = var.memory_size
  timeout     = var.timeout
  
  environment {
    variables = {
      STAGE       = "dev"
      LOG_LEVEL   = "INFO"
    }
  }

  tags = {
    Name        = var.function_name
    Environment = "dev"
    ManagedBy   = "terraform"
  }
  
  # Ensure log group exists before Lambda
  depends_on = [
    aws_cloudwatch_log_group.lambda,
    aws_iam_role_policy_attachment.lambda_basic
  ]
}

# Lambda function URL (optional - for HTTP access)
resource "aws_lambda_function_url" "processor" {
  function_name      = aws_lambda_function.processor.function_name
  authorization_type = "NONE"  # Change to "AWS_IAM" for production
  
  cors {
    allow_origins = ["*"]
    allow_methods = ["POST", "GET"]
    max_age       = 86400
  }
}

# Outputs
output "lambda_arn" {
  description = "Lambda function ARN"
  value       = aws_lambda_function.processor.arn
}

output "lambda_name" {
  description = "Lambda function name"
  value       = aws_lambda_function.processor.function_name
}

output "lambda_url" {
  description = "Lambda function URL"
  value       = aws_lambda_function_url.processor.function_url
}

output "lambda_role_arn" {
  description = "Lambda IAM role ARN"
  value       = aws_iam_role.lambda.arn
}

output "log_group_name" {
  description = "CloudWatch log group name"
  value       = aws_cloudwatch_log_group.lambda.name
}
