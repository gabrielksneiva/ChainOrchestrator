package dtos_test

import (
	"encoding/json"
	"testing"

	"github.com/gabrielksneiva/ChainOrchestrator/internal/application/dtos"
	"github.com/stretchr/testify/assert"
)

func TestPublishTransactionRequest_JSON(t *testing.T) {
	req := &dtos.PublishTransactionRequest{
		ChainType:     "EVM",
		OperationType: "TRANSFER",
		Payload: map[string]interface{}{
			"from":   "0x123",
			"to":     "0x456",
			"amount": "1000",
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(req)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Unmarshal back
	var decoded dtos.PublishTransactionRequest
	err = json.Unmarshal(jsonData, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, req.ChainType, decoded.ChainType)
	assert.Equal(t, req.OperationType, decoded.OperationType)
}

func TestPublishTransactionResponse_JSON(t *testing.T) {
	resp := &dtos.PublishTransactionResponse{
		OperationID: "550e8400-e29b-41d4-a716-446655440000",
		ChainType:   "EVM",
		Status:      "QUEUED",
		Message:     "Transaction queued for processing",
		CreatedAt:   "2025-12-03T10:30:00Z",
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Unmarshal back
	var decoded dtos.PublishTransactionResponse
	err = json.Unmarshal(jsonData, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, resp.OperationID, decoded.OperationID)
	assert.Equal(t, resp.ChainType, decoded.ChainType)
	assert.Equal(t, resp.Status, decoded.Status)
}
