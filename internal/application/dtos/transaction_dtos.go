package dtos

// PublishTransactionRequest representa a requisição para publicar transação no SNS
type PublishTransactionRequest struct {
	ChainType     string                 `json:"chain_type" validate:"required,oneof=EVM BTC TRON SOL"`
	OperationType string                 `json:"operation_type" validate:"required"`
	Payload       map[string]interface{} `json:"payload" validate:"required"`
}

// PublishTransactionResponse resposta quando transação é publicada no SNS
type PublishTransactionResponse struct {
	OperationID string `json:"operation_id"`
	ChainType   string `json:"chain_type"`
	Status      string `json:"status"` // Sempre "QUEUED"
	Message     string `json:"message"`
	CreatedAt   string `json:"created_at"`
}
