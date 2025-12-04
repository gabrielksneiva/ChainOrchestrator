package usecases

import (
	"context"

	"github.com/gabrielksneiva/ChainOrchestrator/internal/application/dtos"
)

// TransactionPublisher interface para publicar transações
type TransactionPublisher interface {
	Execute(ctx context.Context, req *dtos.PublishTransactionRequest) (*dtos.PublishTransactionResponse, error)
}
