package ai

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"tala/internal/fileops"
	"time"
)

// Tool represents a function that the AI can call
type Tool struct {
	Name        string                                   `json:"name"`
	Description string                                   `json:"description"`
	Parameters  map[string]interface{}                   `json:"parameters"`
	Execute     func(args map[string]interface{}) string `json:"-"`
}

// ToolCall represents a request from AI to execute a tool
type ToolCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// ToolResult represents the result of executing a tool
type ToolResult struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	Success bool   `json:"success"`
}

// GetAvailableTools returns all tools available to the AI
func GetAvailableTools() []Tool {
	return []Tool{
		{
			Name:        "list_files",
			Description: "List files and directories in the current directory or specified path",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "Optional path to list. If not provided, lists current directory",
					},
				},
			},
			Execute: func(args map[string]interface{}) string {
				path := ""
				if p, ok := args["path"].(string); ok {
					path = p
				}
				result := fileops.ListDirectory(path)
				return result.Message
			},
		},
		{
			Name:        "read_file",
			Description: "Read the contents of a file",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"filename": map[string]interface{}{
						"type":        "string",
						"description": "Name of the file to read",
					},
				},
				"required": []string{"filename"},
			},
			Execute: func(args map[string]interface{}) string {
				filename, ok := args["filename"].(string)
				if !ok {
					return "Error: filename is required"
				}
				result := fileops.ReadFile(filename)
				return result.Message
			},
		},
		{
			Name:        "create_file",
			Description: "Create a new file with specified content",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"filename": map[string]interface{}{
						"type":        "string",
						"description": "Name of the file to create",
					},
					"content": map[string]interface{}{
						"type":        "string",
						"description": "Content to write to the file",
					},
				},
				"required": []string{"filename", "content"},
			},
			Execute: func(args map[string]interface{}) string {
				filename, ok1 := args["filename"].(string)
				content, ok2 := args["content"].(string)
				if !ok1 || !ok2 {
					return "Error: filename and content are required"
				}
				result := fileops.CreateFile(filename, content)
				return result.Message
			},
		},
		{
			Name:        "update_file",
			Description: "Update an existing file with new content",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"filename": map[string]interface{}{
						"type":        "string",
						"description": "Name of the file to update",
					},
					"content": map[string]interface{}{
						"type":        "string",
						"description": "New content for the file",
					},
				},
				"required": []string{"filename", "content"},
			},
			Execute: func(args map[string]interface{}) string {
				filename, ok1 := args["filename"].(string)
				content, ok2 := args["content"].(string)
				if !ok1 || !ok2 {
					return "Error: filename and content are required"
				}
				result := fileops.UpdateFile(filename, content)
				return result.Message
			},
		},
		{
			Name:        "delete_file",
			Description: "Delete a file",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"filename": map[string]interface{}{
						"type":        "string",
						"description": "Name of the file to delete",
					},
				},
				"required": []string{"filename"},
			},
			Execute: func(args map[string]interface{}) string {
				filename, ok := args["filename"].(string)
				if !ok {
					return "Error: filename is required"
				}
				result := fileops.DeleteFile(filename)
				return result.Message
			},
		},
		{
			Name:        "create_directory",
			Description: "Create a new directory",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"dirname": map[string]interface{}{
						"type":        "string",
						"description": "Name of the directory to create",
					},
				},
				"required": []string{"dirname"},
			},
			Execute: func(args map[string]interface{}) string {
				dirname, ok := args["dirname"].(string)
				if !ok {
					return "Error: dirname is required"
				}
				result := fileops.CreateDirectory(dirname)
				return result.Message
			},
		},
		{
			Name:        "delete_directory",
			Description: "Delete a directory and its contents",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"dirname": map[string]interface{}{
						"type":        "string",
						"description": "Name of the directory to delete",
					},
				},
				"required": []string{"dirname"},
			},
			Execute: func(args map[string]interface{}) string {
				dirname, ok := args["dirname"].(string)
				if !ok {
					return "Error: dirname is required"
				}
				result := fileops.DeleteDirectory(dirname)
				return result.Message
			},
		},
		{
			Name:        "copy_file",
			Description: "Copy a file from source to destination",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"source": map[string]interface{}{
						"type":        "string",
						"description": "Source file path",
					},
					"destination": map[string]interface{}{
						"type":        "string",
						"description": "Destination file path",
					},
				},
				"required": []string{"source", "destination"},
			},
			Execute: func(args map[string]interface{}) string {
				source, ok1 := args["source"].(string)
				destination, ok2 := args["destination"].(string)
				if !ok1 || !ok2 {
					return "Error: source and destination are required"
				}
				result := fileops.CopyFile(source, destination)
				return result.Message
			},
		},
		{
			Name:        "move_file",
			Description: "Move/rename a file from source to destination",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"source": map[string]interface{}{
						"type":        "string",
						"description": "Source file path",
					},
					"destination": map[string]interface{}{
						"type":        "string",
						"description": "Destination file path",
					},
				},
				"required": []string{"source", "destination"},
			},
			Execute: func(args map[string]interface{}) string {
				source, ok1 := args["source"].(string)
				destination, ok2 := args["destination"].(string)
				if !ok1 || !ok2 {
					return "Error: source and destination are required"
				}
				result := fileops.MoveFile(source, destination)
				return result.Message
			},
		},
		{
			Name:        "get_working_directory",
			Description: "Get the current working directory",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{},
			},
			Execute: func(args map[string]interface{}) string {
				result := fileops.GetWorkingDirectory()
				return result.Message
			},
		},
		{
			Name:        "change_directory",
			Description: "Change the current working directory",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "Path to change to",
					},
				},
				"required": []string{"path"},
			},
			Execute: func(args map[string]interface{}) string {
				path, ok := args["path"].(string)
				if !ok {
					return "Error: path is required"
				}
				result := fileops.ChangeDirectory(path)
				return result.Message
			},
		},
		{
			Name:        "execute_command",
			Description: "Execute a shell/bash command in the current directory",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"command": map[string]interface{}{
						"type":        "string",
						"description": "The shell command to execute",
					},
					"timeout": map[string]interface{}{
						"type":        "number",
						"description": "Optional timeout in seconds (default: 30)",
					},
				},
				"required": []string{"command"},
			},
			Execute: func(args map[string]interface{}) string {
				command, ok := args["command"].(string)
				if !ok {
					return "Error: command is required"
				}
				
				// Get timeout (default 30 seconds)
				timeout := 30.0
				if t, ok := args["timeout"].(float64); ok && t > 0 {
					timeout = t
				}
				
				result := ExecuteShellCommand(command, time.Duration(timeout)*time.Second)
				return result
			},
		},
		{
			Name:        "list_processes",
			Description: "List running processes",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"filter": map[string]interface{}{
						"type":        "string",
						"description": "Optional filter to search for specific processes",
					},
				},
			},
			Execute: func(args map[string]interface{}) string {
				filter := ""
				if f, ok := args["filter"].(string); ok {
					filter = f
				}
				
				var cmd *exec.Cmd
				if runtime.GOOS == "windows" {
					if filter != "" {
						// Sanitize filter input to prevent command injection
						sanitizedFilter := strings.ReplaceAll(filter, "*", "")
						sanitizedFilter = strings.ReplaceAll(sanitizedFilter, "&", "")
						sanitizedFilter = strings.ReplaceAll(sanitizedFilter, "|", "")
						sanitizedFilter = strings.ReplaceAll(sanitizedFilter, ";", "")
						if sanitizedFilter != "" {
							cmd = exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq *%s*", sanitizedFilter))
						} else {
							cmd = exec.Command("tasklist")
						}
					} else {
						cmd = exec.Command("tasklist")
					}
				} else {
					if filter != "" {
						cmd = exec.Command("ps", "aux")
					} else {
						cmd = exec.Command("ps", "aux")
					}
				}
				
				output, err := cmd.Output()
				if err != nil {
					return fmt.Sprintf("Error listing processes: %v", err)
				}
				
				result := string(output)
				if filter != "" && runtime.GOOS != "windows" {
					// Filter results on Unix-like systems
					lines := strings.Split(result, "\n")
					var filtered []string
					for _, line := range lines {
						if strings.Contains(strings.ToLower(line), strings.ToLower(filter)) {
							filtered = append(filtered, line)
						}
					}
					result = strings.Join(filtered, "\n")
				}
				
				return result
			},
		},
		{
			Name:        "get_system_info",
			Description: "Get system information (OS, architecture, etc.)",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{},
			},
			Execute: func(args map[string]interface{}) string {
				info := fmt.Sprintf("OS: %s\nArchitecture: %s\nCPUs: %d",
					runtime.GOOS, runtime.GOARCH, runtime.NumCPU())
				
				// Add additional info based on OS
				if runtime.GOOS == "windows" {
					cmd := exec.Command("systeminfo")
					if output, err := cmd.Output(); err == nil {
						// Extract key info from systeminfo
						lines := strings.Split(string(output), "\n")
						for _, line := range lines {
							if strings.Contains(line, "OS Name") ||
								strings.Contains(line, "System Type") ||
								strings.Contains(line, "Total Physical Memory") {
								info += "\n" + strings.TrimSpace(line)
							}
						}
					}
				} else {
					// Unix-like systems
					if output, err := exec.Command("uname", "-a").Output(); err == nil {
						info += "\nKernel: " + strings.TrimSpace(string(output))
					}
				}
				
				return info
			},
		},
	}
}

// ExecuteTool executes a tool with the given arguments
func ExecuteTool(toolName string, args map[string]interface{}) ToolResult {
	tools := GetAvailableTools()
	
	for _, tool := range tools {
		if tool.Name == toolName {
			content := tool.Execute(args)
			// Determine success based on whether the content indicates an error
			contentStr := content
			success := !strings.HasPrefix(contentStr, "Error") && 
					  !strings.HasPrefix(contentStr, "Failed")
			
			return ToolResult{
				Name:    toolName,
				Content: content,
				Success: success,
			}
		}
	}
	
	return ToolResult{
		Name:    toolName,
		Content: fmt.Sprintf("Unknown tool: %s", toolName),
		Success: false,
	}
}

// ParseToolCalls attempts to parse tool calls from AI response text
func ParseToolCalls(responseText string) []ToolCall {
	var toolCalls []ToolCall
	
	// Simple parsing - look for JSON objects that match tool call pattern
	// This is a basic implementation - more sophisticated parsing could be added
	
	// For now, return empty slice - providers will need to implement their own
	// tool calling mechanisms based on their specific formats
	
	return toolCalls
}

// FormatToolsForPrompt formats available tools for inclusion in AI prompts
func FormatToolsForPrompt() string {
	tools := GetAvailableTools()
	
	prompt := "You have access to the following file system tools:\n\n"
	
	for _, tool := range tools {
		prompt += fmt.Sprintf("- %s: %s\n", tool.Name, tool.Description)
		
		// Add parameter info
		if params, ok := tool.Parameters["properties"].(map[string]interface{}); ok {
			prompt += "  Parameters:\n"
			for paramName, paramInfo := range params {
				if info, ok := paramInfo.(map[string]interface{}); ok {
					if desc, ok := info["description"].(string); ok {
						prompt += fmt.Sprintf("    - %s: %s\n", paramName, desc)
					}
				}
			}
		}
		prompt += "\n"
	}
	
	prompt += "To use these tools, mention what you want to do and I will execute the appropriate operations.\n"
	prompt += "For example:\n"
	prompt += "- 'create a hello.txt file with Hello World' - I will create that file\n"
	prompt += "- 'run ls command' - I will execute the ls command\n"
	prompt += "- 'show running processes' - I will list processes\n"
	prompt += "- 'get system information' - I will show system details\n"
	
	return prompt
}

// ExecuteShellCommand executes a shell command with timeout and security checks
func ExecuteShellCommand(command string, timeout time.Duration) string {
	// Security check: block dangerous commands
	if !isCommandSafe(command) {
		return "Error: Command blocked for security reasons"
	}
	
	var cmd *exec.Cmd
	
	// Choose shell based on OS
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	
	// Set up timeout (default max 30 seconds)
	if timeout <= 0 || timeout > 30*time.Second {
		timeout = 30 * time.Second
	}
	
	// Execute with timeout
	done := make(chan error, 1)
	var output []byte
	var err error
	
	go func() {
		output, err = cmd.CombinedOutput()
		done <- err
	}()
	
	select {
	case <-time.After(timeout):
		if cmd.Process != nil {
			if err := cmd.Process.Kill(); err != nil {
				// Log error but continue - process might already be dead
			}
		}
		return fmt.Sprintf("Command timed out after %v", timeout)
	case execErr := <-done:
		if execErr != nil {
			return fmt.Sprintf("Command failed: %v\nOutput: %s", execErr, string(output))
		}
	}
	
	// Limit output size to prevent memory issues
	result := string(output)
	if len(result) > 10000 {
		result = result[:10000] + "\n... (output truncated)"
	}
	
	return result
}

// isCommandSafe checks if a command is safe to execute
func isCommandSafe(command string) bool {
	command = strings.ToLower(strings.TrimSpace(command))
	
	// Block empty commands
	if command == "" {
		return false
	}
	
	// List of dangerous command patterns
	dangerousPatterns := []string{
		"rm -rf",
		"rm -r /",
		"mkfs",
		"dd if=",
		":(){ :|:& };:",  // fork bomb
		"curl", "wget",   // network access (can be dangerous)
		"sudo",
		"su ",
		"passwd",
		"useradd",
		"userdel",
		"chmod 777",
		"chown root",
		"systemctl",
		"service",
		"reboot",
		"shutdown",
		"halt",
		"poweroff",
		"mount",
		"umount",
		"fdisk",
		"parted",
		"format",
		"> /dev/",
		"nc ", "netcat",  // network tools
		"ssh",
		"scp",
		"rsync",
		"crontab",
		"at ",
		"killall",
		"pkill",
		"kill -9",
		"python -c",
		"perl -e",
		"ruby -e",
		"node -e",
		"eval",
		"exec",
		"/bin/bash",
		"/bin/sh",
		"bash -c",
		"sh -c",
	}
	
	// Check against dangerous patterns
	for _, pattern := range dangerousPatterns {
		if strings.Contains(command, pattern) {
			return false
		}
	}
	
	// Block commands with potentially dangerous characters
	dangerousChars := []string{
		";",     // command chaining
		"&&",    // command chaining
		"||",    // command chaining
		"|",     // piping (can be dangerous)
		">",     // redirection
		">>",    // redirection
		"<",     // redirection
		"`",     // command substitution
		"$(",    // command substitution
		"$()",   // command substitution
		"../",   // directory traversal
		"./",    // current directory execution
	}
	
	for _, char := range dangerousChars {
		if strings.Contains(command, char) {
			return false
		}
	}
	
	// Allow only specific safe commands
	safeCommands := []string{
		"ls", "dir", "pwd", "cd", "echo", "cat", "head", "tail",
		"grep", "find", "which", "where", "type", "file",
		"date", "whoami", "id", "uptime", "uname", "hostname",
		"ps", "top", "df", "du", "free", "lscpu", "lsblk",
		"env", "printenv", "history", "alias",
		"wc", "sort", "uniq", "cut", "awk", "sed",
		"git status", "git log", "git branch", "git diff",
		"go version", "go list", "go mod",
		"npm list", "npm version",
		"python --version", "python3 --version",
		"node --version", "php --version",
		"java -version", "javac -version",
		"gcc --version", "clang --version",
	}
	
	// Extract the base command (first word)
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return false
	}
	
	baseCommand := parts[0]
	
	// Check if base command is in safe list
	for _, safe := range safeCommands {
		if strings.HasPrefix(safe, baseCommand) {
			return true
		}
	}
	
	// Special case: allow version checks
	if strings.Contains(command, "--version") || strings.Contains(command, "-version") {
		return true
	}
	
	// Special case: allow help commands
	if strings.Contains(command, "--help") || strings.Contains(command, "-h") {
		return true
	}
	
	return false
}