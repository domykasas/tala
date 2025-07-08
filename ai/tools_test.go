package ai

import (
	"context"
	"os"
	"testing"
)

func setupTestDir(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "tala-ai-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	return tmpDir
}

func cleanupTestDir(t *testing.T, dir string) {
	if err := os.RemoveAll(dir); err != nil {
		t.Errorf("Failed to cleanup temp dir: %v", err)
	}
}

func TestGetAvailableTools(t *testing.T) {
	tools := GetAvailableTools()
	
	if len(tools) == 0 {
		t.Error("GetAvailableTools() should return non-empty slice")
	}
	
	// Check that essential tools are present
	expectedTools := []string{
		"list_files", "read_file", "create_file", "update_file", "delete_file",
		"create_directory", "delete_directory", "copy_file", "move_file",
		"get_working_directory", "change_directory",
	}
	
	toolMap := make(map[string]bool)
	for _, tool := range tools {
		toolMap[tool.Name] = true
		
		// Verify tool has required fields
		if tool.Name == "" {
			t.Error("Tool should have a name")
		}
		if tool.Description == "" {
			t.Error("Tool should have a description")
		}
		if tool.Execute == nil {
			t.Error("Tool should have an execute function")
		}
	}
	
	for _, expectedTool := range expectedTools {
		if !toolMap[expectedTool] {
			t.Errorf("Expected tool %s not found", expectedTool)
		}
	}
}

func TestExecuteTool(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer cleanupTestDir(t, tmpDir)
	
	// Change to temp directory
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)
	
	tests := []struct {
		name     string
		toolName string
		args     map[string]interface{}
		wantErr  bool
	}{
		{
			name:     "get working directory",
			toolName: "get_working_directory",
			args:     map[string]interface{}{},
			wantErr:  false,
		},
		{
			name:     "execute simple command",
			toolName: "execute_command",
			args: map[string]interface{}{
				"command": "echo hello",
			},
			wantErr: false,
		},
		{
			name:     "get system info",
			toolName: "get_system_info",
			args:     map[string]interface{}{},
			wantErr:  false,
		},
		{
			name:     "list files",
			toolName: "list_files",
			args:     map[string]interface{}{},
			wantErr:  false,
		},
		{
			name:     "create file",
			toolName: "create_file",
			args: map[string]interface{}{
				"filename": "test.txt",
				"content":  "Hello World",
			},
			wantErr: false,
		},
		{
			name:     "read file",
			toolName: "read_file",
			args: map[string]interface{}{
				"filename": "test.txt",
			},
			wantErr: false,
		},
		{
			name:     "create directory",
			toolName: "create_directory",
			args: map[string]interface{}{
				"dirname": "testdir",
			},
			wantErr: false,
		},
		{
			name:     "unknown tool",
			toolName: "unknown_tool",
			args:     map[string]interface{}{},
			wantErr:  true,
		},
		{
			name:     "create file missing args",
			toolName: "create_file",
			args: map[string]interface{}{
				"filename": "test2.txt",
				// missing content
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExecuteTool(tt.toolName, tt.args)
			
			if tt.wantErr && result.Success {
				t.Errorf("ExecuteTool() expected error but got success")
			}
			if !tt.wantErr && !result.Success {
				t.Errorf("ExecuteTool() expected success but got error: %s", result.Content)
			}
			
			if result.Name != tt.toolName {
				t.Errorf("ExecuteTool() result name = %v, want %v", result.Name, tt.toolName)
			}
			
			if result.Content == "" {
				t.Error("ExecuteTool() should always return content")
			}
		})
	}
}

func TestFormatToolsForPrompt(t *testing.T) {
	prompt := FormatToolsForPrompt()
	
	if prompt == "" {
		t.Error("FormatToolsForPrompt() should return non-empty string")
	}
	
	// Check that it contains tool information
	if !contains(prompt, "file system tools") {
		t.Error("FormatToolsForPrompt() should mention file system tools")
	}
	
	// Check that it contains some tool names
	expectedTools := []string{"list_files", "create_file", "read_file"}
	for _, tool := range expectedTools {
		if !contains(prompt, tool) {
			t.Errorf("FormatToolsForPrompt() should contain tool: %s", tool)
		}
	}
}

func TestOllamaProviderDetectFileOperations(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer cleanupTestDir(t, tmpDir)
	
	// Change to temp directory
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)
	
	provider := NewOllamaProvider("test-model", 0.7, 1000, "")
	
	tests := []struct {
		name           string
		prompt         string
		expectedTools  int
		expectedTool   string
	}{
		{
			name:          "create file request",
			prompt:        "create a file called test.txt with hello content",
			expectedTools: 1,
			expectedTool:  "create_file",
		},
		{
			name:          "list files request",
			prompt:        "list all files in the directory",
			expectedTools: 1,
			expectedTool:  "list_files",
		},
		{
			name:          "working directory request",
			prompt:        "what is my current working directory?",
			expectedTools: 1,
			expectedTool:  "get_working_directory",
		},
		{
			name:          "create directory request",
			prompt:        "make a directory called testfolder",
			expectedTools: 1,
			expectedTool:  "create_directory",
		},
		{
			name:          "no file operation",
			prompt:        "what is the weather today?",
			expectedTools: 0,
			expectedTool:  "",
		},
		{
			name:          "create file with quotes",
			prompt:        "create \"testfile.txt\" with \"Hello World!\" in it",
			expectedTools: 1,
			expectedTool:  "create_file",
		},
		{
			name:          "execute command request",
			prompt:        "run ls command",
			expectedTools: 1,
			expectedTool:  "execute_command",
		},
		{
			name:          "system info request",
			prompt:        "show me system information",
			expectedTools: 1,
			expectedTool:  "get_system_info",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use the intent detection system instead
			detector := NewIntentDetector(provider)
			intents, err := detector.DetectIntent(context.Background(), tt.prompt)
			
			if err != nil {
				// If intent detection fails, use fallback
				intents = detector.fallbackPatternMatching(tt.prompt)
			}
			
			if len(intents) != tt.expectedTools {
				t.Errorf("Intent detection returned %d tools, expected %d", len(intents), tt.expectedTools)
			}
			
			if tt.expectedTools > 0 && len(intents) > 0 {
				if intents[0].Tool != tt.expectedTool {
					t.Errorf("Intent detection used tool %s, expected %s", intents[0].Tool, tt.expectedTool)
				}
			}
		})
	}
}

func TestProviderSupportsTools(t *testing.T) {
	providers := []Provider{
		NewOpenAIProvider("test", "test", 0.7, 1000),
		NewAnthropicProvider("test", "test", 0.7, 1000),
		NewOllamaProvider("test", 0.7, 1000, ""),
	}
	
	for _, provider := range providers {
		t.Run(provider.GetName(), func(t *testing.T) {
			if !provider.SupportsTools() {
				t.Errorf("Provider %s should support tools", provider.GetName())
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		containsHelper(s, substr))))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}