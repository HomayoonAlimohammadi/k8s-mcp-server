# K8s MCP Server

A Model Context Protocol (MCP) server that provides Kubernetes cluster access tools for AI assistants.

## Overview

This MCP server allows AI assistants to interact with Kubernetes clusters by providing tools to:
- Get specific Kubernetes resources
- List resources with filtering capabilities

The server uses the Kubernetes Go client library and supports both in-cluster and external cluster configurations.

## Features

- **Resource Access**: Get and list Kubernetes resources by kind, name, and namespace
- **Filtering**: Support for label selectors when listing resources
- **Flexible Authentication**: Works with kubeconfig files or in-cluster service account
- **Structured Logging**: Configurable logging levels and formats
- **Graceful Shutdown**: Proper signal handling and shutdown procedures

## Installation

### Prerequisites

- Go 1.23.2 or later
- Access to a Kubernetes cluster
- Valid kubeconfig file (if running outside the cluster)

### Build from source

```bash
git clone https://github.com/HomayoonAlimohammadi/k8s-mcp-server
cd k8s-mcp-server
make build
```

## Configuration

The server is configured through environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_NAME` | `k8s-mcp-server` | Server name identifier |
| `SERVER_VERSION` | `1.0.0` | Server version |
| `KUBECONFIG` | `~/.kube/config` | Path to kubeconfig file |
| `K8S_IN_CLUSTER` | `false` | Use in-cluster configuration |
| `LOG_LEVEL` | `info` | Logging level (debug, info, warn, error) |
| `LOG_FORMAT` | `text` | Log format (text, json) |

## Usage

### Running the server

```bash
# Using default configuration
./bin/k8s-mcp-server

# With custom configuration
SERVER_NAME=my-k8s-server LOG_LEVEL=debug ./bin/k8s-mcp-server
```

### Available Tools

#### getResource
Get a specific Kubernetes resource.

**Parameters:**
- `kind` (required): Resource kind (e.g., "Pod", "Service", "Deployment")
- `name` (required): Resource name
- `namespace` (optional): Resource namespace

#### listResources
List Kubernetes resources with optional filtering.

**Parameters:**
- `kind` (required): Resource kind to list
- `namespace` (optional): Namespace to search in
- `labelSelector` (optional): Label selector for filtering

### Example Usage with MCP Client

```typescript
// Get a specific pod
await mcpClient.callTool("getResource", {
  kind: "Pod",
  name: "my-pod",
  namespace: "default"
});

// List all deployments in a namespace
await mcpClient.callTool("listResources", {
  kind: "Deployment",
  namespace: "production"
});

// List pods with label selector
await mcpClient.callTool("listResources", {
  kind: "Pod",
  labelSelector: "app=nginx"
});
```

## Development

### Running in development

```bash
make run
```

### Testing

```bash
make test
```

### Code formatting and linting

```bash
make fmt
make lint
```

## Docker

### Build image

```bash
make docker-build
```

### Run container

```bash
make docker-run
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linting
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
