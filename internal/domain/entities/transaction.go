package entities

import (
	"time"

	"github.com/gabrielksneiva/ChainOrchestrator/internal/domain/valueobjects"
)

// Transaction representa uma transação blockchain no domínio
type Transaction struct {
	operationID   valueobjects.OperationID
	chainType     valueobjects.ChainType
	operationType valueobjects.OperationType
	payload       map[string]interface{}
	createdAt     time.Time
	status        TransactionStatus
}

// TransactionStatus status da transação no orquestrador
type TransactionStatus string

const (
	TransactionStatusQueued    TransactionStatus = "QUEUED"
	TransactionStatusPublished TransactionStatus = "PUBLISHED"
	TransactionStatusFailed    TransactionStatus = "FAILED"
)

// NewTransaction cria uma nova transação
func NewTransaction(
	chainType valueobjects.ChainType,
	operationType valueobjects.OperationType,
	payload map[string]interface{},
) (*Transaction, error) {
	operationID := valueobjects.NewOperationID()

	return &Transaction{
		operationID:   operationID,
		chainType:     chainType,
		operationType: operationType,
		payload:       payload,
		createdAt:     time.Now(),
		status:        TransactionStatusQueued,
	}, nil
}

// OperationID retorna o ID da operação
func (t *Transaction) OperationID() valueobjects.OperationID {
	return t.operationID
}

// ChainType retorna o tipo de blockchain
func (t *Transaction) ChainType() valueobjects.ChainType {
	return t.chainType
}

// OperationType retorna o tipo de operação
func (t *Transaction) OperationType() valueobjects.OperationType {
	return t.operationType
}

// Payload retorna o payload da transação
func (t *Transaction) Payload() map[string]interface{} {
	return t.payload
}

// CreatedAt retorna a data de criação
func (t *Transaction) CreatedAt() time.Time {
	return t.createdAt
}

// Status retorna o status da transação
func (t *Transaction) Status() TransactionStatus {
	return t.status
}

// MarkAsPublished marca a transação como publicada
func (t *Transaction) MarkAsPublished() {
	t.status = TransactionStatusPublished
}

// MarkAsFailed marca a transação como falha
func (t *Transaction) MarkAsFailed() {
	t.status = TransactionStatusFailed
}
