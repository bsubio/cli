package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// Config represents the CLI configuration
type Config struct {
	APIKey  string `json:"api_key"`
	BaseURL string `json:"base_url"`
}

// getConfigPath returns the path to the config file
func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "bsubio")
	return filepath.Join(configDir, "config.json"), nil
}

// loadConfig loads the configuration from disk
func loadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found. Run 'bsubio config' to set up")
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// saveConfig saves the configuration to disk
func saveConfig(config *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write config file with restricted permissions (user read/write only)
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// runConfig implements the config command
func runConfig(args []string) error {
	reader := bufio.NewReader(os.Stdin)

	// Get API key
	fmt.Print("Enter your BSUB.IO API key: ")
	apiKeyBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("failed to read API key: %w", err)
	}
	fmt.Println() // Print newline after password input
	apiKey := strings.TrimSpace(string(apiKeyBytes))

	if apiKey == "" {
		return fmt.Errorf("API key cannot be empty")
	}

	// Get base URL (optional)
	fmt.Print("Enter base URL [https://app.bsub.io]: ")
	baseURL, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read base URL: %w", err)
	}
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		baseURL = "https://app.bsub.io"
	}

	// Save configuration
	config := &Config{
		APIKey:  apiKey,
		BaseURL: baseURL,
	}

	if err := saveConfig(config); err != nil {
		return err
	}

	configPath, _ := getConfigPath()
	fmt.Printf("Configuration saved to %s\n", configPath)

	return nil
}
