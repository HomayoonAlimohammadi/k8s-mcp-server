package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/HomayoonAlimohammadi/k8s-mcp-server/internal/k8s"
	"github.com/mark3labs/mcp-go/mcp"
)

// ToolHandlers handles MCP tool calls
type ToolHandlers struct {
	k8sClient *k8s.Client
	logger    *slog.Logger
}

// NewToolHandlers creates a new ToolHandlers instance
func NewToolHandlers(k8sClient *k8s.Client, logger *slog.Logger) *ToolHandlers {
	return &ToolHandlers{
		k8sClient: k8sClient,
		logger:    logger,
	}
}

// GetTools returns the list of available tools
func (h *ToolHandlers) GetTools() []mcp.Tool {
	return []mcp.Tool{
		{
			Name:        "list-pods",
			Description: "List pods in a namespace",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]any{
					"namespace": map[string]any{
						"type":        "string",
						"description": "Kubernetes namespace (default: default)",
					},
				},
			},
		},
		{
			Name:        "get-pod",
			Description: "Get detailed information about a specific pod",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]any{
					"name": map[string]any{
						"type":        "string",
						"description": "Pod name",
					},
					"namespace": map[string]any{
						"type":        "string",
						"description": "Kubernetes namespace (default: default)",
					},
				},
				Required: []string{"name"},
			},
		},
		{
			Name:        "get-pod-logs",
			Description: "Get logs from a specific pod",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]any{
					"name": map[string]any{
						"type":        "string",
						"description": "Pod name",
					},
					"namespace": map[string]any{
						"type":        "string",
						"description": "Kubernetes namespace (default: default)",
					},
					"tail": map[string]any{
						"type":        "integer",
						"description": "Number of lines to tail from the end of the log",
					},
				},
				Required: []string{"name"},
			},
		},
		{
			Name:        "list-services",
			Description: "List services in a namespace",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]any{
					"namespace": map[string]any{
						"type":        "string",
						"description": "Kubernetes namespace (default: default)",
					},
				},
			},
		},
		{
			Name:        "get-service",
			Description: "Get detailed information about a specific service",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]any{
					"name": map[string]any{
						"type":        "string",
						"description": "Service name",
					},
					"namespace": map[string]any{
						"type":        "string",
						"description": "Kubernetes namespace (default: default)",
					},
				},
				Required: []string{"name"},
			},
		},
		{
			Name:        "list-deployments",
			Description: "List deployments in a namespace",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]any{
					"namespace": map[string]any{
						"type":        "string",
						"description": "Kubernetes namespace (default: default)",
					},
				},
			},
		},
		{
			Name:        "get-deployment",
			Description: "Get detailed information about a specific deployment",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]any{
					"name": map[string]any{
						"type":        "string",
						"description": "Deployment name",
					},
					"namespace": map[string]any{
						"type":        "string",
						"description": "Kubernetes namespace (default: default)",
					},
				},
				Required: []string{"name"},
			},
		},
		{
			Name:        "list-namespaces",
			Description: "List all namespaces in the cluster",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
			},
		},
	}
}

// HandleTool routes tool calls to appropriate handlers
func (h *ToolHandlers) HandleTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Convert arguments to map[string]interface{}
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		args = make(map[string]interface{})
	}

	switch request.Params.Name {
	case "list-pods":
		return h.ListPods(ctx, args)
	case "get-pod":
		return h.GetPod(ctx, args)
	case "get-pod-logs":
		return h.GetPodLogs(ctx, args)
	case "list-services":
		return h.ListServices(ctx, args)
	case "get-service":
		return h.GetService(ctx, args)
	case "list-deployments":
		return h.ListDeployments(ctx, args)
	case "get-deployment":
		return h.GetDeployment(ctx, args)
	case "list-namespaces":
		return h.ListNamespaces(ctx, args)
	default:
		return nil, fmt.Errorf("unknown tool: %s", request.Params.Name)
	}
}

// ListPods handles the list-pods tool
func (h *ToolHandlers) ListPods(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	namespace := h.getNamespace(args)
	
	h.logger.Debug("Listing pods", "namespace", namespace)
	
	pods, err := h.k8sClient.ListPods(ctx, namespace)
	if err != nil {
		h.logger.Error("Failed to list pods", "namespace", namespace, "error", err)
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	podsJSON, err := json.MarshalIndent(pods, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal pods: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(fmt.Sprintf("Pods in namespace '%s':\n%s", namespace, string(podsJSON))),
		},
	}, nil
}

// GetPod handles the get-pod tool
func (h *ToolHandlers) GetPod(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("pod name is required")
	}

	namespace := h.getNamespace(args)
	
	h.logger.Debug("Getting pod", "name", name, "namespace", namespace)

	pod, err := h.k8sClient.GetPod(ctx, namespace, name)
	if err != nil {
		h.logger.Error("Failed to get pod", "name", name, "namespace", namespace, "error", err)
		return nil, fmt.Errorf("failed to get pod: %w", err)
	}

	podJSON, err := json.MarshalIndent(pod, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal pod: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(fmt.Sprintf("Pod '%s' in namespace '%s':\n%s", name, namespace, string(podJSON))),
		},
	}, nil
}

// GetPodLogs handles the get-pod-logs tool
func (h *ToolHandlers) GetPodLogs(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("pod name is required")
	}

	namespace := h.getNamespace(args)
	
	var tailLines *int64
	if tail, ok := args["tail"]; ok {
		switch v := tail.(type) {
		case float64:
			lines := int64(v)
			tailLines = &lines
		case int:
			lines := int64(v)
			tailLines = &lines
		case string:
			if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
				tailLines = &parsed
			}
		}
	}
	
	h.logger.Debug("Getting pod logs", "name", name, "namespace", namespace, "tail", tailLines)

	logs, err := h.k8sClient.GetPodLogs(ctx, namespace, name, tailLines)
	if err != nil {
		h.logger.Error("Failed to get pod logs", "name", name, "namespace", namespace, "error", err)
		return nil, fmt.Errorf("failed to get pod logs: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(fmt.Sprintf("Logs for pod '%s' in namespace '%s':\n%s", name, namespace, logs)),
		},
	}, nil
}

// ListServices handles the list-services tool
func (h *ToolHandlers) ListServices(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	namespace := h.getNamespace(args)
	
	h.logger.Debug("Listing services", "namespace", namespace)

	services, err := h.k8sClient.ListServices(ctx, namespace)
	if err != nil {
		h.logger.Error("Failed to list services", "namespace", namespace, "error", err)
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	servicesJSON, err := json.MarshalIndent(services, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal services: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(fmt.Sprintf("Services in namespace '%s':\n%s", namespace, string(servicesJSON))),
		},
	}, nil
}

// GetService handles the get-service tool
func (h *ToolHandlers) GetService(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("service name is required")
	}

	namespace := h.getNamespace(args)
	
	h.logger.Debug("Getting service", "name", name, "namespace", namespace)

	service, err := h.k8sClient.GetService(ctx, namespace, name)
	if err != nil {
		h.logger.Error("Failed to get service", "name", name, "namespace", namespace, "error", err)
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	serviceJSON, err := json.MarshalIndent(service, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal service: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(fmt.Sprintf("Service '%s' in namespace '%s':\n%s", name, namespace, string(serviceJSON))),
		},
	}, nil
}

// ListDeployments handles the list-deployments tool
func (h *ToolHandlers) ListDeployments(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	namespace := h.getNamespace(args)
	
	h.logger.Debug("Listing deployments", "namespace", namespace)

	deployments, err := h.k8sClient.ListDeployments(ctx, namespace)
	if err != nil {
		h.logger.Error("Failed to list deployments", "namespace", namespace, "error", err)
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}

	deploymentsJSON, err := json.MarshalIndent(deployments, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal deployments: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(fmt.Sprintf("Deployments in namespace '%s':\n%s", namespace, string(deploymentsJSON))),
		},
	}, nil
}

// GetDeployment handles the get-deployment tool
func (h *ToolHandlers) GetDeployment(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("deployment name is required")
	}

	namespace := h.getNamespace(args)
	
	h.logger.Debug("Getting deployment", "name", name, "namespace", namespace)

	deployment, err := h.k8sClient.GetDeployment(ctx, namespace, name)
	if err != nil {
		h.logger.Error("Failed to get deployment", "name", name, "namespace", namespace, "error", err)
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	deploymentJSON, err := json.MarshalIndent(deployment, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal deployment: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(fmt.Sprintf("Deployment '%s' in namespace '%s':\n%s", name, namespace, string(deploymentJSON))),
		},
	}, nil
}

// ListNamespaces handles the list-namespaces tool
func (h *ToolHandlers) ListNamespaces(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	h.logger.Debug("Listing namespaces")

	namespaces, err := h.k8sClient.ListNamespaces(ctx)
	if err != nil {
		h.logger.Error("Failed to list namespaces", "error", err)
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	namespacesJSON, err := json.MarshalIndent(namespaces, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal namespaces: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(fmt.Sprintf("Namespaces:\n%s", string(namespacesJSON))),
		},
	}, nil
}

// getNamespace extracts the namespace from arguments, defaulting to "default"
func (h *ToolHandlers) getNamespace(args map[string]interface{}) string {
	if ns, ok := args["namespace"].(string); ok && ns != "" {
		return ns
	}
	return "default"
}