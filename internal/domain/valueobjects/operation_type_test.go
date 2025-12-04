package valueobjects_test

import (
	"testing"

	"github.com/gabrielksneiva/ChainOrchestrator/internal/domain/valueobjects"
)

func TestNewOperationType_Valid(t *testing.T) {
	validTypes := []string{
		"TRANSFER", "DEPLOY", "CALL", "APPROVE", "SWAP",
		"STAKE", "UNSTAKE", "WITHDRAW", "MINT", "BURN",
	}

	for _, validType := range validTypes {
		ot, err := valueobjects.NewOperationType(validType)
		if err != nil {
			t.Errorf("NewOperationType should accept valid type %s, got error: %v", validType, err)
		}

		if ot.String() != validType {
			t.Errorf("Expected OperationType %s, got %s", validType, ot.String())
		}
	}
}

func TestNewOperationType_Invalid(t *testing.T) {
	invalidTypes := []string{
		"",
		"INVALID",
		"transfer",
		"Send",
		"RECEIVE",
		"DELETE",
	}

	for _, invalidType := range invalidTypes {
		_, err := valueobjects.NewOperationType(invalidType)
		if err == nil {
			t.Errorf("NewOperationType should reject invalid type: %s", invalidType)
		}
	}
}

func TestOperationTypeIsValid(t *testing.T) {
	tests := []struct {
		name     string
		value    valueobjects.OperationType
		expected bool
	}{
		{"TRANSFER is valid", valueobjects.OperationTypeTransfer, true},
		{"DEPLOY is valid", valueobjects.OperationTypeDeploy, true},
		{"CALL is valid", valueobjects.OperationTypeCall, true},
		{"APPROVE is valid", valueobjects.OperationTypeApprove, true},
		{"SWAP is valid", valueobjects.OperationTypeSwap, true},
		{"STAKE is valid", valueobjects.OperationTypeStake, true},
		{"UNSTAKE is valid", valueobjects.OperationTypeUnstake, true},
		{"WITHDRAW is valid", valueobjects.OperationTypeWithdraw, true},
		{"MINT is valid", valueobjects.OperationTypeMint, true},
		{"BURN is valid", valueobjects.OperationTypeBurn, true},
		{"Invalid type", valueobjects.OperationType("INVALID"), false},
		{"Empty type", valueobjects.OperationType(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value.IsValid() != tt.expected {
				t.Errorf("Expected IsValid() to return %v for %s", tt.expected, tt.value)
			}
		})
	}
}

func TestOperationTypeEquals(t *testing.T) {
	transfer1 := valueobjects.OperationTypeTransfer
	transfer2 := valueobjects.OperationTypeTransfer
	deploy := valueobjects.OperationTypeDeploy

	if !transfer1.Equals(transfer2) {
		t.Error("Equal OperationTypes should return true")
	}

	if transfer1.Equals(deploy) {
		t.Error("Different OperationTypes should return false")
	}
}

func TestOperationTypeRequiresAmount(t *testing.T) {
	tests := []struct {
		name     string
		opType   valueobjects.OperationType
		expected bool
	}{
		{"TRANSFER requires amount", valueobjects.OperationTypeTransfer, true},
		{"SWAP requires amount", valueobjects.OperationTypeSwap, true},
		{"STAKE requires amount", valueobjects.OperationTypeStake, true},
		{"WITHDRAW requires amount", valueobjects.OperationTypeWithdraw, true},
		{"MINT requires amount", valueobjects.OperationTypeMint, true},
		{"DEPLOY doesn't require amount", valueobjects.OperationTypeDeploy, false},
		{"CALL doesn't require amount", valueobjects.OperationTypeCall, false},
		{"APPROVE doesn't require amount", valueobjects.OperationTypeApprove, false},
		{"UNSTAKE doesn't require amount", valueobjects.OperationTypeUnstake, false},
		{"BURN doesn't require amount", valueobjects.OperationTypeBurn, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.opType.RequiresAmount() != tt.expected {
				t.Errorf("Expected RequiresAmount() to return %v for %s", tt.expected, tt.opType)
			}
		})
	}
}

func TestOperationTypeRequiresRecipient(t *testing.T) {
	tests := []struct {
		name     string
		opType   valueobjects.OperationType
		expected bool
	}{
		{"TRANSFER requires recipient", valueobjects.OperationTypeTransfer, true},
		{"MINT requires recipient", valueobjects.OperationTypeMint, true},
		{"DEPLOY doesn't require recipient", valueobjects.OperationTypeDeploy, false},
		{"CALL doesn't require recipient", valueobjects.OperationTypeCall, false},
		{"APPROVE doesn't require recipient", valueobjects.OperationTypeApprove, false},
		{"SWAP doesn't require recipient", valueobjects.OperationTypeSwap, false},
		{"STAKE doesn't require recipient", valueobjects.OperationTypeStake, false},
		{"UNSTAKE doesn't require recipient", valueobjects.OperationTypeUnstake, false},
		{"WITHDRAW doesn't require recipient", valueobjects.OperationTypeWithdraw, false},
		{"BURN doesn't require recipient", valueobjects.OperationTypeBurn, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.opType.RequiresRecipient() != tt.expected {
				t.Errorf("Expected RequiresRecipient() to return %v for %s", tt.expected, tt.opType)
			}
		})
	}
}

func TestOperationTypeConstants(t *testing.T) {
	if valueobjects.OperationTypeTransfer != "TRANSFER" {
		t.Error("OperationTypeTransfer constant should be 'TRANSFER'")
	}
	if valueobjects.OperationTypeDeploy != "DEPLOY" {
		t.Error("OperationTypeDeploy constant should be 'DEPLOY'")
	}
	if valueobjects.OperationTypeCall != "CALL" {
		t.Error("OperationTypeCall constant should be 'CALL'")
	}
}
