# ChainOrchestrator — Readme Oficial

O **ChainOrchestrator** é o serviço responsável por orquestrar operações multi-chain.  
Ele atua como camada de entrada, validação, roteamento e padronização das requisições que chegam pela API Gateway, enviando eventos para processamento assíncrono pelos serviços especializados.

---

## 🧱 Arquitetura Obrigatória

O ChainOrchestrator deve seguir integralmente:

- **Clean Architecture**  
- **Domain-Driven Design (DDD)**  
- **Event-Driven Architecture**  
- **Logs estruturados via ZAP**  
- **Dependency Injection com FX**  
- **Golang + Fiber**  
- Código totalmente **testável, modular e desacoplado**  
- **Nenhuma dependência direta de blockchains**  

---

## 🎯 Função e Propósito

O ChainOrchestrator **não executa lógica on-chain**.  
Ele é a camada de coordenação e direcionamento do sistema blockchain.

### Funções principais:

1. **Receber requisições HTTP via API Gateway:**
   - `POST /transaction`  
   - `GET /walletbalance`  
   - `GET /transaction-status`  

2. **Validar entradas** (payloads, tipos de operação, integridade e regras de negócio)

3. **Decidir para qual blockchain a operação deve ser enviada**  
   - Exemplo: EVM → publica evento no SNS Topic "Transactions"

4. **Publicar eventos padronizados** para processamento assíncrono

5. **Garantir:**
   - Idempotência  
   - Rastreabilidade  
   - Logs estruturados  
   - Normalização de erros  

---

## 🔁 Fluxo Arquitetural

```
1. API Gateway → ChainOrchestrator
2. Validação + normalização da operação
3. Roteamento (EVM, TRON, BTC, SOL — conforme implementado)
4. Publicação no SNS Topic "Transactions"
5. O serviço especializado (ex: ChainEVM) assume o processamento
```

---

## 🏗️ Estrutura do Projeto

```
ChainOrchestrator/
├── cmd/
│   ├── lambda/              # Handler Lambda (se necessário)
│   └── server/              # Servidor HTTP (Fiber)
│       └── main.go
├── internal/
│   ├── application/         # Camada de Aplicação
│   │   ├── dtos/           # Data Transfer Objects
│   │   └── usecases/       # Casos de Uso
│   ├── domain/             # Camada de Domínio (Entities, Value Objects)
│   ├── infrastructure/     # Camada de Infraestrutura
│   │   ├── eventbus/       # SNS Publisher
│   │   ├── http/           # Router (Fiber)
│   │   └── logger/         # Logger (Zap)
│   └── interfaces/         # Camada de Interface
│       ├── handlers/       # HTTP Handlers
│       └── middleware/     # Middlewares (logging, error handling)
├── pkg/
│   ├── config/             # Configurações
│   └── errors/             # Erros customizados
├── terraform/              # Infraestrutura como Código
├── docs/                   # Documentação
├── go.mod
├── Makefile
└── README.md
```

---

## 🛠️ Stack Tecnológica

- **Linguagem:** Go 1.24+
- **Framework HTTP:** Fiber v2
- **Logging:** Zap (estruturado)
- **Dependency Injection:** Fx
- **Validação:** go-playground/validator
- **Event Bus:** AWS SNS
- **Cloud:** AWS (Lambda, SNS, API Gateway)
- **IaC:** Terraform

---

## 🚀 Como Executar

### Pré-requisitos

- Go 1.24+
- AWS CLI configurado (para SNS)
- Credenciais AWS com permissões para SNS

### Configuração

Configure as variáveis de ambiente:

```bash
export ENVIRONMENT=development
export SNS_TOPIC_ARN=arn:aws:sns:us-east-1:123456789012:Transactions
export PORT=3000
```

### Executar localmente

```bash
# Instalar dependências
go mod download

# Executar servidor
go run cmd/server/main.go
```

### Build

```bash
# Build do binário
make build

# Executar binário
./bin/chainorchestrator
```

---

## 📡 Endpoints

### POST /transaction

Publica uma transação para processamento assíncrono.

**Request:**
```json
{
  "chain_type": "EVM",
  "operation_type": "TRANSFER",
  "payload": {
    "from": "0x123...",
    "to": "0x456...",
    "amount": "1000000000000000000",
    "token": "USDT"
  }
}
```

**Response:**
```json
{
  "operation_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "published",
  "created_at": "2025-12-03T10:30:00Z"
}
```

### GET /walletbalance

Consulta saldo de carteira (futura implementação).

### GET /transaction-status

Consulta status de transação (futura implementação).

---

## 🧪 Regras de Testes

- **TDD obrigatório**  
- **Coverage ≥ 90%**  
- **Sem TODO / not implemented / stubs artificiais**  
- Testes unitários e integração isolados  

### Executar testes

```bash
# Rodar todos os testes
make test

# Rodar com coverage
make test-coverage

# Rodar testes específicos
go test ./internal/application/usecases/... -v
```

---

## ⚠️ O que o ChainOrchestrator NÃO faz

O ChainOrchestrator é **agnóstico a blockchains** e **não executa operações on-chain**.

❌ **NÃO assina transações**  
❌ **NÃO consulta RPC**  
❌ **NÃO envia transações**  
❌ **NÃO acessa blockchains**  
❌ **NÃO persiste dados** (apenas coordena)  
❌ **NÃO executa lógica de domínio de EVM ou TRON**  

---

## ✅ O que o ChainOrchestrator FAZ

✅ **Recebe e valida requisições HTTP**  
✅ **Normaliza e padroniza payloads**  
✅ **Roteia para o blockchain correto**  
✅ **Publica eventos no SNS**  
✅ **Garante rastreabilidade e logs estruturados**  
✅ **Aplica regras de negócio de orquestração**  

---

## 🎯 Objetivo

Prover um **gateway inteligente** para todo o ecossistema blockchain, mantendo desacoplamento total entre API e lógica on-chain.

O ChainOrchestrator é o **ponto de entrada único** para todas as operações multi-chain, garantindo:

- **Padronização** de contratos de API
- **Resiliência** através de arquitetura event-driven
- **Escalabilidade** horizontal
- **Observabilidade** com logs estruturados
- **Testabilidade** total com injeção de dependências

---

## 📚 Documentação Adicional

- [Arquitetura](docs/architecture.md) *(em breve)*
- [Blockchains Suportados](docs/blockchains/) *(em breve)*
- [Event Contracts](docs/events.md) *(em breve)*

---

## 📄 Licença

Este projeto é proprietário e confidencial.

---

## 👨‍💻 Mantido por

Gabriel Neiva  
[@gabrielksneiva](https://github.com/gabrielksneiva)
