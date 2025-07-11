package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// Intent represents a detected user intention
type Intent struct {
	Action     string                 `json:"action"`
	Tool       string                 `json:"tool"`
	Parameters map[string]interface{} `json:"parameters"`
	Confidence float64                `json:"confidence"`
}

// IntentDetector uses AI to detect user intentions
type IntentDetector struct {
	provider Provider
}

// NewIntentDetector creates a new intent detector
func NewIntentDetector(provider Provider) *IntentDetector {
	return &IntentDetector{
		provider: provider,
	}
}

// DetectIntent analyzes user input and returns detected intentions
func (detector *IntentDetector) DetectIntent(ctx context.Context, userInput string) ([]Intent, error) {
	// Create a prompt for intent detection
	prompt := detector.createIntentDetectionPrompt(userInput)
	
	// Get AI response
	response, err := detector.provider.GenerateResponse(ctx, prompt)
	if err != nil {
		// Fallback to simple pattern matching if AI fails
		return detector.fallbackPatternMatching(userInput), nil
	}
	
	// Parse the AI response to extract intents
	intents := detector.parseIntentResponse(response)
	
	// If no intents detected by AI, use fallback
	if len(intents) == 0 {
		return detector.fallbackPatternMatching(userInput), nil
	}
	
	return intents, nil
}

// createIntentDetectionPrompt creates a prompt for AI intent detection
func (detector *IntentDetector) createIntentDetectionPrompt(userInput string) string {
	availableTools := GetAvailableTools()
	
	prompt := `You are a conservative intent detection system. Only detect tool usage when the user explicitly requests file operations, commands, or system actions.

DO NOT detect intents for:
- Greetings (hi, hello, hey)
- General questions
- Casual conversation
- Abstract requests

Available tools and their purposes:
`
	
	for _, tool := range availableTools {
		prompt += fmt.Sprintf("- %s: %s\n", tool.Name, tool.Description)
	}
	
	prompt += `
Respond with JSON containing an array of detected intents. Each intent should have:
- "action": brief description of what the user wants to do
- "tool": the tool name to use (from the list above)
- "parameters": object with the parameters for the tool
- "confidence": number between 0-1 indicating confidence

Only detect intents with confidence > 0.8 when user explicitly mentions:
- File operations (create, read, write, delete specific files)
- Directory operations (list, create, delete directories)
- System commands (run specific commands)

For general conversation, greetings, or questions, respond with: []

User input: "` + userInput + `"

JSON response:`
	
	return prompt
}

// parseIntentResponse parses AI response to extract intents
func (detector *IntentDetector) parseIntentResponse(response string) []Intent {
	var intents []Intent
	
	// Find JSON in the response
	jsonStart := strings.Index(response, "[")
	jsonEnd := strings.LastIndex(response, "]")
	
	if jsonStart == -1 || jsonEnd == -1 || jsonEnd <= jsonStart {
		return intents
	}
	
	jsonStr := response[jsonStart : jsonEnd+1]
	
	// Try to parse JSON
	err := json.Unmarshal([]byte(jsonStr), &intents)
	if err != nil {
		// Try to extract individual intent objects
		return detector.extractIntentsFromText(response)
	}
	
	return intents
}

// extractIntentsFromText tries to extract intents from non-JSON text
func (detector *IntentDetector) extractIntentsFromText(text string) []Intent {
	var intents []Intent
	
	// Look for mentions of tool names in the response
	availableTools := GetAvailableTools()
	text = strings.ToLower(text)
	
	for _, tool := range availableTools {
		if strings.Contains(text, tool.Name) || strings.Contains(text, strings.ReplaceAll(tool.Name, "_", " ")) {
			intent := Intent{
				Action:     "detected from text",
				Tool:       tool.Name,
				Parameters: make(map[string]interface{}),
				Confidence: 0.6,
			}
			
			// Try to extract simple parameters
			if tool.Name == "execute_command" {
				if command := detector.extractCommand(text); command != "" {
					intent.Parameters["command"] = command
					intent.Confidence = 0.8
				}
			}
			
			intents = append(intents, intent)
		}
	}
	
	return intents
}

// extractCommand tries to extract a command from text
func (detector *IntentDetector) extractCommand(text string) string {
	// Look for common command patterns
	patterns := []string{
		"run ",
		"execute ",
		"command ",
		"bash ",
		"shell ",
	}
	
	for _, pattern := range patterns {
		if idx := strings.Index(text, pattern); idx != -1 {
			remainder := text[idx+len(pattern):]
			// Extract until next sentence or period
			if endIdx := strings.IndexAny(remainder, ".!?\n"); endIdx != -1 {
				return strings.TrimSpace(remainder[:endIdx])
			}
			return strings.TrimSpace(remainder)
		}
	}
	
	return ""
}

// fallbackPatternMatching provides simple pattern matching as fallback
func (detector *IntentDetector) fallbackPatternMatching(userInput string) []Intent {
	var intents []Intent
	input := strings.ToLower(userInput)
	
	// File operations
	if (strings.Contains(input, "create") || strings.Contains(input, "make")) && strings.Contains(input, "file") {
		intent := Intent{
			Action:     "create file",
			Tool:       "create_file",
			Parameters: detector.extractFileParams(userInput),
			Confidence: 0.7,
		}
		intents = append(intents, intent)
	}
	
	// Directory operations
	if (strings.Contains(input, "create") || strings.Contains(input, "make")) && 
		(strings.Contains(input, "directory") || strings.Contains(input, "folder")) {
		intent := Intent{
			Action:     "create directory",
			Tool:       "create_directory",
			Parameters: detector.extractDirParams(userInput),
			Confidence: 0.7,
		}
		intents = append(intents, intent)
	}
	
	// List files
	if (strings.Contains(input, "list") || strings.Contains(input, "show")) && 
		(strings.Contains(input, "file") || strings.Contains(input, "directory")) {
		intent := Intent{
			Action:     "list files",
			Tool:       "list_files",
			Parameters: make(map[string]interface{}),
			Confidence: 0.8,
		}
		intents = append(intents, intent)
	}
	
	// Execute command
	if strings.Contains(input, "run") || strings.Contains(input, "execute") || 
		strings.Contains(input, "command") || strings.Contains(input, "bash") {
		command := detector.extractCommand(input)
		if command != "" {
			intent := Intent{
				Action:     "execute command",
				Tool:       "execute_command",
				Parameters: map[string]interface{}{"command": command},
				Confidence: 0.8,
			}
			intents = append(intents, intent)
		}
	}
	
	// Working directory
	if strings.Contains(input, "working directory") || strings.Contains(input, "current directory") || 
		strings.Contains(input, "where am i") {
		intent := Intent{
			Action:     "get working directory",
			Tool:       "get_working_directory",
			Parameters: make(map[string]interface{}),
			Confidence: 0.9,
		}
		intents = append(intents, intent)
	}
	
	// System info
	if strings.Contains(input, "system") && strings.Contains(input, "info") {
		intent := Intent{
			Action:     "get system info",
			Tool:       "get_system_info",
			Parameters: make(map[string]interface{}),
			Confidence: 0.9,
		}
		intents = append(intents, intent)
	}
	
	// List processes
	if strings.Contains(input, "process") || strings.Contains(input, "task") {
		intent := Intent{
			Action:     "list processes",
			Tool:       "list_processes",
			Parameters: make(map[string]interface{}),
			Confidence: 0.8,
		}
		intents = append(intents, intent)
	}
	
	return intents
}

// extractFileParams extracts filename and content from user input
func (detector *IntentDetector) extractFileParams(userInput string) map[string]interface{} {
	params := make(map[string]interface{})
	words := strings.Fields(userInput)
	
	// Look for filename - prioritize specific patterns
	for _, word := range words {
		// First check for file with extension
		if strings.Contains(word, ".") {
			filename := strings.Trim(word, "\"'")
			params["filename"] = filename
			break
		}
	}
	
	// If no filename found with extension, look for patterns
	if _, exists := params["filename"]; !exists {
		for i, word := range words {
			// Special case: look for "file called X"
			if i > 1 && words[i-1] == "called" && words[i-2] == "file" {
				filename := strings.Trim(word, "\"'")
				if !strings.Contains(filename, ".") {
					filename += ".txt"
				}
				params["filename"] = filename
				break
			}
			// Look for "file named X"
			if i > 1 && words[i-1] == "named" && words[i-2] == "file" {
				filename := strings.Trim(word, "\"'")
				if !strings.Contains(filename, ".") {
					filename += ".txt"
				}
				params["filename"] = filename
				break
			}
		}
	}
	
	// Look for content
	content := "Hello World!" // Default
	for i, word := range words {
		if word == "with" || word == "containing" {
			if i+1 < len(words) {
				remainingWords := words[i+1:]
				contentStr := strings.Join(remainingWords, " ")
				contentStr = strings.TrimSuffix(contentStr, " in it")
				contentStr = strings.TrimSuffix(contentStr, " to it")
				contentStr = strings.Trim(contentStr, "\"'")
				if contentStr != "" {
					content = contentStr
				}
				break
			}
		}
	}
	params["content"] = content
	
	return params
}

// extractDirParams extracts directory name from user input
func (detector *IntentDetector) extractDirParams(userInput string) map[string]interface{} {
	params := make(map[string]interface{})
	words := strings.Fields(userInput)
	
	for i, word := range words {
		if i > 0 && (words[i-1] == "directory" || words[i-1] == "folder" || 
			words[i-1] == "called" || words[i-1] == "named") {
			dirname := strings.Trim(word, "\"'")
			params["dirname"] = dirname
			break
		}
	}
	
	return params
}