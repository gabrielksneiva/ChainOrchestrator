package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/gabrielksneiva/ChainOrchestrator/internal/application/usecases"
	"github.com/gabrielksneiva/ChainOrchestrator/internal/infrastructure/eventbus"
	"github.com/gabrielksneiva/ChainOrchestrator/internal/infrastructure/http"
	"github.com/gabrielksneiva/ChainOrchestrator/internal/infrastructure/logger"
	"github.com/gabrielksneiva/ChainOrchestrator/internal/interfaces/handlers"
	pkgconfig "github.com/gabrielksneiva/ChainOrchestrator/pkg/config"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	app := fx.New(
		fx.Provide(
			pkgconfig.LoadConfig,
			provideLogger,
			provideAWSConfig,
			provideSNSClient,
			provideSNSPublisher,
			provideValidator,
			usecases.NewPublishTransactionUseCase,
			handlers.NewTransactionHandler,
			http.NewRouter,
		),
		fx.Invoke(runServer),
	)
	app.Run()
}

func provideLogger(cfg *pkgconfig.Config) (*zap.Logger, error) {
	return logger.NewLogger(cfg.Environment)
}

func provideAWSConfig() (aws.Config, error) {
	return config.LoadDefaultConfig(context.Background())
}

func provideSNSClient(awsCfg aws.Config) *sns.Client {
	return sns.NewFromConfig(awsCfg)
}

func provideSNSPublisher(
	snsClient *sns.Client,
	cfg *pkgconfig.Config,
	logger *zap.Logger,
) *eventbus.SNSPublisher {
	return eventbus.NewSNSPublisher(snsClient, cfg.SNSTopicARN, logger)
}

func provideValidator() *validator.Validate {
	return validator.New()
}

func runServer(
	lc fx.Lifecycle,
	router *http.Router,
	cfg *pkgconfig.Config,
	logger *zap.Logger,
) {
	router.Setup()
	app := router.GetApp()

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("starting ChainOrchestrator server",
				zap.String("port", cfg.Port),
				zap.String("environment", cfg.Environment),
			)
			go func() {
				if err := app.Listen(fmt.Sprintf(":%s", cfg.Port)); err != nil {
					logger.Fatal("failed to start server", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("shutting down ChainOrchestrator server")
			return app.Shutdown()
		},
	})
}
