package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	APIKey       string            `json:"api_key"`
	Provider     string            `json:"provider"`
	Model        string            `json:"model"`
	Temperature  float64           `json:"temperature"`
	MaxTokens    int               `json:"max_tokens"`
	SystemPrompt string            `json:"system_prompt"`
	
	// Global settings
	EnableStreaming bool              `json:"enable_streaming"`
	DefaultMode     string            `json:"default_mode"` // "tui", "gui", "headless"
	CustomPrompts   map[string]string `json:"custom_prompts"`
	Aliases         map[string]string `json:"aliases"`
	
	// UI preferences
	ShowTimestamps  bool   `json:"show_timestamps"`
	ShowTokens      bool   `json:"show_tokens"`
	CompactMode     bool   `json:"compact_mode"`
	Theme           string `json:"theme"` // "default", "minimal", "colorful"
	
	// Session settings
	SaveHistory     bool   `json:"save_history"`
	HistoryLimit    int    `json:"history_limit"`
	AutoSave        bool   `json:"auto_save"`
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
		Model:        "llama3.2:1b", // Use faster model by default
		Temperature:  0.7,
		MaxTokens:    0, // 0 means no token limit
		SystemPrompt: "You are a helpful AI assistant.",
		
		// Global settings
		EnableStreaming: true,
		DefaultMode:     "tui",
		CustomPrompts:   make(map[string]string),
		Aliases:         make(map[string]string),
		
		// UI preferences
		ShowTimestamps:  false,
		ShowTokens:      true,
		CompactMode:     false,
		Theme:           "default",
		
		// Session settings
		SaveHistory:     true,
		HistoryLimit:    1000,
		AutoSave:        true,
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

// Custom prompt management
func (c *Config) AddCustomPrompt(name, prompt string) {
	if c.CustomPrompts == nil {
		c.CustomPrompts = make(map[string]string)
	}
	c.CustomPrompts[name] = prompt
}

func (c *Config) GetCustomPrompt(name string) (string, bool) {
	if c.CustomPrompts == nil {
		return "", false
	}
	prompt, exists := c.CustomPrompts[name]
	return prompt, exists
}

func (c *Config) RemoveCustomPrompt(name string) {
	if c.CustomPrompts != nil {
		delete(c.CustomPrompts, name)
	}
}

func (c *Config) ListCustomPrompts() []string {
	if c.CustomPrompts == nil {
		return []string{}
	}
	var names []string
	for name := range c.CustomPrompts {
		names = append(names, name)
	}
	return names
}

// Alias management
func (c *Config) AddAlias(alias, command string) {
	if c.Aliases == nil {
		c.Aliases = make(map[string]string)
	}
	c.Aliases[alias] = command
}

func (c *Config) GetAlias(alias string) (string, bool) {
	if c.Aliases == nil {
		return "", false
	}
	command, exists := c.Aliases[alias]
	return command, exists
}

func (c *Config) RemoveAlias(alias string) {
	if c.Aliases != nil {
		delete(c.Aliases, alias)
	}
}

func (c *Config) ListAliases() []string {
	if c.Aliases == nil {
		return []string{}
	}
	var aliases []string
	for alias := range c.Aliases {
		aliases = append(aliases, alias)
	}
	return aliases
}

