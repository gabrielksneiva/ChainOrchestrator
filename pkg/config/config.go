package config

import (
	"os"
	"strconv"
	"time"
)

// Config configuração da aplicação
type Config struct {
	Environment string
	Port        string

	// AWS SNS
	SNSTopicARN string
	AWSRegion   string

	// Redis (opcional - para cache de consultas)
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	// Timeouts
	RequestTimeout time.Duration
}

// LoadConfig carrega configuração a partir de variáveis de ambiente
func LoadConfig() *Config {
	requestTimeout, _ := strconv.Atoi(getEnv("REQUEST_TIMEOUT_SECONDS", "30"))
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8080"),

		SNSTopicARN: getEnv("SNS_TOPIC_ARN", ""),
		AWSRegion:   getEnv("AWS_REGION", "us-east-1"),

		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,

		RequestTimeout: time.Duration(requestTimeout) * time.Second,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
