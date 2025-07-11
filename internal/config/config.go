package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	APIKey       string `json:"api_key"`
	Provider     string `json:"provider"`
	Model        string `json:"model"`
	Temperature  float64 `json:"temperature"`
	MaxTokens    int    `json:"max_tokens"`
	SystemPrompt string `json:"system_prompt"`
}

// Getter methods for provider creation
func (c *Config) GetProvider() string {
	return c.Provider
}

func (c *Config) GetAPIKey() string {
	return c.APIKey
}

func (c *Config) GetModel() string {
	return c.Model
}

func (c *Config) GetTemperature() float64 {
	return c.Temperature
}

func (c *Config) GetMaxTokens() int {
	return c.MaxTokens
}

func DefaultConfig() *Config {
	return &Config{
		Provider:     "ollama",
		Model:        "deepseek-r1",
		Temperature:  0.7,
		MaxTokens:    0, // 0 means no token limit
		SystemPrompt: "You are a helpful AI assistant.",
	}
}

func Load() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := DefaultConfig()
		if err := config.Save(); err != nil {
			return nil, err
		}
		return config, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) Save() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0750); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}

var getConfigPath = func() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".config", "tala", "config.json"), nil
}

func (c *Config) Validate() error {
	if c.Provider != "ollama" && c.APIKey == "" {
		return fmt.Errorf("API key is required for provider: %s", c.Provider)
	}
	if c.Provider == "" {
		return fmt.Errorf("provider is required")
	}
	if c.Model == "" {
		return fmt.Errorf("model is required")
	}
	return nil
}