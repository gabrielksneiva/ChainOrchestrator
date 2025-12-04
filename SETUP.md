# 🚀 Setup Completo do ChainOrchestrator com Terraform

## 📋 Pré-requisitos

- AWS CLI configurado
- Terraform 1.7+
- Conta AWS: 490873503238
- GitHub repo: gabrielksneiva/ChainOrchestrator

## 🎯 Setup Inicial (Execute UMA VEZ)

### 1️⃣ Criar Terraform Backend (S3 + DynamoDB)

```bash
# Executar na raiz do projeto
./.github/setup-terraform-backend.sh
```

Isso cria:
- ✅ S3 bucket `chainorchestrator-terraform-state` (versionado + criptografado)
- ✅ DynamoDB table `terraform-state-lock` (para lock do state)

### 2️⃣ Atualizar Permissões do GitHub Actions Role

```bash
# Atualizar a policy do role ChainOrchestrator-GitHub
aws iam put-role-policy \
  --role-name ChainOrchestrator-GitHub \
  --policy-name TerraformFullAccess \
  --policy-document file://.github/github-role-policy-terraform.json

# Verificar
aws iam get-role-policy \
  --role-name ChainOrchestrator-GitHub \
  --policy-name TerraformFullAccess
```

### 3️⃣ Deploy Inicial do Terraform (Manual)

```bash
cd terraform

# Inicializar
terraform init

# Validar
terraform validate
terraform fmt -check

# Revisar mudanças
terraform plan

# Aplicar (cria toda infraestrutura)
terraform apply
```

Isso cria:
- ✅ Lambda function: `chainorchestrator-production`
- ✅ API Gateway HTTP API
- ✅ SNS Topic: `chainorchestrator-Transactions-production`
- ✅ SQS Queues: EVM, Bitcoin, TRON, Solana (+ DLQs)
- ✅ IAM Roles e Policies
- ✅ CloudWatch Log Groups

### 4️⃣ Commit e Push

```bash
git add .
git commit -m "feat: add Terraform infrastructure with API Gateway, SNS and SQS"
git push origin develop
```

## 🔄 Workflow Automático (após setup)

Agora, **todo push** para `main` ou `develop` vai:

1. ✅ Lint (golangci-lint)
2. ✅ Security scans (Trivy, Gosec, Nancy, CodeQL)
3. ✅ Tests (90%+ coverage)
4. ✅ Build (Go arm64 binary)
5. ✅ Build Docker image
6. ✅ **Terraform Plan** (preview changes)
7. ✅ **Terraform Apply** (deploy infrastructure)
8. ✅ Update Lambda code (novo binário)
9. ✅ Health check (Lambda + API Gateway)

## 🧪 Testar Deploy

### Via Lambda Diretamente

```bash
aws lambda invoke \
  --function-name chainorchestrator-production \
  --payload '{"httpMethod":"GET","path":"/health"}' \
  response.json && cat response.json
```

### Via API Gateway

```bash
# Pegar URL do API Gateway
API_URL=$(cd terraform && terraform output -raw api_endpoint)

# Test health
curl $API_URL/health

# Test transaction (POST)
curl -X POST $API_URL/transaction \
  -H "Content-Type: application/json" \
  -d '{"blockchain":"EVM","amount":"1.5"}'
```

### Testar SNS → SQS

```bash
SNS_TOPIC=$(cd terraform && terraform output -raw sns_topic_arn)

# Publicar mensagem (vai para EVM queue)
aws sns publish \
  --topic-arn "$SNS_TOPIC" \
  --message '{"type":"transfer","amount":"10"}' \
  --message-attributes '{"blockchain":{"DataType":"String","StringValue":"EVM"}}'

# Ver mensagens na fila EVM
aws sqs receive-message \
  --queue-url $(cd terraform && terraform output -json sqs_queues | jq -r '.evm')
```

## 🎯 Fluxo de Desenvolvimento

### Mudança de Código

```bash
# 1. Editar código Go
vim cmd/lambda/main.go

# 2. Commit e push
git add .
git commit -m "feat: add new feature"
git push origin develop

# 3. CI/CD faz TUDO automaticamente:
#    - Tests, build, deploy Terraform, atualiza Lambda
```

### Mudança de Infraestrutura

```bash
# 1. Editar Terraform
vim terraform/lambda.tf

# 2. Testar localmente (opcional)
cd terraform
terraform plan

# 3. Commit e push
git add .
git commit -m "feat: increase Lambda memory to 1024MB"
git push origin develop

# 4. CI/CD aplica as mudanças via Terraform
```

## 📊 Outputs do Terraform

```bash
cd terraform

terraform output api_endpoint        # URL do API Gateway
terraform output lambda_function_name # Nome do Lambda
terraform output sns_topic_arn       # ARN do SNS
terraform output sqs_queues          # URLs das filas
```

## 🔧 Comandos Úteis

### Ver Logs do Lambda

```bash
aws logs tail /aws/lambda/chainorchestrator-production --follow
```

### Ver State do Terraform

```bash
cd terraform
terraform show
```

### Forçar Re-deploy do Lambda (sem mudar código)

```bash
cd terraform
terraform taint aws_lambda_function.orchestrator
terraform apply
```

## ⚠️ Importante

1. **Não delete o Lambda manualmente** - Terraform gerencia
2. **Não edite resources pela console** - Use Terraform
3. **State está no S3** - Nunca commite `terraform.tfstate`
4. **CI/CD atualiza CÓDIGO**, Terraform atualiza INFRA

## 🗑️ Destruir Tudo (Cleanup)

```bash
cd terraform
terraform destroy
```

⚠️ **CUIDADO**: Isso deleta TUDO (Lambda, API Gateway, SNS, SQS, etc)

## ❓ FAQ

**Q: Quando preciso rodar Terraform manualmente?**  
A: Nunca! CI/CD roda automaticamente. Só rode localmente para testar.

**Q: Como atualizar só o código do Lambda?**  
A: Só fazer push. CI/CD builda e faz deploy.

**Q: Como adicionar nova fila SQS?**  
A: Editar `terraform/sns.tf`, commit, push. CI/CD aplica.

**Q: E se houver conflito no state?**  
A: DynamoDB lock previne. Se travar, use `terraform force-unlock <ID>`

**Q: Posso ter staging + production?**  
A: Sim! Crie workspaces Terraform ou duplique `terraform.tfvars`
