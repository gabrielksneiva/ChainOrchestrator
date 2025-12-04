# 🧪 Testing Guide

## Automated Integration Tests (CI/CD)

Toda vez que você faz push para `develop` ou `main`, a pipeline roda automaticamente **7 testes de integração**:

### Pipeline Flow:
```
Lint → Security → Unit Tests → Build → Deploy → Integration Tests → Summary
```

### Integration Tests:
1. ✅ **Health Check** - Verifica se API Gateway está respondendo
2. ✅ **EVM Transaction** - Envia transação Ethereum e valida resposta
3. ✅ **EVM Queue** - Verifica se mensagem chegou na fila SQS correta
4. ✅ **Bitcoin Transaction** - Envia transação Bitcoin
5. ✅ **Bitcoin Queue** - Verifica roteamento para fila Bitcoin
6. ✅ **Validation** - Testa rejeição de payloads inválidos (HTTP 400)
7. ✅ **Infrastructure** - Verifica Lambda, SNS, SQS, API Gateway

## Manual Testing

### Teste Rápido (Health Check)

```bash
curl https://lj2vv3cmj1.execute-api.us-east-1.amazonaws.com/health
```

Resposta esperada:
```json
{"status":"healthy","service":"ChainOrchestrator"}
```

### Teste Completo (Script)

Execute o script de integração completo:

```bash
./test-integration.sh
```

Esse script testa:
- ✅ Health endpoint
- ✅ POST /transaction (EVM + Bitcoin)
- ✅ SNS → SQS message routing
- ✅ Validation
- ✅ CloudWatch logs
- ✅ Infrastructure components

### Teste Manual de Transação

#### EVM (Ethereum)
```bash
curl -X POST https://lj2vv3cmj1.execute-api.us-east-1.amazonaws.com/transaction \
  -H "Content-Type: application/json" \
  -d '{
    "blockchain": "EVM",
    "network": "ethereum",
    "transaction_type": "transfer",
    "from_address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
    "to_address": "0x8626f6940E2eb28930eFb4CeF49B2d1F2C9C1199",
    "amount": "1.5",
    "currency": "ETH"
  }'
```

Resposta esperada:
```json
{
  "message": "transaction published successfully",
  "message_id": "a1b2c3d4-...",
  "blockchain": "EVM"
}
```

#### Bitcoin
```bash
curl -X POST https://lj2vv3cmj1.execute-api.us-east-1.amazonaws.com/transaction \
  -H "Content-Type: application/json" \
  -d '{
    "blockchain": "Bitcoin",
    "network": "mainnet",
    "transaction_type": "transfer",
    "from_address": "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
    "to_address": "3J98t1WpEZ73CNmYviecrnyiWrnqRhWNLy",
    "amount": "0.001",
    "currency": "BTC"
  }'
```

#### TRON
```bash
curl -X POST https://lj2vv3cmj1.execute-api.us-east-1.amazonaws.com/transaction \
  -H "Content-Type: application/json" \
  -d '{
    "blockchain": "TRON",
    "network": "mainnet",
    "transaction_type": "transfer",
    "from_address": "TXYZopYRdj2D9XRtbG411XZZ3kM5VkAeBf",
    "to_address": "TLPkjzj4u8ej7fHRqKVQGTqCBqxR6b5ZxY",
    "amount": "100",
    "currency": "TRX"
  }'
```

#### Solana
```bash
curl -X POST https://lj2vv3cmj1.execute-api.us-east-1.amazonaws.com/transaction \
  -H "Content-Type: application/json" \
  -d '{
    "blockchain": "Solana",
    "network": "mainnet",
    "transaction_type": "transfer",
    "from_address": "7v91N7iZ9mNicL8WfG6cgSCKyRXydQjLh6UYBWwm6y1Q",
    "to_address": "HN7cABqLq46Es1jh92dQQisAq662SmxELLLsHHe4YWrH",
    "amount": "1.0",
    "currency": "SOL"
  }'
```

### Verificar Mensagens nas Filas SQS

```bash
# EVM Queue
aws sqs receive-message \
  --queue-url $(aws sqs list-queues --queue-name-prefix chainorchestrator-evm-queue --query 'QueueUrls[0]' --output text) \
  --max-number-of-messages 1

# Bitcoin Queue
aws sqs receive-message \
  --queue-url $(aws sqs list-queues --queue-name-prefix chainorchestrator-btc-queue --query 'QueueUrls[0]' --output text) \
  --max-number-of-messages 1

# TRON Queue
aws sqs receive-message \
  --queue-url $(aws sqs list-queues --queue-name-prefix chainorchestrator-tron-queue --query 'QueueUrls[0]' --output text) \
  --max-number-of-messages 1

# Solana Queue
aws sqs receive-message \
  --queue-url $(aws sqs list-queues --queue-name-prefix chainorchestrator-sol-queue --query 'QueueUrls[0]' --output text) \
  --max-number-of-messages 1
```

### Verificar Logs

```bash
# Lambda logs (últimos 5 minutos)
aws logs tail /aws/lambda/chainorchestrator-production --since 5m --follow

# API Gateway logs
aws logs tail /aws/apigateway/chainorchestrator-production --since 5m --follow
```

## Teste de Validação (Erro Esperado)

```bash
# Payload inválido (sem campos obrigatórios)
curl -X POST https://lj2vv3cmj1.execute-api.us-east-1.amazonaws.com/transaction \
  -H "Content-Type: application/json" \
  -d '{"blockchain": "EVM"}'
```

Resposta esperada (HTTP 400):
```json
{
  "error": "validation failed",
  "details": ["field 'network' is required", ...]
}
```

## GitHub Actions - Ver Resultados

1. Acesse: https://github.com/gabrielksneiva/ChainOrchestrator/actions
2. Clique no último workflow run
3. Veja o job "Integration Tests"
4. Todos os 7 testes devem estar ✅ verdes

## Métricas de Sucesso

- ✅ Health check retorna HTTP 200
- ✅ Transactions retornam HTTP 200/202 com `message_id`
- ✅ Mensagens aparecem nas filas SQS corretas (filtradas por blockchain)
- ✅ Validação rejeita payloads inválidos com HTTP 400
- ✅ Todos os componentes AWS acessíveis (Lambda, SNS, SQS, API Gateway)

## Troubleshooting

### API retorna 502/503
```bash
# Verificar status do Lambda
aws lambda get-function --function-name chainorchestrator-production

# Ver logs de erro
aws logs tail /aws/lambda/chainorchestrator-production --since 10m
```

### Mensagens não chegam na fila
```bash
# Verificar SNS topic
aws sns get-topic-attributes \
  --topic-arn arn:aws:sns:us-east-1:490873503238:chainorchestrator-Transactions-production

# Verificar subscriptions
aws sns list-subscriptions-by-topic \
  --topic-arn arn:aws:sns:us-east-1:490873503238:chainorchestrator-Transactions-production
```

### Pipeline falhando
1. Veja qual job falhou no Actions
2. Leia os logs do step específico
3. Se for integration test, pode ser timeout ou permissões AWS
4. Verifique se AWS credentials estão configuradas no GitHub Secrets
