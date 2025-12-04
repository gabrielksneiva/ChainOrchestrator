package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gabrielksneiva/ChainOrchestrator/internal/application/dtos"
	"github.com/gabrielksneiva/ChainOrchestrator/internal/infrastructure/logger"
	"github.com/gabrielksneiva/ChainOrchestrator/internal/interfaces/handlers"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPublishUseCase struct {
	mock.Mock
}

func (m *MockPublishUseCase) Execute(ctx context.Context, req *dtos.PublishTransactionRequest) (*dtos.PublishTransactionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dtos.PublishTransactionResponse), args.Error(1)
}

func TestTransactionHandler_PostTransaction_Success(t *testing.T) {
	mockUseCase := new(MockPublishUseCase)
	validate := validator.New()
	log := logger.NewNopLogger()

	handler := handlers.NewTransactionHandler(mockUseCase, validate, log)

	app := fiber.New()
	app.Post("/transaction", handler.PostTransaction)

	requestBody := map[string]interface{}{
		"chain_type":     "EVM",
		"operation_type": "TRANSFER",
		"payload": map[string]interface{}{
			"from": "0x123",
			"to":   "0x456",
		},
	}

	expectedResponse := &dtos.PublishTransactionResponse{
		OperationID: "123",
		ChainType:   "EVM",
		Status:      "QUEUED",
	}

	mockUseCase.On("Execute", mock.Anything, mock.Anything).Return(expectedResponse, nil)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/transaction", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 202, resp.StatusCode)
	mockUseCase.AssertExpectations(t)
}

func TestTransactionHandler_PostTransaction_InvalidJSON(t *testing.T) {
	mockUseCase := new(MockPublishUseCase)
	validate := validator.New()
	log := logger.NewNopLogger()

	handler := handlers.NewTransactionHandler(mockUseCase, validate, log)

	app := fiber.New()
	app.Post("/transaction", handler.PostTransaction)

	req := httptest.NewRequest("POST", "/transaction", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestTransactionHandler_PostTransaction_ValidationFailed(t *testing.T) {
	mockUseCase := new(MockPublishUseCase)
	validate := validator.New()
	log := logger.NewNopLogger()

	handler := handlers.NewTransactionHandler(mockUseCase, validate, log)

	app := fiber.New()
	app.Post("/transaction", handler.PostTransaction)

	requestBody := map[string]interface{}{
		"chain_type": "INVALID",
		// missing operation_type and payload
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/transaction", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestTransactionHandler_PostTransaction_UseCaseError(t *testing.T) {
	mockUseCase := new(MockPublishUseCase)
	validate := validator.New()
	log := logger.NewNopLogger()

	handler := handlers.NewTransactionHandler(mockUseCase, validate, log)

	app := fiber.New()
	app.Post("/transaction", handler.PostTransaction)

	requestBody := map[string]interface{}{
		"chain_type":     "EVM",
		"operation_type": "TRANSFER",
		"payload": map[string]interface{}{
			"from": "0x123",
		},
	}

	mockUseCase.On("Execute", mock.Anything, mock.Anything).Return(nil, errors.New("use case error"))

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/transaction", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)
	mockUseCase.AssertExpectations(t)
}

func TestTransactionHandler_HandleLambda_Success(t *testing.T) {
	mockUseCase := new(MockPublishUseCase)
	validate := validator.New()
	log := logger.NewNopLogger()

	handler := handlers.NewTransactionHandler(mockUseCase, validate, log)

	requestBody := dtos.PublishTransactionRequest{
		ChainType:     "EVM",
		OperationType: "TRANSFER",
		Payload:       map[string]interface{}{"from": "0x123"},
	}

	expectedResponse := &dtos.PublishTransactionResponse{
		OperationID: "123",
		ChainType:   "EVM",
		Status:      "QUEUED",
	}

	mockUseCase.On("Execute", mock.Anything, mock.Anything).Return(expectedResponse, nil)

	bodyJSON, _ := json.Marshal(requestBody)
	lambdaReq := events.APIGatewayV2HTTPRequest{
		Body: string(bodyJSON),
	}

	resp, err := handler.HandleLambda(context.Background(), lambdaReq)

	assert.NoError(t, err)
	assert.Equal(t, 202, resp.StatusCode)
	mockUseCase.AssertExpectations(t)
}

func TestTransactionHandler_HandleLambda_InvalidJSON(t *testing.T) {
	mockUseCase := new(MockPublishUseCase)
	validate := validator.New()
	log := logger.NewNopLogger()

	handler := handlers.NewTransactionHandler(mockUseCase, validate, log)

	lambdaReq := events.APIGatewayV2HTTPRequest{
		Body: "invalid json",
	}

	resp, err := handler.HandleLambda(context.Background(), lambdaReq)

	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestTransactionHandler_HandleLambda_ValidationFailed(t *testing.T) {
	mockUseCase := new(MockPublishUseCase)
	validate := validator.New()
	log := logger.NewNopLogger()

	handler := handlers.NewTransactionHandler(mockUseCase, validate, log)

	requestBody := dtos.PublishTransactionRequest{
		ChainType: "INVALID",
		// missing required fields
	}

	bodyJSON, _ := json.Marshal(requestBody)
	lambdaReq := events.APIGatewayV2HTTPRequest{
		Body: string(bodyJSON),
	}

	resp, err := handler.HandleLambda(context.Background(), lambdaReq)

	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestTransactionHandler_HandleLambda_UseCaseError(t *testing.T) {
	mockUseCase := new(MockPublishUseCase)
	validate := validator.New()
	log := logger.NewNopLogger()

	handler := handlers.NewTransactionHandler(mockUseCase, validate, log)

	requestBody := dtos.PublishTransactionRequest{
		ChainType:     "EVM",
		OperationType: "TRANSFER",
		Payload:       map[string]interface{}{},
	}

	mockUseCase.On("Execute", mock.Anything, mock.Anything).Return(nil, errors.New("use case error"))

	bodyJSON, _ := json.Marshal(requestBody)
	lambdaReq := events.APIGatewayV2HTTPRequest{
		Body: string(bodyJSON),
	}

	resp, err := handler.HandleLambda(context.Background(), lambdaReq)

	assert.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)
	mockUseCase.AssertExpectations(t)
}

func TestNewTransactionHandler(t *testing.T) {
	mockUseCase := new(MockPublishUseCase)
	validate := validator.New()
	log := logger.NewNopLogger()

	handler := handlers.NewTransactionHandler(mockUseCase, validate, log)

	assert.NotNil(t, handler)
}
