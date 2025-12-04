# SNS Topic for Transactions
resource "aws_sns_topic" "transactions" {
  name              = "${var.project_name}-${var.sns_topic_name}-${var.environment}"
  display_name      = "Blockchain Transactions Topic"
  fifo_topic        = false
  
  tags = {
    Name = "${var.project_name}-transactions"
  }
}

# SNS Topic Policy
resource "aws_sns_topic_policy" "transactions" {
  arn = aws_sns_topic.transactions.arn

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "AllowPublishFromLambda"
        Effect = "Allow"
        Principal = {
          AWS = aws_iam_role.lambda_execution_role.arn
        }
        Action   = "SNS:Publish"
        Resource = aws_sns_topic.transactions.arn
      }
    ]
  })
}

# SNS Topic Subscription Filter Policy para EVM
resource "aws_sns_topic_subscription" "evm_subscription" {
  topic_arn = aws_sns_topic.transactions.arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.evm_queue.arn

  filter_policy = jsonencode({
    chain_type = ["EVM"]
  })
}

# SNS Topic Subscription Filter Policy para TRON
resource "aws_sns_topic_subscription" "tron_subscription" {
  topic_arn = aws_sns_topic.transactions.arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.tron_queue.arn

  filter_policy = jsonencode({
    chain_type = ["TRON"]
  })
}

# SNS Topic Subscription Filter Policy para BTC
resource "aws_sns_topic_subscription" "btc_subscription" {
  topic_arn = aws_sns_topic.transactions.arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.btc_queue.arn

  filter_policy = jsonencode({
    chain_type = ["BTC"]
  })
}

# SNS Topic Subscription Filter Policy para SOL
resource "aws_sns_topic_subscription" "sol_subscription" {
  topic_arn = aws_sns_topic.transactions.arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.sol_queue.arn

  filter_policy = jsonencode({
    chain_type = ["SOL"]
  })
}

# SQS Queues para cada blockchain
resource "aws_sqs_queue" "evm_queue" {
  name                      = "${var.project_name}-evm-queue-${var.environment}"
  delay_seconds             = 0
  max_message_size          = 262144
  message_retention_seconds = 1209600 # 14 dias
  receive_wait_time_seconds = 10
  visibility_timeout_seconds = 300

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.evm_dlq.arn
    maxReceiveCount     = 3
  })

  tags = {
    Name = "${var.project_name}-evm-queue"
    ChainType = "EVM"
  }
}

resource "aws_sqs_queue" "evm_dlq" {
  name                      = "${var.project_name}-evm-dlq-${var.environment}"
  message_retention_seconds = 1209600 # 14 dias

  tags = {
    Name = "${var.project_name}-evm-dlq"
  }
}

resource "aws_sqs_queue" "tron_queue" {
  name                      = "${var.project_name}-tron-queue-${var.environment}"
  delay_seconds             = 0
  max_message_size          = 262144
  message_retention_seconds = 1209600
  receive_wait_time_seconds = 10
  visibility_timeout_seconds = 300

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.tron_dlq.arn
    maxReceiveCount     = 3
  })

  tags = {
    Name = "${var.project_name}-tron-queue"
    ChainType = "TRON"
  }
}

resource "aws_sqs_queue" "tron_dlq" {
  name                      = "${var.project_name}-tron-dlq-${var.environment}"
  message_retention_seconds = 1209600

  tags = {
    Name = "${var.project_name}-tron-dlq"
  }
}

resource "aws_sqs_queue" "btc_queue" {
  name                      = "${var.project_name}-btc-queue-${var.environment}"
  delay_seconds             = 0
  max_message_size          = 262144
  message_retention_seconds = 1209600
  receive_wait_time_seconds = 10
  visibility_timeout_seconds = 300

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.btc_dlq.arn
    maxReceiveCount     = 3
  })

  tags = {
    Name = "${var.project_name}-btc-queue"
    ChainType = "BTC"
  }
}

resource "aws_sqs_queue" "btc_dlq" {
  name                      = "${var.project_name}-btc-dlq-${var.environment}"
  message_retention_seconds = 1209600

  tags = {
    Name = "${var.project_name}-btc-dlq"
  }
}

resource "aws_sqs_queue" "sol_queue" {
  name                      = "${var.project_name}-sol-queue-${var.environment}"
  delay_seconds             = 0
  max_message_size          = 262144
  message_retention_seconds = 1209600
  receive_wait_time_seconds = 10
  visibility_timeout_seconds = 300

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.sol_dlq.arn
    maxReceiveCount     = 3
  })

  tags = {
    Name = "${var.project_name}-sol-queue"
    ChainType = "SOL"
  }
}

resource "aws_sqs_queue" "sol_dlq" {
  name                      = "${var.project_name}-sol-dlq-${var.environment}"
  message_retention_seconds = 1209600

  tags = {
    Name = "${var.project_name}-sol-dlq"
  }
}

# SQS Queue Policies
resource "aws_sqs_queue_policy" "evm_queue_policy" {
  queue_url = aws_sqs_queue.evm_queue.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "sns.amazonaws.com"
        }
        Action   = "SQS:SendMessage"
        Resource = aws_sqs_queue.evm_queue.arn
        Condition = {
          ArnEquals = {
            "aws:SourceArn" = aws_sns_topic.transactions.arn
          }
        }
      }
    ]
  })
}

resource "aws_sqs_queue_policy" "tron_queue_policy" {
  queue_url = aws_sqs_queue.tron_queue.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "sns.amazonaws.com"
        }
        Action   = "SQS:SendMessage"
        Resource = aws_sqs_queue.tron_queue.arn
        Condition = {
          ArnEquals = {
            "aws:SourceArn" = aws_sns_topic.transactions.arn
          }
        }
      }
    ]
  })
}

resource "aws_sqs_queue_policy" "btc_queue_policy" {
  queue_url = aws_sqs_queue.btc_queue.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "sns.amazonaws.com"
        }
        Action   = "SQS:SendMessage"
        Resource = aws_sqs_queue.btc_queue.arn
        Condition = {
          ArnEquals = {
            "aws:SourceArn" = aws_sns_topic.transactions.arn
          }
        }
      }
    ]
  })
}

resource "aws_sqs_queue_policy" "sol_queue_policy" {
  queue_url = aws_sqs_queue.sol_queue.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "sns.amazonaws.com"
        }
        Action   = "SQS:SendMessage"
        Resource = aws_sqs_queue.sol_queue.arn
        Condition = {
          ArnEquals = {
            "aws:SourceArn" = aws_sns_topic.transactions.arn
          }
        }
      }
    ]
  })
}
