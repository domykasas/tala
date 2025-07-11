package fileops

import (
	"os"
	"strings"
	"testing"
)

func TestExecuteCommand(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer cleanupTestDir(t, tmpDir)

	// Change to temp directory
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	tests := []struct {
		name      string
		input     string
		wantErr   bool
		wantSuccess bool
	}{
		{
			name:        "help command",
			input:       "/help",
			wantErr:     false,
			wantSuccess: true,
		},
		{
			name:        "help command without slash",
			input:       "help",
			wantErr:     false,
			wantSuccess: true,
		},
		{
			name:        "pwd command",
			input:       "/pwd",
			wantErr:     false,
			wantSuccess: true,
		},
		{
			name:        "ls command",
			input:       "/ls",
			wantErr:     false,
			wantSuccess: true,
		},
		{
			name:        "create file command",
			input:       "/create test.txt Hello World",
			wantErr:     false,
			wantSuccess: true,
		},
		{
			name:        "cat command",
			input:       "/cat test.txt",
			wantErr:     false,
			wantSuccess: true,
		},
		{
			name:        "mkdir command",
			input:       "/mkdir testdir",
			wantErr:     false,
			wantSuccess: true,
		},
		{
			name:        "unknown command",
			input:       "/unknown",
			wantErr:     false,
			wantSuccess: false,
		},
		{
			name:        "empty command",
			input:       "/",
			wantErr:     false,
			wantSuccess: false,
		},
		{
			name:        "command with missing args",
			input:       "/cat",
			wantErr:     false,
			wantSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExecuteCommand(tt.input)
			
			if (result.Error != nil) != tt.wantErr {
				t.Errorf("ExecuteCommand() error = %v, wantErr %v", result.Error, tt.wantErr)
				return
			}
			
			if result.Success != tt.wantSuccess {
				t.Errorf("ExecuteCommand() success = %v, wantSuccess %v", result.Success, tt.wantSuccess)
				return
			}
			
			if result.Message == "" {
				t.Errorf("ExecuteCommand() should always return a message")
			}
		})
	}
}

func TestGetHelpText(t *testing.T) {
	helpText := GetHelpText()
	
	if helpText == "" {
		t.Error("GetHelpText() should return non-empty string")
	}
	
	// Check that help text contains expected commands
	expectedCommands := []string{"ls", "cat", "create", "mkdir", "rm", "pwd", "cd", "cp", "mv"}
	
	for _, cmd := range expectedCommands {
		if !strings.Contains(helpText, cmd) {
			t.Errorf("GetHelpText() should contain command: %s", cmd)
		}
	}
	
	// Check that help text contains usage examples
	if !strings.Contains(helpText, "Example usage:") {
		t.Error("GetHelpText() should contain usage examples")
	}
}

func TestGetCommands(t *testing.T) {
	commands := GetCommands()
	
	if len(commands) == 0 {
		t.Error("GetCommands() should return non-empty map")
	}
	
	// Check that essential commands are present
	expectedCommands := []string{"ls", "cat", "create", "mkdir", "rm", "pwd", "cd", "cp", "mv", "write", "update", "rmdir"}
	
	for _, cmdName := range expectedCommands {
		if cmd, exists := commands[cmdName]; !exists {
			t.Errorf("GetCommands() should contain command: %s", cmdName)
		} else {
			// Check that command has required fields
			if cmd.Name == "" {
				t.Errorf("Command %s should have a name", cmdName)
			}
			if cmd.Description == "" {
				t.Errorf("Command %s should have a description", cmdName)
			}
			if cmd.Usage == "" {
				t.Errorf("Command %s should have usage", cmdName)
			}
			if cmd.Execute == nil {
				t.Errorf("Command %s should have execute function", cmdName)
			}
		}
	}
}

func TestCommandExecution(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer cleanupTestDir(t, tmpDir)

	// Change to temp directory
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	commands := GetCommands()
	
	// Test pwd command
	if cmd, exists := commands["pwd"]; exists {
		result := cmd.Execute([]string{})
		if !result.Success {
			t.Errorf("pwd command should succeed")
		}
	}
	
	// Test ls command
	if cmd, exists := commands["ls"]; exists {
		result := cmd.Execute([]string{})
		if !result.Success {
			t.Errorf("ls command should succeed")
		}
	}
	
	// Test create command
	if cmd, exists := commands["create"]; exists {
		result := cmd.Execute([]string{"test.txt", "Hello", "World"})
		if !result.Success {
			t.Errorf("create command should succeed: %s", result.Message)
		}
		
		// Verify file was created
		if _, err := os.Stat("test.txt"); os.IsNotExist(err) {
			t.Error("create command should have created file")
		}
	}
	
	// Test cat command
	if cmd, exists := commands["cat"]; exists {
		result := cmd.Execute([]string{"test.txt"})
		if !result.Success {
			t.Errorf("cat command should succeed: %s", result.Message)
		}
	}
	
	// Test mkdir command
	if cmd, exists := commands["mkdir"]; exists {
		result := cmd.Execute([]string{"testdir"})
		if !result.Success {
			t.Errorf("mkdir command should succeed: %s", result.Message)
		}
		
		// Verify directory was created
		if _, err := os.Stat("testdir"); os.IsNotExist(err) {
			t.Error("mkdir command should have created directory")
		}
	}
	
	// Test invalid command arguments
	if cmd, exists := commands["cat"]; exists {
		result := cmd.Execute([]string{}) // No filename
		if result.Success {
			t.Error("cat command should fail with no arguments")
		}
	}
	
	if cmd, exists := commands["cp"]; exists {
		result := cmd.Execute([]string{"onefile"}) // Missing destination
		if result.Success {
			t.Error("cp command should fail with insufficient arguments")
		}
	}
}