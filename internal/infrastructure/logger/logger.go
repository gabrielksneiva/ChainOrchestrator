package logger

import (
	"go.uber.org/zap"
)

// NewLogger cria uma nova instância do logger
func NewLogger(env string) (*zap.Logger, error) {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	logger, err := config.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	if err != nil {
		return nil, err
	}

	return logger, nil
}

// NewNopLogger cria um logger nop para testes
func NewNopLogger() *zap.Logger {
	return zap.NewNop()
}
