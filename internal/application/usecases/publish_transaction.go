package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gabrielksneiva/ChainOrchestrator/internal/application/dtos"
	"github.com/gabrielksneiva/ChainOrchestrator/internal/infrastructure/eventbus"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// PublishTransactionUseCase caso de uso para publicar transação no SNS
type PublishTransactionUseCase struct {
	snsPublisher eventbus.Publisher
	logger       *zap.Logger
}

// NewPublishTransactionUseCase cria uma nova instância do caso de uso
func NewPublishTransactionUseCase(
	snsPublisher eventbus.Publisher,
	logger *zap.Logger,
) *PublishTransactionUseCase {
	return &PublishTransactionUseCase{
		snsPublisher: snsPublisher,
		logger:       logger,
	}
}

// Execute executa o caso de uso de publicação de transação
func (uc *PublishTransactionUseCase) Execute(
	ctx context.Context,
	req *dtos.PublishTransactionRequest,
) (*dtos.PublishTransactionResponse, error) {
	uc.logger.Info("publishing transaction to SNS",
		zap.String("chain_type", req.ChainType),
		zap.String("operation_type", req.OperationType),
	)

	// Gerar ID para a operação
	operationID := uuid.New()
	createdAt := time.Now()

	// Criar mensagem para SNS
	message := map[string]interface{}{
		"operation_id":   operationID.String(),
		"chain_type":     req.ChainType,
		"operation_type": req.OperationType,
		"payload":        req.Payload,
		"created_at":     createdAt.Format(time.RFC3339),
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		uc.logger.Error("failed to marshal message", zap.Error(err))
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	// Publicar no SNS Topic "Transactions"
	if err := uc.snsPublisher.Publish(ctx, string(messageJSON), req.ChainType); err != nil {
		uc.logger.Error("failed to publish transaction to SNS", zap.Error(err))
		return nil, fmt.Errorf("failed to publish transaction: %w", err)
	}

	uc.logger.Info("transaction published to SNS successfully",
		zap.String("operation_id", operationID.String()),
		zap.String("chain_type", req.ChainType),
	)

	// Retorna resposta rápida (async-first)
	return &dtos.PublishTransactionResponse{
		OperationID: operationID.String(),
		ChainType:   req.ChainType,
		Status:      "QUEUED",
		Message:     "Transaction queued for processing",
		CreatedAt:   createdAt.Format(time.RFC3339),
	}, nil
}
