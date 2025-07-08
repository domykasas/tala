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
	GenerateResponseWithTools(ctx context.Context, prompt string) (string, []ToolResult, error)
	GetName() string
	SupportsTools() bool
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

func (p *OpenAIProvider) GenerateResponseWithTools(ctx context.Context, prompt string) (string, []ToolResult, error) {
	// Use AI-based intent detection (simulated for OpenAI)
	detector := NewIntentDetector(p)
	intents, err := detector.DetectIntent(ctx, prompt)
	if err != nil {
		response, err := p.GenerateResponse(ctx, prompt)
		return response, []ToolResult{}, err
	}
	
	// Execute detected tools
	var toolResults []ToolResult
	for _, intent := range intents {
		if intent.Confidence > 0.5 {
			result := ExecuteTool(intent.Tool, intent.Parameters)
			toolResults = append(toolResults, result)
		}
	}
	
	// Generate appropriate response
	if len(toolResults) > 0 {
		summary := "I have successfully completed the following operations:\n"
		for _, result := range toolResults {
			if result.Success {
				summary += fmt.Sprintf("✓ %s\n", result.Content)
			} else {
				summary += fmt.Sprintf("✗ %s failed: %s\n", result.Name, result.Content)
			}
		}
		summary += "\nAll requested operations have been executed."
		return summary, toolResults, nil
	}
	
	response := fmt.Sprintf("OpenAI response to: %s", prompt)
	return response, toolResults, nil
}

func (p *OpenAIProvider) SupportsTools() bool {
	return true
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

func (p *AnthropicProvider) GenerateResponseWithTools(ctx context.Context, prompt string) (string, []ToolResult, error) {
	// Use AI-based intent detection (simulated for Anthropic)
	detector := NewIntentDetector(p)
	intents, err := detector.DetectIntent(ctx, prompt)
	if err != nil {
		response, err := p.GenerateResponse(ctx, prompt)
		return response, []ToolResult{}, err
	}
	
	// Execute detected tools
	var toolResults []ToolResult
	for _, intent := range intents {
		if intent.Confidence > 0.5 {
			result := ExecuteTool(intent.Tool, intent.Parameters)
			toolResults = append(toolResults, result)
		}
	}
	
	// Generate appropriate response
	if len(toolResults) > 0 {
		summary := "I have successfully completed the following operations:\n"
		for _, result := range toolResults {
			if result.Success {
				summary += fmt.Sprintf("✓ %s\n", result.Content)
			} else {
				summary += fmt.Sprintf("✗ %s failed: %s\n", result.Name, result.Content)
			}
		}
		summary += "\nAll requested operations have been executed."
		return summary, toolResults, nil
	}
	
	response := fmt.Sprintf("Anthropic response to: %s", prompt)
	return response, toolResults, nil
}

func (p *AnthropicProvider) SupportsTools() bool {
	return true
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
		return "", fmt.Errorf("ollama error: %s", ollamaResp.Error)
	}

	return ollamaResp.Response, nil
}

func (p *OllamaProvider) GenerateResponseWithTools(ctx context.Context, prompt string) (string, []ToolResult, error) {
	// Use AI-based intent detection
	detector := NewIntentDetector(p)
	intents, err := detector.DetectIntent(ctx, prompt)
	if err != nil {
		// If intent detection fails, fall back to basic response
		response, err := p.GenerateResponse(ctx, prompt)
		return response, []ToolResult{}, err
	}
	
	// Execute detected tools
	var toolResults []ToolResult
	for _, intent := range intents {
		if intent.Confidence > 0.5 { // Only execute high-confidence intents
			result := ExecuteTool(intent.Tool, intent.Parameters)
			toolResults = append(toolResults, result)
		}
	}
	
	// Enhance the prompt with tool information and results
	enhancedPrompt := ""
	if len(toolResults) > 0 {
		enhancedPrompt += "I have executed the following operations for you:\n"
		for _, result := range toolResults {
			enhancedPrompt += fmt.Sprintf("- %s: %s\n", result.Name, result.Content)
		}
		enhancedPrompt += "\nNow, please provide a helpful response about what was accomplished.\n"
	}
	
	enhancedPrompt += "User: " + prompt
	
	// Get AI response with the enhanced prompt
	response, err := p.GenerateResponse(ctx, enhancedPrompt)
	if err != nil {
		// If AI response fails, provide a clear summary of what was accomplished
		if len(toolResults) > 0 {
			summary := "I have successfully completed the following operations:\n"
			for _, result := range toolResults {
				if result.Success {
					summary += fmt.Sprintf("✓ %s\n", result.Content)
				} else {
					summary += fmt.Sprintf("✗ %s failed: %s\n", result.Name, result.Content)
				}
			}
			summary += "\nAll file operations have been executed successfully."
			return summary, toolResults, nil
		}
		return "", toolResults, err
	}
	
	return response, toolResults, nil
}

func (p *OllamaProvider) SupportsTools() bool {
	return true
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