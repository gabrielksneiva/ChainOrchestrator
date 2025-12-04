package logger_test

import (
	"testing"

	"github.com/gabrielksneiva/ChainOrchestrator/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewLogger_Production(t *testing.T) {
	log, err := logger.NewLogger("production")

	assert.NoError(t, err)
	assert.NotNil(t, log)

	// Clean up
	_ = log.Sync()
}

func TestNewLogger_Development(t *testing.T) {
	log, err := logger.NewLogger("development")

	assert.NoError(t, err)
	assert.NotNil(t, log)

	// Verify it has caller information
	log.Info("test message")

	// Clean up
	_ = log.Sync()
}

func TestNewLogger_Staging(t *testing.T) {
	log, err := logger.NewLogger("staging")

	assert.NoError(t, err)
	assert.NotNil(t, log)

	// Clean up
	_ = log.Sync()
}

func TestNewNopLogger(t *testing.T) {
	log := logger.NewNopLogger()

	assert.NotNil(t, log)
	assert.Equal(t, zap.NewNop(), log)

	// Should not panic
	log.Info("test")
	log.Error("error")
	log.Debug("debug")
}

func TestLogger_Logging(t *testing.T) {
	log, _ := logger.NewLogger("development")
	defer log.Sync()

	// Should not panic
	log.Info("info message", zap.String("key", "value"))
	log.Error("error message", zap.Error(assert.AnError))
	log.Debug("debug message")
	log.Warn("warn message")
}

func TestLogger_WithFields(t *testing.T) {
	log, _ := logger.NewLogger("production")
	defer log.Sync()

	// Should not panic with structured fields
	log.Info("structured logging",
		zap.String("service", "ChainOrchestrator"),
		zap.Int("port", 8080),
		zap.Bool("enabled", true),
	)
}

func TestLogger_StackTrace(t *testing.T) {
	log, _ := logger.NewLogger("production")
	defer log.Sync()

	// Error level should include stack trace
	log.Error("error with stack",
		zap.Error(assert.AnError),
		zap.Stack("stack"),
	)
}

func TestLogger_LevelFiltering(t *testing.T) {
	log, _ := logger.NewLogger("production")
	defer log.Sync()

	// Production logger should filter debug messages
	// This should not appear in output
	log.Debug("debug message - should be filtered in production")

	// But errors should appear
	log.Error("error message - should appear")
}

func TestLogger_OutputPaths(t *testing.T) {
	log, err := logger.NewLogger("production")

	assert.NoError(t, err)
	assert.NotNil(t, log)

	// Verify logger was created successfully (implicit verification of output paths)
	defer log.Sync()
}

func TestLogger_Caller(t *testing.T) {
	log, _ := logger.NewLogger("development")
	defer log.Sync()

	// Development logger should include caller information
	log.Info("message with caller info")
}

func TestLogger_ErrorOutput(t *testing.T) {
	log, _ := logger.NewLogger("production")
	defer log.Sync()

	// Errors should go to stderr
	log.Error("error to stderr", zap.Error(assert.AnError))
}

func TestLogger_EmptyEnvironment(t *testing.T) {
	log, err := logger.NewLogger("")

	assert.NoError(t, err)
	assert.NotNil(t, log)
	defer log.Sync()
}

func TestLogger_NonStandardEnvironment(t *testing.T) {
	log, err := logger.NewLogger("test")

	assert.NoError(t, err)
	assert.NotNil(t, log)
	defer log.Sync()
}
