package errors

import "fmt"

// AppError erro customizado da aplicação
type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewAppError cria um novo erro da aplicação
func NewAppError(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Erros comuns
var (
	ErrInvalidInput     = &AppError{Code: "INVALID_INPUT", Message: "invalid input"}
	ErrValidationFailed = &AppError{Code: "VALIDATION_FAILED", Message: "validation failed"}
	ErrPublishFailed    = &AppError{Code: "PUBLISH_FAILED", Message: "failed to publish event"}
	ErrNotImplemented   = &AppError{Code: "NOT_IMPLEMENTED", Message: "feature not implemented"}
)
