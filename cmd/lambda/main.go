package main

// Lambda handler for transaction orchestration
// Updated: 2025-12-04 - Testing CI/CD pipeline with updated IAM

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/gabrielksneiva/ChainOrchestrator/internal/application/usecases"
	"github.com/gabrielksneiva/ChainOrchestrator/internal/infrastructure/eventbus"
	"github.com/gabrielksneiva/ChainOrchestrator/internal/infrastructure/logger"
	"github.com/gabrielksneiva/ChainOrchestrator/internal/interfaces/handlers"
	pkgconfig "github.com/gabrielksneiva/ChainOrchestrator/pkg/config"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var (
	cfg                *pkgconfig.Config
	log                *zap.Logger
	transactionHandler *handlers.TransactionHandler
)

func init() {
	var err error

	// Load config
	cfg = pkgconfig.LoadConfig()

	// Initialize logger
	log, err = logger.NewLogger(cfg.Environment)
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}

	// Initialize AWS config
	awsCfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal("failed to load AWS config", zap.Error(err))
	}

	// Initialize SNS client
	snsClient := sns.NewFromConfig(awsCfg)

	// Initialize SNS publisher
	snsPublisher := eventbus.NewSNSPublisher(snsClient, cfg.SNSTopicARN, log)

	// Initialize validator
	validate := validator.New()

	// Initialize use cases
	publishTxUseCase := usecases.NewPublishTransactionUseCase(snsPublisher, log)

	// Initialize handlers
	transactionHandler = handlers.NewTransactionHandler(publishTxUseCase, validate, log)

	log.Info("Lambda function initialized successfully",
		zap.String("environment", cfg.Environment),
		zap.String("sns_topic_arn", cfg.SNSTopicARN),
	)
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Info("received request",
		zap.String("method", request.RequestContext.HTTP.Method),
		zap.String("path", request.RequestContext.HTTP.Path),
		zap.String("request_id", request.RequestContext.RequestID),
	)

	// Route handling
	route := request.RequestContext.HTTP.Method + " " + request.RawPath

	switch route {
	case "GET /health":
		return handleHealthCheck(ctx)

	case "POST /transaction":
		return handlePostTransaction(ctx, request)

	case "GET /walletbalance":
		return handleWalletBalance(ctx, request)

	case "GET /transaction-status":
		return handleTransactionStatus(ctx, request)

	default:
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 404,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"error":"route not found"}`,
		}, nil
	}
}

func handleHealthCheck(ctx context.Context) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{"status":"healthy","service":"ChainOrchestrator"}`,
	}, nil
}

func handlePostTransaction(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// Use the transaction handler logic
	return transactionHandler.HandleLambda(ctx, request)
}

func handleWalletBalance(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Info("wallet balance endpoint called - not implemented yet")
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 501,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{"error":"not implemented","message":"wallet balance endpoint is not yet implemented"}`,
	}, nil
}

func handleTransactionStatus(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Info("transaction status endpoint called - not implemented yet")
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 501,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{"error":"not implemented","message":"transaction status endpoint is not yet implemented"}`,
	}, nil
}
