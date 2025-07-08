package ai

import (
	"context"
	"fmt"
	"testing"
)

func TestIntentDetector_FallbackPatternMatching(t *testing.T) {
	detector := &IntentDetector{}

	tests := []struct {
		name           string
		input          string
		expectedIntent string
		expectedTool   string
	}{
		{
			name:           "create file intent",
			input:          "create a test.txt file",
			expectedIntent: "create file",
			expectedTool:   "create_file",
		},
		{
			name:           "execute command intent",
			input:          "run ls command",
			expectedIntent: "execute command",
			expectedTool:   "execute_command",
		},
		{
			name:           "list files intent",
			input:          "show me all files",
			expectedIntent: "list files",
			expectedTool:   "list_files",
		},
		{
			name:           "system info intent",
			input:          "get system info",
			expectedIntent: "get system info",
			expectedTool:   "get_system_info",
		},
		{
			name:           "no clear intent",
			input:          "what is the weather?",
			expectedIntent: "",
			expectedTool:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			intents := detector.fallbackPatternMatching(tt.input)
			
			if tt.expectedTool == "" {
				if len(intents) != 0 {
					t.Errorf("Expected no intents, got %d", len(intents))
				}
				return
			}
			
			if len(intents) == 0 {
				t.Errorf("Expected at least one intent, got none")
				return
			}
			
			found := false
			for _, intent := range intents {
				if intent.Tool == tt.expectedTool {
					found = true
					if intent.Action != tt.expectedIntent {
						t.Errorf("Expected action '%s', got '%s'", tt.expectedIntent, intent.Action)
					}
					if intent.Confidence <= 0 || intent.Confidence > 1 {
						t.Errorf("Expected confidence between 0-1, got %f", intent.Confidence)
					}
					break
				}
			}
			
			if !found {
				t.Errorf("Expected tool '%s' not found in intents", tt.expectedTool)
			}
		})
	}
}

func TestIntentDetector_ExtractFileParams(t *testing.T) {
	detector := &IntentDetector{}

	tests := []struct {
		name             string
		input            string
		expectedFilename string
		expectedContent  string
	}{
		{
			name:             "filename with extension",
			input:            "create test.txt with hello content",
			expectedFilename: "test.txt",
			expectedContent:  "hello content",
		},
		{
			name:             "filename in quotes",
			input:            "create \"myfile.py\" with print('hello')",
			expectedFilename: "myfile.py",
			expectedContent:  "print('hello')",
		},
		{
			name:             "filename without extension",
			input:            "create file called myfile containing test data",
			expectedFilename: "myfile.txt",
			expectedContent:  "test data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := detector.extractFileParams(tt.input)
			
			if filename, ok := params["filename"].(string); ok {
				if filename != tt.expectedFilename {
					t.Errorf("Expected filename '%s', got '%s'", tt.expectedFilename, filename)
				}
			} else {
				t.Errorf("Expected filename parameter not found")
			}
			
			if content, ok := params["content"].(string); ok {
				if content != tt.expectedContent {
					t.Errorf("Expected content '%s', got '%s'", tt.expectedContent, content)
				}
			} else {
				t.Errorf("Expected content parameter not found")
			}
		})
	}
}

func TestIntentDetector_ExtractCommand(t *testing.T) {
	detector := &IntentDetector{}

	tests := []struct {
		name            string
		input           string
		expectedCommand string
	}{
		{
			name:            "run command",
			input:           "run ls -la",
			expectedCommand: "ls -la",
		},
		{
			name:            "execute command",
			input:           "execute echo hello world",
			expectedCommand: "echo hello world",
		},
		{
			name:            "bash command",
			input:           "bash pwd",
			expectedCommand: "pwd",
		},
		{
			name:            "no command",
			input:           "what is the weather",
			expectedCommand: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command := detector.extractCommand(tt.input)
			
			if command != tt.expectedCommand {
				t.Errorf("Expected command '%s', got '%s'", tt.expectedCommand, command)
			}
		})
	}
}

func TestIntentDetector_ParseIntentResponse(t *testing.T) {
	detector := &IntentDetector{}

	tests := []struct {
		name     string
		response string
		expected int // number of intents expected
	}{
		{
			name: "valid JSON response",
			response: `[
				{
					"action": "create file",
					"tool": "create_file",
					"parameters": {"filename": "test.txt", "content": "hello"},
					"confidence": 0.9
				}
			]`,
			expected: 1,
		},
		{
			name:     "empty array",
			response: "[]",
			expected: 0,
		},
		{
			name:     "no JSON",
			response: "I'll help you with that task",
			expected: 0,
		},
		{
			name: "JSON with extra text",
			response: `Here's what I'll do:
			[
				{
					"action": "list files",
					"tool": "list_files",
					"parameters": {},
					"confidence": 0.8
				}
			]
			Let me execute this for you.`,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			intents := detector.parseIntentResponse(tt.response)
			
			if len(intents) != tt.expected {
				t.Errorf("Expected %d intents, got %d", tt.expected, len(intents))
			}
			
			for _, intent := range intents {
				if intent.Tool == "" {
					t.Errorf("Intent tool should not be empty")
				}
				if intent.Confidence < 0 || intent.Confidence > 1 {
					t.Errorf("Intent confidence should be between 0-1, got %f", intent.Confidence)
				}
			}
		})
	}
}

// Mock provider for testing
type mockProvider struct {
	response string
	err      error
}

func (m *mockProvider) GenerateResponse(ctx context.Context, prompt string) (string, error) {
	return m.response, m.err
}

func (m *mockProvider) GenerateResponseWithTools(ctx context.Context, prompt string) (string, []ToolResult, error) {
	return m.response, []ToolResult{}, m.err
}

func (m *mockProvider) GetName() string {
	return "Mock"
}

func (m *mockProvider) SupportsTools() bool {
	return true
}

func TestIntentDetector_DetectIntent(t *testing.T) {
	tests := []struct {
		name          string
		mockResponse  string
		mockError     error
		userInput     string
		expectedTools int
	}{
		{
			name: "AI detects intent successfully",
			mockResponse: `[
				{
					"action": "create file",
					"tool": "create_file",
					"parameters": {"filename": "test.txt", "content": "hello"},
					"confidence": 0.9
				}
			]`,
			mockError:     nil,
			userInput:     "create test.txt with hello",
			expectedTools: 1,
		},
		{
			name:          "AI fails, fallback to pattern matching",
			mockResponse:  "",
			mockError:     fmt.Errorf("AI error"),
			userInput:     "run ls command",
			expectedTools: 1, // fallback should detect this
		},
		{
			name:          "No clear intent",
			mockResponse:  "[]",
			mockError:     nil,
			userInput:     "what is the weather?",
			expectedTools: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProv := &mockProvider{
				response: tt.mockResponse,
				err:      tt.mockError,
			}
			
			detector := NewIntentDetector(mockProv)
			intents, err := detector.DetectIntent(context.Background(), tt.userInput)
			
			if err != nil {
				t.Errorf("DetectIntent should not return error: %v", err)
			}
			
			if len(intents) != tt.expectedTools {
				t.Errorf("Expected %d intents, got %d", tt.expectedTools, len(intents))
			}
		})
	}
}