package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/HomayoonAlimohammadi/k8s-mcp-server/internal/config"
	"github.com/HomayoonAlimohammadi/k8s-mcp-server/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	var logLevel *slog.Level
	if err := logLevel.UnmarshalText([]byte(cfg.Logging.Level)); err != nil {
		log.Fatalf("Failed to parse log level: %v\n", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))

	srv, err := server.New(cfg, logger)
	if err != nil {
		logger.Error("Failed to create server", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		logger.Info("Received shutdown signal")
		cancel()
	}()

	if err := srv.Run(ctx); err != nil {
		logger.Error("Server failed", "error", err)
		os.Exit(1)
	}

	logger.Info("Server shutdown gracefully")
}
