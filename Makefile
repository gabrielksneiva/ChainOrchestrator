.PHONY: help test coverage build run clean deploy

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test: ## Run all tests
	go test -v -race ./...

coverage: ## Run tests with coverage
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out

build: ## Build the application
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o bin/bootstrap cmd/lambda/main.go

run: ## Run the application locally
	go run cmd/server/main.go

clean: ## Clean build artifacts
	rm -rf bin/ dist/ coverage.out coverage.html lambda.zip

lambda-zip: build ## Build and package Lambda
	cd bin && zip ../lambda.zip bootstrap

deploy: lambda-zip ## Deploy to AWS Lambda using Terraform
	cd terraform && terraform apply -auto-approve

fmt: ## Format code
	go fmt ./...
	goimports -w .

lint: ## Run linter
	golangci-lint run ./...

deps: ## Download dependencies
	go mod download
	go mod verify

all: clean deps fmt lint test build ## Run all checks and build
