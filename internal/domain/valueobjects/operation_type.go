package valueobjects

import "fmt"

// OperationType representa o tipo de operação blockchain
type OperationType string

const (
	OperationTypeTransfer OperationType = "TRANSFER"
	OperationTypeDeploy   OperationType = "DEPLOY"
	OperationTypeCall     OperationType = "CALL"
	OperationTypeApprove  OperationType = "APPROVE"
	OperationTypeSwap     OperationType = "SWAP"
	OperationTypeStake    OperationType = "STAKE"
	OperationTypeUnstake  OperationType = "UNSTAKE"
	OperationTypeWithdraw OperationType = "WITHDRAW"
	OperationTypeMint     OperationType = "MINT"
	OperationTypeBurn     OperationType = "BURN"
)

// NewOperationType cria e valida um novo OperationType
func NewOperationType(value string) (OperationType, error) {
	ot := OperationType(value)
	if !ot.IsValid() {
		return "", fmt.Errorf("invalid operation type: %s", value)
	}
	return ot, nil
}

// IsValid verifica se o operation type é válido
func (o OperationType) IsValid() bool {
	switch o {
	case OperationTypeTransfer, OperationTypeDeploy, OperationTypeCall,
		OperationTypeApprove, OperationTypeSwap, OperationTypeStake,
		OperationTypeUnstake, OperationTypeWithdraw, OperationTypeMint,
		OperationTypeBurn:
		return true
	default:
		return false
	}
}

// String retorna a representação em string
func (o OperationType) String() string {
	return string(o)
}

// Equals verifica igualdade
func (o OperationType) Equals(other OperationType) bool {
	return o == other
}

// RequiresAmount verifica se a operação requer amount
func (o OperationType) RequiresAmount() bool {
	switch o {
	case OperationTypeTransfer, OperationTypeSwap, OperationTypeStake,
		OperationTypeWithdraw, OperationTypeMint:
		return true
	default:
		return false
	}
}

// RequiresRecipient verifica se a operação requer destinatário
func (o OperationType) RequiresRecipient() bool {
	switch o {
	case OperationTypeTransfer, OperationTypeMint:
		return true
	default:
		return false
	}
}
