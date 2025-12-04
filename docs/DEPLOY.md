# Deploy ChainOrchestrator to AWS Lambda

Este guia descreve como fazer o deploy do ChainOrchestrator usando Terraform.

## Pré-requisitos

1. **AWS CLI configurado** com credenciais válidas
2. **Terraform** instalado (>= 1.0)
3. **Go** 1.24+ instalado
4. **Make** instalado

## Deploy Inicial

### 1. Build do binário Lambda

```bash
make build
```

Isso irá criar o binário `bin/bootstrap` otimizado para AWS Lambda (arm64).

### 2. Empacotar Lambda

```bash
make lambda-zip
```

Isso cria o arquivo `lambda.zip` pronto para upload.

### 3. Configurar Backend do Terraform (primeira vez)

Crie o bucket S3 e tabela DynamoDB para o state do Terraform:

```bash
# Criar bucket S3
aws s3 mb s3://chainorchestrator-terraform-state --region us-east-1

# Habilitar versionamento
aws s3api put-bucket-versioning \
  --bucket chainorchestrator-terraform-state \
  --versioning-configuration Status=Enabled

# Habilitar encriptação
aws s3api put-bucket-encryption \
  --bucket chainorchestrator-terraform-state \
  --server-side-encryption-configuration '{
    "Rules": [{
      "ApplyServerSideEncryptionByDefault": {
        "SSEAlgorithm": "AES256"
      }
    }]
  }'

# Criar tabela DynamoDB para lock
aws dynamodb create-table \
  --table-name terraform-state-lock \
  --attribute-definitions AttributeName=LockID,AttributeType=S \
  --key-schema AttributeName=LockID,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --region us-east-1
```

### 4. Inicializar Terraform

```bash
cd terraform
terraform init
```

### 5. Review do Plano

```bash
terraform plan
```

### 6. Deploy

```bash
terraform apply
```

Ou use o Makefile:

```bash
make deploy
```

## Outputs Importantes

Após o deploy, o Terraform mostrará:

- **api_gateway_url**: URL pública da API
- **lambda_function_name**: Nome da função Lambda
- **sns_topic_arn**: ARN do SNS Topic "Transactions"

Exemplo:
```
Outputs:

api_gateway_url = "https://abc123xyz.execute-api.us-east-1.amazonaws.com"
lambda_function_name = "chainorchestrator-production"
sns_topic_arn = "arn:aws:sns:us-east-1:123456789012:chainorchestrator-Transactions-production"
```

## Testar o Deploy

```bash
# Health check
curl https://abc123xyz.execute-api.us-east-1.amazonaws.com/health

# POST transaction
curl -X POST https://abc123xyz.execute-api.us-east-1.amazonaws.com/transaction \
  -H "Content-Type: application/json" \
  -d '{
    "chain_type": "EVM",
    "operation_type": "TRANSFER",
    "payload": {
      "from": "0x123...",
      "to": "0x456...",
      "amount": "1000000000000000000"
    }
  }'
```

## Atualizar Lambda

Após fazer mudanças no código:

```bash
# Build + package + deploy
make build
make lambda-zip
cd terraform && terraform apply
```

Ou simplesmente:

```bash
make deploy
```

## Monitoramento

### CloudWatch Logs

```bash
# Logs da Lambda
aws logs tail /aws/lambda/chainorchestrator-production --follow

# Logs do API Gateway
aws logs tail /aws/apigateway/chainorchestrator-production --follow
```

### Métricas

Acesse o CloudWatch Console para ver:
- Lambda invocations
- Lambda errors
- Lambda duration
- API Gateway requests
- API Gateway 4XX/5XX errors

### X-Ray Tracing

O X-Ray está habilitado. Acesse o console do X-Ray para ver traces distribuídos.

## Variáveis de Ambiente

As seguintes variáveis são configuradas automaticamente no Lambda:

- `ENVIRONMENT`: production/staging
- `SNS_TOPIC_ARN`: ARN do SNS Topic
- `AWS_REGION`: Região AWS

## Custos Estimados

Para uso típico:

- **Lambda**: ~$0.20/milhão de requisições (256MB, 30s timeout)
- **API Gateway**: $1.00/milhão de requisições
- **SNS**: $0.50/milhão de publicações
- **CloudWatch Logs**: ~$0.50/GB de logs

**Estimativa mensal** para 1M de transações: ~$2-5

## Segurança

- Lambda tem permissões mínimas (apenas SNS:Publish)
- X-Ray tracing habilitado para debugging
- CloudWatch Logs com retenção de 7 dias
- SNS Topic com policy restritiva

## Limpeza (Destruir Infraestrutura)

```bash
cd terraform
terraform destroy
```

⚠️ **CUIDADO**: Isso irá deletar TODOS os recursos criados.

## Troubleshooting

### Lambda não consegue publicar no SNS

Verifique as permissões IAM:
```bash
aws iam get-role-policy \
  --role-name chainorchestrator-lambda-execution-role-production \
  --policy-name chainorchestrator-sns-publish-policy
```

### API Gateway retorna 500

Verifique os logs:
```bash
aws logs tail /aws/lambda/chainorchestrator-production --follow
```

### Timeout na Lambda

Aumente o timeout no `terraform/variables.tf`:
```hcl
variable "lambda_timeout" {
  default = 60  # aumentar de 30 para 60 segundos
}
```

## Ambientes Múltiplos

Para criar múltiplos ambientes (staging, production):

```bash
# Staging
terraform workspace new staging
terraform apply -var="environment=staging"

# Production
terraform workspace new production
terraform apply -var="environment=production"
```
