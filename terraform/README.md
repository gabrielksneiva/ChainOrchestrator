# ChainOrchestrator - Terraform Infrastructure

Infrastructure as Code for ChainOrchestrator using Terraform.

## 📦 Components Provisioned

- **AWS Lambda**: Go-based serverless function (arm64, provided.al2023)
- **API Gateway**: HTTP API with CORS and routes
- **SNS Topic**: Blockchain event distribution
- **SQS Queues**: Per-blockchain message queues (EVM, Bitcoin, TRON, Solana) with DLQs
- **IAM Roles**: Lambda execution role with proper permissions
- **CloudWatch**: Logs and monitoring

## 🏗️ Architecture

```
API Gateway → Lambda → SNS Topic → SQS Queues (filtered by blockchain)
                                    ├─ EVM Queue
                                    ├─ Bitcoin Queue
                                    ├─ TRON Queue
                                    └─ Solana Queue
```

## 🚀 Initial Setup

### 1. Create Terraform Backend (One-Time)

Run the setup script to create S3 bucket and DynamoDB table for state management:

```bash
./.github/setup-terraform-backend.sh
```

This creates:
- S3 bucket: `chainorchestrator-terraform-state` (versioned, encrypted)
- DynamoDB table: `terraform-state-lock` (for state locking)

### 2. Initialize Terraform

```bash
cd terraform
terraform init
```

### 3. Review Configuration

Check `terraform.tfvars` and modify if needed:

```hcl
aws_region         = "us-east-1"
environment        = "production"
lambda_timeout     = 60
lambda_memory      = 512
log_retention_days = 7
```

### 4. Plan and Apply

```bash
# Preview changes
terraform plan

# Apply changes
terraform apply
```

## 🔄 CI/CD Integration

The GitHub Actions workflow automatically:
1. Builds Lambda binary (Go arm64)
2. Creates `lambda.zip` package
3. Runs `terraform plan`
4. Applies infrastructure changes on push to `main`/`develop`
5. Updates Lambda function code
6. Runs health checks

## 📋 Outputs

After deployment, Terraform provides:

```bash
terraform output
```

- `api_endpoint`: API Gateway URL
- `lambda_function_name`: Lambda function name
- `lambda_function_arn`: Lambda ARN
- `sns_topic_arn`: SNS topic ARN
- `sqs_queues`: Map of SQS queue URLs

## 🔧 Manual Operations

### Update Lambda Code Only

```bash
# Build binary
GOOS=linux GOARCH=arm64 go build -o bootstrap cmd/lambda/main.go
zip lambda.zip bootstrap

# Update via AWS CLI
aws lambda update-function-code \
  --function-name chainorchestrator-production \
  --zip-file fileb://lambda.zip
```

### View Infrastructure State

```bash
terraform show
```

### Import Existing Resources (if needed)

```bash
# Example: import existing Lambda
terraform import aws_lambda_function.orchestrator chainorchestrator-production
```

## 🧪 Testing

### Test Lambda Directly

```bash
aws lambda invoke \
  --function-name $(terraform output -raw lambda_function_name) \
  --payload '{"httpMethod":"GET","path":"/health"}' \
  response.json && cat response.json
```

### Test via API Gateway

```bash
API_URL=$(terraform output -raw api_endpoint)
curl $API_URL/health
```

### Publish Test Message to SNS

```bash
SNS_TOPIC=$(terraform output -raw sns_topic_arn)

aws sns publish \
  --topic-arn "$SNS_TOPIC" \
  --message '{"transaction":"test"}' \
  --message-attributes '{"blockchain":{"DataType":"String","StringValue":"EVM"}}'
```

## 📂 File Structure

```
terraform/
├── main.tf           # Provider and backend config
├── variables.tf      # Input variables
├── terraform.tfvars  # Variable values
├── lambda.tf         # Lambda function
├── api_gateway.tf    # API Gateway
├── sns.tf            # SNS + SQS resources
├── iam.tf            # IAM roles and policies
└── outputs.tf        # Output values
```

## 🔐 State Management

State is stored in:
- **S3**: `s3://chainorchestrator-terraform-state/orchestrator/terraform.tfstate`
- **Lock**: DynamoDB table `terraform-state-lock`

**Important**: Never commit `terraform.tfstate` to git!

## 🧹 Cleanup

To destroy all infrastructure:

```bash
terraform destroy
```

**Warning**: This will delete:
- Lambda function
- API Gateway
- SNS topic
- All SQS queues (including messages)
- CloudWatch logs

## 📊 Cost Estimation

Approximate monthly costs (us-east-1):

- Lambda: $0.00 - $5 (depends on invocations, 512MB/60s)
- API Gateway: $1.00 per million requests
- SNS: $0.50 per million publishes
- SQS: $0.40 per million requests (first 1M free)
- CloudWatch Logs: $0.50/GB ingested
- S3 (state): ~$0.02/month

**Estimated total**: $2-10/month for low-medium traffic

## 🆘 Troubleshooting

### Terraform State Locked

```bash
# View lock info
aws dynamodb get-item \
  --table-name terraform-state-lock \
  --key '{"LockID":{"S":"chainorchestrator-terraform-state/orchestrator/terraform.tfstate"}}'

# Force unlock (use carefully!)
terraform force-unlock <LOCK_ID>
```

### Lambda Update Conflicts

If Lambda is being updated via both Terraform and CI/CD:

1. CI/CD updates **code only** (zip file)
2. Terraform manages **configuration** (memory, timeout, env vars)

Both can coexist safely.

### Import Manual Resources

If resources were created manually:

```bash
# Lambda
terraform import aws_lambda_function.orchestrator chainorchestrator-production

# SNS
terraform import aws_sns_topic.transactions arn:aws:sns:us-east-1:ACCOUNT:topic-name

# SQS
terraform import aws_sqs_queue.evm_queue https://sqs.us-east-1.amazonaws.com/ACCOUNT/queue-name
```

## 📚 References

- [Terraform AWS Provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [AWS Lambda](https://docs.aws.amazon.com/lambda/)
- [API Gateway HTTP APIs](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api.html)
- [SNS Message Filtering](https://docs.aws.amazon.com/sns/latest/dg/sns-message-filtering.html)
