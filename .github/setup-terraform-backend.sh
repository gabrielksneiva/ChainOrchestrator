#!/bin/bash
set -e

BUCKET_NAME="chainorchestrator-terraform-state"
TABLE_NAME="terraform-state-lock"
REGION="us-east-1"

echo "🚀 Setting up Terraform Backend Infrastructure"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Check if bucket exists
if aws s3 ls "s3://$BUCKET_NAME" 2>/dev/null; then
  echo "✅ S3 bucket already exists: $BUCKET_NAME"
else
  echo "📦 Creating S3 bucket: $BUCKET_NAME"
  aws s3api create-bucket \
    --bucket "$BUCKET_NAME" \
    --region "$REGION"
  
  # Enable versioning
  aws s3api put-bucket-versioning \
    --bucket "$BUCKET_NAME" \
    --versioning-configuration Status=Enabled
  
  # Enable encryption
  aws s3api put-bucket-encryption \
    --bucket "$BUCKET_NAME" \
    --server-side-encryption-configuration '{
      "Rules": [{
        "ApplyServerSideEncryptionByDefault": {
          "SSEAlgorithm": "AES256"
        }
      }]
    }'
  
  # Block public access
  aws s3api put-public-access-block \
    --bucket "$BUCKET_NAME" \
    --public-access-block-configuration \
      BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true
  
  echo "✅ S3 bucket created and configured"
fi

# Check if DynamoDB table exists
if aws dynamodb describe-table --table-name "$TABLE_NAME" --region "$REGION" 2>/dev/null; then
  echo "✅ DynamoDB table already exists: $TABLE_NAME"
else
  echo "🔐 Creating DynamoDB table: $TABLE_NAME"
  aws dynamodb create-table \
    --table-name "$TABLE_NAME" \
    --attribute-definitions AttributeName=LockID,AttributeType=S \
    --key-schema AttributeName=LockID,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --region "$REGION"
  
  echo "⏳ Waiting for table to be active..."
  aws dynamodb wait table-exists \
    --table-name "$TABLE_NAME" \
    --region "$REGION"
  
  echo "✅ DynamoDB table created"
fi

echo ""
echo "✅ Terraform backend infrastructure ready!"
echo ""
echo "Backend configuration:"
echo "  S3 Bucket:      $BUCKET_NAME"
echo "  DynamoDB Table: $TABLE_NAME"
echo "  Region:         $REGION"
echo ""
echo "Next steps:"
echo "  1. cd terraform"
echo "  2. terraform init"
echo "  3. terraform plan"
echo "  4. terraform apply"
