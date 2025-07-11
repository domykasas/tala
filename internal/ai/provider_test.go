package ai

import (
	"context"
	"strings"
	"testing"
)

func TestCreateProvider(t *testing.T) {
	tests := []struct {
		name         string
		providerType string
		apiKey       string
		model        string
		temperature  float64
		maxTokens    int
		expectError  bool
		expectedName string
	}{
		{
			name:         "OpenAI provider",
			providerType: "openai",
			apiKey:       "test-key",
			model:        "gpt-3.5-turbo",
			temperature:  0.7,
			maxTokens:    1000,
			expectError:  false,
			expectedName: "OpenAI",
		},
		{
			name:         "Anthropic provider",
			providerType: "anthropic",
			apiKey:       "test-key",
			model:        "claude-3-sonnet",
			temperature:  0.7,
			maxTokens:    1000,
			expectError:  false,
			expectedName: "Anthropic",
		},
		{
			name:         "Ollama provider",
			providerType: "ollama",
			apiKey:       "",
			model:        "llama2",
			temperature:  0.7,
			maxTokens:    1000,
			expectError:  false,
			expectedName: "Ollama",
		},
		{
			name:         "Unsupported provider",
			providerType: "unsupported",
			apiKey:       "test-key",
			model:        "test-model",
			temperature:  0.7,
			maxTokens:    1000,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := CreateProvider(tt.providerType, tt.apiKey, tt.model, tt.temperature, tt.maxTokens)
			
			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if provider.GetName() != tt.expectedName {
				t.Errorf("Expected provider name %s, got %s", tt.expectedName, provider.GetName())
			}
		})
	}
}

func TestOpenAIProvider(t *testing.T) {
	provider := NewOpenAIProvider("test-key", "gpt-3.5-turbo", 0.7, 1000)
	
	if provider.GetName() != "OpenAI" {
		t.Errorf("Expected provider name 'OpenAI', got %s", provider.GetName())
	}
	
	ctx := context.Background()
	response, err := provider.GenerateResponse(ctx, "test prompt")
	
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	if !strings.Contains(response, "test prompt") {
		t.Error("Response should contain the original prompt")
	}
}

func TestAnthropicProvider(t *testing.T) {
	provider := NewAnthropicProvider("test-key", "claude-3-sonnet", 0.7, 1000)
	
	if provider.GetName() != "Anthropic" {
		t.Errorf("Expected provider name 'Anthropic', got %s", provider.GetName())
	}
	
	ctx := context.Background()
	response, err := provider.GenerateResponse(ctx, "test prompt")
	
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	if !strings.Contains(response, "test prompt") {
		t.Error("Response should contain the original prompt")
	}
}

func TestOllamaProvider(t *testing.T) {
	provider := NewOllamaProvider("llama2", 0.7, 1000, "")
	
	if provider.GetName() != "Ollama" {
		t.Errorf("Expected provider name 'Ollama', got %s", provider.GetName())
	}
	
	if provider.BaseURL != "http://localhost:11434" {
		t.Errorf("Expected default base URL 'http://localhost:11434', got %s", provider.BaseURL)
	}
}

func TestOllamaProviderCustomURL(t *testing.T) {
	customURL := "http://custom:8080"
	provider := NewOllamaProvider("llama2", 0.7, 1000, customURL)
	
	if provider.BaseURL != customURL {
		t.Errorf("Expected custom base URL %s, got %s", customURL, provider.BaseURL)
	}
}