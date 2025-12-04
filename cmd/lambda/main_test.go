package main

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandler_HealthRoute(t *testing.T) {
	request := events.APIGatewayV2HTTPRequest{
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "GET",
				Path:   "/health",
			},
			RequestID: "test-request-id",
		},
		RawPath: "/health",
	}

	response, err := handler(context.Background(), request)

	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)
	assert.Contains(t, response.Body, "healthy")
}

func TestHandleHealthCheck(t *testing.T) {
	response, err := handleHealthCheck(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)
	assert.Contains(t, response.Body, "healthy")
	assert.Contains(t, response.Body, "ChainOrchestrator")
}

func TestHandler_PostTransactionRoute_ValidationError(t *testing.T) {
	reqBody := map[string]interface{}{
		"chain_type":     "EVM",
		"operation_type": "DEPOSIT",
		"amount":         100.5,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	request := events.APIGatewayV2HTTPRequest{
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "POST",
				Path:   "/transaction",
			},
			RequestID: "test-request-id",
		},
		RawPath: "/transaction",
		Body:    string(bodyBytes),
	}

	response, err := handler(context.Background(), request)

	assert.NoError(t, err)
	// Espera 400 porque está faltando o campo payload
	assert.Equal(t, 400, response.StatusCode)
}

func TestHandler_WalletBalanceRoute(t *testing.T) {
	request := events.APIGatewayV2HTTPRequest{
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "GET",
				Path:   "/walletbalance",
			},
			RequestID: "test-request-id",
		},
		RawPath: "/walletbalance",
	}

	response, err := handler(context.Background(), request)

	assert.NoError(t, err)
	assert.Equal(t, 501, response.StatusCode)
	assert.Contains(t, response.Body, "not implemented")
}

func TestHandleWalletBalance(t *testing.T) {
	request := events.APIGatewayV2HTTPRequest{}
	response, err := handleWalletBalance(context.Background(), request)

	assert.NoError(t, err)
	assert.Equal(t, 501, response.StatusCode)
}

func TestHandler_TransactionStatusRoute(t *testing.T) {
	request := events.APIGatewayV2HTTPRequest{
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "GET",
				Path:   "/transaction-status",
			},
			RequestID: "test-request-id",
		},
		RawPath: "/transaction-status",
	}

	response, err := handler(context.Background(), request)

	assert.NoError(t, err)
	assert.Equal(t, 501, response.StatusCode)
	assert.Contains(t, response.Body, "not implemented")
}

func TestHandleTransactionStatus(t *testing.T) {
	request := events.APIGatewayV2HTTPRequest{}
	response, err := handleTransactionStatus(context.Background(), request)

	assert.NoError(t, err)
	assert.Equal(t, 501, response.StatusCode)
}

func TestHandler_UnknownRoute(t *testing.T) {
	request := events.APIGatewayV2HTTPRequest{
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "GET",
				Path:   "/unknown",
			},
			RequestID: "test-request-id",
		},
		RawPath: "/unknown",
	}

	response, err := handler(context.Background(), request)

	assert.NoError(t, err)
	assert.Equal(t, 404, response.StatusCode)
	assert.Contains(t, response.Body, "route not found")
}
