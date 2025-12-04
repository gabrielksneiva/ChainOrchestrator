package errors_test

import (
	"errors"
	"testing"

	pkgerrors "github.com/gabrielksneiva/ChainOrchestrator/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestAppError_Error(t *testing.T) {
	err := pkgerrors.NewAppError("TEST_CODE", "test message", nil)

	assert.Equal(t, "TEST_CODE: test message", err.Error())
}

func TestAppError_ErrorWithWrappedError(t *testing.T) {
	wrappedErr := errors.New("wrapped error")
	err := pkgerrors.NewAppError("TEST_CODE", "test message", wrappedErr)

	assert.Contains(t, err.Error(), "TEST_CODE")
	assert.Contains(t, err.Error(), "test message")
	assert.Contains(t, err.Error(), "wrapped error")
}

func TestAppError_Fields(t *testing.T) {
	wrappedErr := errors.New("wrapped")
	err := pkgerrors.NewAppError("CODE", "message", wrappedErr)

	assert.Equal(t, "CODE", err.Code)
	assert.Equal(t, "message", err.Message)
	assert.Equal(t, wrappedErr, err.Err)
}

func TestPredefinedErrors(t *testing.T) {
	assert.Equal(t, "INVALID_INPUT", pkgerrors.ErrInvalidInput.Code)
	assert.Equal(t, "VALIDATION_FAILED", pkgerrors.ErrValidationFailed.Code)
	assert.Equal(t, "PUBLISH_FAILED", pkgerrors.ErrPublishFailed.Code)
	assert.Equal(t, "NOT_IMPLEMENTED", pkgerrors.ErrNotImplemented.Code)
}
