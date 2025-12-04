package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/gabrielksneiva/ChainOrchestrator/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_Defaults(t *testing.T) {
	// Clear all env vars
	os.Clearenv()

	cfg := config.LoadConfig()

	assert.NotNil(t, cfg)
	assert.Equal(t, "development", cfg.Environment)
	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "us-east-1", cfg.AWSRegion)
	assert.Equal(t, time.Duration(30)*time.Second, cfg.RequestTimeout)
}

func TestLoadConfig_CustomValues(t *testing.T) {
	os.Setenv("ENVIRONMENT", "production")
	os.Setenv("PORT", "3000")
	os.Setenv("SNS_TOPIC_ARN", "arn:aws:sns:us-east-1:123:Topic")
	os.Setenv("AWS_REGION", "us-west-2")
	os.Setenv("REQUEST_TIMEOUT_SECONDS", "60")
	os.Setenv("REDIS_ADDR", "redis:6379")
	os.Setenv("REDIS_PASSWORD", "secret")
	os.Setenv("REDIS_DB", "1")

	defer os.Clearenv()

	cfg := config.LoadConfig()

	assert.Equal(t, "production", cfg.Environment)
	assert.Equal(t, "3000", cfg.Port)
	assert.Equal(t, "arn:aws:sns:us-east-1:123:Topic", cfg.SNSTopicARN)
	assert.Equal(t, "us-west-2", cfg.AWSRegion)
	assert.Equal(t, time.Duration(60)*time.Second, cfg.RequestTimeout)
	assert.Equal(t, "redis:6379", cfg.RedisAddr)
	assert.Equal(t, "secret", cfg.RedisPassword)
	assert.Equal(t, 1, cfg.RedisDB)
}

func TestLoadConfig_InvalidTimeout(t *testing.T) {
	os.Setenv("REQUEST_TIMEOUT_SECONDS", "invalid")
	defer os.Clearenv()

	cfg := config.LoadConfig()

	// Should use default value (0) when parsing fails
	assert.Equal(t, time.Duration(0)*time.Second, cfg.RequestTimeout)
}

func TestLoadConfig_InvalidRedisDB(t *testing.T) {
	os.Setenv("REDIS_DB", "invalid")
	defer os.Clearenv()

	cfg := config.LoadConfig()

	// Should use default value (0) when parsing fails
	assert.Equal(t, 0, cfg.RedisDB)
}
