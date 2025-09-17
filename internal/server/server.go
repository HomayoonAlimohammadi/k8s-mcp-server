package server

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/HomayoonAlimohammadi/k8s-mcp-server/internal/config"
	"github.com/HomayoonAlimohammadi/k8s-mcp-server/internal/handlers"
	"github.com/HomayoonAlimohammadi/k8s-mcp-server/internal/k8s"
	"github.com/mark3labs/mcp-go/server"
)

// Server represents the main MCP server
type Server struct {
	config       *config.Config
	logger       *slog.Logger
	mcpServer    *server.MCPServer
	toolHandlers *handlers.ToolHandlers
	stdioServer  *server.StdioServer
}

// New creates a new server instance
func New(cfg *config.Config, logger *slog.Logger) (*Server, error) {
	// Create Kubernetes client
	k8sClient, err := k8s.NewClient(cfg.Kubernetes.KubeConfig, cfg.Kubernetes.InCluster)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	// Create tool handlers
	toolHandlers := handlers.NewToolHandlers(k8sClient, logger)

	// Create MCP server with capabilities
	mcpServer := server.NewMCPServer(
		cfg.Server.Name,
		cfg.Server.Version,
		server.WithToolCapabilities(false),            // Tools list does not change
		server.WithResourceCapabilities(false, false), // No resource subscription or list changes
	)

	// Register all tools
	for _, tool := range toolHandlers.GetTools() {
		mcpServer.AddTool(tool, toolHandlers.HandleTool)
	}

	// Create stdio server
	stdioServer := server.NewStdioServer(mcpServer)

	return &Server{
		config:       cfg,
		logger:       logger,
		mcpServer:    mcpServer,
		toolHandlers: toolHandlers,
		stdioServer:  stdioServer,
	}, nil
}

// Run starts the server
func (s *Server) Run(ctx context.Context) error {
	s.logger.Info("Starting K8s MCP Server",
		"name", s.config.Server.Name,
		"version", s.config.Server.Version,
	)

	// Start the stdio server
	return s.stdioServer.Listen(ctx, os.Stdin, os.Stdout)
}
