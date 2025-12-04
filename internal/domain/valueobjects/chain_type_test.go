package valueobjects_test

import (
	"testing"

	"github.com/gabrielksneiva/ChainOrchestrator/internal/domain/valueobjects"
)

func TestNewChainType_Valid(t *testing.T) {
	validTypes := []string{"EVM", "BTC", "TRON", "SOL"}

	for _, validType := range validTypes {
		ct, err := valueobjects.NewChainType(validType)
		if err != nil {
			t.Errorf("NewChainType should accept valid type %s, got error: %v", validType, err)
		}

		if ct.String() != validType {
			t.Errorf("Expected ChainType %s, got %s", validType, ct.String())
		}
	}
}

func TestNewChainType_Invalid(t *testing.T) {
	invalidTypes := []string{
		"",
		"INVALID",
		"evm",
		"Ethereum",
		"ETH",
		"Bitcoin",
		"POLYGON",
	}

	for _, invalidType := range invalidTypes {
		_, err := valueobjects.NewChainType(invalidType)
		if err == nil {
			t.Errorf("NewChainType should reject invalid type: %s", invalidType)
		}
	}
}

func TestChainTypeIsValid(t *testing.T) {
	tests := []struct {
		name     string
		value    valueobjects.ChainType
		expected bool
	}{
		{"EVM is valid", valueobjects.ChainTypeEVM, true},
		{"BTC is valid", valueobjects.ChainTypeBTC, true},
		{"TRON is valid", valueobjects.ChainTypeTRON, true},
		{"SOL is valid", valueobjects.ChainTypeSOL, true},
		{"Invalid type", valueobjects.ChainType("INVALID"), false},
		{"Empty type", valueobjects.ChainType(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value.IsValid() != tt.expected {
				t.Errorf("Expected IsValid() to return %v for %s", tt.expected, tt.value)
			}
		})
	}
}

func TestChainTypeEquals(t *testing.T) {
	evm1 := valueobjects.ChainTypeEVM
	evm2 := valueobjects.ChainTypeEVM
	btc := valueobjects.ChainTypeBTC

	if !evm1.Equals(evm2) {
		t.Error("Equal ChainTypes should return true")
	}

	if evm1.Equals(btc) {
		t.Error("Different ChainTypes should return false")
	}
}

func TestChainTypeSNSTopicAttribute(t *testing.T) {
	tests := []struct {
		chainType valueobjects.ChainType
		expected  string
	}{
		{valueobjects.ChainTypeEVM, "EVM"},
		{valueobjects.ChainTypeBTC, "BTC"},
		{valueobjects.ChainTypeTRON, "TRON"},
		{valueobjects.ChainTypeSOL, "SOL"},
	}

	for _, tt := range tests {
		t.Run(string(tt.chainType), func(t *testing.T) {
			if tt.chainType.SNSTopicAttribute() != tt.expected {
				t.Errorf("Expected SNSTopicAttribute() to return %s, got %s",
					tt.expected, tt.chainType.SNSTopicAttribute())
			}
		})
	}
}

func TestChainTypeConstants(t *testing.T) {
	if valueobjects.ChainTypeEVM != "EVM" {
		t.Error("ChainTypeEVM constant should be 'EVM'")
	}
	if valueobjects.ChainTypeBTC != "BTC" {
		t.Error("ChainTypeBTC constant should be 'BTC'")
	}
	if valueobjects.ChainTypeTRON != "TRON" {
		t.Error("ChainTypeTRON constant should be 'TRON'")
	}
	if valueobjects.ChainTypeSOL != "SOL" {
		t.Error("ChainTypeSOL constant should be 'SOL'")
	}
}
