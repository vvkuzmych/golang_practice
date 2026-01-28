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

variable "queue_name" {
  description = "SQS queue name"
  type        = string
  default     = "my-queue"
}

variable "visibility_timeout" {
  description = "Visibility timeout in seconds"
  type        = number
  default     = 300
}

variable "message_retention_period" {
  description = "Message retention in seconds (4 days)"
  type        = number
  default     = 345600
}

variable "max_receive_count" {
  description = "Max receives before DLQ"
  type        = number
  default     = 3
}

# Dead Letter Queue
resource "aws_sqs_queue" "dlq" {
  name = "${var.queue_name}-dlq"
  
  # Longer retention for DLQ (14 days)
  message_retention_seconds = 1209600
  
  tags = {
    Name        = "${var.queue_name}-dlq"
    Environment = "dev"
    ManagedBy   = "terraform"
  }
}

# Main Queue
resource "aws_sqs_queue" "main" {
  name = var.queue_name
  
  # Visibility timeout (how long message is hidden after receive)
  visibility_timeout_seconds = var.visibility_timeout
  
  # Message retention (how long message stays in queue)
  message_retention_seconds = var.message_retention_period
  
  # Long polling (reduce costs, faster delivery)
  receive_wait_time_seconds = 20
  
  # Dead Letter Queue configuration
  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.dlq.arn
    maxReceiveCount     = var.max_receive_count
  })
  
  tags = {
    Name        = var.queue_name
    Environment = "dev"
    ManagedBy   = "terraform"
  }
}

# CloudWatch Alarms for monitoring

# Alarm: DLQ not empty
resource "aws_cloudwatch_metric_alarm" "dlq_not_empty" {
  alarm_name          = "${var.queue_name}-dlq-not-empty"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 1
  metric_name         = "ApproximateNumberOfMessagesVisible"
  namespace           = "AWS/SQS"
  period              = 60
  statistic           = "Average"
  threshold           = 0
  alarm_description   = "Alert when DLQ has messages"
  
  dimensions = {
    QueueName = aws_sqs_queue.dlq.name
  }
  
  tags = {
    Name        = "${var.queue_name}-dlq-alarm"
    Environment = "dev"
    ManagedBy   = "terraform"
  }
}

# Alarm: High queue depth
resource "aws_cloudwatch_metric_alarm" "high_queue_depth" {
  alarm_name          = "${var.queue_name}-high-depth"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 2
  metric_name         = "ApproximateNumberOfMessagesVisible"
  namespace           = "AWS/SQS"
  period              = 300
  statistic           = "Average"
  threshold           = 1000
  alarm_description   = "Alert when queue has too many messages"
  
  dimensions = {
    QueueName = aws_sqs_queue.main.name
  }
  
  tags = {
    Name        = "${var.queue_name}-depth-alarm"
    Environment = "dev"
    ManagedBy   = "terraform"
  }
}

# Alarm: Old messages
resource "aws_cloudwatch_metric_alarm" "old_messages" {
  alarm_name          = "${var.queue_name}-old-messages"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 1
  metric_name         = "ApproximateAgeOfOldestMessage"
  namespace           = "AWS/SQS"
  period              = 300
  statistic           = "Maximum"
  threshold           = 600  # 10 minutes
  alarm_description   = "Alert when messages are too old"
  
  dimensions = {
    QueueName = aws_sqs_queue.main.name
  }
  
  tags = {
    Name        = "${var.queue_name}-age-alarm"
    Environment = "dev"
    ManagedBy   = "terraform"
  }
}

# Outputs
output "queue_url" {
  description = "Main queue URL"
  value       = aws_sqs_queue.main.url
}

output "queue_arn" {
  description = "Main queue ARN"
  value       = aws_sqs_queue.main.arn
}

output "queue_name" {
  description = "Main queue name"
  value       = aws_sqs_queue.main.name
}

output "dlq_url" {
  description = "DLQ URL"
  value       = aws_sqs_queue.dlq.url
}

output "dlq_arn" {
  description = "DLQ ARN"
  value       = aws_sqs_queue.dlq.arn
}

output "dlq_name" {
  description = "DLQ name"
  value       = aws_sqs_queue.dlq.name
}
