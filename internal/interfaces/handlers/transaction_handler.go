package handlers

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gabrielksneiva/ChainOrchestrator/internal/application/dtos"
	"github.com/gabrielksneiva/ChainOrchestrator/internal/application/usecases"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// TransactionHandler handler para operações de transação
type TransactionHandler struct {
	publishTxUseCase usecases.TransactionPublisher
	validator        *validator.Validate
	logger           *zap.Logger
}

// NewTransactionHandler cria um novo handler de transação
func NewTransactionHandler(
	publishTxUseCase usecases.TransactionPublisher,
	validator *validator.Validate,
	logger *zap.Logger,
) *TransactionHandler {
	return &TransactionHandler{
		publishTxUseCase: publishTxUseCase,
		validator:        validator,
		logger:           logger,
	}
}

// PostTransaction endpoint POST /transaction para Fiber
// Recebe requisição e publica no SNS Topic "Transactions"
func (h *TransactionHandler) PostTransaction(c *fiber.Ctx) error {
	var req dtos.PublishTransactionRequest

	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		h.logger.Error("validation failed", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "validation failed",
			"details": err.Error(),
		})
	}

	// Publicar no SNS (async-first)
	resp, err := h.publishTxUseCase.Execute(c.Context(), &req)
	if err != nil {
		h.logger.Error("failed to publish transaction", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to publish transaction",
		})
	}

	// Retorna resposta rápida (202 Accepted)
	return c.Status(fiber.StatusAccepted).JSON(resp)
}

// HandleLambda handler para AWS Lambda
func (h *TransactionHandler) HandleLambda(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var req dtos.PublishTransactionRequest

	// Parse request body
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		h.logger.Error("failed to parse request body", zap.Error(err))
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"error":"invalid request body"}`,
		}, nil
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		h.logger.Error("validation failed", zap.Error(err))
		errorBody, _ := json.Marshal(map[string]interface{}{
			"error":   "validation failed",
			"details": err.Error(),
		})
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: string(errorBody),
		}, nil
	}

	// Execute use case
	resp, err := h.publishTxUseCase.Execute(ctx, &req)
	if err != nil {
		h.logger.Error("failed to publish transaction", zap.Error(err))
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"error":"failed to publish transaction"}`,
		}, nil
	}

	// Return success response
	responseBody, _ := json.Marshal(resp)
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 202,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(responseBody),
	}, nil
}
