package http_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gabrielksneiva/ChainOrchestrator/internal/application/dtos"
	"github.com/gabrielksneiva/ChainOrchestrator/internal/infrastructure/http"
	"github.com/gabrielksneiva/ChainOrchestrator/internal/infrastructure/logger"
	"github.com/gabrielksneiva/ChainOrchestrator/internal/interfaces/handlers"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTransactionPublisher struct {
	mock.Mock
}

func (m *MockTransactionPublisher) Execute(ctx context.Context, req *dtos.PublishTransactionRequest) (*dtos.PublishTransactionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dtos.PublishTransactionResponse), args.Error(1)
}

func TestNewRouter(t *testing.T) {
	mockPublisher := new(MockTransactionPublisher)
	log := logger.NewNopLogger()
	validate := validator.New()

	// Create a mock handler - in real scenario this would use the actual handler
	// but for this test we just need to verify the router is created
	handler := handlers.NewTransactionHandler(mockPublisher, validate, log)

	router := http.NewRouter(handler, log)

	assert.NotNil(t, router)
	assert.NotNil(t, router.GetApp())
}

func TestRouter_Setup(t *testing.T) {
	mockPublisher := new(MockTransactionPublisher)
	log := logger.NewNopLogger()
	validate := validator.New()

	handler := handlers.NewTransactionHandler(mockPublisher, validate, log)
	router := http.NewRouter(handler, log)

	router.Setup()
	app := router.GetApp()

	assert.NotNil(t, app)
}

func TestRouter_HealthEndpoint(t *testing.T) {
	mockPublisher := new(MockTransactionPublisher)
	log := logger.NewNopLogger()
	validate := validator.New()

	handler := handlers.NewTransactionHandler(mockPublisher, validate, log)
	router := http.NewRouter(handler, log)
	router.Setup()

	app := router.GetApp()

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestRouter_GetApp(t *testing.T) {
	mockPublisher := new(MockTransactionPublisher)
	log := logger.NewNopLogger()
	validate := validator.New()

	handler := handlers.NewTransactionHandler(mockPublisher, validate, log)
	router := http.NewRouter(handler, log)

	app1 := router.GetApp()
	app2 := router.GetApp()

	assert.NotNil(t, app1)
	assert.NotNil(t, app2)
	assert.Equal(t, app1, app2) // Should return the same instance
}
