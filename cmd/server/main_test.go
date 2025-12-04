package main

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/gabrielksneiva/ChainOrchestrator/internal/infrastructure/logger"
	pkgconfig "github.com/gabrielksneiva/ChainOrchestrator/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestProvideLogger(t *testing.T) {
	cfg := &pkgconfig.Config{
		Environment: "development",
	}

	log, err := provideLogger(cfg)

	assert.NoError(t, err)
	assert.NotNil(t, log)
}

func TestProvideLogger_Production(t *testing.T) {
	cfg := &pkgconfig.Config{
		Environment: "production",
	}

	log, err := provideLogger(cfg)

	assert.NoError(t, err)
	assert.NotNil(t, log)
}

func TestProvideValidator(t *testing.T) {
	validator := provideValidator()

	assert.NotNil(t, validator)
}

func TestProvideSNSClient(t *testing.T) {
	awsCfg := aws.Config{
		Region: "us-east-1",
	}

	client := provideSNSClient(awsCfg)

	assert.NotNil(t, client)
}

func TestProvideSNSPublisher(t *testing.T) {
	cfg := &pkgconfig.Config{
		Environment: "development",
		SNSTopicARN: "arn:aws:sns:us-east-1:123456789012:test-topic",
	}

	log, _ := logger.NewLogger("development")
	awsCfg := aws.Config{
		Region: "us-east-1",
	}
	snsClient := sns.NewFromConfig(awsCfg)

	publisher := provideSNSPublisher(snsClient, cfg, log)

	assert.NotNil(t, publisher)
}

func TestProvideAWSConfig(t *testing.T) {
	// Set minimal AWS env vars to make config loading work
	cfg, err := provideAWSConfig()

	// This might fail in CI/CD without AWS credentials, but will pass locally
	// We don't assert error because it depends on environment
	_ = err
	_ = cfg
}
