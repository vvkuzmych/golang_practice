# Configure Terraform and AWS provider
terraform {
  required_version = ">= 1.0"
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
  
  # Uncomment for production
  # backend "s3" {
  #   bucket         = "my-terraform-state"
  #   key            = "full-stack/terraform.tfstate"
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

variable "project_name" {
  description = "Project name"
  type        = string
  default     = "order-processor"
}

variable "environment" {
  description = "Environment (dev, staging, prod)"
  type        = string
  default     = "dev"
}

# ========================================
# SQS Queues
# ========================================

# Dead Letter Queue
resource "aws_sqs_queue" "dlq" {
  name                      = "${var.project_name}-dlq"
  message_retention_seconds = 1209600  # 14 days

  tags = {
    Name        = "${var.project_name}-dlq"
    Environment = var.environment
    ManagedBy   = "terraform"
  }
}

# Main Queue
resource "aws_sqs_queue" "main" {
  name                       = "${var.project_name}-queue"
  visibility_timeout_seconds = 300
  message_retention_seconds  = 345600  # 4 days
  receive_wait_time_seconds  = 20       # Long polling

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.dlq.arn
    maxReceiveCount     = 3
  })

  tags = {
    Name        = "${var.project_name}-queue"
    Environment = var.environment
    ManagedBy   = "terraform"
  }
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
    Environment = var.environment
    ManagedBy   = "terraform"
  }
}

# Basic Lambda execution (CloudWatch Logs)
resource "aws_iam_role_policy_attachment" "lambda_basic" {
  role       = aws_iam_role.lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# SQS policy for Lambda
resource "aws_iam_policy" "lambda_sqs" {
  name        = "${var.project_name}-lambda-sqs-policy"
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
        Resource = [
          aws_sqs_queue.main.arn,
          aws_sqs_queue.dlq.arn
        ]
      }
    ]
  })

  tags = {
    Name        = "${var.project_name}-lambda-sqs-policy"
    Environment = var.environment
    ManagedBy   = "terraform"
  }
}

resource "aws_iam_role_policy_attachment" "lambda_sqs" {
  role       = aws_iam_role.lambda.name
  policy_arn = aws_iam_policy.lambda_sqs.arn
}

# ========================================
# CloudWatch Log Group
# ========================================

resource "aws_cloudwatch_log_group" "lambda" {
  name              = "/aws/lambda/${var.project_name}"
  retention_in_days = 7

  tags = {
    Name        = "${var.project_name}-logs"
    Environment = var.environment
    ManagedBy   = "terraform"
  }
}

# ========================================
# Lambda Function
# ========================================

resource "aws_lambda_function" "processor" {
  filename      = "function.zip"
  function_name = var.project_name
  role          = aws_iam_role.lambda.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"

  source_code_hash = filebase64sha256("function.zip")

  memory_size = 1024
  timeout     = 60

  environment {
    variables = {
      ENVIRONMENT = var.environment
      QUEUE_URL   = aws_sqs_queue.main.url
      DLQ_URL     = aws_sqs_queue.dlq.url
    }
  }

  tags = {
    Name        = var.project_name
    Environment = var.environment
    ManagedBy   = "terraform"
  }

  depends_on = [
    aws_cloudwatch_log_group.lambda,
    aws_iam_role_policy_attachment.lambda_basic,
    aws_iam_role_policy_attachment.lambda_sqs
  ]
}

# ========================================
# Lambda Event Source Mapping (SQS -> Lambda)
# ========================================

resource "aws_lambda_event_source_mapping" "sqs" {
  event_source_arn = aws_sqs_queue.main.arn
  function_name    = aws_lambda_function.processor.arn

  batch_size                         = 10
  maximum_batching_window_in_seconds = 5

  # Partial batch failure handling
  function_response_types = ["ReportBatchItemFailures"]

  enabled = true
}

# ========================================
# CloudWatch Alarms
# ========================================

# Alarm: DLQ not empty
resource "aws_cloudwatch_metric_alarm" "dlq_not_empty" {
  alarm_name          = "${var.project_name}-dlq-not-empty"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 1
  metric_name         = "ApproximateNumberOfMessagesVisible"
  namespace           = "AWS/SQS"
  period              = 60
  statistic           = "Average"
  threshold           = 0
  alarm_description   = "DLQ has poison messages"

  dimensions = {
    QueueName = aws_sqs_queue.dlq.name
  }

  tags = {
    Name        = "${var.project_name}-dlq-alarm"
    Environment = var.environment
    ManagedBy   = "terraform"
  }
}

# Alarm: Lambda errors
resource "aws_cloudwatch_metric_alarm" "lambda_errors" {
  alarm_name          = "${var.project_name}-lambda-errors"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 1
  metric_name         = "Errors"
  namespace           = "AWS/Lambda"
  period              = 300
  statistic           = "Sum"
  threshold           = 5
  alarm_description   = "Lambda function errors"

  dimensions = {
    FunctionName = aws_lambda_function.processor.function_name
  }

  tags = {
    Name        = "${var.project_name}-lambda-errors"
    Environment = var.environment
    ManagedBy   = "terraform"
  }
}

# ========================================
# Outputs
# ========================================

output "queue_url" {
  description = "Main queue URL"
  value       = aws_sqs_queue.main.url
}

output "queue_arn" {
  description = "Main queue ARN"
  value       = aws_sqs_queue.main.arn
}

output "dlq_url" {
  description = "DLQ URL"
  value       = aws_sqs_queue.dlq.url
}

output "lambda_arn" {
  description = "Lambda function ARN"
  value       = aws_lambda_function.processor.arn
}

output "lambda_name" {
  description = "Lambda function name"
  value       = aws_lambda_function.processor.function_name
}

output "lambda_role_arn" {
  description = "Lambda IAM role ARN"
  value       = aws_iam_role.lambda.arn
}

output "log_group_name" {
  description = "CloudWatch log group"
  value       = aws_cloudwatch_log_group.lambda.name
}
