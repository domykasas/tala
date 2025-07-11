package fileops

import (
	"fmt"
	"strings"
)

// Command represents a file operation command
type Command struct {
	Name        string
	Description string
	Usage       string
	Execute     func(args []string) *FileOperation
}

// GetCommands returns all available file operation commands
func GetCommands() map[string]*Command {
	return map[string]*Command{
		"ls": {
			Name:        "ls",
			Description: "List files and directories",
			Usage:       "ls [path]",
			Execute: func(args []string) *FileOperation {
				path := ""
				if len(args) > 0 {
					path = args[0]
				}
				return ListDirectory(path)
			},
		},
		"cat": {
			Name:        "cat",
			Description: "Display file content",
			Usage:       "cat <filename>",
			Execute: func(args []string) *FileOperation {
				if len(args) == 0 {
					return &FileOperation{
						Success: false,
						Message: "Usage: cat <filename>",
					}
				}
				return ReadFile(args[0])
			},
		},
		"create": {
			Name:        "create",
			Description: "Create a new file",
			Usage:       "create <filename> [content]",
			Execute: func(args []string) *FileOperation {
				if len(args) == 0 {
					return &FileOperation{
						Success: false,
						Message: "Usage: create <filename> [content]",
					}
				}
				content := ""
				if len(args) > 1 {
					content = strings.Join(args[1:], " ")
				}
				return CreateFile(args[0], content)
			},
		},
		"write": {
			Name:        "write",
			Description: "Write content to a file (create or update)",
			Usage:       "write <filename> <content>",
			Execute: func(args []string) *FileOperation {
				if len(args) < 2 {
					return &FileOperation{
						Success: false,
						Message: "Usage: write <filename> <content>",
					}
				}
				filename := args[0]
				content := strings.Join(args[1:], " ")
				return CreateFile(filename, content) // CreateFile will overwrite if exists
			},
		},
		"update": {
			Name:        "update",
			Description: "Update an existing file",
			Usage:       "update <filename> <content>",
			Execute: func(args []string) *FileOperation {
				if len(args) < 2 {
					return &FileOperation{
						Success: false,
						Message: "Usage: update <filename> <content>",
					}
				}
				filename := args[0]
				content := strings.Join(args[1:], " ")
				return UpdateFile(filename, content)
			},
		},
		"rm": {
			Name:        "rm",
			Description: "Remove a file",
			Usage:       "rm <filename>",
			Execute: func(args []string) *FileOperation {
				if len(args) == 0 {
					return &FileOperation{
						Success: false,
						Message: "Usage: rm <filename>",
					}
				}
				return DeleteFile(args[0])
			},
		},
		"mkdir": {
			Name:        "mkdir",
			Description: "Create a directory",
			Usage:       "mkdir <dirname>",
			Execute: func(args []string) *FileOperation {
				if len(args) == 0 {
					return &FileOperation{
						Success: false,
						Message: "Usage: mkdir <dirname>",
					}
				}
				return CreateDirectory(args[0])
			},
		},
		"rmdir": {
			Name:        "rmdir",
			Description: "Remove a directory",
			Usage:       "rmdir <dirname>",
			Execute: func(args []string) *FileOperation {
				if len(args) == 0 {
					return &FileOperation{
						Success: false,
						Message: "Usage: rmdir <dirname>",
					}
				}
				return DeleteDirectory(args[0])
			},
		},
		"cp": {
			Name:        "cp",
			Description: "Copy a file",
			Usage:       "cp <source> <destination>",
			Execute: func(args []string) *FileOperation {
				if len(args) < 2 {
					return &FileOperation{
						Success: false,
						Message: "Usage: cp <source> <destination>",
					}
				}
				return CopyFile(args[0], args[1])
			},
		},
		"mv": {
			Name:        "mv",
			Description: "Move/rename a file",
			Usage:       "mv <source> <destination>",
			Execute: func(args []string) *FileOperation {
				if len(args) < 2 {
					return &FileOperation{
						Success: false,
						Message: "Usage: mv <source> <destination>",
					}
				}
				return MoveFile(args[0], args[1])
			},
		},
		"pwd": {
			Name:        "pwd",
			Description: "Print working directory",
			Usage:       "pwd",
			Execute: func(args []string) *FileOperation {
				return GetWorkingDirectory()
			},
		},
		"cd": {
			Name:        "cd",
			Description: "Change directory",
			Usage:       "cd <path>",
			Execute: func(args []string) *FileOperation {
				if len(args) == 0 {
					return &FileOperation{
						Success: false,
						Message: "Usage: cd <path>",
					}
				}
				return ChangeDirectory(args[0])
			},
		},
	}
}

// GetHelpText returns formatted help text for all commands
func GetHelpText() string {
	var help strings.Builder
	help.WriteString("Available file operations:\n\n")
	
	commands := GetCommands()
	for _, cmd := range commands {
		help.WriteString(fmt.Sprintf("  %s - %s\n", cmd.Usage, cmd.Description))
	}
	
	help.WriteString("\nExample usage:\n")
	help.WriteString("  /ls                    # List current directory\n")
	help.WriteString("  /cat myfile.txt        # Display file content\n")
	help.WriteString("  /create hello.txt Hello World!  # Create file with content\n")
	help.WriteString("  /mkdir newfolder       # Create directory\n")
	help.WriteString("  /cp file1.txt file2.txt  # Copy file\n")
	help.WriteString("  /help                  # Show this help\n")
	
	return help.String()
}

// ExecuteCommand parses and executes a file operation command
func ExecuteCommand(input string) *FileOperation {
	// Remove leading slash if present
	input = strings.TrimPrefix(input, "/")
	
	// Special case for help
	if input == "help" || input == "h" {
		return &FileOperation{
			Success: true,
			Message: GetHelpText(),
		}
	}
	
	// Parse command and arguments
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return &FileOperation{
			Success: false,
			Message: "No command specified. Type '/help' for available commands.",
		}
	}
	
	commandName := parts[0]
	args := parts[1:]
	
	commands := GetCommands()
	if cmd, exists := commands[commandName]; exists {
		return cmd.Execute(args)
	}
	
	return &FileOperation{
		Success: false,
		Message: fmt.Sprintf("Unknown command: %s. Type '/help' for available commands.", commandName),
	}
}