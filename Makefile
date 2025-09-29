# K8s MCP Server Makefile

# Variables
BINARY_NAME=k8s-mcp-server
BINARY_PATH=./bin/$(BINARY_NAME)
MAIN_PATH=./cmd/server
DOCKER_IMAGE=k8s-mcp-server
DOCKER_TAG=latest

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt

# Build flags
LDFLAGS=-ldflags "-X main.version=$(shell git describe --tags --always --dirty)"

.PHONY: all build clean test deps fmt lint run docker-build docker-run help

# Default target
all: clean deps fmt lint test build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH) $(MAIN_PATH)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf bin/

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .

# Lint code (requires golangci-lint to be installed)
lint:
	@echo "Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, skipping lint check"; \
		echo "Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	$(BINARY_PATH)

# Run in development mode with hot reload (requires air to be installed)
dev:
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "air not found, falling back to regular run"; \
		echo "Install air with: go install github.com/cosmtrek/air@latest"; \
		$(MAKE) run; \
	fi

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run --rm -p 8080:8080 \
		-v ~/.kube/config:/root/.kube/config:ro \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

# Install development tools
install-tools:
	@echo "Installing development tools..."
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOGET) github.com/cosmtrek/air@latest

# Check for vulnerabilities
security:
	@echo "Checking for security vulnerabilities..."
	@if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./...; \
	else \
		echo "govulncheck not found, installing..."; \
		$(GOGET) golang.org/x/vuln/cmd/govulncheck@latest; \
		govulncheck ./...; \
	fi

# Generate code coverage report
coverage:
	@echo "Generating coverage report..."
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Show help
help:
	@echo "Available targets:"
	@echo "  all          - Run clean, deps, fmt, lint, test, and build"
	@echo "  build        - Build the binary"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  deps         - Download and tidy dependencies"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code (requires golangci-lint)"
	@echo "  run          - Build and run the application"
	@echo "  dev          - Run with hot reload (requires air)"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  install-tools- Install development tools"
	@echo "  security     - Check for security vulnerabilities"
	@echo "  coverage     - Generate test coverage report"
	@echo "  help         - Show this help message"