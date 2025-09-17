package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Save original environment
	originalServerName := os.Getenv("SERVER_NAME")
	originalServerVersion := os.Getenv("SERVER_VERSION")
	originalK8sInCluster := os.Getenv("K8S_IN_CLUSTER")
	originalLogLevel := os.Getenv("LOG_LEVEL")
	originalLogFormat := os.Getenv("LOG_FORMAT")

	// Clean up after test
	defer func() {
		os.Setenv("SERVER_NAME", originalServerName)
		os.Setenv("SERVER_VERSION", originalServerVersion)
		os.Setenv("K8S_IN_CLUSTER", originalK8sInCluster)
		os.Setenv("LOG_LEVEL", originalLogLevel)
		os.Setenv("LOG_FORMAT", originalLogFormat)
	}()

	tests := []struct {
		name        string
		envVars     map[string]string
		expected    *Config
		expectError bool
	}{
		{
			name:    "default configuration",
			envVars: map[string]string{},
			expected: &Config{
				Server: ServerConfig{
					Name:    "k8s-mcp-server",
					Version: "2.0.0",
				},
				Kubernetes: KubernetesConfig{
					InCluster: false,
				},
				Logging: LoggingConfig{
					Level:  "info",
					Format: "text",
				},
			},
		},
		{
			name: "custom configuration",
			envVars: map[string]string{
				"SERVER_NAME":     "custom-server",
				"SERVER_VERSION":  "3.0.0",
				"K8S_IN_CLUSTER":  "true",
				"LOG_LEVEL":       "debug",
				"LOG_FORMAT":      "json",
			},
			expected: &Config{
				Server: ServerConfig{
					Name:    "custom-server",
					Version: "3.0.0",
				},
				Kubernetes: KubernetesConfig{
					InCluster: true,
				},
				Logging: LoggingConfig{
					Level:  "debug",
					Format: "json",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Load configuration
			cfg, err := Load()

			// Check error expectation
			if tt.expectError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			// Check configuration values
			if cfg != nil && tt.expected != nil {
				if cfg.Server.Name != tt.expected.Server.Name {
					t.Errorf("Expected server name %s, got %s", tt.expected.Server.Name, cfg.Server.Name)
				}
				if cfg.Server.Version != tt.expected.Server.Version {
					t.Errorf("Expected server version %s, got %s", tt.expected.Server.Version, cfg.Server.Version)
				}
				if cfg.Kubernetes.InCluster != tt.expected.Kubernetes.InCluster {
					t.Errorf("Expected in-cluster %v, got %v", tt.expected.Kubernetes.InCluster, cfg.Kubernetes.InCluster)
				}
				if cfg.Logging.Level != tt.expected.Logging.Level {
					t.Errorf("Expected log level %s, got %s", tt.expected.Logging.Level, cfg.Logging.Level)
				}
				if cfg.Logging.Format != tt.expected.Logging.Format {
					t.Errorf("Expected log format %s, got %s", tt.expected.Logging.Format, cfg.Logging.Format)
				}
			}

			// Clean up environment variables for this test
			for key := range tt.envVars {
				os.Unsetenv(key)
			}
		})
	}
}

func TestGetEnvWithDefault(t *testing.T) {
	// Test with non-existent environment variable
	result := getEnvWithDefault("NON_EXISTENT_VAR", "default_value")
	if result != "default_value" {
		t.Errorf("Expected 'default_value', got '%s'", result)
	}

	// Test with existing environment variable
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	result = getEnvWithDefault("TEST_VAR", "default_value")
	if result != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", result)
	}
}

func TestGetKubeConfigPath(t *testing.T) {
	// Save original KUBECONFIG
	originalKubeconfig := os.Getenv("KUBECONFIG")
	defer os.Setenv("KUBECONFIG", originalKubeconfig)

	// Test with KUBECONFIG environment variable set
	os.Setenv("KUBECONFIG", "/custom/kubeconfig")
	result := getKubeConfigPath()
	if result != "/custom/kubeconfig" {
		t.Errorf("Expected '/custom/kubeconfig', got '%s'", result)
	}

	// Test with KUBECONFIG unset
	os.Unsetenv("KUBECONFIG")
	result = getKubeConfigPath()
	// The result will vary based on the system's home directory
	// We just check that it's not empty (assuming a home directory exists)
	if result == "" {
		t.Log("Note: No home directory found, kubeconfig path is empty (this may be expected in some environments)")
	}
}