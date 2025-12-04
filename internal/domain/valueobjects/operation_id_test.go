package valueobjects_test

import (
	"testing"

	"github.com/gabrielksneiva/ChainOrchestrator/internal/domain/valueobjects"
	"github.com/google/uuid"
)

func TestNewOperationID(t *testing.T) {
	id := valueobjects.NewOperationID()

	if id.String() == "" {
		t.Error("NewOperationID should generate a non-empty ID")
	}

	// Validate UUID format
	if _, err := uuid.Parse(id.String()); err != nil {
		t.Errorf("NewOperationID should generate a valid UUID, got: %s", id.String())
	}
}

func TestOperationIDFromString_Valid(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	id, err := valueobjects.OperationIDFromString(validUUID)

	if err != nil {
		t.Errorf("OperationIDFromString should accept valid UUID, got error: %v", err)
	}

	if id.String() != validUUID {
		t.Errorf("Expected ID %s, got %s", validUUID, id.String())
	}
}

func TestOperationIDFromString_Invalid(t *testing.T) {
	invalidUUIDs := []string{
		"",
		"invalid-uuid",
		"123",
		"550e8400-e29b-41d4-a716",
		"not-a-uuid-at-all",
	}

	for _, invalid := range invalidUUIDs {
		_, err := valueobjects.OperationIDFromString(invalid)
		if err == nil {
			t.Errorf("OperationIDFromString should reject invalid UUID: %s", invalid)
		}
	}
}

func TestOperationIDEquals(t *testing.T) {
	id1, _ := valueobjects.OperationIDFromString("550e8400-e29b-41d4-a716-446655440000")
	id2, _ := valueobjects.OperationIDFromString("550e8400-e29b-41d4-a716-446655440000")
	id3, _ := valueobjects.OperationIDFromString("550e8400-e29b-41d4-a716-446655440001")

	if !id1.Equals(id2) {
		t.Error("Equal OperationIDs should return true")
	}

	if id1.Equals(id3) {
		t.Error("Different OperationIDs should return false")
	}
}

func TestOperationIDString(t *testing.T) {
	expected := "550e8400-e29b-41d4-a716-446655440000"
	id, _ := valueobjects.OperationIDFromString(expected)

	if id.String() != expected {
		t.Errorf("Expected String() to return %s, got %s", expected, id.String())
	}
}
