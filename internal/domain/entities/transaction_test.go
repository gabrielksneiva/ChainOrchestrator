package entities_test

import (
	"testing"
	"time"

	"github.com/gabrielksneiva/ChainOrchestrator/internal/domain/entities"
	"github.com/gabrielksneiva/ChainOrchestrator/internal/domain/valueobjects"
)

func TestNewTransaction(t *testing.T) {
	chainType := valueobjects.ChainTypeEVM
	operationType := valueobjects.OperationTypeTransfer
	payload := map[string]interface{}{
		"from":   "0x123",
		"to":     "0x456",
		"amount": "1000",
	}

	tx, err := entities.NewTransaction(chainType, operationType, payload)

	if err != nil {
		t.Errorf("NewTransaction should not return error for valid inputs, got: %v", err)
	}

	if tx == nil {
		t.Fatal("NewTransaction should return a transaction instance")
	}

	if tx.ChainType() != chainType {
		t.Errorf("Expected ChainType %s, got %s", chainType, tx.ChainType())
	}

	if tx.OperationType() != operationType {
		t.Errorf("Expected OperationType %s, got %s", operationType, tx.OperationType())
	}

	if tx.Status() != entities.TransactionStatusQueued {
		t.Errorf("Expected initial status QUEUED, got %s", tx.Status())
	}

	if tx.OperationID().String() == "" {
		t.Error("Transaction should have a valid operation ID")
	}

	if time.Since(tx.CreatedAt()) > time.Second {
		t.Error("Transaction CreatedAt should be recent")
	}
}

func TestTransactionPayload(t *testing.T) {
	payload := map[string]interface{}{
		"from":   "0x123",
		"to":     "0x456",
		"amount": "1000",
		"token":  "USDT",
	}

	tx, _ := entities.NewTransaction(
		valueobjects.ChainTypeEVM,
		valueobjects.OperationTypeTransfer,
		payload,
	)

	retrievedPayload := tx.Payload()

	if len(retrievedPayload) != len(payload) {
		t.Errorf("Expected payload length %d, got %d", len(payload), len(retrievedPayload))
	}

	for key, value := range payload {
		if retrievedPayload[key] != value {
			t.Errorf("Expected payload[%s] = %v, got %v", key, value, retrievedPayload[key])
		}
	}
}

func TestTransactionMarkAsPublished(t *testing.T) {
	tx, _ := entities.NewTransaction(
		valueobjects.ChainTypeEVM,
		valueobjects.OperationTypeTransfer,
		map[string]interface{}{},
	)

	if tx.Status() != entities.TransactionStatusQueued {
		t.Errorf("Initial status should be QUEUED, got %s", tx.Status())
	}

	tx.MarkAsPublished()

	if tx.Status() != entities.TransactionStatusPublished {
		t.Errorf("After MarkAsPublished, status should be PUBLISHED, got %s", tx.Status())
	}
}

func TestTransactionMarkAsFailed(t *testing.T) {
	tx, _ := entities.NewTransaction(
		valueobjects.ChainTypeEVM,
		valueobjects.OperationTypeTransfer,
		map[string]interface{}{},
	)

	if tx.Status() != entities.TransactionStatusQueued {
		t.Errorf("Initial status should be QUEUED, got %s", tx.Status())
	}

	tx.MarkAsFailed()

	if tx.Status() != entities.TransactionStatusFailed {
		t.Errorf("After MarkAsFailed, status should be FAILED, got %s", tx.Status())
	}
}

func TestTransactionStatusTransitions(t *testing.T) {
	tx, _ := entities.NewTransaction(
		valueobjects.ChainTypeEVM,
		valueobjects.OperationTypeTransfer,
		map[string]interface{}{},
	)

	// QUEUED -> PUBLISHED -> FAILED
	tx.MarkAsPublished()
	if tx.Status() != entities.TransactionStatusPublished {
		t.Error("Expected status PUBLISHED")
	}

	tx.MarkAsFailed()
	if tx.Status() != entities.TransactionStatusFailed {
		t.Error("Expected status FAILED")
	}
}

func TestTransactionWithDifferentChainTypes(t *testing.T) {
	chainTypes := []valueobjects.ChainType{
		valueobjects.ChainTypeEVM,
		valueobjects.ChainTypeBTC,
		valueobjects.ChainTypeTRON,
		valueobjects.ChainTypeSOL,
	}

	for _, ct := range chainTypes {
		t.Run(string(ct), func(t *testing.T) {
			tx, err := entities.NewTransaction(
				ct,
				valueobjects.OperationTypeTransfer,
				map[string]interface{}{},
			)

			if err != nil {
				t.Errorf("Should create transaction for chain type %s", ct)
			}

			if tx.ChainType() != ct {
				t.Errorf("Expected chain type %s, got %s", ct, tx.ChainType())
			}
		})
	}
}

func TestTransactionWithDifferentOperationTypes(t *testing.T) {
	operationTypes := []valueobjects.OperationType{
		valueobjects.OperationTypeTransfer,
		valueobjects.OperationTypeDeploy,
		valueobjects.OperationTypeCall,
		valueobjects.OperationTypeSwap,
		valueobjects.OperationTypeStake,
	}

	for _, ot := range operationTypes {
		t.Run(string(ot), func(t *testing.T) {
			tx, err := entities.NewTransaction(
				valueobjects.ChainTypeEVM,
				ot,
				map[string]interface{}{},
			)

			if err != nil {
				t.Errorf("Should create transaction for operation type %s", ot)
			}

			if tx.OperationType() != ot {
				t.Errorf("Expected operation type %s, got %s", ot, tx.OperationType())
			}
		})
	}
}

func TestTransactionStatusConstants(t *testing.T) {
	if entities.TransactionStatusQueued != "QUEUED" {
		t.Error("TransactionStatusQueued should be 'QUEUED'")
	}
	if entities.TransactionStatusPublished != "PUBLISHED" {
		t.Error("TransactionStatusPublished should be 'PUBLISHED'")
	}
	if entities.TransactionStatusFailed != "FAILED" {
		t.Error("TransactionStatusFailed should be 'FAILED'")
	}
}
