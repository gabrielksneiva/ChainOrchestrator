# Dockerfile para build do binário Lambda
FROM public.ecr.aws/docker/library/golang:1.24-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build Lambda binary for arm64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build \
    -ldflags="-s -w" \
    -o bootstrap \
    cmd/lambda/main.go

# Final stage - minimal runtime
FROM scratch

COPY --from=builder /build/bootstrap /bootstrap

# Lambda expects the handler to be called "bootstrap"
ENTRYPOINT ["/bootstrap"]
