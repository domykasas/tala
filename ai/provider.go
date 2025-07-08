package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Provider interface {
	GenerateResponse(ctx context.Context, prompt string) (string, error)
	GetName() string
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIProvider struct {
	APIKey      string
	Model       string
	Temperature float64
	MaxTokens   int
}

func NewOpenAIProvider(apiKey, model string, temperature float64, maxTokens int) *OpenAIProvider {
	return &OpenAIProvider{
		APIKey:      apiKey,
		Model:       model,
		Temperature: temperature,
		MaxTokens:   maxTokens,
	}
}

func (p *OpenAIProvider) GenerateResponse(ctx context.Context, prompt string) (string, error) {
	return fmt.Sprintf("OpenAI response to: %s", prompt), nil
}

func (p *OpenAIProvider) GetName() string {
	return "OpenAI"
}

type AnthropicProvider struct {
	APIKey      string
	Model       string
	Temperature float64
	MaxTokens   int
}

func NewAnthropicProvider(apiKey, model string, temperature float64, maxTokens int) *AnthropicProvider {
	return &AnthropicProvider{
		APIKey:      apiKey,
		Model:       model,
		Temperature: temperature,
		MaxTokens:   maxTokens,
	}
}

func (p *AnthropicProvider) GenerateResponse(ctx context.Context, prompt string) (string, error) {
	return fmt.Sprintf("Anthropic response to: %s", prompt), nil
}

func (p *AnthropicProvider) GetName() string {
	return "Anthropic"
}

type OllamaProvider struct {
	Model       string
	Temperature float64
	MaxTokens   int
	BaseURL     string
	client      *http.Client
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
	Error    string `json:"error,omitempty"`
}

func NewOllamaProvider(model string, temperature float64, maxTokens int, baseURL string) *OllamaProvider {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	return &OllamaProvider{
		Model:       model,
		Temperature: temperature,
		MaxTokens:   maxTokens,
		BaseURL:     baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (p *OllamaProvider) GenerateResponse(ctx context.Context, prompt string) (string, error) {
	reqBody := OllamaRequest{
		Model:  p.Model,
		Prompt: prompt,
		Stream: false,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.BaseURL+"/api/generate", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if ollamaResp.Error != "" {
		return "", fmt.Errorf("Ollama error: %s", ollamaResp.Error)
	}

	return ollamaResp.Response, nil
}

func (p *OllamaProvider) GetName() string {
	return "Ollama"
}

func CreateProvider(providerType, apiKey, model string, temperature float64, maxTokens int) (Provider, error) {
	switch providerType {
	case "openai":
		return NewOpenAIProvider(apiKey, model, temperature, maxTokens), nil
	case "anthropic":
		return NewAnthropicProvider(apiKey, model, temperature, maxTokens), nil
	case "ollama":
		return NewOllamaProvider(model, temperature, maxTokens, ""), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", providerType)
	}
}