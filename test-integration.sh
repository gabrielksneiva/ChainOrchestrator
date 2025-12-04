#!/bin/bash
set -e

API_URL="https://lj2vv3cmj1.execute-api.us-east-1.amazonaws.com"
SNS_TOPIC="arn:aws:sns:us-east-1:490873503238:chainorchestrator-Transactions-production"
LAMBDA_NAME="chainorchestrator-production"

echo "🧪 ChainOrchestrator - Integration Test Suite"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

pass() {
  echo -e "${GREEN}✅ PASS${NC} - $1"
}

fail() {
  echo -e "${RED}❌ FAIL${NC} - $1"
  exit 1
}

info() {
  echo -e "${YELLOW}ℹ️  INFO${NC} - $1"
}

# ============================================================
# TEST 1: Health Check
# ============================================================
echo "📋 Test 1: API Gateway Health Check"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

HEALTH_RESPONSE=$(curl -s -w "\n%{http_code}" "${API_URL}/health")
HTTP_CODE=$(echo "$HEALTH_RESPONSE" | tail -n1)
BODY=$(echo "$HEALTH_RESPONSE" | head -n1)

if [ "$HTTP_CODE" -eq 200 ]; then
  if echo "$BODY" | grep -q "healthy"; then
    pass "Health endpoint returned 200 and 'healthy' status"
    info "Response: $BODY"
  else
    fail "Health endpoint returned 200 but unexpected body: $BODY"
  fi
else
  fail "Health endpoint returned HTTP $HTTP_CODE instead of 200"
fi

echo ""

# ============================================================
# TEST 2: POST Transaction (EVM)
# ============================================================
echo "📋 Test 2: POST /transaction (EVM blockchain)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

TRANSACTION_PAYLOAD='{
  "chain_type": "EVM",
  "operation_type": "transfer",
  "payload": {
    "network": "ethereum",
    "from_address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
    "to_address": "0x8626f6940E2eb28930eFb4CeF49B2d1F2C9C1199",
    "amount": "1.5",
    "currency": "ETH"
  }
}'

TX_RESPONSE=$(curl -s -w "\n%{http_code}" \
  -X POST "${API_URL}/transaction" \
  -H "Content-Type: application/json" \
  -d "$TRANSACTION_PAYLOAD")

TX_HTTP_CODE=$(echo "$TX_RESPONSE" | tail -n1)
TX_BODY=$(echo "$TX_RESPONSE" | head -n1)

if [ "$TX_HTTP_CODE" -eq 200 ] || [ "$TX_HTTP_CODE" -eq 202 ]; then
  if echo "$TX_BODY" | grep -q "operation_id"; then
    pass "Transaction submitted successfully"
    info "Response: $TX_BODY"
    OPERATION_ID=$(echo "$TX_BODY" | grep -o '"operation_id":"[^"]*"' | cut -d'"' -f4)
    info "Operation ID: $OPERATION_ID"
  else
    fail "Transaction returned success but no operation_id: $TX_BODY"
  fi
else
  fail "Transaction endpoint returned HTTP $TX_HTTP_CODE: $TX_BODY"
fi

echo ""

# ============================================================
# TEST 3: Check SQS Queue (EVM)
# ============================================================
echo "📋 Test 3: Verify message in EVM SQS Queue"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

EVM_QUEUE_URL=$(aws sqs list-queues --queue-name-prefix "chainorchestrator-evm-queue" --query 'QueueUrls[0]' --output text)

if [ -z "$EVM_QUEUE_URL" ]; then
  fail "Could not find EVM queue"
fi

info "Waiting 5 seconds for message propagation..."
sleep 5

# Try to receive message
SQS_MESSAGE=$(aws sqs receive-message \
  --queue-url "$EVM_QUEUE_URL" \
  --max-number-of-messages 1 \
  --wait-time-seconds 10 \
  --output json)

if echo "$SQS_MESSAGE" | grep -q "Messages"; then
  pass "Message received in EVM queue"
  
  # Extract and display message
  MESSAGE_BODY=$(echo "$SQS_MESSAGE" | jq -r '.Messages[0].Body' | jq -r '.Message')
  info "Message content:"
  echo "$MESSAGE_BODY" | jq .
  
  # Delete message to clean up
  RECEIPT_HANDLE=$(echo "$SQS_MESSAGE" | jq -r '.Messages[0].ReceiptHandle')
  aws sqs delete-message \
    --queue-url "$EVM_QUEUE_URL" \
    --receipt-handle "$RECEIPT_HANDLE" 2>/dev/null || true
  
  info "Message deleted from queue (cleanup)"
else
  fail "No message found in EVM queue after 10 seconds"
fi

echo ""

# ============================================================
# TEST 4: POST Transaction (Bitcoin)
# ============================================================
echo "📋 Test 4: POST /transaction (Bitcoin blockchain)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

BTC_PAYLOAD='{
  "chain_type": "BTC",
  "operation_type": "transfer",
  "payload": {
    "network": "mainnet",
    "from_address": "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
    "to_address": "3J98t1WpEZ73CNmYviecrnyiWrnqRhWNLy",
    "amount": "0.001",
    "currency": "BTC"
  }
}'

BTC_RESPONSE=$(curl -s -w "\n%{http_code}" \
  -X POST "${API_URL}/transaction" \
  -H "Content-Type: application/json" \
  -d "$BTC_PAYLOAD")

BTC_HTTP_CODE=$(echo "$BTC_RESPONSE" | tail -n1)
BTC_BODY=$(echo "$BTC_RESPONSE" | head -n1)

if [ "$BTC_HTTP_CODE" -eq 200 ] || [ "$BTC_HTTP_CODE" -eq 202 ]; then
  pass "Bitcoin transaction submitted successfully"
  info "Response: $BTC_BODY"
else
  fail "Bitcoin transaction returned HTTP $BTC_HTTP_CODE: $BTC_BODY"
fi

echo ""

# ============================================================
# TEST 5: Check Bitcoin SQS Queue
# ============================================================
echo "📋 Test 5: Verify message in Bitcoin SQS Queue"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

BTC_QUEUE_URL=$(aws sqs list-queues --queue-name-prefix "chainorchestrator-btc-queue" --query 'QueueUrls[0]' --output text)

info "Waiting 5 seconds for message propagation..."
sleep 5

BTC_SQS_MESSAGE=$(aws sqs receive-message \
  --queue-url "$BTC_QUEUE_URL" \
  --max-number-of-messages 1 \
  --wait-time-seconds 10 \
  --output json)

if echo "$BTC_SQS_MESSAGE" | grep -q "Messages"; then
  pass "Message received in Bitcoin queue"
  
  BTC_MESSAGE_BODY=$(echo "$BTC_SQS_MESSAGE" | jq -r '.Messages[0].Body' | jq -r '.Message')
  info "Message content:"
  echo "$BTC_MESSAGE_BODY" | jq .
  
  # Cleanup
  BTC_RECEIPT=$(echo "$BTC_SQS_MESSAGE" | jq -r '.Messages[0].ReceiptHandle')
  aws sqs delete-message \
    --queue-url "$BTC_QUEUE_URL" \
    --receipt-handle "$BTC_RECEIPT" 2>/dev/null || true
else
  fail "No message found in Bitcoin queue"
fi

echo ""

# ============================================================
# TEST 6: Invalid Transaction (Validation)
# ============================================================
echo "📋 Test 6: Invalid transaction (missing required fields)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

INVALID_PAYLOAD='{
  "chain_type": "EVM"
}'

INVALID_RESPONSE=$(curl -s -w "\n%{http_code}" \
  -X POST "${API_URL}/transaction" \
  -H "Content-Type: application/json" \
  -d "$INVALID_PAYLOAD")

INVALID_HTTP_CODE=$(echo "$INVALID_RESPONSE" | tail -n1)
INVALID_BODY=$(echo "$INVALID_RESPONSE" | head -n1)

if [ "$INVALID_HTTP_CODE" -eq 400 ]; then
  pass "Validation correctly rejected invalid transaction (HTTP 400)"
  info "Error: $INVALID_BODY"
else
  fail "Expected HTTP 400 for invalid transaction, got $INVALID_HTTP_CODE"
fi

echo ""

# ============================================================
# TEST 7: Lambda Direct Invocation
# ============================================================
echo "📋 Test 7: Lambda CloudWatch Logs"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

info "Checking recent Lambda logs..."
RECENT_LOGS=$(aws logs tail /aws/lambda/${LAMBDA_NAME} --since 5m --format short 2>/dev/null | head -20)

if [ -n "$RECENT_LOGS" ]; then
  pass "Lambda is logging correctly"
  info "Recent log entries:"
  echo "$RECENT_LOGS" | head -5
else
  fail "No recent logs found for Lambda"
fi

echo ""

# ============================================================
# TEST 8: Infrastructure Health
# ============================================================
echo "📋 Test 8: Infrastructure Components"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Check Lambda exists
if aws lambda get-function --function-name "$LAMBDA_NAME" &>/dev/null; then
  pass "Lambda function exists and is accessible"
else
  fail "Lambda function not found"
fi

# Check SNS Topic
if aws sns get-topic-attributes --topic-arn "$SNS_TOPIC" &>/dev/null; then
  pass "SNS topic exists and is accessible"
else
  fail "SNS topic not found"
fi

# Check SQS Queues
QUEUE_COUNT=$(aws sqs list-queues --queue-name-prefix "chainorchestrator" --query 'length(QueueUrls)' --output text)
if [ "$QUEUE_COUNT" -ge 4 ]; then
  pass "All SQS queues exist ($QUEUE_COUNT found)"
else
  fail "Expected at least 4 SQS queues, found $QUEUE_COUNT"
fi

# Check API Gateway
API_ID="lj2vv3cmj1"
if aws apigatewayv2 get-api --api-id "$API_ID" &>/dev/null; then
  pass "API Gateway exists and is accessible"
else
  fail "API Gateway not found"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}🎉 ALL TESTS PASSED!${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Summary:"
echo "  ✅ Health check working"
echo "  ✅ Transaction API working"
echo "  ✅ SNS publishing working"
echo "  ✅ SQS message routing working (EVM, Bitcoin)"
echo "  ✅ Validation working"
echo "  ✅ All infrastructure components healthy"
echo ""
echo "Your ChainOrchestrator is fully operational! 🚀"
