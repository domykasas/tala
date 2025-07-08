package config

import (
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	
	if cfg.Provider != "ollama" {
		t.Errorf("Expected provider 'ollama', got %s", cfg.Provider)
	}
	
	if cfg.Model != "deepseek-r1" {
		t.Errorf("Expected model 'deepseek-r1', got %s", cfg.Model)
	}
	
	if cfg.Temperature != 0.7 {
		t.Errorf("Expected temperature 0.7, got %f", cfg.Temperature)
	}
	
	if cfg.MaxTokens != 0 {
		t.Errorf("Expected max tokens 0, got %d", cfg.MaxTokens)
	}
}

func TestConfigSaveAndLoad(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")
	
	originalGetConfigPath := getConfigPath
	defer func() {
		getConfigPath = originalGetConfigPath
	}()
	
	getConfigPath = func() (string, error) {
		return configPath, nil
	}
	
	cfg := &Config{
		APIKey:      "test-key",
		Provider:    "test-provider",
		Model:       "test-model",
		Temperature: 0.5,
		MaxTokens:   500,
	}
	
	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}
	
	loadedCfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	if loadedCfg.APIKey != cfg.APIKey {
		t.Errorf("Expected API key %s, got %s", cfg.APIKey, loadedCfg.APIKey)
	}
	
	if loadedCfg.Provider != cfg.Provider {
		t.Errorf("Expected provider %s, got %s", cfg.Provider, loadedCfg.Provider)
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		hasErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				APIKey:   "test-key",
				Provider: "openai",
				Model:    "gpt-3.5-turbo",
			},
			hasErr: false,
		},
		{
			name: "missing API key",
			config: &Config{
				Provider: "openai",
				Model:    "gpt-3.5-turbo",
			},
			hasErr: true,
		},
		{
			name: "missing provider",
			config: &Config{
				APIKey: "test-key",
				Model:  "gpt-3.5-turbo",
			},
			hasErr: true,
		},
		{
			name: "missing model",
			config: &Config{
				APIKey:   "test-key",
				Provider: "openai",
			},
			hasErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.hasErr && err == nil {
				t.Error("Expected validation error, got nil")
			}
			if !tt.hasErr && err != nil {
				t.Errorf("Expected no validation error, got %v", err)
			}
		})
	}
}