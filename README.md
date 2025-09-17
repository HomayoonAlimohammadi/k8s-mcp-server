# Kubernetes MCP Server

A professional Model Context Protocol (MCP) server that provides Kubernetes cluster interaction capabilities. This server allows MCP clients to inspect and manage Kubernetes resources through a standardized interface.

## Features

- **Pod Management**: List pods, get detailed pod information, and retrieve pod logs
- **Service Management**: List and inspect Kubernetes services
- **Deployment Management**: Monitor and inspect Kubernetes deployments
- **Namespace Management**: List all available namespaces
- **Comprehensive Logging**: Structured logging with configurable levels
- **Error Handling**: Robust error handling with meaningful error messages
- **Testing**: Comprehensive unit tests with good coverage
- **Docker Support**: Containerized deployment with multi-stage builds

## Tools Available

The server provides the following MCP tools:

1. **list-pods** - List all pods in a namespace
2. **get-pod** - Get detailed information about a specific pod
3. **get-pod-logs** - Retrieve logs from a specific pod (with optional tail limit)
4. **list-services** - List all services in a namespace
5. **get-service** - Get detailed information about a specific service
6. **list-deployments** - List all deployments in a namespace
7. **get-deployment** - Get detailed information about a specific deployment
8. **list-namespaces** - List all namespaces in the cluster

## Quick Start

### Prerequisites

- Go 1.23.2 or later
- Access to a Kubernetes cluster
- kubectl configured or running inside a Kubernetes cluster

### Installation

```bash
# Clone the repository
git clone <repository-url>
cd k8s-mcp-server

# Build the server
make build

# Or run directly
make run
```

### Configuration

The server can be configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_NAME` | Name of the MCP server | `k8s-mcp-server` |
| `SERVER_VERSION` | Version of the server | `2.0.0` |
| `KUBECONFIG` | Path to kubeconfig file | `~/.kube/config` |
| `K8S_IN_CLUSTER` | Whether running inside K8s cluster | `false` |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` |
| `LOG_FORMAT` | Log format (text, json) | `text` |

### Usage

#### Running Locally

```bash
# Run with default configuration
./bin/k8s-mcp-server

# Run with custom configuration
SERVER_NAME=my-k8s-server LOG_LEVEL=debug ./bin/k8s-mcp-server
```

#### Running with Docker

```bash
# Build the Docker image
make docker-build

# Run the container
docker run -v ~/.kube:/root/.kube k8s-mcp-server:latest
```

#### Running in Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: k8s-mcp-server
spec:
  replicas: 1
  selector:
    matchLabels:
      app: k8s-mcp-server
  template:
    metadata:
      labels:
        app: k8s-mcp-server
    spec:
      serviceAccountName: k8s-mcp-server
      containers:
      - name: k8s-mcp-server
        image: k8s-mcp-server:latest
        env:
        - name: K8S_IN_CLUSTER
          value: "true"
        - name: LOG_LEVEL
          value: "info"
```

## Development

### Building

```bash
# Build the binary
make build

# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Run linter
make lint

# Clean build artifacts
make clean
```

### Project Structure

```
├── cmd/
│   └── server/          # Main application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── handlers/        # MCP tool handlers
│   ├── k8s/            # Kubernetes client wrapper
│   └── server/         # MCP server implementation
├── Dockerfile          # Multi-stage Docker build
├── Makefile           # Build automation
├── go.mod             # Go module definition
└── README.md          # This file
```

### Architecture

The server follows a clean architecture pattern:

- **cmd/server**: Application entry point and initialization
- **internal/config**: Configuration loading and validation
- **internal/k8s**: Kubernetes API client abstraction
- **internal/handlers**: MCP tool implementations
- **internal/server**: MCP server setup and lifecycle

### Testing

The project includes comprehensive unit tests:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
make test-coverage
```

### Adding New Tools

1. Define the tool schema in `internal/handlers/tools.go`
2. Implement the handler function
3. Add the tool to the `GetTools()` method
4. Add a case to the `HandleTool()` method
5. Write unit tests

## MCP Integration

This server implements the Model Context Protocol (MCP) specification and can be used with any MCP-compatible client. It communicates via JSON-RPC over stdio.

### Example MCP Client Configuration

```json
{
  "mcpServers": {
    "k8s": {
      "command": "/path/to/k8s-mcp-server",
      "args": [],
      "env": {
        "LOG_LEVEL": "info"
      }
    }
  }
}
```

## Security Considerations

- The server requires appropriate Kubernetes RBAC permissions
- When running in-cluster, use a service account with minimal required permissions
- Logs may contain sensitive cluster information - configure log levels appropriately
- The server does not modify cluster state, only reads information

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite and linter
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Changelog

### v2.0.0

- Complete rewrite with modern Go best practices
- Updated to latest Kubernetes client libraries (v0.31.1)
- Improved error handling and logging
- Added comprehensive test suite
- Added Docker support
- Improved configuration management
- Added Makefile for build automation