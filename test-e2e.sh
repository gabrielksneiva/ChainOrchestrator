#!/bin/bash
# Don't use set -e, we want to run all tests even if some fail
# Exit code will be determined by test results at the end

echo "🧪 Starting ChainOrchestrator E2E Tests..."
echo "=========================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
API_URL="${API_GATEWAY_URL:-https://z78d7j5h9h.execute-api.us-east-1.amazonaws.com}"
SNS_TOPIC_ARN="${SNS_TOPIC_ARN}"
AWS_REGION="${AWS_REGION:-us-east-1}"

# Test counters
TESTS_PASSED=0
TESTS_FAILED=0

# Helper functions
pass() {
    echo -e "${GREEN}✓${NC} $1"
    ((TESTS_PASSED++))
}

fail() {
    echo -e "${RED}✗${NC} $1"
    ((TESTS_FAILED++))
}

info() {
    echo -e "${YELLOW}ℹ${NC} $1"
}

# Test 1: Health Check
echo ""
info "Test 1: API Gateway Health Check"
HEALTH_RESPONSE=$(curl -s -w "\n%{http_code}" "${API_URL}/health")
HTTP_CODE=$(echo "$HEALTH_RESPONSE" | tail -n1)
BODY=$(echo "$HEALTH_RESPONSE" | head -n-1)

if [ "$HTTP_CODE" = "200" ]; then
    if echo "$BODY" | grep -q "healthy"; then
        pass "Health endpoint returned 200 and healthy status"
    else
        fail "Health endpoint returned 200 but unexpected body: $BODY"
    fi
else
    fail "Health endpoint returned HTTP $HTTP_CODE"
fi

# Test 2: Lambda Function Exists
echo ""
info "Test 2: Lambda Function Status"
LAMBDA_STATUS=$(aws lambda get-function \
    --function-name chainorchestrator \
    --region $AWS_REGION \
    --query 'Configuration.State' \
    --output text 2>/dev/null || echo "NOT_FOUND")

if [ "$LAMBDA_STATUS" = "Active" ]; then
    pass "Lambda function is Active"
else
    fail "Lambda function status: $LAMBDA_STATUS"
fi

# Test 3: SNS Topic Exists
echo ""
info "Test 3: SNS Topic Configuration"
SNS_EXISTS=$(aws sns get-topic-attributes \
    --topic-arn "arn:aws:sns:${AWS_REGION}:$(aws sts get-caller-identity --query Account --output text):chainorchestrator-Transactions" \
    --region $AWS_REGION \
    --query 'Attributes.TopicArn' \
    --output text 2>/dev/null || echo "NOT_FOUND")

if [ "$SNS_EXISTS" != "NOT_FOUND" ]; then
    pass "SNS Topic exists and is accessible"
    SNS_TOPIC_ARN="$SNS_EXISTS"
else
    fail "SNS Topic not found or not accessible"
fi

# Test 4: SQS Queues Exist
echo ""
info "Test 4: SQS Queues Configuration"
QUEUES=("evm-queue" "tron-queue" "btc-queue" "sol-queue")
QUEUES_OK=0

for queue in "${QUEUES[@]}"; do
    QUEUE_URL=$(aws sqs get-queue-url \
        --queue-name "$queue" \
        --region $AWS_REGION \
        --query 'QueueUrl' \
        --output text 2>/dev/null || echo "NOT_FOUND")
    
    if [ "$QUEUE_URL" != "NOT_FOUND" ]; then
        ((QUEUES_OK++))
    fi
done

if [ $QUEUES_OK -eq 4 ]; then
    pass "All 4 blockchain queues exist"
else
    fail "Only $QUEUES_OK/4 queues found"
fi

# Test 5: SNS Subscriptions
echo ""
info "Test 5: SNS to SQS Subscriptions"
if [ "$SNS_TOPIC_ARN" != "NOT_FOUND" ]; then
    SUBSCRIPTION_COUNT=$(aws sns list-subscriptions-by-topic \
        --topic-arn "$SNS_TOPIC_ARN" \
        --region $AWS_REGION \
        --query 'Subscriptions | length(@)' \
        --output text 2>/dev/null || echo "0")
    
    if [ "$SUBSCRIPTION_COUNT" -ge 4 ]; then
        pass "SNS has $SUBSCRIPTION_COUNT subscriptions (expected 4+)"
    else
        fail "SNS has only $SUBSCRIPTION_COUNT subscriptions"
    fi
else
    fail "Cannot test subscriptions - SNS topic not found"
fi

# Test 6: Post Transaction (EVM)
echo ""
info "Test 6: POST /transaction endpoint (EVM chain)"
TRANSACTION_PAYLOAD=$(cat <<EOF
{
  "operation_id": "test-e2e-$(date +%s)",
  "chain_type": "EVM",
  "operation_type": "TRANSFER",
  "from_address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
  "to_address": "0x8626f6940E2eb28930eFb4CeF49B2d1F2C9C1199",
  "amount": "0.001",
  "chain_id": "11155111",
  "network": "sepolia"
}
EOF
)

POST_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_URL}/transaction" \
    -H "Content-Type: application/json" \
    -d "$TRANSACTION_PAYLOAD")

POST_HTTP_CODE=$(echo "$POST_RESPONSE" | tail -n1)
POST_BODY=$(echo "$POST_RESPONSE" | head -n-1)

if [ "$POST_HTTP_CODE" = "200" ] || [ "$POST_HTTP_CODE" = "201" ] || [ "$POST_HTTP_CODE" = "202" ]; then
    if echo "$POST_BODY" | grep -q "operation_id"; then
        pass "POST transaction accepted with HTTP $POST_HTTP_CODE"
        OPERATION_ID=$(echo "$POST_BODY" | grep -o '"operation_id":"[^"]*"' | cut -d'"' -f4)
    else
        fail "POST transaction returned HTTP $POST_HTTP_CODE but unexpected body"
    fi
else
    fail "POST transaction failed with HTTP $POST_HTTP_CODE: $POST_BODY"
fi

# Test 7: Check message in SQS
echo ""
info "Test 7: Message routing to EVM queue"
sleep 2  # Wait for message to arrive
EVM_QUEUE_URL=$(aws sqs get-queue-url \
    --queue-name "evm-queue" \
    --region $AWS_REGION \
    --query 'QueueUrl' \
    --output text 2>/dev/null)

if [ -n "$EVM_QUEUE_URL" ] && [ "$EVM_QUEUE_URL" != "NOT_FOUND" ]; then
    MESSAGES=$(aws sqs get-queue-attributes \
        --queue-url "$EVM_QUEUE_URL" \
        --attribute-names ApproximateNumberOfMessages \
        --region $AWS_REGION \
        --query 'Attributes.ApproximateNumberOfMessages' \
        --output text 2>/dev/null || echo "0")
    
    if [ "$MESSAGES" -gt 0 ]; then
        pass "Message found in EVM queue (count: $MESSAGES)"
    else
        info "No messages in EVM queue yet (may have been processed)"
    fi
else
    fail "Could not access EVM queue"
fi

# Test 8: CloudWatch Logs
echo ""
info "Test 8: Lambda CloudWatch Logs"
LOG_STREAMS=$(aws logs describe-log-streams \
    --log-group-name "/aws/lambda/chainorchestrator" \
    --region $AWS_REGION \
    --order-by LastEventTime \
    --descending \
    --max-items 1 \
    --query 'logStreams[0].logStreamName' \
    --output text 2>/dev/null || echo "NOT_FOUND")

if [ "$LOG_STREAMS" != "NOT_FOUND" ] && [ -n "$LOG_STREAMS" ]; then
    pass "CloudWatch logs are being generated"
else
    fail "No CloudWatch log streams found"
fi

# Test 9: IAM Role
echo ""
info "Test 9: Lambda IAM Role Configuration"
ROLE_ARN=$(aws lambda get-function \
    --function-name chainorchestrator \
    --region $AWS_REGION \
    --query 'Configuration.Role' \
    --output text 2>/dev/null || echo "NOT_FOUND")

if echo "$ROLE_ARN" | grep -q "chainorchestrator-lambda-execution-role"; then
    pass "Lambda has correct IAM role"
else
    fail "Lambda IAM role unexpected: $ROLE_ARN"
fi

# Test 10: API Gateway Integration
echo ""
info "Test 10: API Gateway Lambda Integration"
API_ID=$(aws apigatewayv2 get-apis \
    --region $AWS_REGION \
    --query 'Items[?Name==`chainorchestrator-api`].ApiId' \
    --output text 2>/dev/null)

if [ -n "$API_ID" ]; then
    INTEGRATIONS=$(aws apigatewayv2 get-integrations \
        --api-id "$API_ID" \
        --region $AWS_REGION \
        --query 'Items | length(@)' \
        --output text 2>/dev/null || echo "0")
    
    if [ "$INTEGRATIONS" -gt 0 ]; then
        pass "API Gateway has Lambda integration"
    else
        fail "No integrations found for API Gateway"
    fi
else
    fail "API Gateway not found"
fi

# Summary
echo ""
echo "=========================================="
echo "Test Summary:"
echo -e "${GREEN}Passed: $TESTS_PASSED${NC}"
echo -e "${RED}Failed: $TESTS_FAILED${NC}"
echo "=========================================="

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ All E2E tests passed!${NC}"
    exit 0
else
    echo -e "${RED}✗ Some tests failed${NC}"
    exit 1
fi
