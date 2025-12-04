package usecases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gabrielksneiva/ChainOrchestrator/internal/application/dtos"
	"github.com/gabrielksneiva/ChainOrchestrator/internal/application/usecases"
	"github.com/gabrielksneiva/ChainOrchestrator/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSNSPublisher mock do SNS Publisher
type MockSNSPublisher struct {
	mock.Mock
}

func (m *MockSNSPublisher) Publish(ctx context.Context, message string, chainType string) error {
	args := m.Called(ctx, message, chainType)
	return args.Error(0)
}

func TestPublishTransactionUseCase_Execute_Success(t *testing.T) {
	mockSNS := new(MockSNSPublisher)
	log := logger.NewNopLogger()
	useCase := usecases.NewPublishTransactionUseCase(mockSNS, log)

	req := &dtos.PublishTransactionRequest{
		ChainType:     "EVM",
		OperationType: "TRANSFER",
		Payload: map[string]interface{}{
			"from":   "0x123",
			"to":     "0x456",
			"amount": "1000",
		},
	}

	mockSNS.On("Publish", mock.Anything, mock.Anything, "EVM").Return(nil)

	resp, err := useCase.Execute(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "EVM", resp.ChainType)
	assert.Equal(t, "QUEUED", resp.Status)
	assert.NotEmpty(t, resp.OperationID)
	assert.NotEmpty(t, resp.CreatedAt)
	assert.Contains(t, resp.Message, "queued for processing")

	mockSNS.AssertExpectations(t)
	mockSNS.AssertCalled(t, "Publish", mock.Anything, mock.Anything, "EVM")
}

func TestPublishTransactionUseCase_Execute_SNSError(t *testing.T) {
	mockSNS := new(MockSNSPublisher)
	log := logger.NewNopLogger()
	useCase := usecases.NewPublishTransactionUseCase(mockSNS, log)

	req := &dtos.PublishTransactionRequest{
		ChainType:     "BTC",
		OperationType: "TRANSFER",
		Payload:       map[string]interface{}{},
	}

	expectedError := errors.New("SNS publish failed")
	mockSNS.On("Publish", mock.Anything, mock.Anything, "BTC").Return(expectedError)

	resp, err := useCase.Execute(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to publish transaction")

	mockSNS.AssertExpectations(t)
}

func TestPublishTransactionUseCase_Execute_DifferentChainTypes(t *testing.T) {
	chainTypes := []string{"EVM", "BTC", "TRON", "SOL"}

	for _, chainType := range chainTypes {
		t.Run(chainType, func(t *testing.T) {
			mockSNS := new(MockSNSPublisher)
			log := logger.NewNopLogger()
			useCase := usecases.NewPublishTransactionUseCase(mockSNS, log)

			req := &dtos.PublishTransactionRequest{
				ChainType:     chainType,
				OperationType: "TRANSFER",
				Payload:       map[string]interface{}{},
			}

			mockSNS.On("Publish", mock.Anything, mock.Anything, chainType).Return(nil)

			resp, err := useCase.Execute(context.Background(), req)

			assert.NoError(t, err)
			assert.Equal(t, chainType, resp.ChainType)
			mockSNS.AssertExpectations(t)
		})
	}
}

func TestPublishTransactionUseCase_Execute_ComplexPayload(t *testing.T) {
	mockSNS := new(MockSNSPublisher)
	log := logger.NewNopLogger()
	useCase := usecases.NewPublishTransactionUseCase(mockSNS, log)

	req := &dtos.PublishTransactionRequest{
		ChainType:     "EVM",
		OperationType: "SWAP",
		Payload: map[string]interface{}{
			"from_token": "USDT",
			"to_token":   "USDC",
			"amount_in":  "1000000",
			"amount_out": "999000",
			"slippage":   "0.5",
			"deadline":   "1234567890",
			"router":     "0xabc...",
			"path":       []string{"0x123", "0x456"},
		},
	}

	mockSNS.On("Publish", mock.Anything, mock.Anything, "EVM").Return(nil)

	resp, err := useCase.Execute(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	mockSNS.AssertExpectations(t)
}

func TestPublishTransactionUseCase_Execute_ContextCancellation(t *testing.T) {
	mockSNS := new(MockSNSPublisher)
	log := logger.NewNopLogger()
	useCase := usecases.NewPublishTransactionUseCase(mockSNS, log)

	req := &dtos.PublishTransactionRequest{
		ChainType:     "EVM",
		OperationType: "TRANSFER",
		Payload:       map[string]interface{}{},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	mockSNS.On("Publish", mock.Anything, mock.Anything, "EVM").Return(context.Canceled)

	resp, err := useCase.Execute(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestNewPublishTransactionUseCase(t *testing.T) {
	mockSNS := new(MockSNSPublisher)
	log := logger.NewNopLogger()

	useCase := usecases.NewPublishTransactionUseCase(mockSNS, log)

	assert.NotNil(t, useCase)
}

func TestPublishTransactionUseCase_Execute_UnmarshalablePayload(t *testing.T) {
	mockSNS := new(MockSNSPublisher)
	log := logger.NewNopLogger()
	useCase := usecases.NewPublishTransactionUseCase(mockSNS, log)

	// Payload com canal (não serializável em JSON)
	req := &dtos.PublishTransactionRequest{
		ChainType:     "EVM",
		OperationType: "TRANSFER",
		Payload: map[string]interface{}{
			"channel": make(chan int),
		},
	}

	resp, err := useCase.Execute(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to marshal")
}
